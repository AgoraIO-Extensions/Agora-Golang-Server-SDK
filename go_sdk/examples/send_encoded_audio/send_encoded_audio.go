package main

// #cgo pkg-config: libavformat libavcodec libavutil
// #include <libavformat/avformat.h>
// #include <libavutil/avutil.h>
// #include <libavcodec/avcodec.h>
import "C"
import (
	"fmt"
	"os"
	"os/signal"
	"time"
	"unsafe"

	agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/rtc"
	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
)

func openMediaFile(file string) *C.struct_AVFormatContext {
	var pFormatContext *C.struct_AVFormatContext = nil
	fn := C.CString(file)
	defer C.free(unsafe.Pointer(fn))
	if C.avformat_open_input(&pFormatContext, fn, nil, nil) != 0 {
		fmt.Printf("Unable to open file\n")
		return nil
	}

	// Retrieve stream information
	if C.avformat_find_stream_info(pFormatContext, nil) < 0 {
		fmt.Println("Couldn't find stream information")
		return nil
	}

	return pFormatContext
}

func getStreamInfo(pFormatContext *C.struct_AVFormatContext) *C.struct_AVStream {
	streams := unsafe.Slice((**C.struct_AVStream)(unsafe.Pointer(pFormatContext.streams)), pFormatContext.nb_streams)
	return streams[0]
}

func closeMediaFile(pFormatContext **C.struct_AVFormatContext) {
	C.avformat_close_input(pFormatContext)
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
	argus := os.Args
	if len(argus) < 3 {
		fmt.Println("Please input appid, channel name")
		return
	}
	appid := argus[1]
	channelName := argus[2]

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
	
	var con *agoraservice.RtcConnection = nil

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
		OnAIQoSCapabilityMissing: func(con *agoraservice.RtcConnection, defaultFallbackSenario int) int {
			fmt.Printf("onAIQoSCapabilityMissing, defaultFallbackSenario: %d\n", defaultFallbackSenario)
			return int(agoraservice.AudioScenarioDefault)
		},
	}
	audioObserver := &agoraservice.AudioFrameObserver{
		OnPlaybackAudioFrameBeforeMixing: func(localUser *agoraservice.LocalUser, channelId string, userId string, frame *agoraservice.AudioFrame, vadResultStat agoraservice.VadState, vadResultFrame *agoraservice.AudioFrame) bool {
			// do something
			fmt.Printf("Playback audio frame before mixing, from userId %s\n", userId)
			return true
		},
	}
	scenario := agoraservice.AudioScenarioAiServer
	conCfg := &agoraservice.RtcConnectionConfig{
		AutoSubscribeAudio: true,
		AutoSubscribeVideo: false,
		ClientRole:         agoraservice.ClientRoleBroadcaster,
		ChannelProfile:     agoraservice.ChannelProfileLiveBroadcasting,
	}
	publishConfig := agoraservice.NewRtcConPublishConfig()
	publishConfig.AudioScenario = scenario
	publishConfig.IsPublishAudio = true
	publishConfig.IsPublishVideo = false
	publishConfig.AudioPublishType = agoraservice.AudioPublishTypeEncodedPcm
	publishConfig.VideoPublishType = agoraservice.VideoPublishTypeNoPublish

	con = agoraservice.NewRtcConnection(conCfg, publishConfig)
	

	localUser := con.GetLocalUser()
	localUser.SetPlaybackAudioFrameBeforeMixingParameters(1, 16000)
	con.RegisterObserver(conHandler)
	con.RegisterAudioFrameObserver(audioObserver, 0, nil)

	
	con.Connect(token, channelName, userId)
	<-conSignal

	con.PublishAudio()

	pFormatContext := openMediaFile("../test_data/send_audio_16k.aac")
	if pFormatContext == nil {
		return
	}
	defer closeMediaFile(&pFormatContext)

	packet := C.av_packet_alloc()
	defer C.av_packet_free(&packet)
	streamInfo := getStreamInfo(pFormatContext)
	codecParam := (*C.struct_AVCodecParameters)(unsafe.Pointer(streamInfo.codecpar))
	tb := streamInfo.time_base


	sendAudioDuration := 0
	// firstSendTime := time.Now()
	for !(*bStop) {
		// shouldSendMs := int(time.Since(firstSendTime).Milliseconds()) - sendAudioDuration
		// for i := 0; i < shouldSendMs; {
		ret := int(C.av_read_frame(pFormatContext, packet))
		if ret < 0 {
			fmt.Println("Finished reading file:", ret)
			// file.Seek(0, 0)
			closeMediaFile(&pFormatContext)
			pFormatContext = openMediaFile("../test_data/send_audio_16k.aac")
			streamInfo = getStreamInfo(pFormatContext)
			codecParam = (*C.struct_AVCodecParameters)(unsafe.Pointer(streamInfo.codecpar))
			continue
		}
		fmt.Printf("Read frame duration %d, tb %d/%d, samples %d, channels %d\n",
			packet.duration, tb.num, tb.den,
			codecParam.frame_size, codecParam.ch_layout.nb_channels)
		duration := int(packet.duration) * int(tb.num) * 1000 / int(tb.den)
		sendAudioDuration += duration
		// i += duration
		data := C.GoBytes(unsafe.Pointer(packet.data), packet.size)
		if data[0] != 0xFF || (data[1] != 0xF1 && data[1] != 0xF9) {
			fmt.Printf("Invalid aac frame\n")
		}
		ret = con.PushAudioEncodedData(data, &agoraservice.EncodedAudioFrameInfo{
			Speech:            false,
			Codec:             agoraservice.AudioCodecAacLc,
			SampleRateHz:      int(codecParam.sample_rate),
			SamplesPerChannel: int(codecParam.frame_size / codecParam.ch_layout.nb_channels),
			SendEvenIfEmpty:   true,
			NumberOfChannels:  int(codecParam.ch_layout.nb_channels),
		})
		fmt.Printf("SendEncodedAudioFrame %d ret: %d\n", sendAudioDuration, ret)
		C.av_packet_unref(packet)
		// }
		// fmt.Printf("Sent %d frames this time\n", shouldSendMs)
		//TODO: sleep time should be calculated based on the audio frame duration
		time.Sleep(21 * time.Millisecond)
	}
	
	con.Disconnect()

	con.Release()

	agoraservice.Release()
	fmt.Println("Application terminated")
}
