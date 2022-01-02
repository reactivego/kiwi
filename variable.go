package kiwi

type Variable struct {
	Name  string
	Value float64
}

var _ Constrainer = &Variable{}

func Var(n string, v ...float64) *Variable {
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

func (v *Variable) LessThanOrEqualsConstant(constant float64) *Constraint {
	return Term{Variable: v, Coefficient: 1.0}.LessThanOrEqualsConstant(constant)
}

func (v *Variable) LessThanOrEqualsVariable(variable *Variable) *Constraint {
	return Term{Variable: v, Coefficient: 1.0}.LessThanOrEqualsVariable(variable)
}

func (v *Variable) LessThanOrEqualsTerm(term Term) *Constraint {
	return Term{Variable: v, Coefficient: 1.0}.LessThanOrEqualsTerm(term)
}

func (v *Variable) LessThanOrEqualsExpression(expression Expression) *Constraint {
	return Term{Variable: v, Coefficient: 1.0}.LessThanOrEqualsExpression(expression)
}

func (v *Variable) GreaterThanOrEqualsConstant(constant float64) *Constraint {
	return Term{Variable: v, Coefficient: 1.0}.GreaterThanOrEqualsConstant(constant)
}

func (v *Variable) GreaterThanOrEqualsVariable(variable *Variable) *Constraint {
	return Term{Variable: v, Coefficient: 1.0}.GreaterThanOrEqualsVariable(variable)
}

func (v *Variable) GreaterThanOrEqualsTerm(term Term) *Constraint {
	return Term{Variable: v, Coefficient: 1.0}.GreaterThanOrEqualsTerm(term)
}

func (v *Variable) GreaterThanOrEqualsExpression(expression Expression) *Constraint {
	return Term{Variable: v, Coefficient: 1.0}.GreaterThanOrEqualsExpression(expression)
}

func (v *Variable) String() string {
	return v.Name
}
