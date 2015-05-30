// chris 052915

package databox

import (
	"bytes"
	"testing"

	cryptorand "crypto/rand"
	"io/ioutil"
	mathrand "math/rand"
)

func testRwEmpty(t *testing.T) {
	m := NewMarshaller(bytes.NewBuffer([]byte{}), 0)
	var bm, bu []byte
	var err error
	bm, err = ioutil.ReadAll(m)
	if len(bm) != 8 {
		t.Error(bm)
	}
	if err != nil {
		t.Error(err)
		return
	}
	u := NewUnmarshaller(bytes.NewBuffer(bm))
	bu, err = ioutil.ReadAll(u)
	if !bytes.Equal(bm[8:], bu) {
		t.Fail()
	}
}

func testRwRandom(t *testing.T) {
	size := mathrand.Intn(100)
	m := NewMarshaller(cryptorand.Reader, int64(size))
	var bm, bu []byte
	var err error
	bm, err = ioutil.ReadAll(m)
	if len(bm) != 8+size {
		t.Error(bm)
	}
	if err != nil {
		t.Error(err)
		return
	}
	u := NewUnmarshaller(bytes.NewBuffer(bm))
	bu, err = ioutil.ReadAll(u)
	if !bytes.Equal(bm[8:], bu) {
		t.Fail()
	}
}

func TestRw(t *testing.T) {
	testRwEmpty(t)

	for i := 0; i < 1000; i++ {
		testRwRandom(t)
	}
}
