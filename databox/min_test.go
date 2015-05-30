// chris 052915

package databox

import (
	"bytes"
	"io"
	"testing"

	"crypto/rand"
)

func testMinReaderSize(t *testing.T, size int) {
	b := make([]byte, size)
	var n int
	var err error
	_, err = rand.Read(b)
	if err != nil {
		t.FailNow()
	}
	r := bytes.NewBuffer(b)
	lim := &io.LimitedReader{R: r, N: int64(size)}
	min := newMinReader(lim)
	p := make([]byte, size)
	n, err = min.Read(p)
	if n != size {
		t.FailNow()
	}
	if err == io.ErrUnexpectedEOF {
		t.Errorf("unexpected eof, size = %v, n = %v, err = %v", size, n, err)
	}
}

func testMinReaderEarlyEOF(t *testing.T, size int) {
	b := make([]byte, size/2)
	var n int
	var err error
	_, err = rand.Read(b)
	if err != nil {
		t.FailNow()
	}
	r := bytes.NewBuffer(b)
	lim := &io.LimitedReader{R: r, N: int64(size)}
	min := newMinReader(lim)
	p := make([]byte, size)
	n, err = min.Read(p)
	if n == size {
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
	testMinReaderSize(t, 0)
	testMinReaderSize(t, 1)
	testMinReaderSize(t, 100)
	testMinReaderEarlyEOF(t, 1)
	testMinReaderEarlyEOF(t, 100)
}
