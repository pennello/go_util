// chris 090415

// Package flag implements fixes or missing features from the Go
// standard library.
package flag

import (
	"fmt"
)

// SliceStrings implements flag.Value and provides the ability to
// specify a flag that collects the string values passed in across
// multiple specifications of a command-line flag into a slice of
// strings.
//
// For example:
//
//	strings := new(SliceStrings)
//	flag.Var(strings, "string", "can be specified multiple times")
//	flag.Parse()
//
//	go run whatever.go -string a -string b
//
//	// *strings now has many strings in it!
//
type SliceStrings []string

// String returns the value of the slice of strings in a default format
// ("%v").
func (ss *SliceStrings) String() string {
	return fmt.Sprintf("%v", *ss)
}

// Set adds the given string to the slice of strings.
func (ss *SliceStrings) Set(x string) error {
	*ss = append(*ss, x)
	return nil
}
