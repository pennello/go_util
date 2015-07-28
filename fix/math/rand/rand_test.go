// chris 0502115

package rand

import (
	"testing"

	"math/rand"
)

func TestUint64(t *testing.T) {
	// Just make sure we can call it.
	t.Log(Uint64())
}

func testChooseInt32(t *testing.T, begin, end int32) {
	x := ChooseInt32(begin, end)
	if !(begin <= x && x <= end) {
		t.Fail()
	}
}

func testChooseInt32Panic(t *testing.T, begin, end int32) {
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()
	ChooseInt32(begin, end)
}

func TestChooseInt32(t *testing.T) {
	testChooseInt32(t, 0, 0)
	testChooseInt32(t, 0, 100)
	testChooseInt32(t, 100, 100)
	testChooseInt32(t, -100, 100)
	testChooseInt32(t, -200, -100)

	for i := 0; i < 1000; i++ {
		begin := int32(rand.Uint32())
		var end int32
		for {
			end = int32(rand.Uint32())
			if begin <= end {
				break
			}
		}
		testChooseInt32(t, begin, end)
	}

	testChooseInt32Panic(t, 1, 0)
	testChooseInt32Panic(t, 1, -1)
	testChooseInt32Panic(t, -1, -2)
	testChooseInt32Panic(t, 2, 1)
}

func testChooseStringErr(t *testing.T, begin, end string) {
	_, err := ChooseString(begin, end)
	if err == nil {
		t.Fail()
	}
}

func TestChooseString(t *testing.T) {
	testChooseStringErr(t, "z", "a")
	testChooseStringErr(t, "az", "ba")
}

func TestBool(t *testing.T) {
	// Just make sure we can call it.
	t.Log(Bool())
}
