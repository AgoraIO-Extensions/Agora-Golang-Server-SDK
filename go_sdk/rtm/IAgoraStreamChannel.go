package agorartm

/*


#include "C_IAgoraStreamChannel.h"
#include <stdlib.h>

*/
import "C"
import "unsafe"

// #region agora

// #region agora::rtm

/**
 * The qos of rtm message.
 */
type RTM_MESSAGE_QOS C.enum_C_RTM_MESSAGE_QOS

const (
	/**
	 * Will not ensure that messages arrive in order.
	 */
	RTM_MESSAGE_QOS_UNORDERED RTM_MESSAGE_QOS = C.RTM_MESSAGE_QOS_UNORDERED
	/**
	 * Will ensure that messages arrive in order.
	 */
	RTM_MESSAGE_QOS_ORDERED RTM_MESSAGE_QOS = C.RTM_MESSAGE_QOS_ORDERED
)

/**
 * The priority of rtm message.
 */
type RTM_MESSAGE_PRIORITY C.enum_C_RTM_MESSAGE_PRIORITY

const (
	/**
	 * The highest priority
	 */
	RTM_MESSAGE_PRIORITY_HIGHEST RTM_MESSAGE_PRIORITY = C.RTM_MESSAGE_PRIORITY_HIGHEST
	/**
	 * The high priority
	 */
	RTM_MESSAGE_PRIORITY_HIGH RTM_MESSAGE_PRIORITY = C.RTM_MESSAGE_PRIORITY_HIGH
	/**
	 * The normal priority (Default)
	 */
	RTM_MESSAGE_PRIORITY_NORMAL RTM_MESSAGE_PRIORITY = C.RTM_MESSAGE_PRIORITY_NORMAL
	/**
	 * The low priority
	 */
	RTM_MESSAGE_PRIORITY_LOW RTM_MESSAGE_PRIORITY = C.RTM_MESSAGE_PRIORITY_LOW
)

/**
* Join channel options.
 */
type JoinChannelOptions C.struct_C_JoinChannelOptions

// #region JoinChannelOptions

/**
 * Token used to join channel.
 */
func (this_ *JoinChannelOptions) GetToken() string {
	return C.GoString(this_.token)
}

/**
 * Token used to join channel.
 */
func (this_ *JoinChannelOptions) SetToken(token string) {
	this_.token = C.CString(token)
}

/**
 * Whether to subscribe channel metadata information
 */
func (this_ *JoinChannelOptions) GetWithMetadata() bool {
	return bool(this_.withMetadata)
}

/**
 * Whether to subscribe channel metadata information
 */
func (this_ *JoinChannelOptions) SetWithMetadata(withMetadata bool) {
	this_.withMetadata = C.bool(withMetadata)
}

/**
 * Whether to subscribe channel with user presence
 */
func (this_ *JoinChannelOptions) GetWithPresence() bool {
	return bool(this_.withPresence)
}

/**
 * Whether to subscribe channel with user presence
 */
func (this_ *JoinChannelOptions) SetWithPresence(withPresence bool) {
	this_.withPresence = C.bool(withPresence)
}

/**
 * Whether to subscribe channel with lock
 */
func (this_ *JoinChannelOptions) GetWithLock() bool {
	return bool(this_.withLock)
}

/**
 * Whether to subscribe channel with lock
 */
func (this_ *JoinChannelOptions) SetWithLock(withLock bool) {
	this_.withLock = C.bool(withLock)
}

func NewJoinChannelOptions() *JoinChannelOptions {
	return (*JoinChannelOptions)(C.C_JoinChannelOptions_New())
}

func (this_ *JoinChannelOptions) Delete() {
	C.C_JoinChannelOptions_Delete((*C.struct_C_JoinChannelOptions)(this_))
}

// #endregion JoinChannelOptions

/**
* Join topic options.
 */
type JoinTopicOptions C.struct_C_JoinTopicOptions

// #region JoinTopicOptions

/**
 * The qos of rtm message.
 */
func (this_ *JoinTopicOptions) GetQos() RTM_MESSAGE_QOS {
	return RTM_MESSAGE_QOS(this_.qos)
}

/**
 * The qos of rtm message.
 */
func (this_ *JoinTopicOptions) SetQos(qos RTM_MESSAGE_QOS) {
	this_.qos = C.enum_C_RTM_MESSAGE_QOS(qos)
}

/**
 * The priority of rtm message.
 */
