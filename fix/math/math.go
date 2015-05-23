// chris 052315

package math

import "math"

// Round returns x, rounded to the nearest integer.
//
// Special cases are:
//	Round(±Inf) = ±Inf
//	Round(NaN) = NaN
func Round(x float64) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return x
	}
	return math.Trunc(x + math.Copysign(.5, x))
}
