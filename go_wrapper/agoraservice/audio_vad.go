package agoraservice

/*
#cgo CFLAGS: -I../../agora_sdk/include_c/api2 -I../../agora_sdk/include_c/base
#cgo darwin LDFLAGS: -Wl,-rpath,../../agora_sdk_mac -L../../agora_sdk_mac -luap_aed
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

// cVadCfg.fftSz = C.int(1024)
// cVadCfg.anaWindowSz = C.int(768)
// cVadCfg.hopSz = C.int(160)
// cVadCfg.frqInputAvailableFlag = C.int(0)
// cVadCfg.useCVersionAIModule = C.int(0)
// cVadCfg.voiceProbThr = C.float(0.7)
// cVadCfg.rmsThr = C.float(-40.0)
// cVadCfg.jointThr = C.float(0.0)
// cVadCfg.aggressive = C.float(2.0)
type AudioVadConfig struct {
	FftSz                 int     // fft size, default value is 1024
	AnaWindowSz           int     // analysis window size, default value is 768
	HopSz                 int     // hop size, default value is 160
	FrqInputAvailableFlag int     // frequency input available flag, default value is 0
	UseCVersionAIModule   int     // use C version AI module, default value is 0
	VoiceProbThr          float32 // voice probability threshold, default value is 0.7
	RmsThr                float32 // root mean square threshold, default value is -40.0
	JointThr              float32 // joint threshold, default value is 0.0
	Aggressive            float32 // aggressive, default value is 2.0

	StartRecognizeCount    int     // start recognize count, buffer size for 10ms 16KHz 16bit 1channel PCM, default value is 10
	StopRecognizeCount     int     // max recognize count, buffer size for 10ms 16KHz 16bit 1channel PCM, default value is 6
	PreStartRecognizeCount int     // pre start recognize count, buffer size for 10ms 16KHz 16bit 1channel PCM, default value is 10
	ActivePercent          float32 // active percent, if over this percent, will be recognized as speaking, default value is 0.6
	InactivePercent        float32 // inactive percent, if below this percent, will be recognized as non-speaking, default value is 0.2
}

type AudioVad struct {
	vadCfg    *AudioVadConfig
	cVad      unsafe.Pointer
	lastOutTs int64
	// lastStatus int
}

func NewAudioVad(cfg *AudioVadConfig) *AudioVad {
	if cfg == nil {
		cfg = &AudioVadConfig{
			FftSz:                 1024,
			AnaWindowSz:           768,
			HopSz:                 160,
			FrqInputAvailableFlag: 0,
			UseCVersionAIModule:   0,
			VoiceProbThr:          0.7,
			RmsThr:                -40.0,
			JointThr:              0.0,
			Aggressive:            2.0,

			StartRecognizeCount:    10,
			StopRecognizeCount:     6,
			PreStartRecognizeCount: 10,
			ActivePercent:          0.6,
			InactivePercent:        0.2,
		}
	}
	vad := &AudioVad{
		vadCfg:    cfg,
		lastOutTs: 0,
		// lastStatus: VAD_WAIT_SPEEKING,
	}
	cVadCfg := C.struct_Vad_Config_{}
	C.memset((unsafe.Pointer)(&cVadCfg), 0, C.sizeof_struct_Vad_Config_)
	cVadCfg.fftSz = C.int(cfg.FftSz)
	cVadCfg.anaWindowSz = C.int(cfg.AnaWindowSz)
	cVadCfg.hopSz = C.int(cfg.HopSz)
	cVadCfg.frqInputAvailableFlag = C.int(cfg.FrqInputAvailableFlag)
	cVadCfg.useCVersionAIModule = C.int(cfg.UseCVersionAIModule)
	cVadCfg.voiceProbThr = C.float(cfg.VoiceProbThr)
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
	outFrame := &PcmAudioFrame{
		Data:              outData,
		Timestamp:         vad.lastOutTs,
		SamplesPerChannel: samplesPerChannel,
		BytesPerSample:    2,
		NumberOfChannels:  1,
		SampleRate:        16000,
	}
	vad.lastOutTs += int64(frameDuration)

	return outFrame, int(vadState)
}
