#pragma once
#include <stdint.h>
#include <libavformat/avformat.h>
#include <libavutil/error.h>

#define AVERROR_EAGAIN AVERROR(EAGAIN)

typedef struct _MediaFrame {
  // common fields
  int stream_index;
  // AVMEDIA_TYPE_VIDEO or AVMEDIA_TYPE_AUDIO
  int frame_type;
  // time in ms
  int64_t pts;
  uint8_t *buffer;
  int buffer_size;
  // video pixel format or audio sample format
  int format;

  // video fields
  int width;
  int height;
  int stride;

  // audio fields
  int samples;
  int channels;
  int sample_rate;
  int bytes_per_sample;
}MediaFrame;

typedef struct _MediaPacket {
  AVPacket *pkt;
  int media_type;
  int64_t pts;

  int width;
  int height;
  int framerate_num;
  int framerate_den;

  struct _MediaPacket *next;
}MediaPacket;

// open file for decode or demux, return decoder handle
extern void * open_media_file(const char *file_name);
// demux and alloc packet
extern int get_packet(void *decoder, MediaPacket **packet);
// free packet
extern int free_packet(MediaPacket **packet);
// transcode h264 packet to annexb format
extern int h264_to_annexb(void *decoder, MediaPacket **packet);
// decode a demuxed packet to I420 YUV or 16khz 1channel pcm frame
extern int decode_packet(void *decoder, MediaPacket *packet, MediaFrame *frame);
// decode frame directly from file, this is combine of get_packet and decode_packet
// NOTICE: this function must not be used when you are using get_packet and decode_packet function
extern int get_frame(void *decoder, MediaFrame *frame);
// close file and free decoder
extern void close_media_file(void *decoder);