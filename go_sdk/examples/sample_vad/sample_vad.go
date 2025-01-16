package main

import (
	"fmt"
	//"go/printer"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"
	"unsafe"

	agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/agoraservice"
	//"google.golang.org/protobuf/types/known/sourcecontextpb"

	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
)

// use unsafe to conver byte array to uint16 array
func bytesToUInt16Array(data []byte) []uint16 {
	// 通过 unsafe 包将 []byte 转换为 []int16
	//return *(*[]int16)(unsafe.Pointer(&data))
	return *(*[]uint16)(unsafe.Pointer(&data))
}

func main() {
	bStop := new(bool)
	*bStop = false

	// start pprof
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

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
	// change audio senario
	svcCfg.AudioScenario = agoraservice.AudioScenarioGameStreaming

	agoraservice.Initialize(svcCfg)

	mediaNodeFactory := agoraservice.NewMediaNodeFactory()

	sender := mediaNodeFactory.NewAudioPcmDataSender()

	track := agoraservice.NewCustomAudioTrackPcm(sender)

	//agoraservice.EnableExtension("agora.builtin", "agora_audio_label_generator", "", true)

	// Recommended  configurations:
	// For not-so-noisy environments, use this configuration: (16, 30, 50, 0.7, 0.5, 70, -50, 70, -50)
	// For noisy environments, use this configuration: (16, 30, 50, 0.7, 0.5, 70, -40, 70, -40)
	// For high-noise environments, use this configuration: (16, 30, 50, 0.7, 0.5, 70, -30, 70, -30)
	

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
	// for debuging, can do vad dump but never recommended for production
	dumpFile, err := os.OpenFile("./source_dump.pcm", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
	    fmt.Println("Failed to create dump file: ", err)
		
	}
	exceptFile, _ := os.OpenFile("./except_dump.pcm", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)

	audioObserver := &agoraservice.AudioFrameObserver{
		OnPlaybackAudioFrameBeforeMixing: func(localUser *agoraservice.LocalUser, channelId string, userId string, frame *agoraservice.AudioFrame, vadResultState agoraservice.VadState, vadResultFraem *agoraservice.AudioFrame) bool {
			// do something
			dumpFile.Write(frame.Buffer)

			// do loop test
			
			

			// a simple energy check for vad
			var leftValue uint64 = 0
			var rightValue uint64 = 0
			var curValue uint16 = 0
			uint16Array := bytesToUInt16Array(frame.Buffer)
			len := len(uint16Array) 
			for i := 0; i < len; i++ {
				// calc left & right channel energy
				curValue = uint16Array[i]


				if curValue == 0xFFFF {
				    curValue = 0
				}
				if i%2 == 0 {
				    leftValue += (uint64(curValue))
				} else {
				    rightValue += (uint64(curValue))
				}
			}

			len /= 2
			leftValue = leftValue / uint64(len)
			rightValue = rightValue / uint64(len)
			fmt.Printf("left: = %d, right: = %d\n", leftValue, rightValue)

			if leftValue > 100 {
				exceptFile.Write(frame.Buffer)
			    
			}

			sender.SendAudioPcmData(frame)


			// vad process here! and you can get the vad result, then send vadResult to ASR/STT service
			//vadResult, state := vad.Process(frame)
			// for debuging, can do vad dump but never recommended for production
			//NOTE：if enable VAD in LocalUser::RegisterAudioFrameObserver, the vad result will be returned by this callback
			// and never use frame for ARS/STT service, you better to use vadResultFraem for ASR/STT service

			//vadDump.Write(frame, vadResultFraem, vadResultState)

			//fmt.Printf("Playback from userId %s, far field flag %d, rms %d, pitch %d, state=%d\n", userId, frame.FarFieldFlag, frame.Rms, frame.Pitch, int(vadResultState))

			return true
		},
	}
	con := agoraservice.NewRtcConnection(&conCfg)

	localUser := con.GetLocalUser()

	// change audio senario, by wei for stero encodeing
	localUser.SetAudioScenario(agoraservice.AudioScenarioGameStreaming)
	localUser.SetAudioEncoderConfiguration(&agoraservice.AudioEncoderConfiguration{AudioProfile: int(agoraservice.AudioProfileMusicHighQualityStereo)})

	// fill pirvate parameter
	agoraParameterHandler := agoraservice.GetAgoraParameter()
	agoraParameterHandler.SetParameters("{\"che.audio.aec.enable\":false}")
	agoraParameterHandler.SetParameters("{\"che.audio.ans.enable\":false}")
	agoraParameterHandler.SetParameters("{\"che.audio.agc.enable\":false}")
	agoraParameterHandler.SetParameters("{\"che.audio.custom_payload_type\":78}")
	agoraParameterHandler.SetParameters("{\"che.audio.custom_bitrate\":128000}")
	
	

	localUserObserver := &agoraservice.LocalUserObserver{
		OnStreamMessage: func(localUser *agoraservice.LocalUser, uid string, streamId int, data []byte) {
			// do something
			fmt.Printf("*****Stream message, from userId %s\n", uid)
			//con.SendStreamMessage()
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
		OnAudioMetaDataReceived: func(localUser *agoraservice.LocalUser, uid string, metaData []byte) {
			fmt.Printf("*****User audio meta data received, uid %s, meta: %s\n", uid, string(metaData))
			localUser.SendAudioMetaData(metaData)
		},
	}
	localUser.SetPlaybackAudioFrameBeforeMixingParameters(2, 16000)
	con.RegisterObserver(conHandler)
	vadConfigure := &agoraservice.AudioVadConfigV2{
		PreStartRecognizeCount: 16,
		StartRecognizeCount:    30,
		StopRecognizeCount:     50,
		ActivePercent:          0.7,
		InactivePercent:        0.5,
		StartVoiceProb:         70,
		StartRms:               -50.0,
		StopVoiceProb:          70,
		StopRms:                -50.0,
	}
	localUser.RegisterAudioFrameObserver(audioObserver, 0, vadConfigure)
	localUser.RegisterLocalUserObserver(localUserObserver)

	// sender := con.NewPcmSender()
	// defer sender.Release()
	

	//localUser.SetAudioScenario(agoraservice.AudioScenarioChorus)
	con.Connect(token, channelName, userId)
	//<-conSignal

	// for test
	track.SetSendDelayMs(10)
	localUser.SetAudioScenario(agoraservice.AudioScenarioChorus)

	track.SetEnabled(true)
	localUser.PublishAudio(track)

	track.AdjustPublishVolume(100)

	for !(*bStop) {
		time.Sleep(50 * time.Millisecond)
		//curTime := time.Now()
		//timeStr := curTime.Format("2006-01-02 15:04:05.000")
		//localUser.SendAudioMetaData([]byte(timeStr))
	}

	// release ...
	dumpFile.Close()

	localUser.UnpublishAudio(track)
	track.SetEnabled(false)
	localUser.UnregisterAudioFrameObserver()
	localUser.UnregisterLocalUserObserver()

	con.Disconnect()

	con.UnregisterObserver()

	con.Release()

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
