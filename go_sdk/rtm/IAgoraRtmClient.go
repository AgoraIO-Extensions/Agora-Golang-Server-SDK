package agorartm

/*
#include "C_IAgoraRtmClient.h"
#include "C_AgoraRtmBase.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

// #region agora

// #region agora::rtm

/**
 *  Configurations for RTM Client.
 */
type RtmConfig struct {
	AppId             string
	UserId            string
	AreaCode          RtmAreaCode
	ProtocolType      uint32
	PresenceTimeout   uint32
	HeartbeatInterval uint32
	Context           unsafe.Pointer
	UseStringUserId   bool
	Multipath         bool
	EventHandler      *RtmEventHandler
	LogConfig         *RtmLogConfig
	ProxyConfig       *RtmProxyConfig
	EncryptionConfig  *RtmEncryptionConfig
	PrivateConfig     *RtmPrivateConfig
}

/**
 * - For Android, it is the context of Activity or Application.
 * - For Windows, it is the window handle of app. Once set, this parameter enables you to plug
 * or unplug the video devices while they are powered.
 */

func NewRtmConfig() *RtmConfig {
	config := &RtmConfig{
		AppId:             "",
		UserId:            "",
		AreaCode:          RtmAreaCodeGLOB,
		ProtocolType:      0,
		HeartbeatInterval: 0,
		Context:           nil,
		UseStringUserId:   false,
		Multipath:         false,
		EventHandler: nil,
		LogConfig:         nil,
		ProxyConfig:       nil,
		EncryptionConfig:  nil,
		PrivateConfig:     nil,
		PresenceTimeout:   30,
	}

	return config
}

/**
 * The IRtmEventHandler class.
 *
 * The SDK uses this class to send callback event notifications to the app, and the app inherits
 * the methods in this class to retrieve these event notifications.
 *
 * All methods in this class have their default (empty)  implementations, and the app can inherit
 * only some of the required events instead of all. In the callback methods, the app should avoid
 * time-consuming tasks or calling blocking APIs, otherwise the SDK may not work properly.
 */
// old IRtmEventHandler is deleted, please use new RtmEventHandler interface

// new user friendly event handler interface design

// RtmEventHandler define the event handler interface that user can implement
// user only need to implement the needed methods, the methods that are not implemented will be ignored by SDK
//
// commonly used callback methods:
//   - OnLoginResult: login result callback
//   - OnLogoutResult: logout result callback
//   - OnMessageEvent: message event callback
//   - OnPresenceEvent: online status event callback
//   - OnSubscribeResult: subscribe result callback
//   - OnPublishResult: publish result callback
//
// usage one (object oriented):
//
//	type MyEventHandler struct{}
//	func (h *MyEventHandler) OnLoginResult(requestId uint64, errorCode RTM_ERROR_CODE) {
//	    // handle login result
//	}
//	rtmConfig.SetEventHandler(&MyEventHandler{})
//
// usage two (functional):
//
//	handler := &RtmEventHandlerConfig{
//	    OnLoginResult: func(requestId uint64, errorCode RTM_ERROR_CODE) {
//	        // handle login result
//	    },
//	}






// #region MessageEvent
type MessageEvent struct {
	ChannelType  RtmChannelType
	MessageType  RtmMessageType
	ChannelName  string
	ChannelTopic string
	Message      []byte
	Publisher    string
	CustomType   string
}

func NewMessageEvent() *MessageEvent {
	event := &MessageEvent{
		ChannelType:  RtmChannelTypeNONE,
		MessageType:  RtmMessageTypeSTRING,
		ChannelName:  "",
		ChannelTopic: "",
		Message:      make([]byte, 0),
		Publisher:    "",
		CustomType:   "",
	}

	return event
}

func (this_ *MessageEvent) fromC(cEvent *C.struct_C_MessageEvent) {
	if cEvent == nil {
		return
	}

	if !IsValidMemory(unsafe.Pointer(cEvent)) {
		return
	}

	this_.ChannelType = RtmChannelType(cEvent.channelType)
	this_.MessageType = RtmMessageType(cEvent.messageType)

	if cEvent.channelName != nil {
		this_.ChannelName = C.GoString(cEvent.channelName)
	}
	if cEvent.channelTopic != nil {
		this_.ChannelTopic = C.GoString(cEvent.channelTopic)
	}
	if cEvent.publisher != nil {
		this_.Publisher = C.GoString(cEvent.publisher)
	}
	if cEvent.customType != nil {
		this_.CustomType = C.GoString(cEvent.customType)
	}

	if cEvent.message != nil && cEvent.messageLength > 0 {
		this_.Message = C.GoBytes(unsafe.Pointer(cEvent.message), C.int(cEvent.messageLength))
	}
}

