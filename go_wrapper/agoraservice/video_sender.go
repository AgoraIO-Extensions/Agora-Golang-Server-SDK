package agoraservice

/*
#cgo CFLAGS: -I../../agora_sdk/include_c/api2 -I../../agora_sdk/include_c/base

#include <string.h>
#include "agora_local_user.h"
#include "agora_rtc_conn.h"
#include "agora_service.h"
#include "agora_media_base.h"
#include "agora_video_track.h"
*/
import "C"
import "unsafe"

type VideoEncoderConfig struct {
	CodecType         int
	Width             int
	Height            int
	Framerate         int
	Bitrate           int
	MinBitrate        int
	OrientationMode   int
	DegradePreference int
}

type VideoSender struct {
	// config       *PcmSenderConfig
	cLocalUser   unsafe.Pointer
	cVideoTrack  unsafe.Pointer
	cVideoSender unsafe.Pointer
}

func (con *RtcConnection) GetVideoSender() *VideoSender {
	if con.videoSender != nil {
		return con.videoSender
	}
	ret := &VideoSender{
		cLocalUser:   con.cLocalUser,
		cVideoSender: C.agora_media_node_factory_create_video_frame_sender(agoraService.mediaFactory),
	}
	ret.cVideoTrack = C.agora_service_create_custom_video_track_frame(agoraService.service, ret.cVideoSender)
	con.videoSender = ret
	return ret
}

func (con *RtcConnection) ReleaseVideoSender() {
	sender := con.videoSender
	if sender == nil {
		return
	}
	if sender.cVideoSender == nil {
		con.videoSender = nil
		return
	}
	C.agora_local_video_track_destroy(sender.cVideoTrack)
	C.agora_video_frame_sender_destroy(sender.cVideoSender)
	con.videoSender = nil
}

func (sender *VideoSender) SetVideoEncoderConfig(cfg *VideoEncoderConfig) int {
	cCfg := C.struct__video_encoder_config{}
	C.memset(unsafe.Pointer(&cCfg), 0, C.sizeof_struct__video_encoder_config)
	cCfg.codec_type = C.int(cfg.CodecType)
	cCfg.dimensions.width = C.int(cfg.Width)
	cCfg.dimensions.height = C.int(cfg.Height)
	cCfg.frame_rate = C.int(cfg.Framerate)
	cCfg.bitrate = C.int(cfg.Bitrate)
	cCfg.min_bitrate = C.int(cfg.MinBitrate)
	cCfg.orientation_mode = C.int(cfg.OrientationMode)
	cCfg.degradation_preference = C.int(cfg.DegradePreference)
	return int(C.agora_local_video_track_set_video_encoder_config(sender.cVideoTrack, &cCfg))
}

func (sender *VideoSender) Start() int {
	C.agora_local_video_track_set_enabled(sender.cVideoTrack, C.int(1))
	return int(C.agora_local_user_publish_video(sender.cLocalUser, sender.cVideoTrack))
}

func (sender *VideoSender) Stop() int {
	ret := int(C.agora_local_user_unpublish_video(sender.cLocalUser, sender.cVideoTrack))
	C.agora_local_video_track_set_enabled(sender.cVideoTrack, C.int(0))
	return ret
}

func (sender *VideoSender) SendVideoFrame(frame *VideoFrame) int {
	cData := C.CBytes(frame.Buffer)
	defer C.free(cData)
	cFrame := C.struct__external_video_frame{}
	C.memset(unsafe.Pointer(&cFrame), 0, C.sizeof_struct__external_video_frame)
	cFrame._type = 1
	cFrame.format = 1
	cFrame.buffer = cData
	cFrame.stride = C.int(frame.YStride)
	cFrame.height = C.int(frame.Height)
	cFrame.timestamp = C.longlong(frame.Timestamp)
	return int(C.agora_video_frame_sender_send(sender.cVideoSender, &cFrame))
}