func (this_ *JoinTopicOptions) GetPriority() RTM_MESSAGE_PRIORITY {
	return RTM_MESSAGE_PRIORITY(this_.priority)
}

/**
 * The priority of rtm message.
 */
func (this_ *JoinTopicOptions) SetPriority(priority RTM_MESSAGE_PRIORITY) {
	this_.priority = C.enum_C_RTM_MESSAGE_PRIORITY(priority)
}

/**
 * The metaData of topic.
 */
func (this_ *JoinTopicOptions) GetMeta() string {
	return C.GoString(this_.meta)
}

/**
 * The metaData of topic.
 */
func (this_ *JoinTopicOptions) SetMeta(meta string) {
	this_.meta = C.CString(meta)
}

/**
 * The rtm data will sync with media
 */
func (this_ *JoinTopicOptions) GetSyncWithMedia() bool {
	return bool(this_.syncWithMedia)
}

/**
 * The rtm data will sync with media
 */
func (this_ *JoinTopicOptions) SetSyncWithMedia(syncWithMedia bool) {
	this_.syncWithMedia = C.bool(syncWithMedia)
}

func NewJoinTopicOptions() *JoinTopicOptions {
	return (*JoinTopicOptions)(C.C_JoinTopicOptions_New())
}

func (this_ *JoinTopicOptions) Delete() {
	C.C_JoinTopicOptions_Delete((*C.struct_C_JoinTopicOptions)(this_))
}

// #endregion JoinTopicOptions

/**
* Topic options.
 */
type TopicOptions C.struct_C_TopicOptions

// #region TopicOptions

/**
 * The list of users to subscribe.
 */
func (this_ *TopicOptions) GetUsers() []string {
	count := this_.GetUserCount()
	cStrArr := unsafe.Slice(this_.users, count)
	users := make([]string, 0, count)
	for _, cStr := range cStrArr {
		users = append(users, C.GoString(cStr))
	}
	return users
}

/**
 * The list of users to subscribe.
 */
func (this_ *TopicOptions) SetUsers(users []string) {
	cStrArr := make([]*C.char, 0, len(users))
	for _, goStr := range users {
		cStrArr = append(cStrArr, C.CString(goStr))
	}
	this_.users = unsafe.SliceData(cStrArr)
}

/**
 * The number of users.
 */
func (this_ *TopicOptions) GetUserCount() uint {
	return uint(this_.userCount)
}

/**
 * The number of users.
 */
func (this_ *TopicOptions) SetUserCount(count uint) {
	this_.userCount = C.size_t(count)
}

func NewTopicOptions() *TopicOptions {
	return (*TopicOptions)(C.C_TopicOptions_New())
}

func (this_ *TopicOptions) Delete() {
	C.C_TopicOptions_Delete((*C.struct_C_TopicOptions)(this_))
}

// #endregion TopicOptions

/**
* The IStreamChannel class.
*
* This class provides the stream channel methods that can be invoked by your app.
 */
type IStreamChannel struct {
	ptr *C.C_IStreamChannel //same to void*  in c
}

// #region IStreamChannel

/**
* Join the channel.
*
* @param [in] options join channel options.
* @return
* - 0: Success.
* - < 0: Failure.
 */
func (this_ *IStreamChannel) Join(options *JoinChannelOptions, requestId *uint64) int {
	C.agora_rtm_stream_channel_join(unsafe.Pointer(this_),
		(*C.struct_C_JoinChannelOptions)(options),
		(*C.uint64_t)(requestId),
	)
	return 0
}

/**
* Renews the token. Once a token is enabled and used, it expires after a certain period of time.
* You should generate a new token on your server, call this method to renew it.
*
* @param [in] token Token used renew.
* - 0: Success.
* - < 0: Failure.
 */
func (this_ *IStreamChannel) RenewToken(token string) int {
	cToken := C.CString(token)
	var requestId uint64
	C.agora_rtm_stream_channel_renew_token(unsafe.Pointer(this_),
		cToken,
		(*C.uint64_t)(&requestId),
	)
	C.free(unsafe.Pointer(cToken))
	return 0
}

/**
* Leave the channel.
*
* @return
* - 0: Success.
* - < 0: Failure.
 */
func (this_ *IStreamChannel) Leave(requestId *uint64) int {
	C.agora_rtm_stream_channel_leave(unsafe.Pointer(this_),
		(*C.uint64_t)(requestId),
	)
	return 0
}

