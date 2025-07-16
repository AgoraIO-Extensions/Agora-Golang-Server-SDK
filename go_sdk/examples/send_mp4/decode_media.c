#include <libavformat/avformat.h>
#include <libavcodec/avcodec.h>
#include <libavutil/log.h>
#include <libavutil/imgutils.h>
#include <libavutil/samplefmt.h>
#include <libavutil/avutil.h>
#include <libavutil/opt.h>
#include <libavutil/channel_layout.h>
#include <libswresample/swresample.h>
#include <libswscale/swscale.h>
#include <libavcodec/bsf.h>
#include <stdint.h>
#include <stdlib.h>
#include "decode_media.h"


// history
// 2025-06-13: fix memory leak when pop packet, by zhourui@agora
// 2025-07-01: fix crash when close media file, by zhourui@agora

#define MAX_AUDIO_CHANNELS 10
#define AVSYNC_MAX_AUDIO_SIZE 5000
#define AVSYNC_MAX_VIDEO_SIZE 5000

typedef struct _DecodeContext {
  int stream_index;
  AVCodecContext *codec_ctx;
  const AVCodecParameters *codec_par;
  const AVCodec *codec;
  AVFrame *frame;
  int is_eof;

  uint8_t *buffer;
  int buffer_size;
  // for audio
  uint8_t *samples[MAX_AUDIO_CHANNELS];

  // for audio resample
  struct SwrContext *swr_ctx;
  AVChannelLayout dst_ch_layout;
  int dst_sample_rate;
  enum AVSampleFormat dst_sample_fmt;
  int dst_nb_samples;
  int actual_buffer_size;

  // for video resize
  struct SwsContext *sws_ctx;
  enum AVPixelFormat dst_pix_fmt;
  int dst_width;
  int dst_height;
  uint8_t *dst_data[4];
  int dst_linesize[4];
} DecodeContext;

typedef struct _MediaDecoder {
  AVFormatContext *fmt_ctx;
  DecodeContext video_ctx;
  DecodeContext audio_ctx;
  AVPacket *pkt;

  // for avsync
  MediaPacket head_pkt;
  MediaPacket *video_tail_pkt;
  MediaPacket *audio_tail_pkt;
  int read_error;

  // for bitstream filter
  AVBSFContext *bsf;
} MediaDecoder;

int init_swr(DecodeContext *decode_ctx) {
  AVChannelLayout stereo_layout = AV_CHANNEL_LAYOUT_STEREO;
  AVCodecContext *codec_ctx = decode_ctx->codec_ctx;
  decode_ctx->dst_ch_layout = codec_ctx->ch_layout;
  if (decode_ctx->dst_ch_layout.nb_channels > 2) {
    decode_ctx->dst_ch_layout = stereo_layout;
  }
  decode_ctx->dst_sample_rate = codec_ctx->sample_rate;
  if (decode_ctx->dst_sample_rate > 48000) {
    decode_ctx->dst_sample_rate = 48000;
  }
  decode_ctx->dst_sample_fmt = AV_SAMPLE_FMT_S16;
  if (codec_ctx->sample_fmt == decode_ctx->dst_sample_fmt &&
    codec_ctx->sample_rate == decode_ctx->dst_sample_rate &&
    codec_ctx->ch_layout.nb_channels == decode_ctx->dst_ch_layout.nb_channels) {
    // resample is not needed
    return 0;
  }

  struct SwrContext *swr_ctx = swr_alloc();
  if (!swr_ctx) {
    av_log(NULL, AV_LOG_ERROR, "Can't allocate swr context\n");
    return AVERROR(ENOMEM);
  }

  av_opt_set_chlayout(swr_ctx, "in_chlayout", &codec_ctx->ch_layout, 0);
  av_opt_set_chlayout(swr_ctx, "out_chlayout", &decode_ctx->dst_ch_layout, 0);
  av_opt_set_int(swr_ctx, "in_sample_rate", codec_ctx->sample_rate, 0);
  av_opt_set_int(swr_ctx, "out_sample_rate", decode_ctx->dst_sample_rate, 0);
  av_opt_set_sample_fmt(swr_ctx, "in_sample_fmt", codec_ctx->sample_fmt, 0);
  av_opt_set_sample_fmt(swr_ctx, "out_sample_fmt", decode_ctx->dst_sample_fmt, 0);
  int result = swr_init(swr_ctx);
  if (result < 0) {
    av_log(NULL, AV_LOG_ERROR, "Can't init swr context\n");
    swr_free(&swr_ctx);
    return result;
  }
  decode_ctx->swr_ctx = swr_ctx;
  return 0;
}

