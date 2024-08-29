package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"agora.io/agoraservice"

	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
)

type ConnectionRole int

const (
	ConnectionRolePublisher  ConnectionRole = 0
	ConnectionRoleSubscriber ConnectionRole = 1
)

type GlobalContext struct {
	ctx               *context.Context
	cancel            *context.CancelFunc
	termSingal        chan os.Signal
	appId             string
	cert              string
	channelNamePrefix string
	mediaNodeFactory  *agoraservice.MediaNodeFactory
	waitTasks         *sync.WaitGroup
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

func (globalCtx *GlobalContext) startTask(id int) {
	defer globalCtx.waitTasks.Done()

	ctx := globalCtx.ctx
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
	}
	// senderLocalUserObs := &agoraservice.LocalUserObserver{}
	senderCon := agoraservice.NewRtcConnection(&agoraservice.RtcConnectionConfig{
		AutoSubscribeAudio: false,
		AutoSubscribeVideo: false,
		ClientRole:         agoraservice.ClientRoleBroadcaster,
		ChannelProfile:     agoraservice.ChannelProfileLiveBroadcasting,
	})
	defer senderCon.Release()

	pcmSender := globalCtx.mediaNodeFactory.NewAudioPcmDataSender()
	defer pcmSender.Release()
	audioTrack := agoraservice.NewCustomAudioTrackPcm(pcmSender)
	defer audioTrack.Release()

	senderCon.RegisterObserver(senderConObs)
	senderCon.Connect(token1, channelName, senderId)
	defer senderCon.Disconnect()

	select {
	case <-conSignal:
	case <-time.After(5 * time.Second):
		fmt.Printf("sender failed to connect, task %d\n", id)
		return
	}

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

			case <-(*ctx).Done():
				return
			}
		}
	}()

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
	// receiverLocalUserObs := &agoraservice.LocalUserObserver{}
	recieverConAudioFrameObs := &agoraservice.AudioFrameObserver{
		OnPlaybackAudioFrameBeforeMixing: func(localUser *agoraservice.LocalUser, channelId string, userId string, frame *agoraservice.PcmAudioFrame) bool {
			// do something
			fmt.Printf("Playback audio frame before mixing, from userId %s\n", userId)
			return true
		},
	}
	// receiverConVideoFrameObs := &agoraservice.VideoFrameObserver{}
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
	localUser.RegisterAudioFrameObserver(recieverConAudioFrameObs)

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
		ctx:               &ctx,
		cancel:            &cancel,
		termSingal:        c,
		appId:             appid,
		cert:              cert,
		channelNamePrefix: channelName,
		mediaNodeFactory:  mediaNodeFactory,
		waitTasks:         &sync.WaitGroup{},
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

	for i := 0; i < 1; i++ {
		globalCtx.waitTasks.Add(1)
		go globalCtx.startTask(i)
	}

	<-globalCtx.termSingal
	(*globalCtx.cancel)()
	globalCtx.waitTasks.Wait()
}
