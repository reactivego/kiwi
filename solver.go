package kiwi

import "math"

type Solver struct {
	cns                 map[*Constraint]tag
	rows                map[*Symbol]*Row
	vars                map[*Variable]*Symbol
	infeasibleRows      []*Symbol
	objective           *Row
	artificialObjective *Row
}

func NewSolver() *Solver {
	return &Solver{
		cns:       map[*Constraint]tag{},
		rows:      map[*Symbol]*Row{},
		vars:      map[*Variable]*Symbol{},
		objective: NewRow(),
	}
}

/* Add a constraint to the solver.

Returns
------
DuplicateConstraint
	The given constraint has already been added to the solver.

UnsatisfiableConstraint
	The given constraint is required and cannot be satisfied.

*/
func (s *Solver) AddConstraint(constraint *Constraint, options ...ConstraintOption) error {

	if s.HasConstraint(constraint) {
		return DuplicateConstraintException
	}

	for _, option := range options {
		option(constraint)
	}

	row, tag := s.createRow(constraint)
	subject := row.ChooseSubject(tag)

	if subject.IsInvalid() && row.AllDummies() {
		if !NearZero(row.Constant) {
			return UnsatisfiableConstraintException{constraint}
		} else {
			subject = tag.marker
		}

	}

	if subject.IsInvalid() {
		if !s.addWithArtificialVariable(row) {
			return UnsatisfiableConstraintException{constraint}
		}
	} else {
		row.SolveFor(subject)
		s.substitute(subject, row)
		s.rows[subject] = row
	}

	s.cns[constraint] = tag

	return s.optimize(s.objective)
}

/* Remove a constraint from the solver.

Returns
------
UnknownConstraint
	The given constraint has not been added to the solver.

*/
func (s *Solver) RemoveConstraint(constraint *Constraint) error {
	tag, present := s.cns[constraint]
	if !present {
		return UnknownConstraintException{constraint}
	}

	delete(s.cns, constraint)

	// Remove the error effects from the objective function
	// *before* pivoting, or substitutions into the objective
	// will lead to incorrect solver results.
	s.removeConstraintEffects(constraint, tag)

	// If the marker is basic, simply drop the row. Otherwise,
	// pivot the marker into the basis and then drop the row.
	if _, present := s.rows[tag.marker]; present {
		delete(s.rows, tag.marker)
	} else {
		leaving, present := s.getMarkerLeavingRow(tag.marker)
		if !present {
			return FailedToFindLeavingRow
		}
		row := s.rows[leaving]
		delete(s.rows, leaving)
		row.SolveForPair(leaving, tag.marker)
		s.substitute(tag.marker, row)
	}

	return s.optimize(s.objective)
}

func (s *Solver) removeConstraintEffects(constraint *Constraint, tag tag) {
	if constraint == nil {
		panic("constraint is nil")
	}

	if tag.marker.IsError() {
		s.removeMarkerEffects(tag.marker, constraint.Strength)
	} else if tag.other.IsError() {
		s.removeMarkerEffects(tag.other, constraint.Strength)
	}
}

func (s *Solver) removeMarkerEffects(marker *Symbol, strength Strength) {
	if row, present := s.rows[marker]; present {
		s.objective.InsertRowWithCoefficient(row, float64(-strength))
	} else {
		s.objective.InsertSymbolWithCoefficient(marker, float64(-strength))
	}
}

func (s *Solver) getMarkerLeavingRow(marker *Symbol) (*Symbol, bool) {
	dmax := math.MaxFloat64
	r1 := dmax
	r2 := dmax

	var first, second, third *Symbol

	for sym, candidateRow := range s.rows {
		c := candidateRow.CoefficientFor(marker)
		if c == 0.0 {
			continue
		}
		if sym.IsExternal() {
			third = sym
		} else if c < 0.0 {
			r := -candidateRow.Constant / c
			if r < r1 {
				r1 = r
				first = sym
			}
		} else {
			r := candidateRow.Constant / c
			if r < r2 {
				r2 = r
				second = sym
			}
		}
	}

	if first != nil {
		return first, true
	}
	if second != nil {
		return second, true
	}
	if third != nil {
		return third, true
	}
	return nil, false
}

/* Test whether a constraint has been added to the solver.

 */
func (s *Solver) HasConstraint(constraint *Constraint) bool {
	_, present := s.cns[constraint]
	return present
}

