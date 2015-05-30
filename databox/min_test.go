// chris 052915

package databox

import (
	"bytes"
	"io"
	"testing"

	"crypto/rand"
)

func testMinReaderLength(t *testing.T, length int) {
	b := make([]byte, length)
	var n int
	var err error
	_, err = rand.Read(b)
	if err != nil {
		t.FailNow()
	}
	r := bytes.NewBuffer(b)
	lim := &io.LimitedReader{R: r, N: int64(length)}
	min := newMinReader(lim)
	p := make([]byte, length)
	n, err = min.Read(p)
	if n != length {
		t.FailNow()
	}
	if err == io.ErrUnexpectedEOF {
		t.Errorf("unexpected eof, length = %v, n = %v, err = %v", length, n, err)
	}
}

func testMinReaderEarlyEOF(t *testing.T, length int) {
	b := make([]byte, length/2)
	var n int
	var err error
	_, err = rand.Read(b)
	if err != nil {
		t.FailNow()
	}
	r := bytes.NewBuffer(b)
	lim := &io.LimitedReader{R: r, N: int64(length)}
	min := newMinReader(lim)
	p := make([]byte, length)
	n, err = min.Read(p)
	if n == length {
		t.FailNow()
	}
	n, err = min.Read(p)
	if n != 0 {
		t.FailNow()
	}
	if err != io.ErrUnexpectedEOF {
		t.FailNow()
	}
}

func TestMinReader(t *testing.T) {
	testMinReaderLength(t, 0)
	testMinReaderLength(t, 1)
	testMinReaderLength(t, 100)
	testMinReaderEarlyEOF(t, 1)
	testMinReaderEarlyEOF(t, 100)
}
