//go:build darwin

package binding

import (
	"fmt"
	"runtime"

	"github.com/ebitengine/purego"
)

const (
	nsApplicationActivationPolicyRegular = 0

	nsWindowStyleMaskBorderless     = 0
	nsWindowStyleMaskTitled         = 1 << 0
	nsWindowStyleMaskClosable       = 1 << 1
	nsWindowStyleMaskMiniaturizable = 1 << 2
	nsWindowStyleMaskResizable      = 1 << 3
	nsWindowStyleMaskFullScreen     = 1 << 14

	nsBackingStoreBuffered = 2

	nsEventMaskAny = ^uint64(0)

	nsEventTypeLeftMouseDown      = 1
	nsEventTypeLeftMouseUp        = 2
	nsEventTypeRightMouseDown     = 3
	nsEventTypeRightMouseUp       = 4
	nsEventTypeMouseMoved         = 5
	nsEventTypeLeftMouseDragged   = 6
	nsEventTypeRightMouseDragged  = 7
	nsEventTypeKeyDown            = 10
	nsEventTypeKeyUp              = 11
	nsEventTypeFlagsChanged       = 12
	nsEventTypeScrollWheel        = 22
	nsEventTypeOtherMouseDown     = 25
	nsEventTypeOtherMouseUp       = 26
	nsEventTypeOtherMouseDragged  = 27

	nsEventModifierFlagCapsLock = 1 << 16
	nsEventModifierFlagShift    = 1 << 17
	nsEventModifierFlagControl  = 1 << 18
	nsEventModifierFlagOption   = 1 << 19
	nsEventModifierFlagCommand  = 1 << 20
)

type cgPoint struct {
	X float64
	Y float64
}

type cgSize struct {
	Width  float64
	Height float64
}

type cgRect struct {
	Origin cgPoint
	Size   cgSize
}

