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

func (track *RemoteVideoTrack) RegisterVideoEncodedImageReceiver() int {
	if track.cRemoteVideoTrack == nil {
		return -1
	}
	return int(C.agora_remote_video_track_register_video_encoded_image_receiver(track.cRemoteVideoTrack))
}
