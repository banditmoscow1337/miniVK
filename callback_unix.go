//go:build linux || darwin

package binding

import "github.com/ebitengine/purego"

func floatCallback(callback func(uintptr, float64, float64)) (uintptr, func(), error) {
	pointer := purego.NewCallback(func(window uintptr, x, y float64) uintptr { callback(window, x, y); return 0 })
	return pointer, nil, nil
}
