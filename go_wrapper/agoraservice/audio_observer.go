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

//export goOnPlaybackAudioFrameBeforeMixing
func goOnPlaybackAudioFrameBeforeMixing(cLocalUser unsafe.Pointer, channelId *C.char, uid *C.char, frame *C.struct__audio_frame) {
	agoraService.connectionRWMutex.RLock()
	con := agoraService.consByCLocalUser[cLocalUser]
	agoraService.connectionRWMutex.RUnlock()
	if con == nil || con.audioObserver == nil || con.audioObserver.OnPlaybackAudioFrameBeforeMixing == nil {
		return
	}
	goChannelId := C.GoString(channelId)
	goUid := C.GoString(uid)
	goFrame := GoPcmAudioFrame(frame)
	con.audioObserver.OnPlaybackAudioFrameBeforeMixing(con, goChannelId, goUid, goFrame)
}
