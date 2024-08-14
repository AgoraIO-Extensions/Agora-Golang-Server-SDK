package agoraservice

// #cgo CFLAGS: -I${SRCDIR}/../../agora_sdk/include_c/api2 -I${SRCDIR}/../../agora_sdk/include_c/base
// #cgo darwin LDFLAGS: -Wl,-rpath,../../agora_sdk_mac -L../../agora_sdk_mac -lAgoraRtcKit -lAgorafdkaac -lAgoraffmpeg
// #cgo linux LDFLAGS: -L../../agora_sdk/ -lagora_rtc_sdk -lagora-fdkaac -lagora-core
// #include "agora_local_user.h"
// #include "agora_rtc_conn.h"
// #include "agora_service.h"
// #include "agora_media_base.h"
import "C"
import (
	"sync"
	"unsafe"
)

const (
	/**
	 * 0: (Recommended) The default audio scenario.
	 */
	AudioScenarioDefault = 0
	/**
	 * 3: (Recommended) The live gaming scenario, which needs to enable gaming
	 * audio effects in the speaker. Choose this scenario to achieve high-fidelity
	 * music playback.
	 */
	AudioScenarioGameStreaming = 3
	/**
	 * 5: The chatroom scenario, which needs to keep recording when setClientRole to audience.
	 * Normally, app developer can also use mute api to achieve the same result,
	 * and we implement this 'non-orthogonal' behavior only to make API backward compatible.
	 */
	AudioScenarioChatRoom = 5
	/**
	 * 7: Chorus
	 */
	AudioScenarioChorus = 7
	/**
	 * 8: Meeting
	 */
	AudioScenarioMeeting = 8
)

type AgoraServiceConfig struct {
	AppId         string
	AudioScenario int
	LogPath       string
	LogSize       int
}

type AgoraService struct {
	inited               bool
	service              unsafe.Pointer
	mediaFactory         unsafe.Pointer
	consByCCon           map[unsafe.Pointer]*RtcConnection
	consByCLocalUser     map[unsafe.Pointer]*RtcConnection
	consByCVideoObserver map[unsafe.Pointer]*RtcConnection
	connectionRWMutex    *sync.RWMutex
}

func newAgoraService() *AgoraService {
	return &AgoraService{
		inited:               false,
		service:              nil,
		mediaFactory:         nil,
		consByCCon:           make(map[unsafe.Pointer]*RtcConnection),
		consByCLocalUser:     make(map[unsafe.Pointer]*RtcConnection),
		consByCVideoObserver: make(map[unsafe.Pointer]*RtcConnection),
		connectionRWMutex:    &sync.RWMutex{},
	}
}

var agoraService *AgoraService = newAgoraService()

func Initialize(cfg *AgoraServiceConfig) int {
	if agoraService.inited {
		return 0
	}
	if agoraService.service == nil {
		agoraService.service = C.agora_service_create()
		if agoraService.service == nil {
			return -1
		}
	}

	ccfg := CAgoraServiceConfig(cfg)
	defer FreeCAgoraServiceConfig(ccfg)

	ret := int(C.agora_service_initialize(agoraService.service, ccfg))
	if ret != 0 {
		return ret
	}

	agoraService.mediaFactory = C.agora_service_create_media_node_factory(agoraService.service)

	if cfg.LogPath != "" {
		logPath := C.CString(cfg.LogPath)
		defer C.free(unsafe.Pointer(logPath))
		logSize := 512 * 1024
		if cfg.LogSize > 0 {
			logSize = cfg.LogSize
		}
		C.agora_service_set_log_file(agoraService.service, logPath,
			C.uint(logSize))
	}
	agoraService.inited = true
	return 0
}

func Release() int {
	if !agoraService.inited {
		return 0
	}
	if agoraService.service != nil {
		ret := int(C.agora_service_release(agoraService.service))
		if ret != 0 {
			return ret
		}
		agoraService.service = nil
	}
	agoraService.inited = false
	return 0
}

