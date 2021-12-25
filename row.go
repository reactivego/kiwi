package kiwi

import "math"

type Row struct {
	Constant float64
	Cells    map[*Symbol]float64
}

func NewRow() *Row {
	return &Row{Cells: map[*Symbol]float64{}}
}

func NewRowWithConstant(constant float64) *Row {
	return &Row{Constant: constant, Cells: make(map[*Symbol]float64)}
}

func CopyRow(other *Row) *Row {
	// Maps are reference types, we can't just reference the other cell from the other row.
	cells := make(map[*Symbol]float64)
	for s, c := range other.Cells {
		cells[s] = c
	}
	return &Row{Constant: other.Constant, Cells: cells}
}

/**
 * Add a constant value to the row constant.
 *
 * @return The new value of the constant
 */
func (r *Row) add(value float64) float64 {
	r.Constant += value
	return r.Constant
}

/**
 * Insert a symbol into the row with a given coefficient.
 *
 * If the symbol already exists in the row, the coefficient will be
 * added to the existing coefficient. If the resulting coefficient
 * is zero, the symbol will be removed from the row
 */
func (r *Row) InsertSymbolWithCoefficient(sym *Symbol, coeff float64) {
	coeff += r.Cells[sym]
	if NearZero(coeff) {
		delete(r.Cells, sym)
	} else {
		r.Cells[sym] = coeff
	}
}

/**
 * Insert a symbol into the row with coefficient 1.0.
 *
 * If the symbol already exists in the row, the coefficient will be
 * added to the existing coefficient. If the resulting coefficient
 * is zero, the symbol will be removed from the row
 */
func (r *Row) InsertSymbol(sym *Symbol) {
	r.InsertSymbolWithCoefficient(sym, 1.0)
}

/**
 * Insert a row into this row with a given coefficient.
 * The constant and the cells of the other row will be multiplied by
 * the coefficient and added to this row. Any cell with a resulting
 * coefficient of zero will be removed from the row.
 *
 * @param other
 * @param coefficient
 */
func (r *Row) InsertRowWithCoefficient(other *Row, coeff float64) {
	r.Constant += other.Constant * coeff
	for otherSym, otherCoeff := range other.Cells {
		combinedCoeff := r.Cells[otherSym] + otherCoeff*coeff
		if NearZero(combinedCoeff) {
			delete(r.Cells, otherSym)
		} else {
			r.Cells[otherSym] = combinedCoeff
		}
	}
}

/**
 * Insert a row into this row with coefficient 1.0.
 * The constant and the cells of the other row will be multiplied by
 * the coefficient and added to this row. Any cell with a resulting
 * coefficient of zero will be removed from the row.
 *
 * @param other
 */
func (r *Row) InsertRow(other *Row) {
	r.InsertRowWithCoefficient(other, 1.0)
}

/**
 * Remove the given symbol from the row.
 */
func (r *Row) RemoveSymbol(sym *Symbol) {
	delete(r.Cells, sym)
}

/**
 * Reverse the sign of the constant and all cells in the row.
 */
func (r *Row) ReverseSign() {
	r.Constant = -r.Constant
	cells := make(map[*Symbol]float64)
	for s, c := range r.Cells {
		cells[s] = -c
	}
	r.Cells = cells
}

/**
 * Choose the subject for solving for the row
 *
 * This method will choose the best subject for using as the solve
 * target for the row. An invalid symbol will be returned if there
 * is no valid target.
 * The symbols are chosen according to the following precedence:
 * 1) The first symbol representing an external variable.
 * 2) A negative slack or error tag variable.
 * If a subject cannot be found, an invalid symbol will be returned.
 */
func (r *Row) ChooseSubject(tag tag) *Symbol {
	for sym := range r.Cells {
		if sym.IsExternal() {
			return sym
		}
	}

	if tag.marker.IsSlack() || tag.marker.IsError() {
		if r.CoefficientFor(tag.marker) < 0.0 {
			return tag.marker
		}
	}

	if tag.other != nil && (tag.other.IsSlack() || tag.other.IsError()) {
		if r.CoefficientFor(tag.other) < 0.0 {
			return tag.other
		}
	}

	return NewInvalidSymbol()
}

