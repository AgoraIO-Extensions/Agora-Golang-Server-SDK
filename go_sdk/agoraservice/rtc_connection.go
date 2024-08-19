package agoraservice

/*
#cgo CFLAGS: -I../../agora_sdk/include_c/api2 -I../../agora_sdk/include_c/base

#include "agora_local_user.h"
#include "agora_rtc_conn.h"
#include "agora_service.h"
#include "agora_media_base.h"
#include "agora_parameter.h"
*/
import "C"
import "unsafe"

const (
	/**
	* 0: The user has muted the audio.
	 */
	UserMediaInfoMuteAudio = 0
	/**
	* 1: The user has muted the video.
	 */
	UserMediaInfoMuteVideo = 1
	/**
	* 4: The user has enabled the video, which includes video capturing and encoding.
	 */
	UserMediaInfoEnableVideo = 4
	/**
	* 8: The user has enabled the local video capturing.
	 */
	UserMediaInfoEnableLocalVideo = 8
)

const (
	/**
	 * 0: The default audio profile.
	 * - In the Communication profile, it represents a sample rate of 16 kHz, music encoding, mono, and a bitrate
	 * of up to 16 Kbps.
	 * - In the Live-broadcast profile, it represents a sample rate of 48 kHz, music encoding, mono, and a bitrate
	 * of up to 64 Kbps.
	 */
	AudioProfileDefault = 0
	/**
	 * 1: A sample rate of 16 kHz, audio encoding, mono, and a bitrate up to 18 Kbps.
	 */
	AudioProfileSpeechStandard = 1
	/**
	 * 2: A sample rate of 48 kHz, music encoding, mono, and a bitrate of up to 64 Kbps.
	 */
	AudioProfileMusicStandard = 2
	/**
	 * 3: A sample rate of 48 kHz, music encoding, stereo, and a bitrate of up to 80
	 * Kbps.
	 */
	AudioProfileMusicStandardStereo = 3
	/**
	 * 4: A sample rate of 48 kHz, music encoding, mono, and a bitrate of up to 96 Kbps.
	 */
	AudioProfileMusicHighQuality = 4
	/**
	 * 5: A sample rate of 48 kHz, music encoding, stereo, and a bitrate of up to 128 Kbps.
	 */
	AudioProfileMusicHighQualityStereo = 5
	/**
	 * 6: A sample rate of 16 kHz, audio encoding, mono, and a bitrate of up to 64 Kbps.
	 */
	AudioProfileIot = 6
)

type RtcConnectionInfo struct {
	ConnectionId uint
	/**
	 * ID of the target channel. NULL if you did not call the connect
	 * method.
	 */
	ChannelId string
	/**
	 * The state of the current connection: #CONNECTION_STATE_TYPE.
	 */
	State int
	/**
	 * ID of the local user.
	 */
	LocalUserId string
	/**
	 * Internal use only.
	 */
	InternalUid uint
}

type PcmAudioFrame struct {
	Data              []byte
	Timestamp         int64
	SamplesPerChannel int
	BytesPerSample    int
	NumberOfChannels  int
	SampleRate        int
}

// support YUV I420 only
type VideoFrame struct {
	Buffer    []byte
	Width     int
	Height    int
	YStride   int
	UStride   int
	VStride   int
	Timestamp int64
}

