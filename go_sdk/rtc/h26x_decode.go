package agoraservice

// #cgo pkg-config: libavcodec libavutil
// #include "h26x_decode_frame.h"
import "C"
import (
	"fmt"
	"sync"
	"time"
	"unsafe"
)

// global state for managing decoder
type VideoDecoder struct {
	decoder    *C.H264Decoder
	frameCount int
	mutex      sync.Mutex
}

// create video decoder
func NewVideoDecoder() (*VideoDecoder, error) {
	decoder := C.h264_decoder_init()
	if decoder == nil {
		return nil, fmt.Errorf("failed to initialize H264 decoder")
	}

	return &VideoDecoder{
		decoder:    decoder,
		frameCount: 0,
	}, nil
}

// decode and save one frame
func (vd *VideoDecoder) DecodeAndSave(data []byte, width, height, codecType int) (*ExternalVideoFrame, int) {
	vd.mutex.Lock()
	defer vd.mutex.Unlock()

	var yuvFrame *C.YUVFrame

	var outframe *ExternalVideoFrame = nil
	var ret int = 0

	// use h264_decode_stream to automatically handle NAL units
	ret = int(C.h264_decode_stream(
		vd.decoder,
		(*C.uint8_t)(unsafe.Pointer(&data[0])),
		C.int(len(data)),
		C.int(width),
		C.int(height),
		&yuvFrame,
	))

	if ret == 0 && yuvFrame != nil {
		// successfully decoded one frame, write to YUV file
		defer C.yuv_frame_free(yuvFrame)

		// write Y plane
		yData := unsafe.Slice((*byte)(unsafe.Pointer(yuvFrame.y_data)), yuvFrame.y_size)

		// write U plane
		uData := unsafe.Slice((*byte)(unsafe.Pointer(yuvFrame.u_data)), yuvFrame.uv_size)

		// write V plane
		vData := unsafe.Slice((*byte)(unsafe.Pointer(yuvFrame.v_data)), yuvFrame.uv_size)

		// arrange data to outframe.Buffer
		outframe = &ExternalVideoFrame{
			Type:      VideoBufferRawData,
			Format:    VideoPixelI420,
			Buffer:    nil,
			Stride:    width,
			Height:    height,
			Timestamp: 0,
		}
		bufsize := yuvFrame.width * yuvFrame.height * 3 / 2
		uvlen := yuvFrame.width * yuvFrame.height / 4
		outframe.Buffer = make([]byte, bufsize)
		copy(outframe.Buffer, yData)
		copy(outframe.Buffer[yuvFrame.width*yuvFrame.height:], uData)
		copy(outframe.Buffer[yuvFrame.width*yuvFrame.height+uvlen:], vData)
		outframe.Stride = int(yuvFrame.width)
		outframe.Height = int(yuvFrame.height)
		outframe.Timestamp = 0

		vd.frameCount++
		fmt.Printf("Decoded frame %d: %dx%d, YUV size: %d bytes\n",
			vd.frameCount, yuvFrame.width, yuvFrame.height, yuvFrame.total_size)
	}

	return outframe, ret
}

// close decoder
func (vd *VideoDecoder) Close() {
	vd.mutex.Lock()
	defer vd.mutex.Unlock()

	if vd.decoder != nil {
		C.h264_decoder_free(vd.decoder)
		vd.decoder = nil
	}
	fmt.Printf("Total decoded frames: %d\n", vd.frameCount)
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
	if fps <= 0 {
		fps = 15
	}
	if fps > 30 {
		fps = 30
	}
	intervalInMs := 1000 / fps
	ret := &TranscodingWorker{
		isRunning:        false,
		encodeDataQueue:  make(chan *TransCodeInData, 1000),
		decodeVideoFrame: nil,
		timer:            time.NewTicker(time.Duration(intervalInMs) * time.Millisecond),
		conn:             conn,
		decoder:          nil,
		stopChan:         make(chan bool),
		mutex:            sync.Mutex{},
		wg:               sync.WaitGroup{},
		stopOnce:         sync.Once{},
	}
	decoder, err := NewVideoDecoder()
	if err != nil {
		fmt.Printf("NewTranscodingWorker failed to create video decoder: %v\n", err)
	}
	ret.decoder = decoder
	fmt.Printf("NewTranscodingWorker created with fps: %d, intervalInMs: %d\n", fps, intervalInMs)
	return ret
}
func (worker *TranscodingWorker) Start() int {
	worker.mutex.Lock()
	defer worker.mutex.Unlock()
	if worker.isRunning {
		return -1000
	}
	worker.isRunning = true
	// increase WaitGroup count before starting goroutine
	worker.wg.Add(1)
	// start a goroutine to run the worker
	go worker.run()
	return 0
}
func (worker *TranscodingWorker) run() {
	// ensure Done() is called when the function returns
	defer worker.wg.Done()

	for {
		select {
		case data := <-worker.encodeDataQueue:
			//should decoding all the data in the queue, Never never drop any data, even if the decoder is busy.
			yuvframe, ret := worker.decoder.DecodeAndSave(data.data, data.frameInfo.Width, data.frameInfo.Height, int(data.frameInfo.CodecType))
			if yuvframe != nil {
				//fmt.Printf("TranscodingWorker decode video frame failed: %d\n", ret)
				// replace and assign to decodeVideoFrame
				worker.mutex.Lock()
				worker.decodeVideoFrame = yuvframe
				worker.mutex.Unlock()
			} else if ret != 0 {
				fmt.Printf("TranscodingWorker decode video frame failed: %d\n", ret)
			}
		case <-worker.timer.C:
			// do processing decoded video frame now
			worker.mutex.Lock()
			worker.conn.handleDecodedVideoFrameForTranscode(worker.decodeVideoFrame)
			worker.mutex.Unlock()
		case <-worker.stopChan:
			// received stop signal, exit loop
			worker.mutex.Lock()
			worker.isRunning = false
			worker.mutex.Unlock()
			return // use return instead of break, ensure exit function
		}
	}
}
func (worker *TranscodingWorker) PushEncodedData(data []byte, frameInfo *EncodedVideoFrameInfo) {
	worker.encodeDataQueue <- &TransCodeInData{
		data: data,
		frameInfo: frameInfo,
		timestamp: time.Now().UnixMilli(),
	}
}
func (worker *TranscodingWorker) Stop() {
	worker.mutex.Lock()
	if !worker.isRunning {
		worker.mutex.Unlock()
		return // if already stopped, return directly
	}
	worker.mutex.Unlock()

	// stop timer
	worker.timer.Stop()

	// use sync.Once to ensure stopChan is only closed once (prevent repeated calls to stop)
	worker.stopOnce.Do(func() {
		close(worker.stopChan)
	})

	// wait for run() goroutine to complete
	worker.wg.Wait()

	// clean up resources
	worker.mutex.Lock()
	worker.isRunning = false
	if worker.decoder != nil {
		worker.decoder.Close()
		worker.decoder = nil
	}
	worker.mutex.Unlock()
}