type AudioPcmDataSender struct {
	cSender unsafe.Pointer
}

type AudioEncodedFrameSender struct {
	cSender unsafe.Pointer
}

type LocalAudioTrack struct {
	cTrack unsafe.Pointer
}

type VideoFrameSender struct {
	cSender unsafe.Pointer
}

type VideoEncodedImageSender struct {
	cSender unsafe.Pointer
}

type LocalVideoTrack struct {
	cTrack unsafe.Pointer
}

const (
	VideoSendCcEnabled  = 0
	VideoSendCcDisabled = 1
)

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

type VideoEncodedImageSenderOptions struct {
	CcMode    int
	CodecType int
	/**
	 * Target bitrate (Kbps) for video encoding.
	 */
	TargetBitrate int
}

func NewCustomPcmAudioTrack(pcmSender *AudioPcmDataSender) *LocalAudioTrack {
	cTrack := C.agora_service_create_custom_audio_track_pcm(agoraService.service, pcmSender.cSender)
	if cTrack == nil {
		return nil
	}
	return &LocalAudioTrack{
		cTrack: cTrack,
	}
}

const (
	AudioTrackMixEnabled  = 0
	AudioTrackMixDisabled = 1
)

func NewCustomEncodedAudioTrack(encodedAudioSender *AudioEncodedFrameSender, mixMode int) *LocalAudioTrack {
	cTrack := C.agora_service_create_custom_audio_track_encoded(agoraService.service, encodedAudioSender.cSender, mixMode)
	if cTrack == nil {
		return nil
	}
	return &LocalAudioTrack{
		cTrack: cTrack,
	}
}

func (track *LocalAudioTrack) Release() {
	if track.cTrack == nil {
		return
	}
	C.agora_local_audio_track_destroy(track.cTrack)
	track.cTrack = nil
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

func NewAudioPcmDataSender() *AudioPcmDataSender {
	sender := C.agora_media_node_factory_create_audio_pcm_data_sender(agoraService.mediaFactory)
	if sender == nil {
		return nil
	}
	return &AudioPcmDataSender{
		cSender: sender,
	}
}

func (sender *AudioPcmDataSender) Release() {
	if sender.cSender == nil {
		return
	}
	C.agora_audio_pcm_data_sender_destroy(sender.cSender)
	sender.cSender = nil
}

func NewAudioEncodedFrameSender() *AudioEncodedFrameSender {
	sender := C.agora_media_node_factory_create_audio_encoded_frame_sender(agoraService.mediaFactory)
	if sender == nil {
		return nil
	}
	return &AudioEncodedFrameSender{
		cSender: sender,
	}
}

func (sender *AudioEncodedFrameSender) Release() {
	if sender.cSender == nil {
		return
	}
	C.agora_audio_encoded_frame_sender_destroy(sender.cSender)
	sender.cSender = nil
}

func NewVideoFrameSender() *VideoFrameSender {
	sender := C.agora_media_node_factory_create_video_frame_sender(agoraService.mediaFactory)
	if sender == nil {
		return nil
	}
	return &VideoFrameSender{
		cSender: sender,
	}
}

func (sender *VideoFrameSender) Release() {
	if sender.cSender == nil {
		return
	}
	C.agora_video_frame_sender_destroy(sender.cSender)
	sender.cSender = nil
}

func NewVideoEncodedImageSender() *VideoEncodedImageSender {
	sender := C.agora_media_node_factory_create_video_encoded_image_sender(agoraService.mediaFactory)
	if sender == nil {
		return nil
	}
	return &VideoEncodedImageSender{
		cSender: sender,
	}
}

func (sender *VideoEncodedImageSender) Release() {
	if sender.cSender == nil {
		return
	}
	C.agora_video_encoded_image_sender_destroy(sender.cSender)
	sender.cSender = nil
}
