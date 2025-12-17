package agoraservice

/*
#cgo CFLAGS: -I../../agora_sdk/include/c/api2 -I../../agora_sdk/include/c/base
#include <stdlib.h>
#include <string.h>
#include "agora_local_user.h"
#include "agora_rtc_conn.h"
#include "agora_service.h"
#include "agora_media_base.h"
#include "agora_parameter.h"
*/
import "C"
import (
	"fmt"
	"strconv"
	"time"
	"unsafe"
	//"sync/atomic"
)

type RtcConnectionInfo struct {
	ConnectionId uint
	/**
	 * ID of the target channel. NULL if you did not call the connect
	 * method.
	 */
	ChannelId string
	/**
	 * The state of the current connection: #CONNECTION_STATE_TYPE.
	 */
	State int
	/**
	 * ID of the local user.
	 */
	LocalUserId string
	/**
	 * Internal use only.
	 */
	InternalUid uint
}

type AudioFrameObserverAudioParams struct {
	SampleRate     int
	Channels       int
	Mode           RawAudioFrameOpModeType
	SamplesPerCall int
}

type RtcConnectionObserver struct {
	OnConnected                func(con *RtcConnection, conInfo *RtcConnectionInfo, reason int)
	OnDisconnected             func(con *RtcConnection, conInfo *RtcConnectionInfo, reason int)
	OnConnecting               func(con *RtcConnection, conInfo *RtcConnectionInfo, reason int)
	OnReconnecting             func(con *RtcConnection, conInfo *RtcConnectionInfo, reason int)
	OnReconnected              func(con *RtcConnection, conInfo *RtcConnectionInfo, reason int)
	OnConnectionLost           func(con *RtcConnection, conInfo *RtcConnectionInfo)
	OnConnectionFailure        func(con *RtcConnection, conInfo *RtcConnectionInfo, errCode int)
	OnTokenPrivilegeWillExpire func(con *RtcConnection, token string)
	OnTokenPrivilegeDidExpire  func(con *RtcConnection)
	OnUserJoined               func(con *RtcConnection, uid string)
	OnUserLeft                 func(con *RtcConnection, uid string, reason int)
	OnError                    func(con *RtcConnection, err int, msg string)
	OnStreamMessageError       func(con *RtcConnection, uid string, streamId int, errCode int, missed int, cached int)
	OnEncryptionError          func(con *RtcConnection, err int)
	//date: 2025-06-27
	// Triggered when the following two conditions are both met:
	// 1. The developer sets the connection's scenario to AudioScenarioAiServer.
	// 2. The version of the SDK on the client side does not support aiqos.
	// If triggered, it means the client-side version does not support aiqos, and the developer needs to decide to reset the server-side scenario.
	// How should the developer handle it when the callback is triggered?
	// 1. If set return value to -1, it means the SDK internally does not handle the scenario incompatibility.
	// 2. If set return value to a valid scenario, it means the SDK internally automatically falls back to the scenario returned, ensuring compatibility.
	// how to use it: can ref to examples/ai_send_recv_pcm/ai_send_recv_pcm.go
	OnAIQoSCapabilityMissing func(con *RtcConnection, defaultFallbackSenario int) int
}

// struct for local audio track statistics
type LocalAudioTrackStats struct {
	/**
	 * The number of channels.
	 */
	NumChannels int
	/**
	 * The sample rate (Hz).
	 */
	SentSampleRate int
	/**
	 * The average sending bitrate (Kbps).
	 */
	SentBitrate int
	/**
	 * The internal payload type
	 */
	InternalCodec int
	/**
	 * Voice pitch frequency in Hz
	 */
	VoicePitch float64
}
type RemoteAudioTrackStats struct {
	/**
	 * User ID of the remote user sending the audio streams.
	 */
	Uid uint
	/**
	 * Audio quality received by the user: #QUALITY_TYPE.
	 */
	Quality int
	/**
	 * @return Network delay (ms) from the sender to the receiver.
	 */
	NetworkTransportDelay int
	/**
	 * @return Delay (ms) from the receiver to the jitter buffer.
	 */
	JitterBufferDelay int
	/**
	 * The audio frame loss rate in the reported interval.
	 */
	AudioLossRate int
	/**
	 * The number of channels.
	 */
	NumChannels int
	/**
	 * The sample rate (Hz) of the received audio stream in the reported interval.
	 */
	ReceivedSampleRate int
	/**
	 * The average bitrate (Kbps) of the received audio stream in the reported interval.
	 * */
	ReceivedBitrate int
	/**
	 * The total freeze time (ms) of the remote audio stream after the remote user joins the channel.
	 * In a session, audio freeze occurs when the audio frame loss rate reaches 4%.
	 * Agora uses 2 seconds as an audio piece unit to calculate the audio freeze time.
	 * The total audio freeze time = The audio freeze number &times; 2 seconds
	 */
	TotalFrozenTime int
	/**
	 * The total audio freeze time as a percentage (%) of the total time when the audio is available.
	 * */
	FrozenRate int
	/**
	 * The number of audio bytes received.
	 */
	ReceivedBytes int64
	/**
	 * The MOS value of the received audio stream.
	 */
	MosValue int
	/**
	 * The total time (ms) when the remote us	er neither stops sending the audio
	 * stream nor disables the audio module after joining the channel.
	 */
	TotalActiveTime int
	/**
	 * The total publish duration (ms) of the remote audio stream.
	 */
	PublishDuration int
}
type LocalVideoTrackStats struct {
	NumberOfStreams           uint64
	BytesMajorStream          uint64
	BytesMinorStream          uint64
	FramesEncoded             uint32
	SSRCMajorStream           uint32
	SSRCMinorStream           uint32
	CaptureFrameRate          int
	RegulatedCaptureFrameRate int
	InputFrameRate            int
	EncodeFrameRate           int
	RenderFrameRate           int
	TargetMediaBitrateBps     int
	MediaBitrateBps           int
	TotalBitrateBps           int
	CaptureWidth              int
	CaptureHeight             int
	RegulatedCaptureWidth     int
	RegulatedCaptureHeight    int
	Width                     int
	Height                    int
	EncoderType               uint32
	UplinkCostTimeMs          int
	QualityAdaptIndication    int
}

type RemoteVideoTrackStats struct {
	/**
	  User ID of the remote user sending the video streams.
	*/
	Uid uint
	/** **DEPRECATED** Time delay (ms).
	 */
	Delay int
	/**
	  Width (pixels) of the video stream.
	*/
	Width int
	/**
	Height (pixels) of the video stream.
	*/
	Height int
	/**
	Bitrate (Kbps) received since the last count.
	*/
	ReceivedBitrate int
	/** The decoder output frame rate (fps) of the remote video.
	 */
	DecoderOutputFrameRate int
	/** The render output frame rate (fps) of the remote video.
	 */
	RendererOutputFrameRate int
	/** The video frame loss rate (%) of the remote video stream in the reported interval.
	 */
	FrameLossRate int
	/** Packet loss rate (%) of the remote video stream after using the anti-packet-loss method.
	 */
	PacketLossRate int
	RxStreamType   int
	/**
	  The total freeze time (ms) of the remote video stream after the remote user joins the channel.
	  In a video session where the frame rate is set to no less than 5 fps, video freeze occurs when
	  the time interval between two adjacent renderable video frames is more than 500 ms.
	*/
	TotalFrozenTime int
	/**
	  The total video freeze time as a percentage (%) of the total time when the video is available.
	*/
	FrozenRate int
	/**
	  The total video decoded frames.
	*/
	TotalDecodedFrames uint32
	/**
	   The offset (ms) between audio and video stream. A positive value indicates the audio leads the
	  video, and a negative value indicates the audio lags the video.
	*/
	AvSyncTimeMs int
	/**
	   The average offset(ms) between receive first packet which composite the frame and  the frame
	  ready to render.
	*/
	DownlinkProcessTimeMs int
	/**
	  The average time cost in renderer.
	*/
	FrameRenderDelayMs int
	/**
	   The total time (ms) when the remote user neither stops sending the video
	  stream nor disables the video module after joining the channel.
	*/
	TotalActiveTime uint64
	/**
	  The total publish duration (ms) of the remote video stream.
	*/
	PublishDuration uint64
}

