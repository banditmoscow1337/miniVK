//go:build windows

package binding

import (
	"fmt"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

const (
	wsOverlapped       = 0x00000000
	wsCaption          = 0x00C00000
	wsSysMenu          = 0x00080000
	wsThickFrame       = 0x00040000
	wsMinimizeBox      = 0x00020000
	wsMaximizeBox      = 0x00010000
	wsPopup            = 0x80000000
	wsVisible          = 0x10000000
	wsClipChildren     = 0x02000000
	wsClipSiblings     = 0x04000000
	wsOverlappedWindow = wsOverlapped | wsCaption | wsSysMenu | wsThickFrame | wsMinimizeBox | wsMaximizeBox

	wsExAppWindow  = 0x00040000
	wsExWindowEdge = 0x00000100
	wsExAcceptFiles = 0x00000010

	swShow        = 5
	cwUseDefault  = -0x80000000
	pmRemove      = 0x0001
	qsAllInput    = 0x04FF
	mwmoInputAvailable = 0x0004

	wmDestroy       = 0x0002
	wmMove          = 0x0003
	wmSize          = 0x0005
	wmSetFocus      = 0x0007
	wmKillFocus     = 0x0008
	wmClose         = 0x0010
	wmQuit          = 0x0012
	wmInput         = 0x00FF
	wmKeyDown       = 0x0100
	wmKeyUp         = 0x0101
	wmChar          = 0x0102
	wmSysKeyDown    = 0x0104
	wmSysKeyUp      = 0x0105
	wmMouseMove     = 0x0200
	wmLButtonDown   = 0x0201
	wmLButtonUp     = 0x0202
	wmRButtonDown   = 0x0204
	wmRButtonUp     = 0x0205
	wmMButtonDown   = 0x0207
	wmMButtonUp     = 0x0208
	wmMouseWheel    = 0x020A
	wmMouseHWheel   = 0x020E
	wmXButtonDown   = 0x020B
	wmXButtonUp     = 0x020C
	wmDropFiles     = 0x0233

	idcArrow = 32512
	csHRedraw = 0x0002
	csVRedraw = 0x0001
	csOwnDC   = 0x0020

	enumCurrentSettings = 0xFFFFFFFF
	cdsFullscreen       = 0x00000004
	dmPelsWidth         = 0x00080000
	dmPelsHeight        = 0x00100000
	dmDisplayFrequency  = 0x00400000

	swpNoZOrder       = 0x0004
	swpNoOwnerZOrder  = 0x0200
	swpFrameChanged   = 0x0020
	swpShowWindow     = 0x0040
	hwndTop           = 0

	ridInput      = 0x10000003
	rimTypeMouse  = 0
	ridnevRemove  = 0x00000001

	vkShift    = 0x10
	vkControl  = 0x11
	vkMenu     = 0x12
	vkLWin     = 0x5B
	vkRWin     = 0x5C
	vkCapital  = 0x14
	vkNumLock  = 0x90
)

type winPoint struct {
	X, Y int32
}

type winRect struct {
	Left, Top, Right, Bottom int32
}

type winMsg struct {
	HWnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      winPoint
}

type wndClassExW struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     uintptr
	HIcon         uintptr
	HCursor       uintptr
	HbrBackground uintptr
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       uintptr
}

type devModeW struct {
	DmDeviceName       [32]uint16
	DmSpecVersion      uint16
	DmDriverVersion    uint16
	DmSize             uint16
	DmDriverExtra      uint16
	DmFields           uint32
	DmUnion1           [16]byte
	DmColor            int16
	DmDuplex           int16
	DmYResolution      int16
	DmTTOption         int16
	DmCollate          int16
	DmFormName         [32]uint16
	DmLogPixels        uint16
	DmBitsPerPel       uint32
	DmPelsWidth        uint32
	DmPelsHeight       uint32
	DmDisplayFlags     uint32
	DmDisplayFrequency uint32
	DmICMMethod        uint32
	DmICMIntent        uint32
	DmMediaType        uint32
	DmDitherType       uint32
	DmReserved1        uint32
	DmReserved2        uint32
	DmPanningWidth     uint32
	DmPanningHeight    uint32
}

