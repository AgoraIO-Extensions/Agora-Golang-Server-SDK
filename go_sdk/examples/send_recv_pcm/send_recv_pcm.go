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
	/*
		argus := os.Args
		if len(argus) < 3 {
			fmt.Println("Please input appid, channel name")
			return
		}
		appid := argus[1]
		channelName := argus[2]
	*/
	appid := "aab8b8f5a8cd4469a63042fcfafe7063"
	channelName := "wei129"

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
	mediaNodeFactory := agoraservice.NewMediaNodeFactory()
	defer mediaNodeFactory.Release()

	conCfg := agoraservice.RtcConnectionConfig{
		AutoSubscribeAudio: true,
		AutoSubscribeVideo: false,
		ClientRole:         agoraservice.ClientRoleBroadcaster,
		ChannelProfile:     agoraservice.ChannelProfileLiveBroadcasting,
	}
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
			fmt.Println("user left, " + uid)
		},
	}
	audioObserver := &agoraservice.AudioFrameObserver{
		OnPlaybackAudioFrameBeforeMixing: func(localUser *agoraservice.LocalUser, channelId string, userId string, frame *agoraservice.AudioFrame) bool {
			// do something
			fmt.Printf("Playback audio frame before mixing, from userId %s\n", userId)
			return true
		},
	}

	con := agoraservice.NewRtcConnection(&conCfg)
	defer con.Release()

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

	con.RegisterObserver(conHandler)

	//end

	// sender := con.NewPcmSender()
	// defer sender.Release()
	sender := mediaNodeFactory.NewAudioPcmDataSender()
	defer sender.Release()
	track := agoraservice.NewCustomAudioTrackPcm(sender)
	defer track.Release()

	localUser := con.GetLocalUser()
	localUser.SetAudioScenario(agoraservice.AudioScenarioChorus)
	con.Connect(token, channelName, userId)
	<-conSignal

	localUser = con.GetLocalUser()
	localUser.SetPlaybackAudioFrameBeforeMixingParameters(1, 16000)
	localUser.SetAudioVolumeIndicationParameters(300, 1, true)
	localUser.RegisterLocalUserObserver(localUserObserver)

	localUser.RegisterAudioFrameObserver(audioObserver)

	track.SetEnabled(true)
	localUser.PublishAudio(track)

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

	track.AdjustPublishVolume(100)

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

	firstSendTime := time.Now()
	for !(*bStop) {
		shouldSendCount := int(time.Since(firstSendTime).Milliseconds()/10) - (sendCount - 18)
		for i := 0; i < shouldSendCount; i++ {
			dataLen, err := file.Read(frame.Buffer)
			if err != nil || dataLen < 320 {
				fmt.Println("Finished reading file:", err)
				file.Seek(0, 0)
				i--
				continue
			}

			sendCount++
			ret := sender.SendAudioPcmData(frame)
			fmt.Printf("SendAudioPcmData %d ret: %d\n", sendCount, ret)
		}
		fmt.Printf("Sent %d frames this time\n", shouldSendCount)
		time.Sleep(50 * time.Millisecond)
	}

	//release operation:cancel defer release,try manual release

	localUser.UnpublishAudio(track)
	track.SetEnabled(false)
	localUser.UnregisterAudioFrameObserver()
	localUser.UnregisterAudioFrameObserver()
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
	audioObserver = nil
	localUserObserver = nil
	localUser = nil
	conHandler = nil
	con = nil

	fmt.Printf("Disconnected, cost %d ms\n", time.Now().UnixMilli()-start_disconnect)
}
