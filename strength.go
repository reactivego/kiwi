// SPDX-License-Identifier: BSD-3-Clause

package kiwi

import (
	"math"
	"strconv"
)

type Strength float64

const (
	OPTIONAL Strength = iota
	WEAK
	MEDIUM   = (1000 * WEAK)
	STRONG   = (1000 * MEDIUM)
	REQUIRED = (1000 * STRONG) + STRONG + MEDIUM
)

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
		switch {
		case Weak(1) <= s && s <= Weak(1000):
			return "Weak(" + strconv.FormatFloat(float64(s), 'f', -1, 64) + ")"
		case Medium(1) <= s && s <= Medium(1000):
			return "Medium(" + strconv.FormatFloat(float64(s/1000), 'f', -1, 64) + ")"
		case Strong(1) <= s && s <= Strong(1000):
			return "Strong(" + strconv.FormatFloat(float64(s/1000000), 'f', -1, 64) + ")"
		default:
			return strconv.FormatFloat(float64(s), 'f', -1, 64)
		}
	}
}

func (s Strength) Base() Strength {
	switch {
	case Weak(1) <= s && s <= Weak(1000):
		return WEAK
	case Medium(1) <= s && s <= Medium(1000):
		return MEDIUM
	case Strong(1) <= s && s <= Strong(1000):
		return STRONG
	default:
		return s
	}
}

func (s Strength) WithWeight(weight float64) Strength {
	switch {
	case Weak(1) <= s && s <= Weak(1000):
		return Weak(weight)
	case Medium(1) <= s && s <= Medium(1000):
		return Medium(weight)
	case Strong(1) <= s && s <= Strong(1000):
		return Strong(weight)
	default:
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
