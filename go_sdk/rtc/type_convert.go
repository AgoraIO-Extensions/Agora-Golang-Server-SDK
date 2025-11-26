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
	ret.domain_limit = C.int(cfg.DomainLimit)

	// set log  related parameters
	if cfg.LogPath != "" {
		ret.log_file_path = C.CString(cfg.LogPath)
	}
	
	ret.log_file_size_kb = C.uint32_t(cfg.LogSize)
	ret.log_level = C.int(cfg.LogLevel)
	
	if cfg.ConfigDir != "" {
		ret.config_dir = C.CString(cfg.ConfigDir)
	}
	if cfg.DataDir != "" {
		ret.data_dir = C.CString(cfg.DataDir)
	}


	return ret
}

func FreeCAgoraServiceConfig(cfg *C.struct__agora_service_config) {
	C.free(unsafe.Pointer(cfg.app_id))
	if cfg.log_file_path != nil {
	C.free(unsafe.Pointer(cfg.log_file_path))
	}
	if cfg.config_dir != nil {
		C.free(unsafe.Pointer(cfg.config_dir))
	}
	if cfg.data_dir != nil {
		C.free(unsafe.Pointer(cfg.data_dir))
	}
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
	ret.on_encryption_error = (*[0]byte)(C.cgo_on_encryption_error)
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
	ret.on_audio_publish_state_changed = (*[0]byte)(C.cgo_on_audio_publish_state_changed)
	// todo: check if need to use this, still can not work now,should check why can not work
	ret.on_audio_volume_indication = (*[0]byte)(C.cgo_on_audio_volume_indication)
	// not finished ,only a title
	// void (*on_audio_meta_data_received)(AGORA_HANDLE agora_local_user, user_id_t userId, const char* meta_data, size_t length);
	ret.on_audio_meta_data_received = (*[0]byte)(C.cgo_on_audio_meta_data_received)
	// version 2.2.2 to expore statistics
	ret.on_local_audio_track_statistics = (*[0]byte)(C.cgo_on_local_audio_track_statistics)
	ret.on_remote_audio_track_statistics = (*[0]byte)(C.cgo_on_remote_audio_track_statistics)
	ret.on_local_video_track_statistics = (*[0]byte)(C.cgo_on_local_video_track_statistics)
	ret.on_remote_video_track_statistics = (*[0]byte)(C.cgo_on_remote_video_track_statistics)
	// added on 2025-06-09
	ret.on_audio_track_publish_success = (*[0]byte)(C.cgo_on_audio_track_publish_success)
	ret.on_audio_track_unpublished = (*[0]byte)(C.cgo_on_audio_track_unpublished)
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
		PresentTimeMs:     int64(frame.presentation_ms), // NOTE: next version, should include pts in audio_frame in c api layer!!??
	}
	return ret
}

