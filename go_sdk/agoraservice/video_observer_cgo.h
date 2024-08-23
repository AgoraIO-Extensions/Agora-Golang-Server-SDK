#pragma once

#include <stdint.h>
#include <stddef.h>
#include "agora_base.h"
#include "agora_media_base.h"

// typedef struct _video_frame_observer2 {
//   void (*on_frame)(AGORA_HANDLE agora_video_frame_observer2, const char* channel_id, user_id_t remote_uid, const video_frame* frame);
// } video_frame_observer2;

extern int cgo_on_video_frame(AGORA_HANDLE agora_video_frame_observer2, const char* channelId, user_id_t uid, const video_frame* frame);

extern int cgo_on_encoded_video_frame(AGORA_HANDLE agora_video_encoded_frame_observer, uid_t uid, const uint8_t* image_buffer, size_t length,
                                const encoded_video_frame_info* video_encoded_frame_info);