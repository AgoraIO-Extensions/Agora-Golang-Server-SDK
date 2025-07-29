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
type RtmConfig C.struct_C_RtmConfig

// #region RtmConfig

/**
 * The App ID of your project.
 */
func (this_ *RtmConfig) GetAppId() string {
	return C.GoString(this_.appId)
}

/**
 * The App ID of your project.
 */
func (this_ *RtmConfig) SetAppId(appId string) {
	this_.appId = C.CString(appId)
}

/**
 * The ID of the user.
 */
func (this_ *RtmConfig) GetUserId() string {
	return C.GoString(this_.userId)
}

/**
 * The ID of the user.
 */
func (this_ *RtmConfig) SetUserId(userId string) {
	this_.userId = C.CString(userId)
}

/**
 * The region for connection. This advanced feature applies to scenarios that
 * have regional restrictions.
 *
 * For the regions that Agora supports, see #AREA_CODE.
 *
 * After specifying the region, the SDK connects to the Agora servers within
 * that region.
 */
func (this_ *RtmConfig) GetAreaCode() RTM_AREA_CODE {
	return RTM_AREA_CODE(this_.areaCode)
}

/**
 * The region for connection. This advanced feature applies to scenarios that
 * have regional restrictions.
 *
 * For the regions that Agora supports, see #AREA_CODE.
 *
 * After specifying the region, the SDK connects to the Agora servers within
 * that region.
 */
func (this_ *RtmConfig) SetAreaCode(areaCode RTM_AREA_CODE) {
	this_.areaCode = C.enum_C_RTM_AREA_CODE(areaCode)
}

/**
 * Presence timeout in seconds, specify the timeout value when you lost connection between sdk
 * and rtm service.
 */
func (this_ *RtmConfig) GetPresenceTimeout() uint32 {
	return uint32(this_.presenceTimeout)
}

/**
 * Presence timeout in seconds, specify the timeout value when you lost connection between sdk
 * and rtm service.
 */
func (this_ *RtmConfig) SetPresenceTimeout(presenceTimeout uint32) {
	this_.presenceTimeout = C.uint32_t(presenceTimeout)
}

/**
 * - For Android, it is the context of Activity or Application.
 * - For Windows, it is the window handle of app. Once set, this parameter enables you to plug
 * or unplug the video devices while they are powered.
 */
func (this_ *RtmConfig) GetContext() unsafe.Pointer {
	return this_.context
}

/**
 * - For Android, it is the context of Activity or Application.
 * - For Windows, it is the window handle of app. Once set, this parameter enables you to plug
 * or unplug the video devices while they are powered.
 */
func (this_ *RtmConfig) SetContext(context unsafe.Pointer) {
	this_.context = context
}

/**
 * Whether to use String user IDs, if you are using RTC products with Int user IDs,
 * set this value as 'false'. Otherwise errors might occur.
 */
func (this_ *RtmConfig) GetUseStringUserId() bool {
	return bool(this_.useStringUserId)
}

/**
 * Whether to use String user IDs, if you are using RTC products with Int user IDs,
 * set this value as 'false'. Otherwise errors might occur.
 */
func (this_ *RtmConfig) SetUseStringUserId(useStringUserId bool) {
	this_.useStringUserId = C.bool(useStringUserId)
}

/**
 * The callbacks handler
 */
func (this_ *RtmConfig) GetEventHandler() *IRtmEventHandler {
	return (*IRtmEventHandler)(this_.eventHandler)
}

/**
 * The callbacks handler
 */
func (this_ *RtmConfig) SetEventHandler(eventHandler *IRtmEventHandler) {
	this_.eventHandler = unsafe.Pointer(eventHandler)
}

/**
 * The config for customer set log path, log size and log level.
 */
func (this_ *RtmConfig) GetLogConfig() RtmLogConfig {
	return (RtmLogConfig)(this_.logConfig)
}

/**
 * The config for customer set log path, log size and log level.
 */
func (this_ *RtmConfig) SetLogConfig(logConfig RtmLogConfig) {
	this_.logConfig = (C.struct_C_RtmLogConfig)(logConfig)
}

/**
 * The config for proxy setting
 */
func (this_ *RtmConfig) GetProxyConfig() RtmProxyConfig {
	return (RtmProxyConfig)(this_.proxyConfig)
}

/**
 * The config for proxy setting
 */
func (this_ *RtmConfig) SetProxyConfig(proxyConfig RtmProxyConfig) {
	this_.proxyConfig = (C.struct_C_RtmProxyConfig)(proxyConfig)
}

/**
 * The config for encryption setting
 */
func (this_ *RtmConfig) GetEncryptionConfig() RtmEncryptionConfig {
	return (RtmEncryptionConfig)(this_.encryptionConfig)
}

/**
 * The config for encryption setting
 */
func (this_ *RtmConfig) SetEncryptionConfig(encryptionConfig RtmEncryptionConfig) {
	this_.encryptionConfig = (C.struct_C_RtmEncryptionConfig)(encryptionConfig)
}

func NewRtmConfig() *RtmConfig {
	return (*RtmConfig)(C.C_RtmConfig_New())
}
func (this_ *RtmConfig) Delete() {
	C.C_RtmConfig_Delete((*C.struct_C_RtmConfig)(this_))
}

// #endregion RtmConfig

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
type IRtmEventHandler C.C_IRtmEventHandler

// #region IRtmEventHandler

type MessageEvent C.struct_C_MessageEvent

// #region MessageEvent

/**
 * Which channel type, RTM_CHANNEL_TYPE_STREAM or RTM_CHANNEL_TYPE_MESSAGE
 */
func (this_ *MessageEvent) GetChannelType() RTM_CHANNEL_TYPE {
	return RTM_CHANNEL_TYPE(this_.channelType)
}

/**
 * Which channel type, RTM_CHANNEL_TYPE_STREAM or RTM_CHANNEL_TYPE_MESSAGE
 */
func (this_ *MessageEvent) SetChannelType(channelType RTM_CHANNEL_TYPE) {
	this_.channelType = C.enum_C_RTM_CHANNEL_TYPE(channelType)
}

/**
 * Message type
 */
func (this_ *MessageEvent) GetMessageType() RTM_MESSAGE_TYPE {
	return RTM_MESSAGE_TYPE(this_.messageType)
}

/**
 * Message type
 */
func (this_ *MessageEvent) SetMessageType(messageType RTM_MESSAGE_TYPE) {
	this_.messageType = C.enum_C_RTM_MESSAGE_TYPE(messageType)
}

/**
 * The channel which the message was published
 */
func (this_ *MessageEvent) GetChannelName() string {
	return C.GoString(this_.channelName)
}

