package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/rtc"
)

// RemoteAudioFilterProcessorWrapper manages remote audio filtering (using RemoteAudioProcessor wrapper)
type RemoteAudioFilterProcessorWrapper struct {
	// Sender connection (to publish audio)
	connSend          *agoraservice.RtcConnection
	localUserSend     *agoraservice.LocalUser
	pcmSender         *agoraservice.AudioPcmDataSender
	localAudioTrack   *agoraservice.LocalAudioTrack
	
	// Receiver connection (to subscribe and process audio)
	connRecv             *agoraservice.RtcConnection
	localUserRecv        *agoraservice.LocalUser
	remoteAudioTrack     *agoraservice.RemoteAudioTrack
	remoteAudioProcessor *agoraservice.RemoteAudioProcessor  // wrapper that simplifies remote audio processing
	
	channelId         string
	senderUid         string
	receiverUid       string
	
	isRunning         bool
	mu                sync.Mutex
}

// NewRemoteAudioFilterProcessor creates a new remote audio filter processor
func NewRemoteAudioFilterProcessorWrapper() *RemoteAudioFilterProcessorWrapper {
	return &RemoteAudioFilterProcessorWrapper{
		isRunning: true,
	}
}

// Initialize initializes the RTC connections (sender and receiver)
func (p *RemoteAudioFilterProcessorWrapper) Initialize(appid, channelId, senderUid, receiverUid string) error {
	p.channelId = channelId
	p.senderUid = senderUid
	p.receiverUid = receiverUid

	rootDir := "/Volumes/ZR/Agora/SERVER/SDK/Agora-Golang-Server-SDK/go_sdk/examples/apm_filter_remote_wrapper"

	fmt.Println("\n=== Initialize AgoraService ===")

	// Create service configuration
	svcCfg := agoraservice.NewAgoraServiceConfig()
	svcCfg.AppId = appid
	svcCfg.EnableAudioProcessor = true
	svcCfg.EnableAudioDevice = false
	svcCfg.LogPath = rootDir + "/output/agora_rtc_log/agorasdk.log"
	svcCfg.ConfigDir = rootDir + "/agora_rtc_log"
	svcCfg.DataDir = rootDir + "/agora_rtc_log"

	// Initialize service
	ret := agoraservice.Initialize(svcCfg)
	if ret != 0 {
		return fmt.Errorf("failed to initialize AgoraService, error code: %d", ret)
	}
	fmt.Println("✓ AgoraService initialized")

	// Enable remote playback extension
	ret = agoraservice.EnableExtension("agora.builtin", "audio_processing_remote_playback", "", true)
	if ret != 0 {
		fmt.Printf("⚠️  Failed to enable audio_processing_remote_playback extension, error: %d\n", ret)
	} else {
		fmt.Println("✓ Audio processing extension enabled")
	}

	// Create audio PCM sender and track for sender
	p.pcmSender = agoraservice.NewAudioPcmDataSender()
	if p.pcmSender == nil {
		return fmt.Errorf("failed to create AudioPcmDataSender")
	}

	p.localAudioTrack = agoraservice.NewCustomAudioTrackPcm(p.pcmSender, 9) // AUDIO_SCENARIO_GAME_STREAMING
	if p.localAudioTrack == nil {
		return fmt.Errorf("failed to create custom audio track")
	}

	p.localAudioTrack.SetEnabled(true)
	fmt.Println("✓ Audio track created")

	// Create sender connection
	if err := p.createSenderConnection(appid); err != nil {
		return err
	}

	// Wait for sender connection to stabilize
	fmt.Println("\n=== Waiting for sender to stabilize ===")
	time.Sleep(3 * time.Second)

	// Create receiver connection
	if err := p.createReceiverConnection(); err != nil {
		return err
	}

	// Wait for receiver connection
	fmt.Println("\n=== Waiting for receiver connection ===")
	time.Sleep(3 * time.Second)

	return nil
}

