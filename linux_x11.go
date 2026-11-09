//go:build linux

package binding

import (
	"fmt"
	"net/url"
	"runtime"
	"strings"
	"unsafe"
)

const (
	x11KeyPress         = 2
	x11KeyRelease       = 3
	x11ButtonPress      = 4
	x11ButtonRelease    = 5
	x11MotionNotify     = 6
	x11FocusIn          = 9
	x11FocusOut         = 10
	x11ConfigureNotify  = 22
	x11ClientMessage    = 33
	x11SelectionNotify  = 31

	x11KeyPressMask        = 1 << 0
	x11KeyReleaseMask      = 1 << 1
	x11ButtonPressMask     = 1 << 2
	x11ButtonReleaseMask   = 1 << 3
	x11PointerMotionMask   = 1 << 6
	x11StructureNotifyMask = 1 << 17
	x11FocusChangeMask     = 1 << 21

	x11CWBorderPixel = 1 << 3
	x11CWEventMask   = 1 << 11

	x11InputOutput = 1
	x11PropModeReplace = 0
	x11AnyPropertyType = 0
	x11CurrentTime     = 0
	x11GrabModeAsync   = 1

	x11ShiftMask   = 1 << 0
	x11LockMask    = 1 << 1
	x11ControlMask = 1 << 2
	x11Mod1Mask    = 1 << 3
	x11Mod2Mask    = 1 << 4
	x11Mod4Mask    = 1 << 6
)

type xSetWindowAttributes struct {
	BackgroundPixmap  uintptr
	BackgroundPixel   uintptr
	BorderPixmap      uintptr
	BorderPixel       uintptr
	BitGravity        int32
	WinGravity        int32
	BackingStore      int32
	BackingPlanes     uintptr
	BackingPixel      uintptr
	SaveUnder         int32
	EventMask         int64
	DoNotPropagateMask int64
	OverrideRedirect  int32
	Colormap          uintptr
	Cursor            uintptr
}

type xWindowAttributes struct {
	X, Y                int32
	Width, Height       int32
	BorderWidth         int32
	Depth               int32
	Visual              uintptr
	Root                uintptr
	Class               int32
	BitGravity          int32
	WinGravity          int32
	BackingStore        int32
	BackingPlanes       uintptr
	BackingPixel        uintptr
	SaveUnder           int32
	Colormap            uintptr
	MapInstalled        int32
	MapState            int32
	AllEventMasks       int64
	YourEventMask       int64
	DoNotPropagateMask  int64
	OverrideRedirect    int32
	Screen              uintptr
}

type xEvent struct {
	Data [24]uintptr
}

type xKeyEvent struct {
	Type        int32
	_           int32
	Serial      uintptr
	SendEvent   int32
	_           int32
	Display     uintptr
	Window      uintptr
	Root        uintptr
	Subwindow   uintptr
	Time        uintptr
	X, Y        int32
	XRoot, YRoot int32
	State       uint32
	Keycode     uint32
	SameScreen  int32
	_           int32
}

type xButtonEvent struct {
	Type        int32
	_           int32
	Serial      uintptr
	SendEvent   int32
	_           int32
	Display     uintptr
	Window      uintptr
	Root        uintptr
	Subwindow   uintptr
	Time        uintptr
	X, Y        int32
	XRoot, YRoot int32
	State       uint32
	Button      uint32
	SameScreen  int32
	_           int32
}

type xMotionEvent struct {
	Type        int32
	_           int32
	Serial      uintptr
	SendEvent   int32
	_           int32
	Display     uintptr
	Window      uintptr
	Root        uintptr
	Subwindow   uintptr
	Time        uintptr
	X, Y        int32
	XRoot, YRoot int32
	State       uint32
	IsHint      byte
	_           [3]byte
	SameScreen  int32
}

type xConfigureEvent struct {
	Type              int32
	_                 int32
	Serial            uintptr
	SendEvent         int32
	_                 int32
	Display           uintptr
	Event             uintptr
	Window            uintptr
	X, Y              int32
	Width, Height     int32
	BorderWidth       int32
	_                 int32
	Above             uintptr
	OverrideRedirect  int32
	_                 int32
}

type xClientMessageEvent struct {
	Type        int32
	_           int32
	Serial      uintptr
	SendEvent   int32
	_           int32
	Display     uintptr
	Window      uintptr
	MessageType uintptr
	Format      int32
	_           int32
	Data        [5]uintptr
}

type xSelectionEvent struct {
	Type      int32
	_         int32
	Serial    uintptr
	SendEvent int32
	_         int32
	Display   uintptr
	Requestor uintptr
	Selection uintptr
	Target    uintptr
	Property  uintptr
	Time      uintptr
}