// #endregion MessageEvent

type IntervalInfo struct {
	JoinUserList    *UserList
	LeaveUserList   *UserList
	TimeoutUserList *UserList
	UserStateList   []*UserState
	UserStateCount  uint
}

func NewIntervalInfo() *IntervalInfo {
	info := &IntervalInfo{
		JoinUserList:    NewUserList(),
		LeaveUserList:   NewUserList(),
		TimeoutUserList: NewUserList(),
		UserStateList:   make([]*UserState, 0),
		UserStateCount:  0,
	}

	return info
}

// #endregion IntervalInfo

type SnapshotInfo struct {
	UserStateList []*UserState
	UserCount     uint
}

func NewSnapshotInfo() *SnapshotInfo {
	info := &SnapshotInfo{
		UserStateList: make([]*UserState, 0),
		UserCount:     0,
	}

	return info
}

// #endregion SnapshotInfo

type PresenceEvent struct {
	Type           int
	ChannelType    RtmChannelType
	ChannelName    string
	Publisher      string
	StateItems     []*StateItem
	StateItemCount uint
	Interval       *IntervalInfo
	Snapshot       *SnapshotInfo
}

func NewPresenceEvent() *PresenceEvent {
	event := &PresenceEvent{
		Type:           0,
		ChannelType:    0,
		ChannelName:    "",
		Publisher:      "",
		StateItems:     make([]*StateItem, 0),
		StateItemCount: 0,
		Interval:       nil,
		Snapshot:       nil,
	}

	return event
}

