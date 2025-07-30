package main

import (
	//"bytes"
	//"bufio"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"time"
	"unsafe"

	agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/rtc"

	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
)

func PushFileToConsumer(file *os.File, con *agoraservice.RtcConnection, samplerate int) {
	buffer := make([]byte, samplerate*2) // 1s data
	for {
		readLen, err := file.Read(buffer)
		if err != nil {
			fmt.Printf("read up to EOF,cur read: %d", readLen)
			file.Seek(0, 0)
			return
		}
		// round to integer of chunk
		packLen := readLen / samplerate
		con.PushAudioPcmData(buffer[:packLen*samplerate], samplerate, 1)
		fmt.Printf("PushPCMData done:%d\n", readLen)
	}
}
func ReadFileToConsumer(file *os.File, con *agoraservice.RtcConnection, interval int, done chan bool, samplerate int) {
	for {
		select {
		case <-done:
			fmt.Println("ReadFileToConsumer done")
			return	
		default:
			if con != nil {
				isPushCompleted := con.IsPushToRtcCompleted()
				if isPushCompleted {
					PushFileToConsumer(file, con, samplerate)
				}
			}
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	}
}


func LoopbackAudio(audioQueue *agoraservice.Queue, audioConsumer *agoraservice.AudioConsumer, done chan bool) {
	for {
		select {
		case <-done:
			fmt.Println("LoopbackAudio done")
			return
		default:
			AudioFrame := audioQueue.Dequeue()
			if AudioFrame != nil {
				//fmt.Printf("AudioFrame: %d\n", time.Now().UnixMilli())
				if frame, ok := AudioFrame.(*agoraservice.AudioFrame); ok {
					frame.RenderTimeMs = 0
					audioConsumer.PushPCMData(frame.Buffer)
				}
			}
			//time.Sleep(10 * time.Millisecond)
		}
	}
}

func SendTTSDataToClient(samplerate int, audioConsumer *agoraservice.AudioConsumer, file *os.File, done chan bool, audioSendEvent chan struct{}, fallackEvent chan struct{}, localUser *agoraservice.LocalUser, track *agoraservice.LocalAudioTrack) {
	for {
		select {
		case <-done:
			fmt.Println("SendAudioToClient done")
			return
		case <-fallackEvent:

		case <-audioSendEvent:
			// read 1s data from file
			buffer := make([]byte, samplerate*2*2) // 2s data
			readLen, err := file.Read(buffer)
			if err != nil {
				fmt.Printf("read up to EOF,cur read: %d", readLen)
				file.Seek(0, 0)
				continue
			}
			audioConsumer.PushPCMData(buffer[:readLen])
			// and seek to the begin of the file
			file.Seek(0, 0)
			fmt.Println("SendTTSDataToClient done")
		default:
			time.Sleep(40 * time.Millisecond)
			audioConsumer.Consume()
		}
	}
}

func calculateEnergyFast(data []byte) uint64 {
    var energy uint64
    samples := unsafe.Slice((*int16)(unsafe.Pointer(&data[0])), len(data)/2)
    for _, s := range samples {
        energy += uint64(s) * uint64(s)
    }
    return energy
}

func main() {
	bStop := new(bool)
	*bStop = false
	// start pprof
	go func() {
		runtime.SetBlockProfileRate(1)
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

	filepath := "../test_data/send_audio_16k_1ch.pcm"
	if len(argus) > 3 {
		filepath = argus[3]
	}
	//default samplerate to 16k
	samplerate := 16000
	if len(argus) > 4 {
		samplerate, _ = strconv.Atoi(argus[4]) // strconv is in the "strconv" package, which is a standard package in Go's library.
	}

	// mode: 1 loopback, 0 playout music
	mode := 0
	if len(argus) > 5 {
		mode, _ = strconv.Atoi(argus[5])
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
	svcCfg.ConfigDir = "./agora_rtc_log"
	svcCfg.DataDir = "./agora_rtc_log"

	// about AudioScenario: default is AudioScenarioAiServer
	// if want to use other scenario, pls contact us and make sure the scenario is much apdated for your business
	

	agoraservice.Initialize(svcCfg)
	

	
	

	conCfg := &agoraservice.RtcConnectionConfig{
		AutoSubscribeAudio: true,
		AutoSubscribeVideo: false,
		ClientRole:         agoraservice.ClientRoleBroadcaster,
		ChannelProfile:     agoraservice.ChannelProfileLiveBroadcasting,
	}
	conSignal := make(chan struct{})
	OnDisconnectedSign := make(chan struct{})


	//NOTE: you can set senario here, and every connection has its own senario, which can diff from the service config
	// and can diff from each other
	// but recommend to use the same senario for a connection and related audio track
	

	publishConfig := agoraservice.NewRtcConPublishConfig()
	
	publishConfig.IsPublishAudio = true
	publishConfig.IsPublishVideo = false
	publishConfig.AudioPublishType = agoraservice.AudioPublishTypePcm
	publishConfig.VideoPublishType = agoraservice.VideoPublishTypeNoPublish
	publishConfig.AudioScenario = agoraservice.AudioScenarioAiServer
	publishConfig.AudioProfile = agoraservice.AudioProfileDefault
	

	con := agoraservice.NewRtcConnection(conCfg, publishConfig)
	

	audioQueue := agoraservice.NewQueue(10)

	conHandler := &agoraservice.RtcConnectionObserver{
		OnConnected: func(con *agoraservice.RtcConnection, info *agoraservice.RtcConnectionInfo, reason int) {
			// do something
			fmt.Printf("Connected, reason %d\n", reason)
			//NOTE： Must  unpublish,and then update the track, and then publish the track!!!!

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
		OnAIQoSCapabilityMissing: func(con *agoraservice.RtcConnection, defaultFallbackSenario int) int {
			fmt.Printf("onAIQoSCapabilityMissing, defaultFallbackSenario: %d\n", defaultFallbackSenario)
			return int(agoraservice.AudioScenarioChorus)
		},
	}
	framecount := 0
	var frame_diff int64
	last_frame_time := time.Now().UnixMilli()
	start_frame_time := time.Now().UnixMilli()
	audioSendEvent := make(chan struct{})
	fallackEvent := make(chan struct{})
	interruptEvent := make(chan struct{})
	

	filename := "./ai_server_recv_pcm.pcm"
	pcm_file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer pcm_file.Close()

	is_square_wave_inpeak := false
	square_wave_count := 0
	

	audioObserver := &agoraservice.AudioFrameObserver{
		OnPlaybackAudioFrameBeforeMixing: func(localUser *agoraservice.LocalUser, channelId string, userId string, frame *agoraservice.AudioFrame, vadResulatState agoraservice.VadState, vadResultFrame *agoraservice.AudioFrame) bool {
			// do something
			fmt.Printf("Playback audio frame before mixing, from userId %s, far :%d,rms:%d, pitch: %d\n", userId, frame.FarFieldFlag, frame.Rms, frame.Pitch)
			//energy := calculateEnergyFast(frame.Buffer)
			//fmt.Printf("energy: %d, rms: %d, ravg: %f, framecount: %d\n", energy, frame.Rms,float64(energy)/float64(frame.SamplesPerChannel),framecount)
			now := time.Now().UnixMilli()
			if mode != 3 {
				if framecount == 0 {
					last_frame_time = time.Now().UnixMilli()
					start_frame_time = last_frame_time
				}
				framecount++
				frame_diff = now - last_frame_time
				last_frame_time = now
				//fmt.Printf("frame_diff: %d\n", frame_diff)
				if framecount%100 == 0 { // evry 100 frames
					frame_diff = now - start_frame_time
					fmt.Printf("******** frame :%d, duration: %d, avg frame time: %f\n", framecount, frame_diff, float64(frame_diff)/100)
					start_frame_time = now
				}
				if framecount % 4000 == 0  && mode != 3 {
					//audioSendEvent <- struct{}{}
				}

				pcm_file.Write(frame.Buffer)
			}

			if mode == 1 {
				frame.RenderTimeMs = 0
				//sender.SendAudioPcmData(frame)
				con.PushAudioPcmData(frame.Buffer, frame.SamplesPerSec, frame.Channels)
				// loopback
				//audioQueue.Enqueue(frame)

			}
			// 3: 做方波信号的echo
			threshold_value := -25
			if mode == 3 {
				if frame.Rms > threshold_value {
					// from trough to peak​​ now
					if is_square_wave_inpeak == false {
						fmt.Printf("????????#####$$$$$$$ square_wave IN peak now: %d, through count: %d\n", now, square_wave_count)
						is_square_wave_inpeak = true
						// just from trough to peak, should trigger a event
						con.InterruptAudio()
						audioSendEvent <- struct{}{}
						square_wave_count = 0
					}
					square_wave_count++
				} else {
					if is_square_wave_inpeak == true {
						fmt.Printf("????????#####$$$$$$$square_wave OUT  now: %d, peak count: %d\n", now, square_wave_count)
						is_square_wave_inpeak = false
						square_wave_count = 0
					}
					square_wave_count++
				}
			}

			return true
		},
	}
	start := time.Now().UnixMilli()

	
	agoraservice.GetAgoraParameter()
	

	event_count := 0

	//added by wei for localuser observer
	localUserObserver := &agoraservice.LocalUserObserver{
		OnStreamMessage: func(localUser *agoraservice.LocalUser, uid string, streamId int, data []byte) {
			// do something
			fmt.Printf("*****Stream message, from userId %s\n", uid)
			//con.SendStreamMessage(streamId, data)
			con.SendStreamMessage(data)
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
			fmt.Printf("*****User audio meta data received, uid %s, meta: %s, event_count: %d\n", uid, string(metaData), event_count)

			event_count++
			/*
			if event_count == 10 {
				fallackEvent <- struct{}{}
			} else */if event_count%2 == 0 {

				interruptEvent <- struct{}{}
			} else { // simulate to interrupt audio
				audioSendEvent <- struct{}{}
			}

			con.SendAudioMetaData(metaData)
		},
		OnAudioTrackPublishSuccess: func(localUser *agoraservice.LocalUser, audioTrack *agoraservice.LocalAudioTrack) {
			fmt.Printf("*****Audio track publish success, time %d\n", time.Now().UnixMilli())
		},
		OnAudioTrackUnpublished: func(localUser *agoraservice.LocalUser, audioTrack *agoraservice.LocalAudioTrack) {
			fmt.Printf("*****Audio track unpublished, time %d\n", time.Now().UnixMilli())
		},
	}

	con.RegisterObserver(conHandler)

	//end

	// sender := con.NewPcmSender()
	// defer sender.Release()

	localUser := con.GetLocalUser()

	localUser.SetPlaybackAudioFrameBeforeMixingParameters(1, 16000)
	con.RegisterLocalUserObserver(localUserObserver)

	vadConfig := &agoraservice.AudioVadConfigV2{
		PreStartRecognizeCount: 16,
		StartRecognizeCount: 30,
		StopRecognizeCount: 50,
		ActivePercent: 0.7,
		InactivePercent: 0.5,
	}

	con.RegisterAudioFrameObserver(audioObserver, 1, vadConfig)


	con.Connect(token, channelName, userId)
	<-conSignal

	end := time.Now().UnixMilli()
	fmt.Printf("Connect cost %d ms\n", end-start)

	start_publish := time.Now().UnixMilli()
	con.PublishAudio()
	end_publish := time.Now().UnixMilli()
	fmt.Printf("Publish audio cost %d ms\n", end_publish-start_publish)
	//time.Sleep(1000 * time.Millisecond)
	//con.UnpublishAudio()
	start_publish = time.Now().UnixMilli()
	fmt.Printf("Unpublish audio cost %d ms\n", start_publish-end_publish)

	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("NewError opening file: %v\n", err)
		return
	}
	defer file.Close()

	

	done := make(chan bool)
	// new method for push
	/*

			#下面的操作：只是模拟生产的数据。
			# - 在sample中，为了确保生产产生的数据能够一直播放，需要生产足够多的数据，所以用这样的方式来模拟
			# - 在实际使用中，数据是实时产生的，所以不需要这样的操作。只需要在TTS生产数据的时候，调用AudioConsumer.push_pcm_data()
			 # 我们启动2个task
		    # 一个task，用来模拟从TTS接收到语音，然后将语音push到audio_consumer
		    # 另一个task，用来模拟播放语音：从audio_consumer中取出语音播放
		    # 在实际应用中，可以是TTS返回的时候，直接将语音push到audio_consumer
		    # 然后在另外一个“timer”的触发函数中，调用audio_consumer.consume()。
		    # 推荐：
		    # .Timer的模式；也可以和业务已有的timer结合在一起使用，都可以。只需要在timer 触发的函数中，调用audio_consumer.consume()即可
		    # “Timer”的触发间隔，可以和业务已有的timer间隔一致，也可以根据业务需求调整，推荐在40～80ms之间

	*/
	if mode == 0 {
		go ReadFileToConsumer(file, con, 50, done, samplerate)
	} else if mode == 2 || mode == 3 {
		//go SendTTSDataToClient(samplerate, audioConsumer, file, done, audioSendEvent, fallackEvent, localUser, track)
		go func() {
			//consturt a audio frame
			leninsecond := 20
			if mode == 3 {
				leninsecond = 5
			}
			buffer := make([]byte, samplerate*2*leninsecond)     // max to 20s data
			bytesPerFrame := (samplerate / 100) * 2 * 1 // 10ms , mono
			frame := &agoraservice.AudioFrame{
				Buffer:            nil,
				RenderTimeMs:      0,
				SamplesPerChannel: samplerate / 100,
				BytesPerSample:    2,
				Channels:          1,
				SamplesPerSec:     samplerate,
				Type:              agoraservice.AudioFrameTypePCM16,
			}
			for {
				select {
				case <-done:
					fmt.Println("SendAudioToClient done")
					return
				case <-fallackEvent:
					fmt.Println("?????? fallackEvent")

					con.UpdateAudioSenario(agoraservice.AudioScenarioDefault)

				case <-interruptEvent:
					fmt.Println("?????? interruptEvent")
					con.InterruptAudio()

				case <-audioSendEvent:
					// read 1s data from file
					//con.PublishAudio()

					readLen, err := file.Read(buffer)
					if err != nil {
						fmt.Printf("read up to EOF,cur read: %d", readLen)
						file.Seek(0, 0)
						continue
					}
					frame.Buffer = buffer[:readLen]
					packnum := readLen / bytesPerFrame
					frame.SamplesPerChannel = (samplerate / 100) * packnum
					ret := con.PushAudioPcmData(buffer[:readLen], samplerate, 1)
					

					//audioConsumer.PushPCMData(frame.Buffer)
					// and seek to the begin of the file
					file.Seek(0, 0)
					fmt.Printf("SendTTSDataToClient done, ret: %d\n", ret)
				default:
					time.Sleep(40 * time.Millisecond)
				}
			}
		}()
	} else {
		//go LoopbackAudio(audioQueue, audioConsumer, done)
		fmt.Printf("len: %d\n", audioQueue.Size())
	}

	//release operation:cancel defer release,try manual release
	for !(*bStop) {
		time.Sleep(100 * time.Millisecond)
	}
	close(done)

	

	start_disconnect := time.Now().UnixMilli()
	con.Disconnect()
	//<-OnDisconnectedSign
	

	con.Release()

	
	agoraservice.Release()

	
	audioObserver = nil
	localUserObserver = nil
	localUser = nil
	conHandler = nil
	con = nil

	fmt.Printf("Disconnected, cost %d ms\n", time.Now().UnixMilli()-start_disconnect)
}
