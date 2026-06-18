package agoraservice


import (
	"fmt"
	"sync"
	"os"
)


// transcoding worker
type AudioTransCodeInData struct {
	data       []byte
	codecType  int
	sampleRate int
	channels   int
}


type AudioTranscodingWorker struct {
	isRunning        bool
	encodeDataQueue  chan *AudioTransCodeInData // for en
	stopChan         chan bool
	conn             *RtcConnection
	mutex            sync.Mutex
	wg               sync.WaitGroup
	stopOnce         sync.Once
	decoder          *OpusDecoder
	out_channels     int
	out_sampleRate   int
}

func NewAudioTranscodingWorker(conn *RtcConnection, in_channels int) *AudioTranscodingWorker {
	
	intervalInMs := 20

	// for this case: just fix out_channels and out_sampleRate to 1 and 16000
	ret := &AudioTranscodingWorker{
		isRunning:        false,
		encodeDataQueue:  make(chan *AudioTransCodeInData, 1000),
		conn:             conn,
		stopChan:         make(chan bool),
		mutex:            sync.Mutex{},
		wg:               sync.WaitGroup{},
		stopOnce:         sync.Once{},
		decoder:          nil,
		out_channels:     1,
		out_sampleRate:   16000,
	}
	out_channels := ret.out_channels
	out_sampleRate := ret.out_sampleRate
	audioDecoder := NewOpusDecoder(in_channels, out_sampleRate, out_channels)
	if audioDecoder == nil {
		fmt.Printf("NewAudioTranscodingWorker failed to create opus decoder\n", )
	}
	ret.decoder = audioDecoder
	fmt.Printf("NewAudioTranscodingWorker created with in_channels: %d, intervalInMs: %d\n", in_channels, intervalInMs)
	return ret
}
func (worker *AudioTranscodingWorker) Start() int {
	worker.mutex.Lock()
	defer worker.mutex.Unlock()
	if worker.isRunning {
		return -1000
	}
	//check has 
	if worker.decoder == nil {
		return -10001
	}
	worker.isRunning = true
	// increase WaitGroup count before starting goroutine
	worker.wg.Add(1)
	// start a goroutine to run the worker
	go worker.run()
	return 0
}
func (worker *AudioTranscodingWorker) run() {
	// ensure Done() is called when the function returns
	defer worker.wg.Done()

	// get out_channels and out_sampleRate from worker
	out_channels := worker.out_channels
	out_sampleRate := worker.out_sampleRate

	// for debug
	debug := false
	var debugFile *os.File = nil

	if debug {
		debugPath := "./opus_transcode_worker_debug.pcm"
		debugFile, _ = os.OpenFile(debugPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if debugFile == nil {
			debug = false
		}
		defer debugFile.Close()
	}
	for {
		select {
		case inData := <-worker.encodeDataQueue:
			//should decoding all the data in the queue, Never never drop any data, even if the decoder is busy.
			rawAudioData, _ := worker.decoder.Decode(inData.data)
			if rawAudioData != nil  {
				// for debug
				if debug {
					_, err := debugFile.Write(rawAudioData)
					if err != nil {
						fmt.Printf("Failed to write debug file: %v\n", err)
					}
				}
				// and call conn.PushAudioPcmData to push the decoded audio data to the RTC
				worker.conn.PushAudioPcmData(rawAudioData, out_sampleRate, out_channels, 0)
			} 
			
		case <-worker.stopChan:
			// received stop signal, exit loop
			worker.mutex.Lock()
			worker.isRunning = false
			worker.mutex.Unlock()
			return // use return instead of break, ensure exit function
		}
	}
}
func (worker *AudioTranscodingWorker) PushEncodedData(data []byte, channels int, codecType int) {
	if worker.isRunning == false {
		return
	}
	worker.encodeDataQueue <- &AudioTransCodeInData{
		data:      data,
		channels:  channels,
		codecType: codecType,
		sampleRate: 16000, // no need for opus
	}
}
func (worker *AudioTranscodingWorker) Stop() {
	worker.mutex.Lock()
	if !worker.isRunning {
		worker.mutex.Unlock()
		return // if already stopped, return directly
	}
	worker.mutex.Unlock()

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
		worker.decoder.Release()
		worker.decoder = nil
	}
	worker.mutex.Unlock()
}