type xrrScreenSize struct {
	Width   int32
	Height  int32
	MWidth  int32
	MHeight int32
}

var (
	x11Lib    uintptr
	xrandrLib uintptr
	xfixesLib uintptr

	xOpenDisplay          func(displayName *byte) uintptr
	xCloseDisplay         func(dpy uintptr) int32
	xDefaultScreen        func(dpy uintptr) int32
	xRootWindow           func(dpy uintptr, screen int32) uintptr
	xDefaultVisual        func(dpy uintptr, screen int32) uintptr
	xDefaultDepth         func(dpy uintptr, screen int32) int32
	xDisplayWidth         func(dpy uintptr, screen int32) int32
	xDisplayHeight        func(dpy uintptr, screen int32) int32
	xConnectionNumber     func(dpy uintptr) int32
	xCreateWindow         func(dpy, parent uintptr, x, y int32, width, height, borderWidth uint32, depth int32, class uint32, visual uintptr, valuemask uintptr, attributes *xSetWindowAttributes) uintptr
	xDestroyWindow        func(dpy, w uintptr) int32
	xMapWindow            func(dpy, w uintptr) int32
	xUnmapWindow          func(dpy, w uintptr) int32
	xStoreName            func(dpy, w uintptr, name *byte) int32
	xInternAtom           func(dpy uintptr, name *byte, onlyIfExists int32) uintptr
	xSetWMProtocols       func(dpy, w uintptr, protocols *uintptr, count int32) int32
	xChangeProperty       func(dpy, w, property, propType uintptr, format, mode int32, data uintptr, nelements int32) int32
	xGetWindowProperty    func(dpy, w, property uintptr, longOffset, longLength int64, deleteProp int32, reqType uintptr, actualType *uintptr, actualFormat *int32, nitems, bytesAfter *uintptr, prop **byte) int32
	xConvertSelection     func(dpy, selection, target, property, requestor, time uintptr) int32
	xPending              func(dpy uintptr) int32
	xNextEvent            func(dpy uintptr, ev *xEvent) int32
	xSendEvent            func(dpy, w uintptr, propagate int32, eventMask int64, ev *xEvent) int32
	xFlush                func(dpy uintptr) int32
	xSync                 func(dpy uintptr, discard int32) int32
	xFree                 func(data uintptr) int32
	xLookupKeysym         func(keyEvent *xKeyEvent, index int32) uintptr
	xLookupString         func(keyEvent *xKeyEvent, buffer *byte, bytesBuffer int32, keysym *uintptr, status uintptr) int32
	xGetWindowAttributes  func(dpy, w uintptr, attr *xWindowAttributes) int32
	xTranslateCoordinates func(dpy, srcW, destW uintptr, srcX, srcY int32, destX, destY *int32, childReturn *uintptr) int32
	xQueryPointer         func(dpy, w uintptr, rootReturn, childReturn *uintptr, rootX, rootY, winX, winY *int32, maskReturn *uint32) int32
	xWarpPointer          func(dpy, srcW, destW uintptr, srcX, srcY int32, srcWidth, srcHeight uint32, destX, destY int32) int32
	xMoveResizeWindow     func(dpy, w uintptr, x, y int32, width, height uint32) int32
	xGrabPointer          func(dpy, grabW uintptr, ownerEvents int32, eventMask uint32, pointerMode, keyboardMode int32, confineTo, cursor, time uintptr) int32
	xUngrabPointer        func(dpy, time uintptr) int32

	xrrGetScreenInfo            func(dpy, w uintptr) uintptr
	xrrFreeScreenConfigInfo     func(config uintptr)
	xrrConfigSizes              func(config uintptr, nsizes *int32) *xrrScreenSize
	xrrConfigRates              func(config uintptr, sizeID int32, nrates *int32) *int16
	xrrConfigCurrentConfig      func(config uintptr, rotation *uint16) uint16
	xrrConfigCurrentRate        func(config uintptr) int16
	xrrSetScreenConfigAndRate   func(dpy, config, draw uintptr, sizeIndex int32, rotation uint16, rate int16, timestamp uintptr) int32

	xfixesHideCursor func(dpy, w uintptr)
	xfixesShowCursor func(dpy, w uintptr)

	x11Dpy             uintptr
	x11Screen          int32
	x11Root            uintptr
	x11AtomWMProtocols uintptr
	x11AtomWMDeleteWin uintptr
	x11AtomNetWMState  uintptr
	x11AtomNetWMFull   uintptr
	x11AtomXdndAware   uintptr
	x11AtomXdndEnter   uintptr
	x11AtomXdndPos     uintptr
	x11AtomXdndStatus  uintptr
	x11AtomXdndActionCopy uintptr
	x11AtomXdndDrop    uintptr
	x11AtomXdndFin     uintptr
	x11AtomXdndSel     uintptr
	x11AtomUriList     uintptr
	x11AtomXdndProp    uintptr
	x11XdndSource      uintptr
)

