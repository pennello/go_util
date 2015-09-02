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

	lrc, err2 := LockRm(gf.Name(), lf.Name())
	if err2 != nil {
		t.Error(err2)
		return
	}
	defer func() {
		// Extra strict unlock error detection since this is a
		// test.
		err := lrc.Unlock()
		if err != nil {
			t.Error(err)
		}
	}()

	_, err3 := LockRm(gf.Name(), lf.Name())
	if err3 == nil {
		lrc.Unlock()
		t.Error(err3)
	}
}
