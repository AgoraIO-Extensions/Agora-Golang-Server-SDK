package agoraservice

// #cgo CFLAGS: -I${SRCDIR}/../../agora_sdk/include/c/api2 -I${SRCDIR}/../../agora_sdk/include/c/base
// #include "agora_media_node_factory.h"
import "C"
import "unsafe"

type EncodedAudioFrameInfo struct {
	// Speech determines whether the audio frame source is a speech.
	// - 1: (Default) The audio frame source is a speech.
	// - 0: The audio frame source is not a speech.
	Speech bool

	// Codec is the audio codec: AudioCodecType.
	Codec AudioCodecType

	// SampleRateHz is the sample rate (Hz) of the audio frame.
	SampleRateHz int

	// SamplesPerChannel is the number of samples per audio channel.
	// If this value is not set, it is 1024 for AAC, 960 for OPUS default.
	SamplesPerChannel int

	// SendEvenIfEmpty determines whether to send the audio frame even when it is empty.
	// - 1: (Default) Send the audio frame even when it is empty.
	// - 0: Do not send the audio frame when it is empty.
	SendEvenIfEmpty bool

	// NumberOfChannels is the number of channels of the audio frame.
	NumberOfChannels int
}

type AudioEncodedFrameSender struct {
	cSender unsafe.Pointer
}

func (mediaNodeFactory *MediaNodeFactory) NewAudioEncodedFrameSender() *AudioEncodedFrameSender {
	sender := C.agora_media_node_factory_create_audio_encoded_frame_sender(mediaNodeFactory.cFactory)
	if sender == nil {
		return nil
	}
	return &AudioEncodedFrameSender{
		cSender: sender,
	}
}

func (sender *AudioEncodedFrameSender) Release() {
	if sender.cSender == nil {
		return
	}
	C.agora_audio_encoded_frame_sender_destroy(sender.cSender)
	sender.cSender = nil
}

func (sender *AudioEncodedFrameSender) SendEncodedAudioFrame(payload []byte, frameInfo *EncodedAudioFrameInfo) int {
	cData, pinner := unsafeCBytes(payload)
	defer pinner.Unpin()
	cFrameInfo := &C.struct__encoded_audio_frame_info{
		speech:              CIntFromBool(frameInfo.Speech),
		codec:               C.int(frameInfo.Codec),
		sample_rate_hz:      C.int(frameInfo.SampleRateHz),
		samples_per_channel: C.int(frameInfo.SamplesPerChannel),
		send_even_if_empty:  CIntFromBool(frameInfo.SendEvenIfEmpty),
		number_of_channels:  C.int(frameInfo.NumberOfChannels),
	}

	return int(C.agora_audio_encoded_frame_sender_send(
		sender.cSender, (*C.uint8_t)(cData),
		C.uint32_t(len(payload)), cFrameInfo))
}
