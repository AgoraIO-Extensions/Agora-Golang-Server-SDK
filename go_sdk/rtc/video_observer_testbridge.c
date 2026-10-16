//go:build agora_video_observer_test

#include <stdlib.h>
#include <string.h>
#include "agora_local_user.h"
#include "agora_rtc_conn.h"

/* Explicit test double: never compiled into production builds. */
static int create_fail, set_result, unset_result;
static int counts[6]; /* create, destroy, set, unset, connection destroy, unsafe destroy */
static void* attached_observer;

void test_video_observer_reset(void) {
  create_fail = set_result = unset_result = 0;
  attached_observer = NULL;
  memset(counts, 0, sizeof(counts));
}
void test_video_observer_fail(int fail, int set, int unset) {
  create_fail = fail;
  set_result = set;
  unset_result = unset;
}
int test_video_observer_count(int operation) { return counts[operation]; }

AGORA_HANDLE agora_video_frame_observer_create(video_frame_observer* observer) {
  ++counts[0];
  if (create_fail) return NULL;
  video_frame_observer* copy = malloc(sizeof(*copy));
  if (copy) *copy = *observer;
  return copy;
}
void agora_video_frame_observer_destroy(AGORA_HANDLE observer) {
  ++counts[1];
  if (observer == attached_observer) ++counts[5];
  free(observer);
}
int agora_local_user_set_video_frame_observer(AGORA_HANDLE local_user, AGORA_HANDLE observer) {
  (void)local_user;
  ++counts[2];
  /* Real registration may query configuration before returning. */
  video_frame_observer* callbacks = observer;
  callbacks->get_video_frame_process_mode(observer);
  callbacks->get_video_format_preference(observer);
  callbacks->get_observed_frame_position(observer);
  if (set_result == 0) attached_observer = observer;
  return set_result;
}
int agora_local_user_unset_video_frame_observer(AGORA_HANDLE local_user, AGORA_HANDLE observer) {
  (void)local_user;
  (void)observer;
  ++counts[3];
  if (unset_result == 0 && attached_observer == observer) attached_observer = NULL;
  return unset_result;
}
void agora_rtc_conn_destroy(AGORA_HANDLE connection) {
  ++counts[4];
  attached_observer = NULL;
  free(connection);
}
int test_video_observer_query(void* observer, int slot) {
  video_frame_observer* callbacks = observer;
  switch (slot) {
    case 0: return callbacks->get_video_frame_process_mode(observer);
    case 1: return callbacks->get_video_format_preference(observer);
    case 2: return callbacks->get_rotation_applied(observer);
    case 3: return callbacks->get_mirror_applied(observer);
    case 4: return callbacks->get_observed_frame_position(observer);
  }
  return -99;
}
video_frame* test_video_frame_create(int format) {
  /* Odd height exercises chroma-plane rounding; padding detects over-copying. */
  video_frame* frame = calloc(1, sizeof(*frame) + 128);
  uint8_t* storage = (uint8_t*)(frame + 1);
  frame->type = format;
  frame->width = 2;
  frame->height = 3;
  frame->y_stride = 4;
  frame->u_stride = frame->v_stride = 2;
  frame->y_buffer = storage;
  frame->u_buffer = storage + 32;
  frame->v_buffer = storage + 64;
  if (format == 2 || format == 4) {
    frame->y_stride = 24; /* packed RGBA/BGRA uses the total data length */
    frame->u_stride = frame->v_stride = 0;
    frame->u_buffer = frame->v_buffer = NULL;
  }
  if (format == 3 || format == 8) {
    frame->v_stride = 0;
    frame->v_buffer = NULL;
  }
  frame->metadata_buffer = storage + 96;
  frame->metadata_size = 2;
  frame->alpha_buffer = storage + 100;
  frame->alpha_mode = 2;
  memset(storage, 7, 128);
  return frame;
}
int test_video_observer_frame(void* observer, int slot, unsigned int uid, video_frame* frame) {
  video_frame_observer* callbacks = observer;
  switch (slot) {
    case 0: return callbacks->on_capture_video_frame(observer, 4, frame);
    case 1: return callbacks->on_pre_encode_video_frame(observer, 4, frame);
    case 2: return callbacks->on_media_player_video_frame(observer, frame, 42);
    case 3: return callbacks->on_render_video_frame(observer, "test-channel", uid, frame);
    case 4: return callbacks->on_transcoded_video_frame(observer, frame);
  }
  return -99;
}
