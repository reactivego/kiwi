package kiwi

type SolverError string

func (e SolverError) Error() string { return string(e) }

const DuplicateConstraintException = SolverError("Duplicate Constraint")
const InternalSolverError = SolverError("Internal Solver Error")
const UnboundedObjectiveError = SolverError("Objective is Unbounded")
const BadRequiredStrength = SolverError("Bad Required Strength")
const FailedToFindLeavingRow = SolverError("Failed to find Leaving Row")

type UnsatisfiableConstraintException struct{ *Constraint }

func (UnsatisfiableConstraintException) Error() string { return "Unsatisfiable Constraint" }

type UnknownConstraintException struct{ *Constraint }

func (UnknownConstraintException) Error() string { return "Unknown Constraint" }

type UnknownEditVariable struct{ *Variable }

func (UnknownEditVariable) Error() string { return "Unknown Edit Variable" }

type DuplicateEditVariable struct{ *Variable }

func (DuplicateEditVariable) Error() string { return "Duplicate Edit Variable" }