int resample_audio(DecodeContext *decode_ctx, AVFrame *frame) {
  AVCodecContext *codec_ctx = decode_ctx->codec_ctx;
  struct SwrContext *swr_ctx = decode_ctx->swr_ctx;
  int result = 0;
  int dst_nb_samples = frame->nb_samples;
  if (swr_ctx) {
    // compute the number of samples after resample
    dst_nb_samples = av_rescale_rnd(
        swr_get_delay(swr_ctx, codec_ctx->sample_rate) + frame->nb_samples,
        decode_ctx->dst_sample_rate, codec_ctx->sample_rate, AV_ROUND_UP);
    // if (dst_nb_samples > frame->nb_samples) {
    //   av_log(NULL, AV_LOG_ERROR, "dst_nb_samples %d > frame->nb_samples %d\n", dst_nb_samples, frame->nb_samples);
    //   return -1;
    // }
  }
  decode_ctx->dst_nb_samples = dst_nb_samples;
  int buf_size = av_samples_get_buffer_size(NULL, decode_ctx->dst_ch_layout.nb_channels, dst_nb_samples, decode_ctx->dst_sample_fmt, 1);
  if (decode_ctx->buffer == NULL || decode_ctx->buffer_size < buf_size) {
    if (decode_ctx->buffer) {
      av_free(decode_ctx->buffer);
    }
    decode_ctx->buffer_size = buf_size;
    decode_ctx->buffer = (uint8_t *)av_malloc(decode_ctx->buffer_size);
    av_samples_fill_arrays(decode_ctx->samples, NULL, decode_ctx->buffer, decode_ctx->dst_ch_layout.nb_channels, dst_nb_samples, decode_ctx->dst_sample_fmt, 1);
  }
  decode_ctx->actual_buffer_size = buf_size;
  if (!swr_ctx) {
    // just copy audio data
    result = av_samples_copy(decode_ctx->samples, frame->data, 0, 0, frame->nb_samples, codec_ctx->ch_layout.nb_channels, codec_ctx->sample_fmt);
    if (result < 0) {
      av_log(NULL, AV_LOG_ERROR, "Can't copy audio samples, %d\n", result);
      return result;
    }
    return 0;
  }

  // resample audio data
  result = swr_convert(swr_ctx, decode_ctx->samples, dst_nb_samples, (const uint8_t **)frame->data, frame->nb_samples);
  if (result < 0) {
    av_log(NULL, AV_LOG_ERROR, "Can't resample audio\n");
    return result;
  }
  return 0;
}

int deinit_swr(DecodeContext *decode_ctx) {
  if (decode_ctx->swr_ctx) {
    swr_free(&decode_ctx->swr_ctx);
  }
  return 0;
}

int init_sws(DecodeContext *decode_ctx) {
  AVCodecContext *codec_ctx = decode_ctx->codec_ctx;
  decode_ctx->dst_pix_fmt = AV_PIX_FMT_YUV420P;
  decode_ctx->dst_width = codec_ctx->width;
  decode_ctx->dst_height = codec_ctx->height;
  int ret = av_image_alloc(decode_ctx->dst_data, decode_ctx->dst_linesize, 
      decode_ctx->dst_width, decode_ctx->dst_height, decode_ctx->dst_pix_fmt, 1);
  if (ret < 0) {
    av_log(NULL, AV_LOG_ERROR, "Can't allocate image\n");
    return ret;
  }
  decode_ctx->buffer = decode_ctx->dst_data[0];
  decode_ctx->buffer_size = av_image_get_buffer_size(decode_ctx->dst_pix_fmt, decode_ctx->dst_width, decode_ctx->dst_height, 1);

  if (decode_ctx->dst_pix_fmt == codec_ctx->pix_fmt &&
    decode_ctx->dst_width == codec_ctx->width &&
    decode_ctx->dst_height == codec_ctx->height) {
    return 0;
  }
  struct SwsContext *sws_ctx = sws_getContext(
    codec_ctx->width, codec_ctx->height, codec_ctx->pix_fmt,
    decode_ctx->dst_width, decode_ctx->dst_height, decode_ctx->dst_pix_fmt, SWS_BILINEAR, NULL, NULL, NULL);
  if (!sws_ctx) {
    av_log(NULL, AV_LOG_ERROR, "Can't allocate sws context\n");
    return AVERROR(ENOMEM);
  }
  decode_ctx->sws_ctx = sws_ctx;
  return 0;
}

