package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"time"

	agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/rtc"
)

// AudioFileWriter is used to write the audio data to the file
type AudioFileWriter struct {
	receivedFrames int
	outputFile     *os.File
	saveToFile     bool
	totalDataSize  int64
	sampleRate     int
	channels       int
	mu             sync.Mutex

	// audio_label files
	rmsFile        *os.File
	voiceProbFile  *os.File
	musicProbFile  *os.File
	pitchFile      *os.File
	saveAudioLabel bool
}

// NewAudioFileWriter creates a new audio file writer
func NewAudioFileWriter() *AudioFileWriter {
	return &AudioFileWriter{
		receivedFrames: 0,
		saveToFile:     false,
		saveAudioLabel: false,
	}
}

// EnableSaveToFile enables to save the audio data to the file
func (writer *AudioFileWriter) EnableSaveToFile(filepath string, sampleRate, channels int) error {
	writer.mu.Lock()
	defer writer.mu.Unlock()

	writer.CloseOutputFile()

	file, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("Failed to create output file: %s, error: %v\n", filepath, err)
		return err
	}

	writer.outputFile = file
	writer.saveToFile = true
	writer.totalDataSize = 0
	writer.sampleRate = sampleRate
	writer.channels = channels

	fmt.Printf("Start saving output audio (raw PCM) to: %s\n", filepath)
	return nil
}

// EnableSaveAudioLabel enables to save the audio label data to the file
func (writer *AudioFileWriter) EnableSaveAudioLabel(basePath string, sampleRate, channels int) error {
	writer.mu.Lock()
	defer writer.mu.Unlock()

	writer.CloseAudioLabelFiles()

	// create audio_label files for each indicator
	rmsPath := basePath + "rms.pcm"
	voiceProbPath := basePath + "voice_prob.pcm"
	musicProbPath := basePath + "music_prob.pcm"
	pitchPath := basePath + "pitch.pcm"

	var err error
	writer.rmsFile, err = os.Create(rmsPath)
	if err != nil {
		fmt.Printf("Failed to create RMS file: %v\n", err)
		return err
	}

	writer.voiceProbFile, _ = os.Create(voiceProbPath)
	writer.musicProbFile, _ = os.Create(musicProbPath)
	writer.pitchFile, _ = os.Create(pitchPath)

	writer.saveAudioLabel = true

	fmt.Printf("Start saving audio label data to:\n")
	fmt.Printf("  %s\n", rmsPath)
	fmt.Printf("  %s\n", voiceProbPath)
	fmt.Printf("  %s\n", musicProbPath)
	fmt.Printf("  %s\n", pitchPath)

	return nil
}

// OnAudioFrame handles the received audio frame
func (writer *AudioFileWriter) OnAudioFrame(frame *agoraservice.AudioFrame) bool {
	writer.mu.Lock()
	defer writer.mu.Unlock()

	writer.receivedFrames++

	if writer.receivedFrames <= 10 || writer.receivedFrames%100 == 0 {
		fmt.Printf("AudioSink received the %dth frame: samples=%d, channels=%d, sampleRate=%d, bytesPerSample=%d\n",
			writer.receivedFrames, frame.SamplesPerChannel, frame.Channels, frame.SamplesPerSec, frame.BytesPerSample)
	}

	// write audio label data
	if writer.saveAudioLabel && len(frame.Buffer) > 0 {
		numSamples := frame.SamplesPerChannel * frame.Channels

		// RMS: 0-127
		rmsValue := uint8(frame.Rms)
		// voiceProb: 0 or 1
		voiceProbValue := uint8(0)
		if frame.VoiceProb > 0 {
			voiceProbValue = 0x7f
		}
		// musicProb: 0-255
		musicProbValue := uint8(frame.MusicProb)
		// pitch: int16
		pitchValue := int16(frame.Pitch)

		// write the same value to each sample point
		for i := 0; i < numSamples; i++ {
			if writer.rmsFile != nil {
				writer.rmsFile.Write([]byte{rmsValue})
			}
			if writer.voiceProbFile != nil {
				writer.voiceProbFile.Write([]byte{voiceProbValue})
			}
			if writer.musicProbFile != nil {
				writer.musicProbFile.Write([]byte{musicProbValue})
			}
			if writer.pitchFile != nil {
				// write int16 (2 bytes)
				pitchBytes := []byte{byte(pitchValue), byte(pitchValue >> 8)}
				writer.pitchFile.Write(pitchBytes)
			}
		}
	}

	// save the audio data to the file
	if writer.receivedFrames == 1 {
		fmt.Printf("[OnAudioFrame] the first frame - saveToFile=%v, outputFile=%v, bufferLen=%d\n",
			writer.saveToFile, writer.outputFile != nil, len(frame.Buffer))
	}

	if writer.saveToFile && writer.outputFile != nil && len(frame.Buffer) > 0 {
		written, err := writer.outputFile.Write(frame.Buffer)
		if err != nil {
			fmt.Printf("Failed to write to file: %v\n", err)
		} else {
			writer.totalDataSize += int64(written)
			if writer.receivedFrames <= 3 {
				fmt.Printf("[OnAudioFrame] successfully wrote %d bytes, total %d bytes\n", written, writer.totalDataSize)
			}
		}
	}

	return true
}

