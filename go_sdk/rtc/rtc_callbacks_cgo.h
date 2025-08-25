#pragma once

#include "agora_rtc_conn.h"
#include "agora_local_user.h"
// typedef struct _rtc_conn_observer {
//   void (*on_connected)(AGORA_HANDLE agora_rtc_conn /* pointer to RefPtrHolder */, const rtc_conn_info* conn_info, int reason);
//   void (*on_disconnected)(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info, int reason);
//   void (*on_connecting)(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info, int reason);
//   void (*on_reconnecting)(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info, int reason);
//   void (*on_reconnected)(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info, int reason);
//   void (*on_connection_lost)(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info);

//   void (*on_lastmile_quality)(AGORA_HANDLE agora_rtc_conn, int quality);
//   void (*on_lastmile_probe_result)(AGORA_HANDLE agora_rtc_conn, const lastmile_probe_result* result);
//   void (*on_token_privilege_will_expire)(AGORA_HANDLE agora_rtc_conn, const char* token);
//   void (*on_token_privilege_did_expire)(AGORA_HANDLE agora_rtc_conn);
//   void (*on_connection_license_validation_failure)(AGORA_HANDLE agora_rtc_conn, int reason);
//   void (*on_connection_failure)(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info, int reason); 
//   void (*on_user_joined)(AGORA_HANDLE agora_rtc_conn, user_id_t user_id);
//   void (*on_user_left)(AGORA_HANDLE agora_rtc_conn, user_id_t user_id, int reason);
//   void (*on_transport_stats)(AGORA_HANDLE agora_rtc_conn, const rtc_stats* stats);
//   void (*on_change_role_success)(AGORA_HANDLE agora_rtc_conn, int old_role, int new_role);
//   void (*on_change_role_failure)(AGORA_HANDLE agora_rtc_conn, int reason, int current_role);
//   void (*on_user_network_quality)(AGORA_HANDLE agora_rtc_conn, user_id_t user_id, int tx_quality, int rx_quality);
//   void (*on_network_type_changed)(AGORA_HANDLE agora_rtc_conn, int type);
//   void (*on_api_call_executed)(AGORA_HANDLE agora_rtc_conn, int err, const char* api, const char* result);
//   void (*on_content_inspect_result)(AGORA_HANDLE agora_rtc_conn, int result);
//   void (*on_snapshot_taken)(AGORA_HANDLE agora_rtc_conn, const char* channel, uid_t uid, const char* file_path, int width, int height, int err_code);
//   void (*on_error)(AGORA_HANDLE agora_rtc_conn, int error, const char* msg);
//   void (*on_warning)(AGORA_HANDLE agora_rtc_conn, int warning, const char* msg);
//   void (*on_channel_media_relay_state_changed)(AGORA_HANDLE agora_rtc_conn, int state, int code);
//   void (*on_local_user_registered)(AGORA_HANDLE agora_rtc_conn, uid_t uid, const char* user_account);
//   void (*on_user_account_updated)(AGORA_HANDLE agora_rtc_conn, uid_t uid, const char* user_account);
//   void (*on_stream_message_error)(AGORA_HANDLE agora_rtc_conn, user_id_t user_id, int stream_id, int code, int missed, int cached);
//   void (*on_encryption_error)(AGORA_HANDLE agora_rtc_conn, int error_type);
//   void (*on_upload_log_result)(AGORA_HANDLE agora_rtc_conn, const char* request_id, int success, int reason);
// } rtc_conn_observer;
// sonme fucc defined in local user observer 
//  void (*on_audio_volume_indication)(AGORA_HANDLE agora_local_user, const audio_volume_info* speakers, unsigned int speaker_number, int total_volume); 

//connection
extern void cgo_on_connected(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info, int reason);
extern void cgo_on_disconnected(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info, int reason);
extern void cgo_on_connecting(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info, int reason);
extern void cgo_on_reconnecting(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info, int reason);
extern void cgo_on_reconnected(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info, int reason);
extern void cgo_on_connection_lost(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info);
extern void cgo_on_connection_failure(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info, int reason);

//token
extern void cgo_on_token_privilege_will_expire(AGORA_HANDLE agora_rtc_conn, const char* token);
extern void cgo_on_token_privilege_did_expire(AGORA_HANDLE agora_rtc_conn);

