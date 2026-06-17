//go:build !avcodec
// +build !avcodec

package agoraservice

// #cgo CFLAGS: -DUSE_AVCODEC
// #cgo pkg-config: libavformat libavcodec libavutil libswresample
// #include <string.h>
// #include <stdlib.h>
// #include "opus_decode.h"
import "C"
import (
	"fmt"
)


type OpusDecoder struct {
	handle      *any
	cframe      any
	sample_rate int
	bits        int
	channel     int
}

// NewOpusDecoder creates a new OpusDecoder instance.
// in_channels: channel count of the opus stream (1 or 2). pass 0 to default to 1.
// out_sample_rate: desired output PCM sample rate, e.g. 16000/48000. pass 0 to keep 48000.
// out_channels:    desired output PCM channel count. pass 0 to keep in_channels.
// note: the bits is fixed to 16 bits for this version and also for opus decoder.
func NewOpusDecoder(in_channels int, out_sample_rate int, out_channels int) *OpusDecoder {
	return nil
}

// Decode decodes one raw opus packet into a PCM S16 MediaFrame.
// returns 0 on success (frame filled), AVERROR(EAGAIN) when more data is needed,
// or a negative error code on failure.
func (decoder *OpusDecoder) Decode(packet []byte) ([]byte, error) {
	return nil, fmt.Errorf("opusdecoder: not implemented")
}

func (decoder *OpusDecoder) Release() {
	return
}