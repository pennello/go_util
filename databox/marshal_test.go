// chris 052915

package databox

import (
	"bytes"
	"io"
	"testing"

	"io/ioutil"
)

func testMarshaller(t *testing.T, size int) {
	if size > 0xff {
		panic("size too big")
	}
	b := make([]byte, size)
	r := bytes.NewBuffer(b)
	m := NewMarshaller(r, int64(size))
	br, err := ioutil.ReadAll(m)
	expect := make([]byte, HeaderSize+size)
	copy(expect[0:HeaderSize], []byte{0, 0, 0, 0, 0, 0, 0, byte(size)})
	if !bytes.Equal(br, expect) {
		t.Error(br)
		return
	}
	if err != nil {
		t.Error(err)
	}
}

func testMarshallerBuf(t *testing.T, size int) {
	if size > 0xff {
		panic("size too big")
	}
	b := make([]byte, size)
	r := bytes.NewBuffer(b)
	m := NewMarshaller(r, -1)
	br, err := ioutil.ReadAll(m)
	expect := make([]byte, HeaderSize+size)
	copy(expect[0:HeaderSize], []byte{0, 0, 0, 0, 0, 0, 0, byte(size)})
	if !bytes.Equal(br, expect) {
		t.Error(br)
		return
	}
	if err != nil {
		t.Error(err)
	}
}

func testMarshallerShort(t *testing.T, size int) {
	if size > 0xff {
		panic("size too big")
	}
	b := make([]byte, size)
	r := bytes.NewBuffer(b)
	m := NewMarshaller(r, int64(size)+1)
	br, err := ioutil.ReadAll(m)
	expect := make([]byte, HeaderSize+size)
	copy(expect[0:HeaderSize], []byte{0, 0, 0, 0, 0, 0, 0, byte(size)})
	if bytes.Equal(br, expect) {
		t.Error(br)
		return
	}
	if err != io.ErrUnexpectedEOF {
		t.Error(err)
	}
}

func TestMarshaller(t *testing.T) {
	testMarshaller(t, 0)
	testMarshaller(t, 1)
	testMarshaller(t, 2)
	testMarshaller(t, 100)
	testMarshaller(t, 0xff)

	testMarshallerBuf(t, 0)
	testMarshallerBuf(t, 1)
	testMarshallerBuf(t, 2)
	testMarshallerBuf(t, 100)
	testMarshallerBuf(t, 0xff)

	testMarshallerShort(t, 0)
	testMarshallerShort(t, 1)
	testMarshallerShort(t, 2)
	testMarshallerShort(t, 100)
	testMarshallerShort(t, 0xff)
}
