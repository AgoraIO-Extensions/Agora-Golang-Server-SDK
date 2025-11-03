package agoraservice

/*
#cgo CFLAGS: -I../../agora_sdk/include/c/api2 -I../../agora_sdk/include/c/base


#include <string.h>
#include "agora_local_user.h"
#include "agora_rtc_conn.h"
#include "agora_service.h"
#include "agora_media_base.h"
*/
import "C"
import (
	"unsafe"
	"fmt"
)

//export goOnConnected
func goOnConnected(cCon unsafe.Pointer, cConInfo *C.struct__rtc_conn_info, reason C.int) {
	//validity check
	if cCon == nil {
		return
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cCon, ConTypeCCon)
	if con == nil || con.handler == nil || con.handler.OnConnected == nil {
		return
	}
	GoRtcConnectionInfo(cConInfo, &con.connInfo)
	con.handler.OnConnected(con, &con.connInfo, int(reason))
}

//export goOnDisconnected
func goOnDisconnected(cCon unsafe.Pointer, cConInfo *C.struct__rtc_conn_info, reason C.int) {
	//validity check
	if cCon == nil {
		return
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cCon, ConTypeCCon)
	if con == nil || con.handler == nil || con.handler.OnDisconnected == nil {
		return
	}
	GoRtcConnectionInfo(cConInfo, &con.connInfo)
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnDisconnected(con, &con.connInfo, int(reason))

}

//export goOnConnecting
func goOnConnecting(cCon unsafe.Pointer, cConInfo *C.struct__rtc_conn_info, reason C.int) {
	//validity check
	if cCon == nil {
		return
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cCon, ConTypeCCon)
	if con == nil || con.handler == nil || con.handler.OnConnecting == nil {
		return
	}
	GoRtcConnectionInfo(cConInfo, &con.connInfo)
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnConnecting(con, &con.connInfo, int(reason))
}

//export goOnReconnecting
func goOnReconnecting(cCon unsafe.Pointer, cConInfo *C.struct__rtc_conn_info, reason C.int) {
	//validity check
	if cCon == nil {
		return
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cCon, ConTypeCCon)
	if con == nil || con.handler == nil || con.handler.OnReconnecting == nil {
		return
	}
	GoRtcConnectionInfo(cConInfo, &con.connInfo)
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnReconnecting(con, &con.connInfo, int(reason))
}

//export goOnReconnected
func goOnReconnected(cCon unsafe.Pointer, cConInfo *C.struct__rtc_conn_info, reason C.int) {

	//validity check
	if cCon == nil {
		return
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cCon, ConTypeCCon)
	if con == nil || con.handler == nil || con.handler.OnReconnected == nil {
		return
	}
	GoRtcConnectionInfo(cConInfo, &con.connInfo)
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnReconnected(con, &con.connInfo, int(reason))
}

//export goOnConnectionLost
func goOnConnectionLost(cCon unsafe.Pointer, cConInfo *C.struct__rtc_conn_info) {
	//validity check
	if cCon == nil {
		return
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cCon, ConTypeCCon)
	if con == nil || con.handler == nil || con.handler.OnConnectionLost == nil {
		return
	}
	GoRtcConnectionInfo(cConInfo, &con.connInfo)
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnConnectionLost(con, &con.connInfo)
}

//export goOnConnectionFailure
func goOnConnectionFailure(cCon unsafe.Pointer, cConInfo *C.struct__rtc_conn_info, reason C.int) {
	//validity check
	if cCon == nil {
		return
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cCon, ConTypeCCon)
	if con == nil || con.handler == nil || con.handler.OnConnectionFailure == nil {
		return
	}
	GoRtcConnectionInfo(cConInfo, &con.connInfo)
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnConnectionFailure(con, &con.connInfo, int(reason))
}

//export goOnTokenPrivilegeWillExpire
func goOnTokenPrivilegeWillExpire(cCon unsafe.Pointer, ctoken *C.char) {
	//validity check
	if cCon == nil {
		return
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cCon, ConTypeCCon)
	if con == nil || con.handler == nil || con.handler.OnTokenPrivilegeWillExpire == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnTokenPrivilegeWillExpire(con, C.GoString(ctoken))
}

//export goOnTokenPrivilegeDidExpire
func goOnTokenPrivilegeDidExpire(cCon unsafe.Pointer) {
	//validity check
	if cCon == nil {
		return
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cCon, ConTypeCCon)
	if con == nil || con.handler == nil || con.handler.OnTokenPrivilegeDidExpire == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnTokenPrivilegeDidExpire(con)
}

