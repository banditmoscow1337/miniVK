//go:build linux

package binding

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

type linuxBackendType int

const (
	linuxBackendAuto linuxBackendType = iota
	linuxBackendX11
	linuxBackendWayland
	linuxBackendHeadless
)

var (
	requestedLinuxBackend = linuxBackendAuto
	activeLinuxBackend    = linuxBackendAuto
)

// SetLinuxBackend explicitly selects the Linux display backend ("x11", "wayland", "headless", or "auto").
func SetLinuxBackend(name string) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "x11", "xlib":
		requestedLinuxBackend = linuxBackendX11
	case "wayland", "wl":
		requestedLinuxBackend = linuxBackendWayland
	case "headless", "dummy", "mock":
		requestedLinuxBackend = linuxBackendHeadless
	default:
		requestedLinuxBackend = linuxBackendAuto
	}
}

// ActiveLinuxBackend returns the currently active Linux backend ("x11", "wayland", "headless", or "none").
func ActiveLinuxBackend() string {
	switch activeLinuxBackend {
	case linuxBackendX11:
		return "x11"
	case linuxBackendWayland:
		return "wayland"
	case linuxBackendHeadless:
		return "headless"
	default:
		return "none"
	}
}

func platformInit() error {
	pref := requestedLinuxBackend
	if pref == linuxBackendAuto {
		switch strings.ToLower(strings.TrimSpace(os.Getenv("DER_ATOM_BACKEND"))) {
		case "x11", "xlib":
			pref = linuxBackendX11
		case "wayland", "wl":
			pref = linuxBackendWayland
		case "headless", "dummy":
			pref = linuxBackendHeadless
		}
	}

	switch pref {
	case linuxBackendX11:
		if err := x11Init(); err != nil {
			return err
		}
		activeLinuxBackend = linuxBackendX11
		return nil
	case linuxBackendWayland:
		if err := waylandInit(); err != nil {
			return err
		}
		activeLinuxBackend = linuxBackendWayland
		return nil
	case linuxBackendHeadless:
		activeLinuxBackend = linuxBackendHeadless
		return nil
	default:
		if os.Getenv("WAYLAND_DISPLAY") != "" {
			if err := waylandInit(); err == nil {
				activeLinuxBackend = linuxBackendWayland
				return nil
			}
		}
		if os.Getenv("DISPLAY") != "" {
			if err := x11Init(); err == nil {
				activeLinuxBackend = linuxBackendX11
				return nil
			}
		}
		if err := waylandInit(); err == nil {
			activeLinuxBackend = linuxBackendWayland
			return nil
		}
		if err := x11Init(); err == nil {
			activeLinuxBackend = linuxBackendX11
			return nil
		}
		// Fallback to headless state when running unit tests in an environment without X11/Wayland sockets
		activeLinuxBackend = linuxBackendHeadless
		return nil
	}
}

func platformTerminate() {
	switch activeLinuxBackend {
	case linuxBackendX11:
		x11Terminate()
	case linuxBackendWayland:
		waylandTerminate()
	}
	activeLinuxBackend = linuxBackendAuto
}

func platformGetPrimaryMonitor() uintptr {
	switch activeLinuxBackend {
	case linuxBackendX11:
		return x11GetPrimaryMonitor()
	case linuxBackendWayland:
		return waylandGetPrimaryMonitor()
	default:
		return 1
	}
}

func platformGetVideoMode(monitor uintptr) GLFWvidmode {
	switch activeLinuxBackend {
	case linuxBackendX11:
		return x11GetVideoMode(monitor)
	case linuxBackendWayland:
		return waylandGetVideoMode(monitor)
	default:
		return GLFWvidmode{Width: 1920, Height: 1080, RedBits: 8, GreenBits: 8, BlueBits: 8, RefreshRate: 60}
	}
}

func platformGetVideoModes(monitor uintptr) []GLFWvidmode {
	switch activeLinuxBackend {
	case linuxBackendX11:
		return x11GetVideoModes(monitor)
	case linuxBackendWayland:
		return waylandGetVideoModes(monitor)
	default:
		return []GLFWvidmode{
			{Width: 1280, Height: 720, RedBits: 8, GreenBits: 8, BlueBits: 8, RefreshRate: 60},
			{Width: 1920, Height: 1080, RedBits: 8, GreenBits: 8, BlueBits: 8, RefreshRate: 60},
		}
	}
}

func platformCreateWindow(w *Window, width, height int, title string, monitor uintptr) error {
	switch activeLinuxBackend {
	case linuxBackendX11:
		return x11CreateWindow(w, width, height, title, monitor)
	case linuxBackendWayland:
		return waylandCreateWindow(w, width, height, title, monitor)
	case linuxBackendHeadless:
		w.width = width
		w.height = height
		w.fbWidth = width
		w.fbHeight = height
		return nil
	default:
		return fmt.Errorf("linux display backend not initialized")
	}
}

