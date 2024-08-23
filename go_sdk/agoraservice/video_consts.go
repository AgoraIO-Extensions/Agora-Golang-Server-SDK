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