var (
	objcLib       uintptr
	appKitLib     uintptr
	quartzCoreLib uintptr
	coreGraphLib  uintptr
	foundationLib uintptr

	objcGetClass    func(name *byte) uintptr
	selRegisterName func(name *byte) uintptr
	objcMsgSendPtr  uintptr

	msgSend0        func(id, sel uintptr) uintptr
	msgSend1        func(id, sel, a1 uintptr) uintptr
	msgSend2        func(id, sel, a1, a2 uintptr) uintptr
	msgSendBool     func(id, sel uintptr) bool
	msgSendSetBool  func(id, sel uintptr, v bool)
	msgSendFloat    func(id, sel uintptr) float64
	msgSendSetFloat func(id, sel uintptr, v float64)
	msgSendPoint    func(id, sel uintptr) cgPoint
	msgSendSetSize  func(id, sel uintptr, size cgSize)
	msgSendSetFrame func(id, sel uintptr, rect cgRect, display bool)
	msgSendInitWin  func(id, sel uintptr, rect cgRect, style uintptr, backing uintptr, deferFlag bool) uintptr
	msgSendNextEv   func(id, sel uintptr, mask uint64, untilDate, mode uintptr, dequeue bool) uintptr
	msgSendDate     func(id, sel uintptr, secs float64) uintptr
	msgSendUTF8     func(id, sel uintptr, str *byte) uintptr
	msgSendCStr     func(id, sel uintptr) *byte

	cgMainDisplayID                        func() uint32
	cgDisplayPixelsWide                    func(display uint32) uintptr
	cgDisplayPixelsHigh                    func(display uint32) uintptr
	cgDisplayCopyDisplayMode               func(display uint32) uintptr
	cgDisplayCopyAllDisplayModes           func(display uint32, options uintptr) uintptr
	cgDisplayModeGetWidth                  func(mode uintptr) uintptr
	cgDisplayModeGetHeight                 func(mode uintptr) uintptr
	cgDisplayModeGetRefreshRate            func(mode uintptr) float64
	cgDisplayModeRelease                   func(mode uintptr)
	cgDisplayHideCursor                    func(display uint32) int32
	cgDisplayShowCursor                    func(display uint32) int32
	cgAssociateMouseAndMouseCursorPosition func(connected int32) int32
	cfArrayGetCount                        func(arr uintptr) int
	cfArrayGetValueAtIndex                 func(arr uintptr, idx int) uintptr
	cfRelease                              func(cf uintptr)

	nsApp                uintptr
	nsDefaultRunLoopMode uintptr

	clsNSApplication uintptr
	clsNSWindow      uintptr
	clsNSString      uintptr
	clsNSDate        uintptr
	clsNSEvent       uintptr
	clsCAMetalLayer  uintptr

	selSharedApplication            uintptr
	selSetActivationPolicy          uintptr
	selActivateIgnoringOtherApps    uintptr
	selFinishLaunching              uintptr
	selAlloc                        uintptr
	selInit                         uintptr
	selRelease                      uintptr
	selInitWithContentRect          uintptr
	selSetTitle                     uintptr
	selMakeKeyAndOrderFront         uintptr
	selClose                        uintptr
	selContentView                  uintptr
	selSetWantsLayer                uintptr
	selSetLayer                     uintptr
	selBackingScaleFactor           uintptr
	selSetContentsScale             uintptr
	selSetDrawableSize              uintptr
	selIsKeyWindow                  uintptr
	selIsVisible                    uintptr
	selToggleFullScreen             uintptr
	selSetFrameDisplay              uintptr
	selNextEventMatchingMask        uintptr
	selSendEvent                    uintptr
	selUpdateWindows                uintptr
	selType                         uintptr
	selWindow                       uintptr
	selLocationInWindow             uintptr
	selDeltaX                       uintptr
	selDeltaY                       uintptr
	selScrollingDeltaX              uintptr
	selScrollingDeltaY              uintptr
	selKeyCode                      uintptr
	selModifierFlags                uintptr
	selCharacters                   uintptr
	selButtonNumber                 uintptr
	selStringWithUTF8String         uintptr
	selUTF8String                   uintptr
	selDateWithTimeIntervalSinceNow uintptr
	selDistantPast                  uintptr
)

func objcClass(name string) uintptr {
	b := CString(name)
	id := objcGetClass(&b[0])
	runtime.KeepAlive(b)
	return id
}

func objcSel(name string) uintptr {
	b := CString(name)
	s := selRegisterName(&b[0])
	runtime.KeepAlive(b)
	return s
}

func nsString(s string) uintptr {
	b := CString(s)
	res := msgSendUTF8(clsNSString, selStringWithUTF8String, &b[0])
	runtime.KeepAlive(b)
	return res
}

