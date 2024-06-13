package agoraservice

type AreaCode int

const (
	// AreaCodeCN is Mainland China.
	AreaCodeCN AreaCode = 0x00000001
	// AreaCodeNA is North America.
	AreaCodeNA AreaCode = 0x00000002
	// AreaCodeEU is Europe.
	AreaCodeEU AreaCode = 0x00000004
	// AreaCodeAS is Asia, excluding Mainland China.
	AreaCodeAS AreaCode = 0x00000008
	// AreaCodeJP is Japan.
	AreaCodeJP AreaCode = 0x00000010
	// AreaCodeIN is India.
	AreaCodeIN AreaCode = 0x00000020
	// AreaCodeGlob is (Default) Global.
	AreaCodeGlob AreaCode = (0xFFFFFFFF)
)

type ClientRole int

const (
	// Broadcaster. A broadcaster can both send and receive streams.
	ClientRoleBroadcaster ClientRole = 1
	// Audience. An audience can only receive streams.
	ClientRoleAudience ClientRole = 2
)

type ChannelProfile int

const (
	// Communication.
	// This profile prioritizes smoothness and applies to the one-to-one scenario.
	ChannelProfileCommunication ChannelProfile = 0
	// (Default) Live Broadcast.
	// This profile prioritizes supporting a large audience in a live broadcast channel.
	ChannelProfileLiveBroadcasting ChannelProfile = 1
)

type UserMediaInfo int

const (
	/**
	* 0: The user has muted the audio.
	 */
	UserMediaInfoMuteAudio UserMediaInfo = 0
	/**
	* 1: The user has muted the video.
	 */
	UserMediaInfoMuteVideo UserMediaInfo = 1
	/**
	* 4: The user has enabled the video, which includes video capturing and encoding.
	 */
	UserMediaInfoEnableVideo UserMediaInfo = 4
	/**
	* 8: The user has enabled the local video capturing.
	 */
	UserMediaInfoEnableLocalVideo UserMediaInfo = 8
)
