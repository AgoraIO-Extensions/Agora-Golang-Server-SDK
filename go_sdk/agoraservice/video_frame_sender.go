package agoraservice

// #cgo CFLAGS: -I${SRCDIR}/../../agora_sdk/include/c/api2 -I${SRCDIR}/../../agora_sdk/include/c/base
// #include <string.h>
// #include <stdint.h>
// #include "agora_media_node_factory.h"
import "C"
import "unsafe"


  type ColorSpaceType struct {
	// The indices are equal to the values specified in T-REC H.273 Table 2.
	PrimariesId int
	// The indices are equal to the values specified in T-REC H.273 Table 3.
	TransferId  int
	// The indices are equal to the values specified in T-REC H.273 Table 4.
	MatrixId    int
	// The indices are equal to the values specified at
	// https://www.webmproject.org/docs/container/#colour for the element Range.
	RangeId     int
}

// ExternalVideoFrame represents an external video frame.
type ExternalVideoFrame struct {
	// Type is the buffer type.
	Type VideoBufferType
	// Format is the pixel format.
	Format VideoPixelFormat
	// Buffer is the video buffer.
	Buffer []byte
	// Stride is the line spacing of the incoming video frame, which must be in pixels instead of bytes.
	// For textures, it is the width of the texture.
	Stride int
	// Height is the height of the incoming video frame.
	Height int
	// CropLeft is the number of pixels trimmed from the left. The default value is 0.
	CropLeft int
	// CropTop is the number of pixels trimmed from the top. The default value is 0.
	CropTop int
	// CropRight is the number of pixels trimmed from the right. The default value is 0.
	CropRight int
	// CropBottom is the number of pixels trimmed from the bottom. The default value is 0.
	CropBottom int
	// Rotation is the clockwise rotation of the video frame. You can set the rotation angle as 0, 90, 180, or 270.
	// The default value is 0.
	Rotation VideoOrientation
	// Timestamp is the timestamp of the incoming video frame (ms). An incorrect timestamp results in frame loss or
	// unsynchronized audio and video.
	Timestamp int64
	// EGLContext is the EGL context used by the video frame.
	EGLContext unsafe.Pointer
	// EGLType is the texture ID used by the video frame.
	EGLType int
	// TextureID is the incoming 4x4 transformational matrix. The typical value is a unit matrix.
	TextureID int
	// Matrix is the incoming 4x4 transformational matrix. The typical value is a unit matrix.
	Matrix [16]float32
	// MetadataBuffer is the metadata buffer.
	MetadataBuffer []byte
	// AlphaBuffer is the alpha buffer.
	AlphaBuffer []byte
	// [For bgra or rgba only] Extract alphaBuffer from bgra or rgba data. Set it true if you do not explicitly specify the alphabuffer.
	// The default value is false
	FillAlphaBuffer bool

	// The relative position between alphabuffer and the frame.
	// 0: Normal frame;
	// 1: Alphabuffer is above the frame;
	// 2: Alphabuffer is below the frame;
	// 3: Alphabuffer is on the left of frame;
	// 4: Alphabuffer is on the right of frame;
	// The default value is 0.
	AlphaMode int
	/*
  *  The color_space_type
  */
  ColorSpace ColorSpaceType;
}

// VideoFrame represents a video frame.
type VideoFrame struct {
	Type           VideoBufferType  // Type is the buffer type.
	Width          int              // Width is the video pixel width.
	Height         int              // Height is the video pixel height.
	YStride        int              // YStride is the line span of Y buffer in YUV data.
	UStride        int              // UStride is the line span of U buffer in YUV data.
	VStride        int              // VStride is the line span of V buffer in YUV data.
	YBuffer        []byte           // YBuffer is the pointer to the Y buffer pointer in the YUV data.
	UBuffer        []byte           // UBuffer is the pointer to the U buffer pointer in the YUV data.
	VBuffer        []byte           // VBuffer is the pointer to the V buffer pointer in the YUV data.
	Rotation       VideoOrientation // Rotation is the rotation of this frame before rendering the video.
	RenderTimeMs   int64            // RenderTimeMs is the timestamp to render the video stream.
	AVSyncType     int              // AVSyncType is the AV sync type.
	MetadataBuffer []byte           // MetadataBuffer is the metadata buffer.
	SharedContext  unsafe.Pointer   // SharedContext is the EGL context.
	TextureID      int              // TextureID is the texture ID used by the video frame.
	Matrix         [16]float32      // Matrix is the incoming 4x4 transformational matrix.
	AlphaBuffer    []byte           // AlphaBuffer is the alpha buffer.
	// The relative position between alphabuffer and the frame.
	// 0: Normal frame;
	// 1: Alphabuffer is above the frame;
	// 2: Alphabuffer is below the frame;
	// 3: Alphabuffer is on the left of frame;
	// 4: Alphabuffer is on the right of frame;
	// The default value is 0.
	AlphaMode int
}

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

func (sender *VideoFrameSender) SendVideoFrame(frame *ExternalVideoFrame) int {
	cData, pinner := unsafeCBytes(frame.Buffer)
	defer pinner.Unpin()
	cFrame := C.struct__external_video_frame{}
	C.memset(unsafe.Pointer(&cFrame), 0, C.sizeof_struct__external_video_frame)
	cFrame._type = C.int(frame.Type)
	cFrame.format = C.int(frame.Format)
	cFrame.buffer = cData
	cFrame.stride = C.int(frame.Stride)
	cFrame.height = C.int(frame.Height)
	cFrame.crop_left = C.int(frame.CropLeft)
	cFrame.crop_top = C.int(frame.CropTop)
	cFrame.crop_right = C.int(frame.CropRight)
	cFrame.crop_bottom = C.int(frame.CropBottom)
	cFrame.rotation = C.int(frame.Rotation)
	cFrame.timestamp = C.longlong(frame.Timestamp)
	if frame.MetadataBuffer != nil {
		metadata, metadataPinner := unsafeCBytes(frame.MetadataBuffer)
		defer metadataPinner.Unpin()
		cFrame.metadata_buffer = (*C.uint8_t)(metadata)
		cFrame.metadata_size = C.int(len(frame.MetadataBuffer))
	}
	if frame.AlphaBuffer != nil {
		alpha, alphaPinner := unsafeCBytes(frame.AlphaBuffer)
		defer alphaPinner.Unpin()
		cFrame.alpha_buffer = (*C.uint8_t)(alpha)
	}
	if frame.FillAlphaBuffer {
		cFrame.fill_alpha_buffer = C.uint8_t(1)
	}
	cFrame.alpha_mode = C.int(frame.AlphaMode)
	// fill color space: for go,the default value of struct 's member is zero value!
	cFrame.color_space.matrix_id = C.int(frame.ColorSpace.MatrixId)
	cFrame.color_space.primaries_id = C.int(frame.ColorSpace.PrimariesId)
	cFrame.color_space.transfer_id = C.int(frame.ColorSpace.TransferId)
	cFrame.color_space.range_id = C.int(frame.ColorSpace.RangeId)
	return int(C.agora_video_frame_sender_send(sender.cSender, &cFrame))
}
