package agoraservice

import (
	"fmt"
)

// ============================================
// Remote Audio Processor - Downlink Audio Processing
// ============================================
// This file manages audio processing for remote tracks received in RTC connection.
// RemoteAudioTrack is obtained from OnUserAudioTrackSubscribed callback, NOT created manually.
//
// Extension: audio_processing_remote_playback
// Position: playback (2)
// Use case: Process audio received from remote participants
//
// Core Class: RemoteAudioProcessor
//   - Configure audio processing for subscribed remote tracks
//   - Apply BGHVS, AGC, ANS to remote audio
//
// Important Notes:
//   ⚠️  RemoteAudioTrack is obtained from subscription callback
//   ⚠️  Can ONLY configure after OnUserAudioTrackSubscribed is called
//   ⚠️  Must subscribe to remote user's audio track first
//
// Usage Flow:
//   1. Initialize AgoraService and create RTC Connection
//   2. Subscribe to remote user's audio track
//   3. In OnUserAudioTrackSubscribed callback:
//      - Get remoteAudioTrack from callback parameter
//      - Create RemoteAudioProcessor with the track
//      - Enable audio processing configuration
//
// Example:
//   // Set up LocalUserObserver
//   observer := &LocalUserObserver{
//       OnUserAudioTrackSubscribed: func(localUser *LocalUser, uid string, remoteAudioTrack *RemoteAudioTrack) {
//           fmt.Printf("Remote audio track subscribed from user: %s\n", uid)
//
//           // Now we can configure audio processing for this remote track
//           processor := NewRemoteAudioProcessor(remoteAudioTrack)
//           err := processor.Enable(&AudioProcessorConfig{
//               BGHVS: &BGHVSConfig{Enabled: true},
//               AGC:   &AGCConfig{Enabled: true},
//               ANS:   &ANSConfig{Enabled: true},
//           })
//           if err != nil {
//               fmt.Printf("Failed to enable audio processing: %v\n", err)
//           }
//       },
//   }
//
//   localUser := connection.GetLocalUser()
//   localUser.SetObserver(observer)
//
//   // Subscribe to remote user's audio
//   localUser.SubscribeAudio(uid)
// ============================================

// Note: AudioProcessorConfig is defined in local_audio_processor.go
// and is shared between LocalAudioProcessor and RemoteAudioProcessor

// ============================================
// RemoteAudioProcessor - Main Class
// ============================================

// RemoteAudioProcessor manages audio processing extensions for remote tracks
// Extension: audio_processing (position: playback)
type RemoteAudioProcessor struct {
	enabled    bool
	remoteAudioTrack *RemoteAudioTrack
}

// NewRemoteAudioProcessor creates a remote audio processor for a subscribed remote track
//
// ⚠️  Important: This should ONLY be called in OnUserAudioTrackSubscribed callback
// The remoteAudioTrack parameter comes from the subscription callback, not created manually.
//
// Parameters:
//   - remoteAudioTrack: the remote audio track obtained from OnUserAudioTrackSubscribed callback
//
// Returns:
//   - processor: RemoteAudioProcessor instance
//
// Usage:
//   observer := &LocalUserObserver{
//       OnUserAudioTrackSubscribed: func(localUser *LocalUser, uid string, remoteAudioTrack *RemoteAudioTrack) {
//           processor := NewRemoteAudioProcessor(remoteAudioTrack)
//           processor.Enable(config)
//       },
//   }
func NewRemoteAudioProcessor(remoteAudioTrack *RemoteAudioTrack) *RemoteAudioProcessor {
	return &RemoteAudioProcessor{
		enabled:   false,
		remoteAudioTrack: remoteAudioTrack,
	}
}

// Enable enables audio processing for the remote track with configuration
//
// ⚠️  Must be called AFTER remote track is subscribed (in OnUserAudioTrackSubscribed callback)
//
// This will:
//   1. Enable the audio_processing_remote_playback extension globally
//   2. Load AINS model resources (for AI noise suppression)
//   3. Configure audio processing parameters (BGHVS, AGC, ANS) on the remote track
//   4. Enable dump (if requested in config)
//
// Parameters:
//   - config: audio processing configuration (BGHVS, AGC, ANS settings)
//
// Returns:
//   - error: error message, nil indicates success
func (p *RemoteAudioProcessor) Enable(config *AudioProcessorConfig) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}
	
	if p.remoteAudioTrack == nil {
		return fmt.Errorf("remoteAudioTrack is nil, must be called with valid track from OnUserAudioTrackSubscribed")
	}
	
	// Step 1: Enable extension globally
	ret := EnableExtension("agora.builtin", "audio_processing_remote_playback", "", true)
	if ret != 0 {
		return fmt.Errorf("failed to enable extension, error code: %d", ret)
	}
	
	p.enabled = true
	
	// Step 2: Load AINS resources (if needed for AI noise suppression)
	ret = p.remoteAudioTrack.SetFilterProperty("audio_processing_remote_playback", "apm_load_resource", "ains", 2)
	if ret != 0 {
		return fmt.Errorf("failed to load AINS model resources, error code: %d", ret)
	}
	
	// Step 3: Set configuration
	configStr := buildAudioProcessorConfig(config)
	ret = p.remoteAudioTrack.SetFilterProperty("audio_processing_remote_playback", "apm_config", configStr, 2)
	if ret != 0 {
		return fmt.Errorf("failed to set configuration, error code: %d", ret)
	}

	// Step 4: Enable dump if needed
	if config.EnableDump {
		ret = p.remoteAudioTrack.SetFilterProperty("audio_processing_remote_playback", "apm_dump", "true", 2)
		if ret != 0 {
			return fmt.Errorf("failed to enable dump, error code: %d", ret)
		}
	}
	
	return nil
}

// Disable disables the remote audio processing extension
//
// Returns:
//   - error: error message, nil indicates success
func (p *RemoteAudioProcessor) Disable() error {
	if !p.enabled {
		return nil // Already disabled
	}
	
	ret := EnableExtension("agora.builtin", "audio_processing_remote_playback", "", false)
	if ret != 0 {
		return fmt.Errorf("failed to disable extension, error code: %d", ret)
	}
	
	p.enabled = false
	return nil
}

// IsEnabled returns whether the extension is currently enabled
func (p *RemoteAudioProcessor) IsEnabled() bool {
	return p.enabled
}

