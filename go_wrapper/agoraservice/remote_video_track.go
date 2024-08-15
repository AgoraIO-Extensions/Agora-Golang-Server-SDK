package agoraservice

import "unsafe"

type VideoTrackInfo struct {
	IsLocal             bool
	OwnerUid            uint
	TrackId             uint
	ChannelId           string
	StreamType          int
	CodecType           int
	EncodedFrameOnly    bool
	SourceType          int
	ObservationPosition uint
}

type RemoteVideoTrack struct {
	cRemoteVideoTrack unsafe.Pointer
}

func NewRemoteVideoTrack(cRemoteVideoTrack unsafe.Pointer) *RemoteVideoTrack {
	return &RemoteVideoTrack{
		cRemoteVideoTrack: cRemoteVideoTrack,
	}
}
