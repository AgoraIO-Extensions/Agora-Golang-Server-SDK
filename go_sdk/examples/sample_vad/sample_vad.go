package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/agoraservice"
	//"google.golang.org/protobuf/types/known/sourcecontextpb"

	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
)

func main() {
	bStop := new(bool)
	*bStop = false

	// only for debug
	vadDump := agoraservice.NewVadDump("./agora_rtc_log/")
	vadDump.Open()
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

	mediaNodeFactory := agoraservice.NewMediaNodeFactory()

	//agoraservice.EnableExtension("agora.builtin", "agora_audio_label_generator", "", true)

	// Recommended  configurations:
	// For not-so-noisy environments, use this configuration: (16, 30, 20, 0.7, 0.5, 70, -50, 70, -50)
	// For noisy environments, use this configuration: (16, 30, 20, 0.7, 0.5, 70, -40, 70, -40)
	// For high-noise environments, use this configuration: (16, 30, 20, 0.7, 0.5, 70, -30, 70, -30)
	vad := agoraservice.NewAudioVadV2(
		&agoraservice.AudioVadConfigV2{
			PreStartRecognizeCount: 16,
			StartRecognizeCount:    30,
			StopRecognizeCount:     20,
			ActivePercent:          0.7,
			InactivePercent:        0.5,
			StartVoiceProb:         70,
			StartRms:               -50.0,
			StopVoiceProb:          70,
			StopRms:                -50.0,
		})

	conCfg := agoraservice.RtcConnectionConfig{
		AutoSubscribeAudio: true,
		AutoSubscribeVideo: false,
		ClientRole:         agoraservice.ClientRoleBroadcaster,
		ChannelProfile:     agoraservice.ChannelProfileLiveBroadcasting,
	}
	//conSignal := make(chan struct{})
	conHandler := &agoraservice.RtcConnectionObserver{
		OnConnected: func(con *agoraservice.RtcConnection, info *agoraservice.RtcConnectionInfo, reason int) {
			// do something
			fmt.Printf("Connected, reason %d\n", reason)
			//conSignal <- struct{}{}
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

	audioObserver := &agoraservice.AudioFrameObserver{
		OnPlaybackAudioFrameBeforeMixing: func(localUser *agoraservice.LocalUser, channelId string, userId string, frame *agoraservice.AudioFrame, vadResultState agoraservice.VadState, vadResultFraem *agoraservice.AudioFrame) bool {
			// do something

			// vad process here! and you can get the vad result, then send vadResult to ASR/STT service
			//vadResult, state := vad.Process(frame)
			// for debuging, can do vad dump but never recommended for production
			//NOTE：if enable VAD in LocalUser::RegisterAudioFrameObserver, the vad result will be returned by this callback
			// and never use frame for ARS/STT service, you better to use vadResultFraem for ASR/STT service

			vadDump.Write(frame, vadResultFraem, vadResultState)

			fmt.Printf("Playback from userId %s, far field flag %d, rms %d, pitch %d, state=%d-%d\n",
				userId, frame.FarFieldFlag, frame.Rms, frame.Pitch, int(vadResultState))

			return true
		},
	}
	con := agoraservice.NewRtcConnection(&conCfg)

	localUser := con.GetLocalUser()
	localUser.SetPlaybackAudioFrameBeforeMixingParameters(1, 16000)
	con.RegisterObserver(conHandler)
	vadConfigure := &agoraservice.AudioVadConfigV2{
		PreStartRecognizeCount: 16,
		StartRecognizeCount:    30,
		StopRecognizeCount:     20,
		ActivePercent:          0.7,
		InactivePercent:        0.5,
		StartVoiceProb:         70,
		StartRms:               -50.0,
		StopVoiceProb:          70,
		StopRms:                -50.0,
	}
	localUser.RegisterAudioFrameObserver(audioObserver, 1, vadConfigure)
	localUser.RegisterLocalUserObserver(localUserObserver)

	// sender := con.NewPcmSender()
	// defer sender.Release()
	sender := mediaNodeFactory.NewAudioPcmDataSender()

	track := agoraservice.NewCustomAudioTrackPcm(sender)

	localUser.SetAudioScenario(agoraservice.AudioScenarioChorus)
	con.Connect(token, channelName, userId)
	//<-conSignal

	track.SetEnabled(true)
	localUser.PublishAudio(track)

	track.AdjustPublishVolume(100)

	for !(*bStop) {
		time.Sleep(50 * time.Millisecond)
	}

	// release ...

	localUser.UnpublishAudio(track)
	track.SetEnabled(false)
	localUser.UnregisterAudioFrameObserver()
	localUser.UnregisterLocalUserObserver()

	con.Disconnect()

	con.UnregisterObserver()

	con.Release()
	vad.Release()
	vadDump.Close()

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
	fmt.Println("Application terminated")
}
