package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"math/rand"

	"agora.io/agoraservice"

	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
)

const (
	ConnectionCount = 3
)

type GlobalContext struct {
	ctx               context.Context
	cancel            context.CancelFunc
	termSingal        chan os.Signal
	appId             string
	cert              string
	channelNamePrefix string
	mediaNodeFactory  *agoraservice.MediaNodeFactory
	waitTasks         *sync.WaitGroup
	connectionCancels []context.CancelFunc
}

func (globalCtx *GlobalContext) genToken(channelName string, userId string) (string, error) {
	token := ""
	if globalCtx.cert != "" {
		tokenExpirationInSeconds := uint32(24 * 3600)
		privilegeExpirationInSeconds := uint32(24 * 3600)
		var err error
		token, err = rtctokenbuilder.BuildTokenWithUserAccount(globalCtx.appId, globalCtx.cert, channelName, userId,
			rtctokenbuilder.RolePublisher, tokenExpirationInSeconds, privilegeExpirationInSeconds)
		if err != nil {
			fmt.Println("Failed to build token: ", err)
			return token, err
		}
	}
	return token, nil
}

func (globalCtx *GlobalContext) startTask(ctx context.Context, id int) {
	defer globalCtx.waitTasks.Done()

	channelName := fmt.Sprintf("%s%d", globalCtx.channelNamePrefix, id)
	senderId := fmt.Sprintf("%d", id*1000+1)
	receiverId := fmt.Sprintf("%d", id*1000+2)
	token1, err1 := globalCtx.genToken(channelName, senderId)
	token2, err2 := globalCtx.genToken(channelName, receiverId)
	if err1 != nil || err2 != nil {
		fmt.Printf("Failed to generate token, task %d\n", id)
		return
	}
	conSignal := make(chan struct{})
	senderConObs := &agoraservice.RtcConnectionObserver{
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
		OnStreamMessageError: func(con *agoraservice.RtcConnection, uid string, streamId int, errCode int, missed int, cached int) {
			fmt.Printf("send stream message error: %d, channel %s, uid %s\n", errCode, channelName, uid)
		},
	}
	// senderLocalUserObs := &agoraservice.LocalUserObserver{}
	senderCon := agoraservice.NewRtcConnection(&agoraservice.RtcConnectionConfig{
		AutoSubscribeAudio: false,
		AutoSubscribeVideo: false,
		ClientRole:         agoraservice.ClientRoleBroadcaster,
		ChannelProfile:     agoraservice.ChannelProfileLiveBroadcasting,
	})
	defer senderCon.Release()
	// create audio track
	pcmSender := globalCtx.mediaNodeFactory.NewAudioPcmDataSender()
	defer pcmSender.Release()
	audioTrack := agoraservice.NewCustomAudioTrackPcm(pcmSender)
	defer audioTrack.Release()
	// create video track
	yuvSender := globalCtx.mediaNodeFactory.NewVideoFrameSender()
	defer yuvSender.Release()
	videoTrack := agoraservice.NewCustomVideoTrackFrame(yuvSender)
	defer videoTrack.Release()
	// create datastream
	streamId, errCode := senderCon.CreateDataStream(false, false)
	if errCode != 0 {
		fmt.Printf("Failed to create data stream: %d, channel %s\n", errCode, channelName)
	}

	senderCon.RegisterObserver(senderConObs)
	senderCon.Connect(token1, channelName, senderId)
	defer senderCon.Disconnect()

	select {
	case <-conSignal:
	case <-time.After(5 * time.Second):
		fmt.Printf("sender failed to connect, task %d\n", id)
		return
	}

	// send audio
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		audioTrack.SetEnabled(true)
		senderLocalUser := senderCon.GetLocalUser()
		senderLocalUser.PublishAudio(audioTrack)
		defer func() {
			senderLocalUser.UnpublishAudio(audioTrack)
			audioTrack.SetEnabled(false)
		}()

		audioTrack.AdjustPublishVolume(100)

		frame := agoraservice.PcmAudioFrame{
			Data:              make([]byte, 320),
			Timestamp:         0,
			SamplesPerChannel: 160,
			BytesPerSample:    2,
			NumberOfChannels:  1,
			SampleRate:        16000,
		}

		file, err := os.Open("../../../test_data/demo.pcm")
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		ticker := time.NewTicker(50 * time.Millisecond)
		sendCount := 0
		firstSendTime := time.Now()
		for {
			select {
			case <-ticker.C:
				shouldSendCount := int(time.Since(firstSendTime).Milliseconds()/10) - (sendCount - 18)
				for i := 0; i < shouldSendCount; i++ {
					dataLen, err := file.Read(frame.Data)
					if err != nil || dataLen < 320 {
						fmt.Println("Finished reading file:", err)
						file.Seek(0, 0)
						i--
						continue
					}

					sendCount++
					ret := pcmSender.SendAudioPcmData(&frame)
					fmt.Printf("SendAudioPcmData %d ret: %d\n", sendCount, ret)
				}
				fmt.Printf("Sent %d frames this time\n", shouldSendCount)
			case <-ctx.Done():
				fmt.Printf("task %d audio sender finished\n", id)
				return
			}
		}
	}()

	// send video
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		videoTrack.SetVideoEncoderConfiguration(&agoraservice.VideoEncoderConfiguration{
			CodecType:         agoraservice.VideoCodecTypeH264,
			Width:             320,
			Height:            240,
			Framerate:         30,
			Bitrate:           500,
			MinBitrate:        100,
			OrientationMode:   agoraservice.VideoOrientation0,
			DegradePreference: 0,
		})
		videoTrack.SetEnabled(true)
		senderLocalUser := senderCon.GetLocalUser()
		senderLocalUser.PublishVideo(videoTrack)
		defer func() {
			senderLocalUser.UnpublishVideo(videoTrack)
			videoTrack.SetEnabled(false)
		}()

		w := 416
		h := 240
		dataSize := w * h * 3 / 2
		data := make([]byte, dataSize)
		// read yuv from file 103_RaceHorses_416x240p30_300.yuv
		file, err := os.Open("../../../test_data/RaceHorses_416x240p30_300.yuv")
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		ticker := time.NewTicker(33 * time.Millisecond)
		for {
			select {
			case <-ticker.C:
				dataLen, err := file.Read(data)
				if err != nil || dataLen < dataSize {
					file.Seek(0, 0)
					continue
				}
				// senderCon.SendStreamMessage(streamId, data)
				yuvSender.SendVideoFrame(&agoraservice.VideoFrame{
					Buffer:    data,
					Width:     w,
					Height:    h,
					YStride:   w,
					UStride:   w / 2,
					VStride:   w / 2,
					Timestamp: 0,
				})
			case <-ctx.Done():
				fmt.Printf("task %d video sender finished\n", id)
				return
			}
		}
	}()

	// send datastream
	if streamId >= 0 {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			ticker := time.NewTicker(33 * time.Millisecond)
			msg := []byte(fmt.Sprintf("Hello, Agora! from task %d", id))
			for {
				select {
				case <-ticker.C:
					ret := senderCon.SendStreamMessage(streamId, msg)
					fmt.Printf("SendStreamMessage ret: %d, task %d\n", ret, id)
				case <-ctx.Done():
					fmt.Printf("task %d data stream sender finished\n", id)
					return
				}
			}
		}()
	}

	receiverConObs := &agoraservice.RtcConnectionObserver{
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
	receiverLocalUserObs := &agoraservice.LocalUserObserver{
		OnStreamMessage: func(localUser *agoraservice.LocalUser, uid string, streamId int, data []byte) {
			fmt.Printf("recv stream message: %s, channel %s, uid %s\n", string(data), channelName, uid)
		},
	}
	recieverConAudioFrameObs := &agoraservice.AudioFrameObserver{
		OnPlaybackAudioFrameBeforeMixing: func(localUser *agoraservice.LocalUser, channelId string, userId string, frame *agoraservice.PcmAudioFrame) bool {
			// do something
			fmt.Printf("Playback audio frame before mixing, from channel %s, userId %s, audio duration %dms\n",
				channelId, userId, frame.SamplesPerChannel*1000/frame.SampleRate)
			return true
		},
	}
	receiverConVideoFrameObs := &agoraservice.VideoFrameObserver{
		OnFrame: func(localUser *agoraservice.LocalUser, channelId string, userId string, frame *agoraservice.VideoFrame) bool {
			// do something
			fmt.Printf("recv video frame, from channel %s, user %s, video size %dx%d\n",
				channelId, userId, frame.Width, frame.Height)
			return true
		},
	}
	receiverCon := agoraservice.NewRtcConnection(&agoraservice.RtcConnectionConfig{
		AutoSubscribeAudio: true,
		AutoSubscribeVideo: true,
		ClientRole:         agoraservice.ClientRoleBroadcaster,
		ChannelProfile:     agoraservice.ChannelProfileLiveBroadcasting,
	})
	defer receiverCon.Release()

	receiverCon.RegisterObserver(receiverConObs)
	localUser := receiverCon.GetLocalUser()
	localUser.SetPlaybackAudioFrameBeforeMixingParameters(1, 16000)
	localUser.RegisterLocalUserObserver(receiverLocalUserObs)
	localUser.RegisterAudioFrameObserver(recieverConAudioFrameObs)
	localUser.RegisterVideoFrameObserver(receiverConVideoFrameObs)

	receiverCon.Connect(token2, channelName, receiverId)
	select {
	case <-conSignal:
	case <-time.After(5 * time.Second):
		fmt.Printf("receiver failed to connect, task %d\n", id)
	}
	defer receiverCon.Disconnect()

	waitGroup.Wait()
	fmt.Printf("task %d finished\n", id)
}

