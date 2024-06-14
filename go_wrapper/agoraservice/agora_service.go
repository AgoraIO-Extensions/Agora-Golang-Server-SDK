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
import (
	"sync"
	"unsafe"
)

type AgoraServiceConfig struct {
	AppId string
}

type AgoraService struct {
	inited            bool
	service           unsafe.Pointer
	mediaFactory      unsafe.Pointer
	consByCCon        map[unsafe.Pointer]*RtcConnection
	consByCLocalUser  map[unsafe.Pointer]*RtcConnection
	connectionRWMutex *sync.RWMutex
}

func newAgoraService() *AgoraService {
	return &AgoraService{
		inited:            false,
		service:           nil,
		mediaFactory:      nil,
		consByCCon:        make(map[unsafe.Pointer]*RtcConnection),
		consByCLocalUser:  make(map[unsafe.Pointer]*RtcConnection),
		connectionRWMutex: &sync.RWMutex{},
	}
}

var agoraService *AgoraService = newAgoraService()

func Init(cfg *AgoraServiceConfig) int {
	if agoraService.inited {
		return 0
	}
	if agoraService.service == nil {
		agoraService.service = C.agora_service_create()
		if agoraService.service == nil {
			return -1
		}
	}

	ccfg := CAgoraServiceConfig(cfg)
	defer FreeCAgoraServiceConfig(ccfg)

	ret := int(C.agora_service_initialize(agoraService.service, ccfg))
	if ret != 0 {
		return ret
	}

	agoraService.mediaFactory = C.agora_service_create_media_node_factory(agoraService.service)

	logPath := C.CString("./io.agora.rtc_sdk/agorasdk.log")
	defer C.free(unsafe.Pointer(logPath))
	C.agora_service_set_log_file(agoraService.service, logPath,
		C.uint(512*1024))
	agoraService.inited = true
	return 0
}

func Destroy() int {
	if !agoraService.inited {
		return 0
	}
	if agoraService.service != nil {
		ret := int(C.agora_service_release(agoraService.service))
		if ret != 0 {
			return ret
		}
		agoraService.service = nil
	}
	agoraService.inited = false
	return 0
}
