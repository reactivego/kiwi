package kiwi

type Constrainer interface {
	EqualsConstant(float64) *Constraint
	EqualsVariable(*Variable) *Constraint
	EqualsTerm(Term) *Constraint
	EqualsExpression(Expression) *Constraint
	LessThanOrEqualsConstant(float64) *Constraint
	LessThanOrEqualsVariable(*Variable) *Constraint
	LessThanOrEqualsTerm(Term) *Constraint
	LessThanOrEqualsExpression(Expression) *Constraint
	GreaterThanOrEqualsConstant(float64) *Constraint
	GreaterThanOrEqualsVariable(*Variable) *Constraint
	GreaterThanOrEqualsTerm(Term) *Constraint
	GreaterThanOrEqualsExpression(Expression) *Constraint
}
