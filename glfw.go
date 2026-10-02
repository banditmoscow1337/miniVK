package binding

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/ebitengine/purego"
)

const (
	GLFW_CLIENT_API      = 0x00022001
	GLFW_NO_API          = 0
	GLFW_TRUE            = 1
	GLFW_FALSE           = 0
	GLFW_CURSOR          = 0x00033001
	GLFW_CURSOR_NORMAL   = 0x00034001
	GLFW_CURSOR_DISABLED = 0x00034003
	GLFW_DONT_CARE       = -1

	// GLFW Keys
	GLFW_KEY_SPACE        = 32
	GLFW_KEY_BACKSPACE    = 259
	GLFW_KEY_DELETE       = 261
	GLFW_KEY_HOME         = 268
	GLFW_KEY_END          = 269
	GLFW_KEY_TAB          = 258
	GLFW_KEY_A            = 65
	GLFW_KEY_D            = 68
	GLFW_KEY_E            = 69
	GLFW_KEY_Q            = 81
	GLFW_KEY_S            = 83
	GLFW_KEY_W            = 87
	GLFW_KEY_R = 82
	GLFW_KEY_F = 70
	GLFW_KEY_Z = 90
	GLFW_KEY_Y = 89
	GLFW_KEY_F5 = 294
	GLFW_KEY_F6 = 295
	GLFW_KEY_F8 = 297
	GLFW_KEY_ESCAPE       = 256
	GLFW_KEY_ENTER        = 257
	GLFW_KEY_RIGHT        = 262
	GLFW_KEY_LEFT         = 263
	GLFW_KEY_DOWN         = 264
	GLFW_KEY_UP           = 265
	GLFW_KEY_LEFT_CONTROL = 341

	// GLFW Actions & Mouse
	GLFW_PRESS               = 1
	GLFW_RELEASE             = 0
	GLFW_REPEAT              = 2
	GLFW_MOUSE_BUTTON_LEFT   = 0
	GLFW_MOUSE_BUTTON_RIGHT  = 1
	GLFW_MOUSE_BUTTON_MIDDLE = 2
)

type GLFWvidmode struct {
	Width       int32
	Height      int32
	RedBits     int32
	GreenBits   int32
	BlueBits    int32
	RefreshRate int32
}

var (
	GlfwInit                       func() int32
	GlfwTerminate                  func()
	GlfwWindowHint                 func(hint, value int32)
	GlfwGetPrimaryMonitor          func() uintptr
	GlfwGetVideoMode               func(monitor uintptr) *GLFWvidmode
	GlfwGetVideoModes              func(monitor uintptr, count *int32) *GLFWvidmode
	GlfwCreateWindow               func(width, height int32, title *byte, monitor, share uintptr) uintptr
	GlfwDestroyWindow              func(window uintptr)
	GlfwPollEvents                 func()
	GlfwGetTime                    func() float64
	GlfwWindowShouldClose          func(window uintptr) int32
	GlfwSetWindowShouldClose       func(window uintptr, value int32)
	GlfwSetMonitor                 func(window, monitor uintptr, xpos, ypos, width, height, refreshRate int32)
	GlfwGetWindowSize              func(window uintptr, width, height *int32)
	GlfwGetFramebufferSize         func(window uintptr, width, height *int32)
	GlfwSetInputMode               func(window uintptr, mode, value int32)
	GlfwGetCursorPos               func(window uintptr, xpos, ypos *float64)
	GlfwSetFramebufferSizeCallback func(window uintptr, cbfun uintptr) uintptr
	GlfwSetKeyCallback             func(window uintptr, cbfun uintptr) uintptr
	GlfwSetCharCallback            func(window uintptr, cbfun uintptr) uintptr
	GlfwSetCursorPosCallback       func(window uintptr, cbfun uintptr) uintptr
	GlfwSetMouseButtonCallback     func(window uintptr, cbfun uintptr) uintptr
	GlfwSetScrollCallback          func(window uintptr, cbfun uintptr) uintptr
	GlfwSetDropCallback            func(window uintptr, cbfun uintptr) uintptr

	glfwLibHandle uintptr
)

