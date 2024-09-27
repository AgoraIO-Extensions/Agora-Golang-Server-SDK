package agoraservice

// #cgo CFLAGS: -I${SRCDIR}/../../agora_sdk/include/c/api2 -I${SRCDIR}/../../agora_sdk/include/c/base
// #include "agora_media_node_factory.h"
import "C"
import "unsafe"

type AudioFrame struct {
	Type              AudioFrameType
	SamplesPerChannel int    // The number of samples per channel in this frame.
	BytesPerSample    int    // The number of bytes per sample: Two for PCM 16.
	Channels          int    // The number of channels (data is interleaved, if stereo).
	SamplesPerSec     int    // The Sample rate.
	Buffer            []byte // The pointer to the data buffer.
	RenderTimeMs      int64  // The timestamp to render the audio data. Use this member to synchronize the audio renderer while rendering the audio streams.
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
}

func (mediaNodeFactory *MediaNodeFactory) NewAudioPcmDataSender() *AudioPcmDataSender {
	sender := C.agora_media_node_factory_create_audio_pcm_data_sender(mediaNodeFactory.cFactory)
	if sender == nil {
		return nil
	}
	return &AudioPcmDataSender{
		cSender: sender,
	}
}

func (sender *AudioPcmDataSender) Release() {
	if sender.cSender == nil {
		return
	}
	C.agora_audio_pcm_data_sender_destroy(sender.cSender)
	sender.cSender = nil
}

func (sender *AudioPcmDataSender) SendAudioPcmData(frame *AudioFrame) int {
	cData := C.CBytes(frame.Buffer)
	defer C.free(cData)
	return int(C.agora_audio_pcm_data_sender_send(sender.cSender, cData,
		C.uint(frame.RenderTimeMs), C.uint(frame.SamplesPerChannel),
		C.uint(frame.BytesPerSample), C.uint(frame.Channels),
		C.uint(frame.SamplesPerSec)))
}
