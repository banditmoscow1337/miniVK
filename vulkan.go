package binding

import (
	"fmt"
	"os"
	"runtime"
	"unsafe"
)

const (
	VK_SUCCESS                                       = 0
	VK_NOT_READY                                     = 1
	VK_TIMEOUT                                       = 2
	VK_INCOMPLETE                                    = 5
	VK_SUBOPTIMAL_KHR                                = 1000001003
	VK_ERROR_OUT_OF_DATE_KHR                         = -1000001004
	VK_ERROR_SURFACE_LOST_KHR                        = -1000000000
	VK_ERROR_INITIALIZATION_FAILED                   = -3
	VK_INSTANCE_CREATE_ENUMERATE_PORTABILITY_BIT_KHR = 0x00000001

	VK_STRUCTURE_TYPE_APPLICATION_INFO                          = 0
	VK_STRUCTURE_TYPE_INSTANCE_CREATE_INFO                      = 1
	VK_STRUCTURE_TYPE_DEVICE_QUEUE_CREATE_INFO                  = 2
	VK_STRUCTURE_TYPE_DEVICE_CREATE_INFO                        = 3
	VK_STRUCTURE_TYPE_SUBMIT_INFO                               = 4
	VK_STRUCTURE_TYPE_MEMORY_ALLOCATE_INFO                      = 5
	VK_STRUCTURE_TYPE_MAPPED_MEMORY_RANGE                       = 6
	VK_STRUCTURE_TYPE_BIND_SPARSE_INFO                          = 7
	VK_STRUCTURE_TYPE_FENCE_CREATE_INFO                         = 8
	VK_STRUCTURE_TYPE_SEMAPHORE_CREATE_INFO                     = 9
	VK_STRUCTURE_TYPE_EVENT_CREATE_INFO                         = 10
	VK_STRUCTURE_TYPE_QUERY_POOL_CREATE_INFO                    = 11
	VK_STRUCTURE_TYPE_BUFFER_CREATE_INFO                        = 12
	VK_STRUCTURE_TYPE_BUFFER_VIEW_CREATE_INFO                   = 13
	VK_STRUCTURE_TYPE_IMAGE_CREATE_INFO                         = 14
	VK_STRUCTURE_TYPE_IMAGE_VIEW_CREATE_INFO                    = 15
	VK_STRUCTURE_TYPE_SHADER_MODULE_CREATE_INFO                 = 16
	VK_STRUCTURE_TYPE_PIPELINE_CACHE_CREATE_INFO                = 17
	VK_STRUCTURE_TYPE_PIPELINE_SHADER_STAGE_CREATE_INFO         = 18
	VK_STRUCTURE_TYPE_PIPELINE_VERTEX_INPUT_STATE_CREATE_INFO   = 19
	VK_STRUCTURE_TYPE_PIPELINE_INPUT_ASSEMBLY_STATE_CREATE_INFO = 20
	VK_STRUCTURE_TYPE_PIPELINE_VIEWPORT_STATE_CREATE_INFO       = 22
	VK_STRUCTURE_TYPE_PIPELINE_RASTERIZATION_STATE_CREATE_INFO  = 23
	VK_STRUCTURE_TYPE_PIPELINE_MULTISAMPLE_STATE_CREATE_INFO    = 24
	VK_STRUCTURE_TYPE_PIPELINE_DEPTH_STENCIL_STATE_CREATE_INFO  = 25
	VK_STRUCTURE_TYPE_PIPELINE_COLOR_BLEND_STATE_CREATE_INFO    = 26
	VK_STRUCTURE_TYPE_PIPELINE_DYNAMIC_STATE_CREATE_INFO        = 27
	VK_STRUCTURE_TYPE_GRAPHICS_PIPELINE_CREATE_INFO             = 28
	VK_STRUCTURE_TYPE_COMPUTE_PIPELINE_CREATE_INFO              = 29
	VK_STRUCTURE_TYPE_PIPELINE_LAYOUT_CREATE_INFO               = 30
	VK_STRUCTURE_TYPE_SAMPLER_CREATE_INFO                       = 31
	VK_STRUCTURE_TYPE_DESCRIPTOR_SET_LAYOUT_CREATE_INFO         = 32
	VK_STRUCTURE_TYPE_DESCRIPTOR_POOL_CREATE_INFO               = 33
	VK_STRUCTURE_TYPE_DESCRIPTOR_SET_ALLOCATE_INFO              = 34
	VK_STRUCTURE_TYPE_WRITE_DESCRIPTOR_SET                      = 35
	VK_STRUCTURE_TYPE_FRAMEBUFFER_CREATE_INFO                   = 37
	VK_STRUCTURE_TYPE_RENDER_PASS_CREATE_INFO                   = 38
	VK_STRUCTURE_TYPE_COMMAND_POOL_CREATE_INFO                  = 39
	VK_STRUCTURE_TYPE_COMMAND_BUFFER_ALLOCATE_INFO              = 40
	VK_STRUCTURE_TYPE_COMMAND_BUFFER_BEGIN_INFO                 = 42
	VK_STRUCTURE_TYPE_RENDER_PASS_BEGIN_INFO                    = 43
	VK_STRUCTURE_TYPE_MEMORY_BARRIER                            = 46
	VK_STRUCTURE_TYPE_BUFFER_MEMORY_BARRIER                     = 44
	VK_STRUCTURE_TYPE_IMAGE_MEMORY_BARRIER                      = 45
	VK_STRUCTURE_TYPE_PRESENT_INFO_KHR                          = 1000001001

	VK_DESCRIPTOR_TYPE_COMBINED_IMAGE_SAMPLER = 1
	VK_DESCRIPTOR_TYPE_UNIFORM_BUFFER         = 6
	VK_DESCRIPTOR_TYPE_STORAGE_BUFFER         = 7
	VK_DESCRIPTOR_TYPE_UNIFORM_BUFFER_DYNAMIC = 8
	VK_SHADER_STAGE_VERTEX_BIT                = 0x00000001
	VK_SHADER_STAGE_FRAGMENT_BIT              = 0x00000010
	VK_SHADER_STAGE_COMPUTE_BIT               = 0x00000020

	VK_BUFFER_USAGE_TRANSFER_SRC_BIT            = 0x00000001
	VK_BUFFER_USAGE_TRANSFER_DST_BIT            = 0x00000002
	VK_BUFFER_USAGE_UNIFORM_BUFFER_BIT          = 0x00000010
	VK_BUFFER_USAGE_STORAGE_BUFFER_BIT          = 0x00000020
	VK_BUFFER_USAGE_INDEX_BUFFER_BIT            = 0x00000040
	VK_BUFFER_USAGE_VERTEX_BUFFER_BIT           = 0x00000080
	VK_BUFFER_USAGE_INDIRECT_BUFFER_BIT         = 0x00000100
	VK_IMAGE_USAGE_TRANSFER_SRC_BIT             = 0x00000001
	VK_IMAGE_USAGE_TRANSFER_DST_BIT             = 0x00000002
	VK_IMAGE_USAGE_SAMPLED_BIT                  = 0x00000004
	VK_IMAGE_USAGE_COLOR_ATTACHMENT_BIT         = 0x00000010
	VK_IMAGE_USAGE_DEPTH_STENCIL_ATTACHMENT_BIT = 0x00000020
	VK_IMAGE_CREATE_CUBE_COMPATIBLE_BIT         = 0x00000010
	VK_MEMORY_PROPERTY_DEVICE_LOCAL_BIT         = 0x00000001
	VK_MEMORY_PROPERTY_HOST_VISIBLE_BIT         = 0x00000002
	VK_MEMORY_PROPERTY_HOST_COHERENT_BIT        = 0x00000004

	VK_PIPELINE_BIND_POINT_GRAPHICS                   = 0
	VK_PIPELINE_BIND_POINT_COMPUTE                    = 1
	VK_COMMAND_BUFFER_LEVEL_PRIMARY                   = 0
	VK_COMMAND_BUFFER_USAGE_ONE_TIME_SUBMIT_BIT       = 0x00000001
	VK_COMMAND_POOL_CREATE_TRANSIENT_BIT              = 0x00000001
	VK_COMMAND_POOL_CREATE_RESET_COMMAND_BUFFER_BIT   = 0x00000002
	VK_FENCE_CREATE_SIGNALED_BIT                      = 0x00000001
	VK_DESCRIPTOR_POOL_CREATE_FREE_DESCRIPTOR_SET_BIT = 0x00000001
	VK_PIPELINE_STAGE_TOP_OF_PIPE_BIT                 = 0x00000001
	VK_PIPELINE_STAGE_DRAW_INDIRECT_BIT               = 0x00000002
	VK_PIPELINE_STAGE_VERTEX_INPUT_BIT                = 0x00000004
	VK_PIPELINE_STAGE_VERTEX_SHADER_BIT               = 0x00000008
	VK_PIPELINE_STAGE_TRANSFER_BIT                    = 0x00001000
	VK_PIPELINE_STAGE_FRAGMENT_SHADER_BIT             = 0x00000080
	VK_PIPELINE_STAGE_EARLY_FRAGMENT_TESTS_BIT        = 0x00000100
	VK_PIPELINE_STAGE_LATE_FRAGMENT_TESTS_BIT         = 0x00000200
	VK_PIPELINE_STAGE_COLOR_ATTACHMENT_OUTPUT_BIT     = 0x00000400
	VK_PIPELINE_STAGE_COMPUTE_SHADER_BIT              = 0x00000800
	VK_PIPELINE_STAGE_BOTTOM_OF_PIPE_BIT              = 0x00002000
	VK_ACCESS_INDIRECT_COMMAND_READ_BIT               = 0x00000001
	VK_ACCESS_INDEX_READ_BIT                          = 0x00000002
	VK_ACCESS_VERTEX_ATTRIBUTE_READ_BIT               = 0x00000004
	VK_ACCESS_TRANSFER_READ_BIT                       = 0x00000800
	VK_ACCESS_TRANSFER_WRITE_BIT                      = 0x00001000
	VK_ACCESS_SHADER_READ_BIT                         = 0x00000020
	VK_ACCESS_SHADER_WRITE_BIT                        = 0x00000040
	VK_ACCESS_COLOR_ATTACHMENT_READ_BIT               = 0x00000080
	VK_ACCESS_DEPTH_STENCIL_ATTACHMENT_READ_BIT       = 0x00000200
	VK_ACCESS_COLOR_ATTACHMENT_WRITE_BIT              = 0x00000100
	VK_ACCESS_DEPTH_STENCIL_ATTACHMENT_WRITE_BIT      = 0x00000400
	VK_PIPELINE_STAGE_HOST_BIT = 0x00004000
	VK_ACCESS_HOST_READ_BIT = 0x00002000
	VK_ACCESS_MEMORY_READ_BIT                         = 0x00008000

	VK_IMAGE_LAYOUT_UNDEFINED                        = 0
	VK_IMAGE_LAYOUT_GENERAL                          = 1
	VK_IMAGE_LAYOUT_COLOR_ATTACHMENT_OPTIMAL         = 2
	VK_IMAGE_LAYOUT_DEPTH_STENCIL_ATTACHMENT_OPTIMAL = 3
	VK_IMAGE_LAYOUT_DEPTH_STENCIL_READ_ONLY_OPTIMAL  = 4
	VK_IMAGE_LAYOUT_SHADER_READ_ONLY_OPTIMAL         = 5
	VK_IMAGE_LAYOUT_TRANSFER_SRC_OPTIMAL             = 6
	VK_IMAGE_LAYOUT_TRANSFER_DST_OPTIMAL             = 7
	VK_IMAGE_LAYOUT_PRESENT_SRC_KHR                  = 1000001002

	VK_QUEUE_GRAPHICS_BIT       = 0x00000001
	VK_QUEUE_COMPUTE_BIT        = 0x00000002
	VK_QUEUE_TRANSFER_BIT       = 0x00000004
	VK_DEPENDENCY_BY_REGION_BIT = 0x00000001
	VK_SUBPASS_EXTERNAL         = ^uint32(0)
	VK_QUEUE_FAMILY_IGNORED     = ^uint32(0)

	VK_IMAGE_ASPECT_COLOR_BIT   = 0x00000001
	VK_IMAGE_ASPECT_DEPTH_BIT   = 0x00000002
	VK_IMAGE_ASPECT_STENCIL_BIT = 0x00000004

	VK_IMAGE_TYPE_2D           = 1
	VK_INDEX_TYPE_UINT32       = 1
	VK_IMAGE_VIEW_TYPE_2D      = 1
	VK_IMAGE_VIEW_TYPE_CUBE    = 3
	VK_IMAGE_TILING_OPTIMAL    = 0
	VK_SHARING_MODE_EXCLUSIVE  = 0
	VK_SHARING_MODE_CONCURRENT = 1

	VK_FORMAT_R8G8B8A8_UNORM       = 37
	VK_FORMAT_R8G8B8A8_SRGB        = 43
	VK_FORMAT_B8G8R8A8_UNORM       = 44
	VK_FORMAT_B8G8R8A8_SRGB        = 50
	VK_FORMAT_R16G16B16A16_SFLOAT  = 97
	VK_FORMAT_R32G32_SFLOAT        = 103
	VK_FORMAT_R32G32B32_SFLOAT     = 106
	VK_FORMAT_R32G32B32A32_SFLOAT  = 109
	VK_FORMAT_D32_SFLOAT           = 126
	VK_FORMAT_BC1_RGBA_UNORM_BLOCK = 133
	VK_FORMAT_BC1_RGBA_SRGB_BLOCK  = 134
	VK_FORMAT_BC3_UNORM_BLOCK      = 137
	VK_FORMAT_BC3_SRGB_BLOCK       = 138
	VK_FORMAT_BC5_UNORM_BLOCK      = 141
	VK_FORMAT_BC7_UNORM_BLOCK      = 145
	VK_FORMAT_BC7_SRGB_BLOCK       = 146
	VK_FORMAT_ASTC_4X4_UNORM_BLOCK = 157
	VK_FORMAT_ASTC_4X4_SRGB_BLOCK  = 158

	VK_FORMAT_FEATURE_SAMPLED_IMAGE_BIT               = 0x00000001
	VK_FORMAT_FEATURE_COLOR_ATTACHMENT_BIT            = 0x00000080
	VK_FORMAT_FEATURE_COLOR_ATTACHMENT_BLEND_BIT      = 0x00000100
	VK_FORMAT_FEATURE_DEPTH_STENCIL_ATTACHMENT_BIT    = 0x00000200
	VK_FORMAT_FEATURE_BLIT_SRC_BIT                    = 0x00000400
	VK_FORMAT_FEATURE_BLIT_DST_BIT                    = 0x00000800
	VK_FORMAT_FEATURE_SAMPLED_IMAGE_FILTER_LINEAR_BIT = 0x00001000
	VK_FORMAT_FEATURE_TRANSFER_SRC_BIT                = 0x00004000
	VK_FORMAT_FEATURE_TRANSFER_DST_BIT                = 0x00008000

	VK_COLOR_SPACE_SRGB_NONLINEAR_KHR          = 0
	VK_PRESENT_MODE_IMMEDIATE_KHR              = 0
	VK_PRESENT_MODE_MAILBOX_KHR                = 1
	VK_PRESENT_MODE_FIFO_KHR                   = 2
	VK_PRESENT_MODE_FIFO_RELAXED_KHR           = 3
	VK_COMPOSITE_ALPHA_OPAQUE_BIT_KHR          = 0x00000001
	VK_COMPOSITE_ALPHA_PRE_MULTIPLIED_BIT_KHR  = 0x00000002
	VK_COMPOSITE_ALPHA_POST_MULTIPLIED_BIT_KHR = 0x00000004
	VK_COMPOSITE_ALPHA_INHERIT_BIT_KHR         = 0x00000008

	VK_SAMPLE_COUNT_1_BIT            = 0x00000001
	VK_ATTACHMENT_LOAD_OP_LOAD       = 0
	VK_ATTACHMENT_LOAD_OP_CLEAR      = 1
	VK_ATTACHMENT_LOAD_OP_DONT_CARE  = 2
	VK_ATTACHMENT_STORE_OP_STORE     = 0
	VK_ATTACHMENT_STORE_OP_DONT_CARE = 1

	VK_BLEND_FACTOR_ZERO                = 0
	VK_BLEND_FACTOR_ONE                 = 1
	VK_BLEND_FACTOR_SRC_ALPHA           = 6
	VK_BLEND_FACTOR_ONE_MINUS_SRC_ALPHA = 7
	VK_BLEND_OP_ADD                     = 0
)

