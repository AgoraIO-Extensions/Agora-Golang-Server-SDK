#pragma once

#include "agora_media_node_factory.h"

// Audio sink callback wrapper
extern int cgo_onSinkAudioFrameCallback(AGORA_HANDLE agora_audio_sink, const audio_pcm_frame* frame);

