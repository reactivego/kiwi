package kiwi

import "math"

var (
	REQUIRED = STRENGTH(1000, 1000, 1000)
	STRONG   = STRENGTH(1, 0, 0)
	MEDIUM   = STRENGTH(0, 1, 0)
	WEAK     = STRENGTH(0, 0, 1)
	OPTIONAL = STRENGTH(0, 0, 0)
)

// STRENGTH encodes 3 strength bands named 'Strong', 'Medium' and 'Weak'.
// There are 2 special values 'Required' and 'Optional' that are at the extreme
// edges of the strength spectrum.
// 	Required > Strong(1 .. <1000) > Medium(1 .. <1000) > Weak(1 .. <1000) > Optional
func STRENGTH(strong, medium, weak float64) float64 {
	var strength float64
	strength += math.Max(0, math.Min(1000, strong)) * 1000000
	strength += math.Max(0, math.Min(1000, medium)) * 1000
	strength += math.Max(0, math.Min(1000, weak))
	return strength
}

// Strong returns a strong strength with a weight in the range [1 .. 1000)
func Strong(weight ...float64) float64 {
	w := append(weight, 1.0)[0]
	return 1000000 * math.Max(1, math.Min(w, 1000-math.SmallestNonzeroFloat64))
}

// Medium returns a medium strength with a weight in the range [1 .. 1000)
func Medium(weight ...float64) float64 {
	w := append(weight, 1.0)[0]
	return 1000 * math.Max(1, math.Min(w, 1000-math.SmallestNonzeroFloat64))
}

// Weak returns a weak strength with a weight in the range [1 .. 1000)
func Weak(weight ...float64) float64 {
	w := append(weight, 1.0)[0]
	return 1 * math.Max(1, math.Min(w, 1000-math.SmallestNonzeroFloat64))
}
