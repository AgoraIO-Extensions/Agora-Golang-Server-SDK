#pragma once

#include "agora_media_base.h"

// typedef struct _video_frame_observer2 {
//   void (*on_frame)(AGORA_HANDLE agora_video_frame_observer2, const char* channel_id, user_id_t remote_uid, const video_frame* frame);
// } video_frame_observer2;

extern int cgo_on_video_frame(AGORA_HANDLE agora_video_frame_observer2, const char* channelId, user_id_t uid, const video_frame* frame);