/**
 * The channel which the message was published
 */
func (this_ *MessageEvent) SetChannelName(channelName string) {
	this_.channelName = C.CString(channelName)
}

/**
 * If the channelType is RTM_CHANNEL_TYPE_STREAM, which topic the message came from. only for RTM_CHANNEL_TYPE_STREAM
 */
func (this_ *MessageEvent) GetChannelTopic() string {
	return C.GoString(this_.channelTopic)
}

/**
 * If the channelType is RTM_CHANNEL_TYPE_STREAM, which topic the message came from. only for RTM_CHANNEL_TYPE_STREAM
 */
func (this_ *MessageEvent) SetChannelTopic(channelTopic string) {
	this_.channelTopic = C.CString(channelTopic)
}

/**
 * The payload
 */
func (this_ *MessageEvent) GetMessage() []byte {
	return C.GoBytes(
		unsafe.Pointer(this_.message),
		C.int(this_.GetMessageLength()),
	)
}

/**
 * The payload
 */
func (this_ *MessageEvent) SetMessage(message []byte) {
	this_.message = (*C.char)(C.CBytes(message))
}

/**
 * The payload length
 */
func (this_ *MessageEvent) GetMessageLength() uint {
	return uint(this_.messageLength)
}

/**
 * The payload length
 */
func (this_ *MessageEvent) SetMessageLength(messageLength uint) {
	this_.messageLength = C.size_t(messageLength)
}

/**
 * The publisher
 */
func (this_ *MessageEvent) GetPublisher() string {
	return C.GoString(this_.publisher)
}

/**
 * The publisher
 */
func (this_ *MessageEvent) SetPublisher(publisher string) {
	this_.publisher = C.CString(publisher)
}

/**
 * The custom type of the message
 */
func (this_ *MessageEvent) GetCustomType() string {
	return C.GoString(this_.customType)
}

/**
 * The publisher
 */
func (this_ *MessageEvent) SetCustomType(customType string) {
	this_.customType = C.CString(customType)
}

func NewMessageEvent() *MessageEvent {
	return (*MessageEvent)(C.C_MessageEvent_New())
}
func (this_ *MessageEvent) Delete() {
	C.C_MessageEvent_Delete((*C.struct_C_MessageEvent)(this_))
}

// #endregion MessageEvent

type IntervalInfo C.struct_C_IntervalInfo

// #region IntervalInfo

/**
 * Joined users during this interval
 */
func (this_ *IntervalInfo) GetJoinUserList() UserList {
	return (UserList)(this_.joinUserList)
}

/**
 * Joined users during this interval
 */
func (this_ *IntervalInfo) SetJoinUserList(joinUserList UserList) {
	this_.joinUserList = C.struct_C_UserList(joinUserList)
}

/**
 * Left users during this interval
 */
func (this_ *IntervalInfo) GetLeaveUserList() UserList {
	return (UserList)(this_.leaveUserList)
}

/**
 * Left users during this interval
 */
func (this_ *IntervalInfo) SetLeaveUserList(leaveUserList UserList) {
	this_.leaveUserList = C.struct_C_UserList(leaveUserList)
}

/**
 * Timeout users during this interval
 */
func (this_ *IntervalInfo) GetTimeoutUserList() UserList {
	return (UserList)(this_.timeoutUserList)
}

/**
 * Timeout users during this interval
 */
func (this_ *IntervalInfo) SetTimeoutUserList(timeoutUserList UserList) {
	this_.timeoutUserList = C.struct_C_UserList(timeoutUserList)
}

/**
 * The user state changed during this interval
 */
func (this_ *IntervalInfo) GetUserStateList() []UserState {
	count := this_.GetUserStateCount()
	return unsafe.Slice((*UserState)(this_.userStateList), count)
}

/**
 * The user state changed during this interval
 */
func (this_ *IntervalInfo) SetUserStateList(userStateList []UserState) {
	this_.userStateList = (*C.struct_C_UserState)(unsafe.SliceData(userStateList))
}

/**
 * The user count
 */
func (this_ *IntervalInfo) GetUserStateCount() uint {
	return uint(this_.userStateCount)
}

/**
 * The user count
 */
func (this_ *IntervalInfo) SetUserStateCount(userStateCount uint) {
	this_.userStateCount = C.size_t(userStateCount)
}

func NewIntervalInfo() *IntervalInfo {
	return (*IntervalInfo)(C.C_IntervalInfo_New())
}
func (this_ *IntervalInfo) Delete() {
	C.C_IntervalInfo_Delete((*C.struct_C_IntervalInfo)(this_))
}

// #endregion IntervalInfo

type SnapshotInfo C.struct_C_SnapshotInfo

// #region SnapshotInfo

/**
 * The user state in this snapshot event
 */
func (this_ *SnapshotInfo) GetUserStateList() []UserState {
	count := this_.GetUserCount()
	return unsafe.Slice((*UserState)(this_.userStateList), count)
}

/**
 * The user state in this snapshot event
 */
func (this_ *SnapshotInfo) SetUserStateList(userStateList []UserState) {
	this_.userStateList = (*C.struct_C_UserState)(unsafe.SliceData(userStateList))
}

/**
 * The user count
 */
func (this_ *SnapshotInfo) GetUserCount() uint {
	return uint(this_.userCount)
}

/**
 * The user count
 */
func (this_ *SnapshotInfo) SetUserCount(userCount uint) {
	this_.userCount = C.size_t(userCount)
}

func NewSnapshotInfo() *SnapshotInfo {
	return (*SnapshotInfo)(C.C_SnapshotInfo_New())
}
func (this_ *SnapshotInfo) Delete() {
	C.C_SnapshotInfo_Delete((*C.struct_C_SnapshotInfo)(this_))
}

// #endregion SnapshotInfo

type PresenceEvent C.struct_C_PresenceEvent

// #region PresenceEvent

/**
 * Indicate presence event type
 */
func (this_ *PresenceEvent) GetType() RTM_PRESENCE_EVENT_TYPE {
	return RTM_PRESENCE_EVENT_TYPE(this_._type)
}

/**
 * Indicate presence event type
 */
func (this_ *PresenceEvent) SetType(_type RTM_PRESENCE_EVENT_TYPE) {
	this_._type = C.enum_C_RTM_PRESENCE_EVENT_TYPE(_type)
}

/**
 * Which channel type, RTM_CHANNEL_TYPE_STREAM or RTM_CHANNEL_TYPE_MESSAGE
 */
func (this_ *PresenceEvent) GetChannelType() RTM_CHANNEL_TYPE {
	return RTM_CHANNEL_TYPE(this_.channelType)
}

