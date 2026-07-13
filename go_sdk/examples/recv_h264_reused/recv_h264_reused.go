package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/rtc"
	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
)

type ReusedEncodedVideoData struct {
	uid       string
	imageData []byte
	frameInfo *agoraservice.EncodedVideoFrameInfo
}

type ReusedVideoFileWriter struct {
	files map[string]*os.File
	mutex sync.Mutex
}

func NewReusedVideoFileWriter() *ReusedVideoFileWriter {
	return &ReusedVideoFileWriter{files: make(map[string]*os.File)}
}

func (w *ReusedVideoFileWriter) WriteVideoData(username string, imageData []byte) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	file, exists := w.files[username]
	if !exists {
		filename := fmt.Sprintf("./recv_%s.h264", username)
		var err error
		file, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		w.files[username] = file
	}
	_, err := file.Write(imageData)
	return err
}

func (w *ReusedVideoFileWriter) CloseAll() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	for _, file := range w.files {
		file.Close()
	}
	w.files = make(map[string]*os.File)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		cancel()
	}()

	args := os.Args
	if len(args) < 3 {
		fmt.Println("Please input appid, channel name")
		return
	}
	appid := args[1]
	channelName := args[2]
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
	svcCfg.LogPath = "./agora_rtc_log/agrasdk.log"
	agoraservice.Initialize(svcCfg)

	writer := NewReusedVideoFileWriter()
	defer writer.CloseAll()

	conSignal := make(chan struct{})
	conHandler := &agoraservice.RtcConnectionObserver{
		OnConnected: func(con *agoraservice.RtcConnection, info *agoraservice.RtcConnectionInfo, reason int) {
			fmt.Printf("Connected, reason %d\n", reason)
			conSignal <- struct{}{}
		},
		OnConnectionFailure: func(con *agoraservice.RtcConnection, conInfo *agoraservice.RtcConnectionInfo, errCode int) {
			fmt.Printf("Connection failure, error code %d\n", errCode)
		},
	}

	videoChan := make(chan *ReusedEncodedVideoData, 100)
	lastKeyFrameTime := time.Now().UnixMilli()
	encodedVideoObserver := &agoraservice.VideoEncodedFrameObserver{
		OnReusedEncodedVideoFrame: func(uid string, imageBuffer []byte, frameInfo *agoraservice.EncodedVideoFrameInfo) bool {
			fmt.Printf("user %s encoded video received, frame type %d\n", uid, frameInfo.FrameType)
			if frameInfo.FrameType == agoraservice.VideoFrameTypeKeyFrame {
				fmt.Printf("key frame received, time %d, codec type %d\n", time.Now().UnixMilli()-lastKeyFrameTime, frameInfo.CodecType)
				lastKeyFrameTime = time.Now().UnixMilli()
			}
			imageCopy := append([]byte(nil), imageBuffer...)
			videoChan <- &ReusedEncodedVideoData{
				uid:       uid,
				imageData: imageCopy,
				frameInfo: frameInfo,
			}
			return true
		},
	}

	conCfg := &agoraservice.RtcConnectionConfig{
		AutoSubscribeAudio: true,
		AutoSubscribeVideo: true,
		ClientRole:         agoraservice.ClientRoleBroadcaster,
		ChannelProfile:     agoraservice.ChannelProfileLiveBroadcasting,
	}
	publishConfig := agoraservice.NewRtcConPublishConfig()
	publishConfig.AudioScenario = agoraservice.AudioScenarioDefault
	publishConfig.IsPublishAudio = true
	publishConfig.IsPublishVideo = true
	publishConfig.AudioProfile = agoraservice.AudioProfileDefault
	publishConfig.AudioPublishType = agoraservice.AudioPublishTypeNoPublish
	publishConfig.VideoPublishType = agoraservice.VideoPublishTypeEncodedImage
	publishConfig.AudioPublishType = agoraservice.AudioPublishTypePcm

	con := agoraservice.NewRtcConnection(conCfg, publishConfig)
	subvideoopt := &agoraservice.VideoSubscriptionOptions{
		StreamType:       agoraservice.VideoStreamHigh,
		EncodedFrameOnly: true,
	}
	con.GetLocalUser().SubscribeAllVideo(subvideoopt)
	con.RegisterObserver(conHandler)
	con.RegisterVideoEncodedFrameObserver(encodedVideoObserver)
	con.Connect(token, channelName, userId)
	<-conSignal

	stop := false
	for !stop {
		select {
		case videoData := <-videoChan:
			_ = writer.WriteVideoData(videoData.uid, videoData.imageData)
		case <-ctx.Done():
			stop = true
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}

	con.Disconnect()
	con.Release()
	agoraservice.Release()
}