func (this_ *PresenceEvent) fromC(cEvent *C.struct_C_PresenceEvent) {
	if cEvent == nil {
		return
	}

	if !IsValidMemory(unsafe.Pointer(cEvent)) {
		return
	}

	this_.Type = int(cEvent._type)
	this_.ChannelType = RtmChannelType(cEvent.channelType)

	if cEvent.channelName != nil {
		this_.ChannelName = C.GoString(cEvent.channelName)
	}
	if cEvent.publisher != nil {
		this_.Publisher = C.GoString(cEvent.publisher)
	}

	if cEvent.stateItems != nil && cEvent.stateItemCount > 0 {
		if IsValidMemory(unsafe.Pointer(cEvent.stateItems)) {
			itemCount := int(cEvent.stateItemCount)
			if itemCount > 0 {
				this_.StateItems = make([]*StateItem, itemCount)
				this_.StateItemCount = uint(itemCount)

				for i := 0; i < itemCount; i++ {
					cItem := (*C.struct_C_StateItem)(unsafe.Pointer(uintptr(unsafe.Pointer(cEvent.stateItems)) + uintptr(i)*unsafe.Sizeof(C.struct_C_StateItem{})))
					if cItem != nil && IsValidMemory(unsafe.Pointer(cItem)) {
						stateItem := NewStateItem()
						if stateItem != nil {
							if cItem.key != nil {
								stateItem.Key = FastSafeCGoString(cItem.key)
							}
							if cItem.value != nil {
								stateItem.Value = FastSafeCGoString(cItem.value)
							}
							this_.StateItems[i] = stateItem
						}
					}
				}
			} else {
				this_.StateItems = make([]*StateItem, 0)
				this_.StateItemCount = 0
			}
		} else {
			this_.StateItems = make([]*StateItem, 0)
			this_.StateItemCount = 0
		}
	} else {
		this_.StateItems = make([]*StateItem, 0)
		this_.StateItemCount = 0
	}

	this_.Interval = NewIntervalInfo()
	if this_.Interval != nil {
		if cEvent.interval.userStateList != nil && cEvent.interval.userStateCount > 0 {
			if IsValidMemory(unsafe.Pointer(cEvent.interval.userStateList)) {
				userCount := int(cEvent.interval.userStateCount)
				if userCount > 0 {
					this_.Interval.UserStateList = make([]*UserState, userCount)
					this_.Interval.UserStateCount = uint(userCount)

					for i := 0; i < userCount; i++ {
						cUserState := (*C.struct_C_UserState)(unsafe.Pointer(uintptr(unsafe.Pointer(cEvent.interval.userStateList)) + uintptr(i)*unsafe.Sizeof(C.struct_C_UserState{})))
						if cUserState != nil && IsValidMemory(unsafe.Pointer(cUserState)) {
							userState := NewUserState()
							if userState != nil {
								if cUserState.userId != nil {
									userState.UserId = FastSafeCGoString(cUserState.userId)
								}
								if cUserState.states != nil && cUserState.statesCount > 0 {
									if IsValidMemory(unsafe.Pointer(cUserState.states)) {
										stateCount := int(cUserState.statesCount)
										if stateCount > 0 {
											userState.States = make([]StateItem, stateCount)
											//userState.StatesCount = uint(stateCount)

											for j := 0; j < stateCount; j++ {
												cState := (*C.struct_C_StateItem)(unsafe.Pointer(uintptr(unsafe.Pointer(cUserState.states)) + uintptr(j)*unsafe.Sizeof(C.struct_C_StateItem{})))
												if cState != nil && IsValidMemory(unsafe.Pointer(cState)) {
													stateItem := StateItem{}
													if cState.key != nil {
														stateItem.Key = FastSafeCGoString(cState.key)
													}
													if cState.value != nil {
														stateItem.Value = FastSafeCGoString(cState.value)
													}
													userState.States[j] = stateItem
												}
											}
										}
									}
								}
								this_.Interval.UserStateList[i] = userState
							}
						}
					}
				}
			}
		}
	}

	this_.Snapshot = NewSnapshotInfo()
	if this_.Snapshot != nil {
		if cEvent.snapshot.userStateList != nil && cEvent.snapshot.userCount > 0 {
			if IsValidMemory(unsafe.Pointer(cEvent.snapshot.userStateList)) {
				userCount := int(cEvent.snapshot.userCount)
				if userCount > 0 {
					this_.Snapshot.UserStateList = make([]*UserState, userCount)
					this_.Snapshot.UserCount = uint(userCount)

					for i := 0; i < userCount; i++ {
						cUserState := (*C.struct_C_UserState)(unsafe.Pointer(uintptr(unsafe.Pointer(cEvent.snapshot.userStateList)) + uintptr(i)*unsafe.Sizeof(C.struct_C_UserState{})))
						if cUserState != nil && IsValidMemory(unsafe.Pointer(cUserState)) {
							userState := NewUserState()
							if userState != nil {
								if cUserState.userId != nil {
									userState.UserId = FastSafeCGoString(cUserState.userId)
								}
								if cUserState.states != nil && cUserState.statesCount > 0 {
									if IsValidMemory(unsafe.Pointer(cUserState.states)) {
										stateCount := int(cUserState.statesCount)
										if stateCount > 0 {
											userState.States = make([]StateItem, stateCount)
											//userState.StatesCount = uint(stateCount)

											for j := 0; j < stateCount; j++ {
												cState := (*C.struct_C_StateItem)(unsafe.Pointer(uintptr(unsafe.Pointer(cUserState.states)) + uintptr(j)*unsafe.Sizeof(C.struct_C_StateItem{})))
												if cState != nil && IsValidMemory(unsafe.Pointer(cState)) {
													stateItem := StateItem{}
													if cState.key != nil {
														stateItem.Key = FastSafeCGoString(cState.key)
													}
													if cState.value != nil {
														stateItem.Value = FastSafeCGoString(cState.value)
													}
													userState.States[j] = stateItem
												}
											}
										}
									}
								}
								this_.Snapshot.UserStateList[i] = userState
							}
						}
					}
				}
			}
		}
	}
}

// #endregion PresenceEvent

type TopicEvent struct {
	Type           int
	ChannelName    string
	Publisher      string
	TopicInfos     []*TopicInfo
	TopicInfoCount uint
}

func NewTopicEvent() *TopicEvent {
	event := &TopicEvent{
		Type:           0,
		ChannelName:    "",
		Publisher:      "",
		TopicInfos:     make([]*TopicInfo, 0),
		TopicInfoCount: 0,
	}

	return event
}

