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

// AudioFilterProcessor is the audio filter processor
type LocalAudioFilterProcessor struct {
	audioFileWriter  *AudioFileWriter               // the writer for saving the file
	audioSinkContext *agoraservice.AudioSinkContext // the sink for receiving the audio from the C layer
	enableFilter     bool
	pcmSender        *agoraservice.AudioPcmDataSender
	audioTrack       *agoraservice.LocalAudioTrack
}

// NewAudioFilterProcessor creates a new audio filter processor
func NewLocalAudioFilterProcessor(enableFilter bool) *LocalAudioFilterProcessor {
	return &LocalAudioFilterProcessor{
		audioFileWriter: NewAudioFileWriter(),
		enableFilter:    enableFilter,
	}
}

// Initialize initializes the processor
func (processor *LocalAudioFilterProcessor) Initialize(appid string) error {
	fmt.Println("Start testing audio filter...")

	rootDir := "/Volumes/ZR/Agora/SERVER/SDK/Agora-Golang-Server-SDK/go_sdk/examples/apm_filter_local"

	// 1. initialize the Agora Service
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
	fmt.Println("AgoraService initialized successfully")

	if processor.enableFilter {
		// 2. enable audio_processing_pcm extension
		ret = agoraservice.EnableExtension("agora.builtin", "audio_processing_pcm", "", true)
		if ret != 0 {
			fmt.Printf("failed to enable audio_processing_pcm extension, error: %d\n", ret)
		} else {
			fmt.Println("successfully enabled audio_processing_pcm extension")
		}
	}

	// 3. create PCM Sender
	processor.pcmSender = agoraservice.NewAudioPcmDataSender()
	if processor.pcmSender == nil {
		return fmt.Errorf("failed to create AudioPcmDataSender")
	}
	fmt.Println("AudioPcmDataSender created successfully")

	// 4. create audio track
	audioScenario := agoraservice.AudioScenarioChorus
	processor.audioTrack = agoraservice.NewCustomAudioTrackPcm(processor.pcmSender, audioScenario)
	if processor.audioTrack == nil {
		return fmt.Errorf("failed to create LocalAudioTrack")
	}
	fmt.Println("LocalAudioTrack created successfully")

	// 5. configure audio filter (configure before SetEnabled)
	if processor.enableFilter {
		fmt.Println("\n=== configure audio filter (noise reduction) ===")

		// step 2: enable filter on the track
		fmt.Println("step 2: enable filter on the AudioTrack...")
		ret = processor.audioTrack.EnableAudioFilter("audio_processing_pcm", true, 3)
		if ret != 0 {
			fmt.Printf("failed to enable filter on the AudioTrack, error: %d\n", ret)
		} else {
			fmt.Println("successfully enabled filter on the AudioTrack")
		}

		// step 3: load AINS resource
		fmt.Println("step 3: load AINS model resource...")
		ret = processor.audioTrack.SetFilterProperty("audio_processing_pcm", "apm_load_resource", "ains", 3)
		if ret != 0 {
			fmt.Printf("failed to load AINS model resource, error: %d\n", ret)
		} else {
			fmt.Println("successfully loaded AINS model resource")
		}

		// step 4: configure noise reduction parameters (enable AINS AI noise reduction)
		fmt.Println("step 4: configure noise reduction parameters (enable AINS AI noise reduction)...")
		config := "{\"aec\":{\"split_srate_for_48k\":16000},\"ans\":{\"enabled\":true},\"sf_st_cfg\":{\"enabled\":true,\"ainsModelPref\":10},\"sf_ext_cfg\":{\"nsngAlgRoute\":12,\"nsngPredefAgg\":11}}"
		ret = processor.audioTrack.SetFilterProperty("audio_processing_pcm", "apm_config", config, 3)
		if ret != 0 {
			fmt.Printf("failed to configure noise reduction parameters, error: %d\n", ret)
		} else {
			fmt.Println("successfully configured noise reduction parameters")
		}

		// step 5: enable dump (for debugging)
		fmt.Println("step 5: enable apm_dump...")
		ret = processor.audioTrack.SetFilterProperty("audio_processing_pcm", "apm_dump", "true", 3)
		if ret != 0 {
			fmt.Printf("failed to enable apm_dump, error: %d\n", ret)
		} else {
			fmt.Println("successfully enabled apm_dump")
		}

		fmt.Println("=== audio filter configured successfully ===")
	} else {
		fmt.Println("disable audio filter (test original audio)")
	}

	// 6. enable audio track
	processor.audioTrack.SetEnabled(true)

	// 7. create AudioSinkContext for receiving the processed audio
	processor.audioSinkContext = agoraservice.NewAudioSink(func(frame *agoraservice.AudioFrame) bool {
		// forward the received audio data to audioSink for saving
		if processor.audioFileWriter == nil {
			fmt.Println("[Lambda callback] processor.audioFileWriter is nil!")
			return false
		}
		result := processor.audioFileWriter.OnAudioFrame(frame)
		return result
	})
	if processor.audioSinkContext == nil {
		return fmt.Errorf("failed to create AudioSinkContext")
	}
	fmt.Println("AudioSinkContext created successfully")

	return nil
}

