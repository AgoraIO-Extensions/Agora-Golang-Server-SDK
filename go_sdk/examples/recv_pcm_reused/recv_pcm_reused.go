package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/rtc"

	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
)

func main() {
	stop := false
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		stop = true
	}()

	args := os.Args
	if len(args) < 3 {
		fmt.Println("usage: recv_pcm_reused <appid> <channel_name> [output.pcm]")
		return
	}
	appid := args[1]
	channelName := args[2]
	outputPath := "./recv_reused_playback.pcm"
	if len(args) > 3 {
		outputPath = args[3]
	}

	if appid == "" {
		appid = os.Getenv("AGORA_APP_ID")
	}
	cert := os.Getenv("AGORA_APP_CERTIFICATE")
	userId := "reused_pcm_recv"
	if appid == "" {
		fmt.Println("Please set AGORA_APP_ID or pass it as the first argument")
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
			fmt.Println("Failed to build token:", err)
			return
		}
	}

	outFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Failed to open output file:", err)
		return
	}
	defer outFile.Close()

	svcCfg := agoraservice.NewAgoraServiceConfig()
	svcCfg.AppId = appid
	svcCfg.LogPath = "./agora_rtc_log/agorasdk.log"
	svcCfg.ConfigDir = "./agora_rtc_log"
	svcCfg.DataDir = "./agora_rtc_log"

	agoraservice.Initialize(svcCfg)
	defer agoraservice.Release()

	conCfg := &agoraservice.RtcConnectionConfig{
		AutoSubscribeAudio: true,
		AutoSubscribeVideo: false,
		ClientRole:         agoraservice.ClientRoleAudience,
		ChannelProfile:     agoraservice.ChannelProfileLiveBroadcasting,
	}
	publishConfig := agoraservice.NewRtcConPublishConfig()
	publishConfig.AudioPublishType = agoraservice.AudioPublishTypeNoPublish
	publishConfig.AudioScenario = agoraservice.AudioScenarioDefault
	publishConfig.IsPublishAudio = false
	publishConfig.IsPublishVideo = false
	publishConfig.AudioProfile = agoraservice.AudioProfileDefault

	con := agoraservice.NewRtcConnection(conCfg, publishConfig)
	if con == nil {
		fmt.Println("Failed to create RTC connection")
		return
	}
	defer con.Release()

	conSignal := make(chan struct{})
	audioFrameCount := 0
	lastLogTime := time.Now()

	conHandler := &agoraservice.RtcConnectionObserver{
		OnConnected: func(con *agoraservice.RtcConnection, info *agoraservice.RtcConnectionInfo, reason int) {
			fmt.Printf("Connected, reason %d\n", reason)
			conSignal <- struct{}{}
		},
		OnDisconnected: func(con *agoraservice.RtcConnection, info *agoraservice.RtcConnectionInfo, reason int) {
			fmt.Printf("Disconnected, reason %d\n", reason)
		},
		OnUserJoined: func(con *agoraservice.RtcConnection, uid string) {
			fmt.Println("user joined,", uid)
		},
		OnUserLeft: func(con *agoraservice.RtcConnection, uid string, reason int) {
			fmt.Println("user left,", uid)
		},
	}

	audioObserver := &agoraservice.AudioFrameObserver{
		OnReusedPlaybackAudioFrame: func(localUser *agoraservice.LocalUser, channelId string, frame *agoraservice.AudioFrame) bool {
			if frame == nil || len(frame.Buffer) == 0 {
				return true
			}
			pcmCopy := append([]byte(nil), frame.Buffer...)
			if _, err := outFile.Write(pcmCopy); err != nil {
				fmt.Println("write pcm failed:", err)
				return false
			}
			audioFrameCount++
			if time.Since(lastLogTime) >= time.Second {
				fmt.Printf("reused playback audio ok, frames=%d bytes_per_frame=%d channel=%s sample_rate=%d\n",
					audioFrameCount, len(pcmCopy), channelId, frame.SamplesPerSec)
				audioFrameCount = 0
				lastLogTime = time.Now()
			}
			return true
		},
		OnGetPlaybackAudioFrameParam: func(localUser *agoraservice.LocalUser) agoraservice.AudioFrameObserverAudioParams {
			return agoraservice.AudioFrameObserverAudioParams{
				SampleRate:     16000,
				Channels:       1,
				Mode:           agoraservice.RawAudioFrameOpModeReadOnly,
				SamplesPerCall: 320,
			}
		},
	}

	localUserObserver := &agoraservice.LocalUserObserver{
		OnUserAudioTrackSubscribed: func(localUser *agoraservice.LocalUser, uid string, remoteAudioTrack *agoraservice.RemoteAudioTrack) {
			fmt.Printf("user %s audio subscribed\n", uid)
		},
	}

	con.RegisterObserver(conHandler)
	con.RegisterLocalUserObserver(localUserObserver)
	con.RegisterAudioFrameObserver(audioObserver, 0, nil)

	con.Connect(token, channelName, userId)
	<-conSignal

	for !stop {
		time.Sleep(100 * time.Millisecond)
	}

	con.Disconnect()
}
