package kiwi

import (
	"math"
	"strconv"
)

const (
	// REQUIRED is the strongest strength value indicating the contraint must be satisfied.
	REQUIRED Strength = 1001001000
	// STRONG has the same value as Strong(1)
	STRONG Strength = 1000000
	// MEDIUM has the same value as Medium(1)
	MEDIUM Strength = 1000
	// WEAK has the same value as Weak(1)
	WEAK Strength = 1
	// OPTIONAL is the weakest strength value indicating the contraint is optional.
	OPTIONAL Strength = 0
)

type Strength float64

func (s Strength) String() string {
	switch s {
	case REQUIRED:
		return "REQUIRED"
	case STRONG:
		return "STRONG"
	case MEDIUM:
		return "MEDIUM"
	case WEAK:
		return "WEAK"
	case OPTIONAL:
		return "OPTIONAL"
	default:
		if Strong(1000) >= s && s >= Strong(1) {
			return "Strong(" + strconv.FormatFloat(float64(s/1000000), 'f', -1, 64) + ")"
		} else if Medium(1000) >= s && s >= Medium(1) {
			return "Medium(" + strconv.FormatFloat(float64(s/1000), 'f', -1, 64) + ")"
		} else if Weak(1000) >= s && s >= Weak(1) {
			return "Weak(" + strconv.FormatFloat(float64(s), 'f', -1, 64) + ")"
		}
		return strconv.FormatFloat(float64(s), 'f', -1, 64)
	}
}

func (s Strength) Base() Strength {
	if Strong(1000) >= s && s >= Strong(1) {
		return STRONG
	} else if Medium(1000) >= s && s >= Medium(1) {
		return MEDIUM
	} else if Weak(1000) >= s && s >= Weak(1) {
		return WEAK
	}
	return s
}

func (s Strength) WithWeight(weight float64) Strength {
	if Strong(1000) >= s && s >= Strong(1) {
		return Strong(weight)
	} else if Medium(1000) >= s && s >= Medium(1) {
		return Medium(weight)
	} else if Weak(1000) >= s && s >= Weak(1) {
		return Weak(weight)
	}
	return s
}

// Strong returns a strong strength with a weight in the range [1 .. 1000)
func Strong(weight ...float64) Strength {
	w := append(weight, 1.0)[0]
	return Strength(1000000 * math.Max(1, math.Min(w, 999.9999999999999)))
}

// Medium returns a medium strength with a weight in the range [1 .. 1000)
func Medium(weight ...float64) Strength {
	w := append(weight, 1.0)[0]
	return Strength(1000 * math.Max(1, math.Min(w, 999.9999999999999)))
}

// Weak returns a weak strength with a weight in the range [1 .. 1000)
func Weak(weight ...float64) Strength {
	w := append(weight, 1.0)[0]
	return Strength(1 * math.Max(1, math.Min(w, 999.9999999999999)))
}