type RtcConnectionEventHandler struct {
	OnConnected                func(con *RtcConnection, conInfo *RtcConnectionInfo, reason int)
	OnDisconnected             func(con *RtcConnection, conInfo *RtcConnectionInfo, reason int)
	OnReconnecting             func(con *RtcConnection, conInfo *RtcConnectionInfo, reason int)
	OnReconnected              func(con *RtcConnection, conInfo *RtcConnectionInfo, reason int)
	OnConnectionLost           func(con *RtcConnection, conInfo *RtcConnectionInfo)
	OnConnectionFailure        func(con *RtcConnection, conInfo *RtcConnectionInfo, errCode int)
	OnTokenPrivilegeWillExpire func(con *RtcConnection, token string)
	OnTokenPrivilegeDidExpire  func(con *RtcConnection)
	OnUserJoined               func(con *RtcConnection, uid string)
	OnUserLeft                 func(con *RtcConnection, uid string, reason int)
	OnStreamMessageError       func(con *RtcConnection, uid string, streamId int, errCode int, missed int, cached int)
	OnStreamMessage            func(con *RtcConnection, uid string, streamId int, data []byte)
	// userMediaInfo: UserMediaInfoXxx
	// val: 0 for false, 1 for true
	OnUserInfoUpdated            func(con *RtcConnection, uid string, userMediaInfo int, val int)
	OnUserAudioTrackSubscribed   func(con *RtcConnection, uid string, remoteAudioTrack *RemoteAudioTrack)
	OnUserVideoTrackSubscribed   func(con *RtcConnection, uid string, info *VideoTrackInfo, remoteVideoTrack *RemoteVideoTrack)
	OnUserAudioTrackStateChanged func(con *RtcConnection, uid string, remoteAudioTrack *RemoteAudioTrack, state int, reason int, elapsed int)
	OnUserVideoTrackStateChanged func(con *RtcConnection, uid string, remoteAudioTrack *RemoteVideoTrack, state int, reason int, elapsed int)
}

type AudioFrameObserver struct {
	OnPlaybackAudioFrameBeforeMixing func(con *RtcConnection, channelId string, uid string, frame *PcmAudioFrame)
}

type VideoFrameObserver struct {
	OnFrame func(con *RtcConnection, channelId string, userId string, frame *VideoFrame)
}

type AudioEncoderConfiguration struct {
	AudioProfile int
}

// type SubscribeAudioConfig struct {
// 	SampleRate int
// 	Channels   int
// }

/**
 * Configurations for the RTC connection.
 */
type RtcConnectionConfig struct {
	/**
	 * Determines whether to subscribe to all audio streams automatically.
	 * - 1: (Default) Subscribe to all audio streams automatically.
	 * - 0: Do not subscribe to any audio stream automatically.
	 */
	AutoSubscribeAudio bool
	/**
	 * Determines whether to subscribe to all video streams automatically.
	 * - 1: (Default) Subscribe to all video streams automatically.
	 * - 0: Do not subscribe to any video stream automatically.
	 */
	AutoSubscribeVideo bool
	/**
	 * Determines whether to enable audio recording or playout.
	 * - true: It's used to publish audio and mix microphone, or subscribe audio and playout
	 * - false: It's used to publish extenal audio frame only without mixing microphone, or no need audio device to playout audio either
	 */
	EnableAudioRecordingOrPlayout bool
	/**
	 * The maximum sending bitrate.
	 */
	MaxSendBitrate int
	/**
	 * The minimum port.
	 */
	MinPort int
	/**
	 * The maximum port.
	 */
	MaxPort int
	/**
	 * The role of the user. The default user role is ClientRoleAudience.
	 */
	ClientRole     int
	ChannelProfile int
	/**
	 * Determines whether to receive audio media packet or not.
	 */
	AudioRecvMediaPacket bool
	/**
	 * Determines whether to receive video media packet or not.
	 */
	VideoRecvMediaPacket bool
}

// type RtcConnectionConfig struct {
// 	SubAudio       bool
// 	SubVideo       bool
// 	ClientRole     int
// 	ChannelProfile int

// 	SubAudioConfig     *SubscribeAudioConfig
// 	ConnectionHandler  *RtcConnectionEventHandler
// 	AudioFrameObserver *AudioFrameObserver
// 	VideoFrameObserver *VideoFrameObserver
// }

