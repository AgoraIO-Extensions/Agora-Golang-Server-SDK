package agoraservice

/*
#cgo CFLAGS: -I../../agora_sdk/include/c/api2 -I../../agora_sdk/include/c/base

#include <string.h>
#include "agora_local_user.h"
#include "agora_rtc_conn.h"
#include "agora_service.h"
#include "agora_media_base.h"
*/
import "C"
import (
	"unsafe"
)

//export goOnRecordAudioFrame
func goOnRecordAudioFrame(cLocalUser unsafe.Pointer, channelId *C.char, frame *C.struct__audio_frame) C.int {
	//validity check
	if cLocalUser == nil {
		return C.int(0)
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cLocalUser, ConTypeCLocalUser)
	if con == nil || con.audioObserver == nil {
		return C.int(0)
	}
	goChannelId := C.GoString(channelId)
	return dispatchAudioFrameWithChannel(
		con.GetLocalUser(),
		goChannelId,
		frame,
		con.audioObserver.OnReusedRecordAudioFrame,
		con.audioObserver.OnRecordAudioFrame,
	)
}

//export goOnPlaybackAudioFrame
func goOnPlaybackAudioFrame(cLocalUser unsafe.Pointer, channelId *C.char, frame *C.struct__audio_frame) C.int {
	//validity check
	if cLocalUser == nil {
		return C.int(0)
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cLocalUser, ConTypeCLocalUser)

	if con == nil || con.audioObserver == nil {
		return C.int(0)
	}
	goChannelId := C.GoString(channelId)
	return dispatchAudioFrameWithChannel(
		con.GetLocalUser(),
		goChannelId,
		frame,
		con.audioObserver.OnReusedPlaybackAudioFrame,
		con.audioObserver.OnPlaybackAudioFrame,
	)
}

//export goOnMixedAudioFrame
func goOnMixedAudioFrame(cLocalUser unsafe.Pointer, channelId *C.char, frame *C.struct__audio_frame) C.int {
	//validity check
	if cLocalUser == nil {
		return C.int(0)
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cLocalUser, ConTypeCLocalUser)

	if con == nil || con.audioObserver == nil {
		return C.int(0)
	}
	goChannelId := C.GoString(channelId)
	return dispatchAudioFrameWithChannel(
		con.GetLocalUser(),
		goChannelId,
		frame,
		con.audioObserver.OnReusedMixedAudioFrame,
		con.audioObserver.OnMixedAudioFrame,
	)
}

//export goOnEarMonitoringAudioFrame
func goOnEarMonitoringAudioFrame(cLocalUser unsafe.Pointer, frame *C.struct__audio_frame) C.int {
	//validity check
	if cLocalUser == nil {
		return C.int(0)
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cLocalUser, ConTypeCLocalUser)

	if con == nil || con.audioObserver == nil {
		return C.int(0)
	}
	return dispatchAudioFrame(
		con.GetLocalUser(),
		frame,
		con.audioObserver.OnReusedEarMonitoringAudioFrame,
		con.audioObserver.OnEarMonitoringAudioFrame,
	)
}

//export goOnPlaybackAudioFrameBeforeMixing
func goOnPlaybackAudioFrameBeforeMixing(cLocalUser unsafe.Pointer, channelId *C.char, uid *C.char, frame *C.struct__audio_frame) C.int {
	//validity check
	if cLocalUser == nil {
		return C.int(0)
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cLocalUser, ConTypeCLocalUser)
	if con == nil || con.audioObserver == nil || con.audioObserver.OnPlaybackAudioFrameBeforeMixing == nil {
		return C.int(0)
	}
	goChannelId := C.GoString(channelId)
	goUid := C.GoString(uid)
	goFrame := GoPcmAudioFrame(frame)
	// add vad manager here
	var ret bool = false
	var vadResultFrame *AudioFrame = nil
	var vadResultStat VadState = VadStateInvalid

	if con.audioVadManager != nil {
		vadResultFrame, vadResultStat = con.audioVadManager.Process(goChannelId, goUid, goFrame)
	}
	ret = con.audioObserver.OnPlaybackAudioFrameBeforeMixing(con.GetLocalUser(), goChannelId, goUid, goFrame, vadResultStat, vadResultFrame)

	if ret {
		return C.int(1)
	}
	return C.int(0)
}

//export goOnGetAudioFramePosition
func goOnGetAudioFramePosition(cLocalUser unsafe.Pointer) C.int {
	//validity check
	if cLocalUser == nil {
		return C.int(0)
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cLocalUser, ConTypeCLocalUser)

	if con == nil || con.audioObserver == nil || con.audioObserver.OnGetAudioFramePosition == nil {
		return C.int(0)
	}
	return C.int(con.audioObserver.OnGetAudioFramePosition(con.GetLocalUser()))
}