/**
 * Which channel type, RTM_CHANNEL_TYPE_STREAM or RTM_CHANNEL_TYPE_MESSAGE
 */
func (this_ *PresenceEvent) SetChannelType(_type RTM_CHANNEL_TYPE) {
	this_._type = C.enum_C_RTM_CHANNEL_TYPE(_type)
}

/**
 * The channel which the presence event was triggered
 */
func (this_ *PresenceEvent) GetChannelName() string {
	return C.GoString(this_.channelName)
}

/**
 * The channel which the presence event was triggered
 */
func (this_ *PresenceEvent) SetChannelName(channelName string) {
	this_.channelName = C.CString(channelName)
}

/**
 * The user who triggered this event.
 */
func (this_ *PresenceEvent) GetPublisher() string {
	return C.GoString(this_.publisher)
}

/**
 * The user who triggered this event.
 */
func (this_ *PresenceEvent) SetPublisher(publisher string) {
	this_.publisher = C.CString(publisher)
}

/**
 * The user states
 */
func (this_ *PresenceEvent) GetStateItems() []StateItem {
	count := this_.GetStateItemCount()
	return unsafe.Slice((*StateItem)(this_.stateItems), count)
}

/**
 * The user states
 */
func (this_ *PresenceEvent) SetStateItems(stateItems []StateItem) {
	this_.stateItems = (*C.struct_C_StateItem)(unsafe.SliceData(stateItems))
}

/**
 * The states count
 */
func (this_ *PresenceEvent) GetStateItemCount() uint {
	return uint(this_.stateItemCount)
}

/**
 * The states count
 */
func (this_ *PresenceEvent) SetStateItemCount(stateItemCount uint) {
	this_.stateItemCount = C.size_t(stateItemCount)
}

/**
 * Only valid when in interval mode
 */
func (this_ *PresenceEvent) GetInterval() IntervalInfo {
	return IntervalInfo(this_.interval)
}

/**
 * Only valid when in interval mode
 */
func (this_ *PresenceEvent) SetInterval(interval IntervalInfo) {
	this_.interval = (C.struct_C_IntervalInfo)(interval)
}

/**
 * Only valid when receive snapshot event
 */
func (this_ *PresenceEvent) GetSnapshot() SnapshotInfo {
	return SnapshotInfo(this_.snapshot)
}

/**
 * Only valid when receive snapshot event
 */
func (this_ *PresenceEvent) SetSnapshot(snapshot SnapshotInfo) {
	this_.snapshot = (C.struct_C_SnapshotInfo)(snapshot)
}

func NewPresenceEvent() *PresenceEvent {
	return (*PresenceEvent)(C.C_PresenceEvent_New())
}
func (this_ *PresenceEvent) Delete() {
	C.C_PresenceEvent_Delete((*C.struct_C_PresenceEvent)(this_))
}

// #endregion PresenceEvent

type TopicEvent C.struct_C_TopicEvent

// #region TopicEvent

/**
 * Indicate topic event type
 */
func (this_ *TopicEvent) GetType() RTM_TOPIC_EVENT_TYPE {
	return RTM_TOPIC_EVENT_TYPE(this_._type)
}

/**
 * Indicate topic event type
 */
func (this_ *TopicEvent) SetType(_type RTM_TOPIC_EVENT_TYPE) {
	this_._type = C.enum_C_RTM_TOPIC_EVENT_TYPE(_type)
}

/**
 * The channel which the topic event was triggered
 */
func (this_ *TopicEvent) GetChannelName() string {
	return C.GoString(this_.channelName)
}

/**
 * The channel which the topic event was triggered
 */
func (this_ *TopicEvent) SetChannelName(channelName string) {
	this_.channelName = C.CString(channelName)
}

/**
 * The user who triggered this event.
 */
func (this_ *TopicEvent) GetPublisher() string {
	return C.GoString(this_.publisher)
}

/**
 * The user who triggered this event.
 */
func (this_ *TopicEvent) SetPublisher(publisher string) {
	this_.publisher = C.CString(publisher)
}

/**
 * Topic information array.
 */
func (this_ *TopicEvent) GetTopicInfos() []TopicInfo {
	count := this_.GetTopicInfoCount()
	return unsafe.Slice((*TopicInfo)(this_.topicInfos), count)
}

/**
 * Topic information array.
 */
func (this_ *TopicEvent) SetTopicInfos(topicInfos []TopicInfo) {
	this_.topicInfos = (*C.struct_C_TopicInfo)(unsafe.SliceData(topicInfos))
}

/**
 * The count of topicInfos.
 */
func (this_ *TopicEvent) GetTopicInfoCount() uint {
	return uint(this_.topicInfoCount)
}

/**
 * The count of topicInfos.
 */
func (this_ *TopicEvent) SetTopicInfoCount(topicInfoCount uint) {
	this_.topicInfoCount = C.size_t(topicInfoCount)
}

func NewTopicEvent() *TopicEvent {
	return (*TopicEvent)(C.C_TopicEvent_New())
}
func (this_ *TopicEvent) Delete() {
	C.C_TopicEvent_Delete((*C.struct_C_TopicEvent)(this_))
}

// #endregion TopicEvent

type LockEvent C.struct_C_LockEvent

// #region LockEvent

/**
 * Which channel type, RTM_CHANNEL_TYPE_STREAM or RTM_CHANNEL_TYPE_MESSAGE
 */
func (this_ *LockEvent) GetChannelType() RTM_CHANNEL_TYPE {
	return RTM_CHANNEL_TYPE(this_.channelType)
}

/**
 * Which channel type, RTM_CHANNEL_TYPE_STREAM or RTM_CHANNEL_TYPE_MESSAGE
 */
func (this_ *LockEvent) SetChannelType(channelType RTM_CHANNEL_TYPE) {
	this_.channelType = C.enum_C_RTM_CHANNEL_TYPE(channelType)
}

/**
 * Lock event type, indicate lock states
 */
func (this_ *LockEvent) GetEventType() RTM_LOCK_EVENT_TYPE {
	return RTM_LOCK_EVENT_TYPE(this_.eventType)
}

/**
 * Lock event type, indicate lock states
 */
func (this_ *LockEvent) SetEventType(eventType RTM_LOCK_EVENT_TYPE) {
	this_.eventType = C.enum_C_RTM_LOCK_EVENT_TYPE(eventType)
}

/**
 * The channel which the lock event was triggered
 */
func (this_ *LockEvent) GetChannelName() string {
	return C.GoString(this_.channelName)
}

/**
 * The channel which the lock event was triggered
 */
func (this_ *LockEvent) SetChannelName(channelName string) {
	this_.channelName = C.CString(channelName)
}

