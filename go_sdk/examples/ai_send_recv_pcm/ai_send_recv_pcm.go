package main

import (
	//"bytes"
	//"bufio"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"time"
	"unsafe"

	agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/rtc"

	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
)
/*
date:2025-09-15
author: weihongqin
description：aiqos2.0 协议测试, V7
结果：
1、会在一个发送的audioframe包的pts基础上，每一个frame做自增10（对应10ms的数据）
2、前后2个frame的pts，可以无关联，不需要保障连续frame的pts 是也是连续或者自增的。
比如sendaudiodata（data=1s,pts=1000），那么下一个sendaudiodata（data=1s,pts=1000），是可以的。

分析：考虑到pts在接收端会从pts0 开始自增，为了不掩盖用户的协议，考虑到每一个chunk实际可能大小，设计协议如下：
低16位留给sdk：每次从0开始，也就是可以支持65536ms，也就是65.536s，足够了。一个chunk不可能这么大
|2位(agora or customer)|4位（version）｜6位（sessionid）｜12位（sentenceid）｜5位（chunkid）｜1位（isSessionEnd）|18位（reserved）|16位（basepts）
2位(agora or customer)：00表示agora，01表示customer
4位（version）：0-15
6位（sessionid）：0-63
12位（sentenceid）：0-4095
5位（chunkid）：0-31
1位（isSessionEnd）：0表示不结束，1表示结束
18位（reserved）：0-262143
16位（basepts）：0-65535

sessonend的时候，带上duration信息，来表示这个session的持续时间，另外的协议来做
sesson end可以额外的100ms 静音包



如何标识session结束？
在结束的时候，发送10个静音包。静音包里面的sessio/seentce/chunk都保持不变，但是isSessionEnd为true

发送3个而不是在协议里面设置chunknum的原因？是因为audioframe会丢失，对一个10%丢包率低网络来说，
3个包都接收到的概率是72.9%。如果是5个包，概率就会直线下降到：59%
但如果发送3个静音包，接收到任意一个静音包，就表示已经结束。其成功率为：99.9%

所以采用额外发送静音包的方式。

提议：
目的：
为了开始/结束/sentence/chunk的方式
打断那个sesison是需要放在frame里面的--》通过trackinfo来做，就是trackinfo里面带上最近的一个pts
需要验证非brust模式的情况？

basepts用getAgoraCurrentMonotonicTimeInMs来做；
然后每次调用用duration来增加

测试结果：
如果最高位是1， 接收到的pts都会是0
丢包率上行/下行设置为60%的时候，avc已经不能工作；但aiqos看起来还没有丢包。

V8:
面对的问题/目标：
1、session时间的通知：开始/结束
2、session的播放时长
3、chunk好像没有必要？？应该有必要
比如：https://cloud.tencent.com/document/product/1073/37995
有TTS的返回中，会带有subtitles，这个subtitles就是chunkid。
做字幕对齐的时候，可以通过rtm发送{chunkid:subtitles}这样的映射关系和subtitles信息给端侧，端侧则根据这个来做字幕的渲染
subtitles格式：
 "Subtitles": [
            {
                "BeginIndex": 0,
                "BeginTime": 250,
                "EndIndex": 1,
                "EndTime": 430,
                "Phoneme": "ni2",
                "Text": "你"
            },
            {
                "BeginIndex": 1,
                "BeginTime": 430,
                "EndIndex": 2,
                "EndTime": 670,
                "Phoneme": "hao3",
                "Text": "好"
            }
        ]

|2位(agora or customer)|4位（version）｜1位（isSessionEnd）｜6位（sessionid）｜12位（sentenceid）｜5位（chunkid）|18位（reserved）|16位（basepts）

如果isSessionEnd为true，则表示session结束，协议变更为：
|2位(agora or customer)|4位（version）｜1位（isSessionEnd）｜23位(session durationinms/10)|18位（reserved）|16位（basepts）
2位(agora or customer)：00表示agora，01表示customer
4位（version）：0-15



*/

