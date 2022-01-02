package kiwi

type Constrainer interface {
	EqualsConstant(constant float64) *Constraint
	EqualsVariable(variable *Variable) *Constraint
	EqualsTerm(term Term) *Constraint
	EqualsExpression(expression Expression) *Constraint
	LessThanOrEqualsConstant(constant float64) *Constraint
	LessThanOrEqualsVariable(variable *Variable) *Constraint
	LessThanOrEqualsTerm(term Term) *Constraint
	LessThanOrEqualsExpression(expression Expression) *Constraint
	GreaterThanOrEqualsConstant(constant float64) *Constraint
	GreaterThanOrEqualsVariable(variable *Variable) *Constraint
	GreaterThanOrEqualsTerm(term Term) *Constraint
	GreaterThanOrEqualsExpression(expression Expression) *Constraint
}
