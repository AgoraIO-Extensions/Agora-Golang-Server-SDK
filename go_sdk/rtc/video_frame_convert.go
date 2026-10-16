package agoraservice

/*
#include "agora_media_base.h"
*/
import "C"
import "unsafe"

// Native VideoFrame uses byte strides for YUV and a total buffer length for
// packed RGBA/BGRA. Chroma planes round up when the height is odd.
func videoFrameBufferLengths(frame *C.video_frame) (y, u, v int) {
	height := int(frame.height)
	chromaHeight := (height + 1) / 2
	switch VideoPixelFormat(frame._type) {
	case VideoPixelBGRA, VideoPixelRGBA:
		return int(frame.y_stride), 0, 0
	case VideoPixelI422:
		return int(frame.y_stride) * height, int(frame.u_stride) * height, int(frame.v_stride) * height
	case VideoPixelNV12, VideoPixelNV21:
		return int(frame.y_stride) * height, int(frame.u_stride) * chromaHeight, 0
	case VideoPixelDefault, VideoPixelI420:
		return int(frame.y_stride) * height, int(frame.u_stride) * chromaHeight, int(frame.v_stride) * chromaHeight
	case VideoPixelI010:
		// I010 strides count 16-bit samples.
		return 2 * int(frame.y_stride) * height, 2 * int(frame.u_stride) * chromaHeight, 2 * int(frame.v_stride) * chromaHeight
	default:
		// Texture/CVPixelBuffer storage is native-owned, not a CPU byte buffer.
		return 0, 0, 0
	}
}

func copyVideoFrameBuffer(buffer *C.uint8_t, length int) []byte {
	if buffer == nil || length <= 0 {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(buffer), C.int(length))
}
