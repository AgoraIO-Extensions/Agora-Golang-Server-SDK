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

	agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/rtc"

)


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
	

	println("Start to audio processor\nusage:\n	./example_externalaudio_processor <appid> <channel_name> <outsamplerate> <outputchannels> \n	press ctrl+c to exit\n")

	// get parameter from arguments： appid, channel_name

	argus := os.Args
	if len(argus) < 3 {
		fmt.Println("Please input appid, channel name")
		return
	}
	appid := argus[1]
	channelName := argus[2]
	filepath := ""
	inputsamplerate := 16000
	inputchannels := 1

	if len(argus) > 3 {
		filepath = argus[3]
	}
	if len(argus) > 4 {
		inputsamplerate, _ = strconv.Atoi(argus[4])
	}
	if len(argus) > 5 {
		inputchannels, _ = strconv.Atoi(argus[5])
	}
	
	fmt.Println("appid:", appid, "channelName:", channelName, "filepath:", filepath, "inputsamplerate:", inputsamplerate, "inputchannels:", inputchannels)

	//NOTE: only verify statick token appid,for non-static token,should build token by yourself!!
	// how to build dynamic token, please refer to the official documentation
	

	//step1: initialize the agora service config
	svcCfg := agoraservice.NewAgoraServiceConfig()
	svcCfg.AppId = appid
	svcCfg.LogPath = "./agora_rtc_log/agorasdk.log"
	svcCfg.ConfigDir = "./agora_rtc_log"
	svcCfg.DataDir = "./agora_rtc_log"
	svcCfg.IdleMode = true
	svcCfg.EnableAPM = true
	svcCfg.APMConfig = agoraservice.NewAPMConfig()
	svcCfg.APMConfig.EnableDump = true

	//step2: initialize the agora service
	agoraservice.Initialize(svcCfg)

	//step3: create a external audio processor
	externalAudioProcessor := agoraservice.NewExternalAudioProcessor()

	//step4: initialize the external audio processor
	ret := externalAudioProcessor.Initialize(inputsamplerate, inputchannels)
	if ret != 0 {
		fmt.Printf("Failed to initialize external audio processor, error code: %d\n", ret)
		return
	}

	//step5: push audio pcm data to the external audio processor
	
	
	// open file and push to external audio processor
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("Failed to open file: %v, filepath: %s\n", err, filepath)
		return
	}
	defer file.Close()
	
	// 
	lenInMs := inputsamplerate*inputchannels*2/1000
	Num := 10
	buffer := make([]byte, lenInMs*Num)
	totalLen := 0
	totalInPacks := 0

	startTime := time.Now().UnixMilli()

	/*
	1、可以不需要间隔时间推送
	2、回调是固定的10ms的一个包
	2、每次10ms的包是可以的，只要是10x的包也是可以的;
	3、只要是整数倍的就可以
	*/
	

	//release operation:cancel defer release,try manual release
	for !(*bStop) {
		// read file and push to external audio processor
		readLen, err := file.Read(buffer)
		if err != nil {
			fmt.Printf("Failed to read file: %v, filepath: %s\n", err, filepath)
			return
		}
		if readLen < lenInMs*Num {
			break;
		}
		totalLen += readLen
		totalInPacks ++
		ret := externalAudioProcessor.PushAudioPcmData(buffer[:readLen], inputsamplerate, inputchannels, 0)
		if ret != 0 {
			fmt.Printf("Failed to push audio pcm data: %d\n", ret)
			return
		}
		//time.Sleep(10 * time.Millisecond)
	}
	fmt.Printf("Total length: %d, packs: %d, total time: %dms\n", totalLen, totalInPacks, time.Now().UnixMilli() - startTime)

	externalAudioProcessor.Release()

	agoraservice.Release()

	fmt.Println("Application terminated")
}
