package agoraservice

import (
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
	waitSenderStop.Add(2)
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

	waitSenderJoin := make(chan struct{}, 1)
	waitForAudio := make(chan struct{}, 1)
	waitForData := make(chan struct{}, 1)
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
				waitSenderJoin <- struct{}{}
			},
			OnUserLeft: func(con *RtcConnection, uid string, reason int) {
				t.Log("user left, " + uid)
			},
			OnStreamMessage: func(con *RtcConnection, uid string, streamId int, data []byte) {
				t.Log("stream message")
				waitForAudio <- struct{}{}
			},
			OnStreamMessageError: func(con *RtcConnection, uid string, streamId int, errCode int, missed int, cached int) {
				t.Log("stream message error")
			},
		},
		AudioFrameObserver: &RtcConnectionAudioFrameObserver{
			OnPlaybackAudioFrameBeforeMixing: func(con *RtcConnection, channelId string, userId string, frame *PcmAudioFrame) {
				t.Log("Playback audio frame before mixing")
				waitForData <- struct{}{}
			},
		},
	}
	recvCon := NewConnection(&recvCfg)
	defer recvCon.Release()
	recvCon.Connect("", "lhzuttest", "222")
	timer := time.NewTimer(10 * time.Second)
	recvAudio := false
	recvData := false
	select {
	case <-waitSenderJoin:
	case <-waitForAudio:
		recvAudio = true
	case <-waitForData:
		recvData = true
	case <-timer.C:
		t.Error("wait audio or data timeout, recvAudio: ", recvAudio, ", recvData: ", recvData)
		t.Fail()
	}
	*stopSend = true
	waitSenderStop.Wait()
	senderCon.Disconnect()
	recvCon.Disconnect()
}
