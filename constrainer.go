package kiwi

type Constrainer interface {
	EqualsConstant(constant float64) *Constraint
	EqualsVariable(variable *Variable) *Constraint
	EqualsTerm(term Term) *Constraint
	EqualsExpression(expression Expression) *Constraint
	LessThanOrEqualToConstant(constant float64) *Constraint
	LessThanOrEqualToVariable(variable *Variable) *Constraint
	LessThanOrEqualToTerm(term Term) *Constraint
	LessThanOrEqualToExpression(expression Expression) *Constraint
	GreaterThanOrEqualToConstant(constant float64) *Constraint
	GreaterThanOrEqualToVariable(variable *Variable) *Constraint
	GreaterThanOrEqualToTerm(term Term) *Constraint
	GreaterThanOrEqualToExpression(expression Expression) *Constraint
}
