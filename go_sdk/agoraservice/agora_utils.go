package agoraservice

// #cgo CFLAGS: -I${SRCDIR}/../../agora_sdk/include/c/api2 -I${SRCDIR}/../../agora_sdk/include/c/base
// #include <stdlib.h>
// #include "agora_parameter.h"
import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"
)

// AudioConsumer provides utility functions for the Agora SDK.
// const def 
const (
	MinPacketsToSend = 10 // change min packets to send from 18 to 10
	RtcE2EDelay = 200 //e2e delay 90ms for iphone, 120ms for android;150ms for web. so we use 200ms here
)


// AudioConsumer handles PCM data consumption and sending
type AudioConsumer struct {
	mu              sync.Mutex
	startTime       int64 // in ms
	buffer          *bytes.Buffer
	consumedPackets int
	pcmSender       *AudioPcmDataSender
	frame           *AudioFrame

	// Audio parameters
	bytesPerFrame     int
	samplesPerChannel int

	// State
	isInitialized bool
	lastConsumeedTime int64 // in ms
}

// NewAudioConsumer creates a new AudioConsumer instance
func NewAudioConsumer(pcmSender *AudioPcmDataSender, sampleRate int, channels int) *AudioConsumer {
	if pcmSender == nil {
		return nil
	}

	bytesPerFrame := (sampleRate / 100) * channels * 2 // 2 bytes per sample

	consumer := &AudioConsumer{
		buffer:    bytes.NewBuffer(make([]byte, 0, bytesPerFrame*20)), // Pre-allocate buffer
		pcmSender: pcmSender,
		frame: &AudioFrame{
			SamplesPerSec:  sampleRate,
			Channels:       channels,
			BytesPerSample: 2,
			Buffer:         make([]byte, 0, bytesPerFrame), // Pre-allocate frame buffer
		},
		bytesPerFrame:     bytesPerFrame,
		samplesPerChannel: sampleRate / 100 ,
		isInitialized:     true,
		lastConsumeedTime: 0,
	}

	return consumer
}

// PushPCMData adds PCM data to the buffer
func (ac *AudioConsumer) PushPCMData(data []byte) {
	if !ac.isInitialized || data == nil || len(data) == 0 {
		return
	}

	ac.mu.Lock()
	defer ac.mu.Unlock()
	ac.buffer.Write(data)
}

// reset resets the consumer's timing state
func (ac *AudioConsumer) reset() {
	if !ac.isInitialized {
		return
	}

	ac.startTime = (time.Now().UnixMilli())
	ac.consumedPackets = 0
	ac.lastConsumeedTime = ac.startTime
}

// Consume processes and sends audio data
func (ac *AudioConsumer) Consume() int {
	if !ac.isInitialized {
		return -1
	}

	now := time.Now().UnixMilli()
	elapsedTime := now - ac.startTime
	expectedTotalPackets := int(elapsedTime / 10) //change type from
	toBeSentPackets := expectedTotalPackets - ac.consumedPackets

	dataLen := ac.buffer.Len()
	if dataLen > 0 {
	    ac.lastConsumeedTime = now
	}

	// Handle underflow
	if toBeSentPackets > MinPacketsToSend && dataLen/ac.bytesPerFrame < MinPacketsToSend {
		return -2 // should wait for more data
	}

	// Reset state if necessary
	if toBeSentPackets > MinPacketsToSend {
		ac.reset()
		toBeSentPackets = min(MinPacketsToSend, dataLen/ac.bytesPerFrame)
		ac.consumedPackets = (-toBeSentPackets)
	}

	// Calculate actual packets to send
	actualPackets := min(toBeSentPackets, dataLen/ac.bytesPerFrame)
	if actualPackets < 1 {
		return -3
	}

	// Prepare and send frame
	bytesToSend := ac.bytesPerFrame * actualPackets

	ac.mu.Lock()
	frameData := make([]byte, bytesToSend)
	n, _ := ac.buffer.Read(frameData)
	ac.mu.Unlock()

	if n > 0 {
		ac.frame.Buffer = frameData[:n]

		ac.frame.SamplesPerChannel = ac.samplesPerChannel * actualPackets

		ac.consumedPackets += actualPackets

		ret := ac.pcmSender.SendAudioPcmData(ac.frame)
		if ret == 0 {
			ret = ac.consumedPackets // return the actual consumed packets in this round
		}
		// for quickly release buffer
		ac.frame.Buffer = nil
		frameData	= nil
		return ret
	}

	return -5
}

// Len returns the current buffer length
func (ac *AudioConsumer) Len() int {
	return ac.buffer.Len()
}

// Clear empties the buffer
func (ac *AudioConsumer) Clear() {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	ac.buffer.Reset()
}
/*判断AudioConsumer中的数据是否已经完全推送给了RTC 频道
//因为audioconsumer内部有一定的缓存机制，所以当get_remaining_data_size 返回是0的时候，还有数据没有推送给
rtc 频道。如果要判断数据是否完全推送给了rtc 频道，需要调用这个api来做判断。
return value：1--push to rtc completed, 0--push to rtc not completed   -1--error
*/
func (ac *AudioConsumer) IsPushToRtcCompleted() int {
	if !ac.isInitialized {
		return -1
	}
	// no need to lock, because this function is only called in the main thread
	
	remain_size := ac.buffer.Len()
	if remain_size == 0 {
		now := time.Now().UnixMilli()
		diff := now - ac.lastConsumeedTime
		if diff > MinPacketsToSend*10 + RtcE2EDelay {
			return 1
		}
	}
	return 0
}