func platformInit() error {
	runtime.LockOSThread()
	if objcLib != 0 {
		return nil
	}
	var err error
	objcLib, err = openLibrary("/usr/lib/libobjc.A.dylib")
	if err != nil {
		return fmt.Errorf("load libobjc: %w", err)
	}
	appKitLib, err = openLibrary("/System/Library/Frameworks/AppKit.framework/AppKit")
	if err != nil {
		return fmt.Errorf("load AppKit: %w", err)
	}
	quartzCoreLib, err = openLibrary("/System/Library/Frameworks/QuartzCore.framework/QuartzCore")
	if err != nil {
		return fmt.Errorf("load QuartzCore: %w", err)
	}
	coreGraphLib, err = openLibrary("/System/Library/Frameworks/CoreGraphics.framework/CoreGraphics")
	if err != nil {
		return fmt.Errorf("load CoreGraphics: %w", err)
	}
	foundationLib, _ = openLibrary("/System/Library/Frameworks/Foundation.framework/Foundation")

	if err = bindLibrary(&objcGetClass, objcLib, "objc_getClass"); err != nil {
		return err
	}
	if err = bindLibrary(&selRegisterName, objcLib, "sel_registerName"); err != nil {
		return err
	}
	objcMsgSendPtr, err = lookupLibrary(objcLib, "objc_msgSend")
	if err != nil {
		return err
	}

	purego.RegisterFunc(&msgSend0, objcMsgSendPtr)
	purego.RegisterFunc(&msgSend1, objcMsgSendPtr)
	purego.RegisterFunc(&msgSend2, objcMsgSendPtr)
	purego.RegisterFunc(&msgSendBool, objcMsgSendPtr)
	purego.RegisterFunc(&msgSendSetBool, objcMsgSendPtr)
	purego.RegisterFunc(&msgSendFloat, objcMsgSendPtr)
	purego.RegisterFunc(&msgSendSetFloat, objcMsgSendPtr)
	purego.RegisterFunc(&msgSendPoint, objcMsgSendPtr)
	purego.RegisterFunc(&msgSendSetSize, objcMsgSendPtr)
	purego.RegisterFunc(&msgSendSetFrame, objcMsgSendPtr)
	purego.RegisterFunc(&msgSendInitWin, objcMsgSendPtr)
	purego.RegisterFunc(&msgSendNextEv, objcMsgSendPtr)
	purego.RegisterFunc(&msgSendDate, objcMsgSendPtr)
	purego.RegisterFunc(&msgSendUTF8, objcMsgSendPtr)
	purego.RegisterFunc(&msgSendCStr, objcMsgSendPtr)

	_ = bindLibrary(&cgMainDisplayID, coreGraphLib, "CGMainDisplayID")
	_ = bindLibrary(&cgDisplayPixelsWide, coreGraphLib, "CGDisplayPixelsWide")
	_ = bindLibrary(&cgDisplayPixelsHigh, coreGraphLib, "CGDisplayPixelsHigh")
	_ = bindLibrary(&cgDisplayCopyDisplayMode, coreGraphLib, "CGDisplayCopyDisplayMode")
	_ = bindLibrary(&cgDisplayCopyAllDisplayModes, coreGraphLib, "CGDisplayCopyAllDisplayModes")
	_ = bindLibrary(&cgDisplayModeGetWidth, coreGraphLib, "CGDisplayModeGetWidth")
	_ = bindLibrary(&cgDisplayModeGetHeight, coreGraphLib, "CGDisplayModeGetHeight")
	_ = bindLibrary(&cgDisplayModeGetRefreshRate, coreGraphLib, "CGDisplayModeGetRefreshRate")
	_ = bindLibrary(&cgDisplayModeRelease, coreGraphLib, "CGDisplayModeRelease")
	_ = bindLibrary(&cgDisplayHideCursor, coreGraphLib, "CGDisplayHideCursor")
	_ = bindLibrary(&cgDisplayShowCursor, coreGraphLib, "CGDisplayShowCursor")
	_ = bindLibrary(&cgAssociateMouseAndMouseCursorPosition, coreGraphLib, "CGAssociateMouseAndMouseCursorPosition")
	_ = bindLibrary(&cfArrayGetCount, coreGraphLib, "CFArrayGetCount")
	_ = bindLibrary(&cfArrayGetValueAtIndex, coreGraphLib, "CFArrayGetValueAtIndex")
	_ = bindLibrary(&cfRelease, coreGraphLib, "CFRelease")

	clsNSApplication = objcClass("NSApplication")
	clsNSWindow = objcClass("NSWindow")
	clsNSString = objcClass("NSString")
	clsNSDate = objcClass("NSDate")
	clsNSEvent = objcClass("NSEvent")
	clsCAMetalLayer = objcClass("CAMetalLayer")

	selSharedApplication = objcSel("sharedApplication")
	selSetActivationPolicy = objcSel("setActivationPolicy:")
	selActivateIgnoringOtherApps = objcSel("activateIgnoringOtherApps:")
	selFinishLaunching = objcSel("finishLaunching")
	selAlloc = objcSel("alloc")
	selInit = objcSel("init")
	selRelease = objcSel("release")
	selInitWithContentRect = objcSel("initWithContentRect:styleMask:backing:defer:")
	selSetTitle = objcSel("setTitle:")
	selMakeKeyAndOrderFront = objcSel("makeKeyAndOrderFront:")
	selClose = objcSel("close")
	selContentView = objcSel("contentView")
	selSetWantsLayer = objcSel("setWantsLayer:")
	selSetLayer = objcSel("setLayer:")
	selBackingScaleFactor = objcSel("backingScaleFactor")
	selSetContentsScale = objcSel("setContentsScale:")
	selSetDrawableSize = objcSel("setDrawableSize:")
	selIsKeyWindow = objcSel("isKeyWindow")
	selIsVisible = objcSel("isVisible")
	selToggleFullScreen = objcSel("toggleFullScreen:")
	selSetFrameDisplay = objcSel("setFrame:display:")
	selNextEventMatchingMask = objcSel("nextEventMatchingMask:untilDate:inMode:dequeue:")
	selSendEvent = objcSel("sendEvent:")
	selUpdateWindows = objcSel("updateWindows")
	selType = objcSel("type")
	selWindow = objcSel("window")
	selLocationInWindow = objcSel("locationInWindow")
	selDeltaX = objcSel("deltaX")
	selDeltaY = objcSel("deltaY")
	selScrollingDeltaX = objcSel("scrollingDeltaX")
	selScrollingDeltaY = objcSel("scrollingDeltaY")
	selKeyCode = objcSel("keyCode")
	selModifierFlags = objcSel("modifierFlags")
	selCharacters = objcSel("characters")
	selButtonNumber = objcSel("buttonNumber")
	selStringWithUTF8String = objcSel("stringWithUTF8String:")
	selUTF8String = objcSel("UTF8String")
	selDateWithTimeIntervalSinceNow = objcSel("dateWithTimeIntervalSinceNow:")
	selDistantPast = objcSel("distantPast")

	nsDefaultRunLoopMode = nsString("kCFRunLoopDefaultMode")
	nsApp = msgSend0(clsNSApplication, selSharedApplication)
	if nsApp != 0 {
		msgSend1(nsApp, selSetActivationPolicy, nsApplicationActivationPolicyRegular)
		msgSendSetBool(nsApp, selActivateIgnoringOtherApps, true)
		msgSend0(nsApp, selFinishLaunching)
	}
	return nil
}

