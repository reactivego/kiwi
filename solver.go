package kiwi

import (
	"math"

	"github.com/reactivego/kiwi/op"
	"github.com/reactivego/kiwi/strength"
	"github.com/reactivego/kiwi/symbol"
)

type Solver interface {

	/* Add a constraint to the solver.

	Returns
	------
	DuplicateConstraint
		The given constraint has already been added to the solver.

	UnsatisfiableConstraint
		The given constraint is required and cannot be satisfied.

	*/
	AddConstraint(constraint Constraint) error

	/* Remove a constraint from the solver.

	Returns
	------
	UnknownConstraint
		The given constraint has not been added to the solver.

	*/
	RemoveConstraint(constraint Constraint) error

	/* Test whether a constraint has been added to the solver.

	 */
	HasConstraint(constraint Constraint) bool

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
	AddEditVariable(variable Variable, strength float64) error

	/* Remove an edit variable from the solver.

	Returns
	------
	UnknownEditVariable
		The given edit variable has not been added to the solver.

	*/
	RemoveEditVariable(variable Variable) error

	/* Test whether an edit variable has been added to the solver.

	 */
	HasEditVariable(variable Variable) bool

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
	SuggestValue(variable Variable, value float64) error

	/* Update the values of the external solver variables.

	 */
	UpdateVariables()

	/* Reset the solver to the empty starting condition.

	This method resets the internal solver state to the empty starting
	condition, as if no constraints or edit variables have been added.
	This can be faster than deleting the solver and creating a new one
	when the entire system must change, since it can avoid unecessary
	heap (de)allocations.

	*/
	Reset()

	/* Dump a representation of the solver internals to stdout.

	 */
	Dump()
}

func NewSolver() Solver {
	return &solver{
		cns:       map[Constraint]tag{},
		rows:      map[Symbol]Row{},
		vars:      map[Variable]Symbol{},
		objective: NewRow(),
	}
}

type solver struct {
	cns                 map[Constraint]tag
	rows                map[Symbol]Row
	vars                map[Variable]Symbol
	infeasibleRows      []Symbol
	objective           Row
	artificialObjective Row
}

type tag struct {
	marker Symbol
	other  Symbol
}