int resize_video(DecodeContext *decode_ctx, AVFrame *fr) {
  AVCodecContext *ctx = decode_ctx->codec_ctx;
  struct SwsContext *sws_ctx = decode_ctx->sws_ctx;
  if (!sws_ctx) {
    int ret = av_image_copy_to_buffer(decode_ctx->buffer, decode_ctx->buffer_size,
                            (const uint8_t* const *)fr->data, (const int*) fr->linesize,
                            ctx->pix_fmt, ctx->width, ctx->height, 1);
    if (ret < 0) {
      av_log(NULL, AV_LOG_ERROR, "Error copying image to buffer\n");
    }
    return ret;
  }
  int ret = sws_scale(sws_ctx, (const uint8_t* const *)fr->data, fr->linesize, 
      0, ctx->height, (uint8_t * const *)decode_ctx->dst_data, decode_ctx->dst_linesize);
  if (ret < 0) {
    av_log(NULL, AV_LOG_ERROR, "Error scaling image\n");
  }
  return ret;
}

int deinit_sws(DecodeContext *decode_ctx) {
  if (decode_ctx->sws_ctx) {
    sws_freeContext(decode_ctx->sws_ctx);
    decode_ctx->sws_ctx = NULL;
  }
  return 0;
}

int deinit_decoder(DecodeContext *decode_ctx) {
  deinit_swr(decode_ctx);
  deinit_sws(decode_ctx);
  avcodec_free_context(&decode_ctx->codec_ctx);
  av_frame_free(&decode_ctx->frame);
  if (decode_ctx->buffer) {
    av_free(decode_ctx->buffer);
    decode_ctx->buffer = NULL;
    decode_ctx->buffer_size = 0;
  }
  decode_ctx->stream_index = -1;
  decode_ctx->is_eof = 1;
  return 0;
}

int init_decoder(MediaDecoder *decoder, int media_type) {
  AVFormatContext *fmt_ctx = decoder->fmt_ctx;
  DecodeContext *decode_ctx = NULL;
  if (media_type == AVMEDIA_TYPE_VIDEO) {
    decode_ctx = &decoder->video_ctx;
  } else if (media_type == AVMEDIA_TYPE_AUDIO) {
    decode_ctx = &decoder->audio_ctx;
  } else {
    return -1;
  }
  int result = 0;

  decode_ctx->stream_index = av_find_best_stream(fmt_ctx, media_type, -1, -1, NULL, 0);
  if (decode_ctx->stream_index < 0) {
    av_log(NULL, AV_LOG_ERROR, "Can't find video stream in input file\n");
    deinit_decoder(decode_ctx);
    return -1;
  }

  AVCodecParameters *origin_par = fmt_ctx->streams[decode_ctx->stream_index]->codecpar;
  decode_ctx->codec_par = origin_par;

  const AVCodec *codec = avcodec_find_decoder(origin_par->codec_id);
  if (!codec) {
      av_log(NULL, AV_LOG_ERROR, "Can't find decoder\n");
      deinit_decoder(decode_ctx);
      return -1;
  }
  decode_ctx->codec = codec;

  decode_ctx->codec_ctx = avcodec_alloc_context3(codec);
  if (!decode_ctx->codec_ctx) {
      av_log(NULL, AV_LOG_ERROR, "Can't allocate decoder context\n");
      deinit_decoder(decode_ctx);
      return AVERROR(ENOMEM);
  }
  AVCodecContext *codec_ctx = decode_ctx->codec_ctx;

  result = avcodec_parameters_to_context(codec_ctx, origin_par);
  if (result) {
      av_log(NULL, AV_LOG_ERROR, "Can't copy decoder context\n");
      deinit_decoder(decode_ctx);
      return result;
  }

  // ctx->draw_horiz_band = draw_horiz_band;
  codec_ctx->thread_count = 1;

  result = avcodec_open2(codec_ctx, codec, NULL);
  if (result < 0) {
      av_log(codec_ctx, AV_LOG_ERROR, "Can't open decoder\n");
      deinit_decoder(decode_ctx);
      return result;
  }

  decode_ctx->frame = av_frame_alloc();
  if (!decode_ctx->frame) {
      av_log(NULL, AV_LOG_ERROR, "Can't allocate frame\n");
      deinit_decoder(decode_ctx);
      return AVERROR(ENOMEM);
  }

  if (media_type == AVMEDIA_TYPE_VIDEO) {
    av_log(NULL, AV_LOG_INFO, "stream index %d, video codec: %s, pix_fmt %s, width %d, height %d\n", 
      decode_ctx->stream_index, codec->name, av_get_pix_fmt_name(codec_ctx->pix_fmt),
      codec_ctx->width, codec_ctx->height);
    init_sws(decode_ctx);
  } else if (media_type == AVMEDIA_TYPE_AUDIO) {
    av_log(NULL, AV_LOG_INFO, "stream index %d, audio codec: %s, sample_fmt %s, sample_rate %d, channels %d, frame_size %d\n", 
      decode_ctx->stream_index, codec->name, av_get_sample_fmt_name(codec_ctx->sample_fmt),
      codec_ctx->sample_rate, codec_ctx->ch_layout.nb_channels, codec_ctx->frame_size);
    init_swr(decode_ctx);
  }

  return 0;
}

