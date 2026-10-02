//go:build windows

package binding

import "syscall"

func openLibrary(name string) (uintptr, error) {
	h, err := syscall.LoadLibrary(name)
	return uintptr(h), err
}
func lookupLibrary(handle uintptr, name string) (uintptr, error) {
	return syscall.GetProcAddress(syscall.Handle(handle), name)
}
func closeLibrary(handle uintptr) error { return syscall.FreeLibrary(syscall.Handle(handle)) }
