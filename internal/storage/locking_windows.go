//go:build windows

package storage

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	procLockFileEx   = kernel32.NewProc("LockFileEx")
	procUnlockFileEx = kernel32.NewProc("UnlockFileEx")
)

func platformLock(f *os.File) error {
	h := f.Fd()
	var ol syscall.Overlapped
	r, _, err := procLockFileEx.Call(
		h,
		2, // LOCKFILE_EXCLUSIVE_LOCK
		0,
		1,
		0,
		uintptr(unsafe.Pointer(&ol)),
	)
	if r == 0 {
		return err
	}
	return nil
}

func platformUnlock(f *os.File) error {
	h := f.Fd()
	var ol syscall.Overlapped
	r, _, err := procUnlockFileEx.Call(
		h,
		0,
		1,
		0,
		uintptr(unsafe.Pointer(&ol)),
	)
	if r == 0 {
		return err
	}
	return nil
}