//export goOnUserJoined
func goOnUserJoined(cCon unsafe.Pointer, uid *C.char) {
	//validity check
	if cCon == nil {
		return
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cCon, ConTypeCCon)
	if con == nil || con.handler == nil || con.handler.OnUserJoined == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnUserJoined(con, C.GoString(uid))
}

//export goOnUserOffline
func goOnUserOffline(cCon unsafe.Pointer, uid *C.char, reason C.int) {
	//validity check
	if cCon == nil {
		return
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cCon, ConTypeCCon)
	if con == nil || con.handler == nil || con.handler.OnUserLeft == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnUserLeft(con, C.GoString(uid), int(reason))
}

//export goOnError
func goOnError(cCon unsafe.Pointer, err C.int, msg *C.char) {
	//validity check
	if cCon == nil {
		return
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cCon, ConTypeCCon)
	if con == nil || con.handler == nil || con.handler.OnError == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnError(con, int(err), C.GoString(msg))
}

//export goOnStreamMessageError
func goOnStreamMessageError(cCon unsafe.Pointer, uid *C.char, streamId C.int, err C.int, missed C.int, cached C.int) {
	//validity check
	if cCon == nil {
		return
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cCon, ConTypeCCon)
	if con == nil || con.handler == nil || con.handler.OnStreamMessageError == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.handler.OnStreamMessageError(con, C.GoString(uid), int(streamId), int(err), int(missed), int(cached))
}

//export goOnStreamMessage
func goOnStreamMessage(cLocalUser unsafe.Pointer, uid *C.char, streamId C.int, data *C.char, length C.size_t) {
	//validity check
	if cLocalUser == nil {
		return
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cLocalUser, ConTypeCLocalUser)
	if con == nil || con.localUserObserver == nil || con.localUserObserver.OnStreamMessage == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.localUserObserver.OnStreamMessage(con.GetLocalUser(), C.GoString(uid), int(streamId), C.GoBytes(unsafe.Pointer(data), C.int(length)))
}

//export goOnUserInfoUpdated
func goOnUserInfoUpdated(cLocalUser unsafe.Pointer, uid *C.char, msg C.int, val C.int) {
	//validity check
	if cLocalUser == nil {
		return
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cLocalUser, ConTypeCLocalUser)
	if con == nil || con.localUserObserver == nil || con.localUserObserver.OnUserInfoUpdated == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.localUserObserver.OnUserInfoUpdated(con.GetLocalUser(), C.GoString(uid), int(msg), int(val))
}

//export goOnUserAudioTrackSubscribed
func goOnUserAudioTrackSubscribed(cLocalUser unsafe.Pointer, uid *C.char, cRemoteAudioTrack unsafe.Pointer) {
	//validity check
	if cLocalUser == nil {
		return
	}
	
	// get conn from handle
	con := agoraService.getConFromHandle(cLocalUser, ConTypeCLocalUser)
	// date: 20251028 for set apm filter properties
	// if con set to enable 3a, then open apm filter else do nothing! 
	// Whether the OnUserAudioTrackSubscribed method is registered or not, the following operations should be performed!
	if con != nil {
		// open apm filter
		con.setApmFilterProperties(uid, cRemoteAudioTrack)
	} 
	if con == nil || con.localUserObserver == nil || con.localUserObserver.OnUserAudioTrackSubscribed == nil {
		return
	}

	// note： best practise is never reelase handler until app is exiting
	con.localUserObserver.OnUserAudioTrackSubscribed(con.GetLocalUser(), C.GoString(uid), NewRemoteAudioTrack(cRemoteAudioTrack))
}

//export goOnUserVideoTrackSubscribed
func goOnUserVideoTrackSubscribed(cLocalUser unsafe.Pointer, uid *C.char, info *C.struct__video_track_info, cRemoteVideoTrack unsafe.Pointer) {
	//validity check
	if cLocalUser == nil {
		return
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cLocalUser, ConTypeCLocalUser)
	if con == nil || con.localUserObserver == nil || con.localUserObserver.OnUserVideoTrackSubscribed == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.localUserObserver.OnUserVideoTrackSubscribed(con.GetLocalUser(), C.GoString(uid), GoVideoTrackInfo(info), con.NewRemoteVideoTrack(cRemoteVideoTrack))
}

//export goOnUserAudioTrackStateChanged
func goOnUserAudioTrackStateChanged(cLocalUser unsafe.Pointer, uid *C.char, cRemoteAudioTrack unsafe.Pointer, state C.int, reason C.int, elapsed C.int) {
	//validity check
	if cLocalUser == nil {
		return
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cLocalUser, ConTypeCLocalUser)
	if con == nil || con.localUserObserver == nil || con.localUserObserver.OnUserAudioTrackStateChanged == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.localUserObserver.OnUserAudioTrackStateChanged(con.GetLocalUser(), C.GoString(uid), NewRemoteAudioTrack(cRemoteAudioTrack), int(state), int(reason), int(elapsed))
}