/**
 * The detail information of locks
 */
func (this_ *LockEvent) GetLockDetailList() []LockDetail {
	count := this_.GetCount()
	return unsafe.Slice((*LockDetail)(this_.lockDetailList), count)
}

/**
 * The detail information of locks
 */
func (this_ *LockEvent) SetLockDetailList(lockDetailList []LockDetail) {
	this_.lockDetailList = (*C.struct_C_LockDetail)(unsafe.SliceData(lockDetailList))
}

/**
 * The count of locks
 */
func (this_ *LockEvent) GetCount() uint {
	return uint(this_.count)
}

/**
 * The count of locks
 */
func (this_ *LockEvent) SetCount(count uint) {
	this_.count = C.size_t(count)
}

func NewLockEvent() *LockEvent {
	return (*LockEvent)(C.C_LockEvent_New())
}
func (this_ *LockEvent) Delete() {
	C.C_LockEvent_Delete((*C.struct_C_LockEvent)(this_))
}

// #endregion LockEvent

type StorageEvent C.struct_C_StorageEvent

// #region StorageEvent

/**
 * Which channel type, RTM_CHANNEL_TYPE_STREAM or RTM_CHANNEL_TYPE_MESSAGE
 */
func (this_ *StorageEvent) GetChannelType() RTM_CHANNEL_TYPE {
	return RTM_CHANNEL_TYPE(this_.channelType)
}

/**
 * Which channel type, RTM_CHANNEL_TYPE_STREAM or RTM_CHANNEL_TYPE_MESSAGE
 */
func (this_ *StorageEvent) SetChannelType(channelType RTM_CHANNEL_TYPE) {
	this_.channelType = C.enum_C_RTM_CHANNEL_TYPE(channelType)
}

/**
 * Storage type, RTM_STORAGE_TYPE_USER or RTM_STORAGE_TYPE_CHANNEL
 */
func (this_ *StorageEvent) GetStorageType() RTM_STORAGE_TYPE {
	return RTM_STORAGE_TYPE(this_.storageType)
}

/**
 * Storage type, RTM_STORAGE_TYPE_USER or RTM_STORAGE_TYPE_CHANNEL
 */
func (this_ *StorageEvent) SetStorageType(storageType RTM_STORAGE_TYPE) {
	this_.storageType = C.enum_C_RTM_STORAGE_TYPE(storageType)
}

/**
 * Indicate storage event type
 */
func (this_ *StorageEvent) GetEventType() RTM_STORAGE_EVENT_TYPE {
	return RTM_STORAGE_EVENT_TYPE(this_.eventType)
}

/**
 * Indicate storage event type
 */
func (this_ *StorageEvent) SetEventType(eventType RTM_STORAGE_EVENT_TYPE) {
	this_.eventType = C.enum_C_RTM_STORAGE_EVENT_TYPE(eventType)
}

/**
 * The target name of user or channel, depends on the RTM_STORAGE_TYPE
 */
func (this_ *StorageEvent) GetTarget() string {
	return C.GoString(this_.target)
}

/**
 * The target name of user or channel, depends on the RTM_STORAGE_TYPE
 */
func (this_ *StorageEvent) SetTarget(target string) {
	this_.target = C.CString(target)
}

/**
 * The metadata information
 */
func (this_ *StorageEvent) GetData() *IMetadata {
	return (*IMetadata)(this_.data)
}

/**
 * The metadata information
 */
func (this_ *StorageEvent) SetData(data *IMetadata) {
	this_.data = (*C.struct_C_Metadata)(unsafe.Pointer(data))
}

func NewStorageEvent() *StorageEvent {
	return (*StorageEvent)(C.C_StorageEvent_New())
}
func (this_ *StorageEvent) Delete() {
	C.C_StorageEvent_Delete((*C.struct_C_StorageEvent)(this_))
}

// #endregion StorageEvent

/**
 * Occurs when receive a message.
 *
 * @param event details of message event.
 */
func (this_ *IRtmEventHandler) OnMessageEvent(event *MessageEvent) {
	C.C_IRtmEventHandler_onMessageEvent(unsafe.Pointer(this_), (*C.struct_C_MessageEvent)(event))
}

/**
 * Occurs when remote user presence changed
 *
 * @param event details of presence event.
 */
func (this_ *IRtmEventHandler) OnPresenceEvent(event *PresenceEvent) {
	C.C_IRtmEventHandler_onPresenceEvent(unsafe.Pointer(this_), (*C.struct_C_PresenceEvent)(event))
}

/**
 * Occurs when remote user join/leave topic or when user first join this channel,
 * got snapshot of topics in this channel
 *
 * @param event details of topic event.
 */
func (this_ *IRtmEventHandler) OnTopicEvent(event *TopicEvent) {
	C.C_IRtmEventHandler_onTopicEvent(unsafe.Pointer(this_), (*C.struct_C_TopicEvent)(event))
}

/**
 * Occurs when lock state changed
 *
 * @param event details of lock event.
 */
func (this_ *IRtmEventHandler) OnLockEvent(event *LockEvent) {
	C.C_IRtmEventHandler_onLockEvent(unsafe.Pointer(this_), (*C.struct_C_LockEvent)(event))
}

/**
 * Occurs when receive storage event
 *
 * @param event details of storage event.
 */
func (this_ *IRtmEventHandler) OnStorageEvent(event *StorageEvent) {
	C.C_IRtmEventHandler_onStorageEvent(unsafe.Pointer(this_), (*C.struct_C_StorageEvent)(event))
}

