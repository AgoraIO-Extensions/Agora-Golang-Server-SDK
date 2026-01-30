/**
 * H.264 realtime decoder implementation
 * Function: decode H.264 NAL unit to YUV format (used for realtime stream decoding)
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>

#include <libavcodec/avcodec.h>
#include <libavutil/error.h>

#include "h26x_decode_frame.h"

// Decoder structure definition
struct H264Decoder {
    AVCodecContext *codec_ctx;  // decoder context
    const AVCodec *codec;        // decoder
    AVFrame *frame;              // decoded frame
    AVPacket *pkt;               // packet
    int frame_count;             // number of decoded frames
    int flush_sent;              // whether the flush packet has been sent
    uint8_t *buffer;             // buffer
    int buffer_size;             // buffer size
    int buffer_pos;              // current buffer position
};

/**
 * Initialize H.264 realtime decoder
 */
H264Decoder* h264_decoder_init(void) {
    H264Decoder *decoder = (H264Decoder*)calloc(1, sizeof(H264Decoder));
    if (!decoder) {
        fprintf(stderr, "Failed to allocate decoder memory\n");
        return NULL;
    }
    // should register all the codecs
    #if LIBAVCODEC_VERSION_MAJOR < 58
    // FFmpeg 3.x and earlier versions need to register all the codecs
        avcodec_register_all();
    #endif

    // Find H.264 decoder
    decoder->codec = avcodec_find_decoder(AV_CODEC_ID_H264);
    if (!decoder->codec) {
        fprintf(stderr, "H.264 decoder not found\n");
        free(decoder);
        return NULL;
    }

    // Allocate decoder context
    decoder->codec_ctx = avcodec_alloc_context3(decoder->codec);
    if (!decoder->codec_ctx) {
        fprintf(stderr, "Failed to allocate decoder context\n");
        free(decoder);
        return NULL;
    }

    // Set decoder parameters (used for realtime decoding)
    decoder->codec_ctx->flags |= AV_CODEC_FLAG_LOW_DELAY;
    decoder->codec_ctx->flags2 |= AV_CODEC_FLAG2_FAST;

    // Open decoder
    if (avcodec_open2(decoder->codec_ctx, decoder->codec, NULL) < 0) {
        fprintf(stderr, "Failed to open decoder\n");
        avcodec_free_context(&decoder->codec_ctx);
        free(decoder);
        return NULL;
    }

    // Allocate frame and packet
    decoder->frame = av_frame_alloc();
    decoder->pkt = av_packet_alloc();
    if (!decoder->frame || !decoder->pkt) {
        fprintf(stderr, "Failed to allocate frame or packet memory\n");
        h264_decoder_free(decoder);
        return NULL;
    }

    decoder->frame_count = 0;
    decoder->flush_sent = 0;
    
    // Initialize stream buffer (64KB)
    decoder->buffer_size = 64 * 1024;
    decoder->buffer = (uint8_t*)malloc(decoder->buffer_size);
    decoder->buffer_pos = 0;
    if (!decoder->buffer) {
        fprintf(stderr, "Failed to allocate stream buffer memory\n");
        h264_decoder_free(decoder);
        return NULL;
    }
    
    return decoder;
}

/**
 * Find the start position of the next NAL unit
 * NAL unit starts with 00 00 00 01 or 00 00 01
 * @return NAL start position, or -1 if not found
 */
static int find_nal_start(const uint8_t *data, int size, int start_pos) {
    for (int i = start_pos; i < size - 4; i++) {
        // Find 00 00 00 01
        if (data[i] == 0x00 && data[i+1] == 0x00 && 
            data[i+2] == 0x00 && data[i+3] == 0x01) {
            return i;
        }
        // Find 00 00 01 (but not in 00 00 00 01)
        if (i < size - 3 && data[i] == 0x00 && 
            data[i+1] == 0x00 && data[i+2] == 0x01 &&
            (i == 0 || data[i-1] != 0x00)) {
            return i;
        }
    }
    return -1;  // Not found
}

/**
 * Convert AVFrame to YUV format
 */
