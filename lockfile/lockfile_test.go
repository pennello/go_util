// chris 090115

package lockfile

import (
	"os"
	"sync"
	"testing"
	"time"

	"io/ioutil"
)

func testTempFileName(t *testing.T) string {
	f, err := ioutil.TempFile(".", "lockfile_test_")
	if err != nil {
		t.Error(err)
		return ""
	}
	f.Close()
	return f.Name()
}

func TestLock(t *testing.T) {
	name := testTempFileName(t)
	if name == "" {
		return
	}
	lc, err := Lock(name)
	if err != nil {
		t.Error(err)
		return
	}

	locked := true
	sub := false
	done := make(chan struct{})
	mu := new(sync.Mutex)

	go func() {
		defer close(done)
		_, err := Lock(name)
		if err != nil {
			t.Error(err)
			return
		}
		mu.Lock()
		defer mu.Unlock()
		if locked {
			t.Error("parent-goroutine still locked")
			return
		}
		sub = true
	}()

	time.Sleep(10 * time.Millisecond)
	mu.Lock()
	lc.Unlock()
	locked = false
	if sub {
		t.Error("sub-goroutine was able to acquire lock")
	}
	mu.Unlock()

	<-done
}

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
		t.Error("double-lock failed to fail")
	}
}
