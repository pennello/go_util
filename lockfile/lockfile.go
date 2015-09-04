// chris 090115

// Package lockfile implements convenient lock file utilities for
// Unix-based systems.
//
// Removable Lock Files - The Difficulty
//
// Removable lock files are notoriously difficult to get right.  On BSD
// systems, doing this in a race-free manner is trivial, via O_EXLOCK in
// an open(2) call, but this is generally unavailable on Unix-like
// systems due to Linux's lack of support for this option.
//
// Consider several processes' hypothetical sequences of open, lock,
// close (thus implicitly removing the kernel advisory lock), and
// unlink.
//
//	A       B       C
//	open
//	        open
//	lock
//	close
//	        lock
//	unlink
//	                open
//	                lock
//
// Now B thinks it's got a lock, but its lock file has been removed.  C
// has opened the same lock file name, thus creating a new file, and
// locked it.  So now B and C both think they've got the lock.  Game
// over.
//
// You might attempt re-arranging the close and unlink calls, but the
// problem remains.  In general, if B opens the same file as A, and B
// locks after A closes, then B can have a lock on a dangling file
// descriptor.
//
// The general problem is that the close and unlink steps are not
// atomic.
//
// Removable Lock Files - A Solution
//
// One solution is to guard the two halves of the removable lock file
// operations, open and lock, and close and unlink, with another lock
// file that is itself not removed.  This is the approach that this
// package takes with LockRm.  Using this approach, removable lock files
// may be implemented in a race-free manner.
package lockfile

// NB: Does not compile on Windows.

// Close calls' errors are not handled explicitly.  The error conditions
// all end up with the file descriptor being closed anyway, so there is
// nothing special to handle.

import (
	"os"

	"golang.org/x/sys/unix"
)

// LockContext represents a locked file, and is obtained by calling Lock
// or LockNb.
type LockContext struct {
	f *os.File
}

// LockRmContext represents a locked file that can be removed on Unlock,
// and is obtained by calling LockRm.
type LockRmContext struct {
	globalname string

	local *LockContext
}

// lock is the internal implementation for Lock and LockNb.  You merely
// specify whether or not you want the flock call to block by passing
// the block boolean.
func lock(filename string, block bool) (*LockContext, error) {
	f, err := os.OpenFile(filename, os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	how := unix.LOCK_EX
	if !block {
		how = how | unix.LOCK_NB
	}
	if err := unix.Flock(int(f.Fd()), how); err != nil {
		f.Close()
		return nil, err
	}

	return &LockContext{f}, nil
}

// Lock locks on the filename given and returns a new LockContext, on
// which you can later call Unlock.  This implementation blocks, and
// does not clean up the lock file on Unlock.
func Lock(filename string) (*LockContext, error) {
	return lock(filename, true)
}

// LockNb locks on the filename given and returns a new LockContext, on
// which you can later call Unlock.  This implementation does not block,
// and does not clean up the lock file on Unlock.
func LockNb(filename string) (*LockContext, error) {
	return lock(filename, false)
}

// Unlock unlocks the lock file represented by the LockContext.
func (lc *LockContext) Unlock() {
	// Close implicitly releases any kernel advisory locks.
	lc.f.Close()
}

// globalCtx wraps an inner function with a blocking Lock on a global
// lock file.  This is race-free since the global lock file is not
// removed.
func globalCtx(globalname string, inner func() error) error {
	glc, err := Lock(globalname)
	if err != nil {
		return err
	}
	defer glc.Unlock()
	return inner()
}

// LockRm implements a removable lock file, specified by localname.
// This implementation does not block, and removes the lock file on
// Unlock.
//
// On BSD systems, doing this in a race-free manner is trivial, via
// O_EXLOCK in an open(2) call, but this is generally unavailable on
// Unix-like systems due to Linux's lack of support for this option.
//
// With the normal facilities provided, removing a lock file on unlock
// creates race conditions.  However, if the "local" lock file
// operations are secured by use of a "global" lock file, which is
// itself not removed, this can be implemented in a race-free manner.
func LockRm(globalname, localname string) (*LockRmContext, error) {
	var lrc *LockRmContext
	err := globalCtx(globalname, func() error {
		llc, err := LockNb(localname)
		if err != nil {
			return err
		}
		lrc = &LockRmContext{
			globalname: globalname,
			local:      llc,
		}
		return nil
	})
	return lrc, err
}

// Unlock unlocks and removes the lock file represented by the
// LockRmContext.
func (lrc *LockRmContext) Unlock() error {
	return globalCtx(lrc.globalname, func() error {
		lrc.local.Unlock()
		return os.Remove(lrc.local.f.Name())
	})
}
