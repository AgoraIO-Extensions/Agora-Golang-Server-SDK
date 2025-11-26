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
	"encoding/binary"
	"unsafe"
	"sync/atomic"
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
	isDirectMode bool  // default to false, if true, means the pcmSender is in direct mode, and the data will be sent directly to the rtc channel
	directDataLen int // the length of the data in the direct mode
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
		directDataLen: 0,
	}
	

	fmt.Printf("NewAudioConsumer, audioScenario: %d, isDirectMode: %d\n", pcmSender.audioScenario, consumer.IsDirectMode())

	return consumer
}
func (ac *AudioConsumer) IsDirectMode() bool {
	return ac.pcmSender.audioScenario == AudioScenarioAiServer
}
// 暂时不处理非10ms整数的情况
func (ac *AudioConsumer) directPush(data []byte) {
	ac.frame.Buffer = data
	actualPackets := len(data) / ac.bytesPerFrame

	ac.frame.SamplesPerChannel = ac.samplesPerChannel * actualPackets
	ac.mu.Lock()
	ac.directDataLen += len(data) // not equal to the actual packets, because the data may not be 10ms aligned	
	ac.mu.Unlock()

	//fmt.Printf("directPush data len: %d, directDataLen: %d\n", len(data),ac.directDataLen)
	
	ac.pcmSender.SendAudioPcmData(ac.frame)
}


// PushPCMData adds PCM data to the buffer
func (ac *AudioConsumer) undirectPush(data []byte) {

	ac.mu.Lock()
	defer ac.mu.Unlock()
	ac.buffer.Write(data)
}

func (ac *AudioConsumer) PushPCMData(data []byte) {
	//fmt.Printf("PushPCMData data len: %d, isInitialized: %d, isDirectMode: %d\n", len(data), ac.isInitialized, ac.IsDirectMode())
	if !ac.isInitialized || data == nil || len(data) == 0 {
		return
	}
	if ac.IsDirectMode() {
		ac.directPush(data)
	} else {
		ac.undirectPush(data)
	}	
}

// reset resets the consumer's timing state
func (ac *AudioConsumer) reset() {
	if !ac.isInitialized {
		return
	}

	ac.startTime = (time.Now().UnixMilli())
	ac.consumedPackets = 0
	ac.lastConsumeedTime = ac.startTime
	ac.directDataLen = 0
}

func (ac *AudioConsumer) calculateCurWantPackets() int {
	
	now := time.Now().UnixMilli()
	elapsedTime := now - ac.startTime
	expectedTotalPackets := int(elapsedTime / 10) //change type from
	toBeSentPackets := expectedTotalPackets - ac.consumedPackets

	dataLen := 0
	if ac.IsDirectMode() {
		dataLen = ac.directDataLen
	} else {
		dataLen = ac.buffer.Len()
	}

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

	return actualPackets
}

	

