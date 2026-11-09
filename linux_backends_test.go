//go:build linux

package binding

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"
	"unsafe"

	"github.com/ebitengine/purego"
)

// TestLinuxWaylandLiveServer spins up a real in-process Wayland compositor via
// libwayland-server.so.0 using purego and exercises the Wayland client backend
// (linux_wayland.go) + VK_KHR_wayland_surface over a live Unix domain socket.
func TestLinuxWaylandLiveServer(t *testing.T) {
	srvLib, err := loadLibrary([]string{"libwayland-server.so.0", "libwayland-server.so"})
	if err != nil {
		t.Skipf("libwayland-server.so.0 not available: %v", err)
	}
	defer closeLibrary(srvLib)

	var (
		wlDisplayCreate         func() uintptr
		wlDisplayDestroy        func(dpy uintptr)
		wlDisplayAddSocketAuto  func(dpy uintptr) *byte
		wlDisplayGetEventLoop   func(dpy uintptr) uintptr
		wlEventLoopDispatch     func(loop uintptr, timeout int32) int32
		wlDisplayFlushClients   func(dpy uintptr)
		wlGlobalCreate          func(dpy uintptr, iface uintptr, version int32, data uintptr, bind uintptr) uintptr
		wlResourceCreate        func(client uintptr, iface uintptr, version int32, id uint32) uintptr
		wlResourceSetDispatcher func(resource uintptr, dispatcher uintptr, impl uintptr, data uintptr, destroy uintptr)
	)
	for _, b := range []struct {
		fn   any
		name string
	}{
		{&wlDisplayCreate, "wl_display_create"},
		{&wlDisplayDestroy, "wl_display_destroy"},
		{&wlDisplayAddSocketAuto, "wl_display_add_socket_auto"},
		{&wlDisplayGetEventLoop, "wl_display_get_event_loop"},
		{&wlEventLoopDispatch, "wl_event_loop_dispatch"},
		{&wlDisplayFlushClients, "wl_display_flush_clients"},
		{&wlGlobalCreate, "wl_global_create"},
		{&wlResourceCreate, "wl_resource_create"},
		{&wlResourceSetDispatcher, "wl_resource_set_dispatcher"},
	} {
		if err := bindLibrary(b.fn, srvLib, b.name); err != nil {
			t.Fatalf("bind server %s: %v", b.name, err)
		}
	}

	srvCompIface, err := lookupLibrary(srvLib, "wl_compositor_interface")
	if err != nil {
		t.Fatalf("lookup wl_compositor_interface: %v", err)
	}
	srvSurfIface, err := lookupLibrary(srvLib, "wl_surface_interface")
	if err != nil {
		t.Fatalf("lookup wl_surface_interface: %v", err)
	}

	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	srvDpy := wlDisplayCreate()
	if srvDpy == 0 {
		t.Fatalf("wl_display_create failed")
	}
	defer wlDisplayDestroy(srvDpy)

	var lastClient uintptr
	surfDisp := purego.NewCallback(func(impl, res uintptr, opcode uint32, msg uintptr, args *uintptr) uintptr {
		return 0
	})
	compDisp := purego.NewCallback(func(impl, res uintptr, opcode uint32, msg uintptr, rawArgs *uintptr) uintptr {
		if opcode == 0 && rawArgs != nil && lastClient != 0 { // create_surface(id)
			a := unsafe.Slice(rawArgs, 2)
			newID := uint32(a[0])
			sRes := wlResourceCreate(lastClient, srvSurfIface, 4, newID)
			if sRes != 0 {
				wlResourceSetDispatcher(sRes, surfDisp, 0, lastClient, 0)
			}
		}
		return 0
	})
	compBind := purego.NewCallback(func(client, data uintptr, version, id uint32) uintptr {
		lastClient = client
		res := wlResourceCreate(client, srvCompIface, int32(version), id)
		if res != 0 {
			wlResourceSetDispatcher(res, compDisp, 0, client, 0)
		}
		return 0
	})
	wlGlobalCreate(srvDpy, srvCompIface, 4, 0, compBind)

	sockNamePtr := wlDisplayAddSocketAuto(srvDpy)
	if sockNamePtr == nil {
		t.Fatalf("wl_display_add_socket_auto failed")
	}
	sockName := String(sockNamePtr)
	t.Setenv("WAYLAND_DISPLAY", sockName)

	stopCh := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		loop := wlDisplayGetEventLoop(srvDpy)
		for {
			select {
			case <-stopCh:
				return
			default:
				wlEventLoopDispatch(loop, 5)
				wlDisplayFlushClients(srvDpy)
			}
		}
	}()
	defer func() {
		close(stopCh)
		wg.Wait()
	}()

	SetLinuxBackend("wayland")
	defer SetLinuxBackend("auto")

	if GlfwInit() != GLFW_TRUE {
		t.Fatalf("GlfwInit(wayland) failed: %v", LastGLFWError("wayland init"))
	}
	defer GlfwTerminate()

	if ActiveLinuxBackend() != "wayland" {
		t.Fatalf("expected ActiveLinuxBackend() == wayland, got %s", ActiveLinuxBackend())
	}

	var extCount uint32
	extPtr := GlfwGetRequiredInstanceExtensions(&extCount)
	if extPtr == nil || extCount != 2 {
		t.Fatalf("expected 2 Wayland Vulkan extensions, got %d", extCount)
	}
	exts := unsafe.Slice(extPtr, int(extCount))
	if String(exts[0]) != "VK_KHR_surface" || String(exts[1]) != "VK_KHR_wayland_surface" {
		t.Fatalf("unexpected Wayland extensions: %s, %s", String(exts[0]), String(exts[1]))
	}

	win, err := CreateWindow(640, 480, "Wayland PureGo Test", 0, 0)
	if err != nil {
		t.Fatalf("CreateWindow(wayland) failed: %v", err)
	}
	w, h := win.Size()
	if w != 640 || h != 480 {
		t.Fatalf("expected 640x480, got %dx%d", w, h)
	}

	// Test VkSurfaceKHR creation via GlfwCreateWindowSurface if Vulkan ICD is available
	if InitVulkanBindings() == nil {
		appInfo := VkApplicationInfo{
			SType:      VK_STRUCTURE_TYPE_APPLICATION_INFO,
			ApiVersion: VkMakeVersion(1, 0, 0),
		}
		createInfo := VkInstanceCreateInfo{
			SType:                   VK_STRUCTURE_TYPE_INSTANCE_CREATE_INFO,
			PApplicationInfo:        &appInfo,
			EnabledExtensionCount:   extCount,
			PpEnabledExtensionNames: extPtr,
		}
		var instance uintptr
		if VkCreateInstance(&createInfo, 0, &instance) == VK_SUCCESS && instance != 0 {
			_ = LoadInstance(instance)
			var surface uint64
			if res := GlfwCreateWindowSurface(instance, win.Handle(), 0, &surface); res != VK_SUCCESS || surface == 0 {
				t.Fatalf("GlfwCreateWindowSurface(wayland) failed with VkResult=%d", res)
			}
			if VkDestroySurfaceKHR != nil {
				VkDestroySurfaceKHR(instance, surface, 0)
			}
			if VkDestroyInstance != nil {
				VkDestroyInstance(instance, 0)
			}
		}
	}

	GlfwPollEvents()
	GlfwWaitEventsTimeout(0.005)
	win.Destroy()
	runtime.KeepAlive(compBind)
	runtime.KeepAlive(compDisp)
	runtime.KeepAlive(surfDisp)
}

