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
	conn              *RtcConnection
}

func (conn *RtcConnection) NewRemoteVideoTrack(cRemoteVideoTrack unsafe.Pointer) *RemoteVideoTrack {
	return &RemoteVideoTrack{
		cRemoteVideoTrack: cRemoteVideoTrack,
		conn:              conn,
	}
}

func (track *RemoteVideoTrack) RegisterVideoEncodedImageReceiver(receiver *VideoEncodedImageReceiver) int {
	if track.cRemoteVideoTrack == nil {
		return -1
	}
	receiverInner := track.conn.newVideoEncodedImageReceiverInner(receiver)
	if receiverInner == nil {
		return -1
	}
	agoraService.remoteVideoRWMutex.Lock()
	agoraService.remoteEncodedVideoReceivers[receiverInner.cReceiver] = receiverInner
	agoraService.remoteVideoRWMutex.Unlock()
	return int(C.agora_remote_video_track_register_video_encoded_image_receiver(track.cRemoteVideoTrack, receiverInner.cReceiver))
}

func (track *RemoteVideoTrack) UnregisterVideoEncodedImageReceiver(receiver *VideoEncodedImageReceiver) int {
	if track.cRemoteVideoTrack == nil {
		return -1
	}
	var targetReceiver unsafe.Pointer = nil
	agoraService.remoteVideoRWMutex.RLock()
	for cReceiver, receiverInner := range agoraService.remoteEncodedVideoReceivers {
		if receiverInner.receiver == receiver {
			targetReceiver = cReceiver
			break
		}
	}
	agoraService.remoteVideoRWMutex.RUnlock()
	ret := -1
	if targetReceiver != nil {
		ret = int(C.agora_remote_video_track_unregister_video_encoded_image_receiver(track.cRemoteVideoTrack, targetReceiver))
		agoraService.remoteVideoRWMutex.Lock()
		delete(agoraService.remoteEncodedVideoReceivers, targetReceiver)
		agoraService.remoteVideoRWMutex.Unlock()
	}
	return ret
}
