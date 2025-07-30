#include "rtc_callbacks_cgo.h"

extern void goOnConnected(void* agora_rtc_conn, struct _rtc_conn_info* conn_info, int reason);
void cgo_on_connected(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info, int reason) {
  goOnConnected(agora_rtc_conn, (struct _rtc_conn_info*)conn_info, reason);
}

extern void goOnDisconnected(void* agora_rtc_conn, struct _rtc_conn_info* conn_info, int reason);
void cgo_on_disconnected(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info, int reason) {
  goOnDisconnected(agora_rtc_conn, (struct _rtc_conn_info*)conn_info, reason);
}

extern void goOnConnecting(void* agora_rtc_conn, struct _rtc_conn_info* conn_info, int reason);
void cgo_on_connecting(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info, int reason) {
  goOnConnecting(agora_rtc_conn, (struct _rtc_conn_info*)conn_info, reason);
}

extern void goOnReconnecting(void* agora_rtc_conn, struct _rtc_conn_info* conn_info, int reason);
void cgo_on_reconnecting(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info, int reason) {
  goOnReconnecting(agora_rtc_conn, (struct _rtc_conn_info*)conn_info, reason);
}

extern void goOnReconnected(void* agora_rtc_conn, struct _rtc_conn_info* conn_info, int reason);
void cgo_on_reconnected(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info, int reason)  {
  goOnReconnected(agora_rtc_conn, (struct _rtc_conn_info*)conn_info, reason);
}

extern void goOnConnectionLost(void* agora_rtc_conn, struct _rtc_conn_info* conn_info);
void cgo_on_connection_lost(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info) {
  goOnConnectionLost(agora_rtc_conn, (struct _rtc_conn_info*)conn_info);
}

extern void goOnConnectionFailure(AGORA_HANDLE agora_rtc_conn, struct _rtc_conn_info* conn_info, int reason);
void cgo_on_connection_failure(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info, int reason) {
  goOnConnectionFailure(agora_rtc_conn, (struct _rtc_conn_info*)conn_info, reason);
}

//token
extern void goOnTokenPrivilegeWillExpire(void* agora_rtc_conn, const char* token);
void cgo_on_token_privilege_will_expire(AGORA_HANDLE agora_rtc_conn, const char* token) {
  goOnTokenPrivilegeWillExpire(agora_rtc_conn, token);
}

extern void goOnTokenPrivilegeDidExpire(void* agora_rtc_conn);
void cgo_on_token_privilege_did_expire(AGORA_HANDLE agora_rtc_conn) {
  goOnTokenPrivilegeDidExpire(agora_rtc_conn);
}

//user state
extern void goOnUserJoined(void* agora_rtc_conn, user_id_t user_id);
void cgo_on_user_joined(AGORA_HANDLE agora_rtc_conn, user_id_t user_id) {
  goOnUserJoined(agora_rtc_conn, user_id);
}

extern void goOnUserOffline(void* agora_rtc_conn, user_id_t user_id, int reason);
void cgo_on_user_left(AGORA_HANDLE agora_rtc_conn, user_id_t user_id, int reason) {
  goOnUserOffline(agora_rtc_conn, user_id, reason);
}

extern void goOnError(void* agora_rtc_conn, int error, const char* msg);
void cgo_on_error(AGORA_HANDLE agora_rtc_conn, int error, const char* msg) {
  goOnError(agora_rtc_conn, error, msg);
}

//steam message 
extern void goOnStreamMessageError(void* agora_rtc_conn, user_id_t user_id, int stream_id, int code, int missed, int cached);
void cgo_on_stream_message_error(AGORA_HANDLE agora_rtc_conn, user_id_t user_id, int stream_id, int code, int missed, int cached) {
  goOnStreamMessageError(agora_rtc_conn, user_id, stream_id, code, missed, cached);
}

extern void goOnStreamMessage(void* agora_local_user, user_id_t user_id, int stream_id, const char* data, size_t length);
void cgo_on_stream_message(AGORA_HANDLE agora_local_user, user_id_t user_id, int stream_id, const char* data, size_t length) {
  goOnStreamMessage(agora_local_user, user_id, stream_id, data, length);
}

extern void goOnUserInfoUpdated(void* agora_local_user, user_id_t user_id, int msg, int val);
void cgo_on_user_info_updated(AGORA_HANDLE agora_local_user, user_id_t user_id, int msg, int val) {
  goOnUserInfoUpdated(agora_local_user, user_id, msg, val);
}

extern void goOnUserAudioTrackSubscribed(void* agora_local_user, user_id_t user_id, void* agora_remote_audio_track);
void cgo_on_user_audio_track_subscribed(AGORA_HANDLE agora_local_user, user_id_t user_id, AGORA_HANDLE agora_remote_audio_track) {
  goOnUserAudioTrackSubscribed(agora_local_user, user_id, agora_remote_audio_track);
}

extern void goOnUserVideoTrackSubscribed(void* agora_local_user, user_id_t user_id, struct _video_track_info* info, void* agora_remote_video_track);
void cgo_on_user_video_track_subscribed(AGORA_HANDLE agora_local_user, user_id_t user_id, const video_track_info* info, AGORA_HANDLE agora_remote_video_track) {
  goOnUserVideoTrackSubscribed(agora_local_user, user_id, (struct _video_track_info*)info, agora_remote_video_track);
}

