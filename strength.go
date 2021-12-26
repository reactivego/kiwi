package kiwi

import "math"

var (
	REQUIRED = STRENGTH(1000, 1000, 1000, 1)
	STRONG   = STRENGTH(1, 0, 0, 1)
	MEDIUM   = STRENGTH(0, 1, 0, 1)
	WEAK     = STRENGTH(0, 0, 1, 1)
	OPTIONAL = STRENGTH(0, 0, 0, 1)
)

// STRENGTH encodes a higest level Required, a lowest level Optional and 3 strength bands
// inbetween named 'Strong', 'Medium' and 'Weak'.
// Optional() < Weak(1 .. <1000) < Medium(1 .. <1000) < Strong(1 .. <1000) < Required()
func STRENGTH(strong, medium, weak, weight float64) float64 {
	var strength float64
	strength += math.Max(0, math.Min(1000, strong*weight)) * 1000000
	strength += math.Max(0, math.Min(1000, medium*weight)) * 1000
	strength += math.Max(0, math.Min(1000, weak*weight))
	return strength
}

// Strength is a constraint option to set the strength of the constraint
func Strength(strength float64) ConstraintOption {
	return func(c *Constraint) {
		c.Strength = math.Max(0, math.Min(strength, REQUIRED))
	}
}

// RequiredStrength constraint strength is 1001001000.
// RequiredStrength() is the same as Strength(REQUIRED).
func RequiredStrength() ConstraintOption {
	return func(c *Constraint) {
		c.Strength = 1001001000
	}
}

// StrongStrength constraint strength with specific weight in the range [1..1000).
// StrongStrength(1) is the same as Strength(STRONG).
func StrongStrength(weight float64) ConstraintOption {
	return func(c *Constraint) {
		c.Strength = 1000000 * math.Max(1, math.Min(weight, 1000-math.SmallestNonzeroFloat64))
	}
}

// MediumStrength constraint strength with specific weight in the range [1..1000).
// MediumStrength(1) is the same as Strength(MEDIUM).
func MediumStrength(weight float64) ConstraintOption {
	return func(c *Constraint) {
		c.Strength = 1000 * math.Max(1, math.Min(weight, 1000-math.SmallestNonzeroFloat64))
	}
}

// WeakStrength constraint strength with specific weight in the range [1..1000).
// WeakStrength(1) is the same as Strength(WEAK).
func WeakStrength(weight float64) ConstraintOption {
	return func(c *Constraint) {
		c.Strength = 1 * math.Max(1, math.Min(weight, 1000-math.SmallestNonzeroFloat64))
	}
}

// OptionalStrength constraint strength is 0.0.
// OptionalStrength() is the same as Strength(OPTIONAL).
func OptionalStrength() ConstraintOption {
	return func(c *Constraint) {
		c.Strength = 0
	}
}
