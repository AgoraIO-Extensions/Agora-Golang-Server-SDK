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
#include "video_observer_cgo.h"
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

func CRtcConnectionEventHandler(handler *RtcConnectionEventHandler) (*C.struct__rtc_conn_observer, *C.struct__local_user_observer) {
	ret := (*C.struct__rtc_conn_observer)(C.malloc(C.sizeof_struct__rtc_conn_observer))
	C.memset(unsafe.Pointer(ret), 0, C.sizeof_struct__rtc_conn_observer)
	ret.on_connected = (*[0]byte)(C.cgo_on_connected)
	ret.on_disconnected = (*[0]byte)(C.cgo_on_disconnected)
	ret.on_reconnecting = (*[0]byte)(C.cgo_on_reconnecting)
	ret.on_reconnected = (*[0]byte)(C.cgo_on_reconnected)
	ret.on_token_privilege_will_expire = (*[0]byte)(C.cgo_on_token_privilege_will_expire)
	ret.on_token_privilege_did_expire = (*[0]byte)(C.cgo_on_token_privilege_did_expire)
	ret.on_user_joined = (*[0]byte)(C.cgo_on_user_joined)
	ret.on_user_left = (*[0]byte)(C.cgo_on_user_left)
	ret.on_stream_message_error = (*[0]byte)(C.cgo_on_stream_message_error)

	ret1 := (*C.struct__local_user_observer)(C.malloc(C.sizeof_struct__local_user_observer))
	C.memset(unsafe.Pointer(ret1), 0, C.sizeof_struct__local_user_observer)
	ret1.on_stream_message = (*[0]byte)(C.cgo_on_stream_message)
	return ret, ret1
}

func FreeCRtcConnectionEventHandler(handler *C.struct__rtc_conn_observer) {
	C.free(unsafe.Pointer(handler))
}

func FreeCLocalUserObserver(handler *C.struct__local_user_observer) {
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

func CVideoFrameObserver(observer *RtcConnectionVideoFrameObserver) unsafe.Pointer {
	// ret := (*C.struct__video_frame_observer2)(C.malloc(C.sizeof_struct__video_frame_observer2))
	ret := C.struct__video_frame_observer2{}
	C.memset(unsafe.Pointer(&ret), 0, C.sizeof_struct__video_frame_observer2)
	ret.on_frame = (*[0]byte)(C.cgo_on_video_frame)
	obs := C.agora_video_frame_observer2_create(&ret)
	return obs
}

func FreeCVideoFrameObserver(observer unsafe.Pointer) {
	C.agora_video_frame_observer2_destroy(observer)
}

func Uint8PtrToUintptr(p *C.uint8_t) uintptr {
	return uintptr(unsafe.Pointer(p))
}

func GoVideoFrame(frame *C.struct__video_frame) *VideoFrame {
	var buf []byte = nil
	bufLen := frame.y_stride*frame.height + frame.u_stride*frame.height/2 + frame.v_stride*frame.height/2
	uStart := frame.y_stride * frame.height
	vStart := uStart + frame.u_stride*frame.height/2
	if (Uint8PtrToUintptr(frame.u_buffer)-Uint8PtrToUintptr(frame.y_buffer)) == uintptr(uStart) &&
		(Uint8PtrToUintptr(frame.v_buffer)-Uint8PtrToUintptr(frame.y_buffer)) == uintptr(vStart) {
		buf = C.GoBytes(unsafe.Pointer(frame.y_buffer), bufLen)
	} else {
		buf = make([]byte, bufLen)
		copy(buf, unsafe.Slice((*byte)(unsafe.Pointer(frame.y_buffer)), frame.y_stride*frame.height))
		copy(buf[uStart:], unsafe.Slice((*byte)(unsafe.Pointer(frame.u_buffer)), frame.u_stride*frame.height/2))
		copy(buf[vStart:], unsafe.Slice((*byte)(unsafe.Pointer(frame.v_buffer)), frame.v_stride*frame.height/2))
	}
	ret := &VideoFrame{
		Buffer:    buf,
		Width:     int(frame.width),
		Height:    int(frame.height),
		YStride:   int(frame.y_stride),
		UStride:   int(frame.u_stride),
		VStride:   int(frame.v_stride),
		Timestamp: int64(frame.render_time_ms),
	}
	return ret
}
