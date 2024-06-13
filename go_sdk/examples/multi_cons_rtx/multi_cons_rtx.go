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

	agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/agoraservice"

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
	mediaNodeFactory := agoraservice.NewMediaNodeFactory()
	return &GlobalContext{
		ctx:               ctx,
		cancel:            cancel,
		termSingal:        c,
		appId:             appid,
		cert:              cert,
		channelNamePrefix: channelName,
		mediaNodeFactory:  mediaNodeFactory,
		waitTasks:         &sync.WaitGroup{},
		tasks:             nil,
		taskStopSignal:    make(chan int, 100),
	}
}

func (ctx *GlobalContext) release() {
	ctx.mediaNodeFactory.Release()
	agoraservice.Release()
}

func main() {
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

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
	)

	flag.Parse()

	globalCtx := globalInit()
	if globalCtx == nil {
		return
	}
	defer globalCtx.release()
	globalCtx.channelNamePrefix = *channelName
	globalCtx.tasks = make([]*TaskContext, *taskCount)

	if *enableAudioLabel {
		agoraservice.EnableExtension("agora.builtin", "agora_audio_label_generator", "", true)
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

	<-globalCtx.termSingal
	globalCtx.cancel()
	globalCtx.waitTasks.Wait()
}
