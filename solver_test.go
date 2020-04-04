package kiwi

import (
	"kiwi/strength"
	"runtime"
	"testing"
)

// assertEqualsFloat64() calls testing.T.Error() with the given message if
// the given float64s are not equal.
func assertEqualsFloat64(t *testing.T, got, expect float64, message string) {
	_, _, line, _ := runtime.Caller(1)
	if !NearZero(got - expect) {
		t.Errorf("line %d: %s expected %g, got %g", line, message, expect, got)
	}
}

func TestSimpleNew(t *testing.T) {

	solver := NewSolver()
	x := NewVariable("x")

	// x + 2 = 20
	solver.AddConstraint(x.AddConstant(2).EqualsConstant(20))

	solver.UpdateVariables()

	assertEqualsFloat64(t, x.GetValue(), 18, "x =")
}

func TestSimple0(t *testing.T) {

	solver := NewSolver()
	x := NewVariable("x")
	y := NewVariable("y")

	// x = 20
	solver.AddConstraint(x.EqualsConstant(20))
	// x + 2 = y + 10
	solver.AddConstraint(x.AddConstant(2).Equals(y.AddConstant(10)))

	solver.UpdateVariables()

	assertEqualsFloat64(t, y.GetValue(), 12, "y =")
	assertEqualsFloat64(t, x.GetValue(), 20, "x =")
}

func TestSimple1(t *testing.T) {

	solver := NewSolver()
	x := NewVariable("x")
	y := NewVariable("y")

	solver.AddConstraint(x.Equals(y))
	solver.UpdateVariables()
	assertEqualsFloat64(t, x.GetValue(), y.GetValue(), "x = y =")
}

func TestSimple2(t *testing.T) {

	solver := NewSolver()
	x := NewVariable("x")
	y := NewVariable("y")

	// x = 27
	solver.AddConstraint(x.EqualsConstant(27))
	// 10 x = 5 y
	solver.AddConstraint(x.Multiply(10).Equals(y.Multiply(5)))

	solver.UpdateVariables()
	assertEqualsFloat64(t, x.GetValue(), 27, "x =")
	assertEqualsFloat64(t, y.GetValue(), 54, "y =")
}

func TestCasso0(t *testing.T) {
	x := NewVariable("x")
	solver := NewSolver()

	// x <= 10
	err := solver.AddConstraint(x.LessThanOrEqualToConstant(10.0))
	if err != nil {
		t.Error(err)
	}

	// x = 5
	err = solver.AddConstraint(x.EqualsConstant(5.0))
	if err != nil {
		t.Error(err)
	}

	solver.UpdateVariables()

	assertEqualsFloat64(t, x.GetValue(), 5, "x =")

}

func TestCasso1(t *testing.T) {
	x := NewVariable("x")
	y := NewVariable("y")
	solver := NewSolver()

	// x <= y
	err := solver.AddConstraint(x.LessThanOrEqualTo(y))
	if err != nil {
		t.Error(err)
	}
	// y = x + 3.0
	err = solver.AddConstraint(y.EqualsExpression(x.AddConstant(3.0)))
	if err != nil {
		t.Error(err)
	}
	// x = 10.0 [WEAK]
	err = solver.AddConstraint(x.EqualsConstant(10.0).ModifyStrength(strength.WEAK))
	if err != nil {
		t.Error(err)
	}
	// y = 10.0 [WEAK]
	err = solver.AddConstraint(y.EqualsConstant(10.0).ModifyStrength(strength.WEAK))
	if err != nil {
		t.Error(err)
	}

	solver.UpdateVariables()

	if NearZero(x.GetValue() - 10.0) {
		assertEqualsFloat64(t, x.GetValue(), 10, "x =")
		assertEqualsFloat64(t, y.GetValue(), 13, "y =")
	} else {
		assertEqualsFloat64(t, x.GetValue(), 7, "x =")
		assertEqualsFloat64(t, y.GetValue(), 10, "y =")
	}
}

