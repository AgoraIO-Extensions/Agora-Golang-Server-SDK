package agoraservice

// #cgo CFLAGS: -I${SRCDIR}/../../agora_sdk/include/c/api2 -I${SRCDIR}/../../agora_sdk/include/c/base
// #include "agora_service.h"
// #include "agora_audio_track.h"
import "C"
import "unsafe"

type LocalAudioTrack struct {
	cTrack unsafe.Pointer
}

func NewCustomAudioTrackPcm(pcmSender *AudioPcmDataSender) *LocalAudioTrack {
	cTrack := C.agora_service_create_custom_audio_track_pcm(agoraService.service, pcmSender.cSender)
	if cTrack == nil {
		return nil
	}
	return &LocalAudioTrack{
		cTrack: cTrack,
	}
}

func NewCustomAudioTrackEncoded(encodedAudioSender *AudioEncodedFrameSender, mixMode AudioTrackMixingState) *LocalAudioTrack {
	cTrack := C.agora_service_create_custom_audio_track_encoded(agoraService.service, encodedAudioSender.cSender, C.int(mixMode))
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

func (track *LocalAudioTrack) SetEnabled(enable bool) {
	if track.cTrack == nil {
		return
	}
	cEnable := 0
	if enable {
		cEnable = 1
	}
	C.agora_local_audio_track_set_enabled(track.cTrack, C.int(cEnable))
}

func (track *LocalAudioTrack) AdjustPublishVolume(volume int) int {
	if track.cTrack == nil {
		return -1
	}
	return int(C.agora_local_audio_track_adjust_publish_volume(track.cTrack, C.int(volume)))
}

// NOTICE: these interface below is temporary, may be removed in the future
// size is the number of 10ms audio frames
// the default value of this param is 30, ie. 300ms
func (track *LocalAudioTrack) SetMaxBufferedAudioFrameNumber(frameNum int) {
	if track.cTrack == nil {
		return
	}
	C.agora_local_audio_track_set_max_bufferd_frame_number(track.cTrack, C.int(frameNum))
}

func (track *LocalAudioTrack) ClearSenderBuffer() int {
	if track.cTrack == nil {
		return -1
	}
	return int(C.agora_local_audio_track_clear_sender_buffer(track.cTrack))
}

func (track *LocalAudioTrack) SetSendDelayMs(delayMs int) int {
	if track.cTrack == nil {
		return -1
	}
	C.agora_local_audio_track_set_send_delay_ms(track.cTrack, C.int(delayMs))
	return 0
}
