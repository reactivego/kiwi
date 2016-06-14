package kiwi

const EPS = 1.0e-8

func NearZero(value float64) bool {
	if value < 0.0 {
		return -value < EPS
	} else {
		return value < EPS
	}
}
