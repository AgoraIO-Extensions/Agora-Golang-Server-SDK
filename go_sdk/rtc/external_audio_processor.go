package agoraservice

/*
#cgo CFLAGS: -I../../agora_sdk/include/c/api2 -I../../agora_sdk/include/c/base
#include <string.h>
#include <stdlib.h>
#include "agora_audio_track.h"
#include "agora_media_node_factory.h"
#include "audio_sink_callback_cgo.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
	"sync/atomic"
)

//export goOnSinkAudioFrame
func goOnSinkAudioFrame(sink unsafe.Pointer, frame unsafe.Pointer) C.int {
	
	goFrame := GoSinkAudioFrame((*C.struct__audio_pcm_frame)(frame))
	// restore external audio processor instance from user data
	processor := (*ExternalAudioProcessor)(unsafe.Pointer(sink))
	if processor == nil {
		fmt.Printf("[ExternalAudioProcessor] failed to restore external audio processor instance from user data\n")
		return -1
	}
	// call to external audio processor instance to process the audio frame
	result := processor.doResultFrame(goFrame)
	return C.int(result)
}
//date: 2025-11-25 to add external audio processor observer
type ExternalAudioProcessorObserver struct {
	OnProcessedAudioFrame func(processor *ExternalAudioProcessor,frame *AudioFrame, vadResultStat VadState, vadResultFrame *AudioFrame)
}

/*
audio sink related:
*/
type AudioSink struct {
	cSinkCallback *C.audio_sink
	cSink         unsafe.Pointer
}

type ExternalAudioProcessor struct {
	pcmSender   *AudioPcmDataSender
	audioTrack  *LocalAudioTrack
	audioSinks  *AudioSink
	initialized bool
	vadInstance *AudioVadV2
	observer *ExternalAudioProcessorObserver
	InStreamInMs atomic.Int64
	ProcessedStreamInMs atomic.Int64
}

func NewExternalAudioProcessor() *ExternalAudioProcessor {
	processor := &ExternalAudioProcessor{
		pcmSender:   nil,
		audioTrack:  nil,
		audioSinks:  nil,
		initialized: false,
		vadInstance: nil,
		observer: nil,
		InStreamInMs: atomic.Int64{},
		ProcessedStreamInMs: atomic.Int64{},
	}

	//and initialize the pcmSender and audioTrack
	factory := agoraService.mediaFactory
	processor.pcmSender = factory.NewAudioPcmDataSender()
	processor.audioTrack = NewCustomAudioTrackPcm(processor.pcmSender, AudioScenarioDefault, false)
	processor.audioSinks = newAudioSink(processor)
	processor.initialized = false
	return processor
}

