#include "rtc_callbacks_cgo.h"

extern void goOnConnected(void* agora_rtc_conn, struct _rtc_conn_info* conn_info, int reason);
void cgo_on_connected(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info, int reason) {
  goOnConnected(agora_rtc_conn, conn_info, reason);
}

extern void goOnDisconnected(void* agora_rtc_conn, struct _rtc_conn_info* conn_info, int reason);
void cgo_on_disconnected(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info, int reason) {
  goOnDisconnected(agora_rtc_conn, conn_info, reason);
}

extern void goOnReconnecting(void* agora_rtc_conn, struct _rtc_conn_info* conn_info, int reason);
void cgo_on_reconnecting(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info, int reason) {
  goOnReconnecting(agora_rtc_conn, conn_info, reason);
}

extern void goOnReconnected(void* agora_rtc_conn, struct _rtc_conn_info* conn_info, int reason);
void cgo_on_reconnected(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info, int reason)  {
  goOnReconnected(agora_rtc_conn, conn_info, reason);
}

extern void goOnConnectionLost(void* agora_rtc_conn, struct _rtc_conn_info* conn_info);
void cgo_on_connection_lost(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info) {
  goOnConnectionLost(agora_rtc_conn, conn_info);
}

extern void goOnConnectionFailure(AGORA_HANDLE agora_rtc_conn, struct _rtc_conn_info* conn_info, int reason);
void cgo_on_connection_failure(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info, int reason) {
  goOnConnectionFailure(agora_rtc_conn, conn_info, reason);
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