type LocalUserObserver struct {
	OnStreamMessage func(localUser *LocalUser, uid string, streamId int, data []byte)
	// userMediaInfo: UserMediaInfoXxx
	// val: 0 for false, 1 for true
	OnUserInfoUpdated          func(localUser *LocalUser, uid string, userMediaInfo int, val int)
	OnUserAudioTrackSubscribed func(localUser *LocalUser, uid string, remoteAudioTrack *RemoteAudioTrack)
	OnUserVideoTrackSubscribed func(localUser *LocalUser, uid string, info *VideoTrackInfo, remoteVideoTrack *RemoteVideoTrack)
	/*
		for Mute/Unmute
		(state== 0 and reason == 5): mute
		(state== 2 and reason == 6): unmute
	*/
	OnUserAudioTrackStateChanged func(localUser *LocalUser, uid string, remoteAudioTrack *RemoteAudioTrack, state int, reason int, elapsed int)
	OnUserVideoTrackStateChanged func(localUser *LocalUser, uid string, remoteAudioTrack *RemoteVideoTrack, state int, reason int, elapsed int)
	//  void (*on_audio_publish_state_changed)(AGORA_HANDLE agora_local_user, const char* channel, int old_state, int new_state, int elapse_since_last_state);
	OnAudioPublishStateChanged func(localUser *LocalUser, channelId string, oldState int, newState int, elapsed int)
	OnAudioVolumeIndication    func(localUser *LocalUser, audioVolumeInfo []*AudioVolumeInfo, speakerNumber int, totalVolume int)
	OnAudioMetaDataReceived    func(localUser *LocalUser, uid string, metaData []byte)
	// for version 2.2.2
	OnLocalAudioTrackStatistics  func(localUser *LocalUser, stats *LocalAudioTrackStats)
	OnRemoteAudioTrackStatistics func(localUser *LocalUser, uid string, stats *RemoteAudioTrackStats)
	OnLocalVideoTrackStatistics  func(localUser *LocalUser, stats *LocalVideoTrackStats)
	OnRemoteVideoTrackStatistics func(localUser *LocalUser, uid string, stats *RemoteVideoTrackStats)
	// added on 2025-06-09
	OnAudioTrackPublishSuccess func(localUser *LocalUser, audioTrack *LocalAudioTrack)
	OnAudioTrackUnpublished    func(localUser *LocalUser, audioTrack *LocalAudioTrack)
}

type AudioFrameObserver struct {
	OnRecordAudioFrame                func(localUser *LocalUser, channelId string, frame *AudioFrame) bool
	OnPlaybackAudioFrame              func(localUser *LocalUser, channelId string, frame *AudioFrame) bool
	OnMixedAudioFrame                 func(localUser *LocalUser, channelId string, frame *AudioFrame) bool
	OnEarMonitoringAudioFrame         func(localUser *LocalUser, frame *AudioFrame) bool
	OnPlaybackAudioFrameBeforeMixing  func(localUser *LocalUser, channelId string, uid string, frame *AudioFrame, vadResultStat VadState, vadResultFrame *AudioFrame) bool
	OnGetAudioFramePosition           func(localUser *LocalUser) int
	OnGetPlaybackAudioFrameParam      func(localUser *LocalUser) AudioFrameObserverAudioParams
	OnGetRecordAudioFrameParam        func(localUser *LocalUser) AudioFrameObserverAudioParams
	OnGetMixedAudioFrameParam         func(localUser *LocalUser) AudioFrameObserverAudioParams
	OnGetEarMonitoringAudioFrameParam func(localUser *LocalUser) AudioFrameObserverAudioParams
}

type VideoFrameObserver struct {
	OnFrame func(channelId string, userId string, frame *VideoFrame) bool
}

type VideoEncodedFrameObserver struct {
	OnEncodedVideoFrame func(uid string, imageBuffer []byte,
		frameInfo *EncodedVideoFrameInfo) bool
}

type AudioEncoderConfiguration struct {
	AudioProfile int
}

// type SubscribeAudioConfig struct {
// 	SampleRate int
// 	Channels   int
// }

/**
 * Configurations for the RTC connection.
 */
type RtcConnectionConfig struct {
	/**
	 * Determines whether to subscribe to all audio streams automatically.
	 * - true: (Default) Subscribe to all audio streams automatically.
	 * - false: Do not subscribe to any audio stream automatically.
	 */
	AutoSubscribeAudio bool
	/**
	 * Determines whether to subscribe to all video streams automatically.
	 * - true: (Default) Subscribe to all video streams automatically.
	 * - false: Do not subscribe to any video stream automatically.
	 */
	AutoSubscribeVideo bool
	/**
	 * Determines whether to enable audio recording or playout.
	 * - true: It's used to publish audio and mix microphone, or subscribe audio and playout
	 * - false: It's used to publish extenal audio frame only without mixing microphone, or no need audio device to playout audio either
	 */
	EnableAudioRecordingOrPlayout bool
	/**
	 * The maximum sending bitrate.
	 */
	MaxSendBitrate int
	/**
	 * The minimum port.
	 */
	MinPort int
	/**
	 * The maximum port.
	 */
	MaxPort int
	/**
	 * The role of the user. The default user role is ClientRoleAudience.
	 */
	ClientRole     ClientRole
	ChannelProfile ChannelProfile
	/**
	 * Determines whether to receive audio media packet or not.
	 */
	AudioRecvMediaPacket bool
	/**
	 * Determines whether to receive video media packet or not.
	 */
	VideoRecvMediaPacket bool
}

type RtcConnectionPublishConfig struct {
	AudioProfile                   AudioProfile
	AudioScenario                  AudioScenario
	IsPublishAudio                 bool             //default to true
	IsPublishVideo                 bool             // default to false
	AudioPublishType               AudioPublishType // 0: no publish, 1: pcm, 2: encoded pcm. default to 1
	VideoPublishType               VideoPublishType // 0: no publish, 1: yuv, 2: encoded image. default to 0
	VideoEncodedImageSenderOptions *VideoEncodedImageSenderOptions
	//only for send external audio and limited send speek for ai scenario . default to false
	// if want to use this feature, should set the IsSendExternalAudioForAI to true, and set the UnlimitedMs and Speed
	// and consult tech support for the details.
	// for this case, recommend to ExternalAudioSendMs to 500, and Speed to 2
	SendExternalAudioParameters *SendExternalAudioParameters
}
// note: DeliverMuteDataForFakeAdm can only set to rtc engine level, can not
// set to connection level
// so if once a connection has set to true, wihich will affect all the connections,
type SendExternalAudioParameters struct {
	Enabled                   bool
	SendMs                    int
	SendSpeed                 int
	DeliverMuteDataForFakeAdm bool
}

