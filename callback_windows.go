//go:build windows

package binding

import (
	"encoding/binary"
	"fmt"
	"math"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/ebitengine/purego"
)

func floatCallback(callback func(uintptr, float64, float64)) (uintptr, func(), error) {
	target := purego.NewCallback(func(window, x, y uintptr) uintptr {
		callback(window, math.Float64frombits(uint64(x)), math.Float64frombits(uint64(y)))
		return 0
	})
	var code []byte
	switch runtime.GOARCH {
	case "amd64":
		// movq rdx,xmm1; movq r8,xmm2; movabs rax,target; jmp rax
		code = []byte{0x66, 0x48, 0x0f, 0x7e, 0xca, 0x66, 0x49, 0x0f, 0x7e, 0xd0, 0x48, 0xb8}
		address := make([]byte, 8)
		binary.LittleEndian.PutUint64(address, uint64(target))
		code = append(code, address...)
		code = append(code, 0xff, 0xe0)
	case "arm64":
		// fmov x1,d0; fmov x2,d1; ldr x16,8; br x16; .quad target
		code = make([]byte, 24)
		for n, word := range []uint32{0x9e660001, 0x9e660022, 0x58000050, 0xd61f0200} {
			binary.LittleEndian.PutUint32(code[n*4:], word)
		}
		binary.LittleEndian.PutUint64(code[16:], uint64(target))
	default:
		return 0, nil, fmt.Errorf("unsupported Windows callback architecture %s", runtime.GOARCH)
	}
	kernel := syscall.NewLazyDLL("kernel32.dll")
	allocate := kernel.NewProc("VirtualAlloc")
	protect := kernel.NewProc("VirtualProtect")
	free := kernel.NewProc("VirtualFree")
	flush := kernel.NewProc("FlushInstructionCache")
	address, _, allocationErr := allocate.Call(0, uintptr(len(code)), 0x1000|0x2000, 0x04)
	if address == 0 {
		return 0, nil, fmt.Errorf("allocate callback trampoline: %w", allocationErr)
	}
	release := func() { _, _, _ = free.Call(address, 0, 0x8000) }
	copy(unsafe.Slice((*byte)(unsafe.Pointer(address)), len(code)), code)
	var old uint32
	ok, _, protectErr := protect.Call(address, uintptr(len(code)), 0x20, uintptr(unsafe.Pointer(&old)))
	if ok == 0 {
		release()
		return 0, nil, fmt.Errorf("protect callback trampoline: %w", protectErr)
	}
	ok, _, flushErr := flush.Call(^uintptr(0), address, uintptr(len(code)))
	if ok == 0 {
		release()
		return 0, nil, fmt.Errorf("flush callback instructions: %w", flushErr)
	}
	runtime.KeepAlive(code)
	return address, release, nil
}
