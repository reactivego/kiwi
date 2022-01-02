package kiwi

type SolverError string

func (e SolverError) Error() string { return string(e) }

const InternalSolverError = SolverError("Internal Solver Error")
const UnboundedObjective = SolverError("Objective is Unbounded")
const BadRequiredStrength = SolverError("Bad Required Strength")
const FailedToFindLeavingRow = SolverError("Failed to find Leaving Row")

type DuplicateConstraint struct{ *Constraint }

func (DuplicateConstraint) Error() string { return "Duplicate Constraint" }

type UnsatisfiableConstraint struct{ *Constraint }

func (UnsatisfiableConstraint) Error() string { return "Unsatisfiable Constraint" }

type UnknownConstraint struct{ *Constraint }

func (UnknownConstraint) Error() string { return "Unknown Constraint" }

type UnknownEditVariable struct{ *Variable }

func (UnknownEditVariable) Error() string { return "Unknown Edit Variable" }

type DuplicateEditVariable struct{ *Variable }

func (DuplicateEditVariable) Error() string { return "Duplicate Edit Variable" }
