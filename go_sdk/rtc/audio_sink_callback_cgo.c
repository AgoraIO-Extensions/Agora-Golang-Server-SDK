#include "audio_sink_callback_cgo.h"

// Go callback 声明
extern int goOnSinkAudioFrame(void* sink, void* frame);

// C 包装函数实现
int cgo_onSinkAudioFrameCallback(AGORA_HANDLE agora_audio_sink, const audio_pcm_frame* frame) {
    return goOnSinkAudioFrame(agora_audio_sink, (void*)frame);
}

