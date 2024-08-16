#include "video_observer_cgo.h"

// this function declaration must be strictly same with the function exported by go
extern int goOnVideoFrame(void* agora_video_frame_observer2, const char* channelId, const char* uid, const struct _video_frame* frame);
int cgo_on_video_frame(AGORA_HANDLE agora_video_frame_observer2, const char* channelId, user_id_t uid, const video_frame* frame) {
  return goOnVideoFrame(agora_video_frame_observer2, channelId, uid, frame);
}