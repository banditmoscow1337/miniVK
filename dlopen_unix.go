//go:build linux || darwin

package binding

import "github.com/ebitengine/purego"

func openLibrary(name string) (uintptr, error) {
	return purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_GLOBAL)
}
func lookupLibrary(handle uintptr, name string) (uintptr, error) { return purego.Dlsym(handle, name) }
func closeLibrary(handle uintptr) error                          { return purego.Dlclose(handle) }
