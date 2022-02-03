// SPDX-License-Identifier: BSD-3-Clause

package kiwi

import "strconv"

type Term struct {
	Variable    *Variable
	Coefficient float64
}

var _ Constrainer = Term{}

func (t Term) GetValue() float64 {
	return t.Coefficient * t.Variable.Value
}

func (t Term) Multiply(coefficient float64) Term {
	return Term{Variable: t.Variable, Coefficient: t.Coefficient * coefficient}
}

func (t Term) Divide(denominator float64) Term {
	return Term{Variable: t.Variable, Coefficient: t.Coefficient / denominator}
}

func (t Term) Negate() Term {
	return Term{Variable: t.Variable, Coefficient: -t.Coefficient}
}

func (t Term) AddConstant(constant float64) Expression {
	return Expression{Terms: []Term{t}, Constant: constant}
}

func (t Term) AddVariable(variable *Variable) Expression {
	other := Term{Variable: variable, Coefficient: 1.0}
	return Expression{Terms: []Term{t, other}, Constant: 0.0}
}

func (t Term) AddTerm(term Term) Expression {
	return Expression{Terms: []Term{t, term}, Constant: 0.0}
}

func (t Term) AddExpression(expression Expression) Expression {
	terms := make([]Term, 1+len(expression.Terms))
	terms[0] = t
	copy(terms[1:], expression.Terms)
	return Expression{Terms: terms, Constant: expression.Constant}
}

func (t Term) EqualsConstant(constant float64) *Constraint {
	return NewConstraint(t.AddConstant(-constant), EQ)
}

func (t Term) EqualsVariable(variable *Variable) *Constraint {
	return NewConstraint(t.AddTerm(variable.Negate()), EQ)
}

func (t Term) EqualsTerm(term Term) *Constraint {
	return NewConstraint(t.AddTerm(term.Negate()), EQ)
}

func (t Term) EqualsExpression(expression Expression) *Constraint {
	return expression.EqualsTerm(t)
}

func (t Term) LessThanOrEqualsConstant(constant float64) *Constraint {
	return NewConstraint(t.AddConstant(-constant), LE)
}

func (t Term) LessThanOrEqualsVariable(variable *Variable) *Constraint {
	return NewConstraint(t.AddTerm(variable.Negate()), LE)
}

func (t Term) LessThanOrEqualsTerm(term Term) *Constraint {
	return NewConstraint(t.AddTerm(term.Negate()), LE)
}

func (t Term) LessThanOrEqualsExpression(expression Expression) *Constraint {
	return NewConstraint(t.AddExpression(expression.Negate()), LE)
}

func (t Term) GreaterThanOrEqualsConstant(constant float64) *Constraint {
	return NewConstraint(t.AddConstant(-constant), GE)
}

func (t Term) GreaterThanOrEqualsVariable(variable *Variable) *Constraint {
	return NewConstraint(t.AddTerm(variable.Negate()), GE)
}

func (t Term) GreaterThanOrEqualsTerm(term Term) *Constraint {
	return NewConstraint(t.AddTerm(term.Negate()), GE)
}

func (t Term) GreaterThanOrEqualsExpression(expression Expression) *Constraint {
	return NewConstraint(t.AddExpression(expression.Negate()), GE)
}

func (t Term) String() string {
	if t.Coefficient == 1.0 {
		return t.Variable.String()
	}
	return strconv.FormatFloat(t.Coefficient, 'f', -1, 64) + " * " + t.Variable.String()
}