func (this_ *TopicEvent) fromC(cEvent *C.struct_C_TopicEvent) {
	if cEvent == nil {
		return
	}

	if !IsValidMemory(unsafe.Pointer(cEvent)) {
		return
	}

	this_.Type = int(cEvent._type)

	if cEvent.channelName != nil {
		this_.ChannelName = C.GoString(cEvent.channelName)
	}
	if cEvent.publisher != nil {
		this_.Publisher = C.GoString(cEvent.publisher)
	}

	if cEvent.topicInfos != nil && cEvent.topicInfoCount > 0 {
		if IsValidMemory(unsafe.Pointer(cEvent.topicInfos)) {
			infoCount := int(cEvent.topicInfoCount)
			if infoCount > 0 {
				this_.TopicInfos = make([]*TopicInfo, infoCount)
				this_.TopicInfoCount = uint(infoCount)

				for i := 0; i < infoCount; i++ {
					cTopicInfo := (*C.struct_C_TopicInfo)(unsafe.Pointer(uintptr(unsafe.Pointer(cEvent.topicInfos)) + uintptr(i)*unsafe.Sizeof(C.struct_C_TopicInfo{})))
					if cTopicInfo != nil && IsValidMemory(unsafe.Pointer(cTopicInfo)) {
						topicInfo := NewTopicInfo()
						if topicInfo != nil {
							if cTopicInfo.topic != nil {
								topicInfo.Topic = FastSafeCGoString(cTopicInfo.topic)
							}
							if cTopicInfo.publishers != nil && cTopicInfo.publisherCount > 0 {
								if IsValidMemory(unsafe.Pointer(cTopicInfo.publishers)) {
									pubCount := int(cTopicInfo.publisherCount)
									if pubCount > 0 {
										topicInfo.Publishers = make([]PublisherInfo, pubCount)
										//topicInfo.PublisherCount = uint(pubCount)

										for j := 0; j < pubCount; j++ {
											cPublisher := (*C.struct_C_PublisherInfo)(unsafe.Pointer(uintptr(unsafe.Pointer(cTopicInfo.publishers)) + uintptr(j)*unsafe.Sizeof(C.struct_C_PublisherInfo{})))
											if cPublisher != nil && IsValidMemory(unsafe.Pointer(cPublisher)) {
												publisherInfo := PublisherInfo{}
												if cPublisher.publisherUserId != nil {
													publisherInfo.UserId = FastSafeCGoString(cPublisher.publisherUserId)
												}
												if cPublisher.publisherMeta != nil {
													publisherInfo.Meta = FastSafeCGoString(cPublisher.publisherMeta)
												}
												topicInfo.Publishers[j] = publisherInfo
											}
										}
									}
								} else {
									topicInfo.Publishers = make([]PublisherInfo, 0)
									//topicInfo.PublisherCount = 0
								}
							} else {
								topicInfo.Publishers = make([]PublisherInfo, 0)
								//topicInfo.PublisherCount = 0
							}
							this_.TopicInfos[i] = topicInfo
						}
					}
				}
			}
		}
	} else {
		this_.TopicInfos = make([]*TopicInfo, 0)
		this_.TopicInfoCount = 0
	}
}

// #endregion TopicEvent

type LockEvent struct {
	ChannelType    RtmChannelType
	EventType      int
	ChannelName    string
	LockDetailList []*LockDetail
	Count          uint
}

func NewLockEvent() *LockEvent {
	event := &LockEvent{
		ChannelType:    RtmChannelTypeNONE,
		EventType:      0,
		ChannelName:    "",
		LockDetailList: make([]*LockDetail, 0),
		Count:          0,
	}

	return event
}

func (this_ *LockEvent) fromC(cEvent *C.struct_C_LockEvent) {
	if cEvent == nil {
		return
	}

	if !IsValidMemory(unsafe.Pointer(cEvent)) {
		return
	}

	this_.ChannelType = RtmChannelType(cEvent.channelType)
	this_.EventType = int(cEvent.eventType)

	if cEvent.channelName != nil {
		this_.ChannelName = C.GoString(cEvent.channelName)
	}

	if cEvent.lockDetailList != nil && cEvent.count > 0 {
		if IsValidMemory(unsafe.Pointer(cEvent.lockDetailList)) {
			detailCount := int(cEvent.count)
			if detailCount > 0 {
				this_.LockDetailList = make([]*LockDetail, detailCount)
				this_.Count = uint(detailCount)

				for i := 0; i < detailCount; i++ {
					cLockDetail := (*C.struct_C_LockDetail)(unsafe.Pointer(uintptr(unsafe.Pointer(cEvent.lockDetailList)) + uintptr(i)*unsafe.Sizeof(C.struct_C_LockDetail{})))
					if cLockDetail != nil && IsValidMemory(unsafe.Pointer(cLockDetail)) {
						lockDetail := NewLockDetail()
						if lockDetail != nil {
							if cLockDetail.lockName != nil {
								lockDetail.LockName = (FastSafeCGoString(cLockDetail.lockName))
							}
							if cLockDetail.owner != nil {
								lockDetail.Owner = (FastSafeCGoString(cLockDetail.owner))
							}
							lockDetail.Ttl = uint32(cLockDetail.ttl)
							this_.LockDetailList[i] = lockDetail
						}
					}
				}
			}
		}
	} else {
		this_.LockDetailList = make([]*LockDetail, 0)
		this_.Count = 0
	}
}

