#include <stdint.h>
#include "video_observer_cgo.h"

// this function declaration must be strictly same with the function exported by go
extern int goOnVideoFrame(void* agora_video_frame_observer2, const char* channelId, const char* uid, const struct _video_frame* frame);
int cgo_on_video_frame(AGORA_HANDLE agora_video_frame_observer2, const char* channelId, user_id_t uid, const video_frame* frame) {
  return goOnVideoFrame(agora_video_frame_observer2, channelId, uid, frame);
}

extern int goOnEncodedVideoFrame(void* agora_video_encoded_frame_observer, uint32_t uid, const uint8_t* image_buffer, size_t length,
                                const struct _encoded_video_frame_info* video_encoded_frame_info);
int cgo_on_encoded_video_frame(AGORA_HANDLE agora_video_encoded_frame_observer, uint32_t uid, const uint8_t* image_buffer, size_t length,
                                const encoded_video_frame_info* video_encoded_frame_info) {
  return goOnEncodedVideoFrame(agora_video_encoded_frame_observer, uid, image_buffer, length, video_encoded_frame_info);
}