func TestAddDelete1(t *testing.T) {
	x := NewVariable("x")
	solver := NewSolver()

	solver.AddConstraint(x.LessThanOrEqualToConstant(100.0).ModifyStrength(strength.WEAK))
	solver.UpdateVariables()

	assertEqualsFloat64(t, x.GetValue(), 100, "x =")

	c10 := x.LessThanOrEqualToConstant(10.0)
	c20 := x.LessThanOrEqualToConstant(20.0)

	solver.AddConstraint(c10)
	solver.AddConstraint(c20)
	solver.UpdateVariables()

	assertEqualsFloat64(t, x.GetValue(), 10, "x =")

	solver.RemoveConstraint(c10)
	solver.UpdateVariables()

	assertEqualsFloat64(t, x.GetValue(), 20, "x =")

	solver.RemoveConstraint(c20)
	solver.UpdateVariables()

	assertEqualsFloat64(t, x.GetValue(), 100, "x =")

	c10again := x.LessThanOrEqualToConstant(10.0)

	solver.AddConstraint(c10again)
	solver.AddConstraint(c10)
	solver.UpdateVariables()

	assertEqualsFloat64(t, x.GetValue(), 10, "x =")

	solver.RemoveConstraint(c10)
	solver.UpdateVariables()

	assertEqualsFloat64(t, x.GetValue(), 10, "x =")

	solver.RemoveConstraint(c10again)
	solver.UpdateVariables()

	assertEqualsFloat64(t, x.GetValue(), 100, "x =")
}

func TestInconsistent1(t *testing.T) {
	x := NewVariable("x")
	solver := NewSolver()

	err := solver.AddConstraint(x.EqualsConstant(10.0))
	if err != nil {
		t.Errorf("expected err == nil, got err != nil")
	}
	err = solver.AddConstraint(x.EqualsConstant(5.0))
	if err == nil {
		t.Errorf("expected err != nil, got err == nil")
	}
	if _, typematch := err.(UnsatisfiableConstraintException); !typematch {
		t.Errorf("expected typematch == true got false")
	}

	solver.UpdateVariables()
}

func TestInconsistent2(t *testing.T) {
	x := NewVariable("x")
	solver := NewSolver()

	err := solver.AddConstraint(x.GreaterThanOrEqualToConstant(10.0))
	if err != nil {
		t.Errorf("expected err == nil, got err != nil")
	}
	err = solver.AddConstraint(x.LessThanOrEqualToConstant(5.0))
	if err == nil {
		t.Errorf("expected err != nil, got err == nil")
	}
	if _, typematch := err.(UnsatisfiableConstraintException); !typematch {
		t.Errorf("expected UnsatisfiableConstraintException got something else")
	}

	solver.UpdateVariables()
}

func TestInconsistent3(t *testing.T) {
	w := NewVariable("w")
	x := NewVariable("x")
	y := NewVariable("y")
	z := NewVariable("z")
	solver := NewSolver()

	err := solver.AddConstraint(w.GreaterThanOrEqualToConstant(10.0))
	if err != nil {
		t.Errorf("expected err == nil, got err != nil")
	}
	err = solver.AddConstraint(x.GreaterThanOrEqualTo(w))
	if err != nil {
		t.Errorf("expected err == nil, got err != nil")
	}
	err = solver.AddConstraint(y.GreaterThanOrEqualTo(x))
	if err != nil {
		t.Errorf("expected err == nil, got err != nil")
	}
	err = solver.AddConstraint(z.GreaterThanOrEqualTo(y))
	if err != nil {
		t.Errorf("expected err == nil, got err != nil")
	}
	err = solver.AddConstraint(z.GreaterThanOrEqualToConstant(8.0))
	if err != nil {
		t.Errorf("expected err == nil, got err != nil")
	}
	err = solver.AddConstraint(z.LessThanOrEqualToConstant(4.0))
	if err == nil {
		t.Errorf("expected err != nil, got err == nil")
	}
	if _, typematch := err.(UnsatisfiableConstraintException); !typematch {
		t.Errorf("expected UnsatisfiableConstraintException got something else")
	}

	solver.UpdateVariables()
}