func PushFileToConsumer(file *os.File, con *agoraservice.RtcConnection, samplerate int) {
	buffer := make([]byte, samplerate*2) // 1s data
	for {
		readLen, err := file.Read(buffer)
		if err != nil {
			fmt.Printf("read up to EOF,cur read: %d", readLen)
			file.Seek(0, 0)
			return
		}
		// round to integer of chunk
		packLen := readLen / samplerate
		con.PushAudioPcmData(buffer[:packLen*samplerate], samplerate, 1, 0)
		fmt.Printf("PushPCMData done:%d\n", readLen)
	}
}
func ReadFileToConsumer(file *os.File, con *agoraservice.RtcConnection, interval int, done chan bool, samplerate int) {
	for {
		select {
		case <-done:
			fmt.Println("ReadFileToConsumer done")
			return
		default:
			if con != nil {
				isPushCompleted := con.IsPushToRtcCompleted()
				if isPushCompleted {
					PushFileToConsumer(file, con, samplerate)
				}
			}
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	}
}

func LoopbackAudio(audioQueue *agoraservice.Queue, audioConsumer *agoraservice.AudioConsumer, done chan bool) {
	for {
		select {
		case <-done:
			fmt.Println("LoopbackAudio done")
			return
		default:
			AudioFrame := audioQueue.Dequeue()
			if AudioFrame != nil {
				//fmt.Printf("AudioFrame: %d\n", time.Now().UnixMilli())
				if frame, ok := AudioFrame.(*agoraservice.AudioFrame); ok {
					frame.RenderTimeMs = 0
					audioConsumer.PushPCMData(frame.Buffer)
				}
			}
			//time.Sleep(10 * time.Millisecond)
		}
	}
}

func SendTTSDataToClient(samplerate int, audioConsumer *agoraservice.AudioConsumer, file *os.File, done chan bool, audioSendEvent chan struct{}, fallackEvent chan struct{}, localUser *agoraservice.LocalUser, track *agoraservice.LocalAudioTrack) {
	for {
		select {
		case <-done:
			fmt.Println("SendAudioToClient done")
			return
		case <-fallackEvent:

		case <-audioSendEvent:
			// read 1s data from file
			buffer := make([]byte, samplerate*2*2) // 2s data
			readLen, err := file.Read(buffer)
			if err != nil {
				fmt.Printf("read up to EOF,cur read: %d", readLen)
				file.Seek(0, 0)
				continue
			}
			audioConsumer.PushPCMData(buffer[:readLen])
			// and seek to the begin of the file
			file.Seek(0, 0)
			fmt.Println("SendTTSDataToClient done")
		default:
			time.Sleep(40 * time.Millisecond)
			audioConsumer.Consume()
		}
	}
}

func calculateEnergyFast(data []byte) uint64 {
	var energy uint64
	samples := unsafe.Slice((*int16)(unsafe.Pointer(&data[0])), len(data)/2)
	for _, s := range samples {
		energy += uint64(s) * uint64(s)
	}
	return energy
}
func mixAudio(data1 []byte, data2 []byte) []byte {
	// check if the length of data1 and data2 is the same
	if len(data1) != len(data2) {
		return nil
	}
	int16len := len(data1)/2
	// allocate a new buffer
	buffer := make([]byte, len(data1))
	var ret int32 = 0
	
	// algorithm: Saturation Normalization
	sample1 := unsafe.Slice((*int16)(unsafe.Pointer(&data1[0])), int16len)
	sample2 := unsafe.Slice((*int16)(unsafe.Pointer(&data2[0])), int16len)
	dst := unsafe.Slice((*int16)(unsafe.Pointer(&buffer[0])), int16len)
	for i := 0; i < int16len; i++ {
		ret = int32(sample1[i]) + int32(sample2[i])
		if ret > 32767 {
			dst[i] = 32767
		} else if ret < -32768 {
			dst[i] = -32768
		} else {
			dst[i] = int16(ret)
		}
	}
	return buffer
}
func NonblockNotiyEvent(event chan struct{}) int {
	select {
	case event <- struct{}{}:
		return 0
	default:
		return -1
	}
}

/*
date:2025-08-14
author: weihongqin
description: nonblock notiy event
param: event: event channel
return: 0 if success, -1 if failed

mode: 0 loop send file mode
mode： 1 echo mode, just echo back rmeote user 's audio
mode： 2是根据audio meta 信息来触发发送文件和打断
mode：3 是用来做方波信号回环延迟的测试
mode：4， 也是根据audio meta 来做chuany的测试，也就是：
收到audiometa，做打断；在收到后，发送一个新的句子。并且用更新pts
*/
/*
// v2 版本：
改动1: 引入chunk的概念，模拟每调用一次sendpcm的api，就自增加一次。
改动2: 将basepts从int32修改为uint16,并且默认位0
改动3: 不在需要用确保同一个句子分段调用api的时候，需要人工计算basepts的麻烦
设计依据：
1分钟是6Wms，也就是：0xEA60，对24000的音频来说，是2MB数据，这个对一个chunk来说，是完全足够的。因此pts用uint16来表示
限制：
版本最多支持7个版本：0x0-0x7
支持的session次数是65536（一轮对话最少1s，实际远大于。就可以支持65536s的对话周期，也就是18H的对话，足够）
一个session的sentence最多是65536次
一个sentence最多有1023次chunk
一个chunk有6Wms的语音：也就是1分钟的语音，2MB的数据

位分布：高4位（version:不能超过0x7)|中间16位（sessionid)|中间16位（sentenceid)|低10位(chunkid）｜2位（session是否结束的标记)｜低16位（basepts）
*/
// CombineToInt64 合并六个字段为 int64
// 位分布：高4位(version) | 16位(sessionid) | 16位(sentenceid) | 10位(chunkid) | 2位(isSessionEnd) | 16位(basepts)

/*
V4: from client to servereV4: from client to servere
// 位分布：高3位(version) | 18位(sessionid) |  10位(last_chunk_durationinms)|1位(isSessionEnd) | 32位(basepts)
*/

func CombineToInt64(version int8, sessionid uint16, sentenceid uint16, chunkid uint16, issessionend bool) int64 {
	// 1. 验证 version 范围（0-7）
	if version < 0 || version > 0x7 {
		panic("version must be between 0 and 7")
	}
	versionPart := int64(version) << 60 // 高4位（bits 60-63）

	// 2. sessionid（16位，bits 44-59）
	sessionPart := int64(sessionid) << 44

	// 3. sentenceid（16位，bits 28-43）
	sentencePart := int64(sentenceid) << 28

	// 4. chunkid（10位，bits 18-27）
	chunkMask := uint16(0x3FF) // 二进制 0011 1111 1111（10位掩码）
	chunkPart := int64(chunkid&chunkMask) << 18

	// 5. issessionend（2位，bits 16-17）
	isend := int8(0)
	if issessionend {
		isend = 1
	}
	endPart := int64(isend) << 16

	// 6. basepts（16位，bits 0-15）
	basepts := uint16(0)                   // 默认是0,不需要用户设置
	baseptsPart := int64(basepts) & 0xFFFF // 确保仅保留低16位

	// 合并所有部分（按位或）
	return versionPart | sessionPart | sentencePart | chunkPart | endPart | baseptsPart
}

type SessionInfo struct {
	version                 uint8
	sessionid               uint32
	lastChunkDuration       uint16
	isSessionEnd            bool
	basepts                 uint32
	recvedLastChunkDuration uint16
}
type SessionParser struct {
	sessionInfo      SessionInfo
	lastCallbackTime int64
	lastSessionId    uint32
	exitChan         chan struct{}
	callbackFunc     func(sessionId uint32, isFirstOrLastFrame bool) // true: frist frame; false: last frame
	isInited         bool
	isStop           bool
}
type SessionEndReason int

const (
	SessionEndReasonLastChunk SessionEndReason = 0
	SessionEndReasonTimeout   SessionEndReason = 1
	SessionEndReasonInterrupt SessionEndReason = 2
	SessionEndReasonOther     SessionEndReason = 3
)

// isFist: 1-- 开始；0--结束
func NewSessionParser(callbackFunc func(sessionId uint32, isFirstOrLastFrame bool)) *SessionParser {
	return &SessionParser{
		sessionInfo:      SessionInfo{},
		lastCallbackTime: 0,
		lastSessionId:    0,
		exitChan:         make(chan struct{}),
		callbackFunc:     callbackFunc,
		isInited:         false,
		isStop:           true,
	}
}

func (sp *SessionParser) Start() int {
	// start a timer.Ticker to check last callback time
	if sp.isInited {
		return 0
	}
	sp.isInited = true
	sp.isStop = false
	go func() {
		var interval int = 100
		diff := int64(200)
		ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-sp.exitChan:
				return
			case <-ticker.C:
				//diff : 75~125ms to do timedout
				if sp.sessionInfo.recvedLastChunkDuration > 0 {
					diff = int64(200)
				} else {
					diff = int64(500)
				}
				if time.Now().UnixMilli()-sp.lastCallbackTime > diff {
					// callback to notiy current sesion is end, and a new sesion is coming
					sp.doEnd(sp.sessionInfo.sessionid, SessionEndReasonTimeout)
				}
			}

		}
	}()
	return 0
}
func (sp *SessionParser) End() int {
	if sp.isStop {
		return 0
	}
	sp.isStop = true
	sp.exitChan <- struct{}{}
	return 0
}
func (sp *SessionParser) parseInt64V4(value int64) (version uint8, sessionid uint32, lastChunkDuration uint16, isSessionEnd int, basepts uint32) {

	// 1. 提取高3位 version (bits 61-63)
	version = uint8(value>>61) & 0x07 // 0x07 = 二进制 0111（3位掩码）

	// 2. 提取18位 sessionid (bits 43-60)
	sessionid = uint32((value >> 43) & 0x3FFFF) // 0x3FFFF = 二进制 0011 1111 1111 1111 1111（18位掩码）

	// 3. 提取10位 last_chunk_durationinms (bits 33-42)
	lastChunkDuration = uint16((value >> 33) & 0x3FF) // 0x3FF = 二进制 0000 0011 1111 1111（10位掩码）

	// 4. 提取1位 isSessionEnd (bit 32)
	isSessionEnd = int((value >> 32) & 0x01) // 0x01 = 二进制 0000 0001（1位掩码）

	// 5. 提取低32位 basepts (bits 0-31)
	basepts = uint32(value & 0xFFFFFFFF) // 0xFFFFFFFF = 二进制 1111 1111 1111 1111 1111 1111 1111 1111（32位掩码）

	return
}

func (sp *SessionParser) Parse(value int64) {
	version, sessionid, lastChunkDuration, isSessionEnd, basepts := sp.parseInt64V4(value)
	//compare cur frame's sesison id to session's id
	// 需要判断是否是当前新的sesion：如果是第一次，就触发doFist
	// 如果有sesion不一致，就触发doEnd,并且assign 给当前的sessionid
	if sessionid != sp.sessionInfo.sessionid {
		//判断是否是结束？
		// 还是说是开始新的session？
		if sp.sessionInfo.sessionid != 0 {
			sp.doEnd(sp.sessionInfo.sessionid, SessionEndReasonInterrupt)
		}
		sp.doFirst(sessionid)
	}

	//
	sp.sessionInfo.version = version
	sp.sessionInfo.sessionid = sessionid
	sp.sessionInfo.lastChunkDuration = lastChunkDuration
	sp.sessionInfo.isSessionEnd = isSessionEnd == 1
	sp.sessionInfo.basepts = basepts
	sp.lastCallbackTime = time.Now().UnixMilli()
	//check if the session is end
	if sp.sessionInfo.isSessionEnd {
		// update lastcyunk
		sp.sessionInfo.recvedLastChunkDuration += 10
		// for frame with isSessionEnd is true, the lastChunkValue is same
		// but under poor net, the value maybe changed. so only allow to become biger
		if lastChunkDuration > sp.sessionInfo.lastChunkDuration {
			sp.sessionInfo.lastChunkDuration = lastChunkDuration
		}

		if sp.sessionInfo.recvedLastChunkDuration >= sp.sessionInfo.lastChunkDuration {
			//call back:
			// and reset
			sp.doEnd(sp.sessionInfo.sessionid, SessionEndReasonLastChunk)
		}
	}
}

func (sp *SessionParser) doEnd(sessionId uint32, reason SessionEndReason) {

	//fmt.Printf("doEnd, sessionId: %d, reason: %d, lastSessionId: %d, currentSessionId: %d\n", sessionId, reason, sp.lastSessionId, sp.sessionInfo.sessionid)

	if sp.callbackFunc != nil && sp.sessionInfo.sessionid != 0 && sp.lastSessionId != sp.sessionInfo.sessionid {
		sp.callbackFunc(sessionId, false)
	}
	// set to current session id
	sp.lastSessionId = sp.sessionInfo.sessionid

	//empty sesion
	sp.sessionInfo.lastChunkDuration = 0
	sp.sessionInfo.isSessionEnd = false
	sp.sessionInfo.basepts = 0
	sp.sessionInfo.recvedLastChunkDuration = 0
	sp.lastCallbackTime = 0

}
func (sp *SessionParser) doFirst(sessionid uint32) {
	if sp.callbackFunc != nil {
		sp.callbackFunc(sessionid, true)
	}
}

