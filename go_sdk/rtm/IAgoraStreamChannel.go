package agorartm

/*


#include "C_IAgoraStreamChannel.h"
#include <stdlib.h>

*/
import "C"
import (
	"unsafe"
)

// #region agora

// #region agora::rtm

/**
 * The qos of rtm message.
 */
type RtmMessageQos C.enum_C_RTM_MESSAGE_QOS

const (
	/**
	 * Will not ensure that messages arrive in order.
	 */
	RtmMessageQosUNORDERED RtmMessageQos = C.RTM_MESSAGE_QOS_UNORDERED
	/**
	 * Will ensure that messages arrive in order.
	 */
	RtmMessageQosORDERED RtmMessageQos = C.RTM_MESSAGE_QOS_ORDERED
)

/**
 * The priority of rtm message.
 */
type RtmMessagePriority C.enum_C_RTM_MESSAGE_PRIORITY

const (
	/**
	 * The highest priority
	 */
	RtmMessagePriorityHIGHEST RtmMessagePriority = C.RTM_MESSAGE_PRIORITY_HIGHEST
	/**
	 * The high priority
	 */
	RtmMessagePriorityHIGH RtmMessagePriority = C.RTM_MESSAGE_PRIORITY_HIGH
	/**
	 * The normal priority (Default)
	 */
	RtmMessagePriorityNORMAL RtmMessagePriority = C.RTM_MESSAGE_PRIORITY_NORMAL
	/**
	 * The low priority
	 */
	RtmMessagePriorityLOW RtmMessagePriority = C.RTM_MESSAGE_PRIORITY_LOW
)

/**
* Join channel options.
 */
type JoinChannelOptions struct {
	Token        string
	WithMetadata bool
	WithPresence bool
	WithLock     bool
	BeQuiet      bool
}

// #region JoinChannelOptions

/**
 * Token used to join channel.
 */
func (this_ *JoinChannelOptions) GetToken() string {
	return this_.Token
}

/**
 * Token used to join channel.
 */
func (this_ *JoinChannelOptions) SetToken(token string) {
	this_.Token = token
}

/**
 * Whether to subscribe channel metadata information
 */
func (this_ *JoinChannelOptions) GetWithMetadata() bool {
	return this_.WithMetadata
}

/**
 * Whether to subscribe channel metadata information
 */
func (this_ *JoinChannelOptions) SetWithMetadata(withMetadata bool) {
	this_.WithMetadata = withMetadata
}

/**
 * Whether to subscribe channel with user presence
 */
func (this_ *JoinChannelOptions) GetWithPresence() bool {
	return this_.WithPresence
}

/**
 * Whether to subscribe channel with user presence
 */
func (this_ *JoinChannelOptions) SetWithPresence(withPresence bool) {
	this_.WithPresence = withPresence
}

/**
 * Whether to subscribe channel with lock
 */
func (this_ *JoinChannelOptions) GetWithLock() bool {
	return this_.WithLock
}

/**
 * Whether to subscribe channel with lock
 */
func (this_ *JoinChannelOptions) SetWithLock(withLock bool) {
	this_.WithLock = withLock
}

/**
 * Whether to be quiet when joining channel
 */
func (this_ *JoinChannelOptions) GetBeQuiet() bool {
	return this_.BeQuiet
}

/**
 * Whether to be quiet when joining channel
 */
func (this_ *JoinChannelOptions) SetBeQuiet(beQuiet bool) {
	this_.BeQuiet = beQuiet
}

func NewJoinChannelOptions() *JoinChannelOptions {
	options := &JoinChannelOptions{
		Token:        "",
		WithMetadata: false,
		WithPresence: true,
		WithLock:     false,
		BeQuiet:      false,
	}

	return options
}

// #endregion JoinChannelOptions

/**
* Join topic options.
 */
type JoinTopicOptions struct {
	qos           RtmMessageQos
	priority      RtmMessagePriority
	meta          string
	syncWithMedia bool
}

// #region JoinTopicOptions

/**
 * The qos of rtm message.
 */
func (this_ *JoinTopicOptions) GetQos() RtmMessageQos {
	return this_.qos
}