func platformTerminate() {
	if cgAssociateMouseAndMouseCursorPosition != nil {
		cgAssociateMouseAndMouseCursorPosition(1)
	}
	if cgDisplayShowCursor != nil && cgMainDisplayID != nil {
		cgDisplayShowCursor(cgMainDisplayID())
	}
}

func platformGetPrimaryMonitor() uintptr {
	if cgMainDisplayID != nil {
		id := cgMainDisplayID()
		if id != 0 {
			return uintptr(id)
		}
	}
	return 1
}

func platformGetVideoMode(monitor uintptr) GLFWvidmode {
	if cgMainDisplayID == nil {
		return GLFWvidmode{Width: 1920, Height: 1080, RedBits: 8, GreenBits: 8, BlueBits: 8, RefreshRate: 60}
	}
	disp := uint32(monitor)
	if disp == 0 || disp == 1 {
		disp = cgMainDisplayID()
	}
	if cgDisplayCopyDisplayMode != nil {
		mode := cgDisplayCopyDisplayMode(disp)
		if mode != 0 {
			w := int32(cgDisplayModeGetWidth(mode))
			h := int32(cgDisplayModeGetHeight(mode))
			r := int32(cgDisplayModeGetRefreshRate(mode))
			if r <= 0 {
				r = 60
			}
			cgDisplayModeRelease(mode)
			if w > 0 && h > 0 {
				return GLFWvidmode{Width: w, Height: h, RedBits: 8, GreenBits: 8, BlueBits: 8, RefreshRate: r}
			}
		}
	}
	w := int32(cgDisplayPixelsWide(disp))
	h := int32(cgDisplayPixelsHigh(disp))
	if w <= 0 || h <= 0 {
		w, h = 1920, 1080
	}
	return GLFWvidmode{Width: w, Height: h, RedBits: 8, GreenBits: 8, BlueBits: 8, RefreshRate: 60}
}

