package main

import (
	"fmt"
	"os"
	"time"

	"agora.io/agoraservice"
)

func main() {
	svcCfg := agoraservice.AgoraServiceConfig{
		AppId:         "aab8b8f5a8cd4469a63042fcfafe7063",
		AudioScenario: agoraservice.AUDIO_SCENARIO_CHORUS,
		LogPath:       "./agora_rtc_log/agorasdk.log",
		LogSize:       512 * 1024,
	}
	agoraservice.Init(&svcCfg)
	conCfg := agoraservice.RtcConnectionConfig{
		SubAudio:       true,
		SubVideo:       false,
		ClientRole:     1,
		ChannelProfile: 1,

		SubAudioConfig: &agoraservice.SubscribeAudioConfig{
			SampleRate: 16000,
			Channels:   1,
		},
	}
	conSignal := make(chan struct{})
	conHandler := agoraservice.RtcConnectionEventHandler{
		OnConnected: func(con *agoraservice.RtcConnection, info *agoraservice.RtcConnectionInfo, reason int) {
			// do something
			fmt.Println("Connected")
			conSignal <- struct{}{}
		},
		OnDisconnected: func(con *agoraservice.RtcConnection, info *agoraservice.RtcConnectionInfo, reason int) {
			// do something
			fmt.Println("Disconnected")
		},
		OnUserJoined: func(con *agoraservice.RtcConnection, uid string) {
			fmt.Println("user joined, " + uid)
		},
		OnUserLeft: func(con *agoraservice.RtcConnection, uid string, reason int) {
			fmt.Println("user left, " + uid)
		},
	}
	conCfg.ConnectionHandler = &conHandler
	conCfg.AudioFrameObserver = &agoraservice.RtcConnectionAudioFrameObserver{
		OnPlaybackAudioFrameBeforeMixing: func(con *agoraservice.RtcConnection, channelId string, userId string, frame *agoraservice.PcmAudioFrame) {
			// do something
			fmt.Printf("Playback audio frame before mixing, from userId %s\n", userId)
		},
	}
	con := agoraservice.NewConnection(&conCfg)
	defer con.Release()
	sender := con.NewPcmSender()
	defer sender.Release()
	con.Connect("", "lhztest1", "0")
	<-conSignal
	sender.Start()

	frame := agoraservice.PcmAudioFrame{
		Data:              make([]byte, 320),
		Timestamp:         0,
		SamplesPerChannel: 160,
		BytesPerSample:    2,
		NumberOfChannels:  1,
		SampleRate:        16000,
	}

	file, err := os.Open("../test_data/demo.pcm")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	sender.AdjustVolume(100)
	// sender.SetSendBufferSize(30)

	sendCount := 0
	// send 180ms audio data
	for i := 0; i < 18; i++ {
		dataLen, err := file.Read(frame.Data)
		if err != nil || dataLen < 320 {
			fmt.Println("Finished reading file:", err)
			break
		}
		sendCount++
		ret := sender.SendPcmData(&frame)
		fmt.Printf("SendPcmData %d ret: %d\n", sendCount, ret)
	}

	bStop := false
	firstSendTime := time.Now()
	for !bStop {
		shouldSendCount := int(time.Since(firstSendTime).Milliseconds()/10) - (sendCount - 18)
		for i := 0; i < shouldSendCount; i++ {
			dataLen, err := file.Read(frame.Data)
			if err != nil || dataLen < 320 {
				fmt.Println("Finished reading file:", err)
				bStop = true
				break
			}

			sendCount++
			ret := sender.SendPcmData(&frame)
			fmt.Printf("SendPcmData %d ret: %d\n", sendCount, ret)
		}
		fmt.Printf("Sent %d frames this time\n", shouldSendCount)
		time.Sleep(50 * time.Millisecond)
	}
	sender.Stop()
	con.Disconnect()

	agoraservice.Destroy()
}
