package agoraservice

// #cgo darwin CFLAGS: -I../../agora_sdk/include/c/api2 -I../../agora_sdk/include/c/base
// #cgo darwin LDFLAGS: -Wl,-rpath,${SRCDIR}/../../agora_sdk_mac -L../../agora_sdk_mac -lAgoraRtcKit -lAgorafdkaac -lAgoraffmpeg
// #cgo linux CFLAGS: -D__linux__=1 -I../../agora_sdk/include/c/api2 -I../../agora_sdk/include/c/base
// #cgo linux LDFLAGS: -Wl,-rpath,${SRCDIR}/../../agora_sdk -L../../agora_sdk/ -lagora_rtc_sdk -lagora-fdkaac -laosl
// #include "agora_local_user.h"
// #include "agora_rtc_conn.h"
// #include "agora_service.h"
// #include "agora_media_base.h"
// #include "agora_parameter.h"
//
// #ifndef agora_service_load_extension_provider
// AGORA_API_C_INT agora_service_load_extension_provider(AGORA_HANDLE agora_svc, const char* path, unsigned int unload_after_use) {
// 	 return -1;
// }
// #endif
import "C"
import (
	"sync"
	"unsafe"
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
	AreaCode AreaCode
	// ChannelProfile is the channel profile.
	ChannelProfile ChannelProfile
	// AudioScenario is the audio scenario.
	AudioScenario AudioScenario
	// UseStringUid determines whether to enable string uid.
	UseStringUid bool

	// LogPath is the rtc log path.
	// - Default path on linux is in ~/.agora/.
	// - Default path on mac is in ～/Library/Logs/agorasdk.log, if sandbox was turned off;
	//   or in /Users/<username>/Library/Containers/<AppBundleIdentifier>/Data/Library/Logs/agorasdk.log, if sandbox was turned on.
	LogPath string
	// LogSize is the rtc log file size in Byte.
	LogSize int
	// LogLevel is the rtc log level.
	LogLevel int
	// version 2.1.6
	DomainLimit int
	// version 2.2.1
	// if >0, when remote user muted itself, the onplaybackbeforemixing will be still called badk with active pacakage
	// if <=0, when remote user muted itself, the onplaybackbeforemixing will be no longer called back
	// default to 0, i.e when muted, no callback will be triggered
	ShouldCallbackWhenMuted int
	// from version 2.3.0, the name is misleading, it is not about stero encode mode, it is about audio label generator
	EnableSteroEncodeMode int
	// version 2.2.9 and later, if not set, use default path	
	ConfigDir string
	// version 2.2.9 and later, if not set, use default path
	DataDir string
}

// const def for map type
const (
	ConTypeCCon                  = 0
	ConTypeCLocalUser            = 1
	ConTypeCVideoObserver        = 2
	ConTypeCEncodedVideoObserver = 3
)

// AgoraService is the core service of the Agora SDK.
// It manages the lifecycle of the Agora SDK and provides methods to interact with the SDK.
// It is a singleton, so you should not create multiple instances of it.
type AgoraService struct {
	inited  bool
	service unsafe.Pointer
	
	//isSteroEncodeMode bool
	//audioScenario AudioScenario
	// mediaFactory         unsafe.Pointer
	consByCCon                  sync.Map
	consByCLocalUser            sync.Map
	consByCVideoObserver        sync.Map
	consByCEncodedVideoObserver sync.Map
	mediaFactory *MediaNodeFactory
}

// / newAgoraService creates a new instance of AgoraService
func newAgoraService() *AgoraService {
	return &AgoraService{
		inited:  false,
		service: nil,
		//isSteroEncodeMode: false,
		//audioScenario: AudioScenarioChorus,
		mediaFactory: nil,
	}
}

var agoraService *AgoraService = newAgoraService()

func NewAgoraServiceConfig() *AgoraServiceConfig {
	return &AgoraServiceConfig{
		EnableAudioProcessor: true,
		EnableAudioDevice:    false,
		EnableVideo:          false,
		AppId:                "",
		AreaCode:             AreaCodeGlob,
		ChannelProfile:       ChannelProfileLiveBroadcasting,
		// for AI Scenario:
		// default is AudioScenarioAiServer, if want to use other scenario, pls contact us and make sure the scenario is more optimized for your business
		// for other business, you can use AudioScenarioChorus, AudioScenarioDefault, etc. but recommend to contact us and make sure the scenario is more optimized for your business
		AudioScenario:        AudioScenarioAiServer,
		UseStringUid:         false,
		LogPath:              "",  // format like: "./agora_rtc_log/agorasdk.log"
		LogSize:              5 * 1024, // default to: 5MB
		LogLevel: 0,
		DomainLimit: 0, // default to 0
		ShouldCallbackWhenMuted: 0, // default to 0, no callback when muted
		EnableSteroEncodeMode: 0, // default to 0,i.e default to mono encode mode
		ConfigDir: "",   // format like: "./agora_rtc_log"
		DataDir: "",     // format like: "./agora_rtc_log", should ensure the directory exists
	}
}

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

	// agoraService.mediaFactory = C.agora_service_create_media_node_factory(agoraService.service)
	cParam := C.agora_service_get_agora_parameter(agoraService.service)
	cParamStr := C.CString("rtc.set_app_type")
	defer C.free(unsafe.Pointer(cParamStr))
	C.agora_parameter_set_int(cParam, cParamStr, C.int(17))

	agoraParam := GetAgoraParameter()
	
	if (cfg.EnableSteroEncodeMode < 1) { // disable stereo encode mode,can enabel audio label
		// enable audio label generator
		EnableExtension("agora.builtin", "agora_audio_label_generator", "", true)

		// enable vad v2 model
		agoraParam.SetParameters("{\"che.audio.label.enable\": true}")
	}

	// from version 2.2.1
	if cfg.ShouldCallbackWhenMuted > 0 {
		agoraParam.SetParameters("{\"rtc.audio.enable_user_silence_packet\": true}")
	}
	// date: 2025-09-09 
	// to disable av1 resolution limitation: for any resolution, 
	// it will be encoded as av1 if config is av1 or it only work for resolution >= 360p
	agoraParam.SetParameters("{\"che.video.min_enc_level\": 0}")

	agoraService.mediaFactory = newMediaNodeFactory()

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
	// cleanup go layer resources
	agoraService.cleanup()

	if agoraService.mediaFactory != nil {
		agoraService.mediaFactory.release()
		agoraService.mediaFactory = nil
	}

	// and release c layer resources
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