type RtcConnection struct {
	cConnection unsafe.Pointer
	connInfo    RtcConnectionInfo
	localUser   *LocalUser
	parameter   *AgoraParameter
	// cLocalUser  unsafe.Pointer
	// subAudioConfig     *SubscribeAudioConfig
	handler               *RtcConnectionObserver
	cHandler              *C.struct__rtc_conn_observer
	localUserObserver     *LocalUserObserver
	cLocalUserObserver    *C.struct__local_user_observer
	audioObserver         *AudioFrameObserver
	cAudioObserver        *C.struct__audio_frame_observer
	videoObserver         *VideoFrameObserver
	cVideoObserver        unsafe.Pointer
	encodedVideoObserver  *VideoEncodedFrameObserver
	cEncodedVideoObserver unsafe.Pointer

	// remoteVideoRWMutex          *sync.RWMutex
	// remoteEncodedVideoReceivers map[*VideoEncodedImageReceiver]*videoEncodedImageReceiverInner
	// vad related for the connection
	enableVad       int
	audioVadManager *AudioVadManager

	// capabilities observer
	cCapObserverHandle    unsafe.Pointer
	cCapabilitiesObserver *C.struct__capabilites_observer

	// audio scenario
	audioScenario AudioScenario
	audioProfile  AudioProfile

	// publish option
	publishConfig *RtcConnectionPublishConfig

	// track & sender
	audioTrack         *LocalAudioTrack
	videoTrack         *LocalVideoTrack
	audioSender        *AudioPcmDataSender
	videoSender        *VideoFrameSender
	encodedAudioSender *AudioEncodedFrameSender
	encodedVideoSender *VideoEncodedImageSender

	// pcm consumption stats for raw pcm data only
	pcmConsumeStats *PcmConsumeStats

	// stream id for data stream： no need to call createDataStream manually, it is created by the sdk automatically
	// and just use it for sendStreamMessage
	dataStreamId int

	// for ai scenario send external audio parameters,default to nil
	sendExternalAudioParameters *SendExternalAudioParameters
}

// for pcm consumption stats
type PcmConsumeStats struct {
	startTime   int64 // in ms
	totalLength int64 // in bytes
	duration    int   // in ms
}

func NewRtcConPublishConfig() *RtcConnectionPublishConfig {
	return &RtcConnectionPublishConfig{
		AudioProfile:     AudioProfileDefault,
		AudioScenario:    AudioScenarioAiServer,
		IsPublishAudio:   true,
		IsPublishVideo:   false,
		AudioPublishType: AudioPublishTypePcm,
		VideoPublishType: VideoPublishTypeNoPublish,
		VideoEncodedImageSenderOptions: &VideoEncodedImageSenderOptions{
			CcMode:        VideoSendCcEnabled, // should check todo???
			CodecType:     VideoCodecTypeH264,
			TargetBitrate: 5000,
		},
		SendExternalAudioParameters: &SendExternalAudioParameters{
			Enabled:                   false,
			SendMs:                    0,
			SendSpeed:                 0,
			DeliverMuteDataForFakeAdm: false,
		},
	}
}

// NOTE: date：2025-06-27
// add audioScenario param, to set the audio scenario for a connection
// and it can diff from the service config, and can diff from each other
func NewRtcConnection(cfg *RtcConnectionConfig, publishConfig *RtcConnectionPublishConfig) *RtcConnection {
	cCfg := CRtcConnectionConfig(cfg)
	defer FreeCRtcConnectionConfig(cCfg)

	audioScenario := publishConfig.AudioScenario
	audioProfile := publishConfig.AudioProfile

	ret := &RtcConnection{
		cConnection: C.agora_rtc_conn_create(agoraService.service, cCfg),
		// subAudioConfig: cfg.SubAudioConfig,
		handler:              nil,
		localUserObserver:    nil,
		audioObserver:        nil,
		videoObserver:        nil,
		encodedVideoObserver: nil,
		// remoteVideoRWMutex:          &sync.RWMutex{},
		// remoteEncodedVideoReceivers: make(map[*VideoEncodedImageReceiver]*videoEncodedImageReceiverInner),
		enableVad:                   0,
		audioVadManager:             nil,
		audioScenario:               audioScenario,
		audioProfile:                audioProfile,
		publishConfig:               publishConfig,
		dataStreamId:                -1,
		sendExternalAudioParameters: nil,
	}

	if isSupportExternalAudio(publishConfig) {
		ret.sendExternalAudioParameters = &SendExternalAudioParameters{
			Enabled:                   publishConfig.SendExternalAudioParameters.Enabled,
			SendMs:                    publishConfig.SendExternalAudioParameters.SendMs,
			SendSpeed:                 publishConfig.SendExternalAudioParameters.SendSpeed,
			DeliverMuteDataForFakeAdm: publishConfig.SendExternalAudioParameters.DeliverMuteDataForFakeAdm,
		}
	}

	ret.localUser = &LocalUser{
		cLocalUser:  C.agora_rtc_conn_get_local_user(ret.cConnection),
		publishFlag: false,
	}
	ret.parameter = &AgoraParameter{
		cParameter: C.agora_rtc_conn_get_agora_parameter(ret.cConnection),
	}

	ret.pcmConsumeStats = &PcmConsumeStats{
		startTime:   0,
		totalLength: 0,
		duration:    0,
	}

	// re set audio scenario now
	ret.localUser.SetAudioEncoderConfiguration(&AudioEncoderConfiguration{AudioProfile: int(audioProfile)})

	ret.localUser.SetAudioScenario(audioScenario)
	fmt.Printf("______set audio scenario to %d, audio profile to %d\n", audioScenario, audioProfile)

	// for stero encoding mode: from 2.3.0, developer can set the codectype & bitrate through setparameter api and do
	// in app layer!
	/* if agoraService.isSteroEncodeMode {
		ret.enableSteroEncodeMode()
	 }*/

	// save to sync map
	agoraService.setConFromHandle(ret.cConnection, ret, ConTypeCCon)
	agoraService.setConFromHandle(ret.localUser.cLocalUser, ret, ConTypeCLocalUser)

	// for any senario, we can handle the capabilities changed event
	// register capabilities observer for any senario: that means we need to handle the capabilities changed event
	// and the senario for each connection can be updated not only fallback from ai server to default but also
	// can upgrade from default to ai server
	// 2025-06-13, weihongqin@agora
	// register capabilities observer only for audio scenario is AudioScenarioAiServer
	// and it is an inner observer, so developer can not unregister it
	if ret.audioScenario == AudioScenarioAiServer {
		ret.cCapabilitiesObserver = CCapatilitiesObserver()
		ret.cCapObserverHandle = C.agora_local_user_capabilities_observer_create(ret.cCapabilitiesObserver)
		C.agora_local_user_register_capabilities_observer(ret.localUser.cLocalUser, ret.cCapObserverHandle)
	}

	//fmt.Printf("______register capabilities observer: clocaluser %v, cconHandle %v, capHandle %v\n", ret.localUser.cLocalUser, ret.cConnection, ret.cCapObserverHandle)

	// create  track
	if publishConfig.IsPublishAudio {
		// check publish type is not no publish
		if publishConfig.AudioPublishType == AudioPublishTypeNoPublish {
			fmt.Printf("WARN:publish audio is no publish, so no audio track created\n")
		}
		if publishConfig.AudioPublishType == AudioPublishTypePcm {
			ret.audioSender = agoraService.mediaFactory.NewAudioPcmDataSender()
			var isSendExternalAudioForAI bool = false
			if ret.sendExternalAudioParameters != nil && ret.sendExternalAudioParameters.Enabled == true {
				isSendExternalAudioForAI = true
			}
			ret.audioTrack = NewCustomAudioTrackPcm(ret.audioSender, ret.audioScenario, isSendExternalAudioForAI)
		} else if publishConfig.AudioPublishType == AudioPublishTypeEncodedPcm {
			ret.encodedAudioSender = agoraService.mediaFactory.NewAudioEncodedFrameSender()
			ret.audioTrack = NewCustomAudioTrackEncoded(ret.encodedAudioSender, AudioTrackMixDisabled)
		}
		// check if the track is not nil
		if ret.audioTrack != nil {
			ret.audioTrack.SetEnabled(true)
		}
	}
	if publishConfig.IsPublishVideo {
		// check publish type is not no publish
		if publishConfig.VideoPublishType == VideoPublishTypeNoPublish {
			fmt.Printf("WARN:publish video is no publish, so no video track created\n")
		}
		if publishConfig.VideoPublishType == VideoPublishTypeYuv {
			ret.videoSender = agoraService.mediaFactory.NewVideoFrameSender()
			ret.videoTrack = NewCustomVideoTrackFrame(ret.videoSender)
		} else if publishConfig.VideoPublishType == VideoPublishTypeEncodedImage {
			ret.encodedVideoSender = agoraService.mediaFactory.NewVideoEncodedImageSender()
			ret.videoTrack = NewCustomVideoTrackEncoded(ret.encodedVideoSender, ret.publishConfig.VideoEncodedImageSenderOptions)
		}
		if ret.videoTrack != nil {
			ret.videoTrack.SetEnabled(true)
		}
	}

	// auto create data stream
	ret.dataStreamId, _ = ret.createDataStream(false, false)

	ret.setExtraSendFrameSpeed(publishConfig.SendExternalAudioParameters)
	fmt.Printf("______auto create data stream, id: %d\n", ret.dataStreamId)

	return ret
}

