package agoraservice

import (
	"container/list"
	"fmt"
)

type AudioVadConfigV2 struct {
	PreStartRecognizeCount int     // pre start recognize count, buffer size for 10ms 16KHz 16bit 1channel PCM, default value is 16
	StartRecognizeCount    int     // start recognize count, buffer size for 10ms 16KHz 16bit 1channel PCM, default value is 30
	StopRecognizeCount     int     // max recognize count, buffer size for 10ms 16KHz 16bit 1channel PCM, default value is 50
	ActivePercent          float32 // active percent, if over this percent, will be recognized as speaking, default value is 0.7
	InactivePercent        float32 // inactive percent, if below this percent, will be recognized as non-speaking, default value is 0.5
	StartVoiceProb         int     // start voice prob, default value is 70
	StartRms               int     // start rms, default value is -50
	StopVoiceProb          int     // stop voice prob, default value is 70
	StopRms                int     // stop rms, default value is -50
	EnableAdaptiveRmsThreshold bool    // enable adaptive threshold, default value is false
	AdaptiveRmsThresholdFactor float32 // default to : 0.67.i.e 2/3
}

type VadFrame struct {
	frame    *AudioFrame
	isActive bool
}

type VadBuffer struct {
	queue   *list.List
	maxSize int
}

type VadFrameFormat struct {
	BytesPerSample int // The number of bytes per sample: Two for PCM 16.
	Channels       int // The number of channels (data is interleaved, if stereo).
	SamplesPerSec  int // The Sample rate.
}

type AudioVadV2 struct {
	config       *AudioVadConfigV2
	expectFormat *VadFrameFormat
	isSpeaking   bool
	startBuffer  *VadBuffer
	stopBuffer   *VadBuffer
	voiceCount   int
	silenceCount int
	totalVoiceRms  int // range from 0 to 127, respond to db: -127db, to 0db
	refAvgRmsInLastSesseion   int // range from 0 to 127, respond to db: -127db, to 0db
}

func newVadFrame(frame *AudioFrame, isActive bool) *VadFrame {
	return &VadFrame{
		frame:    frame,
		isActive: isActive,
	}
}

func newVadBuffer(maxSize int) *VadBuffer {
	return &VadBuffer{
		queue:   list.New(),
		maxSize: maxSize,
	}
}

func (buf *VadBuffer) pushBack(v *VadFrame) bool {
	l := buf.queue
	if l.Len() >= buf.maxSize {
		l.Remove(l.Front())
	}
	l.PushBack(v)
	return l.Len() >= buf.maxSize
}

func (buf *VadBuffer) clear() {
	buf.queue.Init()
}

func (buf *VadBuffer) getActivePercent(lastN int) float32 {
	l := buf.queue
	count := 0
	curIndex := 0
	startElem := l.Len() - lastN
	for e := l.Front(); e != nil; e = e.Next() {
		if curIndex < startElem {
			curIndex++
			continue
		}
		v := e.Value.(*VadFrame)
		if v.isActive {
			count++
		}
		curIndex++
	}
	// fmt.Printf("[vad] getActivePercent: %d, %d, %f\n", count, lastN, float32(count)/float32(lastN))
	return float32(count) / float32(lastN)
}

func (buf *VadBuffer) flushAudio() *AudioFrame {
	l := buf.queue
	if l.Len() == 0 {
		return nil
	}
	// copy a frame
	samplesCount := 0
	ret := *(l.Front().Value.(*VadFrame).frame)
	data := make([]byte, 0, l.Len()*ret.SamplesPerChannel*ret.BytesPerSample*ret.Channels)
	for e := l.Front(); e != nil; e = e.Next() {
		v := e.Value.(*VadFrame)
		data = append(data, v.frame.Buffer...)
		samplesCount += v.frame.SamplesPerChannel
	}
	ret.Buffer = data
	ret.SamplesPerChannel = samplesCount
	l.Init()
	return &ret
}

func NewAudioVadV2(cfg *AudioVadConfigV2) *AudioVadV2 {
	
	if cfg == nil {
		cfg = &AudioVadConfigV2{
			PreStartRecognizeCount: 16,
			StartRecognizeCount:    30,
			StopRecognizeCount:     65,
			ActivePercent:          0.7,
			InactivePercent:        0.5,
			StartVoiceProb:         70,
			StartRms:               -70.0,
			StopVoiceProb:          70,
			StopRms:                -70.0,
			EnableAdaptiveRmsThreshold: true,
			AdaptiveRmsThresholdFactor: 0.67,
		}
	}
	if cfg.StartRecognizeCount <= 0 {
		cfg.StartRecognizeCount = 30
	}
	if cfg.StopRecognizeCount <= 0 {
		cfg.StopRecognizeCount = 65
	}
	// fmt.Printf("[vad] NewAudioVadV2: %v\n", cfg)
	startQueueSize := cfg.StartRecognizeCount + cfg.PreStartRecognizeCount
	ret := &AudioVadV2{
		config:       cfg,
		expectFormat: nil,
		isSpeaking:   false,
		startBuffer:  newVadBuffer(startQueueSize),
		stopBuffer:   newVadBuffer(cfg.StopRecognizeCount),
		voiceCount:   0,
		silenceCount: 0,
		totalVoiceRms:  0,
		refAvgRmsInLastSesseion: 0,
	}
	// for this version, convert rms from db to rms,i.e convert from -127db to 0db to 0 to 127
	ret.config.StartRms = ret.config.StartRms + 127
	ret.config.StopRms = ret.config.StopRms + 127

	fmt.Printf("[vad] NewAudioVadV2: %v\n", ret.config)
	return ret
}

