//go:build linux

package binding

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"syscall"
	"unsafe"

	"github.com/ebitengine/purego"
)

const (
	wlMarshalFlagDestroy = 1

	wlDisplayGetRegistry = 1
	wlRegistryBind       = 0

	wlCompositorCreateSurface = 0
	wlSurfaceDestroy          = 0
	wlSurfaceCommit           = 6
	wlSurfaceSetBufferScale   = 8

	wlSeatGetPointer  = 0
	wlSeatGetKeyboard = 1

	wlPointerSetCursor = 0
	wlPointerRelease   = 1

	wlKeyboardRelease = 0

	wlDataDeviceManagerGetDataDevice = 1
	wlDataOfferAccept                = 0
	wlDataOfferReceive               = 1
	wlDataOfferDestroy               = 2
	wlDataOfferFinish                = 3

	xdgWmBaseDestroy       = 0
	xdgWmBaseGetXdgSurface = 2
	xdgWmBasePong          = 3

	xdgSurfaceDestroy      = 0
	xdgSurfaceGetToplevel  = 1
	xdgSurfaceAckConfigure = 4

	xdgToplevelDestroy         = 0
	xdgToplevelSetTitle        = 2
	xdgToplevelSetAppID        = 3
	xdgToplevelSetFullscreen   = 11
	xdgToplevelUnsetFullscreen = 12

	zwpRelativePointerManagerDestroy            = 0
	zwpRelativePointerManagerGetRelativePointer = 1
	zwpRelativePointerDestroy                   = 0

	zwpPointerConstraintsDestroy     = 0
	zwpPointerConstraintsLockPointer = 1
	zwpLockedPointerDestroy          = 0
	zwpPointerConstraintsLifetimePersistent = 2

	wlSeatCapabilityPointer  = 1
	wlSeatCapabilityKeyboard = 2
)

type wlMessage struct {
	Name      *byte
	Signature *byte
	Types     *uintptr
}

type wlInterface struct {
	Name        *byte
	Version     int32
	MethodCount int32
	Methods     *wlMessage
	EventCount  int32
	Events      *wlMessage
}

type wlArray struct {
	Size  uintptr
	Alloc uintptr
	Data  uintptr
}

var (
	waylandLib uintptr

	wlDisplayConnect           func(name *byte) uintptr
	wlDisplayDisconnect        func(display uintptr)
	wlDisplayDispatch          func(display uintptr) int32
	wlDisplayDispatchPending   func(display uintptr) int32
	wlDisplayFlush             func(display uintptr) int32
	wlDisplayRoundtrip         func(display uintptr) int32
	wlDisplayGetFd             func(display uintptr) int32
	wlProxyMarshalArrayFlags   func(proxy uintptr, opcode uint32, iface uintptr, version uint32, flags uint32, args *uintptr) uintptr
	wlProxyAddDispatcher       func(proxy uintptr, dispatcher uintptr, impl uintptr, data uintptr) int32
	wlProxyDestroy             func(proxy uintptr)
	wlProxyGetVersion          func(proxy uintptr) uint32

	wlRegistryInterface          uintptr
	wlCompositorInterface        uintptr
	wlSurfaceInterface           uintptr
	wlSeatInterface              uintptr
	wlPointerInterface           uintptr
	wlKeyboardInterface          uintptr
	wlOutputInterface            uintptr
	wlDataDeviceManagerInterface uintptr
	wlDataDeviceInterface        uintptr
	wlDataOfferInterface         uintptr

	// Statically allocated custom protocol tables (xdg-shell, relative-pointer, pointer-constraints)
	wlProtocolInitOnce sync.Once
	emptyTypes         [16]uintptr

	xdgWmBaseIface          wlInterface
	xdgSurfaceIface         wlInterface
	xdgToplevelIface        wlInterface
	zwpRelPtrMgrIface       wlInterface
	zwpRelPtrIface          wlInterface
	zwpPtrConstraintsIface  wlInterface
	zwpLockedPtrIface       wlInterface

	xdgWmBaseMethods         [4]wlMessage
	xdgWmBaseEvents          [1]wlMessage
	xdgSurfaceMethods        [5]wlMessage
	xdgSurfaceEvents         [1]wlMessage
	xdgToplevelMethods       [14]wlMessage
	xdgToplevelEvents        [4]wlMessage
	zwpRelPtrMgrMethods      [2]wlMessage
	zwpRelPtrMethods         [1]wlMessage
	zwpRelPtrEvents          [1]wlMessage
	zwpPtrConstraintsMethods [3]wlMessage
	zwpLockedPtrMethods      [3]wlMessage
	zwpLockedPtrEvents       [2]wlMessage

	xdgSurfaceTypes  [4]uintptr
	xdgToplevelTypes [4]uintptr
	zwpRelPtrTypes   [4]uintptr
	zwpLockPtrTypes  [8]uintptr

	// Keep byte slices alive forever
	wlStaticStrings [][]byte

	wlDispatcherCallback uintptr

	// Active Wayland connection state
	wlDisplay            uintptr
	wlRegistry           uintptr
	wlCompositor         uintptr
	wlSeat               uintptr
	wlPointer            uintptr
	wlKeyboard           uintptr
	wlOutput             uintptr
	wlDataDevMgr         uintptr
	wlDataDev            uintptr
	wlXdgWmBase          uintptr
	wlRelPtrMgr          uintptr
	wlPtrConstraints     uintptr

	wlPointerSerial      uint32
	wlActivePointerWin   *Window
	wlActiveKeyboardWin  *Window
	wlActiveDataOffer    uintptr
	wlModes              []GLFWvidmode
	wlCurrentMode        = GLFWvidmode{Width: 1920, Height: 1080, RedBits: 8, GreenBits: 8, BlueBits: 8, RefreshRate: 60}
	wlOutputScale        = 1
	wlModState           int
)

