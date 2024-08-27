package agoraservice

// #cgo CFLAGS: -I../../agora_sdk/include_c/api2 -I../../agora_sdk/include_c/base
// #include "agora_media_node_factory.h"
// #include "video_observer_cgo.h"
import "C"
import "unsafe"

type VideoEncodedImageReceiver struct {
	OnEncodedVideoFrame func(receiver *VideoEncodedImageReceiver, uid string, imageBuffer []byte,
		frameInfo *EncodedVideoFrameInfo) bool
}

type videoEncodedImageReceiverInner struct {
	cReceiver unsafe.Pointer
	receiver  *VideoEncodedImageReceiver
}

func (conn *RtcConnection) newVideoEncodedImageReceiverInner(receiver *VideoEncodedImageReceiver) *videoEncodedImageReceiverInner {
	observer := C.struct__video_encoded_frame_observer{
		on_encoded_video_frame: (*[0]byte)(cgo_on_encoded_video_frame),
	}
	cReceiver := C.agora_video_encoded_image_receiver_create(&observer)
	if cReceiver == nil {
		return nil
	}
	receiveInner := &videoEncodedImageReceiverInner{
		cReceiver: cReceiver,
		receiver:  receiver,
	}
	conn.remoteVideoRWMutex.Lock()
	conn.remoteEncodedVideoReceivers[cReceiver] = receiveInner
	conn.remoteVideoRWMutex.Unlock()
	return receiveInner
}

func (receiver *videoEncodedImageReceiverInner) release() {
	if receiver.cReceiver == nil {
		return
	}
	C.agora_video_encoded_image_receiver_destroy(receiver.cReceiver)
	receiver.cReceiver = nil
}