/* Add an edit variable to the solver.

This method should be called before the `suggestValue` method is
used to supply a suggested value for the given edit variable.

Returns
------
DuplicateEditVariable
	The given edit variable has already been added to the solver.

BadRequiredStrength
	The given strength is >= required.

*/
func (s *Solver) AddEditVariable(variable *Variable, strength float64) error {
	return nil
}

/* Remove an edit variable from the solver.

Returns
------
UnknownEditVariable
	The given edit variable has not been added to the solver.

*/
func (s *Solver) RemoveEditVariable(variable *Variable) error {
	return nil
}

/* Test whether an edit variable has been added to the solver.

 */
func (s *Solver) HasEditVariable(variable *Variable) bool {
	return false
}

/* Suggest a value for the given edit variable.

This method should be used after an edit variable as been added to
the solver in order to suggest the value for that variable. After
all suggestions have been made, the `solve` method can be used to
update the values of all variables.

Returns
------
UnknownEditVariable
	The given edit variable has not been added to the solver.

*/
func (s *Solver) SuggestValue(variable *Variable, value float64) error {
	return nil
}

/* Update the values of the external solver variables.

 */
func (s *Solver) UpdateVariables() {
	for variable, symbol := range s.vars {
		if row, present := s.rows[symbol]; present {
			variable.Value = row.Constant
		} else {
			variable.Value = 0.0
		}
	}
}

/* Reset the solver to the empty starting condition.

This method resets the internal solver state to the empty starting
condition, as if no constraints or edit variables have been added.
This can be faster than deleting the solver and creating a new one
when the entire system must change, since it can avoid unecessary
heap (de)allocations.

*/
func (s *Solver) Reset() {

}

/* Dump a representation of the solver internals to stdout.

 */
func (s *Solver) Dump() {

}

/**
 * Create a new Row object for the given constraint.
 *
 * The terms in the constraint will be converted to cells in the row.
 * Any term in the constraint with a coefficient of zero is ignored.
 * This method uses the `GetVarSymbol` method to get the symbol for
 * the variables added to the row. If the symbol for a given cell
 * variable is basic, the cell variable will be substituted with the
 * basic row.
 *
 * The necessary slack and error variables will be added to the row.
 * If the constant for the row is negative, the sign for the row
 * will be inverted so the constant becomes positive.
 *
 * The tag will be updated with the marker and error symbols to use
 * for tracking the movement of the constraint in the tableau.
 */
func (s *Solver) createRow(constraint *Constraint) (row *Row, tag tag) {
	expression := constraint.Expression
	row = NewRow(WithConstant(expression.Constant))
	for _, term := range expression.Terms {
		if NearZero(term.Coefficient) {
			continue
		}

		variable := term.Variable
		sym, present := s.vars[variable]
		if !present {
			sym = NewSymbol(EXTERNAL)
			s.vars[variable] = sym
		}

		if otherRow, present := s.rows[sym]; present {
			row.InsertRowWithCoefficient(otherRow, term.Coefficient)
		} else {
			row.InsertSymbolWithCoefficient(sym, term.Coefficient)
		}
	}

	switch constraint.Operator {
	case LE, GE:
		coeff := -1.0
		if constraint.Operator == LE {
			coeff = 1.0
		}
		slack := NewSymbol(SLACK)
		tag.marker = slack
		row.InsertSymbolWithCoefficient(slack, coeff)
		if constraint.Strength < REQUIRED {
			error := NewSymbol(ERROR)
			tag.other = error
			row.InsertSymbolWithCoefficient(error, -coeff)
			s.objective.InsertSymbolWithCoefficient(error, float64(constraint.Strength))
		}
	case EQ:
		if constraint.Strength < REQUIRED {
			errplus := NewSymbol(ERROR)
			errminus := NewSymbol(ERROR)
			tag.marker = errplus
			tag.other = errminus
			row.InsertSymbolWithCoefficient(errplus, -1.0) // v = eplus - eminus
			row.InsertSymbolWithCoefficient(errminus, 1.0) // v - eplus + eminus = 0
			s.objective.InsertSymbolWithCoefficient(errplus, float64(constraint.Strength))
			s.objective.InsertSymbolWithCoefficient(errminus, float64(constraint.Strength))
		} else {
			dummy := NewSymbol(DUMMY)
			tag.marker = dummy
			row.InsertSymbol(dummy)
		}
	}

	// Ensure the row has a positive constant.
	if row.Constant < 0.0 {
		row.ReverseSign()
	}

	// Ensure the tag.other symbol is not nil
	if tag.other == nil {
		tag.other = NewInvalidSymbol()
	}

	return
}

