package agoraservice

/*
#cgo CFLAGS: -I../../agora_sdk/include/c/api2 -I../../agora_sdk/include/c/base

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
	GoRtcConnectionInfo(cConInfo, &con.connInfo)
	con.handler.OnConnected(con, &con.connInfo, int(reason))
}

//export goOnDisconnected
func goOnDisconnected(cCon unsafe.Pointer, cConInfo *C.struct__rtc_conn_info, reason C.int) {
	agoraService.connectionRWMutex.RLock()
	con := agoraService.consByCCon[cCon]
	agoraService.connectionRWMutex.RUnlock()
	if con == nil || con.handler == nil || con.handler.OnDisconnected == nil {
		return
	}
	GoRtcConnectionInfo(cConInfo, &con.connInfo)
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnDisconnected(con, &con.connInfo, int(reason))

}

//export goOnConnecting
func goOnConnecting(cCon unsafe.Pointer, cConInfo *C.struct__rtc_conn_info, reason C.int) {
	agoraService.connectionRWMutex.RLock()
	con := agoraService.consByCCon[cCon]
	agoraService.connectionRWMutex.RUnlock()
	if con == nil || con.handler == nil || con.handler.OnConnecting == nil {
		return
	}
	GoRtcConnectionInfo(cConInfo, &con.connInfo)
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnConnecting(con, &con.connInfo, int(reason))
}

//export goOnReconnecting
func goOnReconnecting(cCon unsafe.Pointer, cConInfo *C.struct__rtc_conn_info, reason C.int) {
	agoraService.connectionRWMutex.RLock()
	con := agoraService.consByCCon[cCon]
	agoraService.connectionRWMutex.RUnlock()
	if con == nil || con.handler == nil || con.handler.OnReconnecting == nil {
		return
	}
	GoRtcConnectionInfo(cConInfo, &con.connInfo)
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnReconnecting(con, &con.connInfo, int(reason))
}

//export goOnReconnected
func goOnReconnected(cCon unsafe.Pointer, cConInfo *C.struct__rtc_conn_info, reason C.int) {

	agoraService.connectionRWMutex.RLock()
	con := agoraService.consByCCon[cCon]
	agoraService.connectionRWMutex.RUnlock()
	if con == nil || con.handler == nil || con.handler.OnReconnected == nil {
		return
	}
	GoRtcConnectionInfo(cConInfo, &con.connInfo)
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnReconnected(con, &con.connInfo, int(reason))
}

//export goOnConnectionLost
func goOnConnectionLost(cCon unsafe.Pointer, cConInfo *C.struct__rtc_conn_info) {
	agoraService.connectionRWMutex.RLock()
	con := agoraService.consByCCon[cCon]
	agoraService.connectionRWMutex.RUnlock()
	if con == nil || con.handler == nil || con.handler.OnConnectionLost == nil {
		return
	}
	GoRtcConnectionInfo(cConInfo, &con.connInfo)
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnConnectionLost(con, &con.connInfo)
}

//export goOnConnectionFailure
func goOnConnectionFailure(cCon unsafe.Pointer, cConInfo *C.struct__rtc_conn_info, reason C.int) {

	agoraService.connectionRWMutex.RLock()
	con := agoraService.consByCCon[cCon]
	agoraService.connectionRWMutex.RUnlock()
	if con == nil || con.handler == nil || con.handler.OnConnectionFailure == nil {
		return
	}
	GoRtcConnectionInfo(cConInfo, &con.connInfo)
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnConnectionFailure(con, &con.connInfo, int(reason))
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

//export goOnError
func goOnError(cCon unsafe.Pointer, err C.int, msg *C.char) {
	agoraService.connectionRWMutex.RLock()
	con := agoraService.consByCCon[cCon]
	agoraService.connectionRWMutex.RUnlock()
	if con == nil || con.handler == nil || con.handler.OnError == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnError(con, int(err), C.GoString(msg))
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
	if con == nil || con.localUserObserver == nil || con.localUserObserver.OnStreamMessage == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.localUserObserver.OnStreamMessage(con.GetLocalUser(), C.GoString(uid), int(streamId), C.GoBytes(unsafe.Pointer(data), C.int(length)))
}

//export goOnUserInfoUpdated
func goOnUserInfoUpdated(cLocalUser unsafe.Pointer, uid *C.char, msg C.int, val C.int) {
	agoraService.connectionRWMutex.RLock()
	con := agoraService.consByCLocalUser[cLocalUser]
	agoraService.connectionRWMutex.RUnlock()
	if con == nil || con.localUserObserver == nil || con.localUserObserver.OnUserInfoUpdated == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.localUserObserver.OnUserInfoUpdated(con.GetLocalUser(), C.GoString(uid), int(msg), int(val))
}

//export goOnUserAudioTrackSubscribed
func goOnUserAudioTrackSubscribed(cLocalUser unsafe.Pointer, uid *C.char, cRemoteAudioTrack unsafe.Pointer) {
	agoraService.connectionRWMutex.RLock()
	con := agoraService.consByCLocalUser[cLocalUser]
	agoraService.connectionRWMutex.RUnlock()
	if con == nil || con.localUserObserver == nil || con.localUserObserver.OnUserAudioTrackSubscribed == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.localUserObserver.OnUserAudioTrackSubscribed(con.GetLocalUser(), C.GoString(uid), NewRemoteAudioTrack(cRemoteAudioTrack))
}

//export goOnUserVideoTrackSubscribed
func goOnUserVideoTrackSubscribed(cLocalUser unsafe.Pointer, uid *C.char, info *C.struct__video_track_info, cRemoteVideoTrack unsafe.Pointer) {
	agoraService.connectionRWMutex.RLock()
	con := agoraService.consByCLocalUser[cLocalUser]
	agoraService.connectionRWMutex.RUnlock()
	if con == nil || con.localUserObserver == nil || con.localUserObserver.OnUserVideoTrackSubscribed == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.localUserObserver.OnUserVideoTrackSubscribed(con.GetLocalUser(), C.GoString(uid), GoVideoTrackInfo(info), con.NewRemoteVideoTrack(cRemoteVideoTrack))
}

//export goOnUserAudioTrackStateChanged
func goOnUserAudioTrackStateChanged(cLocalUser unsafe.Pointer, uid *C.char, cRemoteAudioTrack unsafe.Pointer, state C.int, reason C.int, elapsed C.int) {
	agoraService.connectionRWMutex.RLock()
	con := agoraService.consByCLocalUser[cLocalUser]
	agoraService.connectionRWMutex.RUnlock()
	if con == nil || con.localUserObserver == nil || con.localUserObserver.OnUserAudioTrackStateChanged == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.localUserObserver.OnUserAudioTrackStateChanged(con.GetLocalUser(), C.GoString(uid), NewRemoteAudioTrack(cRemoteAudioTrack), int(state), int(reason), int(elapsed))
}

//export goOnUserVideoTrackStateChanged
func goOnUserVideoTrackStateChanged(cLocalUser unsafe.Pointer, uid *C.char, cRemoteVideoTrack unsafe.Pointer, state C.int, reason C.int, elapsed C.int) {
	agoraService.connectionRWMutex.RLock()
	con := agoraService.consByCLocalUser[cLocalUser]
	agoraService.connectionRWMutex.RUnlock()
	if con == nil || con.localUserObserver == nil || con.localUserObserver.OnUserVideoTrackStateChanged == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.localUserObserver.OnUserVideoTrackStateChanged(con.GetLocalUser(), C.GoString(uid), con.NewRemoteVideoTrack(cRemoteVideoTrack), int(state), int(reason), int(elapsed))
}