type RtcConnection struct {
	cConnection unsafe.Pointer
	cLocalUser  unsafe.Pointer
	// subAudioConfig     *SubscribeAudioConfig
	handler            *RtcConnectionEventHandler
	cHandler           *C.struct__rtc_conn_observer
	cLocalUserObserver *C.struct__local_user_observer
	audioObserver      *AudioFrameObserver
	cAudioObserver     *C.struct__audio_frame_observer
	videoObserver      *VideoFrameObserver
	cVideoObserver     unsafe.Pointer

	// videoSender *VideoSender
}

func NewConnection(cfg *RtcConnectionConfig) *RtcConnection {
	cCfg := CRtcConnectionConfig(cfg)
	defer FreeCRtcConnectionConfig(cCfg)

	ret := &RtcConnection{
		cConnection: C.agora_rtc_conn_create(agoraService.service, cCfg),
		// subAudioConfig: cfg.SubAudioConfig,
		handler:       nil,
		audioObserver: nil,
		videoObserver: nil,
	}
	ret.cLocalUser = C.agora_rtc_conn_get_local_user(ret.cConnection)
	// if ret.handler != nil {
	// ret.cHandler, ret.cLocalUserObserver = CRtcConnectionEventHandler()
	// C.agora_rtc_conn_register_observer(ret.cConnection, ret.cHandler)
	// C.agora_local_user_register_observer(ret.cLocalUser, ret.cLocalUserObserver)
	// }
	// if ret.subAudioConfig == nil {
	// 	ret.subAudioConfig = &SubscribeAudioConfig{
	// 		SampleRate: 16000,
	// 		Channels:   1,
	// 	}
	// }
	// C.agora_local_user_set_playback_audio_frame_before_mixing_parameters(
	// 	ret.cLocalUser, C.uint(ret.subAudioConfig.Channels), C.uint(ret.subAudioConfig.SampleRate))

	// if ret.audioObserver != nil {
	// ret.cAudioObserver = CAudioFrameObserver()
	// C.agora_local_user_register_audio_frame_observer(ret.cLocalUser, ret.cAudioObserver)
	// }

	// if ret.videoObserver != nil {
	// ret.cVideoObserver = CVideoFrameObserver()
	// C.agora_local_user_register_video_frame_observer(ret.cLocalUser, ret.cVideoObserver)
	// }

	agoraService.connectionRWMutex.Lock()
	agoraService.consByCCon[ret.cConnection] = ret
	agoraService.consByCLocalUser[ret.cLocalUser] = ret
	// agoraService.consByCVideoObserver[ret.cVideoObserver] = ret
	agoraService.connectionRWMutex.Unlock()
	return ret
}

func (conn *RtcConnection) Release() {
	if conn.cConnection == nil {
		return
	}
	agoraService.connectionRWMutex.Lock()
	delete(agoraService.consByCCon, conn.cConnection)
	delete(agoraService.consByCLocalUser, conn.cLocalUser)
	if conn.cVideoObserver != nil {
		delete(agoraService.consByCVideoObserver, conn.cVideoObserver)
	}
	agoraService.connectionRWMutex.Unlock()
	if conn.cAudioObserver != nil {
		C.agora_local_user_unregister_audio_frame_observer(conn.cLocalUser)
	}
	if conn.cVideoObserver != nil {
		C.agora_local_user_unregister_video_frame_observer(conn.cLocalUser, conn.cVideoObserver)
	}
	if conn.cLocalUserObserver != nil {
		C.agora_local_user_unregister_observer(conn.cLocalUser)
	}
	if conn.cHandler != nil {
		C.agora_rtc_conn_unregister_observer(conn.cConnection)
	}
	C.agora_rtc_conn_destroy(conn.cConnection)
	conn.cConnection = nil
	if conn.cAudioObserver != nil {
		FreeCAudioFrameObserver(conn.cAudioObserver)
		conn.cAudioObserver = nil
	}
	if conn.cVideoObserver != nil {
		FreeCVideoFrameObserver(conn.cVideoObserver)
		conn.cVideoObserver = nil
	}
	if conn.cLocalUserObserver != nil {
		FreeCLocalUserObserver(conn.cLocalUserObserver)
		conn.cLocalUserObserver = nil
	}
	if conn.cHandler != nil {
		FreeCRtcConnectionEventHandler(conn.cHandler)
		conn.cHandler = nil
	}
	conn.handler = nil
	conn.audioObserver = nil
	conn.videoObserver = nil
}