func platformGetVideoModes(monitor uintptr) []GLFWvidmode {
	if cgMainDisplayID == nil || cgDisplayCopyAllDisplayModes == nil {
		return []GLFWvidmode{platformGetVideoMode(monitor)}
	}
	disp := uint32(monitor)
	if disp == 0 || disp == 1 {
		disp = cgMainDisplayID()
	}
	arr := cgDisplayCopyAllDisplayModes(disp, 0)
	if arr == 0 {
		return []GLFWvidmode{platformGetVideoMode(monitor)}
	}
	defer cfRelease(arr)
	n := cfArrayGetCount(arr)
	var out []GLFWvidmode
	seen := make(map[[3]int32]bool)
	for i := 0; i < n; i++ {
		m := cfArrayGetValueAtIndex(arr, i)
		if m == 0 {
			continue
		}
		w := int32(cgDisplayModeGetWidth(m))
		h := int32(cgDisplayModeGetHeight(m))
		r := int32(cgDisplayModeGetRefreshRate(m))
		if r <= 0 {
			r = 60
		}
		if w <= 0 || h <= 0 {
			continue
		}
		k := [3]int32{w, h, r}
		if !seen[k] {
			seen[k] = true
			out = append(out, GLFWvidmode{Width: w, Height: h, RedBits: 8, GreenBits: 8, BlueBits: 8, RefreshRate: r})
		}
	}
	if len(out) == 0 {
		out = append(out, platformGetVideoMode(monitor))
	}
	return out
}

func platformCreateWindow(w *Window, width, height int, title string, monitor uintptr) error {
	style := uintptr(nsWindowStyleMaskTitled | nsWindowStyleMaskClosable | nsWindowStyleMaskMiniaturizable | nsWindowStyleMaskResizable)
	if monitor != 0 {
		style = nsWindowStyleMaskBorderless
	}
	rect := cgRect{
		Origin: cgPoint{X: 100, Y: 100},
		Size:   cgSize{Width: float64(width), Height: float64(height)},
	}
	rawWin := msgSend0(clsNSWindow, selAlloc)
	nsWin := msgSendInitWin(rawWin, selInitWithContentRect, rect, style, nsBackingStoreBuffered, false)
	if nsWin == 0 {
		return fmt.Errorf("failed to create NSWindow")
	}
	msgSend1(nsWin, selSetTitle, nsString(title))
	view := msgSend0(nsWin, selContentView)
	layer := msgSend0(msgSend0(clsCAMetalLayer, selAlloc), selInit)
	scale := msgSendFloat(nsWin, selBackingScaleFactor)
	if scale <= 0 {
		scale = 1.0
	}
	fbW := int(float64(width) * scale)
	fbH := int(float64(height) * scale)
	msgSendSetFloat(layer, selSetContentsScale, scale)
	msgSendSetSize(layer, selSetDrawableSize, cgSize{Width: float64(fbW), Height: float64(fbH)})
	msgSendSetBool(view, selSetWantsLayer, true)
	msgSend1(view, selSetLayer, layer)
	msgSend1(nsWin, selMakeKeyAndOrderFront, 0)
	if nsApp != 0 {
		msgSendSetBool(nsApp, selActivateIgnoringOtherApps, true)
	}

	w.platformHandle = nsWin
	w.platformSurface = view
	w.platformLayer = layer
	w.handle = nsWin
	w.xpos = 100
	w.ypos = 100
	w.fbWidth = fbW
	w.fbHeight = fbH
	w.focused = true
	if msgSendBool(nsWin, selIsKeyWindow) {
		w.platformExtra = 1
	}
	registerWindow(w)
	return nil
}

