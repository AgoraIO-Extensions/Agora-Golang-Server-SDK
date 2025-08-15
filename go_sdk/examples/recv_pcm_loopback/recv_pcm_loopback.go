package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"

	agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/rtc"

	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
)

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

	println("Start to send and receive PCM data\nusage:\n	./send_recv_pcm <appid> <channel_name>\n	press ctrl+c to exit\n")

	// get parameter from arguments： appid, channel_name

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
	svcCfg.AppId = appid

	agoraservice.Initialize(svcCfg)
	defer agoraservice.Release()

	var conn *agoraservice.RtcConnection = nil

	audioModel := 1 // 0: direct, 1: channel

	pcmQueue := agoraservice.NewQueue(10)
	audioRoutine := func() {
		for !*bStop {
			AudioFrame := pcmQueue.Dequeue()
			if AudioFrame != nil {
				//fmt.Printf("AudioFrame: %d\n", time.Now().UnixMilli())
				if frame, ok := AudioFrame.(*agoraservice.AudioFrame); ok {
					frame.RenderTimeMs = 0
					ret := conn.PushAudioPcmData(frame.Buffer, frame.SamplesPerSec, frame.Channels, 0)
					if ret != 0 {
						fmt.Printf("Send audio pcm data failed, error code %d\n", ret)
					}
				}
			}
		}
		fmt.Printf("AudioRoutine end\n")
	}

	// a go routine to send audio data to channel

	go audioRoutine()

	scenario := agoraservice.AudioScenarioAiServer
	conCfg := agoraservice.RtcConnectionConfig{
		AutoSubscribeAudio: true,
		AutoSubscribeVideo: false,
		ClientRole:         agoraservice.ClientRoleBroadcaster,
		ChannelProfile:     agoraservice.ChannelProfileLiveBroadcasting,
	}
	publishConfig := agoraservice.NewRtcConPublishConfig()
	publishConfig.AudioPublishType = agoraservice.AudioPublishTypePcm
	publishConfig.AudioScenario = scenario
	publishConfig.IsPublishAudio = true
	publishConfig.IsPublishVideo = false
	publishConfig.AudioProfile = agoraservice.AudioProfileDefault

	conSignal := make(chan struct{})
	OnDisconnectedSign := make(chan struct{})
	conHandler := &agoraservice.RtcConnectionObserver{
		OnConnected: func(con *agoraservice.RtcConnection, info *agoraservice.RtcConnectionInfo, reason int) {
			// do something
			fmt.Printf("Connected, reason %d\n", reason)
			conSignal <- struct{}{}
		},
		OnDisconnected: func(con *agoraservice.RtcConnection, info *agoraservice.RtcConnectionInfo, reason int) {
			// do something
			fmt.Printf("Disconnected, reason %d\n", reason)
			OnDisconnectedSign <- struct{}{}
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
			fmt.Printf("user left: %s, reason = %d\n", uid, reason)
		},
		OnAIQoSCapabilityMissing: func(con *agoraservice.RtcConnection, defaultFallbackSenario int) int {
			fmt.Printf("onAIQoSCapabilityMissing, defaultFallbackSenario: %d\n", defaultFallbackSenario)
			return int(agoraservice.AudioScenarioDefault)
		},
	}
	audioObserver := &agoraservice.AudioFrameObserver{
		OnPlaybackAudioFrameBeforeMixing: func(localUser *agoraservice.LocalUser, channelId string, userId string, frame *agoraservice.AudioFrame, vadResultState agoraservice.VadState, vadResultFrame *agoraservice.AudioFrame) bool {
			// do something
			//fmt.Printf("Playback audio frame before mixing, from userId %s\n", userId)
			if audioModel == 1 {
				pcmQueue.Enqueue(frame)
				fmt.Printf("Enqueue audio frame, size %d\n", pcmQueue.Size())
			}
			if audioModel == 0 {
				if conn != nil {
					conn.PushAudioPcmData(frame.Buffer, frame.SamplesPerSec, frame.Channels,0)
				}
			}
			return true
		},
	}

	conn = agoraservice.NewRtcConnection(&conCfg, publishConfig)

	//added by wei for localuser observer
	localUserObserver := &agoraservice.LocalUserObserver{
		OnStreamMessage: func(localUser *agoraservice.LocalUser, uid string, streamId int, data []byte) {
			// do something
			fmt.Printf("*****Stream message, from userId %s\n", uid)
		},

		OnAudioVolumeIndication: func(localUser *agoraservice.LocalUser, audioVolumeInfo []*agoraservice.AudioVolumeInfo, speakerNumber int, totalVolume int) {
			// do something
			fmt.Printf("*****Audio volume indication, speaker number %d\n", speakerNumber)
		},
		OnAudioPublishStateChanged: func(localUser *agoraservice.LocalUser, channelId string, oldState int, newState int, elapse_since_last_state int) {
			fmt.Printf("*****Audio publish state changed, old state %d, new state %d\n", oldState, newState)
		},
		OnUserInfoUpdated: func(localUser *agoraservice.LocalUser, uid string, userMediaInfo int, val int) {
			fmt.Printf("*****User info updated, uid %s\n", uid)
		},
		OnUserAudioTrackSubscribed: func(localUser *agoraservice.LocalUser, uid string, remoteAudioTrack *agoraservice.RemoteAudioTrack) {
			fmt.Printf("*****User audio track subscribed, uid %s\n", uid)
		},
		OnUserVideoTrackSubscribed: func(localUser *agoraservice.LocalUser, uid string, info *agoraservice.VideoTrackInfo, remoteVideoTrack *agoraservice.RemoteVideoTrack) {

		},
		OnUserAudioTrackStateChanged: func(localUser *agoraservice.LocalUser, uid string, remoteAudioTrack *agoraservice.RemoteAudioTrack, state int, reason int, elapsed int) {
			fmt.Printf("*****User audio track state changed, uid %s\n", uid)
		},
		OnUserVideoTrackStateChanged: func(localUser *agoraservice.LocalUser, uid string, remoteAudioTrack *agoraservice.RemoteVideoTrack, state int, reason int, elapsed int) {
			fmt.Printf("*****User video track state changed, uid %s\n", uid)
		},
	}

	conn.RegisterObserver(conHandler)

	localUser := conn.GetLocalUser()
	//localUser.SetAudioScenario(agoraservice.AudioScenarioChorus)
	conn.Connect(token, channelName, userId)
	<-conSignal

	localUser = conn.GetLocalUser()
	localUser.SetPlaybackAudioFrameBeforeMixingParameters(1, 16000)
	conn.RegisterLocalUserObserver(localUserObserver)

	conn.RegisterAudioFrameObserver(audioObserver, 0, nil)

	conn.PublishAudio()

	/*
		// disalbe pre-load audio data from version 2.1.x, by wei
		// all use build-in function for low-latency

		frame := &agoraservice.AudioFrame{
			Type:              agoraservice.AudioFrameTypePCM16,
			SamplesPerChannel: 160,
			BytesPerSample:    2,
			Channels:          1,
			SamplesPerSec:     16000,
			Buffer:            make([]byte, 320),
			RenderTimeMs:      0,
		}

		file, err := os.Open("../test_data/send_audio_16k_1ch.pcm")
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		//track.AdjustPublishVolume(100)


		sendCount := 0

		// send 180ms audio data
		for i := 0; i < 18; i++ {
			dataLen, err := file.Read(frame.Buffer)
			if err != nil || dataLen < 320 {
				fmt.Println("Finished reading file:", err)
				break
			}
			sendCount++
			ret := sender.SendAudioPcmData(frame)
			fmt.Printf("SendAudioPcmData %d ret: %d\n", sendCount, ret)
		}
	*/

	//added by wei for loop back
	for !(*bStop) {
		time.Sleep(40 * time.Millisecond)
	}
	//end of added by wei for loop back
	fmt.Printf("stop now.....\n")

	//release operation:cancel defer release,try manual release

	start_disconnect := time.Now().UnixMilli()
	conn.Disconnect()
	<-OnDisconnectedSign

	fmt.Printf("Disconnect, cost %d ms\n", time.Now().UnixMilli()-start_disconnect)
	//a, b, c, d := agoraservice.GetMapInfo()
	//fmt.Printf("mapinfo:: %d,%d,%d,%d\n", a, b, c, d)

	conn.Release()
	//a, b, c, d = agoraservice.GetMapInfo()
	//fmt.Printf("mapinfo:: %d,%d,%d,%d\n", a, b, c, d)

	agoraservice.Release()

	audioObserver = nil
	localUserObserver = nil
	localUser = nil
	conHandler = nil
	conn = nil

	fmt.Printf("Disconnected, cost %d ms\n", time.Now().UnixMilli()-start_disconnect)
}
