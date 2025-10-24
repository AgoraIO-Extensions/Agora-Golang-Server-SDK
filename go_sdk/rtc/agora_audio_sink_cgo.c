#include "agora_audio_sink_cgo.h"

extern int goOnAudioFrame(void* agora_audio_sink, const audio_pcm_frame* frame);

audio_sink* create_audio_sink_callbacks() {
    audio_sink* sink = (audio_sink*)malloc(sizeof(audio_sink));
    if (sink) {
        memset(sink, 0, sizeof(audio_sink));
        sink->on_audio_frame = goOnAudioFrame;
    }
    return sink;
}


