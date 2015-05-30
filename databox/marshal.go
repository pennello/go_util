// chris 052815

package databox

import (
	"bytes"
	"io"

	"encoding/binary"
	"io/ioutil"
)

// Marshaller encapsulates marshalling an arbitrary io.Reader into a
// databox.  Implements io.Reader.
type Marshaller struct {
	r io.Reader
	// length is an int64 so we can use io.LimitedReader on
	// unmarshal.  It's also so we can use -1 as a convention for an
	// unknown length.
	length      int64
	headerReady bool
}

// NewMarshaller returns a fresh Marshaller, ready to marshal the given
// reader with its expected length into a databox.  Pass length = -1 if
// the length is unknown, and the Marshaller will read the entirety of
// the reader's data into an in-memory buffer and compute the length
// itself based on that.
//
// Panics if you pass a length less than -1.
func NewMarshaller(r io.Reader, length int64) *Marshaller {
	if length < -1 {
		panic("unexpected negative length")
	}
	return &Marshaller{r: r, length: length, headerReady: false}
}

// Read will first give you the header, then the data from the reader.
// If there is less data in the reader than the length indicated, then
// an io.ErrUnexpectedEOF will be returned.
func (m *Marshaller) Read(p []byte) (n int, err error) {
	if m.length == -1 {
		// The user indicated the length is unknown.  We
		// therefore buffer the whole thing into memory to
		// figure it out.
		data, err := ioutil.ReadAll(m.r)
		if err != nil {
			return 0, err
		}
		m.r = bytes.NewBuffer(data)
		m.length = int64(len(data))
	}

	if !m.headerReady {
		headerBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(headerBytes, uint64(m.length))
		header := bytes.NewBuffer(headerBytes)
		lim := &io.LimitedReader{R: m.r, N: m.length}
		min := newMinReader(lim)
		m.r = io.MultiReader(header, min)
		m.headerReady = true
	}

	return m.r.Read(p)
}