//export goOnUserVideoTrackStateChanged
func goOnUserVideoTrackStateChanged(cLocalUser unsafe.Pointer, uid *C.char, cRemoteVideoTrack unsafe.Pointer, state C.int, reason C.int, elapsed C.int) {
	//validity check
	if cLocalUser == nil {
		return
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cLocalUser, ConTypeCLocalUser)
	if con == nil || con.localUserObserver == nil || con.localUserObserver.OnUserVideoTrackStateChanged == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.localUserObserver.OnUserVideoTrackStateChanged(con.GetLocalUser(), C.GoString(uid), con.NewRemoteVideoTrack(cRemoteVideoTrack), int(state), int(reason), int(elapsed))
}


//export goOnAudioVolumeIndication
func goOnAudioVolumeIndication(cLocalUser unsafe.Pointer, Volumes *C.struct__audio_volume_info, speakerNumber C.uint, totalVolume C.int) {
	//validity check
	if cLocalUser == nil {
		return
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cLocalUser, ConTypeCLocalUser)
	if con == nil || con.localUserObserver == nil || con.localUserObserver.OnAudioVolumeIndication == nil {
		return
	}


	// item := C.GoBytes(unsafe.Pointer(Volumes), C.int(unsafe.Sizeof(Volumes))) and assign to a list
	//todo?? move to file end
	frames := make([]*AudioVolumeInfo, int(speakerNumber))
	c_elementSize := C.sizeof_struct__audio_volume_info
	for i := 0; i < int(speakerNumber); i++ {
		c_element := (*C.struct__audio_volume_info)(unsafe.Pointer(uintptr(unsafe.Pointer(Volumes)) + uintptr(c_elementSize)*uintptr(i)))
		//element := (*C._audio_volume_info)(unsafe.Pointer(uintptr(unsafe.Pointer(volumes)) + uintptr(C.sizeof__audio_volume_info)*uintptr(i)))

		volume := GoAudioVolumeInfo(c_element)

		frames[i] = volume

	}

	con.localUserObserver.OnAudioVolumeIndication(con.GetLocalUser(), frames, int(speakerNumber), int(totalVolume))
}


//export goOnAudioPublishStateChanged
func goOnAudioPublishStateChanged(cLocalUser unsafe.Pointer, channel *C.char, oldState C.int, newState C.int, elapseSinceLastState C.int) {
	//fmt.Printf("goOnAudioPublishStateChanged: %d, %d, %d\n", oldState, newState, elapseSinceLastState)
	//validity check
	if cLocalUser == nil {
		return
	}
	
	// get conn from handle
	con := agoraService.getConFromHandle(cLocalUser, ConTypeCLocalUser)
	if con == nil || con.localUserObserver == nil || con.localUserObserver.OnAudioPublishStateChanged == nil {
		return
	}
	con.localUserObserver.OnAudioPublishStateChanged(con.GetLocalUser(), C.GoString(channel), int(oldState), int(newState), int(elapseSinceLastState))
}
//export goOnAudioMetadataReceived
func goOnAudioMetadataReceived(cLocalUser unsafe.Pointer, uid *C.char, metaData *C.char, length C.size_t) {
	//validity check
	if cLocalUser == nil {
		return
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cLocalUser, ConTypeCLocalUser)
	if con == nil || con.localUserObserver == nil || con.localUserObserver.OnAudioMetaDataReceived == nil {
		return
	}
	// note： best practise is never reelase handler until app is exiting
	con.localUserObserver.OnAudioMetaDataReceived(con.GetLocalUser(), C.GoString(uid), C.GoBytes(unsafe.Pointer(metaData),C.int(length)))
}
//export goOnLocalAudioTrackStatistics
func goOnLocalAudioTrackStatistics(cLocalUser unsafe.Pointer, stats *C.struct__local_audio_stats) {
	//validity check
	if cLocalUser == nil {
		return
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cLocalUser, ConTypeCLocalUser)
	if con == nil || con.localUserObserver == nil || con.localUserObserver.OnLocalAudioTrackStatistics == nil {
		return
	}
	con.localUserObserver.OnLocalAudioTrackStatistics(con.GetLocalUser(), GoLocalAudioStats(stats))
}
//export goOnRemoteAudioTrackStatistics
func goOnRemoteAudioTrackStatistics(cLocalUser unsafe.Pointer, uid *C.char, stats *C.struct__remote_audio_stats) {
	//validity check
	if cLocalUser == nil {
		return
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cLocalUser, ConTypeCLocalUser)
	if con == nil || con.localUserObserver == nil || con.localUserObserver.OnRemoteAudioTrackStatistics == nil {
		return
	}
	con.localUserObserver.OnRemoteAudioTrackStatistics(con.GetLocalUser(), C.GoString(uid), GoRemoteAudioStats(stats))
}
//export goOnLocalVideoTrackStatistics
func goOnLocalVideoTrackStatistics(cLocalUser unsafe.Pointer, stats *C.struct__local_video_track_stats) {
	//validity check
	if cLocalUser == nil {
		return
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cLocalUser, ConTypeCLocalUser)
	if con == nil || con.localUserObserver == nil || con.localUserObserver.OnLocalVideoTrackStatistics == nil {
		return
	}
	con.localUserObserver.OnLocalVideoTrackStatistics(con.GetLocalUser(), GoLocalVideoStats(stats))
}
//export goOnRemoteVideoTrackStatistics
func goOnRemoteVideoTrackStatistics(cLocalUser unsafe.Pointer, uid *C.char, stats *C.struct__remote_video_track_stats) {
	//validity check
	if cLocalUser == nil {
		return
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cLocalUser, ConTypeCLocalUser)
	if con == nil || con.localUserObserver == nil || con.localUserObserver.OnRemoteVideoTrackStatistics == nil {
		return
	}
	con.localUserObserver.OnRemoteVideoTrackStatistics(con.GetLocalUser(), C.GoString(uid), GoRemoteVideoStats(stats))
}
//export goOnEncryptionError
func goOnEncryptionError(cCon unsafe.Pointer, errorType C.int) {
	//validity check
	if cCon == nil {
		return
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cCon, ConTypeCCon)
	if con == nil || con.handler == nil || con.handler.OnEncryptionError == nil {
		return
	}
	con.handler.OnEncryptionError(con, int(errorType))
}