const (
	wlTagRegistry = iota + 1
	wlTagOutput
	wlTagSeat
	wlTagPointer
	wlTagKeyboard
	wlTagDataDevice
	wlTagDataOffer
	wlTagXdgWmBase
	wlTagXdgSurface
	wlTagXdgToplevel
	wlTagRelPointer
	wlTagLockedPointer
)

func cStrStatic(s string) *byte {
	b := CString(s)
	wlStaticStrings = append(wlStaticStrings, b)
	return &b[0]
}

func makeWlMsg(name, sig string, types *uintptr) wlMessage {
	if types == nil {
		types = &emptyTypes[0]
	}
	return wlMessage{
		Name:      cStrStatic(name),
		Signature: cStrStatic(sig),
		Types:     types,
	}
}

func initWaylandProtocolTables() {
	wlProtocolInitOnce.Do(func() {
		xdgSurfaceTypes[0] = uintptr(unsafe.Pointer(&xdgSurfaceIface))
		xdgSurfaceTypes[1] = wlSurfaceInterface
		xdgToplevelTypes[0] = uintptr(unsafe.Pointer(&xdgToplevelIface))
		zwpRelPtrTypes[0] = uintptr(unsafe.Pointer(&zwpRelPtrIface))
		zwpRelPtrTypes[1] = wlPointerInterface
		zwpLockPtrTypes[0] = uintptr(unsafe.Pointer(&zwpLockedPtrIface))
		zwpLockPtrTypes[1] = wlSurfaceInterface
		zwpLockPtrTypes[2] = wlPointerInterface

		xdgWmBaseMethods = [4]wlMessage{
			makeWlMsg("destroy", "", nil),
			makeWlMsg("create_positioner", "n", nil),
			makeWlMsg("get_xdg_surface", "no", &xdgSurfaceTypes[0]),
			makeWlMsg("pong", "u", nil),
		}
		xdgWmBaseEvents = [1]wlMessage{
			makeWlMsg("ping", "u", nil),
		}
		xdgWmBaseIface = wlInterface{
			Name:        cStrStatic("xdg_wm_base"),
			Version:     5,
			MethodCount: 4,
			Methods:     &xdgWmBaseMethods[0],
			EventCount:  1,
			Events:      &xdgWmBaseEvents[0],
		}

		xdgSurfaceMethods = [5]wlMessage{
			makeWlMsg("destroy", "", nil),
			makeWlMsg("get_toplevel", "n", &xdgToplevelTypes[0]),
			makeWlMsg("get_popup", "n?oo", nil),
			makeWlMsg("set_window_geometry", "iiii", nil),
			makeWlMsg("ack_configure", "u", nil),
		}
		xdgSurfaceEvents = [1]wlMessage{
			makeWlMsg("configure", "u", nil),
		}
		xdgSurfaceIface = wlInterface{
			Name:        cStrStatic("xdg_surface"),
			Version:     5,
			MethodCount: 5,
			Methods:     &xdgSurfaceMethods[0],
			EventCount:  1,
			Events:      &xdgSurfaceEvents[0],
		}

		xdgToplevelMethods = [14]wlMessage{
			makeWlMsg("destroy", "", nil),
			makeWlMsg("set_parent", "?o", nil),
			makeWlMsg("set_title", "s", nil),
			makeWlMsg("set_app_id", "s", nil),
			makeWlMsg("show_window_menu", "ouii", nil),
			makeWlMsg("move", "ou", nil),
			makeWlMsg("resize", "ouu", nil),
			makeWlMsg("set_max_size", "ii", nil),
			makeWlMsg("set_min_size", "ii", nil),
			makeWlMsg("set_maximized", "", nil),
			makeWlMsg("unset_maximized", "", nil),
			makeWlMsg("set_fullscreen", "?o", nil),
			makeWlMsg("unset_fullscreen", "", nil),
			makeWlMsg("set_minimized", "", nil),
		}
		xdgToplevelEvents = [4]wlMessage{
			makeWlMsg("configure", "iia", nil),
			makeWlMsg("close", "", nil),
			makeWlMsg("configure_bounds", "4ii", nil),
			makeWlMsg("wm_capabilities", "5a", nil),
		}
		xdgToplevelIface = wlInterface{
			Name:        cStrStatic("xdg_toplevel"),
			Version:     5,
			MethodCount: 14,
			Methods:     &xdgToplevelMethods[0],
			EventCount:  4,
			Events:      &xdgToplevelEvents[0],
		}

		zwpRelPtrMgrMethods = [2]wlMessage{
			makeWlMsg("destroy", "", nil),
			makeWlMsg("get_relative_pointer", "no", &zwpRelPtrTypes[0]),
		}
		zwpRelPtrMgrIface = wlInterface{
			Name:        cStrStatic("zwp_relative_pointer_manager_v1"),
			Version:     1,
			MethodCount: 2,
			Methods:     &zwpRelPtrMgrMethods[0],
		}

		zwpRelPtrMethods = [1]wlMessage{
			makeWlMsg("destroy", "", nil),
		}
		zwpRelPtrEvents = [1]wlMessage{
			makeWlMsg("relative_motion", "uuffff", nil),
		}
		zwpRelPtrIface = wlInterface{
			Name:        cStrStatic("zwp_relative_pointer_v1"),
			Version:     1,
			MethodCount: 1,
			Methods:     &zwpRelPtrMethods[0],
			EventCount:  1,
			Events:      &zwpRelPtrEvents[0],
		}

		zwpPtrConstraintsMethods = [3]wlMessage{
			makeWlMsg("destroy", "", nil),
			makeWlMsg("lock_pointer", "noo?ou", &zwpLockPtrTypes[0]),
			makeWlMsg("confine_pointer", "noo?ou", nil),
		}
		zwpPtrConstraintsIface = wlInterface{
			Name:        cStrStatic("zwp_pointer_constraints_v1"),
			Version:     1,
			MethodCount: 3,
			Methods:     &zwpPtrConstraintsMethods[0],
		}

		zwpLockedPtrMethods = [3]wlMessage{
			makeWlMsg("destroy", "", nil),
			makeWlMsg("set_cursor_position_hint", "ff", nil),
			makeWlMsg("set_region", "?o", nil),
		}
		zwpLockedPtrEvents = [2]wlMessage{
			makeWlMsg("locked", "", nil),
			makeWlMsg("unlocked", "", nil),
		}
		zwpLockedPtrIface = wlInterface{
			Name:        cStrStatic("zwp_locked_pointer_v1"),
			Version:     1,
			MethodCount: 3,
			Methods:     &zwpLockedPtrMethods[0],
			EventCount:  2,
			Events:      &zwpLockedPtrEvents[0],
		}
	})
}

