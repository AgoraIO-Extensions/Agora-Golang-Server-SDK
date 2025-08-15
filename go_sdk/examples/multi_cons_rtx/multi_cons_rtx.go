package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"time"

	agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/rtc"

	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
)

type GlobalContext struct {
	ctx               context.Context
	cancel            context.CancelFunc
	termSingal        chan os.Signal
	appId             string
	cert              string
	channelNamePrefix string
	mediaNodeFactory  *agoraservice.MediaNodeFactory
	waitTasks         *sync.WaitGroup
	tasks             []*TaskContext
	taskStopSignal    chan int
	audioSenario      agoraservice.AudioScenario
}

type TaskConfig struct {
	// send options
	sendYuv          bool
	sendEncodedVideo bool
	sendPcm          bool
	sendEncodedAudio bool
	sendData         bool
	// recv options
	enableAudioLabel bool
	recvYuv          bool
	recvEncodedVideo bool
	recvPcm          bool
	recvData         bool

	dumpPcm          bool
	dumpYuv          bool
	dumpEncodedVideo bool

	taskTime int64
	pcmFilePath      string
	encodedAudioFilePath string
	yuvFilePath          string
	encodedVideoFilePath string
	role                 int
	sendVideoFps         int
	sendVideoMinBitrate int
	sendVideoBitrate     int
	sendYuvWidth         int
	sendYuvHeight        int
}

func (globalCtx *GlobalContext) genToken(channelName string, userId string) (string, error) {
	token := ""
	if globalCtx.cert != "" {
		tokenExpirationInSeconds := uint32(24 * 3600)
		privilegeExpirationInSeconds := uint32(24 * 3600)
		var err error
		token, err = rtctokenbuilder.BuildTokenWithUserAccount(globalCtx.appId, globalCtx.cert, channelName, userId,
			rtctokenbuilder.RolePublisher, tokenExpirationInSeconds, privilegeExpirationInSeconds)
		if err != nil {
			fmt.Println("Failed to build token: ", err)
			return token, err
		}
	}
	return token, nil
}

func globalInit() *GlobalContext {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	// catch ternimal signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// get environment variable
	appid := os.Getenv("AGORA_APP_ID")
	cert := os.Getenv("AGORA_APP_CERTIFICATE")
	channelName := "gosdktest"
	if appid == "" {
		fmt.Println("Please set AGORA_APP_ID environment variable, and AGORA_APP_CERTIFICATE if needed")
		return nil
	}
	svcCfg := agoraservice.NewAgoraServiceConfig()
	svcCfg.EnableVideo = true
	svcCfg.AppId = appid

	agoraservice.Initialize(svcCfg)
	
	return &GlobalContext{
		ctx:               ctx,
		cancel:            cancel,
		termSingal:        c,
		appId:             appid,
		cert:              cert,
		channelNamePrefix: channelName,
		
		waitTasks:         &sync.WaitGroup{},
		tasks:             nil,
		taskStopSignal:    make(chan int, 100),
		audioSenario:      svcCfg.AudioScenario,
	}
}

func (ctx *GlobalContext) release() {
	
	agoraservice.Release()
}

// uage like:
// ./multi_cons_rtx --channelName=weiqi --sendYuv=true --sendEncodedVideo=true --sendEncodedAudio=true --sendPcm=true --sendData=true --recvYuv=true --recvEncodedVideo=true --recvPcm=true --recvData=true --dumpPcm=true --dumpYuv=true --dumpEncodedVideo=true --taskCount=2 --randTask=true --role=1 --pcmFilePath=./test.pcm --encodedAudioFilePath=./test.aac --yuvFilePath=./test.yuv --encodedVideoFilePath=./test.h264 --sendVideoFps=15 --sendVideoMinBitrate=100 --sendVideoBitrate=500 --sendYuvWidth=640 --sendYuvHeight=360