void * open_media_file(const char *file_name) {
    MediaDecoder *decoder = (MediaDecoder *)malloc(sizeof(MediaDecoder));
    memset(decoder, 0, sizeof(MediaDecoder));
    decoder->video_tail_pkt = &decoder->head_pkt;
    decoder->audio_tail_pkt = &decoder->head_pkt;

    int result = 0;

    // av_log_set_level(AV_LOG_DEBUG);
    result = avformat_open_input(&decoder->fmt_ctx, file_name, NULL, NULL);
    if (result < 0) {
        av_log(NULL, AV_LOG_ERROR, "Can't open file\n");
        close_media_file(decoder);
        return NULL;
    }

    result = avformat_find_stream_info(decoder->fmt_ctx, NULL);
    if (result < 0) {
        av_log(NULL, AV_LOG_ERROR, "Can't get stream info\n");
        close_media_file(decoder);
        return NULL;
    }

    AVPacket *pkt = av_packet_alloc();
    if (!pkt) {
        av_log(NULL, AV_LOG_ERROR, "Cannot allocate packet\n");
        close_media_file(decoder);
        return NULL;
    }
    decoder->pkt = pkt;

    init_decoder(decoder, AVMEDIA_TYPE_VIDEO);
    init_decoder(decoder, AVMEDIA_TYPE_AUDIO);
    return decoder;
}

void push_packet(MediaDecoder *d, MediaPacket *pkt) {
  MediaPacket **tail = NULL;
  if (pkt->media_type == AVMEDIA_TYPE_VIDEO) {
    tail = &(d->video_tail_pkt);
  } else if (pkt->media_type == AVMEDIA_TYPE_AUDIO) {
    tail = &(d->audio_tail_pkt);
  } else {
    return;
  }
  MediaPacket *last_pkt = *tail;
  MediaPacket *cur_pkt = last_pkt->next;
  while (cur_pkt) {
    if (cur_pkt->pts > pkt->pts) {
      break;
    }
    last_pkt = cur_pkt;
    cur_pkt = last_pkt->next;
  }
  last_pkt->next = pkt;
  *tail = pkt;
}