func wlFixedToFloat(v int32) float64 {
	return float64(v) / 256.0
}

func wlCall(proxy uintptr, opcode uint32, retIface uintptr, version uint32, flags uint32, args ...uintptr) uintptr {
	var arr [8]uintptr
	for i := 0; i < len(args) && i < len(arr); i++ {
		arr[i] = args[i]
	}
	if version == 0 && proxy != 0 && wlProxyGetVersion != nil {
		version = wlProxyGetVersion(proxy)
		if version == 0 {
			version = 1
		}
	}
	res := wlProxyMarshalArrayFlags(proxy, opcode, retIface, version, flags, &arr[0])
	runtime.KeepAlive(arr)
	return res
}

func waylandInit() error {
	if wlDisplay != 0 {
		return nil
	}
	var err error
	if waylandLib == 0 {
		waylandLib, err = loadLibrary([]string{"libwayland-client.so.0", "libwayland-client.so"})
		if err != nil {
			return fmt.Errorf("load libwayland-client: %w", err)
		}
		for _, item := range []struct {
			fn   any
			name string
		}{
			{&wlDisplayConnect, "wl_display_connect"},
			{&wlDisplayDisconnect, "wl_display_disconnect"},
			{&wlDisplayDispatch, "wl_display_dispatch"},
			{&wlDisplayDispatchPending, "wl_display_dispatch_pending"},
			{&wlDisplayFlush, "wl_display_flush"},
			{&wlDisplayRoundtrip, "wl_display_roundtrip"},
			{&wlDisplayGetFd, "wl_display_get_fd"},
			{&wlProxyMarshalArrayFlags, "wl_proxy_marshal_array_flags"},
			{&wlProxyAddDispatcher, "wl_proxy_add_dispatcher"},
			{&wlProxyDestroy, "wl_proxy_destroy"},
			{&wlProxyGetVersion, "wl_proxy_get_version"},
		} {
			if err := bindLibrary(item.fn, waylandLib, item.name); err != nil {
				return err
			}
		}
		wlRegistryInterface, _ = lookupLibrary(waylandLib, "wl_registry_interface")
		wlCompositorInterface, _ = lookupLibrary(waylandLib, "wl_compositor_interface")
		wlSurfaceInterface, _ = lookupLibrary(waylandLib, "wl_surface_interface")
		wlSeatInterface, _ = lookupLibrary(waylandLib, "wl_seat_interface")
		wlPointerInterface, _ = lookupLibrary(waylandLib, "wl_pointer_interface")
		wlKeyboardInterface, _ = lookupLibrary(waylandLib, "wl_keyboard_interface")
		wlOutputInterface, _ = lookupLibrary(waylandLib, "wl_output_interface")
		wlDataDeviceManagerInterface, _ = lookupLibrary(waylandLib, "wl_data_device_manager_interface")
		wlDataDeviceInterface, _ = lookupLibrary(waylandLib, "wl_data_device_interface")
		wlDataOfferInterface, _ = lookupLibrary(waylandLib, "wl_data_offer_interface")
	}

	initWaylandProtocolTables()
	if wlDispatcherCallback == 0 {
		wlDispatcherCallback = purego.NewCallback(waylandUniversalDispatcher)
	}

	wlDisplay = wlDisplayConnect(nil)
	if wlDisplay == 0 {
		return fmt.Errorf("wl_display_connect failed (check WAYLAND_DISPLAY)")
	}

	wlModes = nil
	wlRegistry = wlCall(wlDisplay, wlDisplayGetRegistry, wlRegistryInterface, 1, 0, 0)
	if wlRegistry == 0 {
		wlDisplayDisconnect(wlDisplay)
		wlDisplay = 0
		return fmt.Errorf("wl_display_get_registry failed")
	}
	wlProxyAddDispatcher(wlRegistry, wlDispatcherCallback, wlTagRegistry, 0)
	wlDisplayRoundtrip(wlDisplay)
	wlDisplayRoundtrip(wlDisplay)

	if wlCompositor == 0 {
		waylandTerminate()
		return fmt.Errorf("Wayland compositor missing wl_compositor global")
	}
	if wlDataDevMgr != 0 && wlSeat != 0 && wlDataDev == 0 {
		wlDataDev = wlCall(wlDataDevMgr, wlDataDeviceManagerGetDataDevice, wlDataDeviceInterface, 1, 0, 0, wlSeat)
		if wlDataDev != 0 {
			wlProxyAddDispatcher(wlDataDev, wlDispatcherCallback, wlTagDataDevice, 0)
		}
	}
	return nil
}

