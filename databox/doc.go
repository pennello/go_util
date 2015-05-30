// chris 052915

// Package databox defines an extremely simple data encapsulation format
// consisting of an int64 length header followed by arbitrary data.
//
// Given an io.Reader with arbitrary data, wrap it with a Marshaller.
// The Marshaller will first emit the header data, then the rest of the
// data from the given io.Reader.  You must either know the length up
// front of the data to be boxed, or pass in a length of -1 to instruct
// the Marshaller to read everything it can from the io.Reader into
// memory and compute the length itself based on that.
//
// Given an io.Reader with databox'd data, wrap it with an Unmarshaller.
// The Unmarshaller will first transparently read the header data, then
// emit the data indicated by the header.
package databox
