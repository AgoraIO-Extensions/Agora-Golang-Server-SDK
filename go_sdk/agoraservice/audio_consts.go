package agoraservice

// AudioScenario represents the audio scenario.
type AudioScenario int

const (
	// AudioScenarioDefault is the default audio scenario.
	AudioScenarioDefault AudioScenario = 0
	// AudioScenarioGameStreaming is the live gaming scenario, which needs to enable gaming
	// audio effects in the speaker. Choose this scenario to achieve high-fidelity
	// music playback.
	AudioScenarioGameStreaming AudioScenario = 3
	// AudioScenarioChatRoom is the chatroom scenario, which needs to keep recording when setClientRole to audience.
	// Normally, app developer can also use mute api to achieve the same result,
	// and we implement this 'non-orthogonal' behavior only to make API backward compatible.
	AudioScenarioChatRoom AudioScenario = 5
	// AudioScenarioChorus is the chorus scenario.
	AudioScenarioChorus AudioScenario = 7
	// AudioScenarioMeeting is the meeting scenario.
	AudioScenarioMeeting AudioScenario = 8
	// AudioScenarioAiServer is the AI server scenario.
	AudioScenarioAiServer AudioScenario = 9
	// AudioScenarioAiClient is the AI client scenario.
	AudioScenarioAiClient AudioScenario = 10
	// AudioScenarioNum is the number of audio scenarios.
	AudioScenarioNum AudioScenario = 11
)

// AudioProfile represents the audio profile.
type AudioProfile int

const (
	// AudioProfileDefault is the default audio profile.
	// - In the Communication profile, it represents a sample rate of 16 kHz, music encoding, mono, and a bitrate
	// of up to 16 Kbps.
	// - In the Live-broadcast profile, it represents a sample rate of 48 kHz, music encoding, mono, and a bitrate
	// of up to 64 Kbps.
	AudioProfileDefault AudioProfile = 0
	// AudioProfileSpeechStandard represents a sample rate of 16 kHz, audio encoding, mono, and a bitrate up to 18 Kbps.
	AudioProfileSpeechStandard AudioProfile = 1
	// AudioProfileMusicStandard represents a sample rate of 48 kHz, music encoding, mono, and a bitrate of up to 64 Kbps.
	AudioProfileMusicStandard AudioProfile = 2
	// AudioProfileMusicStandardStereo represents a sample rate of 48 kHz, music encoding, stereo, and a bitrate of up to 80 Kbps.
	AudioProfileMusicStandardStereo AudioProfile = 3
	// AudioProfileMusicHighQuality represents a sample rate of 48 kHz, music encoding, mono, and a bitrate of up to 96 Kbps.
	AudioProfileMusicHighQuality AudioProfile = 4
	// AudioProfileMusicHighQualityStereo represents a sample rate of 48 kHz, music encoding, stereo, and a bitrate of up to 128 Kbps.
	AudioProfileMusicHighQualityStereo AudioProfile = 5
	// AudioProfileIot represents a sample rate of 16 kHz, audio encoding, mono, and a bitrate of up to 64 Kbps.
	AudioProfileIot AudioProfile = 6
)

// VadState represents the VAD (Voice Activity Detection) state.
type VadState int

const (
	// VadStateInvalid represents an invalid VAD state.
	VadStateInvalid VadState = -1
	// VadStateWaitSpeeking represents the state when waiting for speaking.
	VadStateNoSpeeking VadState = 0
	// VadStateStartSpeeking represents the state when speaking starts.
	VadStateStartSpeeking VadState = 1
	// VadStateIsSpeeking represents the state when currently speaking.
	VadStateSpeeking VadState = 2
	// VadStateStopSpeeking represents the state when speaking stops.
	VadStateStopSpeeking VadState = 3
)

// AudioTrackMixingState represents the audio track mixing state.
type AudioTrackMixingState int

const (
	// AudioTrackMixEnabled represents the state when audio track mixing is enabled.
	AudioTrackMixEnabled AudioTrackMixingState = 0
	// AudioTrackMixDisabled represents the state when audio track mixing is disabled.
	AudioTrackMixDisabled AudioTrackMixingState = 1
)

// AudioCodecType represents the audio codec types.
type AudioCodecType int