func waylandTerminate() {
	if wlPointer != 0 {
		wlProxyDestroy(wlPointer)
		wlPointer = 0
	}
	if wlKeyboard != 0 {
		wlProxyDestroy(wlKeyboard)
		wlKeyboard = 0
	}
	if wlDataDev != 0 {
		wlProxyDestroy(wlDataDev)
		wlDataDev = 0
	}
	if wlDataDevMgr != 0 {
		wlProxyDestroy(wlDataDevMgr)
		wlDataDevMgr = 0
	}
	if wlSeat != 0 {
		wlProxyDestroy(wlSeat)
		wlSeat = 0
	}
	if wlOutput != 0 {
		wlProxyDestroy(wlOutput)
		wlOutput = 0
	}
	if wlXdgWmBase != 0 {
		wlProxyDestroy(wlXdgWmBase)
		wlXdgWmBase = 0
	}
	if wlRelPtrMgr != 0 {
		wlProxyDestroy(wlRelPtrMgr)
		wlRelPtrMgr = 0
	}
	if wlPtrConstraints != 0 {
		wlProxyDestroy(wlPtrConstraints)
		wlPtrConstraints = 0
	}
	if wlCompositor != 0 {
		wlProxyDestroy(wlCompositor)
		wlCompositor = 0
	}
	if wlRegistry != 0 {
		wlProxyDestroy(wlRegistry)
		wlRegistry = 0
	}
	if wlDisplay != 0 {
		wlDisplayDisconnect(wlDisplay)
		wlDisplay = 0
	}
}

func waylandGetPrimaryMonitor() uintptr {
	if wlOutput != 0 {
		return wlOutput
	}
	return 1
}

func waylandGetVideoMode(monitor uintptr) GLFWvidmode {
	return wlCurrentMode
}

func waylandGetVideoModes(monitor uintptr) []GLFWvidmode {
	if len(wlModes) > 0 {
		out := make([]GLFWvidmode, len(wlModes))
		copy(out, wlModes)
		return out
	}
	return []GLFWvidmode{wlCurrentMode}
}

func waylandCreateWindow(w *Window, width, height int, title string, monitor uintptr) error {
	if wlDisplay == 0 {
		if err := waylandInit(); err != nil {
			return err
		}
	}
	surf := wlCall(wlCompositor, wlCompositorCreateSurface, wlSurfaceInterface, 4, 0, 0)
	if surf == 0 {
		return fmt.Errorf("wl_compositor_create_surface failed")
	}
	w.platformDisplay = wlDisplay
	w.platformSurface = surf
	w.platformHandle = surf
	w.handle = surf
	if wlOutputScale > 1 {
		wlCall(surf, wlSurfaceSetBufferScale, 0, 3, 0, uintptr(wlOutputScale))
		w.fbWidth = width * wlOutputScale
		w.fbHeight = height * wlOutputScale
	}
	registerWindow(w)

	if wlXdgWmBase != 0 {
		xdgSurf := wlCall(wlXdgWmBase, xdgWmBaseGetXdgSurface, uintptr(unsafe.Pointer(&xdgSurfaceIface)), 1, 0, 0, surf)
		if xdgSurf != 0 {
			w.platformLayer = xdgSurf
			wlProxyAddDispatcher(xdgSurf, wlDispatcherCallback, wlTagXdgSurface, surf)
			toplevel := wlCall(xdgSurf, xdgSurfaceGetToplevel, uintptr(unsafe.Pointer(&xdgToplevelIface)), 1, 0, 0)
			if toplevel != 0 {
				w.platformExtra = toplevel
				wlProxyAddDispatcher(toplevel, wlDispatcherCallback, wlTagXdgToplevel, surf)
				cTitle := CString(title)
				wlCall(toplevel, xdgToplevelSetTitle, 0, 1, 0, uintptr(unsafe.Pointer(&cTitle[0])))
				runtime.KeepAlive(cTitle)
				if monitor != 0 {
					wlCall(toplevel, xdgToplevelSetFullscreen, 0, 1, 0, 0)
				}
			}
		}
	}

	wlCall(surf, wlSurfaceCommit, 0, 1, 0)
	wlDisplayRoundtrip(wlDisplay)
	return nil
}

func waylandDestroyWindow(w *Window) {
	if w == nil || wlDisplay == 0 {
		return
	}
	if w.cursorMode == GLFW_CURSOR_DISABLED {
		waylandSetCursorMode(w, GLFW_CURSOR_NORMAL)
	}
	if wlActivePointerWin == w {
		wlActivePointerWin = nil
	}
	if wlActiveKeyboardWin == w {
		wlActiveKeyboardWin = nil
	}
	if w.platformExtra != 0 {
		wlCall(w.platformExtra, xdgToplevelDestroy, 0, 1, wlMarshalFlagDestroy)
		w.platformExtra = 0
	}
	if w.platformLayer != 0 {
		wlCall(w.platformLayer, xdgSurfaceDestroy, 0, 1, wlMarshalFlagDestroy)
		w.platformLayer = 0
	}
	if w.platformSurface != 0 {
		wlCall(w.platformSurface, wlSurfaceDestroy, 0, 1, wlMarshalFlagDestroy)
		w.platformSurface = 0
		w.platformHandle = 0
	}
	wlDisplayFlush(wlDisplay)
}

