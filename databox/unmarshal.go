// chris 052915

package databox

import (
	"errors"
	"io"

	"encoding/binary"
)

// ErrBadHeader is returned from Unmarshaller.Read when the header data
// is readable, but contains a negative length.
var ErrBadHeader = errors.New("bad header")

// Unmarshaller encapsulates unmarshalling a databox.  Implements
// io.Reader.
type Unmarshaller struct {
	r      io.Reader
	length int64
}

// NewUnmarshaller returns a fresh Unmarshaller, ready to unmarshal a
// databox from the given reader.
func NewUnmarshaller(r io.Reader) *Unmarshaller {
	return &Unmarshaller{r: r, length: -1}
}

// Read will transparently absorb the header and then actually start
// giving back the databox'd data.  If a header can't be read, you'll
// get back an io.ErrUnexpectedEOF.  If the header data contains a
// negative length, you'll get back an ErrBadHeader.  If the reader
// contains less data than specified by the header, then you'll get back
// an io.ErrUnexpectedEOF.
func (u *Unmarshaller) Read(p []byte) (n int, err error) {
	if u.length == -1 {
		headerBytes := make([]byte, 8)
		_, err = io.ReadFull(u.r, headerBytes)
		if err != nil {
			if err == io.EOF {
				err = io.ErrUnexpectedEOF
			}
			return 0, err
		}
		u.length = int64(binary.BigEndian.Uint64(headerBytes))
		if u.length < 0 {
			return 0, ErrBadHeader
		}
		lim := &io.LimitedReader{R: u.r, N: u.length}
		u.r = newMinReader(lim)
	}
	return u.r.Read(p)
}
