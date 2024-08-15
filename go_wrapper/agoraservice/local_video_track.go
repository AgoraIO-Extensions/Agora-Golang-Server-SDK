package agoraservice

// #cgo CFLAGS: -I${SRCDIR}/../../agora_sdk/include_c/api2 -I${SRCDIR}/../../agora_sdk/include_c/base
// #include "agora_service.h"
// #include "agora_video_track.h"
import "C"
import "unsafe"

const (
	VideoCodecTypeNone        = 0
	VideoCodecTypeVp8         = 1
	VideoCodecTypeH264        = 2
	VideoCodecTypeH265        = 3
	VideoCodecTypeGeneric     = 6
	VideoCodecTypeGenericH264 = 7
	VideoCodecTypeAv1         = 12
	VideoCodecTypeVp9         = 13
	VideoCodecTypeGenericJpeg = 20
)

const (
	VideoSendCcEnabled  = 0
	VideoSendCcDisabled = 1
)

type VideoEncoderConfiguration struct {
	// 1: VP8, 2: H264
	CodecType int
	Width     int
	Height    int
	Framerate int
	// kbps
	Bitrate int
	// kbps
	MinBitrate int
	// 0: adaptive, 1: fixed landscape, 2: fixed portrait
	OrientationMode int
	// 0: maintain, 1: maintain frame rate, 2: maintain quality
	DegradePreference int
}

type VideoEncodedImageSenderOptions struct {
	CcMode    int
	CodecType int
	/**
	 * Target bitrate (Kbps) for video encoding.
	 */
	TargetBitrate int
}

type LocalVideoTrack struct {
	cTrack unsafe.Pointer
}

func NewCustomVideoTrack(videoSender *VideoFrameSender) *LocalVideoTrack {
	cTrack := C.agora_service_create_custom_video_track_frame(agoraService.service, videoSender.sender)
	if cTrack == nil {
		return nil
	}
	return &LocalVideoTrack{
		cTrack: cTrack,
	}
}

func NewCustomEncodedVideoTrack(videoSender *VideoEncodedImageSender, senderOptions *VideoEncodedImageSenderOptions) *LocalVideoTrack {
	cSenderOptions := C.sender_options{}
	cptrSenderOptions := &cSenderOptions
	if senderOptions != nil {
		cSenderOptions.cc_mode = C.int(senderOptions.CcMode)
		cSenderOptions.codec_type = C.int(senderOptions.CodecType)
		cSenderOptions.target_bitrate = C.int(senderOptions.TargetBitrate)
	} else {
		cptrSenderOptions = nil
	}
	cTrack := C.agora_service_create_custom_video_track_encoded(agoraService.service, videoSender.cSender, cptrSenderOptions)
	if cTrack == nil {
		return nil
	}
	return &LocalVideoTrack{
		cTrack: cTrack,
	}
}

func (track *LocalVideoTrack) Release() {
	if track.cTrack == nil {
		return
	}
	C.agora_local_video_track_destroy(track.cTrack)
	track.cTrack = nil
}

func (track *LocalVideoTrack) SetEnabled(enable bool) {
	if track.cTrack == nil {
		return
	}
	cEnable := 0
	if enable {
		cEnable = 1
	}
	C.agora_local_video_track_set_enabled(track.cTrack, C.int(cEnable))
}

func (track *LocalVideoTrack) SetVideoEncoderConfiguration(cfg *VideoEncoderConfiguration) int {
	cCfg := C.struct__video_encoder_config{}
	C.memset(unsafe.Pointer(&cCfg), 0, C.sizeof_struct__video_encoder_config)
	cCfg.codec_type = C.int(cfg.CodecType)
	cCfg.dimensions.width = C.int(cfg.Width)
	cCfg.dimensions.height = C.int(cfg.Height)
	cCfg.frame_rate = C.int(cfg.Framerate)
	cCfg.bitrate = C.int(cfg.Bitrate * 1000)
	cCfg.min_bitrate = C.int(cfg.MinBitrate * 1000)
	cCfg.orientation_mode = C.int(cfg.OrientationMode)
	cCfg.degradation_preference = C.int(cfg.DegradePreference)
	return int(C.agora_local_video_track_set_video_encoder_config(track.cTrack, &cCfg))
}
