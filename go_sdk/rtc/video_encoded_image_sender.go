package agoraservice

// #cgo CFLAGS: -I${SRCDIR}/../../agora_sdk/include/c/api2 -I${SRCDIR}/../../agora_sdk/include/c/base
// #include "agora_media_node_factory.h"
import "C"
import "unsafe"

type EncodedVideoFrameInfo struct {
	/**
	 * The video codec: #VideoCodecTypeXxxx.
	 */
	CodecType VideoCodecType
	/**
	 * The width (px) of the video.
	 */
	Width int
	/**
	 * The height (px) of the video.
	 */
	Height int
	/**
	 * The number of video frames per second.
	 * This value will be used for calculating timestamps of the encoded image.
	 * If framesPerSecond equals zero, then real timestamp will be used.
	 * Otherwise, timestamp will be adjusted to the value of framesPerSecond set.
	 */
	FramesPerSecond int
	/**
	 * The frame type of the encoded video frame: #VIDEO_FRAME_TYPE.
	 */
	FrameType VideoFrameType
	/**
	 * The rotation information of the encoded video frame: #VIDEO_ORIENTATION.
	 */
	Rotation VideoOrientation
	/**
	 * The track ID of the video frame.
	 */
	TrackId int // This can be reserved for multiple video tracks, we need to create different ssrc
	// and additional payload for later implementation.
	/**
	 * This is a input parameter which means the timestamp for capturing the video.
	 */
	CaptureTimeMs int64
	/**
	 * The timestamp for decoding the video.
	 */
	DecodeTimeMs int64
	/**
	 * ID of the user.
	 */
	Uid uint32
	/**
	 * The stream type of video frame.
	 */
	StreamType int
	PresentTimeMs int64
}

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

func (sender *VideoEncodedImageSender) SendEncodedVideoImage(payload []byte, frameInfo *EncodedVideoFrameInfo) int {
	cData, pinner := unsafeCBytes(payload)
	defer pinner.Unpin()
	cFrameInfo := &C.struct__encoded_video_frame_info{
		codec_type:        C.int(frameInfo.CodecType),
		width:             C.int(frameInfo.Width),
		height:            C.int(frameInfo.Height),
		frames_per_second: C.int(frameInfo.FramesPerSecond),
		frame_type:        C.int(frameInfo.FrameType),
		rotation:          C.int(frameInfo.Rotation),
		track_id:          C.int(frameInfo.TrackId),
		capture_time_ms:   C.int64_t(frameInfo.CaptureTimeMs),
		decode_time_ms:    C.int64_t(frameInfo.DecodeTimeMs),
		uid:               C.uint(frameInfo.Uid),
		stream_type:       C.int(frameInfo.StreamType),
		presentation_ms:   C.int64_t(frameInfo.PresentTimeMs),
	}
	return int(C.agora_video_encoded_image_sender_send(sender.cSender, (*C.uint8_t)(cData), C.uint32_t(len(payload)), cFrameInfo))
}
