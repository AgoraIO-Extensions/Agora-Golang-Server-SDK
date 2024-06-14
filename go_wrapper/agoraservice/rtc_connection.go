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

type RtcConnectionEventHandler struct {
	OnConnected                func(*RtcConnection, *RtcConnectionInfo, int)
	OnDisconnected             func(*RtcConnection, *RtcConnectionInfo, int)
	OnReconnecting             func(*RtcConnection, *RtcConnectionInfo, int)
	onReconnected              func(*RtcConnection, *RtcConnectionInfo)
	OnTokenPrivilegeDidExpire  func(*RtcConnection, *RtcConnectionInfo)
	OnTokenPrivilegeWillExpire func(*RtcConnection, *RtcConnectionInfo)
	OnUserJoined               func(*RtcConnection, *RtcConnectionInfo, int)
	OnUserOffline              func(*RtcConnection, *RtcConnectionInfo, int, int)
	OnStreamMessageError       func(*RtcConnection, *RtcConnectionInfo, int, int, int)
}

type RtcConnectionConfig struct {
	SubAudio       bool
	SubVideo       bool
	ClientRole     int
	ChannelProfile int
}

type RtcConnection struct {
	cConnection unsafe.Pointer
	cLocalUser  unsafe.Pointer
	handler     *RtcConnectionEventHandler
	cHandler    *C.struct__rtc_conn_observer
}

func NewConnection(cfg *RtcConnectionConfig, handler *RtcConnectionEventHandler) *RtcConnection {
	cCfg := CRtcConnectionConfig(cfg)
	defer FreeCRtcConnectionConfig(cCfg)

	ret := &RtcConnection{
		cConnection: C.agora_rtc_conn_create(agoraService.service, cCfg),
		handler:     handler,
	}
	ret.cLocalUser = C.agora_rtc_conn_get_local_user(ret.cConnection)
	ret.cHandler = CRtcConnectionEventHandler(handler)
	C.agora_rtc_conn_register_observer(ret.cConnection, ret.cHandler)
	C.agora_local_user_set_playback_audio_frame_before_mixing_parameters(ret.cLocalUser, 1, 16000)

	agoraService.connectionMutex.Lock()
	agoraService.connections[ret.cConnection] = ret
	agoraService.connectionMutex.Unlock()
	return ret
}

func (conn *RtcConnection) Release() {
	if conn.cConnection == nil {
		return
	}
	agoraService.connectionMutex.Lock()
	delete(agoraService.connections, conn.cConnection)
	agoraService.connectionMutex.Unlock()
	C.agora_rtc_conn_destroy(conn.cConnection)
	if conn.cHandler != nil {
		FreeCRtcConnectionEventHandler(conn.cHandler)
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