static int convert_frame_to_yuv(AVFrame *frame, int target_width, int target_height, YUVFrame **yuv_frame) {
    if (!frame || !yuv_frame) {
        return AVERROR(EINVAL);
    }

    int width = (target_width > 0) ? target_width : frame->width;
    int height = (target_height > 0) ? target_height : frame->height;

    YUVFrame *yuv = (YUVFrame*)calloc(1, sizeof(YUVFrame));
    if (!yuv) {
        return AVERROR(ENOMEM);
    }

    yuv->width = width;
    yuv->height = height;
    yuv->y_size = width * height;
    yuv->uv_size = (width / 2) * (height / 2);
    yuv->total_size = yuv->y_size + yuv->uv_size * 2;

    yuv->y_data = (uint8_t*)malloc(yuv->y_size);
    yuv->u_data = (uint8_t*)malloc(yuv->uv_size);
    yuv->v_data = (uint8_t*)malloc(yuv->uv_size);

    if (!yuv->y_data || !yuv->u_data || !yuv->v_data) {
        yuv_frame_free(yuv);
        return AVERROR(ENOMEM);
    }

    if (width == frame->width && height == frame->height) {
        // Size matches, copy directly
        // Copy Y plane
        for (int y = 0; y < height; y++) {
            memcpy(yuv->y_data + y * width,
                   frame->data[0] + y * frame->linesize[0],
                   width);
        }
        // Copy U plane
        for (int y = 0; y < height / 2; y++) {
            memcpy(yuv->u_data + y * (width / 2),
                   frame->data[1] + y * frame->linesize[1],
                   width / 2);
        }
        // Copy V plane
        for (int y = 0; y < height / 2; y++) {
            memcpy(yuv->v_data + y * (width / 2),
                   frame->data[2] + y * frame->linesize[2],
                   width / 2);
        }
    } else {
        // Need to scale (simplified version, using nearest neighbor sampling)
        // Scale Y plane
        for (int y = 0; y < height; y++) {
            int src_y = (y * frame->height) / height;
            for (int x = 0; x < width; x++) {
                int src_x = (x * frame->width) / width;
                yuv->y_data[y * width + x] = frame->data[0][src_y * frame->linesize[0] + src_x];
            }
        }
        // Scale U and V planes
        for (int y = 0; y < height / 2; y++) {
            int src_y = (y * frame->height / 2) / (height / 2);
            for (int x = 0; x < width / 2; x++) {
                int src_x = (x * frame->width / 2) / (width / 2);
                yuv->u_data[y * (width / 2) + x] = frame->data[1][src_y * frame->linesize[1] + src_x];
                yuv->v_data[y * (width / 2) + x] = frame->data[2][src_y * frame->linesize[2] + src_x];
            }
        }
    }

    *yuv_frame = yuv;
    return 0;
}

/**
 * Decode a H.264 NAL unit and output YUV format
 */
int h264_decode_frame(H264Decoder *decoder,
                      const uint8_t *nal_data,
                      int nal_size,
                      int width,
                      int height,
                      YUVFrame **yuv_frame) {
    int ret;

    if (!decoder || !yuv_frame) {
        return AVERROR(EINVAL);
    }

    *yuv_frame = NULL;

    // Send NAL unit to decoder
    if (nal_data && nal_size > 0) {
        decoder->pkt->data = (uint8_t*)nal_data;
        decoder->pkt->size = nal_size;
        
        ret = avcodec_send_packet(decoder->codec_ctx, decoder->pkt);
        if (ret < 0 && ret != AVERROR(EAGAIN)) {
            return ret;  // Error
        }
        // If EAGAIN is returned, it means the decoder needs to receive a frame first, continue to try to receive
    }

    // Receive decoded frame
    ret = avcodec_receive_frame(decoder->codec_ctx, decoder->frame);
    if (ret == AVERROR(EAGAIN)) {
        return 1;  // Need more data (decoder internal buffer)
    } else if (ret == AVERROR_EOF) {
        return AVERROR_EOF;  // Ended
    } else if (ret < 0) {
        return ret;  // Error
    }

    // Successfully decoded one frame, convert to YUV
    decoder->frame_count++;
    ret = convert_frame_to_yuv(decoder->frame, width, height, yuv_frame);
    if (ret < 0) {
        return ret;
    }

    return 0;
}

/**
 * Flush the decoder, get the remaining frames
 */
int h264_decode_frame_flush(H264Decoder *decoder,
                             int width,
                             int height,
                             YUVFrame **yuv_frame) {
    int ret;

    if (!decoder || !yuv_frame) {
        return AVERROR(EINVAL);
    }

    *yuv_frame = NULL;

    // Send flush packet (only send once)
    if (!decoder->flush_sent) {
        ret = avcodec_send_packet(decoder->codec_ctx, NULL);
        decoder->flush_sent = 1;
        if (ret < 0 && ret != AVERROR_EOF) {
            return ret;
        }
    }

    // Receive decoded frame
    ret = avcodec_receive_frame(decoder->codec_ctx, decoder->frame);
    if (ret == AVERROR_EOF) {
        return AVERROR_EOF;  // No more frames
    } else if (ret < 0) {
        return ret;  // Error
    }

    // Successfully decoded one frame, convert to YUV
    decoder->frame_count++;
    ret = convert_frame_to_yuv(decoder->frame, width, height, yuv_frame);
    if (ret < 0) {
        return ret;
    }

    return 0;
}

/**
* Decode H.264 stream data (automatically handle NAL unit splitting)
 * 
 * This function receives H.264 stream data of any size, automatically finds the NAL unit boundary and splits it,
 * then decodes each one. The user does not need to manually handle NAL unit splitting.
 * 
 * @param decoder decoder context
 * @param stream_data H.264 stream data (can be any size)
 * @param stream_size stream data size
 * @param width expected width (if 0, use the width detected by the decoder)
 * @param height expected height (if 0, use the height detected by the decoder)
 * @param yuv_frame output YUV frame data (caller is responsible for freeing memory, use yuv_frame_free)
 * @return 0 successfully decoded one frame, 1 needs more data (no complete NAL unit), <0 error
 */