// #endregion LockEvent

type StorageEvent struct {
	ChannelType RtmChannelType
	StorageType RtmStorageType
	EventType   int
	Target      string
	Data        *IMetadata
}

func NewStorageEvent() *StorageEvent {
	event := &StorageEvent{
		ChannelType: RtmChannelTypeNONE,
		StorageType: RtmStorageTypeNONE,
		EventType:   0,
		Target:      "",
		Data:        nil,
	}

	return event
}

func (this_ *StorageEvent) fromC(cEvent *C.struct_C_StorageEvent) {
	if cEvent == nil {
		return
	}

	if !IsValidMemory(unsafe.Pointer(cEvent)) {
		return
	}

	this_.ChannelType = RtmChannelType(cEvent.channelType)
	this_.StorageType = RtmStorageType(cEvent.storageType)
	this_.EventType = int(cEvent.eventType)

	if cEvent.target != nil {
		this_.Target = C.GoString(cEvent.target)
	}

	if cEvent.data != nil {
		this_.Data = CMetadataToIMetadata(cEvent.data)
	}
}

type IRtmClient struct {
	rtmClient  unsafe.Pointer
	handler     *RtmEventHandler
	history    *IRtmHistory
	presence   *IRtmPresence
	lock       *IRtmLock
	storage    *IRtmStorage
	isLoggedIn bool
	cEventHandler *C.struct_C_IRtmEventHandler
}

/**
 * Initializes the rtm client instance.
 *
 * @param [in] config The configurations for RTM Client.
 * @param [in] eventHandler .
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
/**
 * Creates the rtm client object and returns the pointer.
 *
 * @return Pointer of the rtm client object.
 */
func NewRtmClient(config *RtmConfig) *IRtmClient {
	if config == nil {
		return nil
	}

	cConfig := C.C_RtmConfig_New()
	if cConfig == nil {
		return nil
	}
	defer C.C_RtmConfig_Delete(cConfig)

	cConfig.appId = C.CString(config.AppId)
	defer C.free(unsafe.Pointer(cConfig.appId))
	cConfig.userId = C.CString(config.UserId)
	defer C.free(unsafe.Pointer(cConfig.userId))
	cConfig.areaCode = C.enum_C_RTM_AREA_CODE(config.AreaCode)
	cConfig.protocolType = config.ProtocolType
	cConfig.presenceTimeout = C.uint32_t(config.PresenceTimeout)
	cConfig.heartbeatInterval = C.uint32_t(config.HeartbeatInterval)
	cConfig.multipath = C.bool(config.Multipath)
	cConfig.context = config.Context
	cConfig.useStringUserId = C.bool(config.UseStringUserId)

	// add log config
	cLogConfig := (*C.struct_C_RtmLogConfig)(nil)
	if config.LogConfig != nil {
		cLogConfig = (*C.struct_C_RtmLogConfig)(config.LogConfig.toC())
		cConfig.logConfig.filePath = cLogConfig.filePath
		cConfig.logConfig.fileSizeInKB = cLogConfig.fileSizeInKB
		cConfig.logConfig.level = cLogConfig.level
	}

	defer freeRtmLogConfig(unsafe.Pointer(cLogConfig))

	// add encryption config
	cEncryptionConfig := (*C.struct_C_RtmEncryptionConfig)(nil)
	if config.EncryptionConfig != nil {
		cEncryptionConfig = (*C.struct_C_RtmEncryptionConfig)(config.EncryptionConfig.toC())
		cConfig.encryptionConfig.encryptionMode = cEncryptionConfig.encryptionMode
		cConfig.encryptionConfig.encryptionKey = cEncryptionConfig.encryptionKey
		cConfig.encryptionConfig.encryptionSalt = cEncryptionConfig.encryptionSalt
	}
	defer freeRtmEncryptionConfig(unsafe.Pointer(cEncryptionConfig))

	// add proxy config
	cProxyConfig := (*C.struct_C_RtmProxyConfig)(nil)
	if config.ProxyConfig != nil {
		cProxyConfig = (*C.struct_C_RtmProxyConfig)(config.ProxyConfig.toC())
		cConfig.proxyConfig.proxyType = cProxyConfig.proxyType
		cConfig.proxyConfig.server = cProxyConfig.server
	}
	defer freeRtmProxyConfig(unsafe.Pointer(cProxyConfig))

	// add private config
	cPrivateConfig := (*C.struct_C_RtmPrivateConfig)(nil)
	if config.PrivateConfig != nil {
		cPrivateConfig = (*C.struct_C_RtmPrivateConfig)(config.PrivateConfig.toC())
		cConfig.privateConfig.serviceType = cPrivateConfig.serviceType
		cConfig.privateConfig.accessPointHosts = cPrivateConfig.accessPointHosts
		cConfig.privateConfig.accessPointHostsCount = cPrivateConfig.accessPointHostsCount
	}
	defer freeRtmPrivateConfig(unsafe.Pointer(cPrivateConfig))

	
	var cEventHandler *C.struct_C_IRtmEventHandler = nil
	if config.EventHandler != nil {
		// allocate a c event handler, and keep it alive
		//userData := unsafe.Pointer(config.EventHandlerConfig)
		cEventHandler = CRtmEventHandler()
		cConfig.eventHandler = cEventHandler
	} else {
		cConfig.eventHandler = nil
	}

	client := &IRtmClient{
		rtmClient:  nil,
		handler:    config.EventHandler,
		cEventHandler: cEventHandler,
		history:    nil,
		presence:   nil,
		lock:       nil,
		storage:    nil,
		isLoggedIn: false,
	}

	cEventHandler.userData = unsafe.Pointer(client)


	var errorCode C.int
	rtmClient := C.agora_rtm_client_create(cConfig, &errorCode)

	if rtmClient == nil {
		return nil
	}

	client.rtmClient = unsafe.Pointer(rtmClient)

	//note : cEventHandler.userData will be equal to client!!
	// assign userdata
	


	// get storage
	cStorage := C.agora_rtm_client_get_storage(client.rtmClient)
	if cStorage == nil {
		return nil
	}
	client.storage = &IRtmStorage{rtmStorage: unsafe.Pointer(cStorage)}

	// get lock
	cLock := C.agora_rtm_client_get_lock(client.rtmClient)
	if cLock == nil {
		return nil
	}
	client.lock = &IRtmLock{rtmLock: unsafe.Pointer(cLock)}

	// get presence
	cPresence := C.agora_rtm_client_get_presence(client.rtmClient)
	if cPresence == nil {
		return nil
	}
	client.presence = &IRtmPresence{rtmPresence: unsafe.Pointer(cPresence)}

	// get history
	cHistory := C.agora_rtm_client_get_history(client.rtmClient)
	if cHistory == nil {
		return nil
	}
	client.history = &IRtmHistory{rtmHistory: unsafe.Pointer(cHistory)}

	return client
}

