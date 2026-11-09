package binding

const (
	VK_STRUCTURE_TYPE_XLIB_SURFACE_CREATE_INFO_KHR    = 1000004000
	VK_STRUCTURE_TYPE_XCB_SURFACE_CREATE_INFO_KHR     = 1000005000
	VK_STRUCTURE_TYPE_WAYLAND_SURFACE_CREATE_INFO_KHR = 1000006000
	VK_STRUCTURE_TYPE_WIN32_SURFACE_CREATE_INFO_KHR   = 1000009000
	VK_STRUCTURE_TYPE_MACOS_SURFACE_CREATE_INFO_MVK   = 1000123000
	VK_STRUCTURE_TYPE_METAL_SURFACE_CREATE_INFO_EXT   = 1000217000
)

type VkWin32SurfaceCreateInfoKHR struct {
	SType     uint32
	PNext     uintptr
	Flags     uint32
	Hinstance uintptr
	Hwnd      uintptr
}

type VkMetalSurfaceCreateInfoEXT struct {
	SType  uint32
	PNext  uintptr
	Flags  uint32
	PLayer uintptr
}

type VkMacOSSurfaceCreateInfoMVK struct {
	SType uint32
	PNext uintptr
	Flags uint32
	PView uintptr
}

type VkXlibSurfaceCreateInfoKHR struct {
	SType  uint32
	PNext  uintptr
	Flags  uint32
	Dpy    uintptr
	Window uintptr
}

type VkXcbSurfaceCreateInfoKHR struct {
	SType      uint32
	PNext      uintptr
	Flags      uint32
	Connection uintptr
	Window     uint32
}

type VkWaylandSurfaceCreateInfoKHR struct {
	SType   uint32
	PNext   uintptr
	Flags   uint32
	Display uintptr
	Surface uintptr
}

var (
	VkCreateWin32SurfaceKHR   func(instance uintptr, pCreateInfo *VkWin32SurfaceCreateInfoKHR, pAllocator uintptr, pSurface *uint64) int32
	VkCreateMetalSurfaceEXT   func(instance uintptr, pCreateInfo *VkMetalSurfaceCreateInfoEXT, pAllocator uintptr, pSurface *uint64) int32
	VkCreateMacOSSurfaceMVK   func(instance uintptr, pCreateInfo *VkMacOSSurfaceCreateInfoMVK, pAllocator uintptr, pSurface *uint64) int32
	VkCreateXlibSurfaceKHR    func(instance uintptr, pCreateInfo *VkXlibSurfaceCreateInfoKHR, pAllocator uintptr, pSurface *uint64) int32
	VkCreateXcbSurfaceKHR     func(instance uintptr, pCreateInfo *VkXcbSurfaceCreateInfoKHR, pAllocator uintptr, pSurface *uint64) int32
	VkCreateWaylandSurfaceKHR func(instance uintptr, pCreateInfo *VkWaylandSurfaceCreateInfoKHR, pAllocator uintptr, pSurface *uint64) int32
)