int h264_decode_stream(H264Decoder *decoder,
                       const uint8_t *stream_data,
                       int stream_size,
                       int width,
                       int height,
                       YUVFrame **yuv_frame) {
    int ret;

    if (!decoder || !yuv_frame) {
        return AVERROR(EINVAL);
    }

    *yuv_frame = NULL;

    // If there is new data, append to the buffer
    if (stream_data && stream_size > 0) {
        // Append new data to the buffer
        if (decoder->buffer_pos + stream_size > decoder->buffer_size) {
            // Buffer too small, need to resize
            int new_size = decoder->buffer_size * 2;
            while (new_size < decoder->buffer_pos + stream_size) {
                new_size *= 2;
            }
            uint8_t *new_buffer = (uint8_t*)realloc(decoder->buffer, new_size);
            if (!new_buffer) {
                return AVERROR(ENOMEM);
            }
            decoder->buffer = new_buffer;
            decoder->buffer_size = new_size;
        }

        memcpy(decoder->buffer + decoder->buffer_pos, stream_data, stream_size);
        decoder->buffer_pos += stream_size;
    }

    // If there is no buffer data, return needs more data
    if (decoder->buffer_pos == 0) {
        return 1;  // Need more data
    }

    // Find and process complete NAL units
    int pos = 0;
    while (pos < decoder->buffer_pos) {
        // Find the start position of the next NAL unit
        int nal_start = find_nal_start(decoder->buffer, decoder->buffer_pos, pos);
        if (nal_start < 0) {
            // NAL start code not found, possible reasons:
            // 1. Data is incomplete (NAL unit crosses data packet transmission) - need to wait for more data
            // 2. Data has been processed completely
            // Keep data waiting for more data
            if (pos > 0) {
                // Move unprocessed data to the beginning of the buffer
                int remaining = decoder->buffer_pos - pos;
                memmove(decoder->buffer, decoder->buffer + pos, remaining);
                decoder->buffer_pos = remaining;
            }
            return 1;  // Need more data
        }

        // Find the next NAL unit, calculate the size of the current NAL unit
        int next_nal = find_nal_start(decoder->buffer, decoder->buffer_pos, nal_start + 4);
        int nal_size;
        if (next_nal < 0) {
            // This is the last NAL unit, use all remaining data
            nal_size = decoder->buffer_pos - nal_start;
        } else {
            nal_size = next_nal - nal_start;
        }

        // Decode current NAL unit
        ret = h264_decode_frame(decoder, 
                               decoder->buffer + nal_start, 
                               nal_size,
                               width, height, yuv_frame);

        if (ret == 0 && *yuv_frame) {
            // Successfully decoded one frame, clean up processed data
            pos = nal_start + nal_size;
            if (pos < decoder->buffer_pos) {
                // There is still unprocessed data, move to the beginning of the buffer
                int remaining = decoder->buffer_pos - pos;
                memmove(decoder->buffer, decoder->buffer + pos, remaining);
                decoder->buffer_pos = remaining;
            } else {
                // All data has been processed
                decoder->buffer_pos = 0;
            }
            return 0;  // Successfully decoded one frame
        } else if (ret == 1) {
            // Need more data (decoder internal buffer), continue to process the next NAL unit
            pos = nal_start + nal_size;
        } else if (ret < 0) {
            // Decoding error (可能是非视频帧的 NAL 单元，如 SPS/PPS/SEI）
            // Skip this NAL unit, continue to process the next NAL unit
            pos = nal_start + nal_size;
            // Do not return error, continue to try the next NAL unit
        } else {
            // Other cases, continue to process the next NAL unit
            pos = nal_start + nal_size;
        }
    }

    // All data has been processed, but no frame has been decoded
    // Clean up processed data
    if (pos > 0 && pos < decoder->buffer_pos) {
        int remaining = decoder->buffer_pos - pos;
        memmove(decoder->buffer, decoder->buffer + pos, remaining);
        decoder->buffer_pos = remaining;
    } else if (pos >= decoder->buffer_pos) {
        decoder->buffer_pos = 0;
    }
    
    return 1;  // Need more data
}

/**
 * Free YUV frame data
 */
void yuv_frame_free(YUVFrame *yuv_frame) {
    if (!yuv_frame) return;
    if (yuv_frame->y_data) free(yuv_frame->y_data);
    if (yuv_frame->u_data) free(yuv_frame->u_data);
    if (yuv_frame->v_data) free(yuv_frame->v_data);
    free(yuv_frame);
}

/**
 * Get the number of decoded frames
 */
int h264_decode_frame_get_count(H264Decoder *decoder) {
    if (!decoder) return 0;
    return decoder->frame_count;
}

/**
 * Free decoder resources
 */
void h264_decoder_free(H264Decoder *decoder) {
    if (!decoder) return;

    if (decoder->codec_ctx) {
        avcodec_free_context(&decoder->codec_ctx);
    }
    if (decoder->frame) {
        av_frame_free(&decoder->frame);
    }
    if (decoder->pkt) {
        av_packet_free(&decoder->pkt);
    }
    if (decoder->buffer) {
        free(decoder->buffer);
    }
    free(decoder);
}
