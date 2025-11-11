package agorartm

/*
//引入Agora C封装
#cgo linux CFLAGS: -I${SRCDIR}/../../agora_sdk/agora_rtm_sdk_c/include -I${SRCDIR}/../../agora_sdk/agora_rtm_sdk_c/agora_rtm_sdk/high_level_api/include
#cgo darwin CFLAGS: -I${SRCDIR}/../../agora_sdk/agora_rtm_sdk_c/include -I${SRCDIR}/../../agora_sdk/agora_rtm_sdk_c/agora_rtm_sdk/high_level_api/include

//链接AgoraRTM SDK
#cgo linux LDFLAGS: -L${SRCDIR}/../../agora_sdk -lagora_rtm_sdk -laosl -lagora_rtm_sdk_c
#cgo darwin LDFLAGS: -L${SRCDIR}/../../agora_sdk_mac -lAgoraRtmKit -laosl -lagora_rtm_sdk_c
#include <stdlib.h>
#include <string.h>
#include "C_AgoraRtmBase.h"
*/
import "C"
import "unsafe"

// #region agora

// #region agora::rtm

const DEFAULT_LOG_SIZE_IN_KB = 1024

/**
 * IP areas.
 */
type RtmAreaCode C.enum_C_RTM_AREA_CODE

const (
	/**
	 * Mainland China.
	 */
	RtmAreaCodeCN RtmAreaCode = C.RTM_AREA_CODE_CN
	/**
	 * North America.
	 */
	RtmAreaCodeNA RtmAreaCode = C.RTM_AREA_CODE_NA
	/**
	 * Europe.
	 */
	RtmAreaCodeEU RtmAreaCode = C.RTM_AREA_CODE_EU
	/**
	 * Asia, excluding Mainland China.
	 */
	RtmAreaCodeAS RtmAreaCode = C.RTM_AREA_CODE_AS
	/**
	 * Japan.
	 */
	RtmAreaCodeJP RtmAreaCode = C.RTM_AREA_CODE_JP
	/**
	 * India.
	 */
	RtmAreaCodeIN RtmAreaCode = C.RTM_AREA_CODE_IN
	/**
	 * (Default) Global.
	 */
	RtmAreaCodeGLOB RtmAreaCode = C.RTM_AREA_CODE_GLOB
)

/**
 * The log level for rtm sdk.
 */
type RtmLogLevel C.enum_C_RTM_LOG_LEVEL

const (
	/**
	 * 0x0000: No logging.
	 */
	RtmLogLevelNONE RtmLogLevel = C.RTM_LOG_LEVEL_NONE
	/**
	 * 0x0001: Informational messages.
	 */
	RtmLogLevelINFO RtmLogLevel = C.RTM_LOG_LEVEL_INFO
	/**
	 * 0x0002: Warnings.
	 */
	RtmLogLevelWARN RtmLogLevel = C.RTM_LOG_LEVEL_WARN
	/**
	 * 0x0004: Errors.
	 */
	RtmLogLevelERROR RtmLogLevel = C.RTM_LOG_LEVEL_ERROR
	/**
	 * 0x0008: Critical errors that may lead to program termination.
	 */
	RtmLogLevelFATAL RtmLogLevel = C.RTM_LOG_LEVEL_FATAL
)

/**
 * The encryption mode.
 */
type RtmEncryptionMode C.enum_C_RTM_ENCRYPTION_MODE

const (
	/**
	 * Disable message encryption.
	 */
	RtmEncryptionModeNONE RtmEncryptionMode = C.RTM_ENCRYPTION_MODE_NONE
	/**
	 * 128-bit AES encryption, GCM mode.
	 */
	RtmEncryptionModeAES128GCM RtmEncryptionMode = C.RTM_ENCRYPTION_MODE_AES_128_GCM
	/**
	 * 256-bit AES encryption, GCM mode.
	 */
	RtmEncryptionModeAES256GCM RtmEncryptionMode = C.RTM_ENCRYPTION_MODE_AES_256_GCM
)

/**
 * The error codes of rtm client.
 * ref to rtm sdk
 */

/**
* Connection states between rtm sdk and agora server.
ref to rtm sdk
*/

/**
* Reasons for connection state change.
ref to rtm sdk
*/

/**
 * RTM channel type.
 */
type RtmChannelType C.enum_C_RTM_CHANNEL_TYPE