type VkDeviceQueueCreateInfo struct {
	SType            uint32
	PNext            uintptr
	Flags            uint32
	QueueFamilyIndex uint32
	QueueCount       uint32
	PQueuePriorities *float32
}

type VkDeviceCreateInfo struct {
	SType                   uint32
	PNext                   uintptr
	Flags                   uint32
	QueueCreateInfoCount    uint32
	PQueueCreateInfos       *VkDeviceQueueCreateInfo
	EnabledLayerCount       uint32
	PpEnabledLayerNames     **byte
	EnabledExtensionCount   uint32
	PpEnabledExtensionNames **byte
	PEnabledFeatures        uintptr
}

type VkPushConstantRange struct {
	StageFlags uint32
	Offset     uint32
	Size       uint32
}

type VkPipelineLayoutCreateInfo struct {
	SType                  uint32
	PNext                  uintptr
	Flags                  uint32
	SetLayoutCount         uint32
	PSetLayouts            *uint64
	PushConstantRangeCount uint32
	PPushConstantRanges    *VkPushConstantRange
}

type VkApplicationInfo struct {
	SType              uint32
	PNext              uintptr
	PApplicationName   *byte
	ApplicationVersion uint32
	PEngineName        *byte
	EngineVersion      uint32
	ApiVersion         uint32
}

