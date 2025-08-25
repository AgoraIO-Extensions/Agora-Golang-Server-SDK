package agoraservice

// #cgo CFLAGS: -I${SRCDIR}/../../agora_sdk/include/c/api2 -I${SRCDIR}/../../agora_sdk/include/c/base
// #include "agora_media_node_factory.h"
import "C"
import (
	"sync"
	"unsafe"
)

type AudioFrame struct {
	Type              AudioFrameType
	SamplesPerChannel int    // The number of samples per channel in this frame.
	BytesPerSample    int    // The number of bytes per sample: Two for PCM 16.
	Channels          int    // The number of channels (data is interleaved, if stereo).
	SamplesPerSec     int    // The Sample rate.
	Buffer            []byte // The pointer to the data buffer.
	RenderTimeMs      int64  // The timestamp to render the audio data. Use this member to synchronize the audio renderer while rendering the audio streams.
	PresentTimeMs     int64  // The timestamp to present the audio data. Use this member to synchronize the audio renderer while rendering the audio streams.
	AvsyncType        int

	// these field below are only used for audio observer.
	FarFieldFlag int
	Rms          int
	VoiceProb    int
	MusicProb    int
	Pitch        int
}

type AudioPcmDataSender struct {
	cSender unsafe.Pointer
	mu      sync.RWMutex
	closed  bool
	audioScenario AudioScenario
}
type AudioVolumeInfo struct {
	UserId     string
	Volume     uint32
	VAD        uint32
	VoicePitch float64
}

func (mediaNodeFactory *MediaNodeFactory) NewAudioPcmDataSender() *AudioPcmDataSender {
	if mediaNodeFactory == nil || mediaNodeFactory.cFactory == nil {
		return nil
	}
	sender := C.agora_media_node_factory_create_audio_pcm_data_sender(mediaNodeFactory.cFactory)
	if sender == nil {
		return nil
	}
	return &AudioPcmDataSender{
		cSender: sender,
		closed:  false,
		audioScenario: AudioScenarioChorus,  //default to chorus mode
	}
}

func (sender *AudioPcmDataSender) Release() {
	if sender.closed {
		return
	}
	sender.mu.Lock()
	defer sender.mu.Unlock()
	sender.closed = true
	if sender.cSender == nil {
		return
	}
	C.agora_audio_pcm_data_sender_destroy(sender.cSender)
	sender.cSender = nil
}

func (sender *AudioPcmDataSender) SendAudioPcmData(frame *AudioFrame) int {
	if sender.closed || sender.cSender == nil || frame == nil  {
		return -1
	}
	sender.mu.RLock()
	defer sender.mu.RUnlock()
	if sender.closed || sender.cSender == nil || frame == nil || len(frame.Buffer) == 0 {
		return -1
	}
	cData, pinner := unsafeCBytes(frame.Buffer)
	defer pinner.Unpin()
	return int(C.agora_audio_pcm_data_sender_send(sender.cSender, cData,
		C.uint(frame.RenderTimeMs), C.int64_t(frame.PresentTimeMs),
		C.uint(frame.SamplesPerChannel),
		C.uint(frame.BytesPerSample), C.uint(frame.Channels),
		C.uint(frame.SamplesPerSec)))
}
