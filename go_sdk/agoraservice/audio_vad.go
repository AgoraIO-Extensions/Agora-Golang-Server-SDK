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
	"fmt"
	"os"
	"unsafe"
)

type AudioVadConfig struct {
	StartRecognizeCount    int     // start recognize count, buffer size for 10ms 16KHz 16bit 1channel PCM, default value is 10
	StopRecognizeCount     int     // max recognize count, buffer size for 10ms 16KHz 16bit 1channel PCM, default value is 6
	PreStartRecognizeCount int     // pre start recognize count, buffer size for 10ms 16KHz 16bit 1channel PCM, default value is 10
	ActivePercent          float32 // active percent, if over this percent, will be recognized as speaking, default value is 0.6
	InactivePercent        float32 // inactive percent, if below this percent, will be recognized as non-speaking, default value is 0.2
	VoiceProb   float32 // voice probability threshold, default value is 0.7; range from 0 to 1
	RmsThr float32 // rms threshold, default value is -40; range from -100 to 0
	JointThr float32 // joint threshold, default value is 0.0; range from 0 to 1
	Aggressive float32 // default value is 2.0; range from 0 to 3
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
		cVad: nil,
		// lastStatus: VAD_WAIT_SPEEKING,
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
/*
* for stero vad
*/
/*
* 考虑建立一个针对双声道的vad，用来做双声道的处理
* 可以放在audio_vad中
* 在那个里面，可以有函数：
* steroAudioVad
* stero pcm--> mono pcm
* 还是说暂时放在sample里面？但提供stero pcm convert to mono pcm?
*/
type SteroAudioVad struct {
	LeftVadInstance *AudioVad
	RightVadInstance *AudioVad
	LeftVadConfigure *AudioVadConfig
	RightVadConfigure *AudioVadConfig
    
}
func NewSteroVad(leftVadConfig *AudioVadConfig, rightVadConfig *AudioVadConfig) *SteroAudioVad{
	return &SteroAudioVad{
		LeftVadInstance: NewAudioVad(leftVadConfig),
		RightVadInstance: NewAudioVad(rightVadConfig),
		LeftVadConfigure: leftVadConfig,
		RightVadConfigure: rightVadConfig,
		}
	}
// chagne frame to unsafe pointer， private function, only for internal use
func bytesToInt16Array(data []byte) []int16 {
	// 通过 unsafe 包将 []byte 转换为 []int16
	//return *(*[]int16)(unsafe.Pointer(&data))
	return *(*[]int16)(unsafe.Pointer(&data))
}
// return value: 1. left vad state, 2. right vad state
// test for dump stereo pcm to mono
var (
	LeftFile  *os.File = nil
	RightFile *os.File = nil
	DebugMonoPcm int = 0
)
func (vad *SteroAudioVad) ProcessAudioFrame(inFrame *AudioFrame) (*AudioFrame, int, *AudioFrame, int) {
	// 0.validity check: only support 16k 2 channel 16 bit pcm
	if inFrame == nil  || inFrame.Buffer == nil || inFrame.SamplesPerSec != 16000 || inFrame.Channels != 2 || inFrame.BytesPerSample != 2 {
		fmt.Printf("invalid: samplesPerSec: %d, channels: %d, bytesPerSample: %d\n", inFrame.SamplesPerSec, inFrame.Channels, inFrame.BytesPerSample)
		return nil, 0, nil, 0
	}
	if DebugMonoPcm > 0 && ( nil == LeftFile || nil == RightFile) {
	    LeftFile, _ = os.OpenFile("./left.pcm", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		RightFile, _ = os.OpenFile("./right.pcm", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

	}
	// split stero pcm to 2 mono pcm
	// process vad for each mono pcm, and return the vad state
	// 1. allocate 2 new AudioFrame for 16khz mono and 16bits pcm
	//1. split stero pcm to 2 mono pcm through unsafe pointer mode to save time
	inLength  := len(inFrame.Buffer)
	channelDataLen := inLength/2
    leftBuffer := make([]byte, channelDataLen)
	rightBuffer := make([]byte, channelDataLen)
	dataLen := channelDataLen/2

	//fmt.Printf("info: samplesPerSec: %d, channels: %d, bytesPerSample: %d, len: %d\n", 
	//inFrame.SamplesPerSec, inFrame.Channels, inFrame.BytesPerSample, inLength)


	lefeFrame := &AudioFrame{
		Type: inFrame.Type,
		SamplesPerChannel: 160,
		BytesPerSample: 2,
		Channels: 1,
		SamplesPerSec: 16000,
		Buffer:       leftBuffer,
		RenderTimeMs: 0,
	}
	rightFrame := &AudioFrame{
		Type: inFrame.Type,
		SamplesPerChannel: 160,
		BytesPerSample: 2,
		Channels: 1,
		SamplesPerSec: 16000,
		Buffer:       rightBuffer,
		RenderTimeMs: 0,
	}
	// 2. split stero pcm to 2 mono pcm
	// 2.1 convert byte to uint16
	ptrInframe := bytesToInt16Array(inFrame.Buffer)
	ptrLeftFrame := bytesToInt16Array(lefeFrame.Buffer)
	ptrRightFrame := bytesToInt16Array(rightFrame.Buffer)

	//2.2 assign stereo pcm to mono pcm
	
	i := 0
	for j := 0; j < dataLen; j++ {
		ptrLeftFrame[j] = ptrInframe[i]
		ptrRightFrame[j] = ptrInframe[i+1]
		i += 2
	}

	//2.3 for debug mono pcm
	if DebugMonoPcm > 0 && LeftFile != nil && RightFile != nil {
		LeftFile.Write(lefeFrame.Buffer)
		RightFile.Write(rightFrame.Buffer)
	}
	
	

	// 3. do vad test for each mono pcm
	leftVadResultFrame, leftVadState := vad.LeftVadInstance.ProcessPcmFrame(lefeFrame)
	rightVadResultFrame, rightVadStat := vad.RightVadInstance.ProcessPcmFrame(rightFrame)
	//fmt.Printf("left vad result: %v, left vad state: %v, right vad result: %v, right vad state: %v\n",leftVadResultFrame, leftVadState, rightVadResultFrame, rightVadStat)
	// 4. do vad test for each mono pcm
	return leftVadResultFrame, leftVadState, rightVadResultFrame, rightVadStat
}
func (vad *SteroAudioVad) Release() {
	if vad.LeftVadInstance != nil {
		vad.LeftVadInstance.Release()
	}
	if vad.RightVadInstance != nil {
		vad.RightVadInstance.Release()
	}
	vad.LeftVadInstance = nil
	vad.RightVadInstance = nil
}
    
