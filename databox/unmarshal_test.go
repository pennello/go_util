// chris 052915

package databox

import (
	"bytes"
	"io"
	"testing"

	"io/ioutil"
)

func testUnmarshallerNoHeader(t *testing.T, size int) {
	b := make([]byte, size)
	r := bytes.NewBuffer(b)
	u := NewUnmarshaller(r)
	_, err := ioutil.ReadAll(u)
	if err != io.ErrUnexpectedEOF {
		t.Error(err)
	}
}

func testUnmarshallerBadHeader(t *testing.T) {
	b := []byte{0xff, 0, 0, 0, 0, 0, 0, 1}
	r := bytes.NewBuffer(b)
	u := NewUnmarshaller(r)
	_, err := ioutil.ReadAll(u)
	if err != ErrBadHeader {
		t.Error(err)
	}
}

func testUnmarshaller(t *testing.T, size int) {
	if size > 0xff {
		panic("test size too big")
	}
	b := make([]byte, HeaderSize+size)
	b[7] = byte(size)
	r := bytes.NewBuffer(b)
	u := NewUnmarshaller(r)
	br, err := ioutil.ReadAll(u)
	if len(br) != size {
		t.Error(br)
	}
	if err != nil {
		t.Error(err)
	}
}

func TestUnmarshaller(t *testing.T) {
	testUnmarshallerNoHeader(t, 0)
	testUnmarshallerNoHeader(t, 1)
	testUnmarshallerNoHeader(t, 7)

	testUnmarshallerBadHeader(t)

	testUnmarshaller(t, 0)
	testUnmarshaller(t, 1)
	testUnmarshaller(t, 2)
	testUnmarshaller(t, 100)
}
