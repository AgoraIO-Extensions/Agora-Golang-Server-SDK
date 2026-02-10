//go:build !avcodec
// +build !avcodec
package agoraservice


import (
	"fmt"
	"sync"
	"time"
)

// global state for managing decoder
type VideoDecoder struct {
	decoder    *any
	frameCount int
	mutex      sync.Mutex
}

// create video decoder
func NewVideoDecoder() (*VideoDecoder, error) {
	
	return nil, fmt.Errorf("0202:200318_v1: failed to initialize H264 decoder: %d", -10001)
}

// decode and save one frame
func (vd *VideoDecoder) DecodeAndSave(data []byte, width, height, codecType int) (*ExternalVideoFrame, int) {
	
	return nil, -10002
}

// close decoder
func (vd *VideoDecoder) Close() {
	
}

// transcoding worker
type TransCodeInData struct {
	data      []byte
	frameInfo *EncodedVideoFrameInfo
	// need to keep timestamp for calc processing time??
	timestamp int64
}
type TranscodingWorker struct {
	isRunning        bool
	encodeDataQueue  chan *TransCodeInData // for en
	decodeVideoFrame *ExternalVideoFrame   // for decode video frame, only one frame is allowed to be decoded at a time.
	timer            *time.Ticker          // for timer
	conn             *RtcConnection
	decoder          *VideoDecoder
	stopChan         chan bool
	mutex            sync.Mutex
	wg               sync.WaitGroup // for waiting run() goroutine to complete
	stopOnce         sync.Once      // ensure stopChan is only closed once
}

func NewTranscodingWorker(conn *RtcConnection, fps int) *TranscodingWorker {
	fmt.Printf("0202:200318_v1: NewTranscodingWorker, not implemented\n")
	return nil
}
func (worker *TranscodingWorker) Start() int {
	fmt.Printf("0202:200318_v1: TranscodingWorker Start, not implemented\n")
	return -10003
}
func (worker *TranscodingWorker) run() {
	fmt.Printf("0202:200318_v1: TranscodingWorker run, not implemented\n")
	return
}
func (worker *TranscodingWorker) PushEncodedData(data []byte, frameInfo *EncodedVideoFrameInfo) {
	fmt.Printf("0202:200318_v1: TranscodingWorker PushEncodedData, not implemented\n")
}
func (worker *TranscodingWorker) Stop() {
	fmt.Printf("0202:200318_v1: TranscodingWorker Stop, not implemented\n")
	return
}
