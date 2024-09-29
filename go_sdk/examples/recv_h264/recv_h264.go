package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"agora.io/agoraservice"

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
	svcCfg := agoraservice.NewAgoraServiceConfig()
	svcCfg.EnableVideo = true
	svcCfg.AppId = appid

	agoraservice.Initialize(svcCfg)
	defer agoraservice.Release()
	mediaNodeFactory := agoraservice.NewMediaNodeFactory()
	defer mediaNodeFactory.Release()

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

	videoChan := make(chan *EncodedVideoData, 100)
	localUserObserver := &agoraservice.LocalUserObserver{
		OnUserVideoTrackSubscribed: func(localUser *agoraservice.LocalUser, uid string, info *agoraservice.VideoTrackInfo, remoteVideoTrack *agoraservice.RemoteVideoTrack) {
			fmt.Printf("user %s video subscribed\n", uid)
			remoteVideoTrack.RegisterVideoEncodedImageReceiver(&agoraservice.VideoEncodedImageReceiver{
				OnEncodedVideoFrame: func(receiver *agoraservice.VideoEncodedImageReceiver, uid string, imageBuffer []byte, frameInfo *agoraservice.EncodedVideoFrameInfo) bool {
					// fmt.Printf("user %s encoded video received\n", uid)
					videoChan <- &EncodedVideoData{
						uid:       uid,
						imageData: imageBuffer,
						frameInfo: frameInfo,
					}
					return true
				},
			})
		},
	}
	con := agoraservice.NewRtcConnection(&conCfg)
	defer con.Release()

	localUser := con.GetLocalUser()
	con.RegisterObserver(conHandler)
	ret := localUser.RegisterLocalUserObserver(localUserObserver)
	if ret != 0 {
		fmt.Println("RegisterLocalUserObserver failed, ret ", ret)
		return
	}

	con.Connect(token, channelName, userId)
	<-conSignal

	file, err := os.OpenFile("./recv.264", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if file == nil {
		fmt.Println("Error opening file: ", err)
		return
	}
	defer file.Close()

	stop := false
	for !stop {
		select {
		case videoData := <-videoChan:
			fmt.Printf("user %s encoded video received, frame len %d\n", videoData.uid, len(videoData.imageData))
			file.Write(videoData.imageData)
		case <-ctx.Done():
			stop = true
		}
	}
	con.Disconnect()
}