/**
 * Occurs when user join a stream channel.
 *
 * @param channelName The name of the channel.
 * @param userId The id of the user.
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnJoinResult(requestId uint64, channelName string, userId string, errorCode RTM_ERROR_CODE) {
	cChannelName := C.CString(channelName)
	cUserId := C.CString(userId)
	C.C_IRtmEventHandler_onJoinResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		cChannelName, cUserId,
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
	C.free(unsafe.Pointer(cChannelName))
	C.free(unsafe.Pointer(cUserId))
}

/**
 * Occurs when user leave a stream channel.
 *
 * @param channelName The name of the channel.
 * @param userId The id of the user.
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnLeaveResult(requestId uint64, channelName string, userId string, errorCode RTM_ERROR_CODE) {
	cChannelName := C.CString(channelName)
	cUserId := C.CString(userId)
	C.C_IRtmEventHandler_onLeaveResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		cChannelName, cUserId,
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
	C.free(unsafe.Pointer(cChannelName))
	C.free(unsafe.Pointer(cUserId))
}

/**
 * Occurs when user join topic.
 *
 * @param channelName The name of the channel.
 * @param userId The id of the user.
 * @param topic The name of the topic.
 * @param meta The meta of the topic.
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnJoinTopicResult(requestId uint64, channelName string, userId string, topic string, meta string, errorCode RTM_ERROR_CODE) {
	cChannelName := C.CString(channelName)
	cUserId := C.CString(userId)
	cTopic := C.CString(topic)
	cMeta := C.CString(meta)
	C.C_IRtmEventHandler_onJoinTopicResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		cChannelName, cUserId,
		cTopic,
		cMeta,
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
	C.free(unsafe.Pointer(cChannelName))
	C.free(unsafe.Pointer(cUserId))
	C.free(unsafe.Pointer(cTopic))
	C.free(unsafe.Pointer(cMeta))
}

/**
 * Occurs when user leave topic.
 *
 * @param channelName The name of the channel.
 * @param userId The id of the user.
 * @param topic The name of the topic.
 * @param meta The meta of the topic.
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnLeaveTopicResult(requestId uint64, channelName string, userId string, topic string, meta string, errorCode RTM_ERROR_CODE) {
	cChannelName := C.CString(channelName)
	cUserId := C.CString(userId)
	cTopic := C.CString(topic)
	cMeta := C.CString(meta)
	C.C_IRtmEventHandler_onLeaveTopicResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		cChannelName,
		cUserId,
		cTopic,
		cMeta,
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
	C.free(unsafe.Pointer(cChannelName))
	C.free(unsafe.Pointer(cUserId))
	C.free(unsafe.Pointer(cTopic))
	C.free(unsafe.Pointer(cMeta))
}

/**
 * Occurs when user subscribe topic.
 *
 * @param channelName The name of the channel.
 * @param userId The id of the user.
 * @param topic The name of the topic.
 * @param succeedUsers The subscribed users.
 * @param failedUser The failed to subscribe users.
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnSubscribeTopicResult(requestId uint64, channelName string, userId string, topic string, succeedUsers UserList, failedUsers UserList, errorCode RTM_ERROR_CODE) {
	cChannelName := C.CString(channelName)
	cUserId := C.CString(userId)
	cTopic := C.CString(topic)
	C.C_IRtmEventHandler_onSubscribeTopicResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		cChannelName,
		cUserId,
		cTopic,
		C.struct_C_UserList(succeedUsers),
		C.struct_C_UserList(failedUsers),
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
	C.free(unsafe.Pointer(cChannelName))
	C.free(unsafe.Pointer(cUserId))
	C.free(unsafe.Pointer(cTopic))
}

/**
 * Occurs when the connection state changes between rtm sdk and agora service.
 *
 * @param channelName The name of the channel.
 * @param state The new connection state.
 * @param reason The reason for the connection state change.
 */
func (this_ *IRtmEventHandler) OnConnectionStateChanged(channelName string, state RTM_CONNECTION_STATE, reason RTM_CONNECTION_CHANGE_REASON) {
	cChannelName := C.CString(channelName)
	C.C_IRtmEventHandler_onConnectionStateChanged(unsafe.Pointer(this_),
		cChannelName,
		C.enum_C_RTM_CONNECTION_STATE(state),
		C.enum_C_RTM_CONNECTION_CHANGE_REASON(reason),
	)
	C.free(unsafe.Pointer(cChannelName))
}

/**
 * Occurs when token will expire in 30 seconds.
 *
 * @param channelName The name of the channel.
 */
func (this_ *IRtmEventHandler) OnTokenPrivilegeWillExpire(channelName string) {
	cChannelName := C.CString(channelName)
	C.C_IRtmEventHandler_onTokenPrivilegeWillExpire(unsafe.Pointer(this_),
		cChannelName,
	)
	C.free(unsafe.Pointer(cChannelName))
}

