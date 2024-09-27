package agoraservice

type VideoCodecType int

const (
	VideoCodecTypeNone        VideoCodecType = 0
	VideoCodecTypeVp8         VideoCodecType = 1
	VideoCodecTypeH264        VideoCodecType = 2
	VideoCodecTypeH265        VideoCodecType = 3
	VideoCodecTypeGeneric     VideoCodecType = 6
	VideoCodecTypeGenericH264 VideoCodecType = 7
	VideoCodecTypeAv1         VideoCodecType = 12
	VideoCodecTypeVp9         VideoCodecType = 13
	VideoCodecTypeGenericJpeg VideoCodecType = 20
)

type VideoSendCcState int

const (
	VideoSendCcEnabled  VideoSendCcState = 0
	VideoSendCcDisabled VideoSendCcState = 1
)

type VideoStreamType int

const (
	/**
	 * 0: The high-quality video stream, which has a higher resolution and bitrate.
	 */
	VideoStreamHigh VideoStreamType = 0
	/**
	 * 1: The low-quality video stream, which has a lower resolution and bitrate.
	 */
	VideoStreamLow VideoStreamType = 1
)

/**
 * Types of the video frame.
 */
type VideoFrameType int

const (
	/**
	 * (Default) Blank frame
	 */
	VideoFrameTypeBlankFrame VideoFrameType = 0
	/**
	 * (Default) Key frame
	 */
	VideoFrameTypeKeyFrame VideoFrameType = 3
	/**
	 * (Default) Delta frame
	 */
	VideoFrameTypeDeltaFrame VideoFrameType = 4
	/**
	 * (Default) B frame
	 */
	VideoFrameTypeBFrame VideoFrameType = 5
	/**
	 * (Default) Droppable frame
	 */
	VideoFrameTypeDroppableFrame VideoFrameType = 6
	/**
	 * (Default) Unknown frame type
	 */
	VideoFrameTypeUnknown VideoFrameType = -1
)

// OrientationMode represents video output orientation modes.
type OrientationMode int

const (
	// OrientationModeAdaptive indicates that the output video always follows the orientation of the captured video.
	// The receiver takes the rotational information passed on from the video encoder. This mode applies to scenarios
	// where video orientation can be adjusted on the receiver:
	// - If the captured video is in landscape mode, the output video is in landscape mode.
	// - If the captured video is in portrait mode, the output video is in portrait mode.
	OrientationModeAdaptive OrientationMode = 0

	// OrientationModeFixedLandscape indicates landscape mode. In this mode, the SDK always outputs videos in landscape (horizontal) mode.
	// If the captured video is in portrait mode, the video encoder crops it to fit the output. Applies to situations where the receiving end
	// cannot process the rotational information. For example, CDN live streaming.
	OrientationModeFixedLandscape OrientationMode = 1

	// OrientationModeFixedPortrait indicates portrait mode. In this mode, the SDK always outputs video in portrait (portrait) mode.
	// If the captured video is in landscape mode, the video encoder crops it to fit the output. Applies to situations where the receiving end
	// cannot process the rotational information. For example, CDN live streaming.
	OrientationModeFixedPortrait OrientationMode = 2
)

// DegradationPreference represents video degradation preferences under limited bandwidth.
type DegradationPreference int

const (
	// DegradeMaintainQuality prefers to reduce the video frame rate while maintaining video quality during video
	// encoding under limited bandwidth. This degradation preference is suitable for scenarios where
	// video quality is prioritized.
	//
	// Note: In the COMMUNICATION channel profile, the resolution of the video sent may change, so
	// remote users need to handle this issue.
	DegradeMaintainQuality DegradationPreference = 0

	// DegradeMaintainFramerate prefers to reduce the video quality while maintaining the video frame rate during video
	// encoding under limited bandwidth. This degradation preference is suitable for scenarios where
	// smoothness is prioritized and video quality is allowed to be reduced.
	DegradeMaintainFramerate DegradationPreference = 1

	// DegradeMaintainBalanced reduces the video frame rate and video quality simultaneously during video encoding under
	// limited bandwidth. MaintainBalanced has a lower reduction than MaintainQuality and MaintainFramerate,
	// and this preference is suitable for scenarios where both smoothness and video quality are a
	// priority.
	DegradeMaintainBalanced DegradationPreference = 2

	// DegradeMaintainResolution degrades framerate in order to maintain resolution.
	DegradeMaintainResolution DegradationPreference = 3

	// DegradeDisabled disables VQC adjustment.
	DegradeDisabled DegradationPreference = 100
)