func waylandPollEvents() {
	if wlDisplay == 0 {
		return
	}
	wlDisplayFlush(wlDisplay)
	wlDisplayDispatchPending(wlDisplay)
	fd := int(wlDisplayGetFd(wlDisplay))
	for pollReadable(fd, 0) {
		if wlDisplayDispatch(wlDisplay) < 0 {
			break
		}
	}
}

func waylandWaitEventsTimeout(timeout float64) {
	if wlDisplay == 0 {
		return
	}
	wlDisplayFlush(wlDisplay)
	if wlDisplayDispatchPending(wlDisplay) > 0 {
		return
	}
	fd := int(wlDisplayGetFd(wlDisplay))
	if pollReadable(fd, timeout) {
		wlDisplayDispatch(wlDisplay)
	}
}

func waylandSetWindowMonitor(w *Window, monitor uintptr, xpos, ypos, width, height, refreshRate int) {
	if w == nil {
		return
	}
	if w.platformExtra != 0 && wlDisplay != 0 {
		if monitor != 0 {
			wlCall(w.platformExtra, xdgToplevelSetFullscreen, 0, 1, 0, 0)
		} else {
			wlCall(w.platformExtra, xdgToplevelUnsetFullscreen, 0, 1, 0)
		}
		if w.platformSurface != 0 {
			wlCall(w.platformSurface, wlSurfaceCommit, 0, 1, 0)
		}
		wlDisplayFlush(wlDisplay)
	}
	w.xpos = xpos
	w.ypos = ypos
	if width > 0 && height > 0 {
		w.dispatchWindowSize(width, height)
		w.dispatchFramebufferSize(width*wlOutputScale, height*wlOutputScale)
	}
}

func waylandGetWindowSize(w *Window) (int, int) {
	if w == nil {
		return 0, 0
	}
	return w.width, w.height
}

func waylandGetFramebufferSize(w *Window) (int, int) {
	if w == nil {
		return 0, 0
	}
	return w.fbWidth, w.fbHeight
}

func waylandGetWindowPos(w *Window) (int, int) {
	if w == nil {
		return 0, 0
	}
	return w.xpos, w.ypos
}

var (
	wlRelPointer    uintptr
	wlLockedPointer uintptr
)

func waylandSetCursorMode(w *Window, mode int) {
	if w == nil || wlDisplay == 0 {
		return
	}
	if mode == GLFW_CURSOR_DISABLED {
		if wlPointer != 0 {
			wlCall(wlPointer, wlPointerSetCursor, 0, 1, 0, uintptr(wlPointerSerial), 0, 0, 0)
		}
		if wlRelPtrMgr != 0 && wlPointer != 0 && wlRelPointer == 0 {
			wlRelPointer = wlCall(wlRelPtrMgr, zwpRelativePointerManagerGetRelativePointer, uintptr(unsafe.Pointer(&zwpRelPtrIface)), 1, 0, 0, wlPointer)
			if wlRelPointer != 0 {
				wlProxyAddDispatcher(wlRelPointer, wlDispatcherCallback, wlTagRelPointer, w.platformSurface)
			}
		}
		if wlPtrConstraints != 0 && wlPointer != 0 && w.platformSurface != 0 && wlLockedPointer == 0 {
			wlLockedPointer = wlCall(wlPtrConstraints, zwpPointerConstraintsLockPointer, uintptr(unsafe.Pointer(&zwpLockedPtrIface)), 1, 0, 0, w.platformSurface, wlPointer, 0, zwpPointerConstraintsLifetimePersistent)
			if wlLockedPointer != 0 {
				wlProxyAddDispatcher(wlLockedPointer, wlDispatcherCallback, wlTagLockedPointer, w.platformSurface)
			}
		}
		wlDisplayFlush(wlDisplay)
	} else {
		if wlRelPointer != 0 {
			wlCall(wlRelPointer, zwpRelativePointerDestroy, 0, 1, wlMarshalFlagDestroy)
			wlRelPointer = 0
		}
		if wlLockedPointer != 0 {
			wlCall(wlLockedPointer, zwpLockedPointerDestroy, 0, 1, wlMarshalFlagDestroy)
			wlLockedPointer = 0
		}
		wlDisplayFlush(wlDisplay)
	}
}

func waylandGetCursorPos(w *Window) (float64, float64) {
	if w == nil {
		return 0, 0
	}
	return w.cursorX, w.cursorY
}

func waylandCreateWindowSurface(instance uintptr, w *Window, allocator uintptr, surface *uint64) (int32, error) {
	if err := bindInstance(&VkCreateWaylandSurfaceKHR, instance, "vkCreateWaylandSurfaceKHR"); err != nil {
		return VK_ERROR_INITIALIZATION_FAILED, err
	}
	info := VkWaylandSurfaceCreateInfoKHR{
		SType:   VK_STRUCTURE_TYPE_WAYLAND_SURFACE_CREATE_INFO_KHR,
		Display: wlDisplay,
		Surface: w.platformSurface,
	}
	res := VkCreateWaylandSurfaceKHR(instance, &info, allocator, surface)
	if res != VK_SUCCESS {
		return res, fmt.Errorf("vkCreateWaylandSurfaceKHR returned %d", res)
	}
	return VK_SUCCESS, nil
}

