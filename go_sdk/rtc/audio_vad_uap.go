//go:build vad_uap
// +build vad_uap

package agoraservice

/*
#cgo CFLAGS: -I../../agora_sdk/include/c/api2 -I../../agora_sdk/include/c/base
#cgo darwin LDFLAGS: -L../../agora_sdk_mac -luap_aed
#cgo linux LDFLAGS: -L../../agora_sdk/ -lagora_uap_aed
#include <string.h>
#include "vad.h"
*/
import "C"
import (
	"unsafe"
)

func VadUAPSupported() bool {
	return true
}

func NewAudioVad(cfg *AudioVadConfig) *AudioVad {
	if cfg == nil {
		cfg = &AudioVadConfig{
			StartRecognizeCount:    30,
			StopRecognizeCount:     48,
			PreStartRecognizeCount: 16,
			ActivePercent:          0.8,
			InactivePercent:        0.2,
			RmsThr:                 -40.0,
			JointThr:               0.0,
			Aggressive:             2.0,
			VoiceProb:              0.7,
		}
	}
	vad := &AudioVad{
		vadCfg:    cfg,
		lastOutTs: 0,
		cVad:      nil,
	}
	cVadCfg := C.struct_Vad_Config_{}
	C.memset((unsafe.Pointer)(&cVadCfg), 0, C.sizeof_struct_Vad_Config_)
	cVadCfg.fftSz = C.int(1024)
	cVadCfg.anaWindowSz = C.int(768)
	cVadCfg.hopSz = C.int(160)
	cVadCfg.frqInputAvailableFlag = C.int(0)
	cVadCfg.useCVersionAIModule = C.int(0)
	cVadCfg.voiceProbThr = C.float(cfg.VoiceProb)
	cVadCfg.rmsThr = C.float(cfg.RmsThr)
	cVadCfg.jointThr = C.float(cfg.JointThr)
	cVadCfg.aggressive = C.float(cfg.Aggressive)
	cVadCfg.startRecognizeCount = C.int(cfg.StartRecognizeCount)
	cVadCfg.stopRecognizeCount = C.int(cfg.StopRecognizeCount)
	cVadCfg.preStartRecognizeCount = C.int(cfg.PreStartRecognizeCount)
	cVadCfg.activePercent = C.float(cfg.ActivePercent)
	cVadCfg.inactivePercent = C.float(cfg.InactivePercent)
	ret := int(C.Agora_UAP_VAD_Create(&vad.cVad, &cVadCfg))
	if ret != 0 {
		return nil
	}

	return vad
}

func (vad *AudioVad) Release() {
	if vad.cVad == nil {
		return
	}
	C.Agora_UAP_VAD_Destroy(&vad.cVad)
	vad.cVad = nil
}

func (vad *AudioVad) ProcessPcmFrame(frame *AudioFrame) (*AudioFrame, int) {
	if frame.SamplesPerSec != 16000 || frame.Channels != 1 || frame.BytesPerSample != 2 {
		return nil, -1
	}
	if vad.cVad == nil {
		return nil, -1
	}
	cData := C.CBytes(frame.Buffer)
	defer C.free(cData)
	in := C.Vad_AudioData{
		audioData: (unsafe.Pointer)(cData),
		size:      C.int(len(frame.Buffer)),
	}
	var vadState C.enum_VAD_STATE = C.enum_VAD_STATE(0)
	var out C.Vad_AudioData
	C.memset((unsafe.Pointer)(&out), 0, C.sizeof_struct_Vad_AudioData_)
	ret := int(C.Agora_UAP_VAD_Proc(vad.cVad, &in, &out, &vadState))
	if ret < 0 {
		return nil, ret
	}
	samplesPerChannel := int(out.size) / 2 / 1
	frameDuration := 1000 * samplesPerChannel / 16000
	outData := C.GoBytes(out.audioData, out.size)
	outFrame := &AudioFrame{
		Type:              AudioFrameTypePCM16,
		Buffer:            outData,
		RenderTimeMs:      vad.lastOutTs,
		SamplesPerChannel: samplesPerChannel,
		BytesPerSample:    2,
		Channels:          1,
		SamplesPerSec:     16000,
	}
	vad.lastOutTs += int64(frameDuration)

	return outFrame, int(vadState)
}