// VideoMirrorModeType represents the video mirror mode types.
type VideoMirrorModeType int

const (
	// VideoMirrorModeAuto represents the mirror mode determined by the SDK.
	VideoMirrorModeAuto VideoMirrorModeType = 0
	// VideoMirrorModeEnabled enables the mirror mode.
	VideoMirrorModeEnabled VideoMirrorModeType = 1
	// VideoMirrorModeDisabled disables the mirror mode.
	VideoMirrorModeDisabled VideoMirrorModeType = 2
)

/**
 * The rotation information.
 */
type VideoOrientation int

const (
	/**
	 * 0: Rotate the video by 0 degree clockwise.
	 */
	VideoOrientation0 VideoOrientation = 0
	/**
	 * 90: Rotate the video by 90 degrees clockwise.
	 */
	VideoOrientation90 VideoOrientation = 90
	/**
	 * 180: Rotate the video by 180 degrees clockwise.
	 */
	VideoOrientation180 VideoOrientation = 180
	/**
	 * 270: Rotate the video by 270 degrees clockwise.
	 */
	VideoOrientation270 VideoOrientation = 270
)

/**
 * Video buffer types.
 */
type VideoBufferType int

const (
	/**
	 * 1: Raw data.
	 */
	VideoBufferRawData VideoBufferType = 1
	/**
	 * 2: The same as VIDEO_BUFFER_RAW_DATA.
	 */
	VideoBufferArray VideoBufferType = 2
	/**
	 * 3: The video buffer in the format of texture.
	 */
	VideoBufferTexture VideoBufferType = 3
)

/**
 * Video pixel formats.
 */
type VideoPixelFormat int

const (
	/**
	 * 0: Default format.
	 */
	VideoPixelDefault VideoPixelFormat = 0
	/**
	 * 1: I420.
	 */
	VideoPixelI420 VideoPixelFormat = 1
	/**
	 * 2: BGRA.
	 */
	VideoPixelBGRA VideoPixelFormat = 2
	/**
	 * 3: NV21.
	 */
	VideoPixelNV21 VideoPixelFormat = 3
	/**
	 * 4: RGBA.
	 */
	VideoPixelRGBA VideoPixelFormat = 4
	/**
	 * 8: NV12.
	 */
	VideoPixelNV12 VideoPixelFormat = 8
	/**
	 * 10: GL_TEXTURE_2D
	 */
	VideoTexture2D VideoPixelFormat = 10
	/**
	 * 11: GL_TEXTURE_OES
	 */
	VideoTextureOES VideoPixelFormat = 11
	/*
		12: pixel format for iOS CVPixelBuffer NV12
	*/
	VideoCVPixelNV12 VideoPixelFormat = 12
	/*
		13: pixel format for iOS CVPixelBuffer I420
	*/
	VideoCVPixelI420 VideoPixelFormat = 13
	/*
		14: pixel format for iOS CVPixelBuffer BGRA
	*/
	VideoCVPixelBGRA VideoPixelFormat = 14
	/**
	15: pixel format for iOS CVPixelBuffer P010(10bit NV12)
	*/
	VideoCVPixelP010 VideoPixelFormat = 15
	/**
	 * 16: I422.
	 */
	VideoPixelI422 VideoPixelFormat = 16
	/**
	 * 17: ID3D11Texture2D, only support DXGI_FORMAT_B8G8R8A8_UNORM, DXGI_FORMAT_B8G8R8A8_TYPELESS, DXGI_FORMAT_NV12 texture format
	 */
	VideoTextureID3D11Texture2D VideoPixelFormat = 17
	/**
	 * 18: I010. 10bit I420 data.
	 * @technical preview
	 */
	VideoPixelI010 VideoPixelFormat = 18
)