// Consume processes and sends audio data
func (ac *AudioConsumer) undirectConsume() int {

	// Calculate actual packets to send
	actualPackets := ac.calculateCurWantPackets()
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

func (ac *AudioConsumer) directConsume() int {
	
	actualPackets := ac.calculateCurWantPackets()
	if actualPackets < 1 {
		return -3
	}

	// Prepare and send frame
	ac.mu.Lock()
	bytesToSend := ac.bytesPerFrame * actualPackets
	ac.directDataLen -= bytesToSend
	defer ac.mu.Unlock()

	

	if bytesToSend > 0 {
		

		ac.consumedPackets += actualPackets

		
		return 0
	}

	return -5
}

func (ac *AudioConsumer) Consume() int {
	if !ac.isInitialized {
		return -1
	}
	if ac.IsDirectMode() {
		return ac.directConsume()
	} else {
		return ac.undirectConsume()
	}
}

// Len returns the current buffer length
func (ac *AudioConsumer) Len() int {
	if ac.IsDirectMode() {
		return ac.directDataLen;
	} else {
		return ac.buffer.Len()
	}
}

// Clear empties the buffer
func (ac *AudioConsumer) Clear() {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	ac.directDataLen = 0
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
	remain_size := 0
	if ac.IsDirectMode() {
		remain_size = ac.directDataLen
	} else {
		remain_size = ac.buffer.Len()
	}
	
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
	ac.directDataLen = 0
	ac.isDirectMode = false
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
	// date: 2025-11-03, add voice prob, rms, pitch file for debug
	voiceProbFile *os.File
	rmsFile *os.File
	pitchFile *os.File
	// data for viceprob, rms, pitch file
	itemData []byte
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
		voiceProbFile: nil,
		rmsFile: nil,
		pitchFile: nil,
		// 960 bytes or 10ms is big enough for voice prob, rms, pitch data, i.e can support up to 48khz
		itemData: nil,

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
	//open voice prob file
	v.voiceProbFile, err = os.OpenFile(fmt.Sprintf("%s/voice_prob.txt", v.path), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to create voice prob file: ", err)
	}
	//open rms file
	v.rmsFile, err = os.OpenFile(fmt.Sprintf("%s/rms.txt", v.path), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to create rms file: ", err)
	}
	//open pitch file
	v.pitchFile, err = os.OpenFile(fmt.Sprintf("%s/pitch.txt", v.path), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to create pitch file: ", err)
	}

	return 0
}
func (v *VadDump) fillItemData(value int16, leninBytes int) {
	if v.itemData == nil || len(v.itemData) < leninBytes {
		v.itemData = make([]byte, leninBytes)
	}
	// chang byte to int16 ptr
	int16Ptr := (*int16)(unsafe.Pointer(&v.itemData[0]))
	count := leninBytes/2;
	for i := 0; i < count; i++ {
		elem := (*int16)(unsafe.Pointer(uintptr(unsafe.Pointer(int16Ptr)) + uintptr(i)*2))
		*elem = int16(value)
	}
}
func (v *VadDump) Write(frame *AudioFrame, vadFrame *AudioFrame, state VadState) int {
	//len := len(frame.Buffer)

	// write source
	if v.sourceFile != nil {
		if _, err := v.sourceFile.Write(frame.Buffer); err != nil {
			fmt.Println("Failed to write dump file: ", err)
		}
	}
	//write voice prob info: date: 2025-11-03, add voice prob info to file
	if v.voiceProbFile != nil {
		v.fillItemData(int16(frame.VoiceProb*127*127), len(frame.Buffer))
		v.voiceProbFile.Write(v.itemData)
	}
	//write rms info
	if v.rmsFile != nil {
		v.fillItemData(int16(frame.Rms*127), len(frame.Buffer))
		v.rmsFile.Write(v.itemData)
	}
	//write pitch info
	if v.pitchFile != nil {
		v.fillItemData(int16(frame.Pitch), len(frame.Buffer))
		v.pitchFile.Write(v.itemData)
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
	if v.voiceProbFile != nil {
		v.voiceProbFile.Close()
	}
	if v.rmsFile != nil {
		v.rmsFile.Close()
	}
	if v.pitchFile != nil {
		v.pitchFile.Close()
	}

	//assign to nil
	v.sourceFile = nil
	v.vadFile = nil
	v.labelFile = nil
	v.voiceProbFile = nil
	v.rmsFile = nil
	v.pitchFile = nil
	v.itemData = nil
	v.isOpen = false
	v.frameCount = 0
	v.count = 0
	return 0
}
//why use queue instead of chan?
// 1. chan is not thread-safe, so we need to use mutex to protect the queue;
// 2. chan is blocking, so we need to use select to avoid blocking;
// 3. chan is not flexible, so we need to use queue to replace chan;
// 4. chan is not scalable, so we need to use queue to replace chan;
// 5. chan is not easy to use, so we need to use queue to replace chan;
// 6. use interface{} to avoid type casting
// a thread-safe with timeout and non-blocking chan notify implementation

// user-defined queue: thread-safe with timeout and non-blocking chan notify implementation
// Queue 是一个线程安全的队列实现，可以替换chan；
// 优点：可以避免chan的阻塞问题，可以设置超时时间
// 推荐：如果用在音频的处理上，可以设置超时间为10ms；如果用在视频处理上，可以设置为20ms
// recommend: if used in audio processing, set timeout to 10ms; if used in video processing, set to 20ms
// sample: ref to send_recv_yuv_pcm.go
type Queue struct {
	items   []interface{}
	mutex   sync.Mutex
	timeout int // in ms
	notify  chan struct{}
}

// NewQueue 创建一个新的队列
func NewQueue(timeout int) *Queue {
	return &Queue{
		items:   make([]interface{}, 0),
		timeout: timeout,
		notify:  make(chan struct{}),
	}
}

// Enqueue 将元素添加到队列尾部
func (q *Queue) Enqueue(item interface{}) {
	q.mutex.Lock()
	q.items = append(q.items, item)
	q.mutex.Unlock()

	//q.items = append(q.items, item)
	// notify the dequeue routine: non-blocking mode
	select {
	case q.notify <- struct{}{}:
		//fmt.Println("notify the dequeue routine")
	default:
		//fmt.Println("--no notify the dequeue routine")
	}
}

// Dequeue 从队列头部移除并返回元素
func (q *Queue) Dequeue() interface{} {
	q.mutex.Lock()
	size := len(q.items)
	q.mutex.Unlock()

	//fmt.Printf("Dequeue size: %d\n", size)

	// if size is 0, wait notify signal or timeout
	if size == 0 {
		select {
		case <-time.After(time.Duration(q.timeout) * time.Millisecond):
			return nil
		case <-q.notify:
			// do nothing,just run to next step
		}
	}

	// get the signal item
	q.mutex.Lock()
	defer q.mutex.Unlock()
	size = len(q.items)
	if (size > 0) {
		item := q.items[0]
		q.items = q.items[1:]
		return item
	}
	return nil
}

// Peek 返回队列头部元素但不移除
func (q *Queue) Peek() interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.items) == 0 {
		return nil
	}

	return q.items[0]
}

