package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"agora.io/agoraservice"

	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
)

func main() {
	bStop := new(bool)
	*bStop = false
	// catch ternimal signal
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		*bStop = true
		fmt.Println("Application terminated")
		os.Exit(0)
	}()

	// get environment variable
	appid := os.Getenv("AGORA_APP_ID")
	cert := os.Getenv("AGORA_APP_CERTIFICATE")
	channelName := "gosdktest"
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
	svcCfg := agoraservice.AgoraServiceConfig{
		EnableAudioProcessor: true,
		EnableAudioDevice:    false,
		EnableVideo:          true,

		AppId:          appid,
		ChannelProfile: agoraservice.ChannelProfileLiveBroadcasting,
		AudioScenario:  agoraservice.AudioScenarioChorus,
		UseStringUid:   false,
		LogPath:        "./agora_rtc_log/agorasdk.log",
		LogSize:        512 * 1024,
	}
	agoraservice.Initialize(&svcCfg)
	defer agoraservice.Release()

	conCfg := agoraservice.RtcConnectionConfig{
		AutoSubscribeAudio: true,
		AutoSubscribeVideo: true,
		ClientRole:         agoraservice.ClientRoleBroadcaster,
		ChannelProfile:     agoraservice.ChannelProfileLiveBroadcasting,
	}
	conSignal := make(chan struct{})
	conHandler := &agoraservice.RtcConnectionEventHandler{
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
	videoObserver := &agoraservice.VideoFrameObserver{
		OnFrame: func(con *agoraservice.RtcConnection, channelId string, userId string, frame *agoraservice.VideoFrame) {
			// do something
			fmt.Printf("recv video frame, from channel %s, user %s\n", channelId, userId)
		},
	}
	con := agoraservice.NewConnection(&conCfg)
	defer con.Release()

	con.RegisterObserver(conHandler)
	con.RegisterVideoFrameObserver(videoObserver)

	sender := agoraservice.NewVideoFrameSender()
	defer sender.Release()
	track := agoraservice.NewCustomVideoTrack(sender)
	defer track.Release()

	con.Connect(token, channelName, userId)
	<-conSignal

	track.SetVideoEncoderConfiguration(&agoraservice.VideoEncoderConfiguration{
		CodecType:         2,
		Width:             320,
		Height:            240,
		Framerate:         30,
		Bitrate:           500,
		MinBitrate:        100,
		OrientationMode:   0,
		DegradePreference: 0,
	})
	track.SetEnabled(true)
	con.PublishVideo(track)

	w := 416
	h := 240
	dataSize := w * h * 3 / 2
	data := make([]byte, dataSize)
	// read yuv from file 103_RaceHorses_416x240p30_300.yuv
	file, err := os.Open("../../../test_data/103_RaceHorses_416x240p30_300.yuv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	for !*bStop {
		dataLen, err := file.Read(data)
		if err != nil || dataLen < dataSize {
			file.Seek(0, 0)
			continue
		}
		// senderCon.SendStreamMessage(streamId, data)
		sender.SendVideoFrame(&agoraservice.VideoFrame{
			Buffer:    data,
			Width:     w,
			Height:    h,
			YStride:   w,
			UStride:   w / 2,
			VStride:   w / 2,
			Timestamp: 0,
		})
		time.Sleep(33 * time.Millisecond)
	}
	con.UnpublishVideo(track)
	track.SetEnabled(false)
	con.Disconnect()
}
