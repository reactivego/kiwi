package kiwi

type DuplicateConstraintException struct{}

func (DuplicateConstraintException) Error() string {
	return "DuplicateConstraintException"
}

type InternalSolverError struct {
	msg string
}

func (e InternalSolverError) Error() string {
	return e.msg
}

type UnsatisfiableConstraintException struct {
	constraint Constraint
}

func (e UnsatisfiableConstraintException) Error() string {
	return "UnsatisfiableConstraintException"
}

type UnknownConstraintException struct {
	constraint Constraint
}

func (e UnknownConstraintException) Error() string {
	return "UnknownConstraintException"
}
