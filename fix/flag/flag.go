// chris 090415

// Package flag implements fixes or missing features from the Go
// standard library.
package flag

import (
	"fmt"
)

// Strings implements flag.Value and provides the ability to specify a
// flag that collects the string values passed in across multiple
// specifications of a command-line flag into a slice of strings.
//
// For example:
//
//	strings := new(Strings)
//	flag.Var(strings, "string", "can be specified multiple times")
//	flag.Parse()
//
//	go run whatever.go -string a -string b
//
//	// *strings now has many strings in it!
//
type Strings []string

// String returns the value of the slice of strings in a default format
// ("%v").
func (s *Strings) String() string {
	return fmt.Sprintf("%v", *s)
}

// Set adds the given string to the slice of strings.
func (s *Strings) Set(x string) error {
	*s = append(*s, x)
	return nil
}