const (
	/**
	 * 0: Unknown channel type.
	 */
	RtmChannelTypeNONE RtmChannelType = C.RTM_CHANNEL_TYPE_NONE
	/**
	 * 1: Message channel.
	 */
	RtmChannelTypeMESSAGE RtmChannelType = C.RTM_CHANNEL_TYPE_MESSAGE
	/**
	 * 2: Stream channel.
	 */
	RtmChannelTypeSTREAM RtmChannelType = C.RTM_CHANNEL_TYPE_STREAM
	/**
	 * 3: User.
	 */
	RtmChannelTypeUSER RtmChannelType = C.RTM_CHANNEL_TYPE_USER
)

/*
*
@brief Message type when user publish message to channel or topic
*/
type RtmMessageType C.enum_C_RTM_MESSAGE_TYPE

const (
	/**
	  0: The binary message.
	*/
	RtmMessageTypeBINARY RtmMessageType = C.RTM_MESSAGE_TYPE_BINARY
	/**
	  1: The ascii message.
	*/
	RtmMessageTypeSTRING RtmMessageType = C.RTM_MESSAGE_TYPE_STRING
)

/*
*
@brief Storage type indicate the storage event was triggered by user or channel
*/
type RtmStorageType C.enum_C_RTM_STORAGE_TYPE

const (
	/**
	  0: Unknown type.
	*/
	RtmStorageTypeNONE RtmStorageType = C.RTM_STORAGE_TYPE_NONE
	/**
	  1: The user storage event.
	*/
	RtmStorageTypeUSER RtmStorageType = C.RTM_STORAGE_TYPE_USER
	/**
	  2: The channel storage event.
	*/
	RtmStorageTypeCHANNEL RtmStorageType = C.RTM_STORAGE_TYPE_CHANNEL
)

/**
* The storage event type, indicate storage operation
ref to rtm sdk
*/

/**
* The lock event type, indicate lock operation
ref to rtm sdk
*/

/**
 * The proxy type
 */
type RtmProxyType C.enum_C_RTM_PROXY_TYPE

const (
	/**
	 * 0: Link without proxy
	 */
	RtmProxyTypeNONE RtmProxyType = C.RTM_PROXY_TYPE_NONE
	/**
	 * 1: Link with http proxy
	 */
	RtmProxyTypeHTTP RtmProxyType = C.RTM_PROXY_TYPE_HTTP
	/**
	 * 2: Link with tcp cloud proxy
	 */
	RtmProxyTypeCLOUD_TCP RtmProxyType = C.RTM_PROXY_TYPE_CLOUD_TCP
)

/*
*
@brief Topic event type
ref to rtm sdk
*/

/*
*
@brief Presence event type
ref to rtm sdk
*/

/*
*
@brief Rtm link operation
*/

/**
 * Definition of LogConfiguration
 */
type RtmLogConfig struct {
	FilePath     string
	FileSizeInKB uint32
	Level        RtmLogLevel
}

func NewRtmLogConfig() *RtmLogConfig {
	config := &RtmLogConfig{
		FileSizeInKB: DEFAULT_LOG_SIZE_IN_KB,
		Level:        RtmLogLevelINFO,
	}

	return config
}
func (config *RtmLogConfig) toC() unsafe.Pointer {
	if config == nil {
		return nil
	}
	cConfig := C.C_RtmLogConfig_New()
	if cConfig == nil {
		return nil
	}
	cConfig.filePath = C.CString(config.FilePath)
	cConfig.fileSizeInKB = C.uint32_t(config.FileSizeInKB)
	cConfig.level = C.enum_C_RTM_LOG_LEVEL(config.Level)
	return unsafe.Pointer(cConfig)
}
func freeRtmLogConfig(cHandler unsafe.Pointer) {
	cConfig := (*C.struct_C_RtmLogConfig)(cHandler)
	if cConfig == nil {
		return
	}

	if cConfig.filePath != nil {
		C.free(unsafe.Pointer(cConfig.filePath))
		cConfig.filePath = nil
	}
	if cConfig != nil {
		C.C_RtmLogConfig_Delete((*C.struct_C_RtmLogConfig)(cConfig))
	}
}

// #endregion RtmLogConfig

/**
 * User list.
 */
type UserList struct {
	Users []string
}

func NewUserList() *UserList {
	userList := &UserList{
		Users: make([]string, 0),
	}

	return userList
}

/*
*
@brief Topic publisher information
*/
type PublisherInfo struct {
	UserId string
	Meta   string
}

func NewPublisherInfo() *PublisherInfo {
	publisher := &PublisherInfo{
		UserId: "",
		Meta:   "",
	}

	return publisher
}

/*
*
@brief Topic information
*/
type TopicInfo struct {
	Topic      string
	Publishers []PublisherInfo
}

