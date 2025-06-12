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
	"fmt"
	"unsafe"
)

type LocalUser struct {
	connection *RtcConnection
	cLocalUser unsafe.Pointer
	audioTrack *LocalAudioTrack
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

func (localUser *LocalUser) RegisterAudioFrameObserver(observer *AudioFrameObserver,enableVad int, vadConfigure *AudioVadConfigV2) int {
	return localUser.connection.registerAudioFrameObserver(observer, enableVad, vadConfigure)
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
// Implement the publication of audio track


func (localUser *LocalUser) PublishAudio(track *LocalAudioTrack) int {
	if localUser.cLocalUser == nil {
		return -1
	}
	if localUser.audioTrack == nil {
		localUser.audioTrack = track
	} 
	fmt.Printf("______publish audio track id %d, old id %d, cTrack %p, old cTrack %p\n", track.id, localUser.audioTrack.id, track.cTrack, localUser.audioTrack.cTrack)
	
	// 将内部保持的cTrack赋值给track
	track.cTrack = localUser.audioTrack.cTrack
	return int(C.agora_local_user_publish_audio(localUser.cLocalUser, track.cTrack))
	
}

//TODO: 这样修改后，应用层也会被同步修改！ 
//ToDo: 提供一个打断api，设置在localUser::InterruptAudio()
//
// 内部只是做原子+正交化原则设计api
// 应用层需要做：打断操作！
// 然后调用UpdateAudioTrack
// 然后在做pcmsender的更新senario
// 然后调用PublishAudio
func (localUser *LocalUser) UpdateAudioTrack(senario AudioScenario) int {
	if localUser.cLocalUser == nil {
		return -1
	}
	// check if the track is the same as the internal one
	if localUser.audioTrack == nil {	
		return 0
	}
	// 如果senario是一样的，则不做任何操作
	if localUser.audioTrack.audioScenario == senario {
		return 0
	}
	// 如果senario不一样，则释放cTrack,重新创建一个新的cTrack，然后做assign
	// never change the id of the track
	localUser.audioTrack.Release()
	
	localUser.audioTrack.audioScenario = senario
	localUser.audioTrack.cTrack = nil
	//ToDo：回退交给客户来做会更好，用最佳实践的方式来
	//原因是因为客户选择的回退策略，不一定是chorus，而是其他策略！！这样我们固定策略，会限制客户
	var cTrack unsafe.Pointer = nil
	if senario == AudioScenarioAiServer {
		cTrack  = C.agora_service_create_direct_custom_audio_track_pcm(agoraService.service, localUser.audioTrack.pcmSender.cSender)
	} else {
		cTrack = C.agora_service_create_custom_audio_track_pcm(agoraService.service, localUser.audioTrack.pcmSender.cSender)
	}
	localUser.audioTrack.cTrack = cTrack

	localUser.audioTrack.SetEnabled(true)
	localUser.audioTrack.SetSendDelayMs(10)
	localUser.SetAudioScenario(senario)
	//update pcmsender's info
	localUser.audioTrack.pcmSender.audioScenario = senario
	//update agoraService's info
	agoraService.audioScenario = senario
	
	
	fmt.Printf("______update audio track id %d, old id %d, cTrack %p, old cTrack %p\n", localUser.audioTrack.id, localUser.audioTrack.id, localUser.audioTrack.cTrack, localUser.audioTrack.cTrack)
	
	
	return 0
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
func (localUser *LocalUser) SetAudioVolumeIndicationParameters(intervalInMs int, smooth int, reportVad bool) int {
	if localUser.cLocalUser == nil {
		return -1
	}
	
	ret := C.agora_local_user_set_audio_volume_indication_parameters(localUser.cLocalUser, C.int(intervalInMs), C.int(smooth), C.bool(reportVad))
	return int(ret)
}
func (localUser *LocalUser) SendAudioMetaData(metaData []byte) int {
	if localUser.cLocalUser == nil {
		return -1
	}
	cMetaData := C.CBytes(metaData)
	defer C.free(cMetaData)
	
	ret := C.agora_local_user_send_audio_meta_data(localUser.cLocalUser, (*C.char)(cMetaData), (C.size_t)(len(metaData)))
	return int(ret)
}
// add a api to interrupt the audio
func (localUser *LocalUser) InterruptAudio(audioConsumer *AudioConsumer) int {
	if localUser.cLocalUser == nil {
		return -1
	}
	if audioConsumer != nil {
		audioConsumer.Clear()
	}
	localUser.UnpublishAudio(localUser.audioTrack)
	
	// 根据不同的audioScenario, 做不同的处理
	if localUser.audioTrack.audioScenario == AudioScenarioAiServer {
		// 如果是aiServer, 则需要做打断操作
		
	} else {
		// 如果是chorus, 则需要做打断操作
		//localUser.audioTrack.ClearSenderBuffer()
	}
	return 0
}