// IsEmpty 检查队列是否为空
func (q *Queue) IsEmpty() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return len(q.items) == 0
}

// Size 返回队列中的元素个数
func (q *Queue) Size() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return len(q.items)
}

// Clear 清空队列
func (q *Queue) Clear() {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.items = make([]interface{}, 0)
}

// to generate wave header
// generate wave header, default to 16bit pcm data to wav file
// for 16bit pcm data to wav file
func GenerateWAVHeader(sampleRate int, channels int, pcmDataSizeInBytes int) []byte {
	totalSize := pcmDataSizeInBytes + 36 // 数据块+36字节头
	header := make([]byte, 44)

	// RIFF块
	copy(header[0:4], "RIFF")
	binary.LittleEndian.PutUint32(header[4:8], uint32(totalSize))
	copy(header[8:12], "WAVE")

	// fmt子块
	copy(header[12:16], "fmt ")
	binary.LittleEndian.PutUint32(header[16:20], 16)                               // fmt块大小
	binary.LittleEndian.PutUint16(header[20:22], 1)                                // PCM格式, fixed type
	binary.LittleEndian.PutUint16(header[22:24], uint16(channels))                 // 单声道
	binary.LittleEndian.PutUint32(header[24:28], uint32(sampleRate))               // 采样率
	binary.LittleEndian.PutUint32(header[28:32], uint32(sampleRate*channels*16/8)) // 字节率（16000 * 1 * 16/8）
	binary.LittleEndian.PutUint16(header[32:34], 2)                                // 块对齐（1 * 16/8）
	binary.LittleEndian.PutUint16(header[34:36], 16)                               // 位深度

	// data块
	copy(header[36:40], "data")
	binary.LittleEndian.PutUint32(header[40:44], uint32(pcmDataSizeInBytes))
	return header
}

type PrecisionTimer struct {
    ticker   *time.Ticker
    stopChan chan bool
}

