package agoraservice

// #cgo CFLAGS: -I${SRCDIR}/../../agora_sdk/include/c/api2 -I${SRCDIR}/../../agora_sdk/include/c/base
// #include "agora_service.h"
// #include "agora_audio_track.h"
import "C"
import (
	"fmt"
	"unsafe"
)

type LocalAudioTrack struct {
	cTrack unsafe.Pointer
}

// NOTE: date：2025-06-27
// add audioScenario_of_connection param, to set the audio scenario for the audio track
// recommend to use the same audio scenario for the connection and related audio track
func NewCustomAudioTrackPcm(pcmSender *AudioPcmDataSender, audioScenario_of_connection AudioScenario, isSendExternalAudioForAI bool) *LocalAudioTrack {
	if agoraService == nil || agoraService.service == nil {
		return nil
	}
	var cTrack unsafe.Pointer = nil
	audioScenario := audioScenario_of_connection

	fmt.Printf("NewCustomAudioTrackPcm, audioScenario: %d, pcmSender.audioScenario: %d\n", audioScenario, pcmSender.audioScenario)
	if audioScenario == AudioScenarioAiServer && isSendExternalAudioForAI == false {
		cTrack  = C.agora_service_create_direct_custom_audio_track_pcm(agoraService.service, pcmSender.cSender)
	} else {
		cTrack = C.agora_service_create_custom_audio_track_pcm(agoraService.service, pcmSender.cSender)
	}
	
	if cTrack == nil {
		return nil
	}
	audioTrack := &LocalAudioTrack{
		cTrack: cTrack,
	}
	pcmSender.audioScenario = audioScenario

	fmt.Printf("NewCustomAudioTrackPcm, scenario: %d\n", audioScenario)

	// set send delay ms to 10ms, to avoid audio delay. NOTE: do not set it to 0, otherwise, it would set to default value: 260ms
	if audioTrack.cTrack != nil {
		audioTrack.SetSendDelayMs(10)
		audioTrack.SetEnabled(true)
	}
	// anywan, set the max buffered audio frame number to 100000
	
	audioTrack.SetMaxBufferedAudioFrameNumber(100000)
	
	return audioTrack
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
	if track == nil || track.cTrack == nil {
		return
	}
	cEnable := 0
	if enable {
		cEnable = 1
	}
	C.agora_local_audio_track_set_enabled(track.cTrack, C.int(cEnable))
}

func (track *LocalAudioTrack) AdjustPublishVolume(volume int) int {
	if track == nil || track.cTrack == nil {
		return -1
	}
	return int(C.agora_local_audio_track_adjust_publish_volume(track.cTrack, C.int(volume)))
}

// NOTICE: these interface below is temporary, may be removed in the future
// size is the number of 10ms audio frames
// the default value of this param is 30, ie. 300ms
func (track *LocalAudioTrack) SetMaxBufferedAudioFrameNumber(frameNum int) {
	if track == nil || track.cTrack == nil {
		return
	}
	C.agora_local_audio_track_set_max_bufferd_frame_number(track.cTrack, C.int(frameNum))
}

func (track *LocalAudioTrack) ClearSenderBuffer() int {
	if track == nil || track.cTrack == nil {
		return -1
	}
	return int(C.agora_local_audio_track_clear_sender_buffer(track.cTrack))
}

func (track *LocalAudioTrack) SetSendDelayMs(delayMs int) int {
	if track == nil || track.cTrack == nil {
		return -1
	}
	C.agora_local_audio_track_set_send_delay_ms(track.cTrack, C.int(delayMs))
	return 0
}
// to create direct audio track,no need to set send delay ms
// its will push all data to sd-rtn. and it will encode possible more data to sd-rtn.
func NewDirectCustomAudioTrackPcm(pcmSender *AudioPcmDataSender) *LocalAudioTrack {
	
	cTrack := C.agora_service_create_direct_custom_audio_track_pcm(agoraService.service, pcmSender.cSender)
	if cTrack == nil {
		return nil
	}

	return &LocalAudioTrack{
		cTrack: cTrack,
	}
}