/**
* Return the channel name of this stream channel.
*
* @return The channel name.
 */
func (this_ *IStreamChannel) GetChannelName() string {
	ret := C.GoString(C.agora_rtm_stream_channel_get_channel_name(unsafe.Pointer(this_)))
	return ret
}

/**
* Join a topic.
*
* @param [in] topic The name of the topic.
* @param [in] options The options of the topic.
* @return
* - 0: Success.
* - < 0: Failure.
 */
func (this_ *IStreamChannel) JoinTopic(topic string, options *JoinTopicOptions, requestId *uint64) int {
	cTopic := C.CString(topic)
	C.agora_rtm_stream_channel_join_topic(unsafe.Pointer(this_),
		cTopic,
		(*C.struct_C_JoinTopicOptions)(options),
		(*C.uint64_t)(requestId),
	)
	C.free(unsafe.Pointer(cTopic))
	return 0
}

/**
* Publish a message in the topic.
*
* @param [in] topic The name of the topic.
* @param [in] message The content of the message.
* @param [in] length The length of the message.
* @param [in] option The option of the message.
* @return
* - 0: Success.
* - < 0: Failure.
 */
func (this_ *IStreamChannel) PublishTopicMessage(topic string, message string, length uint, option *TopicMessageOptions) int {
	cTopic := C.CString(topic)
	cMessage := C.CString(message)
	var requestId uint64
	C.agora_rtm_stream_channel_publish_topic_message(unsafe.Pointer(this_),
		cTopic,
		cMessage,
		C.size_t(length),
		(*C.struct_C_TopicMessageOptions)(option),
		(*C.uint64_t)(&requestId),
	)
	C.free(unsafe.Pointer(cTopic))
	C.free(unsafe.Pointer(cMessage))
	return 0
}

/**
* Leave the topic.
*
* @param [in] topic The name of the topic.
* @return
* - 0: Success.
* - < 0: Failure.
 */
func (this_ *IStreamChannel) LeaveTopic(topic string, requestId *uint64) int {
	cTopic := C.CString(topic)
	C.agora_rtm_stream_channel_leave_topic(unsafe.Pointer(this_),
		cTopic,
		(*C.uint64_t)(requestId),
	)
	C.free(unsafe.Pointer(cTopic))
	return 0
}

/**
* Subscribe a topic.
*
* @param [in] topic The name of the topic.
* @param [in] options The options of subscribe the topic.
* @return
* - 0: Success.
* - < 0: Failure.
 */
func (this_ *IStreamChannel) SubscribeTopic(topic string, options *TopicOptions, requestId *uint64) int {
	cTopic := C.CString(topic)
	C.agora_rtm_stream_channel_subscribe_topic(unsafe.Pointer(this_),
		cTopic,
		(*C.struct_C_TopicOptions)(options),
		(*C.uint64_t)(requestId),
	)
	C.free(unsafe.Pointer(cTopic))
	return 0
}

/**
* Unsubscribe a topic.
*
* @param [in] topic The name of the topic.
* @return
* - 0: Success.
* - < 0: Failure.
 */
func (this_ *IStreamChannel) UnsubscribeTopic(topic string, options *TopicOptions) int {
	cTopic := C.CString(topic)
	var requestId uint64
	C.agora_rtm_stream_channel_unsubscribe_topic(unsafe.Pointer(this_),
		cTopic,
		(*C.struct_C_TopicOptions)(options),
		(*C.uint64_t)(&requestId),
	)
	C.free(unsafe.Pointer(cTopic))
	return 0
}

/**
* Get subscribed user list
*
* @param [in] topic The name of the topic.
* @param [out] users The list of subscribed users.
* @return
* - 0: Success.
* - < 0: Failure.
 */
func (this_ *IStreamChannel) GetSubscribedUserList(topic string) int {
	cTopic := C.CString(topic)
	var requestId uint64
	C.agora_rtm_stream_channel_get_subscribed_user_list(unsafe.Pointer(this_),
		cTopic,
		(*C.uint64_t)(&requestId),
	)
	C.free(unsafe.Pointer(cTopic))
	return 0
}

/**
* Release the stream channel instance.
*
* @return
* - 0: Success.
* - < 0: Failure.
 */
func (this_ *IStreamChannel) Release() int {
	C.agora_rtm_stream_channel_release(unsafe.Pointer(this_))
	return 0
}

// #endregion IStreamChannel

// #endregion agora::rtm

// #endregion agora