func platformDestroyWindow(w *Window) {
	switch activeLinuxBackend {
	case linuxBackendX11:
		x11DestroyWindow(w)
	case linuxBackendWayland:
		waylandDestroyWindow(w)
	}
}

func platformPollEvents() {
	switch activeLinuxBackend {
	case linuxBackendX11:
		x11PollEvents()
	case linuxBackendWayland:
		waylandPollEvents()
	}
}

func platformWaitEventsTimeout(timeout float64) {
	switch activeLinuxBackend {
	case linuxBackendX11:
		x11WaitEventsTimeout(timeout)
	case linuxBackendWayland:
		waylandWaitEventsTimeout(timeout)
	default:
		if timeout > 0 && timeout < 0.05 {
			time.Sleep(time.Duration(timeout * float64(time.Second)))
		}
	}
}

func platformSetWindowMonitor(w *Window, monitor uintptr, xpos, ypos, width, height, refreshRate int) {
	switch activeLinuxBackend {
	case linuxBackendX11:
		x11SetWindowMonitor(w, monitor, xpos, ypos, width, height, refreshRate)
	case linuxBackendWayland:
		waylandSetWindowMonitor(w, monitor, xpos, ypos, width, height, refreshRate)
	default:
		w.xpos = xpos
		w.ypos = ypos
		if width > 0 && height > 0 {
			w.dispatchWindowSize(width, height)
			w.dispatchFramebufferSize(width, height)
		}
	}
}

func platformGetWindowSize(w *Window) (int, int) {
	switch activeLinuxBackend {
	case linuxBackendX11:
		return x11GetWindowSize(w)
	case linuxBackendWayland:
		return waylandGetWindowSize(w)
	default:
		return w.width, w.height
	}
}

func platformGetFramebufferSize(w *Window) (int, int) {
	switch activeLinuxBackend {
	case linuxBackendX11:
		return x11GetFramebufferSize(w)
	case linuxBackendWayland:
		return waylandGetFramebufferSize(w)
	default:
		return w.fbWidth, w.fbHeight
	}
}

func platformGetWindowPos(w *Window) (int, int) {
	switch activeLinuxBackend {
	case linuxBackendX11:
		return x11GetWindowPos(w)
	case linuxBackendWayland:
		return waylandGetWindowPos(w)
	default:
		return w.xpos, w.ypos
	}
}

func platformWindowFocused(w *Window) bool {
	if w == nil {
		return false
	}
	return w.focused
}

func platformSetCursorMode(w *Window, mode int) {
	switch activeLinuxBackend {
	case linuxBackendX11:
		x11SetCursorMode(w, mode)
	case linuxBackendWayland:
		waylandSetCursorMode(w, mode)
	}
}

func platformGetCursorPos(w *Window) (float64, float64) {
	switch activeLinuxBackend {
	case linuxBackendX11:
		return x11GetCursorPos(w)
	case linuxBackendWayland:
		return waylandGetCursorPos(w)
	default:
		return w.cursorX, w.cursorY
	}
}

func platformRequiredInstanceExtensions() []string {
	switch activeLinuxBackend {
	case linuxBackendWayland:
		return []string{"VK_KHR_surface", "VK_KHR_wayland_surface"}
	case linuxBackendX11:
		return []string{"VK_KHR_surface", "VK_KHR_xlib_surface"}
	default:
		if os.Getenv("WAYLAND_DISPLAY") != "" {
			return []string{"VK_KHR_surface", "VK_KHR_wayland_surface"}
		}
		return []string{"VK_KHR_surface", "VK_KHR_xlib_surface"}
	}
}

func platformCreateWindowSurface(instance uintptr, w *Window, allocator uintptr, surface *uint64) (int32, error) {
	switch activeLinuxBackend {
	case linuxBackendWayland:
		return waylandCreateWindowSurface(instance, w, allocator, surface)
	case linuxBackendX11:
		return x11CreateWindowSurface(instance, w, allocator, surface)
	default:
		return VK_ERROR_INITIALIZATION_FAILED, fmt.Errorf("no active X11 or Wayland display surface")
	}
}

type unixPollFd struct {
	Fd      int32
	Events  int16
	Revents int16
}

func pollReadable(fd int, timeoutSeconds float64) bool {
	if fd < 0 {
		return false
	}
	ms := int(timeoutSeconds * 1000.0)
	if ms < 0 {
		ms = 0
	}
	pfd := unixPollFd{
		Fd:     int32(fd),
		Events: 0x0001, // POLLIN
	}
	ts := syscall.NsecToTimespec(int64(ms) * int64(time.Millisecond))
	// Use ppoll (SYS_PPOLL) which exists on both linux/amd64 and linux/arm64
	_, _, _ = syscall.Syscall6(
		syscall.SYS_PPOLL,
		uintptr(unsafe.Pointer(&pfd)),
		1,
		uintptr(unsafe.Pointer(&ts)),
		0,
		0,
		0,
	)
	return (pfd.Revents & 0x0001) != 0
}