func platformDestroyWindow(w *Window) {
	if w != nil && w.platformHandle != 0 {
		msgSend0(w.platformHandle, selClose)
		w.platformHandle = 0
		w.platformLayer = 0
		w.platformSurface = 0
	}
}

func platformPollEvents() {
	if nsApp == 0 {
		return
	}
	distantPast := msgSend0(clsNSDate, selDistantPast)
	for {
		ev := msgSendNextEv(nsApp, selNextEventMatchingMask, nsEventMaskAny, distantPast, nsDefaultRunLoopMode, true)
		if ev == 0 {
			break
		}
		handleCocoaEvent(ev)
		msgSend1(nsApp, selSendEvent, ev)
	}
	msgSend0(nsApp, selUpdateWindows)

	// Detect real OS key-window transitions without overwriting manual focus state
	glfwMu.RLock()
	wins := make([]*Window, 0, len(glfwWindows))
	for _, w := range glfwWindows {
		wins = append(wins, w)
	}
	glfwMu.RUnlock()
	for _, w := range wins {
		if w.platformHandle == 0 {
			continue
		}
		isKey := msgSendBool(w.platformHandle, selIsKeyWindow)
		var keyFlag uintptr
		if isKey {
			keyFlag = 1
		}
		if keyFlag != w.platformExtra {
			w.platformExtra = keyFlag
			w.dispatchFocus(isKey)
		}
	}
}

func platformWaitEventsTimeout(timeout float64) {
	if nsApp == 0 {
		return
	}
	untilDate := msgSendDate(clsNSDate, selDateWithTimeIntervalSinceNow, timeout)
	ev := msgSendNextEv(nsApp, selNextEventMatchingMask, nsEventMaskAny, untilDate, nsDefaultRunLoopMode, true)
	if ev != 0 {
		handleCocoaEvent(ev)
		msgSend1(nsApp, selSendEvent, ev)
	}
	platformPollEvents()
}

func platformSetWindowMonitor(w *Window, monitor uintptr, xpos, ypos, width, height, refreshRate int) {
	if w == nil {
		return
	}
	w.xpos = xpos
	w.ypos = ypos
	if width > 0 && height > 0 {
		w.dispatchWindowSize(width, height)
	}
	scale := 1.0
	if w.platformHandle != 0 {
		if width > 0 && height > 0 {
			rect := cgRect{
				Origin: cgPoint{X: float64(xpos), Y: float64(ypos)},
				Size:   cgSize{Width: float64(width), Height: float64(height)},
			}
			msgSendSetFrame(w.platformHandle, selSetFrameDisplay, rect, true)
		}
		if s := msgSendFloat(w.platformHandle, selBackingScaleFactor); s > 0 {
			scale = s
		}
	}
	if width > 0 && height > 0 {
		fbW := int(float64(width) * scale)
		fbH := int(float64(height) * scale)
		if w.platformLayer != 0 {
			msgSendSetSize(w.platformLayer, selSetDrawableSize, cgSize{Width: float64(fbW), Height: float64(fbH)})
		}
		w.dispatchFramebufferSize(fbW, fbH)
	}
}

func platformGetWindowSize(w *Window) (int, int) {
	if w == nil {
		return 0, 0
	}
	return w.width, w.height
}

func platformGetFramebufferSize(w *Window) (int, int) {
	if w == nil {
		return 0, 0
	}
	if w.platformHandle != 0 && w.fbWidth == 0 && w.fbHeight == 0 {
		if scale := msgSendFloat(w.platformHandle, selBackingScaleFactor); scale > 0 {
			w.fbWidth = int(float64(w.width) * scale)
			w.fbHeight = int(float64(w.height) * scale)
		}
	}
	return w.fbWidth, w.fbHeight
}

