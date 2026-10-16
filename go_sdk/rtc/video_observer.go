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
	"unsafe"
)

func videoObserverForHandle(handle unsafe.Pointer) *VideoFrameObserver {
	conn := agoraService.getConFromHandle(handle, ConTypeCVideoObserver)
	if conn == nil {
		return nil
	}
	return conn.videoObserver.Load()
}

func dispatchVideoFrame(frame *C.video_frame, callback func(*VideoFrame) bool) C.int {
	if frame == nil {
		return 0
	}
	if callback == nil {
		return 1
	}
	if !callback(GoVideoFrame(frame)) {
		return 0
	}
	return 1
}

//export goOnCaptureVideoFrame
func goOnCaptureVideoFrame(handle unsafe.Pointer, sourceType C.int, frame *C.video_frame) C.int {
	observer := videoObserverForHandle(handle)
	if observer == nil {
		return 0
	}
	var callback func(*VideoFrame) bool
	if observer.OnCaptureVideoFrame != nil {
		callback = func(frame *VideoFrame) bool { return observer.OnCaptureVideoFrame(int(sourceType), frame) }
	}
	return dispatchVideoFrame(frame, callback)
}

//export goOnPreEncodeVideoFrame
func goOnPreEncodeVideoFrame(handle unsafe.Pointer, sourceType C.int, frame *C.video_frame) C.int {
	observer := videoObserverForHandle(handle)
	if observer == nil {
		return 0
	}
	var callback func(*VideoFrame) bool
	if observer.OnPreEncodeVideoFrame != nil {
		callback = func(frame *VideoFrame) bool { return observer.OnPreEncodeVideoFrame(int(sourceType), frame) }
	}
	return dispatchVideoFrame(frame, callback)
}

//export goOnMediaPlayerVideoFrame
func goOnMediaPlayerVideoFrame(handle unsafe.Pointer, frame *C.video_frame, mediaPlayerId C.int) C.int {
	observer := videoObserverForHandle(handle)
	if observer == nil {
		return 0
	}
	var callback func(*VideoFrame) bool
	if observer.OnMediaPlayerVideoFrame != nil {
		callback = func(frame *VideoFrame) bool { return observer.OnMediaPlayerVideoFrame(frame, int(mediaPlayerId)) }
	}
	return dispatchVideoFrame(frame, callback)
}

//export goOnRenderVideoFrame
func goOnRenderVideoFrame(handle unsafe.Pointer, channelId *C.char, uid C.uid_t, frame *C.video_frame) C.int {
	observer := videoObserverForHandle(handle)
	if observer == nil {
		return 0
	}
	var callback func(*VideoFrame) bool
	if observer.OnRenderVideoFrame != nil {
		callback = func(frame *VideoFrame) bool {
			return observer.OnRenderVideoFrame(C.GoString(channelId), uint32(uid), frame)
		}
	}
	return dispatchVideoFrame(frame, callback)
}

//export goOnTranscodedVideoFrame
func goOnTranscodedVideoFrame(handle unsafe.Pointer, frame *C.video_frame) C.int {
	observer := videoObserverForHandle(handle)
	if observer == nil {
		return 0
	}
	return dispatchVideoFrame(frame, observer.OnTranscodedVideoFrame)
}

//export goGetVideoFrameProcessMode
func goGetVideoFrameProcessMode(handle unsafe.Pointer) C.int {
	observer := videoObserverForHandle(handle)
	if observer == nil || observer.GetVideoFrameProcessMode == nil {
		return C.int(VideoFrameProcessModeReadOnly)
	}
	return C.int(observer.GetVideoFrameProcessMode())
}

//export goGetVideoFormatPreference
func goGetVideoFormatPreference(handle unsafe.Pointer) C.int {
	observer := videoObserverForHandle(handle)
	if observer == nil || observer.GetVideoFormatPreference == nil {
		return C.int(VideoPixelI420)
	}
	return C.int(observer.GetVideoFormatPreference())
}

//export goGetRotationApplied
func goGetRotationApplied(handle unsafe.Pointer) C.int {
	observer := videoObserverForHandle(handle)
	if observer != nil && observer.GetRotationApplied != nil && observer.GetRotationApplied() {
		return 1
	}
	return 0
}

//export goGetMirrorApplied
func goGetMirrorApplied(handle unsafe.Pointer) C.int {
	observer := videoObserverForHandle(handle)
	if observer != nil && observer.GetMirrorApplied != nil && observer.GetMirrorApplied() {
		return 1
	}
	return 0
}

//export goGetObservedFramePosition
func goGetObservedFramePosition(handle unsafe.Pointer) C.uint32_t {
	observer := videoObserverForHandle(handle)
	if observer == nil || observer.GetObservedFramePosition == nil {
		return C.uint32_t(VideoPositionPostCapturer | VideoPositionPreRenderer)
	}
	return C.uint32_t(observer.GetObservedFramePosition())
}

//export goOnEncodedVideoFrame
func goOnEncodedVideoFrame(observer unsafe.Pointer, channelId *C.char, uid *C.char, imageBuffer *C.uint8_t, length C.size_t,
	video_encoded_frame_info *C.struct__encoded_video_frame_info) C.int {
	// validity check
	if observer == nil {
		return C.int(0)
	}
	con := agoraService.getConFromHandle(observer, ConTypeCEncodedVideoObserver)
	if con == nil || con.encodedVideoObserver == nil || con.encodedVideoObserver.OnEncodedVideoFrame == nil {
		return C.int(0)
	}
	goUid := C.GoString(uid)
	goChannelId := C.GoString(channelId)
	//fmt.Printf("goOnEncodedVideoFrame, uid: %s, length: %d, uid: %d\n", goUid, length, uid)
	//fmt.Printf("goOnEncodedVideoFrame, video_encoded_frame_info: %v\n", video_encoded_frame_info)
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
	if con.encodedVideoObserver.OnEncodedVideoFrame(goChannelId, goUid, goImageBuffer, goFrameInfo) {
		return C.int(1)
	}
	return C.int(0)
}
