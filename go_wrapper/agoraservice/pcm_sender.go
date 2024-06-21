package agoraservice

/*
#cgo CFLAGS: -I../../agora_sdk/include_c/api2 -I../../agora_sdk/include_c/base
#cgo LDFLAGS: -L../../agora_sdk/ -lagora_rtc_sdk -lagora-fdkaac -lagora-ffmpeg

#include "agora_local_user.h"
#include "agora_rtc_conn.h"
#include "agora_service.h"
#include "agora_media_base.h"
*/
import "C"
import "unsafe"

// type PcmSenderConfig struct {
// 	SampleRate     int
// 	Channels       int
// 	BytesPerSample int
// }

type PcmSender struct {
	// config       *PcmSenderConfig
	cLocalUser   unsafe.Pointer
	cAudioTrack  unsafe.Pointer
	cAudioSender unsafe.Pointer
}

func (con *RtcConnection) NewPcmSender() *PcmSender {
	ret := &PcmSender{
		cLocalUser:   con.cLocalUser,
		cAudioSender: C.agora_media_node_factory_create_audio_pcm_data_sender(agoraService.mediaFactory),
	}
	ret.cAudioTrack = C.agora_service_create_custom_audio_track_pcm(agoraService.service, ret.cAudioSender)
	return ret
}

func (sender *PcmSender) Release() {
	if sender.cAudioSender == nil {
		return
	}
	C.agora_local_audio_track_destroy(sender.cAudioTrack)
	C.agora_audio_pcm_data_sender_destroy(sender.cAudioSender)
}

func (sender *PcmSender) Start() int {
	C.agora_local_audio_track_set_enabled(sender.cAudioTrack, C.int(1))
	return int(C.agora_local_user_publish_audio(sender.cLocalUser, sender.cAudioTrack))
}

func (sender *PcmSender) Stop() int {
	ret := int(C.agora_local_user_unpublish_audio(sender.cLocalUser, sender.cAudioTrack))
	C.agora_local_audio_track_set_enabled(sender.cAudioTrack, C.int(0))
	return ret
}

func (sender *PcmSender) SendPcmData(frame *PcmAudioFrame) int {
	cData := C.CBytes(frame.Data)
	defer C.free(cData)
	return int(C.agora_audio_pcm_data_sender_send(sender.cAudioSender, cData,
		C.uint(frame.Timestamp), C.uint(frame.SamplesPerChannel),
		C.uint(frame.BytesPerSample), C.uint(frame.NumberOfChannels),
		C.uint(frame.SampleRate)))
}

func (sender *PcmSender) AdjustVolume(volume int) int {
	return int(C.agora_local_audio_track_adjust_publish_volume(sender.cAudioTrack, C.int(volume)))
}

// NOTICE: these interface below is temporary, may be removed in the future
// size is the number of 10ms audio frames
// the default value of this param is 30, ie. 300ms
// func (sender *PcmSender) SetSendBufferSize(size int) int {
// 	return int(C.agora_local_audio_track_set_max_buffer_audio_frame_number(sender.cAudioTrack, C.int(size)))
// }

// func (sender *PcmSender) ClearSendBuffer() int {
// 	return int(C.agora_local_audio_track_clear_buffer(sender.cAudioTrack))
// }
