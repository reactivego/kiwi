package kiwi

import (
	"fmt"
	"math"
)

type Constraint struct {
	Expression Expression
	Operator   Operator
	Strength   Strength
}

type ConstraintOption func(*Constraint)

// WithStrength is a constraint option to set the strength of the constraint
func WithStrength(strength Strength) ConstraintOption {
	return func(c *Constraint) {
		optional, strength, required := float64(OPTIONAL), float64(strength), float64(REQUIRED)
		c.Strength = Strength(math.Max(optional, math.Min(strength, required)))
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

func (c *Constraint) String() string {
	return fmt.Sprint(c.Expression, c.Operator, " 0 | Strength = ", c.Strength)
}
