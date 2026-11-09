package binding

import (
	"fmt"
	"runtime"
	"sync"
	"time"
	"unsafe"

	"github.com/ebitengine/purego"
)

const (
	GLFW_CLIENT_API      = 0x00022001
	GLFW_NO_API          = 0
	GLFW_TRUE            = 1
	GLFW_FALSE           = 0
	GLFW_FOCUSED         = 0x00020001
	GLFW_ICONIFIED       = 0x00020002
	GLFW_RESIZABLE       = 0x00020003
	GLFW_VISIBLE         = 0x00020004
	GLFW_DECORATED       = 0x00020005
	GLFW_AUTO_ICONIFY    = 0x00020006
	GLFW_FLOATING        = 0x00020007
	GLFW_MAXIMIZED       = 0x00020008
	GLFW_CURSOR          = 0x00033001
	GLFW_CURSOR_NORMAL   = 0x00034001
	GLFW_CURSOR_HIDDEN   = 0x00034002
	GLFW_CURSOR_DISABLED = 0x00034003
	GLFW_DONT_CARE       = -1

	// Error codes
	GLFW_NO_ERROR          = 0
	GLFW_NOT_INITIALIZED   = 0x00010001
	GLFW_INVALID_ENUM      = 0x00010003
	GLFW_INVALID_VALUE     = 0x00010004
	GLFW_API_UNAVAILABLE   = 0x00010006
	GLFW_PLATFORM_ERROR    = 0x00010008
	GLFW_FORMAT_UNAVAILABLE = 0x00010009

	// GLFW Keys
	GLFW_KEY_UNKNOWN      = -1
	GLFW_KEY_SPACE        = 32
	GLFW_KEY_APOSTROPHE   = 39
	GLFW_KEY_COMMA        = 44
	GLFW_KEY_MINUS        = 45
	GLFW_KEY_PERIOD       = 46
	GLFW_KEY_SLASH        = 47
	GLFW_KEY_0            = 48
	GLFW_KEY_1            = 49
	GLFW_KEY_2            = 50
	GLFW_KEY_3            = 51
	GLFW_KEY_4            = 52
	GLFW_KEY_5            = 53
	GLFW_KEY_6            = 54
	GLFW_KEY_7            = 55
	GLFW_KEY_8            = 56
	GLFW_KEY_9            = 57
	GLFW_KEY_SEMICOLON    = 59
	GLFW_KEY_EQUAL        = 61
	GLFW_KEY_A            = 65
	GLFW_KEY_B            = 66
	GLFW_KEY_C            = 67
	GLFW_KEY_D            = 68
	GLFW_KEY_E            = 69
	GLFW_KEY_F            = 70
	GLFW_KEY_G            = 71
	GLFW_KEY_H            = 72
	GLFW_KEY_I            = 73
	GLFW_KEY_J            = 74
	GLFW_KEY_K            = 75
	GLFW_KEY_L            = 76
	GLFW_KEY_M            = 77
	GLFW_KEY_N            = 78
	GLFW_KEY_O            = 79
	GLFW_KEY_P            = 80
	GLFW_KEY_Q            = 81
	GLFW_KEY_R            = 82
	GLFW_KEY_S            = 83
	GLFW_KEY_T            = 84
	GLFW_KEY_U            = 85
	GLFW_KEY_V            = 86
	GLFW_KEY_W            = 87
	GLFW_KEY_X            = 88
	GLFW_KEY_Y            = 89
	GLFW_KEY_Z            = 90
	GLFW_KEY_LEFT_BRACKET  = 91
	GLFW_KEY_BACKSLASH     = 92
	GLFW_KEY_RIGHT_BRACKET = 93
	GLFW_KEY_GRAVE_ACCENT  = 96
	GLFW_KEY_ESCAPE       = 256
	GLFW_KEY_ENTER        = 257
	GLFW_KEY_TAB          = 258
	GLFW_KEY_BACKSPACE    = 259
	GLFW_KEY_INSERT       = 260
	GLFW_KEY_DELETE       = 261
	GLFW_KEY_RIGHT        = 262
	GLFW_KEY_LEFT         = 263
	GLFW_KEY_DOWN         = 264
	GLFW_KEY_UP           = 265
	GLFW_KEY_PAGE_UP      = 266
	GLFW_KEY_PAGE_DOWN    = 267
	GLFW_KEY_HOME         = 268
	GLFW_KEY_END          = 269
	GLFW_KEY_CAPS_LOCK    = 280
	GLFW_KEY_SCROLL_LOCK  = 281
	GLFW_KEY_NUM_LOCK     = 282
	GLFW_KEY_PRINT_SCREEN = 283
	GLFW_KEY_PAUSE        = 284
	GLFW_KEY_F1           = 290
	GLFW_KEY_F2           = 291
	GLFW_KEY_F3           = 292
	GLFW_KEY_F4           = 293
	GLFW_KEY_F5           = 294
	GLFW_KEY_F6           = 295
	GLFW_KEY_F7           = 296
	GLFW_KEY_F8           = 297
	GLFW_KEY_F9           = 298
	GLFW_KEY_F10          = 299
	GLFW_KEY_F11          = 300
	GLFW_KEY_F12          = 301
	GLFW_KEY_LEFT_SHIFT   = 340
	GLFW_KEY_LEFT_CONTROL = 341
	GLFW_KEY_LEFT_ALT     = 342
	GLFW_KEY_LEFT_SUPER   = 343
	GLFW_KEY_RIGHT_SHIFT   = 344
	GLFW_KEY_RIGHT_CONTROL = 345
	GLFW_KEY_RIGHT_ALT     = 346
	GLFW_KEY_RIGHT_SUPER   = 347
	GLFW_KEY_MENU          = 348

	// Modifiers
	GLFW_MOD_SHIFT     = 0x0001
	GLFW_MOD_CONTROL   = 0x0002
	GLFW_MOD_ALT       = 0x0004
	GLFW_MOD_SUPER     = 0x0008
	GLFW_MOD_CAPS_LOCK = 0x0010
	GLFW_MOD_NUM_LOCK  = 0x0020

	// GLFW Actions & Mouse
	GLFW_PRESS               = 1
	GLFW_RELEASE             = 0
	GLFW_REPEAT              = 2
	GLFW_MOUSE_BUTTON_LEFT   = 0
	GLFW_MOUSE_BUTTON_RIGHT  = 1
	GLFW_MOUSE_BUTTON_MIDDLE = 2
	GLFW_MOUSE_BUTTON_4      = 3
	GLFW_MOUSE_BUTTON_5      = 4
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

// Internal pure-Go state replacing libglfw
var (
	glfwMu          sync.RWMutex
	glfwInitialized bool
	glfwStartTime   = time.Now()
	glfwLastErrCode int32
	glfwLastErrText []byte
	glfwHints       = map[int32]int32{
		GLFW_CLIENT_API: GLFW_NO_API,
		GLFW_RESIZABLE:  GLFW_TRUE,
		GLFW_VISIBLE:    GLFW_TRUE,
		GLFW_DECORATED:  GLFW_TRUE,
		GLFW_FOCUSED:    GLFW_TRUE,
	}
	glfwWindows      = map[uintptr]*Window{}
	glfwNextHandleID = uintptr(0x1000)

	// Pinned slices for C-like return pointers
	pinnedVideoMode   GLFWvidmode
	pinnedVideoModes  []GLFWvidmode
	pinnedExtBytes    [][]byte
	pinnedExtPointers []*byte
)

type Window struct {
	handle      uintptr
	callbacks   []uintptr
	cleanup     []func()
	title       string
	width       int
	height      int
	fbWidth     int
	fbHeight    int
	xpos        int
	ypos        int
	shouldClose bool
	focused     bool
	cursorMode  int
	cursorX     float64
	cursorY     float64
	lastRawX    float64
	lastRawY    float64
	hasLastRaw  bool
	monitor     uintptr
	refreshRate int
	keyState    map[int]int

	// Platform-specific native handles
	platformHandle  uintptr
	platformDisplay uintptr
	platformSurface uintptr
	platformLayer   uintptr
	platformExtra   uintptr

	// Direct Go callbacks (zero C-trampoline overhead)
	fbSizeCB      func(w *Window, width, height int)
	keyCB         func(w *Window, key, scancode, action, mods int)
	charCB        func(w *Window, character rune)
	cursorPosCB   func(w *Window, x, y float64)
	mouseButtonCB func(w *Window, button, action, mods int)
	scrollCB      func(w *Window, x, y float64)
	dropCB        func(w *Window, paths []string)
	focusCB       func(focused bool)

	// Raw C callback pointers if registered via GlfwSet*Callback directly
	rawFbSizeCB      uintptr
	rawKeyCB         uintptr
	rawCharCB        uintptr
	rawCursorPosCB   uintptr
	rawMouseButtonCB uintptr
	rawScrollCB      uintptr
	rawDropCB        uintptr
	rawFocusCB       uintptr
}

func init() {
	bindNativeGlfwFunctions()
}

func setGLFWError(code int32, msg string) {
	glfwMu.Lock()
	defer glfwMu.Unlock()
	glfwLastErrCode = code
	glfwLastErrText = CString(msg)
}

func clearGLFWError() {
	glfwMu.Lock()
	defer glfwMu.Unlock()
	glfwLastErrCode = GLFW_NO_ERROR
	glfwLastErrText = nil
}

func lookupWindow(handle uintptr) *Window {
	glfwMu.RLock()
	defer glfwMu.RUnlock()
	return glfwWindows[handle]
}

func registerWindow(w *Window) uintptr {
	glfwMu.Lock()
	defer glfwMu.Unlock()
	if w.handle == 0 {
		if w.platformHandle != 0 {
			w.handle = w.platformHandle
		} else {
			glfwNextHandleID++
			w.handle = glfwNextHandleID
		}
	}
	glfwWindows[w.handle] = w
	return w.handle
}

func unregisterWindow(handle uintptr) *Window {
	glfwMu.Lock()
	defer glfwMu.Unlock()
	w := glfwWindows[handle]
	delete(glfwWindows, handle)
	return w
}

func bindNativeGlfwFunctions() {
	GlfwInit = nativeGlfwInit
	GlfwTerminate = nativeGlfwTerminate
	GlfwWindowHint = nativeGlfwWindowHint
	GlfwGetPrimaryMonitor = nativeGlfwGetPrimaryMonitor
	GlfwGetVideoMode = nativeGlfwGetVideoMode
	GlfwGetVideoModes = nativeGlfwGetVideoModes
	GlfwCreateWindow = nativeGlfwCreateWindow
	GlfwDestroyWindow = nativeGlfwDestroyWindow
	GlfwPollEvents = nativeGlfwPollEvents
	GlfwGetTime = nativeGlfwGetTime
	GlfwWindowShouldClose = nativeGlfwWindowShouldClose
	GlfwSetWindowShouldClose = nativeGlfwSetWindowShouldClose
	GlfwSetMonitor = nativeGlfwSetMonitor
	GlfwGetWindowSize = nativeGlfwGetWindowSize
	GlfwGetFramebufferSize = nativeGlfwGetFramebufferSize
	GlfwSetInputMode = nativeGlfwSetInputMode
	GlfwGetCursorPos = nativeGlfwGetCursorPos
	GlfwSetFramebufferSizeCallback = nativeGlfwSetFramebufferSizeCallback
	GlfwSetKeyCallback = nativeGlfwSetKeyCallback
	GlfwSetCharCallback = nativeGlfwSetCharCallback
	GlfwSetCursorPosCallback = nativeGlfwSetCursorPosCallback
	GlfwSetMouseButtonCallback = nativeGlfwSetMouseButtonCallback
	GlfwSetDropCallback = nativeGlfwSetDropCallback
	GlfwGetError = nativeGlfwGetError
	GlfwWaitEventsTimeout = nativeGlfwWaitEventsTimeout
	GlfwGetWindowPos = nativeGlfwGetWindowPos
	GlfwGetWindowAttrib = nativeGlfwGetWindowAttrib
	GlfwSetWindowFocusCallback = nativeGlfwSetWindowFocusCallback
	GlfwInitVulkanLoader = nativeGlfwInitVulkanLoader
	GlfwVulkanSupported = nativeGlfwVulkanSupported
	GlfwGetRequiredInstanceExtensions = nativeGlfwGetRequiredInstanceExtensions
	GlfwCreateWindowSurface = nativeGlfwCreateWindowSurface
	bindExtraGlfwFunctions()
}

func InitPuregoBindings() error {
	if runtime.GOARCH != "amd64" && runtime.GOARCH != "arm64" {
		return fmt.Errorf("native ABI supports amd64 and arm64 only")
	}
	switch runtime.GOOS {
	case "windows", "darwin", "linux":
	default:
		return fmt.Errorf("unsupported OS %s", runtime.GOOS)
	}
	bindNativeGlfwFunctions()
	glfwLibHandle = 1 // mark windowing bindings initialized
	if err := InitVulkanBindings(); err != nil {
		// Allow windowing/input tests to run even if a system lacks a Vulkan ICD,
		// while returning the Vulkan loader error to callers that require Vulkan.
		return err
	}
	if address, err := lookupLibrary(vkLibHandle, "vkGetInstanceProcAddr"); err == nil && address != 0 {
		GlfwInitVulkanLoader(address)
	}
	return nil
}

func CloseBindings() {
	if glfwInitialized {
		nativeGlfwTerminate()
	}
	if vkLibHandle != 0 {
		_ = closeLibrary(vkLibHandle)
		vkLibHandle = 0
	}
	glfwLibHandle = 0
}

func LastGLFWError(operation string) error {
	var description *byte
	code := GlfwGetError(&description)
	return fmt.Errorf("%s: GLFW %d: %s", operation, code, String(description))
}

func nativeGlfwInit() int32 {
	glfwMu.Lock()
	if glfwInitialized {
		glfwMu.Unlock()
		return GLFW_TRUE
	}
	glfwStartTime = time.Now()
	glfwMu.Unlock()

	if err := platformInit(); err != nil {
		setGLFWError(GLFW_PLATFORM_ERROR, err.Error())
		return GLFW_FALSE
	}

	glfwMu.Lock()
	glfwInitialized = true
	glfwMu.Unlock()
	clearGLFWError()
	return GLFW_TRUE
}

func nativeGlfwTerminate() {
	glfwMu.Lock()
	if !glfwInitialized {
		glfwMu.Unlock()
		return
	}
	handles := make([]uintptr, 0, len(glfwWindows))
	for h := range glfwWindows {
		handles = append(handles, h)
	}
	glfwMu.Unlock()

	for _, h := range handles {
		nativeGlfwDestroyWindow(h)
	}
	platformTerminate()

	glfwMu.Lock()
	glfwInitialized = false
	glfwHints = map[int32]int32{
		GLFW_CLIENT_API: GLFW_NO_API,
		GLFW_RESIZABLE:  GLFW_TRUE,
		GLFW_VISIBLE:    GLFW_TRUE,
		GLFW_DECORATED:  GLFW_TRUE,
		GLFW_FOCUSED:    GLFW_TRUE,
	}
	glfwMu.Unlock()
}

func nativeGlfwWindowHint(hint, value int32) {
	glfwMu.Lock()
	glfwHints[hint] = value
	glfwMu.Unlock()
}

func nativeGlfwGetPrimaryMonitor() uintptr {
	if !glfwInitialized {
		nativeGlfwInit()
	}
	return platformGetPrimaryMonitor()
}

func nativeGlfwGetVideoMode(monitor uintptr) *GLFWvidmode {
	if !glfwInitialized {
		nativeGlfwInit()
	}
	mode := platformGetVideoMode(monitor)
	glfwMu.Lock()
	pinnedVideoMode = mode
	ptr := &pinnedVideoMode
	glfwMu.Unlock()
	return ptr
}

func nativeGlfwGetVideoModes(monitor uintptr, count *int32) *GLFWvidmode {
	if !glfwInitialized {
		nativeGlfwInit()
	}
	modes := platformGetVideoModes(monitor)
	if len(modes) == 0 {
		modes = []GLFWvidmode{platformGetVideoMode(monitor)}
	}
	glfwMu.Lock()
	pinnedVideoModes = make([]GLFWvidmode, len(modes))
	copy(pinnedVideoModes, modes)
	if count != nil {
		*count = int32(len(pinnedVideoModes))
	}
	ptr := &pinnedVideoModes[0]
	glfwMu.Unlock()
	return ptr
}

func nativeGlfwCreateWindow(width, height int32, title *byte, monitor, share uintptr) uintptr {
	if width <= 0 || height <= 0 {
		setGLFWError(GLFW_INVALID_VALUE, "invalid window dimensions")
		return 0
	}
	if !glfwInitialized {
		if nativeGlfwInit() == GLFW_FALSE {
			return 0
		}
	}
	w := &Window{
		title:      String(title),
		width:      int(width),
		height:     int(height),
		fbWidth:    int(width),
		fbHeight:   int(height),
		focused:    true,
		cursorMode: GLFW_CURSOR_NORMAL,
		monitor:    monitor,
		keyState:   make(map[int]int),
	}
	if err := platformCreateWindow(w, int(width), int(height), w.title, monitor); err != nil {
		setGLFWError(GLFW_PLATFORM_ERROR, err.Error())
		return 0
	}
	clearGLFWError()
	return registerWindow(w)
}

func nativeGlfwDestroyWindow(window uintptr) {
	w := unregisterWindow(window)
	if w == nil {
		return
	}
	platformDestroyWindow(w)
	w.handle = 0
	w.platformHandle = 0
	w.callbacks = nil
	for _, release := range w.cleanup {
		release()
	}
	w.cleanup = nil
}

func nativeGlfwPollEvents() {
	if !glfwInitialized {
		return
	}
	platformPollEvents()
}

func nativeGlfwWaitEventsTimeout(timeout float64) {
	if !glfwInitialized {
		return
	}
	if timeout <= 0 {
		platformPollEvents()
		return
	}
	platformWaitEventsTimeout(timeout)
}

func nativeGlfwGetTime() float64 {
	glfwMu.RLock()
	start := glfwStartTime
	glfwMu.RUnlock()
	return time.Since(start).Seconds()
}

func nativeGlfwWindowShouldClose(window uintptr) int32 {
	w := lookupWindow(window)
	if w == nil || !w.shouldClose {
		return GLFW_FALSE
	}
	return GLFW_TRUE
}

func nativeGlfwSetWindowShouldClose(window uintptr, value int32) {
	w := lookupWindow(window)
	if w == nil {
		return
	}
	w.shouldClose = value != 0
}

func nativeGlfwSetMonitor(window, monitor uintptr, xpos, ypos, width, height, refreshRate int32) {
	w := lookupWindow(window)
	if w == nil {
		return
	}
	w.monitor = monitor
	w.refreshRate = int(refreshRate)
	platformSetWindowMonitor(w, monitor, int(xpos), int(ypos), int(width), int(height), int(refreshRate))
}

func nativeGlfwGetWindowSize(window uintptr, width, height *int32) {
	w := lookupWindow(window)
	if w == nil {
		if width != nil {
			*width = 0
		}
		if height != nil {
			*height = 0
		}
		return
	}
	ww, hh := platformGetWindowSize(w)
	if width != nil {
		*width = int32(ww)
	}
	if height != nil {
		*height = int32(hh)
	}
}

func nativeGlfwGetFramebufferSize(window uintptr, width, height *int32) {
	w := lookupWindow(window)
	if w == nil {
		if width != nil {
			*width = 0
		}
		if height != nil {
			*height = 0
		}
		return
	}
	fw, fh := platformGetFramebufferSize(w)
	if width != nil {
		*width = int32(fw)
	}
	if height != nil {
		*height = int32(fh)
	}
}

func nativeGlfwGetWindowPos(window uintptr, xpos, ypos *int32) {
	w := lookupWindow(window)
	if w == nil {
		if xpos != nil {
			*xpos = 0
		}
		if ypos != nil {
			*ypos = 0
		}
		return
	}
	x, y := platformGetWindowPos(w)
	if xpos != nil {
		*xpos = int32(x)
	}
	if ypos != nil {
		*ypos = int32(y)
	}
}

func nativeGlfwGetWindowAttrib(window uintptr, attrib int32) int32 {
	w := lookupWindow(window)
	if w == nil {
		return GLFW_FALSE
	}
	switch attrib {
	case GLFW_FOCUSED:
		if platformWindowFocused(w) {
			return GLFW_TRUE
		}
		return GLFW_FALSE
	case GLFW_VISIBLE:
		return GLFW_TRUE
	default:
		glfwMu.RLock()
		val, ok := glfwHints[attrib]
		glfwMu.RUnlock()
		if ok {
			return val
		}
		return GLFW_FALSE
	}
}

func nativeGlfwSetInputMode(window uintptr, mode, value int32) {
	w := lookupWindow(window)
	if w == nil {
		return
	}
	if mode == GLFW_CURSOR {
		w.cursorMode = int(value)
		w.hasLastRaw = false
		platformSetCursorMode(w, int(value))
	}
}

func nativeGlfwGetCursorPos(window uintptr, xpos, ypos *float64) {
	w := lookupWindow(window)
	if w == nil {
		if xpos != nil {
			*xpos = 0
		}
		if ypos != nil {
			*ypos = 0
		}
		return
	}
	cx, cy := platformGetCursorPos(w)
	if xpos != nil {
		*xpos = cx
	}
	if ypos != nil {
		*ypos = cy
	}
}

func nativeGlfwSetFramebufferSizeCallback(window uintptr, cbfun uintptr) uintptr {
	w := lookupWindow(window)
	if w == nil {
		return 0
	}
	prev := w.rawFbSizeCB
	w.rawFbSizeCB = cbfun
	if cbfun != 0 {
		var fn func(win uintptr, width, height int32)
		purego.RegisterFunc(&fn, cbfun)
		w.fbSizeCB = func(win *Window, width, height int) {
			fn(win.handle, int32(width), int32(height))
		}
	} else {
		w.fbSizeCB = nil
	}
	return prev
}

func nativeGlfwSetKeyCallback(window uintptr, cbfun uintptr) uintptr {
	w := lookupWindow(window)
	if w == nil {
		return 0
	}
	prev := w.rawKeyCB
	w.rawKeyCB = cbfun
	if cbfun != 0 {
		var fn func(win uintptr, key, scancode, action, mods int32)
		purego.RegisterFunc(&fn, cbfun)
		w.keyCB = func(win *Window, key, scancode, action, mods int) {
			fn(win.handle, int32(key), int32(scancode), int32(action), int32(mods))
		}
	} else {
		w.keyCB = nil
	}
	return prev
}

func nativeGlfwSetCharCallback(window uintptr, cbfun uintptr) uintptr {
	w := lookupWindow(window)
	if w == nil {
		return 0
	}
	prev := w.rawCharCB
	w.rawCharCB = cbfun
	if cbfun != 0 {
		var fn func(win uintptr, codepoint uint32)
		purego.RegisterFunc(&fn, cbfun)
		w.charCB = func(win *Window, character rune) {
			fn(win.handle, uint32(character))
		}
	} else {
		w.charCB = nil
	}
	return prev
}

func nativeGlfwSetCursorPosCallback(window uintptr, cbfun uintptr) uintptr {
	w := lookupWindow(window)
	if w == nil {
		return 0
	}
	prev := w.rawCursorPosCB
	w.rawCursorPosCB = cbfun
	if cbfun != 0 {
		var fn func(win uintptr, x, y float64)
		purego.RegisterFunc(&fn, cbfun)
		w.cursorPosCB = func(win *Window, x, y float64) {
			fn(win.handle, x, y)
		}
	} else {
		w.cursorPosCB = nil
	}
	return prev
}

func nativeGlfwSetMouseButtonCallback(window uintptr, cbfun uintptr) uintptr {
	w := lookupWindow(window)
	if w == nil {
		return 0
	}
	prev := w.rawMouseButtonCB
	w.rawMouseButtonCB = cbfun
	if cbfun != 0 {
		var fn func(win uintptr, button, action, mods int32)
		purego.RegisterFunc(&fn, cbfun)
		w.mouseButtonCB = func(win *Window, button, action, mods int) {
			fn(win.handle, int32(button), int32(action), int32(mods))
		}
	} else {
		w.mouseButtonCB = nil
	}
	return prev
}

func nativeGlfwSetScrollCallbackPtr(window uintptr, cbfun uintptr) uintptr {
	w := lookupWindow(window)
	if w == nil {
		return 0
	}
	prev := w.rawScrollCB
	w.rawScrollCB = cbfun
	if cbfun != 0 {
		var fn func(win uintptr, x, y float64)
		purego.RegisterFunc(&fn, cbfun)
		w.scrollCB = func(win *Window, x, y float64) {
			fn(win.handle, x, y)
		}
	} else {
		w.scrollCB = nil
	}
	return prev
}

func nativeGlfwSetDropCallback(window uintptr, cbfun uintptr) uintptr {
	w := lookupWindow(window)
	if w == nil {
		return 0
	}
	prev := w.rawDropCB
	w.rawDropCB = cbfun
	if cbfun != 0 {
		var fn func(win uintptr, count int32, paths **byte)
		purego.RegisterFunc(&fn, cbfun)
		w.dropCB = func(win *Window, paths []string) {
			if len(paths) == 0 {
				fn(win.handle, 0, nil)
				return
			}
			cstrings := make([][]byte, len(paths))
			cpointers := make([]*byte, len(paths))
			for i, p := range paths {
				cstrings[i] = CString(p)
				cpointers[i] = &cstrings[i][0]
			}
			fn(win.handle, int32(len(paths)), &cpointers[0])
			runtime.KeepAlive(cstrings)
			runtime.KeepAlive(cpointers)
		}
	} else {
		w.dropCB = nil
	}
	return prev
}

func nativeGlfwSetWindowFocusCallback(window uintptr, cbfun uintptr) uintptr {
	w := lookupWindow(window)
	if w == nil {
		return 0
	}
	prev := w.rawFocusCB
	w.rawFocusCB = cbfun
	if cbfun != 0 {
		var fn func(win uintptr, focused int32)
		purego.RegisterFunc(&fn, cbfun)
		w.focusCB = func(focused bool) {
			v := int32(GLFW_FALSE)
			if focused {
				v = GLFW_TRUE
			}
			fn(w.handle, v)
		}
	} else {
		w.focusCB = nil
	}
	return prev
}

func nativeGlfwGetError(description **byte) int32 {
	glfwMu.Lock()
	defer glfwMu.Unlock()
	code := glfwLastErrCode
	if description != nil {
		if len(glfwLastErrText) > 0 {
			*description = &glfwLastErrText[0]
		} else {
			*description = nil
		}
	}
	glfwLastErrCode = GLFW_NO_ERROR
	return code
}

func nativeGlfwInitVulkanLoader(loader uintptr) {
	if loader != 0 {
		_ = bind(&VkGetInstanceProcAddr, loader, "vkGetInstanceProcAddr")
	}
}

func nativeGlfwVulkanSupported() int32 {
	if vkLibHandle == 0 && VkGetInstanceProcAddr == nil {
		if err := InitVulkanBindings(); err != nil {
			return GLFW_FALSE
		}
	}
	if VkGetInstanceProcAddr != nil {
		return GLFW_TRUE
	}
	return GLFW_FALSE
}

func nativeGlfwGetRequiredInstanceExtensions(count *uint32) **byte {
	if !glfwInitialized {
		nativeGlfwInit()
	}
	exts := platformRequiredInstanceExtensions()
	glfwMu.Lock()
	defer glfwMu.Unlock()
	if len(exts) == 0 {
		if count != nil {
			*count = 0
		}
		return nil
	}
	pinnedExtBytes = make([][]byte, len(exts))
	pinnedExtPointers = make([]*byte, len(exts))
	for i, ext := range exts {
		pinnedExtBytes[i] = CString(ext)
		pinnedExtPointers[i] = &pinnedExtBytes[i][0]
	}
	if count != nil {
		*count = uint32(len(pinnedExtPointers))
	}
	return &pinnedExtPointers[0]
}

func nativeGlfwCreateWindowSurface(instance uintptr, window uintptr, allocator uintptr, surface *uint64) int32 {
	if surface == nil {
		return VK_ERROR_INITIALIZATION_FAILED
	}
	w := lookupWindow(window)
	if w == nil {
		setGLFWError(GLFW_INVALID_VALUE, "invalid window handle for surface creation")
		return VK_ERROR_INITIALIZATION_FAILED
	}
	if instance == 0 {
		setGLFWError(GLFW_INVALID_VALUE, "invalid VkInstance for surface creation")
		return VK_ERROR_INITIALIZATION_FAILED
	}
	res, err := platformCreateWindowSurface(instance, w, allocator, surface)
	if err != nil {
		setGLFWError(GLFW_PLATFORM_ERROR, err.Error())
	}
	return res
}

func CreateWindow(w, h int, title string, monitor uintptr, share uintptr) (*Window, error) {
	if w <= 0 || h <= 0 {
		setGLFWError(GLFW_INVALID_VALUE, "invalid window dimensions")
		return nil, LastGLFWError("create window")
	}
	if !glfwInitialized {
		if nativeGlfwInit() == GLFW_FALSE {
			return nil, LastGLFWError("create window")
		}
	}
	win := &Window{
		title:      title,
		width:      w,
		height:     h,
		fbWidth:    w,
		fbHeight:   h,
		focused:    true,
		cursorMode: GLFW_CURSOR_NORMAL,
		monitor:    monitor,
		keyState:   make(map[int]int),
	}
	if err := platformCreateWindow(win, w, h, title, monitor); err != nil {
		setGLFWError(GLFW_PLATFORM_ERROR, err.Error())
		return nil, LastGLFWError("create window")
	}
	registerWindow(win)
	return win, nil
}

func (w *Window) ShouldClose() bool {
	if w == nil {
		return true
	}
	return w.shouldClose
}

func (w *Window) Destroy() {
	if w != nil && w.handle != 0 {
		nativeGlfwDestroyWindow(w.handle)
	}
}

func (w *Window) SetShouldClose(v bool) {
	if w != nil {
		w.shouldClose = v
	}
}

func (w *Window) SetMonitor(monitor uintptr, x, y, width, height, refreshRate int) {
	if w != nil {
		nativeGlfwSetMonitor(w.handle, monitor, int32(x), int32(y), int32(width), int32(height), int32(refreshRate))
	}
}

func (w *Window) FramebufferSize() (int, int) {
	if w == nil {
		return 0, 0
	}
	return platformGetFramebufferSize(w)
}

func (w *Window) Size() (int, int) {
	if w == nil {
		return 0, 0
	}
	return platformGetWindowSize(w)
}

func (w *Window) SetInputMode(mode, value int) {
	if w != nil {
		nativeGlfwSetInputMode(w.handle, int32(mode), int32(value))
	}
}

func (w *Window) SetFramebufferSizeCallback(cb func(w *Window, width, height int)) {
	if w != nil {
		w.fbSizeCB = cb
	}
}

func (w *Window) SetKeyCallback(cb func(w *Window, key, scancode, action, mods int)) {
	if w != nil {
		w.keyCB = cb
	}
}

func (w *Window) SetCharCallback(cb func(w *Window, character rune)) {
	if w != nil {
		w.charCB = cb
	}
}

func (w *Window) SetCursorPosCallback(cb func(w *Window, x, y float64)) error {
	if w != nil {
		w.cursorPosCB = cb
	}
	return nil
}

func (w *Window) SetMouseButtonCallback(cb func(w *Window, button, action, mods int)) {
	if w != nil {
		w.mouseButtonCB = cb
	}
}

func (w *Window) SetScrollCallback(cb func(w *Window, x, y float64)) error {
	if w != nil {
		w.scrollCB = cb
	}
	return nil
}

func (w *Window) SetDropCallback(cb func(w *Window, paths []string)) {
	if w != nil {
		w.dropCB = cb
	}
}

func (w *Window) GetCursorPos() (float64, float64) {
	if w == nil {
		return 0, 0
	}
	return platformGetCursorPos(w)
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

func (w *Window) Handle() uintptr {
	if w == nil {
		return 0
	}
	return w.handle
}

func (w *Window) Position() (int, int) {
	if w == nil {
		return 0, 0
	}
	return platformGetWindowPos(w)
}

func (w *Window) Focused() bool {
	if w == nil {
		return false
	}
	return platformWindowFocused(w)
}

func (w *Window) SetFocusCallback(callback func(bool)) {
	if w != nil {
		w.focusCB = callback
	}
}

// Event dispatch helpers invoked by platform backends (and event tests)
func (w *Window) dispatchFramebufferSize(width, height int) {
	if w == nil {
		return
	}
	w.fbWidth = width
	w.fbHeight = height
	if w.fbSizeCB != nil {
		w.fbSizeCB(w, width, height)
	}
}

func (w *Window) dispatchWindowSize(width, height int) {
	if w == nil {
		return
	}
	w.width = width
	w.height = height
}

func (w *Window) dispatchKey(key, scancode, action, mods int) {
	if w == nil {
		return
	}
	if w.keyState != nil && key > 0 {
		if action == GLFW_PRESS && w.keyState[key] == GLFW_PRESS {
			action = GLFW_REPEAT
		}
		if action == GLFW_RELEASE {
			w.keyState[key] = GLFW_RELEASE
		} else {
			w.keyState[key] = GLFW_PRESS
		}
	}
	if w.keyCB != nil {
		w.keyCB(w, key, scancode, action, mods)
	}
}

func (w *Window) dispatchChar(r rune) {
	if w == nil || r < 32 || (r >= 0x7f && r <= 0x9f) {
		return
	}
	if w.charCB != nil {
		w.charCB(w, r)
	}
}

func (w *Window) dispatchCursorPos(x, y float64) {
	if w == nil {
		return
	}
	w.cursorX = x
	w.cursorY = y
	if w.cursorPosCB != nil {
		w.cursorPosCB(w, x, y)
	}
}

func (w *Window) dispatchRelativeCursorDelta(dx, dy float64) {
	if w == nil {
		return
	}
	w.cursorX += dx
	w.cursorY += dy
	if w.cursorPosCB != nil {
		w.cursorPosCB(w, w.cursorX, w.cursorY)
	}
}

func (w *Window) dispatchMouseButton(button, action, mods int) {
	if w == nil {
		return
	}
	if w.mouseButtonCB != nil {
		w.mouseButtonCB(w, button, action, mods)
	}
}

func (w *Window) dispatchScroll(dx, dy float64) {
	if w == nil {
		return
	}
	if w.scrollCB != nil {
		w.scrollCB(w, dx, dy)
	}
}

func (w *Window) dispatchDrop(paths []string) {
	if w == nil {
		return
	}
	if w.dropCB != nil {
		w.dropCB(w, paths)
	}
}

func (w *Window) dispatchFocus(focused bool) {
	if w == nil {
		return
	}
	w.focused = focused
	if w.focusCB != nil {
		w.focusCB(focused)
	}
}