// ProcessPcmFile processes the PCM file
func (processor *LocalAudioFilterProcessor) ProcessPcmFile(filepath string, sampleRate, channels int) error {
	fmt.Printf("start processing PCM file: %s\n", filepath)
	fmt.Printf("PCM parameters:\n")
	fmt.Printf("  sample rate: %d Hz\n", sampleRate)
	fmt.Printf("  channels: %d\n", channels)
	fmt.Printf("  bit depth: 16 bits (int16)\n")

	// add AudioSinkContext to audioTrack
	ret := processor.audioTrack.AddAudioSink(processor.audioSinkContext.GetHandle(), 16000, 1)
	if ret != 0 {
		fmt.Printf("failed to add AudioSinkContext to audioTrack, error: %d\n", ret)
	} else {
		fmt.Println("successfully added AudioSinkContext to audioTrack")
	}

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
	outputDir := "/Volumes/ZR/Agora/SERVER/SDK/Agora-Golang-Server-SDK/go_sdk/examples/apm_filter_local/output"
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

		// create audio frame
		audioFrame := &agoraservice.AudioFrame{
			Buffer:            buffer[:n],
			SamplesPerChannel: samplesPerChannel,
			BytesPerSample:    2, // 16bit = 2 bytes
			Channels:          1,
			SamplesPerSec:     sampleRate,
			RenderTimeMs:      0,
			PresentTimeMs:     0,
		}

		// send audio data through PCM Sender
		ret := processor.pcmSender.SendAudioPcmData(audioFrame)

		// print every 1000 frames
		if (frameCount+1)%1000 == 0 {
			fmt.Printf("prepare to send the %dth frame...\n", frameCount+1)
			fmt.Printf("the %dth frame sent return: %d\n", frameCount+1, ret)
		}

		if ret == 0 {
			successCount++
			frameCount++

			// print every 100 frames
			if frameCount%100 == 0 {
				fmt.Printf("sent %d frames...\n", frameCount)
			}
		} else {
			failCount++
			frameCount++
			fmt.Printf("failed to send the %dth frame, error: %d\n", frameCount, ret)
			return fmt.Errorf("failed to send the %dth frame, error: %d", frameCount, ret)
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

// Cleanup clean up resources
func (processor *LocalAudioFilterProcessor) Cleanup() {
	if processor.audioFileWriter != nil {
		processor.audioFileWriter.Close()
	}

	// release AudioSinkContext
	if processor.audioSinkContext != nil {
		processor.audioSinkContext.Release()
		processor.audioSinkContext = nil
		fmt.Println("AudioSinkContext released")
	}

	// release Audio Track
	if processor.audioTrack != nil {
		processor.audioTrack.Release()
		processor.audioTrack = nil
		fmt.Println("LocalAudioTrack released")
	}

	// release PCM Sender
	if processor.pcmSender != nil {
		processor.pcmSender.Release()
		processor.pcmSender = nil
		fmt.Println("AudioPcmDataSender released")
	}

	agoraservice.Release()
	fmt.Println("AgoraService released")
}

// global processor variable, used to clean up resources when signal processing
var globalProcessor *LocalAudioFilterProcessor

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
	pcmFilePath := "/Volumes/ZR/Agora/SERVER/SDK/Agora-Golang-Server-SDK/go_sdk/examples/apm_filter_local/noise.pcm"
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
	processor := NewLocalAudioFilterProcessor(enableFilter)
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
