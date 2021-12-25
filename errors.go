package kiwi

type SolverException string

func (e SolverException) Error() string { return string(e) }

const DuplicateConstraintException = SolverException("DuplicateConstraintException")
const InternalSolverError = SolverException("internal solver error")
const UnboundedObjectiveError = SolverException("The objective is unbounded.")

type UnsatisfiableConstraintException struct{ *Constraint }

func (UnsatisfiableConstraintException) Error() string { return "UnsatisfiableConstraintException" }

type UnknownConstraintException struct{ *Constraint }

func (UnknownConstraintException) Error() string { return "UnknownConstraintException" }
