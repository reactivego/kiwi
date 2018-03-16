package kiwi

import op "kiwi/operator"

type Term interface {
	GetVariable() Variable
	GetCoefficient() float64
	GetValue() float64

	Multiply(coefficient float64) Term
	Divide(denominator float64) Term
	Negate() Term

	AddConstant(constant float64) Expression
	AddVariable(variable Variable) Expression
	Add(term Term) Expression
	AddExpression(expression Expression) Expression

	EqualsConstant(constant float64) Constraint
	EqualsVariable(variable Variable) Constraint
	Equals(term Term) Constraint
	EqualsExpression(expression Expression) Constraint
	LessThanOrEqualToConstant(constant float64) Constraint
	LessThanOrEqualToVariable(variable Variable) Constraint
	LessThanOrEqualTo(term Term) Constraint
	LessThanOrEqualToExpression(expression Expression) Constraint
	GreaterThanOrEqualToConstant(constant float64) Constraint
	GreaterThanOrEqualToVariable(variable Variable) Constraint
	GreaterThanOrEqualTo(term Term) Constraint
	GreaterThanOrEqualToExpression(expression Expression) Constraint
}

func NewTermFromVariableAndCoefficient(variable Variable, coefficient float64) Term {
	return &term{variable: variable, coefficient: coefficient}
}

func NewTermFromVariable(variable Variable) Term {
	return NewTermFromVariableAndCoefficient(variable, 1.0)
}

type term struct {
	variable    Variable
	coefficient float64
}

func (t term) GetVariable() Variable {
	return t.variable
}

func (t term) GetCoefficient() float64 {
	return t.coefficient
}

func (t term) GetValue() float64 {
	return t.coefficient * t.variable.GetValue()
}

func (t term) Multiply(coefficient float64) Term {
	return NewTermFromVariableAndCoefficient(t.variable, t.coefficient*coefficient)
}

func (t term) Divide(denominator float64) Term {
	return NewTermFromVariableAndCoefficient(t.variable, t.coefficient/denominator)
}

func (t term) Negate() Term {
	return NewTermFromVariableAndCoefficient(t.variable, -t.coefficient)
}

func (t term) AddConstant(constant float64) Expression {
	return NewExpressionFromTermAndConstant(t, constant)
}

func (t term) AddVariable(variable Variable) Expression {
	return t.Add(NewTermFromVariable(variable))
}

func (t term) Add(term Term) Expression {
	return NewExpressionFromTermsAndConstant([]Term{t, term}, 0.0)
}

func (t term) AddExpression(expression Expression) Expression {
	return NewExpressionFromTermsAndConstant(append([]Term{t}, expression.GetTerms()...), expression.GetConstant())
}

func (t term) EqualsConstant(constant float64) Constraint {
	return NewConstraint(t.AddConstant(-constant), op.EQ)
}

func (t term) EqualsVariable(variable Variable) Constraint {
	return NewConstraint(t.Add(variable.Negate()), op.EQ)
}

func (t term) Equals(term Term) Constraint {
	return NewConstraint(t.Add(term.Negate()), op.EQ)
}

func (t term) EqualsExpression(expression Expression) Constraint {
	return expression.EqualsTerm(t)
}

func (t term) LessThanOrEqualToConstant(constant float64) Constraint {
	return NewConstraint(t.AddConstant(-constant), op.LE)
}

func (t term) LessThanOrEqualToVariable(variable Variable) Constraint {
	return NewConstraint(t.Add(variable.Negate()), op.LE)
}

func (t term) LessThanOrEqualTo(term Term) Constraint {
	return NewConstraint(t.Add(term.Negate()), op.LE)
}

func (t term) LessThanOrEqualToExpression(expression Expression) Constraint {
	return NewConstraint(t.AddExpression(expression.Negate()), op.LE)
}

func (t term) GreaterThanOrEqualToConstant(constant float64) Constraint {
	return NewConstraint(t.AddConstant(-constant), op.GE)
}

func (t term) GreaterThanOrEqualToVariable(variable Variable) Constraint {
	return NewConstraint(t.Add(variable.Negate()), op.GE)
}

func (t term) GreaterThanOrEqualTo(term Term) Constraint {
	return NewConstraint(t.Add(term.Negate()), op.GE)
}

func (t term) GreaterThanOrEqualToExpression(expression Expression) Constraint {
	return NewConstraint(t.AddExpression(expression.Negate()), op.GE)
}