/*
date:2025-08-21
author: weihongqin
description: byte manager for pcm raw data
*/

// ByteManager: the manager for pcm raw data
type PcmRawDataManager struct {
	data       []byte
	start      int // the start position of the buffer
	lock       sync.Mutex
	sampleRate int // the sample rate of the pcm data
	channels   int // the channels of the pcm data
	bytesInMs  int // the bytes in ms of the pcm data
}

// NewByteManager: create a new byte manager
func NewPcmRawDataManager(sampleRate int, channels int) *PcmRawDataManager {
	bm := &PcmRawDataManager{
		data:       make([]byte, 0),
		sampleRate: 24000,
		channels:   1,
		bytesInMs:  0,
	}
	bm.sampleRate = sampleRate
	bm.channels = channels
	bm.bytesInMs = sampleRate * 2 / 1000 * channels
	return bm
}

// UpdateParameters: update the parameters of the byte manager
// return 0: success
// return 1: if the parameters are the same as the current parameters, do nothing
func (bm *PcmRawDataManager) UpdateParameters(sampleRate int, channels int) int {
	bm.lock.Lock()
	defer bm.lock.Unlock()
	if sampleRate != bm.sampleRate || channels != bm.channels {
		bm.sampleRate = sampleRate
		bm.channels = channels
		bm.bytesInMs = sampleRate * 2 / 1000 * channels
		// and reset the data
		bm.Reset()
		return 0
	}
	return 1
}

// Push: push data to the end of the buffer without copy
// always push the data to the end of the buffer
func (bm *PcmRawDataManager) Push(data []byte) {
	bm.lock.Lock()
	defer bm.lock.Unlock()
	bm.data = append(bm.data, data...)
}

// Pop: pop data from the begin of the buffer without copy
// the returned slice is the view of the internal data, and the caller should not modify it
// if the available data is not round to the nearest multiple of bytesInMs, return nil
func (bm *PcmRawDataManager) Pop() []byte {
	bm.lock.Lock()
	defer bm.lock.Unlock()

	available := len(bm.data) - bm.start
	// round to the nearest multiple of bytesInMs
	length := available / bm.bytesInMs * bm.bytesInMs
	// if length is 0, return empty slice
	if length == 0 {
		return nil
	}

	// return the data view of the requested length
	result := bm.data[bm.start : bm.start+length]
	bm.start += length

	// 如果已弹出超过一半的数据，进行内存整理
	if bm.start > len(bm.data)/2 {
		bm.data = bm.data[bm.start:]
		bm.start = 0
	}

	return result
}

// Len : return the data length in the buffer
func (bm *PcmRawDataManager) Len() int {
	bm.lock.Lock()
	defer bm.lock.Unlock()
	return len(bm.data) - bm.start
}

// Reset: reset the buffer or empty the buffer
// and return the remain data length
func (bm *PcmRawDataManager) Reset() int {
	bm.lock.Lock()
	defer bm.lock.Unlock()
	remain := len(bm.data) - bm.start
	bm.data = bm.data[:0]
	bm.start = 0
	return remain
}

/*
date:2025-08-21
author: weihongqin
description: PTS allocator for V5 protocol, to simplify the PTS calculation
v5:
// 位分布：高3位(version) | 12位(sessionid) | 16位(last_chunk_durationinms) | 1位(isSessionEnd) | 32位(basepts)
*/
type PTSAllocator struct {
	version    int16
	sessionId  uint16
	basePts    uint32
	sampleRate int
	bytesInMs  int
}

func NewPTSAllocatorV5(sampleRate int) *PTSAllocator {
	pa := &PTSAllocator{
		version:   0,
		sessionId: 1, // MUST start from 1, and limitation is 12 bits
		basePts:   0,
	}
	pa.sampleRate = sampleRate
	pa.bytesInMs = sampleRate * 2 / 1000
	return pa
}

func (pa *PTSAllocator) Allocate(curPushDataLen int, isSessionEnd bool) int64 {
	// combination to 64 bits:
	curDuration := curPushDataLen / pa.bytesInMs
	curPts := pa.basePts
	lastFrameNumber := uint16(0)

	if isSessionEnd {
		lastFrameNumber = uint16(curDuration/10) // 10ms per frame
	}

	combined := pa.combineToInt64V5(pa.version, pa.sessionId, lastFrameNumber, isSessionEnd, curPts)

	// increase the base pts
	pa.basePts += uint32(curDuration)

	// 1. check if the session is end
	if isSessionEnd {
		pa.sessionId++
		// limitation check: should be less than 12 bits
		if pa.sessionId > 0xFFF {
			pa.sessionId = 1
		}

		pa.basePts = 0
	}

	return combined
}


// CombineToInt64V4 合并字段到int64
// V5: 
// 位分布：高14位(version) | 10位(sessionid) | 16位(last_frame_number) | 1位(isSessionEnd) | 21位(basepts)

func (pa *PTSAllocator) combineToInt64V5(version int16, sessionid uint16, last_frame_number uint16, isend bool, basepts uint32) int64 {
	// 参数校验
	if version < -8192 || version > 8191 { // int16范围，但只使用14位（有符号：-8192~8191）
		panic("version exceeds 14 bits")
	}
	if sessionid > 0x3FF { // 10位最大值1023 (0x3FF)
		panic("sessionid exceeds 10 bits")
	}
	if last_frame_number > 0xFFFF { // 16位最大值65535
		panic("last_duration exceeds 16 bits")
	}
	if basepts > 0x1FFFFF { // 21位最大值2097151 (0x1FFFFF)
		panic("basepts exceeds 21 bits")
	}

	// 合并字段（注意处理version的符号位）
	var result int64

	// version (14位，左移50位)
	// 将int16转为无符号14位：version & 0x3FFF
	result |= int64(uint16(version)&0x3FFF) << 50

	// sessionid (10位，左移40位)
	result |= int64(sessionid&0x3FF) << 40

	// last_duration (16位，左移24位)
	result |= int64(last_frame_number) << 24

	// isend (1位，左移23位)
	if isend {
		result |= 1 << 23
	}

	// basepts (21位，最低位)
	result |= int64(basepts & 0x1FFFFF)

	return result
}
/*
date:2025-08-22
author: weihongqin
description: V6: session结束的时候，固定是发送4个包，在server内部自己实现
V6:note: v6 is not work
/ V6: session结束的时候，固定是发送4个包，在server内部自己实现。也就是说当开发者发送isSessionEnd为true的时候，server内部会自动发送4个静音包，来标记session结束。
// 规则：sesion结束的时候，固定是发送4个静音包，在server内部自己实现！或者开发者自己来实现
// 位分布：高16位(version) | 4位(sessionid) | 11位(sentenceid) |1位(isSessionEnd) | 32位(basepts)
不标记sentenct是否结束：是因为在一个session内的sentence，是连续的，所以不需要标记。通过后一个sentenceid的变化来通知
session结束：通过isSessionEnd来标记，当isSessionEnd为true的时候，server内部会自动发送4个静音包，来标记session结束。
用来解决的场景：通常当作是一个Flag，允许用户自己来传递，可以设计为UINT16的值，来传递。
1、用来做字幕对齐
2、用来做句子结束的标记
3、用来做session结束的标记
4、用来做其他事件的标记
需要做的验证：
*/
// CombineToInt64V6合并字段到int64
// 位分布：高16位(version) | 4位(sessionid) | 11位(sentenceid) | 1位(isSessionEnd) | 32位(basepts)
func CombineToInt64V6(version int16, sessionid uint16, sentenceid uint16, isend bool, basepts uint32) int64 {
	// 参数校验
	if sessionid > 0xF { // 4位最大值15 (0xF)
		panic("sessionid exceeds 4 bits")
	}
	if sentenceid > 0x7FF { // 11位最大值2047 (0x7FF)
		panic("sentenceid exceeds 11 bits")
	}
	if basepts > 0xFFFFFFFF { // 32位最大值
		panic("basepts exceeds 32 bits")
	}

	// 合并字段（注意处理version的符号位）
	var result int64

	// version (16位，左移48位)
	// 直接将int16转为uint16保留原始位模式
	result |= int64(uint16(version)) << 48

	// sessionid (4位，左移44位)
	result |= int64(sessionid&0xF) << 44

	// sentenceid (11位，左移33位)
	result |= int64(sentenceid&0x7FF) << 33

	// isend (1位，左移32位)
	if isend {
		result |= 1 << 32
	}

	// basepts (32位，最低位)
	result |= int64(basepts & 0xFFFFFFFF)

	return result
}