/**
 * The qos of rtm message.
 */
func (this_ *JoinTopicOptions) SetQos(qos RtmMessageQos) {
	this_.qos = qos
}

/**
 * The priority of rtm message.
 */
func (this_ *JoinTopicOptions) GetPriority() RtmMessagePriority {
	return this_.priority
}

/**
 * The priority of rtm message.
 */
func (this_ *JoinTopicOptions) SetPriority(priority RtmMessagePriority) {
	this_.priority = priority
}

/**
 * The metaData of topic.
 */
func (this_ *JoinTopicOptions) GetMeta() string {
	return this_.meta
}

/**
 * The metaData of topic.
 */
func (this_ *JoinTopicOptions) SetMeta(meta string) {
	this_.meta = meta
}

/**
 * The rtm data will sync with media
 */
func (this_ *JoinTopicOptions) GetSyncWithMedia() bool {
	return this_.syncWithMedia
}

/**
 * The rtm data will sync with media
 */
func (this_ *JoinTopicOptions) SetSyncWithMedia(syncWithMedia bool) {
	this_.syncWithMedia = syncWithMedia
}

func NewJoinTopicOptions() *JoinTopicOptions {
	return &JoinTopicOptions{
		qos:           RtmMessageQosUNORDERED,
		priority:      RtmMessagePriorityNORMAL,
		meta:          "",
		syncWithMedia: false,
	}
}

// #endregion JoinTopicOptions

/**
* Topic options.
 */
type TopicOptions struct {
	users     []string
	userCount uint
}

// #region TopicOptions

/**
 * The list of users to subscribe.
 */
func (this_ *TopicOptions) GetUsers() []string {
	return this_.users
}

/**
 * The list of users to subscribe.
 */
func (this_ *TopicOptions) SetUsers(users []string) {
	this_.users = make([]string, len(users))
	copy(this_.users, users)
	this_.userCount = uint(len(users))
}

/**
 * The number of users.
 */
func (this_ *TopicOptions) GetUserCount() uint {
	return this_.userCount
}

/**
 * The number of users.
 */
func (this_ *TopicOptions) SetUserCount(count uint) {
	this_.userCount = count
}

func NewTopicOptions() *TopicOptions {
	return &TopicOptions{
		users:     make([]string, 0),
		userCount: 0,
	}
}

// #endregion TopicOptions

/**
* The IStreamChannel class.
*
* This class provides the stream channel methods that can be invoked by your app.
 */
