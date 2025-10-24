package agoraservice

// #cgo CFLAGS: -I${SRCDIR}/../../agora_sdk/include/c/api2 -I${SRCDIR}/../../agora_sdk/include/c/base
// #include "agora_audio_sink_cgo.h"
import "C"
import (
	"fmt"
	"sync"
	"unsafe"
)

// AudioSinkContext audio sink context
type AudioSinkContext struct {
	cSink      unsafe.Pointer
	cAudioSink *C.struct__audio_sink
	OnFrame    func(frame *AudioFrame) bool
	mu         sync.RWMutex
}

// global singleton AudioSinkContext
var globalAudioSinkContext *AudioSinkContext
var globalAudioSinkMutex sync.RWMutex

//export goOnAudioFrame
func goOnAudioFrame(agora_audio_sink unsafe.Pointer, frame *C.struct__audio_pcm_frame) C.int {

	if frame == nil {
		fmt.Printf("[goOnAudioFrame] frame is nil\n")
		return C.int(0)
	}

	globalAudioSinkMutex.RLock()
	ctx := globalAudioSinkContext
	globalAudioSinkMutex.RUnlock()

	if ctx == nil {
		fmt.Printf("[goOnAudioFrame] ctx is nil\n")
		return C.int(0)
	}

	if ctx.OnFrame == nil {
		fmt.Printf("[goOnAudioFrame] ctx.OnFrame is nil\n")
		return C.int(0)
	}

	ctx.mu.RLock()
	defer ctx.mu.RUnlock()

	dataSize := int(frame.samples_per_channel * frame.num_channels)
	if dataSize <= 0 {
		fmt.Printf("[goOnAudioFrame] dataSize is less than 0\n")
		return C.int(0)
	}
	byteSize := dataSize * 2 // int16 = 2 bytes
	buffer := make([]byte, byteSize)
	for i := 0; i < dataSize; i++ {
		sample := int16(frame.data[i])
		buffer[i*2] = byte(sample)
		buffer[i*2+1] = byte(sample >> 8)
	}

	goFrame := &AudioFrame{
		SamplesPerChannel: int(frame.samples_per_channel),
		BytesPerSample:    int(frame.bytes_per_sample),
		Channels:          int(frame.num_channels),
		SamplesPerSec:     int(frame.sample_rate_hz),
		Buffer:            buffer,
		RenderTimeMs:      int64(frame.capture_timestamp),
		PresentTimeMs:     int64(frame.capture_timestamp),
		FarFieldFlag:      int(frame.audio_label.far_filed_flag),
		Rms:               int(frame.audio_label.rms),
		VoiceProb:         int(frame.audio_label.voice_prob),
		MusicProb:         int(frame.audio_label.music_prob),
		Pitch:             int(frame.audio_label.pitch),
	}

	var result bool
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("[goOnAudioFrame] PANIC caught: %v\n", r)
				result = false
			}
		}()
		result = ctx.OnFrame(goFrame)
	}()

	if result {
		return C.int(1)
	}
	return C.int(0)
}

// NewAudioSink creates a new audio sink
func NewAudioSink(onFrame func(frame *AudioFrame) bool) *AudioSinkContext {
	if onFrame == nil {
		return nil
	}

	ctx := &AudioSinkContext{
		OnFrame: onFrame,
	}

	ctx.cAudioSink = C.create_audio_sink_callbacks()
	if ctx.cAudioSink == nil {
		return nil
	}

	ctx.cSink = C.agora_audio_sink_create(ctx.cAudioSink)
	if ctx.cSink == nil {
		C.free(unsafe.Pointer(ctx.cAudioSink))
		return nil
	}

	globalAudioSinkMutex.Lock()
	if globalAudioSinkContext != nil {
		fmt.Println("[NewAudioSink] warning: AudioSink instance already exists, will be replaced")
	}
	globalAudioSinkContext = ctx
	globalAudioSinkMutex.Unlock()

	fmt.Printf("[NewAudioSink] created successfully: cAudioSink=%p, cSink=%p\n", ctx.cAudioSink, ctx.cSink)

	return ctx
}

// Release releases the audio sink
func (ctx *AudioSinkContext) Release() {
	if ctx == nil {
		return
	}

	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	globalAudioSinkMutex.Lock()
	if globalAudioSinkContext == ctx {
		globalAudioSinkContext = nil
	}
	globalAudioSinkMutex.Unlock()

	if ctx.cSink != nil {
		C.agora_audio_sink_destroy(ctx.cSink)
		ctx.cSink = nil
	}

	if ctx.cAudioSink != nil {
		C.free(unsafe.Pointer(ctx.cAudioSink))
		ctx.cAudioSink = nil
	}
}

// GetHandle gets the C handle
func (ctx *AudioSinkContext) GetHandle() unsafe.Pointer {
	if ctx == nil {
		return nil
	}
	return ctx.cSink
}