// CloseOutputFile closes the output file
func (writer *AudioFileWriter) CloseOutputFile() {
	if writer.outputFile != nil {
		writer.outputFile.Sync()
		writer.outputFile.Close()

		if writer.saveToFile {
			fmt.Printf("Output PCM saved completed:\n")
			fmt.Printf("  data size: %d bytes\n", writer.totalDataSize)
			fmt.Printf("  sample rate: %d Hz\n", writer.sampleRate)
			fmt.Printf("  channels: %d\n", writer.channels)
			if writer.sampleRate > 0 && writer.channels > 0 {
				duration := float64(writer.totalDataSize) / float64(writer.sampleRate*writer.channels*2)
				fmt.Printf("  duration: %.2f seconds\n", duration)
			}
		}

		writer.outputFile = nil
		writer.saveToFile = false
	}
}

// CloseAudioLabelFiles closes the audio label files
func (writer *AudioFileWriter) CloseAudioLabelFiles() {
	hadFiles := (writer.rmsFile != nil || writer.voiceProbFile != nil ||
		writer.musicProbFile != nil || writer.pitchFile != nil)

	if writer.rmsFile != nil {
		writer.rmsFile.Close()
		writer.rmsFile = nil
	}
	if writer.voiceProbFile != nil {
		writer.voiceProbFile.Close()
		writer.voiceProbFile = nil
	}
	if writer.musicProbFile != nil {
		writer.musicProbFile.Close()
		writer.musicProbFile = nil
	}
	if writer.pitchFile != nil {
		writer.pitchFile.Close()
		writer.pitchFile = nil
	}

	writer.saveAudioLabel = false

	if hadFiles {
		fmt.Println("audio label files saved completed")
	}
}

// Close closes all files
func (writer *AudioFileWriter) Close() {
	writer.mu.Lock()
	defer writer.mu.Unlock()
	writer.CloseOutputFile()
	writer.CloseAudioLabelFiles()
}

// AudioFilterProcessor is the audio filter processor (using LocalAudioProcessor wrapper)
type LocalAudioFilterProcessorWrapper struct {
	audioFileWriter     *AudioFileWriter                     // the writer for saving the file
	enableFilter        bool
	localAudioProcessor *agoraservice.LocalAudioProcessor    // wrapper that manages PCM sender, audio track and audio sink
}

// NewAudioFilterProcessor creates a new audio filter processor
func NewLocalAudioFilterProcessorWrapper(enableFilter bool) *LocalAudioFilterProcessorWrapper {
	return &LocalAudioFilterProcessorWrapper{
		audioFileWriter: NewAudioFileWriter(),
		enableFilter:    enableFilter,
	}
}

