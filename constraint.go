package kiwi

import (
	"github.com/reactivego/kiwi/op"
	"github.com/reactivego/kiwi/strength"
)

type Constraint interface {
	GetExpression() Expression
	GetOp() op.Type
	GetStrength() strength.Value
	ModifyStrength(strength strength.Value) Constraint
}

func NewConstraintWithStrength(expr Expression, op op.Type, strength strength.Value) Constraint {
	// reduce: c + pv + qv + rw -> c + (p+q)v + rw
	var vars = make(map[Variable]float64)
	for _, t := range expr.GetTerms() {
		vars[t.GetVariable()] = vars[t.GetVariable()] + t.GetCoefficient()
	}
	var terms []Term
	for v, c := range vars {
		terms = append(terms, NewTermFromVariableAndCoefficient(v, c))
	}
	expr = NewExpressionFromTermsAndConstant(terms, expr.GetConstant())
	return &constraint{expression: expr, op: op, strength: strength.Clip()}
}

func NewConstraint(expr Expression, op op.Type) Constraint {
	return NewConstraintWithStrength(expr, op, strength.REQUIRED)
}

// func CopyConstraintWithStrength(other Constraint, strength Strength) Constraint {
// 	return other.ModifyStrength(strength)
// }

type constraint struct {
	expression Expression
	op         op.Type
	strength   strength.Value
}

func (c constraint) GetExpression() Expression {
	return c.expression
}

func (c constraint) GetOp() op.Type {
	return c.op
}

func (c constraint) GetStrength() strength.Value {
	return c.strength
}

func (c constraint) ModifyStrength(strength strength.Value) Constraint {
	return &constraint{expression: c.expression, op: c.op, strength: strength.Clip()}
}
