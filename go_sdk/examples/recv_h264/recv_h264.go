package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"
	"sync"
	agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/rtc"

	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
)

type EncodedVideoData struct {
	uid       string
	imageData []byte
	frameInfo *agoraservice.EncodedVideoFrameInfo
}
// VideoFileWriter 管理每个用户的视频文件写入
type VideoFileWriter struct {
	files map[string]*os.File
	mutex sync.Mutex
}

// NewVideoFileWriter 创建新的视频文件写入器
func NewVideoFileWriter() *VideoFileWriter {
	return &VideoFileWriter{
		files: make(map[string]*os.File),
	}
}

// WriteVideoData 根据用户名写入视频数据
func (w *VideoFileWriter) WriteVideoData(username string, imageData []byte) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// 检查该用户的文件是否已打开
	file, exists := w.files[username]
	if !exists {
		// 创建新文件
		filename := fmt.Sprintf("./recv_%s.h264", username)
		var err error
		file, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("failed to open file for user %s: %w", username, err)
		}
		w.files[username] = file
		fmt.Printf("创建文件: %s\n", filename)
	}

	// 写入数据
	_, err := file.Write(imageData)
	if err != nil {
		return fmt.Errorf("failed to write data for user %s: %w", username, err)
	}

	return nil
}

// CloseAll 关闭所有打开的文件
func (w *VideoFileWriter) CloseAll() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	for username, file := range w.files {
		file.Close()
		fmt.Printf("关闭文件: %s.h264\n", username)
	}
	w.files = make(map[string]*os.File)
}

// CloseUser 关闭特定用户的文件
func (w *VideoFileWriter) CloseUser(username string) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	file, exists := w.files[username]
	if !exists {
		return fmt.Errorf("no file open for user %s", username)
	}

	err := file.Close()
	delete(w.files, username)
	fmt.Printf("关闭文件: %s.h264\n", username)
	return err
}

var videoFileWriter *VideoFileWriter
func main() {
	bStop := new(bool)
	*bStop = false
	videoFileWriter = NewVideoFileWriter()
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

	cert = ""
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
	parseSei := false
	encodedVideoObserver := &agoraservice.VideoEncodedFrameObserver{
		OnEncodedVideoFrame: func(uid string, imageBuffer []byte, frameInfo *agoraservice.EncodedVideoFrameInfo) bool {
			// fmt.Printf("user %s encoded video received\n", uid)
			//fmt.Printf("user %s encoded video received, frame type %d\n", uid, frameInfo.FrameType)
			if frameInfo.FrameType == agoraservice.VideoFrameTypeKeyFrame {
				fmt.Printf("key frame received, time %d, codec type %d\n", time.Now().UnixMilli()-lastKeyFrameTime, frameInfo.CodecType)
				lastKeyFrameTime = time.Now().UnixMilli()
			}

				//fmt.Printf("imageBuffer len: %x \n", imageBuffer)
			
			// anyway ,parse sei
			if parseSei {
				seiData, seiLength := agoraservice.FindSEI(imageBuffer, frameInfo.CodecType)
				if seiData != nil {
					fmt.Printf("**********sei data: %s, sei length: %d\n", seiData, seiLength)
				} else {
					//fmt.Printf("failed to parse sei data, sei length: %d, codec type: %d,len: %d\n", seiLength, frameInfo.CodecType, len(imageBuffer))
				}
			}
			
			// end
			videoChan <- &EncodedVideoData{
				uid:       uid,
				imageData: imageBuffer,
				frameInfo: frameInfo,
			}
			return true
		},
	}
	scenario := agoraservice.AudioScenarioDefault
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
	subvideoopt := &agoraservice.VideoSubscriptionOptions{
		StreamType:       agoraservice.VideoStreamHigh,
		EncodedFrameOnly: true,
	}
	con.GetLocalUser().SubscribeAllVideo(subvideoopt)
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
			//file.Write(videoData.imageData)
			videoFileWriter.WriteVideoData(videoData.uid, videoData.imageData)
			srcUid = videoData.uid
		case <-ctx.Done():
			stop = true
		default:
			time.Sleep(100 * time.Millisecond)
			curTime := time.Now().UnixMilli()
			if curTime - lastFrameTime > 1000 {
				fmt.Printf("send intra request to user %s, time %d\n", srcUid, curTime - lastFrameTime)
				//con.SendIntraRequest(srcUid)
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