func globalInit() *GlobalContext {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	// catch ternimal signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// get environment variable
	appid := os.Getenv("AGORA_APP_ID")
	cert := os.Getenv("AGORA_APP_CERTIFICATE")
	channelName := "gosdktest"
	if appid == "" {
		fmt.Println("Please set AGORA_APP_ID environment variable, and AGORA_APP_CERTIFICATE if needed")
		return nil
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
		LogSize:        1024 * 1024,
	}
	agoraservice.Initialize(&svcCfg)
	mediaNodeFactory := agoraservice.NewMediaNodeFactory()
	return &GlobalContext{
		ctx:               ctx,
		cancel:            cancel,
		termSingal:        c,
		appId:             appid,
		cert:              cert,
		channelNamePrefix: channelName,
		mediaNodeFactory:  mediaNodeFactory,
		waitTasks:         &sync.WaitGroup{},
		connectionCancels: make([]context.CancelFunc, ConnectionCount),
	}
}

func (ctx *GlobalContext) release() {
	ctx.mediaNodeFactory.Release()
	agoraservice.Release()
}

func main() {
	globalCtx := globalInit()
	if globalCtx == nil {
		return
	}
	defer globalCtx.release()

	for i := 0; i < ConnectionCount; i++ {
		globalCtx.waitTasks.Add(1)
		ctx, cancel := context.WithCancel(context.Background())
		globalCtx.connectionCancels[i] = cancel
		go globalCtx.startTask(ctx, i)
	}

	globalCtx.waitTasks.Add(1)
	go func() {
		defer globalCtx.waitTasks.Done()
		stop := false
		for !stop {
			// random select a connection to stop and start new task
			randTime := time.Duration(5+rand.Intn(10)) * time.Second
			randIndex := rand.Intn(ConnectionCount - 1)
			select {
			case <-time.After(randTime):
				globalCtx.connectionCancels[randIndex]()
				ctx, cancel := context.WithCancel(context.Background())
				globalCtx.connectionCancels[randIndex] = cancel
				globalCtx.waitTasks.Add(1)
				go globalCtx.startTask(ctx, randIndex)
			case <-globalCtx.ctx.Done():
				stop = true
				for _, cancel := range globalCtx.connectionCancels {
					cancel()
				}
			}
		}
	}()

	<-globalCtx.termSingal
	globalCtx.cancel()
	globalCtx.waitTasks.Wait()
}
