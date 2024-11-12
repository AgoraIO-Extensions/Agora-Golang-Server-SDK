package agoraservice

/*
#cgo CFLAGS: -I../../agora_sdk/include/c/api2 -I../../agora_sdk/include/c/base

#include <stdlib.h>
#include <string.h>
#include "agora_base.h"
#include "agora_local_user.h"
#include "agora_rtc_conn.h"
#include "agora_service.h"
#include "agora_media_base.h"
#include "rtc_callbacks_cgo.h"
#include "audio_observer_cgo.h"
#include "video_observer_cgo.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

func CAgoraServiceConfig(cfg *AgoraServiceConfig) *C.struct__agora_service_config {
	ret := (*C.struct__agora_service_config)(C.malloc(C.sizeof_struct__agora_service_config))
	C.memset(unsafe.Pointer(ret), 0, C.sizeof_struct__agora_service_config)
	ret.enable_audio_processor = CIntFromBool(cfg.EnableAudioProcessor)
	ret.enable_audio_device = CIntFromBool(cfg.EnableAudioDevice)
	ret.enable_video = CIntFromBool(cfg.EnableVideo)

	ret.app_id = C.CString(cfg.AppId)
	ret.area_code = C.uint(cfg.AreaCode)
	ret.channel_profile = C.int(cfg.ChannelProfile)
	ret.audio_scenario = C.int(cfg.AudioScenario)
	ret.use_string_uid = CIntFromBool(cfg.UseStringUid)
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
	ret.auto_subscribe_audio = CIntFromBool(cfg.AutoSubscribeAudio)
	ret.auto_subscribe_video = CIntFromBool(cfg.AutoSubscribeVideo)
	ret.enable_audio_recording_or_playout = CIntFromBool(cfg.EnableAudioRecordingOrPlayout)
	ret.max_send_bitrate = C.int(cfg.MaxSendBitrate)
	ret.min_port = C.int(cfg.MinPort)
	ret.max_port = C.int(cfg.MaxPort)
	// ret.audio_subs_options
	ret.client_role_type = C.int(cfg.ClientRole)
	ret.channel_profile = C.int(cfg.ChannelProfile)
	ret.audio_recv_media_packet = CIntFromBool(cfg.AudioRecvMediaPacket)
	ret.video_recv_media_packet = CIntFromBool(cfg.VideoRecvMediaPacket)
	return ret
}

func FreeCRtcConnectionConfig(cfg *C.struct__rtc_conn_config) {
	C.free(unsafe.Pointer(cfg))
}

func GoRtcConnectionInfo(cInfo *C.struct__rtc_conn_info, info *RtcConnectionInfo) {
	info.ConnectionId = uint(cInfo.id)
	info.ChannelId = C.GoString(cInfo.channel_id)
	info.State = int(cInfo.state)
	info.LocalUserId = C.GoString(cInfo.local_user_id)
	info.InternalUid = uint(cInfo.internal_uid)
}

func CRtcConnectionObserver() *C.struct__rtc_conn_observer {
	ret := (*C.struct__rtc_conn_observer)(C.malloc(C.sizeof_struct__rtc_conn_observer))
	C.memset(unsafe.Pointer(ret), 0, C.sizeof_struct__rtc_conn_observer)
	ret.on_connected = (*[0]byte)(C.cgo_on_connected)
	ret.on_disconnected = (*[0]byte)(C.cgo_on_disconnected)
	ret.on_connecting = (*[0]byte)(C.cgo_on_connecting)
	ret.on_reconnecting = (*[0]byte)(C.cgo_on_reconnecting)
	ret.on_reconnected = (*[0]byte)(C.cgo_on_reconnected)
	ret.on_connection_lost = (*[0]byte)(C.cgo_on_connection_lost)
	ret.on_connection_failure = (*[0]byte)(C.cgo_on_connection_failure)
	ret.on_token_privilege_will_expire = (*[0]byte)(C.cgo_on_token_privilege_will_expire)
	ret.on_token_privilege_did_expire = (*[0]byte)(C.cgo_on_token_privilege_did_expire)
	ret.on_user_joined = (*[0]byte)(C.cgo_on_user_joined)
	ret.on_user_left = (*[0]byte)(C.cgo_on_user_left)
	ret.on_error = (*[0]byte)(C.cgo_on_error)
	ret.on_stream_message_error = (*[0]byte)(C.cgo_on_stream_message_error)
	return ret
}

func FreeCRtcConnectionObserver(handler *C.struct__rtc_conn_observer) {
	C.free(unsafe.Pointer(handler))
}

func CLocalUserObserver() *C.struct__local_user_observer {
	ret := (*C.struct__local_user_observer)(C.malloc(C.sizeof_struct__local_user_observer))
	C.memset(unsafe.Pointer(ret), 0, C.sizeof_struct__local_user_observer)
	ret.on_stream_message = (*[0]byte)(C.cgo_on_stream_message)
	ret.on_user_info_updated = (*[0]byte)(C.cgo_on_user_info_updated)
	ret.on_user_audio_track_subscribed = (*[0]byte)(C.cgo_on_user_audio_track_subscribed)
	ret.on_user_video_track_subscribed = (*[0]byte)(C.cgo_on_user_video_track_subscribed)
	ret.on_user_audio_track_state_changed = (*[0]byte)(C.cgo_on_user_audio_track_state_changed)
	ret.on_user_video_track_state_changed = (*[0]byte)(C.cgo_on_user_video_track_state_changed)
	ret.on_audio_volume_indication = (*[0]byte)(C.cgo_on_audio_volume_indication)
	return ret
}

func FreeCLocalUserObserver(handler *C.struct__local_user_observer) {
	C.free(unsafe.Pointer(handler))
}

func CAudioFrameObserver() *C.struct__audio_frame_observer {
	ret := (*C.struct__audio_frame_observer)(C.malloc(C.sizeof_struct__audio_frame_observer))
	C.memset(unsafe.Pointer(ret), 0, C.sizeof_struct__audio_frame_observer)
	ret.on_record_audio_frame = (*[0]byte)(C.cgo_on_record_audio_frame)
	ret.on_playback_audio_frame = (*[0]byte)(C.cgo_on_playback_audio_frame)
	ret.on_mixed_audio_frame = (*[0]byte)(C.cgo_on_mixed_audio_frame)
	ret.on_ear_monitoring_audio_frame = (*[0]byte)(C.cgo_on_ear_monitoring_audio_frame)
	ret.on_playback_audio_frame_before_mixing = (*[0]byte)(C.cgo_on_playback_audio_frame_before_mixing)
	ret.on_get_audio_frame_position = (*[0]byte)(C.cgo_on_get_audio_frame_position)
	ret.on_get_playback_audio_frame_param = (*[0]byte)(C.cgo_on_get_playback_audio_frame_param)
	ret.on_get_record_audio_frame_param = (*[0]byte)(C.cgo_on_get_record_audio_frame_param)
	ret.on_get_mixed_audio_frame_param = (*[0]byte)(C.cgo_on_get_mixed_audio_frame_param)
	ret.on_get_ear_monitoring_audio_frame_param = (*[0]byte)(C.cgo_on_get_ear_monitoring_audio_frame_param)
	return ret
}

func FreeCAudioFrameObserver(observer *C.struct__audio_frame_observer) {
	C.free(unsafe.Pointer(observer))
}

func GoPcmAudioFrame(frame *C.struct__audio_frame) *AudioFrame {
	bufferLen := frame.samples_per_channel * frame.bytes_per_sample * frame.channels
	ret := &AudioFrame{
		Type:              AudioFrameType(frame._type),
		SamplesPerChannel: int(frame.samples_per_channel),
		BytesPerSample:    int(frame.bytes_per_sample),
		Channels:          int(frame.channels),
		SamplesPerSec:     int(frame.samples_per_sec),
		Buffer:            C.GoBytes(unsafe.Pointer(frame.buffer), bufferLen),
		RenderTimeMs:      int64(frame.render_time_ms),
		AvsyncType:        int(frame.avsync_type),
		FarFieldFlag:      int(frame.far_filed_flag),
		Rms:               int(frame.rms),
		VoiceProb:         int(frame.voice_prob),
		MusicProb:         int(frame.music_prob),
		Pitch:             int(frame.pitch),
	}
	return ret
}

func CVideoFrameObserver() unsafe.Pointer {
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

func CVideoEncodedFrameObserver() unsafe.Pointer {
	ret := C.struct__video_encoded_frame_observer{}
	C.memset(unsafe.Pointer(&ret), 0, C.sizeof_struct__video_encoded_frame_observer)
	ret.on_encoded_video_frame = (*[0]byte)(C.cgo_on_encoded_video_frame)
	obs := C.agora_video_encoded_frame_observer_create(&ret)
	return obs
}

func FreeCEncodedVideoFrameObserver(observer unsafe.Pointer) {
	C.agora_video_encoded_frame_observer_destroy(observer)
}

func Uint8PtrToUintptr(p *C.uint8_t) uintptr {
	return uintptr(unsafe.Pointer(p))
}

func GoVideoFrame(frame *C.struct__video_frame) *VideoFrame {
	// var buf []byte = nil
	// bufLen := frame.y_stride*frame.height + frame.u_stride*frame.height/2 + frame.v_stride*frame.height/2
	// uStart := frame.y_stride * frame.height
	// vStart := uStart + frame.u_stride*frame.height/2
	// if (Uint8PtrToUintptr(frame.u_buffer)-Uint8PtrToUintptr(frame.y_buffer)) == uintptr(uStart) &&
	// 	(Uint8PtrToUintptr(frame.v_buffer)-Uint8PtrToUintptr(frame.y_buffer)) == uintptr(vStart) {
	// 	buf = C.GoBytes(unsafe.Pointer(frame.y_buffer), bufLen)
	// } else {
	// 	buf = make([]byte, bufLen)
	// 	copy(buf, unsafe.Slice((*byte)(unsafe.Pointer(frame.y_buffer)), frame.y_stride*frame.height))
	// 	copy(buf[uStart:], unsafe.Slice((*byte)(unsafe.Pointer(frame.u_buffer)), frame.u_stride*frame.height/2))
	// 	copy(buf[vStart:], unsafe.Slice((*byte)(unsafe.Pointer(frame.v_buffer)), frame.v_stride*frame.height/2))
	// }
	yLen := frame.y_stride * frame.height
	uLen := frame.u_stride * frame.height / 2
	vLen := frame.v_stride * frame.height / 2
	ret := &VideoFrame{
		Type:           VideoBufferType(frame._type),
		Width:          int(frame.width),
		Height:         int(frame.height),
		YStride:        int(frame.y_stride),
		UStride:        int(frame.u_stride),
		VStride:        int(frame.v_stride),
		YBuffer:        C.GoBytes(unsafe.Pointer(frame.y_buffer), yLen),
		UBuffer:        C.GoBytes(unsafe.Pointer(frame.u_buffer), uLen),
		VBuffer:        C.GoBytes(unsafe.Pointer(frame.v_buffer), vLen),
		Rotation:       VideoOrientation(frame.rotation),
		RenderTimeMs:   int64(frame.render_time_ms),
		AVSyncType:     int(frame.avsync_type),
		MetadataBuffer: C.GoBytes(unsafe.Pointer(frame.metadata_buffer), frame.metadata_size),
		SharedContext:  frame.shared_context,
		TextureID:      int(frame.texture_id),
		Matrix:         [16]float32{},
		AlphaBuffer:    nil,
	}
	for i := 0; i < 16; i++ {
		ret.Matrix[i] = float32(frame.matrix[i])
	}
	if frame.alpha_buffer != nil {
		ret.AlphaBuffer = C.GoBytes(unsafe.Pointer(frame.alpha_buffer), frame.y_stride*frame.height)
	}
	return ret
}

func GoVideoTrackInfo(cInfo *C.struct__video_track_info) *VideoTrackInfo {
	ret := &VideoTrackInfo{
		IsLocal:             (int(cInfo.is_local) != 0),
		OwnerUid:            uint(cInfo.owner_uid),
		TrackId:             uint(cInfo.track_id),
		ChannelId:           C.GoString(cInfo.channel_id),
		StreamType:          int(cInfo.stream_type),
		CodecType:           int(cInfo.codec_type),
		EncodedFrameOnly:    (int(cInfo.encoded_frame_only) != 0),
		SourceType:          int(cInfo.source_type),
		ObservationPosition: uint(cInfo.observation_position),
	}
	return ret
}

func unsafeCBytes(data []byte) (unsafe.Pointer, runtime.Pinner) {
	ptr := unsafe.Pointer(&data[0])

	var pinner runtime.Pinner
	pinner.Pin(ptr)

	return ptr, pinner
}