func x11Atom(name string) uintptr {
	b := CString(name)
	a := xInternAtom(x11Dpy, &b[0], 0)
	runtime.KeepAlive(b)
	return a
}

func x11Init() error {
	if x11Dpy != 0 {
		return nil
	}
	var err error
	if x11Lib == 0 {
		x11Lib, err = loadLibrary([]string{"libX11.so.6", "libX11.so"})
		if err != nil {
			return fmt.Errorf("load libX11: %w", err)
		}
		for _, item := range []struct {
			fn   any
			name string
		}{
			{&xOpenDisplay, "XOpenDisplay"},
			{&xCloseDisplay, "XCloseDisplay"},
			{&xDefaultScreen, "XDefaultScreen"},
			{&xRootWindow, "XRootWindow"},
			{&xDefaultVisual, "XDefaultVisual"},
			{&xDefaultDepth, "XDefaultDepth"},
			{&xDisplayWidth, "XDisplayWidth"},
			{&xDisplayHeight, "XDisplayHeight"},
			{&xConnectionNumber, "XConnectionNumber"},
			{&xCreateWindow, "XCreateWindow"},
			{&xDestroyWindow, "XDestroyWindow"},
			{&xMapWindow, "XMapWindow"},
			{&xUnmapWindow, "XUnmapWindow"},
			{&xStoreName, "XStoreName"},
			{&xInternAtom, "XInternAtom"},
			{&xSetWMProtocols, "XSetWMProtocols"},
			{&xChangeProperty, "XChangeProperty"},
			{&xGetWindowProperty, "XGetWindowProperty"},
			{&xConvertSelection, "XConvertSelection"},
			{&xPending, "XPending"},
			{&xNextEvent, "XNextEvent"},
			{&xSendEvent, "XSendEvent"},
			{&xFlush, "XFlush"},
			{&xSync, "XSync"},
			{&xFree, "XFree"},
			{&xLookupKeysym, "XLookupKeysym"},
			{&xLookupString, "XLookupString"},
			{&xGetWindowAttributes, "XGetWindowAttributes"},
			{&xTranslateCoordinates, "XTranslateCoordinates"},
			{&xQueryPointer, "XQueryPointer"},
			{&xWarpPointer, "XWarpPointer"},
			{&xMoveResizeWindow, "XMoveResizeWindow"},
			{&xGrabPointer, "XGrabPointer"},
			{&xUngrabPointer, "XUngrabPointer"},
		} {
			if err := bindLibrary(item.fn, x11Lib, item.name); err != nil {
				return err
			}
		}
	}

	if xrandrLib == 0 {
		if h, err := loadLibrary([]string{"libXrandr.so.2", "libXrandr.so"}); err == nil {
			xrandrLib = h
			_ = bindLibrary(&xrrGetScreenInfo, xrandrLib, "XRRGetScreenInfo")
			_ = bindLibrary(&xrrFreeScreenConfigInfo, xrandrLib, "XRRFreeScreenConfigInfo")
			_ = bindLibrary(&xrrConfigSizes, xrandrLib, "XRRConfigSizes")
			_ = bindLibrary(&xrrConfigRates, xrandrLib, "XRRConfigRates")
			_ = bindLibrary(&xrrConfigCurrentConfig, xrandrLib, "XRRConfigCurrentConfiguration")
			_ = bindLibrary(&xrrConfigCurrentRate, xrandrLib, "XRRConfigCurrentRate")
			_ = bindLibrary(&xrrSetScreenConfigAndRate, xrandrLib, "XRRSetScreenConfigAndRate")
		}
	}

	if xfixesLib == 0 {
		if h, err := loadLibrary([]string{"libXfixes.so.3", "libXfixes.so"}); err == nil {
			xfixesLib = h
			_ = bindLibrary(&xfixesHideCursor, xfixesLib, "XFixesHideCursor")
			_ = bindLibrary(&xfixesShowCursor, xfixesLib, "XFixesShowCursor")
		}
	}

	x11Dpy = xOpenDisplay(nil)
	if x11Dpy == 0 {
		return fmt.Errorf("XOpenDisplay failed (check DISPLAY environment variable)")
	}
	x11Screen = xDefaultScreen(x11Dpy)
	x11Root = xRootWindow(x11Dpy, x11Screen)

	x11AtomWMProtocols = x11Atom("WM_PROTOCOLS")
	x11AtomWMDeleteWin = x11Atom("WM_DELETE_WINDOW")
	x11AtomNetWMState = x11Atom("_NET_WM_STATE")
	x11AtomNetWMFull = x11Atom("_NET_WM_STATE_FULLSCREEN")
	x11AtomXdndAware = x11Atom("XdndAware")
	x11AtomXdndEnter = x11Atom("XdndEnter")
	x11AtomXdndPos = x11Atom("XdndPosition")
	x11AtomXdndStatus = x11Atom("XdndStatus")
	x11AtomXdndActionCopy = x11Atom("XdndActionCopy")
	x11AtomXdndDrop = x11Atom("XdndDrop")
	x11AtomXdndFin = x11Atom("XdndFinished")
	x11AtomXdndSel = x11Atom("XdndSelection")
	x11AtomUriList = x11Atom("text/uri-list")
	x11AtomXdndProp = x11Atom("MINIVK_XDND_DATA")
	return nil
}

