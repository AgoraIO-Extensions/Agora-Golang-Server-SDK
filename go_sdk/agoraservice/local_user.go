package agoraservice

/*
#cgo CFLAGS: -I../../agora_sdk/include/c/api2 -I../../agora_sdk/include/c/base

#include "agora_local_user.h"
#include "agora_rtc_conn.h"
#include "agora_service.h"
#include "agora_media_base.h"
#include "agora_parameter.h"
*/
import "C"
import "unsafe"

type LocalUser struct {
	connection *RtcConnection
	cLocalUser unsafe.Pointer
}

type VideoSubscriptionOptions struct {
	/**
	 * The type of the video stream to subscribe to: #VideoStreamXxx.
	 *
	 * The default value is VIDEO_STREAM_HIGH, which means the high-resolution and high-bitrate
	 * video stream.
	 */
	StreamType VideoStreamType
	/**
	 * Determines whether to subscribe to encoded video data only:
	 * - true: Subscribe to encoded video data only.
	 * - false: (Default) Do not subscribe to encoded video data only.
	 */
	EncodedFrameOnly bool
}

func (localUser *LocalUser) GetRtcConnection() *RtcConnection {
	return localUser.connection
}

func (localUser *LocalUser) RegisterLocalUserObserver(observer *LocalUserObserver) int {
	return localUser.connection.registerLocalUserObserver(observer)
}

func (localUser *LocalUser) UnregisterLocalUserObserver() int {
	return localUser.connection.unregisterLocalUserObserver()
}

func (localUser *LocalUser) RegisterAudioFrameObserver(observer *AudioFrameObserver) int {
	return localUser.connection.registerAudioFrameObserver(observer)
}

func (localUser *LocalUser) UnregisterAudioFrameObserver() int {
	return localUser.connection.unregisterAudioFrameObserver()
}

func (localUser *LocalUser) RegisterVideoFrameObserver(observer *VideoFrameObserver) int {
	return localUser.connection.registerVideoFrameObserver(observer)
}

func (localUser *LocalUser) UnregisterVideoFrameObserver() int {
	return localUser.connection.unregisterVideoFrameObserver()
}

func (localUser *LocalUser) RegisterVideoEncodedFrameObserver(observer *VideoEncodedFrameObserver) int {
	return localUser.connection.registerVideoEncodedFrameObserver(observer)
}

func (localUser *LocalUser) UnregisterVideoEncodedFrameObserver() int {
	return localUser.connection.unregisterVideoEncodedFrameObserver()
}

func (localUser *LocalUser) SetUserRole(role int) int {
	if localUser.cLocalUser == nil {
		return -1
	}
	C.agora_local_user_set_user_role(localUser.cLocalUser, C.int(role))
	return 0
}

func (localUser *LocalUser) SubscribeAudio(uid string) int {
	if localUser.cLocalUser == nil {
		return -1
	}
	cUid := C.CString(uid)
	defer C.free(unsafe.Pointer(cUid))
	return int(C.agora_local_user_subscribe_audio(localUser.cLocalUser, cUid))
}

func (localUser *LocalUser) UnsubscribeAudio(uid string) int {
	if localUser.cLocalUser == nil {
		return -1
	}
	cUid := C.CString(uid)
	defer C.free(unsafe.Pointer(cUid))
	return int(C.agora_local_user_unsubscribe_audio(localUser.cLocalUser, cUid))
}

func (localUser *LocalUser) SubscribeAllAudio() int {
	if localUser.cLocalUser == nil {
		return -1
	}
	return int(C.agora_local_user_subscribe_all_audio(localUser.cLocalUser))
}

func (localUser *LocalUser) UnsubscribeAllAudio() int {
	if localUser.cLocalUser == nil {
		return -1
	}
	return int(C.agora_local_user_unsubscribe_all_audio(localUser.cLocalUser))
}

func (localUser *LocalUser) SubscribeVideo(uid string, options *VideoSubscriptionOptions) int {
	if localUser.cLocalUser == nil {
		return -1
	}
	cUid := C.CString(uid)
	defer C.free(unsafe.Pointer(cUid))
	cOptions := C.video_subscription_options{}
	po := &cOptions
	if options != nil {
		cOptions._type = C.int(options.StreamType)
		cOptions.encoded_frame_only = C.int(0)
		if options.EncodedFrameOnly {
			cOptions.encoded_frame_only = C.int(1)
		}
	} else {
		po = nil
	}
	return int(C.agora_local_user_subscribe_video(localUser.cLocalUser, cUid, po))
}