MediaPacket *pop_packet(MediaDecoder *d) {
    // If the queue is empty, return NULL
    if (!d->head_pkt.next) {
        return NULL;
    }

    // If read error occurs, pop the packet directly without synchronization check
    if (d->read_error != 0) {
        MediaPacket *pkt_node = d->head_pkt.next;
        d->head_pkt.next = pkt_node->next;
        if (d->video_tail_pkt == pkt_node) {
            d->video_tail_pkt = &d->head_pkt;
        }
        if (d->audio_tail_pkt == pkt_node) {
            d->audio_tail_pkt = &d->head_pkt;
        }
        pkt_node->next = NULL; // Clear the next pointer to prevent accidental access
        return pkt_node;
    }

    // Check if there are video and audio streams
    int has_video_stream = (d->video_ctx.stream_index >= 0);
    int has_audio_stream = (d->audio_ctx.stream_index >= 0);
    
    // If there is only one stream type, pop the packet directly without synchronization check
    if (!has_video_stream || !has_audio_stream) {
        MediaPacket *pkt_node = d->head_pkt.next;
        d->head_pkt.next = pkt_node->next;
        if (d->video_tail_pkt == pkt_node) {
            d->video_tail_pkt = &d->head_pkt;
        }
        if (d->audio_tail_pkt == pkt_node) {
            d->audio_tail_pkt = &d->head_pkt;
        }
        pkt_node->next = NULL;
        av_log(NULL, AV_LOG_DEBUG, "Single stream mode, popping packet directly\n");
        return pkt_node;
    }

    // Double stream synchronization logic
    int64_t video_size = 0, audio_size = 0;
    int video_exist = 0, audio_exist = 0;
    
    if (d->video_tail_pkt != &d->head_pkt) {
        video_exist = 1;
        video_size = d->video_tail_pkt->pts - d->head_pkt.next->pts;
    }
    if (d->audio_tail_pkt != &d->head_pkt) {
        audio_exist = 1;
        audio_size = d->audio_tail_pkt->pts - d->head_pkt.next->pts;
    }
    
    // Only skip when both streams exist and the buffer is insufficient
    if (video_exist && audio_exist) {
        if (video_size < AVSYNC_MAX_VIDEO_SIZE && audio_size < AVSYNC_MAX_AUDIO_SIZE) {
            av_log(NULL, AV_LOG_DEBUG, "Both streams buffer insufficient, video_size %lld, audio_size %lld\n", video_size, audio_size);
            return NULL;
        }
    }

    // Pop the packet
    MediaPacket *pkt_node = d->head_pkt.next;
    d->head_pkt.next = pkt_node->next;
    if (d->video_tail_pkt == pkt_node) {
        d->video_tail_pkt = &d->head_pkt;
    }
    if (d->audio_tail_pkt == pkt_node) {
        d->audio_tail_pkt = &d->head_pkt;
    }
    pkt_node->next = NULL;
    return pkt_node;
}

int get_packet(void *decoder, MediaPacket **packet) {
  MediaDecoder *d = (MediaDecoder *)decoder;
  AVFormatContext *fmt_ctx = d->fmt_ctx;
  AVPacket *pkt = NULL;

  *packet = NULL;
  if (d->read_error == 0) {
    // packet->pkt = pkt;
    // packet->media_type = AVMEDIA_TYPE_UNKNOWN;
    pkt = av_packet_alloc();
    if (!pkt) {
        av_log(NULL, AV_LOG_ERROR, "Cannot allocate packet\n");
        return AVERROR(ENOMEM);
    }

    int result = av_read_frame(fmt_ctx, pkt);
    av_log(NULL, AV_LOG_DEBUG, "read frame result %d, pkg stream index %d\n",
      result, pkt->stream_index);
    if (result >= 0) {
      if (pkt->stream_index == d->video_ctx.stream_index) {
          AVStream *s = fmt_ctx->streams[pkt->stream_index];
          MediaPacket *mpkt = (MediaPacket *)malloc(sizeof(MediaPacket));
          memset(mpkt, 0, sizeof(MediaPacket));
          pkt->time_base = s->time_base;
          mpkt->pkt = pkt;
          mpkt->media_type = AVMEDIA_TYPE_VIDEO;
          mpkt->pts = pkt->pts * 1000 * av_q2d(s->time_base);
          mpkt->width = s->codecpar->width;
          mpkt->height = s->codecpar->height;
          mpkt->framerate_num = s->r_frame_rate.num;
          mpkt->framerate_den = s->r_frame_rate.den;
          push_packet(d, mpkt);
      } else if (pkt->stream_index == d->audio_ctx.stream_index) {
          AVStream *s = fmt_ctx->streams[pkt->stream_index];
          MediaPacket *mpkt = (MediaPacket *)malloc(sizeof(MediaPacket));
          memset(mpkt, 0, sizeof(MediaPacket));
          pkt->time_base = s->time_base;
          mpkt->pkt = pkt;
          mpkt->media_type = AVMEDIA_TYPE_AUDIO;
          mpkt->pts = pkt->pts * 1000 * av_q2d(s->time_base);
          push_packet(d, mpkt);
      } else {
        av_packet_free(&pkt);
        pkt = NULL;
      }
    } else {
      av_packet_free(&pkt);
      pkt = NULL;
      if (result != AVERROR(EAGAIN)) {
        d->read_error = result;
      }
    }
  }
  
  if (d->head_pkt.next == NULL) {
    return d->read_error;
  }

  *packet = pop_packet(d);
  return 0;
}