/**
 * Occurs when subscribe a channel
 *
 * @param channelName The name of the channel.
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnSubscribeResult(requestId uint64, channelName string, errorCode RTM_ERROR_CODE) {
	cChannelName := C.CString(channelName)
	C.C_IRtmEventHandler_onSubscribeResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		cChannelName,
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
	C.free(unsafe.Pointer(cChannelName))
}

/**
 * Occurs when user publish message.
 *
 * @param requestId The related request id when user publish message
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnPublishResult(requestId uint64, errorCode RTM_ERROR_CODE) {
	C.C_IRtmEventHandler_onPublishResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
}

/**
 * Occurs when user login.
 *
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnLoginResult(requestId uint64, errorCode RTM_ERROR_CODE) {

	C.C_IRtmEventHandler_onLoginResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
}

/**
 * Occurs when user setting the channel metadata
 *
 * @param requestId The related request id when user perform this operation
 * @param channelName The name of the channel.
 * @param channelType The type of the channel.
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnSetChannelMetadataResult(requestId uint64, channelName string, channelType RTM_CHANNEL_TYPE, errorCode RTM_ERROR_CODE) {
	cChannelName := C.CString(channelName)
	C.C_IRtmEventHandler_onSetChannelMetadataResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
	C.free(unsafe.Pointer(cChannelName))
}

/**
 * Occurs when user updating the channel metadata
 *
 * @param requestId The related request id when user perform this operation
 * @param channelName The name of the channel.
 * @param channelType The type of the channel.
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnUpdateChannelMetadataResult(requestId uint64, channelName string, channelType RTM_CHANNEL_TYPE, errorCode RTM_ERROR_CODE) {
	cChannelName := C.CString(channelName)
	C.C_IRtmEventHandler_onUpdateChannelMetadataResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
	C.free(unsafe.Pointer(cChannelName))
}

/**
 * Occurs when user removing the channel metadata
 *
 * @param requestId The related request id when user perform this operation
 * @param channelName The name of the channel.
 * @param channelType The type of the channel.
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnRemoveChannelMetadataResult(requestId uint64, channelName string, channelType RTM_CHANNEL_TYPE, errorCode RTM_ERROR_CODE) {
	cChannelName := C.CString(channelName)
	C.C_IRtmEventHandler_onRemoveChannelMetadataResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
	C.free(unsafe.Pointer(cChannelName))
}

/**
 * Occurs when user try to get the channel metadata
 *
 * @param requestId The related request id when user perform this operation
 * @param channelName The name of the channel.
 * @param channelType The type of the channel.
 * @param data The result metadata of getting operation.
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnGetChannelMetadataResult(requestId uint64, channelName string, channelType RTM_CHANNEL_TYPE, data *IMetadata, errorCode RTM_ERROR_CODE) {
	cChannelName := C.CString(channelName)
	C.C_IRtmEventHandler_onGetChannelMetadataResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		(*C.struct_C_Metadata)(unsafe.Pointer(data)),
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
	C.free(unsafe.Pointer(cChannelName))
}

/**
 * Occurs when user setting the user metadata
 *
 * @param requestId The related request id when user perform this operation
 * @param userId The id of the user.
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnSetUserMetadataResult(requestId uint64, userId string, errorCode RTM_ERROR_CODE) {
	cUserId := C.CString(userId)
	C.C_IRtmEventHandler_onSetUserMetadataResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		cUserId,
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
	C.free(unsafe.Pointer(cUserId))
}

/**
 * Occurs when user updating the user metadata
 *
 * @param requestId The related request id when user perform this operation
 * @param userId The id of the user.
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnUpdateUserMetadataResult(requestId uint64, userId string, errorCode RTM_ERROR_CODE) {
	cUserId := C.CString(userId)
	C.C_IRtmEventHandler_onUpdateUserMetadataResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		cUserId,
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
	C.free(unsafe.Pointer(cUserId))
}

/**
 * Occurs when user removing the user metadata
 *
 * @param requestId The related request id when user perform this operation
 * @param userId The id of the user.
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnRemoveUserMetadataResult(requestId uint64, userId string, errorCode RTM_ERROR_CODE) {
	cUserId := C.CString(userId)
	C.C_IRtmEventHandler_onRemoveUserMetadataResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		cUserId,
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
	C.free(unsafe.Pointer(cUserId))
}

/**
 * Occurs when user try to get the user metadata
 *
 * @param requestId The related request id when user perform this operation
 * @param userId The id of the user.
 * @param data The result metadata of getting operation.
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnGetUserMetadataResult(requestId uint64, userId string, data *IMetadata, errorCode RTM_ERROR_CODE) {
	cUserId := C.CString(userId)
	C.C_IRtmEventHandler_onGetUserMetadataResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		cUserId,
		(*C.struct_C_Metadata)(unsafe.Pointer(data)),
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
	C.free(unsafe.Pointer(cUserId))
}

/**
 * Occurs when user subscribe a user metadata
 *
 * @param userId The id of the user.
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnSubscribeUserMetadataResult(requestId uint64, userId string, errorCode RTM_ERROR_CODE) {
	cUserId := C.CString(userId)
	C.C_IRtmEventHandler_onSubscribeUserMetadataResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		cUserId,
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
	C.free(unsafe.Pointer(cUserId))
}

/**
 * Occurs when user set a lock
 *
 * @param channelName The name of the channel.
 * @param channelType The type of the channel.
 * @param lockName The name of the lock.
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnSetLockResult(requestId uint64, channelName string, channelType RTM_CHANNEL_TYPE, lockName string, errorCode RTM_ERROR_CODE) {
	cChannelName := C.CString(channelName)
	cLockName := C.CString(lockName)
	C.C_IRtmEventHandler_onSetLockResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		cLockName,
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
	C.free(unsafe.Pointer(cChannelName))
	C.free(unsafe.Pointer(cLockName))
}

/**
 * Occurs when user delete a lock
 *
 * @param channelName The name of the channel.
 * @param channelType The type of the channel.
 * @param lockName The name of the lock.
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnRemoveLockResult(requestId uint64, channelName string, channelType RTM_CHANNEL_TYPE, lockName string, errorCode RTM_ERROR_CODE) {
	cChannelName := C.CString(channelName)
	cLockName := C.CString(lockName)
	C.C_IRtmEventHandler_onRemoveLockResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		cLockName,
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
	C.free(unsafe.Pointer(cChannelName))
	C.free(unsafe.Pointer(cLockName))
}

/**
 * Occurs when user release a lock
 *
 * @param channelName The name of the channel.
 * @param channelType The type of the channel.
 * @param lockName The name of the lock.
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnReleaseLockResult(requestId uint64, channelName string, channelType RTM_CHANNEL_TYPE, lockName string, errorCode RTM_ERROR_CODE) {
	cChannelName := C.CString(channelName)
	cLockName := C.CString(lockName)
	C.C_IRtmEventHandler_onReleaseLockResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		cLockName,
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
	C.free(unsafe.Pointer(cChannelName))
	C.free(unsafe.Pointer(cLockName))
}

/**
 * Occurs when user acquire a lock
 *
 * @param channelName The name of the channel.
 * @param channelType The type of the channel.
 * @param lockName The name of the lock.
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnAcquireLockResult(requestId uint64, channelName string, channelType RTM_CHANNEL_TYPE, lockName string, errorCode RTM_ERROR_CODE, errorDetails string) {
	cChannelName := C.CString(channelName)
	cLockName := C.CString(lockName)
	cErrorDetails := C.CString(errorDetails)
	C.C_IRtmEventHandler_onAcquireLockResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		cLockName,
		C.enum_C_RTM_ERROR_CODE(errorCode),
		cErrorDetails,
	)
	C.free(unsafe.Pointer(cChannelName))
	C.free(unsafe.Pointer(cLockName))
	C.free(unsafe.Pointer(cErrorDetails))
}

/**
 * Occurs when user revoke a lock
 *
 * @param channelName The name of the channel.
 * @param channelType The type of the channel.
 * @param lockName The name of the lock.
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnRevokeLockResult(requestId uint64, channelName string, channelType RTM_CHANNEL_TYPE, lockName string, errorCode RTM_ERROR_CODE) {
	cChannelName := C.CString(channelName)
	cLockName := C.CString(lockName)
	C.C_IRtmEventHandler_onRevokeLockResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		cLockName,
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
	C.free(unsafe.Pointer(cChannelName))
	C.free(unsafe.Pointer(cLockName))
}

/**
 * Occurs when user try to get locks from the channel
 *
 * @param channelName The name of the channel.
 * @param channelType The type of the channel.
 * @param lockDetailList The details of the locks.
 * @param count The count of the locks.
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnGetLocksResult(requestId uint64, channelName string, channelType RTM_CHANNEL_TYPE, lockDetailList []LockDetail, count uint, errorCode RTM_ERROR_CODE) {
	cChannelName := C.CString(channelName)
	C.C_IRtmEventHandler_onGetLocksResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		(*C.struct_C_LockDetail)(unsafe.SliceData(lockDetailList)),
		C.size_t(count),
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
	C.free(unsafe.Pointer(cChannelName))
}

/**
 * Occurs when query who joined this channel
 *
 * @param requestId The related request id when user perform this operation
 * @param userStatesList The states the users.
 * @param count The user count.
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnWhoNowResult(requestId uint64, userStateList []UserState, count uint, nextPage string, errorCode RTM_ERROR_CODE) {
	cNextPage := C.CString(nextPage)
	C.C_IRtmEventHandler_onWhoNowResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		(*C.struct_C_UserState)(unsafe.SliceData(userStateList)),
		C.size_t(count),
		cNextPage,
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
	C.free(unsafe.Pointer(cNextPage))
}

/**
 * Occurs when query who joined this channel
 *
 * @param requestId The related request id when user perform this operation
 * @param userStatesList The states the users.
 * @param count The user count.
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnGetOnlineUsersResult(requestId uint64, userStateList []UserState, count uint, nextPage string, errorCode RTM_ERROR_CODE) {
	cNextPage := C.CString(nextPage)
	C.C_IRtmEventHandler_onGetOnlineUsersResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		(*C.struct_C_UserState)(unsafe.SliceData(userStateList)),
		C.size_t(count),
		cNextPage,
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
	C.free(unsafe.Pointer(cNextPage))
}

/**
 * Occurs when query which channels the user joined
 *
 * @param requestId The related request id when user perform this operation
 * @param channels The channel informations.
 * @param count The channel count.
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnWhereNowResult(requestId uint64, channels []ChannelInfo, count uint, errorCode RTM_ERROR_CODE) {
	C.C_IRtmEventHandler_onWhereNowResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		(*C.struct_C_ChannelInfo)(unsafe.SliceData(channels)),
		C.size_t(count),
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
}

/**
 * Occurs when query which channels the user joined
 *
 * @param requestId The related request id when user perform this operation
 * @param channels The channel informations.
 * @param count The channel count.
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnGetUserChannelsResult(requestId uint64, channels []ChannelInfo, count uint, errorCode RTM_ERROR_CODE) {
	C.C_IRtmEventHandler_onGetUserChannelsResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		(*C.struct_C_ChannelInfo)(unsafe.SliceData(channels)),
		C.size_t(count),
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
}

/**
 * Occurs when set user presence
 *
 * @param requestId The related request id when user perform this operation
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnPresenceSetStateResult(requestId uint64, errorCode RTM_ERROR_CODE) {
	C.C_IRtmEventHandler_onPresenceSetStateResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
}

/**
 * Occurs when delete user presence
 *
 * @param requestId The related request id when user perform this operation
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnPresenceRemoveStateResult(requestId uint64, errorCode RTM_ERROR_CODE) {
	C.C_IRtmEventHandler_onPresenceRemoveStateResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
}

/**
 * Occurs when get user presence
 *
 * @param requestId The related request id when user perform this operation
 * @param states The user states
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnPresenceGetStateResult(requestId uint64, state *UserState, errorCode RTM_ERROR_CODE) {
	C.C_IRtmEventHandler_onPresenceGetStateResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		(*C.struct_C_UserState)(state),
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
}

/**
 * Occurs when user logout
 *
 * @param requestId The related request id when user perform this operation
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnLogoutResult(requestId uint64, errorCode RTM_ERROR_CODE) {
	C.C_IRtmEventHandler_onLogoutResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
}

/**
 * Occurs when user renew token
 *
 * @param requestId The related request id when user perform this operation
 * @param serverType The type of server.
 * @param channelName The name of the channel.
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnRenewTokenResult(requestId uint64, serverType RTM_SERVICE_TYPE, channelName string, errorCode RTM_ERROR_CODE) {
	C.C_IRtmEventHandler_onRenewTokenResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		C.enum_C_RTM_SERVICE_TYPE(serverType),
		C.CString(channelName),
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
}

/**
 * Occurs when user publish topic message
 *
 * @param requestId The related request id when user perform this operation
 * @param channelName The name of the channel.
 * @param topic The name of the topic.
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnPublishTopicMessageResult(requestId uint64, channelName string, topic string, errorCode RTM_ERROR_CODE) {
	C.C_IRtmEventHandler_onPublishTopicMessageResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		C.CString(channelName),
		C.CString(topic),
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
}

/**
 * Occurs when user unsubscribe topic
 *
 * @param requestId The related request id when user perform this operation
 * @param channelName The name of the channel.
 * @param topic The name of the topic.
 * @param errorCode The error code.
 */

