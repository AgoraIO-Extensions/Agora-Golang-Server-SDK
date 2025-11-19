package agoraservice

/*
#cgo CFLAGS: -I../../agora_sdk/include/c/api2 -I../../agora_sdk/include/c/base
#include <string.h>
#include <stdlib.h>
#include "agora_audio_track.h"
#include "agora_media_node_factory.h"

// Go callback 声明
extern int goOnSinkAudioFrame(void* sink, void* frame);

// C 包装函数 - 直接内联，没有额外开销
static inline int onSinkAudioFrameCallback(void* sink, const audio_pcm_frame* frame) {
    return goOnSinkAudioFrame(sink, (void*)frame);
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

//export goOnSinkAudioFrame
func goOnSinkAudioFrame(sink unsafe.Pointer, frame unsafe.Pointer) C.int {
	// TODO: 实现你的音频帧处理逻辑
	pcmFrame := (*C.audio_pcm_frame)(frame)
	fmt.Printf("[ExternalAudioProcessor] goOnAudioFrame, pcmFrame: %+v\n", pcmFrame)
	return C.int(1) // 返回 1 表示成功
}


/*
audio sink related:
*/
type AudioSink struct {
	cSinkCallback *C.audio_sink
	cSink         unsafe.Pointer
}

type ExternalAudioProcessor struct {
	pcmSender  *AudioPcmDataSender
	audioTrack *LocalAudioTrack
	audioSinks *AudioSink
	initialized bool
}

func NewExternalAudioProcessor() *ExternalAudioProcessor {
	processor := &ExternalAudioProcessor{
		pcmSender:  nil,
		audioTrack: nil,
		audioSinks: nil,
		initialized: false,
	}

	//and initialize the pcmSender and audioTrack
	factory := agoraService.mediaFactory
	processor.pcmSender = factory.NewAudioPcmDataSender()
	processor.audioTrack = NewCustomAudioTrackPcm(processor.pcmSender, AudioScenarioDefault)
	processor.audioSinks = newAudioSink()
	processor.initialized = false
	return processor
}

func (p *ExternalAudioProcessor) Initialize(sampleRate int, channels int) int {
	var ret int = 0
	//1. add sink and set filter properties
	ret = p.addAudioSink(sampleRate, channels)
	if ret != 0 {
		fmt.Printf("[ExternalAudioProcessor] failed to add audio sink, error code: %d\n", ret)
		return ret
	}
	ret = p.setFilterProperties()
	if ret != 0 {
		fmt.Printf("[ExternalAudioProcessor] failed to set filter properties, error code: %d\n", ret)
		return ret
	}

	//3. set track properties
	ret = p.audioTrack.SetSendDelayMs(10)
	if ret != 0 {
		fmt.Printf("[ExternalAudioProcessor] failed to set send delay ms, error code: %d\n", ret)
		return ret
	}
	
	p.audioTrack.SetMaxBufferedAudioFrameNumber(100000) // up to 300 frames, 3000ms
	
	
	// 4 enable & publish audio track
	p.audioTrack.SetEnabled(true)

	// 5 check if initialized successfully
	p.initialized = true
	return 0
}

func (p *ExternalAudioProcessor) PushAudioPcmData(data []byte, sampleRate int, channels int, startPtsInMs int64) int {
	if p == nil || p.audioTrack == nil || p.pcmSender == nil {
		return -2000
	}
	readLen := len(data)
	bytesPerFrameInMs := (sampleRate / 1000) * 2 * channels // 1ms , channels and 16bit
	// validity check: only accepts data with lengths that are integer multiples of 10ms​​ 
	if readLen % bytesPerFrameInMs != 0 {
		fmt.Printf("PushAudioPcmData data length is not integer multiples of 10ms, readLen: %d, bytesPerFrame: %d\n", readLen, bytesPerFrameInMs)
		return -2
	}
	packnumInMs := readLen / bytesPerFrameInMs
	
	
	frame := &AudioFrame{
				Buffer:            nil,
				RenderTimeMs:      0,
				PresentTimeMs:     startPtsInMs,
				SamplesPerChannel: sampleRate / 100,
				BytesPerSample:    2,
				Channels:          channels,
				SamplesPerSec:     sampleRate,
				Type:              AudioFrameTypePCM16,
			}

	frame.Buffer = data
	frame.SamplesPerChannel = (sampleRate / 1000) * packnumInMs
	

	ret := p.pcmSender.SendAudioPcmData(frame)

	return ret
}	