/**
 * Release the rtm client instance.
 *
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmClient) Release() int {
	// validity check
	if this_.rtmClient == nil {
		return -10001
	}

	// do really release
	ret := int(C.agora_rtm_client_release(this_.rtmClient))

	if this_.cEventHandler != nil {
		C.C_IRtmEventHandler_Delete(this_.cEventHandler)
		this_.cEventHandler = nil
	}
	this_.rtmClient = nil
	this_.history = nil
	this_.presence = nil
	this_.lock = nil
	this_.storage = nil
	this_.isLoggedIn = false
	this_.handler = nil

	return ret
}

/**
 * Login the Agora RTM service. The operation result will be notified by \ref agora::rtm::IRtmEventHandler::onLoginResult callback.
 *
 * @param [in] token Token used to login RTM service.
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmClient) Login(token string) (int, uint64) {

	// check if already logged in
	if this_.rtmClient == nil {
		return -10002, 0
	}
	if this_.isLoggedIn {
		return -10003, 0
	}

	// do really login
	var requestId uint64
	cToken := C.CString(token)
	defer C.free(unsafe.Pointer(cToken))
	ret := int(C.agora_rtm_client_login(this_.rtmClient,
		cToken,
		(*C.uint64_t)(unsafe.Pointer(&requestId)),
	))

	// update login status
	if ret == 0 {
		this_.isLoggedIn = true
	} else {
		this_.isLoggedIn = false
	}
	return int(ret), requestId
}

/**
 * Logout the Agora RTM service. Be noticed that this method will break the rtm service including storage/lock/presence.
 *
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmClient) Logout() (int, uint64) {
	// validity check, to avoid repeat logout
	if this_.rtmClient == nil {
		return -10004, 0
	}
	if !this_.isLoggedIn {
		return -10005, 0
	}
	var requestId uint64
	ret := int(C.agora_rtm_client_logout(this_.rtmClient,
		(*C.uint64_t)(unsafe.Pointer(&requestId)),
	))

	// update login status
	if ret == 0 {
		this_.isLoggedIn = false
	}
	return ret, requestId
}

/**
 * Get the storage instance.
 *
 * @return
 * - return NULL if error occurred
 */
