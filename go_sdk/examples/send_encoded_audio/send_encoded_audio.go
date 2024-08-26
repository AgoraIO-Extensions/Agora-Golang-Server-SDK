package main

// #cgo CFLAGS: -I/opt/homebrew/Cellar/ffmpeg/7.0.1/include
// #cgo LDFLAGS: -L/opt/homebrew/Cellar/ffmpeg/7.0.1/lib -lavformat -lavcodec -lavutil
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

	"agora.io/agoraservice"
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

func closeMediaFile(pFormatContext *C.struct_AVFormatContext) {
	C.avformat_close_input(&pFormatContext)
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
		os.Exit(0)
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
	svcCfg := agoraservice.AgoraServiceConfig{
		EnableAudioProcessor: true,
		EnableAudioDevice:    false,
		EnableVideo:          false,

		AppId:          appid,
		ChannelProfile: agoraservice.ChannelProfileLiveBroadcasting,
		AudioScenario:  agoraservice.AudioScenarioChorus,
		UseStringUid:   false,
		LogPath:        "./agora_rtc_log/agorasdk.log",
		LogSize:        512 * 1024,
	}
	agoraservice.Initialize(&svcCfg)
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
	conHandler := &agoraservice.RtcConnectionObserver{
		OnConnected: func(con *agoraservice.RtcConnection, info *agoraservice.RtcConnectionInfo, reason int) {
			// do something
			fmt.Println("Connected")
			conSignal <- struct{}{}
		},
		OnDisconnected: func(con *agoraservice.RtcConnection, info *agoraservice.RtcConnectionInfo, reason int) {
			// do something
			fmt.Println("Disconnected")
		},
		OnUserJoined: func(con *agoraservice.RtcConnection, uid string) {
			fmt.Println("user joined, " + uid)
		},
		OnUserLeft: func(con *agoraservice.RtcConnection, uid string, reason int) {
			fmt.Println("user left, " + uid)
		},
	}
	audioObserver := &agoraservice.AudioFrameObserver{
		OnPlaybackAudioFrameBeforeMixing: func(localUser *agoraservice.LocalUser, channelId string, userId string, frame *agoraservice.PcmAudioFrame) {
			// do something
			fmt.Printf("Playback audio frame before mixing, from userId %s\n", userId)
		},
	}
	con := agoraservice.NewConnection(&conCfg)
	defer con.Release()

	localUser := con.GetLocalUser()
	localUser.SetPlaybackAudioFrameBeforeMixingParameters(1, 16000)
	con.RegisterObserver(conHandler)
	localUser.RegisterAudioFrameObserver(audioObserver)

	// sender := con.NewPcmSender()
	// defer sender.Release()
	sender := mediaNodeFactory.NewAudioEncodedFrameSender() // .NewAudioPcmDataSender()
	defer sender.Release()
	track := agoraservice.NewCustomEncodedAudioTrack(sender, agoraservice.AudioTrackMixDisabled) // .NewCustomPcmAudioTrack(sender)
	defer track.Release()

	con.Connect(token, channelName, userId)
	<-conSignal

	track.SetEnabled(true)
	localUser.PublishAudio(track)

	// frame := agoraservice.PcmAudioFrame{
	// 	Data:              make([]byte, 320),
	// 	Timestamp:         0,
	// 	SamplesPerChannel: 160,
	// 	BytesPerSample:    2,
	// 	NumberOfChannels:  1,
	// 	SampleRate:        16000,
	// }

	// file, err := os.Open("../../../test_data/test.aac")
	// if err != nil {
	// 	fmt.Println("Error opening file:", err)
	// 	return
	// }
	// defer file.Close()

	// Open video file
	pFormatContext := openMediaFile("../../../test_data/test.aac")
	packet := C.av_packet_alloc()
	streamInfo := getStreamInfo(pFormatContext)
	codecParam := (*C.struct_AVCodecParameters)(unsafe.Pointer(streamInfo.codecpar))
	tb := streamInfo.time_base

	track.AdjustPublishVolume(100)

	sendAudioDuration := 0
	firstSendTime := time.Now()
	for !(*bStop) {
		shouldSendMs := int(time.Since(firstSendTime).Milliseconds()) - sendAudioDuration
		for i := 0; i < shouldSendMs; {
			ret := int(C.av_read_frame(pFormatContext, packet))
			if ret < 0 {
				fmt.Println("Finished reading file:", ret)
				// file.Seek(0, 0)
				closeMediaFile(pFormatContext)
				pFormatContext = openMediaFile("../../../test_data/test.aac")
				streamInfo = getStreamInfo(pFormatContext)
				codecParam = (*C.struct_AVCodecParameters)(unsafe.Pointer(streamInfo.codecpar))
				continue
			}
			fmt.Printf("Read frame duration %d, tb %d/%d, samples %d, channels %d\n",
				packet.duration, tb.num, tb.den,
				codecParam.frame_size, codecParam.ch_layout.nb_channels)
			duration := int(packet.duration) * int(tb.num) * 1000 / int(tb.den)
			sendAudioDuration += duration
			i += duration
			data := C.GoBytes(unsafe.Pointer(packet.data), packet.size)
			ret = sender.SendEncodedAudioFrame(data, &agoraservice.EncodedAudioFrameInfo{
				Speech:            true,
				Codec:             agoraservice.AudioCodecAacLc,
				SampleRateHz:      int(codecParam.sample_rate),
				SamplesPerChannel: int(codecParam.frame_size / codecParam.ch_layout.nb_channels),
				SendEvenIfEmpty:   true,
				NumberOfChannels:  int(codecParam.ch_layout.nb_channels),
			})
			fmt.Printf("SendEncodedAudioFrame %d ret: %d\n", sendAudioDuration, ret)
			C.av_packet_unref(packet)
		}
		fmt.Printf("Sent %d frames this time\n", shouldSendMs)
		time.Sleep(50 * time.Millisecond)
	}
	closeMediaFile(pFormatContext)
	localUser.UnpublishAudio(track)
	track.SetEnabled(false)
	con.Disconnect()
}