type VkInstanceCreateInfo struct {
	SType                   uint32
	PNext                   uintptr
	Flags                   uint32
	PApplicationInfo        *VkApplicationInfo
	EnabledLayerCount       uint32
	PpEnabledLayerNames     **byte
	EnabledExtensionCount   uint32
	PpEnabledExtensionNames **byte
}

type VkSwapchainCreateInfoKHR struct {
	SType                 uint32
	PNext                 uintptr
	Flags                 uint32
	Surface               uint64
	MinImageCount         uint32
	ImageFormat           uint32
	ImageColorSpace       uint32
	ImageExtent           VkExtent2D
	ImageArrayLayers      uint32
	ImageUsage            uint32
	ImageSharingMode      uint32
	QueueFamilyIndexCount uint32
	PQueueFamilyIndices   *uint32
	PreTransform          uint32
	CompositeAlpha        uint32
	PresentMode           uint32
	Clipped               uint32
	OldSwapchain          uint64
}

type VkQueueFamilyProperties struct {
	QueueFlags                  uint32
	QueueCount                  uint32
	TimestampValidBits          uint32
	MinImageTransferGranularity VkExtent3D
}

type VkSurfaceCapabilitiesKHR struct {
	MinImageCount           uint32
	MaxImageCount           uint32
	CurrentExtent           VkExtent2D
	MinImageExtent          VkExtent2D
	MaxImageExtent          VkExtent2D
	MaxImageArrayLayers     uint32
	SupportedTransforms     uint32
	CurrentTransform        uint32
	SupportedCompositeAlpha uint32
	SupportedUsageFlags     uint32
}

type VkSurfaceFormatKHR struct {
	Format     uint32
	ColorSpace uint32
}

type VkSemaphoreCreateInfo struct {
	SType uint32
	PNext uintptr
	Flags uint32
}

type VkFenceCreateInfo struct {
	SType uint32
	PNext uintptr
	Flags uint32
}

type VkExtent2D struct {
	Width  uint32
	Height uint32
}

type VkExtent3D struct {
	Width  uint32
	Height uint32
	Depth  uint32
}

type VkOffset2D struct {
	X int32
	Y int32
}

type VkOffset3D struct {
	X int32
	Y int32
	Z int32
}

type VkRect2D struct {
	Offset VkOffset2D
	Extent VkExtent2D
}

type VkClearColorValue struct {
	Float32 [4]float32
}

type VkClearDepthStencilValue struct {
	Depth   float32
	Stencil uint32
}

type VkClearValue struct {
	Color VkClearColorValue
}

type VkAttachmentDescription struct {
	Flags          uint32
	Format         uint32
	Samples        uint32
	LoadOp         uint32
	StoreOp        uint32
	StencilLoadOp  uint32
	StencilStoreOp uint32
	InitialLayout  uint32
	FinalLayout    uint32
}

type VkCommandPoolCreateInfo struct {
	SType            uint32
	PNext            uintptr
	Flags            uint32
	QueueFamilyIndex uint32
}

type VkAttachmentReference struct {
	Attachment uint32
	Layout     uint32
}

type VkSubpassDescription struct {
	Flags                   uint32
	PipelineBindPoint       uint32
	InputAttachmentCount    uint32
	PInputAttachments       *VkAttachmentReference
	ColorAttachmentCount    uint32
	PColorAttachments       *VkAttachmentReference
	PResolveAttachments     *VkAttachmentReference
	PDepthStencilAttachment *VkAttachmentReference
	PreserveAttachmentCount uint32
	PPreserveAttachments    *uint32
}

type VkSubpassDependency struct {
	SrcSubpass      uint32
	DstSubpass      uint32
	SrcStageMask    uint32
	DstStageMask    uint32
	SrcAccessMask   uint32
	DstAccessMask   uint32
	DependencyFlags uint32
}

type VkRenderPassCreateInfo struct {
	SType           uint32
	PNext           uintptr
	Flags           uint32
	AttachmentCount uint32
	PAttachments    *VkAttachmentDescription
	SubpassCount    uint32
	PSubpasses      *VkSubpassDescription
	DependencyCount uint32
	PDependencies   *VkSubpassDependency
}

type VkRenderPassBeginInfo struct {
	SType           uint32
	PNext           uintptr
	RenderPass      uint64
	Framebuffer     uint64
	RenderArea      VkRect2D
	ClearValueCount uint32
	PClearValues    *VkClearValue
}

type VkCommandBufferAllocateInfo struct {
	SType              uint32
	PNext              uintptr
	CommandPool        uint64
	Level              uint32
	CommandBufferCount uint32
}

type VkCommandBufferBeginInfo struct {
	SType            uint32
	PNext            uintptr
	Flags            uint32
	PInheritanceInfo uintptr
}

type VkSubmitInfo struct {
	SType                uint32
	PNext                uintptr
	WaitSemaphoreCount   uint32
	PWaitSemaphores      *uint64
	PWaitDstStageMask    *uint32
	CommandBufferCount   uint32
	PCommandBuffers      *uintptr
	SignalSemaphoreCount uint32
	PSignalSemaphores    *uint64
}

type VkPresentInfoKHR struct {
	SType              uint32
	PNext              uintptr
	WaitSemaphoreCount uint32
	PWaitSemaphores    *uint64
	SwapchainCount     uint32
	PSwapchains        *uint64
	PImageIndices      *uint32
	PResults           *uint32
}

type VkBufferCreateInfo struct {
	SType                 uint32
	PNext                 uintptr
	Flags                 uint32
	Size                  uint64
	Usage                 uint32
	SharingMode           uint32
	QueueFamilyIndexCount uint32
	PQueueFamilyIndices   *uint32
}

type VkMemoryRequirements struct {
	Size           uint64
	Alignment      uint64
	MemoryTypeBits uint32
}

type VkMemoryAllocateInfo struct {
	SType           uint32
	PNext           uintptr
	AllocationSize  uint64
	MemoryTypeIndex uint32
}

type VkMemoryType struct {
	PropertyFlags uint32
	HeapIndex     uint32
}

type VkMemoryHeap struct {
	Size  uint64
	Flags uint32
}

type VkPhysicalDeviceMemoryProperties struct {
	MemoryTypeCount uint32
	MemoryTypes     [32]VkMemoryType
	MemoryHeapCount uint32
	MemoryHeaps     [16]VkMemoryHeap
}

type VkPhysicalDeviceFeatures struct {
	RobustBufferAccess                      uint32
	FullDrawIndexUint32                     uint32
	ImageCubeArray                          uint32
	IndependentBlend                        uint32
	GeometryShader                          uint32
	TessellationShader                      uint32
	SampleRateShading                       uint32
	DualSrcBlend                            uint32
	LogicOp                                 uint32
	MultiDrawIndirect                       uint32
	DrawIndirectFirstInstance               uint32
	DepthClamp                              uint32
	DepthBiasClamp                          uint32
	FillModeNonSolid                        uint32
	DepthBounds                             uint32
	WideLines                               uint32
	LargePoints                             uint32
	AlphaToOne                              uint32
	MultiViewport                           uint32
	SamplerAnisotropy                       uint32
	TextureCompressionETC2                  uint32
	TextureCompressionASTCLDR               uint32
	TextureCompressionBC                    uint32
	OcclusionQueryPrecise                   uint32
	PipelineStatisticsQuery                 uint32
	VertexPipelineStoresAndAtomics          uint32
	FragmentStoresAndAtomics                uint32
	ShaderTessellationAndGeometryPointSize  uint32
	ShaderImageGatherExtended               uint32
	ShaderStorageImageExtendedFormats       uint32
	ShaderStorageImageMultisample           uint32
	ShaderStorageImageReadWithoutFormat     uint32
	ShaderStorageImageWriteWithoutFormat    uint32
	ShaderUniformBufferArrayDynamicIndexing uint32
	ShaderSampledImageArrayDynamicIndexing  uint32
	ShaderStorageBufferArrayDynamicIndexing uint32
	ShaderStorageImageArrayDynamicIndexing  uint32
	ShaderClipDistance                      uint32
	ShaderCullDistance                      uint32
	ShaderFloat64                           uint32
	ShaderInt64                             uint32
	ShaderInt16                             uint32
	ShaderResourceResidency                 uint32
	ShaderResourceMinLod                    uint32
	SparseBinding                           uint32
	SparseResidencyBuffer                   uint32
	SparseResidencyImage2D                  uint32
	SparseResidencyImage3D                  uint32
	SparseResidency2Samples                 uint32
	SparseResidency4Samples                 uint32
	SparseResidency8Samples                 uint32
	SparseResidency16Samples                uint32
	SparseResidencyAliased                  uint32
	VariableMultisampleRate                 uint32
	InheritedQueries                        uint32
}

