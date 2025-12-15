package agoraservice

// #cgo darwin CFLAGS: -I../../agora_sdk/include/c/api2 -I../../agora_sdk/include/c/base
// #cgo darwin LDFLAGS: -Wl,-rpath,${SRCDIR}/../../agora_sdk_mac -L../../agora_sdk_mac -lAgoraRtcKit -lAgorafdkaac -lAgoraffmpeg -lAgoraAiNoiseSuppressionExtension
// #cgo linux CFLAGS: -D__linux__=1 -I../../agora_sdk/include/c/api2 -I../../agora_sdk/include/c/base
// #cgo linux LDFLAGS: -Wl,-rpath,${SRCDIR}/../../agora_sdk -L../../agora_sdk/ -lagora_rtc_sdk -lagora-fdkaac -laosl
// #include "agora_local_user.h"
// #include "agora_rtc_conn.h"
// #include "agora_service.h"
// #include "agora_media_base.h"
// #include "agora_parameter.h"
// #include "agora_audio_track.h"
//
// #ifndef __linux__
// // agora_service_load_extension_provider is only available on Linux
// // Provide a stub implementation for macOS that returns ERR_NOT_SUPPORTED
// static inline int agora_service_load_extension_provider(void* service, const char* path, unsigned int unload_after_use) {
//     return -7; // ERR_NOT_SUPPORTED
// }
// #endif
import "C"
import (
	"fmt"
	"sync"
	"unsafe"
	"runtime"
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

	// 20251028 for apm filter related config
	// apm model: 0: disable apm, default to disable, 1: enable apm. def to int not bool for future extension
	APMModel int 
	APMConfig *APMConfig

	// date: 2025-11-03, idle mode, if true, the connection will be released when idle for a period of time
	IdleMode bool
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
	apmConfig *APMConfig
	//timer related
	timer	*PrecisionTimer
	idleQueue []*IdleItem
	idleQueueMutex sync.Mutex
	idleMode bool
	apmModel int
}

