// chris 0502115

package rand

// Uint64 returns a pseudo-random 32-bit value as a uint64 from the
// default Source.  Calls Uint32, shifts, and ors with another call to
// Uint32.
// https://groups.google.com/d/topic/golang-nuts/Kle874lT1Eo/discussion
func Uint64() uint64 {
  return uint64(rand.Uint32())<<32 | uint64(rand.Uint32())
}
