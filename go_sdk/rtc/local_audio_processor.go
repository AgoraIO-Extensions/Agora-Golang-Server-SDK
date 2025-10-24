package agoraservice

import (
	"encoding/json"
	"fmt"
)

// ============================================
// Local Audio Processor - Uplink Audio Processing
// ============================================
// This file processes audio sent from local client using pcm_source extension.
// It manages PCM Sender and LocalAudioTrack lifecycle.
//
// Extension: audio_processing_pcm_source
// Position: pcm_source
// Use case: Process local captured/sent audio (microphone input)
//
// Core Class: LocalAudioProcessor
//   - Automatically handles Extension initialization
//   - Automatically creates PCM Sender and AudioTrack
//   - Configures audio processing (BGHVS, AGC, ANS, AINS)
//   - Simplifies audio data transmission
//   - Simplifies AudioSink addition
//
// Usage Flow:
//   1. Initialize AgoraService
//   2. Create LocalAudioProcessor
//   3. Add AudioSink (optional - to monitor processed audio)
//   4. Send PCM data
//
// Example:
//   processor, _ := agoraservice.NewLocalAudioProcessor(
//       agoraservice.AudioScenarioChorus,
//       &agoraservice.AudioProcessorConfig{
//           BGHVS: &agoraservice.BGHVSConfig{Enabled: true},
//           AGC:   &agoraservice.AGCConfig{Enabled: true},
//           ANS:   &agoraservice.ANSConfig{Enabled: true},
//       },
//   )
//   defer processor.Release()
//   processor.SendPCM(pcmData, 48000, 1)
// ============================================

// AudioProcessorConfig configuration for audio processing (both local and remote)
// This config applies to audio processing module (APM) including:
//   - BGHVS: Background/Human Voice Separation
//   - AGC: Automatic Gain Control
//   - ANS: Acoustic Noise Suppression
type AudioProcessorConfig struct {
	// BGHVS (Background/Human Voice Separation) configuration
	BGHVS *BGHVSConfig
	
	// AGC (Automatic Gain Control) configuration
	AGC *AGCConfig
	
	// ANS (Acoustic Noise Suppression) configuration
	ANS *ANSConfig
	
	// AEC (Acoustic Echo Cancellation) configuration
	AEC *AECConfig
	
	// EnableDump enables audio dump for debugging, default false
	EnableDump bool
	
	// CustomConfig custom configuration string (advanced usage)
	// If set, individual configs (BGHVS, AGC, ANS) will be merged with this
	// Example: `{"aec":{"split_srate_for_48k":16000},...}`
	CustomConfig string
}

// BGHVSConfig Background/Human Voice Separation configuration
type BGHVSConfig struct {
	// Enabled enables BGHVS
	Enabled bool
}

// AGCConfig Automatic Gain Control configuration
type AGCConfig struct {
	// Enabled enables AGC
	Enabled bool
}

// ANSConfig Acoustic Noise Suppression configuration
type ANSConfig struct {
	// Enabled enables ANS
	Enabled bool
}

// AECConfig Acoustic Echo Cancellation configuration
type AECConfig struct {
	// Enabled enables AEC
	Enabled bool
}

// ============================================
// Internal Functions
// ============================================

// initLocalAudioProcessingExtension initializes the local audio processing extension
func initLocalAudioProcessingExtension() error {
	ret := EnableExtension("agora.builtin", "audio_processing_pcm_source", "", true)
	if ret != 0 {
		return fmt.Errorf("[LocalAudioProcessor] failed to enable audio_processing_pcm_source extension, error code: %d", ret)
	}
	return nil
}

// configureLocalAudioFilter configures the audio filter on local track
func configureLocalAudioFilter(audioTrack *LocalAudioTrack, config *AudioProcessorConfig) error {
	if audioTrack == nil {
		return fmt.Errorf("[LocalAudioProcessor] audioTrack cannot be nil")
	}
	
	if config == nil {
		return fmt.Errorf("[LocalAudioProcessor] config cannot be nil")
	}
	
	var ret int
	
	// Step 1: Enable filter on Track
	ret = audioTrack.EnableAudioFilter("audio_processing_pcm_source", true, 3)
	if ret != 0 {
		return fmt.Errorf("[LocalAudioProcessor] failed to enable filter on Track, error code: %d", ret)
	}
	
	// Step 2: Load AINS resources (if needed)
	ret = audioTrack.SetFilterProperty("audio_processing_pcm_source", "apm_load_resource", "ains", 3)
	if ret != 0 {
		return fmt.Errorf("[LocalAudioProcessor] failed to load AINS model resources, error code: %d", ret)
	}
	
	// Step 3: Build and apply configuration
	configStr := buildAudioProcessorConfig(config)
	ret = audioTrack.SetFilterProperty("audio_processing_pcm_source", "apm_config", configStr, 3)
	if ret != 0 {
		return fmt.Errorf("[LocalAudioProcessor] failed to configure audio processing parameters, error code: %d", ret)
	}
	
	// Step 4: Enable dump (if debugging is needed)
	if config.EnableDump {
		ret = audioTrack.SetFilterProperty("audio_processing_pcm_source", "apm_dump", "true", 3)
		if ret != 0 {
			fmt.Printf("[LocalAudioProcessor] failed to enable dump, error code: %d (non-critical)\n", ret)
		}
	}
	
	return nil
}

