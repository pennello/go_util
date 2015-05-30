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
	// size is an int64 so we can use io.LimitedReader on
	// unmarshal.  It's also so we can use -1 as a convention for an
	// unknown size.
	size        int64
	headerReady bool
}

// NewMarshaller returns a fresh Marshaller, ready to marshal the given
// reader with its expected size into a databox.  Pass size = -1 if the
// size is unknown, and the Marshaller will read the entirety of the
// reader's data into an in-memory buffer and compute the size itself
// based on that.
//
// Panics if you pass a size less than -1.
func NewMarshaller(r io.Reader, size int64) *Marshaller {
	if size < -1 {
		panic("unexpected negative size")
	}
	return &Marshaller{r: r, size: size, headerReady: false}
}

// Read will first give you the header, then the data from the reader.
// If there is less data in the reader than the size indicated, then an
// io.ErrUnexpectedEOF will be returned.
func (m *Marshaller) Read(p []byte) (n int, err error) {
	if m.size == -1 {
		// The user indicated the size is unknown.  We therefore
		// buffer the whole thing into memory to figure it out.
		data, err := ioutil.ReadAll(m.r)
		if err != nil {
			return 0, err
		}
		m.r = bytes.NewBuffer(data)
		m.size = int64(len(data))
	}

	if !m.headerReady {
		headerBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(headerBytes, uint64(m.size))
		header := bytes.NewBuffer(headerBytes)
		lim := &io.LimitedReader{R: m.r, N: m.size}
		min := newMinReader(lim)
		m.r = io.MultiReader(header, min)
		m.headerReady = true
	}

	return m.r.Read(p)
}
