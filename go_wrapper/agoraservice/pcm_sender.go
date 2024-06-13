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

type PcmAudioFrame struct {
	Data              []byte
	CaptureTimestamp  int64
	SamplesPerChannel int
	BytesPerSample    int
	NumberOfChannels  int
	SampleRate        int
}

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
	return int(C.agora_local_user_publish_audio(sender.cLocalUser, sender.cAudioTrack))
}

func (sender *PcmSender) Stop() int {
	return int(C.agora_local_user_unpublish_audio(sender.cLocalUser, sender.cAudioTrack))
}

func (sender *PcmSender) SendPcmData(frame *PcmAudioFrame) int {
	cData := C.CBytes(frame.Data)
	defer C.free(cData)
	return int(C.agora_audio_pcm_data_sender_send(sender.cAudioSender, cData,
		C.uint(frame.CaptureTimestamp), C.uint(frame.SamplesPerChannel),
		C.uint(frame.BytesPerSample), C.uint(frame.NumberOfChannels),
		C.uint(frame.SampleRate)))
}