func (conn *RtcConnection) Release() {
	if conn.cConnection == nil {
		return
	}
	conn.unregisterObserver()
	// delete from sync map
	agoraService.deleteConFromHandle(conn.cConnection, ConTypeCCon)
	agoraService.deleteConFromHandle(conn.localUser.cLocalUser, ConTypeCLocalUser)
	if conn.cVideoObserver != nil {
		agoraService.deleteConFromHandle(conn.cVideoObserver, ConTypeCVideoObserver)
	}
	if conn.cEncodedVideoObserver != nil {
		agoraService.deleteConFromHandle(conn.cEncodedVideoObserver, ConTypeCEncodedVideoObserver)
	}

	// get all receiverInners
	// encodedVideoReceiversInners := make([]*videoEncodedImageReceiverInner, 0, 10)
	// conn.remoteVideoRWMutex.RLock()
	// for _, receiverInner := range conn.remoteEncodedVideoReceivers {
	// 	encodedVideoReceiversInners = append(encodedVideoReceiversInners, receiverInner)
	// }
	// conn.remoteVideoRWMutex.RUnlock()
	// // remove all receiverInners from service
	// agoraService.remoteVideoRWMutex.Lock()
	// for _, receiverInner := range encodedVideoReceiversInners {
	// 	delete(agoraService.remoteEncodedVideoReceivers, receiverInner.cReceiver)
	// }
	// agoraService.remoteVideoRWMutex.Unlock()

	localUser := conn.localUser

	if conn.cCapObserverHandle != nil {
		C.agora_local_user_unregister_capabilities_observer(localUser.cLocalUser, conn.cCapObserverHandle)
	}
	if conn.cAudioObserver != nil {
		C.agora_local_user_unregister_audio_frame_observer(localUser.cLocalUser)
	}
	if conn.cVideoObserver != nil {
		C.agora_local_user_unregister_video_frame_observer(localUser.cLocalUser, conn.cVideoObserver)
	}
	if conn.cEncodedVideoObserver != nil {
		C.agora_local_user_unregister_video_encoded_frame_observer(localUser.cLocalUser, conn.cEncodedVideoObserver)
	}
	if conn.cLocalUserObserver != nil {
		C.agora_local_user_unregister_observer(localUser.cLocalUser)
	}
	if conn.cHandler != nil {
		C.agora_rtc_conn_unregister_observer(conn.cConnection)
	}
	if agoraService.idleMode {
		addIdleItem(conn.cConnection, 2000) // set to delay 2s
	} else {
		C.agora_rtc_conn_destroy(conn.cConnection)
	}

	// clear all receiverInners
	// conn.remoteVideoRWMutex.Lock()
	// conn.remoteEncodedVideoReceivers = make(map[*VideoEncodedImageReceiver]*videoEncodedImageReceiverInner)
	// conn.remoteVideoRWMutex.Unlock()

	// for _, receiverInner := range encodedVideoReceiversInners {
	// 	receiverInner.release()
	// }
	// encodedVideoReceiversInners = nil

	conn.cConnection = nil
	if conn.cAudioObserver != nil {
		FreeCAudioFrameObserver(conn.cAudioObserver)
		conn.cAudioObserver = nil
	}
	if conn.cVideoObserver != nil {
		FreeCVideoFrameObserver(conn.cVideoObserver)
		conn.cVideoObserver = nil
	}
	if conn.cEncodedVideoObserver != nil {
		FreeCEncodedVideoFrameObserver(conn.cEncodedVideoObserver)
		conn.cEncodedVideoObserver = nil
	}
	if conn.cLocalUserObserver != nil {
		FreeCLocalUserObserver(conn.cLocalUserObserver)
		conn.cLocalUserObserver = nil
	}
	if conn.cHandler != nil {
		FreeCRtcConnectionObserver(conn.cHandler)
		conn.cHandler = nil
	}
	if conn.cCapabilitiesObserver != nil {
		FreeCCapatilitiesObserver(conn.cCapabilitiesObserver)
		conn.cCapabilitiesObserver = nil
	}
	if conn.cCapObserverHandle != nil {
		//C.agora_local_user_unregister_capabilities_observer(conn.localUser.cLocalUser, conn.cCapObserverHandle)
		C.agora_local_user_capabilities_observer_destory(conn.cCapObserverHandle)
		conn.cCapObserverHandle = nil
	}
	conn.parameter = nil
	conn.localUser = nil
	localUser.cLocalUser = nil
	localUser = nil
	conn.handler = nil
	conn.localUserObserver = nil
	conn.audioObserver = nil
	conn.videoObserver = nil

	// finally to release all tracks
	if conn.audioTrack != nil {
		conn.audioTrack.Release()
		conn.audioTrack = nil
	}
	if conn.videoTrack != nil {
		conn.videoTrack.Release()
		conn.videoTrack = nil
	}

	// finally to release all senders
	if conn.audioSender != nil {
		conn.audioSender.Release()
		conn.audioSender = nil
	}
	if conn.videoSender != nil {
		conn.videoSender.Release()
		conn.videoSender = nil
	}
	if conn.encodedAudioSender != nil {
		conn.encodedAudioSender.Release()
		conn.encodedAudioSender = nil
	}
	if conn.encodedVideoSender != nil {
		conn.encodedVideoSender.Release()
		conn.encodedVideoSender = nil
	}

	// and set all to nil
	conn.audioTrack = nil
	conn.videoTrack = nil
	conn.audioSender = nil
	conn.videoSender = nil
	conn.encodedAudioSender = nil
	conn.encodedVideoSender = nil
}

func (conn *RtcConnection) GetLocalUser() *LocalUser {
	return conn.localUser
}

func (conn *RtcConnection) GetAgoraParameter() *AgoraParameter {
	return conn.parameter
}

func (conn *RtcConnection) GetConnectionInfo() *RtcConnectionInfo {
	return &conn.connInfo
}

func (conn *RtcConnection) Connect(token string, channel string, uid string) int {
	if conn.cConnection == nil {
		return -1
	}
	conn.connInfo.ChannelId = channel
	conn.connInfo.LocalUserId = uid
	uidInt, _ := strconv.Atoi(uid)
	conn.connInfo.InternalUid = uint(uidInt)
	cChannel := C.CString(channel)
	cToken := C.CString(token)
	cUid := C.CString(uid)
	defer C.free(unsafe.Pointer(cChannel))
	defer C.free(unsafe.Pointer(cToken))
	defer C.free(unsafe.Pointer(cUid))
	return int(C.agora_rtc_conn_connect(conn.cConnection, cToken, cChannel, cUid))
}

