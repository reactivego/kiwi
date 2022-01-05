// SPDX-License-Identifier: BSD-3-Clause

package kiwi

import (
	"fmt"
	"math"
	"strings"
)

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

/*
AddConstraint adds a constraint to the solver.

Returns

	DuplicateConstraint
The given constraint has already been added to the solver.
	UnsatisfiableConstraint
The given constraint is required and cannot be satisfied.
*/
func (s *Solver) AddConstraint(constraint *Constraint, options ...ConstraintOption) error {
	_, present := s.cns[constraint]
	if present {
		return DuplicateConstraint{constraint}
	}

	for _, option := range options {
		option(constraint)
	}

	// Creating a row causes symbols to be reserved for the variables
	// in the constraint. If this method exits with an exception,
	// then its possible those variables will linger in the var map.
	// Since its likely that those variables will be used in other
	// constraints and since exceptional conditions are uncommon,
	// i'm not too worried about aggressive cleanup of the var map.
	row, tag := s.createRow(constraint)
	subject := row.chooseSubject(tag)

	// If chooseSubject could not find a valid entering symbol, one
	// last option is available if the entire row is composed of
	// dummy variables. If the constant of the row is zero, then
	// this represents redundant constraints and the new dummy
	// marker can enter the basis. If the constant is non-zero,
	// then it represents an unsatisfiable constraint.
	if subject.is(INVALID) && row.allDummies() {
		if !NearZero(row.constant) {
			return UnsatisfiableConstraint{constraint}
		} else {
			subject = tag.marker
		}

	}

	// If an entering symbol still isn't found, then the row must
	// be added using an artificial variable. If that fails, then
	// the row represents an unsatisfiable constraint.
	if subject.is(INVALID) {
		if !s.addWithArtificialVariable(row) {
			return UnsatisfiableConstraint{constraint}
		}
	} else {
		row.solveFor(subject)
		s.substitute(subject, row)
		s.rows[subject] = row
	}

	s.cns[constraint] = tag

	// Optimizing after each constraint is added performs less
	// aggregate work due to a smaller average system size. It
	// also ensures the solver remains in a consistent state.
	return s.optimize(s.objective)
}

