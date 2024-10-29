package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"time"

	agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/agoraservice"

	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
)

type AudioLabel struct {
	BufferSize        int `json:"buffer_size"`
	SamplesPerChannel int `json:"samples_per_channel"`
	BytesPerSample    int `json:"bytes_per_sample"`
	Channels          int `json:"channels"`
	SampleRate        int `json:"sample_rate"`
	FarFieldFlag      int `json:"far_field_flag"`
	VoiceProb         int `json:"voice_prob"`
	Rms               int `json:"rms"`
	Pitch             int `json:"pitch"`
}

func AudioFrameToString(frame *agoraservice.AudioFrame) string {
	al := AudioLabel{
		BufferSize:        len(frame.Buffer),
		SamplesPerChannel: frame.SamplesPerChannel,
		BytesPerSample:    frame.BytesPerSample,
		Channels:          frame.Channels,
		SampleRate:        frame.SamplesPerSec,
		FarFieldFlag:      frame.FarFieldFlag,
		VoiceProb:         frame.VoiceProb,
		Rms:               frame.Rms,
		Pitch:             frame.Pitch,
	}
	alStr, _ := json.Marshal(al)
	return string(alStr)
}

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

	// get environment variable
	appid := os.Getenv("AGORA_APP_ID")
	cert := os.Getenv("AGORA_APP_CERTIFICATE")
	channelName := "gosdktest"
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

	agoraservice.EnableExtension("agora.builtin", "agora_audio_label_generator", "", true)
	agoraservice.GetAgoraParameter().SetParameters("{\"che.audio.label.enable\": true}")

	vad := agoraservice.NewAudioVadV2(
		&agoraservice.AudioVadConfigV2{
			StartVoiceProb:         70,
			StartRms:               -50.0,
			StopVoiceProb:          70,
			StopRms:                -50.0,
			StartRecognizeCount:    30,
			StopRecognizeCount:     20,
			PreStartRecognizeCount: 16,
			ActivePercent:          0.7,
			InactivePercent:        0.5,
		})
	defer vad.Release()

	conCfg := agoraservice.RtcConnectionConfig{
		AutoSubscribeAudio: true,
		AutoSubscribeVideo: false,
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
	var preVadDump *os.File = nil
	var vadDump *os.File = nil
	defer func() {
		if preVadDump != nil {
			preVadDump.Close()
		}
		if vadDump != nil {
			vadDump.Close()
		}
	}()
	var vadCount *int = new(int)
	*vadCount = 0
	audioObserver := &agoraservice.AudioFrameObserver{
		OnPlaybackAudioFrameBeforeMixing: func(localUser *agoraservice.LocalUser, channelId string, userId string, frame *agoraservice.AudioFrame) bool {
			// do something
			// fmt.Printf("Playback audio frame before mixing, from userId %s\n", userId)
			if preVadDump == nil {
				var err error
				preVadDump, err = os.OpenFile(fmt.Sprintf("./pre_vad_%s_%v.pcm", userId, time.Now().Format("2006-01-02-15-04-05")), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
				if err != nil {
					fmt.Println("Failed to create dump file: ", err)
				}
			}
			if preVadDump != nil {
				fmt.Printf("PreVad: %s\n", AudioFrameToString(frame))
				preVadDump.Write(frame.Buffer)
			}
			// vad
			vadResult, state := vad.ProcessAudioFrame(frame)
			duration := 0
			if vadResult != nil {
				duration = vadResult.SamplesPerChannel / 16
			}
			if state == agoraservice.VadStateIsSpeeking || state == agoraservice.VadStateStartSpeeking {
				fmt.Printf("Vad result: state: %v, duration: %v\n", state, duration)
				if vadDump == nil {
					*vadCount++
					var err error
					vadDump, err = os.OpenFile(fmt.Sprintf("./vad_%d_%s_%v.pcm", *vadCount, userId, time.Now().Format("2006-01-02-15-04-05")), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
					if err != nil {
						fmt.Println("Failed to create dump file: ", err)
					}
				}
				if vadDump != nil {
					vadDump.Write(vadResult.Buffer)
				}
			} else {
				if vadDump != nil {
					vadDump.Close()
					vadDump = nil
				}
			}
			return true
		},
	}
	con := agoraservice.NewRtcConnection(&conCfg)
	defer con.Release()

	localUser := con.GetLocalUser()
	localUser.SetPlaybackAudioFrameBeforeMixingParameters(1, 16000)
	con.RegisterObserver(conHandler)
	localUser.RegisterAudioFrameObserver(audioObserver)

	// sender := con.NewPcmSender()
	// defer sender.Release()
	sender := mediaNodeFactory.NewAudioPcmDataSender()
	defer sender.Release()
	track := agoraservice.NewCustomAudioTrackPcm(sender)
	defer track.Release()

	localUser.SetAudioScenario(agoraservice.AudioScenarioChorus)
	con.Connect(token, channelName, userId)
	<-conSignal

	track.SetEnabled(true)
	localUser.PublishAudio(track)

	track.AdjustPublishVolume(100)

	for !(*bStop) {
		time.Sleep(1 * time.Second)
	}
	localUser.UnpublishAudio(track)
	track.SetEnabled(false)
	con.Disconnect()
}
