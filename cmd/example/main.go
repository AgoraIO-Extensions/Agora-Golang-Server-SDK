package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	agrtm "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/rtm"
)

/*
1. 不在需要MyRtmEventHandler，直接使用agrtm.RtmEventHandler
2. event_handler_adapter.go 不在需要！
*/
func logWithTime(format string, args ...interface{}) {
	fmt.Printf("[%s] %s\n",
		time.Now().Format("2006-01-02 15:04:05.000"),
		fmt.Sprintf(format, args...))
}
func main() {
	// start pprof
	go func() {
		// listen on all interfaces, if you want to listen on localhost, use http.ListenAndServe("localhost:6060", nil)
		// but local host is not accessible from outside!!
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()
	// rtm start
	appId := os.Getenv("APPID")
	userId := os.Getenv("USER_ID")
	token := os.Getenv("TOKEN")
	channelName := os.Getenv("CHANNEL_NAME")
	lenArgs := len(os.Args)
	if lenArgs < 4 {
		fmt.Println("Usage: process <appid> <channelname> <userid> <token_optional>")
		os.Exit(1)
		logWithTime("Usage: process <appid> <channelname> <userid> <token_optional>")
	}

	appId = os.Args[1]
	channelName = os.Args[2]
	userId = os.Args[3]

	if lenArgs >= 5 {
		token = os.Args[4]
	}
	logWithTime("appId: %s, channelName: %s, userId: %s, token: %s\n", appId, channelName, userId, token)

	// 检查参数
	if appId == "" || channelName == "" || userId == "" {
		fmt.Println("参数错误")
		os.Exit(1)
	}
	ret := int(0)
	var requestId uint64

	sign := make(chan struct{})
	msgChan := make(chan struct{})
	var data []byte = make([]byte, 0)
	var rtmClient *agrtm.IRtmClient = nil

	rtmConfig := agrtm.NewRtmConfig()
	//defer rtmConfig.Delete()
	rtmConfig.AppId = appId
	rtmConfig.UserId = userId
	rtmConfig.EventHandler = &agrtm.RtmEventHandler{
		OnLoginResult: func(requestId uint64, errorCode int) {
			fmt.Printf("onLoginResult: requestId=%d, errorCode=%d\n", requestId, errorCode)
			sign <- struct{}{}
		},
		OnLogoutResult: func(requestId uint64, errorCode int) {
			fmt.Printf("onLogoutResult: requestId=%d, errorCode=%d\n", requestId, errorCode)
		},
		OnMessageEvent: func(event *agrtm.MessageEvent) {
			fmt.Printf("onMessageEvent: event=%v\n", event)
			data = event.Message
			logWithTime("send channel message recved: %s", string(data))
			msgChan <- struct{}{}

		},
		OnLinkStateEvent: func(event *agrtm.LinkStateEvent) {
			fmt.Printf("onLinkStateEvent: event=%v\n", event)
		},
		OnSubscribeResult: func(requestId uint64, channelName string, errorCode int) {
			fmt.Printf("onSubscribeResult: requestId=%d, channelName=%s, errorCode=%d\n", requestId, channelName, errorCode)
			sign <- struct{}{}
		},
		OnConnectionStateChanged: func(channelName string, state int, reason int) {
			fmt.Printf("onConnectionStateChanged: channelName=%s, state=%d, reason=%d\n", channelName, state, reason)
		},
	}

	logConfig := agrtm.NewRtmLogConfig()
	logConfig.FilePath = "./logs/rtm.log"
	logConfig.FileSizeInKB = 1024
	logConfig.Level = agrtm.RtmLogLevelINFO
	rtmConfig.LogConfig = logConfig
	fmt.Printf("NewRtmConfig: %+v\n", rtmConfig) //DEBUG

	rtmClient = agrtm.NewRtmClient(rtmConfig)
	logWithTime("CreateAgoraRtmClient: %p\n", rtmClient) //DEBUG

	// set user channel info to event handler

	logWithTime("Login Start: %d\n", ret)
	if token == "" {
		token = appId
	}
	ret, requestId = rtmClient.Login(token)

	fmt.Printf("Login ret: %d, requestId: %d, token: %s\n", ret, requestId, token)

	if ret != 0 {
		panic(ret)
	}
	// wait for login result and timedout to 3 seconds
	select {
	case <-sign:
	case <-time.After(time.Second * 3):
		panic("login timeout")
	}
	logWithTime("login success")

	opt := agrtm.NewSubscribeOptions()

	logWithTime("Subscribe start: %d\n", ret)
	ret, requestId = rtmClient.Subscribe(channelName, opt)
	fmt.Printf("Subscribe ret: %d, requestId: %d\n", ret, requestId)

	if ret != 0 {
		panic(ret)
	}
	// wait for subscribe result and timedout to 3 seconds
	select {
	case <-sign:
	case <-time.After(time.Second * 3):
		panic("subscribe timeout")
	}
	logWithTime("subscribe success")

	//阻塞直到有信号传入
	c := make(chan os.Signal, 1)
	signal.Notify(c)
	logWithTime("rtm client start to work")

waitSignal:
	for { // blocking mdel
		select {
		case signal := <-c:
			if signal == os.Interrupt ||
				signal == os.Kill ||
				signal == syscall.SIGABRT ||
				signal == syscall.SIGTERM {
				logWithTime("exit signal: %v", signal)
				break waitSignal
			}
		case <-msgChan:
			rtmClient.SendChannelMessage(channelName, data)
			logWithTime("send channel message send: %s", string(data))
		}
	}

	//clean
	rtmClient.Logout()
	// wait for logout
	//	time.Sleep(time.Second * 3)
	//unregister event handler

	//release
	rtmClient.Release()
	rtmClient = nil

	// release myEventHandler

	rtmConfig = nil
}