/*
RemoveConstraint removes a constraint from the solver.

Returns

	UnknownConstraint
The given constraint has not been added to the solver.
*/
func (s *Solver) RemoveConstraint(constraint *Constraint) error {
	tag, present := s.cns[constraint]
	if !present {
		return UnknownConstraint{constraint}
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

func (s *Solver) removeMarkerEffects(marker *symbol, strength Strength) {
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

/*
HasConstraint tests whether a constraint has been added to the solver.
*/
func (s *Solver) HasConstraint(constraint *Constraint) bool {
	_, present := s.cns[constraint]
	return present
}

/*
AddEditVariable adds an edit variable to the solver.

This method should be called before the `suggestValue` method is
used to supply a suggested value for the given edit variable.

Returns

	DuplicateEditVariable
The given edit variable has already been added to the solver.
	BadRequiredStrength
The given strength is >= required.
*/
func (s *Solver) AddEditVariable(variable *Variable, options ...ConstraintOption) error {
	if _, present := s.edits[variable]; present {
		return DuplicateEditVariable{variable}
	}
	e := Expression{[]Term{{variable, 1.0}}, 0.0}
	constraint := NewConstraint(e, EQ, options...)
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

/*
RemoveEditVariable removes an edit variable from the solver.

Returns

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

/*
HasEditVariable tests whether an edit variable has been added to the solver.
*/
func (s *Solver) HasEditVariable(variable *Variable) bool {
	_, present := s.edits[variable]
	return present
}

/*
SuggestValue suggests a value for the given edit variable.

This method should be used after an edit variable as been added to
the solver in order to suggest the value for that variable. After
all suggestions have been made, the `solve` method can be used to
update the values of all variables.

Returns

	UnknownEditVariable
The given edit variable has not been added to the solver.
*/
func (s *Solver) SuggestValue(variable *Variable, value float64) error {
	info, present := s.edits[variable]
	if !present {
		return UnknownEditVariable{variable}
	}
	delta := value - info.constant
	info.constant = value

	if row, present := s.rows[info.tag.marker]; present {
		// Check first if the positive error variable is basic.
		if row.add(-delta) < 0.0 {
			s.infeasibleRows = append(s.infeasibleRows, info.tag.marker)
		}
	} else if row, present = s.rows[info.tag.other]; present {
		// Check next if the negative error variable is basic.
		if row.add(delta) < 0.0 {
			s.infeasibleRows = append(s.infeasibleRows, info.tag.other)
		}
	} else {
		// Otherwise update each row where the error variables exist.
		for sym, row := range s.rows {
			coeff := row.coefficientFor(info.tag.marker)
			if coeff != 0.0 && row.add(delta*coeff) < 0.0 && !sym.is(EXTERNAL) {
				s.infeasibleRows = append(s.infeasibleRows, sym)
			}
		}
	}

	return s.dualOptimize()
}

/*
UpdateVariables updates the values of the external solver variables.
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

/*
Reset resets the solver to the empty starting condition.

This method resets the internal solver state to the empty starting
condition, as if no constraints or edit variables have been added.
This can be faster than deleting the solver and creating a new one
when the entire system must change, since it can avoid unecessary
heap (de)allocations.
*/
func (s *Solver) Reset() {
	for k := range s.cns {
		delete(s.cns, k)
	}
	for k := range s.rows {
		delete(s.rows, k)
	}
	for k := range s.vars {
		delete(s.vars, k)
	}
	for k := range s.edits {
		delete(s.edits, k)
	}
	s.infeasibleRows = nil
	s.objective = newRow()
	s.artificialObjective = nil
}

/*
createRow createRow creates a new row object for the given constraint.

The terms in the constraint will be converted to cells in the row.
Any term in the constraint with a coefficient of zero is ignored.
This method uses the `GetVarSymbol` method to get the symbol for
the variables added to the row. If the symbol for a given cell
variable is basic, the cell variable will be substituted with the
basic row.

The necessary slack and error variables will be added to the row.
If the constant for the row is negative, the sign for the row
will be inverted so the constant becomes positive.

The tag will be updated with the marker and error symbols to use
for tracking the movement of the constraint in the tableau.
*/
func (s *Solver) createRow(constraint *Constraint) (row *row, tag tag) {
	expression := constraint.Expression
	row = newRow(withConstant(expression.Constant))
	for _, term := range expression.Terms {
		if NearZero(term.Coefficient) {
			continue
		}

		sym, present := s.vars[term.Variable]
		if !present {
			sym = newSymbol(EXTERNAL)
			s.vars[term.Variable] = sym
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
		tag.marker = newSymbol(SLACK)
		row.insertSymbolWithCoefficient(tag.marker, coeff)
		if constraint.Strength < REQUIRED {
			tag.other = newSymbol(ERROR)
			row.insertSymbolWithCoefficient(tag.other, -coeff)
			s.objective.insertSymbolWithCoefficient(tag.other, float64(constraint.Strength))
		}
	case EQ:
		if constraint.Strength < REQUIRED {
			tag.marker = newSymbol(ERROR)                     // errplus
			tag.other = newSymbol(ERROR)                      // errminus
			row.insertSymbolWithCoefficient(tag.marker, -1.0) // v = eplus - eminus
			row.insertSymbolWithCoefficient(tag.other, 1.0)   // v - eplus + eminus = 0
			s.objective.insertSymbolWithCoefficient(tag.marker, float64(constraint.Strength))
			s.objective.insertSymbolWithCoefficient(tag.other, float64(constraint.Strength))
		} else {
			tag.marker = newSymbol(DUMMY)
			row.insertSymbol(tag.marker)
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

	return row, tag
}

/*
addWithArtificialVariable adds the row to the tableau using an artificial variable.

This will return false if the constraint cannot be satisfied.
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

/*
optimize optimizes the system for the given objective function.

This method performs iterations of Phase 2 of the simplex method
until the objective function reaches a minimum.

Returns

	InternalSolverError
The value of the objective function is unbounded.
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
			return UnboundedObjective
		}
		// pivot the entering symbol into the basis
		delete(s.rows, exitSym)
		exitRow.solveForPair(exitSym, enterSym)
		s.substitute(enterSym, exitRow)
		s.rows[enterSym] = exitRow
	}
}

/*
dualOptimize optimizes the system using the dual of the simplex method.

The current state of the system should be such that the objective
function is optimal, but not feasible. This method will perform
an iteration of the dual simplex method to make the solution both
optimal and feasible.

Returns

	InternalSolverError
The system cannot be dual optimized.
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

/*
substitute substitutes the parametric symbol with the given row.

This method will substitute all instances of the parametric symbol
in the tableau and the objective function with the given row.
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

func (s Solver) String() string {
	sb := strings.Builder{}
	fmt.Fprintln(&sb, "Objective")
	fmt.Fprintln(&sb, "---------")
	fmt.Fprintln(&sb, s.objective)
	fmt.Fprintln(&sb)
	fmt.Fprintln(&sb, "Tableau")
	fmt.Fprintln(&sb, "-------")
	for s, r := range s.rows {
		fmt.Fprintln(&sb, s, "|", r)
	}
	fmt.Fprintln(&sb)
	fmt.Fprintln(&sb, "Infeasible")
	fmt.Fprintln(&sb, "----------")
	for _, s := range s.infeasibleRows {
		fmt.Fprintln(&sb, s)
	}
	fmt.Fprintln(&sb)
	fmt.Fprintln(&sb, "Variables")
	fmt.Fprintln(&sb, "---------")
	for v, s := range s.vars {
		fmt.Fprintln(&sb, v, " = ", s)
	}
	fmt.Fprintln(&sb)
	fmt.Fprintln(&sb, "Edit Variables")
	fmt.Fprintln(&sb, "--------------")
	for e := range s.edits {
		fmt.Fprintln(&sb, e)
	}
	fmt.Fprintln(&sb)
	fmt.Fprintln(&sb, "Constraints")
	fmt.Fprintln(&sb, "-----------")
	for c := range s.cns {
		fmt.Fprintln(&sb, c)
	}
	return sb.String()
}