// Initialize sets up the ExternalAudioProcessor by configuring the underlying audio sink and its filter properties,
// setting audio track parameters, and enabling the audio track for publishing.
// outputSampleRate: the desired output sample rate in Hz.
// outputChannels:   the number of output audio channels.
// Returns 0 on success, or a non-zero error code on failure.
func (p *ExternalAudioProcessor) Initialize(apmConfig *APMConfig,outputSampleRate int, outputChannels int, vadConfig *AudioVadConfigV2, observer *ExternalAudioProcessorObserver) int {
	var ret int = 0

	//check apm model: must support amp but can only do vad
	if isSupportExternalAudioProcessor(agoraService.apmModel)==false {
		fmt.Printf("[ExternalAudioProcessor] apm model is not supported, error code: %d, model: %d\n", ret, agoraService.apmModel)
		return -3000
	}
	if apmConfig == nil && vadConfig == nil {
	    fmt.Printf("[ExternalAudioProcessor] apm config and vad config are both nil, error code: %d\n", ret)
		return -3002
	}
	// check vad config and create vad instance if needed
	if vadConfig != nil {
		p.vadInstance = NewAudioVadV2(vadConfig)
		if p.vadInstance == nil {
			fmt.Printf("[ExternalAudioProcessor] failed to create vad instance, error code: %d\n", ret)
			return -3001
		}
	}
	

	ret = p.setFilterProperties(apmConfig)
	if ret != 0 {
		fmt.Printf("[ExternalAudioProcessor] failed to set filter properties, error code: %d\n", ret)
		return ret
	}

	//1. add sink and set filter properties
	ret = p.addAudioSink(outputSampleRate, outputChannels)
	if ret != 0 {
		fmt.Printf("[ExternalAudioProcessor] failed to add audio sink, error code: %d\n", ret)
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

	// 4. register observer
	p.observer = observer

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
	if readLen%bytesPerFrameInMs != 0 {
		fmt.Printf("PushAudioPcmData data length is not integer multiples of 10ms, readLen: %d, bytesPerFrame: %d\n", readLen, bytesPerFrameInMs)
		return -2001
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
	p.InStreamInMs.Add(int64(packnumInMs))

	return ret
}


func (p *ExternalAudioProcessor) Release() {

	// release track
	// release pcmSender
	// release audioSinks
	if p.audioSinks != nil {
		p.removeAudioSink()
	}
	if p.audioTrack != nil {
		p.audioTrack.SetEnabled(false)
		p.audioTrack.ClearSenderBuffer()
		p.audioTrack.Release()
		p.audioTrack = nil
	}

	if p.audioSinks != nil {
		p.audioSinks.release()
		p.audioSinks = nil
	}

	if p.pcmSender != nil {
		p.pcmSender.Release()
		p.pcmSender = nil
	}

	if p.vadInstance != nil {
		p.vadInstance.Release()
		p.vadInstance = nil
	}

	p.observer = nil
	p.initialized = false

}

//
// private methods
//

func (p *ExternalAudioProcessor) setFilterProperties(apmConfig *APMConfig) int {
	if p.audioTrack == nil || p.audioTrack.cTrack == nil {
		return -2002
	}

	if apmConfig == nil {
		// in this case, no need to set filter properties.
		return 0
	}
	cLocalTrackHandle := p.audioTrack.cTrack

	// Step 1: Enable filter on Track
	ret := enableAudioFilterByTrack(cLocalTrackHandle, "audio_processing_pcm_source", true, true)
	if ret != 0 {
		fmt.Printf("[LocalAudioProcessor] failed to enable filter on Track, error code: %d\n", ret)
		return -2003
	}

	// Step 2: Load AINS resources (if needed)
	ret = setFilterPropertyByTrack(cLocalTrackHandle, "audio_processing_pcm_source", "apm_load_resource", "ains", true)
	if ret != 0 {
		fmt.Printf("[LocalAudioProcessor] failed to load AINS model resources, error code: %d\n", ret)
		return -2004
	}

	// Step 3: Build and apply configuration
	apmConfigJSON := apmConfig.toJson()
	ret = setFilterPropertyByTrack(cLocalTrackHandle, "audio_processing_pcm_source", "apm_config", apmConfigJSON, true)
	if ret != 0 {
		fmt.Printf("[LocalAudioProcessor] failed to configure audio processing parameters, error code: %d\n", ret)
		return -2005
	}
	fmt.Printf("[LocalAudioProcessor] apmConfigJSON: %s\n", apmConfigJSON)

	// Step 4: Enable dump (if debugging is needed)
	if apmConfig.EnableDump {
		ret = setFilterPropertyByTrack(cLocalTrackHandle, "audio_processing_pcm_source", "apm_dump", "true", true)
		if ret != 0 {
			fmt.Printf("[LocalAudioProcessor] failed to enable dump, error code: %d (non-critical)\n", ret)
			return -2006
		}
	}

	return ret
}

func newAudioSink(processor *ExternalAudioProcessor) *AudioSink {
	sink := &AudioSink{
		cSinkCallback: nil,
		cSink:         nil,
	}
	// allocate audio_sink	 structure
	csink_callback := (*C.audio_sink)(C.malloc(C.sizeof_audio_sink))
	//C.memset(unsafe.Pointer(csink_callback), 0, C.sizeof_audio_sink)

	// set callback function pointer - using C wrapper function
	csink_callback.on_audio_frame = (*[0]byte)(C.cgo_onSinkAudioFrameCallback)
	csink_callback.user_data = unsafe.Pointer(processor)

	// create audio sink
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
	if p.audioTrack == nil || p.audioTrack.cTrack == nil || p.audioSinks == nil || p.audioSinks.cSink == nil {
		return -1000
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
		return -1001
	}
	return int(C.agora_audio_track_remove_audio_sink(p.audioTrack.cTrack, p.audioSinks.cSink))
}
func (p *ExternalAudioProcessor) doResultFrame(frame *AudioFrame) int	 {
	// TODO: implement the logic to process the audio frame
	//fmt.Printf("[ExternalAudioProcessor] doResultFrame, voice prob: %d, rms: %d, pitch: %d\n", frame.VoiceProb, frame.Rms, frame.Pitch)
	var vadResultState VadState = VadStateInvalid
	var vadResultFrame *AudioFrame = nil
	if p.vadInstance != nil {
		vadResultFrame, vadResultState = p.vadInstance.Process(frame)
	}
	p.ProcessedStreamInMs.Add(int64(10))
	if p.observer != nil {
		p.observer.OnProcessedAudioFrame(p, frame, vadResultState, vadResultFrame)
	}
	return 0
}
