package main

// #cgo pkg-config: libavformat libavcodec libavutil
// #include <libavformat/avformat.h>
// #include <libavutil/avutil.h>
// #include <libavcodec/avcodec.h>
import "C"
import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"
	"unsafe"

	agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/rtc"
)

type TaskContext struct {
	id        int
	cfg       *TaskConfig
	ctx       context.Context
	globalCtx *GlobalContext

	con *agoraservice.RtcConnection

	streamId int

	dumpPcmFile          *os.File
	dumpYuvFile          *os.File
	dumpEncodedVideoFile *os.File
}

const (
	SendPcmPath = "../test_data/send_audio_16k_1ch.pcm"
	//SendYuvWidth         = 640
	//SendYuvHeight        = 360
	//SendYuvFps           = 15
	//SendYuvBitrate       = 500
	//SendYuvMinBitrate    = 100
	SendYuvPath          = "../test_data/360p_I420.yuv"
	SendEncodedAudioPath = "../test_data/send_audio_16k.aac"
	SendEncodedVideoPath = "../test_data/send_video.h264"
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

func (globalCtx *GlobalContext) newTask(id int, cfg *TaskConfig) *TaskContext {
	taskCtx := &TaskContext{
		id:        id,
		cfg:       cfg,
		globalCtx: globalCtx,
		streamId:  -1,
	}
	return taskCtx
}

func (taskCtx *TaskContext) sendPcm(taskCfg *TaskConfig) {
	ctx := taskCtx.ctx

	con := taskCtx.con

	con.PublishAudio()
	// defer func() {
	// 	senderLocalUser.UnpublishAudio(audioTrack)
	// 	audioTrack.SetEnabled(false)
	// }()

	frame := agoraservice.AudioFrame{
		Type:              agoraservice.AudioFrameTypePCM16,
		SamplesPerChannel: 160,
		BytesPerSample:    2,
		Channels:          1,
		SamplesPerSec:     16000,
		Buffer:            make([]byte, 320),
		RenderTimeMs:     0,
		PresentTimeMs:    0,
	}
	filePath := taskCfg.pcmFilePath

	if filePath == "" {
		fmt.Printf("task %d No pcm file\n", taskCtx.id)
		filePath = SendPcmPath
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("task %d Error opening file: %s\n", taskCtx.id, err.Error())
		return
	}
	defer file.Close()

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	sendCount := 0
	firstSendTime := time.Now()
	for {
		select {
		case <-ticker.C:
			shouldSendCount := int(time.Since(firstSendTime).Milliseconds()/10) - (sendCount - 18)
			for i := 0; i < shouldSendCount; i++ {
				dataLen, err := file.Read(frame.Buffer)
				if err != nil || dataLen < 320 {
					errMsg := ""
					if err != nil {
						errMsg = err.Error()
					}
					fmt.Printf("task %d Finished reading file: %s\n", taskCtx.id, errMsg)
					file.Seek(0, 0)
					i--
					continue
				}

				sendCount++
				con.PushAudioPcmData(frame.Buffer, frame.SamplesPerSec, frame.Channels, 0)
				// fmt.Printf("SendAudioPcmData %d ret: %d\n", sendCount, ret)
			}
			// fmt.Printf("Sent %d frames this time\n", shouldSendCount)
		case <-ctx.Done():
			fmt.Printf("task %d audio sender finished\n", taskCtx.id)
			return
		}
	}
}

func (taskCtx *TaskContext) sendEncodedAudio(taskCfg *TaskConfig) {
	ctx := taskCtx.ctx

	con := taskCtx.con

	con.PublishAudio()

	// defer func() {
	// 	senderLocalUser.UnpublishAudio(encodedAudioTrack)
	// 	encodedAudioTrack.SetEnabled(false)
	// }()

	filePath := taskCfg.encodedAudioFilePath
	if filePath == "" {
		fmt.Printf("task %d No encoded audio file\n", taskCtx.id)
		filePath = SendEncodedAudioPath
	}

	pFormatContext := openMediaFile(filePath)
	if pFormatContext == nil {
		fmt.Printf("task %d Failed to open media file\n", taskCtx.id)
		return
	}
	defer closeMediaFile(&pFormatContext)

	packet := C.av_packet_alloc()
	defer C.av_packet_free(&packet)
	streamInfo := getStreamInfo(pFormatContext)
	codecParam := (*C.struct_AVCodecParameters)(unsafe.Pointer(streamInfo.codecpar))
	tb := streamInfo.time_base

	var sendAudioDuration int64 = 0
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	firstSendTime := time.Now()
	for {
		select {
		case <-ticker.C:
			shouldSendMs := int64(time.Since(firstSendTime).Milliseconds())
			for {
				sendAudioDurationMs := sendAudioDuration * int64(tb.num) * int64(1000) / int64(tb.den)
				if sendAudioDurationMs >= shouldSendMs {
					break
				}
				ret := int(C.av_read_frame(pFormatContext, packet))
				if ret < 0 {
					fmt.Printf("task %d Finished reading file: %d\n", taskCtx.id, ret)
					closeMediaFile(&pFormatContext)
					pFormatContext = openMediaFile(filePath)
					streamInfo = getStreamInfo(pFormatContext)
					codecParam = (*C.struct_AVCodecParameters)(unsafe.Pointer(streamInfo.codecpar))
					continue
				}
				sendAudioDuration += int64(packet.duration)
				data := C.GoBytes(unsafe.Pointer(packet.data), packet.size)
				// if data[0] != 0xFF || (data[1] != 0xF1 && data[1] != 0xF9) {
				// 	fmt.Printf("Invalid aac frame\n")
				// }
				ret = con.PushAudioEncodedData(data, &agoraservice.EncodedAudioFrameInfo{
					Speech:            false,
					Codec:             agoraservice.AudioCodecAacLc,
					SampleRateHz:      int(codecParam.sample_rate),
					SamplesPerChannel: int(codecParam.frame_size / codecParam.ch_layout.nb_channels),
					SendEvenIfEmpty:   true,
					NumberOfChannels:  int(codecParam.ch_layout.nb_channels),
				})
				// fmt.Printf("task %d SendEncodedAudioFrame %d ret: %d\n", taskCtx.id, sendAudioDuration, ret)
				C.av_packet_unref(packet)
			}
		case <-ctx.Done():
			fmt.Printf("task %d encoded audio sender finished\n", taskCtx.id)
			return
		}
	}
}

func (taskCtx *TaskContext) sendYuv(taskCfg *TaskConfig) {
	ctx := taskCtx.ctx

	con := taskCtx.con

	con.SetVideoEncoderConfiguration(&agoraservice.VideoEncoderConfiguration{
		CodecType:         agoraservice.VideoCodecTypeH264,
		Width:             taskCfg.sendYuvWidth,
		Height:            taskCfg.sendYuvHeight,
		Framerate:         taskCfg.sendVideoFps,
		Bitrate:           taskCfg.sendVideoBitrate,
		MinBitrate:        taskCfg.sendVideoMinBitrate,
		OrientationMode:   agoraservice.OrientationModeAdaptive,
		DegradePreference: 0,
	})
	con.PublishVideo()
	// defer func() {
	// 	senderLocalUser.UnpublishVideo(videoTrack)
	// 	videoTrack.SetEnabled(false)
	// }()

	filePath := taskCfg.yuvFilePath
	if filePath == "" {
		fmt.Printf("task %d No yuv file\n", taskCtx.id)
		filePath = SendYuvPath
	}

	w := taskCfg.sendYuvWidth
	h := taskCfg.sendYuvHeight
	dataSize := w * h * 3 / 2
	data := make([]byte, dataSize)
	// read yuv from file 103_RaceHorses_416x240p30_300.yuv
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("task %d Error opening file: %s\n", taskCtx.id, err.Error())
		return
	}
	defer file.Close()

	iterval := 1000 / taskCfg.sendVideoFps
	fmt.Println("iterval:", iterval)

	ticker := time.NewTicker(time.Duration(iterval) * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			dataLen, err := file.Read(data)
			if err != nil || dataLen < dataSize {
				file.Seek(0, 0)
				continue
			}
			// senderCon.SendStreamMessage(streamId, data)
			con.PushVideoFrame(&agoraservice.ExternalVideoFrame{
				Type:      agoraservice.VideoBufferRawData,
				Format:    agoraservice.VideoPixelI420,
				Buffer:    data,
				Stride:    w,
				Height:    h,
				Timestamp: 0,
			})
		case <-ctx.Done():
			fmt.Printf("task %d video sender finished\n", taskCtx.id)
			return
		}
	}
}

