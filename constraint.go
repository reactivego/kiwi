package kiwi

type Constraint struct {
	Expression Expression
	Op         Operator
	Strength   Strength
}

func NewConstraintWithStrength(expr Expression, op Operator, strength Strength) *Constraint {
	// reduce: c + pv + qv + rw -> c + (p+q)v + rw
	var vars = make(map[*Variable]float64)
	for _, t := range expr.Terms {
		vars[t.Variable] = vars[t.Variable] + t.Coefficient
	}
	terms := make([]Term, 0, len(vars))
	for v, c := range vars {
		terms = append(terms, Term{Variable: v, Coefficient: c})
	}
	expr = Expression{Terms: terms, Constant: expr.Constant}
	return &Constraint{Expression: expr, Op: op, Strength: strength.Clip()}
}

func NewConstraint(expr Expression, op Operator) *Constraint {
	return NewConstraintWithStrength(expr, op, REQUIRED)
}

func (c *Constraint) ModifyStrength(strength Strength) *Constraint {
	c.Strength = strength.Clip()
	return c
}