// Release frees resources
func (ac *AudioConsumer) Release() {
	if !ac.isInitialized {
		return
	}

	ac.isInitialized = false

	ac.mu.Lock()
	defer ac.mu.Unlock()

	// Clear references to allow GC
	ac.buffer = nil
	ac.frame = nil
	ac.pcmSender = nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

/*
* VadDump provides utility functions for the Agora SDK to dump valid inforamtion for VAD result
* its helpful to debug the VAD result, but do not recommend to use it in production environment
* samle ref to : sample_vad.go
* usage:
* vadDump := NewVadDump("./dump_log_name")
* vadDump.Open()
* vadDump.Write(frame, vadFrame, state) // frame: the audio frame, vadFrame: the vad result, state: the vad state
* vadDump.Close()
 */
type VadDump struct {
	mode       int    // 0: dump source; 1 dump source + far; 2 dump source + far + pitch;3 dump source + far + pitch + rms
	path       string // dir path of the dump file, and the file name is path/vad/source_.pcm
	count      int    // current count of the vad section dump file
	sourceFile *os.File
	vadFile    *os.File
	labelFile  *os.File
	isOpen     bool
	frameCount int
}

func NewVadDump(path string) *VadDump {
	//check path is a dir or not, and is writable
	if info, err := os.Stat(path); os.IsNotExist(err) || !info.IsDir() || !info.Mode().IsDir() {
		fmt.Println("Path does not exist: ", path)
		return nil
	}
	// create sub dir according to YYYYMMDD
	now := time.Now()
	vadPath := fmt.Sprintf("%s/%04d%02d%02d%02d%02d%02d", path, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	if _, err := os.Stat(vadPath); os.IsNotExist(err) {

		//create sub dir if not exist
		if err := os.MkdirAll(vadPath, 0755); err != nil {
			fmt.Println("Failed to create dir: ", err)
			return nil
		}
	}
	ret := &VadDump{
		mode:       1,
		path:       vadPath,
		count:      0, // count of the vad section dump file
		sourceFile: nil,
		vadFile:    nil,
		labelFile:  nil,
		isOpen:     false,
		frameCount: 0, // count of audio frame
	}
	return ret
}

func (v *VadDump) newVadFile() *os.File {
	var err error
	v.vadFile, err = os.OpenFile(fmt.Sprintf("%s/vad_%d.pcm", v.path, v.count), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to create dump file: ", err)
		return nil
	}
	return v.vadFile
}
func (v *VadDump) closeVadFile() int {
	if v.vadFile != nil {
		v.vadFile.Close()
		v.vadFile = nil

		return 0
	}

	return -1
}

func (v *VadDump) Open() int {

	var err error

	if v.isOpen {
		return 1
	}
	v.isOpen = true
	v.count = 0
	v.frameCount = 0

	// open source file
	v.sourceFile, err = os.OpenFile(fmt.Sprintf("%s/source.pcm", v.path), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to create dump file: ", err)
	}

	//open label file
	v.labelFile, err = os.OpenFile(fmt.Sprintf("%s/label.txt", v.path), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to create label file: ", err)
	}

	return 0
}
func (v *VadDump) Write(frame *AudioFrame, vadFrame *AudioFrame, state VadState) int {
	//len := len(frame.Buffer)

	// write source
	if v.sourceFile != nil {
		if _, err := v.sourceFile.Write(frame.Buffer); err != nil {
			fmt.Println("Failed to write dump file: ", err)
		}
	}

	// write label info
	if v.labelFile != nil {
		info := fmt.Sprintf("ct:%d fct:%d state:%d far:%d vop:%d rms:%d pitch:%d mup:%d\n", v.count, v.frameCount, state,
			frame.FarFieldFlag, frame.VoiceProb, frame.Rms, frame.Pitch, frame.MusicProb)
		v.labelFile.WriteString(info)

		// increment frame count
		v.frameCount++

	}

	//check vad state
	if vadFrame != nil {

		if state == VadStateStartSpeeking {
			//open a new vad file
			v.vadFile = v.newVadFile()
			// write current bytes to file
			v.vadFile.Write(vadFrame.Buffer)
			// and increment count
			v.count++
		} else if state == VadStateSpeeking {
			// write current bytes to file
			v.vadFile.Write(vadFrame.Buffer)

		} else if state == VadStateStopSpeeking {
			// close current vad file
			v.closeVadFile()
		}
	}
	return 0
}

func (v *VadDump) Close() int {

	//close all file
	if v.sourceFile != nil {
		v.sourceFile.Close()
	}
	if v.vadFile != nil {
		v.vadFile.Close()
	}
	if v.labelFile != nil {
		v.labelFile.Close()

	}

	//assign to nil
	v.sourceFile = nil
	v.vadFile = nil
	v.labelFile = nil

	v.isOpen = false
	v.frameCount = 0
	v.count = 0
	return 0
}