func translateEvdevKey(code uint32) (int, rune) {
	switch code {
	case 1:
		return GLFW_KEY_ESCAPE, 0
	case 2:
		return GLFW_KEY_1, '1'
	case 3:
		return GLFW_KEY_2, '2'
	case 4:
		return GLFW_KEY_3, '3'
	case 5:
		return GLFW_KEY_4, '4'
	case 6:
		return GLFW_KEY_5, '5'
	case 7:
		return GLFW_KEY_6, '6'
	case 8:
		return GLFW_KEY_7, '7'
	case 9:
		return GLFW_KEY_8, '8'
	case 10:
		return GLFW_KEY_9, '9'
	case 11:
		return GLFW_KEY_0, '0'
	case 14:
		return GLFW_KEY_BACKSPACE, 0
	case 15:
		return GLFW_KEY_TAB, 0
	case 16:
		return GLFW_KEY_Q, 'q'
	case 17:
		return GLFW_KEY_W, 'w'
	case 18:
		return GLFW_KEY_E, 'e'
	case 19:
		return GLFW_KEY_R, 'r'
	case 20:
		return GLFW_KEY_T, 't'
	case 21:
		return GLFW_KEY_Y, 'y'
	case 22:
		return GLFW_KEY_U, 'u'
	case 23:
		return GLFW_KEY_I, 'i'
	case 24:
		return GLFW_KEY_O, 'o'
	case 25:
		return GLFW_KEY_P, 'p'
	case 28:
		return GLFW_KEY_ENTER, 0
	case 29:
		return GLFW_KEY_LEFT_CONTROL, 0
	case 30:
		return GLFW_KEY_A, 'a'
	case 31:
		return GLFW_KEY_S, 's'
	case 32:
		return GLFW_KEY_D, 'd'
	case 33:
		return GLFW_KEY_F, 'f'
	case 34:
		return GLFW_KEY_G, 'g'
	case 35:
		return GLFW_KEY_H, 'h'
	case 36:
		return GLFW_KEY_J, 'j'
	case 37:
		return GLFW_KEY_K, 'k'
	case 38:
		return GLFW_KEY_L, 'l'
	case 42:
		return GLFW_KEY_LEFT_SHIFT, 0
	case 44:
		return GLFW_KEY_Z, 'z'
	case 45:
		return GLFW_KEY_X, 'x'
	case 46:
		return GLFW_KEY_C, 'c'
	case 47:
		return GLFW_KEY_V, 'v'
	case 48:
		return GLFW_KEY_B, 'b'
	case 49:
		return GLFW_KEY_N, 'n'
	case 50:
		return GLFW_KEY_M, 'm'
	case 54:
		return GLFW_KEY_RIGHT_SHIFT, 0
	case 56:
		return GLFW_KEY_LEFT_ALT, 0
	case 57:
		return GLFW_KEY_SPACE, ' '
	case 59:
		return GLFW_KEY_F1, 0
	case 60:
		return GLFW_KEY_F2, 0
	case 61:
		return GLFW_KEY_F3, 0
	case 62:
		return GLFW_KEY_F4, 0
	case 63:
		return GLFW_KEY_F5, 0
	case 64:
		return GLFW_KEY_F6, 0
	case 65:
		return GLFW_KEY_F7, 0
	case 66:
		return GLFW_KEY_F8, 0
	case 67:
		return GLFW_KEY_F9, 0
	case 68:
		return GLFW_KEY_F10, 0
	case 87:
		return GLFW_KEY_F11, 0
	case 88:
		return GLFW_KEY_F12, 0
	case 97:
		return GLFW_KEY_RIGHT_CONTROL, 0
	case 100:
		return GLFW_KEY_RIGHT_ALT, 0
	case 102:
		return GLFW_KEY_HOME, 0
	case 103:
		return GLFW_KEY_UP, 0
	case 104:
		return GLFW_KEY_PAGE_UP, 0
	case 105:
		return GLFW_KEY_LEFT, 0
	case 106:
		return GLFW_KEY_RIGHT, 0
	case 107:
		return GLFW_KEY_END, 0
	case 108:
		return GLFW_KEY_DOWN, 0
	case 109:
		return GLFW_KEY_PAGE_DOWN, 0
	case 110:
		return GLFW_KEY_INSERT, 0
	case 111:
		return GLFW_KEY_DELETE, 0
	case 125:
		return GLFW_KEY_LEFT_SUPER, 0
	case 126:
		return GLFW_KEY_RIGHT_SUPER, 0
	default:
		return GLFW_KEY_UNKNOWN, 0
	}
}

