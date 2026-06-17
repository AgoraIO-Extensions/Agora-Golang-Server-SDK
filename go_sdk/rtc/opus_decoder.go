package agoraservice

// #cgo pkg-config: libavformat libavcodec libavutil libswresample
// #include <string.h>
// #include <stdlib.h>
// #include "opus_decode.h"
import "C"
import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
	"unsafe"

	agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/rtc"

	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
)

// output PCM params (opus is always decoded at 48k internally, then resampled here)
const (
	outSampleRate = 16000
	outChannels   = 1
)

// opusReceiver decodes per-user opus streams into PCM and writes each to a file.
// A separate decoder is kept per uid because the opus decoder carries internal state.
type opusReceiver struct {
	mu       sync.Mutex
	decoders map[string]unsafe.Pointer
	files    map[string]*os.File
}

func newOpusReceiver() *opusReceiver {
	return &opusReceiver{
		decoders: make(map[string]unsafe.Pointer),
		files:    make(map[string]*os.File),
	}
}

// onEncodedAudioFrame is the OnEncodedAudioFrameReceived callback body.
// packet is a raw opus payload, exactly what decode_opus expects.
func (r *opusReceiver) onEncodedAudioFrame(uid string, packet []byte, sendTs int64, codec int) {
	// only handle OPUS (codec == 1); skip PCMA/PCMU/AAC/...
	if codec != int(agoraservice.AudioCodecOpus) {
		fmt.Printf("skip non-opus frame from %s, codec %d\n", uid, codec)
		return
	}
	if len(packet) == 0 {
		return
	}

	r.mu.Lock()
	dec, ok := r.decoders[uid]
	if !ok {
		dec = C.open_opus_decoder(C.int(1), C.int(outSampleRate), C.int(outChannels))
		r.decoders[uid] = dec
		f, err := os.Create(fmt.Sprintf("opus_%s.pcm", uid))
		if err != nil {
			fmt.Printf("create pcm file failed for %s: %v\n", uid, err)
		}
		r.files[uid] = f
		fmt.Printf("opus decoder created for uid %s\n", uid)
	}
	f := r.files[uid]
	r.mu.Unlock()

	if dec == nil {
		return
	}

	cFrame := C.struct__MediaFrame{}
	C.memset(unsafe.Pointer(&cFrame), 0, C.sizeof_struct__MediaFrame)

	ret := C.decode_opus(dec,
		(*C.uint8_t)(unsafe.Pointer(&packet[0])),
		C.int(len(packet)),
		&cFrame)
	if ret != 0 {
		// AVERROR(EAGAIN) etc., just wait for more packets
		return
	}

	// copy out the PCM; cFrame.buffer is reused on the next decode call
	pcm := C.GoBytes(unsafe.Pointer(cFrame.buffer), C.int(cFrame.buffer_size))
	if f != nil {
		f.Write(pcm)
	}
	fmt.Printf("decoded opus from %s: ts=%d samples=%d rate=%d ch=%d bytes=%d\n",
		uid, sendTs, int(cFrame.samples), int(cFrame.sample_rate),
		int(cFrame.channels), int(cFrame.buffer_size))
}

func (r *opusReceiver) release(uid string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if dec, ok := r.decoders[uid]; ok {
		C.close_opus_decoder(dec)
		delete(r.decoders, uid)
	}
	if f, ok := r.files[uid]; ok {
		if f != nil {
			f.Close()
		}
		delete(r.files, uid)
	}
}

func (r *opusReceiver) releaseAll() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for uid, dec := range r.decoders {
		C.close_opus_decoder(dec)
		delete(r.decoders, uid)
	}
	for uid, f := range r.files {
		if f != nil {
			f.Close()
		}
		delete(r.files, uid)
	}
}

func main() {
	bStop := new(bool)
	*bStop = false
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		*bStop = true
		fmt.Println("Application terminated")
	}()

	argus := os.Args
	if len(argus) < 3 {
		fmt.Println("Usage: recv_opus <appid> <channel_name>")
		return
	}
	appid := argus[1]
	channelName := argus[2]

	cert := os.Getenv("AGORA_APP_CERTIFICATE")
	userId := "0"
	if appid == "" {
		fmt.Println("Please set appid, and AGORA_APP_CERTIFICATE if needed")
		return
	}
	token := ""
	if cert != "" {
		tokenExpirationInSeconds := uint32(3600)
		privilegeExpirationInSeconds := uint32(3600)
		var err error
		token, err = rtctokenbuilder.BuildTokenWithUserAccount(appid, cert, channelName, userId,
			rtctokenbuilder.RolePublisher, tokenExpirationInSeconds, privilegeExpirationInSeconds)
		if err != nil {
			fmt.Println("Failed to build token: ", err)
			return
		}
	}

	svcCfg := agoraservice.NewAgoraServiceConfig()
	svcCfg.AppId = appid
	agoraservice.Initialize(svcCfg)
	defer agoraservice.Release()

	receiver := newOpusReceiver()
	defer receiver.releaseAll()

	conSignal := make(chan struct{}, 1)
	conHandler := agoraservice.RtcConnectionObserver{
		OnConnected: func(con *agoraservice.RtcConnection, info *agoraservice.RtcConnectionInfo, reason int) {
			fmt.Println("Connected")
			select {
			case conSignal <- struct{}{}:
			default:
			}
		},
		OnDisconnected: func(con *agoraservice.RtcConnection, info *agoraservice.RtcConnectionInfo, reason int) {
			fmt.Println("Disconnected")
		},
		OnUserJoined: func(con *agoraservice.RtcConnection, uid string) {
			fmt.Println("user joined, " + uid)
		},
		OnUserLeft: func(con *agoraservice.RtcConnection, uid string, reason int) {
			fmt.Println("user left, " + uid)
			receiver.release(uid)
		},
	}

	conCfg := &agoraservice.RtcConnectionConfig{
		AutoSubscribeAudio: true,
		AutoSubscribeVideo: false,
		ClientRole:         agoraservice.ClientRoleBroadcaster,
		ChannelProfile:     agoraservice.ChannelProfileLiveBroadcasting,
	}
	// receive-only: default publish config (publishes nothing)
	publishConfig := agoraservice.NewRtcConPublishConfig()
	con := agoraservice.NewRtcConnection(conCfg, publishConfig)
	defer con.Release()

	con.RegisterObserver(&conHandler)

	// register the encoded audio frame observer to receive raw opus packets
	encObserver := &agoraservice.AudioEncodedFrameObserver{
		OnEncodedAudioFrameReceived: receiver.onEncodedAudioFrame,
	}
	con.RegisterAudioEncodedFrameObserver(encObserver)

	if ret := con.Connect(token, channelName, userId, ""); ret != 0 {
		fmt.Printf("Connect failed, ret %d\n", ret)
		return
	}
	<-conSignal

	fmt.Println("Waiting for opus encoded audio... press Ctrl+C to stop")
	for !(*bStop) {
		// the callback runs on the SDK thread; just keep the process alive
		time.Sleep(100 * time.Millisecond)
	}

	con.Disconnect()
	fmt.Println("done, pcm written to opus_<uid>.pcm")
}