/**
 * Test whether a row is composed of all dummy variables.
 */
func (r *Row) AllDummies() bool {
	for sym := range r.Cells {
		if !sym.IsDummy() {
			return false
		}
	}
	return true
}

/**
 * Solve the row for the given symbol.
 *
 * This method assumes the row is of the form a * x + b * y + c = 0
 * and (assuming solve for x) will modify the row to represent the
 * right hand side of x = -b/a * y - c / a. The target symbol will
 * be removed from the row, and the constant and other cells will
 * be multiplied by the negative inverse of the target coefficient.
 * The given symbol *must* exist in the row.
 *
 * @param symbol
 */
func (r *Row) SolveFor(sym *Symbol) {
	coeff := -1.0 / r.Cells[sym]
	delete(r.Cells, sym)
	r.Constant *= coeff
	cells := make(map[*Symbol]float64)
	for s, c := range r.Cells {
		cells[s] = c * coeff
	}
	r.Cells = cells
}

/**
 * Solve the row for the given symbols.
 *
 * This method assumes the row is of the form x = b * y + c and will
 * solve the row such that y = x / b - c / b. The rhs symbol will be
 * removed from the row, the lhs added, and the result divided by the
 * negative inverse of the rhs coefficient.
 * The lhs symbol *must not* exist in the row, and the rhs symbol
 * must* exist in the row.
 *
 * @param lhs
 * @param rhs
 */
func (r *Row) SolveForPair(lhs, rhs *Symbol) {
	r.InsertSymbolWithCoefficient(lhs, -1.0)
	r.SolveFor(rhs)
}

/**
 * Get the coefficient for the given symbol.
 *
 * If the symbol does not exist in the row, zero will be returned.
 *
 * @return
 */
func (r Row) CoefficientFor(sym *Symbol) float64 {
	if coeff, present := r.Cells[sym]; present {
		return coeff
	} else {
		return 0.0
	}
}

/**
 * Substitute a symbol with the data from another row.
 *
 * Given a row of the form a * x + b and a substitution of the
 * form x = 3 * y + c the row will be updated to reflect the
 * expression 3 * a * y + a * c + b.
 * If the symbol does not exist in the row, this is a no-op.
 */
func (r *Row) Substitute(sym *Symbol, row *Row) {
	if coeff, present := r.Cells[sym]; present {
		delete(r.Cells, sym)
		r.InsertRowWithCoefficient(row, coeff)
	}
}

/**
 * Get the first Slack or Error symbol in the row.
 *
 * If no such symbol is present, and Invalid symbol will be returned.
 */
func (r *Row) AnyPivotableSymbol() *Symbol {
	for sym := range r.Cells {
		if sym.IsSlack() || sym.IsError() {
			return sym
		}
	}
	return NewInvalidSymbol()
}

/**
 * Compute the entering variable for a pivot operation.
 *
 * This method will return first symbol in the objective function which
 * is non-dummy and has a coefficient less than zero. If no symbol meets
 * the criteria, it means the objective function is at a minimum, and an
 * invalid symbol is returned.
 */
func (r *Row) GetEnteringSymbol() *Symbol {
	objective := r
	for sym, coeff := range objective.Cells {
		if !sym.IsDummy() && coeff < 0.0 {
			return sym
		}
	}
	return NewInvalidSymbol()
}

/**
 * Compute the entering symbol for the dual optimize operation.
 *
 * This method will return the symbol in the row which has a positive
 * coefficient and yields the minimum ratio for its respective symbol
 * in the objective function. The provided row *must* be infeasible.
 * If no symbol is found which meats the criteria, an invalid symbol
 * is returned.
 */
func (r *Row) GetDualEnteringSymbol(row *Row) *Symbol {
	objective := r
	ratio := math.MaxFloat64
	entering := NewInvalidSymbol()
	for sym, coeff := range row.Cells {
		if !sym.IsDummy() && coeff > 0.0 {
			r := objective.CoefficientFor(sym) / coeff
			if r < ratio {
				ratio = r
				entering = sym
			}
		}
	}
	return entering
}
