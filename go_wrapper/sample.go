package main

import (
	"fmt"
	"os"
	"time"

	"agora.io/agoraservice"
)

func main() {
	svcCfg := agoraservice.AgoraServiceConfig{
		AppId: "aab8b8f5a8cd4469a63042fcfafe7063",
	}
	agoraservice.Init(&svcCfg)
	conCfg := agoraservice.RtcConnectionConfig{
		SubAudio:       true,
		SubVideo:       false,
		ClientRole:     1,
		ChannelProfile: 1,
	}
	conSignal := make(chan struct{})
	conHandler := agoraservice.RtcConnectionEventHandler{
		OnConnected: func(con *agoraservice.RtcConnection, info *agoraservice.RtcConnectionInfo, reason int) {
			// do something
			fmt.Println("Connected")
			conSignal <- struct{}{}
		},
	}
	con := agoraservice.NewConnection(&conCfg, &conHandler)
	sender := con.NewPcmSender()
	con.Connect("", "lhztest", "0")
	<-conSignal
	sender.Start()

	frame := agoraservice.PcmAudioFrame{
		Data:              make([]byte, 320),
		CaptureTimestamp:  0,
		SamplesPerChannel: 160,
		BytesPerSample:    2,
		NumberOfChannels:  1,
		SampleRate:        16000,
	}

	file, err := os.Open("../agora_sdk/demo.pcm")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	for i := 0; i < 500; i++ {
		_, err := file.Read(frame.Data)
		if err != nil {
			fmt.Println("Error reading file:", err)
			break
		}

		sender.SendPcmData(&frame)
		time.Sleep(10 * time.Millisecond)
	}
	sender.Stop()
	sender.Release()
	con.Disconnect()
	con.Release()
	agoraservice.Destroy()
}