int h264_to_annexb(void *decoder, MediaPacket **packet) {
  MediaDecoder *d = (MediaDecoder *)decoder;
  AVFormatContext *fmt_ctx = d->fmt_ctx;
  AVPacket *pkt = NULL;
  
  if (*packet) {
    pkt = (*packet)->pkt;
  }

  if (!d->bsf && !pkt) {
    av_log(NULL, AV_LOG_ERROR, "bsf not initilized when flushing\n");
    return -1;
  }

  if (!d->bsf/* && pkt*/) {
    const AVBitStreamFilter *bsf = av_bsf_get_by_name("h264_mp4toannexb");
    if (!bsf) {
      av_log(NULL, AV_LOG_ERROR, "Can't find bitstream filter\n");
      return -1;
    }
    int result = av_bsf_alloc(bsf, &d->bsf);
    if (result < 0) {
      av_log(NULL, AV_LOG_ERROR, "Can't allocate bitstream filter\n");
      return result;
    }
    result = avcodec_parameters_copy(d->bsf->par_in, fmt_ctx->streams[pkt->stream_index]->codecpar);
    if (result < 0) {
      av_log(NULL, AV_LOG_ERROR, "Can't copy codec parameters\n");
      return result;
    }
    d->bsf->time_base_in = fmt_ctx->streams[pkt->stream_index]->time_base;
    result = av_bsf_init(d->bsf);
    if (result < 0) {
      av_log(NULL, AV_LOG_ERROR, "Can't init bitstream filter\n");
      return result;
    }
  }

  if (pkt)
      av_packet_rescale_ts(pkt, pkt->time_base, d->bsf->time_base_in);

  int ret = av_bsf_send_packet(d->bsf, pkt);
  if (ret < 0) {
      free_packet(packet);
      av_log(NULL, AV_LOG_ERROR, "Error submitting a packet for filtering: %s\n",
              av_err2str(ret));
      return ret;
  }

  AVPacket *pkt_bsf = av_packet_alloc();
  ret = av_bsf_receive_packet(d->bsf, pkt_bsf);
  if (ret == AVERROR(EAGAIN)) {
    av_packet_free(&pkt_bsf);
    pkt_bsf = NULL;
    free_packet(packet);
    return 0;
  } else if (ret < 0) {
      if (ret != AVERROR_EOF)
          av_log(NULL, AV_LOG_ERROR,
                  "Error applying bitstream filters to a packet: %s\n",
                  av_err2str(ret));
      av_packet_free(&pkt_bsf);
      pkt_bsf = NULL;
      free_packet(packet);
      return ret;
  }

  pkt_bsf->time_base = d->bsf->time_base_out;

  // set output packet
  if (!(*packet)) {
    *packet = (MediaPacket *)malloc(sizeof(MediaPacket));
    memset(*packet, 0, sizeof(MediaPacket));
  }
  if (pkt) {
    av_packet_free(&pkt);
    pkt = NULL;
  }
  (*packet)->pkt = pkt_bsf;
  (*packet)->pts = pkt_bsf->pts * 1000 * av_q2d(pkt_bsf->time_base);
  return 0;
}

int free_packet(MediaPacket **packet) {
  if (!(*packet)) {
    return 0;
  }
  if ((*packet)->pkt) {
    av_packet_free(&(*packet)->pkt);
    (*packet)->pkt = NULL;
  }
  free(*packet);
  *packet = NULL;
  return 0;
}

