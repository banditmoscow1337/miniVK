package binding

import (
	"math"
	"reflect"
	"runtime"
	"testing"
	"time"
	"unsafe"

	"github.com/ebitengine/purego"
)

func TestAllGlfwFunctionsBound(t *testing.T) {
	if err := InitPuregoBindings(); err != nil {
		t.Logf("InitPuregoBindings returned error (Vulkan ICD optional in headless CI): %v", err)
	}
	fns := map[string]any{
		"GlfwInit":                          GlfwInit,
		"GlfwTerminate":                     GlfwTerminate,
		"GlfwWindowHint":                    GlfwWindowHint,
		"GlfwGetPrimaryMonitor":             GlfwGetPrimaryMonitor,
		"GlfwGetVideoMode":                  GlfwGetVideoMode,
		"GlfwGetVideoModes":                 GlfwGetVideoModes,
		"GlfwCreateWindow":                  GlfwCreateWindow,
		"GlfwDestroyWindow":                 GlfwDestroyWindow,
		"GlfwPollEvents":                    GlfwPollEvents,
		"GlfwGetTime":                       GlfwGetTime,
		"GlfwWindowShouldClose":             GlfwWindowShouldClose,
		"GlfwSetWindowShouldClose":          GlfwSetWindowShouldClose,
		"GlfwSetMonitor":                    GlfwSetMonitor,
		"GlfwGetWindowSize":                 GlfwGetWindowSize,
		"GlfwGetFramebufferSize":            GlfwGetFramebufferSize,
		"GlfwSetInputMode":                  GlfwSetInputMode,
		"GlfwGetCursorPos":                  GlfwGetCursorPos,
		"GlfwSetFramebufferSizeCallback":    GlfwSetFramebufferSizeCallback,
		"GlfwSetKeyCallback":                GlfwSetKeyCallback,
		"GlfwSetCharCallback":               GlfwSetCharCallback,
		"GlfwSetCursorPosCallback":          GlfwSetCursorPosCallback,
		"GlfwSetMouseButtonCallback":        GlfwSetMouseButtonCallback,
		"GlfwSetScrollCallback":             GlfwSetScrollCallback,
		"GlfwSetDropCallback":               GlfwSetDropCallback,
		"GlfwGetError":                      GlfwGetError,
		"GlfwWaitEventsTimeout":             GlfwWaitEventsTimeout,
		"GlfwGetWindowPos":                  GlfwGetWindowPos,
		"GlfwGetWindowAttrib":               GlfwGetWindowAttrib,
		"GlfwSetWindowFocusCallback":        GlfwSetWindowFocusCallback,
		"GlfwInitVulkanLoader":              GlfwInitVulkanLoader,
		"GlfwVulkanSupported":               GlfwVulkanSupported,
		"GlfwGetRequiredInstanceExtensions": GlfwGetRequiredInstanceExtensions,
		"GlfwCreateWindowSurface":           GlfwCreateWindowSurface,
	}
	for name, fn := range fns {
		if reflect.ValueOf(fn).IsNil() {
			t.Fatalf("expected function %s to be bound and non-nil", name)
		}
	}
}

