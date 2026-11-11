//go:build darwin

package binding

import (
	"fmt"
	"runtime"

	"github.com/ebitengine/purego"
)

const (
	nsApplicationActivationPolicyRegular = 0

	nsApplicationPresentationDefault                    = 0
	nsApplicationPresentationAutoHideDock               = 1 << 0
	nsApplicationPresentationHideDock                   = 1 << 1
	nsApplicationPresentationAutoHideMenuBar            = 1 << 2
	nsApplicationPresentationHideMenuBar                = 1 << 3

	nsWindowStyleMaskBorderless     = 0
	nsWindowStyleMaskTitled         = 1 << 0
	nsWindowStyleMaskClosable       = 1 << 1
	nsWindowStyleMaskMiniaturizable = 1 << 2
	nsWindowStyleMaskResizable      = 1 << 3
	nsWindowStyleMaskFullScreen     = 1 << 14

	nsWindowCollectionBehaviorManaged           = 1 << 2
	nsWindowCollectionBehaviorFullScreenPrimary = 1 << 7

	nsNormalWindowLevel   = 0
	nsMainMenuWindowLevel = 24

	nsBackingStoreBuffered = 2

	nsEventMaskAny = ^uint64(0)

	nsEventTypeLeftMouseDown     = 1
	nsEventTypeLeftMouseUp       = 2
	nsEventTypeRightMouseDown    = 3
	nsEventTypeRightMouseUp      = 4
	nsEventTypeMouseMoved        = 5
	nsEventTypeLeftMouseDragged  = 6
	nsEventTypeRightMouseDragged = 7
	nsEventTypeKeyDown           = 10
	nsEventTypeKeyUp             = 11
	nsEventTypeFlagsChanged      = 12
	nsEventTypeScrollWheel       = 22
	nsEventTypeOtherMouseDown    = 25
	nsEventTypeOtherMouseUp      = 26
	nsEventTypeOtherMouseDragged = 27

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

	objcGetClass          func(name *byte) uintptr
	selRegisterName       func(name *byte) uintptr
	objcAllocateClassPair func(superClass uintptr, name *byte, extraBytes uintptr) uintptr
	objcRegisterClassPair func(cls uintptr)
	classAddMethod        func(cls uintptr, name uintptr, imp uintptr, types *byte) bool
	objcMsgSendPtr        uintptr

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
	cgDisplaySetDisplayMode                func(display uint32, mode uintptr, options uintptr) int32
	cgDisplayModeGetWidth                  func(mode uintptr) uintptr
	cgDisplayModeGetHeight                 func(mode uintptr) uintptr
	cgDisplayModeGetRefreshRate            func(mode uintptr) float64
	cgDisplayModeRelease                   func(mode uintptr)
	cgShieldingWindowLevel                 func() int32
	cgDisplayHideCursor                    func(display uint32) int32
	cgDisplayShowCursor                    func(display uint32) int32
	cgAssociateMouseAndMouseCursorPosition func(connected int32) int32
	cfArrayGetCount                        func(arr uintptr) int
	cfArrayGetValueAtIndex                 func(arr uintptr, idx int) uintptr
	cfRelease                              func(cf uintptr)

	nsApp                uintptr
	nsDefaultRunLoopMode uintptr
	savedDesktopMode     uintptr

	clsNSApplication uintptr
	clsNSWindow      uintptr
	clsNSView        uintptr
	clsMiniVKWindow  uintptr
	clsMiniVKView    uintptr
	clsNSString      uintptr
	clsNSDate        uintptr
	clsNSEvent       uintptr
	clsCAMetalLayer  uintptr

	selSharedApplication            uintptr
	selSetActivationPolicy          uintptr
	selSetPresentationOptions       uintptr
	selActivateIgnoringOtherApps    uintptr
	selFinishLaunching              uintptr
	selAlloc                        uintptr
	selInit                         uintptr
	selRelease                      uintptr
	selInitWithContentRect          uintptr
	selSetTitle                     uintptr
	selSetStyleMask                 uintptr
	selStyleMask                    uintptr
	selSetLevel                     uintptr
	selSetCollectionBehavior        uintptr
	selSetHidesOnDeactivate         uintptr
	selMakeKeyAndOrderFront         uintptr
	selMakeFirstResponder           uintptr
	selClose                        uintptr
	selContentView                  uintptr
	selSetContentView               uintptr
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

	// Keep Objective-C IMP callbacks alive
	impReturnYes uintptr
	impNoOpEvent uintptr
	darwinPrevMods uintptr
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

func registerCustomCocoaClasses() {
	if clsMiniVKWindow != 0 {
		return
	}
	impReturnYes = purego.NewCallback(func(self, cmd uintptr) uintptr {
		return 1
	})
	impNoOpEvent = purego.NewCallback(func(self, cmd, event uintptr) uintptr {
		return 0
	})

	boolTypes := CString("c@:")
	voidEventTypes := CString("v@:@")
	boolEventTypes := CString("c@:@")

	// Subclass NSWindow -> miniVKWindow so borderless fullscreen windows can become key/main
	// and unhandled key events never trigger NSBeep().
	winName := CString("miniVKWindow")
	clsMiniVKWindow = objcAllocateClassPair(clsNSWindow, &winName[0], 0)
	if clsMiniVKWindow != 0 {
		classAddMethod(clsMiniVKWindow, objcSel("canBecomeKeyWindow"), impReturnYes, &boolTypes[0])
		classAddMethod(clsMiniVKWindow, objcSel("canBecomeMainWindow"), impReturnYes, &boolTypes[0])
		classAddMethod(clsMiniVKWindow, objcSel("keyDown:"), impNoOpEvent, &voidEventTypes[0])
		classAddMethod(clsMiniVKWindow, objcSel("keyUp:"), impNoOpEvent, &voidEventTypes[0])
		objcRegisterClassPair(clsMiniVKWindow)
	} else {
		clsMiniVKWindow = objcClass("miniVKWindow")
	}

	// Subclass NSView -> miniVKView so keyDown:/keyUp:/performKeyEquivalent: never call NSBeep().
	viewName := CString("miniVKView")
	clsMiniVKView = objcAllocateClassPair(clsNSView, &viewName[0], 0)
	if clsMiniVKView != 0 {
		classAddMethod(clsMiniVKView, objcSel("acceptsFirstResponder"), impReturnYes, &boolTypes[0])
		classAddMethod(clsMiniVKView, objcSel("canBecomeKeyView"), impReturnYes, &boolTypes[0])
		classAddMethod(clsMiniVKView, objcSel("keyDown:"), impNoOpEvent, &voidEventTypes[0])
		classAddMethod(clsMiniVKView, objcSel("keyUp:"), impNoOpEvent, &voidEventTypes[0])
		classAddMethod(clsMiniVKView, objcSel("performKeyEquivalent:"), impReturnYes, &boolEventTypes[0])
		objcRegisterClassPair(clsMiniVKView)
	} else {
		clsMiniVKView = objcClass("miniVKView")
	}

	runtime.KeepAlive(boolTypes)
	runtime.KeepAlive(voidEventTypes)
	runtime.KeepAlive(boolEventTypes)
	runtime.KeepAlive(winName)
	runtime.KeepAlive(viewName)
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
	_ = bindLibrary(&objcAllocateClassPair, objcLib, "objc_allocateClassPair")
	_ = bindLibrary(&objcRegisterClassPair, objcLib, "objc_registerClassPair")
	_ = bindLibrary(&classAddMethod, objcLib, "class_addMethod")

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
	_ = bindLibrary(&cgDisplaySetDisplayMode, coreGraphLib, "CGDisplaySetDisplayMode")
	_ = bindLibrary(&cgDisplayModeGetWidth, coreGraphLib, "CGDisplayModeGetWidth")
	_ = bindLibrary(&cgDisplayModeGetHeight, coreGraphLib, "CGDisplayModeGetHeight")
	_ = bindLibrary(&cgDisplayModeGetRefreshRate, coreGraphLib, "CGDisplayModeGetRefreshRate")
	_ = bindLibrary(&cgDisplayModeRelease, coreGraphLib, "CGDisplayModeRelease")
	_ = bindLibrary(&cgShieldingWindowLevel, coreGraphLib, "CGShieldingWindowLevel")
	_ = bindLibrary(&cgDisplayHideCursor, coreGraphLib, "CGDisplayHideCursor")
	_ = bindLibrary(&cgDisplayShowCursor, coreGraphLib, "CGDisplayShowCursor")
	_ = bindLibrary(&cgAssociateMouseAndMouseCursorPosition, coreGraphLib, "CGAssociateMouseAndMouseCursorPosition")
	_ = bindLibrary(&cfArrayGetCount, coreGraphLib, "CFArrayGetCount")
	_ = bindLibrary(&cfArrayGetValueAtIndex, coreGraphLib, "CFArrayGetValueAtIndex")
	_ = bindLibrary(&cfRelease, coreGraphLib, "CFRelease")

	clsNSApplication = objcClass("NSApplication")
	clsNSWindow = objcClass("NSWindow")
	clsNSView = objcClass("NSView")
	clsNSString = objcClass("NSString")
	clsNSDate = objcClass("NSDate")
	clsNSEvent = objcClass("NSEvent")
	clsCAMetalLayer = objcClass("CAMetalLayer")

	selSharedApplication = objcSel("sharedApplication")
	selSetActivationPolicy = objcSel("setActivationPolicy:")
	selSetPresentationOptions = objcSel("setPresentationOptions:")
	selActivateIgnoringOtherApps = objcSel("activateIgnoringOtherApps:")
	selFinishLaunching = objcSel("finishLaunching")
	selAlloc = objcSel("alloc")
	selInit = objcSel("init")
	selRelease = objcSel("release")
	selInitWithContentRect = objcSel("initWithContentRect:styleMask:backing:defer:")
	selSetTitle = objcSel("setTitle:")
	selSetStyleMask = objcSel("setStyleMask:")
	selStyleMask = objcSel("styleMask")
	selSetLevel = objcSel("setLevel:")
	selSetCollectionBehavior = objcSel("setCollectionBehavior:")
	selSetHidesOnDeactivate = objcSel("setHidesOnDeactivate:")
	selMakeKeyAndOrderFront = objcSel("makeKeyAndOrderFront:")
	selMakeFirstResponder = objcSel("makeFirstResponder:")
	selClose = objcSel("close")
	selContentView = objcSel("contentView")
	selSetContentView = objcSel("setContentView:")
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

	registerCustomCocoaClasses()

	nsDefaultRunLoopMode = nsString("kCFRunLoopDefaultMode")
	nsApp = msgSend0(clsNSApplication, selSharedApplication)
	if nsApp != 0 {
		msgSend1(nsApp, selSetActivationPolicy, nsApplicationActivationPolicyRegular)
		msgSendSetBool(nsApp, selActivateIgnoringOtherApps, true)
		msgSend0(nsApp, selFinishLaunching)
	}
	return nil
}

func restoreDesktopDisplayMode() {
	if savedDesktopMode != 0 && cgDisplaySetDisplayMode != nil && cgMainDisplayID != nil {
		cgDisplaySetDisplayMode(cgMainDisplayID(), savedDesktopMode, 0)
		if cgDisplayModeRelease != nil {
			cgDisplayModeRelease(savedDesktopMode)
		}
		savedDesktopMode = 0
	}
	if nsApp != 0 && selSetPresentationOptions != 0 {
		msgSend1(nsApp, selSetPresentationOptions, nsApplicationPresentationDefault)
	}
}

func platformTerminate() {
	restoreDesktopDisplayMode()
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

func resolveDisplayID(monitor uintptr) uint32 {
	disp := uint32(monitor)
	if (disp == 0 || disp == 1) && cgMainDisplayID != nil {
		disp = cgMainDisplayID()
	}
	return disp
}

func platformGetVideoMode(monitor uintptr) GLFWvidmode {
	if cgMainDisplayID == nil {
		return GLFWvidmode{Width: 1920, Height: 1080, RedBits: 8, GreenBits: 8, BlueBits: 8, RefreshRate: 60}
	}
	disp := resolveDisplayID(monitor)
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
	disp := resolveDisplayID(monitor)
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

func applyFullscreenOnMac(w *Window, monitor uintptr, width, height, refreshRate int) {
	if w == nil || w.platformHandle == 0 {
		return
	}
	disp := resolveDisplayID(monitor)

	// Optionally switch hardware display mode if a matching CGDisplayMode exists
	if width > 0 && height > 0 && cgDisplayCopyAllDisplayModes != nil && cgDisplaySetDisplayMode != nil {
		if savedDesktopMode == 0 && cgDisplayCopyDisplayMode != nil {
			savedDesktopMode = cgDisplayCopyDisplayMode(disp)
		}
		arr := cgDisplayCopyAllDisplayModes(disp, 0)
		if arr != 0 {
			n := cfArrayGetCount(arr)
			var bestMode uintptr
			for i := 0; i < n; i++ {
				m := cfArrayGetValueAtIndex(arr, i)
				if m == 0 {
					continue
				}
				mw := int(cgDisplayModeGetWidth(m))
				mh := int(cgDisplayModeGetHeight(m))
				mr := int(cgDisplayModeGetRefreshRate(m))
				if mw == width && mh == height {
					bestMode = m
					if refreshRate <= 0 || mr == refreshRate {
						break
					}
				}
			}
			if bestMode != 0 {
				cgDisplaySetDisplayMode(disp, bestMode, 0)
			}
			cfRelease(arr)
		}
	}

	// Hide Menu Bar and Dock for true fullscreen
	if nsApp != 0 && selSetPresentationOptions != 0 {
		opts := uintptr(nsApplicationPresentationHideDock | nsApplicationPresentationHideMenuBar)
		msgSend1(nsApp, selSetPresentationOptions, opts)
	}

	// Switch window to borderless above the menu bar covering the full display bounds
	msgSend1(w.platformHandle, selSetStyleMask, nsWindowStyleMaskBorderless)
	level := uintptr(nsMainMenuWindowLevel + 1)
	if cgShieldingWindowLevel != nil {
		if sh := cgShieldingWindowLevel(); sh > 0 {
			level = uintptr(sh)
		}
	}
	msgSend1(w.platformHandle, selSetLevel, level)
	msgSendSetBool(w.platformHandle, selSetHidesOnDeactivate, true)

	screenW := int(cgDisplayPixelsWide(disp))
	screenH := int(cgDisplayPixelsHigh(disp))
	if screenW <= 0 || screenH <= 0 {
		screenW = width
		screenH = height
	}
	fullRect := cgRect{
		Origin: cgPoint{X: 0, Y: 0},
		Size:   cgSize{Width: float64(screenW), Height: float64(screenH)},
	}
	msgSendSetFrame(w.platformHandle, selSetFrameDisplay, fullRect, true)
	msgSend1(w.platformHandle, selMakeKeyAndOrderFront, 0)
	if w.platformSurface != 0 {
		msgSend1(w.platformHandle, selMakeFirstResponder, w.platformSurface)
	}

	scale := msgSendFloat(w.platformHandle, selBackingScaleFactor)
	if scale <= 0 {
		scale = 1.0
	}
	winW := width
	winH := height
	if winW <= 0 || winH <= 0 {
		winW = screenW
		winH = screenH
	}
	fbW := int(float64(winW) * scale)
	fbH := int(float64(winH) * scale)
	if w.platformLayer != 0 {
		msgSendSetFloat(w.platformLayer, selSetContentsScale, scale)
		msgSendSetSize(w.platformLayer, selSetDrawableSize, cgSize{Width: float64(fbW), Height: float64(fbH)})
	}
	w.xpos = 0
	w.ypos = 0
	w.dispatchWindowSize(winW, winH)
	w.dispatchFramebufferSize(fbW, fbH)
}

func applyWindowedOnMac(w *Window, xpos, ypos, width, height int) {
	if w == nil || w.platformHandle == 0 {
		return
	}
	restoreDesktopDisplayMode()

	style := uintptr(nsWindowStyleMaskTitled | nsWindowStyleMaskClosable | nsWindowStyleMaskMiniaturizable | nsWindowStyleMaskResizable)
	msgSend1(w.platformHandle, selSetStyleMask, style)
	msgSend1(w.platformHandle, selSetLevel, nsNormalWindowLevel)
	msgSendSetBool(w.platformHandle, selSetHidesOnDeactivate, false)
	msgSend1(w.platformHandle, selSetTitle, nsString(w.title))

	if width <= 0 {
		width = 800
	}
	if height <= 0 {
		height = 600
	}
	rect := cgRect{
		Origin: cgPoint{X: float64(xpos), Y: float64(ypos)},
		Size:   cgSize{Width: float64(width), Height: float64(height)},
	}
	msgSendSetFrame(w.platformHandle, selSetFrameDisplay, rect, true)
	msgSend1(w.platformHandle, selMakeKeyAndOrderFront, 0)
	if w.platformSurface != 0 {
		msgSend1(w.platformHandle, selMakeFirstResponder, w.platformSurface)
	}

	scale := msgSendFloat(w.platformHandle, selBackingScaleFactor)
	if scale <= 0 {
		scale = 1.0
	}
	fbW := int(float64(width) * scale)
	fbH := int(float64(height) * scale)
	if w.platformLayer != 0 {
		msgSendSetFloat(w.platformLayer, selSetContentsScale, scale)
		msgSendSetSize(w.platformLayer, selSetDrawableSize, cgSize{Width: float64(fbW), Height: float64(fbH)})
	}
	w.xpos = xpos
	w.ypos = ypos
	w.dispatchWindowSize(width, height)
	w.dispatchFramebufferSize(fbW, fbH)
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
	winClass := clsMiniVKWindow
	if winClass == 0 {
		winClass = clsNSWindow
	}
	rawWin := msgSend0(winClass, selAlloc)
	nsWin := msgSendInitWin(rawWin, selInitWithContentRect, rect, style, nsBackingStoreBuffered, false)
	if nsWin == 0 {
		return fmt.Errorf("failed to create NSWindow")
	}
	msgSend1(nsWin, selSetTitle, nsString(title))
	msgSend1(nsWin, selSetCollectionBehavior, nsWindowCollectionBehaviorFullScreenPrimary|nsWindowCollectionBehaviorManaged)

	// Install custom miniVKView so key presses never trigger NSBeep()
	var view uintptr
	if clsMiniVKView != 0 {
		view = msgSend0(msgSend0(clsMiniVKView, selAlloc), selInit)
		if view != 0 {
			msgSend1(nsWin, selSetContentView, view)
		}
	}
	if view == 0 {
		view = msgSend0(nsWin, selContentView)
	}

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
	msgSend1(nsWin, selMakeFirstResponder, view)
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

	if monitor != 0 {
		applyFullscreenOnMac(w, monitor, width, height, w.refreshRate)
	}

	if msgSendBool(nsWin, selIsKeyWindow) {
		w.platformExtra = 1
	}
	registerWindow(w)
	return nil
}

func platformDestroyWindow(w *Window) {
	if w != nil && w.platformHandle != 0 {
		if w.monitor != 0 {
			restoreDesktopDisplayMode()
		}
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
		evType := handleCocoaEvent(ev)
		// Do not pass KeyDown/KeyUp events to [NSApp sendEvent:] unless Command is held,
		// preventing NSResponder's default keyDown: from playing the macOS alert beep sound.
		if evType == nsEventTypeKeyDown || evType == nsEventTypeKeyUp {
			flags := msgSend0(ev, selModifierFlags)
			if (flags & nsEventModifierFlagCommand) == 0 {
				continue
			}
		}
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
		evType := handleCocoaEvent(ev)
		if evType != nsEventTypeKeyDown && evType != nsEventTypeKeyUp {
			msgSend1(nsApp, selSendEvent, ev)
		}
	}
	platformPollEvents()
}

func platformSetWindowMonitor(w *Window, monitor uintptr, xpos, ypos, width, height, refreshRate int) {
	if w == nil {
		return
	}
	if monitor != 0 {
		applyFullscreenOnMac(w, monitor, width, height, refreshRate)
	} else {
		applyWindowedOnMac(w, xpos, ypos, width, height)
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

// macKeyCodes maps Cocoa virtual key codes using GLFW's US ANSI key positions.
var macKeyCodes = [...]int{
	0x00: GLFW_KEY_A,
	0x01: GLFW_KEY_S,
	0x02: GLFW_KEY_D,
	0x03: GLFW_KEY_F,
	0x04: GLFW_KEY_H,
	0x05: GLFW_KEY_G,
	0x06: GLFW_KEY_Z,
	0x07: GLFW_KEY_X,
	0x08: GLFW_KEY_C,
	0x09: GLFW_KEY_V,
	0x0B: GLFW_KEY_B,
	0x0C: GLFW_KEY_Q,
	0x0D: GLFW_KEY_W,
	0x0E: GLFW_KEY_E,
	0x0F: GLFW_KEY_R,
	0x10: GLFW_KEY_Y,
	0x11: GLFW_KEY_T,
	0x12: GLFW_KEY_1,
	0x13: GLFW_KEY_2,
	0x14: GLFW_KEY_3,
	0x15: GLFW_KEY_4,
	0x16: GLFW_KEY_6,
	0x17: GLFW_KEY_5,
	0x18: GLFW_KEY_EQUAL,
	0x19: GLFW_KEY_9,
	0x1A: GLFW_KEY_7,
	0x1B: GLFW_KEY_MINUS,
	0x1C: GLFW_KEY_8,
	0x1D: GLFW_KEY_0,
	0x1E: GLFW_KEY_RIGHT_BRACKET,
	0x1F: GLFW_KEY_O,
	0x20: GLFW_KEY_U,
	0x21: GLFW_KEY_LEFT_BRACKET,
	0x22: GLFW_KEY_I,
	0x23: GLFW_KEY_P,
	0x24: GLFW_KEY_ENTER,
	0x25: GLFW_KEY_L,
	0x26: GLFW_KEY_J,
	0x27: GLFW_KEY_APOSTROPHE,
	0x28: GLFW_KEY_K,
	0x29: GLFW_KEY_SEMICOLON,
	0x2A: GLFW_KEY_BACKSLASH,
	0x2B: GLFW_KEY_COMMA,
	0x2C: GLFW_KEY_SLASH,
	0x2D: GLFW_KEY_N,
	0x2E: GLFW_KEY_M,
	0x2F: GLFW_KEY_PERIOD,
	0x30: GLFW_KEY_TAB,
	0x31: GLFW_KEY_SPACE,
	0x32: GLFW_KEY_GRAVE_ACCENT,
	0x33: GLFW_KEY_BACKSPACE,
	0x35: GLFW_KEY_ESCAPE,
	0x36: GLFW_KEY_RIGHT_SUPER,
	0x37: GLFW_KEY_LEFT_SUPER,
	0x38: GLFW_KEY_LEFT_SHIFT,
	0x39: GLFW_KEY_CAPS_LOCK,
	0x3A: GLFW_KEY_LEFT_ALT,
	0x3B: GLFW_KEY_LEFT_CONTROL,
	0x3C: GLFW_KEY_RIGHT_SHIFT,
	0x3D: GLFW_KEY_RIGHT_ALT,
	0x3E: GLFW_KEY_RIGHT_CONTROL,
	0x47: GLFW_KEY_NUM_LOCK,
	0x60: GLFW_KEY_F5,
	0x61: GLFW_KEY_F6,
	0x62: GLFW_KEY_F7,
	0x63: GLFW_KEY_F3,
	0x64: GLFW_KEY_F8,
	0x65: GLFW_KEY_F9,
	0x67: GLFW_KEY_F11,
	0x69: GLFW_KEY_PRINT_SCREEN,
	0x6D: GLFW_KEY_F10,
	0x6E: GLFW_KEY_MENU,
	0x6F: GLFW_KEY_F12,
	0x72: GLFW_KEY_INSERT,
	0x73: GLFW_KEY_HOME,
	0x74: GLFW_KEY_PAGE_UP,
	0x75: GLFW_KEY_DELETE,
	0x76: GLFW_KEY_F4,
	0x77: GLFW_KEY_END,
	0x78: GLFW_KEY_F2,
	0x79: GLFW_KEY_PAGE_DOWN,
	0x7A: GLFW_KEY_F1,
	0x7B: GLFW_KEY_LEFT,
	0x7C: GLFW_KEY_RIGHT,
	0x7D: GLFW_KEY_DOWN,
	0x7E: GLFW_KEY_UP,
}

func translateMacKey(keyCode uint16) int {
	if int(keyCode) < len(macKeyCodes) {
		if key := macKeyCodes[keyCode]; key != 0 {
			return key
		}
	}
	return GLFW_KEY_UNKNOWN
}

func handleCocoaEvent(ev uintptr) uintptr {
	if ev == 0 {
		return 0
	}
	evType := msgSend0(ev, selType)
	nsWin := msgSend0(ev, selWindow)
	w := lookupWindow(nsWin)
	if w == nil {
		// Fallback to first active window if event has nil window (e.g. during borderless fullscreen)
		glfwMu.RLock()
		for _, win := range glfwWindows {
			if win.focused {
				w = win
				break
			}
		}
		glfwMu.RUnlock()
		if w == nil {
			return evType
		}
	}
	rawFlags := msgSend0(ev, selModifierFlags)
	mods := cocoaMods(rawFlags)
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
	case nsEventTypeFlagsChanged:
		kc := uint16(msgSend0(ev, selKeyCode))
		key := translateMacKey(kc)
		if key != GLFW_KEY_UNKNOWN {
			action := GLFW_RELEASE
			if rawFlags > darwinPrevMods {
				action = GLFW_PRESS
			}
			darwinPrevMods = rawFlags
			w.dispatchKey(key, int(kc), action, mods)
		}
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
	return evType
}
