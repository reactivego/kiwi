package kiwi

import (
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

	assertEqualsFloat64(t, x.Value, 18, "x =")
}

func TestSimple0(t *testing.T) {

	solver := NewSolver()
	x := NewVariable("x")
	y := NewVariable("y")

	// x = 20
	solver.AddConstraint(x.EqualsConstant(20))
	// x + 2 = y + 10
	solver.AddConstraint(x.AddConstant(2).EqualsExpression(y.AddConstant(10)))

	solver.UpdateVariables()

	assertEqualsFloat64(t, y.Value, 12, "y =")
	assertEqualsFloat64(t, x.Value, 20, "x =")
}

func TestSimple1(t *testing.T) {

	solver := NewSolver()
	x := NewVariable("x")
	y := NewVariable("y")

	solver.AddConstraint(x.EqualsVariable(y))
	solver.UpdateVariables()
	assertEqualsFloat64(t, x.Value, y.Value, "x = y =")
}

func TestSimple2(t *testing.T) {

	solver := NewSolver()
	x := NewVariable("x")
	y := NewVariable("y")

	// x = 27
	solver.AddConstraint(x.EqualsConstant(27))
	// 10 x = 5 y
	solver.AddConstraint(x.Multiply(10).EqualsTerm(y.Multiply(5)))

	solver.UpdateVariables()
	assertEqualsFloat64(t, x.Value, 27, "x =")
	assertEqualsFloat64(t, y.Value, 54, "y =")
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

	assertEqualsFloat64(t, x.Value, 5, "x =")

}

func TestCasso1(t *testing.T) {
	x := NewVariable("x")
	y := NewVariable("y")
	solver := NewSolver()

	// x <= y
	err := solver.AddConstraint(x.LessThanOrEqualToVariable(y))
	if err != nil {
		t.Error(err)
	}
	// y = x + 3.0
	err = solver.AddConstraint(y.EqualsExpression(x.AddConstant(3.0)))
	if err != nil {
		t.Error(err)
	}
	// x = 10.0 [WEAK]
	err = solver.AddConstraint(x.EqualsConstant(10.0), Strength(WEAK))
	if err != nil {
		t.Error(err)
	}
	// y = 10.0 [WEAK]
	err = solver.AddConstraint(y.EqualsConstant(10.0), Strength(WEAK))
	if err != nil {
		t.Error(err)
	}

	solver.UpdateVariables()

	if NearZero(x.Value - 10.0) {
		assertEqualsFloat64(t, x.Value, 10, "x =")
		assertEqualsFloat64(t, y.Value, 13, "y =")
	} else {
		assertEqualsFloat64(t, x.Value, 7, "x =")
		assertEqualsFloat64(t, y.Value, 10, "y =")
	}
}

func TestAddDelete1(t *testing.T) {
	x := NewVariable("x")
	solver := NewSolver()

	solver.AddConstraint(x.LessThanOrEqualToConstant(100.0), Strength(WEAK))
	solver.UpdateVariables()

	assertEqualsFloat64(t, x.Value, 100, "x =")

	c10 := x.LessThanOrEqualToConstant(10.0)
	c20 := x.LessThanOrEqualToConstant(20.0)

	solver.AddConstraint(c10)
	solver.AddConstraint(c20)
	solver.UpdateVariables()

	assertEqualsFloat64(t, x.Value, 10, "x =")

	solver.RemoveConstraint(c10)
	solver.UpdateVariables()

	assertEqualsFloat64(t, x.Value, 20, "x =")

	solver.RemoveConstraint(c20)
	solver.UpdateVariables()

	assertEqualsFloat64(t, x.Value, 100, "x =")

	c10again := x.LessThanOrEqualToConstant(10.0)

	solver.AddConstraint(c10again)
	solver.AddConstraint(c10)
	solver.UpdateVariables()

	assertEqualsFloat64(t, x.Value, 10, "x =")

	solver.RemoveConstraint(c10)
	solver.UpdateVariables()

	assertEqualsFloat64(t, x.Value, 10, "x =")

	solver.RemoveConstraint(c10again)
	solver.UpdateVariables()

	assertEqualsFloat64(t, x.Value, 100, "x =")
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
	if _, typematch := err.(UnsatisfiableConstraint); !typematch {
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
	if _, typematch := err.(UnsatisfiableConstraint); !typematch {
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
	err = solver.AddConstraint(x.GreaterThanOrEqualToVariable(w))
	if err != nil {
		t.Errorf("expected err == nil, got err != nil")
	}
	err = solver.AddConstraint(y.GreaterThanOrEqualToVariable(x))
	if err != nil {
		t.Errorf("expected err == nil, got err != nil")
	}
	err = solver.AddConstraint(z.GreaterThanOrEqualToVariable(y))
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
	if _, typematch := err.(UnsatisfiableConstraint); !typematch {
		t.Errorf("expected UnsatisfiableConstraintException got something else")
	}

	solver.UpdateVariables()
}

func TestStrength(t *testing.T) {
	if REQUIRED != 1001001000 {
		t.Errorf("REQUIRED expected: 1001001000 got: %0.f", REQUIRED)
	}
	if STRONG != 1000000 {
		t.Errorf("STRONG expected:1000000 got: %0.f", STRONG)
	}
	if MEDIUM != 1000 {
		t.Errorf("MEDIUM expected: 1000 got: %0.f", MEDIUM)
	}
	if WEAK != 1 {
		t.Errorf("WEAK expected: 1 got: %0.f", WEAK)
	}
	if OPTIONAL != 0 {
		t.Errorf("OPTIONAL expected: 0 got: %0.f", OPTIONAL)
	}

	if Strong(1) != 1000000 {
		t.Errorf("Strong(1) expected:1000000 got: %0.f", Strong(1))
	}
	if Strong(321) != 321000000 {
		t.Errorf("Strong(321) expected:32100000 got: %0.f", Strong(321))
	}
	if Medium(1) != 1000 {
		t.Errorf("Medium(1) expected: 1000 got: %0.f", Medium(1))
	}
	if Medium(321) != 321000 {
		t.Errorf("Medium(321) expected: 321000 got: %0.f", Medium(321))
	}
	if Weak(1) != 1 {
		t.Errorf("Weak(1) expected: 1 got: %0.f", Weak(1))
	}
	if Weak(321) != 321 {
		t.Errorf("Weak(321) expected: 321 got: %0.f", Weak(321))
	}
}
