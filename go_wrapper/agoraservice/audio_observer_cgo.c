#include "audio_observer_cgo.h"

// this function declaration must be strictly same with the function exported by go
extern int goOnPlaybackAudioFrameBeforeMixing(void* agora_local_user, const char* channelId, const char* uid, const struct _audio_frame* frame);
int cgo_on_playback_audio_frame_before_mixing(AGORA_HANDLE agora_local_user, const char* channelId, user_id_t uid, const audio_frame* frame) {
  return goOnPlaybackAudioFrameBeforeMixing(agora_local_user, channelId, uid, frame);
}