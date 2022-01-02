package kiwi

import "math"

type row struct {
	constant float64
	cells    map[*symbol]float64
}

type rowOption func(*row)

func withConstant(constant float64) rowOption {
	return func(r *row) {
		r.constant = constant
	}
}

func newRow(options ...rowOption) *row {
	r := &row{cells: map[*symbol]float64{}}
	for _, option := range options {
		option(r)
	}
	return r
}

func (r *row) copy() *row {
	cells := make(map[*symbol]float64)
	for s, c := range r.cells {
		cells[s] = c
	}
	return &row{r.constant, cells}
}

/*
add adds a constant value to the row constant.

Returns

The new value of the constant
 */
func (r *row) add(value float64) float64 {
	r.constant += value
	return r.constant
}

/*
insertSymbolWithCoefficient inserts a symbol into the row with a given coefficient.

If the symbol already exists in the row, the coefficient will be
added to the existing coefficient. If the resulting coefficient
is zero, the symbol will be removed from the row
 */
func (r *row) insertSymbolWithCoefficient(sym *symbol, coeff float64) {
	coeff += r.cells[sym]
	if NearZero(coeff) {
		delete(r.cells, sym)
	} else {
		r.cells[sym] = coeff
	}
}

/*
insertSymbol inserts a symbol into the row with coefficient 1.0.

If the symbol already exists in the row, the coefficient will be
added to the existing coefficient. If the resulting coefficient
is zero, the symbol will be removed from the row
 */
func (r *row) insertSymbol(sym *symbol) {
	r.insertSymbolWithCoefficient(sym, 1.0)
}

/*
insertRowWithCoefficient inserts a row into this row with a given coefficient.
The constant and the cells of the other row will be multiplied by
the coefficient and added to this row. Any cell with a resulting
coefficient of zero will be removed from the row.
 */
func (r *row) insertRowWithCoefficient(other *row, coeff float64) {
	r.constant += other.constant * coeff
	for otherSym, otherCoeff := range other.cells {
		combinedCoeff := r.cells[otherSym] + otherCoeff*coeff
		if NearZero(combinedCoeff) {
			delete(r.cells, otherSym)
		} else {
			r.cells[otherSym] = combinedCoeff
		}
	}
}

/*
removeSymbol removes the given symbol from the row.
 */
func (r *row) removeSymbol(sym *symbol) {
	delete(r.cells, sym)
}

/*
reverseSign reverse the sign of the constant and all cells in the row.
 */
func (r *row) reverseSign() {
	r.constant = -r.constant
	cells := make(map[*symbol]float64)
	for s, c := range r.cells {
		cells[s] = -c
	}
	r.cells = cells
}

/*
chooseSubject chooses the subject for solving for the row

This method will choose the best subject for using as the solve
target for the row. An invalid symbol will be returned if there
is no valid target.
The symbols are chosen according to the following precedence:
	1) The first symbol representing an external variable.
	2) A negative slack or error tag variable.
If a subject cannot be found, an invalid symbol will be returned.
 */
func (r *row) chooseSubject(tag tag) *symbol {
	for sym := range r.cells {
		if sym.is(EXTERNAL) {
			return sym
		}
	}

	if tag.marker.is(SLACK) || tag.marker.is(ERROR) {
		if r.coefficientFor(tag.marker) < 0.0 {
			return tag.marker
		}
	}

	if tag.other != nil && (tag.other.is(SLACK) || tag.other.is(ERROR)) {
		if r.coefficientFor(tag.other) < 0.0 {
			return tag.other
		}
	}

	return newSymbol(INVALID)
}

/*
allDummies tests whether a row is composed of all dummy variables.
 */
func (r *row) allDummies() bool {
	for sym := range r.cells {
		if !sym.is(DUMMY) {
			return false
		}
	}
	return true
}

/*
solveFor solves the row for the given symbol.

This method assumes the row is of the form a * x + b * y + c = 0
and (assuming solve for x) will modify the row to represent the
right hand side of x = -b/a * y - c / a. The target symbol will
be removed from the row, and the constant and other cells will
be multiplied by the negative inverse of the target coefficient.
The given symbol *must* exist in the row.
 */
func (r *row) solveFor(sym *symbol) {
	coeff := -1.0 / r.cells[sym]
	delete(r.cells, sym)
	r.constant *= coeff
	cells := make(map[*symbol]float64)
	for s, c := range r.cells {
		cells[s] = c * coeff
	}
	r.cells = cells
}

/*
solveForPair solves the row for the given symbols.

This method assumes the row is of the form x = b * y + c and will
solve the row such that y = x / b - c / b. The rhs symbol will be
removed from the row, the lhs added, and the result divided by the
negative inverse of the rhs coefficient.
The lhs symbol *must not* exist in the row, and the rhs symbol
must* exist in the row.
 */
func (r *row) solveForPair(lhs, rhs *symbol) {
	r.insertSymbolWithCoefficient(lhs, -1.0)
	r.solveFor(rhs)
}

/*
coefficientFor gets the coefficient for the given symbol.

If the symbol does not exist in the row, zero will be returned.
 */
func (r row) coefficientFor(sym *symbol) float64 {
	if coeff, present := r.cells[sym]; present {
		return coeff
	} else {
		return 0.0
	}
}

/*
substitute substitutes a symbol with the data from another row.

Given a row of the form a * x + b and a substitution of the
form x = 3 * y + c the row will be updated to reflect the
expression 3 * a * y + a * c + b.
If the symbol does not exist in the row, this is a no-op.
 */
func (r *row) substitute(sym *symbol, other *row) {
	if coeff, present := r.cells[sym]; present {
		delete(r.cells, sym)
		r.insertRowWithCoefficient(other, coeff)
	}
}

/*
anyPivotableSymbol gets the first Slack or Error symbol in the row.

If no such symbol is present, and Invalid symbol will be returned.
 */
func (r *row) anyPivotableSymbol() *symbol {
	for sym := range r.cells {
		if sym.is(SLACK) || sym.is(ERROR) {
			return sym
		}
	}
	return newSymbol(INVALID)
}

/*
getEnteringSymbol computes the entering variable for a pivot operation.

This method will return first symbol in the objective function which
is non-dummy and has a coefficient less than zero. If no symbol meets
the criteria, it means the objective function is at a minimum, and an
invalid symbol is returned.
 */
func (r *row) getEnteringSymbol() *symbol {
	objective := r
	for sym, coeff := range objective.cells {
		if !sym.is(DUMMY) && coeff < 0.0 {
			return sym
		}
	}
	return newSymbol(INVALID)
}

/*
getDualEnteringSymbol computes the entering symbol for the dual optimize operation.

This method will return the symbol in the row which has a positive
coefficient and yields the minimum ratio for its respective symbol
in the objective function. The provided row *must* be infeasible.
If no symbol is found which meats the criteria, an invalid symbol
is returned.
 */
func (r *row) getDualEnteringSymbol(other *row) *symbol {
	objective := r
	ratio := math.MaxFloat64
	entering := newSymbol(INVALID)
	for sym, coeff := range other.cells {
		if !sym.is(DUMMY) && coeff > 0.0 {
			r := objective.coefficientFor(sym) / coeff
			if r < ratio {
				ratio = r
				entering = sym
			}
		}
	}
	return entering
}