// Initialize initializes the processor using LocalAudioProcessor wrapper
func (processor *LocalAudioFilterProcessorWrapper) Initialize(appid string) error {
	fmt.Println("Start testing audio filter with LocalAudioProcessor wrapper...")

	rootDir := "/Volumes/ZR/Agora/SERVER/SDK/Agora-Golang-Server-SDK/go_sdk/examples/apm_filter_local_wrapper"

	// 1. Initialize the Agora Service
	svcCfg := agoraservice.NewAgoraServiceConfig()
	svcCfg.EnableAudioProcessor = true
	svcCfg.EnableVideo = false
	svcCfg.AppId = appid
	svcCfg.LogPath = rootDir + "/output/agora_rtc_log/agorasdk.log"
	svcCfg.ConfigDir = rootDir + "/agora_rtc_log"
	svcCfg.DataDir = rootDir + "/agora_rtc_log"

	ret := agoraservice.Initialize(svcCfg)
	if ret != 0 {
		return fmt.Errorf("failed to initialize AgoraService, error: %d", ret)
	}
	fmt.Println("✓ AgoraService initialized successfully")

	// 2. Create LocalAudioProcessor with configuration (replaces steps 2-6 from original)
	var config *agoraservice.AudioProcessorConfig
	if processor.enableFilter {
		fmt.Println("\n=== Creating LocalAudioProcessor with audio filter enabled ===")
		
		// Configure audio processing: AINS (AI Noise Suppression)
		config = &agoraservice.AudioProcessorConfig{
			ANS: &agoraservice.ANSConfig{
				Enabled: true,
			},
			EnableDump: true, // Enable dump for debugging
			// Use custom config for advanced AINS settings
			CustomConfig: `{"aec":{"split_srate_for_48k":16000},"ans":{"enabled":true},"sf_st_cfg":{"enabled":true,"ainsModelPref":10},"sf_ext_cfg":{"nsngAlgRoute":12,"nsngPredefAgg":11}}`,
		}
	} else {
		fmt.Println("\n=== Creating LocalAudioProcessor without audio filter ===")
		config = &agoraservice.AudioProcessorConfig{
			ANS: &agoraservice.ANSConfig{
				Enabled: false,
			},
		}
	}

	// Create LocalAudioProcessor (automatically handles extension, sender, track, and configuration)
	audioScenario := agoraservice.AudioScenarioChorus
	localAudioProc, err := agoraservice.NewLocalAudioProcessor(audioScenario, config)
	if err != nil {
		return fmt.Errorf("failed to create LocalAudioProcessor: %w", err)
	}
	processor.localAudioProcessor = localAudioProc
	fmt.Println("✓ LocalAudioProcessor created and configured successfully")
	
	if processor.enableFilter {
		fmt.Println("  - Extension: audio_processing_pcm_source (auto-enabled)")
		fmt.Println("  - AINS model: loaded")
		fmt.Println("  - Configuration: applied")
		fmt.Println("  - Dump: enabled")
	}

	return nil
}

