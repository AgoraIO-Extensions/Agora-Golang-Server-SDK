package agoraservice

import "unsafe"

type RemoteVideoTrack struct {
	cRemoteVideoTrack unsafe.Pointer
}

func NewRemoteVideoTrack(cRemoteVideoTrack unsafe.Pointer) *RemoteVideoTrack {
	return &RemoteVideoTrack{
		cRemoteVideoTrack: cRemoteVideoTrack,
	}
}