/**
 * Add the row to the tableau using an artificial variable.
 *
 * This will return false if the constraint cannot be satisfied.
 */
func (s *Solver) addWithArtificialVariable(row *Row) bool {
	// Create and add the artificial variable to the tableau
	art := NewSymbol(SLACK)
	s.rows[art] = row.Copy()

	// Optimize the artificial objective. This is successful only
	// if the artificial objective could be optimized to zero.
	s.artificialObjective = row.Copy()
	s.optimize(s.artificialObjective)
	success := NearZero(s.artificialObjective.Constant)
	s.artificialObjective = nil

	// If the artificial variable is basic, pivot the row so that
	// it becomes basic. If the row is constant, exit early.
	if rowptr, present := s.rows[art]; present {
		delete(s.rows, art)
		if len(rowptr.Cells) == 0 {
			return success
		}
		entering := rowptr.AnyPivotableSymbol()
		if entering.IsInvalid() {
			return false // unsatisfiable (will this ever happen?)
		}
		rowptr.SolveForPair(art, entering)
		s.substitute(entering, rowptr)
		s.rows[entering] = rowptr
	}

	// Remove the artificial variable from the tableau.
	for _, row := range s.rows {
		row.RemoveSymbol(art)
	}
	s.objective.RemoveSymbol(art)
	return success
}

/**
 * Optimize the system for the given objective function.
 *
 * This method performs iterations of Phase 2 of the simplex method
 * until the objective function reaches a minimum.
 *
 * @returns InternalSolverError The value of the objective function is unbounded.
 */
func (s *Solver) optimize(objective *Row) error {
	for {
		enterSym := objective.GetEnteringSymbol()
		if enterSym.IsInvalid() {
			return nil
		}

		// Compute the row which holds the exit symbol for a pivot.
		ratio := math.MaxFloat64
		var exitSym *Symbol
		var exitRow *Row
		for sym, row := range s.rows {
			if !sym.IsExternal() {
				temp := row.CoefficientFor(enterSym)
				if temp < 0.0 {
					tempRatio := -row.Constant / temp
					if tempRatio < ratio {
						ratio = tempRatio
						exitSym = sym
						exitRow = row
					}
				}
			}
		}

		// If no appropriate exit symbol was found, this indicates that
		// the objective function is unbounded.
		if exitSym == nil || exitRow == nil {
			return UnboundedObjectiveError
		}
		// pivot the entering symbol into the basis
		delete(s.rows, exitSym)
		exitRow.SolveForPair(exitSym, enterSym)
		s.substitute(enterSym, exitRow)
		s.rows[enterSym] = exitRow
	}
}

/**
 * Optimize the system using the dual of the simplex method.
 *
 * The current state of the system should be such that the objective
 * function is optimal, but not feasible. This method will perform
 * an iteration of the dual simplex method to make the solution both
 * optimal and feasible.
 *
 * @returns InternalSolverError The system cannot be dual optimized.
 */
func (s *Solver) dualOptimize() error {
	for len(s.infeasibleRows) > 0 {

		last := len(s.infeasibleRows) - 1
		leaving := s.infeasibleRows[last]
		s.infeasibleRows[last] = nil
		s.infeasibleRows = s.infeasibleRows[:last]
		row := s.rows[leaving]

		if row != nil && row.Constant < 0.0 {
			entering := s.objective.GetDualEnteringSymbol(row)
			if entering.IsInvalid() {
				return InternalSolverError
			}
			delete(s.rows, leaving)

			row.SolveForPair(leaving, entering)
			s.substitute(entering, row)
			s.rows[entering] = row
		}

	}
	return nil
}

/**
 * Substitute the parametric symbol with the given row.
 *
 * This method will substitute all instances of the parametric symbol
 * in the tableau and the objective function with the given row.
 */
func (s *Solver) substitute(sym *Symbol, row *Row) {
	for isym, irow := range s.rows {
		irow.Substitute(sym, row)
		if !isym.IsExternal() && irow.Constant < 0.0 {
			s.infeasibleRows = append(s.infeasibleRows, isym)
		}
	}
	s.objective.Substitute(sym, row)
	if s.artificialObjective != nil {
		s.artificialObjective.Substitute(sym, row)
	}
}
