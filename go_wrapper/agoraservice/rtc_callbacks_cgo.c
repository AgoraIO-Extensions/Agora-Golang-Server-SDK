#include "rtc_callbacks_cgo.h"

extern void goOnConnected(void* agora_rtc_conn, struct _rtc_conn_info* conn_info, int reason);
void cgo_on_connected(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info, int reason) {
  goOnConnected(agora_rtc_conn, conn_info, reason);
}