func InitPuregoBindings() error {
	if runtime.GOARCH != "amd64" && runtime.GOARCH != "arm64" {
		return fmt.Errorf("native ABI supports amd64 and arm64 only")
	}
	if glfwLibHandle != 0 {
		return nil
	}
	var names []string
	switch runtime.GOOS {
	case "windows":
		names = []string{"glfw3.dll"}
	case "darwin":
		names = []string{"libglfw.3.dylib", "libglfw.dylib"}
	case "linux":
		names = []string{"libglfw.so.3", "libglfw.so"}
	default:
		return fmt.Errorf("unsupported OS %s", runtime.GOOS)
	}
	var err error
	glfwLibHandle, err = loadLibrary(libraryCandidates("DER_ATOM_GLFW", names...))
	if err != nil {
		return err
	}
	ok := false
	defer func() {
		if !ok {
			CloseBindings()
		}
	}()
	if err = bindLibrary(&GlfwInit, glfwLibHandle, "glfwInit"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwTerminate, glfwLibHandle, "glfwTerminate"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwWindowHint, glfwLibHandle, "glfwWindowHint"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwGetPrimaryMonitor, glfwLibHandle, "glfwGetPrimaryMonitor"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwGetVideoMode, glfwLibHandle, "glfwGetVideoMode"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwGetVideoModes, glfwLibHandle, "glfwGetVideoModes"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwCreateWindow, glfwLibHandle, "glfwCreateWindow"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwDestroyWindow, glfwLibHandle, "glfwDestroyWindow"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwPollEvents, glfwLibHandle, "glfwPollEvents"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwGetTime, glfwLibHandle, "glfwGetTime"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwWindowShouldClose, glfwLibHandle, "glfwWindowShouldClose"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwSetWindowShouldClose, glfwLibHandle, "glfwSetWindowShouldClose"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwSetMonitor, glfwLibHandle, "glfwSetWindowMonitor"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwGetWindowSize, glfwLibHandle, "glfwGetWindowSize"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwGetFramebufferSize, glfwLibHandle, "glfwGetFramebufferSize"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwSetInputMode, glfwLibHandle, "glfwSetInputMode"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwSetFramebufferSizeCallback, glfwLibHandle, "glfwSetFramebufferSizeCallback"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwSetKeyCallback, glfwLibHandle, "glfwSetKeyCallback"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwSetCharCallback, glfwLibHandle, "glfwSetCharCallback"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwSetCursorPosCallback, glfwLibHandle, "glfwSetCursorPosCallback"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwSetMouseButtonCallback, glfwLibHandle, "glfwSetMouseButtonCallback"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwSetScrollCallback, glfwLibHandle, "glfwSetScrollCallback"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwSetDropCallback, glfwLibHandle, "glfwSetDropCallback"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwGetCursorPos, glfwLibHandle, "glfwGetCursorPos"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwGetError, glfwLibHandle, "glfwGetError"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwWaitEventsTimeout, glfwLibHandle, "glfwWaitEventsTimeout"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwGetWindowPos, glfwLibHandle, "glfwGetWindowPos"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwGetWindowAttrib, glfwLibHandle, "glfwGetWindowAttrib"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwSetWindowFocusCallback, glfwLibHandle, "glfwSetWindowFocusCallback"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwInitVulkanLoader, glfwLibHandle, "glfwInitVulkanLoader"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwVulkanSupported, glfwLibHandle, "glfwVulkanSupported"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwGetRequiredInstanceExtensions, glfwLibHandle, "glfwGetRequiredInstanceExtensions"); err != nil {
		return err
	}
	if err = bindLibrary(&GlfwCreateWindowSurface, glfwLibHandle, "glfwCreateWindowSurface"); err != nil {
		return err
	}
	if err = InitVulkanBindings(); err != nil {
		return err
	}
	address, err := lookupLibrary(vkLibHandle, "vkGetInstanceProcAddr")
	if err != nil {
		return err
	}
	GlfwInitVulkanLoader(address)
	ok = true
	return nil
}
func CloseBindings() {
	if vkLibHandle != 0 {
		_ = closeLibrary(vkLibHandle)
		vkLibHandle = 0
	}
	if glfwLibHandle != 0 {
		_ = closeLibrary(glfwLibHandle)
		glfwLibHandle = 0
	}
}
func LastGLFWError(operation string) error {
	var description *byte
	code := GlfwGetError(&description)
	return fmt.Errorf("%s: GLFW %d: %s", operation, code, String(description))
}

type Window struct {
	handle    uintptr
	callbacks []uintptr
	cleanup   []func()
}

func CreateWindow(w, h int, title string, monitor uintptr, share uintptr) (*Window, error) {
	text := CString(title)
	hWin := GlfwCreateWindow(int32(w), int32(h), &text[0], monitor, share)
	runtime.KeepAlive(text)
	if hWin == 0 {
		return nil, LastGLFWError("create window")
	}
	return &Window{handle: hWin}, nil
}

func (w *Window) ShouldClose() bool { return GlfwWindowShouldClose(w.handle) != 0 }
func (w *Window) Destroy() {
	if w != nil && w.handle != 0 {
		GlfwDestroyWindow(w.handle)
		w.handle = 0
		w.callbacks = nil
		for _, release := range w.cleanup {
			release()
		}
		w.cleanup = nil
	}
}
func (w *Window) SetShouldClose(v bool) {
	val := int32(0)
	if v {
		val = 1
	}
	GlfwSetWindowShouldClose(w.handle, val)
}
func (w *Window) SetMonitor(monitor uintptr, x, y, width, height, refreshRate int) {
	GlfwSetMonitor(w.handle, monitor, int32(x), int32(y), int32(width), int32(height), int32(refreshRate))
}
func (w *Window) FramebufferSize() (int, int) {
	var width, height int32
	GlfwGetFramebufferSize(w.handle, &width, &height)
	return int(width), int(height)
}
func (w *Window) Size() (int, int) {
	var width, height int32
	GlfwGetWindowSize(w.handle, &width, &height)
	return int(width), int(height)
}
func (w *Window) SetInputMode(mode, value int) {
	GlfwSetInputMode(w.handle, int32(mode), int32(value))
}
func (w *Window) SetFramebufferSizeCallback(cb func(w *Window, width, height int)) {
	cbfun := purego.NewCallback(func(win uintptr, width, height int32) uintptr {
		cb(w, int(width), int(height))
		return 0
	})
	w.callbacks = append(w.callbacks, cbfun)
	GlfwSetFramebufferSizeCallback(w.handle, cbfun)
}

