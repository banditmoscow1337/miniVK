package binding

import (
	"fmt"
	"github.com/ebitengine/purego"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func libraryCandidates(variable string, names ...string) []string {
	if override := os.Getenv(variable); override != "" {
		return []string{override}
	}
	var out []string
	executable, _ := os.Executable()
	directory := filepath.Dir(executable)
	for _, name := range names {
		out = append(out, filepath.Join(directory, "lib", name), filepath.Join(directory, name))
	}
	if sdk := os.Getenv("VULKAN_SDK"); sdk != "" && variable == "DER_ATOM_VULKAN" {
		for _, name := range names {
			out = append(out, filepath.Join(sdk, "lib", name), filepath.Join(sdk, "Lib", name), filepath.Join(sdk, "Bin", name))
		}
	}
	out = append(out, names...)
	if runtime.GOOS == "darwin" {
		for _, name := range names {
			out = append(out, filepath.Join("/opt/homebrew/lib", name), filepath.Join("/usr/local/lib", name))
		}
	}
	return out
}
func loadLibrary(candidates []string) (uintptr, error) {
	var problems []string
	for _, name := range candidates {
		handle, err := openLibrary(name)
		if err == nil {
			return handle, nil
		}
		problems = append(problems, name+": "+err.Error())
	}
	return 0, fmt.Errorf("native library not found: %s", strings.Join(problems, "; "))
}
func bind(target any, address uintptr, name string) (err error) {
	if address == 0 {
		return fmt.Errorf("native function %s is unavailable", name)
	}
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("bind %s: %v", name, p)
		}
	}()
	purego.RegisterFunc(target, address)
	return nil
}
func bindLibrary(target any, handle uintptr, name string) error {
	address, err := lookupLibrary(handle, name)
	if err != nil {
		return fmt.Errorf("resolve %s: %w", name, err)
	}
	return bind(target, address, name)
}
func CString(s string) []byte { return append([]byte(s), 0) }
func String(pointer *byte) string {
	if pointer == nil {
		return ""
	}
	var b []byte
	for *pointer != 0 {
		b = append(b, *pointer)
		pointer = addByte(pointer)
	}
	return string(b)
}