// ParseInt64V4 从int64解析字段
func ParseInt64V6(value int64) (version int16, sessionid uint16, sentenceid uint16, isend bool, basepts uint32) {
	// version (16位，先取uint16再转int16保留原始位模式)
	version = int16(uint16(value >> 48))

	// sessionid (4位)
	sessionid = uint16((value >> 44) & 0xF)

	// sentenceid (11位)
	sentenceid = uint16((value >> 33) & 0x7FF)

	// isend (1位)
	isend = (value>>32)&0x1 != 0

	// basepts (32位)
	basepts = uint32(value & 0xFFFFFFFF)

	return
}

/*
date:2025-09-15
author: weihongqin
description: V7: 
for newly V7
*/
func CombineToInt64V7(isAgora bool, version int8, sessionid uint16, sentenceid uint16, chunkid uint16, issessionend bool, reserved uint32, basepts uint16) int64 {
    var packed int64 = 0

    // 1. 设置 isAgora (2位), 放置在最高位 (第63-62位)
    var agoraVal int64
    if isAgora {
        agoraVal = 0x3 // 二进制 11 (2位)
    } else {
        agoraVal = 0x0 // 二进制 00 (2位)
    }
    packed |= (agoraVal & 0x3) << 62 // 0x3 是2位掩码 (0b11)

    // 2. 设置 version (4位), 放置在第58-61位
    packed |= (int64(version) & 0x0F) << 58 // 0x0F 是4位掩码 (0b1111)

    // 3. 设置 sessionid (6位), 放置在第52-57位
    packed |= (int64(sessionid) & 0x3F) << 52 // 0x3F 是6位掩码 (0b111111)

    // 4. 设置 sentenceid (12位), 放置在第40-51位
    packed |= (int64(sentenceid) & 0x0FFF) << 40 // 0x0FFF 是12位掩码

    // 5. 设置 chunkid (5位), 放置在第35-39位
    packed |= (int64(chunkid) & 0x1F) << 35 // 0x1F 是5位掩码 (0b11111)

    // 6. 设置 issessionend (1位), 放置在第34位
    if issessionend {
        packed |= 1 << 34
    }

    // 7. 设置 reserved (18位), 放置在第16-33位
    packed |= (int64(reserved) & 0x3FFFF) << 16 // 0x3FFFF 是18位掩码

    // 8. 设置 basepts (16位), 放置在第0-15位 (最低位)
    packed |= int64(basepts) & 0xFFFF // 0xFFFF 是16位掩码

    return packed
}
func UnpackFromInt64V7(packed int64) (
    isAgora bool,
    version int8,
    sessionid uint16,
    sentenceid uint16,
    chunkid uint16,
    issessionend bool,
    reserved uint32,
    basepts uint16,
) {
    // 1. 提取 isAgora (2位, 第63-62位)
    // 右移62位后，与2位掩码0x3进行按位与操作，结果不等于0则为true
    isAgora = (packed>>62)&0x3 != 0

    // 2. 提取 version (4位, 第58-61位)
    // 右移58位后，与4位掩码0xF进行按位与操作，然后转换为int8
    version = int8((packed >> 58) & 0x0F)

    // 3. 提取 sessionid (6位, 第52-57位)
    // 右移52位后，与6位掩码0x3F进行按位与操作，然后转换为uint16
    sessionid = uint16((packed >> 52) & 0x3F)

    // 4. 提取 sentenceid (12位, 第40-51位)
    // 右移40位后，与12位掩码0xFFF进行按位与操作，然后转换为uint16
    sentenceid = uint16((packed >> 40) & 0x0FFF)

    // 5. 提取 chunkid (5位, 第35-39位)
    // 右移35位后，与5位掩码0x1F进行按位与操作，然后转换为uint16
    chunkid = uint16((packed >> 35) & 0x1F)

    // 6. 提取 issessionend (1位, 第34位)
    // 右移34位后，与1位掩码0x1进行按位与操作，结果不等于0则为true
    issessionend = (packed>>34)&0x1 != 0

    // 7. 提取 reserved (18位, 第16-33位)
    // 右移16位后，与18位掩码0x3FFFF进行按位与操作，然后转换为uint32
    reserved = uint32((packed >> 16) & 0x3FFFF)

    // 8. 提取 basepts (16位, 第0-15位)
    // 直接与16位掩码0xFFFF进行按位与操作，然后转换为uint16
    basepts = uint16(packed & 0xFFFF)

    return isAgora, version, sessionid, sentenceid, chunkid, issessionend, reserved, basepts
}

//v8


// CombineToInt64V8 将多个字段按照指定位布局组合成int64
// 位布局：1位(固定0) | 1位(isAgora) | 4位(version) | 8位(turnId) | 1位(cmdOrDataType) | 12位(sentenceId) | 5位(reserved1) | 16位(reserved2) | 16位(basePts)
func CombineToInt64V8Data(isAgora bool, version uint8, turnId uint8, isCmdOrData bool, sentenceId uint16, reserved1 uint16, reserved2 uint16, basePts uint16) int64 {
    var result int64 = 0
    
    // 1. 固定0位（最高位），不需要操作，默认为0
    
    // 2. isAgora: 第62位（1位）
    if isAgora {
        result |= 1 << 62 // 将isAgora放在第62位[6,7](@ref)
    }
    
    // 3. version: 第58-61位（4位）
    result |= int64(version&0xF) << 58 // 取低4位，左移58位[7,8](@ref)
    
    // 4. turnId: 第50-57位（8位）
    result |= int64(turnId) << 50
    
    // 5. cmdOrDataType: 第49位（1位）
    if isCmdOrData {
        result |= 1 << 49
    }
    
    // 6. sentenceId: 第37-48位（12位）
    result |= int64(sentenceId&0xFFF) << 37 // 取低12位[5](@ref)
    
    // 7. reserved1: 第32-36位（5位）
    result |= int64(reserved1&0x1F) << 32 // 取低5位
    
    // 8. reserved2: 第16-31位（16位）
    result |= int64(reserved2) << 16
    
    // 9. basePts: 第0-15位（16位）
    result |= int64(basePts)
    
    return result
}
// 位布局：1位(固定0) | 1位(isAgora) | 4位(version) | 8位(turnId) | 1位(cmdOrDataType):1 | 5位(cmdType) | 12位(turnDurationInPacks) | 16位(reserved) | 16位(basePts)
// cmd type：1， current turn is end, follow by the turn duration in packs
// cmd type：2， current turn is interrupted
func CombineToInt64V8Cmd(isAgora bool, version uint8, turnId uint8, isCmdOrData bool, cmdType uint8, turnDurationInPacks uint16, reserved uint16, basePts uint16) int64 {
    var result int64 = 0
    
    // 1. 固定0位（最高位），不需要操作，默认为0
    
    // 2. isAgora: 第62位（1位）
    if isAgora {
        result |= 1 << 62
    }
    
    // 3. version: 第58-61位（4位）
    result |= int64(version&0xF) << 58 // 取低4位
    
    // 4. turnId: 第50-57位（8位）
    result |= int64(turnId) << 50
    
    // 5. cmdOrDataType: 第49位（1位）
    if isCmdOrData {
        result |= 1 << 49
    }
    
    // 6. cmdType: 第44-48位（5位）
    result |= int64(cmdType&0x1F) << 44 // 取低5位
    
    // 7. turnDurationInPacks: 第32-43位（12位）
    result |= int64(turnDurationInPacks&0xFFF) << 32 // 取低12位
    
    // 8. reserved: 第16-31位（16位）
    result |= int64(reserved) << 16
    
    // 9. basePts: 第0-15位（16位）
    result |= int64(basePts)
    
    return result
}

type V8Data struct {
	isAgora bool
	version uint8
	turnId uint8
	isCmdOrData bool
	cmdType uint8
	turnDurationInPacks uint16
	reserved uint16
	sentenceId uint16
	basePts uint16
}

// 反向解析函数，用于验证结果
func ParseInt64V8(value int64) V8Data {
    //fmt.Printf("原始值: %d\n", value)
    //fmt.Printf("二进制: %064b\n", uint64(value))
    
    // 提取各个字段
    isAgora := (value >> 62) & 1
    version := (value >> 58) & 0xF
    turnId := (value >> 50) & 0xFF
    cmdOrDataType := (value >> 49) & 1
    sentenceId := (value >> 37) & 0xFFF
    //reserved1 := (value >> 32) & 0x1F
    reserved2 := (value >> 16) & 0xFFFF
    basePts := value & 0xFFFF

	v8Data := V8Data{
		isAgora: isAgora == 1,
		version: uint8(version),
		turnId: uint8(turnId),
		isCmdOrData: cmdOrDataType == 1,
		cmdType: 1,
		turnDurationInPacks: uint16(sentenceId&0xFFF),
		reserved: uint16(reserved2),
		sentenceId: uint16(sentenceId),
		basePts: uint16(basePts),
	}

	//fmt.Printf("v8Data: %v\n", v8Data)
    //fmt.Printf("isAgora: %d,version: %d,turnId: %d,cmdOrDataType: %d,sentenceId: %d,reserved1: %d,reserved2: %d,basePts: %d\n", isAgora, version, turnId, cmdOrDataType, sentenceId, reserved1, reserved2, basePts)

	return v8Data
}


