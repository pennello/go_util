// chris 090115

package lockfile

import (
	"os"
	"testing"

	"io/ioutil"
)

func TestRmlock(t *testing.T) {
	gf, err := ioutil.TempFile(".", "lockfile_test_global_")
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(gf.Name())
	lf, err := ioutil.TempFile(".", "lockfile_test_local_")

	lc, err2 := Lock(gf.Name(), lf.Name())
	if err2 != nil {
		t.Error(err2)
		return
	}
	defer func() {
		// Extra strict unlock error detection since this is a
		// test.
		err := lc.Unlock()
		if err != nil {
			t.Error(err)
		}
	}()

	_, err3 := Lock(gf.Name(), lf.Name())
	if err3 == nil {
		lc.Unlock()
		t.Error(err3)
	}
}
