
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#include "opus_decode.h"


#ifdef USE_AVCODEC
// standalone opus decoder for raw opus packets, self-contained.
#include <libavcodec/avcodec.h>
#include <libavutil/opt.h>
#include <libavutil/channel_layout.h>
#include <libavutil/samplefmt.h>
#include <libavutil/avutil.h>
#include <libavutil/error.h>
#include <libswresample/swresample.h>

#define OPUS_MAX_CHANNELS 8

#if LIBAVCODEC_VERSION_INT >= AV_VERSION_INT(59, 37, 100)
#define OPUS_USE_NEW_CH_LAYOUT 1
#else
#define OPUS_USE_NEW_CH_LAYOUT 0
#endif

typedef struct _OpusDecoder {
  AVCodecContext *codec_ctx;
  AVFrame *frame;

  // resample to interleaved S16
  struct SwrContext *swr_ctx;
#if OPUS_USE_NEW_CH_LAYOUT
  AVChannelLayout dst_ch_layout;
#else
  uint64_t dst_channel_layout;
  int dst_channels;
#endif
  enum AVSampleFormat dst_sample_fmt;
  int dst_sample_rate;
  int dst_nb_samples;
  int swr_inited;

  // output pcm buffer
  uint8_t *buffer;
  int buffer_size;        // allocated size
  int actual_buffer_size; // valid bytes for current frame
  uint8_t *samples[OPUS_MAX_CHANNELS];

  // requested output params
  int out_sample_rate;
  int out_channels;
} OpusDecoder2;

static int opus_ctx_channels(const AVCodecContext *ctx) {
#if OPUS_USE_NEW_CH_LAYOUT
  return ctx->ch_layout.nb_channels;
#else
  return ctx->channels;
#endif
}

static int opus_dst_channels(const OpusDecoder2 *oc) {
#if OPUS_USE_NEW_CH_LAYOUT
  return oc->dst_ch_layout.nb_channels;
#else
  return oc->dst_channels;
#endif
}

static void opus_set_ctx_channels(AVCodecContext *ctx, int channels) {
#if OPUS_USE_NEW_CH_LAYOUT
  av_channel_layout_default(&ctx->ch_layout, channels);
#else
  ctx->channels = channels;
  ctx->channel_layout = av_get_default_channel_layout(channels);
#endif
}

static void opus_set_dst_channels(OpusDecoder2 *oc, int channels) {
#if OPUS_USE_NEW_CH_LAYOUT
  av_channel_layout_default(&oc->dst_ch_layout, channels);
#else
  oc->dst_channels = channels;
  oc->dst_channel_layout = av_get_default_channel_layout(channels);
#endif
}

static int opus_set_swr_chlayout(struct SwrContext *swr_ctx,
                                 AVCodecContext *ctx,
                                 OpusDecoder2 *oc) {
#if OPUS_USE_NEW_CH_LAYOUT
  av_opt_set_chlayout(swr_ctx, "in_chlayout", &ctx->ch_layout, 0);
  av_opt_set_chlayout(swr_ctx, "out_chlayout", &oc->dst_ch_layout, 0);
#else
  av_opt_set_int(swr_ctx, "in_channel_layout", ctx->channel_layout, 0);
  av_opt_set_int(swr_ctx, "out_channel_layout", oc->dst_channel_layout, 0);
#endif
  return 0;
}

// init the resampler lazily after the first frame is decoded,
// because the real sample_fmt is only known once decoding starts.
static int opus_init_swr(OpusDecoder2 *oc) {
  AVCodecContext *ctx = oc->codec_ctx;

  oc->dst_sample_fmt = AV_SAMPLE_FMT_S16;
  oc->dst_sample_rate = oc->out_sample_rate > 0 ? oc->out_sample_rate : ctx->sample_rate;
  int out_ch = oc->out_channels > 0 ? oc->out_channels : opus_ctx_channels(ctx);
  opus_set_dst_channels(oc, out_ch);

  oc->swr_inited = 1;

  // no conversion needed
  if (ctx->sample_fmt == oc->dst_sample_fmt &&
      ctx->sample_rate == oc->dst_sample_rate &&
      opus_ctx_channels(ctx) == opus_dst_channels(oc)) {
    return 0;
  }

  struct SwrContext *swr_ctx = swr_alloc();
  if (!swr_ctx) {
    av_log(NULL, AV_LOG_ERROR, "Can't allocate opus swr context\n");
    return AVERROR(ENOMEM);
  }
  opus_set_swr_chlayout(swr_ctx, ctx, oc);
  av_opt_set_int(swr_ctx, "in_sample_rate", ctx->sample_rate, 0);
  av_opt_set_int(swr_ctx, "out_sample_rate", oc->dst_sample_rate, 0);
  av_opt_set_sample_fmt(swr_ctx, "in_sample_fmt", ctx->sample_fmt, 0);
  av_opt_set_sample_fmt(swr_ctx, "out_sample_fmt", oc->dst_sample_fmt, 0);
  int result = swr_init(swr_ctx);
  if (result < 0) {
    av_log(NULL, AV_LOG_ERROR, "Can't init opus swr context\n");
    swr_free(&swr_ctx);
    return result;
  }
  oc->swr_ctx = swr_ctx;
  return 0;
}