/*
date:2025-08-14
author: weihongqin
description: lixiang_test
param: conn: rtc connection
param: done: done channel
param: audioSendEvent: audio send event
param: interruptEvent: interrupt event
*/
func chuanyin_test(conn *agoraservice.RtcConnection, done chan bool, audioSendEvent chan struct{}, interruptEvent chan struct{}, file *os.File, samplerate int) {

	// allocate buffer
	leninsecond := 20

	buffer := make([]byte, samplerate*2*leninsecond) // max to 20s data
	readLen, _ := file.Read(buffer)
	bytesinms := samplerate * 2 / 100
	readLen = (readLen / bytesinms) * bytesinms

	// 默认session，sentence，chunkid都是从1开始
	// 模拟场景是： 一个sessio有3个句子； 一个句子分3个chunk来发送
	//
	sessionid := uint16(1)
	sentenceid := uint16(1)
	isSessionEnd := false
	chunkid := uint16(1)
	version := int8(4)

	for {
		select {
		case <-done:
			fmt.Println("lixiang_test done")
			return
		case <-audioSendEvent:
			//date: for loop
			for i := 0; i < 100; i++ {
			isSessionEnd = false
			if sentenceid == 3 && chunkid  == 3 {
				isSessionEnd = true
			}
			masknumber := CombineToInt64(version, sessionid*10, sentenceid, chunkid, isSessionEnd)
			if isSessionEnd {
				readLen = (bytesinms )*7
			}
			ret := conn.PushAudioPcmData(buffer[:readLen], samplerate, 1, masknumber)
			fmt.Printf("lixiang_test audioSendEvent, ret: %d, sessionid: %d, sentenceid: %d, chunkid: %d, isend: %d, masknumber: %d, readLen: %d\n", ret, sessionid, sentenceid, chunkid, isSessionEnd, masknumber, readLen)

			chunkid++
			

			if chunkid > 3 {
				chunkid = 1
				sentenceid++
			}
			if sentenceid > 3 {
				sessionid++
				sentenceid = 1
				chunkid = 1
			}
			if isSessionEnd {
				break
			}
		}
		

		case <-interruptEvent:
			fmt.Println("lixiang_test interruptEvent")
			conn.InterruptAudio()
		default:
			time.Sleep(40 * time.Millisecond)
		}
	}
}

func chuanyin_testV4(conn *agoraservice.RtcConnection, done chan bool, audioSendEvent chan struct{}, interruptEvent chan struct{}, file *os.File, samplerate int) {

	// allocate buffer
	leninsecond := 2

	buffer := make([]byte, samplerate*2*leninsecond) // max to 20s data
	readLen, _ := file.Read(buffer)
	bytesinms := samplerate * 2 / 1000
	readLen = (readLen / bytesinms) * bytesinms

	pa := NewPTSAllocatorV5(samplerate)

	// 默认session，sentence，chunkid都是从1开始
	// 模拟场景是： 一个sessio有3个句子； 一个句子分3个chunk来发送
	//
	sessionid := uint16(1)
	sentenceid := uint16(1)
	isSessionEnd := false
	chunkid := uint16(1)
	

	for {
		select {
		case <-done:
			fmt.Println("lixiang_test done")
			return
		case <-audioSendEvent:
			for i := 0; i < 2; i++ {
			isSessionEnd = false
			curLen := readLen
			if chunkid == 2  {
				isSessionEnd = true

				//对结尾的chunk，最多是10个包！
				curLen = (bytesinms * 10)*7
			}
			masknumber := pa.Allocate(curLen, isSessionEnd)
			ret := conn.PushAudioPcmData(buffer[:curLen], samplerate, 1, masknumber)
			fmt.Printf("lixiang_test audioSendEvent, ret: %d, sessionid: %d, sentenceid: %d, chunkid: %d, isend: %d, masknumber: %d, curLen: %d\n", ret, sessionid, sentenceid, chunkid, isSessionEnd, masknumber, curLen)
			//fmt.Printf("lixiang_test audioSendEvent, ret: %d, sessionid: %d, sentenceid: %d, chunkid: %d, isend: %d, masknumber: %d\n", ret, sessionid, sentenceid, chunkid, isSessionEnd, masknumber)

			chunkid++

			
			if chunkid > 2 {
				sessionid++
				sentenceid = 1
				chunkid = 1
			}
			time.Sleep(100 * time.Millisecond)
		}
		now := time.Now().UnixMilli()
		fmt.Printf("lixiang_test Fin, now: %d\n", now)
		time.Sleep(2000 * time.Millisecond)
		now = time.Now().UnixMilli()
		fmt.Printf("lixiang_test Fin unpblish, now: %d\n", now)
		conn.UnpublishAudio()

		case <-interruptEvent:
			fmt.Println("lixiang_test interruptEvent")
			//conn.InterruptAudio()
		default:
			time.Sleep(40 * time.Millisecond)
		}
	}
}

func chuanyin_testV6(conn *agoraservice.RtcConnection, done chan bool, audioSendEvent chan struct{}, interruptEvent chan struct{}, file *os.File, samplerate int) {

	// allocate buffer
	leninsecond := 2

	buffer := make([]byte, samplerate*2*leninsecond) // max to 20s data
	
	readLen, _ := file.Read(buffer)
	bytesinms := samplerate * 2 / 1000
	readLen = (readLen / bytesinms) * bytesinms

	mutebuffer := make([]byte, bytesinms*10*6)

	//pa := NewPTSAllocatorV5(samplerate)

	// 默认session，sentence，chunkid都是从1开始
	// 模拟场景是： 一个sessio有3个句子； 一个句子分3个chunk来发送
	//
	version := int8(4)
	sessionid := uint16(1)
	sentenceid := uint16(1)
	isSessionEnd := false
	chunkid := uint16(1)
	reserved := uint32(0)
	basepts := uint16(0)
	isAgora := false
	ret := 0
	
	

	for {
		select {
		case <-done:
			fmt.Println("lixiang_test done")
			return
		case <-audioSendEvent:
			for i := 1; i < 3; i++ {
			isSessionEnd = false
			curLen := readLen
			if chunkid == 2  {
				isSessionEnd = true

				//对结尾的chunk，最多是10个包！
				curLen = (bytesinms * 10)*7
			}
			var masknumber int64 = 0

			masknumber = CombineToInt64V7(isAgora, version, sessionid, sentenceid, chunkid, isSessionEnd, reserved, basepts)
		
			if isSessionEnd == false {
				basepts = uint16(time.Now().UnixMilli())
				masknumber = CombineToInt64V7(isAgora, version, sessionid, sentenceid, chunkid, isSessionEnd, reserved, basepts)

				ret = conn.PushAudioPcmData(buffer, samplerate, 1, masknumber)
				fmt.Printf("lixiang_test audioSendEvent file, ret: %d, sessionid: %d, sentenceid: %d, chunkid: %d, isend: %v, masknumber: %d, curLen: %d, basepts: %d\n", ret, sessionid, sentenceid, chunkid, isSessionEnd, masknumber, curLen, basepts)

				/*
				chunkid++
				isSessionEnd = true
				basepts += 2000
				masknumber = CombineToInt64V7(isAgora, version, sessionid, sentenceid, chunkid, isSessionEnd, reserved, basepts)
				ret = conn.PushAudioPcmData(mutebuffer, samplerate, 1, masknumber)
				fmt.Printf("lixiang_test audioSendEvent mut2, ret: %d, sessionid: %d, sentenceid: %d, chunkid: %d, isend: %d, masknumber: %d, curLen: %d, basepts: %d\n", ret, sessionid, sentenceid, chunkid, isSessionEnd, masknumber, curLen, basepts)
				*/

			} else {
				curLen = len(mutebuffer)
				ret = conn.PushAudioPcmData(mutebuffer, samplerate, 1, masknumber)
				fmt.Printf("lixiang_test audioSendEvent mut, ret: %d, sessionid: %d, sentenceid: %d, chunkid: %d, isend: %d, masknumber: %d, curLen: %d\n", ret, sessionid, sentenceid, chunkid, isSessionEnd, masknumber, curLen)

				//ret = conn.PushAudioPcmData(mutebuffer, samplerate, 1, masknumber)
				//fmt.Printf("lixiang_test audioSendEvent mut, ret: %d, sessionid: %d, sentenceid: %d, chunkid: %d, isend: %d, masknumber: %d, curLen: %d\n", ret, sessionid, sentenceid, chunkid, isSessionEnd, masknumber, curLen)


			}
			//ret = conn.PushAudioPcmData(buffer[curLen/2:], samplerate, 1, masknumber)
			//ret = conn.PushAudioPcmData(buffer[curLen/2:], samplerate, 1, masknumber)
			//fmt.Printf("lixiang_test audioSendEvent, ret: %d, sessionid: %d, sentenceid: %d, chunkid: %d, isend: %d, masknumber: %d\n", ret, sessionid, sentenceid, chunkid, isSessionEnd, masknumber)

			chunkid++

			
			if chunkid > 2 {
				sessionid++
				sentenceid = 1
				chunkid = 1
			}
			time.Sleep(100 * time.Millisecond)
		}
		
		case <-interruptEvent:
			fmt.Println("??????lixiang_test interruptEvent")
			conn.InterruptAudio()
		default:
			time.Sleep(40 * time.Millisecond)
		}
	}
}

