package agoraservice

import "unsafe"

type RemoteAudioTrack struct {
	cRemoteAudioTrack unsafe.Pointer
}

func NewRemoteAudioTrack(cRemoteAudioTrack unsafe.Pointer) *RemoteAudioTrack {
	return &RemoteAudioTrack{
		cRemoteAudioTrack: cRemoteAudioTrack,
	}
}
