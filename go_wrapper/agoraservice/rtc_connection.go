package agoraservice

/*
#cgo CFLAGS: -I../../agora_sdk/include_c/api2 -I../../agora_sdk/include_c/base
#cgo LDFLAGS: -L../../agora_sdk/ -lagora_rtc_sdk -lagora-fdkaac -lagora-ffmpeg

#include "agora_local_user.h"
#include "agora_rtc_conn.h"
#include "agora_service.h"
#include "agora_media_base.h"
*/
import "C"
import "unsafe"

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

type RtcConnectionEventHandler struct {
	OnConnected                func(con *RtcConnection, conInfo *RtcConnectionInfo, reason int)
	OnDisconnected             func(con *RtcConnection, conInfo *RtcConnectionInfo, reason int)
	OnReconnecting             func(con *RtcConnection, conInfo *RtcConnectionInfo, reason int)
	onReconnected              func(con *RtcConnection, conInfo *RtcConnectionInfo, reason int)
	OnTokenPrivilegeWillExpire func(con *RtcConnection, token string)
	OnTokenPrivilegeDidExpire  func(con *RtcConnection)
	OnUserJoined               func(con *RtcConnection, uid string)
	OnUserLeft                 func(con *RtcConnection, uid string, reason int)
	OnStreamMessageError       func(con *RtcConnection, uid string, streamId int, errCode int, missed int, cached int)
	OnStreamMessage            func(con *RtcConnection, uid string, streamId int, data []byte)
}

type RtcConnectionAudioFrameObserver struct {
	OnPlaybackAudioFrameBeforeMixing func(*RtcConnection, string, string, *PcmAudioFrame)
}

type SubscribeAudioConfig struct {
	SampleRate int
	Channels   int
}

type RtcConnectionConfig struct {
	SubAudio       bool
	SubVideo       bool
	ClientRole     int
	ChannelProfile int

	SubAudioConfig     *SubscribeAudioConfig
	ConnectionHandler  *RtcConnectionEventHandler
	AudioFrameObserver *RtcConnectionAudioFrameObserver
}

type RtcConnection struct {
	cConnection        unsafe.Pointer
	cLocalUser         unsafe.Pointer
	subAudioConfig     *SubscribeAudioConfig
	handler            *RtcConnectionEventHandler
	cHandler           *C.struct__rtc_conn_observer
	cLocalUserObserver *C.struct__local_user_observer
	audioObserver      *RtcConnectionAudioFrameObserver
	cAudioObserver     *C.struct__audio_frame_observer
}

func NewConnection(cfg *RtcConnectionConfig) *RtcConnection {
	cCfg := CRtcConnectionConfig(cfg)
	defer FreeCRtcConnectionConfig(cCfg)

	ret := &RtcConnection{
		cConnection:    C.agora_rtc_conn_create(agoraService.service, cCfg),
		subAudioConfig: cfg.SubAudioConfig,
		handler:        cfg.ConnectionHandler,
		audioObserver:  cfg.AudioFrameObserver,
	}
	ret.cLocalUser = C.agora_rtc_conn_get_local_user(ret.cConnection)
	// C.agora_local_user_subscribe_all_audio(ret.cLocalUser)
	if ret.handler != nil {
		ret.cHandler, ret.cLocalUserObserver = CRtcConnectionEventHandler(ret.handler)
		C.agora_rtc_conn_register_observer(ret.cConnection, ret.cHandler)
		C.agora_local_user_register_observer(ret.cLocalUser, ret.cLocalUserObserver)
	}
	if ret.subAudioConfig == nil {
		ret.subAudioConfig = &SubscribeAudioConfig{
			SampleRate: 16000,
			Channels:   1,
		}
	}
	C.agora_local_user_set_playback_audio_frame_before_mixing_parameters(
		ret.cLocalUser, C.uint(ret.subAudioConfig.Channels), C.uint(ret.subAudioConfig.SampleRate))

	if ret.audioObserver != nil {
		ret.cAudioObserver = CAudioFrameObserver(ret.audioObserver)
		C.agora_local_user_register_audio_frame_observer(ret.cLocalUser, ret.cAudioObserver)
	}

	agoraService.connectionRWMutex.Lock()
	agoraService.consByCCon[ret.cConnection] = ret
	agoraService.consByCLocalUser[ret.cLocalUser] = ret
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
	agoraService.connectionRWMutex.Unlock()
	C.agora_rtc_conn_destroy(conn.cConnection)
	conn.cConnection = nil
	if conn.cHandler != nil {
		FreeCRtcConnectionEventHandler(conn.cHandler)
		conn.cHandler = nil
	}
	if conn.cLocalUserObserver != nil {
		FreeCLocalUserObserver(conn.cLocalUserObserver)
		conn.cLocalUserObserver = nil
	}
	if conn.cAudioObserver != nil {
		FreeCAudioFrameObserver(conn.cAudioObserver)
		conn.cAudioObserver = nil
	}
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
