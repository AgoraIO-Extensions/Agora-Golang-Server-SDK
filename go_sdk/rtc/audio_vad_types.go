package agoraservice

import (
	"errors"
	"fmt"
	"os"
	"unsafe"
)

// ErrVadUAPNotEnabled is returned when UAP VAD (libagora_uap_aed) was not linked at build time.
var ErrVadUAPNotEnabled = errors.New("agora UAP VAD is not enabled; build with -tags vad_uap")

type AudioVadConfig struct {
	StartRecognizeCount    int     // start recognize count, buffer size for 10ms 16KHz 16bit 1channel PCM, default value is 10
	StopRecognizeCount     int     // max recognize count, buffer size for 10ms 16KHz 16bit 1channel PCM, default value is 6
	PreStartRecognizeCount int     // pre start recognize count, buffer size for 10ms 16KHz 16bit 1channel PCM, default value is 10
	ActivePercent          float32 // active percent, if over this percent, will be recognized as speaking, default value is 0.6
	InactivePercent        float32 // inactive percent, if below this percent, will be recognized as non-speaking, default value is 0.2
	VoiceProb              float32 // voice probability threshold, default value is 0.7; range from 0 to 1
	RmsThr                 float32 // rms threshold, default value is -40; range from -100 to 0
	JointThr               float32 // joint threshold, default value is 0.0; range from 0 to 1
	Aggressive             float32 // default value is 2.0; range from 0 to 3
}

type AudioVad struct {
	vadCfg    *AudioVadConfig
	cVad      unsafe.Pointer
	lastOutTs int64
}

type SteroAudioVad struct {
	LeftVadInstance   *AudioVad
	RightVadInstance  *AudioVad
	LeftVadConfigure  *AudioVadConfig
	RightVadConfigure *AudioVadConfig
}

func NewSteroVad(leftVadConfig *AudioVadConfig, rightVadConfig *AudioVadConfig) *SteroAudioVad {
	return &SteroAudioVad{
		LeftVadInstance:   NewAudioVad(leftVadConfig),
		RightVadInstance:  NewAudioVad(rightVadConfig),
		LeftVadConfigure:  leftVadConfig,
		RightVadConfigure: rightVadConfig,
	}
}

func bytesToInt16Array(data []byte) []int16 {
	return *(*[]int16)(unsafe.Pointer(&data))
}

var (
	LeftFile     *os.File = nil
	RightFile    *os.File = nil
	DebugMonoPcm int      = 0
)

func (vad *SteroAudioVad) ProcessAudioFrame(inFrame *AudioFrame) (*AudioFrame, int, *AudioFrame, int) {
	if inFrame == nil || inFrame.Buffer == nil || inFrame.SamplesPerSec != 16000 || inFrame.Channels != 2 || inFrame.BytesPerSample != 2 {
		fmt.Printf("invalid: samplesPerSec: %d, channels: %d, bytesPerSample: %d\n", inFrame.SamplesPerSec, inFrame.Channels, inFrame.BytesPerSample)
		return nil, 0, nil, 0
	}
	if DebugMonoPcm > 0 && (LeftFile == nil || RightFile == nil) {
		LeftFile, _ = os.OpenFile("./left.pcm", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		RightFile, _ = os.OpenFile("./right.pcm", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	}

	inLength := len(inFrame.Buffer)
	channelDataLen := inLength / 2
	leftBuffer := make([]byte, channelDataLen)
	rightBuffer := make([]byte, channelDataLen)
	dataLen := channelDataLen / 2

	leftFrame := &AudioFrame{
		Type:              inFrame.Type,
		SamplesPerChannel: 160,
		BytesPerSample:    2,
		Channels:          1,
		SamplesPerSec:     16000,
		Buffer:            leftBuffer,
		RenderTimeMs:      0,
	}
	rightFrame := &AudioFrame{
		Type:              inFrame.Type,
		SamplesPerChannel: 160,
		BytesPerSample:    2,
		Channels:          1,
		SamplesPerSec:     16000,
		Buffer:            rightBuffer,
		RenderTimeMs:      0,
	}

	ptrInframe := bytesToInt16Array(inFrame.Buffer)
	ptrLeftFrame := bytesToInt16Array(leftFrame.Buffer)
	ptrRightFrame := bytesToInt16Array(rightFrame.Buffer)

	i := 0
	for j := 0; j < dataLen; j++ {
		ptrLeftFrame[j] = ptrInframe[i]
		ptrRightFrame[j] = ptrInframe[i+1]
		i += 2
	}

	if DebugMonoPcm > 0 && LeftFile != nil && RightFile != nil {
		LeftFile.Write(leftFrame.Buffer)
		RightFile.Write(rightFrame.Buffer)
	}

	leftVadResultFrame, leftVadState := vad.LeftVadInstance.ProcessPcmFrame(leftFrame)
	rightVadResultFrame, rightVadStat := vad.RightVadInstance.ProcessPcmFrame(rightFrame)
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
