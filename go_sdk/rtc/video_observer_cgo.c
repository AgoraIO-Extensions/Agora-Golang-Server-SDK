#include <stdint.h>
#include "video_observer_cgo.h"

// Use the generated declarations so Go/C callback signatures cannot drift.
#include "_cgo_export.h"

int cgo_on_capture_video_frame(AGORA_HANDLE observer, int source_type, video_frame* frame) {
  return goOnCaptureVideoFrame(observer, source_type, frame);
}
int cgo_on_pre_encode_video_frame(AGORA_HANDLE observer, int source_type, video_frame* frame) {
  return goOnPreEncodeVideoFrame(observer, source_type, frame);
}
int cgo_on_media_player_video_frame(AGORA_HANDLE observer, video_frame* frame, int media_player_id) {
  return goOnMediaPlayerVideoFrame(observer, frame, media_player_id);
}
int cgo_on_render_video_frame(AGORA_HANDLE observer, const char* channel_id, uid_t uid, video_frame* frame) {
  return goOnRenderVideoFrame(observer, (char*)channel_id, uid, frame);
}
int cgo_on_transcoded_video_frame(AGORA_HANDLE observer, video_frame* frame) {
  return goOnTranscodedVideoFrame(observer, frame);
}
int cgo_get_video_frame_process_mode(AGORA_HANDLE observer) {
  return goGetVideoFrameProcessMode(observer);
}
int cgo_get_video_format_preference(AGORA_HANDLE observer) {
  return goGetVideoFormatPreference(observer);
}
int cgo_get_rotation_applied(AGORA_HANDLE observer) {
  return goGetRotationApplied(observer);
}
int cgo_get_mirror_applied(AGORA_HANDLE observer) {
  return goGetMirrorApplied(observer);
}
uint32_t cgo_get_observed_frame_position(AGORA_HANDLE observer) {
  return goGetObservedFramePosition(observer);
}

int cgo_on_encoded_video_frame(AGORA_HANDLE agora_video_encoded_frame_observer, const char* channelId, const char* uid, const uint8_t* image_buffer, size_t length,
                                const encoded_video_frame_info* video_encoded_frame_info) {
  return goOnEncodedVideoFrame(agora_video_encoded_frame_observer, (char*)channelId, (char*)uid,
                               (uint8_t*)image_buffer, length, (encoded_video_frame_info*)video_encoded_frame_info);
}