func x11Terminate() {
	if x11Dpy != 0 {
		xCloseDisplay(x11Dpy)
		x11Dpy = 0
	}
}

func x11GetPrimaryMonitor() uintptr {
	return 1
}

func x11GetVideoMode(monitor uintptr) GLFWvidmode {
	if x11Dpy == 0 {
		return GLFWvidmode{Width: 1920, Height: 1080, RedBits: 8, GreenBits: 8, BlueBits: 8, RefreshRate: 60}
	}
	if xrrGetScreenInfo != nil && xrrConfigSizes != nil && xrrConfigCurrentConfig != nil {
		cfg := xrrGetScreenInfo(x11Dpy, x11Root)
		if cfg != 0 {
			defer xrrFreeScreenConfigInfo(cfg)
			var nsizes int32
			sizesPtr := xrrConfigSizes(cfg, &nsizes)
			var rot uint16
			idx := int(xrrConfigCurrentConfig(cfg, &rot))
			rate := int32(xrrConfigCurrentRate(cfg))
			if rate <= 0 {
				rate = 60
			}
			if sizesPtr != nil && idx >= 0 && idx < int(nsizes) {
				sizes := unsafe.Slice(sizesPtr, int(nsizes))
				return GLFWvidmode{
					Width:       sizes[idx].Width,
					Height:      sizes[idx].Height,
					RedBits:     8,
					GreenBits:   8,
					BlueBits:    8,
					RefreshRate: rate,
				}
			}
		}
	}
	w := xDisplayWidth(x11Dpy, x11Screen)
	h := xDisplayHeight(x11Dpy, x11Screen)
	if w <= 0 || h <= 0 {
		w, h = 1920, 1080
	}
	return GLFWvidmode{Width: w, Height: h, RedBits: 8, GreenBits: 8, BlueBits: 8, RefreshRate: 60}
}

func x11GetVideoModes(monitor uintptr) []GLFWvidmode {
	if x11Dpy != 0 && xrrGetScreenInfo != nil && xrrConfigSizes != nil && xrrConfigRates != nil {
		cfg := xrrGetScreenInfo(x11Dpy, x11Root)
		if cfg != 0 {
			defer xrrFreeScreenConfigInfo(cfg)
			var nsizes int32
			sizesPtr := xrrConfigSizes(cfg, &nsizes)
			if sizesPtr != nil && nsizes > 0 {
				sizes := unsafe.Slice(sizesPtr, int(nsizes))
				var modes []GLFWvidmode
				for i := 0; i < int(nsizes); i++ {
					var nrates int32
					ratesPtr := xrrConfigRates(cfg, int32(i), &nrates)
					if ratesPtr != nil && nrates > 0 {
						rates := unsafe.Slice(ratesPtr, int(nrates))
						for _, r := range rates {
							rate := int32(r)
							if rate <= 0 {
								rate = 60
							}
							modes = append(modes, GLFWvidmode{
								Width:       sizes[i].Width,
								Height:      sizes[i].Height,
								RedBits:     8,
								GreenBits:   8,
								BlueBits:    8,
								RefreshRate: rate,
							})
						}
					} else {
						modes = append(modes, GLFWvidmode{
							Width:       sizes[i].Width,
							Height:      sizes[i].Height,
							RedBits:     8,
							GreenBits:   8,
							BlueBits:    8,
							RefreshRate: 60,
						})
					}
				}
				if len(modes) > 0 {
					return modes
				}
			}
		}
	}
	return []GLFWvidmode{x11GetVideoMode(monitor)}
}