//user state
extern void cgo_on_user_joined(AGORA_HANDLE agora_rtc_conn, user_id_t user_id);
extern void cgo_on_user_left(AGORA_HANDLE agora_rtc_conn, user_id_t user_id, int reason);

//steam message 
extern void cgo_on_stream_message_error(AGORA_HANDLE agora_rtc_conn, user_id_t user_id, int stream_id, int code, int missed, int cached);
extern void cgo_on_stream_message(AGORA_HANDLE agora_local_user, user_id_t user_id, int stream_id, const char* data, size_t length);

  // /**
  //  * The media information of a specified user.
  //  */
  // enum USER_MEDIA_INFO {
  //   /**
  //    * 0: The user has muted the audio.
  //    */
  //   USER_MEDIA_INFO_MUTE_AUDIO = 0,
  //   /**
  //    * 1: The user has muted the video.
  //    */
  //   USER_MEDIA_INFO_MUTE_VIDEO = 1,
  //   /**
  //    * 4: The user has enabled the video, which includes video capturing and encoding.
  //    */
  //   USER_MEDIA_INFO_ENABLE_VIDEO = 4,
  //   /**
  //    * 8: The user has enabled the local video capturing.
  //    */
  //   USER_MEDIA_INFO_ENABLE_LOCAL_VIDEO = 8,
  // };
extern void cgo_on_user_info_updated(AGORA_HANDLE agora_local_user, user_id_t user_id, int msg, int val);

extern void cgo_on_user_audio_track_subscribed(AGORA_HANDLE agora_local_user, user_id_t user_id, AGORA_HANDLE agora_remote_audio_track);

extern void cgo_on_user_video_track_subscribed(AGORA_HANDLE agora_local_user, user_id_t user_id, const video_track_info* info, AGORA_HANDLE agora_remote_video_track);

extern void cgo_on_user_audio_track_state_changed(AGORA_HANDLE agora_local_user, user_id_t user_id, AGORA_HANDLE agora_remote_audio_track, int state, int reason, int elapsed);

extern void cgo_on_user_video_track_state_changed(AGORA_HANDLE agora_local_user, user_id_t user_id, AGORA_HANDLE agora_remote_video_track, int state, int reason, int elapsed);

extern void cgo_on_error(AGORA_HANDLE agora_rtc_conn, int error, const char* msg);
extern void cgo_on_audio_volume_indication(AGORA_HANDLE agora_local_user, const audio_volume_info* speakers, unsigned int speaker_number, int total_volume);
extern void cgo_on_audio_publish_state_changed(AGORA_HANDLE agora_local_user, const char* channelid, int oldstate, int newstate, int elapseSinceLastState);
extern void cgo_on_audio_volume_indication(AGORA_HANDLE agora_local_user, const audio_volume_info* speakers, unsigned int speaker_number, int total_volume);
extern void cgo_on_audio_meta_data_received(AGORA_HANDLE agora_local_user, user_id_t userId, const char* meta_data, size_t length);
extern void cgo_on_local_audio_track_statistics(AGORA_HANDLE agora_local_user, const local_audio_stats* stats);
extern void cgo_on_remote_audio_track_statistics(AGORA_HANDLE agora_local_user, user_id_t userId, const remote_audio_stats* stats);
extern void cgo_on_local_video_track_statistics(AGORA_HANDLE agora_local_user, const local_video_track_stats* stats);
extern void cgo_on_remote_video_track_statistics(AGORA_HANDLE agora_local_user, user_id_t userId, const remote_video_track_stats* stats);
extern void cgo_on_encryption_error(AGORA_HANDLE agora_rtc_conn, int error_type);

// added on 2025-06-09
extern void cgo_on_audio_track_publish_success(AGORA_HANDLE agora_local_user, AGORA_HANDLE agora_local_audio_track );
extern void cgo_on_audio_track_unpublished(AGORA_HANDLE agora_local_user, AGORA_HANDLE agora_local_audio_track);

// added on 2025-06-13
extern void cgo_on_capabilities_changed(AGORA_HANDLE agora_local_user, const capabilities* caps, int size);