// ProcessPcmFile processes the PCM file using LocalAudioProcessor
func (processor *LocalAudioFilterProcessorWrapper) ProcessPcmFile(filepath string, sampleRate, channels int) error {
	fmt.Printf("start processing PCM file: %s\n", filepath)
	fmt.Printf("PCM parameters:\n")
	fmt.Printf("  sample rate: %d Hz\n", sampleRate)
	fmt.Printf("  channels: %d\n", channels)
	fmt.Printf("  bit depth: 16 bits (int16)\n")

	// Add AudioSink using LocalAudioProcessor wrapper (returns AudioSinkContext internally managed)
	_, err := processor.localAudioProcessor.AddSink(processor.audioFileWriter.OnAudioFrame, 16000, 1)
	if err != nil {
		return fmt.Errorf("failed to add AudioSink: %w", err)
	}
	fmt.Println("✓ AudioSink added successfully")

	// open PCM file
	pcmFile, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open PCM file: %s, error: %v", filepath, err)
	}
	defer pcmFile.Close()

	// get file size
	fileInfo, _ := pcmFile.Stat()
	fmt.Printf("  file size: %d bytes\n", fileInfo.Size())

	// create output directory (same as the executable file)
	outputDir := "/Volumes/ZR/Agora/SERVER/SDK/Agora-Golang-Server-SDK/go_sdk/examples/apm_filter_local_wrapper/output"
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	fmt.Printf("✓ Output directory created: %s\n", outputDir)	

	// enable saving output audio to file
	outputPath := outputDir + "/output.pcm"
	processor.audioFileWriter.EnableSaveToFile(outputPath, sampleRate, channels)
	fmt.Printf("output file path: %s\n", outputPath)

	// enable saving audio_label data
	audioLabelBasePath := outputDir + "/audio_label/"
	err = os.MkdirAll(audioLabelBasePath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create audio label directory: %w", err)
	}
	fmt.Printf("✓ Audio label directory created: %s\n", audioLabelBasePath)
	processor.audioFileWriter.EnableSaveAudioLabel(audioLabelBasePath, sampleRate, channels)
	fmt.Println("\n✅ audio_label saving enabled:")
	fmt.Println("   - rms.pcm: uint8 (1 byte), range 0-127")
	fmt.Println("   - voice_prob.pcm: uint8 (1 byte), 0 or 1")
	fmt.Println("   - music_prob.pcm: uint8 (1 byte), range 0-255")
	fmt.Println("   - pitch.pcm: int16 (2 bytes), original value")

	// save input original data (for comparison)
	inputCopyPath := outputDir + "/input_copy.pcm"
	inputCopyFile, err := os.Create(inputCopyPath)
	if err != nil {
		fmt.Printf("failed to create input copy file: %s, error: %v", inputCopyPath, err)
	} else {
		defer inputCopyFile.Close()
		fmt.Printf("save input copy (raw PCM) to: %s\n", inputCopyPath)
	}

	// calculate the number of samples per frame (10ms)
	samplesPerChannel := sampleRate / 100
	fmt.Printf("  samples per frame: %d\n", samplesPerChannel)

	bytesPerFrame := samplesPerChannel * channels * 2 // 16bit = 2 bytes
	buffer := make([]byte, bytesPerFrame)

	frameCount := 0
	successCount := 0
	failCount := 0
	totalInputDataSize := int64(0)

	for {
		// read one frame data
		n, err := pcmFile.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("failed to read file: %v\n", err)
			break
		}

		if n == 0 {
			break
		}

		// save input data to file
		if inputCopyFile != nil {
			inputCopyFile.Write(buffer[:n])
			totalInputDataSize += int64(n)
		}

		// Send audio data using LocalAudioProcessor wrapper (much simpler!)
		err = processor.localAudioProcessor.SendPCM(buffer[:n], sampleRate, channels)

		// print every 1000 frames
		if (frameCount+1)%1000 == 0 {
			fmt.Printf("prepare to send the %dth frame...\n", frameCount+1)
			if err != nil {
				fmt.Printf("the %dth frame sent failed: %v\n", frameCount+1, err)
			} else {
				fmt.Printf("the %dth frame sent successfully\n", frameCount+1)
			}
		}

		if err == nil {
			successCount++
			frameCount++

			// print every 100 frames
			if frameCount%100 == 0 {
				fmt.Printf("sent %d frames...\n", frameCount)
			}
		} else {
			failCount++
			frameCount++
			fmt.Printf("failed to send the %dth frame, error: %v\n", frameCount, err)
			return fmt.Errorf("failed to send the %dth frame: %w", frameCount, err)
		}

		// wait 10ms
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("audio data sent completed!")
	fmt.Printf("  total frames: %d\n", frameCount)
	fmt.Printf("  success: %d frames\n", successCount)
	fmt.Printf("  failed: %d frames\n", failCount)

	if inputCopyFile != nil {
		inputCopyFile.Sync()
		fmt.Println("\ninput copy PCM saved completed:")
		fmt.Printf("  data size: %d bytes\n", totalInputDataSize)
		fmt.Printf("  sample rate: %d Hz\n", sampleRate)
		fmt.Printf("  channels: %d\n", channels)
		if sampleRate > 0 && channels > 0 {
			duration := float64(totalInputDataSize) / float64(sampleRate*channels*2)
			fmt.Printf("  duration: %.2f seconds\n", duration)
		}
	}

	// wait for a while to ensure all data is processed
	time.Sleep(500 * time.Millisecond)

	processor.audioFileWriter.Close()

	fmt.Println("\n✅ audio_label file saved completed")
	fmt.Println("   - rms.pcm: uint8 (1 byte), 0-127")
	fmt.Println("   - voice_prob.pcm: uint8 (1 byte), 0 or 1")
	fmt.Println("   - music_prob.pcm: uint8 (1 byte), 0-255")
	fmt.Println("   - pitch.pcm: int16 (2 bytes)")

	return nil
}