func (conn *RtcConnection) SetUserRole(role int) int {
	if conn.cConnection == nil {
		return -1
	}
	C.agora_local_user_set_user_role(conn.cLocalUser, C.int(role))
	return 0
}

func (conn *RtcConnection) Connect(token string, channel string, uid string) int {
	if conn.cConnection == nil {
		return -1
	}
	cChannel := C.CString(channel)
	cToken := C.CString(token)
	cUid := C.CString(uid)
	defer C.free(unsafe.Pointer(cChannel))
	defer C.free(unsafe.Pointer(cToken))
	defer C.free(unsafe.Pointer(cUid))
	return int(C.agora_rtc_conn_connect(conn.cConnection, cToken, cChannel, cUid))
}

func (conn *RtcConnection) Disconnect() int {
	if conn.cConnection == nil {
		return -1
	}
	return int(C.agora_rtc_conn_disconnect(conn.cConnection))
}

func (conn *RtcConnection) RenewToken(token string) int {
	if conn.cConnection == nil {
		return -1
	}
	cToken := C.CString(token)
	defer C.free(unsafe.Pointer(cToken))
	return int(C.agora_rtc_conn_renew_token(conn.cConnection, cToken))
}

func (conn *RtcConnection) CreateDataStream(reliable bool, ordered bool) (int, int) {
	if conn.cConnection == nil {
		return -1, -1
	}
	// int* stream_id, int reliable, int ordered
	cStreamId := C.int(-1)
	ret := int(C.agora_rtc_conn_create_data_stream(conn.cConnection, &cStreamId, CIntFromBool(reliable), CIntFromBool(ordered)))
	return int(cStreamId), ret
}

func (conn *RtcConnection) SendStreamMessage(streamId int, msg []byte) int {
	if conn.cConnection == nil {
		return -1
	}
	cMsg := C.CBytes(msg)
	defer C.free(cMsg)
	return int(C.agora_rtc_conn_send_stream_message(conn.cConnection, C.int(streamId), (*C.char)(cMsg), C.uint32_t(len(msg))))
}

func (conn *RtcConnection) SubscribeAudio(uid string) int {
	if conn.cLocalUser == nil {
		return -1
	}
	cUid := C.CString(uid)
	defer C.free(unsafe.Pointer(cUid))
	return int(C.agora_local_user_subscribe_audio(conn.cLocalUser, cUid))
}

func (conn *RtcConnection) UnsubscribeAudio(uid string) int {
	if conn.cLocalUser == nil {
		return -1
	}
	cUid := C.CString(uid)
	defer C.free(unsafe.Pointer(cUid))
	return int(C.agora_local_user_unsubscribe_audio(conn.cLocalUser, cUid))
}

func (conn *RtcConnection) SubscribeAllAudio() int {
	if conn.cLocalUser == nil {
		return -1
	}
	return int(C.agora_local_user_subscribe_all_audio(conn.cLocalUser))
}

func (conn *RtcConnection) UnsubscribeAllAudio() int {
	if conn.cLocalUser == nil {
		return -1
	}
	return int(C.agora_local_user_unsubscribe_all_audio(conn.cLocalUser))
}

func (conn *RtcConnection) SetParameters(parameters string) int {
	if conn.cConnection == nil {
		return -1
	}
	cParamHdl := C.agora_rtc_conn_get_agora_parameter(conn.cConnection)
	if cParamHdl == nil {
		return -1
	}
	cParameters := C.CString(parameters)
	defer C.free(unsafe.Pointer(cParameters))
	return int(C.agora_parameter_set_parameters(cParamHdl, cParameters))
}