int decode_packet(void *decoder, MediaPacket *packet, MediaFrame *frame) {
  MediaDecoder *d = (MediaDecoder *)decoder;
  AVFormatContext *fmt_ctx = d->fmt_ctx;
  AVPacket *pkt = packet->pkt;
  int media_type = packet->media_type;

  DecodeContext *decode_ctx = NULL;
  if (media_type == AVMEDIA_TYPE_VIDEO) {
    decode_ctx = &d->video_ctx;
  } else if (media_type == AVMEDIA_TYPE_AUDIO) {
    decode_ctx = &d->audio_ctx;
  } else {
    // this branch should not be reached
    av_packet_unref(pkt);
    return -1;
  }

  AVCodecContext *ctx = decode_ctx->codec_ctx;
  AVFrame *fr = decode_ctx->frame;

  // pkt will be empty on read error/EOF
  int result = avcodec_send_packet(ctx, pkt);

  av_packet_unref(pkt);

  if (result < 0) {
      av_log(NULL, AV_LOG_ERROR, "Error submitting a packet for decoding\n");
      return result;
  }

  while (result >= 0) {
      result = avcodec_receive_frame(ctx, fr);
      if (result == AVERROR_EOF) {
          av_log(NULL, AV_LOG_INFO, "decode media %d EOF\n", media_type);
          decode_ctx->is_eof = 1;
          break;
      } else if (result == AVERROR(EAGAIN)) {
          break;
      } else if (result < 0) {
          av_log(NULL, AV_LOG_ERROR, "Error decoding frame\n");
          return result;
      }

      if (media_type == AVMEDIA_TYPE_VIDEO) {
        int ret = resize_video(decode_ctx, fr);
        if (ret < 0) {
          av_log(NULL, AV_LOG_ERROR, "Error resize video, code %d\n", ret);
          return ret;
        }
        frame->frame_type = AVMEDIA_TYPE_VIDEO;
        frame->stream_index = decode_ctx->stream_index;
        frame->pts = fr->pts * 1000 * av_q2d(fmt_ctx->streams[decode_ctx->stream_index]->time_base);
        frame->buffer = decode_ctx->buffer;
        frame->buffer_size = decode_ctx->buffer_size;
        frame->format = decode_ctx->dst_pix_fmt;
        frame->width = decode_ctx->dst_width;
        frame->height = decode_ctx->dst_height;
        frame->stride = decode_ctx->dst_width;
      } else if (media_type == AVMEDIA_TYPE_AUDIO) {
        int ret = resample_audio(decode_ctx, fr);
        if (ret < 0) {
          av_log(NULL, AV_LOG_ERROR, "Error resample audio, code %d\n", ret);
          return ret;
        }
        frame->frame_type = AVMEDIA_TYPE_AUDIO;
        frame->stream_index = decode_ctx->stream_index;
        frame->pts = fr->pts * 1000 * av_q2d(fmt_ctx->streams[decode_ctx->stream_index]->time_base);
        frame->buffer = decode_ctx->buffer;
        frame->buffer_size = decode_ctx->actual_buffer_size;
        frame->format = decode_ctx->dst_sample_fmt;
        frame->samples = decode_ctx->dst_nb_samples;
        frame->channels = decode_ctx->dst_ch_layout.nb_channels;
        frame->sample_rate = decode_ctx->dst_sample_rate;
        frame->bytes_per_sample = av_get_bytes_per_sample(decode_ctx->dst_sample_fmt);
      }
      return 0;
  }
  return result;
}


