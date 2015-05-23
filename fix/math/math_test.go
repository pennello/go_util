// chris 052315

package math

import "testing"

func testRound(t *testing.T, x float64, check int) {
	out := Round(x)
	if out != check {
		t.Errorf("Round(%v) != %v (is %v)", x, check, out)
	}
}

func TestRound(t *testing.T) {
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
