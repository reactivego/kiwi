package strength

import "math"

var (
	REQUIRED = Encode(1000.0, 1000.0, 1000.0, 1.0)
	STRONG   = Encode(1.0, 0.0, 0.0, 1.0)
	MEDIUM   = Encode(0.0, 1.0, 0.0, 1.0)
	WEAK     = Encode(0.0, 0.0, 1.0, 1.0)
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