func GoSinkAudioFrame(frame *C.struct__audio_pcm_frame) *AudioFrame {
	bufferLen := int(frame.samples_per_channel) * int(frame.bytes_per_sample) * int(frame.num_channels)
	samplepersec := int(frame.sample_rate_hz) * int(frame.num_channels)
	buffer := C.GoBytes(unsafe.Pointer(&frame.data[0]), C.int(bufferLen))
	ret := &AudioFrame{
		Type:              AudioFrameType(AudioFrameTypePCM16),
		SamplesPerChannel: int(frame.samples_per_channel),
		BytesPerSample:    int(frame.bytes_per_sample),
		Channels:          int(frame.num_channels),
		SamplesPerSec:     samplepersec,
		Buffer:            buffer,
		RenderTimeMs:      int64(frame.capture_timestamp),
		AvsyncType:        int(0),
		FarFieldFlag:      int(-1),
		Rms:               int(frame.audio_label.rms),
		VoiceProb:         int(frame.audio_label.voice_prob),
		MusicProb:         int(frame.audio_label.music_prob),
		Pitch:             int(frame.audio_label.pitch),
		PresentTimeMs:     int64(frame.capture_timestamp), // NOTE: next version, should include pts in audio_frame in c api layer!!??
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

func GoAudioVolumeInfo(frame *C.struct__audio_volume_info) *AudioVolumeInfo {
	if frame == nil {
		return nil
	}

	ret := &AudioVolumeInfo{
		UserId:     C.GoString(frame.user_id),
		Volume:     uint32(frame.volume),
		VAD:        uint32(frame.vad),
		VoicePitch: float64(frame.voicePitch),
	}
	return ret
}



func GoLocalAudioStats(stats *C.struct__local_audio_stats) *LocalAudioTrackStats {
	ret := &LocalAudioTrackStats{
		NumChannels: int(stats.num_channels),
		SentSampleRate: int(stats.sent_sample_rate),
		SentBitrate: int(stats.sent_bitrate),
		InternalCodec: int(stats.internal_codec),
		VoicePitch: float64(stats.voice_pitch),
	}
	return ret
}
func GoRemoteAudioStats(stats *C.struct__remote_audio_stats) *RemoteAudioTrackStats {
	ret := &RemoteAudioTrackStats{
		Uid: uint(stats.uid),
		Quality: int(stats.quality),
		NetworkTransportDelay: int(stats.network_transport_delay),
		JitterBufferDelay: int(stats.jitter_buffer_delay),
		AudioLossRate: int(stats.audio_loss_rate),
		NumChannels: int(stats.num_channels),
		ReceivedSampleRate: int(stats.received_sample_rate),
		ReceivedBitrate: int(stats.received_bitrate),
		TotalFrozenTime: int(stats.total_frozen_time), // ms
		FrozenRate: int(stats.frozen_rate),
		MosValue: int(stats.mos_value),
		TotalActiveTime: int(stats.total_active_time), // ms
		PublishDuration: int(stats.publish_duration), // ms
	}
	return ret	
}
func GoLocalVideoStats(stats *C.struct__local_video_track_stats) *LocalVideoTrackStats {
	ret := &LocalVideoTrackStats{
		NumberOfStreams: uint64(stats.number_of_streams),
		BytesMajorStream: uint64(stats.bytes_major_stream),
		BytesMinorStream: uint64(stats.bytes_minor_stream),
		FramesEncoded: uint32(stats.frames_encoded),
		SSRCMajorStream: uint32(stats.ssrc_major_stream),
		SSRCMinorStream: uint32(stats.ssrc_minor_stream),
		CaptureFrameRate: int(stats.capture_frame_rate),
		RegulatedCaptureFrameRate: int(stats.regulated_capture_frame_rate),
		InputFrameRate: int(stats.input_frame_rate),
		EncodeFrameRate: int(stats.encode_frame_rate),
		RenderFrameRate: int(stats.render_frame_rate),
		TargetMediaBitrateBps: int(stats.target_media_bitrate_bps),
		MediaBitrateBps: int(stats.media_bitrate_bps),
		TotalBitrateBps: int(stats.total_bitrate_bps),
		CaptureWidth: int(stats.capture_width),
		CaptureHeight: int(stats.capture_height),
		RegulatedCaptureWidth: int(stats.regulated_capture_width),
		RegulatedCaptureHeight: int(stats.regulated_capture_height),
		Width: int(stats.width),
		Height: int(stats.height),
		EncoderType: uint32(stats.encoder_type),
		UplinkCostTimeMs: int(stats.uplink_cost_time_ms),
		QualityAdaptIndication: int(stats.quality_adapt_indication),
	}
	return ret
}
func GoRemoteVideoStats(stats *C.struct__remote_video_track_stats) *RemoteVideoTrackStats {
	ret := &RemoteVideoTrackStats{
		Uid: uint(stats.uid),
		Delay: int(stats.delay),
		Width: int(stats.width),
		Height: int(stats.height),
		ReceivedBitrate: int(stats.received_bitrate),
		DecoderOutputFrameRate: int(stats.decoder_output_frame_rate),
		RendererOutputFrameRate: int(stats.renderer_output_frame_rate),
		FrameLossRate: int(stats.frame_loss_rate),
		PacketLossRate: int(stats.packet_loss_rate),
		RxStreamType: int(stats.rx_stream_type),
		TotalFrozenTime: int(stats.total_frozen_time), // ms
		FrozenRate: int(stats.frozen_rate),
		TotalDecodedFrames: uint32(stats.total_decoded_frames),
		AvSyncTimeMs: int(stats.av_sync_time_ms),
		DownlinkProcessTimeMs: int(stats.downlink_process_time_ms),
		FrameRenderDelayMs: int(stats.frame_render_delay_ms),
		TotalActiveTime: uint64(stats.totalActiveTime), // ms
		PublishDuration: uint64(stats.publishDuration), // ms
	}
	return ret
}

func CCapatilitiesObserver() *C.struct__capabilites_observer {
	ret := (*C.struct__capabilites_observer)(C.malloc(C.sizeof_struct__capabilites_observer))
	C.memset(unsafe.Pointer(ret), 0, C.sizeof_struct__capabilites_observer)
	ret.on_capabilities_changed = (*[0]byte)(C.cgo_on_capabilities_changed)
	return ret
}

func FreeCCapatilitiesObserver(observer *C.struct__capabilites_observer) {
	C.free(unsafe.Pointer(observer))
}