func (taskCtx *TaskContext) sendEncodedVideo(taskCfg *TaskConfig) {
	ctx := taskCtx.ctx

	con := taskCtx.con

	con.PublishVideo()
	// defer func() {
	// 	localUser.UnpublishVideo(encodedVideoTrack)
	// 	encodedVideoTrack.SetEnabled(false)
	// }()

	filePath := taskCfg.encodedVideoFilePath
	if filePath == "" {
		fmt.Printf("task %d No encoded video file\n", taskCtx.id)
		filePath = SendEncodedVideoPath
	}

	pFormatContext := openMediaFile(filePath)
	if pFormatContext == nil {
		fmt.Printf("task %d Failed to open media file\n", taskCtx.id)
		return
	}
	defer closeMediaFile(&pFormatContext)

	packet := C.av_packet_alloc()
	defer C.av_packet_free(&packet)
	streamInfo := getStreamInfo(pFormatContext)
	codecParam := (*C.struct_AVCodecParameters)(unsafe.Pointer(streamInfo.codecpar))

	//sendInterval := 1000 * int64(codecParam.framerate.den) / int64(codecParam.framerate.num)
	//change sendInterval to configured value not the file framerate
	sendInterval := 1000 / taskCfg.sendVideoFps
	ticker := time.NewTicker(time.Duration(sendInterval) * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			ret := int(C.av_read_frame(pFormatContext, packet))
			if ret < 0 {
				fmt.Println("Finished reading file:", ret)
				// file.Seek(0, 0)
				closeMediaFile(&pFormatContext)
				pFormatContext = openMediaFile(filePath)
				streamInfo = getStreamInfo(pFormatContext)
				codecParam = (*C.struct_AVCodecParameters)(unsafe.Pointer(streamInfo.codecpar))
				continue
			}

			isKeyFrame := packet.flags&C.AV_PKT_FLAG_KEY != 0
			frameType := agoraservice.VideoFrameTypeKeyFrame
			if !isKeyFrame {
				frameType = agoraservice.VideoFrameTypeDeltaFrame
			}
			data := C.GoBytes(unsafe.Pointer(packet.data), packet.size)
			con.PushVideoEncodedData(data, &agoraservice.EncodedVideoFrameInfo{
				CodecType:       agoraservice.VideoCodecTypeH264,
				Width:           int(codecParam.width),
				Height:          int(codecParam.height),
				FramesPerSecond: int(codecParam.framerate.num / codecParam.framerate.den),
				FrameType:       frameType,
				Rotation:        agoraservice.VideoOrientation0,
			})
			C.av_packet_unref(packet)
			time.Sleep(time.Duration(sendInterval) * time.Millisecond)
		case <-ctx.Done():
			fmt.Printf("task %d encoded video sender finished\n", taskCtx.id)
			return
		}
	}
}

