package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/agoraservice"

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
	}()

	// get parameter from arguments： appid, channel_name
	argus := os.Args
	if len(argus) < 3 {
		fmt.Println("Please input appid, channel name")
		return
	}
	appid := argus[1]
	channelName := argus[2]

	// get environment variable
	//appid := os.Getenv("AGORA_APP_ID")
	cert := os.Getenv("AGORA_APP_CERTIFICATE")
	//channelName := "gosdktest"
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

	agoraservice.Initialize(svcCfg)
	mediaNodeFactory := agoraservice.NewMediaNodeFactory()

	conCfg := agoraservice.RtcConnectionConfig{
		AutoSubscribeAudio: true,
		AutoSubscribeVideo: true,
		ClientRole:         agoraservice.ClientRoleBroadcaster,
		ChannelProfile:     agoraservice.ChannelProfileLiveBroadcasting,
	}
	conSignal := make(chan struct{})
	conHandler := &agoraservice.RtcConnectionObserver{
		OnConnected: func(con *agoraservice.RtcConnection, info *agoraservice.RtcConnectionInfo, reason int) {
			// do something
			fmt.Printf("Connected, reason %d\n", reason)
			conSignal <- struct{}{}
		},
		OnDisconnected: func(con *agoraservice.RtcConnection, info *agoraservice.RtcConnectionInfo, reason int) {
			// do something
			fmt.Printf("Disconnected, reason %d\n", reason)
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
	}
	// for calling onFrame
	frameCount := 0
	frameLastRecvTime := time.Now().UnixMilli()
	videoObserver := &agoraservice.VideoFrameObserver{
		OnFrame: func(channelId string, userId string, frame *agoraservice.VideoFrame) bool {
			if frame == nil {
				return true
			    
			}
			frameCount++
			Now := time.Now().UnixMilli()
			if Now -frameLastRecvTime > 1000 {
				fps := int64(frameCount*1000) / (Now - frameLastRecvTime)
				fmt.Printf("fps, %d fps, %d\n", frameCount, fps)
				frameCount = 0
				frameLastRecvTime = time.Now().UnixMilli()
			}
			// do something
			fmt.Printf("recv video frame, from channel %s, user %s\n", channelId, userId)
			return true
		},
	}
	scenario := svcCfg.AudioScenario
	con := agoraservice.NewRtcConnection(&conCfg, scenario)

	localUser := con.GetLocalUser()
	con.RegisterObserver(conHandler)
	localUser.RegisterVideoFrameObserver(videoObserver)

	sender := mediaNodeFactory.NewVideoFrameSender()
	
	track := agoraservice.NewCustomVideoTrackFrame(sender)


	con.Connect(token, channelName, userId)
	<-conSignal

	track.SetVideoEncoderConfiguration(&agoraservice.VideoEncoderConfiguration{
		CodecType:         agoraservice.VideoCodecTypeH264,
		Width:             320,
		Height:            240,
		Framerate:         30,
		Bitrate:           1000,
		MinBitrate:        400,
		OrientationMode:   agoraservice.OrientationModeAdaptive,
		DegradePreference: 2,
	})
	track.SetEnabled(true)
	localUser.PublishVideo(track)

	// for yuv test
	
	w := 352
	h := 288
	dataSize := w * h * 3 / 2
	data := make([]byte, dataSize)
	// read yuv from file 103_RaceHorses_416x240p30_300.yuv
	file, err := os.Open("../test_data/send_video_cif.yuv")
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
		sender.SendVideoFrame(&agoraservice.ExternalVideoFrame{
			Type:      agoraservice.VideoBufferRawData,
			Format:    agoraservice.VideoPixelI420,
			Buffer:    data,
			Stride:    w,
			Height:    h,
			Timestamp: 0,
		})
		time.Sleep(33 * time.Millisecond)
	}
	/*
	// rgag colos space type test
	w := 360
	h := 720
	// for rgba
	dataSize := w * h * 4
	data := make([]byte, dataSize)
	// read yuv from file 103_RaceHorses_416x240p30_300.yuv
	file, err := os.Open("../test_data/rgba_360_720.data")
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
		sender.SendVideoFrame(&agoraservice.ExternalVideoFrame{
			Type:      agoraservice.VideoBufferRawData,
			Format:    agoraservice.VideoPixelRGBA,
			Buffer:    data,
			Stride:    w,
			Height:    h,
			Timestamp: 0,
			
			// for rgba with pure background color test
			ColorSpace: agoraservice.ColorSpaceType{
				MatrixId:    1,
				PrimariesId: 1,
				RangeId:     2, //or 2,
				TransferId:  1,
			},
		})
		time.Sleep(33 * time.Millisecond)
	}
*/
	//release now
	

	localUser.UnpublishVideo(track)
	track.SetEnabled(false)
	localUser.UnregisterAudioFrameObserver()
	localUser.UnregisterVideoFrameObserver()
	localUser.UnregisterLocalUserObserver()

	start_disconnect := time.Now().UnixMilli()
	con.Disconnect()
	//<-OnDisconnectedSign
	con.UnregisterObserver()

	con.Release()

	track.Release()
	sender.Release()
	mediaNodeFactory.Release()
	agoraservice.Release()

	track = nil
	videoObserver = nil
	
	localUser = nil
	conHandler = nil
	con = nil

	fmt.Printf("Disconnected, cost %d ms\n", time.Now().UnixMilli()-start_disconnect)
}
