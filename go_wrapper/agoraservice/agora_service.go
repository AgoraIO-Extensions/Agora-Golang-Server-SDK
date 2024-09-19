package agoraservice

// #cgo CFLAGS: -I${SRCDIR}/../../agora_sdk/include_c/api2 -I${SRCDIR}/../../agora_sdk/include_c/base
// #cgo darwin LDFLAGS: -Wl,-rpath,../../agora_sdk_mac -L../../agora_sdk_mac -lAgoraRtcKit -lAgorafdkaac -lAgoraffmpeg
// #cgo linux LDFLAGS: -L../../agora_sdk/ -lagora_rtc_sdk -lagora-fdkaac -lagora-core
// #include "agora_local_user.h"
// #include "agora_rtc_conn.h"
// #include "agora_service.h"
// #include "agora_media_base.h"
import "C"
import (
	"sync"
	"unsafe"
)

const (
	/**
	 * 0: (Recommended) The default audio scenario.
	 */
	AUDIO_SCENARIO_DEFAULT = 0
	/**
	 * 3: (Recommended) The live gaming scenario, which needs to enable gaming
	 * audio effects in the speaker. Choose this scenario to achieve high-fidelity
	 * music playback.
	 */
	AUDIO_SCENARIO_GAME_STREAMING = 3
	/**
	 * 5: The chatroom scenario, which needs to keep recording when setClientRole to audience.
	 * Normally, app developer can also use mute api to achieve the same result,
	 * and we implement this 'non-orthogonal' behavior only to make API backward compatible.
	 */
	AUDIO_SCENARIO_CHATROOM = 5
	/**
	 * 7: Chorus
	 */
	AUDIO_SCENARIO_CHORUS = 7
	/**
	 * 8: Meeting
	 */
	AUDIO_SCENARIO_MEETING = 8
	/**
	 * 9: Reserved.
	 */
	AUDIO_SCENARIO_NUM = 9
)

type AgoraServiceConfig struct {
	AppId         string
	AudioScenario int
	LogPath       string
	LogSize       int
	AreaCode      uint
}

type AgoraService struct {
	inited               bool
	service              unsafe.Pointer
	mediaFactory         unsafe.Pointer
	consByCCon           map[unsafe.Pointer]*RtcConnection
	consByCLocalUser     map[unsafe.Pointer]*RtcConnection
	consByCVideoObserver map[unsafe.Pointer]*RtcConnection
	connectionRWMutex    *sync.RWMutex
	cFuncMutex           *sync.Mutex
}

func newAgoraService() *AgoraService {
	return &AgoraService{
		inited:               false,
		service:              nil,
		mediaFactory:         nil,
		consByCCon:           make(map[unsafe.Pointer]*RtcConnection),
		consByCLocalUser:     make(map[unsafe.Pointer]*RtcConnection),
		consByCVideoObserver: make(map[unsafe.Pointer]*RtcConnection),
		connectionRWMutex:    &sync.RWMutex{},
		cFuncMutex:           &sync.Mutex{},
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

	if cfg.LogPath != "" {
		logPath := C.CString(cfg.LogPath)
		defer C.free(unsafe.Pointer(logPath))
		logSize := 512 * 1024
		if cfg.LogSize > 0 {
			logSize = cfg.LogSize
		}
		C.agora_service_set_log_file(agoraService.service, logPath,
			C.uint(logSize))
	}
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
