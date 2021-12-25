package kiwi

import "math"

// Strength encodes 2 strenght levels Required (highest) Optional (lowest) and 3 strength
// bands inbetween these levels named: 'Strong', 'Medium' and 'Weak' that can be weighted
// individually.
//
// Optional() < Weak(1 .. <1000) < Medium(1 .. <1000) < Strong(1 .. <1000) < Required()
type Strength float64

var (
	REQUIRED = Create(1000, 1000, 1000, 1)
	STRONG   = Create(1, 0, 0, 1)
	MEDIUM   = Create(0, 1, 0, 1)
	WEAK     = Create(0, 0, 1, 1)
	OPTIONAL = Create(0, 0, 0, 1)
)

func Create(a, b, c, w float64) Strength {
	var strength float64
	strength += math.Max(0, math.Min(1000, a*w)) * 1000000
	strength += math.Max(0, math.Min(1000, b*w)) * 1000
	strength += math.Max(0, math.Min(1000, c*w))
	return Strength(strength)
}

// Required strength is 1001001000
func Required() Strength {
	return Strength(1001001000)
}

// Strong strength with specific weight in the range [1..1000)
func Strong(weight float64) Strength {
	return Strength(1000000 * math.Max(1, math.Min(weight, 1000-math.SmallestNonzeroFloat64)))
}

// Medium strength with specific weight in the range [1..1000)
func Medium(weight float64) Strength {
	return Strength(1000 * math.Max(1, math.Min(weight, 1000-math.SmallestNonzeroFloat64)))
}

// Weak strength with specific weight in the range [1..1000)
func Weak(weight float64) Strength {
	return Strength(1 * math.Max(1, math.Min(weight, 1000-math.SmallestNonzeroFloat64)))
}

// Optional strength is 0.0
func Optional() Strength {
	return Strength(0)
}

func (s Strength) Clip() Strength {
	return Strength(math.Max(0, math.Min(float64(s), float64(REQUIRED))))
}
