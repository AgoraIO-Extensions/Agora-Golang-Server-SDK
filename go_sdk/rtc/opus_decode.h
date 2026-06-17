#pragma once
#include <stdint.h>

// Minimal media frame description, self-contained for the opus decoder.
typedef struct _MediaFrame {
  int stream_index;
  int frame_type;
  int64_t pts;          // time in ms
  uint8_t *buffer;      // pcm data (interleaved S16 for audio)
  int buffer_size;      // valid bytes in buffer
  int format;           // sample format (S16 = 1)

  // video fields (unused for opus)
  int width;
  int height;
  int stride;

  // audio fields
  int samples;          // samples per channel
  int channels;
  int sample_rate;
  int bytes_per_sample;
} MediaFrame;

// Standalone Opus decoder for raw opus packets (no container/demuxer).
// Use this for bare opus payloads, e.g. Agora's OnEncodedAudioFrameReceived packet.
//
// in_channels:     channel count of the opus stream (1 or 2). pass 0 to default to 1.
// out_sample_rate: desired output PCM sample rate, e.g. 16000/48000. pass 0 to keep 48000.
// out_channels:    desired output PCM channel count. pass 0 to keep in_channels.
// returns an opaque handle, or NULL on error.
extern void * open_opus_decoder(int in_channels, int out_sample_rate, int out_channels);

// decode one raw opus packet into a PCM S16 MediaFrame.
// returns 0 on success (frame filled), AVERROR(EAGAIN) when more data is needed,
// or a negative error code on failure.
extern int decode_opus(void *handle, const uint8_t *data, int data_size, MediaFrame *frame);

// close opus decoder and free resources.
extern void close_opus_decoder(void *handle);
