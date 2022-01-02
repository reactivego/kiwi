package kiwi

import (
	"fmt"
	"strings"
)

type Expression struct {
	Terms    []Term
	Constant float64
}

var _ Constrainer = Expression{}

func (e Expression) GetValue() float64 {
	result := e.Constant
	for _, t := range e.Terms {
		result += t.GetValue()
	}
	return result
}

func (e Expression) IsConstant() bool {
	return len(e.Terms) == 0
}

func (e Expression) Multiply(coefficient float64) Expression {
	terms := make([]Term, len(e.Terms))
	for i := range e.Terms {
		terms[i] = e.Terms[i].Multiply(coefficient)
	}
	return Expression{Terms: terms, Constant: e.Constant * coefficient}
}

func (e Expression) Divide(denominator float64) Expression {
	return e.Multiply(1.0 / denominator)
}

func (e Expression) Negate() Expression {
	return e.Multiply(-1.0)
}

func (e Expression) AddConstant(constant float64) Expression {
	return Expression{Terms: e.Terms, Constant: e.Constant + constant}
}

func (e Expression) AddVariable(variable *Variable) Expression {
	return e.AddTerm(Term{Variable: variable, Coefficient: 1.0})
}

func (e Expression) AddTerm(term Term) Expression {
	terms := make([]Term, len(e.Terms)+1)
	n := copy(terms, e.Terms)
	terms[n] = term
	return Expression{Terms: terms, Constant: e.Constant}
}

func (e Expression) AddExpression(other Expression) Expression {
	terms := make([]Term, len(e.Terms)+len(other.Terms))
	n := copy(terms, e.Terms)
	copy(terms[n:], other.Terms)
	return Expression{Terms: terms, Constant: e.Constant + other.Constant}
}

func (e Expression) EqualsConstant(constant float64) *Constraint {
	return NewConstraint(e.AddConstant(-constant), EQ)
}

func (e Expression) EqualsVariable(variable *Variable) *Constraint {
	return NewConstraint(e.AddTerm(variable.Negate()), EQ)
}

func (e Expression) EqualsTerm(term Term) *Constraint {
	return NewConstraint(e.AddTerm(term.Negate()), EQ)
}

func (e Expression) EqualsExpression(expression Expression) *Constraint {
	return NewConstraint(e.AddExpression(expression.Negate()), EQ)
}

func (e Expression) LessThanOrEqualToConstant(constant float64) *Constraint {
	return NewConstraint(e.AddConstant(-constant), LE)
}

func (e Expression) LessThanOrEqualToVariable(variable *Variable) *Constraint {
	return NewConstraint(e.AddTerm(variable.Negate()), LE)
}

func (e Expression) LessThanOrEqualToTerm(term Term) *Constraint {
	return NewConstraint(e.AddTerm(term.Negate()), LE)
}

func (e Expression) LessThanOrEqualToExpression(expression Expression) *Constraint {
	return NewConstraint(e.AddExpression(expression.Negate()), LE)
}

func (e Expression) GreaterThanOrEqualToConstant(constant float64) *Constraint {
	return NewConstraint(e.AddConstant(-constant), GE)
}

func (e Expression) GreaterThanOrEqualToVariable(variable *Variable) *Constraint {
	return NewConstraint(e.AddTerm(variable.Negate()), GE)
}

func (e Expression) GreaterThanOrEqualToTerm(term Term) *Constraint {
	return NewConstraint(e.AddTerm(term.Negate()), GE)
}

func (e Expression) GreaterThanOrEqualToExpression(expression Expression) *Constraint {
	return NewConstraint(e.AddExpression(expression.Negate()), GE)
}

func (e Expression) String() string {
	var factors []string
	for _, t := range e.Terms {
		factors = append(factors, fmt.Sprint(t))
	}
	factors = append(factors, fmt.Sprint(e.Constant))
	return strings.Join(factors, " + ")
}
