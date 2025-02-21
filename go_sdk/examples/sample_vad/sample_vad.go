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

// dump vad result for debug
var LeftVadFile *os.File = nil
var RightVadFile *os.File = nil
var leftCount int = 0
var rightCount int = 0
func dumpSteroVadResult(isleft int, frame *agoraservice.AudioFrame, result int) {
	if (result == 1 || result == 2 || result == 3) {
		// open the file for dump
		if (isleft == 1) {
			if (LeftVadFile == nil) {
				LeftVadFile, _ = os.Create(fmt.Sprintf("./left_vad_dump_%d.pcm", leftCount))
			}
			LeftVadFile.Write(frame.Buffer)
		} else {
			if (RightVadFile == nil) {
				RightVadFile, _ = os.Create(fmt.Sprintf("./right_vad_dump_%d.pcm", rightCount))
			}
			RightVadFile.Write(frame.Buffer)
		}

		if result == 3 {
		    if (isleft == 1) {
				leftCount++
				LeftVadFile.Close()
				LeftVadFile = nil
		    } else {
				rightCount++
				RightVadFile.Close()
				RightVadFile = nil
		    }
		}
	}
	
    
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

	echoBack := 0

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
	
	//vad v1 for stero 
	vadConfigV1 := &agoraservice.AudioVadConfig{
		StartRecognizeCount:    10,
		StopRecognizeCount:     6,
		PreStartRecognizeCount: 10,
		ActivePercent:          0.6,
		InactivePercent:        0.2,
	}
	// generate stero vad
	steroVadInst := agoraservice.NewSteroVad(vadConfigV1, vadConfigV1)

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
	//exceptFile, _ := os.OpenFile("./except_dump.pcm", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)

	audioObserver := &agoraservice.AudioFrameObserver{
		OnPlaybackAudioFrameBeforeMixing: func(localUser *agoraservice.LocalUser, channelId string, userId string, frame *agoraservice.AudioFrame, vadResultState agoraservice.VadState, vadResultFraem *agoraservice.AudioFrame) bool {
			// do something
			dumpFile.Write(frame.Buffer)
			
			

			// do stero vad process
			/*
			start := time.Now().Local().UnixMilli()
			leftFrame, leftState, rightFrame, rightState := steroVadInst.ProcessAudioFrame(frame)
			end := time.Now().UnixMilli()
			leftLen := 0
			if leftFrame != nil {
				leftLen = len(leftFrame.Buffer)
			}
			rightLen := 0
			if rightFrame != nil {
				rightLen = len(rightFrame.Buffer)
			}
			fmt.Printf("left vad state %d, left len %d, right vad state %d, right len: %d,diff = %d\n", leftState, leftLen, rightState, rightLen, end-start)
			
			// dump vad frame for debug
			dumpSteroVadResult(1, leftFrame, leftState)
			dumpSteroVadResult(0, rightFrame, rightState)
			*/
		
			if echoBack == 1 {
				sender.SendAudioPcmData(frame)
			}
	


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
	agoraParameterHandler := agoraservice.GetAgoraParameter()
	builtInFunc := 1
	if builtInFunc == 0 {
		localUser.SetAudioScenario(agoraservice.AudioScenarioGameStreaming)
		localUser.SetAudioEncoderConfiguration(&agoraservice.AudioEncoderConfiguration{AudioProfile: int(agoraservice.AudioProfileMusicHighQualityStereo)})

		// fill pirvate parameter
		
		agoraParameterHandler.SetParameters("{\"che.audio.aec.enable\":false}")
		agoraParameterHandler.SetParameters("{\"che.audio.ans.enable\":false}")
		agoraParameterHandler.SetParameters("{\"che.audio.agc.enable\":false}")
		agoraParameterHandler.SetParameters("{\"che.audio.custom_payload_type\":78}")
		agoraParameterHandler.SetParameters("{\"che.audio.custom_bitrate\":128000}")
	} else {
	    con.EnableSteroEncodeMode()
	}

	// dump audio
	// set to dump 
	agoraParameterHandler.SetParameters("{\"che.audio.frame_dump\":{\"location\":\"all\",\"action\":\"start\",\"max_size_bytes\":\"100000000\",\"uuid\":\"123456789\", \"duration\": \"150000\"}}")
	// end 
	

	
	

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

	// for test AudioConsumer
	audioConsumer := agoraservice.NewAudioConsumer(sender, 16000, 2)
	defer audioConsumer.Release()
	// read pcm data from file and push to audioConsumer
	sourceFilePath := "../test_data/1.pcm"
	sourceFile, err := os.Open(sourceFilePath)
	defer sourceFile.Close()
	if err != nil {
		cwd, _ := os.Getwd()
		fmt.Printf("open file error: %s, cwd = %s\n", err, cwd)
	} 
	fileData := make([]byte, 640*100)  // 100ms 
	

	

	track.SetEnabled(true)
	localUser.PublishAudio(track)

	track.AdjustPublishVolume(100)

	for !(*bStop) {
		time.Sleep(50 * time.Millisecond)
		//curTime := time.Now()
		//timeStr := curTime.Format("2006-01-02 15:04:05.000")
		//localUser.SendAudioMetaData([]byte(timeStr))

		// check file ' length
		
		if echoBack == 0 && audioConsumer.Len() < 640*100 {
			
			for {
				n, _ := sourceFile.Read(fileData)
				if n < 640 {
					sourceFile.Seek(0, 0)
					break
				}
				audioConsumer.PushPCMData(fileData[:n])
			}
		}
			
		audioConsumer.Consume()
		
		/*
		n, _ := sourceFile.Read(fileData)
		if n < 640 {
			sourceFile.Seek(0, 0)
			continue
		}
		frame := &agoraservice.AudioFrame{}
		frame.Channels = 2
		frame.SamplesPerSec = 16000
		frame.SamplesPerChannel = 160*10
		frame.BytesPerSample = 2
		frame.Buffer = fileData[:n]

		sender.SendAudioPcmData(frame)
		*/
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

	steroVadInst.Release()

	track = nil
	audioObserver = nil
	localUserObserver = nil
	localUser = nil
	conHandler = nil
	con = nil
	fmt.Println("Application terminated")
}