type rawInputDevice struct {
	UsUsagePage uint16
	UsUsage     uint16
	DwFlags     uint32
	HwndTarget  uintptr
}

type rawInputHeader struct {
	DwType  uint32
	DwSize  uint32
	HDevice uintptr
	WParam  uintptr
}

type rawMouse struct {
	UsFlags            uint16
	_                  uint16
	UsButtonFlags      uint16
	UsButtonData       uint16
	UlRawButtons       uint32
	LLastX             int32
	LLastY             int32
	UlExtraInformation uint32
}

type rawInput struct {
	Header rawInputHeader
	Mouse  rawMouse
}

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	shell32  = syscall.NewLazyDLL("shell32.dll")

	procGetModuleHandleW            = kernel32.NewProc("GetModuleHandleW")
	procSetProcessDPIAware          = user32.NewProc("SetProcessDPIAware")
	procRegisterClassExW            = user32.NewProc("RegisterClassExW")
	procUnregisterClassW            = user32.NewProc("UnregisterClassW")
	procCreateWindowExW             = user32.NewProc("CreateWindowExW")
	procDestroyWindow               = user32.NewProc("DestroyWindow")
	procDefWindowProcW              = user32.NewProc("DefWindowProcW")
	procShowWindow                  = user32.NewProc("ShowWindow")
	procUpdateWindow                = user32.NewProc("UpdateWindow")
	procPeekMessageW                = user32.NewProc("PeekMessageW")
	procTranslateMessage            = user32.NewProc("TranslateMessage")
	procDispatchMessageW            = user32.NewProc("DispatchMessageW")
	procMsgWaitForMultipleObjectsEx = user32.NewProc("MsgWaitForMultipleObjectsEx")
	procGetClientRect               = user32.NewProc("GetClientRect")
	procGetWindowRect               = user32.NewProc("GetWindowRect")
	procAdjustWindowRectEx          = user32.NewProc("AdjustWindowRectEx")
	procClientToScreen              = user32.NewProc("ClientToScreen")
	procScreenToClient              = user32.NewProc("ScreenToClient")
	procSetWindowLongPtrW           = user32.NewProc("SetWindowLongPtrW")
	procSetWindowPos                = user32.NewProc("SetWindowPos")
	procGetForegroundWindow         = user32.NewProc("GetForegroundWindow")
	procLoadCursorW                 = user32.NewProc("LoadCursorW")
	procGetCursorPos                = user32.NewProc("GetCursorPos")
	procSetCursorPos                = user32.NewProc("SetCursorPos")
	procShowCursor                  = user32.NewProc("ShowCursor")
	procClipCursor                  = user32.NewProc("ClipCursor")
	procEnumDisplaySettingsW        = user32.NewProc("EnumDisplaySettingsW")
	procChangeDisplaySettingsExW    = user32.NewProc("ChangeDisplaySettingsExW")
	procGetKeyState                 = user32.NewProc("GetKeyState")
	procRegisterRawInputDevices     = user32.NewProc("RegisterRawInputDevices")
	procGetRawInputData             = user32.NewProc("GetRawInputData")

	procDragAcceptFiles = shell32.NewProc("DragAcceptFiles")
	procDragQueryFileW  = shell32.NewProc("DragQueryFileW")
	procDragFinish      = shell32.NewProc("DragFinish")

	winHInstance   uintptr
	winClassAtom   uintptr
	winClassName   = utf16Ptr("miniVK_Window_Class")
	winWndProcCallback uintptr
)

func utf16Ptr(s string) *uint16 {
	u := utf16.Encode([]rune(s + "\x00"))
	return &u[0]
}

func platformInit() error {
	if procSetProcessDPIAware.Find() == nil {
		_, _, _ = procSetProcessDPIAware.Call()
	}
	hInst, _, _ := procGetModuleHandleW.Call(0)
	winHInstance = hInst

	if winWndProcCallback == 0 {
		winWndProcCallback = syscall.NewCallback(winWndProc)
	}
	cursor, _, _ := procLoadCursorW.Call(0, idcArrow)
	wc := wndClassExW{
		CbSize:        uint32(unsafe.Sizeof(wndClassExW{})),
		Style:         csHRedraw | csVRedraw | csOwnDC,
		LpfnWndProc:   winWndProcCallback,
		HInstance:     winHInstance,
		HCursor:       cursor,
		LpszClassName: winClassName,
	}
	atom, _, err := procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))
	if atom == 0 && winClassAtom == 0 {
		return fmt.Errorf("RegisterClassExW failed: %v", err)
	}
	if atom != 0 {
		winClassAtom = atom
	}
	return nil
}