func x11CreateWindow(w *Window, width, height int, title string, monitor uintptr) error {
	if x11Dpy == 0 {
		if err := x11Init(); err != nil {
			return err
		}
	}
	visual := xDefaultVisual(x11Dpy, x11Screen)
	depth := xDefaultDepth(x11Dpy, x11Screen)
	swa := xSetWindowAttributes{
		BorderPixel: 0,
		EventMask:   x11KeyPressMask | x11KeyReleaseMask | x11ButtonPressMask | x11ButtonReleaseMask | x11PointerMotionMask | x11StructureNotifyMask | x11FocusChangeMask,
	}
	xwin := xCreateWindow(
		x11Dpy,
		x11Root,
		0, 0,
		uint32(width), uint32(height),
		0,
		depth,
		x11InputOutput,
		visual,
		x11CWBorderPixel|x11CWEventMask,
		&swa,
	)
	if xwin == 0 {
		return fmt.Errorf("XCreateWindow failed")
	}
	cTitle := CString(title)
	xStoreName(x11Dpy, xwin, &cTitle[0])
	runtime.KeepAlive(cTitle)

	proto := x11AtomWMDeleteWin
	xSetWMProtocols(x11Dpy, xwin, &proto, 1)

	xdndVer := uintptr(5)
	atomType := uintptr(4) // XA_ATOM
	xChangeProperty(x11Dpy, xwin, x11AtomXdndAware, atomType, 32, x11PropModeReplace, uintptr(unsafe.Pointer(&xdndVer)), 1)

	xMapWindow(x11Dpy, xwin)
	if monitor != 0 {
		x11SetFullscreen(xwin, true)
	}
	xFlush(x11Dpy)

	w.platformDisplay = x11Dpy
	w.platformHandle = xwin
	w.handle = xwin
	return nil
}

func x11SetFullscreen(xwin uintptr, fullscreen bool) {
	if x11Dpy == 0 || xwin == 0 {
		return
	}
	action := uintptr(0) // _NET_WM_STATE_REMOVE
	if fullscreen {
		action = 1 // _NET_WM_STATE_ADD
	}
	var ev xClientMessageEvent
	ev.Type = x11ClientMessage
	ev.Window = xwin
	ev.MessageType = x11AtomNetWMState
	ev.Format = 32
	ev.Data[0] = action
	ev.Data[1] = x11AtomNetWMFull
	ev.Data[2] = 0
	ev.Data[3] = 1
	xSendEvent(x11Dpy, x11Root, 0, (1<<20)|(1<<19), (*xEvent)(unsafe.Pointer(&ev)))
	xFlush(x11Dpy)
}

func x11DestroyWindow(w *Window) {
	if w != nil && w.platformHandle != 0 && x11Dpy != 0 {
		if w.cursorMode == GLFW_CURSOR_DISABLED {
			x11SetCursorMode(w, GLFW_CURSOR_NORMAL)
		}
		xDestroyWindow(x11Dpy, w.platformHandle)
		xFlush(x11Dpy)
		w.platformHandle = 0
	}
}

func x11PollEvents() {
	if x11Dpy == 0 {
		return
	}
	var ev xEvent
	for xPending(x11Dpy) > 0 {
		xNextEvent(x11Dpy, &ev)
		x11ProcessEvent(&ev)
	}
}

func x11WaitEventsTimeout(timeout float64) {
	if x11Dpy == 0 {
		return
	}
	if xPending(x11Dpy) == 0 {
		fd := int(xConnectionNumber(x11Dpy))
		pollReadable(fd, timeout)
	}
	x11PollEvents()
}

func x11SetWindowMonitor(w *Window, monitor uintptr, xpos, ypos, width, height, refreshRate int) {
	if w == nil || w.platformHandle == 0 || x11Dpy == 0 {
		return
	}
	if monitor != 0 {
		x11SetFullscreen(w.platformHandle, true)
	} else {
		x11SetFullscreen(w.platformHandle, false)
		xMoveResizeWindow(x11Dpy, w.platformHandle, int32(xpos), int32(ypos), uint32(width), uint32(height))
		xFlush(x11Dpy)
	}
	w.xpos = xpos
	w.ypos = ypos
	if width > 0 && height > 0 {
		w.dispatchWindowSize(width, height)
		w.dispatchFramebufferSize(width, height)
	}
}

func x11GetWindowSize(w *Window) (int, int) {
	if w == nil {
		return 0, 0
	}
	if x11Dpy != 0 && w.platformHandle != 0 {
		var attr xWindowAttributes
		if xGetWindowAttributes(x11Dpy, w.platformHandle, &attr) != 0 && attr.Width > 0 && attr.Height > 0 {
			w.width = int(attr.Width)
			w.height = int(attr.Height)
		}
	}
	return w.width, w.height
}

