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

	agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/agoraservice"

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
	svcCfg.EnableVideo = true

	agoraservice.Initialize(svcCfg)
	defer agoraservice.Release()
	mediaNodeFactory := agoraservice.NewMediaNodeFactory()
	defer mediaNodeFactory.Release()

	conCfg := agoraservice.RtcConnectionConfig{
		AutoSubscribeAudio: true,
		AutoSubscribeVideo: true,
		ClientRole:         agoraservice.ClientRoleBroadcaster,
		ChannelProfile:     agoraservice.ChannelProfileLiveBroadcasting,
	}
	conSignal := make(chan struct{})
	conHandler := agoraservice.RtcConnectionObserver{
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
	// conCfg.AudioFrameObserver = &agoraservice.RtcConnectionAudioFrameObserver{
	// 	OnPlaybackAudioFrameBeforeMixing: func(con *agoraservice.RtcConnection, channelId string, userId string, frame *agoraservice.PcmAudioFrame) {
	// 		// do something
	// 		fmt.Printf("Playback audio frame before mixing, from userId %s\n", userId)
	// 	},
	// }

	scenario := svcCfg.AudioScenario
	con := agoraservice.NewRtcConnection(&conCfg, scenario)
	defer con.Release()

	localUser := con.GetLocalUser()
	con.RegisterObserver(&conHandler)

	localUser.SetPlaybackAudioFrameBeforeMixingParameters(1, 16000)
	audioSender := mediaNodeFactory.NewAudioPcmDataSender()
	defer audioSender.Release()
	audioTrack := agoraservice.NewCustomAudioTrackPcm(audioSender, scenario)
	defer audioTrack.Release()

	videoSender := mediaNodeFactory.NewVideoEncodedImageSender()
	defer videoSender.Release()
	videoTrack := agoraservice.NewCustomVideoTrackEncoded(videoSender, &agoraservice.VideoEncodedImageSenderOptions{
		CcMode:    agoraservice.VideoSendCcDisabled,
		CodecType: agoraservice.VideoCodecTypeH264,
	})
	defer videoTrack.Release()

	con.Connect(token, channelName, userId)
	<-conSignal

	audioTrack.SetEnabled(true)
	localUser.PublishAudio(audioTrack)
	videoTrack.SetEnabled(true)
	localUser.PublishVideo(videoTrack)

	fn := C.CString("../test_data/test_avsync.mp4")
	defer C.free(unsafe.Pointer(fn))
	decoder := C.open_media_file(fn)
	if decoder == nil {
		fmt.Println("Error opening media file")
		return
	}
	defer C.close_media_file(decoder)

	var cPkt *C.struct__MediaPacket = nil
	cFrame := C.struct__MediaFrame{}
	C.memset(unsafe.Pointer(&cFrame), 0, C.sizeof_struct__MediaFrame)

	firstPts := int64(0)
	firstSendTime := time.Now()
	for !(*bStop) {
		totalSendTime := time.Since(firstSendTime).Milliseconds()
		ret := C.get_packet(decoder, &cPkt)
		if ret != 0 {
			fmt.Println("Finished reading file:", ret)
			break
		}
		if cPkt == nil {
			continue
		}
		if cPkt.media_type == C.AVMEDIA_TYPE_UNKNOWN {
			fmt.Println("Unknown media type")
			C.free_packet(&cPkt)
			continue
		}
		//NOTICE: time stamp must be greater than 0
		// if time stamp is 0, system time will be used in sdk,
		// which will cause frame whose time is less than system time, to be dropped
		if cPkt.pts <= 0 {
			cPkt.pts = 1
		}
		if firstPts == 0 {
			firstPts = int64(cPkt.pts)
			firstSendTime = time.Now()
			totalSendTime = 0
			time.Sleep(50 * time.Millisecond)
			fmt.Println("First pts:", firstPts)
		}
		if int64(cPkt.pts)-firstPts > totalSendTime {
			time.Sleep(50 * time.Millisecond)
			fmt.Println("Sleeping for 50ms")
		}
		if cPkt.media_type == C.AVMEDIA_TYPE_AUDIO {
			ret = C.decode_packet(decoder, cPkt, &cFrame)
			C.free_packet(&cPkt)
			if ret != 0 {
				if ret == C.AVERROR_EAGAIN {
					continue
				}
				fmt.Printf("decode audio finished, %d\n", ret)
				break
			}
			if cFrame.format != C.AV_SAMPLE_FMT_S16 {
				fmt.Println("Unsupported audio format")
				continue
			}
			audioFrame := agoraservice.AudioFrame{
				Type:              agoraservice.AudioFrameTypePCM16,
				SamplesPerChannel: int(cFrame.samples),
				BytesPerSample:    int(cFrame.bytes_per_sample),
				Channels:          int(cFrame.channels),
				SamplesPerSec:     int(cFrame.sample_rate),
				Buffer:            unsafe.Slice((*byte)(unsafe.Pointer(cFrame.buffer)), cFrame.buffer_size),
				RenderTimeMs:      int64(cFrame.pts),
			}
			ret := audioSender.SendAudioPcmData(&audioFrame)
			fmt.Printf("SendPcmData %d ret: %d\n", cFrame.pts, ret)
		} else if cPkt.media_type == C.AVMEDIA_TYPE_VIDEO {
			ret = C.h264_to_annexb(decoder, &cPkt)
			if cPkt == nil {
				continue
			}
			isKeyFrame := cPkt.pkt.flags&C.AV_PKT_FLAG_KEY != 0
			frameType := agoraservice.VideoFrameTypeKeyFrame
			if !isKeyFrame {
				frameType = agoraservice.VideoFrameTypeDeltaFrame
			}
			data := unsafe.Slice((*byte)(unsafe.Pointer(cPkt.pkt.data)), cPkt.pkt.size)
			ret := videoSender.SendEncodedVideoImage(data, &agoraservice.EncodedVideoFrameInfo{
				CodecType:       agoraservice.VideoCodecTypeH264,
				Width:           int(cPkt.width),
				Height:          int(cPkt.height),
				FramesPerSecond: int(cPkt.framerate_num / cPkt.framerate_den),
				FrameType:       frameType,
				CaptureTimeMs:   int64(cPkt.pts),
				DecodeTimeMs:    int64(cPkt.pts),
				Rotation:        agoraservice.VideoOrientation0,
			})
			fmt.Printf("SendVideoFrame %d, data size %d, sync header %x%x%x%x, frametype %d ret: %d\n",
				cPkt.pts, cPkt.pkt.size, data[0], data[1], data[2], data[3], frameType, ret)
			C.free_packet(&cPkt)
		}
	}
}
