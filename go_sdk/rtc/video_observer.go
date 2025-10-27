package agoraservice

/*
#cgo CFLAGS: -I../../agora_sdk/include/c/api2 -I../../agora_sdk/include/c/base

#include "agora_local_user.h"
#include "agora_rtc_conn.h"
#include "agora_service.h"
#include "agora_media_base.h"
*/
import "C"
import (
	"strconv"
	"unsafe"
)

//export goOnVideoFrame
func goOnVideoFrame(cObserver unsafe.Pointer, channelId *C.char, uid *C.char, frame *C.struct__video_frame) C.int {
	// validity check
	if cObserver == nil {
		return C.int(0)
	}
	// added by wei to avoid uid=0 i.e local users callback , which is a bug in current 44.3.1 sdk version
	// but fix it from cgo layer is not a good idea
	goUid := C.GoString(uid)
	if goUid == "0" {
		return C.int(0)
	}
	
	con := agoraService.getConFromHandle(cObserver, ConTypeCVideoObserver)
	if con == nil || con.videoObserver == nil || con.videoObserver.OnFrame == nil {
		return C.int(0)
	}
	goChannelId := C.GoString(channelId)
	
	goFrame := GoVideoFrame(frame)
	ret := con.videoObserver.OnFrame(goChannelId, goUid, goFrame)
	if ret {
		return C.int(1)
	}
	return C.int(0)
}

//export goOnEncodedVideoFrame
func goOnEncodedVideoFrame(observer unsafe.Pointer, uid C.uint32_t, imageBuffer *C.uint8_t, length C.size_t,
	video_encoded_frame_info *C.struct__encoded_video_frame_info) C.int {
	// validity check
	if observer == nil {
		return C.int(0)
	}
	con := agoraService.getConFromHandle(observer, ConTypeCEncodedVideoObserver)
	if con == nil || con.encodedVideoObserver == nil || con.encodedVideoObserver.OnEncodedVideoFrame == nil {
		return C.int(0)
	}
	goUid := strconv.FormatUint(uint64(uid), 10)
	goImageBuffer := C.GoBytes(unsafe.Pointer(imageBuffer), C.int(length))
	// GoEncodedVideoFrameInfo(video_encoded_frame_info)
	goFrameInfo := &EncodedVideoFrameInfo{
		CodecType:       VideoCodecType(video_encoded_frame_info.codec_type),
		Width:           int(video_encoded_frame_info.width),
		Height:          int(video_encoded_frame_info.height),
		FramesPerSecond: int(video_encoded_frame_info.frames_per_second),
		FrameType:       VideoFrameType(video_encoded_frame_info.frame_type),
		Rotation:        VideoOrientation(video_encoded_frame_info.rotation),
		TrackId:         int(video_encoded_frame_info.track_id),
		CaptureTimeMs:   int64(video_encoded_frame_info.capture_time_ms),
		DecodeTimeMs:    int64(video_encoded_frame_info.decode_time_ms),
		Uid:             uint32(video_encoded_frame_info.uid),
		StreamType:      int(video_encoded_frame_info.stream_type),
		PresentTimeMs:   int64(video_encoded_frame_info.presentation_ms),
	}
	if con.encodedVideoObserver.OnEncodedVideoFrame(goUid, goImageBuffer, goFrameInfo) {
		return C.int(1)
	}
	return C.int(0)
}