func platformGetWindowPos(w *Window) (int, int) {
	if w == nil {
		return 0, 0
	}
	return w.xpos, w.ypos
}

func platformWindowFocused(w *Window) bool {
	if w == nil {
		return false
	}
	return w.focused
}

func platformSetCursorMode(w *Window, mode int) {
	if cgMainDisplayID == nil {
		return
	}
	disp := cgMainDisplayID()
	if mode == GLFW_CURSOR_DISABLED {
		cgAssociateMouseAndMouseCursorPosition(0)
		cgDisplayHideCursor(disp)
	} else {
		cgAssociateMouseAndMouseCursorPosition(1)
		cgDisplayShowCursor(disp)
	}
}

func platformGetCursorPos(w *Window) (float64, float64) {
	if w == nil {
		return 0, 0
	}
	return w.cursorX, w.cursorY
}

func platformRequiredInstanceExtensions() []string {
	return []string{"VK_KHR_surface", "VK_EXT_metal_surface"}
}

func platformCreateWindowSurface(instance uintptr, w *Window, allocator uintptr, surface *uint64) (int32, error) {
	if err := bindInstance(&VkCreateMetalSurfaceEXT, instance, "vkCreateMetalSurfaceEXT"); err == nil {
		info := VkMetalSurfaceCreateInfoEXT{
			SType:  VK_STRUCTURE_TYPE_METAL_SURFACE_CREATE_INFO_EXT,
			PLayer: w.platformLayer,
		}
		res := VkCreateMetalSurfaceEXT(instance, &info, allocator, surface)
		if res != VK_SUCCESS {
			return res, fmt.Errorf("vkCreateMetalSurfaceEXT returned %d", res)
		}
		return VK_SUCCESS, nil
	}
	if err := bindInstance(&VkCreateMacOSSurfaceMVK, instance, "vkCreateMacOSSurfaceMVK"); err != nil {
		return VK_ERROR_INITIALIZATION_FAILED, err
	}
	info := VkMacOSSurfaceCreateInfoMVK{
		SType: VK_STRUCTURE_TYPE_MACOS_SURFACE_CREATE_INFO_MVK,
		PView: w.platformSurface,
	}
	res := VkCreateMacOSSurfaceMVK(instance, &info, allocator, surface)
	if res != VK_SUCCESS {
		return res, fmt.Errorf("vkCreateMacOSSurfaceMVK returned %d", res)
	}
	return VK_SUCCESS, nil
}

func cocoaMods(flags uintptr) int {
	mods := 0
	if (flags & nsEventModifierFlagShift) != 0 {
		mods |= GLFW_MOD_SHIFT
	}
	if (flags & nsEventModifierFlagControl) != 0 {
		mods |= GLFW_MOD_CONTROL
	}
	if (flags & nsEventModifierFlagOption) != 0 {
		mods |= GLFW_MOD_ALT
	}
	if (flags & nsEventModifierFlagCommand) != 0 {
		mods |= GLFW_MOD_SUPER
	}
	if (flags & nsEventModifierFlagCapsLock) != 0 {
		mods |= GLFW_MOD_CAPS_LOCK
	}
	return mods
}

