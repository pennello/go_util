// chris 052115 Uint64.
// chris 072315 Choose functions.
// chris 072815 Bool.

package rand

import (
	"errors"
	"fmt"

	"math/rand"
)

// ErrBadRange is returned by ChooseString when the input strings do not
// contain the same number of runes.
var ErrBadRange = errors.New("bad range")

// Uint64 returns a pseudo-random 64-bit value as a uint64 from the
// default Source.
//
// Calls Uint32, shifts, and ors with another call to Uint32.
//
// See:
// https://groups.google.com/d/topic/golang-nuts/Kle874lT1Eo/discussion
func Uint64() uint64 {
	return uint64(rand.Uint32())<<32 | uint64(rand.Uint32())
}

// ChooseInt32 returns, as an int32, a pseudo-random integer in
// [begin, end] from the default Source.  It panics if begin > end.
//
// This is bit width-specific because it uses a 64-bit integer to
// contain the difference between begin and end.
func ChooseInt32(begin, end int32) int32 {
	if begin > end {
		panic(fmt.Sprintf("invalid arguments to ChooseInt32: begin %d > end %d", begin, end))
	}
	diff := int64(end) - int64(begin)
	return int32(int64(begin) + rand.Int63n(diff+1))
}

// ChooseRune returns a pseudo-random rune that is lexically and
// inclusively between begin and end.  Specifically, the output x
// satisfies begin <= x <= end.
func ChooseRune(begin, end rune) rune {
	return rune(ChooseInt32(int32(begin), int32(end)))
}

// ChooseString returns a pseudo-random string that is lexically and
// inclusively between begin and end, using the default Source.
// Specifically, the output x satisfies begin <= x <= end.
//
// If the begin and end strings don't contain the same number of runes,
// the empty string is returned with ErrBadRange.  If they do, but at a
// particular index, end contains a rune less than the corresponding
// rune in begin, the empty string is returned with ErrBadRange.
func ChooseString(begin, end string) (string, error) {
	br := []rune(begin)
	er := []rune(end)
	if len(br) != len(er) {
		return "", ErrBadRange
	}
	ret := make([]rune, len(br))
	for i := range ret {
		if int32(br[i]) > int32(er[i]) {
			return "", ErrBadRange
		}
		ret[i] = ChooseRune(br[i], er[i])
	}
	return string(ret), nil
}

// Bool returns a pseudo-random boolean from the default Source.
func Bool() bool {
	// The distribution of the whole is uniform, so the distribution
	// of each bit is uniform.  Thus, the probability that a
	// particular bit is 1 is 50%.

	// We use Int63 since that is the most fundamental source of
	// random bits in the math/rand package.

	if rand.Int63()&1 == 1 {
		return true
	}
	return false
}
