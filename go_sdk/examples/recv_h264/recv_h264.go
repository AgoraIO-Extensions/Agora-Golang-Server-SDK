package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/rtc"

	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
)

type EncodedVideoData struct {
	uid       string
	imageData []byte
	frameInfo *agoraservice.EncodedVideoFrameInfo
}

func main() {
	bStop := new(bool)
	*bStop = false
	ctx, cancel := context.WithCancel(context.Background())
	// catch ternimal signal
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		cancel()
	}()
	argus := os.Args
	if len(argus) < 3 {
		fmt.Println("Please input appid, channel name")
		return
	}
	appid := argus[1]
	channelName := argus[2]

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
	svcCfg.EnableVideo = true
	svcCfg.AppId = appid

	agoraservice.Initialize(svcCfg)
	
	
	
	var con *agoraservice.RtcConnection = nil

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
		OnAIQoSCapabilityMissing: func(con *agoraservice.RtcConnection, defaultFallbackSenario int) int {
			fmt.Printf("onAIQoSCapabilityMissing, defaultFallbackSenario: %d\n", defaultFallbackSenario)
			return int(agoraservice.AudioScenarioDefault)
		},
	}

	videoChan := make(chan *EncodedVideoData, 100)
	localUserObserver := &agoraservice.LocalUserObserver{
		OnUserVideoTrackSubscribed: func(localUser *agoraservice.LocalUser, uid string, info *agoraservice.VideoTrackInfo, remoteVideoTrack *agoraservice.RemoteVideoTrack) {
			fmt.Printf("user %s video subscribed\n", uid)
		},
	}
	lastKeyFrameTime := time.Now().UnixMilli()
	encodedVideoObserver := &agoraservice.VideoEncodedFrameObserver{
		OnEncodedVideoFrame: func(uid string, imageBuffer []byte, frameInfo *agoraservice.EncodedVideoFrameInfo) bool {
			// fmt.Printf("user %s encoded video received\n", uid)
			//fmt.Printf("user %s encoded video received, frame type %d\n", uid, frameInfo.FrameType)
			if frameInfo.FrameType == agoraservice.VideoFrameTypeKeyFrame {
				fmt.Printf("key frame received, time %d\n", time.Now().UnixMilli()-lastKeyFrameTime)
				lastKeyFrameTime = time.Now().UnixMilli()
			}
			videoChan <- &EncodedVideoData{
				uid:       uid,
				imageData: imageBuffer,
				frameInfo: frameInfo,
			}
			return true
		},
	}
	scenario := agoraservice.AudioScenarioAiServer
	conCfg := &agoraservice.RtcConnectionConfig{
		AutoSubscribeAudio: true,
		AutoSubscribeVideo: true,
		ClientRole:         agoraservice.ClientRoleBroadcaster,
		ChannelProfile:     agoraservice.ChannelProfileLiveBroadcasting,
	}
	publishConfig := agoraservice.NewRtcConPublishConfig()
	publishConfig.AudioScenario = scenario
	publishConfig.IsPublishAudio = true
	publishConfig.IsPublishVideo = true
	publishConfig.AudioProfile = agoraservice.AudioProfileDefault
	publishConfig.AudioPublishType = agoraservice.AudioPublishTypeNoPublish
	publishConfig.VideoPublishType = agoraservice.VideoPublishTypeNoPublish

	publishConfig.AudioPublishType = agoraservice.AudioPublishTypePcm

	con = agoraservice.NewRtcConnection(conCfg, publishConfig)
	


	con.RegisterObserver(conHandler)
	ret := con.RegisterLocalUserObserver(localUserObserver)
	if ret != 0 {
		fmt.Println("RegisterLocalUserObserver failed, ret ", ret)
		return
	}
	con.RegisterVideoEncodedFrameObserver(encodedVideoObserver)

	con.Connect(token, channelName, userId)
	<-conSignal

	file, err := os.OpenFile("./recv.264", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if file == nil {
		fmt.Println("Error opening file: ", err)
		return
	}
	defer file.Close()

	stop := false
	lastFrameTime := time.Now().UnixMilli()
	srcUid := ""
	for !stop {
		select {
		case videoData := <-videoChan:
			//fmt.Printf("user %s encoded video received, frame len %d\n", videoData.uid, len(videoData.imageData))
			file.Write(videoData.imageData)
			srcUid = videoData.uid
		case <-ctx.Done():
			stop = true
		default:
			time.Sleep(100 * time.Millisecond)
			curTime := time.Now().UnixMilli()
			if curTime - lastFrameTime > 1000 {
				fmt.Printf("send intra request to user %s, time %d\n", srcUid, curTime - lastFrameTime)
				con.SendIntraRequest(srcUid)
				lastFrameTime = curTime
			}
		}
	}

	// release resource

	
	con.Disconnect()
	//<-OnDisconnectedSign
	

	con.Release()

	

	agoraservice.Release()


	localUserObserver = nil
	conHandler = nil
	con = nil

	fmt.Printf("App exited\n")
	
}