func platformTerminate() {
	_, _, _ = procClipCursor.Call(0)
	if winClassAtom != 0 {
		_, _, _ = procUnregisterClassW.Call(uintptr(unsafe.Pointer(winClassName)), winHInstance)
		winClassAtom = 0
	}
}

func platformGetPrimaryMonitor() uintptr {
	return 1
}

func platformGetVideoMode(monitor uintptr) GLFWvidmode {
	var dm devModeW
	dm.DmSize = uint16(unsafe.Sizeof(dm))
	ret, _, _ := procEnumDisplaySettingsW.Call(0, enumCurrentSettings, uintptr(unsafe.Pointer(&dm)))
	if ret == 0 {
		return GLFWvidmode{Width: 1920, Height: 1080, RedBits: 8, GreenBits: 8, BlueBits: 8, RefreshRate: 60}
	}
	rate := int32(dm.DmDisplayFrequency)
	if rate <= 1 {
		rate = 60
	}
	return GLFWvidmode{
		Width:       int32(dm.DmPelsWidth),
		Height:      int32(dm.DmPelsHeight),
		RedBits:     8,
		GreenBits:   8,
		BlueBits:    8,
		RefreshRate: rate,
	}
}

func platformGetVideoModes(monitor uintptr) []GLFWvidmode {
	var modes []GLFWvidmode
	seen := make(map[[3]int32]bool)
	for i := uint32(0); ; i++ {
		var dm devModeW
		dm.DmSize = uint16(unsafe.Sizeof(dm))
		ret, _, _ := procEnumDisplaySettingsW.Call(0, uintptr(i), uintptr(unsafe.Pointer(&dm)))
		if ret == 0 {
			break
		}
		if dm.DmBitsPerPel < 24 || dm.DmPelsWidth == 0 || dm.DmPelsHeight == 0 {
			continue
		}
		rate := int32(dm.DmDisplayFrequency)
		if rate <= 1 {
			rate = 60
		}
		key := [3]int32{int32(dm.DmPelsWidth), int32(dm.DmPelsHeight), rate}
		if !seen[key] {
			seen[key] = true
			modes = append(modes, GLFWvidmode{
				Width:       key[0],
				Height:      key[1],
				RedBits:     8,
				GreenBits:   8,
				BlueBits:    8,
				RefreshRate: rate,
			})
		}
	}
	if len(modes) == 0 {
		modes = append(modes, platformGetVideoMode(monitor))
	}
	return modes
}

func platformCreateWindow(w *Window, width, height int, title string, monitor uintptr) error {
	style := uint32(wsOverlappedWindow | wsClipChildren | wsClipSiblings)
	exStyle := uint32(wsExAppWindow | wsExWindowEdge | wsExAcceptFiles)
	x := int32(cwUseDefault)
	y := int32(cwUseDefault)
	winW := int32(width)
	winH := int32(height)

	if monitor != 0 {
		style = wsPopup | wsVisible | wsClipChildren | wsClipSiblings
		exStyle = wsExAppWindow | wsExAcceptFiles
		x = 0
		y = 0
		mode := platformGetVideoMode(monitor)
		if width <= 0 {
			winW = mode.Width
		}
		if height <= 0 {
			winH = mode.Height
		}
	} else {
		rect := winRect{Left: 0, Top: 0, Right: int32(width), Bottom: int32(height)}
		_, _, _ = procAdjustWindowRectEx.Call(uintptr(unsafe.Pointer(&rect)), uintptr(style), 0, uintptr(exStyle))
		winW = rect.Right - rect.Left
		winH = rect.Bottom - rect.Top
	}

	titlePtr := utf16Ptr(title)
	hwnd, _, err := procCreateWindowExW.Call(
		uintptr(exStyle),
		uintptr(unsafe.Pointer(winClassName)),
		uintptr(unsafe.Pointer(titlePtr)),
		uintptr(style),
		uintptr(x),
		uintptr(y),
		uintptr(winW),
		uintptr(winH),
		0,
		0,
		winHInstance,
		0,
	)
	if hwnd == 0 {
		return fmt.Errorf("CreateWindowExW failed: %v", err)
	}
	w.platformHandle = hwnd
	w.handle = hwnd
	registerWindow(w)

	_, _, _ = procDragAcceptFiles.Call(hwnd, 1)
	_, _, _ = procShowWindow.Call(hwnd, swShow)
	_, _, _ = procUpdateWindow.Call(hwnd)
	return nil
}