func main() {
	var enablePprof bool = true
	if enablePprof {
		go func() {
			fmt.Println("**********enable pprof on port 6060**********")
			log.Println(http.ListenAndServe(":6060", nil))
		}()
	}
	

	// parse command options
	var (
		channelName      = flag.String("channelName", "gosdktest", "Channel name")
		sendYuv          = flag.Bool("sendYuv", false, "Enable Send YUV data")
		sendEncodedVideo = flag.Bool("sendEncodedVideo", false, "Enable Send encoded video data")
		sendPcm          = flag.Bool("sendPcm", false, "Enable Send PCM audio data")
		sendEncodedAudio = flag.Bool("sendEncodedAudio", false, "Enable Send encoded audio data")
		sendData         = flag.Bool("sendData", false, "Enable Send custom data")

		enableAudioLabel = flag.Bool("enableAudioLabel", false, "Enable Audio Label")
		recvYuv          = flag.Bool("recvYuv", false, "Enable Receive YUV data")
		recvEncodedVideo = flag.Bool("recvEncodedVideo", false, "Enable Receive encoded video data")
		recvPcm          = flag.Bool("recvPcm", false, "Enable Receive PCM audio data")
		recvData         = flag.Bool("recvData", false, "Enable Receive custom data")

		dumpPcm          = flag.Bool("dumpPcm", false, "Enable Dump PCM audio data")
		dumpYuv          = flag.Bool("dumpYuv", false, "Enable Dump YUV data")
		dumpEncodedVideo = flag.Bool("dumpEncodedVideo", false, "Enable Dump encoded video data")

		taskCount = flag.Int("taskCount", 1, "Task count")
		randTask  = flag.Bool("randTask", false, "Enable Randomly restart task")

		// role
		role        = flag.Int("role", 1, "Descprtion: 1 for host, 0 for client")
		// pcm file path
		pcmFilePath = flag.String("pcmFilePath", "", "Descprtion: Pcm file path")
		// encoded audio file path
		encodedAudioFilePath = flag.String("encodedAudioFilePath", "", "Descprtion: Encoded audio file path")
		// yuv file path
		yuvFilePath = flag.String("yuvFilePath", "", "Descprtion: Yuv file path")
		// encoded video file path
		encodedVideoFilePath = flag.String("encodedVideoFilePath", "", "Descprtion: Encoded video file path")	

		// send video parameter: fps, min bitrate, bitrate
		sendVideoFps = flag.Int("sendVideoFps", 15, "Descprtion: Send video fps")
		sendVideoMinBitrate = flag.Int("sendVideoMinBitrate", 100, "Descprtion: Send video min bitrate")
		sendVideoBitrate = flag.Int("sendVideoBitrate", 500, "Descprtion: Send video bitrate")
		sendYuvWidth = flag.Int("sendYuvWidth", 640, "Descprtion: Send yuv width")
		sendYuvHeight = flag.Int("sendYuvHeight", 360, "Descprtion: Send yuv height")
	)

	flag.Parse()
	
	fmt.Println("channelName:", *channelName)
	fmt.Println("sendYuv:", *sendYuv)
	fmt.Println("sendEncodedVideo:", *sendEncodedVideo)
	fmt.Println("sendPcm:", *sendPcm)
	fmt.Println("sendEncodedAudio:", *sendEncodedAudio)
	fmt.Println("sendData:", *sendData)
	fmt.Println("enableAudioLabel:", *enableAudioLabel)
	fmt.Println("recvYuv:", *recvYuv)
	fmt.Println("recvEncodedVideo:", *recvEncodedVideo)
	fmt.Println("recvPcm:", *recvPcm)
	fmt.Println("recvData:", *recvData)
	fmt.Println("dumpPcm:", *dumpPcm)
	fmt.Println("dumpYuv:", *dumpYuv)
	fmt.Println("dumpEncodedVideo:", *dumpEncodedVideo)
	fmt.Println("taskCount:", *taskCount)
	fmt.Println("randTask:", *randTask)
	fmt.Println("role:", *role)
	fmt.Println("pcmFilePath:", *pcmFilePath)
	fmt.Println("encodedAudioFilePath:", *encodedAudioFilePath)
	fmt.Println("yuvFilePath:", *yuvFilePath)
	fmt.Println("encodedVideoFilePath:", *encodedVideoFilePath)
	fmt.Println("sendVideoFps:", *sendVideoFps)
	fmt.Println("sendVideoMinBitrate:", *sendVideoMinBitrate)
	fmt.Println("sendVideoBitrate:", *sendVideoBitrate)
	fmt.Println("sendYuvWidth:", *sendYuvWidth)
	fmt.Println("sendYuvHeight:", *sendYuvHeight)

	//validity check
	if *sendYuv && *yuvFilePath == "" {
		fmt.Println("yuvFilePath is required when sendYuv is true")
		return
	}
	if *sendEncodedVideo && *encodedVideoFilePath == "" {
		fmt.Println("encodedVideoFilePath is required when sendEncodedVideo is true")
		return
	}
	if *sendPcm && *pcmFilePath == "" {
		fmt.Println("pcmFilePath is required when sendPcm is true")
		return
	}	
	if *sendEncodedAudio && *encodedAudioFilePath == "" {
		fmt.Println("encodedAudioFilePath is required when sendEncodedAudio is true")
		return
	}
	
	

	globalCtx := globalInit()
	if globalCtx == nil {
		return
	}
	defer globalCtx.release()
	globalCtx.channelNamePrefix = *channelName
	globalCtx.tasks = make([]*TaskContext, *taskCount)

	if *enableAudioLabel {
		//agoraservice.EnableExtension("agora.builtin", "agora_audio_label_generator", "", true)
	}

	taskCfg := TaskConfig{
		sendYuv:          *sendYuv,
		sendEncodedVideo: *sendEncodedVideo,
		sendPcm:          *sendPcm,
		sendEncodedAudio: *sendEncodedAudio,
		sendData:         *sendData,

		enableAudioLabel: *enableAudioLabel,
		recvYuv:          *recvYuv,
		recvEncodedVideo: *recvEncodedVideo,
		recvPcm:          *recvPcm,
		recvData:         *recvData,

		dumpPcm:          *dumpPcm,
		dumpYuv:          *dumpYuv,
		dumpEncodedVideo: *dumpEncodedVideo,
		pcmFilePath:      *pcmFilePath,
		encodedAudioFilePath: *encodedAudioFilePath,
		yuvFilePath: *yuvFilePath,
		encodedVideoFilePath: *encodedVideoFilePath,
		role:             *role,
		sendVideoFps:         *sendVideoFps,
		sendVideoMinBitrate: *sendVideoMinBitrate,
		sendVideoBitrate:     *sendVideoBitrate,
		sendYuvWidth: *sendYuvWidth,
		sendYuvHeight: *sendYuvHeight,
	}
	for i := 0; i < *taskCount; i++ {
		globalCtx.waitTasks.Add(1)
		tmpCfg := taskCfg
		if *randTask {
			tmpCfg.taskTime = int64(5 + rand.Intn(10))
		}
		task := globalCtx.newTask(i, &tmpCfg)
		globalCtx.tasks[i] = task
		go func(t *TaskContext) {
			defer globalCtx.waitTasks.Done()
			t.startTask()
		}(task)
	}

	globalCtx.waitTasks.Add(1)
	go func() {
		defer globalCtx.waitTasks.Done()
		stop := false
		for !stop {
			select {
			case taskId := <-globalCtx.taskStopSignal:
				if !(*randTask) {
					break
				}
				tmpCfg := taskCfg
				tmpCfg.taskTime = int64(5 + rand.Intn(10))
				task := globalCtx.newTask(taskId, &tmpCfg)
				globalCtx.tasks[taskId] = task
				globalCtx.waitTasks.Add(1)
				go func(t *TaskContext) {
					defer globalCtx.waitTasks.Done()
					t.startTask()
				}(task)
			case <-globalCtx.ctx.Done():
				stop = true
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-globalCtx.ctx.Done():
				return
			case <-ticker.C:
				printStats()
			}
		}
	}()

	<-globalCtx.termSingal
	globalCtx.cancel()
	globalCtx.waitTasks.Wait()
}