// createSenderConnection creates the sender connection
func (p *RemoteAudioFilterProcessorWrapper) createSenderConnection(appid string) error {
	fmt.Println("\n=== Creating Sender Connection ===")

	// Create sender connection configuration
	senderCfg := &agoraservice.RtcConnectionConfig{
		ClientRole:         agoraservice.ClientRoleBroadcaster,
		ChannelProfile:     agoraservice.ChannelProfileLiveBroadcasting,
		AutoSubscribeAudio: false,
		AutoSubscribeVideo: false,
		AudioRecvMediaPacket: false,
	}

	publishCfg := agoraservice.NewRtcConPublishConfig()
	publishCfg.AudioScenario = 9 // AUDIO_SCENARIO_GAME_STREAMING
	publishCfg.IsPublishAudio = false  // 不让 connection 自动创建 track
	publishCfg.AudioPublishType = 0    // AudioPublishTypeNoPublish

	p.connSend = agoraservice.NewRtcConnection(senderCfg, publishCfg)
	if p.connSend == nil {
		return fmt.Errorf("failed to create sender connection")
	}

	p.localUserSend = p.connSend.GetLocalUser()
	if p.localUserSend == nil {
		return fmt.Errorf("failed to get sender local user")
	}

	// Register sender observer to monitor publish success
	senderObserver := &agoraservice.LocalUserObserver{
		OnAudioTrackPublishSuccess: func(localUser *agoraservice.LocalUser, track *agoraservice.LocalAudioTrack) {
			fmt.Println("✅ [Callback] Audio track published successfully")
		},
	}
	p.connSend.RegisterLocalUserObserver(senderObserver)

	// Connect to channel (first parameter is token, use "" for no token)
	fmt.Printf("Connecting sender to channel: %s, uid: %s\n", p.channelId, p.senderUid)
	ret := p.connSend.Connect("", p.channelId, p.senderUid)
	if ret != 0 {
		return fmt.Errorf("failed to connect sender, error code: %d", ret)
	}
	fmt.Println("✓ Sender Connect() called successfully")

	time.Sleep(3 * time.Second)

	// 手动发布我们创建的 audio track
	fmt.Println("Publishing manually created audio track...")
	ret = p.localUserSend.PublishAudioTrack(p.localAudioTrack)
	if ret != 0 {
		return fmt.Errorf("failed to publish audio track, error code: %d", ret)
	}
	fmt.Println("✓ Audio track published successfully")

	time.Sleep(2 * time.Second)

	return nil
}

// createReceiverConnection creates the receiver connection
func (p *RemoteAudioFilterProcessorWrapper) createReceiverConnection() error {
	fmt.Println("\n=== Creating Receiver Connection ===")

	// Create receiver connection configuration
	receiverCfg := &agoraservice.RtcConnectionConfig{
		ClientRole:         agoraservice.ClientRoleAudience,
		ChannelProfile:     agoraservice.ChannelProfileLiveBroadcasting,
		AutoSubscribeAudio: true, // Auto subscribe
		AutoSubscribeVideo: false,
		AudioRecvMediaPacket: false,
	}

	publishCfg := agoraservice.NewRtcConPublishConfig()

	p.connRecv = agoraservice.NewRtcConnection(receiverCfg, publishCfg)
	if p.connRecv == nil {
		return fmt.Errorf("failed to create receiver connection")
	}

	p.localUserRecv = p.connRecv.GetLocalUser()
	if p.localUserRecv == nil {
		return fmt.Errorf("failed to get receiver local user")
	}

	// Set up observer for receiver
	p.setupReceiverObserver()

	// Set audio frame parameters
	ret := p.localUserRecv.SetPlaybackAudioFrameBeforeMixingParameters(1, 48000)
	if ret != 0 {
		fmt.Printf("⚠️  Failed to set audio frame parameters, error: %d\n", ret)
	}

	// Connect to channel (first parameter is token, use "" for no token)
	fmt.Printf("Connecting receiver to channel: %s, uid: %s\n", p.channelId, p.receiverUid)
	ret = p.connRecv.Connect("", p.channelId, p.receiverUid)
	if ret != 0 {
		return fmt.Errorf("failed to connect receiver, error code: %d", ret)
	}

	fmt.Println("✓ Receiver connection created")
	return nil
}