func platformDestroyWindow(w *Window) {
	if w != nil && w.platformHandle != 0 {
		if w.cursorMode == GLFW_CURSOR_DISABLED {
			platformSetCursorMode(w, GLFW_CURSOR_NORMAL)
		}
		_, _, _ = procDestroyWindow.Call(w.platformHandle)
		w.platformHandle = 0
	}
}

func platformPollEvents() {
	var msg winMsg
	for {
		ret, _, _ := procPeekMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0, pmRemove)
		if ret == 0 {
			break
		}
		_, _, _ = procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		_, _, _ = procDispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
	}
}

func platformWaitEventsTimeout(timeout float64) {
	ms := uint32(timeout * 1000.0)
	_, _, _ = procMsgWaitForMultipleObjectsEx.Call(0, 0, uintptr(ms), qsAllInput, mwmoInputAvailable)
	platformPollEvents()
}

func platformSetWindowMonitor(w *Window, monitor uintptr, xpos, ypos, width, height, refreshRate int) {
	if w == nil || w.platformHandle == 0 {
		return
	}
	const gwlStyle = ^uintptr(15) // -16 (GWL_STYLE)
	if monitor != 0 {
		if width > 0 && height > 0 {
			var dm devModeW
			dm.DmSize = uint16(unsafe.Sizeof(dm))
			dm.DmFields = dmPelsWidth | dmPelsHeight
			dm.DmPelsWidth = uint32(width)
			dm.DmPelsHeight = uint32(height)
			if refreshRate > 0 {
				dm.DmFields |= dmDisplayFrequency
				dm.DmDisplayFrequency = uint32(refreshRate)
			}
			_, _, _ = procChangeDisplaySettingsExW.Call(0, uintptr(unsafe.Pointer(&dm)), 0, cdsFullscreen, 0)
		}
		style := uintptr(wsPopup | wsVisible | wsClipChildren | wsClipSiblings)
		_, _, _ = procSetWindowLongPtrW.Call(w.platformHandle, gwlStyle, style)
		_, _, _ = procSetWindowPos.Call(w.platformHandle, hwndTop, 0, 0, uintptr(width), uintptr(height), swpNoOwnerZOrder|swpFrameChanged|swpShowWindow)
	} else {
		_, _, _ = procChangeDisplaySettingsExW.Call(0, 0, 0, 0, 0)
		style := uintptr(wsOverlappedWindow | wsVisible | wsClipChildren | wsClipSiblings)
		_, _, _ = procSetWindowLongPtrW.Call(w.platformHandle, gwlStyle, style)
		rect := winRect{Left: 0, Top: 0, Right: int32(width), Bottom: int32(height)}
		_, _, _ = procAdjustWindowRectEx.Call(uintptr(unsafe.Pointer(&rect)), style, 0, wsExAppWindow|wsExWindowEdge)
		winW := rect.Right - rect.Left
		winH := rect.Bottom - rect.Top
		_, _, _ = procSetWindowPos.Call(w.platformHandle, hwndTop, uintptr(xpos), uintptr(ypos), uintptr(winW), uintptr(winH), swpNoZOrder|swpNoOwnerZOrder|swpFrameChanged|swpShowWindow)
	}
	w.width = width
	w.height = height
	w.fbWidth = width
	w.fbHeight = height
}

