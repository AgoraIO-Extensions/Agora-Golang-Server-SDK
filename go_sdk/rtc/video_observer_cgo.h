#pragma once

#include <stdint.h>
#include <stddef.h>
#include "agora_base.h"
#include "agora_media_base.h"

extern int cgo_on_capture_video_frame(AGORA_HANDLE observer, int source_type, video_frame* frame);
extern int cgo_on_pre_encode_video_frame(AGORA_HANDLE observer, int source_type, video_frame* frame);
extern int cgo_on_media_player_video_frame(AGORA_HANDLE observer, video_frame* frame, int media_player_id);
extern int cgo_on_render_video_frame(AGORA_HANDLE observer, const char* channel_id, uid_t uid, video_frame* frame);
extern int cgo_on_transcoded_video_frame(AGORA_HANDLE observer, video_frame* frame);
extern int cgo_get_video_frame_process_mode(AGORA_HANDLE observer);
extern int cgo_get_video_format_preference(AGORA_HANDLE observer);
extern int cgo_get_rotation_applied(AGORA_HANDLE observer);
extern int cgo_get_mirror_applied(AGORA_HANDLE observer);
extern uint32_t cgo_get_observed_frame_position(AGORA_HANDLE observer);

extern int cgo_on_encoded_video_frame(AGORA_HANDLE agora_video_encoded_frame_observer, const char* channelId, const char* uid, const uint8_t* image_buffer, size_t length,
                                const encoded_video_frame_info* video_encoded_frame_info);