static int opus_resample(OpusDecoder2 *oc, AVFrame *frame) {
  AVCodecContext *ctx = oc->codec_ctx;
  struct SwrContext *swr_ctx = oc->swr_ctx;
  int dst_channels = opus_dst_channels(oc);

  int dst_nb_samples = frame->nb_samples;
  if (swr_ctx) {
    dst_nb_samples = av_rescale_rnd(
        swr_get_delay(swr_ctx, ctx->sample_rate) + frame->nb_samples,
        oc->dst_sample_rate, ctx->sample_rate, AV_ROUND_UP);
  }

  int buf_size = av_samples_get_buffer_size(NULL, dst_channels,
      dst_nb_samples, oc->dst_sample_fmt, 1);
  if (oc->buffer == NULL || oc->buffer_size < buf_size) {
    if (oc->buffer) {
      av_free(oc->buffer);
    }
    oc->buffer_size = buf_size;
    oc->buffer = (uint8_t *)av_malloc(oc->buffer_size);
    if (!oc->buffer) {
      oc->buffer_size = 0;
      return AVERROR(ENOMEM);
    }
    av_samples_fill_arrays(oc->samples, NULL, oc->buffer, dst_channels,
        dst_nb_samples, oc->dst_sample_fmt, 1);
  }

  if (!swr_ctx) {
    // formats already match, just copy
    int result = av_samples_copy(oc->samples, frame->data, 0, 0, frame->nb_samples,
        opus_ctx_channels(ctx), ctx->sample_fmt);
    if (result < 0) {
      av_log(NULL, AV_LOG_ERROR, "Can't copy opus samples, %d\n", result);
      return result;
    }
    oc->dst_nb_samples = frame->nb_samples;
    oc->actual_buffer_size = av_samples_get_buffer_size(NULL, dst_channels,
        oc->dst_nb_samples, oc->dst_sample_fmt, 1);
    return 0;
  }

  int converted = swr_convert(swr_ctx, oc->samples, dst_nb_samples,
      (const uint8_t **)frame->data, frame->nb_samples);
  if (converted < 0) {
    av_log(NULL, AV_LOG_ERROR, "Can't resample opus\n");
    return converted;
  }
  // actual number of samples written
  oc->dst_nb_samples = converted;
  oc->actual_buffer_size = av_samples_get_buffer_size(NULL, dst_channels,
      converted, oc->dst_sample_fmt, 1);
  return 0;
}

void * open_opus_decoder(int in_channels, int out_sample_rate, int out_channels, int *error_code) {
  *error_code = 0;
  OpusDecoder2 *oc = (OpusDecoder2 *)malloc(sizeof(OpusDecoder2));
  if (!oc) {
    *error_code = -10001;
    return NULL;
  }
  memset(oc, 0, sizeof(OpusDecoder2));
  oc->out_sample_rate = out_sample_rate;
  oc->out_channels = out_channels;

  // should register all the codecs
  #if LIBAVCODEC_VERSION_MAJOR < 58
  // FFmpeg 3.x and earlier versions need to register all the codecs
      avcodec_register_all();
  #endif
  *error_code = -2001;

#if LIBAVCODEC_VERSION_INT >= AV_VERSION_INT(58, 0, 0)
  const AVCodec *codec = avcodec_find_decoder(AV_CODEC_ID_OPUS);
#else
  AVCodec *codec = avcodec_find_decoder(AV_CODEC_ID_OPUS);
#endif
  if (!codec) {
    av_log(NULL, AV_LOG_ERROR, "Can't find opus decoder, FFmpeg built without opus?\n");
    free(oc);
    return NULL;
  }
  *error_code = -2002;
  oc->codec_ctx = avcodec_alloc_context3(codec);
  if (!oc->codec_ctx) {
    av_log(NULL, AV_LOG_ERROR, "Can't allocate opus decoder context\n");
    free(oc);
    return NULL;
  }
  AVCodecContext *ctx = oc->codec_ctx;
  // opus always decodes at 48kHz internally
  ctx->sample_rate = 48000;
  opus_set_ctx_channels(ctx, in_channels > 0 ? in_channels : 1);
  ctx->thread_count = 1;

  *error_code = -2003;
  if (avcodec_open2(ctx, codec, NULL) < 0) {
    av_log(NULL, AV_LOG_ERROR, "Can't open opus decoder\n");
    avcodec_free_context(&oc->codec_ctx);
    free(oc);
    return NULL;
  }

  *error_code = -2004;
  oc->frame = av_frame_alloc();
  if (!oc->frame) {
    av_log(NULL, AV_LOG_ERROR, "Can't allocate opus frame\n");
    avcodec_free_context(&oc->codec_ctx);
    free(oc);
    return NULL;
  }
  *error_code = 0;

  av_log(NULL, AV_LOG_INFO, "opus decoder opened, in_channels %d, out_rate %d, out_channels %d\n",
    opus_ctx_channels(ctx), out_sample_rate, out_channels);
  return oc;
}