// TestLinuxX11WireProtocolServer spins up a minimal in-process X11 protocol server
// listening on /tmp/.X11-unix/X<display> so that real libX11.so.6 (XOpenDisplay,
// XInternAtom, XCreateWindow, XMapWindow, XStoreName, XDestroyWindow, XCloseDisplay)
// executes against our purego X11 backend even without an external Xorg server.
func TestLinuxX11WireProtocolServer(t *testing.T) {
	if err := os.MkdirAll("/tmp/.X11-unix", 0777); err != nil {
		t.Skipf("cannot create /tmp/.X11-unix: %v", err)
	}
	const displayNum = 87
	sockPath := fmt.Sprintf("/tmp/.X11-unix/X%d", displayNum)
	_ = os.Remove(sockPath)
	ln, err := net.Listen("unix", sockPath)
	if err != nil {
		t.Skipf("cannot listen on %s: %v", sockPath, err)
	}
	defer func() {
		_ = ln.Close()
		_ = os.Remove(sockPath)
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		serveMockX11Client(conn)
	}()

	t.Setenv("DISPLAY", fmt.Sprintf(":%d", displayNum))
	t.Setenv("XAUTHORITY", filepath.Join(t.TempDir(), "nonexistent_xauth"))
	SetLinuxBackend("x11")
	defer SetLinuxBackend("auto")

	if GlfwInit() != GLFW_TRUE {
		t.Fatalf("GlfwInit(x11) failed: %v", LastGLFWError("x11 init"))
	}
	defer GlfwTerminate()

	if ActiveLinuxBackend() != "x11" {
		t.Fatalf("expected ActiveLinuxBackend() == x11, got %s", ActiveLinuxBackend())
	}

	win, err := CreateWindow(640, 480, "X11 PureGo Test", 0, 0)
	if err != nil {
		t.Fatalf("CreateWindow(x11) failed: %v", err)
	}
	w, h := win.Size()
	if w != 640 || h != 480 {
		t.Fatalf("expected 640x480, got %dx%d", w, h)
	}
	px, py := win.Position()
	if px != 50 || py != 60 {
		t.Fatalf("expected translated position (50, 60), got (%d, %d)", px, py)
	}
	cx, cy := win.GetCursorPos()
	if cx != 120 || cy != 90 {
		t.Fatalf("expected queried cursor pos (120, 90), got (%f, %f)", cx, cy)
	}
	GlfwPollEvents()
	GlfwWaitEventsTimeout(0.005)
	win.Destroy()
	GlfwTerminate()
	wg.Wait()
}

