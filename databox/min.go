// chris 052915

package databox

import "io"

// A minReader enforces a minimum size to read based on a LimitedReader:
// if we get an EOF too early, it's transformed into an
// ErrUnexpectedEOF.

type minReader struct {
	lim *io.LimitedReader
}

func newMinReader(lim *io.LimitedReader) *minReader {
	return &minReader{lim: lim}
}

func (mr *minReader) Read(p []byte) (n int, err error) {
	n, err = mr.lim.Read(p)
	if err == io.EOF && mr.lim.N > 0 {
		err = io.ErrUnexpectedEOF
	}
	return n, err
}
