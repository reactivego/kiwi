package strength

import "math"

var (
	REQUIRED = Encode(1000.0, 1000.0, 1000.0, 1.0)
	STRONG   = Encode(1.0, 0.0, 0.0, 1.0)
	MEDIUM   = Encode(0.0, 1.0, 0.0, 1.0)
	WEAK     = Encode(0.0, 0.0, 1.0, 1.0)
	OPTIONAL = Encode(0.0, 0.0, 0.0, 1.0)
)

type Value float64

func (strength Value) Clip() Value {
	return Value(math.Max(0.0, math.Min(float64(strength), float64(REQUIRED))))
}

func (strength Value) Float64() float64 {
	return float64(strength)
}

func Encode(a, b, c, w float64) Value {
	var strength float64
	strength += math.Max(0.0, math.Min(1000.0, a*w)) * 1000000.0
	strength += math.Max(0.0, math.Min(1000.0, b*w)) * 1000.0
	strength += math.Max(0.0, math.Min(1000.0, c*w))
	return Value(strength)
}

// Strength Value: Optional() < Weak(1 .. <1000) < Medium(1 .. <1000) < Strong(1 .. <1000) < Required()

func Required() Value {
	return Value(1000000000.0)
}

// weight 1..<1000
func Strong(weight float64) Value {
	// 1000000 .. 999999999
	return encode(1000000.0, weight)
}

// weight 1..<1000
func Medium(weight float64) Value {
	// 1000 .. 999999.999
	return encode(1000.0, weight)
}

// weight 1..<1000
func Weak(weight float64) Value {
	// 1 .. 999.9999999
	return encode(1.0, weight)
}

func Optional() Value {
	// 0.0
	return Value(0.0)
}

func encode(base, weight float64) Value {
	return Value(base * math.Max(1.0, math.Min(weight, 999.999999)))
}