// LoadExtensionProvider loads an extension provider.
func LoadExtensionProvider(path string, unloadAfterUse bool) int {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	var cUnloadAfterUse C.uint = 0
	if unloadAfterUse {
		cUnloadAfterUse = 1
	}
	return int(C.agora_service_load_extension_provider(agoraService.service, cPath, cUnloadAfterUse))
}

func EnableExtension(providerName, extensionName, trackId string, autoEnableOnTrack bool) int {
	cProviderName := C.CString(providerName)
	defer C.free(unsafe.Pointer(cProviderName))
	cExtensionName := C.CString(extensionName)
	defer C.free(unsafe.Pointer(cExtensionName))
	cTrackId := C.CString(trackId)
	defer C.free(unsafe.Pointer(cTrackId))
	var cAutoEnableOnTrack C.uint = 0
	if autoEnableOnTrack {
		cAutoEnableOnTrack = 1
	}
	return int(C.agora_service_enable_extension(agoraService.service, cProviderName, cExtensionName, cTrackId, cAutoEnableOnTrack))
}

func DisableExtension(providerName, extensionName, trackId string) int {
	cProviderName := C.CString(providerName)
	defer C.free(unsafe.Pointer(cProviderName))
	cExtensionName := C.CString(extensionName)
	defer C.free(unsafe.Pointer(cExtensionName))
	cTrackId := C.CString(trackId)
	defer C.free(unsafe.Pointer(cTrackId))
	return int(C.agora_service_disable_extension(agoraService.service, cProviderName, cExtensionName, cTrackId))
}

func GetAgoraParameter() *AgoraParameter {
	return &AgoraParameter{
		cParameter: C.agora_service_get_agora_parameter(agoraService.service),
	}
}

// add cleanup method
func (s *AgoraService) cleanup() {
    // Clean all connections and release associated resources
    s.consByCCon.Range(func(key, value interface{}) bool {
        if conn, ok := value.(*RtcConnection); ok && conn != nil {
            // Perform any necessary cleanup for the connection
            
        }
        s.consByCCon.Delete(key)
        return true
    })

    s.consByCLocalUser.Range(func(key, value interface{}) bool {
        if conn, ok := value.(*RtcConnection); ok && conn != nil {
            //conn.Close()  // do not call this method, it should be called by user
        }
        s.consByCLocalUser.Delete(key)
        return true
    })

    s.consByCVideoObserver.Range(func(key, value interface{}) bool {
        if conn, ok := value.(*RtcConnection); ok && conn != nil {
        }
        s.consByCVideoObserver.Delete(key)
        return true
    })

    s.consByCEncodedVideoObserver.Range(func(key, value interface{}) bool {
        if conn, ok := value.(*RtcConnection); ok && conn != nil {
        }
        s.consByCEncodedVideoObserver.Delete(key)
        return true
    })

    // After cleaning up all connections, set maps to nil
    s.consByCCon = sync.Map{}
    s.consByCLocalUser = sync.Map{}
    s.consByCVideoObserver = sync.Map{}
    s.consByCEncodedVideoObserver = sync.Map{}
}

// to get value from sync.Map, use Load method
// to set value to sync.Map, use Store method
// to delete value from sync.Map, use Delete method
func (s *AgoraService) setConFromHandle(handle unsafe.Pointer, con *RtcConnection, conType int) int {
	switch conType {
	case ConTypeCCon:
		s.consByCCon.Store(handle, con)
	case ConTypeCLocalUser:
		s.consByCLocalUser.Store(handle, con)
	case ConTypeCVideoObserver:
		s.consByCVideoObserver.Store(handle, con)
	case ConTypeCEncodedVideoObserver:
		s.consByCEncodedVideoObserver.Store(handle, con)
	default:
		return -1
	}
	return 0
}

func (s *AgoraService) getConFromHandle(handle unsafe.Pointer, conType int) *RtcConnection {
	var value interface{}
	var ok bool

	if handle == nil {
		return nil
	}

	switch conType {
	case ConTypeCCon:
		value, ok = s.consByCCon.Load(handle)
	case ConTypeCLocalUser:
		value, ok = s.consByCLocalUser.Load(handle)
	case ConTypeCVideoObserver:
		value, ok = s.consByCVideoObserver.Load(handle)
	case ConTypeCEncodedVideoObserver:
		value, ok = s.consByCEncodedVideoObserver.Load(handle)
	default:
		return nil
	}
	if !ok {
		return nil
	}
	return value.(*RtcConnection)
}

func (s *AgoraService) deleteConFromHandle(handle unsafe.Pointer, conType int) bool {
	if handle == nil {
		return false
	}
	switch conType {
	case ConTypeCCon:
		s.consByCCon.Delete(handle)
	case ConTypeCLocalUser:
		s.consByCLocalUser.Delete(handle)
	case ConTypeCVideoObserver:
		s.consByCVideoObserver.Delete(handle)
	case ConTypeCEncodedVideoObserver:
		s.consByCEncodedVideoObserver.Delete(handle)
	}
	return true
}