func (conn *RtcConnection) SetAudioEncoderConfiguration(config *AudioEncoderConfiguration) int {
	if conn.cConnection == nil {
		return -1
	}
	cConfig := C.struct__audio_encoder_config{}
	cConfig.audio_profile = C.int(AudioProfileDefault)
	if config != nil {
		cConfig.audio_profile = C.int(config.AudioProfile)
	}
	return int(C.agora_local_user_set_audio_encoder_config(conn.cLocalUser, &cConfig))
}

func (conn *RtcConnection) PublishAudio(track *LocalAudioTrack) int {
	if conn.cConnection == nil {
		return -1
	}
	return int(C.agora_local_user_publish_audio(conn.cLocalUser, track.cTrack))
}

func (conn *RtcConnection) UnpublishAudio(track *LocalAudioTrack) int {
	if conn.cConnection == nil {
		return -1
	}
	return int(C.agora_local_user_unpublish_audio(conn.cLocalUser, track.cTrack))
}

func (conn *RtcConnection) PublishVideo(track *LocalVideoTrack) int {
	if conn.cConnection == nil {
		return -1
	}
	return int(C.agora_local_user_publish_video(conn.cLocalUser, track.cTrack))
}

func (conn *RtcConnection) UnpublishVideo(track *LocalVideoTrack) int {
	if conn.cConnection == nil {
		return -1
	}
	return int(C.agora_local_user_unpublish_video(conn.cLocalUser, track.cTrack))
}

func (conn *RtcConnection) SetPlaybackAudioFrameBeforeMixingParameters(channels int, sampleRate int) int {
	if conn.cConnection == nil {
		return -1
	}
	return int(C.agora_local_user_set_playback_audio_frame_before_mixing_parameters(conn.cLocalUser, C.uint(channels), C.uint(sampleRate)))
}

func (conn *RtcConnection) RegisterObserver(handler *RtcConnectionEventHandler) int {
	if conn.cConnection == nil || handler == nil {
		return -1
	}
	conn.handler = handler
	if conn.cHandler == nil && conn.cLocalUserObserver == nil {
		conn.cHandler, conn.cLocalUserObserver = CRtcConnectionEventHandler()
		C.agora_rtc_conn_register_observer(conn.cConnection, conn.cHandler)
		C.agora_local_user_register_observer(conn.cLocalUser, conn.cLocalUserObserver)
	}
	return 0
}

func (conn *RtcConnection) UnregisterObserver() int {
	conn.handler = nil
	return 0
}

func (conn *RtcConnection) RegisterAudioFrameObserver(observer *AudioFrameObserver) int {
	if conn.cConnection == nil || observer == nil {
		return -1
	}
	conn.audioObserver = observer
	if conn.cAudioObserver == nil {
		conn.cAudioObserver = CAudioFrameObserver()
		C.agora_local_user_register_audio_frame_observer(conn.cLocalUser, conn.cAudioObserver)
	}
	return 0
}

func (conn *RtcConnection) UnregisterAudioFrameObserver() int {
	conn.audioObserver = nil
	return 0
}

func (conn *RtcConnection) RegisterVideoFrameObserver(observer *VideoFrameObserver) int {
	if conn.cConnection == nil || observer == nil {
		return -1
	}
	conn.videoObserver = observer
	if conn.cVideoObserver == nil {
		conn.cVideoObserver = CVideoFrameObserver()
		C.agora_local_user_register_video_frame_observer(conn.cLocalUser, conn.cVideoObserver)
		agoraService.connectionRWMutex.Lock()
		agoraService.consByCVideoObserver[conn.cVideoObserver] = conn
		agoraService.connectionRWMutex.Unlock()
	}
	return 0
}

func (conn *RtcConnection) UnregisterVideoFrameObserver() int {
	conn.videoObserver = nil
	return 0
}