func x11GetFramebufferSize(w *Window) (int, int) {
	ww, hh := x11GetWindowSize(w)
	w.fbWidth = ww
	w.fbHeight = hh
	return ww, hh
}

func x11GetWindowPos(w *Window) (int, int) {
	if w == nil {
		return 0, 0
	}
	if x11Dpy != 0 && w.platformHandle != 0 {
		var dx, dy int32
		var child uintptr
		if xTranslateCoordinates(x11Dpy, w.platformHandle, x11Root, 0, 0, &dx, &dy, &child) != 0 {
			w.xpos = int(dx)
			w.ypos = int(dy)
		}
	}
	return w.xpos, w.ypos
}

func x11SetCursorMode(w *Window, mode int) {
	if w == nil || w.platformHandle == 0 || x11Dpy == 0 {
		return
	}
	if mode == GLFW_CURSOR_DISABLED {
		if xfixesHideCursor != nil {
			xfixesHideCursor(x11Dpy, w.platformHandle)
		}
		xGrabPointer(
			x11Dpy,
			w.platformHandle,
			1,
			x11ButtonPressMask|x11ButtonReleaseMask|x11PointerMotionMask,
			x11GrabModeAsync,
			x11GrabModeAsync,
			w.platformHandle,
			0,
			x11CurrentTime,
		)
		cx := int32(w.width / 2)
		cy := int32(w.height / 2)
		xWarpPointer(x11Dpy, 0, w.platformHandle, 0, 0, 0, 0, cx, cy)
		w.lastRawX = float64(cx)
		w.lastRawY = float64(cy)
		w.hasLastRaw = true
		xFlush(x11Dpy)
	} else {
		xUngrabPointer(x11Dpy, x11CurrentTime)
		if xfixesShowCursor != nil {
			xfixesShowCursor(x11Dpy, w.platformHandle)
		}
		xFlush(x11Dpy)
	}
}

func x11GetCursorPos(w *Window) (float64, float64) {
	if w == nil {
		return 0, 0
	}
	if w.cursorMode == GLFW_CURSOR_DISABLED || x11Dpy == 0 || w.platformHandle == 0 {
		return w.cursorX, w.cursorY
	}
	var root, child uintptr
	var rx, ry, wx, wy int32
	var mask uint32
	if xQueryPointer(x11Dpy, w.platformHandle, &root, &child, &rx, &ry, &wx, &wy, &mask) != 0 {
		w.cursorX = float64(wx)
		w.cursorY = float64(wy)
	}
	return w.cursorX, w.cursorY
}

func x11CreateWindowSurface(instance uintptr, w *Window, allocator uintptr, surface *uint64) (int32, error) {
	if err := bindInstance(&VkCreateXlibSurfaceKHR, instance, "vkCreateXlibSurfaceKHR"); err != nil {
		return VK_ERROR_INITIALIZATION_FAILED, err
	}
	info := VkXlibSurfaceCreateInfoKHR{
		SType:  VK_STRUCTURE_TYPE_XLIB_SURFACE_CREATE_INFO_KHR,
		Dpy:    x11Dpy,
		Window: w.platformHandle,
	}
	res := VkCreateXlibSurfaceKHR(instance, &info, allocator, surface)
	if res != VK_SUCCESS {
		return res, fmt.Errorf("vkCreateXlibSurfaceKHR returned %d", res)
	}
	return VK_SUCCESS, nil
}

func x11Mods(state uint32) int {
	mods := 0
	if (state & x11ShiftMask) != 0 {
		mods |= GLFW_MOD_SHIFT
	}
	if (state & x11ControlMask) != 0 {
		mods |= GLFW_MOD_CONTROL
	}
	if (state & x11Mod1Mask) != 0 {
		mods |= GLFW_MOD_ALT
	}
	if (state & x11Mod4Mask) != 0 {
		mods |= GLFW_MOD_SUPER
	}
	if (state & x11LockMask) != 0 {
		mods |= GLFW_MOD_CAPS_LOCK
	}
	if (state & x11Mod2Mask) != 0 {
		mods |= GLFW_MOD_NUM_LOCK
	}
	return mods
}

