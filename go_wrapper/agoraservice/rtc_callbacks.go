package agoraservice

/*
#cgo CFLAGS: -I../../agora_sdk/include_c/api2 -I../../agora_sdk/include_c/base
#cgo darwin arm64 LDFLAGS: -L../../agora_sdk_mac/arm64 -lAgoraRtcKit -lAgorafdkaac -lAgoraffmpeg
#cgo linux LDFLAGS: -L../../agora_sdk/ -lagora_rtc_sdk -lagora-fdkaac -lagora-core

#include "agora_local_user.h"
#include "agora_rtc_conn.h"
#include "agora_service.h"
#include "agora_media_base.h"
*/
import "C"
import "unsafe"

//export goOnConnected
func goOnConnected(cCon unsafe.Pointer, cConInfo *C.struct__rtc_conn_info, reason C.int) {
	agoraService.connectionRWMutex.RLock()
	con := agoraService.consByCCon[cCon]
	agoraService.connectionRWMutex.RUnlock()
	if con == nil || con.handler == nil || con.handler.OnConnected == nil {
		return
	}
	con.handler.OnConnected(con, GoRtcConnectionInfo(cConInfo), int(reason))
}

//export goOnDisconnected
func goOnDisconnected(cCon unsafe.Pointer, cConInfo *C.struct__rtc_conn_info, reason C.int) {
	agoraService.connectionRWMutex.RLock()
	con := agoraService.consByCCon[cCon]
	agoraService.connectionRWMutex.RUnlock()
	if con == nil || con.handler == nil || con.handler.OnDisconnected == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnDisconnected(con, GoRtcConnectionInfo(cConInfo), int(reason))

}

//export goOnReconnecting
func goOnReconnecting(cCon unsafe.Pointer, cConInfo *C.struct__rtc_conn_info, reason C.int) {
	agoraService.connectionRWMutex.RLock()
	con := agoraService.consByCCon[cCon]
	agoraService.connectionRWMutex.RUnlock()
	if con == nil || con.handler == nil || con.handler.OnReconnecting == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnReconnecting(con, GoRtcConnectionInfo(cConInfo), int(reason))
}

//export goOnReconnected
func goOnReconnected(cCon unsafe.Pointer, cConInfo *C.struct__rtc_conn_info, reason C.int) {

	agoraService.connectionRWMutex.RLock()
	con := agoraService.consByCCon[cCon]
	agoraService.connectionRWMutex.RUnlock()
	if con == nil || con.handler == nil || con.handler.onReconnected == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.handler.onReconnected(con, GoRtcConnectionInfo(cConInfo), int(reason))
}

//export goOnTokenPrivilegeWillExpire
func goOnTokenPrivilegeWillExpire(cCon unsafe.Pointer, ctoken *C.char) {
	agoraService.connectionRWMutex.RLock()
	con := agoraService.consByCCon[cCon]
	agoraService.connectionRWMutex.RUnlock()
	if con == nil || con.handler == nil || con.handler.OnTokenPrivilegeWillExpire == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnTokenPrivilegeWillExpire(con, C.GoString(ctoken))
}

//export goOnTokenPrivilegeDidExpire
func goOnTokenPrivilegeDidExpire(cCon unsafe.Pointer) {

	agoraService.connectionRWMutex.RLock()
	con := agoraService.consByCCon[cCon]
	agoraService.connectionRWMutex.RUnlock()
	if con == nil || con.handler == nil || con.handler.OnTokenPrivilegeDidExpire == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnTokenPrivilegeDidExpire(con)
}

//export goOnUserJoined
func goOnUserJoined(cCon unsafe.Pointer, uid *C.char) {

	agoraService.connectionRWMutex.RLock()
	con := agoraService.consByCCon[cCon]
	agoraService.connectionRWMutex.RUnlock()
	if con == nil || con.handler == nil || con.handler.OnUserJoined == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnUserJoined(con, C.GoString(uid))
}

//export goOnUserOffline
func goOnUserOffline(cCon unsafe.Pointer, uid *C.char, reason C.int) {

	agoraService.connectionRWMutex.RLock()
	con := agoraService.consByCCon[cCon]
	agoraService.connectionRWMutex.RUnlock()
	if con == nil || con.handler == nil || con.handler.OnUserLeft == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnUserLeft(con, C.GoString(uid), int(reason))
}

//export goOnStreamMessageError
func goOnStreamMessageError(cCon unsafe.Pointer, uid *C.char, streamId C.int, err C.int, missed C.int, cached C.int) {
	agoraService.connectionRWMutex.RLock()
	con := agoraService.consByCCon[cCon]
	agoraService.connectionRWMutex.RUnlock()
	if con == nil || con.handler == nil || con.handler.OnStreamMessageError == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnStreamMessageError(con, C.GoString(uid), int(streamId), int(err), int(missed), int(cached))
}

//export goOnStreamMessage
func goOnStreamMessage(cLocalUser unsafe.Pointer, uid *C.char, streamId C.int, data *C.char, length C.size_t) {
	agoraService.connectionRWMutex.RLock()
	con := agoraService.consByCLocalUser[cLocalUser]
	agoraService.connectionRWMutex.RUnlock()
	if con == nil || con.handler == nil || con.handler.OnStreamMessage == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnStreamMessage(con, C.GoString(uid), int(streamId), C.GoBytes(unsafe.Pointer(data), C.int(length)))
}
