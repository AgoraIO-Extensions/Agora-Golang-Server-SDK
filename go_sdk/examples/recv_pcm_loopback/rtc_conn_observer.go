package main

import (
	"fmt"

	agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/agoraservice"
)

var _ agoraservice.RtcConnectionObserver = (*RtcConnectionObserverImpl)(nil)

type RtcConnectionObserverImpl struct {
	conSignal          chan struct{}
	OnDisconnectedSign chan struct{}
}

func NewRtcConnectionObserverImpl() *RtcConnectionObserverImpl {
	return &RtcConnectionObserverImpl{
		conSignal:          make(chan struct{}),
		OnDisconnectedSign: make(chan struct{}),
	}
}

func (r *RtcConnectionObserverImpl) WaitConnected() error {
	<-r.conSignal
	return nil
}

func (r *RtcConnectionObserverImpl) WaitDisconnected() error {
	<-r.OnDisconnectedSign
	return nil
}

// OnConnected implements agoraservice.RtcConnectionObserver.
func (r *RtcConnectionObserverImpl) OnConnected(con *agoraservice.RtcConnection, conInfo *agoraservice.RtcConnectionInfo, reason int) {
	// do something
	fmt.Printf("Connected, reason %d\n", reason)
	r.conSignal <- struct{}{}
}

// OnConnecting implements agoraservice.RtcConnectionObserver.
func (r *RtcConnectionObserverImpl) OnConnecting(con *agoraservice.RtcConnection, conInfo *agoraservice.RtcConnectionInfo, reason int) {
	fmt.Printf("Connecting, reason %d\n", reason)
}

// OnConnectionFailure implements agoraservice.RtcConnectionObserver.
func (r *RtcConnectionObserverImpl) OnConnectionFailure(con *agoraservice.RtcConnection, conInfo *agoraservice.RtcConnectionInfo, errCode int) {
	fmt.Printf("Connection failure, error code %d\n", errCode)
}

// OnConnectionLost implements agoraservice.RtcConnectionObserver.
func (r *RtcConnectionObserverImpl) OnConnectionLost(con *agoraservice.RtcConnection, conInfo *agoraservice.RtcConnectionInfo) {
	fmt.Printf("Connection lost\n")
}

// OnDisconnected implements agoraservice.RtcConnectionObserver.
func (r *RtcConnectionObserverImpl) OnDisconnected(con *agoraservice.RtcConnection, conInfo *agoraservice.RtcConnectionInfo, reason int) {
	// do something
	fmt.Printf("Disconnected, reason %d\n", reason)
	r.OnDisconnectedSign <- struct{}{}
}

// OnEncryptionError implements agoraservice.RtcConnectionObserver.
func (r *RtcConnectionObserverImpl) OnEncryptionError(con *agoraservice.RtcConnection, err int) {
}

// OnError implements agoraservice.RtcConnectionObserver.
func (r *RtcConnectionObserverImpl) OnError(con *agoraservice.RtcConnection, err int, msg string) {
}

// OnReconnected implements agoraservice.RtcConnectionObserver.
func (r *RtcConnectionObserverImpl) OnReconnected(con *agoraservice.RtcConnection, conInfo *agoraservice.RtcConnectionInfo, reason int) {
	fmt.Printf("Reconnected, reason %d\n", reason)
}

// OnReconnecting implements agoraservice.RtcConnectionObserver.
func (r *RtcConnectionObserverImpl) OnReconnecting(con *agoraservice.RtcConnection, conInfo *agoraservice.RtcConnectionInfo, reason int) {
	fmt.Printf("Reconnecting, reason %d\n", reason)
}

// OnStreamMessageError implements agoraservice.RtcConnectionObserver.
func (r *RtcConnectionObserverImpl) OnStreamMessageError(con *agoraservice.RtcConnection, uid string, streamId int, errCode int, missed int, cached int) {
}

// OnTokenPrivilegeDidExpire implements agoraservice.RtcConnectionObserver.
func (r *RtcConnectionObserverImpl) OnTokenPrivilegeDidExpire(con *agoraservice.RtcConnection) {
}

// OnTokenPrivilegeWillExpire implements agoraservice.RtcConnectionObserver.
func (r *RtcConnectionObserverImpl) OnTokenPrivilegeWillExpire(con *agoraservice.RtcConnection, token string) {
}

// OnUserJoined implements agoraservice.RtcConnectionObserver.
func (r *RtcConnectionObserverImpl) OnUserJoined(con *agoraservice.RtcConnection, uid string) {
	fmt.Println("user joined, " + uid)
}

// OnUserLeft implements agoraservice.RtcConnectionObserver.
func (r *RtcConnectionObserverImpl) OnUserLeft(con *agoraservice.RtcConnection, uid string, reason int) {
	fmt.Printf("user left: %s, reason = %d\n", uid, reason)
}