func chuanyin_callback(sessionId uint32, isFirstOrLastFrame bool) {
	fmt.Printf("---------chuanyin_callback, sessionId: %d, isFirstOrLastFrame: %d\n", sessionId, isFirstOrLastFrame)
}

func chuanyin_testV8(conn *agoraservice.RtcConnection, done chan bool, audioSendEvent chan struct{}, interruptEvent chan struct{}, file *os.File, samplerate int) {

	// allocate buffer
	leninsecond := 20

	buffer := make([]byte, samplerate*2*leninsecond) // max to 20s data
	readLen, _ := file.Read(buffer)
	packsize := int(samplerate * 2 / 1000)*10

	fmt.Printf("lixiang_test, readLen: %d, bytesinms: %d\n", readLen, packsize)
	readLen = int((readLen / packsize) * packsize)
	

	fmt.Printf("lixiang_test, After normalize readLen: %d, bytesinms: %d\n", readLen, packsize)


	// 默认session，sentence，chunkid都是从1开始
	// 模拟场景是： 一个sessio有3个句子； 一个句子分3个chunk来发送
	//
	
	sentenceid := uint16(1)
	
	chunkid := uint16(1)
	version := int8(4)
	turnId := uint8(1)
	reserved1 := uint16(0)
	reserved2 := uint16(0)
	basepts := uint16(0)
	for {
		select {
		case <-done:
			fmt.Println("lixiang_test done")
			return
		case <-audioSendEvent:
			//date: for loop
			for i := 0; i < 2; i++ {
			
			turnId += 10
			sendLen := readLen
			masknumber := CombineToInt64V8Data(true, uint8(version), turnId, false, sentenceid, reserved1, reserved2, basepts)
			ret := conn.PushAudioPcmData(buffer[:sendLen], samplerate, 1, masknumber)
			fmt.Printf("lixiang_test, ret: %d, turnId: %d, sentenceid: %d, basepts: %d, masknumber: %d, readLen: %d\n", ret, turnId, sentenceid, basepts,masknumber, sendLen)

			sendLen = (packsize *7)
			// then follow by the cmd type 1: indicate current turn is end;
			masknumber = CombineToInt64V8Cmd(true, uint8(version), turnId, true, 1, uint16(readLen/packsize), reserved1, basepts)
			ret = conn.PushAudioPcmData(buffer[:sendLen], samplerate, 1, masknumber)
			fmt.Printf("lixiang_test, ret: %d, turnId: %d, sentenceid: %d, basepts: %d, masknumber: %d, readLen: %d\n", ret, turnId, sentenceid, basepts,masknumber, readLen)
			// then follow by the cmd type 2: indicate current turn is interrupted;
			masknumber = CombineToInt64V8Cmd(true, uint8(version), turnId, true, 2, uint16(readLen/packsize), reserved1, basepts)
			ret = conn.PushAudioPcmData(buffer[:sendLen], samplerate, 1, masknumber)


			fmt.Printf("lixiang_test, ret: %d, turnId: %d, sentenceid: %d, basepts: %d, masknumber: %d, readLen: %d\n", ret, turnId, sentenceid, basepts,masknumber, readLen)

			chunkid++
			}
		

		case <-interruptEvent:
			fmt.Println("lixiang_test interruptEvent")
			conn.InterruptAudio()
		default:
			time.Sleep(40 * time.Millisecond)
		}
	}
}