func platformGetWindowSize(w *Window) (int, int) {
	if w == nil || w.platformHandle == 0 {
		if w != nil {
			return w.width, w.height
		}
		return 0, 0
	}
	var rc winRect
	ret, _, _ := procGetClientRect.Call(w.platformHandle, uintptr(unsafe.Pointer(&rc)))
	if ret != 0 {
		w.width = int(rc.Right - rc.Left)
		w.height = int(rc.Bottom - rc.Top)
	}
	return w.width, w.height
}

func platformGetFramebufferSize(w *Window) (int, int) {
	ww, hh := platformGetWindowSize(w)
	w.fbWidth = ww
	w.fbHeight = hh
	return ww, hh
}

func platformGetWindowPos(w *Window) (int, int) {
	if w == nil || w.platformHandle == 0 {
		if w != nil {
			return w.xpos, w.ypos
		}
		return 0, 0
	}
	var pt winPoint
	ret, _, _ := procClientToScreen.Call(w.platformHandle, uintptr(unsafe.Pointer(&pt)))
	if ret != 0 {
		w.xpos = int(pt.X)
		w.ypos = int(pt.Y)
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
	if w == nil || w.platformHandle == 0 {
		return
	}
	if mode == GLFW_CURSOR_DISABLED {
		_, _, _ = procShowCursor.Call(0)
		var rc winRect
		_, _, _ = procGetClientRect.Call(w.platformHandle, uintptr(unsafe.Pointer(&rc)))
		tl := winPoint{X: rc.Left, Y: rc.Top}
		br := winPoint{X: rc.Right, Y: rc.Bottom}
		_, _, _ = procClientToScreen.Call(w.platformHandle, uintptr(unsafe.Pointer(&tl)))
		_, _, _ = procClientToScreen.Call(w.platformHandle, uintptr(unsafe.Pointer(&br)))
		clip := winRect{Left: tl.X, Top: tl.Y, Right: br.X, Bottom: br.Y}
		_, _, _ = procClipCursor.Call(uintptr(unsafe.Pointer(&clip)))

		rid := rawInputDevice{
			UsUsagePage: 0x01,
			UsUsage:     0x02,
			DwFlags:     0,
			HwndTarget:  w.platformHandle,
		}
		_, _, _ = procRegisterRawInputDevices.Call(uintptr(unsafe.Pointer(&rid)), 1, unsafe.Sizeof(rid))
	} else {
		rid := rawInputDevice{
			UsUsagePage: 0x01,
			UsUsage:     0x02,
			DwFlags:     ridnevRemove,
			HwndTarget:  0,
		}
		_, _, _ = procRegisterRawInputDevices.Call(uintptr(unsafe.Pointer(&rid)), 1, unsafe.Sizeof(rid))
		_, _, _ = procClipCursor.Call(0)
		_, _, _ = procShowCursor.Call(1)
	}
}

func platformGetCursorPos(w *Window) (float64, float64) {
	if w == nil {
		return 0, 0
	}
	if w.cursorMode == GLFW_CURSOR_DISABLED || w.platformHandle == 0 {
		return w.cursorX, w.cursorY
	}
	var pt winPoint
	ret, _, _ := procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	if ret != 0 {
		_, _, _ = procScreenToClient.Call(w.platformHandle, uintptr(unsafe.Pointer(&pt)))
		w.cursorX = float64(pt.X)
		w.cursorY = float64(pt.Y)
	}
	return w.cursorX, w.cursorY
}

func platformRequiredInstanceExtensions() []string {
	return []string{"VK_KHR_surface", "VK_KHR_win32_surface"}
}

func platformCreateWindowSurface(instance uintptr, w *Window, allocator uintptr, surface *uint64) (int32, error) {
	if err := bindInstance(&VkCreateWin32SurfaceKHR, instance, "vkCreateWin32SurfaceKHR"); err != nil {
		return VK_ERROR_INITIALIZATION_FAILED, err
	}
	info := VkWin32SurfaceCreateInfoKHR{
		SType:     VK_STRUCTURE_TYPE_WIN32_SURFACE_CREATE_INFO_KHR,
		Hinstance: winHInstance,
		Hwnd:      w.platformHandle,
	}
	res := VkCreateWin32SurfaceKHR(instance, &info, allocator, surface)
	if res != VK_SUCCESS {
		return res, fmt.Errorf("vkCreateWin32SurfaceKHR returned %d", res)
	}
	return VK_SUCCESS, nil
}

func winMods() int {
	mods := 0
	if s, _, _ := procGetKeyState.Call(vkShift); int16(s) < 0 {
		mods |= GLFW_MOD_SHIFT
	}
	if s, _, _ := procGetKeyState.Call(vkControl); int16(s) < 0 {
		mods |= GLFW_MOD_CONTROL
	}
	if s, _, _ := procGetKeyState.Call(vkMenu); int16(s) < 0 {
		mods |= GLFW_MOD_ALT
	}
	if s1, _, _ := procGetKeyState.Call(vkLWin); int16(s1) < 0 {
		mods |= GLFW_MOD_SUPER
	} else if s2, _, _ := procGetKeyState.Call(vkRWin); int16(s2) < 0 {
		mods |= GLFW_MOD_SUPER
	}
	if s, _, _ := procGetKeyState.Call(vkCapital); (s & 1) != 0 {
		mods |= GLFW_MOD_CAPS_LOCK
	}
	if s, _, _ := procGetKeyState.Call(vkNumLock); (s & 1) != 0 {
		mods |= GLFW_MOD_NUM_LOCK
	}
	return mods
}

func translateWinKey(vk uintptr) int {
	if vk >= 'A' && vk <= 'Z' {
		return int(vk)
	}
	if vk >= '0' && vk <= '9' {
		return int(vk)
	}
	if vk >= 0x70 && vk <= 0x7B { // VK_F1..VK_F12
		return GLFW_KEY_F1 + int(vk-0x70)
	}
	switch vk {
	case 0x20:
		return GLFW_KEY_SPACE
	case 0x08:
		return GLFW_KEY_BACKSPACE
	case 0x09:
		return GLFW_KEY_TAB
	case 0x0D:
		return GLFW_KEY_ENTER
	case 0x1B:
		return GLFW_KEY_ESCAPE
	case 0x21:
		return GLFW_KEY_PAGE_UP
	case 0x22:
		return GLFW_KEY_PAGE_DOWN
	case 0x23:
		return GLFW_KEY_END
	case 0x24:
		return GLFW_KEY_HOME
	case 0x25:
		return GLFW_KEY_LEFT
	case 0x26:
		return GLFW_KEY_UP
	case 0x27:
		return GLFW_KEY_RIGHT
	case 0x28:
		return GLFW_KEY_DOWN
	case 0x2D:
		return GLFW_KEY_INSERT
	case 0x2E:
		return GLFW_KEY_DELETE
	case 0x10, 0xA0:
		return GLFW_KEY_LEFT_SHIFT
	case 0xA1:
		return GLFW_KEY_RIGHT_SHIFT
	case 0x11, 0xA2:
		return GLFW_KEY_LEFT_CONTROL
	case 0xA3:
		return GLFW_KEY_RIGHT_CONTROL
	case 0x12, 0xA4:
		return GLFW_KEY_LEFT_ALT
	case 0xA5:
		return GLFW_KEY_RIGHT_ALT
	case 0x5B:
		return GLFW_KEY_LEFT_SUPER
	case 0x5C:
		return GLFW_KEY_RIGHT_SUPER
	default:
		return GLFW_KEY_UNKNOWN
	}
}

func winWndProc(hwnd, msg, wParam, lParam uintptr) uintptr {
	w := lookupWindow(hwnd)
	if w == nil {
		ret, _, _ := procDefWindowProcW.Call(hwnd, msg, wParam, lParam)
		return ret
	}
	switch msg {
	case wmClose:
		w.shouldClose = true
		return 0
	case wmDestroy:
		return 0
	case wmMove:
		w.xpos = int(int16(uint16(lParam & 0xFFFF)))
		w.ypos = int(int16(uint16((lParam >> 16) & 0xFFFF)))
		return 0
	case wmSize:
		width := int(uint16(lParam & 0xFFFF))
		height := int(uint16((lParam >> 16) & 0xFFFF))
		w.dispatchWindowSize(width, height)
		w.dispatchFramebufferSize(width, height)
		return 0
	case wmSetFocus:
		w.dispatchFocus(true)
		return 0
	case wmKillFocus:
		w.dispatchFocus(false)
		return 0
	case wmKeyDown, wmSysKeyDown:
		scancode := int((lParam >> 16) & 0x1FF)
		action := GLFW_PRESS
		if (lParam & (1 << 30)) != 0 {
			action = GLFW_REPEAT
		}
		w.dispatchKey(translateWinKey(wParam), scancode, action, winMods())
		return 0
	case wmKeyUp, wmSysKeyUp:
		scancode := int((lParam >> 16) & 0x1FF)
		w.dispatchKey(translateWinKey(wParam), scancode, GLFW_RELEASE, winMods())
		return 0
	case wmChar:
		w.dispatchChar(rune(wParam))
		return 0
	case wmMouseMove:
		if w.cursorMode != GLFW_CURSOR_DISABLED {
			x := float64(int16(uint16(lParam & 0xFFFF)))
			y := float64(int16(uint16((lParam >> 16) & 0xFFFF)))
			w.dispatchCursorPos(x, y)
		}
		return 0
	case wmInput:
		if w.cursorMode == GLFW_CURSOR_DISABLED {
			var raw rawInput
			size := uint32(unsafe.Sizeof(raw))
			ret, _, _ := procGetRawInputData.Call(
				lParam,
				ridInput,
				uintptr(unsafe.Pointer(&raw)),
				uintptr(unsafe.Pointer(&size)),
				unsafe.Sizeof(rawInputHeader{}),
			)
			if ret != ^uintptr(0) && raw.Header.DwType == rimTypeMouse {
				dx := float64(raw.Mouse.LLastX)
				dy := float64(raw.Mouse.LLastY)
				if dx != 0 || dy != 0 {
					w.dispatchRelativeCursorDelta(dx, dy)
				}
			}
		}
		return 0
	case wmLButtonDown:
		w.dispatchMouseButton(GLFW_MOUSE_BUTTON_LEFT, GLFW_PRESS, winMods())
		return 0
	case wmLButtonUp:
		w.dispatchMouseButton(GLFW_MOUSE_BUTTON_LEFT, GLFW_RELEASE, winMods())
		return 0
	case wmRButtonDown:
		w.dispatchMouseButton(GLFW_MOUSE_BUTTON_RIGHT, GLFW_PRESS, winMods())
		return 0
	case wmRButtonUp:
		w.dispatchMouseButton(GLFW_MOUSE_BUTTON_RIGHT, GLFW_RELEASE, winMods())
		return 0
	case wmMButtonDown:
		w.dispatchMouseButton(GLFW_MOUSE_BUTTON_MIDDLE, GLFW_PRESS, winMods())
		return 0
	case wmMButtonUp:
		w.dispatchMouseButton(GLFW_MOUSE_BUTTON_MIDDLE, GLFW_RELEASE, winMods())
		return 0
	case wmMouseWheel:
		delta := float64(int16(uint16((wParam>>16)&0xFFFF))) / 120.0
		w.dispatchScroll(0, delta)
		return 0
	case wmMouseHWheel:
		delta := float64(int16(uint16((wParam>>16)&0xFFFF))) / 120.0
		w.dispatchScroll(delta, 0)
		return 0
	case wmDropFiles:
		hDrop := wParam
		count, _, _ := procDragQueryFileW.Call(hDrop, 0xFFFFFFFF, 0, 0)
		if count > 0 {
			paths := make([]string, 0, count)
			for i := uintptr(0); i < count; i++ {
				length, _, _ := procDragQueryFileW.Call(hDrop, i, 0, 0)
				buf := make([]uint16, length+1)
				_, _, _ = procDragQueryFileW.Call(hDrop, i, uintptr(unsafe.Pointer(&buf[0])), length+1)
				paths = append(paths, syscall.UTF16ToString(buf[:length]))
			}
			w.dispatchDrop(paths)
		}
		_, _, _ = procDragFinish.Call(hDrop)
		return 0
	}
	ret, _, _ := procDefWindowProcW.Call(hwnd, msg, wParam, lParam)
	return ret
}