// buildAudioProcessorConfig builds configuration JSON based on AudioProcessorConfig
// This function is shared by both LocalAudioProcessor and RemoteAudioProcessor
func buildAudioProcessorConfig(config *AudioProcessorConfig) string {
	// Start with custom config if provided, otherwise use default base
	var configMap map[string]interface{}
	
	if config.CustomConfig != "" {
		// Parse custom config as base
		if err := json.Unmarshal([]byte(config.CustomConfig), &configMap); err != nil {
			// If parse fails, start with empty map
			configMap = make(map[string]interface{})
		}
	} else {
		configMap = make(map[string]interface{})
	}
	
	// Helper function: merge structured config into existing config (instead of overwriting)
	// This preserves advanced parameters from CustomConfig while allowing structured config to override specific fields
	mergeConfig := func(key string, structValue map[string]interface{}) {
		if existing, ok := configMap[key]; ok {
			// If key already exists, merge instead of overwrite
			if existingMap, ok := existing.(map[string]interface{}); ok {
				// Merge structured config fields into existing config
				for k, v := range structValue {
					existingMap[k] = v  // Override specific fields only
				}
				return
			}
		}
		// If key doesn't exist, set directly
		configMap[key] = structValue
	}
	
	// Add BGHVS configuration (merge with existing config)
	if config.BGHVS != nil {
		mergeConfig("bghvs", map[string]interface{}{
			"enabled": config.BGHVS.Enabled,
		})
	}
	
	// Add AGC configuration (merge with existing config)
	if config.AGC != nil {
		mergeConfig("agc", map[string]interface{}{
			"enabled": config.AGC.Enabled,
		})
	}
	
	// Add ANS configuration (merge with existing config)
	if config.ANS != nil {
		mergeConfig("ans", map[string]interface{}{
			"enabled": config.ANS.Enabled,
		})
	}
	
	// Add AEC configuration (merge with existing config)
	if config.AEC != nil {
		mergeConfig("aec", map[string]interface{}{
			"enabled": config.AEC.Enabled,
		})
	}
	
	// Convert to JSON string
	jsonBytes, err := json.Marshal(configMap)
	if err != nil {
		// Fallback to basic config
		return `{"bghvs":{"enabled":false},"agc":{"enabled":false},"ans":{"enabled":false},"aec":{"enabled":false}}`
	}
	
	return string(jsonBytes)
}

// ============================================
// LocalAudioProcessor - Main Class
// ============================================
type LocalAudioProcessor struct {
	pcmSender  *AudioPcmDataSender
	audioTrack *LocalAudioTrack
}

// NewLocalAudioProcessor creates a local audio processor for uplink audio processing
func NewLocalAudioProcessor(scenario AudioScenario, config *AudioProcessorConfig) (*LocalAudioProcessor, error) {
	// Step 1: Initialize Extension (must be called before creating AudioTrack)
	err := initLocalAudioProcessingExtension()
	if err != nil {
		return nil, fmt.Errorf("[LocalAudioProcessor] failed to initialize extension: %w", err)
	}
	
	// Step 2: Create PCM Sender
	pcmSender := NewAudioPcmDataSender()
	if pcmSender == nil {
		return nil, fmt.Errorf("[LocalAudioProcessor] failed to create PCM Sender")
	}
	
	// Step 3: Create AudioTrack
	audioTrack := NewCustomAudioTrackPcm(pcmSender, scenario)
	if audioTrack == nil {
		pcmSender.Release()
		return nil, fmt.Errorf("[LocalAudioProcessor] failed to create AudioTrack")
	}
	
	// Step 4: Configure audio processing
	if config != nil {
		err = configureLocalAudioFilter(audioTrack, config)
		if err != nil {
			audioTrack.Release()
			pcmSender.Release()
			return nil, fmt.Errorf("[LocalAudioProcessor] failed to configure audio filter: %w", err)
		}
	}
	
	// Step 5: Enable Track
	audioTrack.SetEnabled(true)
	
	return &LocalAudioProcessor{
		pcmSender:  pcmSender,
		audioTrack: audioTrack,
	}, nil
}

// SendPCM sends PCM audio data to the processor
func (p *LocalAudioProcessor) SendPCM(data []byte, sampleRate, channels int) error {
	if p.pcmSender == nil {
		return fmt.Errorf("[LocalAudioProcessor] PCM Sender not initialized")
	}
	
	samplesPerChannel := len(data) / (channels * 2)
	
	frame := &AudioFrame{
		Buffer:            data,
		SamplesPerChannel: samplesPerChannel,
		BytesPerSample:    2,
		Channels:          channels,
		SamplesPerSec:     sampleRate,
		RenderTimeMs:      0,
		PresentTimeMs:     0,
	}
	
	ret := p.pcmSender.SendAudioPcmData(frame)
	if ret != 0 {
		return fmt.Errorf("[LocalAudioProcessor] failed to send PCM data, error code: %d", ret)
	}
	
	return nil
}

// AddSink adds an audio sink to monitor processed audio
func (p *LocalAudioProcessor) AddSink(callback func(*AudioFrame) bool, outputSampleRate, outputChannels int) (*AudioSinkContext, error) {
	if p.audioTrack == nil {
		return nil, fmt.Errorf("[LocalAudioProcessor] AudioTrack not initialized")
	}
	
	sink := NewAudioSink(callback)
	if sink == nil {
		return nil, fmt.Errorf("[LocalAudioProcessor] failed to create AudioSink")
	}
	
	ret := p.audioTrack.AddAudioSink(sink.GetHandle(), outputSampleRate, outputChannels)
	if ret != 0 {
		sink.Release()
		return nil, fmt.Errorf("[LocalAudioProcessor] failed to add AudioSink, error code: %d", ret)
	}
	
	return sink, nil
}

// Release releases all resources
func (p *LocalAudioProcessor) Release() {
	if p.audioTrack != nil {
		p.audioTrack.Release()
		p.audioTrack = nil
	}
	
	if p.pcmSender != nil {
		p.pcmSender.Release()
		p.pcmSender = nil
	}
}

