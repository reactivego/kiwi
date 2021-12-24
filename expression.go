package kiwi

import "github.com/reactivego/kiwi/op"

type Expression interface {
	GetTerms() []Term
	GetConstant() float64
	GetValue() float64
	IsConstant() bool

	Multiply(coefficient float64) Expression
	Divide(denominator float64) Expression
	Negate() Expression

	AddConstant(constant float64) Expression
	AddVariable(variable Variable) Expression
	AddTerm(term Term) Expression
	Add(expression Expression) Expression

	EqualsConstant(constant float64) Constraint
	EqualsVariable(variable Variable) Constraint
	EqualsTerm(term Term) Constraint
	Equals(expression Expression) Constraint

	LessThanOrEqualToConstant(constant float64) Constraint
	LessThanOrEqualToVariable(variable Variable) Constraint
	LessThanOrEqualToTerm(term Term) Constraint
	LessThanOrEqualTo(expression Expression) Constraint

	GreaterThanOrEqualToConstant(constant float64) Constraint
	GreaterThanOrEqualToVariable(variable Variable) Constraint
	GreaterThanOrEqualToTerm(term Term) Constraint
	GreaterThanOrEqualTo(expression Expression) Constraint
}

func NewExpressionFromTermsAndConstant(terms []Term, constant float64) Expression {
	return &expression{terms: terms, constant: constant}
}

func NewExpressionFromConstant(constant float64) Expression {
	return &expression{constant: constant}
}

func NewExpressionFromTermAndConstant(term Term, constant float64) Expression {
	return NewExpressionFromTermsAndConstant([]Term{term}, constant)
}

func NewExpressionFromTerm(term Term) Expression {
	return NewExpressionFromTermAndConstant(term, 0.0)
}

type expression struct {
	terms    []Term
	constant float64
}

func (e expression) GetTerms() []Term {
	return e.terms
}

func (e expression) GetConstant() float64 {
	return e.constant
}

func (e expression) GetValue() float64 {
	result := e.constant
	for _, t := range e.terms {
		result += t.GetValue()
	}
	return result
}

func (e expression) IsConstant() bool {
	return len(e.terms) == 0
}

func (e expression) Multiply(coefficient float64) Expression {
	var terms []Term
	for _, term := range e.terms {
		terms = append(terms, term.Multiply(coefficient))
	}
	return NewExpressionFromTermsAndConstant(terms, e.constant*coefficient)
}

func (e expression) Divide(denominator float64) Expression {
	return e.Multiply(1.0 / denominator)
}

func (e expression) Negate() Expression {
	return e.Multiply(-1.0)
}

func (e expression) AddConstant(constant float64) Expression {
	return NewExpressionFromTermsAndConstant(e.terms, e.constant+constant)
}

func (e expression) AddVariable(variable Variable) Expression {
	return e.AddTerm(NewTermFromVariable(variable))
}

func (e expression) AddTerm(term Term) Expression {
	var terms []Term
	terms = append(terms, e.terms...)
	terms = append(terms, term)
	return NewExpressionFromTermsAndConstant(terms, e.constant)
}

func (e expression) Add(expression Expression) Expression {
	var terms []Term
	terms = append(terms, e.terms...)
	terms = append(terms, expression.GetTerms()...)
	return NewExpressionFromTermsAndConstant(terms, e.constant+expression.GetConstant())
}

func (e expression) EqualsConstant(constant float64) Constraint {
	return NewConstraint(e.AddConstant(-constant), op.EQ)
}

func (e expression) EqualsVariable(variable Variable) Constraint {
	return NewConstraint(e.AddTerm(variable.Negate()), op.EQ)
}

func (e expression) EqualsTerm(term Term) Constraint {
	return NewConstraint(e.AddTerm(term.Negate()), op.EQ)
}

func (e expression) Equals(expression Expression) Constraint {
	return NewConstraint(e.Add(expression.Negate()), op.EQ)
}

func (e expression) LessThanOrEqualToConstant(constant float64) Constraint {
	return NewConstraint(e.AddConstant(-constant), op.LE)
}

func (e expression) LessThanOrEqualToVariable(variable Variable) Constraint {
	return NewConstraint(e.AddTerm(variable.Negate()), op.LE)
}

func (e expression) LessThanOrEqualToTerm(term Term) Constraint {
	return NewConstraint(e.AddTerm(term.Negate()), op.LE)
}

func (e expression) LessThanOrEqualTo(expression Expression) Constraint {
	return NewConstraint(e.Add(expression.Negate()), op.LE)
}

func (e expression) GreaterThanOrEqualToConstant(constant float64) Constraint {
	return NewConstraint(e.AddConstant(-constant), op.GE)
}

func (e expression) GreaterThanOrEqualToVariable(variable Variable) Constraint {
	return NewConstraint(e.AddTerm(variable.Negate()), op.GE)
}

func (e expression) GreaterThanOrEqualToTerm(term Term) Constraint {
	return NewConstraint(e.AddTerm(term.Negate()), op.GE)
}

func (e expression) GreaterThanOrEqualTo(expression Expression) Constraint {
	return NewConstraint(e.Add(expression.Negate()), op.GE)
}
