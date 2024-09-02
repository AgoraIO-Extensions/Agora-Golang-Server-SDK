package agoraservice

// #cgo CFLAGS: -I${SRCDIR}/../../agora_sdk/include/c/api2 -I${SRCDIR}/../../agora_sdk/include/c/base
// #include <string.h>
// #include "agora_media_node_factory.h"
import "C"
import "unsafe"

type VideoFrameSender struct {
	cSender unsafe.Pointer
}

func (mediaNodeFactory *MediaNodeFactory) NewVideoFrameSender() *VideoFrameSender {
	sender := C.agora_media_node_factory_create_video_frame_sender(mediaNodeFactory.cFactory)
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

func (sender *VideoFrameSender) SendVideoFrame(frame *VideoFrame) int {
	cData := C.CBytes(frame.Buffer)
	defer C.free(cData)
	cFrame := C.struct__external_video_frame{}
	C.memset(unsafe.Pointer(&cFrame), 0, C.sizeof_struct__external_video_frame)
	cFrame._type = 1
	cFrame.format = 1
	cFrame.buffer = cData
	cFrame.stride = C.int(frame.YStride)
	cFrame.height = C.int(frame.Height)
	cFrame.timestamp = C.longlong(frame.Timestamp)
	return int(C.agora_video_frame_sender_send(sender.cSender, &cFrame))
}
