package kiwi

import "fmt"

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

func (t Term) Add(term Term) Expression {
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
	return NewConstraint(t.Add(variable.Negate()), EQ)
}

func (t Term) EqualsTerm(term Term) *Constraint {
	return NewConstraint(t.Add(term.Negate()), EQ)
}

func (t Term) EqualsExpression(expression Expression) *Constraint {
	return expression.EqualsTerm(t)
}

func (t Term) LessThanOrEqualToConstant(constant float64) *Constraint {
	return NewConstraint(t.AddConstant(-constant), LE)
}

func (t Term) LessThanOrEqualToVariable(variable *Variable) *Constraint {
	return NewConstraint(t.Add(variable.Negate()), LE)
}

func (t Term) LessThanOrEqualToTerm(term Term) *Constraint {
	return NewConstraint(t.Add(term.Negate()), LE)
}

func (t Term) LessThanOrEqualToExpression(expression Expression) *Constraint {
	return NewConstraint(t.AddExpression(expression.Negate()), LE)
}

func (t Term) GreaterThanOrEqualToConstant(constant float64) *Constraint {
	return NewConstraint(t.AddConstant(-constant), GE)
}

func (t Term) GreaterThanOrEqualToVariable(variable *Variable) *Constraint {
	return NewConstraint(t.Add(variable.Negate()), GE)
}

func (t Term) GreaterThanOrEqualToTerm(term Term) *Constraint {
	return NewConstraint(t.Add(term.Negate()), GE)
}

func (t Term) GreaterThanOrEqualToExpression(expression Expression) *Constraint {
	return NewConstraint(t.AddExpression(expression.Negate()), GE)
}

func (t Term) String() string {
	return fmt.Sprint(t.Coefficient, " * ", t.Variable)
}
