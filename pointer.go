package binding

import "unsafe"

func addByte(p *byte) *byte { return (*byte)(unsafe.Add(unsafe.Pointer(p), 1)) }
