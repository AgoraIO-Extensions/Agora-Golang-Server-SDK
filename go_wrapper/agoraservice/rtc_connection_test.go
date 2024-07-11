package agoraservice

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

func TestBaseCase(t *testing.T) {
	// Test code here
	t.Log("Test case executed")
	svcCfg := AgoraServiceConfig{
		AppId:         "aab8b8f5a8cd4469a63042fcfafe7063",
		AudioScenario: AUDIO_SCENARIO_CHORUS,
		LogPath:       "./agora_rtc_log/agorasdk.log",
		LogSize:       512 * 1024,
	}
	Init(&svcCfg)
	senderCfg := RtcConnectionConfig{
		SubAudio:       false,
		SubVideo:       false,
		ClientRole:     1,
		ChannelProfile: 1,
		ConnectionHandler: &RtcConnectionEventHandler{
			OnConnected: func(con *RtcConnection, info *RtcConnectionInfo, reason int) {
				t.Log("sender Connected")
			},
			OnDisconnected: func(con *RtcConnection, info *RtcConnectionInfo, reason int) {
				t.Log("sender Disconnected")
			},
		},
	}
	senderCon := NewConnection(&senderCfg)
	defer senderCon.Release()
	sender := senderCon.NewPcmSender()
	defer sender.Release()
	senderCon.Connect("", "lhzuttest", "111")
	sender.Start()
	var stopSend *bool = new(bool)
	*stopSend = false
	waitSenderStop := &sync.WaitGroup{}
	waitSenderStop.Add(1)
	go func() {
		defer waitSenderStop.Done()
		data := make([]byte, 320)
		for !*stopSend {
			sender.SendPcmData(&PcmAudioFrame{
				Data:              data,
				Timestamp:         0,
				SamplesPerChannel: 160,
				BytesPerSample:    2,
				NumberOfChannels:  1,
				SampleRate:        16000,
			})
			time.Sleep(10 * time.Millisecond)
		}
		sender.Stop()
	}()

	waitSenderStop.Add(1)
	go func() {
		defer waitSenderStop.Done()
		streamId, ret := senderCon.CreateDataStream(true, true)
		if ret != 0 {
			t.Log("create stream failed")
			return
		}
		data := make([]byte, 256)
		for !*stopSend {
			senderCon.SendStreamMessage(streamId, data)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	waitSenderStop.Add(1)
	go func() {
		defer func() {
			senderCon.ReleaseVideoSender()
			waitSenderStop.Done()
		}()
		sender := senderCon.GetVideoSender()
		w := 416
		h := 240
		dataSize := w * h * 3 / 2
		data := make([]byte, dataSize)
		// read yuv from file 103_RaceHorses_416x240p30_300.yuv
		file, err := os.Open("../../test_data/103_RaceHorses_416x240p30_300.yuv")
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		sender.SetVideoEncoderConfig(&VideoEncoderConfig{
			CodecType:         2,
			Width:             320,
			Height:            240,
			Framerate:         30,
			Bitrate:           500,
			MinBitrate:        100,
			OrientationMode:   0,
			DegradePreference: 0,
		})
		sender.Start()
		for !*stopSend {
			dataLen, err := file.Read(data)
			if err != nil || dataLen < dataSize {
				file.Seek(0, 0)
				continue
			}
			// senderCon.SendStreamMessage(streamId, data)
			sender.SendVideoFrame(&VideoFrame{
				Buffer:    data,
				Width:     w,
				Height:    h,
				YStride:   w,
				UStride:   w / 2,
				VStride:   w / 2,
				Timestamp: 0,
			})
			time.Sleep(33 * time.Millisecond)
		}
		sender.Stop()
	}()

	waitSenderJoin := make(chan struct{}, 1)
	waitForVideo := make(chan struct{}, 1)
	recvVideo := new(bool)
	*recvVideo = false
	waitForAudio := make(chan struct{}, 1)
	recvAudio := new(bool)
	*recvAudio = false
	waitForData := make(chan struct{}, 1)
	recvData := new(bool)
	*recvData = false
	recvCfg := RtcConnectionConfig{
		SubAudio:       true,
		SubVideo:       true,
		ClientRole:     2,
		ChannelProfile: 1,
		ConnectionHandler: &RtcConnectionEventHandler{
			OnConnected: func(con *RtcConnection, info *RtcConnectionInfo, reason int) {
				t.Log("recver Connected")
			},
			OnDisconnected: func(con *RtcConnection, info *RtcConnectionInfo, reason int) {
				t.Log("recver Disconnected")
			},
			OnUserJoined: func(con *RtcConnection, uid string) {
				t.Log("user joined, " + uid)
				waitSenderJoin <- struct{}{}
			},
			OnUserLeft: func(con *RtcConnection, uid string, reason int) {
				t.Log("user left, " + uid)
			},
			OnStreamMessage: func(con *RtcConnection, uid string, streamId int, data []byte) {
				t.Log("stream message")
				if !*recvData {
					*recvData = true
					waitForData <- struct{}{}
				}
			},
			OnStreamMessageError: func(con *RtcConnection, uid string, streamId int, errCode int, missed int, cached int) {
				t.Log("stream message error")
			},
		},
		AudioFrameObserver: &RtcConnectionAudioFrameObserver{
			OnPlaybackAudioFrameBeforeMixing: func(con *RtcConnection, channelId string, userId string, frame *PcmAudioFrame) {
				t.Log("Playback audio frame before mixing")
				if !*recvAudio {
					*recvAudio = true
					waitForAudio <- struct{}{}
				}
			},
		},
		VideoFrameObserver: &RtcConnectionVideoFrameObserver{
			OnFrame: func(con *RtcConnection, channelId, userId string, frame *VideoFrame) {
				t.Log("on video frame")
				if !*recvVideo {
					defer func() {
						*recvVideo = true
						waitForVideo <- struct{}{}
					}()
					// write frame to file
					file, err := os.Create(fmt.Sprintf("recv_%dx%d.yuv", frame.Width, frame.Height))
					if err != nil {
						fmt.Println("Error opening file:", err)
						return
					}
					defer file.Close()
					file.Write(frame.Buffer)
				}
			},
		},
	}
	recvCon := NewConnection(&recvCfg)
	defer recvCon.Release()
	recvCon.SetParameters("{\"rtc.video.playout_delay_max\": 250}")
	recvCon.SetParameters("{\"rtc.video.broadcaster_playout_delay_max\": 250}")
	recvCon.Connect("", "lhzuttest", "222")
	timer := time.NewTimer(10 * time.Second)
	for *recvAudio == false || *recvData == false || *recvVideo == false {
		select {
		case <-waitSenderJoin:
		case <-waitForAudio:
		case <-waitForData:
		case <-waitForVideo:
		case <-timer.C:
			t.Error("wait video, audio or data timeout, recvVideo: ", *recvVideo, ", recvAudio: ", *recvAudio, ", recvData: ", *recvData)
			t.Fail()
			break
		}
	}
	*stopSend = true
	waitSenderStop.Wait()
	senderCon.Disconnect()
	recvCon.Disconnect()
}

func TestDatastreamCase(t *testing.T) {
	// Test code here
	t.Log("Test case executed")
	svcCfg := AgoraServiceConfig{
		AppId: "aab8b8f5a8cd4469a63042fcfafe7063",
	}
	Init(&svcCfg)
	senderCfg := RtcConnectionConfig{
		SubAudio:       false,
		SubVideo:       false,
		ClientRole:     1,
		ChannelProfile: 1,
		ConnectionHandler: &RtcConnectionEventHandler{
			OnConnected: func(con *RtcConnection, info *RtcConnectionInfo, reason int) {
				t.Log("sender Connected")
			},
			OnDisconnected: func(con *RtcConnection, info *RtcConnectionInfo, reason int) {
				t.Log("sender Disconnected")
			},
		},
	}
	var stopSend *bool = new(bool)
	*stopSend = false
	waitSenderStop := &sync.WaitGroup{}
	const conNum = 10
	senderCons := make([]*RtcConnection, conNum)
	for i := 0; i < conNum; i++ {
		senderCon := NewConnection(&senderCfg)
		senderCons[i] = senderCon
		senderCon.Connect("", fmt.Sprintf("lhzuttest%d", i), fmt.Sprintf("111%d", i))
		for j := 0; j < 10; j++ {
			streamId, ret := senderCon.CreateDataStream(true, true)
			if ret != 0 {
				t.Logf("connection %d create stream %d failed, error %d\n", i, j, ret)
				continue
			}
			waitSenderStop.Add(1)
			go func(con *RtcConnection, streamId int) {
				defer waitSenderStop.Done()
				for !*stopSend {
					dataStr := fmt.Sprintf("connection %p stream %d", con, streamId)
					con.SendStreamMessage(streamId, []byte(dataStr))
					time.Sleep(100 * time.Millisecond)
				}
			}(senderCon, streamId)
		}
	}

	msgesMutex := &sync.Mutex{}
	recvMsgs := make(map[string]struct{}, 0)
	recvCfg := RtcConnectionConfig{
		SubAudio:       true,
		SubVideo:       false,
		ClientRole:     2,
		ChannelProfile: 1,
		ConnectionHandler: &RtcConnectionEventHandler{
			OnConnected: func(con *RtcConnection, info *RtcConnectionInfo, reason int) {
				t.Log("recver Connected")
			},
			OnDisconnected: func(con *RtcConnection, info *RtcConnectionInfo, reason int) {
				t.Log("recver Disconnected")
			},
			OnUserJoined: func(con *RtcConnection, uid string) {
				t.Log("user joined, " + uid)
			},
			OnUserLeft: func(con *RtcConnection, uid string, reason int) {
				t.Log("user left, " + uid)
			},
			OnStreamMessage: func(con *RtcConnection, uid string, streamId int, data []byte) {
				msg := string(data)
				msgesMutex.Lock()
				recvMsgs[msg] = struct{}{}
				msgesMutex.Unlock()
				// t.Log("stream message: ", msg)
			},
			OnStreamMessageError: func(con *RtcConnection, uid string, streamId int, errCode int, missed int, cached int) {
				t.Log("stream message error")
			},
		},
		AudioFrameObserver: &RtcConnectionAudioFrameObserver{
			OnPlaybackAudioFrameBeforeMixing: func(con *RtcConnection, channelId string, userId string, frame *PcmAudioFrame) {
				t.Log("Playback audio frame before mixing")
			},
		},
	}
	recvCons := make([]*RtcConnection, conNum)
	for i := 0; i < conNum; i++ {
		recvCon := NewConnection(&recvCfg)
		recvCons[i] = recvCon
		recvCon.Connect("", fmt.Sprintf("lhzuttest%d", i), fmt.Sprintf("222%d", i))
	}

	time.Sleep(10 * time.Second)
	*stopSend = true
	waitSenderStop.Wait()

	for i := 0; i < conNum; i++ {
		recvCons[i].Disconnect()
		recvCons[i].Release()
		senderCons[i].Disconnect()
		senderCons[i].Release()
	}

	t.Log("recvMsgs count: ", len(recvMsgs))
	t.Log("recvMsgs: ", recvMsgs)
}

func TestVadCase(t *testing.T) {
	// Test code here
	t.Log("Test case executed")
	svcCfg := AgoraServiceConfig{
		AppId: "aab8b8f5a8cd4469a63042fcfafe7063",
	}
	Init(&svcCfg)
	senderCfg := RtcConnectionConfig{
		SubAudio:       false,
		SubVideo:       false,
		ClientRole:     1,
		ChannelProfile: 1,
		ConnectionHandler: &RtcConnectionEventHandler{
			OnConnected: func(con *RtcConnection, info *RtcConnectionInfo, reason int) {
				t.Log("sender Connected")
			},
			OnDisconnected: func(con *RtcConnection, info *RtcConnectionInfo, reason int) {
				t.Log("sender Disconnected")
			},
		},
	}
	senderCon := NewConnection(&senderCfg)
	defer senderCon.Release()
	sender := senderCon.NewPcmSender()
	defer sender.Release()
	senderCon.Connect("", "lhzuttest", "111")
	sender.Start()
	var stopSend *bool = new(bool)
	*stopSend = false
	waitSenderStop := &sync.WaitGroup{}
	waitSenderStop.Add(1)
	go func() {
		defer waitSenderStop.Done()
		file, err := os.Open("../../test_data/demo.raw")
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		data := make([]byte, 320)
		for !*stopSend {
			dataLen, err := file.Read(data)
			if err != nil || dataLen < 320 {
				break
			}
			sender.SendPcmData(&PcmAudioFrame{
				Data:              data,
				Timestamp:         0,
				SamplesPerChannel: 160,
				BytesPerSample:    2,
				NumberOfChannels:  1,
				SampleRate:        16000,
			})
			time.Sleep(10 * time.Millisecond)
		}
		sender.Stop()
	}()

	vad := NewAudioVad(&AudioVadConfig{
		StartRecognizeCount:    10,
		StopRecognizeCount:     6,
		PreStartRecognizeCount: 10,
		ActivePercent:          0.6,
		InactivePercent:        0.2,
	})
	defer vad.Release()
	recvCfg := RtcConnectionConfig{
		SubAudio:       true,
		SubVideo:       false,
		ClientRole:     2,
		ChannelProfile: 1,
		//NOTICE: the input audio format of vad is fixed to 16k, 1 channel, 16bit
		SubAudioConfig: &SubscribeAudioConfig{
			SampleRate: 16000,
			Channels:   1,
		},
		ConnectionHandler: &RtcConnectionEventHandler{
			OnConnected: func(con *RtcConnection, info *RtcConnectionInfo, reason int) {
				t.Log("recver Connected")
			},
			OnDisconnected: func(con *RtcConnection, info *RtcConnectionInfo, reason int) {
				t.Log("recver Disconnected")
			},
			OnUserJoined: func(con *RtcConnection, uid string) {
				t.Log("user joined, " + uid)
			},
			OnUserLeft: func(con *RtcConnection, uid string, reason int) {
				t.Log("user left, " + uid)
			},
			OnStreamMessage: func(con *RtcConnection, uid string, streamId int, data []byte) {
				t.Log("stream message")
			},
			OnStreamMessageError: func(con *RtcConnection, uid string, streamId int, errCode int, missed int, cached int) {
				t.Log("stream message error")
			},
		},
		AudioFrameObserver: &RtcConnectionAudioFrameObserver{
			OnPlaybackAudioFrameBeforeMixing: func(con *RtcConnection, channelId string, userId string, frame *PcmAudioFrame) {
				// t.Log("Playback audio frame before mixing")
				out, ret := vad.ProcessPcmFrame(frame)
				if ret < 0 {
					t.Log("vad process frame failed")
					t.Fail()
				}
				if out != nil {
					t.Logf("vad state %d, out frame time: %d, duration %d\n", ret, out.Timestamp, out.SamplesPerChannel/16)
				} else {
					t.Logf("vad state %d\n", ret)
				}
			},
		},
		VideoFrameObserver: nil,
	}
	recvCon := NewConnection(&recvCfg)
	defer recvCon.Release()
	recvCon.Connect("", "lhzuttest", "222")

	waitSenderStop.Wait()
	time.Sleep(5 * time.Second)
	senderCon.Disconnect()
	recvCon.Disconnect()
}

func TestSubAudio(t *testing.T) {
	// Test code here
	t.Log("Test case executed")
	svcCfg := AgoraServiceConfig{
		AppId: "aab8b8f5a8cd4469a63042fcfafe7063",
	}
	Init(&svcCfg)
	senderCfg := RtcConnectionConfig{
		SubAudio:       false,
		SubVideo:       false,
		ClientRole:     1,
		ChannelProfile: 1,
		ConnectionHandler: &RtcConnectionEventHandler{
			OnConnected: func(con *RtcConnection, info *RtcConnectionInfo, reason int) {
				t.Log("sender Connected")
			},
			OnDisconnected: func(con *RtcConnection, info *RtcConnectionInfo, reason int) {
				t.Log("sender Disconnected")
			},
		},
	}
	var stopSend *bool = new(bool)
	*stopSend = false
	waitSenderStop := &sync.WaitGroup{}
	const conNum = 5
	senderCons := make([]*RtcConnection, conNum)
	for i := 0; i < conNum; i++ {
		senderCon := NewConnection(&senderCfg)
		senderCons[i] = senderCon
		senderCon.Connect("", "lhzuttestsubaudio", fmt.Sprintf("111%d", i))
		waitSenderStop.Add(1)
		go func(con *RtcConnection) {
			defer waitSenderStop.Done()
			file, err := os.Open("../../test_data/demo.pcm")
			if err != nil {
				fmt.Println("Error opening file:", err)
				return
			}
			defer file.Close()
			sender := con.NewPcmSender()
			defer sender.Release()
			sender.Start()

			data := make([]byte, 320)
			for !*stopSend {
				dataLen, err := file.Read(data)
				if err != nil || dataLen < 320 {
					file.Seek(0, 0)
					break
				}
				sender.SendPcmData(&PcmAudioFrame{
					Data:              data,
					Timestamp:         0,
					SamplesPerChannel: 160,
					BytesPerSample:    2,
					NumberOfChannels:  1,
					SampleRate:        16000,
				})
				time.Sleep(10 * time.Millisecond)
			}
			sender.Stop()
		}(senderCon)
	}

	var setRecvState *bool = new(bool)
	*setRecvState = true
	recvedMutex := &sync.Mutex{}
	recvedUsers := make(map[string]struct{}, conNum)
	recvCfg := RtcConnectionConfig{
		SubAudio:       true,
		SubVideo:       false,
		ClientRole:     2,
		ChannelProfile: 1,
		//NOTICE: the input audio format of vad is fixed to 16k, 1 channel, 16bit
		SubAudioConfig: &SubscribeAudioConfig{
			SampleRate: 16000,
			Channels:   1,
		},
		ConnectionHandler: &RtcConnectionEventHandler{
			OnConnected: func(con *RtcConnection, info *RtcConnectionInfo, reason int) {
				t.Log("recver Connected")
			},
			OnDisconnected: func(con *RtcConnection, info *RtcConnectionInfo, reason int) {
				t.Log("recver Disconnected")
			},
			OnUserJoined: func(con *RtcConnection, uid string) {
				t.Log("user joined, " + uid)
			},
			OnUserLeft: func(con *RtcConnection, uid string, reason int) {
				t.Log("user left, " + uid)
			},
			OnStreamMessage: func(con *RtcConnection, uid string, streamId int, data []byte) {
				t.Log("stream message")
			},
			OnStreamMessageError: func(con *RtcConnection, uid string, streamId int, errCode int, missed int, cached int) {
				t.Log("stream message error")
			},
		},
		AudioFrameObserver: &RtcConnectionAudioFrameObserver{
			OnPlaybackAudioFrameBeforeMixing: func(con *RtcConnection, channelId string, userId string, frame *PcmAudioFrame) {
				// t.Log("Playback audio frame before mixing")
				if !*setRecvState {
					return
				}
				recvedMutex.Lock()
				recvedUsers[userId] = struct{}{}
				recvedMutex.Unlock()
			},
		},
		VideoFrameObserver: nil,
	}
	recvCon := NewConnection(&recvCfg)
	defer recvCon.Release()
	recvCon.Connect("", "lhzuttestsubaudio", "222")

	time.Sleep(5 * time.Second)
	if len(recvedUsers) != conNum {
		t.Error("not all users received audio")
		t.Fail()
	} else {
		t.Log("all users received audio")
	}
	*setRecvState = false
	time.Sleep(1 * time.Second)
	recvedMutex.Lock()
	// wrong code: recvedUsers = make(map[string]struct{}, 0)
	for i := 0; i < conNum; i++ {
		delete(recvedUsers, fmt.Sprintf("111%d", i))
	}
	recvedMutex.Unlock()
	recvCon.UnsubscribeAllAudio()
	// make sure no callback after unsubscribe
	time.Sleep(2 * time.Second)
	*setRecvState = true
	time.Sleep(2 * time.Second)
	if len(recvedUsers) != 0 {
		t.Error("after unsubscribe all audio, still received audio")
		t.Fail()
	} else {
		t.Log("after unsubscribe all audio, no audio received")
	}

	recvCon.SubscribeAudio("1110")
	time.Sleep(2 * time.Second)
	if len(recvedUsers) != 1 {
		t.Error("after subscribe audio, not received audio")
		t.Fail()
	} else {
		t.Log("after subscribe 1110 audio, received audio")
	}
	*setRecvState = false
	time.Sleep(1 * time.Second)
	recvedMutex.Lock()
	delete(recvedUsers, "1110")
	recvedMutex.Unlock()
	recvCon.UnsubscribeAudio("1110")
	time.Sleep(2 * time.Second)
	*setRecvState = true
	time.Sleep(2 * time.Second)
	if len(recvedUsers) != 0 {
		t.Error("after unsubscribe 1110 audio, still received audio")
		t.Fail()
	} else {
		t.Log("after unsubscribe 1110 audio, no audio received")
	}

	*stopSend = true
	waitSenderStop.Wait()
	for i := 0; i < conNum; i++ {
		senderCons[i].Disconnect()
		senderCons[i].Release()
	}
	recvCon.Disconnect()
}

func TestReconnect(t *testing.T) {
	// Test code here
	t.Log("Test case executed")
	svcCfg := AgoraServiceConfig{
		AppId: "aab8b8f5a8cd4469a63042fcfafe7063",
	}
	Init(&svcCfg)
	connectedSignal := make(chan struct{}, 1)
	senderCfg := RtcConnectionConfig{
		SubAudio:       true,
		SubVideo:       true,
		ClientRole:     1,
		ChannelProfile: 1,
		ConnectionHandler: &RtcConnectionEventHandler{
			OnConnected: func(con *RtcConnection, info *RtcConnectionInfo, reason int) {
				t.Log("sender Connected")
				connectedSignal <- struct{}{}
			},
			OnDisconnected: func(con *RtcConnection, info *RtcConnectionInfo, reason int) {
				t.Log("sender Disconnected")
			},
		},
	}
	con := NewConnection(&senderCfg)
	defer con.Release()
	con.Connect("", "lhzuttestreconnect", "111")

	select {
	case <-connectedSignal:
	case <-time.After(10 * time.Second):
		t.Error("connect timeout")
		t.Fail()
	}

	con.Disconnect()

	con.Connect("", "lhzuttestreconnect", "111")

	select {
	case <-connectedSignal:
	case <-time.After(10 * time.Second):
		t.Error("connect timeout")
		t.Fail()
	}

	con.Disconnect()
}

func TestVad1(t *testing.T) {
	agoraVad := NewAudioVad(&AudioVadConfig{
		StartRecognizeCount:    10,
		StopRecognizeCount:     6,
		PreStartRecognizeCount: 10,
		ActivePercent:          0.6,
		InactivePercent:        0.2,
	})
	defer agoraVad.Release()
	f, err := os.ReadFile("../../test_data/demo.raw")
	if err != nil {
		t.Error(err)
	}
	for len(f) > 0 {
		dataSize := 320
		if len(f) < dataSize {
			break
		} else {
			data := f[:dataSize]
			f = f[dataSize:]
			out, ret := agoraVad.ProcessPcmFrame(&PcmAudioFrame{
				Data:              data,
				SampleRate:        16000,
				NumberOfChannels:  1,
				BytesPerSample:    2,
				SamplesPerChannel: 160,
			})
			t.Log(ret, len(out.Data))
		}
	}
}

func TestConnLost(t *testing.T) {
	// Test code here
	t.Log("Test case executed")
	svcCfg := AgoraServiceConfig{
		AppId: "aab8b8f5a8cd4469a63042fcfafe6666",
	}
	Init(&svcCfg)
	defer Destroy()

	connectedSignal := make(chan struct{}, 1)
	senderCfg := RtcConnectionConfig{
		SubAudio:       true,
		SubVideo:       true,
		ClientRole:     1,
		ChannelProfile: 1,
		ConnectionHandler: &RtcConnectionEventHandler{
			OnConnected: func(con *RtcConnection, info *RtcConnectionInfo, reason int) {
				t.Log("sender Connected")
				connectedSignal <- struct{}{}
			},
			OnDisconnected: func(con *RtcConnection, info *RtcConnectionInfo, reason int) {
				t.Log("sender Disconnected ", reason)
			},
			OnConnectionLost: func(con *RtcConnection, info *RtcConnectionInfo) {
				t.Log("sender ConnectionLost")
			},
			OnConnectionFailure: func(con *RtcConnection, info *RtcConnectionInfo, reason int) {
				t.Log("sender ConnectionFailure ", reason)
			},
		},
	}
	con := NewConnection(&senderCfg)
	defer con.Release()
	con.Connect("", "lhzuttestreconnect", "111")

	select {
	case <-connectedSignal:
		t.Fail()
	case <-time.After(10 * time.Second):
		t.Log("connect timeout")
	}

	con.Disconnect()
}

func TestUserInfoUpdated(t *testing.T) {
	// Test code here
	t.Log("Test case executed")
	svcCfg := AgoraServiceConfig{
		AppId:         "aab8b8f5a8cd4469a63042fcfafe7063",
		AudioScenario: AUDIO_SCENARIO_CHORUS,
		LogPath:       "./agora_rtc_log/agorasdk.log",
		LogSize:       512 * 1024,
	}
	Init(&svcCfg)
	senderCfg := RtcConnectionConfig{
		SubAudio:       false,
		SubVideo:       false,
		ClientRole:     1,
		ChannelProfile: 1,
		ConnectionHandler: &RtcConnectionEventHandler{
			OnConnected: func(con *RtcConnection, info *RtcConnectionInfo, reason int) {
				t.Log("sender Connected")
			},
			OnDisconnected: func(con *RtcConnection, info *RtcConnectionInfo, reason int) {
				t.Log("sender Disconnected")
			},
		},
	}
	senderCon := NewConnection(&senderCfg)
	defer senderCon.Release()
	sender := senderCon.NewPcmSender()
	defer sender.Release()
	senderCon.Connect("", "lhzuttest", "111")
	sender.Start()
	var stopSend *bool = new(bool)
	*stopSend = false
	waitSenderStop := &sync.WaitGroup{}
	waitSenderStop.Add(1)
	go func() {
		defer waitSenderStop.Done()
		data := make([]byte, 320)
		for !*stopSend {
			sender.SendPcmData(&PcmAudioFrame{
				Data:              data,
				Timestamp:         0,
				SamplesPerChannel: 160,
				BytesPerSample:    2,
				NumberOfChannels:  1,
				SampleRate:        16000,
			})
			time.Sleep(10 * time.Millisecond)
		}
		sender.Stop()
	}()

	waitSenderJoin := make(chan struct{}, 1)
	audioMuteState := make(chan int, 1)
	recvCfg := RtcConnectionConfig{
		SubAudio:       true,
		SubVideo:       true,
		ClientRole:     2,
		ChannelProfile: 1,
		ConnectionHandler: &RtcConnectionEventHandler{
			OnConnected: func(con *RtcConnection, info *RtcConnectionInfo, reason int) {
				t.Log("recver Connected")
			},
			OnDisconnected: func(con *RtcConnection, info *RtcConnectionInfo, reason int) {
				t.Log("recver Disconnected")
			},
			OnUserJoined: func(con *RtcConnection, uid string) {
				t.Log("user joined, " + uid)
				waitSenderJoin <- struct{}{}
			},
			OnUserLeft: func(con *RtcConnection, uid string, reason int) {
				t.Log("user left, " + uid)
			},
			OnUserInfoUpdated: func(con *RtcConnection, uid string, mediaInfo int, val int) {
				t.Logf("user info updated, user %s, info %d, value %d\n", uid, mediaInfo, val)
				if uid == "111" && mediaInfo == USER_MEDIA_INFO_MUTE_AUDIO {
					audioMuteState <- val
				}
			},
		},
	}
	recvCon := NewConnection(&recvCfg)
	defer recvCon.Release()
	recvCon.SetParameters("{\"rtc.video.playout_delay_max\": 250}")
	recvCon.SetParameters("{\"rtc.video.broadcaster_playout_delay_max\": 250}")
	recvCon.Connect("", "lhzuttest", "222")
	select {
	case val, _ := <-audioMuteState:
		if val != 0 {
			t.Fail()
		}
	case <-time.After(10 * time.Second):
		t.Error("wait audio mute state timeout")
		t.Fail()
	}
	*stopSend = true
	waitSenderStop.Wait()

	select {
	case val, _ := <-audioMuteState:
		if val != 1 {
			t.Fail()
		}
	case <-time.After(10 * time.Second):
		t.Error("wait audio mute state timeout")
		t.Fail()
	}

	senderCon.Disconnect()
	recvCon.Disconnect()
}