func (vad *AudioVadV2) Release() {
	vad.startBuffer.clear()
	vad.stopBuffer.clear()
	vad.startBuffer = nil
	vad.stopBuffer = nil
}

func (vad *AudioVadV2) isActive(frame *AudioFrame) bool {
	voiceProb := 0
	rmsProb := 0
	if vad.isSpeaking {
		voiceProb = vad.config.StopVoiceProb
		rmsProb = vad.config.StopRms
	} else {
		voiceProb = vad.config.StartVoiceProb
		rmsProb = vad.config.StartRms
	}

	active := frame.FarFieldFlag == 1 && frame.VoiceProb > voiceProb && frame.Rms > rmsProb
	// date: 2025-10-29 for sdk which support apm filter, and in this version, we don't need to use farfield flag to detect speech, so we don't need to use farfield flag to detect speech
	refRms := vad.getRefActiveAvgRms()
	active = (frame.VoiceProb == 1) || (frame.Rms > refRms)
	//fmt.Printf("[vad] voiceProb: %d, rms: %d, pitch: %d, refRms: %d, active: %v\n", frame.VoiceProb, frame.Rms, frame.Pitch, refRms, active)
	return active
}
func (vad *AudioVadV2) getRefActiveAvgRms() int {
	if vad.config.EnableAdaptiveRmsThreshold && vad.refAvgRmsInLastSesseion > 0 {
		return int(float32(vad.refAvgRmsInLastSesseion)*vad.config.AdaptiveRmsThresholdFactor) // half of avg as threshold value
	}
	// convert from db to rms
	return vad.config.StartRms
}

func (vad *AudioVadV2) Process(frame *AudioFrame) (*AudioFrame, VadState) {
	if vad.expectFormat == nil {
		vad.expectFormat = &VadFrameFormat{
			BytesPerSample: frame.BytesPerSample,
			Channels:       frame.Channels,
			SamplesPerSec:  frame.SamplesPerSec,
		}
	} else {
		if vad.expectFormat.BytesPerSample != frame.BytesPerSample ||
			vad.expectFormat.Channels != frame.Channels ||
			vad.expectFormat.SamplesPerSec != frame.SamplesPerSec {
			return nil, VadStateNoSpeeking
		}
	}
	// if vad.isSpeaking {
	// 	fmt.Printf("[vad] +++++++++++++++++\n")
	// } else {
	// 	fmt.Printf("[vad] -----------------\n")
	// }
	isActive := vad.isActive(frame)
	vadFrame := newVadFrame(frame, isActive)
	if !vad.isSpeaking {
		full := vad.startBuffer.pushBack(vadFrame)
		if isActive {
			vad.totalVoiceRms += frame.Rms
			vad.voiceCount++
		} else {
			vad.voiceCount = 0
			vad.totalVoiceRms = 0
		}
		// todo： 是否需要根据N个isActive来计算avg rms? 而不是开始的时候，直接计算avg rms?
		// fmt.Printf("[vad] isSpeaking: false, startBuffer: %d\n", vad.startBuffer.queue.Len())
		if full {
			//activePercent = float32(vad.voiceCount) / float32(vad.config.StartRecognizeCount)
			//todo: disable activePercent check, use voiceCount instead
			if vad.voiceCount >= vad.config.StartRecognizeCount {
				vad.isSpeaking = true
				vad.stopBuffer.clear()
				ret := vad.startBuffer.flushAudio()

				// update ref rms
				vad.refAvgRmsInLastSesseion = vad.totalVoiceRms / vad.voiceCount
				fmt.Printf("[vad] StartSpeeking: %d, %d, %d, %d\n", vad.voiceCount, vad.totalVoiceRms, vad.refAvgRmsInLastSesseion,vad.config.StartRecognizeCount)

				//reset 
				vad.voiceCount = 0
				vad.totalVoiceRms = 0
				vad.silenceCount = 0

				// return the frame, and the state
				return ret, VadStateStartSpeeking
			}
		}
		return nil, VadStateNoSpeeking
	} else {
		full := vad.stopBuffer.pushBack(vadFrame)
		// fmt.Printf("[vad] isSpeaking: true, stopBuffer: %d\n", vad.stopBuffer.queue.Len())
		if isActive {
			vad.silenceCount = 0
		} else {
			vad.silenceCount++
		}
		if full {
			//inactivePercent := float32(vad.silenceCount) / float32(vad.config.StopRecognizeCount)
			if (vad.silenceCount >= vad.config.StopRecognizeCount) {
				vad.isSpeaking = false
				vad.stopBuffer.clear()

				fmt.Printf("[vad] StopSpeeking: %d, %d, %d, %d\n", vad.silenceCount, vad.totalVoiceRms, vad.refAvgRmsInLastSesseion,vad.config.StopRecognizeCount)

				// reset 
				vad.voiceCount = 0
				vad.totalVoiceRms = 0
				vad.silenceCount = 0

				//return
				return frame, VadStateStopSpeeking
			}
		}
		return frame, VadStateSpeeking
	}
}