// date: 2025-07-04
// add a function to disconnect the connection
// and it will unpublish all tracks and unregister all observers
// and then do really disconnect, no need to call unregister observer manually
func (conn *RtcConnection) Disconnect() int {
	if conn.cConnection == nil {
		return -1
	}
	// date: 2025-07-04
	//1. unpublish all tracks
	conn.UnpublishAudio()
	conn.UnpublishVideo()

	//2. unregister all observers, but except rtc connection observer
	conn.unregisterAudioFrameObserver()
	conn.unregisterVideoFrameObserver()
	conn.unregisterVideoEncodedFrameObserver()
	conn.unregisterAudioEncodedFrameObserver()

	conn.unregisterLocalUserObserver()

	//3 and then do really disconnect
	ret := int(C.agora_rtc_conn_disconnect(conn.cConnection))

	//3. unregister rtc connection observer，yeah

	return ret
}

func (conn *RtcConnection) RenewToken(token string) int {
	if conn.cConnection == nil {
		return -1
	}
	cToken := C.CString(token)
	defer C.free(unsafe.Pointer(cToken))
	return int(C.agora_rtc_conn_renew_token(conn.cConnection, cToken))
}

func (conn *RtcConnection) createDataStream(reliable bool, ordered bool) (int, int) {
	if conn.cConnection == nil {
		return -1, -1
	}
	// int* stream_id, int reliable, int ordered
	cStreamId := C.int(-1)
	ret := int(C.agora_rtc_conn_create_data_stream(conn.cConnection, &cStreamId, CIntFromBool(reliable), CIntFromBool(ordered)))
	return int(cStreamId), ret
}

func (conn *RtcConnection) SendStreamMessage(msg []byte) int {
	if conn == nil || conn.cConnection == nil {
		return -1
	}
	cMsg := C.CBytes(msg)
	defer C.free(cMsg)
	streamId := conn.dataStreamId
	return int(C.agora_rtc_conn_send_stream_message(conn.cConnection, C.int(streamId), (*C.char)(cMsg), C.uint32_t(len(msg))))
}

func (conn *RtcConnection) RegisterObserver(handler *RtcConnectionObserver) int {
	if conn.cConnection == nil || handler == nil {
		return -1
	}
	// avoid re-register observer
	if conn.handler == handler {
		return 0
	}
	// unregister old observer
	conn.unregisterObserver()

	// register new observer
	conn.handler = handler
	if conn.cHandler == nil {
		conn.cHandler = CRtcConnectionObserver()
		C.agora_rtc_conn_register_observer(conn.cConnection, conn.cHandler)
	}
	return 0

}

func (conn *RtcConnection) unregisterObserver() int {
	// check if need to unregister
	if conn.cConnection == nil {
		return 0
	}
	if conn.cHandler != nil {
		C.agora_rtc_conn_unregister_observer(conn.cConnection)
		FreeCRtcConnectionObserver(conn.cHandler)
	}
	conn.cHandler = nil
	conn.handler = nil
	return 0
}

func (conn *RtcConnection) RegisterLocalUserObserver(handler *LocalUserObserver) int {
	if conn.cConnection == nil || handler == nil {
		return -1
	}
	// avoid re-register observer
	if conn.localUserObserver == handler {
		return 0
	}
	// unregister old observer
	conn.unregisterLocalUserObserver()

	// register new observer
	conn.localUserObserver = handler
	var ret int = -100
	if conn.cLocalUserObserver == nil {
		conn.cLocalUserObserver = CLocalUserObserver()
		ret = int(C.agora_local_user_register_observer(conn.localUser.cLocalUser, conn.cLocalUserObserver))
	}
	return int(ret)
}

func (conn *RtcConnection) unregisterLocalUserObserver() int {
	// check if need to unregister
	if conn.cConnection == nil {
		return 0
	}
	if conn.cLocalUserObserver != nil {
		C.agora_local_user_unregister_observer(conn.localUser.cLocalUser)
		FreeCLocalUserObserver(conn.cLocalUserObserver)
	}
	conn.cLocalUserObserver = nil
	conn.localUserObserver = nil
	return 0
}

func (conn *RtcConnection) RegisterAudioFrameObserver(observer *AudioFrameObserver, enableVad int, vadConfigure *AudioVadConfigV2) int {
	if conn.cConnection == nil || observer == nil {
		return -1
	}
	// avoid re-register observer
	if conn.audioObserver == observer {
		return 0
	}
	// unregister old observer
	conn.unregisterAudioFrameObserver()

	// re-assign vad related
	conn.enableVad = enableVad
	if conn.enableVad > 0 && vadConfigure != nil {
		conn.audioVadManager = NewAudioVadManager(vadConfigure)
	}

	conn.audioObserver = observer
	if conn.cAudioObserver == nil {
		conn.cAudioObserver = CAudioFrameObserver()
		C.agora_local_user_register_audio_frame_observer(conn.localUser.cLocalUser, conn.cAudioObserver)
	}
	return 0
}

func (conn *RtcConnection) unregisterAudioFrameObserver() int {
	// check if need to unregister
	if conn.cConnection == nil {
		return 0
	}
	conn.enableVad = 0
	if conn.cAudioObserver != nil {
		C.agora_local_user_unregister_audio_frame_observer(conn.localUser.cLocalUser)
		FreeCAudioFrameObserver(conn.cAudioObserver)
	}
	conn.cAudioObserver = nil
	conn.audioObserver = nil
	if conn.audioVadManager != nil {
		conn.audioVadManager.Release()
		conn.audioVadManager = nil
	}

	conn.enableVad = 0
	return 0
}

func (conn *RtcConnection) RegisterVideoFrameObserver(observer *VideoFrameObserver) int {
	if conn.cConnection == nil || observer == nil {
		return -1
	}
	// avoid re-register observer
	if conn.videoObserver == observer {
		return 0
	}
	// unregister old observer
	conn.unregisterVideoFrameObserver()

	conn.videoObserver = observer
	if conn.cVideoObserver == nil {
		conn.cVideoObserver = CVideoFrameObserver()
		// store to sync map
		agoraService.setConFromHandle(conn.cVideoObserver, conn, ConTypeCVideoObserver)
		C.agora_local_user_register_video_frame_observer(conn.localUser.cLocalUser, conn.cVideoObserver)
	}
	return 0
}

func (conn *RtcConnection) unregisterVideoFrameObserver() int {
	// check if need to unregister
	if conn.cConnection == nil {
		return 0
	}
	if conn.cVideoObserver != nil {
		C.agora_local_user_unregister_video_frame_observer(conn.localUser.cLocalUser, conn.cVideoObserver)
		// delete from sync map
		agoraService.deleteConFromHandle(conn.cVideoObserver, ConTypeCVideoObserver)
		FreeCVideoFrameObserver(conn.cVideoObserver)
		conn.cVideoObserver = nil
	}
	conn.cVideoObserver = nil
	conn.videoObserver = nil
	return 0
}

func (conn *RtcConnection) RegisterVideoEncodedFrameObserver(observer *VideoEncodedFrameObserver) int {
	if conn.cConnection == nil || observer == nil {
		return -1
	}
	// avoid re-register observer
	if conn.encodedVideoObserver == observer {
		return 0
	}
	// unregister old observer
	conn.unregisterVideoEncodedFrameObserver()

	conn.encodedVideoObserver = observer
	if conn.cEncodedVideoObserver == nil {
		conn.cEncodedVideoObserver = CVideoEncodedFrameObserver()
		// store to sync map
		agoraService.setConFromHandle(conn.cEncodedVideoObserver, conn, ConTypeCEncodedVideoObserver)
		C.agora_local_user_register_video_encoded_frame_observer(conn.localUser.cLocalUser, conn.cEncodedVideoObserver)
	}
	return 0
}