func (taskCtx *TaskContext) sendData() {
	ctx := taskCtx.ctx
	con := taskCtx.con
	id := taskCtx.id

	ticker := time.NewTicker(33 * time.Millisecond)
	defer ticker.Stop()
	msg := []byte(fmt.Sprintf("Hello, Agora! from task %d", id))
	for {
		select {
		case <-ticker.C:
			ret := con.SendStreamMessage(msg)
			fmt.Printf("SendStreamMessage ret: %d, task %d\n", ret, id)
		case <-ctx.Done():
			fmt.Printf("task %d data stream sender finished\n", id)
			return
		}
	}
}

func (taskCtx *TaskContext) startTask() {
	id := taskCtx.id
	globalCtx := taskCtx.globalCtx
	taskCtx.ctx = globalCtx.ctx
	// defer globalCtx.waitTasks.Done()

	cfg := taskCtx.cfg

	var channelName string
	if cfg.role == 1 { // for broadcaster
		channelName = fmt.Sprintf("%s%d", globalCtx.channelNamePrefix, id)
	} else {
		channelName = globalCtx.channelNamePrefix
	}
	senderId := "0"
	token1, err1 := globalCtx.genToken(channelName, senderId)
	if err1 != nil {
		fmt.Printf("Failed to generate token, task %d\n", id)
		return
	}

	if cfg.taskTime > 0 {
		ctx, cancel := context.WithTimeout(taskCtx.ctx, time.Duration(cfg.taskTime)*time.Second)
		taskCtx.ctx = ctx
		defer func() {
			cancel()
			taskCtx.globalCtx.taskStopSignal <- id
		}()
	}
	ctx := taskCtx.ctx

	defer func() {
		fmt.Printf("task %d finisheded\n", id)
		if taskCtx.dumpPcmFile != nil {
			taskCtx.dumpPcmFile.Close()
		}
		if taskCtx.dumpYuvFile != nil {
			taskCtx.dumpYuvFile.Close()
		}
		if taskCtx.dumpEncodedVideoFile != nil {
			taskCtx.dumpEncodedVideoFile.Close()
		}
	}()

	// make sure channel not block callback
	conSignal := make(chan struct{}, 1)
	obs := &agoraservice.RtcConnectionObserver{
		OnConnected: func(con *agoraservice.RtcConnection, info *agoraservice.RtcConnectionInfo, reason int) {
			// do something
			fmt.Printf("task %d Connected, reason %d\n", id, reason)
			conSignal <- struct{}{}
		},
		OnDisconnected: func(con *agoraservice.RtcConnection, info *agoraservice.RtcConnectionInfo, reason int) {
			// do something
			fmt.Printf("task %d Disconnected, reason %d\n", id, reason)
		},
		OnConnecting: func(con *agoraservice.RtcConnection, conInfo *agoraservice.RtcConnectionInfo, reason int) {
			fmt.Printf("task %d Connecting, reason %d\n", id, reason)
		},
		OnReconnecting: func(con *agoraservice.RtcConnection, conInfo *agoraservice.RtcConnectionInfo, reason int) {
			fmt.Printf("task %d Reconnecting, reason %d\n", id, reason)
		},
		OnReconnected: func(con *agoraservice.RtcConnection, conInfo *agoraservice.RtcConnectionInfo, reason int) {
			fmt.Printf("task %d Reconnected, reason %d\n", id, reason)
		},
		OnConnectionLost: func(con *agoraservice.RtcConnection, conInfo *agoraservice.RtcConnectionInfo) {
			fmt.Printf("task %d Connection lost\n", id)
		},
		OnConnectionFailure: func(con *agoraservice.RtcConnection, conInfo *agoraservice.RtcConnectionInfo, errCode int) {
			fmt.Printf("task %d Connection failure, error code %d\n", id, errCode)
		},
		OnUserJoined: func(con *agoraservice.RtcConnection, uid string) {
			fmt.Printf("task %d user %s joined\n", id, uid)
		},
		OnUserLeft: func(con *agoraservice.RtcConnection, uid string, reason int) {
			fmt.Printf("task %d user %s left\n", id, uid)
		},
		OnStreamMessageError: func(con *agoraservice.RtcConnection, uid string, streamId int, errCode int, missed int, cached int) {
			fmt.Printf("task %d send stream message error: %d, channel %s, uid %s\n", id, errCode, channelName, uid)
		},
	}
	// senderLocalUserObs := &agoraservice.LocalUserObserver{}
	var role agoraservice.ClientRole
	if cfg.role == 1 {
		role = agoraservice.ClientRoleBroadcaster
	} else {
		role = agoraservice.ClientRoleAudience
	}

	scenario := globalCtx.audioSenario

	var con *agoraservice.RtcConnection = nil
	conCfg := &agoraservice.RtcConnectionConfig{
		AutoSubscribeAudio: cfg.recvPcm,
		AutoSubscribeVideo: cfg.recvYuv || cfg.recvEncodedVideo,
		ClientRole:         role, //agoraservice.ClientRoleBroadcaster,
		ChannelProfile:     agoraservice.ChannelProfileLiveBroadcasting,
	}
	publishConfig := agoraservice.NewRtcConPublishConfig()
	publishConfig.AudioScenario = scenario
	// disable publish audio and video, and assign value according to task type
	publishConfig.IsPublishAudio = false
	publishConfig.IsPublishVideo = false

	fmt.Printf("task %d config.role: %d, rtc role: %d, channelname: %s\n", id, cfg.role, role, channelName)

	// according to task type, fill publishConfig
	if cfg.sendPcm {
		publishConfig.IsPublishAudio = true
		publishConfig.AudioPublishType = agoraservice.AudioPublishTypePcm
	}
	if cfg.sendEncodedAudio {
		publishConfig.IsPublishAudio = true
		publishConfig.AudioPublishType = agoraservice.AudioPublishTypeEncodedPcm
	}
	if cfg.sendYuv {
		publishConfig.IsPublishVideo = true
		publishConfig.VideoPublishType = agoraservice.VideoPublishTypeYuv
	}
	if cfg.sendEncodedVideo {
		publishConfig.IsPublishVideo = true
		publishConfig.VideoPublishType = agoraservice.VideoPublishTypeEncodedImage
		publishConfig.VideoEncodedImageSenderOptions.CcMode = agoraservice.VideoSendCcEnabled
		publishConfig.VideoEncodedImageSenderOptions.CodecType = agoraservice.VideoCodecTypeH264
		publishConfig.VideoEncodedImageSenderOptions.TargetBitrate = cfg.sendVideoBitrate
	}

	con = agoraservice.NewRtcConnection(conCfg, publishConfig)

	taskCtx.con = con
	defer taskCtx.releaseTask()
	// create datastream
	if cfg.sendData {

		//note :no need to create datastream, it will be created automatically when send stream message
	}

	localUser := con.GetLocalUser()
	con.RegisterObserver(obs)
	if cfg.recvPcm {
		recvAudioFrameObs := &agoraservice.AudioFrameObserver{
			OnPlaybackAudioFrameBeforeMixing: func(localUser *agoraservice.LocalUser, channelId string, userId string, frame *agoraservice.AudioFrame, vadResultStat agoraservice.VadState, vadResultFrame *agoraservice.AudioFrame) bool {
				// do something
				if cfg.enableAudioLabel {
					fmt.Printf("task %d OnPlaybackAudioFrameBeforeMixing, from channel %s, "+
						"userId %s, audio duration %dms, farFieldFlag %d, rms %d, voiceProb %d, musicProb %d, pitch %d\n",
						id, channelId, userId, frame.SamplesPerChannel*1000/frame.SamplesPerSec,
						frame.FarFieldFlag, frame.Rms, frame.VoiceProb, frame.MusicProb, frame.Pitch)
				}
				// fmt.Printf("Playback audio frame before mixing, from channel %s, userId %s, audio duration %dms\n",
				// 	channelId, userId, frame.SamplesPerChannel*1000/frame.SamplesPerSec)
				if userId == senderId {
					return true
				}
				if cfg.dumpPcm {
					if taskCtx.dumpPcmFile == nil {
						var err error
						taskCtx.dumpPcmFile, err = os.OpenFile(fmt.Sprintf("./recv%d.pcm", id), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
						if err != nil {
							fmt.Printf("task %d Failed to open dump file, %s", id, err.Error())
							return false
						}
					}
					taskCtx.dumpPcmFile.Write(frame.Buffer)
				}
				return true
			},
		}
		localUser.SetPlaybackAudioFrameBeforeMixingParameters(1, 16000)
		con.RegisterAudioFrameObserver(recvAudioFrameObs, 0, nil)
	}
	if cfg.recvYuv {
		recvVideoFrameObs := &agoraservice.VideoFrameObserver{
			OnFrame: func(channelId string, userId string, frame *agoraservice.VideoFrame) bool {
				// do something
				// fmt.Printf("recv video frame, from channel %s, user %s, video size %dx%d\n",
				// 	channelId, userId, frame.Width, frame.Height)
				if userId == senderId {
					return true
				}
				if cfg.dumpYuv {
					if taskCtx.dumpYuvFile == nil {
						var err error
						taskCtx.dumpYuvFile, err = os.OpenFile(
							fmt.Sprintf("./recv%d_%dx%d.yuv", id, frame.Width, frame.Height),
							os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
						if err != nil {
							fmt.Printf("task %d Failed to open dump file, %s", id, err.Error())
							return false
						}
					}
					taskCtx.dumpYuvFile.Write(frame.YBuffer)
					taskCtx.dumpYuvFile.Write(frame.UBuffer)
					taskCtx.dumpYuvFile.Write(frame.VBuffer)
				}
				return true
			},
		}
		con.RegisterVideoFrameObserver(recvVideoFrameObs)
	}
	localUserObs := &agoraservice.LocalUserObserver{}
	if cfg.recvEncodedVideo {
		localUserObs.OnUserVideoTrackSubscribed = func(localUser *agoraservice.LocalUser, uid string, info *agoraservice.VideoTrackInfo, remoteVideoTrack *agoraservice.RemoteVideoTrack) {
			fmt.Printf("task %d user %s video subscribed\n", id, uid)
		}
		encodedVideoObserver := &agoraservice.VideoEncodedFrameObserver{
			OnEncodedVideoFrame: func(uid string, imageBuffer []byte, frameInfo *agoraservice.EncodedVideoFrameInfo) bool {
				if uid == senderId {
					return true
				}
				// fmt.Printf("user %s encoded video received, frame len %d\n", uid, len(imageBuffer))
				if cfg.dumpEncodedVideo {
					if taskCtx.dumpEncodedVideoFile == nil {
						var err error
						taskCtx.dumpEncodedVideoFile, err = os.OpenFile(fmt.Sprintf("./recv%d.h264", id), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
						if err != nil {
							fmt.Printf("task %d Failed to open dump file, %s", id, err.Error())
							return false
						}
					}
					taskCtx.dumpEncodedVideoFile.Write(imageBuffer)
				}
				return true
			},
		}
		subvideoopt := &agoraservice.VideoSubscriptionOptions{
			StreamType:       agoraservice.VideoStreamHigh,
			EncodedFrameOnly: true,
		}
		localUser.SubscribeAllVideo(subvideoopt)
		con.RegisterVideoEncodedFrameObserver(encodedVideoObserver)
	}
	if cfg.recvData {
		localUserObs.OnStreamMessage = func(localUser *agoraservice.LocalUser, uid string, streamId int, data []byte) {
			fmt.Printf("task %d recv stream message: %s, channel %s, uid %s\n", id, string(data), channelName, uid)
		}
	}
	con.RegisterLocalUserObserver(localUserObs)

	con.Connect(token1, channelName, senderId)
	// defer con.Disconnect()

	select {
	case <-conSignal:
	case <-time.After(25 * time.Second):
		fmt.Printf("task %d failed to connect after 25 seconds\n", id)
		return
	}

	// send audio
	waitGroup := &sync.WaitGroup{}
	if cfg.sendPcm {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			fmt.Printf("task %d send pcm\n", id)
			taskCtx.sendPcm(cfg)
		}()
	}

	// send encoded audio
	if cfg.sendEncodedAudio {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			fmt.Printf("task %d send encoded audio\n", id)
			taskCtx.sendEncodedAudio(cfg)
		}()
	}

	// send video
	if cfg.sendYuv {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			fmt.Printf("task %d send yuv\n", id)
			taskCtx.sendYuv(cfg)
		}()
	}

	// send encoded video
	if cfg.sendEncodedVideo {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			fmt.Printf("task %d send encoded video\n", id)
			taskCtx.sendEncodedVideo(cfg)
		}()
	}

	// send datastream
	if taskCtx.streamId >= 0 {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			taskCtx.sendData()
		}()
	}
	// make sure at least one waiting task
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		select {
		case <-ctx.Done():
			// fmt.Printf("task %d finished\n", id)
		}
	}()
	waitGroup.Wait()
	fmt.Printf("task %d finished\n", id)
}

func (taskCtx *TaskContext) releaseTask() {
	fmt.Printf("task %d release\n", taskCtx.id)
	if taskCtx.con == nil {
		fmt.Printf("task %d release empty connection\n", taskCtx.id)
		return
	}

	// for relase ,only disconnect and release,do not care unregister!!

	taskCtx.con.Disconnect()

	taskCtx.con.Release()
	taskCtx.con = nil
	fmt.Printf("task %d released\n", taskCtx.id)
}
