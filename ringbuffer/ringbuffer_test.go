// chris 00515

package ringbuffer

import (
	"bytes"
	"testing"

	cryptorand "crypto/rand"
	mathrand "math/rand"
)

func testBytesHelp(t *testing.T, b *B, expect []byte) {
	actual := b.Bytes()
	if !bytes.Equal(expect, actual) {
		t.Logf("pos: %v\n", b.pos)
		t.Logf("buf: %v\n", b.buf)
		t.Errorf("incorrect bytes: %v (expected %v)\n", actual, expect)
		return
	}
}

func TestBytes(t *testing.T) {
	var b *B

	b = &B{
		Size: 3,
		buf:  []byte(`abc`),
		pos:  0,

		wrapped: true,
	}

	testBytesHelp(t, b, []byte(`abc`))
	b.pos = 1
	testBytesHelp(t, b, []byte(`bca`))
	b.pos = 2
	testBytesHelp(t, b, []byte(`cab`))

	b = &B{
		Size: 3,
		buf:  []byte(`abc`),
		pos:  0,

		wrapped: false,
	}

	testBytesHelp(t, b, []byte{})
	b.pos = 1
	testBytesHelp(t, b, []byte(`a`))
	b.pos = 2
	testBytesHelp(t, b, []byte(`ab`))
	b.pos = 3
	testBytesHelp(t, b, []byte(`abc`))
}

func testWriteHelp(t *testing.T, b *B, p []byte) (success bool) {
	n, err := b.Write(p)
	if err != nil {
		t.Error("write failed", err)
		return false
	}
	if n != len(p) {
		t.Errorf("didn't write whole slice (only %d)\n", n)
		return false
	}
	return true
}

func TestWriteSmall1(t *testing.T) {
	b := New(5)
	b.buf = []byte(`abcde`)
	if success := testWriteHelp(t, b, []byte(`xy`)); !success {
		return
	}
	expect := []byte(`xy`)
	testBytesHelp(t, b, expect)
}

func TestWriteSmall2(t *testing.T) {
	b := New(5)
	b.buf = []byte(`abcde`)
	if success := testWriteHelp(t, b, []byte(`xyz`)); !success {
		return
	}
	expect := []byte(`xyz`)
	testBytesHelp(t, b, expect)
}

func TestWriteBig1(t *testing.T) {
	b := New(5)
	b.buf = []byte(`abcde`)
	if success := testWriteHelp(t, b, []byte(`xyzwvut`)); !success {
		return
	}
	expect := []byte(`zwvut`)
	testBytesHelp(t, b, expect)
}

func TestWriteBig2(t *testing.T) {
	b := New(5)
	b.buf = []byte(`abcde`)
	if success := testWriteHelp(t, b, []byte(`xyzwvutp`)); !success {
		return
	}
	expect := []byte(`wvutp`)
	testBytesHelp(t, b, expect)
}

func TestWriteExact(t *testing.T) {
	b := New(5)
	b.buf = []byte(`abcde`)
	p := []byte(`xyzwv`)
	if success := testWriteHelp(t, b, p); !success {
		return
	}
	expect := p
	testBytesHelp(t, b, expect)
}

func testWriteRandom(t *testing.T) {
	const max = 16384
	posintn := func(n int) int {
		return mathrand.Intn(max-1) + 1
	}
	b := New(posintn(max))
	p := make([]byte, posintn(max))
	if _, err := cryptorand.Read(p); err != nil {
		t.Error(err)
		return
	}

	if success := testWriteHelp(t, b, p); !success {
		return
	}

	var expect []byte
	if len(p) > b.Size {
		expect = p[len(p)-b.Size:]
	} else {
		expect = p
	}
	testBytesHelp(t, b, expect)
}

func TestWriteRandom(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping random tests in short mode")
	}
	for i := 0; i < 1000; i++ {
		testWriteRandom(t)
	}
}

func TestMultiWrite1(t *testing.T) {
	b := New(5)
	b.buf = []byte(`abcde`)
	if success := testWriteHelp(t, b, []byte(`x`)); !success {
		return
	}
	if success := testWriteHelp(t, b, []byte(`y`)); !success {
		return
	}
	expect := []byte(`xy`)
	testBytesHelp(t, b, expect)
}

func TestMultiWrite2(t *testing.T) {
	b := New(5)
	b.buf = []byte(`abcde`)
	if success := testWriteHelp(t, b, []byte(`x`)); !success {
		return
	}
	if success := testWriteHelp(t, b, []byte(`y`)); !success {
		return
	}
	if success := testWriteHelp(t, b, []byte(`z`)); !success {
		return
	}
	expect := []byte(`xyz`)
	testBytesHelp(t, b, expect)
}

func TestMultiWrite3(t *testing.T) {
	b := New(5)
	b.buf = []byte(`abcde`)
	if success := testWriteHelp(t, b, []byte(`xy`)); !success {
		return
	}
	if success := testWriteHelp(t, b, []byte(`z`)); !success {
		return
	}
	expect := []byte(`xyz`)
	testBytesHelp(t, b, expect)
}

func TestMultiWrite4(t *testing.T) {
	b := New(5)
	b.buf = []byte(`abcde`)
	if success := testWriteHelp(t, b, []byte(`x`)); !success {
		return
	}
	if success := testWriteHelp(t, b, []byte(`yz`)); !success {
		return
	}
	expect := []byte(`xyz`)
	testBytesHelp(t, b, expect)
}

func TestMultiWrite5(t *testing.T) {
	b := New(5)
	b.buf = []byte(`abcde`)
	if success := testWriteHelp(t, b, []byte(`x`)); !success {
		return
	}
	if success := testWriteHelp(t, b, []byte(`y`)); !success {
		return
	}
	if success := testWriteHelp(t, b, []byte(`12345`)); !success {
		return
	}
	expect := []byte(`12345`)
	testBytesHelp(t, b, expect)
}

func TestMultiWrite6(t *testing.T) {
	b := New(5)
	b.buf = []byte(`abcde`)
	if success := testWriteHelp(t, b, []byte(`x`)); !success {
		return
	}
	if success := testWriteHelp(t, b, []byte(`y`)); !success {
		return
	}
	if success := testWriteHelp(t, b, []byte(`z`)); !success {
		return
	}
	if success := testWriteHelp(t, b, []byte(`12345`)); !success {
		return
	}
	expect := []byte(`12345`)
	testBytesHelp(t, b, expect)
}

func TestMultiWrite7(t *testing.T) {
	b := New(5)
	b.buf = []byte(`abcde`)
	if success := testWriteHelp(t, b, []byte(`xy`)); !success {
		return
	}
	if success := testWriteHelp(t, b, []byte(`z`)); !success {
		return
	}
	if success := testWriteHelp(t, b, []byte(`12345`)); !success {
		return
	}
	expect := []byte(`12345`)
	testBytesHelp(t, b, expect)
}

func TestMultiWrite8(t *testing.T) {
	b := New(5)
	b.buf = []byte(`abcde`)
	if success := testWriteHelp(t, b, []byte(`x`)); !success {
		return
	}
	if success := testWriteHelp(t, b, []byte(`12345`)); !success {
		return
	}
	expect := []byte(`12345`)
	testBytesHelp(t, b, expect)
}