func main() {
	bStop := new(bool)
	*bStop = false
	// start pprof
	go func() {
		runtime.SetBlockProfileRate(1)
		http.ListenAndServe("localhost:6060", nil)
	}()
	// catch ternimal signal
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		*bStop = true
		fmt.Println("Application terminated")
	}()
	readLen := 88656
	bytesinms := 480

	fmt.Printf("lixiang_test, readLen: %d, bytesinms: %d\n", readLen, bytesinms)
	readLen = int((readLen / bytesinms) * bytesinms)
	fmt.Printf("lixiang_test, After normalize readLen: %d, bytesinms: %d\n", readLen, bytesinms)




	println("Start to send and receive PCM data\nusage:\n	./send_recv_pcm <appid> <channel_name>\n	press ctrl+c to exit\n")

	// get parameter from arguments： appid, channel_name

	argus := os.Args
	if len(argus) < 3 {
		fmt.Println("Please input appid, channel name")
		return
	}
	appid := argus[1]
	channelName := argus[2]

	filepath := "../test_data/send_audio_16k_1ch.pcm"
	if len(argus) > 3 {
		filepath = argus[3]
	}
	//default samplerate to 16k
	samplerate := 16000
	if len(argus) > 4 {
		samplerate, _ = strconv.Atoi(argus[4]) // strconv is in the "strconv" package, which is a standard package in Go's library.
	}

	// mode: 1 loopback, 0 playout music
	mode := 0
	if len(argus) > 5 {
		mode, _ = strconv.Atoi(argus[5])
	}

	// get environment variable
	if appid == "" {
		appid = os.Getenv("AGORA_APP_ID")
	}

	cert := os.Getenv("AGORA_APP_CERTIFICATE")

	userId := "0"
	if appid == "" {
		fmt.Println("Please set AGORA_APP_ID environment variable, and AGORA_APP_CERTIFICATE if needed")
		return
	}
	token := ""
	if cert != "" {
		tokenExpirationInSeconds := uint32(3600)
		privilegeExpirationInSeconds := uint32(3600)
		var err error
		token, err = rtctokenbuilder.BuildTokenWithUserAccount(appid, cert, channelName, userId,
			rtctokenbuilder.RolePublisher, tokenExpirationInSeconds, privilegeExpirationInSeconds)
		if err != nil {
			fmt.Println("Failed to build token: ", err)
			return
		}
	}
	svcCfg := agoraservice.NewAgoraServiceConfig()
	svcCfg.AppId = appid
	svcCfg.LogPath = "./agora_rtc_log/agorasdk.log"
	svcCfg.ConfigDir = "./agora_rtc_log"
	svcCfg.DataDir = "./agora_rtc_log"
	svcCfg.IdleMode = true

	// about AudioScenario: default is AudioScenarioAiServer
	// if want to use other scenario, pls contact us and make sure the scenario is much apdated for your business

	agoraservice.Initialize(svcCfg)

	// global set audio dump
	agoraParameterHandler := agoraservice.GetAgoraParameter()
	//agoraParameterHandler.SetParameters("{\"che.audio.acm_ptime\":40}")
	agoraParameterHandler.SetParameters("{\"che.audio.custom_bitrate\":32000}")
	//agoraParameterHandler.SetParameters("{\"che.audio.opus_celt_only\":true}")

	

	conCfg := &agoraservice.RtcConnectionConfig{
		AutoSubscribeAudio: true,
		AutoSubscribeVideo: false,
		ClientRole:         agoraservice.ClientRoleBroadcaster,
		ChannelProfile:     agoraservice.ChannelProfileLiveBroadcasting,
	}
	conSignal := make(chan struct{})
	OnDisconnectedSign := make(chan struct{})

	

	//NOTE: you can set senario here, and every connection has its own senario, which can diff from the service config
	// and can diff from each other
	// but recommend to use the same senario for a connection and related audio track

	publishConfig := agoraservice.NewRtcConPublishConfig()
	

	publishConfig.IsPublishAudio = true
	publishConfig.IsPublishVideo = false
	publishConfig.AudioPublishType = agoraservice.AudioPublishTypePcm
	publishConfig.VideoPublishType = agoraservice.VideoPublishTypeNoPublish
	publishConfig.AudioScenario = agoraservice.AudioScenarioAiServer
	publishConfig.AudioProfile = agoraservice.AudioProfileDefault
	
	publishConfig.SendExternalAudioParameters.Enabled = true
	publishConfig.SendExternalAudioParameters.SendMs = 2000
	publishConfig.SendExternalAudioParameters.SendSpeed = 2

	
	con := agoraservice.NewRtcConnection(conCfg, publishConfig)
	
	
	// 目前客户的方法：todo: tmp, debug, delete
	// 0829版本：也需要在NewRtcConnection之后，调用一次SetParameters，否则，会不生效！！
	// 但在每一NewRtcConnecton	后调用，不会产生线程的泄漏！
	//agoraParameterHandler.SetParameters("{\"che.audio.frame_dump\":{\"location\":\"all\",\"action\":\"start\",\"max_size_bytes\":\"100000000\",\"uuid\":\"123456789\", \"duration\": \"150000\"}}")
	fmt.Printf("agoraParameterHandler: %v\n", agoraParameterHandler)
	//推荐用con.GetAgoraParameter()来调用！这个对任意版本都工作，而且没有线程泄漏
	con.GetAgoraParameter().SetParameters("{\"che.audio.frame_dump\":{\"location\":\"all\",\"action\":\"start\",\"max_size_bytes\":\"100000000\",\"uuid\":\"123456789\", \"duration\": \"150000\"}}")

	
	// todo: chuanyin test
	parser := NewSessionParser(chuanyin_callback)
	parser.Start()

	//audioQueue := agoraservice.NewQueue(10)

	conHandler := &agoraservice.RtcConnectionObserver{
		OnConnected: func(con *agoraservice.RtcConnection, info *agoraservice.RtcConnectionInfo, reason int) {
			// do something
			fmt.Printf("Connected, reason %d\n", reason)
			//NOTE： Must  unpublish,and then update the track, and then publish the track!!!!

			conSignal <- struct{}{}
		},
		OnDisconnected: func(con *agoraservice.RtcConnection, info *agoraservice.RtcConnectionInfo, reason int) {
			// do something
			fmt.Printf("Disconnected, reason %d\n", reason)
			time.Sleep(1000 * time.Millisecond)
			OnDisconnectedSign <- struct{}{}
		},
		OnConnecting: func(con *agoraservice.RtcConnection, conInfo *agoraservice.RtcConnectionInfo, reason int) {
			fmt.Printf("Connecting, reason %d\n", reason)
		},
		OnReconnecting: func(con *agoraservice.RtcConnection, conInfo *agoraservice.RtcConnectionInfo, reason int) {
			fmt.Printf("Reconnecting, reason %d\n", reason)
		},
		OnReconnected: func(con *agoraservice.RtcConnection, conInfo *agoraservice.RtcConnectionInfo, reason int) {
			fmt.Printf("Reconnected, reason %d\n", reason)
		},
		OnConnectionLost: func(con *agoraservice.RtcConnection, conInfo *agoraservice.RtcConnectionInfo) {
			fmt.Printf("Connection lost\n")
		},
		OnConnectionFailure: func(con *agoraservice.RtcConnection, conInfo *agoraservice.RtcConnectionInfo, errCode int) {
			fmt.Printf("Connection failure, error code %d\n", errCode)
		},
		OnUserJoined: func(con *agoraservice.RtcConnection, uid string) {
			fmt.Println("user joined, " + uid)
		},
		OnUserLeft: func(con *agoraservice.RtcConnection, uid string, reason int) {
			fmt.Println("user left, " + uid)
		},
		OnAIQoSCapabilityMissing: func(con *agoraservice.RtcConnection, defaultFallbackSenario int) int {
			fmt.Printf("onAIQoSCapabilityMissing, defaultFallbackSenario: %d\n", defaultFallbackSenario)
			return int(agoraservice.AudioScenarioChorus)
		},
	}
	
	audioSendEvent := make(chan struct{})
	fallackEvent := make(chan struct{})
	interruptEvent := make(chan struct{})

	filename := "./ai_server_recv_pcm.pcm"
	pcm_file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer pcm_file.Close()

	is_square_wave_inpeak := false
	square_wave_count := 0

	bm := NewPcmRawDataManager(16000, 1)
	fmt.Printf("bm: %v\n", bm)
	audioObserver := &agoraservice.AudioFrameObserver{
		OnPlaybackAudioFrameBeforeMixing: func(localUser *agoraservice.LocalUser, channelId string, userId string, frame *agoraservice.AudioFrame, vadResulatState agoraservice.VadState, vadResultFrame *agoraservice.AudioFrame) bool {
			// do something
			//fmt.Printf("Playback audio frame before mixing, from userId %s, far :%d,rms:%d, pitch: %d\n", userId, frame.FarFieldFlag, frame.Rms, frame.Pitch)
			//energy := calculateEnergyFast(frame.Buffer)
			//fmt.Printf("energy: %d, rms: %d, ravg: %f, framecount: %d\n", energy, frame.Rms,float64(energy)/float64(frame.SamplesPerChannel),framecount)
			//for test on 2025-08-18,

			/* todo: test parser
			parser.Parse(frame.PresentTimeMs)
			bm.Push(frame.Buffer)
			data := bm.Pop()
			if data != nil {
				con.PushAudioPcmData(data, frame.SamplesPerSec, frame.Channels, 0)
			}
			*/
			fmt.Printf("Playback audio frame before mixing, from userId %s, far :%d,rms:%d, pitch: %d\n", userId, frame.FarFieldFlag, frame.Rms, frame.Pitch)

			//end

			masknumer := frame.PresentTimeMs
			v8Data := ParseInt64V8(masknumer)
			//fmt.Printf("v8Data: %v\n", v8Data)

		

			if mode == 1 {
				frame.RenderTimeMs = 0
				frame.PresentTimeMs = 0
				//sender.SendAudioPcmData(frame)
				masknumer = CombineToInt64V8Data(v8Data.isAgora, v8Data.version, v8Data.turnId, v8Data.isCmdOrData, v8Data.sentenceId, v8Data.reserved, v8Data.reserved, v8Data.basePts)
				con.PushAudioPcmData(frame.Buffer, frame.SamplesPerSec, frame.Channels, masknumer)
				// loopback
				//audioQueue.Enqueue(frame)

			}
			// 3: 做方波信号的echo
			threshold_value := -25
			// from 2.4.1 for new vad algorithm
			threshold_value =  100
			if mode == 3 {
				if frame.Rms > threshold_value {
					// from trough to peak​​ now
					if is_square_wave_inpeak == false {
						fmt.Printf("????????#####$$$$$$$ square_wave IN peak now: %d, through count: %d\n",1, square_wave_count)
						is_square_wave_inpeak = true
						// just from trough to peak, should trigger a event
						con.InterruptAudio()
						audioSendEvent <- struct{}{}
						square_wave_count = 0
					}
					square_wave_count++
				} else {
					if is_square_wave_inpeak == true {
						fmt.Printf("????????#####$$$$$$$square_wave OUT  now: %d, peak count: %d\n", 1, square_wave_count)
						is_square_wave_inpeak = false
						square_wave_count = 0
					}
					square_wave_count++
				}
			}

			return true
		},
	}
	start := time.Now().UnixMilli()

	agoraservice.GetAgoraParameter()

	event_count := 0

	//added by wei for localuser observer
	localUserObserver := &agoraservice.LocalUserObserver{
		OnStreamMessage: func(localUser *agoraservice.LocalUser, uid string, streamId int, data []byte) {
			// do something
			fmt.Printf("*****Stream message, from userId %s\n", uid)
			//con.SendStreamMessage(streamId, data)
			con.SendStreamMessage(data)
			// trigger a event
			event_count++
			/*
				if event_count == 10 {
					fallackEvent <- struct{}{}
				} else */if event_count%2 == 0 {
				//NonblockNotiyEvent(interruptEvent)
				fmt.Printf("lixiang_test interruptEvent, event_count: %d\n", event_count)

				//interruptEvent <- struct{}{}
			} else { // simulate to interrupt audio
				audioSendEvent <- struct{}{}
				//NonblockNotiyEvent(audioSendEvent)
				fmt.Printf("lixiang_test audioSendEvent, event_count: %d\n", event_count)
			}
		},

		OnAudioVolumeIndication: func(localUser *agoraservice.LocalUser, audioVolumeInfo []*agoraservice.AudioVolumeInfo, speakerNumber int, totalVolume int) {
			// do something
			fmt.Printf("*****Audio volume indication, speaker number %d\n", speakerNumber)
		},
		OnAudioPublishStateChanged: func(localUser *agoraservice.LocalUser, channelId string, oldState int, newState int, elapse_since_last_state int) {
			fmt.Printf("*****Audio publish state changed, old state %d, new state %d\n", oldState, newState)
		},
		OnUserInfoUpdated: func(localUser *agoraservice.LocalUser, uid string, userMediaInfo int, val int) {
			fmt.Printf("*****User info updated, uid %s\n", uid)
		},
		OnUserAudioTrackSubscribed: func(localUser *agoraservice.LocalUser, uid string, remoteAudioTrack *agoraservice.RemoteAudioTrack) {
			fmt.Printf("*****User audio track subscribed, uid %s\n", uid)
		},
		OnUserVideoTrackSubscribed: func(localUser *agoraservice.LocalUser, uid string, info *agoraservice.VideoTrackInfo, remoteVideoTrack *agoraservice.RemoteVideoTrack) {

		},
		OnUserAudioTrackStateChanged: func(localUser *agoraservice.LocalUser, uid string, remoteAudioTrack *agoraservice.RemoteAudioTrack, state int, reason int, elapsed int) {
			fmt.Printf("*****User audio track state changed, uid %s\n", uid)
		},
		OnUserVideoTrackStateChanged: func(localUser *agoraservice.LocalUser, uid string, remoteAudioTrack *agoraservice.RemoteVideoTrack, state int, reason int, elapsed int) {
			fmt.Printf("*****User video track state changed, uid %s\n", uid)
		},
		OnAudioMetaDataReceived: func(localUser *agoraservice.LocalUser, uid string, metaData []byte) {
			fmt.Printf("*****User audio meta data received, uid %s, meta: %s, event_count: %d\n", uid, string(metaData), event_count)

			event_count++
			/*
				if event_count == 10 {
					fallackEvent <- struct{}{}
				} else */if event_count%2 == 0 {
				//NonblockNotiyEvent(interruptEvent)
				fmt.Printf("lixiang_test interruptEvent, event_count: %d\n", event_count)

				//interruptEvent <- struct{}{}
			} else { // simulate to interrupt audio
				audioSendEvent <- struct{}{}
				//NonblockNotiyEvent(audioSendEvent)
				fmt.Printf("lixiang_test audioSendEvent, event_count: %d\n", event_count)
			}

			con.SendAudioMetaData(metaData)
		},
		OnAudioTrackPublishSuccess: func(localUser *agoraservice.LocalUser, audioTrack *agoraservice.LocalAudioTrack) {
			fmt.Printf("*****Audio track publish success, time %d\n", time.Now().UnixMilli())
		},
		OnAudioTrackUnpublished: func(localUser *agoraservice.LocalUser, audioTrack *agoraservice.LocalAudioTrack) {
			fmt.Printf("*****Audio track unpublished, time %d\n", time.Now().UnixMilli())
		},
	}

	con.RegisterObserver(conHandler)

	//end

	// sender := con.NewPcmSender()
	// defer sender.Release()

	localUser := con.GetLocalUser()

	localUser.SetPlaybackAudioFrameBeforeMixingParameters(1, 16000)
	con.RegisterLocalUserObserver(localUserObserver)

	vadConfig := &agoraservice.AudioVadConfigV2{
		PreStartRecognizeCount: 16,
		StartRecognizeCount:    30,
		StopRecognizeCount:     50,
		ActivePercent:          0.7,
		InactivePercent:        0.5,
	}

	con.RegisterAudioFrameObserver(audioObserver, 1, vadConfig)

	con.Connect(token, channelName, userId)
	<-conSignal

	end := time.Now().UnixMilli()
	fmt.Printf("Connect cost %d ms\n", end-start)

	start_publish := time.Now().UnixMilli()
	con.PublishAudio()
	end_publish := time.Now().UnixMilli()
	fmt.Printf("Publish audio cost %d ms\n", end_publish-start_publish)
	//time.Sleep(1000 * time.Millisecond)
	//con.UnpublishAudio()
	start_publish = time.Now().UnixMilli()
	fmt.Printf("Unpublish audio cost %d ms\n", start_publish-end_publish)

	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("NewError opening file: %v\n", err)
		return
	}
	defer file.Close()

	done := make(chan bool)
	// new method for push
	/*

			#下面的操作：只是模拟生产的数据。
			# - 在sample中，为了确保生产产生的数据能够一直播放，需要生产足够多的数据，所以用这样的方式来模拟
			# - 在实际使用中，数据是实时产生的，所以不需要这样的操作。只需要在TTS生产数据的时候，调用AudioConsumer.push_pcm_data()
			 # 我们启动2个task
		    # 一个task，用来模拟从TTS接收到语音，然后将语音push到audio_consumer
		    # 另一个task，用来模拟播放语音：从audio_consumer中取出语音播放
		    # 在实际应用中，可以是TTS返回的时候，直接将语音push到audio_consumer
		    # 然后在另外一个“timer”的触发函数中，调用audio_consumer.consume()。
		    # 推荐：
		    # .Timer的模式；也可以和业务已有的timer结合在一起使用，都可以。只需要在timer 触发的函数中，调用audio_consumer.consume()即可
		    # “Timer”的触发间隔，可以和业务已有的timer间隔一致，也可以根据业务需求调整，推荐在40～80ms之间

	*/
	if mode == 0 {
		go ReadFileToConsumer(file, con, 50, done, samplerate)
	} else if mode == 2 || mode == 3 {
		//go SendTTSDataToClient(samplerate, audioConsumer, file, done, audioSendEvent, fallackEvent, localUser, track)
		go func() {
			//consturt a audio frame
			leninsecond := 30
			if mode == 3 {
				leninsecond = 5
			}
			buffer := make([]byte, samplerate*2*leninsecond) // max to 20s data
			bytesPerFrame := (samplerate / 100) * 2 * 1      // 10ms , mono
			bytesInMs := samplerate * 2 / 1000
			frame := &agoraservice.AudioFrame{
				Buffer:            nil,
				RenderTimeMs:      0,
				SamplesPerChannel: samplerate / 100,
				BytesPerSample:    2,
				Channels:          1,
				SamplesPerSec:     samplerate,
				Type:              agoraservice.AudioFrameTypePCM16,
			}
			for {
				select {
				case <-done:
					fmt.Println("SendAudioToClient done")
					return
				case <-fallackEvent:
					fmt.Println("?????? fallackEvent")

					con.UpdateAudioSenario(agoraservice.AudioScenarioDefault)

				case <-interruptEvent:
					fmt.Println("?????? interruptEvent")
					con.InterruptAudio()

				case <-audioSendEvent:
					// read 1s data from file
					//con.PublishAudio()

					readLen, err := file.Read(buffer)
					if err != nil {
						fmt.Printf("read up to EOF,cur read: %d", readLen)
						file.Seek(0, 0)
						continue
					}
					//round readLen to bytesInMs
					readLen = (readLen / bytesInMs) * bytesInMs
					frame.Buffer = buffer[:readLen]
					packnum := readLen / bytesPerFrame
					frame.SamplesPerChannel = (samplerate / 100) * packnum
					ret := con.PushAudioPcmData(buffer[:readLen], samplerate, 1, 0)

					//audioConsumer.PushPCMData(frame.Buffer)
					// and seek to the begin of the file
					file.Seek(0, 0)
					fmt.Printf("SendTTSDataToClient done, ret: %d, samplerate: %d\n", ret, samplerate)
				default:
					time.Sleep(40 * time.Millisecond)
				}
			}
		}()
	} else if mode == 4 {
		//go SendTTSDataToClient(samplerate, audioConsumer, file, done, audioSendEvent, fallackEvent, localUser, track)
		// V2
		go chuanyin_test(con, done, audioSendEvent, interruptEvent, file, samplerate)
		// V8
		//go chuanyin_testV8(con, done, audioSendEvent, interruptEvent, file, samplerate)
	}

	//release operation:cancel defer release,try manual release
	for !(*bStop) {
		time.Sleep(100 * time.Millisecond)
		//simulate to send audio meta data

	}
	close(done)

	start_disconnect := time.Now().UnixMilli()
	con.Disconnect()
	//<-OnDisconnectedSign

	//time.Sleep(5000 * time.Millisecond)

	con.Release()

	time.Sleep(1000 * time.Millisecond)

	agoraservice.Release()

	audioObserver = nil
	localUserObserver = nil
	localUser = nil
	conHandler = nil
	con = nil

	parser.End()

	fmt.Printf("Disconnected, cost %d ms\n", time.Now().UnixMilli()-start_disconnect)
}
