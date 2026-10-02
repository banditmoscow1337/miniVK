package binding

import (
	"fmt"
	"runtime"
)

type OpenAL struct {
	AlcOpenDevice          func(*byte) uintptr
	AlcCloseDevice         func(uintptr) uint8
	AlcCreateContext       func(uintptr, *int32) uintptr
	AlcDestroyContext      func(uintptr)
	AlcMakeContextCurrent  func(uintptr) uint8
	AlGenSources           func(int32, *uint32)
	AlDeleteSources        func(int32, *uint32)
	AlGenBuffers           func(int32, *uint32)
	AlDeleteBuffers        func(int32, *uint32)
	AlBufferData           func(uint32, int32, *byte, int32, int32)
	AlSourcei              func(uint32, int32, int32)
	AlSourcef              func(uint32, int32, float32)
	AlSourcePlay           func(uint32)
	AlSourcePause          func(uint32)
	AlSourceStop           func(uint32)
	AlSourceQueueBuffers   func(uint32, int32, *uint32)
	AlSourceUnqueueBuffers func(uint32, int32, *uint32)
	AlGetSourcei           func(uint32, int32, *int32)
	AlGetError             func() int32
	handle                 uintptr
}

func LoadOpenAL() (*OpenAL, error) {
	if runtime.GOARCH != "amd64" && runtime.GOARCH != "arm64" {
		return nil, fmt.Errorf("OpenAL ABI supports amd64 and arm64 only")
	}
	var names []string
	switch runtime.GOOS {
	case "windows":
		names = []string{"OpenAL32.dll"}
	case "darwin":
		names = []string{"libopenal.1.dylib", "libopenal.dylib", "/System/Library/Frameworks/OpenAL.framework/OpenAL"}
	case "linux":
		names = []string{"libopenal.so.1", "libopenal.so"}
	default:
		return nil, fmt.Errorf("OpenAL is unsupported on %s", runtime.GOOS)
	}
	handle, err := loadLibrary(libraryCandidates("DER_ATOM_OPENAL", names...))
	if err != nil {
		return nil, err
	}
	a := &OpenAL{handle: handle}
	for _, item := range []struct {
		target any
		name   string
	}{
		{&a.AlcOpenDevice, "alcOpenDevice"},
		{&a.AlcCloseDevice, "alcCloseDevice"},
		{&a.AlcCreateContext, "alcCreateContext"},
		{&a.AlcDestroyContext, "alcDestroyContext"},
		{&a.AlcMakeContextCurrent, "alcMakeContextCurrent"},
		{&a.AlGenSources, "alGenSources"},
		{&a.AlDeleteSources, "alDeleteSources"},
		{&a.AlGenBuffers, "alGenBuffers"},
		{&a.AlDeleteBuffers, "alDeleteBuffers"},
		{&a.AlBufferData, "alBufferData"},
		{&a.AlSourcei, "alSourcei"},
		{&a.AlSourcef, "alSourcef"},
		{&a.AlSourcePlay, "alSourcePlay"},
		{&a.AlSourcePause, "alSourcePause"},
		{&a.AlSourceStop, "alSourceStop"},
		{&a.AlSourceQueueBuffers, "alSourceQueueBuffers"},
		{&a.AlSourceUnqueueBuffers, "alSourceUnqueueBuffers"},
		{&a.AlGetSourcei, "alGetSourcei"},
		{&a.AlGetError, "alGetError"},
	} {
		if err := bindLibrary(item.target, handle, item.name); err != nil {
			a.Close()
			return nil, err
		}
	}
	return a, nil
}

func (a *OpenAL) Close() error {
	if a == nil || a.handle == 0 {
		return nil
	}
	err := closeLibrary(a.handle)
	a.handle = 0
	return err
}
