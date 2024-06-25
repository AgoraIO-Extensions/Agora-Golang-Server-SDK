package agoraservice

/*
#cgo CFLAGS: -I../../agora_sdk/include_c/api2 -I../../agora_sdk/include_c/base
#cgo darwin arm64 LDFLAGS: -L../../agora_sdk_mac
#cgo linux LDFLAGS: -L../../agora_sdk/ -lagora_uap_aed
#include <string.h>
#include "vad.h"
*/
import "C"
import "unsafe"

const (
	VAD_WAIT_SPEEKING  = 0
	VAD_START_SPEEKING = 1
	VAD_IS_SPEEKING    = 2
	VAD_STOP_SPEEKING  = 3
)

type AudioVadConfig struct {
	StartRecognizeCount int     // start recognize count, buffer size for 40ms 16KHz 16bit 1channel PCM, default value is 8
	MaxRecognizeCount   int     // max recognize count, buffer size for 40ms 16KHz 16bit 1channel PCM, default value is 20
	ActivePercent       float32 // active percent, if over this percent, will be recognized as speaking, default value is 0.6
	InactivePercent     float32 // inactive percent, if below this percent, will be recognized as non-speaking, default value is 0.2
}

type AudioVad struct {
	vadCfg     *AudioVadConfig
	cVad       unsafe.Pointer
	lastOutTs  int64
	lastStatus int
}

func NewAudioVad(cfg *AudioVadConfig) *AudioVad {
	if cfg == nil {
		cfg = &AudioVadConfig{
			StartRecognizeCount: 8,
			MaxRecognizeCount:   20,
			ActivePercent:       0.6,
			InactivePercent:     0.2,
		}
	}
	vad := &AudioVad{
		vadCfg:     cfg,
		lastOutTs:  0,
		lastStatus: VAD_WAIT_SPEEKING,
	}
	cVadCfg := C.struct_Vad_Config_{}
	C.memset((unsafe.Pointer)(&cVadCfg), 0, C.sizeof_struct_Vad_Config_)
	cVadCfg.fftSz = C.int(1024)
	cVadCfg.anaWindowSz = C.int(768)
	cVadCfg.hopSz = C.int(160)
	cVadCfg.frqInputAvailableFlag = C.int(0)
	cVadCfg.useCVersionAIModule = C.int(0)
	cVadCfg.voiceProbThr = C.float(0.7)
	cVadCfg.rmsThr = C.float(-40.0)
	cVadCfg.jointThr = C.float(0.0)
	cVadCfg.aggressive = C.float(2.0)
	cVadCfg.startRecognizeCount = C.int(cfg.StartRecognizeCount)
	cVadCfg.maxRecognizeCount = C.int(cfg.MaxRecognizeCount)
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
}

func (vad *AudioVad) ProcessPcmFrame(frame *PcmAudioFrame) (*PcmAudioFrame, int) {
	if frame.SampleRate != 16000 || frame.NumberOfChannels != 1 || frame.BytesPerSample != 2 {
		return nil, -1
	}
	cData := C.CBytes(frame.Data)
	defer C.free(cData)
	in := C.Vad_AudioData{
		audioData: (unsafe.Pointer)(cData),
		size:      C.int(len(frame.Data)),
	}
	var out C.Vad_AudioData
	C.memset((unsafe.Pointer)(&out), 0, C.sizeof_struct_Vad_AudioData_)
	ret := int(C.Agora_UAP_VAD_Proc(vad.cVad, &in, &out))
	if ret < 0 {
		return nil, ret
	}
	if out.size == 0 || out.size == 2 {
		if vad.lastStatus == VAD_IS_SPEEKING ||
			vad.lastStatus == VAD_START_SPEEKING {
			vad.lastStatus = VAD_STOP_SPEEKING
			return nil, VAD_STOP_SPEEKING
		}
		vad.lastStatus = VAD_WAIT_SPEEKING
		return nil, VAD_WAIT_SPEEKING
	}
	if vad.lastStatus == VAD_WAIT_SPEEKING ||
		vad.lastStatus == VAD_STOP_SPEEKING {
		vad.lastStatus = VAD_START_SPEEKING
	} else {
		vad.lastStatus = VAD_IS_SPEEKING
	}
	samplesPerChannel := int(out.size) / 2 / 1
	frameDuration := 1000 * samplesPerChannel / 16000
	outData := C.GoBytes(out.audioData, out.size)
	outFrame := &PcmAudioFrame{
		Data:              outData,
		Timestamp:         vad.lastOutTs,
		SamplesPerChannel: samplesPerChannel,
		BytesPerSample:    2,
		NumberOfChannels:  1,
		SampleRate:        16000,
	}
	vad.lastOutTs += int64(frameDuration)

	return outFrame, vad.lastStatus
}