func TestWindowLifecycleAndAllCallbacks(t *testing.T) {
	if GlfwInit() != GLFW_TRUE {
		t.Fatalf("GlfwInit failed: %v", LastGLFWError("GlfwInit"))
	}
	defer GlfwTerminate()

	// Calling GlfwInit a second time should be idempotent
	if GlfwInit() != GLFW_TRUE {
		t.Fatalf("second GlfwInit call failed")
	}

	GlfwWindowHint(GLFW_CLIENT_API, GLFW_NO_API)
	GlfwWindowHint(GLFW_RESIZABLE, GLFW_TRUE)

	t0 := GlfwGetTime()
	time.Sleep(5 * time.Millisecond)
	t1 := GlfwGetTime()
	if t1 <= t0 {
		t.Fatalf("expected GlfwGetTime to advance: t0=%f t1=%f", t0, t1)
	}

	mon := GlfwGetPrimaryMonitor()
	if mon == 0 {
		t.Fatalf("expected non-zero primary monitor handle")
	}
	mode := GetVideoMode(mon)
	if mode == nil || mode.Width <= 0 || mode.Height <= 0 || mode.RefreshRate <= 0 {
		t.Fatalf("unexpected primary video mode: %+v", mode)
	}
	modes := GetVideoModes(mon)
	if len(modes) == 0 {
		t.Fatalf("expected at least one video mode")
	}

	// Verify invalid dimensions trigger GLFW_INVALID_VALUE
	if _, err := CreateWindow(0, -10, "Invalid", 0, 0); err == nil {
		t.Fatalf("expected error creating window with invalid dimensions")
	}
	if h := GlfwCreateWindow(-1, 0, nil, 0, 0); h != 0 {
		t.Fatalf("expected 0 handle from GlfwCreateWindow with invalid dimensions")
	}

	win, err := CreateWindow(800, 600, "miniVK PureGo Test", 0, 0)
	if err != nil {
		t.Fatalf("CreateWindow failed: %v", err)
	}
	defer win.Destroy()

	if win.Handle() == 0 {
		t.Fatalf("expected non-zero window handle")
	}
	w, h := win.Size()
	if w != 800 || h != 600 {
		t.Fatalf("expected window size 800x600, got %dx%d", w, h)
	}
	var cw, ch int32
	GlfwGetWindowSize(win.Handle(), &cw, &ch)
	if cw != 800 || ch != 600 {
		t.Fatalf("GlfwGetWindowSize mismatch: %dx%d", cw, ch)
	}

	fw, fh := win.FramebufferSize()
	if fw <= 0 || fh <= 0 {
		t.Fatalf("expected positive framebuffer size, got %dx%d", fw, fh)
	}
	var cfw, cfh int32
	GlfwGetFramebufferSize(win.Handle(), &cfw, &cfh)
	if int(cfw) != fw || int(cfh) != fh {
		t.Fatalf("GlfwGetFramebufferSize mismatch: %dx%d vs %dx%d", cfw, cfh, fw, fh)
	}

	// ShouldClose / SetShouldClose
	if win.ShouldClose() {
		t.Fatalf("new window should not be marked ShouldClose")
	}
	win.SetShouldClose(true)
	if !win.ShouldClose() || GlfwWindowShouldClose(win.Handle()) != GLFW_TRUE {
		t.Fatalf("expected ShouldClose=true after SetShouldClose(true)")
	}
	GlfwSetWindowShouldClose(win.Handle(), GLFW_FALSE)
	if win.ShouldClose() {
		t.Fatalf("expected ShouldClose=false after GlfwSetWindowShouldClose(GLFW_FALSE)")
	}

	// SetMonitor & Position
	win.SetMonitor(0, 120, 140, 1024, 768, 60)
	px, py := win.Position()
	if px != 120 || py != 140 {
		t.Fatalf("expected window position (120, 140), got (%d, %d)", px, py)
	}
	var cpx, cpy int32
	GlfwGetWindowPos(win.Handle(), &cpx, &cpy)
	if cpx != 120 || cpy != 140 {
		t.Fatalf("GlfwGetWindowPos mismatch: (%d, %d)", cpx, cpy)
	}
	w, h = win.Size()
	if w != 1024 || h != 768 {
		t.Fatalf("expected updated window size 1024x768, got %dx%d", w, h)
	}

	// InputMode & CursorPos
	win.SetInputMode(GLFW_CURSOR, GLFW_CURSOR_DISABLED)
	win.SetInputMode(GLFW_CURSOR, GLFW_CURSOR_NORMAL)

	// High-level Window callbacks
	var (
		gotFBW, gotFBH                     int
		gotKey, gotScancode, gotAct, gotMd int
		gotChar                            rune
		gotCurX, gotCurY                   float64
		gotBtn, gotBtnAct, gotBtnMods      int
		gotScrollX, gotScrollY             float64
		gotDrop                            []string
		gotFocus                           bool
	)
	win.SetFramebufferSizeCallback(func(w *Window, width, height int) {
		gotFBW, gotFBH = width, height
	})
	win.SetKeyCallback(func(w *Window, key, scancode, action, mods int) {
		gotKey, gotScancode, gotAct, gotMd = key, scancode, action, mods
	})
	win.SetCharCallback(func(w *Window, character rune) {
		gotChar = character
	})
	if err := win.SetCursorPosCallback(func(w *Window, x, y float64) {
		gotCurX, gotCurY = x, y
	}); err != nil {
		t.Fatalf("SetCursorPosCallback error: %v", err)
	}
	win.SetMouseButtonCallback(func(w *Window, button, action, mods int) {
		gotBtn, gotBtnAct, gotBtnMods = button, action, mods
	})
	if err := win.SetScrollCallback(func(w *Window, x, y float64) {
		gotScrollX, gotScrollY = x, y
	}); err != nil {
		t.Fatalf("SetScrollCallback error: %v", err)
	}
	win.SetDropCallback(func(w *Window, paths []string) {
		gotDrop = append([]string(nil), paths...)
	})
	win.SetFocusCallback(func(focused bool) {
		gotFocus = focused
	})

	win.dispatchFramebufferSize(1280, 720)
	if gotFBW != 1280 || gotFBH != 720 {
		t.Fatalf("FramebufferSizeCallback mismatch: got %dx%d", gotFBW, gotFBH)
	}

	// Key press + automatic repeat detection + release
	win.dispatchKey(GLFW_KEY_W, 17, GLFW_PRESS, GLFW_MOD_SHIFT)
	if gotKey != GLFW_KEY_W || gotScancode != 17 || gotAct != GLFW_PRESS || gotMd != GLFW_MOD_SHIFT {
		t.Fatalf("KeyCallback press mismatch: key=%d sc=%d act=%d mods=%d", gotKey, gotScancode, gotAct, gotMd)
	}
	win.dispatchKey(GLFW_KEY_W, 17, GLFW_PRESS, GLFW_MOD_SHIFT)
	if gotAct != GLFW_REPEAT {
		t.Fatalf("expected repeated key press to report GLFW_REPEAT, got %d", gotAct)
	}
	win.dispatchKey(GLFW_KEY_W, 17, GLFW_RELEASE, 0)
	if gotAct != GLFW_RELEASE {
		t.Fatalf("expected key release to report GLFW_RELEASE, got %d", gotAct)
	}

	win.dispatchChar('Ж')
	if gotChar != 'Ж' {
		t.Fatalf("CharCallback mismatch: got %q", gotChar)
	}

	win.dispatchCursorPos(321.5, 123.25)
	cx, cy := win.GetCursorPos()
	if math.Abs(gotCurX-321.5) > 1e-6 || math.Abs(gotCurY-123.25) > 1e-6 || math.Abs(cx-321.5) > 1e-6 || math.Abs(cy-123.25) > 1e-6 {
		t.Fatalf("CursorPos mismatch: cb=(%f,%f) get=(%f,%f)", gotCurX, gotCurY, cx, cy)
	}
	var ccx, ccy float64
	GlfwGetCursorPos(win.Handle(), &ccx, &ccy)
	if math.Abs(ccx-321.5) > 1e-6 || math.Abs(ccy-123.25) > 1e-6 {
		t.Fatalf("GlfwGetCursorPos mismatch: (%f,%f)", ccx, ccy)
	}
	win.dispatchRelativeCursorDelta(10.5, -3.25)
	cx, cy = win.GetCursorPos()
	if math.Abs(cx-332.0) > 1e-6 || math.Abs(cy-120.0) > 1e-6 {
		t.Fatalf("RelativeCursorDelta mismatch: got (%f,%f)", cx, cy)
	}

	win.dispatchMouseButton(GLFW_MOUSE_BUTTON_RIGHT, GLFW_PRESS, GLFW_MOD_CONTROL)
	if gotBtn != GLFW_MOUSE_BUTTON_RIGHT || gotBtnAct != GLFW_PRESS || gotBtnMods != GLFW_MOD_CONTROL {
		t.Fatalf("MouseButtonCallback mismatch: btn=%d act=%d mods=%d", gotBtn, gotBtnAct, gotBtnMods)
	}

	win.dispatchScroll(1.5, -2.5)
	if math.Abs(gotScrollX-1.5) > 1e-6 || math.Abs(gotScrollY+2.5) > 1e-6 {
		t.Fatalf("ScrollCallback mismatch: (%f,%f)", gotScrollX, gotScrollY)
	}

	win.dispatchDrop([]string{"/tmp/a.spv", "/tmp/b.png"})
	if len(gotDrop) != 2 || gotDrop[0] != "/tmp/a.spv" || gotDrop[1] != "/tmp/b.png" {
		t.Fatalf("DropCallback mismatch: %+v", gotDrop)
	}

	win.dispatchFocus(false)
	if gotFocus || win.Focused() || GlfwGetWindowAttrib(win.Handle(), GLFW_FOCUSED) != GLFW_FALSE {
		t.Fatalf("FocusCallback(false) mismatch")
	}
	win.dispatchFocus(true)
	if !gotFocus || !win.Focused() || GlfwGetWindowAttrib(win.Handle(), GLFW_FOCUSED) != GLFW_TRUE {
		t.Fatalf("FocusCallback(true) mismatch")
	}

	// Verify raw C-callback pointers across all GlfwSet*Callback functions
	var (
		rawFBW, rawFBH int32
		rawKeyHit      int32
		rawCharHit     uint32
		rawBtnHit      int32
		rawFocusHit    int32
		rawDropped     []string
	)
	rawFBCB := purego.NewCallback(func(w uintptr, width, height int32) uintptr {
		rawFBW, rawFBH = width, height
		return 0
	})
	GlfwSetFramebufferSizeCallback(win.Handle(), rawFBCB)
	win.dispatchFramebufferSize(1600, 900)
	if rawFBW != 1600 || rawFBH != 900 {
		t.Fatalf("GlfwSetFramebufferSizeCallback raw C callback failed: %dx%d", rawFBW, rawFBH)
	}

	rawKeyCB := purego.NewCallback(func(w uintptr, key, scancode, action, mods int32) uintptr {
		rawKeyHit = key
		return 0
	})
	GlfwSetKeyCallback(win.Handle(), rawKeyCB)
	win.dispatchKey(GLFW_KEY_ESCAPE, 9, GLFW_PRESS, 0)
	if rawKeyHit != GLFW_KEY_ESCAPE {
		t.Fatalf("GlfwSetKeyCallback raw C callback failed: got %d", rawKeyHit)
	}

	rawCharCB := purego.NewCallback(func(w uintptr, cp uint32) uintptr {
		rawCharHit = cp
		return 0
	})
	GlfwSetCharCallback(win.Handle(), rawCharCB)
	win.dispatchChar('A')
	if rawCharHit != 'A' {
		t.Fatalf("GlfwSetCharCallback raw C callback failed: got %d", rawCharHit)
	}

	rawMouseCB := purego.NewCallback(func(w uintptr, button, action, mods int32) uintptr {
		rawBtnHit = button
		return 0
	})
	GlfwSetMouseButtonCallback(win.Handle(), rawMouseCB)
	win.dispatchMouseButton(GLFW_MOUSE_BUTTON_MIDDLE, GLFW_PRESS, 0)
	if rawBtnHit != GLFW_MOUSE_BUTTON_MIDDLE {
		t.Fatalf("GlfwSetMouseButtonCallback raw C callback failed: got %d", rawBtnHit)
	}

	rawFocusCB := purego.NewCallback(func(w uintptr, focused int32) uintptr {
		rawFocusHit = focused
		return 0
	})
	GlfwSetWindowFocusCallback(win.Handle(), rawFocusCB)
	win.dispatchFocus(true)
	if rawFocusHit != GLFW_TRUE {
		t.Fatalf("GlfwSetWindowFocusCallback raw C callback failed: got %d", rawFocusHit)
	}

	rawDropCB := purego.NewCallback(func(w uintptr, count int32, rawPaths **byte) uintptr {
		ptrs := unsafe.Slice(rawPaths, int(count))
		for _, p := range ptrs {
			rawDropped = append(rawDropped, String(p))
		}
		return 0
	})
	GlfwSetDropCallback(win.Handle(), rawDropCB)
	win.dispatchDrop([]string{"/path/one", "/path/two"})
	if len(rawDropped) != 2 || rawDropped[0] != "/path/one" || rawDropped[1] != "/path/two" {
		t.Fatalf("GlfwSetDropCallback raw C callback failed: %+v", rawDropped)
	}

	// Test C-style GlfwCreateWindow and GlfwDestroyWindow directly
	cTitle := CString("C-Style Window")
	h2 := GlfwCreateWindow(320, 240, &cTitle[0], 0, 0)
	runtime.KeepAlive(cTitle)
	if h2 == 0 {
		t.Fatalf("GlfwCreateWindow failed: %v", LastGLFWError("GlfwCreateWindow"))
	}
	GlfwDestroyWindow(h2)

	// PollEvents and WaitEventsTimeout
	GlfwPollEvents()
	GlfwWaitEventsTimeout(0.005)

	// Required Vulkan extensions
	var extCount uint32
	extPtrs := GlfwGetRequiredInstanceExtensions(&extCount)
	if extPtrs == nil || extCount < 2 {
		t.Fatalf("expected at least 2 required Vulkan instance extensions, got %d", extCount)
	}
	extSlice := unsafe.Slice(extPtrs, int(extCount))
	gotSurfaceExt := false
	for _, p := range extSlice {
		name := String(p)
		if name == "VK_KHR_surface" {
			gotSurfaceExt = true
		}
	}
	if !gotSurfaceExt {
		t.Fatalf("expected VK_KHR_surface in GlfwGetRequiredInstanceExtensions")
	}
}
