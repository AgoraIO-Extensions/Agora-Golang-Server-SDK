package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/rtc"
	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
)

func main() {
	bStop := new(bool)
	*bStop = false
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		*bStop = true
		fmt.Println("Application terminated")
	}()

	argus := os.Args
	if len(argus) < 3 {
		fmt.Println("Please input appid, channel name")
		return
	}
	appid := argus[1]
	channelName := argus[2]
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
	svcCfg.EnableVideo = true
	svcCfg.AppId = appid
	svcCfg.LogPath = "./agora_rtc_log/agorasdk.log"
	svcCfg.ConfigDir = "./agora_rtc_log"
	svcCfg.DataDir = "./agora_rtc_log"
	agoraservice.Initialize(svcCfg)

	conCfg := agoraservice.RtcConnectionConfig{
		AutoSubscribeAudio: true,
		AutoSubscribeVideo: true,
		ClientRole:         agoraservice.ClientRoleBroadcaster,
		ChannelProfile:     agoraservice.ChannelProfileLiveBroadcasting,
	}
	publishConfig := agoraservice.NewRtcConPublishConfig()
	publishConfig.AudioPublishType = agoraservice.AudioPublishTypePcm
	publishConfig.AudioScenario = agoraservice.AudioScenarioDefault
	publishConfig.IsPublishAudio = true
	publishConfig.IsPublishVideo = true
	publishConfig.VideoPublishType = agoraservice.VideoPublishTypeYuv
	publishConfig.AudioProfile = agoraservice.AudioProfileDefault

	conSignal := make(chan struct{})
	conHandler := &agoraservice.RtcConnectionObserver{
		OnConnected: func(con *agoraservice.RtcConnection, info *agoraservice.RtcConnectionInfo, reason int) {
			fmt.Printf("Connected, reason %d\n", reason)
			conSignal <- struct{}{}
		},
		OnDisconnected: func(con *agoraservice.RtcConnection, info *agoraservice.RtcConnectionInfo, reason int) {
			fmt.Printf("Disconnected, reason %d\n", reason)
		},
		OnUserJoined: func(con *agoraservice.RtcConnection, uid string) {
			fmt.Println("user joined, " + uid)
		},
		OnUserLeft: func(con *agoraservice.RtcConnection, uid string, reason int) {
			fmt.Println("user left, " + uid)
		},
	}

	frameCount := 0
	frameLastRecvTime := time.Now().UnixMilli()
	copiedFrameCount := 0
	videoObserver := &agoraservice.VideoFrameObserver{
		OnReusedFrame: func(channelId string, userId string, frame *agoraservice.VideoFrame) bool {
			if frame == nil {
				return true
			}
			yCopy := append([]byte(nil), frame.YBuffer...)
			uCopy := append([]byte(nil), frame.UBuffer...)
			vCopy := append([]byte(nil), frame.VBuffer...)
			if len(yCopy) == 0 || len(uCopy) == 0 || len(vCopy) == 0 {
				fmt.Printf("invalid reused frame copy, channel=%s user=%s y=%d u=%d v=%d\n",
					channelId, userId, len(yCopy), len(uCopy), len(vCopy))
				return false
			}
			copiedFrameCount++
			frameCount++
			now := time.Now().UnixMilli()
			if now-frameLastRecvTime > 1000 {
				fps := int64(frameCount*1000) / (now - frameLastRecvTime)
				fmt.Printf("reused raw video ok, copied=%d fps=%d y=%d u=%d v=%d\n",
					copiedFrameCount, fps, len(yCopy), len(uCopy), len(vCopy))
				frameCount = 0
				copiedFrameCount = 0
				frameLastRecvTime = now
			}
			return true
		},
	}

	con := agoraservice.NewRtcConnection(&conCfg, publishConfig)
	con.RegisterObserver(conHandler)
	con.RegisterVideoFrameObserver(videoObserver)
	con.Connect(token, channelName, userId)
	<-conSignal

	encoderCfg := agoraservice.NewVideoEncoderConfiguration()
	encoderCfg.Width = 960
	encoderCfg.Height = 720
	encoderCfg.Framerate = 20
	encoderCfg.Bitrate = 1700
	encoderCfg.MinBitrate = -1
	encoderCfg.CodecType = agoraservice.VideoCodecTypeAv1
	con.SetVideoEncoderConfiguration(encoderCfg)
	con.PublishVideo()

	w := 960
	h := 720
	file, err := os.Open("../test_data/960-720.yuv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	con.GetLocalUser().SubscribeAllVideo(&agoraservice.VideoSubscriptionOptions{
		StreamType:       agoraservice.VideoStreamLow,
		EncodedFrameOnly: false,
	})

	dataSize := w * h * 3 / 2
	data := make([]byte, dataSize)
	for !*bStop {
		dataLen, err := file.Read(data)
		if err != nil || dataLen < dataSize {
			file.Seek(0, 0)
			continue
		}
		frame := &agoraservice.ExternalVideoFrame{
			Type:      agoraservice.VideoBufferRawData,
			Format:    agoraservice.VideoPixelI420,
			Buffer:    data,
			Stride:    w,
			Height:    h,
			Timestamp: 0,
		}
		con.PushVideoFrame(frame)
		time.Sleep(33 * time.Millisecond)
	}
}
