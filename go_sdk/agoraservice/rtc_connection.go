package agoraservice

/*
#cgo CFLAGS: -I../../agora_sdk/include/c/api2 -I../../agora_sdk/include/c/base

#include "agora_local_user.h"
#include "agora_rtc_conn.h"
#include "agora_service.h"
#include "agora_media_base.h"
#include "agora_parameter.h"
*/
import "C"
import (
	"strconv"
	"unsafe"
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

type AudioFrameObserverAudioParams struct {
	SampleRate     int
	Channels       int
	Mode           RawAudioFrameOpModeType
	SamplesPerCall int
}

type RtcConnectionObserver struct {
	OnConnected                func(con *RtcConnection, conInfo *RtcConnectionInfo, reason int)
	OnDisconnected             func(con *RtcConnection, conInfo *RtcConnectionInfo, reason int)
	OnConnecting               func(con *RtcConnection, conInfo *RtcConnectionInfo, reason int)
	OnReconnecting             func(con *RtcConnection, conInfo *RtcConnectionInfo, reason int)
	OnReconnected              func(con *RtcConnection, conInfo *RtcConnectionInfo, reason int)
	OnConnectionLost           func(con *RtcConnection, conInfo *RtcConnectionInfo)
	OnConnectionFailure        func(con *RtcConnection, conInfo *RtcConnectionInfo, errCode int)
	OnTokenPrivilegeWillExpire func(con *RtcConnection, token string)
	OnTokenPrivilegeDidExpire  func(con *RtcConnection)
	OnUserJoined               func(con *RtcConnection, uid string)
	OnUserLeft                 func(con *RtcConnection, uid string, reason int)
	OnError                    func(con *RtcConnection, err int, msg string)
	OnStreamMessageError       func(con *RtcConnection, uid string, streamId int, errCode int, missed int, cached int)
}

type LocalUserObserver struct {
	OnStreamMessage func(localUser *LocalUser, uid string, streamId int, data []byte)
	// userMediaInfo: UserMediaInfoXxx
	// val: 0 for false, 1 for true
	OnUserInfoUpdated          func(localUser *LocalUser, uid string, userMediaInfo int, val int)
	OnUserAudioTrackSubscribed func(localUser *LocalUser, uid string, remoteAudioTrack *RemoteAudioTrack)
	OnUserVideoTrackSubscribed func(localUser *LocalUser, uid string, info *VideoTrackInfo, remoteVideoTrack *RemoteVideoTrack)
	/*
		for Mute/Unmute
		(state== 0 and reason == 5): mute
		(state== 2 and reason == 6): unmute
	*/
	OnUserAudioTrackStateChanged func(localUser *LocalUser, uid string, remoteAudioTrack *RemoteAudioTrack, state int, reason int, elapsed int)
	OnUserVideoTrackStateChanged func(localUser *LocalUser, uid string, remoteAudioTrack *RemoteVideoTrack, state int, reason int, elapsed int)
	//  void (*on_audio_publish_state_changed)(AGORA_HANDLE agora_local_user, const char* channel, int old_state, int new_state, int elapse_since_last_state);
	OnAudioPublishStateChanged func(localUser *LocalUser, channelId string, oldState int, newState int, elapsed int)
	OnAudioVolumeIndication    func(localUser *LocalUser, audioVolumeInfo []*AudioVolumeInfo, speakerNumber int, totalVolume int)
}

type AudioFrameObserver struct {
	OnRecordAudioFrame                func(localUser *LocalUser, channelId string, frame *AudioFrame) bool
	OnPlaybackAudioFrame              func(localUser *LocalUser, channelId string, frame *AudioFrame) bool
	OnMixedAudioFrame                 func(localUser *LocalUser, channelId string, frame *AudioFrame) bool
	OnEarMonitoringAudioFrame         func(localUser *LocalUser, frame *AudioFrame) bool
	OnPlaybackAudioFrameBeforeMixing  func(localUser *LocalUser, channelId string, uid string, frame *AudioFrame) bool
	OnGetAudioFramePosition           func(localUser *LocalUser) int
	OnGetPlaybackAudioFrameParam      func(localUser *LocalUser) AudioFrameObserverAudioParams
	OnGetRecordAudioFrameParam        func(localUser *LocalUser) AudioFrameObserverAudioParams
	OnGetMixedAudioFrameParam         func(localUser *LocalUser) AudioFrameObserverAudioParams
	OnGetEarMonitoringAudioFrameParam func(localUser *LocalUser) AudioFrameObserverAudioParams
}

type VideoFrameObserver struct {
	OnFrame func(channelId string, userId string, frame *VideoFrame) bool
}

type VideoEncodedFrameObserver struct {
	OnEncodedVideoFrame func(uid string, imageBuffer []byte,
		frameInfo *EncodedVideoFrameInfo) bool
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
	 * - true: (Default) Subscribe to all audio streams automatically.
	 * - false: Do not subscribe to any audio stream automatically.
	 */
	AutoSubscribeAudio bool
	/**
	 * Determines whether to subscribe to all video streams automatically.
	 * - true: (Default) Subscribe to all video streams automatically.
	 * - false: Do not subscribe to any video stream automatically.
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
	ClientRole     ClientRole
	ChannelProfile ChannelProfile
	/**
	 * Determines whether to receive audio media packet or not.
	 */
	AudioRecvMediaPacket bool
	/**
	 * Determines whether to receive video media packet or not.
	 */
	VideoRecvMediaPacket bool
}

type RtcConnection struct {
	cConnection unsafe.Pointer
	connInfo    RtcConnectionInfo
	localUser   *LocalUser
	parameter   *AgoraParameter
	// cLocalUser  unsafe.Pointer
	// subAudioConfig     *SubscribeAudioConfig
	handler               *RtcConnectionObserver
	cHandler              *C.struct__rtc_conn_observer
	localUserObserver     *LocalUserObserver
	cLocalUserObserver    *C.struct__local_user_observer
	audioObserver         *AudioFrameObserver
	cAudioObserver        *C.struct__audio_frame_observer
	videoObserver         *VideoFrameObserver
	cVideoObserver        unsafe.Pointer
	encodedVideoObserver  *VideoEncodedFrameObserver
	cEncodedVideoObserver unsafe.Pointer

	// remoteVideoRWMutex          *sync.RWMutex
	// remoteEncodedVideoReceivers map[*VideoEncodedImageReceiver]*videoEncodedImageReceiverInner
}

func NewRtcConnection(cfg *RtcConnectionConfig) *RtcConnection {
	cCfg := CRtcConnectionConfig(cfg)
	defer FreeCRtcConnectionConfig(cCfg)

	ret := &RtcConnection{
		cConnection: C.agora_rtc_conn_create(agoraService.service, cCfg),
		// subAudioConfig: cfg.SubAudioConfig,
		handler:              nil,
		localUserObserver:    nil,
		audioObserver:        nil,
		videoObserver:        nil,
		encodedVideoObserver: nil,
		// remoteVideoRWMutex:          &sync.RWMutex{},
		// remoteEncodedVideoReceivers: make(map[*VideoEncodedImageReceiver]*videoEncodedImageReceiverInner),
	}
	ret.localUser = &LocalUser{
		connection: ret,
		cLocalUser: C.agora_rtc_conn_get_local_user(ret.cConnection),
	}
	ret.parameter = &AgoraParameter{
		cParameter: C.agora_rtc_conn_get_agora_parameter(ret.cConnection),
	}

	// save to sync map
	agoraService.setConFromHandle(ret.cConnection, ret, ConTypeCCon)
	agoraService.setConFromHandle(ret.localUser.cLocalUser, ret, ConTypeCLocalUser)

	return ret
}

func (conn *RtcConnection) Release() {
	if conn.cConnection == nil {
		return
	}
	// delete from sync map
	agoraService.deleteConFromHandle(conn.cConnection, ConTypeCCon)
	agoraService.deleteConFromHandle(conn.localUser.cLocalUser, ConTypeCLocalUser)
	if conn.cVideoObserver != nil {
		agoraService.deleteConFromHandle(conn.cVideoObserver, ConTypeCVideoObserver)
	}
	if conn.cEncodedVideoObserver != nil {
		agoraService.deleteConFromHandle(conn.cEncodedVideoObserver, ConTypeCEncodedVideoObserver)
	}

	// get all receiverInners
	// encodedVideoReceiversInners := make([]*videoEncodedImageReceiverInner, 0, 10)
	// conn.remoteVideoRWMutex.RLock()
	// for _, receiverInner := range conn.remoteEncodedVideoReceivers {
	// 	encodedVideoReceiversInners = append(encodedVideoReceiversInners, receiverInner)
	// }
	// conn.remoteVideoRWMutex.RUnlock()
	// // remove all receiverInners from service
	// agoraService.remoteVideoRWMutex.Lock()
	// for _, receiverInner := range encodedVideoReceiversInners {
	// 	delete(agoraService.remoteEncodedVideoReceivers, receiverInner.cReceiver)
	// }
	// agoraService.remoteVideoRWMutex.Unlock()

	localUser := conn.localUser
	if conn.cAudioObserver != nil {
		C.agora_local_user_unregister_audio_frame_observer(localUser.cLocalUser)
	}
	if conn.cVideoObserver != nil {
		C.agora_local_user_unregister_video_frame_observer(localUser.cLocalUser, conn.cVideoObserver)
	}
	if conn.cEncodedVideoObserver != nil {
		C.agora_local_user_unregister_video_encoded_frame_observer(localUser.cLocalUser, conn.cEncodedVideoObserver)
	}
	if conn.cLocalUserObserver != nil {
		C.agora_local_user_unregister_observer(localUser.cLocalUser)
	}
	if conn.cHandler != nil {
		C.agora_rtc_conn_unregister_observer(conn.cConnection)
	}
	C.agora_rtc_conn_destroy(conn.cConnection)

	// clear all receiverInners
	// conn.remoteVideoRWMutex.Lock()
	// conn.remoteEncodedVideoReceivers = make(map[*VideoEncodedImageReceiver]*videoEncodedImageReceiverInner)
	// conn.remoteVideoRWMutex.Unlock()

	// for _, receiverInner := range encodedVideoReceiversInners {
	// 	receiverInner.release()
	// }
	// encodedVideoReceiversInners = nil

	conn.cConnection = nil
	if conn.cAudioObserver != nil {
		FreeCAudioFrameObserver(conn.cAudioObserver)
		conn.cAudioObserver = nil
	}
	if conn.cVideoObserver != nil {
		FreeCVideoFrameObserver(conn.cVideoObserver)
		conn.cVideoObserver = nil
	}
	if conn.cEncodedVideoObserver != nil {
		FreeCEncodedVideoFrameObserver(conn.cEncodedVideoObserver)
		conn.cEncodedVideoObserver = nil
	}
	if conn.cLocalUserObserver != nil {
		FreeCLocalUserObserver(conn.cLocalUserObserver)
		conn.cLocalUserObserver = nil
	}
	if conn.cHandler != nil {
		FreeCRtcConnectionObserver(conn.cHandler)
		conn.cHandler = nil
	}
	conn.parameter = nil
	conn.localUser = nil
	localUser.connection = nil
	localUser.cLocalUser = nil
	localUser = nil
	conn.handler = nil
	conn.localUserObserver = nil
	conn.audioObserver = nil
	conn.videoObserver = nil
}

func (conn *RtcConnection) GetLocalUser() *LocalUser {
	return conn.localUser
}

func (conn *RtcConnection) GetAgoraParameter() *AgoraParameter {
	return conn.parameter
}

func (conn *RtcConnection) GetConnectionInfo() *RtcConnectionInfo {
	return &conn.connInfo
}

func (conn *RtcConnection) Connect(token string, channel string, uid string) int {
	if conn.cConnection == nil {
		return -1
	}
	conn.connInfo.ChannelId = channel
	conn.connInfo.LocalUserId = uid
	uidInt, _ := strconv.Atoi(uid)
	conn.connInfo.InternalUid = uint(uidInt)
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

func (conn *RtcConnection) RegisterObserver(handler *RtcConnectionObserver) int {
	if conn.cConnection == nil || handler == nil {
		return -1
	}
	// avoid re-register observer
	if conn.handler == handler {
		return 0
	}
	// unregister old observer
	conn.UnregisterObserver()

	// register new observer
	conn.handler = handler
	if conn.cHandler == nil {
		conn.cHandler = CRtcConnectionObserver()
		C.agora_rtc_conn_register_observer(conn.cConnection, conn.cHandler)
	}
	return 0

}

func (conn *RtcConnection) UnregisterObserver() int {
	// check if need to unregister
	if conn.cConnection == nil {
		return 0
	}
	if conn.cHandler != nil {
		C.agora_rtc_conn_unregister_observer(conn.cConnection)
		FreeCRtcConnectionObserver(conn.cHandler)
	}
	conn.cHandler = nil
	conn.handler = nil
	return 0
}

func (conn *RtcConnection) registerLocalUserObserver(handler *LocalUserObserver) int {
	if conn.cConnection == nil || handler == nil {
		return -1
	}
	// avoid re-register observer
	if conn.localUserObserver == handler {
		return 0
	}
	// unregister old observer
	conn.unregisterLocalUserObserver()

	// register new observer
	conn.localUserObserver = handler
	if conn.cLocalUserObserver == nil {
		conn.cLocalUserObserver = CLocalUserObserver()
		C.agora_local_user_register_observer(conn.localUser.cLocalUser, conn.cLocalUserObserver)
	}
	return 0
}

func (conn *RtcConnection) unregisterLocalUserObserver() int {
	// check if need to unregister
	if conn.cConnection == nil {
		return 0
	}
	if conn.cLocalUserObserver != nil {
		C.agora_local_user_unregister_observer(conn.localUser.cLocalUser)
		FreeCLocalUserObserver(conn.cLocalUserObserver)
	}
	conn.cLocalUserObserver = nil
	conn.localUserObserver = nil
	return 0
}

func (conn *RtcConnection) registerAudioFrameObserver(observer *AudioFrameObserver) int {
	if conn.cConnection == nil || observer == nil {
		return -1
	}
	// avoid re-register observer
	if conn.audioObserver == observer {
		return 0
	}
	// unregister old observer
	conn.unregisterAudioFrameObserver()

	conn.audioObserver = observer
	if conn.cAudioObserver == nil {
		conn.cAudioObserver = CAudioFrameObserver()
		C.agora_local_user_register_audio_frame_observer(conn.localUser.cLocalUser, conn.cAudioObserver)
	}
	return 0
}

func (conn *RtcConnection) unregisterAudioFrameObserver() int {
	// check if need to unregister
	if conn.cConnection == nil {
		return 0
	}
	if conn.cAudioObserver != nil {
		C.agora_local_user_unregister_audio_frame_observer(conn.localUser.cLocalUser)
		FreeCAudioFrameObserver(conn.cAudioObserver)
	}
	conn.cAudioObserver = nil
	conn.audioObserver = nil
	return 0
}

func (conn *RtcConnection) registerVideoFrameObserver(observer *VideoFrameObserver) int {
	if conn.cConnection == nil || observer == nil {
		return -1
	}
	// avoid re-register observer
	if conn.videoObserver == observer {
		return 0
	}
	// unregister old observer
	conn.unregisterVideoFrameObserver()

	conn.videoObserver = observer
	if conn.cVideoObserver == nil {
		conn.cVideoObserver = CVideoFrameObserver()
		// store to sync map
		agoraService.setConFromHandle(conn.cVideoObserver, conn, ConTypeCVideoObserver)
		C.agora_local_user_register_video_frame_observer(conn.localUser.cLocalUser, conn.cVideoObserver)
	}
	return 0
}

func (conn *RtcConnection) unregisterVideoFrameObserver() int {
	// check if need to unregister
	if conn.cConnection == nil {
		return 0
	}
	if conn.cVideoObserver != nil {
		C.agora_local_user_unregister_video_frame_observer(conn.localUser.cLocalUser, conn.cVideoObserver)
		// delete from sync map
		agoraService.deleteConFromHandle(conn.cVideoObserver, ConTypeCVideoObserver)
		FreeCVideoFrameObserver(conn.cVideoObserver)
		conn.cVideoObserver = nil
	}
	conn.cVideoObserver = nil
	conn.videoObserver = nil
	return 0
}

func (conn *RtcConnection) registerVideoEncodedFrameObserver(observer *VideoEncodedFrameObserver) int {
	if conn.cConnection == nil || observer == nil {
		return -1
	}
	// avoid re-register observer
	if conn.encodedVideoObserver == observer {
		return 0
	}
	// unregister old observer
	conn.unregisterVideoEncodedFrameObserver()

	conn.encodedVideoObserver = observer
	if conn.cEncodedVideoObserver == nil {
		conn.cEncodedVideoObserver = CVideoEncodedFrameObserver()
		// store to sync map
		agoraService.setConFromHandle(conn.cEncodedVideoObserver, conn, ConTypeCEncodedVideoObserver)
		C.agora_local_user_register_video_encoded_frame_observer(conn.localUser.cLocalUser, conn.cEncodedVideoObserver)
	}
	return 0
}

func (conn *RtcConnection) unregisterVideoEncodedFrameObserver() int {
	// check if need to unregister
	if conn.cConnection == nil {
		return 0
	}
	if conn.cEncodedVideoObserver != nil {
		C.agora_local_user_unregister_video_encoded_frame_observer(conn.localUser.cLocalUser, conn.cEncodedVideoObserver)
		// delete from sync map
		agoraService.deleteConFromHandle(conn.cEncodedVideoObserver, ConTypeCEncodedVideoObserver)
		FreeCEncodedVideoFrameObserver(conn.cEncodedVideoObserver)
		conn.cEncodedVideoObserver = nil
	}
	conn.cEncodedVideoObserver = nil
	conn.encodedVideoObserver = nil
	return 0
}