// setupReceiverObserver sets up the receiver local user observer
func (p *RemoteAudioFilterProcessorWrapper) setupReceiverObserver() {
	observer := &agoraservice.LocalUserObserver{
		OnUserAudioTrackSubscribed: func(localUser *agoraservice.LocalUser, uid string, remoteAudioTrack *agoraservice.RemoteAudioTrack) {
			fmt.Printf("\n✅ [Callback] OnUserAudioTrackSubscribed: uid=%s\n", uid)

			p.mu.Lock()
			defer p.mu.Unlock()

			// Save remote audio track
			p.remoteAudioTrack = remoteAudioTrack

			fmt.Println("\n=== Configuring Remote Audio Filter (using wrapper) ===")

			// Create RemoteAudioProcessor
			processor := agoraservice.NewRemoteAudioProcessor(remoteAudioTrack)
			p.remoteAudioProcessor = processor

			// Configure audio processing with simple API
			config := &agoraservice.AudioProcessorConfig{
				ANS: &agoraservice.ANSConfig{
					Enabled: true,
				},
				EnableDump: true,
				// Use custom config for advanced AINS settings
				CustomConfig: `{"aec":{"split_srate_for_48k":16000},"ans":{"enabled":true},"sf_st_cfg":{"enabled":true,"ainsModelPref":10},"sf_ext_cfg":{"nsngAlgRoute":12,"nsngPredefAgg":11}}`,
			}

			// Enable audio processing (automatically handles all steps)
			err := processor.Enable(config)
			if err != nil {
				fmt.Printf("Failed to enable audio processing: %v\n", err)
				return
			}
		},

		OnUserAudioTrackStateChanged: func(localUser *agoraservice.LocalUser, uid string, remoteAudioTrack *agoraservice.RemoteAudioTrack, state, reason, elapsed int) {
			fmt.Printf("[Callback] OnUserAudioTrackStateChanged: uid=%s, state=%d, reason=%d\n", uid, state, reason)
		},
	}

	p.connRecv.RegisterLocalUserObserver(observer)
	fmt.Println("✓ Receiver observer registered")
}

