package kiwi

import (
	"math"
	"strconv"
)

const (
	OPTIONAL Strength = iota
	WEAK
	MEDIUM   = (1000 * WEAK)
	STRONG   = (1000 * MEDIUM)
	REQUIRED = (1000 * STRONG) + STRONG + MEDIUM
)

type Strength float64

func (s Strength) String() string {
	switch s {
	case OPTIONAL:
		return "OPTIONAL"
	case WEAK:
		return "WEAK"
	case MEDIUM:
		return "MEDIUM"
	case STRONG:
		return "STRONG"
	case REQUIRED:
		return "REQUIRED"
	default:
		if Weak(1) <= s && s <= Weak(1000) {
			return "Weak(" + strconv.FormatFloat(float64(s), 'f', -1, 64) + ")"
		} else if Medium(1) <= s && s <= Medium(1000) {
			return "Medium(" + strconv.FormatFloat(float64(s/1000), 'f', -1, 64) + ")"
		} else if Strong(1) <= s && s <= Strong(1000) {
			return "Strong(" + strconv.FormatFloat(float64(s/1000000), 'f', -1, 64) + ")"
		} else {
			return strconv.FormatFloat(float64(s), 'f', -1, 64)
		}
	}
}

func (s Strength) Base() Strength {
	if Weak(1) <= s && s <= Weak(1000) {
		return WEAK
	} else if Medium(1) <= s && s <= Medium(1000) {
		return MEDIUM
	} else if Strong(1) <= s && s <= Strong(1000) {
		return STRONG
	} else {
		return s
	}
}

func (s Strength) WithWeight(weight float64) Strength {
	if Weak(1) <= s && s <= Weak(1000) {
		return Weak(weight)
	} else if Medium(1) <= s && s <= Medium(1000) {
		return Medium(weight)
	} else if Strong(1) <= s && s <= Strong(1000) {
		return Strong(weight)
	} else {
		return s
	}
}

// Weak returns a weak strength with a weight in the range [1 .. 1000)
func Weak(weight ...float64) Strength {
	w := append(weight, 1.0)[0]
	return Strength(1 * math.Max(1, math.Min(w, 999.9999999999999)))
}

// Medium returns a medium strength with a weight in the range [1 .. 1000)
func Medium(weight ...float64) Strength {
	w := append(weight, 1.0)[0]
	return Strength(1000 * math.Max(1, math.Min(w, 999.9999999999999)))
}

// Strong returns a strong strength with a weight in the range [1 .. 1000)
func Strong(weight ...float64) Strength {
	w := append(weight, 1.0)[0]
	return Strength(1000000 * math.Max(1, math.Min(w, 999.9999999999999)))
}