//export goOnAudioTrackPublishSuccess
func goOnAudioTrackPublishSuccess(cLocalUser unsafe.Pointer, cLocalAudioTrack unsafe.Pointer) {
	//validity check
	fmt.Printf("goOnAudioTrackPublishSuccess: %v\n", cLocalUser)
	if cLocalUser == nil {
		return
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cLocalUser, ConTypeCLocalUser)
	if con == nil || con.localUserObserver == nil || con.localUserObserver.OnAudioTrackPublishSuccess == nil {
		return
	}
	con.localUserObserver.OnAudioTrackPublishSuccess(con.GetLocalUser(), nil)
}
//export goOnAudioTrackUnpublished
func goOnAudioTrackUnpublished(cLocalUser unsafe.Pointer, cLocalAudioTrack unsafe.Pointer) {
	fmt.Printf("goOnAudioTrackPublishSuccess: %v\n", cLocalUser)
	//validity check
	if cLocalUser == nil {
		return
	}
	// get conn from handle
	con := agoraService.getConFromHandle(cLocalUser, ConTypeCLocalUser)
	if con == nil || con.localUserObserver == nil || con.localUserObserver.OnAudioTrackUnpublished == nil {
		return
	}
	con.localUserObserver.OnAudioTrackUnpublished(con.GetLocalUser(), nil)
}

//export goOnCapabilitiesChanged
func goOnCapabilitiesChanged(cCapObserverHandle unsafe.Pointer, caps *C.struct__capabilities, size C.int) {
	//validity check
	//fmt.Printf("goOnCapabilitiesChanged, size: %d, cCapObserverHandle: %v\n", size, cCapObserverHandle)
	if cCapObserverHandle == nil {
		return
	}
	var con *RtcConnection = nil
	// get con from  handle
	var found bool = false
	agoraService.consByCCon.Range(func(key, value interface{}) bool {
		con = value.(*RtcConnection)
		//fmt.Printf("goOnCapabilitiesChanged, con: %v, cCapObserverHandle: %v\n", cCapObserverHandle, con.cCapObserverHandle)
		if con.cCapObserverHandle == cCapObserverHandle {
			found = true
			return false
		}
		return true
	})

	//fmt.Printf("goOnCapabilitiesChanged, found: %v, con: %v\n", found, con)
	if !found {
		con = nil
	}
	
	//fmt.Printf("goOnCapabilitiesChanged, size: %d, con: %v\n", size, con)
	if con == nil  {
		return
	}
	con.handleCapabilitiesChanged(caps, size)
}