type VkBufferCopy struct {
	SrcOffset uint64
	DstOffset uint64
	Size      uint64
}

type VkMemoryBarrier struct {
	SType         uint32
	PNext         uintptr
	SrcAccessMask uint32
	DstAccessMask uint32
}

type VkFormatProperties struct {
	LinearTilingFeatures  uint32
	OptimalTilingFeatures uint32
	BufferFeatures        uint32
}

type VkImageBlit struct {
	SrcSubresource VkImageSubresourceLayers
	SrcOffsets     [2]VkOffset3D
	DstSubresource VkImageSubresourceLayers
	DstOffsets     [2]VkOffset3D
}

type VkBufferMemoryBarrier struct {
	SType               uint32
	PNext               uintptr
	SrcAccessMask       uint32
	DstAccessMask       uint32
	SrcQueueFamilyIndex uint32
	DstQueueFamilyIndex uint32
	Buffer              uint64
	Offset              uint64
	Size                uint64
}

type VkImageSubresourceLayers struct {
	AspectMask     uint32
	MipLevel       uint32
	BaseArrayLayer uint32
	LayerCount     uint32
}

type VkImageSubresourceRange struct {
	AspectMask     uint32
	BaseMipLevel   uint32
	LevelCount     uint32
	BaseArrayLayer uint32
	LayerCount     uint32
}

type VkBufferImageCopy struct {
	BufferOffset      uint64
	BufferRowLength   uint32
	BufferImageHeight uint32
	ImageSubresource  VkImageSubresourceLayers
	ImageOffset       VkOffset3D
	ImageExtent       VkExtent3D
}

type VkImageMemoryBarrier struct {
	SType               uint32
	PNext               uintptr
	SrcAccessMask       uint32
	DstAccessMask       uint32
	OldLayout           uint32
	NewLayout           uint32
	SrcQueueFamilyIndex uint32
	DstQueueFamilyIndex uint32
	Image               uint64
	SubresourceRange    VkImageSubresourceRange
}

type VkImageCreateInfo struct {
	SType                 uint32
	PNext                 uintptr
	Flags                 uint32
	ImageType             uint32
	Format                uint32
	Extent                VkExtent3D
	MipLevels             uint32
	ArrayLayers           uint32
	Samples               uint32
	Tiling                uint32
	Usage                 uint32
	SharingMode           uint32
	QueueFamilyIndexCount uint32
	PQueueFamilyIndices   *uint32
	InitialLayout         uint32
}

type VkImageViewCreateInfo struct {
	SType            uint32
	PNext            uintptr
	Flags            uint32
	Image            uint64
	ViewType         uint32
	Format           uint32
	Components       [4]uint32
	SubresourceRange [5]uint32
}

type VkFramebufferCreateInfo struct {
	SType           uint32
	PNext           uintptr
	Flags           uint32
	RenderPass      uint64
	AttachmentCount uint32
	PAttachments    *uint64
	Width           uint32
	Height          uint32
	Layers          uint32
}

type VkShaderModuleCreateInfo struct {
	SType    uint32
	PNext    uintptr
	Flags    uint32
	CodeSize uint64
	PCode    *uint32
}

type VkExtensionProperties struct {
	ExtensionName [256]byte
	SpecVersion   uint32
}

const (
	VK_DYNAMIC_STATE_VIEWPORT = 0
	VK_DYNAMIC_STATE_SCISSOR  = 1

	VK_FILTER_NEAREST = 0
	VK_FILTER_LINEAR  = 1

	VK_SAMPLER_MIPMAP_MODE_NEAREST = 0
	VK_SAMPLER_MIPMAP_MODE_LINEAR  = 1

	VK_SAMPLER_ADDRESS_MODE_REPEAT        = 0
	VK_SAMPLER_ADDRESS_MODE_CLAMP_TO_EDGE = 2

	VK_STRUCTURE_TYPE_SWAPCHAIN_CREATE_INFO_KHR = 1000001000
	VK_PRIMITIVE_TOPOLOGY_TRIANGLE_LIST         = 3
	VK_POLYGON_MODE_FILL                        = 0
	VK_CULL_MODE_NONE                           = 0
	VK_CULL_MODE_BACK_BIT                       = 0x00000002
	VK_FRONT_FACE_COUNTER_CLOCKWISE             = 0
	VK_COMPARE_OP_LESS                          = 1
	VK_COMPARE_OP_ALWAYS                        = 7
)

type VkViewport struct {
	X, Y, Width, Height, MinDepth, MaxDepth float32
}

type VkPipelineShaderStageCreateInfo struct {
	SType               uint32
	PNext               uintptr
	Flags               uint32
	Stage               uint32
	Module              uint64
	PName               *byte
	PSpecializationInfo uintptr
}

type VkVertexInputBindingDescription struct {
	Binding   uint32
	Stride    uint32
	InputRate uint32
}

type VkVertexInputAttributeDescription struct {
	Location uint32
	Binding  uint32
	Format   uint32
	Offset   uint32
}

type VkPipelineVertexInputStateCreateInfo struct {
	SType                           uint32
	PNext                           uintptr
	Flags                           uint32
	VertexBindingDescriptionCount   uint32
	PVertexBindingDescriptions      *VkVertexInputBindingDescription
	VertexAttributeDescriptionCount uint32
	PVertexAttributeDescriptions    *VkVertexInputAttributeDescription
}

type VkPipelineInputAssemblyStateCreateInfo struct {
	SType                  uint32
	PNext                  uintptr
	Flags                  uint32
	Topology               uint32
	PrimitiveRestartEnable uint32
}

type VkPipelineViewportStateCreateInfo struct {
	SType         uint32
	PNext         uintptr
	Flags         uint32
	ViewportCount uint32
	PViewports    *VkViewport
	ScissorCount  uint32
	PScissors     *VkRect2D
}

type VkPipelineRasterizationStateCreateInfo struct {
	SType                   uint32
	PNext                   uintptr
	Flags                   uint32
	DepthClampEnable        uint32
	RasterizerDiscardEnable uint32
	PolygonMode             uint32
	CullMode                uint32
	FrontFace               uint32
	DepthBiasEnable         uint32
	DepthBiasConstantFactor float32
	DepthBiasClamp          float32
	DepthBiasSlopeFactor    float32
	LineWidth               float32
}

type VkPipelineMultisampleStateCreateInfo struct {
	SType                 uint32
	PNext                 uintptr
	Flags                 uint32
	RasterizationSamples  uint32
	SampleShadingEnable   uint32
	MinSampleShading      float32
	PSampleMask           *uint32
	AlphaToCoverageEnable uint32
	AlphaToOneEnable      uint32
}

type VkPipelineDepthStencilStateCreateInfo struct {
	SType                 uint32
	PNext                 uintptr
	Flags                 uint32
	DepthTestEnable       uint32
	DepthWriteEnable      uint32
	DepthCompareOp        uint32
	DepthBoundsTestEnable uint32
	StencilTestEnable     uint32
	Front                 [7]uint32
	Back                  [7]uint32
	MinDepthBounds        float32
	MaxDepthBounds        float32
}

type VkPipelineColorBlendAttachmentState struct {
	BlendEnable         uint32
	SrcColorBlendFactor uint32
	DstColorBlendFactor uint32
	ColorBlendOp        uint32
	SrcAlphaBlendFactor uint32
	DstAlphaBlendFactor uint32
	AlphaBlendOp        uint32
	ColorWriteMask      uint32
}

type VkPipelineColorBlendStateCreateInfo struct {
	SType           uint32
	PNext           uintptr
	Flags           uint32
	LogicOpEnable   uint32
	LogicOp         uint32
	AttachmentCount uint32
	PAttachments    *VkPipelineColorBlendAttachmentState
	BlendConstants  [4]float32
}