func translateMacKey(keyCode uint16) int {
	switch keyCode {
	case 0x00:
		return GLFW_KEY_A
	case 0x01:
		return GLFW_KEY_S
	case 0x02:
		return GLFW_KEY_D
	case 0x03:
		return GLFW_KEY_F
	case 0x06:
		return GLFW_KEY_Z
	case 0x0C:
		return GLFW_KEY_Q
	case 0x0D:
		return GLFW_KEY_W
	case 0x0E:
		return GLFW_KEY_E
	case 0x0F:
		return GLFW_KEY_R
	case 0x10:
		return GLFW_KEY_Y
	case 0x24:
		return GLFW_KEY_ENTER
	case 0x30:
		return GLFW_KEY_TAB
	case 0x31:
		return GLFW_KEY_SPACE
	case 0x33:
		return GLFW_KEY_BACKSPACE
	case 0x35:
		return GLFW_KEY_ESCAPE
	case 0x3B:
		return GLFW_KEY_LEFT_CONTROL
	case 0x60:
		return GLFW_KEY_F5
	case 0x61:
		return GLFW_KEY_F6
	case 0x64:
		return GLFW_KEY_F8
	case 0x73:
		return GLFW_KEY_HOME
	case 0x75:
		return GLFW_KEY_DELETE
	case 0x77:
		return GLFW_KEY_END
	case 0x7B:
		return GLFW_KEY_LEFT
	case 0x7C:
		return GLFW_KEY_RIGHT
	case 0x7D:
		return GLFW_KEY_DOWN
	case 0x7E:
		return GLFW_KEY_UP
	default:
		return GLFW_KEY_UNKNOWN
	}
}

func handleCocoaEvent(ev uintptr) {
	if ev == 0 {
		return
	}
	nsWin := msgSend0(ev, selWindow)
	w := lookupWindow(nsWin)
	if w == nil {
		return
	}
	evType := msgSend0(ev, selType)
	mods := cocoaMods(msgSend0(ev, selModifierFlags))
	switch evType {
	case nsEventTypeKeyDown:
		kc := uint16(msgSend0(ev, selKeyCode))
		w.dispatchKey(translateMacKey(kc), int(kc), GLFW_PRESS, mods)
		charsObj := msgSend0(ev, selCharacters)
		if charsObj != 0 {
			str := String(msgSendCStr(charsObj, selUTF8String))
			for _, r := range str {
				w.dispatchChar(r)
			}
		}
	case nsEventTypeKeyUp:
		kc := uint16(msgSend0(ev, selKeyCode))
		w.dispatchKey(translateMacKey(kc), int(kc), GLFW_RELEASE, mods)
	case nsEventTypeMouseMoved, nsEventTypeLeftMouseDragged, nsEventTypeRightMouseDragged, nsEventTypeOtherMouseDragged:
		if w.cursorMode == GLFW_CURSOR_DISABLED {
			dx := msgSendFloat(ev, selDeltaX)
			dy := msgSendFloat(ev, selDeltaY)
			if dx != 0 || dy != 0 {
				w.dispatchRelativeCursorDelta(dx, dy)
			}
		} else {
			pt := msgSendPoint(ev, selLocationInWindow)
			w.dispatchCursorPos(pt.X, float64(w.height)-pt.Y)
		}
	case nsEventTypeLeftMouseDown:
		w.dispatchMouseButton(GLFW_MOUSE_BUTTON_LEFT, GLFW_PRESS, mods)
	case nsEventTypeLeftMouseUp:
		w.dispatchMouseButton(GLFW_MOUSE_BUTTON_LEFT, GLFW_RELEASE, mods)
	case nsEventTypeRightMouseDown:
		w.dispatchMouseButton(GLFW_MOUSE_BUTTON_RIGHT, GLFW_PRESS, mods)
	case nsEventTypeRightMouseUp:
		w.dispatchMouseButton(GLFW_MOUSE_BUTTON_RIGHT, GLFW_RELEASE, mods)
	case nsEventTypeOtherMouseDown:
		w.dispatchMouseButton(GLFW_MOUSE_BUTTON_MIDDLE, GLFW_PRESS, mods)
	case nsEventTypeOtherMouseUp:
		w.dispatchMouseButton(GLFW_MOUSE_BUTTON_MIDDLE, GLFW_RELEASE, mods)
	case nsEventTypeScrollWheel:
		dx := msgSendFloat(ev, selScrollingDeltaX)
		dy := msgSendFloat(ev, selScrollingDeltaY)
		w.dispatchScroll(dx, dy)
	}
}