func (conn *RtcConnection) unregisterVideoEncodedFrameObserver() int {
	// check if need to unregister
	if conn.cConnection == nil {
		return 0
	}
	if conn.cEncodedVideoObserver != nil {
		C.agora_local_user_unregister_video_encoded_frame_observer(conn.localUser.cLocalUser, conn.cEncodedVideoObserver)
		// delete from sync map
		agoraService.deleteConFromHandle(conn.cEncodedVideoObserver, ConTypeCEncodedVideoObserver)
		FreeCEncodedVideoFrameObserver(conn.cEncodedVideoObserver)
		conn.cEncodedVideoObserver = nil
	}
	conn.cEncodedVideoObserver = nil
	conn.encodedVideoObserver = nil
	return 0
}

/*
* for stero encoded audio mode
* Must be called before con.connect
 */
func (conn *RtcConnection) enableSteroEncodeMode() int {
	if conn.cConnection == nil {
		return -1
	}
	//set private parameter
	localUser := conn.localUser

	// remove set senario to gs here, as it can pass senario as a parameter in NewRtcConnection
	// only force profile to stero encoding profile and force codec to opus here!
	// the default codec for musicstandstero is HAAC, no need for audio only senario
	// the HAAC is valid or better for muisic+audio combination senario
	// localUser.SetAudioScenario(AudioScenarioGameStreaming)
	localUser.SetAudioEncoderConfiguration(&AudioEncoderConfiguration{AudioProfile: int(AudioProfileMusicStandardStereo)})

	// fill pirvate parameter
	agoraParameterHandler := conn.parameter
	agoraParameterHandler.SetParameters("{\"che.audio.aec.enable\":false}")
	agoraParameterHandler.SetParameters("{\"che.audio.ans.enable\":false}")
	agoraParameterHandler.SetParameters("{\"che.audio.agc.enable\":false}")
	//  // "HEAAC_2ch" is case 78,but no need to set it for ai senario so disable it
	// in ai senario, we only want stero encoded audio, but really for speech so disable 78
	agoraParameterHandler.SetParameters("{\"che.audio.custom_payload_type\":122}")
	// and set bitrate to 32000,but not work in here! so move to other place
	//agoraParameterHandler.SetParameters("{\"che.audio.custom_bitrate\":32000}")
	return 0

}

type EncryptionConfig struct {
	EncryptionMode    int
	EncryptionKey     string
	EncryptionKdfSalt []byte
}

// EnableEncryption enables or disables encryption for the RTC connection.
// It sets the encryption mode and configuration for the connection.
// Must be called before RtcConnection.Connect
//
// Parameters:
// - enable: An integer indicating whether to enable (1) or disable (0) encryption.
// - config: A pointer to an EncryptionConfig struct containing the encryption configuration.
//
// Returns:
// - An integer indicating the result of the operation. 0 indicates success, and negative value indicates failure.
func (conn *RtcConnection) EnableEncryption(enable int, config *EncryptionConfig) int {
	if conn.cConnection == nil || config == nil || enable == 0 {
		return -1
	}

	cConfig := C.struct__encryption_config{}
	C.memset(unsafe.Pointer(&cConfig), 0, C.sizeof_struct__encryption_config)
	cConfig.encryption_mode = C.int(config.EncryptionMode)
	if config.EncryptionKey != "" {
		ckey := C.CString(config.EncryptionKey)
		defer C.free(unsafe.Pointer(ckey))
		cConfig.encryption_key = ckey
	}
	saltlen := 0
	if config.EncryptionKdfSalt != nil {
		saltlen = len(config.EncryptionKdfSalt)
	}
	if saltlen > 0 {
		if saltlen > 32 {
			saltlen = 32
		}
		C.memcpy(unsafe.Pointer(&cConfig.encryption_kdf_salt[0]), unsafe.Pointer(&config.EncryptionKdfSalt[0]), C.size_t(saltlen))
	}

	ret := C.agora_rtc_conn_enable_encryption(conn.cConnection, C.int(enable), &cConfig)
	return int(ret)
}
func (conn *RtcConnection) handleCapabilitiesChanged(caps *C.struct__capabilities, size C.int) int {
	if conn == nil || conn.cConnection == nil || conn.cCapabilitiesObserver == nil {
		return -1
	}

	var fallback bool = true

	// 获取数组起始地址
	capsPtr := unsafe.Pointer(caps)
	// 计算每个元素的大小
	elemSize := unsafe.Sizeof(C.struct__capabilities{})

	// 遍历数组
	for i := 0; i < int(size); i++ {
		// 计算当前元素的地址
		curCapPtr := (*C.struct__capabilities)(unsafe.Pointer(uintptr(capsPtr) + uintptr(i)*elemSize))

		// 打印基本信息
		fmt.Printf("Capability[%d] - Type: %d\n", i, curCapPtr.capability_type)

		// 获取并打印 item_map 信息
		if curCapPtr.item_map != nil {
			itemMap := (*C.struct__capability_item_map)(unsafe.Pointer(curCapPtr.item_map))

			// 遍历 item_map 中的每个 item
			for j := 0; j < int(itemMap.size); j++ {
				item := (*C.struct__capability_item)(unsafe.Pointer(uintptr(unsafe.Pointer(itemMap.item)) + uintptr(j)*unsafe.Sizeof(C.struct__capability_item{})))

				// 打印 item 的详细信息
				itemName := ""
				if item.name != nil {
					itemName = C.GoString(item.name)
				}
				fmt.Printf("  Item[%d] - ID: %d, Name: %s\n", j, item.id, itemName)

			}
		}

		// check fabll back or not
		if curCapPtr.capability_type == 19 && curCapPtr.item_map != nil {
			itemMap := (*C.struct__capability_item_map)(unsafe.Pointer(curCapPtr.item_map))
			for j := 0; j < int(itemMap.size); j++ {
				item := (*C.struct__capability_item)(unsafe.Pointer(uintptr(unsafe.Pointer(itemMap.item)) + uintptr(j)*unsafe.Sizeof(C.struct__capability_item{})))

				itemName := ""
				if item.name != nil {
					itemName = C.GoString(item.name)
				}
				if itemName == "SUPPORT" && item.id == 0 {
					fallback = false
				}
			}
		}
	}
	fmt.Printf("*********fallback: %v\n", fallback)
	// call conectionobserver to notify
	if fallback && conn.cHandler != nil && conn.handler.OnAIQoSCapabilityMissing != nil {
		userDefinedSenario := conn.handler.OnAIQoSCapabilityMissing(conn, int(AudioScenarioChorus))
		if userDefinedSenario >= 0 {
			fmt.Printf("onAIQoSCapabilityMissing, userDefinedSenario: %d\n", userDefinedSenario)
			conn.UpdateAudioSenario(AudioScenario(userDefinedSenario))
		}
	}
	return 0
}

// author: weihongqin
// description: push audio pcm data to agora sdk
// param: data: audio data
// param: sampleRate: sample rate
// param: channels: channels
// return: 0: success, -1: error, -2: invalid data
//
//date:2025-07-04 10:00:00
func (conn *RtcConnection) PublishAudio() int {
	if conn == nil || conn.cConnection == nil || conn.audioTrack == nil {
		return -2000
	}
	//conn.audioTrack.SetEnabled(true)
	ret := conn.localUser.publishAudio(conn.audioTrack)
	return int(ret)
}
func (conn *RtcConnection) UnpublishAudio() int {
	if conn == nil || conn.cConnection == nil || conn.audioTrack == nil {
		return -2000
	}

	//reset consumer stats
	if conn.pcmConsumeStats != nil {
		conn.pcmConsumeStats.reset()
	}
	//conn.audioTrack.SetEnabled(false)
	ret := conn.localUser.unpublishAudio(conn.audioTrack)
	return int(ret)
}

func (conn *RtcConnection) PublishVideo() int {
	if conn == nil || conn.cConnection == nil || conn.videoTrack == nil {
		return -2000
	}
	//conn.videoTrack.SetEnabled(true)
	ret := conn.localUser.publishVideo(conn.videoTrack)
	return int(ret)
}

