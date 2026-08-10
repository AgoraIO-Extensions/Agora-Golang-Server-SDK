//go:build !vad_uap
// +build !vad_uap

package agoraservice

func VadUAPSupported() bool {
	return false
}

func NewAudioVad(cfg *AudioVadConfig) *AudioVad {
	return nil
}

func (vad *AudioVad) Release() {
}

func (vad *AudioVad) ProcessPcmFrame(frame *AudioFrame) (*AudioFrame, int) {
	return nil, -1
}
