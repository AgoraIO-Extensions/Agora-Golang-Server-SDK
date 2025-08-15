package main

import (
	"fmt"
	//"go/printer"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"time"
	"unsafe"

	agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/rtc"
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

func DebugSteroPcmSource(filePath string) {
	// only for 16000 16bit 1channel PCM, and 160 samples per channel	
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println("Failed to open file: ", err)
		return
	}	
	defer file.Close()

	
	// open stero vad for debug
	vadConfigV1 := &agoraservice.AudioVadConfig{
		StartRecognizeCount:    30,
		StopRecognizeCount:     45,
		PreStartRecognizeCount: 16,
		ActivePercent:          0.6,
		InactivePercent:        0.2,
		RmsThr:                 -40.0,
		JointThr:               0.0,
		Aggressive:             2.0,
		VoiceProb:              0.7,
	}
	steroVadInst := agoraservice.NewSteroVad(vadConfigV1, vadConfigV1)
	defer steroVadInst.Release()

	sourcfile, err := os.OpenFile("./source_dump.pcm", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	
	defer sourcfile.Close()
	

	buffer := make([]byte, 640) // 100ms
	frame := &agoraservice.AudioFrame{
		SamplesPerSec:  16000,
		Channels:       2,
		BytesPerSample: 2,
		Buffer:         nil, // Pre-allocate frame buffer
	}

	frameCount := 0

	for {
		n, err := file.Read(buffer)

		if err != nil {
			break
		}
		if n < 640 {
			break;
		}
		frameCount++
		readed := n
		
		
		// process stero vad
		frame.Buffer = buffer
		leftFrame, leftState, rightFrame, rightState := steroVadInst.ProcessAudioFrame(frame)
		fmt.Printf("n: %d-%d, left: %d, right: %d   ", readed, len(buffer), leftState, rightState)
		dumpSteroVadResult(1, leftFrame, leftState)
		dumpSteroVadResult(0, rightFrame, rightState)
		
	}
	
	
}

func dumpSteroVadResult(isleft int, frame *agoraservice.AudioFrame, result int) {
	if (result == 1 || result == 2 || result == 3) {
		// open the file for dump
		if (isleft == 1) {
			if (LeftVadFile == nil) {
				LeftVadFile, _ = os.Create(fmt.Sprintf("./left_vad_dump_%d.pcm", leftCount))
			}
			LeftVadFile.Write(frame.Buffer)
		} else {
			if RightVadFile == nil {
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
// for file vad test: stereo vad file test
func file_vad_main() {
	DebugSteroPcmSource("/Users/weihognqin/Downloads/output_case1_1.pcm")
}

// offline vad test
func offline_vad1_test() {
	// get the file list in the current directory
	argus := os.Args
	if len(argus) < 2 {
		fmt.Println("Please input stero vad file path")
		return
	}
	filePath := argus[1]
	channels := 1
	if len(argus) > 2 {
		channels, _ = strconv.Atoi(argus[2])
	}
	// open file for read
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Failed to open file: ", err)
		return
	}
	defer file.Close()

	// set vadv1
	vadConfigV1 := &agoraservice.AudioVadConfig{
		StartRecognizeCount:    30,
		StopRecognizeCount:     45,
		PreStartRecognizeCount: 16,
		ActivePercent:          0.6,
		InactivePercent:        0.2,
		RmsThr:                 -30.0,
		JointThr:               0.0,
		Aggressive:             2.0,
		VoiceProb:              0.7,
	}
	
	var chunk int = 320

	var monoVadInst *agoraservice.AudioVad = nil
	var steroVadInst *agoraservice.SteroAudioVad = nil
	if channels == 1 {
		monoVadInst = agoraservice.NewAudioVad(vadConfigV1)
		chunk = 320
		defer monoVadInst.Release()
	} else {	
		steroVadInst = agoraservice.NewSteroVad(vadConfigV1, vadConfigV1)
		chunk = 640
		defer steroVadInst.Release()
	}
	
	
	//time.Sleep(1000 * time.Millisecond)
	
	// read the file
	buffer := make([]byte, chunk)
	frame := &agoraservice.AudioFrame{
		SamplesPerSec:  16000,
		Channels:       channels,
		BytesPerSample: 2,
		Buffer:         nil, // Pre-allocate frame buffer
	}
	frameCount := 0
	
	vadfile, _ := os.Create("./vad_dump.pcm")
	defer vadfile.Close()
	//go func ()  {
		
	var leftFrame *agoraservice.AudioFrame = nil
	var leftState int = 0
	var rightFrame *agoraservice.AudioFrame = nil
	var rightState int = 0
	
	for {
		n, err := file.Read(buffer)
		if err != nil {
			break
		}
		if n < chunk {
			break;
		}
		
		// process stero vad
		frame.Buffer = buffer[:n]
		
		if channels == 1 {
			leftFrame, leftState = monoVadInst.ProcessPcmFrame(frame)
			if leftState == 1 || leftState == 2 || leftState == 3 {
				vadfile.Write(leftFrame.Buffer)
			}
			fmt.Printf("count: %d-%d, left: %d, right: %d\n", frameCount, n	, leftState, len(leftFrame.Buffer))
			
		} else {
			leftFrame, leftState, rightFrame, rightState = steroVadInst.ProcessAudioFrame(frame)
			fmt.Printf("count: %d-%d, left: %d, right: %d\n", frameCount, n	, leftState, rightState)
			dumpSteroVadResult(1, leftFrame, leftState)
			dumpSteroVadResult(0, rightFrame, rightState)
		}
		
		

		
		
		frameCount++
	}
	//signal <- struct{}{}
	//}()
	//<-signal
}

func main() {
	//offline_vad1_test()
	original_main()
}

func original_main() {
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

	// get parameter from arguments： appid, channel_name, steroencodemoe echoback

	argus := os.Args
	if len(argus) < 3 {
		fmt.Println("Please input appid, channel name")
		return
	}
	appid := argus[1]
	channelName := argus[2]

	echoBack := 0
	steroMode := 0
	if len(argus) > 3 {
		steroMode, _ = strconv.Atoi(argus[3])
	}

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
	svcCfg.LogPath = "./agora_rtc_log/agorasdk.log"
	// setting senario for diff mode
	var steroVadInst *agoraservice.SteroAudioVad = nil
	if steroMode > 0 {
		svcCfg.AudioScenario = agoraservice.AudioScenarioGameStreaming
		//vad v1 for stero
		vadConfigV1 := &agoraservice.AudioVadConfig{
			StartRecognizeCount:    30,
			StopRecognizeCount:     48,
			PreStartRecognizeCount: 16,
			ActivePercent:          0.8,
			InactivePercent:        0.2,
			RmsThr:                 -40.0,
			JointThr:               0.0,
			Aggressive:             2.0,
			VoiceProb:              0.7,
		}
		steroVadInst = agoraservice.NewSteroVad(vadConfigV1, vadConfigV1)
	} else {
		svcCfg.AudioScenario = agoraservice.AudioScenarioChorus
	}
	svcCfg.EnableSteroEncodeMode = steroMode

	agoraservice.Initialize(svcCfg)

	

	
	// generate stero vad

	//agoraservice.EnableExtension("agora.builtin", "agora_audio_label_generator", "", true)

	// Recommended  configurations:
	// For not-so-noisy environments, use this configuration: (16, 30, 50, 0.7, 0.5, 70, -50, 70, -50)
	// For noisy environments, use this configuration: (16, 30, 50, 0.7, 0.5, 70, -40, 70, -40)
	// For high-noise environments, use this configuration: (16, 30, 50, 0.7, 0.5, 70, -30, 70, -30)
	scenario := agoraservice.AudioScenarioChorus
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

	var con *agoraservice.RtcConnection = nil
	
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
		OnAIQoSCapabilityMissing: func(con *agoraservice.RtcConnection, defaultFallbackSenario int) int {
			fmt.Printf("onAIQoSCapabilityMissing, defaultFallbackSenario: %d\n", defaultFallbackSenario)
			return int(agoraservice.AudioScenarioDefault)
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
			if steroMode > 0 {

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
			} else {
				vadDump.Write(frame, vadResultFraem, vadResultState)
				fmt.Printf("Playback from userId %s, far field flag %d, rms %d, pitch %d, state=%d\n", userId, frame.FarFieldFlag, frame.Rms, frame.Pitch, int(vadResultState))

			}

			if echoBack == 1 {
				con.PushAudioPcmData(frame.Buffer, frame.SamplesPerSec, frame.Channels, 0)
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

	// set to dump
	// change audio senario, by wei for stero encodeing
	agoraParameterHandler := agoraservice.GetAgoraParameter()

	//agoraParameterHandler.SetParameters("{\"che.audio.frame_dump\":{\"location\":\"all\",\"action\":\"start\",\"max_size_bytes\":\"100000000\",\"uuid\":\"123456789\", \"duration\": \"150000\"}}")
	// end


	
	con = agoraservice.NewRtcConnection(&conCfg, publishConfig)
	//agoraParameterHandler = con.GetAgoraParameter()
	agoraParameterHandler.SetParameters("{\"che.audio.frame_dump\":{\"location\":\"all\",\"action\":\"start\",\"max_size_bytes\":\"100000000\",\"uuid\":\"123456789\", \"duration\": \"150000\"}}")

	

	


	localUser := con.GetLocalUser()

	
	// dump audio
	// set to dump
	//agoraParameterHandler.SetParameters("{\"che.audio.frame_dump\":{\"location\":\"all\",\"action\":\"start\",\"max_size_bytes\":\"100000000\",\"uuid\":\"123456789\", \"duration\": \"150000\"}}")
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
			con.SendAudioMetaData(metaData)
		},
	}

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
	if steroMode > 0 {
		localUser.SetPlaybackAudioFrameBeforeMixingParameters(2, 16000)
		con.RegisterAudioFrameObserver(audioObserver, 0, nil)
	} else {
		localUser.SetPlaybackAudioFrameBeforeMixingParameters(1, 16000)
		con.RegisterAudioFrameObserver(audioObserver, 1, vadConfigure)
	}

	con.RegisterLocalUserObserver(localUserObserver)

	// sender := con.NewPcmSender()
	// defer sender.Release()

	//localUser.SetAudioScenario(agoraservice.AudioScenarioChorus)
	con.Connect(token, channelName, userId)
	//<-conSignal

	// for test AudioConsumer
	
	// read pcm data from file and push to audioConsumer
	// test 16khz, mono, 2 channel, but some eorror
	sourceFilePath := make([]string, 2)
	channel := make([]int, 2)
	
	// test 16khz, stero
	sourceFilePath[0] = "../test_data/1.pcm"
	channel[0] = 2
	// test 16khz, mono
	sourceFilePath[1] = "../test_data/recv_audio_16k_1ch.pcm"
	channel[1] = 1

	index := 0

	// open file
	sourceFile, err := os.Open(sourceFilePath[index])
	defer sourceFile.Close()
	if err != nil {
		cwd, _ := os.Getwd()
		fmt.Printf("open file error: %s, cwd = %s\n", err, cwd)
	}
	fileData := make([]byte, 320*100*channel[index]) // 100ms

	
	con.PublishAudio()

	

	for !(*bStop) {
		time.Sleep(50 * time.Millisecond)
		//curTime := time.Now()
		//timeStr := curTime.Format("2006-01-02 15:04:05.000")
		//localUser.SendAudioMetaData([]byte(timeStr))

		// check file ' length

		if echoBack == 0 && con.IsPushToRtcCompleted() {

			for {
				n, _ := sourceFile.Read(fileData)
				if n < 640 {
					sourceFile.Seek(0, 0)
					break
				}
				con.PushAudioPcmData(fileData[:n], 16000, channel[index], 0)
			}
		}
	}

	// release ...
	dumpFile.Close()

	

	con.Disconnect()

	

	con.Release()

	vadDump.Close()

	
	agoraservice.Release()
	if steroVadInst != nil {
		steroVadInst.Release()
	}
	steroVadInst = nil

	
	audioObserver = nil
	localUserObserver = nil
	localUser = nil
	conHandler = nil
	con = nil
	fmt.Println("Application terminated")
}