// Cleanup clean up resources using LocalAudioProcessor wrapper
func (processor *LocalAudioFilterProcessorWrapper) Cleanup() {
	if processor.audioFileWriter != nil {
		processor.audioFileWriter.Close()
	}

	// Release LocalAudioProcessor (automatically releases AudioSink, PCM Sender and Audio Track)
	if processor.localAudioProcessor != nil {
		processor.localAudioProcessor.Release()
		processor.localAudioProcessor = nil
		fmt.Println("✓ LocalAudioProcessor released (including AudioSink, PCM Sender and Audio Track)")
	}

	agoraservice.Release()
	fmt.Println("✓ AgoraService released")
}

// global processor variable, used to clean up resources when signal processing
var globalProcessor *LocalAudioFilterProcessorWrapper

func main() {
	bStop := new(bool)
	*bStop = false

	// catch terminal signal - first clean up resources then exit
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		fmt.Println("\nreceived termination signal, cleaning up...")
		*bStop = true

		// immediately close file, ensure data is written to disk
		if globalProcessor != nil && globalProcessor.audioFileWriter != nil {
			fmt.Println("flush and close output file...")
			globalProcessor.audioFileWriter.Close()
		}

		// give the program 1 second to complete other cleanup, then force exit
		time.Sleep(1 * time.Second)
		fmt.Println("force exit")
		os.Exit(0)
	}()

	// get parameters: appid, pcm_file_path
	argus := os.Args
	if len(argus) < 2 {
		fmt.Println("usage: extension <appid> [pcm_file_path]")
		fmt.Println("  appid: Agora App ID")
		fmt.Println("  pcm_file_path: (optional) PCM file path, default is noise.pcm")
		return
	}

	appid := argus[1]

	// PCM file path (default use absolute path)
	pcmFilePath := "/Volumes/ZR/Agora/SERVER/SDK/Agora-Golang-Server-SDK/go_sdk/examples/apm_filter_local_wrapper/noise.pcm"
	if len(argus) >= 3 {
		pcmFilePath = argus[2]
	}

	if appid == "" {
		fmt.Println("please set AGORA_APP_ID parameter")
		return
	}

	fmt.Println("=== Agora audio filter test ===")

	// create processor (enable filter)
	enableFilter := true
	processor := NewLocalAudioFilterProcessorWrapper(enableFilter)
	globalProcessor = processor // save to global variable, used for signal processing

	// initialize
	err := processor.Initialize(appid)
	if err != nil {
		fmt.Printf("failed to initialize: %v\n", err)
		return
	}
	defer processor.Cleanup()

	fmt.Printf("input file path: %s\n", pcmFilePath)

	// PCM parameters
	sampleRate := 48000
	channels := 1

	fmt.Printf("PCM parameters: %d Hz, %d channels\n", sampleRate, channels)

	// process PCM file
	err = processor.ProcessPcmFile(pcmFilePath, sampleRate, channels)
	if err != nil {
		fmt.Printf("failed to process PCM file: %v\n", err)
		return
	}

	// wait for a while
	time.Sleep(1 * time.Second)

	fmt.Println("=== test completed ===")
}

