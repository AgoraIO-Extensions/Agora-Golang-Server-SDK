package main

import (
	//"bytes"
	//"bufio"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"
	_ "net/http/pprof"

	agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/agoraservice"

	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
)

func PushFileToConsumer(file *os.File, audioConsumer *agoraservice.AudioConsumer) {
	buffer := make([]byte, 320)
	for {
		readLen, err := file.Read(buffer)
		if err != nil || readLen < 320 {
			fmt.Printf("read up to EOF,cur read: %d", readLen)
			file.Seek(0, 0)
			break
		}
		audioConsumer.PushPCMData(buffer)
	}
	buffer = nil
}
func ReadFileToConsumer(file *os.File, consumer *agoraservice.AudioConsumer, interval int, done chan bool) {
	for {
		select {
		case <-done:
			fmt.Println("ReadFileToConsumer done")
			return
		default:
			len := consumer.Len()
			if len < 320*interval {
				PushFileToConsumer(file, consumer)
			}
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	}
}


func ConsumeAudio(audioConsumer *agoraservice.AudioConsumer, interval int, done chan bool) {
	for {
		select {
		case <-done:
			fmt.Println("ConsumeAudio done")
			return
		default:
			audioConsumer.Consume()
			time.Sleep(time.Duration(interval) * time.Millisecond)
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
	// eanbe callback when muted
	svcCfg.ShouldCallbackWhenMuted = 1

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
		OnPlaybackAudioFrameBeforeMixing: func(localUser *agoraservice.LocalUser, channelId string, userId string, frame *agoraservice.AudioFrame, vadResulatState agoraservice.VadState, vadResultFrame *agoraservice.AudioFrame) bool {
			// do something
			fmt.Printf("Playback audio frame before mixing, from userId %s, far :%d,rms:%d, pitch: %d\n", userId, frame.FarFieldFlag, frame.Rms, frame.Pitch)
			return true
		},
	}
	start := time.Now().UnixMilli()

	con := agoraservice.NewRtcConnection(&conCfg)
	defer con.Release()

	localStreamId, _ := con.CreateDataStream(false, false)

	//added by wei for localuser observer
	localUserObserver := &agoraservice.LocalUserObserver{
		OnStreamMessage: func(localUser *agoraservice.LocalUser, uid string, streamId int, data []byte) {
			// do something
			fmt.Printf("*****Stream message, from userId %s\n", uid)
			//con.SendStreamMessage(streamId, data)
			con.SendStreamMessage(localStreamId, data)
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
			//ret := "abc"
			//retbytes := []byte(ret)	
			
			localUser.SendAudioMetaData(metaData)
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

	end := time.Now().UnixMilli()
	fmt.Printf("Connect cost %d ms\n", end-start)

	localUser = con.GetLocalUser()
	localUser.SetPlaybackAudioFrameBeforeMixingParameters(1, 16000)
	localUser.SetAudioVolumeIndicationParameters(300, 1, true)
	localUser.RegisterLocalUserObserver(localUserObserver)

	localUser.RegisterAudioFrameObserver(audioObserver, 0, nil)

	track.SetEnabled(true)
	localUser.PublishAudio(track)

	file, err := os.Open("../test_data/send_audio_16k_1ch.pcm")
	if err != nil {
		fmt.Printf("NewError opening file: %v\n", err)
		return
	}
	defer file.Close()

	track.AdjustPublishVolume(100)

	audioConsumer := agoraservice.NewAudioConsumer(sender, 16000, 1)
	defer audioConsumer.Release()

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
	go ReadFileToConsumer(file, audioConsumer, 50, done)
	go ConsumeAudio(audioConsumer, 50, done)


	//release operation:cancel defer release,try manual release
	for !(*bStop) {
		time.Sleep(100 * time.Millisecond)
	}
	close(done)

	audioConsumer.Release()

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
