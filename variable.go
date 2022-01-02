package kiwi

type Variable struct {
	Name  string
	Value float64
}

var _ Constrainer = &Variable{}

func V(n string, v ...float64) *Variable {
	return &Variable{n, append(v, 0.0)[0]}
}

func NewVariable(name string) *Variable {
	return &Variable{Name: name}
}

func (v *Variable) Multiply(coefficient float64) Term {
	return Term{Variable: v, Coefficient: coefficient}
}

func (v *Variable) Divide(denominator float64) Term {
	return Term{Variable: v, Coefficient: 1.0 / denominator}
}

func (v *Variable) Negate() Term {
	return Term{Variable: v, Coefficient: -1.0}
}

func (v *Variable) AddConstant(constant float64) Expression {
	return Term{Variable: v, Coefficient: 1.0}.AddConstant(constant)
}

func (v *Variable) AddVariable(variable *Variable) Expression {
	return Term{Variable: v, Coefficient: 1.0}.AddVariable(variable)
}

func (v *Variable) EqualsConstant(constant float64) *Constraint {
	return Term{Variable: v, Coefficient: 1.0}.EqualsConstant(constant)
}

func (v *Variable) EqualsVariable(variable *Variable) *Constraint {
	return Term{Variable: v, Coefficient: 1.0}.EqualsVariable(variable)
}

func (v *Variable) EqualsTerm(term Term) *Constraint {
	return term.EqualsVariable(v)
}

func (v *Variable) EqualsExpression(expression Expression) *Constraint {
	return expression.EqualsVariable(v)
}

func (v *Variable) LessThanOrEqualToConstant(constant float64) *Constraint {
	return Term{Variable: v, Coefficient: 1.0}.LessThanOrEqualToConstant(constant)
}

func (v *Variable) LessThanOrEqualToVariable(variable *Variable) *Constraint {
	return Term{Variable: v, Coefficient: 1.0}.LessThanOrEqualToVariable(variable)
}

func (v *Variable) LessThanOrEqualToTerm(term Term) *Constraint {
	return Term{Variable: v, Coefficient: 1.0}.LessThanOrEqualToTerm(term)
}

func (v *Variable) LessThanOrEqualToExpression(expression Expression) *Constraint {
	return Term{Variable: v, Coefficient: 1.0}.LessThanOrEqualToExpression(expression)
}

func (v *Variable) GreaterThanOrEqualToConstant(constant float64) *Constraint {
	return Term{Variable: v, Coefficient: 1.0}.GreaterThanOrEqualToConstant(constant)
}

func (v *Variable) GreaterThanOrEqualToVariable(variable *Variable) *Constraint {
	return Term{Variable: v, Coefficient: 1.0}.GreaterThanOrEqualToVariable(variable)
}

func (v *Variable) GreaterThanOrEqualToTerm(term Term) *Constraint {
	return Term{Variable: v, Coefficient: 1.0}.GreaterThanOrEqualToTerm(term)
}

func (v *Variable) GreaterThanOrEqualToExpression(expression Expression) *Constraint {
	return Term{Variable: v, Coefficient: 1.0}.GreaterThanOrEqualToExpression(expression)
}