int decode_opus(void *handle, const uint8_t *data, int data_size, MediaFrame *frame) {
  if (!handle || !data || data_size <= 0 || !frame) {
    return AVERROR(EINVAL);
  }
  OpusDecoder2 *oc = (OpusDecoder2 *)handle;
  AVCodecContext *ctx = oc->codec_ctx;
  AVFrame *fr = oc->frame;
  int result;

#if LIBAVCODEC_VERSION_INT >= AV_VERSION_INT(58, 0, 0)
  AVPacket *pkt = av_packet_alloc();
  if (!pkt) {
    return AVERROR(ENOMEM);
  }
  // wrap the raw opus payload without copying; ffmpeg won't free user buffer
  pkt->data = (uint8_t *)data;
  pkt->size = data_size;

  result = avcodec_send_packet(ctx, pkt);
  pkt->data = NULL;
  pkt->size = 0;
  av_packet_free(&pkt);
#else
  AVPacket pkt;
  av_init_packet(&pkt);
  // wrap the raw opus payload without copying; ffmpeg won't free user buffer
  pkt.data = (uint8_t *)data;
  pkt.size = data_size;

  result = avcodec_send_packet(ctx, &pkt);
  av_packet_unref(&pkt);
#endif
  if (result < 0) {
    av_log(NULL, AV_LOG_ERROR, "opus send packet error %d\n", result);
    return result;
  }

  // one opus packet decodes to one frame
  result = avcodec_receive_frame(ctx, fr);
  if (result < 0) {
    // AVERROR(EAGAIN): need more packets; AVERROR_EOF: flushed
    return result;
  }

  if (!oc->swr_inited) {
    int ret = opus_init_swr(oc);
    if (ret < 0) {
      return ret;
    }
  }

  int ret = opus_resample(oc, fr);
  if (ret < 0) {
    av_log(NULL, AV_LOG_ERROR, "opus resample error %d\n", ret);
    return ret;
  }

  frame->frame_type = AVMEDIA_TYPE_AUDIO;
  frame->stream_index = 0;
  frame->pts = fr->pts;
  frame->buffer = oc->buffer;
  frame->buffer_size = oc->actual_buffer_size;
  frame->format = oc->dst_sample_fmt;
  frame->samples = oc->dst_nb_samples;
  frame->channels = opus_dst_channels(oc);
  frame->sample_rate = oc->dst_sample_rate;
  frame->bytes_per_sample = av_get_bytes_per_sample(oc->dst_sample_fmt);
  return 0;
}

void close_opus_decoder(void *handle) {
  if (!handle) {
    return;
  }
  OpusDecoder2 *oc = (OpusDecoder2 *)handle;
  if (oc->swr_ctx) {
    swr_free(&oc->swr_ctx);
  }
  if (oc->frame) {
    av_frame_free(&oc->frame);
  }
  if (oc->codec_ctx) {
    avcodec_free_context(&oc->codec_ctx);
  }
  if (oc->buffer) {
    av_free(oc->buffer);
    oc->buffer = NULL;
    oc->buffer_size = 0;
  }
  free(oc);
}
#else
void * open_opus_decoder(int in_channels, int out_sample_rate, int out_channels, int *error_code) {
  *error_code = 0;
  return NULL;
}

int decode_opus(void *handle, const uint8_t *data, int data_size, MediaFrame *frame) {
  return -1;
}

void close_opus_decoder(void *handle) {
  return;
}
#endif
