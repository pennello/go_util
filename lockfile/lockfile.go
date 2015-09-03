// chris 090115 Unix removable lock file.

// TODO Note how Close calls errors are not handled.
// TODO Generalize to lockfile library:
//       - Lock
//       - LockNb
//       - LockRm

package lockfile

import (
	"os"

	"golang.org/x/sys/unix"
)

const mode = 0666

type LockContext struct {
	f *os.File
}

type LockRmContext struct {
	globalname string

	local *os.File
}

// TODO document:
// - blocking
// - doesn't remove
func Lock(filename string) (*LockContext, error) {
	f, err := os.OpenFile(filename, os.O_CREATE, mode)
	if err != nil {
		return nil, err
	}

	if err := unix.Flock(int(f.Fd()), unix.LOCK_EX); err != nil {
		f.Close()
		return nil, err
	}

	return &LockContext{f}, nil
}

func (lc *LockContext) Unlock() {
	// Close implicitly releases any kernel advisory locks.
	lc.f.Close()
}

func globalCtx(globalname string, inner func() error) error {
	f, err := os.OpenFile(globalname, os.O_CREATE, mode)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := unix.Flock(int(f.Fd()), unix.LOCK_EX); err != nil {
		return err
	}

	if err := inner(); err != nil {
		return err
	}

	if err = unix.Flock(int(f.Fd()), unix.LOCK_UN); err != nil {
		return err
	}

	return nil
}

func LockRm(globalname, localname string) (*LockRmContext, error) {
	var lrc *LockRmContext

	err := globalCtx(globalname, func() error {
		f, err := os.OpenFile(localname, os.O_CREATE, mode)
		if err != nil {
			return err
		}

		err = unix.Flock(int(f.Fd()), unix.LOCK_EX | unix.LOCK_NB)
		if err != nil {
			f.Close()
			return err
		}

		lrc = &LockRmContext{
			globalname: globalname,
			local: f,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return lrc, nil

}

func (lrc *LockRmContext) Unlock() error {
	return globalCtx(lrc.globalname, func() error {
		lrc.local.Close()
		return os.Remove(lrc.local.Name())
	})
}