//export goOnGetPlaybackAudioFrameParam
func goOnGetPlaybackAudioFrameParam(cLocalUser unsafe.Pointer) C.struct__audio_params {
	cAudioParam := C.struct__audio_params{}
	C.memset(unsafe.Pointer(&cAudioParam), 0, C.sizeof_struct__audio_params)

	//validity check
	if cLocalUser == nil {
		return cAudioParam
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cLocalUser, ConTypeCLocalUser)

	if con == nil || con.audioObserver == nil || con.audioObserver.OnGetPlaybackAudioFrameParam == nil {
		return cAudioParam
	}
	goAudioParam := con.audioObserver.OnGetPlaybackAudioFrameParam(con.GetLocalUser())
	cAudioParam.sample_rate = C.int(goAudioParam.SampleRate)
	cAudioParam.channels = C.int(goAudioParam.Channels)
	cAudioParam.mode = C.RAW_AUDIO_FRAME_OP_MODE_TYPE(goAudioParam.Mode)
	cAudioParam.samples_per_call = C.int(goAudioParam.SamplesPerCall)
	return cAudioParam
}

//export goOnGetRecordAudioFrameParam
func goOnGetRecordAudioFrameParam(cLocalUser unsafe.Pointer) C.struct__audio_params {

	cAudioParam := C.struct__audio_params{}
	C.memset(unsafe.Pointer(&cAudioParam), 0, C.sizeof_struct__audio_params)

	//validity check
	if cLocalUser == nil {
		return cAudioParam
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cLocalUser, ConTypeCLocalUser)
	if con == nil || con.audioObserver == nil || con.audioObserver.OnGetRecordAudioFrameParam == nil {
		return cAudioParam
	}
	goAudioParam := con.audioObserver.OnGetRecordAudioFrameParam(con.GetLocalUser())
	cAudioParam.sample_rate = C.int(goAudioParam.SampleRate)
	cAudioParam.channels = C.int(goAudioParam.Channels)
	cAudioParam.mode = C.RAW_AUDIO_FRAME_OP_MODE_TYPE(goAudioParam.Mode)
	cAudioParam.samples_per_call = C.int(goAudioParam.SamplesPerCall)
	return cAudioParam
}

//export goOnGetMixedAudioFrameParam
func goOnGetMixedAudioFrameParam(cLocalUser unsafe.Pointer) C.struct__audio_params {

	cAudioParam := C.struct__audio_params{}
	C.memset(unsafe.Pointer(&cAudioParam), 0, C.sizeof_struct__audio_params)

	//validity check
	if cLocalUser == nil {
		return cAudioParam
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cLocalUser, ConTypeCLocalUser)
	if con == nil || con.audioObserver == nil || con.audioObserver.OnGetMixedAudioFrameParam == nil {
		return cAudioParam
	}
	goAudioParam := con.audioObserver.OnGetMixedAudioFrameParam(con.GetLocalUser())
	cAudioParam.sample_rate = C.int(goAudioParam.SampleRate)
	cAudioParam.channels = C.int(goAudioParam.Channels)
	cAudioParam.mode = C.RAW_AUDIO_FRAME_OP_MODE_TYPE(goAudioParam.Mode)
	cAudioParam.samples_per_call = C.int(goAudioParam.SamplesPerCall)
	return cAudioParam
}

//export goOnGetEarMonitoringAudioFrameParam
func goOnGetEarMonitoringAudioFrameParam(cLocalUser unsafe.Pointer) C.struct__audio_params {
	cAudioParam := C.struct__audio_params{}
	C.memset(unsafe.Pointer(&cAudioParam), 0, C.sizeof_struct__audio_params)

	//validity check
	if cLocalUser == nil {
		return cAudioParam
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cLocalUser, ConTypeCLocalUser)
	if con == nil || con.audioObserver == nil || con.audioObserver.OnGetEarMonitoringAudioFrameParam == nil {
		return cAudioParam
	}
	goAudioParam := con.audioObserver.OnGetEarMonitoringAudioFrameParam(con.GetLocalUser())
	cAudioParam.sample_rate = C.int(goAudioParam.SampleRate)
	cAudioParam.channels = C.int(goAudioParam.Channels)
	cAudioParam.mode = C.RAW_AUDIO_FRAME_OP_MODE_TYPE(goAudioParam.Mode)
	cAudioParam.samples_per_call = C.int(goAudioParam.SamplesPerCall)
	return cAudioParam
}

func dispatchAudioFrameWithChannel(
	localUser *LocalUser,
	channelId string,
	frame *C.struct__audio_frame,
	reused func(localUser *LocalUser, channelId string, frame *AudioFrame) bool,
	copied func(localUser *LocalUser, channelId string, frame *AudioFrame) bool,
) C.int {
	if reused == nil && copied == nil {
		return C.int(0)
	}
	if reused != nil {
		if reused(localUser, channelId, GoPcmAudioFrameReused(frame)) {
			return C.int(1)
		}
		return C.int(0)
	}
	if copied(localUser, channelId, GoPcmAudioFrame(frame)) {
		return C.int(1)
	}
	return C.int(0)
}

func dispatchAudioFrame(
	localUser *LocalUser,
	frame *C.struct__audio_frame,
	reused func(localUser *LocalUser, frame *AudioFrame) bool,
	copied func(localUser *LocalUser, frame *AudioFrame) bool,
) C.int {
	if reused == nil && copied == nil {
		return C.int(0)
	}
	if reused != nil {
		if reused(localUser, GoPcmAudioFrameReused(frame)) {
			return C.int(1)
		}
		return C.int(0)
	}
	if copied(localUser, GoPcmAudioFrame(frame)) {
		return C.int(1)
	}
	return C.int(0)
}