func NewTopicInfo() *TopicInfo {
	topicInfo := &TopicInfo{
		Topic:      "",
		Publishers: make([]PublisherInfo, 0),
	}

	return topicInfo
}

/*
*
*
@brief User state property
*/
type StateItem struct {
	Key   string
	Value string
}

func NewStateItem() *StateItem {
	stateItem := &StateItem{
		Key:   "",
		Value: "",
	}

	return stateItem
}

/**
*  The information of a Lock.
 */
type LockDetail struct {
	LockName string
	Owner    string
	Ttl      uint32
}

func NewLockDetail() *LockDetail {
	return &LockDetail{
		LockName: "",
		Owner:    "",
		Ttl:      0,
	}
}

// 内部方法：转换为C对象

// #endregion LockDetail

/**
*  The states of user.
 */
type UserState struct {
	UserId string
	States []StateItem
}

func NewUserState() *UserState {
	userState := &UserState{
		UserId: "",
		States: make([]StateItem, 0),
	}

	return userState
}

/**
 *  The subscribe option.
 */
type SubscribeOptions struct {
	WithMessage  bool
	WithMetadata bool
	WithPresence bool
	WithLock     bool
	BeQuiet      bool
}

func NewSubscribeOptions() *SubscribeOptions {
	return &SubscribeOptions{
		WithMessage:  true,
		WithMetadata: false,
		WithPresence: true,
		WithLock:     false,
		BeQuiet:      false,
	}
}

func (this_ *SubscribeOptions) toC() unsafe.Pointer {
	if this_ == nil {
		return nil
	}
	cOpt := C.C_SubscribeOptions_New()
	if cOpt == nil {
		return nil
	}
	cOpt.withMessage = C.bool(this_.WithMessage)
	cOpt.withMetadata = C.bool(this_.WithMetadata)
	cOpt.withPresence = C.bool(this_.WithPresence)
	cOpt.withLock = C.bool(this_.WithLock)
	cOpt.beQuiet = C.bool(this_.BeQuiet)
	return unsafe.Pointer(cOpt)
}

func freeSubscribeOptions(cOpt unsafe.Pointer) {
	if cOpt != nil {
		C.C_SubscribeOptions_Delete((*C.struct_C_SubscribeOptions)(cOpt))
	}
}

// #endregion SubscribeOptions

/**
 *  The channel information.
 */
type ChannelInfo struct {
	ChannelName string
	ChannelType RtmChannelType
}

func NewChannelInfo() *ChannelInfo {
	channelInfo := &ChannelInfo{
		ChannelName: "",
		ChannelType: 0,
	}

	return channelInfo
}

/**
 *  The option to query user presence.
 */
type PresenceOptions struct {
	IncludeUserId bool
	IncludeState  bool
	Page          string
}

func NewPresenceOptions() *PresenceOptions {
	presenceOptions := &PresenceOptions{
		IncludeUserId: true,
		IncludeState:  false,
		Page:          "",
	}

	return presenceOptions
}

// #endregion PresenceOptions

/**
*  The option to query user presence.
 */
type GetOnlineUsersOptions struct {
	IncludeUserId bool
	IncludeState  bool
	Page          string
}

func NewGetOnlineUsersOptions() *GetOnlineUsersOptions {
	getOnlineUsersOptions := &GetOnlineUsersOptions{
		IncludeUserId: false,
		IncludeState:  false,
		Page:          "",
	}

	return getOnlineUsersOptions
}

/*
*

	@brief Publish message option
*/
type PublishOptions struct {
	ChannelType    RtmChannelType
	MessageType    RtmMessageType
	CustomType     string
	StoreInHistory bool
}

func NewPublishOptions() *PublishOptions {
	publishOptions := &PublishOptions{
		ChannelType:    0,
		MessageType:    0,
		CustomType:     "",
		StoreInHistory: false,
	}

	return publishOptions
}

func (this_ *PublishOptions) toC() unsafe.Pointer {
	if this_ == nil {
		return nil
	}

	cOpt := C.C_PublishOptions_New()
	if cOpt == nil {
		return nil
	}

	cOpt.channelType = C.enum_C_RTM_CHANNEL_TYPE(this_.ChannelType)
	cOpt.messageType = C.enum_C_RTM_MESSAGE_TYPE(this_.MessageType)

	if this_.CustomType != "" {
		cOpt.customType = C.CString(this_.CustomType)
	} else {
		cOpt.customType = nil
	}

	cOpt.storeInHistory = C.bool(this_.StoreInHistory)

	return unsafe.Pointer(cOpt)
}

