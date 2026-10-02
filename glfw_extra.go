package binding

var (
	GlfwGetError               func(**byte) int32
	GlfwWaitEventsTimeout      func(float64)
	GlfwGetWindowPos           func(uintptr, *int32, *int32)
	GlfwGetWindowAttrib        func(uintptr, int32) int32
	GlfwSetWindowFocusCallback func(uintptr, uintptr) uintptr
	GlfwInitVulkanLoader       func(uintptr)
)
