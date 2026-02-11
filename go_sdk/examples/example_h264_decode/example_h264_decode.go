package main


import "C"
import (
	"fmt"
	"os"
	"os/signal"
	"time"

	agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/rtc"

	//rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
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

	// get parameter from arguments： appid, channel_name, output_file(optional)

	argus := os.Args
	if len(argus) < 3 {
		fmt.Println("Usage: program <appid> <channel_name> [output_yuv_file]")
		fmt.Println("This program will receive and decode video from the channel")
		fmt.Println("Default output file: ./received_video.yuv")
		return
	}
	appid := argus[1]
	channelName := argus[2]
	
	// 输出文件路径，支持用户自定义或使用默认值
	outputFile := "./received_video.yuv"
	if len(argus) >= 4 {
		outputFile = argus[3]
	}
	fmt.Printf("Output YUV file: %s\n", outputFile)

	// get environment variable

	//cert := os.Getenv("AGORA_APP_CERTIFICATE")

	userId := "0"
	if appid == "" {
		fmt.Println("Please set AGORA_APP_ID environment variable, and AGORA_APP_CERTIFICATE if needed")
		return
	}
	token := ""
	/*
	if cert != "" {
		tokenExpirationInSeconds := uint32(3600)
		privilegeExpirationInSeconds := uint32(3600)
		var err error
		token, err = rtctokenbuilder.BuildTokenWithUserAccount(appid, cert, channelName, userId,
			rtctokenbuilder.RoleSubscriber, tokenExpirationInSeconds, privilegeExpirationInSeconds)
		if err != nil {
			fmt.Println("Failed to build token: ", err)
			return
		}
	}*/
	svcCfg := agoraservice.NewAgoraServiceConfig()
	svcCfg.AppId = appid
	// whether sending or receiving video, we need to set EnableVideo to true!!
	svcCfg.EnableVideo = true
	svcCfg.LogPath = "./agora_rtc_log/agorasdk.log"
	svcCfg.ConfigDir = "./agora_rtc_log"
	svcCfg.DataDir = "./agora_rtc_log"
	

	agoraservice.Initialize(svcCfg)

	var con *agoraservice.RtcConnection = nil
	scenario := agoraservice.AudioScenarioDefault

	// global private parameters
	parameterHandler := agoraservice.GetAgoraParameter()
	parameterHandler.SetParameters("{\"che.video.minQP\":21}")
	parameterHandler.SetParameters("{\"rtc.video.low_stream_enable_hw_encoder\":false}")
	parameterHandler.SetParameters("{\"engine.video.enable_hw_encoder\":false}")
	parameterHandler.SetParameters("{\"rtc.video.enable_minor_stream_intra_request\":true}")
	parameterHandler.SetParameters("{\"rtc.video.enable_pvc\":false}")
	parameterHandler.SetParameters("{\"rtc.video.color_space_enable\":true}")
	parameterHandler.SetParameters("{\"rtc.video.videoFullrange\":1}")
	parameterHandler.SetParameters("{\"rtc.video.matrixCoefficients\":5}")

	
	
	conSignal := make(chan struct{})
	conHandler := agoraservice.RtcConnectionObserver{
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
		OnAIQoSCapabilityMissing: func(con *agoraservice.RtcConnection, defaultFallbackSenario int) int {
			fmt.Printf("onAIQoSCapabilityMissing, defaultFallbackSenario: %d\n", defaultFallbackSenario)
			return int(agoraservice.AudioScenarioDefault)
		},
	}
	// conCfg.AudioFrameObserver = &agoraservice.RtcConnectionAudioFrameObserver{
	// 	OnPlaybackAudioFrameBeforeMixing: func(con *agoraservice.RtcConnection, channelId string, userId string, frame *agoraservice.PcmAudioFrame) {
	// 		// do something
	// 		fmt.Printf("Playback audio frame before mixing, from userId %s\n", userId)
	// 	},
	// }

	publishConfig := agoraservice.NewRtcConPublishConfig()
	publishConfig.AudioPublishType = agoraservice.AudioPublishTypePcm
	publishConfig.AudioScenario = scenario
	publishConfig.IsPublishAudio = true
	publishConfig.IsPublishVideo = true
	publishConfig.VideoPublishType = agoraservice.VideoPublishTypeYuv
	publishConfig.AudioProfile = agoraservice.AudioProfileDefault

	// 配置为接收模式，不发布音视频
	conCfg := &agoraservice.RtcConnectionConfig{
		AutoSubscribeAudio: true,
		AutoSubscribeVideo: false,
		ClientRole:         agoraservice.ClientRoleBroadcaster, // 设置为观众角色，只接收不发送
		ChannelProfile:     agoraservice.ChannelProfileLiveBroadcasting,
	}

	con = agoraservice.NewRtcConnection(conCfg, publishConfig)
	if con == nil {
		fmt.Println("ERROR: Failed to create RTC connection")
		agoraservice.Release()
		return
	}

	con.RegisterObserver(&conHandler)

	// can update in session life cycle
	encoderCfg := agoraservice.NewVideoEncoderConfiguration()
	encoderCfg.Width = 1280
	encoderCfg.Height = 720
	encoderCfg.Framerate = 24

	encoderCfg.Bitrate = 1250
	encoderCfg.MinBitrate = -1
	encoderCfg.DegradePreference = agoraservice.DegradeMaintainFramerate

	// 注册视频编码帧观察者，用于接收编码的视频帧
	encodedVideoObserver := &agoraservice.VideoEncodedFrameObserver{
		OnEncodedVideoFrame: func(uid string, imageBuffer []byte, frameInfo *agoraservice.EncodedVideoFrameInfo) bool {
			fmt.Printf("Received encoded video from uid: %s, size: %d, width: %d, height: %d, frameType: %d, codecType: %d\n",
				uid, len(imageBuffer), frameInfo.Width, frameInfo.Height, frameInfo.FrameType, frameInfo.CodecType)

			//update encoder configuration
			width := frameInfo.Width
			height := frameInfo.Height
			if width != encoderCfg.Width && height != encoderCfg.Height {
				fmt.Printf("update encoder configuration, width: %d, height: %d\n", width, height)
				encoderCfg.Width = width
				encoderCfg.Height = height
				con.SetVideoEncoderConfiguration(encoderCfg)
				
			}
			// checck codec type: only process h264 codec type
			if frameInfo.CodecType == agoraservice.VideoCodecTypeH264 {

				// forward to a channel
				con.PushVideoEncodedDataForTranscode(imageBuffer, frameInfo)
				
			}

		
			return true
		},
	}
	con.RegisterVideoEncodedFrameObserver(encodedVideoObserver)

	con.Connect(token, channelName, userId)
	<-conSignal

	// set to sub and high video stream,and only recv encoded video frame
	// do not do decoding in sdk!
	con.GetLocalUser().SubscribeAllVideo(&agoraservice.VideoSubscriptionOptions{
		StreamType:       agoraservice.VideoStreamHigh,
		EncodedFrameOnly: true,
	})

	
	
	con.SetVideoEncoderConfiguration(encoderCfg)

	// enable dual video stream
	
	con.SetSimulcastStream(true, &agoraservice.SimulcastStreamConfig{
		Width: 384,
		Height: 216,
		Bitrate: 300,
		Framerate: 10,
	})

	con.PublishVideo()

	fmt.Println("Connected to channel, waiting for video streams...")
	fmt.Println("Press Ctrl+C to stop receiving")

	// 只接收视频流，等待用户中断
	for !(*bStop) {
		time.Sleep(100 * time.Millisecond)
	}


	con.Disconnect()
	con.Release()
	agoraservice.Release()
	fmt.Println("Application terminated")

}