func (conn *RtcConnection) UnpublishVideo() int {
	if conn == nil || conn.cConnection == nil || conn.videoTrack == nil {
		return -2000
	}
	//conn.videoTrack.SetEnabled(false)
	ret := conn.localUser.unpublishVideo(conn.videoTrack)
	return int(ret)
}
func (conn *RtcConnection) InterruptAudio() int {
	if conn == nil || conn.cConnection == nil || conn.audioTrack == nil {
		return -2000
	}

	if conn.audioScenario == AudioScenarioAiServer {
		// for aiServer, we need to unpublish the track
		conn.UnpublishAudio()
		// and publish the track again
		conn.PublishAudio()

	} else {
		// and other scenarios, we need to clear the buffer of the track
		conn.audioTrack.ClearSenderBuffer()
	}

	// if has audio consumption, we need to reset the stats
	if conn.pcmConsumeStats != nil {
		conn.pcmConsumeStats.reset()
	}
	return 0
}

/*
date:2025-08-14
author: weihongqin
description: push audio pcm data to agora sdk
param: data: audio data
param: sampleRate: sample rate
param: channels: channels
param: startPtsMs: start present timestamp in ms. and it can pass to onPlaybackBeforeMixing's audioframe.
default value is 0. Note: be carefully to set the pstvalue
for server end, the AudioFrame.RenderTimeMs is the capture timestamp
and the presentTimeMs is really the present timestamp.
*/
func (conn *RtcConnection) PushAudioPcmData(data []byte, sampleRate int, channels int, startPtsInMs int64) int {
	if conn == nil || conn.cConnection == nil || conn.audioSender == nil {
		return -2000
	}
	readLen := len(data)
	bytesPerFrameInMs := (sampleRate / 1000) * 2 * channels // 1ms , channels and 16bit
	// validity check: only accepts data with lengths that are integer multiples of 10ms​​
	if readLen%bytesPerFrameInMs != 0 {
		fmt.Printf("PushAudioPcmData data length is not integer multiples of 10ms, readLen: %d, bytesPerFrame: %d\n", readLen, bytesPerFrameInMs)
		return -2
	}
	packnumInMs := readLen / bytesPerFrameInMs

	frame := &AudioFrame{
		Buffer:            nil,
		RenderTimeMs:      0,
		PresentTimeMs:     startPtsInMs,
		SamplesPerChannel: sampleRate / 100,
		BytesPerSample:    2,
		Channels:          channels,
		SamplesPerSec:     sampleRate,
		Type:              AudioFrameTypePCM16,
	}

	frame.Buffer = data
	frame.SamplesPerChannel = (sampleRate / 1000) * packnumInMs

	//for ai server limited mode only
	// for every new round, should call the function to set the total extra send ms
	// automatically set the total extra send ms when the pcm data is sent
	isNewRound := conn.pcmConsumeStats.isNewRound(sampleRate, channels)
	if isNewRound {
		conn.setTotalExtraSendMs()
	}

	ret := conn.audioSender.SendAudioPcmData(frame)
	if ret == 0 {
		conn.pcmConsumeStats.addPcmData(readLen, sampleRate, channels)
	}
	return ret
}
func (conn *RtcConnection) PushAudioEncodedData(data []byte, frameInfo *EncodedAudioFrameInfo) int {
	if conn == nil || conn.cConnection == nil || conn.encodedAudioSender == nil {
		return -2000
	}
	return conn.encodedAudioSender.SendEncodedAudioFrame(data, frameInfo)
}
func (conn *RtcConnection) PushVideoFrame(frame *ExternalVideoFrame) int {
	if conn == nil || conn.cConnection == nil || conn.videoSender == nil {
		return -2000
	}
	return conn.videoSender.SendVideoFrame(frame)
}
func (conn *RtcConnection) PushVideoEncodedData(data []byte, frameInfo *EncodedVideoFrameInfo) int {
	if conn == nil || conn.cConnection == nil || conn.encodedVideoSender == nil {
		return -2000
	}
	return conn.encodedVideoSender.SendEncodedVideoImage(data, frameInfo)
}

func (conn *RtcConnection) unregisterAudioEncodedFrameObserver() int {
	// todo: ??? not implement

	return -1
}

// to pudate connction's scenario
func (conn *RtcConnection) UpdateAudioSenario(scenario AudioScenario) int {

	//1. validate the connection
	if conn == nil || conn.cConnection == nil {
		return -2000
	}

	//2. if scenario is the same as the internal one, do nothing
	if conn.audioScenario == scenario {
		return 0
	}

	//3. check the track

	// 3. update the connection's senario

	conn.audioScenario = scenario
	conn.localUser.SetAudioScenario(scenario)

	//4.unpublish the track
	conn.UnpublishAudio()

	//5. update the audioScenario
	conn.audioTrack.Release()
	conn.audioTrack.cTrack = nil

	//ToDo：It would be better to leave the fallback to the client, using the best practice approach
	//The reason is that the client's chosen fallback strategy may not be chorus, but other strategies!! This way, we fix the strategy, which will limit the client
	//4. create a new cTrack
	var cTrack unsafe.Pointer = nil
	var isAiServer bool = false
	if scenario == AudioScenarioAiServer {
		isAiServer = true
	}
	// for audio, pcmdatasender and encodedaudiosender, check which one is valid
	var csender unsafe.Pointer = nil
	if conn.audioSender != nil {
		csender = conn.audioSender.cSender
	} else if conn.encodedAudioSender != nil {
		csender = conn.encodedAudioSender.cSender
	}
	if isAiServer {
		cTrack = C.agora_service_create_direct_custom_audio_track_pcm(agoraService.service, csender)
	} else {
		cTrack = C.agora_service_create_custom_audio_track_pcm(agoraService.service, csender)
	}

	//5. assign the new cTrack
	conn.audioTrack.cTrack = cTrack

	//6. set properties to new cTrack
	conn.audioTrack.SetSendDelayMs(10)
	//anyway, to set max buffered audio frame number to 100000
	conn.audioTrack.SetMaxBufferedAudioFrameNumber(100000) //up to 16min,100000 frames
	
	conn.localUser.SetAudioScenario(scenario)

	//7. update the properties of pcmsender
	//update pcmsender's info

	//update connection's info
	conn.audioScenario = scenario

	//8. publish the track
	conn.PublishAudio()
	fmt.Printf("______update audio track \n")
	return 0
}

func (conn *RtcConnection) IsPushToRtcCompleted() bool {
	if conn == nil || conn.pcmConsumeStats == nil {
		return false
	}
	return conn.pcmConsumeStats.isPushCompleted(conn.audioScenario)
}
func (consumer *PcmConsumeStats) addPcmData(len int, samplerate int, channels int) {
	isNewRound := consumer.isNewRound(samplerate, channels)
	if isNewRound {
		consumer.startTime = time.Now().UnixMilli()
		consumer.totalLength = 0
	}

	consumer.totalLength += int64(len)

	// update duration
	consumer.duration = int(consumer.totalLength / (int64(samplerate/1000) * int64(channels) * 2))
	//fmt.Printf("addPcmData, duration: %d, totalLength: %d, startTime: %d\n", consumer.duration, consumer.totalLength, consumer.startTime)
}

func (consumer *PcmConsumeStats) isNewRound(samplerate int, channels int) bool {
	if consumer.startTime == 0 {
		return true
	}
	now := time.Now().UnixMilli()
	diff := now - consumer.startTime

	if diff > int64(consumer.duration) {
		return true
	}
	return false
}
func (consumer *PcmConsumeStats) getCurrentPosition() int {
	if consumer == nil || consumer.startTime == 0 {
		return 0
	}
	now := time.Now().UnixMilli()
	diff := now - consumer.startTime
	return int(diff)
}

