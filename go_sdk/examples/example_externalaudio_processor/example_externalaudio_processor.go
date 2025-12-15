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
	"time"

	agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/rtc"

)



//mode: 0: only vad, 1: apm+vad, 2: apm+vad+dump
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
	

	println("Start to audio processor\nusage:\n	./example_externalaudio_processor <appid> <channel_name> <outsamplerate> <outputchannels> <mode> \n	press ctrl+c to exit\n")

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
	mode := 0

	if len(argus) > 3 {
		filepath = argus[3]
	}
	if len(argus) > 4 {
		inputsamplerate, _ = strconv.Atoi(argus[4])
	}
	if len(argus) > 5 {
		inputchannels, _ = strconv.Atoi(argus[5])
	}
	if len(argus) > 6 {
		mode, _ = strconv.Atoi(argus[6])
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
	svcCfg.APMModel = 1
	svcCfg.APMConfig = agoraservice.NewAPMConfig()
	svcCfg.APMConfig.EnableDump = false

	if mode == 2 {
		svcCfg.APMConfig.EnableDump = true
	}

	//step2: initialize the agora service
	agoraservice.Initialize(svcCfg)

	//step3: create a external audio processor
	externalAudioProcessor := agoraservice.NewExternalAudioProcessor()

	//step4: initialize the external audio processor
	apmConfig := agoraservice.NewAPMConfig()
	apmConfig.EnableDump = false
	apmConfig.AiAecConfig.Enabled = false
	apmConfig.BghvsCConfig.Enabled = true
	apmConfig.AgcConfig.Enabled = false
	apmConfig.AiNsConfig.AiNSEnabled = true
	apmConfig.AiNsConfig.NsEnabled = true
	apmConfig.BghvsCConfig.Enabled = true
	
	vadConfig := &agoraservice.AudioVadConfigV2{
		EnableAdaptiveRmsThreshold: true,
		AdaptiveRmsThresholdFactor: 0.67,
		StartVoiceProb: 70,
		StartRms: -50,
		StopVoiceProb: 70,
		StopRms: -50,
	}

	// add observer
	observer := &agoraservice.ExternalAudioProcessorObserver{
		OnProcessedAudioFrame: func(processor *agoraservice.ExternalAudioProcessor, frame *agoraservice.AudioFrame, vadResultStat agoraservice.VadState, vadResultFrame *agoraservice.AudioFrame) {
			var retBufferLen int = 0
			if vadResultFrame != nil {
				retBufferLen = len(vadResultFrame.Buffer)
			}
			fmt.Printf("OnProcessedAudioFrame, voice prob: %d, rms: %d, pitch: %d, vadResultStat: %d, vadResultFrame len: %d\n", frame.VoiceProb, frame.Rms, frame.Pitch, vadResultStat, retBufferLen)
		},
	}


	// mode: 0: only vad, so set apmconfig to nil
	if mode == 0 {
		apmConfig = nil
	}
	

	ret := externalAudioProcessor.Initialize(apmConfig, inputsamplerate, inputchannels, vadConfig,observer)
	if ret != 0 {
		fmt.Printf("Failed to initialize external audio processor, error code: %d\n", ret)
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
	fmt.Printf("Total length: %d, packs: %d, total time: %dms, inStreamInMs: %d, processedStreamInMs: %d\n", totalLen, totalInPacks, time.Now().UnixMilli() - startTime,
externalAudioProcessor.InStreamInMs, externalAudioProcessor.ProcessedStreamInMs)

	externalAudioProcessor.Release()

	agoraservice.Release()

	fmt.Println("Application terminated")
}
