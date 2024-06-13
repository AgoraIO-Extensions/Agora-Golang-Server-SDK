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
	if con == nil {
		return
	}
	con.handler.OnConnected(con, GoRtcConnectionInfo(cConInfo), int(reason))
}