type VkPipelineDynamicStateCreateInfo struct {
	SType             uint32
	PNext             uintptr
	Flags             uint32
	DynamicStateCount uint32
	PDynamicStates    *uint32
}

type VkGraphicsPipelineCreateInfo struct {
	SType               uint32
	PNext               uintptr
	Flags               uint32
	StageCount          uint32
	PStages             *VkPipelineShaderStageCreateInfo
	PVertexInputState   *VkPipelineVertexInputStateCreateInfo
	PInputAssemblyState *VkPipelineInputAssemblyStateCreateInfo
	PTessellationState  uintptr
	PViewportState      *VkPipelineViewportStateCreateInfo
	PRasterizationState *VkPipelineRasterizationStateCreateInfo
	PMultisampleState   *VkPipelineMultisampleStateCreateInfo
	PDepthStencilState  *VkPipelineDepthStencilStateCreateInfo
	PColorBlendState    *VkPipelineColorBlendStateCreateInfo
	PDynamicState       *VkPipelineDynamicStateCreateInfo
	Layout              uint64
	RenderPass          uint64
	Subpass             uint32
	BasePipelineHandle  uint64
	BasePipelineIndex   int32
}

type VkComputePipelineCreateInfo struct {
	SType              uint32
	PNext              uintptr
	Flags              uint32
	Stage              VkPipelineShaderStageCreateInfo
	Layout             uint64
	BasePipelineHandle uint64
	BasePipelineIndex  int32
}

type VkDrawIndexedIndirectCommand struct {
	IndexCount    uint32
	InstanceCount uint32
	FirstIndex    uint32
	VertexOffset  int32
	FirstInstance uint32
}

type VkDescriptorSetLayoutBinding struct {
	Binding            uint32
	DescriptorType     uint32
	DescriptorCount    uint32
	StageFlags         uint32
	PImmutableSamplers uintptr
}

type VkDescriptorSetLayoutCreateInfo struct {
	SType        uint32
	PNext        uintptr
	Flags        uint32
	BindingCount uint32
	PBindings    *VkDescriptorSetLayoutBinding
}

type VkDescriptorPoolSize struct {
	Type            uint32
	DescriptorCount uint32
}

type VkDescriptorPoolCreateInfo struct {
	SType         uint32
	PNext         uintptr
	Flags         uint32
	MaxSets       uint32
	PoolSizeCount uint32
	PPoolSizes    *VkDescriptorPoolSize
}

type VkDescriptorSetAllocateInfo struct {
	SType              uint32
	PNext              uintptr
	DescriptorPool     uint64
	DescriptorSetCount uint32
	PSetLayouts        *uint64
}

type VkDescriptorBufferInfo struct {
	Buffer uint64
	Offset uint64
	Range  uint64
}

type VkDescriptorImageInfo struct {
	Sampler     uint64
	ImageView   uint64
	ImageLayout uint32
}

type VkWriteDescriptorSet struct {
	SType            uint32
	PNext            uintptr
	DstSet           uint64
	DstBinding       uint32
	DstArrayElement  uint32
	DescriptorCount  uint32
	DescriptorType   uint32
	PImageInfo       *VkDescriptorImageInfo
	PBufferInfo      *VkDescriptorBufferInfo
	PTexelBufferView uintptr
}

type VkSamplerCreateInfo struct {
	SType                   uint32
	PNext                   uintptr
	Flags                   uint32
	MagFilter               uint32
	MinFilter               uint32
	MipmapMode              uint32
	AddressModeU            uint32
	AddressModeV            uint32
	AddressModeW            uint32
	MipLodBias              float32
	AnisotropyEnable        uint32
	MaxAnisotropy           float32
	CompareEnable           uint32
	CompareOp               uint32
	MinLod                  float32
	MaxLod                  float32
	BorderColor             uint32
	UnnormalizedCoordinates uint32
}

