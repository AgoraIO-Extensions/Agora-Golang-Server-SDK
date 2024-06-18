package agoraservice

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestBaseCase(t *testing.T) {
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
		w := 320
		h := 640
		dataSize := w * h * 3 / 2
		data := make([]byte, dataSize)
		sender.Start()
		for !*stopSend {
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
			time.Sleep(100 * time.Millisecond)
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
					*recvVideo = true
					waitForVideo <- struct{}{}
				}
			},
		},
	}
	recvCon := NewConnection(&recvCfg)
	defer recvCon.Release()
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
