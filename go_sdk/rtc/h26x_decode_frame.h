/**
 * h264_decode_frame.h
 * H.264 realtime decoder header file
 * Function: decode H.264 NAL unit to YUV format (used for realtime stream decoding)
 */
#ifndef H26X_DECODE_FRAME_H
#define H26X_DECODE_FRAME_H

#include <stdint.h>

typedef struct H264Decoder H264Decoder;
/**
 * YUV frame data structure
 */
typedef struct {
    uint8_t *y_data;      // Y plane data (brightness)
    uint8_t *u_data;      // U plane data (chroma)
    uint8_t *v_data;      // V plane data (chroma)
    int width;            // width
    int height;           // height
    int y_size;           // Y plane size (width * height)
    int uv_size;          // U/V plane size (width/2 * height/2)
    int total_size;       // total size (y_size + uv_size * 2)
} YUVFrame;

H264Decoder* h264_decoder_init(void);

/**
 * Decode H.264 stream data 
 * 
 * This function receives H.264 stream data of any size, automatically finds the NAL unit boundary and splits it,
 * then decodes each one.
 * 
 * @param decoder decoder context
 * @param stream_data H.264 stream data (can be any size, NULL means no new data, only process buffer)
 * @param stream_size stream data size (ignore when stream_data is NULL)
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
                       YUVFrame **yuv_frame);

/**
 * Decode a H.264 NAL unit and output YUV format
 * 
 * Note: the input h264_data must be a complete NAL unit, including the start code (00 00 00 01 or 00 00 01)
 * 
 * @param decoder decoder context
 * @param nal_data NAL unit data (must include the start code)
 * @param nal_size NAL unit size (including the start code)
 * @param width expected width (if 0, use the width detected by the decoder)
 * @param height expected height (if 0, use the height detected by the decoder)
 * @param yuv_frame output YUV frame data (caller is responsible for freeing memory, use yuv_frame_free)
 * @return 0 successfully decoded one frame, 1 needs more data (decoder internal buffer), <0 error
 */
 int h264_decode_frame(H264Decoder *decoder,
                       const uint8_t *nal_data,
                       int nal_size,
                       int width,
                       int height,
                       YUVFrame **yuv_frame);

/**
 * Flush the decoder, get the remaining frames
 * Call this when the stream ends, to get the last frame from the decoder internal buffer
 * 
 * @param decoder decoder context
 * @param width expected width (if 0, use the width detected by the decoder)
 * @param height expected height (if 0, use the height detected by the decoder)
 * @param yuv_frame output YUV frame data (caller is responsible for freeing memory, use yuv_frame_free)
 * @return 0 successfully decoded one frame, AVERROR_EOF no more frames, <0 error
 */
int h264_decode_frame_flush(H264Decoder *decoder,
                             int width,
                             int height,
                             YUVFrame **yuv_frame);

/**
 * Free YUV frame data
 * @param yuv_frame YUV frame data
 */
void yuv_frame_free(YUVFrame *yuv_frame);

/**
 * Get the number of decoded frames
 * @param decoder decoder context
 * @return the number of decoded frames
 */
int h264_decode_frame_get_count(H264Decoder *decoder);

/**
* Free the decoder resources
 * @param decoder decoder context
 */
void h264_decoder_free(H264Decoder *decoder);

#endif // H264_DECODE_FRAME_H
