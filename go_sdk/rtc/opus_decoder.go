//go:build avcodec
// +build avcodec

package agoraservice

// #cgo CFLAGS: -DUSE_AVCODEC
// #cgo pkg-config: libavformat libavcodec libavutil libswresample
// #include <string.h>
// #include <stdlib.h>
// #include "opus_decode.h"
import "C"
import (
	"fmt"
	"unsafe"
)


type OpusDecoder struct {
	handle      unsafe.Pointer
	cframe      C.struct__MediaFrame
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
	var errorCode C.int
	handle := C.open_opus_decoder(C.int(in_channels), C.int(out_sample_rate), C.int(out_channels), &errorCode)
	if handle == nil {
		fmt.Printf("0202:200318_v1: failed to open opus decoder: %d\n", errorCode)
		return nil
	}
	decoder := &OpusDecoder{handle: handle, sample_rate: out_sample_rate, bits: 16, channel: out_channels}
	decoder.cframe = C.struct__MediaFrame{}
	C.memset(unsafe.Pointer(&decoder.cframe), 0, C.sizeof_struct__MediaFrame)
	return decoder
}

// Decode decodes one raw opus packet into a PCM S16 MediaFrame.
// returns 0 on success (frame filled), AVERROR(EAGAIN) when more data is needed,
// or a negative error code on failure.
func (decoder *OpusDecoder) Decode(packet []byte) ([]byte, error) {
	frame := &decoder.cframe
	C.memset(unsafe.Pointer(frame), 0, C.sizeof_struct__MediaFrame)
	ret := C.decode_opus(decoder.handle, (*C.uint8_t)(unsafe.Pointer(&packet[0])), C.int(len(packet)), frame)
	if ret != 0 {
		return nil, fmt.Errorf("decode opus failed: %d", ret)
	}
	return C.GoBytes(unsafe.Pointer(frame.buffer), C.int(frame.buffer_size)), nil
}

func (decoder *OpusDecoder) Release() {
	C.close_opus_decoder(decoder.handle)
	decoder.handle = nil
}