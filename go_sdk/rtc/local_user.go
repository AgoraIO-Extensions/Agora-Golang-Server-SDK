package agoraservice

/*
#cgo CFLAGS: -I../../agora_sdk/include/c/api2 -I../../agora_sdk/include/c/base

#include "agora_local_user.h"
#include "agora_rtc_conn.h"
#include "agora_service.h"
#include "agora_media_base.h"
#include "agora_parameter.h"
#include <stdio.h>
#include <stdbool.h>
*/
import "C"
import (
	"unsafe"
)

type LocalUser struct {
	cLocalUser unsafe.Pointer
	publishFlag bool
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



func (localUser *LocalUser) SetUserRole(role int) int {
	if localUser.cLocalUser == nil {
		return -1000
	}
	C.agora_local_user_set_user_role(localUser.cLocalUser, C.int(role))
	return 0
}

func (localUser *LocalUser) SubscribeAudio(uid string) int {
	if localUser.cLocalUser == nil {
		return -1000
	}
	cUid := C.CString(uid)
	defer C.free(unsafe.Pointer(cUid))
	return int(C.agora_local_user_subscribe_audio(localUser.cLocalUser, cUid))
}

func (localUser *LocalUser) UnsubscribeAudio(uid string) int {
	if localUser.cLocalUser == nil {
		return -1000
	}
	cUid := C.CString(uid)
	defer C.free(unsafe.Pointer(cUid))
	return int(C.agora_local_user_unsubscribe_audio(localUser.cLocalUser, cUid))
}

func (localUser *LocalUser) SubscribeAllAudio() int {
	if localUser.cLocalUser == nil {
		return -1000
	}
	return int(C.agora_local_user_subscribe_all_audio(localUser.cLocalUser))
}

func (localUser *LocalUser) UnsubscribeAllAudio() int {
	if localUser.cLocalUser == nil {
		return -1000
	}
	return int(C.agora_local_user_unsubscribe_all_audio(localUser.cLocalUser))
}

func (localUser *LocalUser) SubscribeVideo(uid string, options *VideoSubscriptionOptions) int {
	if localUser.cLocalUser == nil {
		return -1000
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
		return -1000
	}
	cUid := C.CString(uid)
	defer C.free(unsafe.Pointer(cUid))
	return int(C.agora_local_user_unsubscribe_video(localUser.cLocalUser, cUid))
}

func (localUser *LocalUser) SubscribeAllVideo(options *VideoSubscriptionOptions) int {
	if localUser.cLocalUser == nil {
		return -1000
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
		return -1000
	}
	return int(C.agora_local_user_unsubscribe_all_video(localUser.cLocalUser))
}

func (localUser *LocalUser) SetAudioEncoderConfiguration(config *AudioEncoderConfiguration) int {
	if localUser.cLocalUser == nil {
		return -1000
	}
	cConfig := C.struct__audio_encoder_config{}
	cConfig.audio_profile = C.int(AudioProfileDefault)
	if config != nil {
		cConfig.audio_profile = C.int(config.AudioProfile)
	}
	return int(C.agora_local_user_set_audio_encoder_config(localUser.cLocalUser, &cConfig))
}
// Implement the publication of audio track


func (localUser *LocalUser) publishAudio(track *LocalAudioTrack) int {
	if localUser.cLocalUser == nil {
		return -1000
	}

	if localUser.publishFlag == true {
	    return 0
	}
	
	ret := int(C.agora_local_user_publish_audio(localUser.cLocalUser, track.cTrack))
	localUser.publishFlag = true
	if ret != 0 {
		localUser.publishFlag = false
	}
	return ret
	
}



func (localUser *LocalUser) unpublishAudio(track *LocalAudioTrack) int {
	if localUser.cLocalUser == nil {
		return -1000
	}
	if localUser.publishFlag == false {
		return 0
	}

	ret := int(C.agora_local_user_unpublish_audio(localUser.cLocalUser, track.cTrack))
	localUser.publishFlag = false
	if ret != 0 {
		localUser.publishFlag = true
	}
	return ret
}

func (localUser *LocalUser) publishVideo(track *LocalVideoTrack) int {
	if localUser.cLocalUser == nil {
		return -1000
	}
	return int(C.agora_local_user_publish_video(localUser.cLocalUser, track.cTrack))
}

func (localUser *LocalUser) unpublishVideo(track *LocalVideoTrack) int {
	if localUser.cLocalUser == nil {
		return -1000
	}
	return int(C.agora_local_user_unpublish_video(localUser.cLocalUser, track.cTrack))
}

func (localUser *LocalUser) SetPlaybackAudioFrameParameters(channels int, sampleRate int, mode int, samplesPerCall int) int {
	if localUser.cLocalUser == nil {
		return -1000
	}
	return int(C.agora_local_user_set_playback_audio_frame_parameters(localUser.cLocalUser, C.uint(channels), C.uint(sampleRate), C.int(mode), C.int(samplesPerCall)))
}

func (localUser *LocalUser) SetRecordingAudioFrameParameters(channels int, sampleRate int, mode int, samplesPerCall int) int {
	if localUser.cLocalUser == nil {
		return -1000
	}
	return int(C.agora_local_user_set_recording_audio_frame_parameters(localUser.cLocalUser, C.uint(channels), C.uint(sampleRate), C.int(mode), C.int(samplesPerCall)))
}

func (localUser *LocalUser) SetMixedAudioFrameParameters(channels int, sampleRate int, samplesPerCall int) int {
	if localUser.cLocalUser == nil {
		return -1000
	}
	return int(C.agora_local_user_set_mixed_audio_frame_parameters(localUser.cLocalUser, C.uint(channels), C.uint(sampleRate), C.int(samplesPerCall)))
}

func (localUser *LocalUser) SetPlaybackAudioFrameBeforeMixingParameters(channels int, sampleRate int) int {
	if localUser.cLocalUser == nil {
		return -1000
	}
	return int(C.agora_local_user_set_playback_audio_frame_before_mixing_parameters(localUser.cLocalUser, C.uint(channels), C.uint(sampleRate)))
}

func (localUser *LocalUser) SetAudioScenario(audioScenario AudioScenario) int {
	if localUser.cLocalUser == nil {
		return -1000
	}
	return int(C.agora_local_user_set_audio_scenario(localUser.cLocalUser, C.int(audioScenario)))
}
func (localUser *LocalUser) SetAudioVolumeIndicationParameters(intervalInMs int, smooth int, reportVad bool) int {
	if localUser.cLocalUser == nil {
		return -1000
	}
	
	ret := C.agora_local_user_set_audio_volume_indication_parameters(localUser.cLocalUser, C.int(intervalInMs), C.int(smooth), C.bool(reportVad))
	return int(ret)
}
func (localUser *LocalUser) sendAudioMetaData(metaData []byte) int {
	if localUser.cLocalUser == nil {
		return -1000
	}
	cMetaData := C.CBytes(metaData)
	defer C.free(cMetaData)
	
	ret := C.agora_local_user_send_audio_meta_data(localUser.cLocalUser, (*C.char)(cMetaData), (C.size_t)(len(metaData)))
	return int(ret)
}