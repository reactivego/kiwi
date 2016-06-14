package kiwi

type Variable interface {
	GetName() string
	GetValue() float64
	SetValue(value float64)

	Multiply(coefficient float64) Term
	Divide(denominator float64) Term
	Negate() Term

	AddConstant(constant float64) Expression
	Add(variable Variable) Expression

	EqualsConstant(constant float64) Constraint
	Equals(variable Variable) Constraint
	EqualsTerm(term Term) Constraint
	EqualsExpression(expression Expression) Constraint

	LessThanOrEqualToConstant(constant float64) Constraint
	LessThanOrEqualTo(variable Variable) Constraint
	LessThanOrEqualToTerm(term Term) Constraint
	LessThanOrEqualToExpression(expression Expression) Constraint
}

func NewVariable(name string) Variable {
	return &variable{name: name}
}

type variable struct {
	name  string
	value float64
}

func (v *variable) GetName() string {
	return v.name
}

func (v *variable) GetValue() float64 {
	return v.value
}

func (v *variable) SetValue(value float64) {
	v.value = value
}

func (v *variable) Multiply(coefficient float64) Term {
	return NewTermFromVariableAndCoefficient(v, coefficient)
}

func (v *variable) Divide(denominator float64) Term {
	return NewTermFromVariableAndCoefficient(v, 1.0/denominator)
}

func (v *variable) Negate() Term {
	return NewTermFromVariableAndCoefficient(v, -1.0)
}

func (v *variable) AddConstant(constant float64) Expression {
	return NewTermFromVariable(v).AddConstant(constant)
}

func (v *variable) Add(variable Variable) Expression {
	return NewTermFromVariable(v).AddVariable(variable)
}

func (v *variable) EqualsConstant(constant float64) Constraint {
	return NewTermFromVariable(v).EqualsConstant(constant)
}

func (v *variable) Equals(variable Variable) Constraint {
	return NewTermFromVariable(v).EqualsVariable(variable)
}

func (v *variable) EqualsTerm(term Term) Constraint {
	return term.EqualsVariable(v)
}

func (v *variable) EqualsExpression(expression Expression) Constraint {
	return expression.EqualsVariable(v)
}

func (v *variable) LessThanOrEqualToConstant(constant float64) Constraint {
	return NewTermFromVariable(v).LessThanOrEqualToConstant(constant)
}

func (v *variable) LessThanOrEqualTo(variable Variable) Constraint {
	return NewTermFromVariable(v).LessThanOrEqualToVariable(variable)
}

func (v *variable) LessThanOrEqualToTerm(term Term) Constraint {
	return NewTermFromVariable(v).LessThanOrEqualTo(term)
}

func (v *variable) LessThanOrEqualToExpression(expression Expression) Constraint {
	return NewTermFromVariable(v).LessThanOrEqualToExpression(expression)
}