func (this_ *IRtmEventHandler) OnUnsubscribeTopicResult(requestId uint64, channelName string, topic string, errorCode RTM_ERROR_CODE) {
	C.C_IRtmEventHandler_onUnsubscribeTopicResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		C.CString(channelName),
		C.CString(topic),
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
}

/**
 * Occurs when user get subscribed user list
 *
 * @param requestId The related request id when user perform this operation
 * @param channelName The name of the channel.
 * @param topic The name of the topic.
 * @param users The subscribed user list.
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnGetSubscribedUserListResult(requestId uint64, channelName string, topic string, user UserList, errorCode RTM_ERROR_CODE) {
	C.C_IRtmEventHandler_onGetSubscribedUserListResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		C.CString(channelName),
		C.CString(topic),
		(C.struct_C_UserList)(user),
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
}

/**
 * Occurs when user get history messages
 *
 * @param requestId The related request id when user perform this operation
 * @param errorCode The error code.
 * @param messageList The history message list.
 * @param count The message count.
 * @param newStart The timestamp of next history message. If newStart is 0, means there are no more history messages
 */

func (this_ *IRtmEventHandler) OnGetHistoryMessagesResult(requestId uint64, messageList []HistoryMessage, count uint, newStart uint64, errorCode RTM_ERROR_CODE) {
	C.C_IRtmEventHandler_onGetHistoryMessagesResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		(*C.struct_C_HistoryMessage)(unsafe.SliceData(messageList)),
		C.size_t(count),
		C.uint64_t(newStart),
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
}

/**
 * Occurs when user unsubscribe user metadata
 *
 * @param requestId The related request id when user perform this operation
 * @param userId The id of the user.
 * @param errorCode The error code.
 */
func (this_ *IRtmEventHandler) OnUnsubscribeUserMetadataResult(requestId uint64, userId string, errorCode RTM_ERROR_CODE) {
	C.C_IRtmEventHandler_onUnsubscribeUserMetadataResult(unsafe.Pointer(this_),
		C.uint64_t(requestId),
		C.CString(userId),
		C.enum_C_RTM_ERROR_CODE(errorCode),
	)
}

// #endregion IRtmEventHandler

/**
 * The IRtmClient class.
 *
 * This class provides the main methods that can be invoked by your app.
 *
 * IRtmClient is the basic interface class of the Agora RTM SDK.
 * Creating an IRtmClient object and then calling the methods of
 * this object enables you to use Agora RTM SDK's functionality.
 */
type IRtmClient C.C_IRtmClient

