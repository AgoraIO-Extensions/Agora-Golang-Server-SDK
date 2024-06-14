#pragma once

#include "agora_media_base.h"

// typedef struct _audio_frame_observer {
//   /* return value stands for a 'bool' in C++: 1 for success, 0 for failure */
//   int (*on_record_audio_frame)(AGORA_HANDLE agora_local_user /* raw pointer */,const char* channelId, const audio_frame* frame);
//   int (*on_playback_audio_frame)(AGORA_HANDLE agora_local_user, const char* channelId, const audio_frame* frame);
//   int (*on_mixed_audio_frame)(AGORA_HANDLE agora_local_user, const char* channelId, const audio_frame* frame);
//   int (*on_ear_monitoring_audio_frame)(AGORA_HANDLE agora_local_user, const audio_frame* frame);
//   int (*on_playback_audio_frame_before_mixing)(AGORA_HANDLE agora_local_user, const char* channelId, user_id_t uid, const audio_frame* frame);
//   int (*on_get_audio_frame_position)(AGORA_HANDLE agora_local_user);
//   audio_params (*on_get_playback_audio_frame_param)(AGORA_HANDLE agora_local_user);
//   audio_params (*on_get_record_audio_frame_param)(AGORA_HANDLE agora_local_user);
//   audio_params (*on_get_mixed_audio_frame_param)(AGORA_HANDLE agora_local_user);
//   audio_params (*on_get_ear_monitoring_audio_frame_param)(AGORA_HANDLE agora_local_user);
// } audio_frame_observer;

extern int cgo_on_playback_audio_frame_before_mixing(AGORA_HANDLE agora_local_user, const char* channelId, user_id_t uid, const audio_frame* frame);