int get_frame(void *decoder, MediaFrame *frame) {
  MediaDecoder *d = (MediaDecoder *)decoder;
  AVPacket *pkt = d->pkt;
  AVFormatContext *fmt_ctx = d->fmt_ctx;

  int result = 0;
  while (result >= 0) {
      result = av_read_frame(fmt_ctx, pkt);

      DecodeContext *decode_ctx = NULL;
      int media_type = AVMEDIA_TYPE_UNKNOWN;
      if (result >= 0) {
        if (pkt->stream_index == d->video_ctx.stream_index) {
            media_type = AVMEDIA_TYPE_VIDEO;
        } else if (pkt->stream_index == d->audio_ctx.stream_index) {
            media_type = AVMEDIA_TYPE_AUDIO;
        } else {
            // skip other streams
            av_packet_unref(pkt);
            continue;
        }
      } else {
        // EOF
        if (d->video_ctx.is_eof && d->audio_ctx.is_eof) {
          return AVERROR_EOF;
        }
        // flush decoder if decoder did not reach EOF
        if (!d->video_ctx.is_eof && d->audio_ctx.is_eof) {
          media_type = AVMEDIA_TYPE_VIDEO;
        } else if (d->video_ctx.is_eof && !d->audio_ctx.is_eof) {
          media_type = AVMEDIA_TYPE_AUDIO;
        } else {
          if (d->video_ctx.frame->pts < d->audio_ctx.frame->pts) {
            media_type = AVMEDIA_TYPE_VIDEO;
          } else {
            media_type = AVMEDIA_TYPE_AUDIO;
          }
        }
      }
      if (media_type == AVMEDIA_TYPE_VIDEO) {
        decode_ctx = &d->video_ctx;
      } else if (media_type == AVMEDIA_TYPE_AUDIO) {
        decode_ctx = &d->audio_ctx;
      } else {
        // this branch should not be reached
        av_packet_unref(pkt);
        continue;
      }

      av_log(NULL, AV_LOG_DEBUG, "read frame result %d, pkg stream index %d, media type %d\n",
       result, pkt->stream_index, media_type);
      AVCodecContext *ctx = decode_ctx->codec_ctx;
      AVFrame *fr = decode_ctx->frame;

      // pkt will be empty on read error/EOF
      result = avcodec_send_packet(ctx, pkt);

      av_packet_unref(pkt);

      if (result < 0) {
          av_log(NULL, AV_LOG_ERROR, "Error submitting a packet for decoding\n");
          return result;
      }

      while (result >= 0) {
          result = avcodec_receive_frame(ctx, fr);
          if (result == AVERROR_EOF) {
              av_log(NULL, AV_LOG_INFO, "decode media %d EOF\n", media_type);
              decode_ctx->is_eof = 1;
              result = 0;
              break;
          } else if (result == AVERROR(EAGAIN)) {
              result = 0;
              break;
          } else if (result < 0) {
              av_log(NULL, AV_LOG_ERROR, "Error decoding frame\n");
              return result;
          }

          if (media_type == AVMEDIA_TYPE_VIDEO) {
            int ret = resize_video(decode_ctx, fr);
            if (ret < 0) {
              av_log(NULL, AV_LOG_ERROR, "Error resize video, code %d\n", ret);
              return ret;
            }
            frame->frame_type = AVMEDIA_TYPE_VIDEO;
            frame->stream_index = decode_ctx->stream_index;
            frame->pts = fr->pts * 1000 * av_q2d(fmt_ctx->streams[decode_ctx->stream_index]->time_base);
            frame->buffer = decode_ctx->buffer;
            frame->buffer_size = decode_ctx->buffer_size;
            frame->format = decode_ctx->dst_pix_fmt;
            frame->width = decode_ctx->dst_width;
            frame->height = decode_ctx->dst_height;
            frame->stride = decode_ctx->dst_width;
          } else if (media_type == AVMEDIA_TYPE_AUDIO) {
            int ret = resample_audio(decode_ctx, fr);
            if (ret < 0) {
              av_log(NULL, AV_LOG_ERROR, "Error resample audio, code %d\n", ret);
              return ret;
            }
            frame->frame_type = AVMEDIA_TYPE_AUDIO;
            frame->stream_index = decode_ctx->stream_index;
            frame->pts = fr->pts * 1000 * av_q2d(fmt_ctx->streams[decode_ctx->stream_index]->time_base);
            frame->buffer = decode_ctx->buffer;
            frame->buffer_size = decode_ctx->actual_buffer_size;
            frame->format = decode_ctx->dst_sample_fmt;
            frame->samples = decode_ctx->dst_nb_samples;
            frame->channels = decode_ctx->dst_ch_layout.nb_channels;
            frame->sample_rate = decode_ctx->dst_sample_rate;
            frame->bytes_per_sample = av_get_bytes_per_sample(decode_ctx->dst_sample_fmt);
          }
          return 0;
      }
  }
  return 0;
}

void close_media_file(void *decoder) {
    MediaDecoder *d = (MediaDecoder *)decoder;
    deinit_decoder(&d->video_ctx);
    deinit_decoder(&d->audio_ctx);
    // free avsync
    int count = 0;
    while (d->head_pkt.next) {
      MediaPacket *pkt = d->head_pkt.next;
      d->head_pkt.next = pkt->next;
      av_packet_free(&pkt->pkt);
      pkt->pkt = NULL;
      free(pkt);
      pkt = NULL;
      count++;
    }
    av_log(NULL, AV_LOG_WARNING, "Free %d packets\n", count);
    // free bitstream filter
    if (d->bsf) {
      av_bsf_free(&d->bsf);
    }
    avformat_close_input(&d->fmt_ctx);
    av_packet_free(&d->pkt);
    d->pkt = NULL;
    free(d);
    d = NULL;
}