var (
	vkLibHandle uintptr

	VkCreateInstance                          func(createInfo *VkInstanceCreateInfo, allocator uintptr, instance *uintptr) int32
	VkDestroyInstance                         func(instance uintptr, allocator uintptr)
	VkEnumeratePhysicalDevices                func(instance uintptr, pPhysicalDeviceCount *uint32, pPhysicalDevices *uintptr) int32
	VkGetPhysicalDeviceQueueFamilyProperties  func(physicalDevice uintptr, pQueueFamilyPropertyCount *uint32, pQueueFamilyProperties *VkQueueFamilyProperties)
	VkGetPhysicalDeviceSurfaceSupportKHR      func(physicalDevice uintptr, queueFamilyIndex uint32, surface uint64, pSupported *uint32) int32
	VkGetPhysicalDeviceSurfaceCapabilitiesKHR func(physicalDevice uintptr, surface uint64, pSurfaceCapabilities *VkSurfaceCapabilitiesKHR) int32
	VkGetPhysicalDeviceSurfaceFormatsKHR      func(physicalDevice uintptr, surface uint64, pSurfaceFormatCount *uint32, pSurfaceFormats *VkSurfaceFormatKHR) int32
	VkGetPhysicalDeviceSurfacePresentModesKHR func(physicalDevice uintptr, surface uint64, pPresentModeCount *uint32, pPresentModes *uint32) int32
	VkCreateDevice                            func(physicalDevice uintptr, pCreateInfo *VkDeviceCreateInfo, pAllocator uintptr, pDevice *uintptr) int32
	VkDestroyDevice                           func(device uintptr, pAllocator uintptr)
	VkDeviceWaitIdle                          func(device uintptr) int32
	VkGetPhysicalDeviceMemoryProperties       func(physicalDevice uintptr, pMemoryProperties *VkPhysicalDeviceMemoryProperties)
	VkGetPhysicalDeviceProperties func(physicalDevice uintptr, properties *VkPhysicalDeviceProperties)
	VkCreateQueryPool func(device uintptr, info *VkQueryPoolCreateInfo, allocator uintptr, pool *uint64) int32
	VkDestroyQueryPool func(device uintptr, pool uint64, allocator uintptr)
	VkCmdResetQueryPool func(commandBuffer uintptr, pool uint64, firstQuery, queryCount uint32)
	VkCmdWriteTimestamp func(commandBuffer uintptr, stage uint32, pool uint64, query uint32)
	VkGetQueryPoolResults func(device uintptr, pool uint64, firstQuery, queryCount uint32, dataSize uintptr, data unsafe.Pointer, stride uint64, flags uint32) int32
	VkGetPhysicalDeviceFeatures               func(physicalDevice uintptr, pFeatures *VkPhysicalDeviceFeatures)
	VkGetPhysicalDeviceFormatProperties       func(physicalDevice uintptr, format uint32, pFormatProperties *VkFormatProperties)

	VkCreateBuffer                func(device uintptr, pCreateInfo *VkBufferCreateInfo, pAllocator uintptr, pBuffer *uint64) int32
	VkDestroyBuffer               func(device uintptr, buffer uint64, pAllocator uintptr)
	VkGetBufferMemoryRequirements func(device uintptr, buffer uint64, pMemoryRequirements *VkMemoryRequirements)
	VkGetImageMemoryRequirements  func(device uintptr, image uint64, pMemoryRequirements *VkMemoryRequirements)
	VkAllocateMemory              func(device uintptr, pAllocateInfo *VkMemoryAllocateInfo, pAllocator uintptr, pMemory *uint64) int32
	VkFreeMemory                  func(device uintptr, memory uint64, pAllocator uintptr)
	VkBindBufferMemory            func(device uintptr, buffer uint64, memory uint64, memoryOffset uint64) int32
	VkBindImageMemory             func(device uintptr, image uint64, memory uint64, memoryOffset uint64) int32
	VkMapMemory                   func(device uintptr, memory uint64, offset uint64, size uint64, flags uint32, ppData *unsafe.Pointer) int32
	VkUnmapMemory                 func(device uintptr, memory uint64)
	VkCreateSwapchainKHR          func(device uintptr, pCreateInfo *VkSwapchainCreateInfoKHR, pAllocator uintptr, pSwapchain *uint64) int32
	VkDestroySwapchainKHR         func(device uintptr, swapchain uint64, pAllocator uintptr)
	VkCreateSemaphore             func(device uintptr, pCreateInfo *VkSemaphoreCreateInfo, pAllocator uintptr, pSemaphore *uint64) int32
	VkDestroySemaphore            func(device uintptr, semaphore uint64, pAllocator uintptr)
	VkCreateFence                 func(device uintptr, pCreateInfo *VkFenceCreateInfo, pAllocator uintptr, pFence *uint64) int32
	VkDestroyFence                func(device uintptr, fence uint64, pAllocator uintptr)

	VkCreateImage        func(device uintptr, pCreateInfo *VkImageCreateInfo, pAllocator uintptr, pImage *uint64) int32
	VkDestroyImage       func(device uintptr, image uint64, pAllocator uintptr)
	VkCreateImageView    func(device uintptr, pCreateInfo *VkImageViewCreateInfo, pAllocator uintptr, pView *uint64) int32
	VkDestroyImageView   func(device uintptr, imageView uint64, pAllocator uintptr)
	VkCreateFramebuffer  func(device uintptr, pCreateInfo *VkFramebufferCreateInfo, pAllocator uintptr, pFramebuffer *uint64) int32
	VkDestroyFramebuffer func(device uintptr, framebuffer uint64, pAllocator uintptr)

	VkCreateShaderModule  func(device uintptr, pCreateInfo *VkShaderModuleCreateInfo, pAllocator uintptr, pShaderModule *uint64) int32
	VkDestroyShaderModule func(device uintptr, shaderModule uint64, pAllocator uintptr)

	VkAcquireNextImageKHR    func(device uintptr, swapchain uint64, timeout uint64, semaphore uint64, fence uint64, pImageIndex *uint32) int32
	VkAllocateCommandBuffers func(device uintptr, pAllocateInfo *VkCommandBufferAllocateInfo, pCommandBuffers *uintptr) int32
	VkFreeCommandBuffers     func(device uintptr, commandPool uint64, commandBufferCount uint32, pCommandBuffers *uintptr)
	VkBeginCommandBuffer     func(commandBuffer uintptr, pBeginInfo *VkCommandBufferBeginInfo) int32
	VkEndCommandBuffer       func(commandBuffer uintptr) int32
	VkQueueSubmit            func(queue uintptr, submitCount uint32, pSubmits *VkSubmitInfo, fence uint64) int32
	VkQueueWaitIdle          func(queue uintptr) int32
	VkQueuePresentKHR        func(queue uintptr, pPresentInfo *VkPresentInfoKHR) int32
	VkCreateRenderPass       func(device uintptr, pCreateInfo *VkRenderPassCreateInfo, pAllocator uintptr, pRenderPass *uint64) int32
	VkDestroyRenderPass      func(device uintptr, renderPass uint64, pAllocator uintptr)

	VkCmdCopyBuffer                        func(commandBuffer uintptr, srcBuffer uint64, dstBuffer uint64, regionCount uint32, pRegions *VkBufferCopy)
	VkCmdCopyBufferToImage                 func(commandBuffer uintptr, srcBuffer uint64, dstImage uint64, dstImageLayout uint32, regionCount uint32, pRegions *VkBufferImageCopy)
	VkCmdCopyImageToBuffer func(commandBuffer uintptr, srcImage uint64, srcImageLayout uint32, dstBuffer uint64, regionCount uint32, pRegions *VkBufferImageCopy)
	VkCmdPipelineBarrier                   func(commandBuffer uintptr, srcStageMask uint32, dstStageMask uint32, dependencyFlags uint32, memoryBarrierCount uint32, pMemoryBarriers *VkMemoryBarrier, bufferMemoryBarrierCount uint32, pBufferMemoryBarriers *VkBufferMemoryBarrier, imageMemoryBarrierCount uint32, pImageMemoryBarriers *VkImageMemoryBarrier)
	VkCmdBeginRenderPass                   func(commandBuffer uintptr, pRenderPassBegin *VkRenderPassBeginInfo, contents uint32)
	VkCmdEndRenderPass                     func(commandBuffer uintptr)
	VkCmdBindPipeline                      func(commandBuffer uintptr, pipelineBindPoint uint32, pipeline uint64)
	VkCmdBindDescriptorSets                func(commandBuffer uintptr, pipelineBindPoint uint32, layout uint64, firstSet uint32, descriptorSetCount uint32, pDescriptorSets *uint64, dynamicOffsetCount uint32, pDynamicOffsets *uint32)
	VkCmdBindVertexBuffers                 func(commandBuffer uintptr, firstBinding uint32, bindingCount uint32, pBuffers *uint64, pOffsets *uint64)
	VkCmdBindIndexBuffer                   func(commandBuffer uintptr, buffer uint64, offset uint64, indexType uint32)
	VkCmdDispatch                          func(commandBuffer uintptr, groupCountX uint32, groupCountY uint32, groupCountZ uint32)
	VkCmdDrawIndexedIndirect               func(commandBuffer uintptr, buffer uint64, offset uint64, drawCount uint32, stride uint32)
	VkCmdDrawIndexed                       func(commandBuffer uintptr, indexCount uint32, instanceCount uint32, firstIndex uint32, vertexOffset int32, firstInstance uint32)
	VkCmdDraw                              func(commandBuffer uintptr, vertexCount uint32, instanceCount uint32, firstVertex uint32, firstInstance uint32)
	VkCmdDrawInstanced                     func(commandBuffer uintptr, vertexCount uint32, instanceCount uint32, firstVertex uint32, firstInstance uint32)
	VkCmdPushConstants                     func(commandBuffer uintptr, layout uint64, stageFlags uint32, offset uint32, size uint32, pValues unsafe.Pointer)
	VkCmdBlitImage                         func(commandBuffer uintptr, srcImage uint64, srcImageLayout uint32, dstImage uint64, dstImageLayout uint32, regionCount uint32, pRegions *VkImageBlit, filter uint32)
	VkCreateCommandPool                    func(device uintptr, pCreateInfo *VkCommandPoolCreateInfo, pAllocator uintptr, pCommandPool *uint64) int32
	VkDestroyCommandPool                   func(device uintptr, commandPool uint64, pAllocator uintptr)
	VkGetDeviceQueue                       func(device uintptr, queueFamilyIndex uint32, queueIndex uint32, pQueue *uintptr)
	VkEnumerateInstanceExtensionProperties func(pLayerName *byte, pPropertyCount *uint32, pProperties *VkExtensionProperties) int32
	VkCreateDescriptorSetLayout            func(device uintptr, pCreateInfo *VkDescriptorSetLayoutCreateInfo, pAllocator uintptr, pSetLayout *uint64) int32
	VkDestroyDescriptorSetLayout           func(device uintptr, descriptorSetLayout uint64, pAllocator uintptr)
	VkCreatePipelineLayout                 func(device uintptr, pCreateInfo *VkPipelineLayoutCreateInfo, pAllocator uintptr, pPipelineLayout *uint64) int32
	VkDestroyPipelineLayout                func(device uintptr, pipelineLayout uint64, pAllocator uintptr)
	VkCreateDescriptorPool                 func(device uintptr, pCreateInfo *VkDescriptorPoolCreateInfo, pAllocator uintptr, pDescriptorPool *uint64) int32
	VkDestroyDescriptorPool                func(device uintptr, descriptorPool uint64, pAllocator uintptr)
	VkAllocateDescriptorSets               func(device uintptr, pAllocateInfo *VkDescriptorSetAllocateInfo, pDescriptorSets *uint64) int32
	VkFreeDescriptorSets                   func(device uintptr, descriptorPool uint64, descriptorSetCount uint32, pDescriptorSets *uint64) int32
	VkUpdateDescriptorSets                 func(device uintptr, descriptorWriteCount uint32, pDescriptorWrites *VkWriteDescriptorSet, descriptorCopyCount uint32, pDescriptorCopies unsafe.Pointer)
	VkCreateGraphicsPipelines              func(device uintptr, pipelineCache uint64, createInfoCount uint32, pCreateInfos *VkGraphicsPipelineCreateInfo, pAllocator uintptr, pPipelines *uint64) int32
	VkCreateComputePipelines               func(device uintptr, pipelineCache uint64, createInfoCount uint32, pCreateInfos *VkComputePipelineCreateInfo, pAllocator uintptr, pPipelines *uint64) int32
	VkDestroyPipeline                      func(device uintptr, pipeline uint64, pAllocator uintptr)
	VkDestroySurfaceKHR                    func(instance uintptr, surface uint64, pAllocator uintptr)

	GlfwVulkanSupported               func() int32
	GlfwGetRequiredInstanceExtensions func(count *uint32) **byte
	GlfwCreateWindowSurface           func(instance uintptr, window uintptr, allocator uintptr, surface *uint64) int32

	VkGetSwapchainImagesKHR              func(device uintptr, swapchain uint64, pSwapchainImageCount *uint32, pSwapchainImages *uint64) int32
	VkCreateSampler                      func(device uintptr, pCreateInfo *VkSamplerCreateInfo, pAllocator uintptr, pSampler *uint64) int32
	VkDestroySampler                     func(device uintptr, sampler uint64, pAllocator uintptr)
	VkCmdSetViewport                     func(commandBuffer uintptr, firstViewport uint32, viewportCount uint32, pViewports *VkViewport)
	VkCmdSetScissor                      func(commandBuffer uintptr, firstScissor uint32, scissorCount uint32, pScissors *VkRect2D)
	VkWaitForFences                      func(device uintptr, fenceCount uint32, pFences *uint64, waitAll uint32, timeout uint64) int32
	VkResetFences                        func(device uintptr, fenceCount uint32, pFences *uint64) int32
	VkResetCommandBuffer                 func(commandBuffer uintptr, flags uint32) int32
	VkEnumerateDeviceExtensionProperties func(physicalDevice uintptr, pLayerName *byte, pPropertyCount *uint32, pProperties *VkExtensionProperties) int32
)

