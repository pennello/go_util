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
	defer os.Remove(name)
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

func TestLockNb(t *testing.T) {
	name := testTempFileName(t)
	if name == "" {
		return
	}
	defer os.Remove(name)
	lc, err := LockNb(name)
	if err != nil {
		t.Error(err)
		return
	}
	defer lc.Unlock()

	if _, err2 := LockNb(name); err2 == nil {
		t.Error(err)
		return
	}
}

func testLockRm(t *testing.T, globalname string) string {
	localname := testTempFileName(t)
	if localname == "" {
		return ""
	}

	lrc, err := LockRm(globalname, localname)
	if err != nil {
		t.Error(err)
		return ""
	}
	defer lrc.Unlock()

	if _, err2 := LockRm(globalname, localname); err2 == nil {
		t.Error("double-lock failed to fail")
	}

	return localname
}

func TestLockRm(t *testing.T) {
	globalname := testTempFileName(t)
	if globalname == "" {
		return
	}
	defer os.Remove(globalname)
	for i := 0; i < 10000; i++ {
		localname := testLockRm(t, globalname)
		if localname == "" {
			return
		}
		if _, err := os.Stat(localname); err == nil {
			t.Error("rmlock didn't clean up")
		}
	}
}