// SendAudioData sends audio data from PCM file
func (p *RemoteAudioFilterProcessorWrapper) SendAudioData(pcmFilePath string, sampleRate, channels int) {
	fmt.Printf("\n=== Sending audio data from: %s ===\n", pcmFilePath)
	fmt.Printf("Sample rate: %d Hz, Channels: %d\n", sampleRate, channels)

	file, err := os.Open(pcmFilePath)
	if err != nil {
		fmt.Printf("Failed to open PCM file: %v\n", err)
		return
	}
	defer file.Close()

	// Get file size
	stat, _ := file.Stat()
	fileSize := stat.Size()
	fmt.Printf("File size: %d bytes\n", fileSize)

	// Calculate samples per frame (10ms)
	samplesPerChannel := sampleRate / 100
	bytesPerFrame := samplesPerChannel * channels * 2 // 2 bytes per sample
	buffer := make([]byte, bytesPerFrame)

	timestamp := uint32(0)
	frameCount := 0
	successCount := 0
	failCount := 0

	fmt.Println("Starting to send audio...")

	for p.isRunning {
		n, err := file.Read(buffer)
		if err != nil || n == 0 {
			// EOF, loop back to start
			file.Seek(0, 0)
			continue
		}

		// Send audio frame
		frame := &agoraservice.AudioFrame{
			Type:              agoraservice.AudioFrameTypePCM16,
			SamplesPerChannel: samplesPerChannel,
			BytesPerSample:    2,
			Channels:          channels,
			SamplesPerSec:     sampleRate,
			Buffer:            buffer[:n],
			RenderTimeMs:      int64(timestamp),
		}
		ret := p.pcmSender.SendAudioPcmData(frame)
		if ret == 0 {
			successCount++
		} else {
			failCount++
			if failCount <= 5 {
				fmt.Printf("Send failed, error: %d\n", ret)
			}
		}

		frameCount++
		if frameCount%100 == 0 {
			fmt.Printf("Sent %d frames (success: %d, failed: %d)\n", frameCount, successCount, failCount)
		}

		timestamp += 10
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Printf("\nAudio sending stopped. Total frames: %d (success: %d, failed: %d)\n", frameCount, successCount, failCount)
}

// Run keeps the program running to receive audio
func (p *RemoteAudioFilterProcessorWrapper) Run(pcmFilePath string, sampleRate, channels, duration int) {
	fmt.Printf("\n=== Running test for %d seconds ===\n", duration)
	
	// Start sending audio in a goroutine immediately
	fmt.Println("Starting audio sender...")
	go p.SendAudioData(pcmFilePath, sampleRate, channels)

	// Wait a bit for the sending to start and receiver to subscribe
	time.Sleep(2 * time.Second)
	fmt.Println("Sender started, waiting for subscription...")

	// Wait for specified duration
	for i := 0; i < duration; i++ {
		time.Sleep(1 * time.Second)
		if (i+1)%10 == 0 {
			fmt.Printf("Running... %d/%d seconds\n", i+1, duration)
		}
	}

	fmt.Println("\n=== Duration completed ===")
}

// Cleanup cleans up resources
func (p *RemoteAudioFilterProcessorWrapper) Cleanup() {
	fmt.Println("\n=== Cleanup resources ===")

	p.isRunning = false

	// Cleanup receiver connection
	if p.connRecv != nil {
		p.connRecv.Disconnect()
		p.connRecv.Release()
		p.connRecv = nil
		fmt.Println("✓ Receiver connection released")
	}

	// Cleanup sender connection
	if p.connSend != nil {
		// 手动 unpublish 我们创建的 track
		if p.localUserSend != nil && p.localAudioTrack != nil {
			p.localUserSend.UnpublishAudioTrack(p.localAudioTrack)
			fmt.Println("✓ Audio track unpublished")
		}
		p.connSend.Disconnect()
		p.connSend.Release()
		p.connSend = nil
		fmt.Println("✓ Sender connection released")
	}

	// Release local audio track (手动创建的)
	if p.localAudioTrack != nil {
		p.localAudioTrack.Release()
		p.localAudioTrack = nil
		fmt.Println("✓ LocalAudioTrack released")
	}

	// Release PCM sender (手动创建的)
	if p.pcmSender != nil {
		p.pcmSender.Release()
		p.pcmSender = nil
		fmt.Println("✓ AudioPcmDataSender released")
	}

	// Release service
	agoraservice.Release()
	fmt.Println("✓ AgoraService released")
}

// global processor variable
var globalProcessor *RemoteAudioFilterProcessorWrapper

func main() {
	// Catch terminal signal
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		fmt.Println("\nReceived termination signal, cleaning up...")

		if globalProcessor != nil {
			globalProcessor.Cleanup()
		}

		time.Sleep(1 * time.Second)
		fmt.Println("Exit")
		os.Exit(0)
	}()

	// Parse command line arguments
	args := os.Args
	if len(args) < 3 {
		fmt.Println("Usage: apm_filter_remote <appid> <channel_id> [duration_seconds]")
		fmt.Println("  appid: Agora App ID")
		fmt.Println("  channel_id: Channel to join")
		fmt.Println("  duration_seconds: (optional) Duration in seconds, default 40")
		fmt.Println("")
		fmt.Println("Example:")
		fmt.Println("  apm_filter_remote your_app_id test_channel 40")
		return
	}

	appid := args[1]
	channelId := args[2]

	duration := 40 // Default 40 seconds
	if len(args) >= 4 {
		fmt.Sscanf(args[3], "%d", &duration)
	}

	if appid == "" {
		fmt.Println("Error: AGORA_APP_ID is required")
		return
	}

	// Use fixed PCM file path
	pcmFilePath := "/Volumes/ZR/Agora/SERVER/SDK/Agora-Golang-Server-SDK/go_sdk/examples/apm_filter_remote_wrapper/noise.pcm"
	
	// Check if PCM file exists
	if _, err := os.Stat(pcmFilePath); os.IsNotExist(err) {
		fmt.Printf("Error: PCM file not found: %s\n", pcmFilePath)
		fmt.Println("Please make sure noise.pcm exists in the apm_filter_remote directory")
		return
	}
	
	fmt.Printf("Using PCM file: %s\n", pcmFilePath)

	// Create processor
	processor := NewRemoteAudioFilterProcessorWrapper()
	globalProcessor = processor

	// Initialize with two different user IDs (sender and receiver)
	senderUid := "10001"    // Broadcaster
	receiverUid := "10002"  // Audience
	err := processor.Initialize(appid, channelId, senderUid, receiverUid)
	if err != nil {
		fmt.Printf("Failed to initialize: %v\n", err)
		return
	}
	defer processor.Cleanup()

	processor.Run(pcmFilePath, 48000, 1, duration)
}
