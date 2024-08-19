package agoraservice

/*
#cgo CFLAGS: -I../../agora_sdk/include_c/api2 -I../../agora_sdk/include_c/base

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

func (localUser *RtcConnection) RegisterAudioFrameObserver(observer *AudioFrameObserver) int {
	return localUser.connection.registerAudioFrameObserver(observer)
}

func (localUser *RtcConnection) UnregisterAudioFrameObserver() int {
	return localUser.connection.unregisterAudioFrameObserver()
}

func (localUser *RtcConnection) RegisterVideoFrameObserver(observer *VideoFrameObserver) int {
	return localUser.connection.registerVideoFrameObserver(observer)
}

func (localUser *RtcConnection) UnregisterVideoFrameObserver() int {
	return localUser.connection.unregisterVideoFrameObserver()
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

func (localUser *LocalUser) SetPlaybackAudioFrameBeforeMixingParameters(channels int, sampleRate int) int {
	if localUser.cLocalUser == nil {
		return -1
	}
	return int(C.agora_local_user_set_playback_audio_frame_before_mixing_parameters(localUser.cLocalUser, C.uint(channels), C.uint(sampleRate)))
}