var VkGetInstanceProcAddr func(uintptr, *byte) uintptr
var VkGetDeviceProcAddr func(uintptr, *byte) uintptr
var VkFlushMappedMemoryRanges func(uintptr, uint32, *VkMappedMemoryRange) int32
var VkInvalidateMappedMemoryRanges func(uintptr, uint32, *VkMappedMemoryRange) int32
var VkEnumerateInstanceLayerProperties func(*uint32, *VkLayerProperties) int32

type VkMappedMemoryRange struct {
	SType        uint32
	PNext        uintptr
	Memory       uint64
	Offset, Size uint64
}
type VkLayerProperties struct {
	LayerName                          [256]byte
	SpecVersion, ImplementationVersion uint32
	Description                        [256]byte
}

func InitVulkanBindings() error {
	var names []string
	switch runtime.GOOS {
	case "windows":
		names = []string{"vulkan-1.dll"}
	case "darwin":
		names = []string{"libvulkan.1.dylib", "libvulkan.dylib", "libMoltenVK.dylib"}
	case "linux":
		names = []string{"libvulkan.so.1", "libvulkan.so"}
	default:
		return fmt.Errorf("unsupported OS %s", runtime.GOOS)
	}
	var err error
	vkLibHandle, err = loadLibrary(libraryCandidates("DER_ATOM_VULKAN", names...))
	if err != nil {
		return err
	}
	if err = bindLibrary(&VkGetInstanceProcAddr, vkLibHandle, "vkGetInstanceProcAddr"); err != nil {
		return err
	}
	if err = bindInstance(&VkCreateInstance, 0, "vkCreateInstance"); err != nil {
		return err
	}
	if err = bindInstance(&VkEnumerateInstanceExtensionProperties, 0, "vkEnumerateInstanceExtensionProperties"); err != nil {
		return err
	}
	return bindInstance(&VkEnumerateInstanceLayerProperties, 0, "vkEnumerateInstanceLayerProperties")
}
func bindInstance(target any, instance uintptr, name string) error {
	text := CString(name)
	address := VkGetInstanceProcAddr(instance, &text[0])
	runtime.KeepAlive(text)
	return bind(target, address, name)
}
func bindDevice(target any, device uintptr, name string) error {
	text := CString(name)
	address := VkGetDeviceProcAddr(device, &text[0])
	runtime.KeepAlive(text)
	return bind(target, address, name)
}
func LoadInstance(instance uintptr) error {
	var err error
	if err = bindInstance(&VkGetDeviceProcAddr, instance, "vkGetDeviceProcAddr"); err != nil {
		return err
	}
	if err = bindInstance(&VkDestroyInstance, instance, "vkDestroyInstance"); err != nil {
		return err
	}
	if err = bindInstance(&VkEnumeratePhysicalDevices, instance, "vkEnumeratePhysicalDevices"); err != nil {
		return err
	}
	if err = bindInstance(&VkGetPhysicalDeviceQueueFamilyProperties, instance, "vkGetPhysicalDeviceQueueFamilyProperties"); err != nil {
		return err
	}
	if err = bindInstance(&VkGetPhysicalDeviceSurfaceSupportKHR, instance, "vkGetPhysicalDeviceSurfaceSupportKHR"); err != nil {
		return err
	}
	if err = bindInstance(&VkGetPhysicalDeviceSurfaceCapabilitiesKHR, instance, "vkGetPhysicalDeviceSurfaceCapabilitiesKHR"); err != nil {
		return err
	}
	if err = bindInstance(&VkGetPhysicalDeviceSurfaceFormatsKHR, instance, "vkGetPhysicalDeviceSurfaceFormatsKHR"); err != nil {
		return err
	}
	if err = bindInstance(&VkGetPhysicalDeviceSurfacePresentModesKHR, instance, "vkGetPhysicalDeviceSurfacePresentModesKHR"); err != nil {
		return err
	}
	if err = bindInstance(&VkCreateDevice, instance, "vkCreateDevice"); err != nil {
		return err
	}
	if err = bindInstance(&VkGetPhysicalDeviceMemoryProperties, instance, "vkGetPhysicalDeviceMemoryProperties"); err != nil {
		return err
	}
	if err = bindInstance(&VkGetPhysicalDeviceFeatures, instance, "vkGetPhysicalDeviceFeatures"); err != nil {
		return err
	}
	if err = bindInstance(&VkGetPhysicalDeviceFormatProperties, instance, "vkGetPhysicalDeviceFormatProperties"); err != nil {
		return err
	}
	if err = bindInstance(&VkDestroySurfaceKHR, instance, "vkDestroySurfaceKHR"); err != nil {
		return err
	}
	if err = bindInstance(&VkEnumerateDeviceExtensionProperties, instance, "vkEnumerateDeviceExtensionProperties"); err != nil {
		return err
	}
	return nil
}
func LoadDevice(device uintptr) error {
	var err error
	if err = bindDevice(&VkCmdCopyImage, device, "vkCmdCopyImage"); err != nil {
		return err
	}
	if err = bindDevice(&VkDestroyDevice, device, "vkDestroyDevice"); err != nil {
		return err
	}
	if err = bindDevice(&VkDeviceWaitIdle, device, "vkDeviceWaitIdle"); err != nil {
		return err
	}
	if err = bindDevice(&VkCreateBuffer, device, "vkCreateBuffer"); err != nil {
		return err
	}
	if err = bindDevice(&VkDestroyBuffer, device, "vkDestroyBuffer"); err != nil {
		return err
	}
	if err = bindDevice(&VkGetBufferMemoryRequirements, device, "vkGetBufferMemoryRequirements"); err != nil {
		return err
	}
	if err = bindDevice(&VkGetImageMemoryRequirements, device, "vkGetImageMemoryRequirements"); err != nil {
		return err
	}
	if err = bindDevice(&VkAllocateMemory, device, "vkAllocateMemory"); err != nil {
		return err
	}
	if err = bindDevice(&VkFreeMemory, device, "vkFreeMemory"); err != nil {
		return err
	}
	if err = bindDevice(&VkBindBufferMemory, device, "vkBindBufferMemory"); err != nil {
		return err
	}
	if err = bindDevice(&VkBindImageMemory, device, "vkBindImageMemory"); err != nil {
		return err
	}
	if err = bindDevice(&VkMapMemory, device, "vkMapMemory"); err != nil {
		return err
	}
	if err = bindDevice(&VkUnmapMemory, device, "vkUnmapMemory"); err != nil {
		return err
	}
	if err = bindDevice(&VkCreateSwapchainKHR, device, "vkCreateSwapchainKHR"); err != nil {
		return err
	}
	if err = bindDevice(&VkDestroySwapchainKHR, device, "vkDestroySwapchainKHR"); err != nil {
		return err
	}
	if err = bindDevice(&VkCreateSemaphore, device, "vkCreateSemaphore"); err != nil {
		return err
	}
	if err = bindDevice(&VkDestroySemaphore, device, "vkDestroySemaphore"); err != nil {
		return err
	}
	if err = bindDevice(&VkCreateFence, device, "vkCreateFence"); err != nil {
		return err
	}
	if err = bindDevice(&VkDestroyFence, device, "vkDestroyFence"); err != nil {
		return err
	}
	if err = bindDevice(&VkCreateImage, device, "vkCreateImage"); err != nil {
		return err
	}
	if err = bindDevice(&VkDestroyImage, device, "vkDestroyImage"); err != nil {
		return err
	}
	if err = bindDevice(&VkCreateImageView, device, "vkCreateImageView"); err != nil {
		return err
	}
	if err = bindDevice(&VkDestroyImageView, device, "vkDestroyImageView"); err != nil {
		return err
	}
	if err = bindDevice(&VkCreateFramebuffer, device, "vkCreateFramebuffer"); err != nil {
		return err
	}
	if err = bindDevice(&VkDestroyFramebuffer, device, "vkDestroyFramebuffer"); err != nil {
		return err
	}
	if err = bindDevice(&VkCreateShaderModule, device, "vkCreateShaderModule"); err != nil {
		return err
	}
	if err = bindDevice(&VkDestroyShaderModule, device, "vkDestroyShaderModule"); err != nil {
		return err
	}
	if err = bindDevice(&VkAllocateCommandBuffers, device, "vkAllocateCommandBuffers"); err != nil {
		return err
	}
	if err = bindDevice(&VkFreeCommandBuffers, device, "vkFreeCommandBuffers"); err != nil {
		return err
	}
	if err = bindDevice(&VkAcquireNextImageKHR, device, "vkAcquireNextImageKHR"); err != nil {
		return err
	}
	if err = bindDevice(&VkBeginCommandBuffer, device, "vkBeginCommandBuffer"); err != nil {
		return err
	}
	if err = bindDevice(&VkEndCommandBuffer, device, "vkEndCommandBuffer"); err != nil {
		return err
	}
	if err = bindDevice(&VkQueueSubmit, device, "vkQueueSubmit"); err != nil {
		return err
	}
	if err = bindDevice(&VkQueueWaitIdle, device, "vkQueueWaitIdle"); err != nil {
		return err
	}
	if err = bindDevice(&VkQueuePresentKHR, device, "vkQueuePresentKHR"); err != nil {
		return err
	}
	if err = bindDevice(&VkCmdCopyBuffer, device, "vkCmdCopyBuffer"); err != nil {
		return err
	}
	if err = bindDevice(&VkCmdCopyBufferToImage, device, "vkCmdCopyBufferToImage"); err != nil {
		return err
	}
	if err = bindDevice(&VkCmdPipelineBarrier, device, "vkCmdPipelineBarrier"); err != nil {
		return err
	}
	if err = bindDevice(&VkCmdBeginRenderPass, device, "vkCmdBeginRenderPass"); err != nil {
		return err
	}
	if err = bindDevice(&VkCmdEndRenderPass, device, "vkCmdEndRenderPass"); err != nil {
		return err
	}
	if err = bindDevice(&VkCmdBindPipeline, device, "vkCmdBindPipeline"); err != nil {
		return err
	}
	if err = bindDevice(&VkCmdBindDescriptorSets, device, "vkCmdBindDescriptorSets"); err != nil {
		return err
	}
	if err = bindDevice(&VkCmdBindVertexBuffers, device, "vkCmdBindVertexBuffers"); err != nil {
		return err
	}
	if err = bindDevice(&VkCmdBindIndexBuffer, device, "vkCmdBindIndexBuffer"); err != nil {
		return err
	}
	if err = bindDevice(&VkCmdDispatch, device, "vkCmdDispatch"); err != nil {
		return err
	}
	if err = bindDevice(&VkCmdDrawIndexedIndirect, device, "vkCmdDrawIndexedIndirect"); err != nil {
		return err
	}
	if err = bindDevice(&VkCmdDrawIndexed, device, "vkCmdDrawIndexed"); err != nil {
		return err
	}
	if err = bindDevice(&VkCmdDraw, device, "vkCmdDraw"); err != nil {
		return err
	}
	if err = bindDevice(&VkCmdDrawInstanced, device, "vkCmdDraw"); err != nil {
		return err
	}
	if err = bindDevice(&VkCmdPushConstants, device, "vkCmdPushConstants"); err != nil {
		return err
	}
	if err = bindDevice(&VkCmdBlitImage, device, "vkCmdBlitImage"); err != nil {
		return err
	}
	if err = bindDevice(&VkCreateRenderPass, device, "vkCreateRenderPass"); err != nil {
		return err
	}
	if err = bindDevice(&VkDestroyRenderPass, device, "vkDestroyRenderPass"); err != nil {
		return err
	}
	if err = bindDevice(&VkCreateCommandPool, device, "vkCreateCommandPool"); err != nil {
		return err
	}
	if err = bindDevice(&VkDestroyCommandPool, device, "vkDestroyCommandPool"); err != nil {
		return err
	}
	if err = bindDevice(&VkGetDeviceQueue, device, "vkGetDeviceQueue"); err != nil {
		return err
	}
	if err = bindDevice(&VkCreateDescriptorSetLayout, device, "vkCreateDescriptorSetLayout"); err != nil {
		return err
	}
	if err = bindDevice(&VkDestroyDescriptorSetLayout, device, "vkDestroyDescriptorSetLayout"); err != nil {
		return err
	}
	if err = bindDevice(&VkCreatePipelineLayout, device, "vkCreatePipelineLayout"); err != nil {
		return err
	}
	if err = bindDevice(&VkDestroyPipelineLayout, device, "vkDestroyPipelineLayout"); err != nil {
		return err
	}
	if err = bindDevice(&VkCreateDescriptorPool, device, "vkCreateDescriptorPool"); err != nil {
		return err
	}
	if err = bindDevice(&VkDestroyDescriptorPool, device, "vkDestroyDescriptorPool"); err != nil {
		return err
	}
	if err = bindDevice(&VkAllocateDescriptorSets, device, "vkAllocateDescriptorSets"); err != nil {
		return err
	}
	if err = bindDevice(&VkFreeDescriptorSets, device, "vkFreeDescriptorSets"); err != nil {
		return err
	}
	if err = bindDevice(&VkUpdateDescriptorSets, device, "vkUpdateDescriptorSets"); err != nil {
		return err
	}
	if err = bindDevice(&VkCreateGraphicsPipelines, device, "vkCreateGraphicsPipelines"); err != nil {
		return err
	}
	if err = bindDevice(&VkCreateComputePipelines, device, "vkCreateComputePipelines"); err != nil {
		return err
	}
	if err = bindDevice(&VkDestroyPipeline, device, "vkDestroyPipeline"); err != nil {
		return err
	}
	if err = bindDevice(&VkGetSwapchainImagesKHR, device, "vkGetSwapchainImagesKHR"); err != nil {
		return err
	}
	if err = bindDevice(&VkCreateSampler, device, "vkCreateSampler"); err != nil {
		return err
	}
	if err = bindDevice(&VkDestroySampler, device, "vkDestroySampler"); err != nil {
		return err
	}
	if err = bindDevice(&VkCmdSetViewport, device, "vkCmdSetViewport"); err != nil {
		return err
	}
	if err = bindDevice(&VkCmdSetScissor, device, "vkCmdSetScissor"); err != nil {
		return err
	}
	if err = bindDevice(&VkWaitForFences, device, "vkWaitForFences"); err != nil {
		return err
	}
	if err = bindDevice(&VkResetFences, device, "vkResetFences"); err != nil {
		return err
	}
	if err = bindDevice(&VkResetCommandBuffer, device, "vkResetCommandBuffer"); err != nil {
		return err
	}
	if err = bindDevice(&VkCmdCopyImageToBuffer, device, "vkCmdCopyImageToBuffer"); err != nil {
		return err
	}
	if err = bindDevice(&VkFlushMappedMemoryRanges, device, "vkFlushMappedMemoryRanges"); err != nil {
		return err
	}
	if err = bindDevice(&VkInvalidateMappedMemoryRanges, device, "vkInvalidateMappedMemoryRanges"); err != nil {
		return err
	}
	return nil
}
func VkMakeVersion(major, minor, patch uint32) uint32 { return (major << 22) | (minor << 12) | patch }