func translateX11Keysym(sym uintptr) int {
	if sym >= 'a' && sym <= 'z' {
		return int(sym - 'a' + 'A')
	}
	if sym >= 'A' && sym <= 'Z' {
		return int(sym)
	}
	if sym >= '0' && sym <= '9' {
		return int(sym)
	}
	switch sym {
	case 0x0020:
		return GLFW_KEY_SPACE
	case 0xff08:
		return GLFW_KEY_BACKSPACE
	case 0xff09:
		return GLFW_KEY_TAB
	case 0xff0d:
		return GLFW_KEY_ENTER
	case 0xff1b:
		return GLFW_KEY_ESCAPE
	case 0xff50:
		return GLFW_KEY_HOME
	case 0xff51:
		return GLFW_KEY_LEFT
	case 0xff52:
		return GLFW_KEY_UP
	case 0xff53:
		return GLFW_KEY_RIGHT
	case 0xff54:
		return GLFW_KEY_DOWN
	case 0xff55:
		return GLFW_KEY_PAGE_UP
	case 0xff56:
		return GLFW_KEY_PAGE_DOWN
	case 0xff57:
		return GLFW_KEY_END
	case 0xff63:
		return GLFW_KEY_INSERT
	case 0xffff:
		return GLFW_KEY_DELETE
	case 0xffbe:
		return GLFW_KEY_F1
	case 0xffbf:
		return GLFW_KEY_F2
	case 0xffc0:
		return GLFW_KEY_F3
	case 0xffc1:
		return GLFW_KEY_F4
	case 0xffc2:
		return GLFW_KEY_F5
	case 0xffc3:
		return GLFW_KEY_F6
	case 0xffc4:
		return GLFW_KEY_F7
	case 0xffc5:
		return GLFW_KEY_F8
	case 0xffc6:
		return GLFW_KEY_F9
	case 0xffc7:
		return GLFW_KEY_F10
	case 0xffc8:
		return GLFW_KEY_F11
	case 0xffc9:
		return GLFW_KEY_F12
	case 0xffe1:
		return GLFW_KEY_LEFT_SHIFT
	case 0xffe2:
		return GLFW_KEY_RIGHT_SHIFT
	case 0xffe3:
		return GLFW_KEY_LEFT_CONTROL
	case 0xffe4:
		return GLFW_KEY_RIGHT_CONTROL
	case 0xffe9:
		return GLFW_KEY_LEFT_ALT
	case 0xffea:
		return GLFW_KEY_RIGHT_ALT
	case 0xffeb:
		return GLFW_KEY_LEFT_SUPER
	case 0xffec:
		return GLFW_KEY_RIGHT_SUPER
	default:
		return GLFW_KEY_UNKNOWN
	}
}

func parseURIList(raw string) []string {
	var paths []string
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(strings.TrimSuffix(line, "\r"))
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "file://") {
			u, err := url.Parse(line)
			if err == nil && u.Path != "" {
				if unescaped, err := url.PathUnescape(u.Path); err == nil {
					paths = append(paths, unescaped)
					continue
				}
				paths = append(paths, u.Path)
				continue
			}
		}
		paths = append(paths, line)
	}
	return paths
}