func freePublishOptions(cOpt unsafe.Pointer) {
	if cOpt != nil {
		cPublishOpt := (*C.struct_C_PublishOptions)(cOpt)
		if cPublishOpt.customType != nil {
			C.free(unsafe.Pointer(cPublishOpt.customType))
		}

		C.C_PublishOptions_Delete((*C.struct_C_PublishOptions)(cOpt))
	}
}

// #endregion PublishOptions

/*
@brief topic message option
*/
type TopicMessageOptions struct {
	MessageType RtmMessageType
	SendTs      uint64
	CustomType  string
}

func NewTopicMessageOptions() *TopicMessageOptions {
	topicMessageOptions := &TopicMessageOptions{
		MessageType: 0,
		SendTs:      0,
		CustomType:  "",
	}

	return topicMessageOptions
}

/*
*
@brief Proxy configuration
*/
type RtmProxyConfig struct {
	ProxyType RtmProxyType
	Server    string
	Port      uint16
	Account   string
	Password  string
}

func NewRtmProxyConfig() *RtmProxyConfig {
	rtmProxyConfig := &RtmProxyConfig{
		ProxyType: 0,
		Server:    "",
		Port:      0,
		Account:   "",
		Password:  "",
	}

	return rtmProxyConfig
}
func (this_ *RtmProxyConfig) toC() unsafe.Pointer {
	if this_ == nil {
		return nil
	}
	cProxyConfig := C.C_RtmProxyConfig_New()
	if cProxyConfig == nil {
		return nil
	}
	cProxyConfig.proxyType = C.enum_C_RTM_PROXY_TYPE(this_.ProxyType)
	cProxyConfig.server = C.CString(this_.Server)
	cProxyConfig.port = C.uint16_t(this_.Port)
	cProxyConfig.account = C.CString(this_.Account)
	cProxyConfig.password = C.CString(this_.Password)
	return unsafe.Pointer(cProxyConfig)
}
func freeRtmProxyConfig(cProxyHandler unsafe.Pointer) {
	cProxyConfig := (*C.struct_C_RtmProxyConfig)(cProxyHandler)
	if cProxyConfig == nil {
		return
	}
	if cProxyConfig.server != nil {
		// free server
		C.free(unsafe.Pointer(cProxyConfig.server))
		cProxyConfig.server = nil
	}
	if cProxyConfig.account != nil {
		C.free(unsafe.Pointer(cProxyConfig.account))
		cProxyConfig.account = nil
	}
	if cProxyConfig.password != nil {
		C.free(unsafe.Pointer(cProxyConfig.password))
		cProxyConfig.password = nil
	}
	C.C_RtmProxyConfig_Delete((*C.struct_C_RtmProxyConfig)(cProxyConfig))
}

/*
*
@brief encryption configuration
*/
type RtmEncryptionConfig struct {
	EncryptionMode RtmEncryptionMode
	EncryptionKey  string
	EncryptionSalt []byte
}

func NewRtmEncryptionConfig() *RtmEncryptionConfig {
	rtmEncryptionConfig := &RtmEncryptionConfig{
		EncryptionMode: RtmEncryptionModeNONE,
		EncryptionKey:  "",
		EncryptionSalt: make([]byte, 32),
	}

	return rtmEncryptionConfig
}
func (this_ *RtmEncryptionConfig) toC() unsafe.Pointer {
	if this_ == nil {
		return nil
	}
	cEncryptionConfig := C.C_RtmEncryptionConfig_New()
	if cEncryptionConfig == nil {
		return nil
	}
	cEncryptionConfig.encryptionMode = C.enum_C_RTM_ENCRYPTION_MODE(this_.EncryptionMode)
	cEncryptionConfig.encryptionKey = C.CString(this_.EncryptionKey)
	// calc salt length from c type: uint8_t[32]
	saltLength := len(this_.EncryptionSalt)
	fixLen := int(unsafe.Sizeof(cEncryptionConfig.encryptionSalt) / unsafe.Sizeof(cEncryptionConfig.encryptionSalt[0]))
	// init salt to 0
	C.memset(unsafe.Pointer(&cEncryptionConfig.encryptionSalt[0]), 0, C.size_t(fixLen))
	if saltLength > fixLen {
		saltLength = fixLen
	}
	if saltLength > 0 {
		C.memcpy(unsafe.Pointer(&cEncryptionConfig.encryptionSalt[0]), unsafe.Pointer(&this_.EncryptionSalt[0]), C.size_t(saltLength))
	}
	return unsafe.Pointer(cEncryptionConfig)
}
func freeRtmEncryptionConfig(cHandler unsafe.Pointer) {
	cEncryptionConfig := (*C.struct_C_RtmEncryptionConfig)(cHandler)
	if cEncryptionConfig == nil {
		return
	}
	if cEncryptionConfig.encryptionKey != nil {
		// free encryption key
		C.free(unsafe.Pointer(cEncryptionConfig.encryptionKey))
		// and set to nil
		cEncryptionConfig.encryptionKey = nil
	}

	C.C_RtmEncryptionConfig_Delete((*C.struct_C_RtmEncryptionConfig)(cEncryptionConfig))

}

