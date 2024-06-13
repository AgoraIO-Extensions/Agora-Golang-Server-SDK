package agoraservice

/*
#cgo CFLAGS: -I../../agora_sdk/include_c/api2 -I../../agora_sdk/include_c/base

#include "agora_local_user.h"
#include "agora_rtc_conn.h"
#include "agora_service.h"
#include "agora_media_base.h"
*/
import "C"
import (
	"unsafe"
)

//export goOnVideoFrame
func goOnVideoFrame(cObserver unsafe.Pointer, channelId *C.char, uid *C.char, frame *C.struct__video_frame) {
	agoraService.connectionRWMutex.RLock()
	con := agoraService.consByCVideoObserver[cObserver]
	agoraService.connectionRWMutex.RUnlock()
	if con == nil || con.videoObserver == nil || con.videoObserver.OnFrame == nil {
		return
	}
	goChannelId := C.GoString(channelId)
	goUid := C.GoString(uid)
	goFrame := GoVideoFrame(frame)
	con.videoObserver.OnFrame(con, goChannelId, goUid, goFrame)
}