// / newAgoraService creates a new instance of AgoraService
func newAgoraService() *AgoraService {
	return &AgoraService{
		inited:  false,
		service: nil,
		//isSteroEncodeMode: false,
		//audioScenario: AudioScenarioChorus,
		mediaFactory: nil,
		apmConfig: nil,
		timer: nil,
		idleQueue: make([]*IdleItem, 0),
		idleQueueMutex: sync.Mutex{},
		idleMode: false,
		apmModel: 0,
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
		APMModel: 0,
		APMConfig: nil,
		IdleMode: true, // default to true for  idle mode
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

		//should load extension for linux only
		if runtime.GOOS != "darwin" {
			ret = LoadExtensionProvider("libagora_ai_noise_suppression_extension.so", true)
			if ret != 0 {
				fmt.Printf("load ains.so failed, ret: %d\n", ret)
			}
		}

		// enable apm filter but disable 3a by default
		EnableExtension("agora.builtin", "audio_processing_remote_playback", "", true)
		if isSupportExternalAudioProcessor(cfg.APMModel) {
			EnableExtension("agora.builtin", "audio_processing_pcm_source", "", true)
		}
	}

	if isEnableAPM(cfg.APMModel) {
		if cfg.APMConfig == nil {
			agoraService.apmConfig = NewAPMConfig()
		}else {
			agoraService.apmConfig = cfg.APMConfig
		}
	}
	agoraService.apmModel = cfg.APMModel
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
	// and start the timer
	if cfg.IdleMode {
		agoraService.idleMode = true
		agoraService.timer = NewPrecisionTimer(50) // 50ms
		agoraService.timer.Start(timerTask)
	}	
	return 0
}
func timerTask() {
	// check the idle queue
	//fmt.Printf("timerTask: idle queue size: %d\n", 1)
	removeIdleItem()
}
func addIdleItem(handle unsafe.Pointer, lifeCycleInMs int) {
	idleItem := newIdleItem(handle, lifeCycleInMs/50) // convert to 50ms interval
	agoraService.idleQueueMutex.Lock()
	agoraService.idleQueue = append(agoraService.idleQueue, idleItem)
	agoraService.idleQueueMutex.Unlock()
}
func removeIdleItem() {
	agoraService.idleQueueMutex.Lock()
	defer agoraService.idleQueueMutex.Unlock()

	//check size
	if len(agoraService.idleQueue) == 0 {
		return
	}

	for _, item := range agoraService.idleQueue {
		item.LifeCycle -=1
		}
	//check the first one, if the LifeCycle is 0, remove it
	item := agoraService.idleQueue[0]
	if item.LifeCycle <= 0 {
		agoraService.idleQueue = agoraService.idleQueue[1:]
		// and release it
		C.agora_rtc_conn_destroy(item.Handle)
		fmt.Printf("remove idle item: %p\n", item.Handle)
		// reset 
		item = nil
	}
	
}
func releaseAllIdleItems() {
	agoraService.idleQueueMutex.Lock()
	defer agoraService.idleQueueMutex.Unlock()

	for _, item := range agoraService.idleQueue {
		C.agora_rtc_conn_destroy(item.Handle)
		fmt.Printf("release all idle item: %p\n", item.Handle)
		item = nil
	}
	agoraService.idleQueue = make([]*IdleItem, 0)
	
}
// Release the Agora service.
// This function must be called once globally.
// After this function is called, you must not call other agora APIs any more.
func Release() int {
	if !agoraService.inited {
		return 0
	}
	// stop timer
	if agoraService.timer != nil {
		agoraService.timer.Stop()
		agoraService.timer = nil

		releaseAllIdleItems()
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
//date: 20251028 for set apm filter related struct:
//AiNSConfig: for AiNS , ns,and sf_st_cfg,sf_ext_cfg
type AiNsConfig struct {
	NsEnabled bool
	AiNSEnabled bool
	AiNSModelPref int
	NsngAlgRoute int
	NsngPredefAgg int
}
type AiAecConfig struct {
	Enabled bool
	SplitSrateFor48k int
}
type BghvsCConfig struct {
	Enabled bool
	VadThr float32
}
type AgcConfig struct {
	Enabled bool
}

type APMConfig struct {
	AiNsConfig *AiNsConfig
	AiAecConfig *AiAecConfig
	BghvsCConfig *BghvsCConfig
	AgcConfig *AgcConfig
	EnableDump bool
}
/*
char apm_config[] = "{\"aec\":{\"split_srate_for_48k\":16000},"\
        "\"bghvs\":{\"enabled\":true, \"vadThr\":0.8},"\
        "\"agc\":{\"enabled\":true}, \"ans\":{\"enabled\":true},"\
        "\"sf_st_cfg\":{\"enabled\":true,\"ainsModelPref\":10},\"sf_ext_cfg\":{\"nsngAlgRoute\":12,\"nsngPredefAgg\":11}}";
*/
func NewAPMConfig() *APMConfig {
	return &APMConfig{
		AiNsConfig: &AiNsConfig{
			AiNSEnabled: true,
			NsEnabled: true,
			AiNSModelPref: 10,
			NsngAlgRoute: 12,
			NsngPredefAgg: 11,
		},
		AiAecConfig: &AiAecConfig{
			Enabled: false,
			SplitSrateFor48k: 16000,
		},
		BghvsCConfig: &BghvsCConfig{
			Enabled: true,
			VadThr: 0.8,
		},
		AgcConfig: &AgcConfig{
			Enabled: false,
		},
		EnableDump: false,
	}
}

func (cfg *APMConfig) toJson() string {
	jsonConfigure := fmt.Sprintf("{\"aec\":{\"enabled\":%t,\"split_srate_for_48k\":%d}," +
		"\"bghvs\":{\"enabled\":%t, \"vadThr\":%f}," +
		"\"agc\":{\"enabled\":%t}, \"ans\":{\"enabled\":%t}," +
		"\"sf_st_cfg\":{\"enabled\":%t,\"ainsModelPref\":%d},\"sf_ext_cfg\":{\"nsngAlgRoute\":%d,\"nsngPredefAgg\":%d}}",
		cfg.AiAecConfig.Enabled, cfg.AiAecConfig.SplitSrateFor48k, cfg.BghvsCConfig.Enabled, cfg.BghvsCConfig.VadThr,
		 cfg.AgcConfig.Enabled, cfg.AiNsConfig.NsEnabled, cfg.AiNsConfig.AiNSEnabled,
		 cfg.AiNsConfig.AiNSModelPref, cfg.AiNsConfig.NsngAlgRoute, cfg.AiNsConfig.NsngPredefAgg)

	return jsonConfigure
}



//note: for remote audio track,no need to enable it! it is enabled by default in service creation
func getAudioFilterPosition(isLocalTrack bool) int {
	if isLocalTrack {
		return int(3)
	}
	return int(2)
}
func  enableAudioFilterByTrack(track unsafe.Pointer, name string, enable bool, isLocalTrack bool) int {
	if track  == nil {
		return -1000
	}

	if !isLocalTrack { // for remote track, no need to enable filter
		return 0
	}
	
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cEnable := C.int(0)
	if enable {
		cEnable = C.int(1)
	}
	Position := getAudioFilterPosition(isLocalTrack)

	ret := int(C.agora_audio_track_enable_audio_filter(track, cName, cEnable, C.int(Position)))
	return ret
}

func setFilterPropertyByTrack(track unsafe.Pointer, name string, key string, value string, isLocalTrack bool) int {
	if track == nil  {
		return -1000
	}
	

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))
	Position := getAudioFilterPosition(isLocalTrack)
	ret := int(C.agora_audio_track_set_filter_property(track, cName, cKey, cValue, C.int(Position)))
	return ret
}
func isSupportExternalAudioProcessor(model int) bool {
	if model < 1 {
		return false
	}
	return true
}
func isEnableAPM(model int) bool {
	if model < 1 {
		return false
	}
	return true
}

