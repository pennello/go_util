// chris 052315

package math

import "math"

func Round(x float64) int {
	return int(x + math.Copysign(.5, x))
}
