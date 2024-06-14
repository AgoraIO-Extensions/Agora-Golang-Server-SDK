package agoraservice

/*
#cgo CFLAGS: -I../../agora_sdk/include_c/api2 -I../../agora_sdk/include_c/base
#cgo LDFLAGS: -L../../agora_sdk/ -lagora_rtc_sdk -lagora-fdkaac -lagora-ffmpeg

#include <stdlib.h>
#include <string.h>
#include "agora_local_user.h"
#include "agora_rtc_conn.h"
#include "agora_service.h"
#include "agora_media_base.h"
#include "rtc_callbacks_cgo.h"
#include "audio_observer_cgo.h"
*/
import "C"
import "unsafe"

func CAgoraServiceConfig(cfg *AgoraServiceConfig) *C.struct__agora_service_config {
	ret := (*C.struct__agora_service_config)(C.malloc(C.sizeof_struct__agora_service_config))
	C.memset(unsafe.Pointer(ret), 0, C.sizeof_struct__agora_service_config)
	ret.app_id = C.CString(cfg.AppId)
	ret.enable_audio_device = C.int(0)
	ret.enable_audio_processor = C.int(1)
	ret.enable_video = C.int(1)
	return ret
}

func FreeCAgoraServiceConfig(cfg *C.struct__agora_service_config) {
	C.free(unsafe.Pointer(cfg.app_id))
	C.free(unsafe.Pointer(cfg))
}

func CIntFromBool(b bool) C.int {
	if b {
		return 1
	}
	return 0
}

func CRtcConnectionConfig(cfg *RtcConnectionConfig) *C.struct__rtc_conn_config {
	ret := (*C.struct__rtc_conn_config)(C.malloc(C.sizeof_struct__rtc_conn_config))
	C.memset(unsafe.Pointer(ret), 0, C.sizeof_struct__rtc_conn_config)
	ret.auto_subscribe_audio = CIntFromBool(cfg.SubAudio)
	ret.auto_subscribe_video = CIntFromBool(cfg.SubVideo)
	ret.client_role_type = C.int(cfg.ClientRole)
	ret.channel_profile = C.int(cfg.ChannelProfile)
	return ret
}

func FreeCRtcConnectionConfig(cfg *C.struct__rtc_conn_config) {
	C.free(unsafe.Pointer(cfg))
}

func GoRtcConnectionInfo(cInfo *C.struct__rtc_conn_info) *RtcConnectionInfo {
	ret := &RtcConnectionInfo{
		ConnectionId: uint(cInfo.id),
		ChannelId:    C.GoString(cInfo.channel_id),
		State:        int(cInfo.state),
		LocalUserId:  C.GoString(cInfo.local_user_id),
		InternalUid:  uint(cInfo.internal_uid),
	}
	return ret
}

func CRtcConnectionEventHandler(handler *RtcConnectionEventHandler) *C.struct__rtc_conn_observer {
	ret := (*C.struct__rtc_conn_observer)(C.malloc(C.sizeof_struct__rtc_conn_observer))
	C.memset(unsafe.Pointer(ret), 0, C.sizeof_struct__rtc_conn_observer)
	ret.on_connected = (*[0]byte)(C.cgo_on_connected)
	return ret
}

func FreeCRtcConnectionEventHandler(handler *C.struct__rtc_conn_observer) {
	C.free(unsafe.Pointer(handler))
}

func CAudioFrameObserver(observer *RtcConnectionAudioFrameObserver) *C.struct__audio_frame_observer {
	ret := (*C.struct__audio_frame_observer)(C.malloc(C.sizeof_struct__audio_frame_observer))
	C.memset(unsafe.Pointer(ret), 0, C.sizeof_struct__audio_frame_observer)
	ret.on_playback_audio_frame_before_mixing = (*[0]byte)(C.cgo_on_playback_audio_frame_before_mixing)
	return ret
}

func FreeCAudioFrameObserver(observer *C.struct__audio_frame_observer) {
	C.free(unsafe.Pointer(observer))
}

func GoPcmAudioFrame(frame *C.struct__audio_frame) *PcmAudioFrame {
	bufferLen := frame.samples_per_channel * frame.bytes_per_sample * frame.channels
	ret := &PcmAudioFrame{
		Data:              C.GoBytes(unsafe.Pointer(frame.buffer), bufferLen),
		Timestamp:         int64(frame.render_time_ms),
		SamplesPerChannel: int(frame.samples_per_channel),
		BytesPerSample:    int(frame.bytes_per_sample),
		NumberOfChannels:  int(frame.channels),
		SampleRate:        int(frame.samples_per_sec),
	}
	return ret
}