type VkImageCopy struct {
	SrcSubresource VkImageSubresourceLayers
	SrcOffset      VkOffset3D
	DstSubresource VkImageSubresourceLayers
	DstOffset      VkOffset3D
	Extent         VkExtent3D
}

var VkCmdCopyImage func(uintptr, uint64, uint32, uint64, uint32, uint32, *VkImageCopy)

func CompileShader(device uintptr, sourcePath string) (uint64, error) {
	code, err := os.ReadFile(sourcePath)
	if err != nil {
		return 0, err
	}
	if len(code) < 20 || len(code)%4 != 0 || *(*uint32)(unsafe.Pointer(&code[0])) != 0x07230203 {
		return 0, fmt.Errorf("%s is not a valid SPIR-V module", sourcePath)
	}
	createInfo := VkShaderModuleCreateInfo{
		SType:    VK_STRUCTURE_TYPE_SHADER_MODULE_CREATE_INFO,
		CodeSize: uint64(len(code)),
		PCode:    (*uint32)(unsafe.Pointer(&code[0])),
	}
	var shaderModule uint64
	if res := VkCreateShaderModule(device, &createInfo, 0, &shaderModule); res != VK_SUCCESS {
		return 0, fmt.Errorf("failed to create shader module: %d", res)
	}
	return shaderModule, nil
}