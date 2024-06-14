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

//export goOnConnected
func goOnConnected(cCon unsafe.Pointer, cConInfo *C.struct__rtc_conn_info, reason C.int) {
	agoraService.connectionMutex.Lock()
	con := agoraService.connections[cCon]
	agoraService.connectionMutex.Unlock()
	if con == nil || con.handler == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnConnected(con, GoRtcConnectionInfo(cConInfo), int(reason))
}

//export goOnDisconnected
func goOnDisconnected(cCon unsafe.Pointer, cConInfo *C.struct__rtc_conn_info, reason C.int) {
	agoraService.connectionMutex.Lock()
	con := agoraService.connections[cCon]
	agoraService.connectionMutex.Unlock()
	if con == nil || con.handler == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnDisconnected(con, GoRtcConnectionInfo(cConInfo), int(reason))

}

//export goOnReconnecting
func goOnReconnecting(cCon unsafe.Pointer, cConInfo *C.struct__rtc_conn_info, reason C.int) {
	agoraService.connectionMutex.Lock()
	con := agoraService.connections[cCon]
	agoraService.connectionMutex.Unlock()
	if con == nil || con.handler == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnReconnecting(con, GoRtcConnectionInfo(cConInfo), int(reason))
}

//export goOnReconnected
func goOnReconnected(cCon unsafe.Pointer, cConInfo *C.struct__rtc_conn_info) {

	agoraService.connectionMutex.Lock()
	con := agoraService.connections[cCon]
	agoraService.connectionMutex.Unlock()
	if con == nil || con.handler == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.handler.onReconnected(con, GoRtcConnectionInfo(cConInfo))
}

//export goOnTokenPrivilegeWillExpire
func goOnTokenPrivilegeWillExpire(cCon unsafe.Pointer, cConInfo *C.struct__rtc_conn_info) {
	agoraService.connectionMutex.Lock()
	con := agoraService.connections[cCon]
	agoraService.connectionMutex.Unlock()
	if con == nil || con.handler == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnTokenPrivilegeWillExpire(con, GoRtcConnectionInfo(cConInfo))
}

//export goOnTokenPrivilegeDidExpire
func goOnTokenPrivilegeDidExpire(cCon unsafe.Pointer, cConInfo *C.struct__rtc_conn_info) {

	agoraService.connectionMutex.Lock()
	con := agoraService.connections[cCon]
	agoraService.connectionMutex.Unlock()
	if con == nil || con.handler == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnTokenPrivilegeDidExpire(con, GoRtcConnectionInfo(cConInfo))
}

//export goOnUserJoined
func goOnUserJoined(cCon unsafe.Pointer, cConInfo *C.struct__rtc_conn_info, uid C.int) {

	agoraService.connectionMutex.Lock()
	con := agoraService.connections[cCon]
	agoraService.connectionMutex.Unlock()
	if con == nil || con.handler == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnUserJoined(con, GoRtcConnectionInfo(cConInfo), int(uid))
}

//export goOnUserOffline
func goOnUserOffline(cCon unsafe.Pointer, cConInfo *C.struct__rtc_conn_info, uid C.int, reason C.int) {

	agoraService.connectionMutex.Lock()
	con := agoraService.connections[cCon]
	agoraService.connectionMutex.Unlock()
	if con == nil || con.handler == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnUserOffline(con, GoRtcConnectionInfo(cConInfo), int(uid), int(reason))
}

//export goOnStreamMessageError
func goOnStreamMessageError(cCon unsafe.Pointer, cConInfo *C.struct__rtc_conn_info, uid C.int, streamId C.int, error C.int) {

	agoraService.connectionMutex.Lock()
	con := agoraService.connections[cCon]
	agoraService.connectionMutex.Unlock()
	if con == nil || con.handler == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnStreamMessageError(con, GoRtcConnectionInfo(cConInfo), int(uid), int(streamId), int(error))
}