// #region RtmPrivateConfig
type RtmPrivateConfig struct {
	ServiceType      uint32
	AccessPointHosts []string
}

func NewRtmPrivateConfig() *RtmPrivateConfig {
	rtmPrivateConfig := &RtmPrivateConfig{
		ServiceType:      0,
		AccessPointHosts: nil,
	}
	return rtmPrivateConfig
}
func (this_ *RtmPrivateConfig) toC() unsafe.Pointer {
	if this_ == nil {
		return nil
	}
	accessPointHostsCount := len(this_.AccessPointHosts)
	cPrivateConfig := (*C.struct_C_RtmPrivateConfig)(C.malloc(C.sizeof_struct_C_RtmPrivateConfig))

	cPrivateConfig.serviceType = C.enum_C_RTM_SERVICE_TYPE(this_.ServiceType)

	cPrivateConfig.accessPointHostsCount = C.size_t(accessPointHostsCount)

	// allocate *char[accessPointHostsCount]
	if accessPointHostsCount > 0 {
		cPrivateConfig.accessPointHosts = (**C.char)(C.malloc(C.size_t(len(this_.AccessPointHosts)) * C.size_t(unsafe.Sizeof(uintptr(0)))))

		for i, host := range this_.AccessPointHosts {
			cStr := C.CString(host)
			// 计算指针位置并赋值
			ptr := (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(cPrivateConfig.accessPointHosts)) + uintptr(i)*unsafe.Sizeof(uintptr(0))))
			*ptr = cStr
		}
	} else {
		cPrivateConfig.accessPointHosts = nil
	}

	return unsafe.Pointer(cPrivateConfig)
}
func freeRtmPrivateConfig(cHandler unsafe.Pointer) {
	if cHandler != nil {
		cPrivateConfig := (*C.struct_C_RtmPrivateConfig)(cHandler)
		count := int(cPrivateConfig.accessPointHostsCount)
		if cPrivateConfig.accessPointHosts != nil && count > 0 {
			for i := 0; i < count; i++ {
				ptr := (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(cPrivateConfig.accessPointHosts)) + uintptr(i)*unsafe.Sizeof(uintptr(0))))
				C.free(unsafe.Pointer(*ptr))
			}
			C.free(unsafe.Pointer(cPrivateConfig.accessPointHosts))
		}
		C.free(unsafe.Pointer(cPrivateConfig))
	}
}

// link state event

type RtmServiceType C.enum_C_RTM_SERVICE_TYPE

const (
	RtmServiceTypeNONE    RtmServiceType = C.RTM_SERVICE_TYPE_NONE
	RtmServiceTypeMESSAGE RtmServiceType = C.RTM_SERVICE_TYPE_MESSAGE
	RtmServiceTypeSTREAM  RtmServiceType = C.RTM_SERVICE_TYPE_STREAM
)

type LinkStateEvent struct {
	CurrentState       int
	PreviousState      int
	ServiceType        RtmServiceType
	Operation          int
	ReasonCode         int
	Reason             string
	AffectedChannels   []string
	UnrestoredChannels []string
	IsResumed          bool
	Timestamp          uint64
}

func NewLinkStateEvent() *LinkStateEvent {
	linkStateEvent := &LinkStateEvent{
		CurrentState:       0,
		PreviousState:      0,
		ServiceType:        RtmServiceTypeNONE,
		Operation:          0,
		ReasonCode:         0,
		Reason:             "",
		AffectedChannels:   make([]string, 0),
		UnrestoredChannels: make([]string, 0),
		IsResumed:          false,
		Timestamp:          0,
	}

	return linkStateEvent
}

// #region LinkStateEvent

type HistoryMessage struct {
	MessageType   RtmMessageType
	Publisher     string
	Message       string
	MessageLength uint
	CustomType    string
	Timestamp     uint64
}

func NewHistoryMessage() *HistoryMessage {
	historyMessage := &HistoryMessage{
		MessageType:   0,
		Publisher:     "",
		Message:       "",
		MessageLength: 0,
		CustomType:    "",
		Timestamp:     0,
	}

	return historyMessage
}
