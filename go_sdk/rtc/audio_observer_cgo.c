#include "audio_observer_cgo.h"

extern int goOnRecordAudioFrame(void* agora_local_user, const char* channelId, const struct _audio_frame* frame);
int cgo_on_record_audio_frame(AGORA_HANDLE agora_local_user,const char* channelId, const audio_frame* frame) {
  return goOnRecordAudioFrame(agora_local_user, channelId, frame);
}

extern int goOnPlaybackAudioFrame(void* agora_local_user, const char* channelId, const struct _audio_frame* frame);
int cgo_on_playback_audio_frame(AGORA_HANDLE agora_local_user, const char* channelId, const audio_frame* frame) {
  return goOnPlaybackAudioFrame(agora_local_user, channelId, frame);
}

extern int goOnMixedAudioFrame(void* agora_local_user, const char* channelId, const struct _audio_frame* frame);
int cgo_on_mixed_audio_frame(AGORA_HANDLE agora_local_user, const char* channelId, const audio_frame* frame) {
  return goOnMixedAudioFrame(agora_local_user, channelId, frame);
}

extern int goOnEarMonitoringAudioFrame(void* agora_local_user, const struct _audio_frame* frame);
int cgo_on_ear_monitoring_audio_frame(AGORA_HANDLE agora_local_user, const audio_frame* frame) {
  return goOnEarMonitoringAudioFrame(agora_local_user, frame);
}

// this function declaration must be strictly same with the function exported by go
extern int goOnPlaybackAudioFrameBeforeMixing(void* agora_local_user, const char* channelId, const char* uid, const struct _audio_frame* frame);
int cgo_on_playback_audio_frame_before_mixing(AGORA_HANDLE agora_local_user, const char* channelId, user_id_t uid, const audio_frame* frame) {
  return goOnPlaybackAudioFrameBeforeMixing(agora_local_user, channelId, uid, frame);
}

extern int goOnGetAudioFramePosition(void* agora_local_user);
int cgo_on_get_audio_frame_position(AGORA_HANDLE agora_local_user) {
  return goOnGetAudioFramePosition(agora_local_user);
}

extern audio_params goOnGetPlaybackAudioFrameParam(void* agora_local_user);
audio_params cgo_on_get_playback_audio_frame_param(AGORA_HANDLE agora_local_user) {
  return goOnGetPlaybackAudioFrameParam(agora_local_user);
}

extern audio_params goOnGetRecordAudioFrameParam(void* agora_local_user);
audio_params cgo_on_get_record_audio_frame_param(AGORA_HANDLE agora_local_user) {
  return goOnGetRecordAudioFrameParam(agora_local_user);
}

extern audio_params goOnGetMixedAudioFrameParam(void* agora_local_user);
audio_params cgo_on_get_mixed_audio_frame_param(AGORA_HANDLE agora_local_user) {
  return goOnGetMixedAudioFrameParam(agora_local_user);
}

extern audio_params goOnGetEarMonitoringAudioFrameParam(void* agora_local_user);
audio_params cgo_on_get_ear_monitoring_audio_frame_param(AGORA_HANDLE agora_local_user) {
  return goOnGetEarMonitoringAudioFrameParam(agora_local_user);
}
