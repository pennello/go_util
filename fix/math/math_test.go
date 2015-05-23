// chris 052315

package math

import (
	"math"
	"testing"
)

func testRound(t *testing.T, x, check float64) {
	out := Round(x)
	if out != check {
		t.Errorf("Round(%v) != %v (is %v)", x, check, out)
	}
}

func testRoundInf(t *testing.T, sign int) {
	x := math.Inf(sign)
	out := Round(x)
	if !math.IsInf(Round(math.Inf(sign)), sign) {
		t.Errorf("Round(%v) \"!=\" %v (is %v)", x, x, out)
	}
}

func testRoundNaN(t *testing.T) {
	x := math.NaN()
	out := Round(x)
	if !math.IsNaN(out) {
		t.Errorf("Round(%v) \"!=\" %v (is %v)", x, x, out)
	}
}

func TestRound(t *testing.T) {
	testRoundInf(t, 1)
	testRoundInf(t, -1)

	testRoundNaN(t)

	testRound(t, -0, -0)
	testRound(t, 0, -0)
	testRound(t, -0, 0)
	testRound(t, 0, 0)

	testRound(t, .25, 0)
	testRound(t, .5, 1)
	testRound(t, .75, 1)
	testRound(t, 1, 1)

	testRound(t, -.25, 0)
	testRound(t, -.5, -1)
	testRound(t, -.75, -1)
	testRound(t, -1, -1)

	testRound(t, 27, 27)
	testRound(t, 27.25, 27)
	testRound(t, 27.5, 28)
	testRound(t, 27.75, 28)
	testRound(t, 28, 28)

	testRound(t, -27, -27)
	testRound(t, -27.25, -27)
	testRound(t, -27.5, -28)
	testRound(t, -27.75, -28)
	testRound(t, -28, -28)
}
