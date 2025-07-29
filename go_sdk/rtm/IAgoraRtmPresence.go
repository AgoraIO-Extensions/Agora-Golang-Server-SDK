package agorartm

/*


#include "C_IAgoraRtmPresence.h"
#include <stdlib.h>

*/
import "C"
import "unsafe"

// #region agora

// #region agora::rtm

/**
 * The IRtmPresence class.
 *
 * This class provides the rtm presence methods that can be invoked by your app.
 */
type IRtmPresence C.C_IRtmPresence

// #region IRtmPresence

/**
 * To query who joined this channel
 *
 * @param [in] channelName The name of the channel.
 * @param [in] channelType The type of the channel.
 * @param [in] options The query option.
 * @param [out] requestId The related request id of this operation.
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmPresence) WhoNow(channelName string, channelType RTM_CHANNEL_TYPE, options *PresenceOptions, requestId *uint64) int {
	cChannelName := C.CString(channelName)
	ret := int(C.agora_rtm_presence_who_now(unsafe.Pointer(this_),
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		(*C.struct_C_PresenceOptions)(options),
		(*C.uint64_t)(requestId),
	))
	C.free(unsafe.Pointer(cChannelName))
	return ret
}

/**
 * To query which channels the user joined
 *
 * @param [in] userId The id of the user.
 * @param [out] requestId The related request id of this operation.
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmPresence) WhereNow(userId string, requestId *uint64) int {
	cUserId := C.CString(userId)
	ret := int(C.agora_rtm_presence_where_now(unsafe.Pointer(this_),
		cUserId,
		(*C.uint64_t)(requestId),
	))
	C.free(unsafe.Pointer(cUserId))
	return ret
}

/**
 * Set user state
 *
 * @param [in] channelName The name of the channel.
 * @param [in] channelType The type of the channel.
 * @param [in] items The states item of user.
 * @param [in] count The count of states item.
 * @param [out] requestId The related request id of this operation.
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmPresence) SetState(channelName string, channelType RTM_CHANNEL_TYPE, items []StateItem, count uint, requestId *uint64) int {
	cChannelName := C.CString(channelName)
	ret := int(C.agora_rtm_presence_set_state(unsafe.Pointer(this_),
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		(*C.struct_C_StateItem)(unsafe.SliceData(items)),
		C.size_t(count),
		(*C.uint64_t)(requestId),
	))
	C.free(unsafe.Pointer(cChannelName))
	return ret
}

/**
 * Delete user state
 *
 * @param [in] channelName The name of the channel.
 * @param [in] channelType The type of the channel.
 * @param [in] keys The keys of state item.
 * @param [in] count The count of keys.
 * @param [out] requestId The related request id of this operation.
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmPresence) RemoveState(channelName string, channelType RTM_CHANNEL_TYPE, keys []string, count uint, requestId *uint64) int {
	cChannelName := C.CString(channelName)
	cKeysArr := make([](*C.char), 0, len(keys))
	for i, key := range keys {
		cKeysArr[i] = C.CString(key)
	}
	ret := int(C.agora_rtm_presence_remove_state(unsafe.Pointer(this_),
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		unsafe.SliceData(cKeysArr),
		C.size_t(count),
		(*C.uint64_t)(requestId),
	))
	C.free(unsafe.Pointer(cChannelName))
	for _, cKey := range cKeysArr {
		C.free(unsafe.Pointer(cKey))
	}
	return ret
}

/**
 * Get user state
 *
 * @param [in] channelName The name of the channel.
 * @param [in] channelType The type of the channel.
 * @param [in] userId The id of the user.
 * @param [out] requestId The related request id of this operation.
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmPresence) GetState(channelName string, channelType RTM_CHANNEL_TYPE, userId string, requestId *uint64) int {
	cChannelName := C.CString(channelName)
	cUserId := C.CString(userId)
	ret := int(C.agora_rtm_presence_get_state(unsafe.Pointer(this_),
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		cUserId,
		(*C.uint64_t)(requestId),
	))
	C.free(unsafe.Pointer(cChannelName))
	C.free(unsafe.Pointer(cUserId))
	return ret
}

/**
 * To query who joined this channel
 *
 * @param [in] channelName The name of the channel.
 * @param [in] channelType The type of the channel.
 * @param [in] options The query option.
 * @param [out] requestId The related request id of this operation.
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmPresence) GetOnlineUsers(channelName string, channelType RTM_CHANNEL_TYPE, options *GetOnlineUsersOptions, requestId *uint64) int {
	cChannelName := C.CString(channelName)
	ret := int(C.agora_rtm_presence_get_online_users(unsafe.Pointer(this_),
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		(*C.struct_C_GetOnlineUsersOptions)(options),
		(*C.uint64_t)(requestId),
	))
	C.free(unsafe.Pointer(cChannelName))
	return ret
}

/**
 * To query which channels the user joined
 *
 * @param [in] userId The id of the user.
 * @param [out] requestId The related request id of this operation.
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmPresence) GetUserChannels(userId string, requestId *uint64) int {
	cUserId := C.CString(userId)
	ret := int(C.agora_rtm_presence_get_user_channels(unsafe.Pointer(this_),
		cUserId,
		(*C.uint64_t)(requestId),
	))
	C.free(unsafe.Pointer(cUserId))
	return ret
}

// #endregion IRtmPresence

// #endregion agora::rtm

// #endregion agora
