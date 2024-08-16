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
	// AudioScenarioDefault is the default audio scenario.
	AudioScenarioDefault = 0
	// AudioScenarioGameStreaming is the live gaming scenario, which needs to enable gaming
	// audio effects in the speaker. Choose this scenario to achieve high-fidelity
	// music playback.
	AudioScenarioGameStreaming = 3
	// AudioScenarioChatRoom is the chatroom scenario, which needs to keep recording when setClientRole to audience.
	// Normally, app developer can also use mute api to achieve the same result,
	// and we implement this 'non-orthogonal' behavior only to make API backward compatible.
	AudioScenarioChatRoom = 5
	// AudioScenarioChorus is the chorus scenario.
	AudioScenarioChorus = 7
	// AudioScenarioMeeting is the meeting scenario.
	AudioScenarioMeeting = 8
)

const (
	// AreaCodeCN is Mainland China.
	AreaCodeCN = 0x00000001
	// AreaCodeNA is North America.
	AreaCodeNA = 0x00000002
	// AreaCodeEU is Europe.
	AreaCodeEU = 0x00000004
	// AreaCodeAS is Asia, excluding Mainland China.
	AreaCodeAS = 0x00000008
	// AreaCodeJP is Japan.
	AreaCodeJP = 0x00000010
	// AreaCodeIN is India.
	AreaCodeIN = 0x00000020
	// AreaCodeGlob is (Default) Global.
	AreaCodeGlob = (0xFFFFFFFF)
)

const (
	// Broadcaster. A broadcaster can both send and receive streams.
	ClientRoleBroadcaster = 1
	// Audience. An audience can only receive streams.
	ClientRoleAudience = 2
)

const (
	// Communication.
	// This profile prioritizes smoothness and applies to the one-to-one scenario.
	ChannelProfileCommunication = 0
	// (Default) Live Broadcast.
	// This profile prioritizes supporting a large audience in a live broadcast channel.
	ChannelProfileLiveBroadcasting = 1
)

// AgoraServiceConfig is used to initialize agora service.
type AgoraServiceConfig struct {
	// EnableAudioProcessor determines whether to enable the audio processor.
	// - `true`: (Default) Enable the audio processor. Once enabled, the underlying
	//   audio processor is initialized in advance.
	// - `false`: Disable the audio processor. Set this member
	//   as `false` if you do not need audio at all.
	EnableAudioProcessor bool
	// EnableAudioDevice determines whether to enable the audio device.
	// - `true`: (Default) Enable the audio device. Once enabled, the underlying
	//   audio device module is initialized in advance to support audio
	//   recording and playback.
	// - `false`: Disable the audio device. Set this member as `false` if
	//   you do not need to record or play the audio.
	//
	// When this member is set as `false`, and `enableAudioProcessor` is set as `true`,
	// you can still pull the PCM audio data.
	EnableAudioDevice bool
	// EnableVideo determines whether to enable video.
	// - `true`: Enable video. Once enabled, the underlying video engine is
	//   initialized in advance.
	// - `false`: (Default) Disable video. Set this parameter as `false` if you
	//   do not need video at all.
	EnableVideo bool

	// AppId is the App ID of your project.
	AppId string
	// AreaCode is the supported area code, default is AreaCodeGlob.
	AreaCode uint
	// ChannelProfile is the channel profile.
	ChannelProfile int
	// AudioScenario is the audio scenario.
	AudioScenario int
	// UseStringUid determines whether to enable string uid.
	UseStringUid bool

	// LogPath is the rtc log path.
	// - Default path on linux is in ~/.agora/.
	// - Default path on mac is in ～/Library/Logs/agorasdk.log, if sandbox was turned off;
	//   or in /Users/<username>/Library/Containers/<AppBundleIdentifier>/Data/Library/Logs/agorasdk.log, if sandbox was turned on.
	LogPath string
	// LogSize is the rtc log file size in Byte.
	LogSize int
}

type AgoraService struct {
	inited               bool
	service              unsafe.Pointer
	mediaFactory         unsafe.Pointer
	consByCCon           map[unsafe.Pointer]*RtcConnection
	consByCLocalUser     map[unsafe.Pointer]*RtcConnection
	consByCVideoObserver map[unsafe.Pointer]*RtcConnection
	connectionRWMutex    *sync.RWMutex
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
	}
}

var agoraService *AgoraService = newAgoraService()

// Initialize the Agora service.
// The Agora service is the core service of the Agora SDK.
// You must call this method before calling any other methods.
// This function must be called once globally.
func Initialize(cfg *AgoraServiceConfig) int {
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

// Release the Agora service.
// This function must be called once globally.
// After this function is called, you must not call other agora APIs any more.
func Release() int {
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