func waylandUniversalDispatcher(implData, target uintptr, opcode uint32, message uintptr, rawArgs *uintptr) uintptr {
	if rawArgs == nil {
		return 0
	}
	args := unsafe.Slice(rawArgs, 12)
	switch implData {
	case wlTagRegistry:
		if opcode == 0 { // global(name: u, interface: s, version: u)
			name := uint32(args[0])
			ifaceName := String((*byte)(unsafe.Pointer(args[1])))
			ver := uint32(args[2])
			cname := CString(ifaceName)
			switch ifaceName {
			case "wl_compositor":
				if ver > 4 {
					ver = 4
				}
				wlCompositor = wlCall(target, wlRegistryBind, wlCompositorInterface, ver, 0, uintptr(name), uintptr(unsafe.Pointer(&cname[0])), uintptr(ver), 0)
			case "wl_seat":
				if ver > 5 {
					ver = 5
				}
				wlSeat = wlCall(target, wlRegistryBind, wlSeatInterface, ver, 0, uintptr(name), uintptr(unsafe.Pointer(&cname[0])), uintptr(ver), 0)
				if wlSeat != 0 {
					wlProxyAddDispatcher(wlSeat, wlDispatcherCallback, wlTagSeat, 0)
				}
			case "wl_output":
				if ver > 3 {
					ver = 3
				}
				wlOutput = wlCall(target, wlRegistryBind, wlOutputInterface, ver, 0, uintptr(name), uintptr(unsafe.Pointer(&cname[0])), uintptr(ver), 0)
				if wlOutput != 0 {
					wlProxyAddDispatcher(wlOutput, wlDispatcherCallback, wlTagOutput, 0)
				}
			case "wl_data_device_manager":
				if ver > 3 {
					ver = 3
				}
				wlDataDevMgr = wlCall(target, wlRegistryBind, wlDataDeviceManagerInterface, ver, 0, uintptr(name), uintptr(unsafe.Pointer(&cname[0])), uintptr(ver), 0)
			case "xdg_wm_base":
				if ver > 2 {
					ver = 2
				}
				wlXdgWmBase = wlCall(target, wlRegistryBind, uintptr(unsafe.Pointer(&xdgWmBaseIface)), ver, 0, uintptr(name), uintptr(unsafe.Pointer(&cname[0])), uintptr(ver), 0)
				if wlXdgWmBase != 0 {
					wlProxyAddDispatcher(wlXdgWmBase, wlDispatcherCallback, wlTagXdgWmBase, 0)
				}
			case "zwp_relative_pointer_manager_v1":
				wlRelPtrMgr = wlCall(target, wlRegistryBind, uintptr(unsafe.Pointer(&zwpRelPtrMgrIface)), 1, 0, uintptr(name), uintptr(unsafe.Pointer(&cname[0])), 1, 0)
			case "zwp_pointer_constraints_v1":
				wlPtrConstraints = wlCall(target, wlRegistryBind, uintptr(unsafe.Pointer(&zwpPtrConstraintsIface)), 1, 0, uintptr(name), uintptr(unsafe.Pointer(&cname[0])), 1, 0)
			}
			runtime.KeepAlive(cname)
		}
	case wlTagOutput:
		switch opcode {
		case 1: // mode(flags: u, width: i, height: i, refresh: i)
			flags := uint32(args[0])
			w := int32(args[1])
			h := int32(args[2])
			refresh := int32(args[3]) / 1000
			if refresh <= 0 {
				refresh = 60
			}
			if w > 0 && h > 0 {
				m := GLFWvidmode{Width: w, Height: h, RedBits: 8, GreenBits: 8, BlueBits: 8, RefreshRate: refresh}
				wlModes = append(wlModes, m)
				if (flags & 1) != 0 { // WL_OUTPUT_MODE_CURRENT
					wlCurrentMode = m
				}
			}
		case 3: // scale(factor: i)
			factor := int(int32(args[0]))
			if factor > 0 {
				wlOutputScale = factor
			}
		}
	case wlTagSeat:
		if opcode == 0 { // capabilities(capabilities: u)
			caps := uint32(args[0])
			if (caps&wlSeatCapabilityPointer) != 0 && wlPointer == 0 {
				wlPointer = wlCall(target, wlSeatGetPointer, wlPointerInterface, 1, 0, 0)
				if wlPointer != 0 {
					wlProxyAddDispatcher(wlPointer, wlDispatcherCallback, wlTagPointer, 0)
				}
			}
			if (caps&wlSeatCapabilityKeyboard) != 0 && wlKeyboard == 0 {
				wlKeyboard = wlCall(target, wlSeatGetKeyboard, wlKeyboardInterface, 1, 0, 0)
				if wlKeyboard != 0 {
					wlProxyAddDispatcher(wlKeyboard, wlDispatcherCallback, wlTagKeyboard, 0)
				}
			}
		}
	case wlTagPointer:
		switch opcode {
		case 0: // enter(serial: u, surface: o, sx: f, sy: f)
			wlPointerSerial = uint32(args[0])
			surf := args[1]
			w := lookupWindow(surf)
			wlActivePointerWin = w
			if w != nil {
				sx := wlFixedToFloat(int32(args[2]))
				sy := wlFixedToFloat(int32(args[3]))
				if w.cursorMode != GLFW_CURSOR_DISABLED {
					w.dispatchCursorPos(sx, sy)
				}
			}
		case 1: // leave(serial: u, surface: o)
			wlPointerSerial = uint32(args[0])
		case 2: // motion(time: u, sx: f, sy: f)
			w := wlActivePointerWin
			if w != nil && w.cursorMode != GLFW_CURSOR_DISABLED {
				sx := wlFixedToFloat(int32(args[1]))
				sy := wlFixedToFloat(int32(args[2]))
				w.dispatchCursorPos(sx, sy)
			}
		case 3: // button(serial: u, time: u, button: u, state: u)
			wlPointerSerial = uint32(args[0])
			w := wlActivePointerWin
			if w != nil {
				btnCode := uint32(args[2])
				state := uint32(args[3])
				action := GLFW_RELEASE
				if state != 0 {
					action = GLFW_PRESS
				}
				btn := GLFW_MOUSE_BUTTON_LEFT
				switch btnCode {
				case 0x110:
					btn = GLFW_MOUSE_BUTTON_LEFT
				case 0x111:
					btn = GLFW_MOUSE_BUTTON_RIGHT
				case 0x112:
					btn = GLFW_MOUSE_BUTTON_MIDDLE
				}
				w.dispatchMouseButton(btn, action, wlModState)
			}
		case 4: // axis(time: u, axis: u, value: f)
			w := wlActivePointerWin
			if w != nil {
				axis := uint32(args[1])
				val := wlFixedToFloat(int32(args[2])) / 10.0
				if axis == 0 { // vertical
					w.dispatchScroll(0, -val)
				} else {
					w.dispatchScroll(val, 0)
				}
			}
		}
	case wlTagKeyboard:
		switch opcode {
		case 0: // keymap(format: u, fd: h, size: u)
			fd := int(int32(args[1]))
			if fd >= 0 {
				_ = syscall.Close(fd)
			}
		case 1: // enter(serial: u, surface: o, keys: a)
			surf := args[1]
			w := lookupWindow(surf)
			wlActiveKeyboardWin = w
			if w != nil {
				w.dispatchFocus(true)
			}
		case 2: // leave(serial: u, surface: o)
			surf := args[1]
			if w := lookupWindow(surf); w != nil {
				w.dispatchFocus(false)
			}
			wlActiveKeyboardWin = nil
		case 3: // key(serial: u, time: u, key: u, state: u)
			w := wlActiveKeyboardWin
			if w != nil {
				evKey := uint32(args[2])
				state := uint32(args[3])
				key, ch := translateEvdevKey(evKey)
				action := GLFW_RELEASE
				if state != 0 {
					action = GLFW_PRESS
				}
				w.dispatchKey(key, int(evKey+8), action, wlModState)
				if action == GLFW_PRESS && ch != 0 {
					if (wlModState & GLFW_MOD_SHIFT) != 0 && ch >= 'a' && ch <= 'z' {
						ch = ch - 'a' + 'A'
					}
					w.dispatchChar(ch)
				}
			}
		case 4: // modifiers(serial: u, depressed: u, latched: u, locked: u, group: u)
			dep := uint32(args[1])
			mods := 0
			if (dep & 1) != 0 {
				mods |= GLFW_MOD_SHIFT
			}
			if (dep & 4) != 0 {
				mods |= GLFW_MOD_CONTROL
			}
			if (dep & 8) != 0 {
				mods |= GLFW_MOD_ALT
			}
			if (dep & 64) != 0 {
				mods |= GLFW_MOD_SUPER
			}
			wlModState = mods
		}
	case wlTagDataDevice:
		switch opcode {
		case 0: // data_offer(id: n)
			offer := args[0]
			if offer != 0 {
				wlActiveDataOffer = offer
				wlProxyAddDispatcher(offer, wlDispatcherCallback, wlTagDataOffer, 0)
			}
		case 1: // enter(serial: u, surface: o, x: f, y: f, id: ?o)
			serial := uint32(args[0])
			surf := args[1]
			offer := args[4]
			wlActivePointerWin = lookupWindow(surf)
			if offer != 0 {
				wlActiveDataOffer = offer
				mime := CString("text/uri-list")
				wlCall(offer, wlDataOfferAccept, 0, 1, 0, uintptr(serial), uintptr(unsafe.Pointer(&mime[0])))
				runtime.KeepAlive(mime)
			}
		case 4: // drop()
			if wlActiveDataOffer != 0 && wlActivePointerWin != nil {
				var p [2]int
				if syscall.Pipe(p[:]) == nil {
					mime := CString("text/uri-list")
					wlCall(wlActiveDataOffer, wlDataOfferReceive, 0, 1, 0, uintptr(unsafe.Pointer(&mime[0])), uintptr(p[1]))
					runtime.KeepAlive(mime)
					wlDisplayFlush(wlDisplay)
					_ = syscall.Close(p[1])
					f := os.NewFile(uintptr(p[0]), "wayland-drop")
					data, _ := io.ReadAll(f)
					_ = f.Close()
					paths := parseURIList(string(data))
					if len(paths) > 0 {
						wlActivePointerWin.dispatchDrop(paths)
					}
				}
				wlCall(wlActiveDataOffer, wlDataOfferDestroy, 0, 1, wlMarshalFlagDestroy)
				wlActiveDataOffer = 0
			}
		}
	case wlTagXdgWmBase:
		if opcode == 0 { // ping(serial: u)
			serial := uint32(args[0])
			wlCall(target, xdgWmBasePong, 0, 1, 0, uintptr(serial))
		}
	case wlTagXdgSurface:
		if opcode == 0 { // configure(serial: u)
			serial := uint32(args[0])
			wlCall(target, xdgSurfaceAckConfigure, 0, 1, 0, uintptr(serial))
		}
	case wlTagXdgToplevel:
		switch opcode {
		case 0: // configure(width: i, height: i, states: a)
			width := int(int32(args[0]))
			height := int(int32(args[1]))
			if width > 0 && height > 0 {
				for _, w := range glfwWindows {
					if w.platformExtra == target {
						w.dispatchWindowSize(width, height)
						w.dispatchFramebufferSize(width*wlOutputScale, height*wlOutputScale)
					}
				}
			}
		case 1: // close()
			for _, w := range glfwWindows {
				if w.platformExtra == target {
					w.shouldClose = true
				}
			}
		}
	case wlTagRelPointer:
		if opcode == 0 { // relative_motion(utime_hi: u, utime_lo: u, dx: f, dy: f, dx_unaccel: f, dy_unaccel: f)
			dx := wlFixedToFloat(int32(args[4]))
			dy := wlFixedToFloat(int32(args[5]))
			if dx == 0 && dy == 0 {
				dx = wlFixedToFloat(int32(args[2]))
				dy = wlFixedToFloat(int32(args[3]))
			}
			if wlActivePointerWin != nil && wlActivePointerWin.cursorMode == GLFW_CURSOR_DISABLED {
				wlActivePointerWin.dispatchRelativeCursorDelta(dx, dy)
			}
		}
	}
	return 0
}