func (this_ *IRtmClient) GetStorage() *IRtmStorage {
	// validity check
	if this_.storage == nil || this_.rtmClient == nil {
		return nil
	}
	return this_.storage
}

/**
 * Get the lock instance.
 *
 * @return
 * - return NULL if error occurred
 */
func (this_ *IRtmClient) GetLock() *IRtmLock {
	// validity check
	if this_.lock == nil || this_.rtmClient == nil {
		return nil
	}
	return this_.lock
}

/**
 * Get the presence instance.
 *
 * @return
 * - return NULL if error occurred
 */
func (this_ *IRtmClient) GetPresence() *IRtmPresence {
	// validity check
	if this_.presence == nil || this_.rtmClient == nil {
		return nil
	}
	return this_.presence
}

/**
 * Get the history instance.
 *
 * @return
 * - return NULL if error occurred
 */
func (this_ *IRtmClient) GetHistory() *IRtmHistory {
	// validity check
	if this_.history == nil || this_.rtmClient == nil {
		return nil
	}
	return this_.history
}

/**
 * Renews the token. Once a token is enabled and used, it expires after a certain period of time.
 * You should generate a new token on your server, call this method to renew it.
 *
 * @param [in] token Token used renew.
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmClient) RenewToken(token string) (int, uint64) {
	// validity check
	if this_.rtmClient == nil || !this_.isLoggedIn {
		return -10006, 0
	}
	// do really renew token
	cToken := C.CString(token)
	defer C.free(unsafe.Pointer(cToken))
	var requestId uint64
	ret := int(C.agora_rtm_client_renew_token(this_.rtmClient,
		cToken,
		(*C.uint64_t)(unsafe.Pointer(&requestId)),
	))
	return int(ret), requestId
}

/**
 * Publish a message in the channel.
 *
 * @param [in] channelName The name of the channel.
 * @param [in] message The content of the message.
 * @param [in] length The length of the message.
 * @param [in] option The option of the message.
 * @param [out] requestId The related request id of this operation.
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmClient) Publish(channelName string, message []byte, option *PublishOptions) (int, uint64) {

	// validity check: only logged in can publish
	if this_.rtmClient == nil || !this_.isLoggedIn || message == nil || len(message) == 0 {
		return -10007, 0
	}

	// do really publish
	length := int(len(message))
	cChannelName := C.CString(channelName)
	defer C.free(unsafe.Pointer(cChannelName))
	cMessage := C.CBytes(message)
	defer C.free(unsafe.Pointer(cMessage))
	var cOption unsafe.Pointer
	if option != nil {
		cOption = option.toC()
		defer freePublishOptions(cOption)
	} else {
		cOption = NewPublishOptions().toC()
		defer freePublishOptions(cOption)
	}

	var requestId uint64
	ret := int(C.agora_rtm_client_publish(this_.rtmClient,
		cChannelName,
		(*C.char)(cMessage),
		C.size_t(length),
		(*C.struct_C_PublishOptions)(cOption),
		(*C.uint64_t)(unsafe.Pointer(&requestId)),
	))
	return ret, requestId
}

/**
 * Send a message to a channel.
 *
 * @param [in] channelName The name of the channel.
 * @param [in] message The content of the message.
 * @param [in] length The length of the message.
 * @param [out] requestId The related request id of this operation.
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmClient) SendChannelMessage(channelName string, message []byte) (int, uint64) {
	// validity check: only logged in can send channel message
	if this_.rtmClient == nil || !this_.isLoggedIn || message == nil || len(message) == 0 {
		return -10008, 0
	}

	// do really send channel message
	length := int(len(message))
	var requestId uint64
	cChannelName := C.CString(channelName)
	defer C.free(unsafe.Pointer(cChannelName))
	cMessage := C.CBytes(message)
	defer C.free(unsafe.Pointer(cMessage))
	opt := NewPublishOptions()
	opt.ChannelType = RtmChannelTypeMESSAGE
	opt.MessageType = RtmMessageTypeBINARY

	cOption := opt.toC()
	defer freePublishOptions(cOption)

	ret := int(C.agora_rtm_client_publish(this_.rtmClient,
		cChannelName,
		(*C.char)(cMessage),
		C.size_t(length),
		(*C.struct_C_PublishOptions)(cOption),
		(*C.uint64_t)(unsafe.Pointer(&requestId)),
	))
	return ret, requestId
}

/**
 * Send a message to a user.
 *
 * @param [in] userId The id of the user.
 * @param [in] message The content of the message.
 * @param [in] length The length of the message.
 * @param [out] requestId The related request id of this operation.
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmClient) SendUserMessage(userId string, message []byte) (int, uint64) {
	// validity check: only logged in can send user message
	if this_.rtmClient == nil || !this_.isLoggedIn || message == nil || len(message) == 0 {
		return -10009, 0
	}

	// do really send user message
	length := int(len(message))
	cUserId := C.CString(userId)
	defer C.free(unsafe.Pointer(cUserId))
	cMessage := C.CBytes(message)
	defer C.free(unsafe.Pointer(cMessage))
	opt := NewPublishOptions()
	opt.ChannelType = RtmChannelTypeUSER
	opt.MessageType = RtmMessageTypeBINARY

	cOption := opt.toC()
	defer freePublishOptions(cOption)
	var requestId uint64

	ret := int(C.agora_rtm_client_publish(this_.rtmClient,
		cUserId,
		(*C.char)(cMessage),
		C.size_t(length),
		(*C.struct_C_PublishOptions)(cOption),
		(*C.uint64_t)(unsafe.Pointer(&requestId)),
	))
	return ret, requestId
}

/**
 * Subscribe a channel.
 *
 * @param [in] channelName The name of the channel.
 * @param [in] options The options of subscribe the channel.
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmClient) Subscribe(channelName string, option *SubscribeOptions) (int, uint64) {
	// validity check: only logged in can subscribe
	if this_.rtmClient == nil || !this_.isLoggedIn || option == nil {
		return -10010, 0
	}

	// do really subscribe
	cChannelName := C.CString(channelName)
	defer C.free(unsafe.Pointer(cChannelName))
	cOption := option.toC()
	defer freeSubscribeOptions(cOption)
	var requestId uint64
	ret := int(C.agora_rtm_client_subscribe(this_.rtmClient,
		cChannelName,
		(*C.struct_C_SubscribeOptions)(cOption),
		(*C.uint64_t)(unsafe.Pointer(&requestId)),
	))
	return ret, requestId
}

/**
 * Unsubscribe a channel.
 *
 * @param [in] channelName The name of the channel.
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmClient) Unsubscribe(channelName string) (int, uint64) {
	// validity check: only logged in can unsubscribe
	if this_.rtmClient == nil || !this_.isLoggedIn {
		return -10011, 0
	}

	// do really unsubscribe
	cChannelName := C.CString(channelName)
	defer C.free(unsafe.Pointer(cChannelName))
	var requestId uint64
	ret := int(C.agora_rtm_client_unsubscribe(this_.rtmClient,
		cChannelName,
		(*C.uint64_t)(unsafe.Pointer(&requestId)),
	))
	return int(ret), requestId
}

/**
 * Create a stream channel instance.
 *
 * @param [in] channelName The Name of the channel.
 * @return
 * - return NULL if error occurred
 */
