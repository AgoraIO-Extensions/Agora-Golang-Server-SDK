package agoraservice

// #cgo CFLAGS: -I${SRCDIR}/../../agora_sdk/include/c/api2 -I${SRCDIR}/../../agora_sdk/include/c/base
// #include <string.h>
// #include "agora_service.h"
// #include "agora_media_node_factory.h"
// #include "agora_video_track.h"
import "C"
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

// func (track *RemoteVideoTrack) RegisterVideoEncodedImageReceiver(receiver *VideoEncodedImageReceiver) int {
// 	if track.cRemoteVideoTrack == nil {
// 		return -1
// 	}
// 	receiverInner := newVideoEncodedImageReceiverInner(receiver)
// 	if receiverInner == nil {
// 		return -1
// 	}
// 	// add receiver in connection
// 	track.conn.remoteVideoRWMutex.Lock()
// 	// the key is for unregister, remoteEncodedVideoReceivers is for release
// 	track.conn.remoteEncodedVideoReceivers[receiver] = receiverInner
// 	track.conn.remoteVideoRWMutex.Unlock()
// 	// add receiver in service
// 	agoraService.remoteVideoRWMutex.Lock()
// 	// the key is for callback
// 	agoraService.remoteEncodedVideoReceivers[receiverInner.cReceiver] = receiverInner
// 	agoraService.remoteVideoRWMutex.Unlock()
// 	// register receiver
// 	return int(C.agora_remote_video_track_register_video_encoded_image_receiver(track.cRemoteVideoTrack, receiverInner.cReceiver))
// }

// func (track *RemoteVideoTrack) UnregisterVideoEncodedImageReceiver(receiver *VideoEncodedImageReceiver) int {
// 	if track.cRemoteVideoTrack == nil {
// 		return -1
// 	}
// 	// find receiver in connection
// 	track.conn.remoteVideoRWMutex.RLock()
// 	receiverInner, ok := track.conn.remoteEncodedVideoReceivers[receiver]
// 	track.conn.remoteVideoRWMutex.RUnlock()
// 	if !ok {
// 		return -1
// 	}
// 	// unregister receiver
// 	ret := int(C.agora_remote_video_track_unregister_video_encoded_image_receiver(track.cRemoteVideoTrack, receiverInner.cReceiver))
// 	// remove receiver in service
// 	agoraService.remoteVideoRWMutex.Lock()
// 	delete(agoraService.remoteEncodedVideoReceivers, receiverInner.cReceiver)
// 	agoraService.remoteVideoRWMutex.Unlock()
// 	// remove receiver in connection
// 	track.conn.remoteVideoRWMutex.Lock()
// 	delete(track.conn.remoteEncodedVideoReceivers, receiver)
// 	track.conn.remoteVideoRWMutex.Unlock()
// 	// release receiver
// 	receiverInner.release()
// 	return ret
// }
