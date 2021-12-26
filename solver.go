package kiwi

import "math"

type Solver struct {
	cns                 map[*Constraint]tag
	rows                map[*symbol]*row
	vars                map[*Variable]*symbol
	edits               map[*Variable]*edit
	infeasibleRows      []*symbol
	objective           *row
	artificialObjective *row
}

func NewSolver() *Solver {
	return &Solver{
		cns:       map[*Constraint]tag{},
		rows:      map[*symbol]*row{},
		vars:      map[*Variable]*symbol{},
		edits:     map[*Variable]*edit{},
		objective: newRow(),
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
	subject := row.chooseSubject(tag)

	if subject.is(INVALID) && row.allDummies() {
		if !NearZero(row.constant) {
			return UnsatisfiableConstraintException{constraint}
		} else {
			subject = tag.marker
		}

	}

	if subject.is(INVALID) {
		if !s.addWithArtificialVariable(row) {
			return UnsatisfiableConstraintException{constraint}
		}
	} else {
		row.solveFor(subject)
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
		row.solveForPair(leaving, tag.marker)
		s.substitute(tag.marker, row)
	}

	return s.optimize(s.objective)
}

func (s *Solver) removeConstraintEffects(constraint *Constraint, tag tag) {
	if tag.marker.is(ERROR) {
		s.removeMarkerEffects(tag.marker, constraint.Strength)
	} else if tag.other.is(ERROR) {
		s.removeMarkerEffects(tag.other, constraint.Strength)
	}
}

func (s *Solver) removeMarkerEffects(marker *symbol, strength float64) {
	if row, present := s.rows[marker]; present {
		s.objective.insertRowWithCoefficient(row, float64(-strength))
	} else {
		s.objective.insertSymbolWithCoefficient(marker, float64(-strength))
	}
}

func (s *Solver) getMarkerLeavingRow(marker *symbol) (*symbol, bool) {
	dmax := math.MaxFloat64
	r1 := dmax
	r2 := dmax

	var first, second, third *symbol

	for sym, candidateRow := range s.rows {
		c := candidateRow.coefficientFor(marker)
		if c == 0.0 {
			continue
		}
		if sym.is(EXTERNAL) {
			third = sym
		} else if c < 0.0 {
			r := -candidateRow.constant / c
			if r < r1 {
				r1 = r
				first = sym
			}
		} else {
			r := candidateRow.constant / c
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
func (s *Solver) AddEditVariable(variable *Variable, options ...ConstraintOption) error {
	if _, present := s.edits[variable]; present {
		return DuplicateEditVariable{variable}
	}
	constraint := variable.EqualsConstant(0.0)
	for _, option := range options {
		option(constraint)
	}
	if constraint.Strength == REQUIRED {
		return BadRequiredStrength
	}
	s.AddConstraint(constraint)
	s.edits[variable] = &edit{
		tag:        s.cns[constraint],
		constraint: constraint,
		constant:   0.0,
	}
	return nil
}

/* Remove an edit variable from the solver.

Returns
------
UnknownEditVariable
	The given edit variable has not been added to the solver.

*/
func (s *Solver) RemoveEditVariable(variable *Variable) error {
	edit, present := s.edits[variable]
	if !present {
		return UnknownEditVariable{variable}
	}
	s.RemoveConstraint(edit.constraint)
	delete(s.edits, variable)
	return nil
}

/* Test whether an edit variable has been added to the solver.

 */
func (s *Solver) HasEditVariable(variable *Variable) bool {
	_, present := s.edits[variable]
	return present
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
	edit, present := s.edits[variable]
	if !present {
		return UnknownEditVariable{variable}
	}
	defer s.dualOptimize()
	delta := value - edit.constant
	edit.constant = value

	// Check first if the positive error variable is basic.
	row, present := s.rows[edit.tag.marker]
	if present {
		if row.add(-delta) < 0.0 {
			s.infeasibleRows = append(s.infeasibleRows, edit.tag.marker)
		}
		return nil
	}

	// Check next if the negative error variable is basic.
	row, present = s.rows[edit.tag.other]
	if present {
		if row.add(delta) < 0.0 {
			s.infeasibleRows = append(s.infeasibleRows, edit.tag.other)
		}
		return nil
	}

	// Otherwise update each row where the error variables exist.
	for sym, row := range s.rows {
		coeff := row.coefficientFor(edit.tag.marker)
		if coeff != 0.0 && row.add(delta*coeff) < 0.0 && !sym.is(EXTERNAL) {
			s.infeasibleRows = append(s.infeasibleRows, sym)
		}
		return nil
	}

	return nil
}

/* Update the values of the external solver variables.

 */
func (s *Solver) UpdateVariables() {
	for variable, symbol := range s.vars {
		if row, present := s.rows[symbol]; present {
			variable.Value = row.constant
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
 * createRow creates a new row object for the given constraint.
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
func (s *Solver) createRow(constraint *Constraint) (row *row, tag tag) {
	expression := constraint.Expression
	row = newRow(withConstant(expression.Constant))
	for _, term := range expression.Terms {
		if NearZero(term.Coefficient) {
			continue
		}

		variable := term.Variable
		sym, present := s.vars[variable]
		if !present {
			sym = newSymbol(EXTERNAL)
			s.vars[variable] = sym
		}

		if otherRow, present := s.rows[sym]; present {
			row.insertRowWithCoefficient(otherRow, term.Coefficient)
		} else {
			row.insertSymbolWithCoefficient(sym, term.Coefficient)
		}
	}

	switch constraint.Operator {
	case LE, GE:
		coeff := -1.0
		if constraint.Operator == LE {
			coeff = 1.0
		}
		slack := newSymbol(SLACK)
		tag.marker = slack
		row.insertSymbolWithCoefficient(slack, coeff)
		if constraint.Strength < REQUIRED {
			error := newSymbol(ERROR)
			tag.other = error
			row.insertSymbolWithCoefficient(error, -coeff)
			s.objective.insertSymbolWithCoefficient(error, float64(constraint.Strength))
		}
	case EQ:
		if constraint.Strength < REQUIRED {
			errplus := newSymbol(ERROR)
			errminus := newSymbol(ERROR)
			tag.marker = errplus
			tag.other = errminus
			row.insertSymbolWithCoefficient(errplus, -1.0) // v = eplus - eminus
			row.insertSymbolWithCoefficient(errminus, 1.0) // v - eplus + eminus = 0
			s.objective.insertSymbolWithCoefficient(errplus, float64(constraint.Strength))
			s.objective.insertSymbolWithCoefficient(errminus, float64(constraint.Strength))
		} else {
			dummy := newSymbol(DUMMY)
			tag.marker = dummy
			row.insertSymbol(dummy)
		}
	}

	// Ensure the row has a positive constant.
	if row.constant < 0.0 {
		row.reverseSign()
	}

	// Ensure the tag.other symbol is not nil
	if tag.other == nil {
		tag.other = newSymbol(INVALID)
	}

	return
}

/**
 * Add the row to the tableau using an artificial variable.
 *
 * This will return false if the constraint cannot be satisfied.
 */
func (s *Solver) addWithArtificialVariable(row *row) bool {
	// Create and add the artificial variable to the tableau
	art := newSymbol(SLACK)
	s.rows[art] = row.copy()

	// Optimize the artificial objective. This is successful only
	// if the artificial objective could be optimized to zero.
	s.artificialObjective = row.copy()
	s.optimize(s.artificialObjective)
	success := NearZero(s.artificialObjective.constant)
	s.artificialObjective = nil

	// If the artificial variable is basic, pivot the row so that
	// it becomes basic. If the row is constant, exit early.
	if rowptr, present := s.rows[art]; present {
		delete(s.rows, art)
		if len(rowptr.cells) == 0 {
			return success
		}
		entering := rowptr.anyPivotableSymbol()
		if entering.is(INVALID) {
			return false // unsatisfiable (will this ever happen?)
		}
		rowptr.solveForPair(art, entering)
		s.substitute(entering, rowptr)
		s.rows[entering] = rowptr
	}

	// Remove the artificial variable from the tableau.
	for _, row := range s.rows {
		row.removeSymbol(art)
	}
	s.objective.removeSymbol(art)
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
func (s *Solver) optimize(objective *row) error {
	for {
		enterSym := objective.getEnteringSymbol()
		if enterSym.is(INVALID) {
			return nil
		}

		// Compute the row which holds the exit symbol for a pivot.
		ratio := math.MaxFloat64
		var exitSym *symbol
		var exitRow *row
		for sym, row := range s.rows {
			if !sym.is(EXTERNAL) {
				temp := row.coefficientFor(enterSym)
				if temp < 0.0 {
					tempRatio := -row.constant / temp
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
		exitRow.solveForPair(exitSym, enterSym)
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
		r := s.rows[leaving]

		if r != nil && r.constant < 0.0 {
			entering := s.objective.getDualEnteringSymbol(r)
			if entering.is(INVALID) {
				return InternalSolverError
			}
			delete(s.rows, leaving)

			r.solveForPair(leaving, entering)
			s.substitute(entering, r)
			s.rows[entering] = r
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
func (s *Solver) substitute(sym *symbol, other *row) {
	for isym, irow := range s.rows {
		irow.substitute(sym, other)
		if !isym.is(EXTERNAL) && irow.constant < 0.0 {
			s.infeasibleRows = append(s.infeasibleRows, isym)
		}
	}
	s.objective.substitute(sym, other)
	if s.artificialObjective != nil {
		s.artificialObjective.substitute(sym, other)
	}
}