func (w *Window) SetKeyCallback(cb func(w *Window, key, scancode, action, mods int)) {
	cbfun := purego.NewCallback(func(win uintptr, key, scancode, action, mods int32) uintptr {
		cb(w, int(key), int(scancode), int(action), int(mods))
		return 0
	})
	w.callbacks = append(w.callbacks, cbfun)
	GlfwSetKeyCallback(w.handle, cbfun)
}

func (w *Window) SetCharCallback(cb func(w *Window, character rune)) {
	cbfun := purego.NewCallback(func(win uintptr, codepoint uint32) uintptr {
		cb(w, rune(codepoint))
		return 0
	})
	w.callbacks = append(w.callbacks, cbfun)
	GlfwSetCharCallback(w.handle, cbfun)
}

func (w *Window) SetCursorPosCallback(cb func(w *Window, x, y float64)) error {
	pointer, release, err := floatCallback(func(window uintptr, x, y float64) { cb(w, x, y) })
	if err != nil {
		return err
	}
	w.callbacks = append(w.callbacks, pointer)
	if release != nil {
		w.cleanup = append(w.cleanup, release)
	}
	GlfwSetCursorPosCallback(w.handle, pointer)
	return nil
}

func (w *Window) SetMouseButtonCallback(cb func(w *Window, button, action, mods int)) {
	cbfun := purego.NewCallback(func(win uintptr, button, action, mods int32) uintptr {
		cb(w, int(button), int(action), int(mods))
		return 0
	})
	w.callbacks = append(w.callbacks, cbfun)
	GlfwSetMouseButtonCallback(w.handle, cbfun)
}

func (w *Window) SetScrollCallback(cb func(w *Window, x, y float64)) error {
	pointer, release, err := floatCallback(func(window uintptr, x, y float64) { cb(w, x, y) })
	if err != nil {
		return err
	}
	w.callbacks = append(w.callbacks, pointer)
	if release != nil {
		w.cleanup = append(w.cleanup, release)
	}
	GlfwSetScrollCallback(w.handle, pointer)
	return nil
}

func (w *Window) SetDropCallback(cb func(w *Window, paths []string)) {
	cbfun := purego.NewCallback(func(win uintptr, count int32, rawPaths **byte) uintptr {
		if count <= 0 || rawPaths == nil {
			cb(w, nil)
			return 0
		}
		pointers := unsafe.Slice(rawPaths, int(count))
		paths := make([]string, 0, len(pointers))
		for _, pointer := range pointers {
			if pointer == nil {
				continue
			}
			var value []byte
			for current := pointer; *current != 0; current = (*byte)(unsafe.Add(unsafe.Pointer(current), 1)) {
				value = append(value, *current)
			}
			paths = append(paths, string(value))
		}
		cb(w, paths)
		return 0
	})
	w.callbacks = append(w.callbacks, cbfun)
	GlfwSetDropCallback(w.handle, cbfun)
}

func (w *Window) GetCursorPos() (float64, float64) {
	var x, y float64
	GlfwGetCursorPos(w.handle, &x, &y)
	return x, y
}

func GetVideoModes(monitor uintptr) []GLFWvidmode {
	var count int32
	ptr := GlfwGetVideoModes(monitor, &count)
	if ptr == nil || count <= 0 {
		return nil
	}
	slice := unsafe.Slice(ptr, int(count))
	res := make([]GLFWvidmode, count)
	copy(res, slice)
	return res
}

func GetVideoMode(monitor uintptr) *GLFWvidmode {
	return GlfwGetVideoMode(monitor)
}

func (w *Window) Handle() uintptr { return w.handle }
func (w *Window) Position() (int, int) {
	var x, y int32
	GlfwGetWindowPos(w.handle, &x, &y)
	return int(x), int(y)
}
func (w *Window) Focused() bool { return GlfwGetWindowAttrib(w.handle, 0x00020001) != 0 }
func (w *Window) SetFocusCallback(callback func(bool)) {
	pointer := purego.NewCallback(func(window uintptr, focused int32) uintptr { callback(focused != 0); return 0 })
	w.callbacks = append(w.callbacks, pointer)
	GlfwSetWindowFocusCallback(w.handle, pointer)
}