func serveMockX11Client(conn net.Conn) {
	_ = conn.SetDeadline(time.Now().Add(5 * time.Second))
	hdr := make([]byte, 12)
	if _, err := io.ReadFull(conn, hdr); err != nil {
		return
	}
	nAuthName := int(binary.LittleEndian.Uint16(hdr[6:8]))
	nAuthData := int(binary.LittleEndian.Uint16(hdr[8:10]))
	pad4 := func(n int) int { return (n + 3) & ^3 }
	extra := pad4(nAuthName) + pad4(nAuthData)
	if extra > 0 {
		buf := make([]byte, extra)
		if _, err := io.ReadFull(conn, buf); err != nil {
			return
		}
	}

	// Build X11 Connection Setup Success reply (1 screen, 1 depth=24, 1 visual)
	// Header (8 bytes) + body (116 bytes = 29 words)
	body := make([]byte, 116)
	binary.LittleEndian.PutUint32(body[0:4], 1)          // release-number
	binary.LittleEndian.PutUint32(body[4:8], 0x00400000) // resource-id-base
	binary.LittleEndian.PutUint32(body[8:12], 0x001fffff) // resource-id-mask
	binary.LittleEndian.PutUint32(body[12:16], 256)       // motion-buffer-size
	binary.LittleEndian.PutUint16(body[16:18], 4)         // vendor-length ("mini")
	binary.LittleEndian.PutUint16(body[18:20], 0xffff)    // max-request-length
	body[20] = 1                                          // number-of-screens
	body[21] = 1                                          // number-of-pixmap-formats
	body[22] = 0                                          // image-byte-order
	body[23] = 0                                          // bitmap-format-bit-order
	body[24] = 32                                         // bitmap-format-scanline-unit
	body[25] = 32                                         // bitmap-format-scanline-pad
	body[26] = 8                                          // min-keycode
	body[27] = 255                                        // max-keycode
	copy(body[32:36], []byte("mini"))
	// Pixmap format (8 bytes at offset 36)
	body[36] = 24
	body[37] = 32
	body[38] = 32
	// Screen (40 bytes at offset 44)
	s := body[44:84]
	binary.LittleEndian.PutUint32(s[0:4], 0x100)  // root window
	binary.LittleEndian.PutUint32(s[4:8], 0x101)  // default-colormap
	binary.LittleEndian.PutUint32(s[8:12], 0xffffff)
	binary.LittleEndian.PutUint32(s[12:16], 0x000000)
	binary.LittleEndian.PutUint16(s[20:22], 1920) // width-in-pixels
	binary.LittleEndian.PutUint16(s[22:24], 1080) // height-in-pixels
	binary.LittleEndian.PutUint16(s[24:26], 508)  // width-in-millimeters
	binary.LittleEndian.PutUint16(s[26:28], 285)  // height-in-millimeters
	binary.LittleEndian.PutUint16(s[28:30], 1)    // min-installed-maps
	binary.LittleEndian.PutUint16(s[30:32], 1)    // max-installed-maps
	binary.LittleEndian.PutUint32(s[32:36], 0x21) // root-visual
	s[36] = 0
	s[37] = 0
	s[38] = 24 // root-depth
	s[39] = 1  // number-of-allowed-depths
	// Depth (8 bytes at offset 84) + VisualType (24 bytes at offset 92)
	d := body[84:92]
	d[0] = 24
	binary.LittleEndian.PutUint16(d[2:4], 1) // visuals-len = 1
	v := body[92:116]
	binary.LittleEndian.PutUint32(v[0:4], 0x21) // visual-id
	v[4] = 4                                    // TrueColor
	v[5] = 8                                    // bits-per-rgb
	binary.LittleEndian.PutUint16(v[6:8], 256)  // colormap-entries
	binary.LittleEndian.PutUint32(v[8:12], 0xff0000)
	binary.LittleEndian.PutUint32(v[12:16], 0x00ff00)
	binary.LittleEndian.PutUint32(v[16:20], 0x0000ff)

	setupHdr := make([]byte, 8)
	setupHdr[0] = 1 // Success
	binary.LittleEndian.PutUint16(setupHdr[2:4], 11)
	binary.LittleEndian.PutUint16(setupHdr[4:6], 0)
	binary.LittleEndian.PutUint16(setupHdr[6:8], uint16(len(body)/4))
	_, _ = conn.Write(setupHdr)
	_, _ = conn.Write(body)

	var seq uint16
	var nextAtom uint32 = 100
	reqHdr := make([]byte, 4)
	for {
		if _, err := io.ReadFull(conn, reqHdr); err != nil {
			return
		}
		seq++
		opcode := reqHdr[0]
		reqWords := int(binary.LittleEndian.Uint16(reqHdr[2:4]))
		if reqWords > 1 {
			payload := make([]byte, (reqWords-1)*4)
			if _, err := io.ReadFull(conn, payload); err != nil {
				return
			}
		}
		switch opcode {
		case 14: // GetGeometry
			reply := make([]byte, 32)
			reply[0] = 1
			reply[1] = 24
			binary.LittleEndian.PutUint16(reply[2:4], seq)
			binary.LittleEndian.PutUint32(reply[8:12], 0x100)
			binary.LittleEndian.PutUint16(reply[16:18], 640)
			binary.LittleEndian.PutUint16(reply[18:20], 480)
			_, _ = conn.Write(reply)
		case 3: // GetWindowAttributes
			reply := make([]byte, 44)
			reply[0] = 1
			reply[1] = 2
			binary.LittleEndian.PutUint16(reply[2:4], seq)
			binary.LittleEndian.PutUint32(reply[4:8], 3) // 3 extra words
			binary.LittleEndian.PutUint32(reply[8:12], 0x21)
			binary.LittleEndian.PutUint16(reply[12:14], 1)
			reply[26] = 2 // IsViewable
			_, _ = conn.Write(reply)
		case 16: // InternAtom
			nextAtom++
			reply := make([]byte, 32)
			reply[0] = 1
			binary.LittleEndian.PutUint16(reply[2:4], seq)
			binary.LittleEndian.PutUint32(reply[8:12], nextAtom)
			_, _ = conn.Write(reply)
		case 20: // GetProperty
			reply := make([]byte, 32)
			reply[0] = 1
			binary.LittleEndian.PutUint16(reply[2:4], seq)
			_, _ = conn.Write(reply)
		case 38: // QueryPointer
			reply := make([]byte, 32)
			reply[0] = 1
			reply[1] = 1
			binary.LittleEndian.PutUint16(reply[2:4], seq)
			binary.LittleEndian.PutUint32(reply[8:12], 0x100)
			binary.LittleEndian.PutUint16(reply[20:22], 120)
			binary.LittleEndian.PutUint16(reply[22:24], 90)
			_, _ = conn.Write(reply)
		case 40: // TranslateCoordinates
			reply := make([]byte, 32)
			reply[0] = 1
			reply[1] = 1
			binary.LittleEndian.PutUint16(reply[2:4], seq)
			binary.LittleEndian.PutUint16(reply[12:14], 50)
			binary.LittleEndian.PutUint16(reply[14:16], 60)
			_, _ = conn.Write(reply)
		case 43: // GetInputFocus
			reply := make([]byte, 32)
			reply[0] = 1
			reply[1] = 1
			binary.LittleEndian.PutUint16(reply[2:4], seq)
			binary.LittleEndian.PutUint32(reply[8:12], 0x00400001)
			_, _ = conn.Write(reply)
		case 98: // QueryExtension
			reply := make([]byte, 32)
			reply[0] = 1
			binary.LittleEndian.PutUint16(reply[2:4], seq)
			reply[8] = 0 // extension not present on mock server
			_, _ = conn.Write(reply)
		}
	}
}