func (s *solver) AddConstraint(constraint Constraint) error {

	if _, present := s.cns[constraint]; present {
		return DuplicateConstraintException{}
	}

	row, tag := s.createRow(constraint)
	subject := row.ChooseSubject(tag)

	if subject.IsInvalid() && row.AllDummies() {
		if !NearZero(row.GetConstant()) {
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

func (s *solver) RemoveConstraint(constraint Constraint) error {
	tag, present := s.cns[constraint]
	if !present {
		return UnknownConstraintException{constraint}
	}

	delete(s.cns, constraint)
	s.RemoveConstraintEffects(constraint, tag)

	if _, present := s.rows[tag.marker]; present {
		delete(s.rows, tag.marker)
	} else {
		row, present := s.GetMarkerLeavingRow(tag.marker)
		if !present {
			return InternalSolverError{"internal solver error"}
		}

		//This looks wrong! changes made below
		//Symbol leaving = tag.marker;
		//rows.remove(tag.marker);

		var leaving Symbol
		for s, v := range s.rows {
			if v == row {
				leaving = s
			}
		}
		if leaving == nil {
			return InternalSolverError{"internal solver error"}
		}

		delete(s.rows, leaving)
		row.SolveForPair(leaving, tag.marker)
		s.substitute(tag.marker, row)
	}

	return s.optimize(s.objective)
}

func (s *solver) RemoveConstraintEffects(constraint Constraint, tag tag) {
	if constraint == nil {
		panic("constraint is nil")
	}

	if tag.marker.IsError() {
		s.RemoveMarkerEffects(tag.marker, constraint.GetStrength())
	} else if tag.other.IsError() {
		s.RemoveMarkerEffects(tag.other, constraint.GetStrength())
	}
}

func (s *solver) RemoveMarkerEffects(marker Symbol, strength strength.Value) {
	if row, present := s.rows[marker]; present {
		s.objective.InsertRowWithCoefficient(row, -strength.Float64())
	} else {
		s.objective.InsertSymbolWithCoefficient(marker, -strength.Float64())
	}
}

func (s *solver) GetMarkerLeavingRow(marker Symbol) (Row, bool) {

	dmax := math.MaxFloat64
	r1 := dmax
	r2 := dmax

	var first, second, third Row

	for s, candidateRow := range s.rows {
		c := candidateRow.CoefficientFor(marker)
		if c == 0.0 {
			continue
		}
		if s.IsExternal() {
			third = candidateRow
		} else if c < 0.0 {
			r := -candidateRow.GetConstant() / c
			if r < r1 {
				r1 = r
				first = candidateRow
			}
		} else {
			r := candidateRow.GetConstant() / c
			if r < r2 {
				r2 = r
				second = candidateRow
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

func (s *solver) HasConstraint(constraint Constraint) bool {
	return false
}

func (s *solver) AddEditVariable(variable Variable, strength float64) error {
	return nil
}

func (s *solver) RemoveEditVariable(variable Variable) error {
	return nil
}

func (s *solver) HasEditVariable(variable Variable) bool {
	return false
}

func (s *solver) SuggestValue(variable Variable, value float64) error {
	return nil
}

func (s *solver) UpdateVariables() {
	for variable, symbol := range s.vars {
		if row, present := s.rows[symbol]; present {
			variable.SetValue(row.GetConstant())
		} else {
			variable.SetValue(0.0)
		}
	}
}

func (s *solver) Reset() {

}

func (s *solver) Dump() {

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
func (s *solver) createRow(constraint Constraint) (row Row, tag tag) {
	expression := constraint.GetExpression()
	row = NewRowWithConstant(expression.GetConstant())
	for _, term := range expression.GetTerms() {
		if NearZero(term.GetCoefficient()) {
			continue
		}

		variable := term.GetVariable()
		sym, present := s.vars[variable]
		if !present {
			sym = NewSymbol(symbol.EXTERNAL)
			s.vars[variable] = sym
		}

		if otherRow, present := s.rows[sym]; present {
			row.InsertRowWithCoefficient(otherRow, term.GetCoefficient())
		} else {
			row.InsertSymbolWithCoefficient(sym, term.GetCoefficient())
		}
	}

	switch constraint.GetOp() {
	case op.LE, op.GE:
		coeff := -1.0
		if constraint.GetOp() == op.LE {
			coeff = 1.0
		}
		slack := NewSymbol(symbol.SLACK)
		tag.marker = slack
		row.InsertSymbolWithCoefficient(slack, coeff)
		if constraint.GetStrength() < strength.REQUIRED {
			error := NewSymbol(symbol.ERROR)
			tag.other = error
			row.InsertSymbolWithCoefficient(error, -coeff)
			s.objective.InsertSymbolWithCoefficient(error, constraint.GetStrength().Float64())
		}
	case op.EQ:
		if constraint.GetStrength() < strength.REQUIRED {
			errplus := NewSymbol(symbol.ERROR)
			errminus := NewSymbol(symbol.ERROR)
			tag.marker = errplus
			tag.other = errminus
			row.InsertSymbolWithCoefficient(errplus, -1.0) // v = eplus - eminus
			row.InsertSymbolWithCoefficient(errminus, 1.0) // v - eplus + eminus = 0
			s.objective.InsertSymbolWithCoefficient(errplus, constraint.GetStrength().Float64())
			s.objective.InsertSymbolWithCoefficient(errminus, constraint.GetStrength().Float64())
		} else {
			dummy := NewSymbol(symbol.DUMMY)
			tag.marker = dummy
			row.InsertSymbol(dummy)
		}
	}

	// Ensure the row has a positive constant.
	if row.GetConstant() < 0.0 {
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
func (s *solver) addWithArtificialVariable(row Row) bool {
	// Create and add the artificial variable to the tableau
	art := NewSymbol(symbol.SLACK)
	s.rows[art] = CopyRow(row)

	// Optimize the artificial objective. This is successful only
	// if the artificial objective could be optimized to zero.
	s.artificialObjective = CopyRow(row)
	s.optimize(s.artificialObjective)
	success := NearZero(s.artificialObjective.GetConstant())
	s.artificialObjective = nil

	// If the artificial variable is basic, pivot the row so that
	// it becomes basic. If the row is constant, exit early.
	if rowptr, present := s.rows[art]; present {
		delete(s.rows, art)
		if len(rowptr.GetCells()) == 0 {
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
func (s *solver) optimize(objective Row) error {
	for {
		enterSym := objective.GetEnteringSymbol()
		if enterSym.IsInvalid() {
			return nil
		}

		// Compute the row which holds the exit symbol for a pivot.
		ratio := math.MaxFloat64
		var exitSym Symbol
		var exitRow Row
		for sym, row := range s.rows {
			if !sym.IsExternal() {
				temp := row.CoefficientFor(enterSym)
				if temp < 0.0 {
					tempRatio := -row.GetConstant() / temp
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
			return InternalSolverError{"The objective is unbounded."}
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
func (s *solver) dualOptimize() error {
	for len(s.infeasibleRows) > 0 {

		last := len(s.infeasibleRows) - 1
		leaving := s.infeasibleRows[last]
		s.infeasibleRows[last] = nil
		s.infeasibleRows = s.infeasibleRows[:last]
		row := s.rows[leaving]

		if row != nil && row.GetConstant() < 0.0 {
			entering := s.objective.GetDualEnteringSymbol(row)
			if entering.IsInvalid() {
				return InternalSolverError{"internal solver error"}
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
func (s *solver) substitute(sym Symbol, row Row) {
	for isym, irow := range s.rows {
		irow.Substitute(sym, row)
		if !isym.IsExternal() && irow.GetConstant() < 0.0 {
			s.infeasibleRows = append(s.infeasibleRows, isym)
		}
	}
	s.objective.Substitute(sym, row)
	if s.artificialObjective != nil {
		s.artificialObjective.Substitute(sym, row)
	}
}
