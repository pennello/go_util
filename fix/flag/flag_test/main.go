// chris 090415

// flag_test tests the contents of its parent package.  This is
// implemented as a separate command due to the flag library using
// global state and it being somewhat difficult to integrate with
// package testing.
package main

import (
	"flag"
	"log"

	fixflag "chrispennello.com/go/util/fix/flag"
)

func main() {
	tests := new(fixflag.SliceStrings)
	flag.Var(tests, "test", "can specify more than once")
	flag.Parse()
	log.Print(*tests)
}
