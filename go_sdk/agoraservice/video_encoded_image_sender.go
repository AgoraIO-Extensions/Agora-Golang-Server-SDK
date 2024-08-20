package agoraservice

// #cgo CFLAGS: -I${SRCDIR}/../../agora_sdk/include_c/api2 -I${SRCDIR}/../../agora_sdk/include_c/base
// #include "agora_media_node_factory.h"
import "C"
import "unsafe"

type VideoEncodedImageSender struct {
	cSender unsafe.Pointer
}

func (mediaNodeFactory *MediaNodeFactory) NewVideoEncodedImageSender() *VideoEncodedImageSender {
	sender := C.agora_media_node_factory_create_video_encoded_image_sender(mediaNodeFactory.cFactory)
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
