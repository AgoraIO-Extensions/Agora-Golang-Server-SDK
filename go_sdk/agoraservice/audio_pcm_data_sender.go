package agoraservice

// #cgo CFLAGS: -I${SRCDIR}/../../agora_sdk/include/c/api2 -I${SRCDIR}/../../agora_sdk/include/c/base
// #include "agora_media_node_factory.h"
import "C"
import "unsafe"

type AudioPcmDataSender struct {
	cSender unsafe.Pointer
}

func (mediaNodeFactory *MediaNodeFactory) NewAudioPcmDataSender() *AudioPcmDataSender {
	sender := C.agora_media_node_factory_create_audio_pcm_data_sender(mediaNodeFactory.cFactory)
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

func (sender *AudioPcmDataSender) SendAudioPcmData(frame *PcmAudioFrame) int {
	cData := C.CBytes(frame.Data)
	defer C.free(cData)
	return int(C.agora_audio_pcm_data_sender_send(sender.cSender, cData,
		C.uint(frame.Timestamp), C.uint(frame.SamplesPerChannel),
		C.uint(frame.BytesPerSample), C.uint(frame.NumberOfChannels),
		C.uint(frame.SampleRate)))
}