func (this_ *IRtmClient) CreateStreamChannel(channelName string) *IStreamChannel {

	// validity check: can't create stream channel if rtm client is not created
	if this_.rtmClient == nil {
		return nil
	}

	// do really create stream channel
	cChannelName := C.CString(channelName)
	defer C.free(unsafe.Pointer(cChannelName))

	var errorCode C.int
	ret := C.agora_rtm_client_create_stream_channel(this_.rtmClient,
		cChannelName,
		&errorCode,
	)

	if ret == nil || errorCode != 0 {
		return nil
	}

	streamChannel := &IStreamChannel{streamChannel: unsafe.Pointer(ret)}
	return streamChannel
}

/**
 * Set parameters of the sdk or engine
 *
 * @param [in] parameters The parameters in json format
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmClient) SetParameters(parameters string) int {

	// validity check: can't set parameters if rtm client is not created
	if this_.rtmClient == nil {
		return -10012
	}

	// do really set parameters
	cParameters := C.CString(parameters)
	defer C.free(unsafe.Pointer(cParameters))

	ret := int(C.agora_rtm_client_set_parameters(this_.rtmClient,
		cParameters,
	))
	return ret
}

// #endregion IRtmClient

/**
 * Convert error code to error string
 *
 * @param [in] errorCode Received error code
 * @return The error reason
 */
func GetErrorReason(errorCode int) string {
	return C.GoString(C.agora_rtm_client_get_error_reason(C.int(errorCode)))
}

/**
 * Get the version info of the Agora RTM SDK.
 *
 * @return The version info of the Agora RTM SDK.
 */
func GetVersion() string {
	return C.GoString(C.agora_rtm_client_get_version())
}

// #endregion agora::rtm

// #endregion agora
