package agoraservice

// #cgo CFLAGS: -I${SRCDIR}/../../agora_sdk/include/c/api2 -I${SRCDIR}/../../agora_sdk/include/c/base
// #include <stdlib.h>
// #include "agora_service.h"
// #include "agora_audio_track.h"
import "C"
import "unsafe"

type RemoteAudioTrack struct {
	cRemoteAudioTrack unsafe.Pointer
}

func NewRemoteAudioTrack(cRemoteAudioTrack unsafe.Pointer) *RemoteAudioTrack {
	return &RemoteAudioTrack{
		cRemoteAudioTrack: cRemoteAudioTrack,
	}
}

func (track *RemoteAudioTrack) EnableAudioFilter(name string, enable bool, position int) int {
	if track == nil || track.cRemoteAudioTrack == nil {
		return -1
	}
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cEnable := C.int(0)
	if enable {
		cEnable = C.int(1)
	}
	cPosition := C.int(position)

	return int(C.agora_audio_track_enable_audio_filter(track.cRemoteAudioTrack, cName, cEnable, cPosition))
}

func (track *RemoteAudioTrack) SetFilterProperty(name string, key string, value string, position int) int {
	if track == nil || track.cRemoteAudioTrack == nil {
		return -1
	}
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))
	cPosition := C.int(position)
	return int(C.agora_audio_track_set_filter_property(track.cRemoteAudioTrack, cName, cKey, cValue, cPosition))
}