type IStreamChannel struct {
	streamChannel unsafe.Pointer
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
	cOptions := C.C_JoinChannelOptions_New()
	defer C.C_JoinChannelOptions_Delete(cOptions)
	cOptions.token = C.CString(options.Token)
	defer C.free(unsafe.Pointer(cOptions.token))
	cOptions.withMetadata = C.bool(options.WithMetadata)
	cOptions.withPresence = C.bool(options.WithPresence)
	cOptions.withLock = C.bool(options.WithLock)
	cOptions.beQuiet = C.bool(options.BeQuiet)
	C.agora_rtm_stream_channel_join(this_.streamChannel,
		cOptions,
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
	defer C.free(unsafe.Pointer(cToken))
	var requestId uint64
	C.agora_rtm_stream_channel_renew_token(this_.streamChannel,
		cToken,
		(*C.uint64_t)(&requestId),
	)
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
	C.agora_rtm_stream_channel_leave(this_.streamChannel,
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
	ret := C.GoString(C.agora_rtm_stream_channel_get_channel_name(this_.streamChannel))
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
	defer C.free(unsafe.Pointer(cTopic))
	cOptions := C.C_JoinTopicOptions_New()
	defer C.C_JoinTopicOptions_Delete(cOptions)
	if options != nil {
		cOptions.qos = C.enum_C_RTM_MESSAGE_QOS(options.qos)
		cOptions.priority = C.enum_C_RTM_MESSAGE_PRIORITY(options.priority)
		cOptions.meta = C.CString(options.meta)
		defer C.free(unsafe.Pointer(cOptions.meta))
		cOptions.syncWithMedia = C.bool(options.syncWithMedia)
	} else {
		cOptions.qos = C.enum_C_RTM_MESSAGE_QOS(RtmMessageQosUNORDERED)
		cOptions.priority = C.enum_C_RTM_MESSAGE_PRIORITY(RtmMessagePriorityNORMAL)
		cOptions.meta = nil
		cOptions.syncWithMedia = C.bool(false)
	}
	C.agora_rtm_stream_channel_join_topic(this_.streamChannel,
		cTopic,
		cOptions,
		(*C.uint64_t)(requestId),
	)
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
	defer C.free(unsafe.Pointer(cTopic))
	cMessage := C.CString(message)
	defer C.free(unsafe.Pointer(cMessage))
	var requestId uint64
	cOption := C.C_TopicMessageOptions_New()
	defer C.C_TopicMessageOptions_Delete(cOption)
	if option != nil {
		cOption.messageType = C.enum_C_RTM_MESSAGE_TYPE(option.MessageType)
		cOption.sendTs = C.uint64_t(option.SendTs)
		cOption.customType = C.CString(option.CustomType)
		defer C.free(unsafe.Pointer(cOption.customType))
	} else {
		cOption.messageType = C.enum_C_RTM_MESSAGE_TYPE(RtmMessageTypeBINARY)
		cOption.sendTs = 0
		cOption.customType = nil
	}
	C.agora_rtm_stream_channel_publish_topic_message(this_.streamChannel,
		cTopic,
		cMessage,
		C.size_t(length),
		cOption,
		(*C.uint64_t)(&requestId),
	)
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
	defer C.free(unsafe.Pointer(cTopic))
	C.agora_rtm_stream_channel_leave_topic(this_.streamChannel,
		cTopic,
		(*C.uint64_t)(requestId),
	)
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
	defer C.free(unsafe.Pointer(cTopic))
	cOptions := C.C_TopicOptions_New()
	defer C.C_TopicOptions_Delete(cOptions)
	if options != nil {
		if len(options.users) > 0 {
			users := make([]*C.char, len(options.users))
			defer func() {
				for _, user := range users {
					C.free(unsafe.Pointer(user))
				}
			}()
			for i, user := range options.users {
				users[i] = C.CString(user)
			}
			cOptions.users = (**C.char)(unsafe.Pointer(&users[0]))
			cOptions.userCount = C.size_t(len(options.users))
		} else {
			cOptions.users = nil
			cOptions.userCount = 0
		}
	} else {
		cOptions.users = nil
		cOptions.userCount = 0
	}
	C.agora_rtm_stream_channel_subscribe_topic(this_.streamChannel,
		cTopic,
		cOptions,
		(*C.uint64_t)(requestId),
	)
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
	defer C.free(unsafe.Pointer(cTopic))
	cOptions := C.C_TopicOptions_New()
	defer C.C_TopicOptions_Delete(cOptions)
	if options != nil {
		if len(options.users) > 0 {
			users := make([]*C.char, len(options.users))
			defer func() {
				for _, user := range users {
					C.free(unsafe.Pointer(user))
				}
			}()
			for i, user := range options.users {
				users[i] = C.CString(user)
			}
			cOptions.users = (**C.char)(unsafe.Pointer(&users[0]))
			cOptions.userCount = C.size_t(len(options.users))
		} else {
			cOptions.users = nil
			cOptions.userCount = 0
		}
	} else {
		cOptions.users = nil
		cOptions.userCount = 0
	}
	var requestId uint64
	C.agora_rtm_stream_channel_unsubscribe_topic(this_.streamChannel,
		cTopic,
		cOptions,
		(*C.uint64_t)(&requestId),
	)
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
	defer C.free(unsafe.Pointer(cTopic))
	var requestId uint64
	C.agora_rtm_stream_channel_get_subscribed_user_list(this_.streamChannel,
		cTopic,
		(*C.uint64_t)(&requestId),
	)
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
	if this_.streamChannel == nil {
		return -1
	}
	ret := C.agora_rtm_stream_channel_release(this_.streamChannel)
	this_.streamChannel = nil
	return int(ret)
}

// #endregion IStreamChannel

// #endregion agora::rtm

// #endregion agora