extern void goOnUserAudioTrackStateChanged(void* agora_local_user, user_id_t user_id, void* agora_remote_audio_track, int state, int reason, int elapsed);
void cgo_on_user_audio_track_state_changed(AGORA_HANDLE agora_local_user, user_id_t user_id, AGORA_HANDLE agora_remote_audio_track, int state, int reason, int elapsed) {
  goOnUserAudioTrackStateChanged(agora_local_user, user_id, agora_remote_audio_track, state, reason, elapsed);
}

extern void goOnUserVideoTrackStateChanged(void* agora_local_user, user_id_t user_id, void* agora_remote_video_track, int state, int reason, int elapsed);
void cgo_on_user_video_track_state_changed(AGORA_HANDLE agora_local_user, user_id_t user_id, AGORA_HANDLE agora_remote_video_track, int state, int reason, int elapsed) {
  goOnUserVideoTrackStateChanged(agora_local_user, user_id, agora_remote_video_track, state, reason, elapsed);
}

extern void goOnAudioPublishStateChanged(void* agora_local_user, const char* channelid, int oldstate, int newstate, int elapseSinceLastState);
void  cgo_on_audio_publish_state_changed(AGORA_HANDLE agora_local_user, const char* channelid, int oldstate, int newstate, int elapseSinceLastState) {
  goOnAudioPublishStateChanged(agora_local_user, channelid, oldstate, newstate, elapseSinceLastState);
  
}

extern void goOnAudioVolumeIndication(void* agora_local_user, struct _audio_volume_info* speakers, unsigned int speaker_number, int total_volume);
void cgo_on_audio_volume_indication(AGORA_HANDLE agora_local_user, const audio_volume_info* speakers, unsigned int speaker_number, int total_volume) {
  goOnAudioVolumeIndication(agora_local_user, (struct _audio_volume_info*)speakers, speaker_number, total_volume);
}
extern void goOnAudioMetadataReceived(void* agora_local_user, user_id_t userId, const char* meta_data, size_t length);
void cgo_on_audio_meta_data_received(AGORA_HANDLE agora_local_user, user_id_t userId, const char* meta_data, size_t length)
{
  goOnAudioMetadataReceived(agora_local_user, userId, meta_data, length);
}

extern void goOnLocalAudioTrackStatistics(void* agora_local_user, struct _local_audio_stats* stats);
void cgo_on_local_audio_track_statistics(AGORA_HANDLE agora_local_user, const local_audio_stats* stats) {
  goOnLocalAudioTrackStatistics(agora_local_user, (struct _local_audio_stats*)stats);
}
extern void goOnRemoteAudioTrackStatistics(void* agora_local_user, user_id_t userId, struct _remote_audio_stats* stats);
void cgo_on_remote_audio_track_statistics(AGORA_HANDLE agora_local_user, user_id_t userId, const remote_audio_stats* stats) {
  goOnRemoteAudioTrackStatistics(agora_local_user, userId, (struct _remote_audio_stats*)stats);
}


extern void goOnLocalVideoTrackStatistics(void* agora_local_user, struct _local_video_track_stats* stats);
void cgo_on_local_video_track_statistics(AGORA_HANDLE agora_local_user, const local_video_track_stats* stats) {
  goOnLocalVideoTrackStatistics(agora_local_user, (struct _local_video_track_stats*)stats);
}


extern void goOnRemoteVideoTrackStatistics(void* agora_local_user, user_id_t userId, struct _remote_video_track_stats* stats);
void cgo_on_remote_video_track_statistics(AGORA_HANDLE agora_local_user, user_id_t userId, const remote_video_track_stats* stats) {
  goOnRemoteVideoTrackStatistics(agora_local_user, userId, (struct _remote_video_track_stats*)stats);
}

extern void goOnEncryptionError(void* agora_rtc_conn, int error_type);
void cgo_on_encryption_error(AGORA_HANDLE agora_rtc_conn, int error_type) {
  goOnEncryptionError(agora_rtc_conn, error_type);
}
extern void goOnAudioTrackPublishSuccess(void* agora_local_user, void* agora_local_audio_track);
void cgo_on_audio_track_publish_success(AGORA_HANDLE agora_local_user, AGORA_HANDLE agora_local_audio_track) {
  goOnAudioTrackPublishSuccess(agora_local_user, agora_local_audio_track);
}
extern void goOnAudioTrackUnpublished(void* agora_local_user, void* agora_local_audio_track);
void cgo_on_audio_track_unpublished(AGORA_HANDLE agora_local_user, AGORA_HANDLE agora_local_audio_track) {
  goOnAudioTrackUnpublished(agora_local_user, agora_local_audio_track);
}

extern void goOnCapabilitiesChanged(void* agora_local_user, struct  _capabilities* caps, int size);
void cgo_on_capabilities_changed(AGORA_HANDLE agora_local_user, const capabilities* caps, int size) {
  goOnCapabilitiesChanged(agora_local_user, (struct _capabilities*)caps, size);
}