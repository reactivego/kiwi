package kiwi

import "math"

type Constraint struct {
	Expression Expression
	Operator   Operator
	Strength   float64
}

type ConstraintOption func(*Constraint)

// Strength is a constraint option to set the strength of the constraint
func Strength(strength float64) ConstraintOption {
	return func(c *Constraint) {
		c.Strength = math.Max(0, math.Min(strength, REQUIRED))
	}
}

func NewConstraint(expr Expression, op Operator, options ...ConstraintOption) *Constraint {
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
	c := &Constraint{Expression: expr, Operator: op, Strength: REQUIRED}
	for _, opt := range options {
		opt(c)
	}
	return c
}
