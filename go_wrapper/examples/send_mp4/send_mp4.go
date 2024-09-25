package main

// #cgo pkg-config: libavformat libavcodec libavutil libswresample libswscale
// #include <string.h>
// #include <stdlib.h>
// #include <libavutil/error.h>
// #include <libavutil/pixfmt.h>
// #include <libavutil/samplefmt.h>
// #include <libavutil/avutil.h>
// #include "decode_media.h"
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
		AppId:         appid,
		AudioScenario: agoraservice.AUDIO_SCENARIO_CHORUS,
		LogPath:       "./agora_rtc_log/agorasdk.log",
		LogSize:       512 * 1024,
	}
	agoraservice.Init(&svcCfg)
	defer agoraservice.Destroy()

	conCfg := agoraservice.RtcConnectionConfig{
		SubAudio:       false,
		SubVideo:       false,
		ClientRole:     1,
		ChannelProfile: 1,

		SubAudioConfig: &agoraservice.SubscribeAudioConfig{
			SampleRate: 16000,
			Channels:   1,
		},
	}
	conSignal := make(chan struct{})
	conHandler := agoraservice.RtcConnectionEventHandler{
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
	conCfg.ConnectionHandler = &conHandler
	// conCfg.AudioFrameObserver = &agoraservice.RtcConnectionAudioFrameObserver{
	// 	OnPlaybackAudioFrameBeforeMixing: func(con *agoraservice.RtcConnection, channelId string, userId string, frame *agoraservice.PcmAudioFrame) {
	// 		// do something
	// 		fmt.Printf("Playback audio frame before mixing, from userId %s\n", userId)
	// 	},
	// }
	con := agoraservice.NewConnection(&conCfg)
	defer con.Release()

	audioSender := con.NewPcmSender()
	defer audioSender.Release()

	videoSender := con.GetVideoSender()
	videoSender.SetVideoEncoderConfig(&agoraservice.VideoEncoderConfig{
		CodecType:         2,
		Width:             360,
		Height:            640,
		Framerate:         30,
		Bitrate:           800,
		MinBitrate:        400,
		OrientationMode:   0,
		DegradePreference: 0,
	})

	con.Connect(token, channelName, userId)
	defer con.Disconnect()
	<-conSignal
	audioSender.Start()
	videoSender.Start()

	fn := C.CString("../../../test_data/demo-1.mp4")
	defer C.free(unsafe.Pointer(fn))
	decoder := C.open_media_file(fn)
	if decoder == nil {
		fmt.Println("Error opening media file")
		return
	}
	defer C.close_media_file(decoder)

	cFrame := C.struct__MediaFrame{}
	C.memset(unsafe.Pointer(&cFrame), 0, C.sizeof_struct__MediaFrame)

	firstPts := int64(0)
	firstSendTime := time.Now()
	for !(*bStop) {
		totalSendTime := time.Since(firstSendTime).Milliseconds()
		ret := C.get_frame(decoder, &cFrame)
		if ret == C.AVERROR_EOF {
			fmt.Println("Finished reading file:", ret)
			break
		}
		if ret != 0 {
			fmt.Println("Error reading frame:", ret)
			continue
		}
		//NOTICE: time stamp must be greater than 0
		// if time stamp is 0, system time will be used in sdk,
		// which will cause frame whose time is less than system time, to be dropped
		if cFrame.pts <= 0 {
			cFrame.pts = 1
		}
		if firstPts == 0 {
			firstPts = int64(cFrame.pts)
			firstSendTime = time.Now()
			totalSendTime = 0
			time.Sleep(50 * time.Millisecond)
			fmt.Println("First pts:", firstPts)
		}
		if int64(cFrame.pts)-firstPts > totalSendTime {
			time.Sleep(50 * time.Millisecond)
			fmt.Println("Sleeping for 50ms")
		}
		if cFrame.frame_type == C.AVMEDIA_TYPE_AUDIO {
			if cFrame.format != C.AV_SAMPLE_FMT_S16 {
				fmt.Println("Unsupported audio format")
				continue
			}
			audioFrame := agoraservice.PcmAudioFrame{
				Data:              unsafe.Slice((*byte)(unsafe.Pointer(cFrame.buffer)), cFrame.buffer_size),
				Timestamp:         int64(cFrame.pts),
				SamplesPerChannel: int(cFrame.samples),
				BytesPerSample:    int(cFrame.bytes_per_sample),
				NumberOfChannels:  int(cFrame.channels),
				SampleRate:        int(cFrame.sample_rate),
			}
			ret := audioSender.SendPcmData(&audioFrame)
			fmt.Printf("SendPcmData %d ret: %d\n", cFrame.pts, ret)
		}
		if cFrame.frame_type == C.AVMEDIA_TYPE_VIDEO {
			if cFrame.format != C.AV_PIX_FMT_YUV420P {
				fmt.Println("Unsupported video format")
				continue
			}
			videoFrame := agoraservice.VideoFrame{
				Buffer:    unsafe.Slice((*byte)(unsafe.Pointer(cFrame.buffer)), cFrame.buffer_size),
				Width:     int(cFrame.width),
				Height:    int(cFrame.height),
				YStride:   int(cFrame.width),
				UStride:   int(cFrame.width / 2),
				VStride:   int(cFrame.width / 2),
				Timestamp: int64(cFrame.pts),
			}
			ret := videoSender.SendVideoFrame(&videoFrame)
			fmt.Printf("SendVideoFrame %d ret: %d\n", cFrame.pts, ret)
		}
	}
}
