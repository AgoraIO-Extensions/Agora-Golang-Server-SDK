package agoraservice

// #cgo CFLAGS: -I${SRCDIR}/../../agora_sdk/include_c/api2 -I${SRCDIR}/../../agora_sdk/include_c/base
// #include "agora_media_node_factory.h"
import "C"
import "unsafe"

type AudioEncodedFrameSender struct {
	cSender unsafe.Pointer
}

func (mediaNodeFactory *MediaNodeFactory) NewAudioEncodedFrameSender() *AudioEncodedFrameSender {
	sender := C.agora_media_node_factory_create_audio_encoded_frame_sender(mediaNodeFactory.cFactory)
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