func x11ProcessEvent(ev *xEvent) {
	evType := *(*int32)(unsafe.Pointer(ev))
	switch evType {
	case x11KeyPress, x11KeyRelease:
		ke := (*xKeyEvent)(unsafe.Pointer(ev))
		w := lookupWindow(ke.Window)
		if w == nil {
			return
		}
		sym := xLookupKeysym(ke, 0)
		key := translateX11Keysym(sym)
		action := GLFW_PRESS
		if evType == x11KeyRelease {
			action = GLFW_RELEASE
		}
		w.dispatchKey(key, int(ke.Keycode), action, x11Mods(ke.State))
		if evType == x11KeyPress {
			var buf [32]byte
			var ksym uintptr
			n := xLookupString(ke, &buf[0], int32(len(buf)), &ksym, 0)
			if n > 0 {
				for _, r := range string(buf[:n]) {
					w.dispatchChar(r)
				}
			}
		}
	case x11ButtonPress, x11ButtonRelease:
		be := (*xButtonEvent)(unsafe.Pointer(ev))
		w := lookupWindow(be.Window)
		if w == nil {
			return
		}
		mods := x11Mods(be.State)
		action := GLFW_PRESS
		if evType == x11ButtonRelease {
			action = GLFW_RELEASE
		}
		switch be.Button {
		case 1:
			w.dispatchMouseButton(GLFW_MOUSE_BUTTON_LEFT, action, mods)
		case 2:
			w.dispatchMouseButton(GLFW_MOUSE_BUTTON_MIDDLE, action, mods)
		case 3:
			w.dispatchMouseButton(GLFW_MOUSE_BUTTON_RIGHT, action, mods)
		case 4:
			if evType == x11ButtonPress {
				w.dispatchScroll(0, 1.0)
			}
		case 5:
			if evType == x11ButtonPress {
				w.dispatchScroll(0, -1.0)
			}
		case 6:
			if evType == x11ButtonPress {
				w.dispatchScroll(1.0, 0)
			}
		case 7:
			if evType == x11ButtonPress {
				w.dispatchScroll(-1.0, 0)
			}
		}
	case x11MotionNotify:
		me := (*xMotionEvent)(unsafe.Pointer(ev))
		w := lookupWindow(me.Window)
		if w == nil {
			return
		}
		if w.cursorMode == GLFW_CURSOR_DISABLED {
			cx := int32(w.width / 2)
			cy := int32(w.height / 2)
			if me.X == cx && me.Y == cy {
				return
			}
			dx := float64(me.X) - float64(cx)
			dy := float64(me.Y) - float64(cy)
			w.dispatchRelativeCursorDelta(dx, dy)
			xWarpPointer(x11Dpy, 0, w.platformHandle, 0, 0, 0, 0, cx, cy)
			xFlush(x11Dpy)
		} else {
			w.dispatchCursorPos(float64(me.X), float64(me.Y))
		}
	case x11ConfigureNotify:
		ce := (*xConfigureEvent)(unsafe.Pointer(ev))
		w := lookupWindow(ce.Window)
		if w == nil {
			return
		}
		w.xpos = int(ce.X)
		w.ypos = int(ce.Y)
		if int(ce.Width) != w.width || int(ce.Height) != w.height {
			w.dispatchWindowSize(int(ce.Width), int(ce.Height))
			w.dispatchFramebufferSize(int(ce.Width), int(ce.Height))
		}
	case x11FocusIn:
		fe := (*xKeyEvent)(unsafe.Pointer(ev))
		if w := lookupWindow(fe.Window); w != nil {
			w.dispatchFocus(true)
		}
	case x11FocusOut:
		fe := (*xKeyEvent)(unsafe.Pointer(ev))
		if w := lookupWindow(fe.Window); w != nil {
			w.dispatchFocus(false)
		}
	case x11ClientMessage:
		cm := (*xClientMessageEvent)(unsafe.Pointer(ev))
		w := lookupWindow(cm.Window)
		if w == nil {
			return
		}
		if cm.MessageType == x11AtomWMProtocols && cm.Data[0] == x11AtomWMDeleteWin {
			w.shouldClose = true
		} else if cm.MessageType == x11AtomXdndEnter {
			x11XdndSource = cm.Data[0]
		} else if cm.MessageType == x11AtomXdndPos {
			x11XdndSource = cm.Data[0]
			var reply xClientMessageEvent
			reply.Type = x11ClientMessage
			reply.Window = x11XdndSource
			reply.MessageType = x11AtomXdndStatus
			reply.Format = 32
			reply.Data[0] = w.platformHandle
			reply.Data[1] = 1
			reply.Data[4] = x11AtomXdndActionCopy
			xSendEvent(x11Dpy, x11XdndSource, 0, 0, (*xEvent)(unsafe.Pointer(&reply)))
			xFlush(x11Dpy)
		} else if cm.MessageType == x11AtomXdndDrop {
			x11XdndSource = cm.Data[0]
			xConvertSelection(x11Dpy, x11AtomXdndSel, x11AtomUriList, x11AtomXdndProp, w.platformHandle, cm.Data[2])
			xFlush(x11Dpy)
		}
	case x11SelectionNotify:
		se := (*xSelectionEvent)(unsafe.Pointer(ev))
		w := lookupWindow(se.Requestor)
		if w == nil || se.Property == 0 {
			return
		}
		var actualType uintptr
		var actualFormat int32
		var nitems, bytesAfter uintptr
		var prop *byte
		if xGetWindowProperty(x11Dpy, w.platformHandle, se.Property, 0, 65536, 1, x11AnyPropertyType, &actualType, &actualFormat, &nitems, &bytesAfter, &prop) == 0 && prop != nil {
			raw := unsafe.Slice(prop, int(nitems))
			paths := parseURIList(string(raw))
			xFree(uintptr(unsafe.Pointer(prop)))
			if len(paths) > 0 {
				w.dispatchDrop(paths)
			}
		}
		if x11XdndSource != 0 {
			var fin xClientMessageEvent
			fin.Type = x11ClientMessage
			fin.Window = x11XdndSource
			fin.MessageType = x11AtomXdndFin
			fin.Format = 32
			fin.Data[0] = w.platformHandle
			fin.Data[1] = 1
			fin.Data[2] = x11AtomXdndActionCopy
			xSendEvent(x11Dpy, x11XdndSource, 0, 0, (*xEvent)(unsafe.Pointer(&fin)))
			xFlush(x11Dpy)
		}
	}
}