func NewPrecisionTimer(interval int) *PrecisionTimer {
    return &PrecisionTimer{
        ticker:   time.NewTicker(time.Duration(interval) * time.Millisecond),
        stopChan: make(chan bool),
    }
}
func (t *PrecisionTimer) Start(taskFunc func()) {
    go func() {
        for {
            select {
            case <-t.ticker.C:
                taskFunc()
            case <-t.stopChan:
                t.ticker.Stop()
                return
            }
        }
    }()
}
func (pt *PrecisionTimer) Stop() {
    pt.ticker.Stop()
    close(pt.stopChan)
}
type IdleItem struct {
	Handle unsafe.Pointer
	LifeCycle int // in ms
}
func newIdleItem(handle unsafe.Pointer, lifeCycle int) *IdleItem {
    return &IdleItem{
        Handle: handle,
        LifeCycle: lifeCycle,
    }
}
/*
lock-freee ring buffer implementation
for single producer and single consumer
*/

// production level lock-free ring buffer
type LockFreeRingBuffer struct {
    buffer   []*AudioFrame
    mask     uint64  // capacity - 1, for bitwise operation instead of modulo
    
    // Cache line padding to avoid false sharing
    _pad0    [8]uint64
    writePos uint64  // only writer access
    _pad1    [8]uint64
    readPos  uint64  // only reader access
    _pad2    [8]uint64
}

// create ring buffer (capacity must be 2^n)
func NewLockFreeRingBuffer(capacity int) *LockFreeRingBuffer {
    // ensure capacity is power of 2
    if capacity&(capacity-1) != 0 {
        panic("capacity must be power of 2")
    }
    
    return &LockFreeRingBuffer{
        buffer: make([]*AudioFrame, capacity),
        mask:   uint64(capacity - 1),
    }
}

// write (in C callback, must be non-blocking)
func (rb *LockFreeRingBuffer) TryWrite(frame *AudioFrame) bool {
    // use relaxed semantic read (better performance)
    w := atomic.LoadUint64(&rb.writePos)
    r := atomic.LoadUint64(&rb.readPos)
    
    // use bitwise operation instead of modulo (performance提升 3-5 倍)
    if (w+1)&rb.mask == r&rb.mask {
        return false // full
    }
    
    // write data
    rb.buffer[w&rb.mask] = frame
    
    // Release semantic ensure write visibility
    atomic.StoreUint64(&rb.writePos, w+1)
    return true
}

// batch write (optional, better performance)
func (rb *LockFreeRingBuffer) TryWriteBatch(frames []*AudioFrame) int {
    written := 0
    for _, frame := range frames {
        if !rb.TryWrite(frame) {
            break
        }
        written++
    }
    return written
}

// read (in Go side)
func (rb *LockFreeRingBuffer) TryRead() (*AudioFrame, bool) {
    r := atomic.LoadUint64(&rb.readPos)
    w := atomic.LoadUint64(&rb.writePos)
    
    if r == w {
        return nil, false // empty
    }
    
    frame := rb.buffer[r&rb.mask]
    atomic.StoreUint64(&rb.readPos, r+1)
    return frame, true
}

// batch read (better performance)
func (rb *LockFreeRingBuffer) TryReadBatch(maxFrames int) []*AudioFrame {
    frames := make([]*AudioFrame, 0, maxFrames)
    
    for i := 0; i < maxFrames; i++ {
        if frame, ok := rb.TryRead(); ok {
            frames = append(frames, frame)
        } else {
            break
        }
    }
    
    return frames
}

// get current data size
func (rb *LockFreeRingBuffer) Size() int {
    w := atomic.LoadUint64(&rb.writePos)
    r := atomic.LoadUint64(&rb.readPos)
    return int((w - r) & rb.mask)
}

// check if empty
func (rb *LockFreeRingBuffer) IsEmpty() bool {
    w := atomic.LoadUint64(&rb.writePos)
    r := atomic.LoadUint64(&rb.readPos)
    return w == r
}

// capacity
func (rb *LockFreeRingBuffer) Capacity() int {
    return len(rb.buffer)
}