func (localUser *LocalUser) UnsubscribeVideo(uid string) int {
	if localUser.cLocalUser == nil {
		return -1
	}
	cUid := C.CString(uid)
	defer C.free(unsafe.Pointer(cUid))
	return int(C.agora_local_user_unsubscribe_video(localUser.cLocalUser, cUid))
}

func (localUser *LocalUser) SubscribeAllVideo(options *VideoSubscriptionOptions) int {
	if localUser.cLocalUser == nil {
		return -1
	}
	cOptions := C.video_subscription_options{}
	po := &cOptions
	if options != nil {
		cOptions._type = C.int(options.StreamType)
		cOptions.encoded_frame_only = C.int(0)
		if options.EncodedFrameOnly {
			cOptions.encoded_frame_only = C.int(1)
		}
	} else {
		po = nil
	}
	return int(C.agora_local_user_subscribe_all_video(localUser.cLocalUser, po))
}

func (localUser *LocalUser) UnsubscribeAllVideo() int {
	if localUser.cLocalUser == nil {
		return -1
	}
	return int(C.agora_local_user_unsubscribe_all_video(localUser.cLocalUser))
}

func (localUser *LocalUser) SetAudioEncoderConfiguration(config *AudioEncoderConfiguration) int {
	if localUser.cLocalUser == nil {
		return -1
	}
	cConfig := C.struct__audio_encoder_config{}
	cConfig.audio_profile = C.int(AudioProfileDefault)
	if config != nil {
		cConfig.audio_profile = C.int(config.AudioProfile)
	}
	return int(C.agora_local_user_set_audio_encoder_config(localUser.cLocalUser, &cConfig))
}

func (localUser *LocalUser) PublishAudio(track *LocalAudioTrack) int {
	if localUser.cLocalUser == nil {
		return -1
	}
	return int(C.agora_local_user_publish_audio(localUser.cLocalUser, track.cTrack))
}

func (localUser *LocalUser) UnpublishAudio(track *LocalAudioTrack) int {
	if localUser.cLocalUser == nil {
		return -1
	}
	return int(C.agora_local_user_unpublish_audio(localUser.cLocalUser, track.cTrack))
}

func (localUser *LocalUser) PublishVideo(track *LocalVideoTrack) int {
	if localUser.cLocalUser == nil {
		return -1
	}
	return int(C.agora_local_user_publish_video(localUser.cLocalUser, track.cTrack))
}

func (localUser *LocalUser) UnpublishVideo(track *LocalVideoTrack) int {
	if localUser.cLocalUser == nil {
		return -1
	}
	return int(C.agora_local_user_unpublish_video(localUser.cLocalUser, track.cTrack))
}

func (localUser *LocalUser) SetPlaybackAudioFrameParameters(channels int, sampleRate int, mode int, samplesPerCall int) int {
	if localUser.cLocalUser == nil {
		return -1
	}
	return int(C.agora_local_user_set_playback_audio_frame_parameters(localUser.cLocalUser, C.uint(channels), C.uint(sampleRate), C.int(mode), C.int(samplesPerCall)))
}

func (localUser *LocalUser) SetRecordingAudioFrameParameters(channels int, sampleRate int, mode int, samplesPerCall int) int {
	if localUser.cLocalUser == nil {
		return -1
	}
	return int(C.agora_local_user_set_recording_audio_frame_parameters(localUser.cLocalUser, C.uint(channels), C.uint(sampleRate), C.int(mode), C.int(samplesPerCall)))
}

func (localUser *LocalUser) SetMixedAudioFrameParameters(channels int, sampleRate int, samplesPerCall int) int {
	if localUser.cLocalUser == nil {
		return -1
	}
	return int(C.agora_local_user_set_mixed_audio_frame_parameters(localUser.cLocalUser, C.uint(channels), C.uint(sampleRate), C.int(samplesPerCall)))
}

func (localUser *LocalUser) SetPlaybackAudioFrameBeforeMixingParameters(channels int, sampleRate int) int {
	if localUser.cLocalUser == nil {
		return -1
	}
	return int(C.agora_local_user_set_playback_audio_frame_before_mixing_parameters(localUser.cLocalUser, C.uint(channels), C.uint(sampleRate)))
}

func (localUser *LocalUser) SetAudioScenario(audioScenario AudioScenario) int {
	if localUser.cLocalUser == nil {
		return -1
	}
	return int(C.agora_local_user_set_audio_scenario(localUser.cLocalUser, C.int(audioScenario)))
}