func (p *ExternalAudioProcessor) Release() {
	
	// release track
	// release pcmSender
	// release audioSinks
	if p.audioSinks != nil {
		p.removeAudioSink()
		p.audioSinks.release()
		p.audioSinks = nil
	}
	if p.audioTrack != nil {
		p.audioTrack.Release()
		p.audioTrack = nil
	}
	if p.pcmSender != nil {
		p.pcmSender.Release()
		p.pcmSender = nil
	}
	
}

//
// private methods
//

func (p *ExternalAudioProcessor) setFilterProperties() int {
	if p.audioTrack == nil || p.audioTrack.cTrack == nil {
		return -1
	}
	cLocalTrackHandle := p.audioTrack.cTrack

	// Step 1: Enable filter on Track
	ret := enableAudioFilterByTrack(cLocalTrackHandle, "audio_processing_pcm_source", true, true)
	if ret != 0 {
		fmt.Printf("[LocalAudioProcessor] failed to enable filter on Track, error code: %d\n", ret)
	}

	// Step 2: Load AINS resources (if needed)
	ret = setFilterPropertyByTrack(cLocalTrackHandle, "audio_processing_pcm_source", "apm_load_resource", "ains", true)
	if ret != 0 {
		fmt.Printf("[LocalAudioProcessor] failed to load AINS model resources, error code: %d\n", ret)
	}

	// Step 3: Build and apply configuration
	apmConfigJSON := agoraService.apmConfig.toJson()
	ret = setFilterPropertyByTrack(cLocalTrackHandle, "audio_processing_pcm_source", "apm_config", apmConfigJSON, true)
	if ret != 0 {
		fmt.Printf("[LocalAudioProcessor] failed to configure audio processing parameters, error code: %d\n", ret)
	}

	// Step 4: Enable dump (if debugging is needed)
	if agoraService.apmConfig.EnableDump {
		ret = setFilterPropertyByTrack(cLocalTrackHandle, "audio_processing_pcm_source", "apm_dump", "true", true)
		if ret != 0 {
			fmt.Printf("[LocalAudioProcessor] failed to enable dump, error code: %d (non-critical)\n", ret)
		}
	}

	return ret
}





func newAudioSink() *AudioSink {
	sink := &AudioSink{
		cSinkCallback: nil,
		cSink:         nil,
	}
	// 分配 sink 回调结构体
	csink_callback := (*C.audio_sink)(C.malloc(C.sizeof_audio_sink))
	C.memset(unsafe.Pointer(csink_callback), 0, C.sizeof_audio_sink)
	// 设置回调函数指针 - 使用 C 包装函数
	csink_callback.on_audio_frame = (*[0]byte)(C.onSinkAudioFrameCallback)

	// 创建 audio sink
	sink.cSink = C.agora_audio_sink_create(csink_callback)
	sink.cSinkCallback = csink_callback
	return sink
}
func (s *AudioSink) release() {
	if s.cSinkCallback != nil {
		C.free(unsafe.Pointer(s.cSinkCallback))
		s.cSinkCallback = nil
	}
	if s.cSink != nil {
		C.agora_audio_sink_destroy(s.cSink)
		s.cSink = nil
	}
}

func (p *ExternalAudioProcessor) addAudioSink(sampleRate int, channels int) int {
	if p.audioTrack == nil || p.audioTrack.cTrack == nil || p.audioSinks == nil {
		return -1
	}

	wants := &C.struct__audio_sink_wants{
		samples_per_sec: C.int(sampleRate),
		channels:        C.uint32_t(channels),
	}

	return int(C.agora_audio_track_add_audio_sink(p.audioTrack.cTrack, p.audioSinks.cSink, wants))
}

// RemoveAudioSink 从 audio track 移除音频接收器
func (p *ExternalAudioProcessor) removeAudioSink() int {
	if p.audioTrack == nil || p.audioTrack.cTrack == nil || p.audioSinks == nil {
		return -1
	}
	return int(C.agora_audio_track_remove_audio_sink(p.audioTrack.cTrack, p.audioSinks.cSink))
}