const (
	// AudioCodecOpus represents the OPUS audio codec.
	AudioCodecOpus AudioCodecType = 1
	// AudioCodecPcma represents the PCMA audio codec.
	AudioCodecPcma AudioCodecType = 3
	// AudioCodecPcmu represents the PCMU audio codec.
	AudioCodecPcmu AudioCodecType = 4
	// AudioCodecG722 represents the G722 audio codec.
	AudioCodecG722 AudioCodecType = 5
	// AudioCodecAacLc represents the AAC LC audio codec.
	AudioCodecAacLc AudioCodecType = 8
	// AudioCodecHeAac represents the HE AAC audio codec.
	AudioCodecHeAac AudioCodecType = 9
	// AudioCodecJc1 represents the JC1 audio codec.
	AudioCodecJc1 AudioCodecType = 10
	// AudioCodecHeAac2 represents the HE AAC 2 audio codec.
	AudioCodecHeAac2 AudioCodecType = 11
	// AudioCodecLpcnet represents the LPCNET audio codec.
	AudioCodecLpcnet AudioCodecType = 12
	// 13: Opus multi channel codec, supporting 3 to 8 channels audio.
	AudioCodecOpusMC AudioCodecType = 13
)

// AudioEncodingType represents the audio encoding types of audio encoded frame observer.
type AudioEncodingType int

const (
	// AudioEncodingTypeAac16000Low represents the AAC audio encoding type with a sample rate of 16000 and low quality.
	AudioEncodingTypeAac16000Low AudioEncodingType = 0x010101
	// AudioEncodingTypeAac16000Medium represents the AAC audio encoding type with a sample rate of 16000 and medium quality.
	AudioEncodingTypeAac16000Medium AudioEncodingType = 0x010102
	// AudioEncodingTypeAac32000Low represents the AAC audio encoding type with a sample rate of 32000 and low quality.
	AudioEncodingTypeAac32000Low AudioEncodingType = 0x010201
	// AudioEncodingTypeAac32000Medium represents the AAC audio encoding type with a sample rate of 32000 and medium quality.
	AudioEncodingTypeAac32000Medium AudioEncodingType = 0x010202
	// AudioEncodingTypeAac32000High represents the AAC audio encoding type with a sample rate of 32000 and high quality.
	AudioEncodingTypeAac32000High AudioEncodingType = 0x010203
	// AudioEncodingTypeAac48000Medium represents the AAC audio encoding type with a sample rate of 48000 and medium quality.
	AudioEncodingTypeAac48000Medium AudioEncodingType = 0x010302
	// AudioEncodingTypeAac48000High represents the AAC audio encoding type with a sample rate of 48000 and high quality.
	AudioEncodingTypeAac48000High AudioEncodingType = 0x010303
	// AudioEncodingTypeOpus16000Low represents the OPUS audio encoding type with a sample rate of 16000 and low quality.
	AudioEncodingTypeOpus16000Low AudioEncodingType = 0x020101
	// AudioEncodingTypeOpus16000Medium represents the OPUS audio encoding type with a sample rate of 16000 and medium quality.
	AudioEncodingTypeOpus16000Medium AudioEncodingType = 0x020102
	// AudioEncodingTypeOpus48000Medium represents the OPUS audio encoding type with a sample rate of 48000 and medium quality.
	AudioEncodingTypeOpus48000Medium AudioEncodingType = 0x020302
	// AudioEncodingTypeOpus48000High represents the OPUS audio encoding type with a sample rate of 48000 and high quality.
	AudioEncodingTypeOpus48000High AudioEncodingType = 0x020303
)

// RawAudioFrameOpModeType represents the raw audio frame operation mode type.
type RawAudioFrameOpModeType int

const (
	// RawAudioFrameOpModeReadOnly represents the read-only mode for raw audio frames.
	// Users only read the audio frame data without modifying anything.
	RawAudioFrameOpModeReadOnly RawAudioFrameOpModeType = 0
	// RawAudioFrameOpModeReadWrite represents the read-write mode for raw audio frames.
	// Users read the data from audio frame, modify it, and then play it.
	RawAudioFrameOpModeReadWrite RawAudioFrameOpModeType = 2
)

type AudioFrameType int

const (
	AudioFrameTypePCM16 AudioFrameType = 0
)

type AudioPublishType int

const (
	AudioPublishTypeNoPublish AudioPublishType = 0
	AudioPublishTypePcm AudioPublishType = 1
	AudioPublishTypeEncodedPcm AudioPublishType = 2
)
type VideoPublishType int

const (
	VideoPublishTypeNoPublish VideoPublishType = 0
	VideoPublishTypeYuv VideoPublishType = 1
	VideoPublishTypeEncodedImage VideoPublishType = 2
)
