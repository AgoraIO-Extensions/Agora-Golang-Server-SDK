package agoraservice

import (
	"container/list"
)

type AudioVadConfigV2 struct {
	PreStartRecognizeCount int     // pre start recognize count, buffer size for 10ms 16KHz 16bit 1channel PCM, default value is 16
	StartRecognizeCount    int     // start recognize count, buffer size for 10ms 16KHz 16bit 1channel PCM, default value is 30
	StopRecognizeCount     int     // max recognize count, buffer size for 10ms 16KHz 16bit 1channel PCM, default value is 20
	ActivePercent          float32 // active percent, if over this percent, will be recognized as speaking, default value is 0.7
	InactivePercent        float32 // inactive percent, if below this percent, will be recognized as non-speaking, default value is 0.5
	StartVoiceProb         int     // start voice prob, default value is 70
	StartRms               int     // start rms, default value is -50
	StopVoiceProb          int     // stop voice prob, default value is 70
	StopRms                int     // stop rms, default value is -50
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
			StopRecognizeCount:     20,
			ActivePercent:          0.7,
			InactivePercent:        0.5,
			StartVoiceProb:         70,
			StartRms:               -50.0,
			StopVoiceProb:          70,
			StopRms:                -50.0,
		}
	}
	if cfg.StartRecognizeCount <= 0 {
		cfg.StartRecognizeCount = 10
	}
	if cfg.StopRecognizeCount <= 0 {
		cfg.StopRecognizeCount = 6
	}
	// fmt.Printf("[vad] NewAudioVadV2: %v\n", cfg)
	startQueueSize := cfg.StartRecognizeCount + cfg.PreStartRecognizeCount
	return &AudioVadV2{
		config:       cfg,
		expectFormat: nil,
		isSpeaking:   false,
		startBuffer:  newVadBuffer(startQueueSize),
		stopBuffer:   newVadBuffer(cfg.StopRecognizeCount),
	}
}

func (vad *AudioVadV2) Release() {
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
	// fmt.Printf("[vad] isActive: %v, isSpeaking: %v, FarFieldFlag: %d, voiceProb: %d, rms: %d, pitch: %d\n",
	// 	active, vad.isSpeaking, frame.FarFieldFlag, frame.VoiceProb, frame.Rms, frame.Pitch)
	return active
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
	vadFrame := newVadFrame(frame, vad.isActive(frame))
	if !vad.isSpeaking {
		full := vad.startBuffer.pushBack(vadFrame)
		// fmt.Printf("[vad] isSpeaking: false, startBuffer: %d\n", vad.startBuffer.queue.Len())
		if full {
			activePercent := vad.startBuffer.getActivePercent(vad.config.StartRecognizeCount)
			if activePercent >= vad.config.ActivePercent {
				vad.isSpeaking = true
				vad.stopBuffer.clear()
				ret := vad.startBuffer.flushAudio()
				return ret, VadStateStartSpeeking
			}
		}
		return nil, VadStateNoSpeeking
	} else {
		full := vad.stopBuffer.pushBack(vadFrame)
		// fmt.Printf("[vad] isSpeaking: true, stopBuffer: %d\n", vad.stopBuffer.queue.Len())
		if full {
			activePercent := vad.stopBuffer.getActivePercent(vad.config.StopRecognizeCount)
			if (1.0 - activePercent) >= vad.config.InactivePercent {
				vad.isSpeaking = false
				vad.stopBuffer.clear()
				return frame, VadStateStopSpeeking
			}
		}
		return frame, VadStateSpeeking
	}
}
