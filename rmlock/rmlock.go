// chris 090115 Unix removable lock file.

// TODO Note how Close calls errors are not handled.

package rmlock

import (
	"os"

	"golang.org/x/sys/unix"
)

const mode = 0666

type LockContext struct {
	globalname string

	local *os.File
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

func Lock(globalname, localname string) (*LockContext, error) {
	var lc *LockContext

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

		lc = &LockContext{
			globalname: globalname,
			local: f,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return lc, nil

}

func (lc *LockContext) Unlock() error {
	return globalCtx(lc.globalname, func() error {
		lc.local.Close()
		return os.Remove(lc.local.Name())
	})
}