// #region IRtmClient

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
func CreateAgoraRtmClient(config *RtmConfig) *IRtmClient {
	err := 0
	return (*IRtmClient)(C.agora_rtm_client_create((*C.struct_C_RtmConfig)(config), (*C.int)(unsafe.Pointer(&err))))
}

/**
 * Release the rtm client instance.
 *
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmClient) Release() int {
	return int(C.agora_rtm_client_release(unsafe.Pointer(this_)))
}

/**
 * Login the Agora RTM service. The operation result will be notified by \ref agora::rtm::IRtmEventHandler::onLoginResult callback.
 *
 * @param [in] token Token used to login RTM service.
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmClient) Login(token string) int {
	var requestId uint64
	cToken := C.CString(token)
	ret := int(C.agora_rtm_client_login(unsafe.Pointer(this_),
		cToken,
		(*C.uint64_t)(unsafe.Pointer(&requestId)),
	))
	C.free(unsafe.Pointer(cToken))
	return int(ret)
}

/**
 * Logout the Agora RTM service. Be noticed that this method will break the rtm service including storage/lock/presence.
 *
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmClient) Logout() int {
	var requestId uint64
	return int(C.agora_rtm_client_logout(unsafe.Pointer(this_),
		(*C.uint64_t)(unsafe.Pointer(&requestId)),
	))
}

/**
 * Get the storage instance.
 *
 * @return
 * - return NULL if error occurred
 */
func (this_ *IRtmClient) GetStorage() *IRtmStorage {
	return (*IRtmStorage)(C.agora_rtm_client_get_storage(unsafe.Pointer(this_)))
}

/**
 * Get the lock instance.
 *
 * @return
 * - return NULL if error occurred
 */
func (this_ *IRtmClient) GetLock() *IRtmLock {
	return (*IRtmLock)(C.agora_rtm_client_get_lock(unsafe.Pointer(this_)))
}

/**
 * Get the presence instance.
 *
 * @return
 * - return NULL if error occurred
 */
func (this_ *IRtmClient) GetPresence() *IRtmPresence {
	return (*IRtmPresence)(C.agora_rtm_client_get_presence(unsafe.Pointer(this_)))
}

/**
 * Renews the token. Once a token is enabled and used, it expires after a certain period of time.
 * You should generate a new token on your server, call this method to renew it.
 *
 * @param [in] token Token used renew.
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmClient) RenewToken(token string) int {
	cToken := C.CString(token)
	var requestId uint64
	ret := int(C.agora_rtm_client_renew_token(unsafe.Pointer(this_),
		cToken,
		(*C.uint64_t)(unsafe.Pointer(&requestId)),
	))
	C.free(unsafe.Pointer(cToken))
	return int(ret)
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
func (this_ *IRtmClient) Publish(channelName string, message []byte, length uint, option *PublishOptions, requestId *uint64) int {
	cChannelName := C.CString(channelName)
	cMessage := C.CBytes(message)
	ret := int(C.agora_rtm_client_publish(unsafe.Pointer(this_),
		cChannelName,
		(*C.char)(cMessage),
		C.size_t(length),
		(*C.struct_C_PublishOptions)(option),
		(*C.uint64_t)(requestId),
	))
	C.free(unsafe.Pointer(cChannelName))
	C.free(unsafe.Pointer(cMessage))
	return ret
}
func (this_ *IRtmClient) SendChannelMessage(channelName string, message []byte, length uint, requestId *uint64) int {
	cChannelName := C.CString(channelName)
	cMessage := C.CBytes(message)
	opt := NewPublishOptions()
	opt.SetChannelType(RTM_CHANNEL_TYPE_MESSAGE)
	opt.SetMessageType(RTM_MESSAGE_TYPE_BINARY)

	ret := int(C.agora_rtm_client_publish(unsafe.Pointer(this_),
		cChannelName,
		(*C.char)(cMessage),
		C.size_t(length),
		(*C.struct_C_PublishOptions)(opt),
		(*C.uint64_t)(requestId),
	))
	C.free(unsafe.Pointer(cChannelName))
	C.free(unsafe.Pointer(cMessage))
	opt.Delete()
	return ret
}
func (this_ *IRtmClient) SendUserMessage(userId string, message []byte, length uint, requestId *uint64) int {
	cUserId := C.CString(userId)
	cMessage := C.CBytes(message)
	opt := NewPublishOptions()
	opt.SetChannelType(RTM_CHANNEL_TYPE_USER)
	opt.SetMessageType(RTM_MESSAGE_TYPE_BINARY)

	ret := int(C.agora_rtm_client_publish(unsafe.Pointer(this_),
		cUserId,
		(*C.char)(cMessage),
		C.size_t(length),
		(*C.struct_C_PublishOptions)(opt),
		(*C.uint64_t)(requestId),
	))
	C.free(unsafe.Pointer(cUserId))
	C.free(unsafe.Pointer(cMessage))
	opt.Delete()
	return ret
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
func (this_ *IRtmClient) Subscribe(channelName string, option *SubscribeOptions, requestId *uint64) int {
	cChannelName := C.CString(channelName)
	ret := int(C.agora_rtm_client_subscribe(unsafe.Pointer(this_),
		cChannelName,
		(*C.struct_C_SubscribeOptions)(option),
		(*C.uint64_t)(requestId),
	))
	C.free(unsafe.Pointer(cChannelName))
	return ret
}

/**
 * Unsubscribe a channel.
 *
 * @param [in] channelName The name of the channel.
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmClient) Unsubscribe(channelName string) int {
	cChannelName := C.CString(channelName)
	var requestId uint64
	ret := int(C.agora_rtm_client_unsubscribe(unsafe.Pointer(this_),
		cChannelName,
		(*C.uint64_t)(unsafe.Pointer(&requestId)),
	))
	C.free(unsafe.Pointer(cChannelName))
	return int(ret)
}

/**
 * Create a stream channel instance.
 *
 * @param [in] channelName The Name of the channel.
 * @return
 * - return NULL if error occurred
 */
func (this_ *IRtmClient) CreateStreamChannel(channelName string) *IStreamChannel {
	cChannelName := C.CString(channelName)
	var errorCode int
	ret := C.agora_rtm_client_create_stream_channel(unsafe.Pointer(this_),
		cChannelName,
		(*C.int)(unsafe.Pointer(&errorCode)),
	)
	C.free(unsafe.Pointer(cChannelName))

	streamChannel := &IStreamChannel{ptr: (*C.C_IStreamChannel)(ret)}
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
	cParameters := C.CString(parameters)
	ret := int(C.agora_rtm_client_set_parameters(unsafe.Pointer(this_),
		cParameters,
	))
	C.free(unsafe.Pointer(cParameters))
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