const E2E_DELAY_MS = 180

func (consumer *PcmConsumeStats) isPushCompleted(scenario AudioScenario) bool {
	now := time.Now().UnixMilli()
	diff := now - consumer.startTime
	delay := E2E_DELAY_MS
	/* 	if scenario == AudioScenarioAiServer {
		delay = E2E_DELAY_MS
	} */
	if diff > int64(consumer.duration+delay) {
		return true
	}
	return false
}
func (consumer *PcmConsumeStats) reset() {
	consumer.startTime = 0
	consumer.totalLength = 0
	consumer.duration = 0
}
func (conn *RtcConnection) SetVideoEncoderConfiguration(cfg *VideoEncoderConfiguration) int {
	// validate the connection
	if conn == nil || conn.cConnection == nil || conn.videoTrack == nil || conn.videoTrack.cTrack == nil {
		return -1
	}

	// validate the cfg

	cCfg := C.struct__video_encoder_config{}
	C.memset(unsafe.Pointer(&cCfg), 0, C.sizeof_struct__video_encoder_config)
	cCfg.codec_type = C.int(cfg.CodecType)
	cCfg.dimensions.width = C.int(cfg.Width)
	cCfg.dimensions.height = C.int(cfg.Height)
	cCfg.frame_rate = C.int(cfg.Framerate)
	cCfg.bitrate = C.int(cfg.Bitrate * 1000)
	cCfg.min_bitrate = C.int(cfg.MinBitrate * 1000)
	cCfg.orientation_mode = C.int(cfg.OrientationMode)
	cCfg.degradation_preference = C.int(cfg.DegradePreference)
	cCfg.mirror_mode = C.int(cfg.MirrorMode)
	cCfg.encode_alpha = CIntFromBool(cfg.EncodeAlpha)
	return int(C.agora_local_video_track_set_video_encoder_config(conn.videoTrack.cTrack, &cCfg))
}
func (conn *RtcConnection) SendAudioMetaData(metaData []byte) int {
	if conn == nil || conn.cConnection == nil || conn.localUser == nil || conn.localUser.cLocalUser == nil {
		return -2000
	}
	return conn.localUser.sendAudioMetaData(metaData)
}

// apm related: date 20251028
// at present, only support service-level apm filter related config, not support connection-level apm filter related config
// later, we will support connection-level apm filter related config,i.e, add a new connection-level api like:
// conneciton::SetApmConfigure(apmconfigure *ApmConf) which can set the apm filter related config for the connection
// but should call before connection.connect
func (conn *RtcConnection) isEnable3A() bool {
	if conn == nil || conn.cConnection == nil || conn.localUser == nil || conn.localUser.cLocalUser == nil {
		return false
	}
	if agoraService.apmConfig == nil {
		return false
	}
	return true
}
func (conn *RtcConnection) setApmFilterProperties(uid *C.char, cRemoteAudioTrack unsafe.Pointer) int {
	if conn == nil || conn.cConnection == nil || conn.localUser == nil || conn.localUser.cLocalUser == nil {
		return -2000
	}
	apmEnabled := conn.isEnable3A()
	fmt.Printf("------setApmFilterProperties, uid: %s, apmEnabled: %t\n", C.GoString(uid), apmEnabled)
	if !apmEnabled {
		return 0
	}

	// Load AINS resource
	ret := setFilterPropertyByTrack(cRemoteAudioTrack, "audio_processing_remote_playback", "apm_load_resource", "ains", false)
	if ret != 0 {
		fmt.Printf("Failed to set apm_load_resource, error: %d\n", ret)
	}

	// Set APM configuration
	configJSON := agoraService.apmConfig.toJson()
	ret = setFilterPropertyByTrack(cRemoteAudioTrack, "audio_processing_remote_playback", "apm_config", configJSON, false)
	if ret != 0 {
		fmt.Printf("Failed to set apm_config, error: %d\n", ret)
		return -1
	}
	fmt.Println("*****APM config set:%s", configJSON)

	// Enable dump
	if agoraService.apmConfig.EnableDump {
		ret = setFilterPropertyByTrack(cRemoteAudioTrack, "audio_processing_remote_playback", "apm_dump", "true", false)
		if ret != 0 {
			fmt.Printf("Failed to enable apm_dump, error: %d\n", ret)
		}
	}

	return 0
}

// should call before connection is connected
func (conn *RtcConnection) setExtraSendFrameSpeed(sendExternalAudioParameters *SendExternalAudioParameters) int {
	if conn == nil || conn.cConnection == nil || conn.localUser == nil || conn.localUser.cLocalUser == nil {
		return -2000
	}
	if sendExternalAudioParameters == nil || sendExternalAudioParameters.Enabled == false || sendExternalAudioParameters.SendMs <= 0 || sendExternalAudioParameters.SendSpeed <= 1 {
		return -2001
	}

	// check the parameters
	speed := sendExternalAudioParameters.SendSpeed

	if speed > 5 {
		speed = 5
	}
	if speed < 1 {
		speed = 1
	}

	// set the parameters to the connection
	params := fmt.Sprintf("{\"che.audio.extra_send_frames_per_interval_for_fake_adm\": %d", speed)
	conn.GetAgoraParameter().SetParameters(params)

	// deliver mute data for fake adm: can only effect at service level, no effect for connection level
	setDeliverMuteData(sendExternalAudioParameters.DeliverMuteDataForFakeAdm)
	return 0
}
func (conn *RtcConnection) setTotalExtraSendMs() int {
	// note: for internal call only, so do not need to check the connection
	// only check the model from parameters
	if conn.sendExternalAudioParameters == nil || conn.sendExternalAudioParameters.Enabled == false {
		return -2001
	}

	ret := C.agora_local_audio_track_set_total_extra_send_ms(conn.audioTrack.cTrack, C.uint64_t(conn.sendExternalAudioParameters.SendMs))

	return int(ret)
}

func isSupportExternalAudio(publishConfig *RtcConnectionPublishConfig) bool {
	if publishConfig.SendExternalAudioParameters != nil && publishConfig.SendExternalAudioParameters.Enabled == true && publishConfig.SendExternalAudioParameters.SendMs > 0 && publishConfig.SendExternalAudioParameters.SendSpeed > 1 {
		return true
	}
	return false
}

// add a closer function to set only once the extra audio send ms the agora parameter
var deliverMuteDataHasSet bool = false
func setDeliverMuteData(deliverMuteData bool) bool {
    
	fmt.Printf("createDeliverMuteDataSetter, deliverMuteData: %t\n", deliverMuteDataHasSet)
    
   
	if !deliverMuteData && !deliverMuteDataHasSet {
		agoraParameterHandler := GetAgoraParameter()
		agoraParameterHandler.SetParameters("{\"che.audio.deliver_mute_data_for_fake_adm\":false}")
		deliverMuteDataHasSet = true // set the flag to ensure only once
		fmt.Printf("createDeliverMuteDataSetter, deliverMuteData: %t, deliverMuteDataHasSet: %t\n", deliverMuteData, deliverMuteDataHasSet)
	}

	return deliverMuteDataHasSet  
}
// send intra request to the remote user to request a new key frame
// limitation: 
// 1. if remtoe user is fixed-interval key frame encoding, the api call will be ignored
// 2. the min-frequency of the api call is 1s, otherwise, it will be ignored  or the vos will do converge
func (conn *RtcConnection) SendIntraRequest(remoteUid string) int {
	if conn == nil || conn.cConnection == nil || conn.localUser == nil || conn.localUser.cLocalUser == nil {
		return -2000
	}
	cUid := C.CString(remoteUid)
	defer C.free(unsafe.Pointer(cUid))
	ret := int(C.agora_local_user_send_intra_request(conn.localUser.cLocalUser, cUid))
	return ret
}
