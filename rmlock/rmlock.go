// chris 090115 Unix removable lock file.

// TODO Note how Close calls errors are not handled.

package rmlock

import (
	"os"

	"golang.org/x/sys/unix"
)

type LockContext struct {
	globalname string

	local *os.File
}

func globalCtx(globalname string, inner func() error) error {
	f, err := os.OpenFile(globalname, os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := unix.Flock(f.Fd(), unix.O_LOCK_EX); err != nil {
		return err
	}

	if err := inner(); err != nil {
		return err
	}

	if err = unix.Flock(f.Fd(), unix.O_LOCK_UN); err != nil {
		return err
	}
}

func Lock(globalname, localname string) (*LockContext, error) {
	var lc LockContext

	err := globalCtx(globalname, func() error {
		f, err := os.Openfile(localname, os.O_CREATE, 0666)
		if err != nil {
			return err
		}

		err = unix.Flock(f.Fd(), unix.O_LOCKEX | unix.O_NONBLOCK)
		if err != nil {
			f.Close()
			return err
		}

		lc = &LockContext{
			globalname: globalname,
			local: f,
		}
	})

	if err != nil {
		return nil, err
	}

	return lc, nil

}

func (lc *LockContext) Unlock() error {
	return globalCtx(lc.globalname, func() error {
		lc.local.Close()
		return os.Remove(lc.local.Name)
	})
}
