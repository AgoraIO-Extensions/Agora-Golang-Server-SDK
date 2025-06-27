package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/agoraservice"

	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
)

// convert video frame to external video frame
func ConvertVideoFrameToExternalVideoFrame(frame *agoraservice.VideoFrame) *agoraservice.ExternalVideoFrame {
	bufsize := frame.Width * frame.Height * 3 / 2
	ysize := frame.Width * frame.Height
	merged := make([]byte, 0, bufsize)
	merged = append(merged, frame.YBuffer[:ysize]...)
	merged = append(merged, frame.UBuffer[:ysize/4]...)
	merged = append(merged, frame.VBuffer[:ysize/4]...)
	extFrame := &agoraservice.ExternalVideoFrame{
		Type:      agoraservice.VideoBufferRawData,
		Format:    agoraservice.VideoPixelI420,
		Buffer:    merged,
		Stride:    frame.Width,
		Height:    frame.Height,
		Timestamp: 0,
	}
	//fmt.Printf("ExtFrame, type: %d, w*h: %d*%d, size %d\n", frame.Type, frame.Width, frame.Height, len(merged))
	return extFrame
}

func PushFileToConsumer(file *os.File, audioConsumer *agoraservice.AudioConsumer) {
	buffer := make([]byte, 320)
	for {
		readLen, err := file.Read(buffer)
		if err != nil || readLen < 320 {
			fmt.Printf("read up to EOF,cur read: %d", readLen)
			file.Seek(0, 0)
			break
		}
		audioConsumer.PushPCMData(buffer)
	}
	buffer = nil
}
func ReadFileToConsumer(file *os.File, consumer *agoraservice.AudioConsumer, interval int, done chan bool) {
	for {
		select {
		case <-done:
			fmt.Println("ReadFileToConsumer done")
			return
		default:
			len := consumer.Len()
			if len < 320*interval {
				PushFileToConsumer(file, consumer)
			}
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	}
}

func ConsumeAudio(audioConsumer *agoraservice.AudioConsumer, interval int, done chan bool) {
	for {
		select {
		case <-done:
			fmt.Println("ConsumeAudio done")
			return
		default:
			audioConsumer.Consume()
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	}
}

// generate wave header
// for 16bit pcm data to wav file
/*
func generateWAVHeader(sampleRate int, channels int, pcmDataSize int) []byte {
	totalSize := pcmDataSize + 36 // 数据块+36字节头
	header := make([]byte, 44)

	// RIFF块
	copy(header[0:4], "RIFF")
	binary.LittleEndian.PutUint32(header[4:8], uint32(totalSize))
	copy(header[8:12], "WAVE")

	// fmt子块
	copy(header[12:16], "fmt ")
	binary.LittleEndian.PutUint32(header[16:20], 16)                               // fmt块大小
	binary.LittleEndian.PutUint16(header[20:22], 1)                                // PCM格式, fixed type
	binary.LittleEndian.PutUint16(header[22:24], uint16(channels))                 // 单声道
	binary.LittleEndian.PutUint32(header[24:28], uint32(sampleRate))               // 采样率
	binary.LittleEndian.PutUint32(header[28:32], uint32(sampleRate*channels*16/8)) // 字节率（16000 * 1 * 16/8）
	binary.LittleEndian.PutUint16(header[32:34], 2)                                // 块对齐（1 * 16/8）
	binary.LittleEndian.PutUint16(header[34:36], 16)                               // 位深度

	// data块
	copy(header[36:40], "data")
	binary.LittleEndian.PutUint32(header[40:44], uint32(pcmDataSize))
	return header
}
*/

// sample to recv and echo back yuv and pcm
func main() {
	bStop := new(bool)
	*bStop = false
	// start pprof
	go func() {
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

	// create a queue for yuv and pcm
	yuvQueue := agoraservice.NewQueue(20)
	pcmQueue := agoraservice.NewQueue(10)

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
	
	videoObserver := &agoraservice.VideoFrameObserver{
		OnFrame: func(channelId string, userId string, frame *agoraservice.VideoFrame) bool {
			if frame == nil {
				return true

			}
			
			//fmt.Printf("recv video frame, from channel %s, user %s, type %d, width %d, height %d, stride %d, ysize %d, usize %d, vsize %d\n",channelId, userId, frame.Type, frame.Width, frame.Height, frame.YStride, len(frame.YBuffer), len(frame.UBuffer), len(frame.VBuffer))

			yuvQueue.Enqueue(frame)
			

			return true
		},
	}
	
	var err error
	audioObserver := &agoraservice.AudioFrameObserver{
		OnPlaybackAudioFrameBeforeMixing: func(localUser *agoraservice.LocalUser, channelId string, userId string, frame *agoraservice.AudioFrame, vadResulatState agoraservice.VadState, vadResultFrame *agoraservice.AudioFrame) bool {
			// do something: for play a file
			

			pcmQueue.Enqueue(frame)
			//fmt.Printf("Playback audio frame before mixing, from userId %s, far :%d,rms:%d, pitch: %d\n", userId, frame.FarFieldFlag, frame.Rms, frame.Pitch)
			return true
		},
	}
	localUserObserver := &agoraservice.LocalUserObserver{
		OnLocalAudioTrackStatistics: func(localUser *agoraservice.LocalUser, stats *agoraservice.LocalAudioTrackStats) {
			fmt.Printf("OnLocalAudioTrackStatistics, stats: %v\n", stats)
		},
		OnRemoteAudioTrackStatistics: func(localUser *agoraservice.LocalUser, uid string, stats *agoraservice.RemoteAudioTrackStats) {
			fmt.Printf("OnRemoteAudioTrackStatistics, stats: %v\n", stats)
		},
		OnLocalVideoTrackStatistics: func(localUser *agoraservice.LocalUser, stats *agoraservice.LocalVideoTrackStats) {
			fmt.Printf("OnLocalVideoTrackStatistics, stats: %v\n", stats)
		},
		OnRemoteVideoTrackStatistics: func(localUser *agoraservice.LocalUser, uid string, stats *agoraservice.RemoteVideoTrackStats) {
			fmt.Printf("OnRemoteVideoTrackStatistics, stats: %v\n", stats)
		},
		OnUserAudioTrackStateChanged: func(localUser *agoraservice.LocalUser, uid string, remoteAudioTrack *agoraservice.RemoteAudioTrack, state int, reason int, elapsed int) {
			fmt.Printf("OnUserAudioTrackStateChanged, uid: %s, state: %d, reason: %d, elapsed: %d\n", uid, state, reason, elapsed)
		},
	}

	yuvsender := mediaNodeFactory.NewVideoFrameSender()
	pcmsender := mediaNodeFactory.NewAudioPcmDataSender()

	// open a pcm file for push
	pcmfile, err := os.Open("../test_data/send_audio_16k_1ch.pcm")
	if err != nil {
		fmt.Printf("NewError opening file: %v\n", err)
		return
	}
	defer pcmfile.Close()

	audioConsumer := agoraservice.NewAudioConsumer(pcmsender, 16000, 1)
	defer audioConsumer.Release()
	done := make(chan bool)

	//go ReadFileToConsumer(pcmfile, audioConsumer, 50, done)
	//go ConsumeAudio(audioConsumer, 50, done)

	// create 2 rouite for process audio and video
	audioRoutine := func() {
		for !*bStop {
			AudioFrame := pcmQueue.Dequeue()
			if AudioFrame != nil {
				//fmt.Printf("AudioFrame: %d\n", time.Now().UnixMilli())
				if frame, ok := AudioFrame.(*agoraservice.AudioFrame); ok {
					frame.RenderTimeMs = 0
					ret := pcmsender.SendAudioPcmData(frame)
					if ret != 0 {
						fmt.Printf("Send audio pcm data failed, error code %d\n", ret)
					}
				}
			}
		}
		fmt.Printf("AudioRoutine end\n")
	}
	videoRoutine := func() {
		for !*bStop {
			videoFrame := yuvQueue.Dequeue()
			if videoFrame != nil {
				if frame, ok := videoFrame.(*agoraservice.VideoFrame); ok {
					extFrame := ConvertVideoFrameToExternalVideoFrame(frame)
					yuvsender.SendVideoFrame(extFrame)
				}
			}
		}
		fmt.Printf("VideoRoutine end\n")
	}

	go audioRoutine()
	go videoRoutine()
	fmt.Printf("start audioRoutine: %v, videoRoutine: %v\n", audioRoutine, videoRoutine)

	// step0: create connection
	scenario := svcCfg.AudioScenario
	con := agoraservice.NewRtcConnection(&conCfg, scenario)

	localUser := con.GetLocalUser()
	con.RegisterObserver(conHandler)

	// step1: register video frame observer and video track
	localUser.RegisterVideoFrameObserver(videoObserver)

	track := agoraservice.NewCustomVideoTrackFrame(yuvsender)

	// step2: register audio frame observer and audio track
	localUser.SetPlaybackAudioFrameBeforeMixingParameters(1, 16000)
	audioTrack := agoraservice.NewCustomAudioTrackPcm(pcmsender, scenario)
	localUser.RegisterAudioFrameObserver(audioObserver, 1, nil)

	//localuserobserver
	localUser.RegisterLocalUserObserver(localUserObserver)

	

	con.Connect(token, channelName, userId)
	<-conSignal

	track.SetVideoEncoderConfiguration(&agoraservice.VideoEncoderConfiguration{
		CodecType:         agoraservice.VideoCodecTypeH264,
		Width:             320,
		Height:            240,
		Framerate:         30,
		Bitrate:           1500,
		MinBitrate:        300,
		OrientationMode:   agoraservice.OrientationModeAdaptive,
		DegradePreference: 0,
	})

	// step4: publish video and audio
	track.SetEnabled(true)
	localUser.PublishVideo(track)
	audioTrack.SetEnabled(true)
	localUser.PublishAudio(audioTrack)

	// for yuv test

	

	
	// rgag colos space type test

	for !*bStop {
		
		
		time.Sleep(33 * time.Millisecond)
	}
	close(done)

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
	yuvsender.Release()
	mediaNodeFactory.Release()
	agoraservice.Release()

	track = nil
	videoObserver = nil

	localUser = nil
	conHandler = nil
	con = nil

	fmt.Printf("Disconnected, cost %d ms\n", time.Now().UnixMilli()-start_disconnect)
}
