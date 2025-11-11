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
type IRtmPresence struct {
	rtmPresence unsafe.Pointer
}

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
func (this_ *IRtmPresence) WhoNow(channelName string, channelType RtmChannelType, options *PresenceOptions, requestId *uint64) int {
	cChannelName := C.CString(channelName)
	defer C.free(unsafe.Pointer(cChannelName))

	cOptions := C.C_PresenceOptions_New()
	defer C.C_PresenceOptions_Delete(cOptions)
	if options != nil {
		cOptions.includeUserId = C.bool(options.IncludeUserId)
		cOptions.includeState = C.bool(options.IncludeState)
		cOptions.page = C.CString(options.Page)
		defer C.free(unsafe.Pointer(cOptions.page))
	} else {
		cOptions.includeUserId = C.bool(false)
		cOptions.includeState = C.bool(false)
		cOptions.page = nil
	}

	ret := int(C.agora_rtm_presence_who_now(this_.rtmPresence,
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		cOptions,
		(*C.uint64_t)(requestId),
	))
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
	defer C.free(unsafe.Pointer(cUserId))

	ret := int(C.agora_rtm_presence_where_now(this_.rtmPresence,
		cUserId,
		(*C.uint64_t)(requestId),
	))
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
func (this_ *IRtmPresence) SetState(channelName string, channelType RtmChannelType, items []*StateItem, count uint, requestId *uint64) int {
	cChannelName := C.CString(channelName)
	defer C.free(unsafe.Pointer(cChannelName))

	actualCount := uint(len(items))
	if count < actualCount {
		actualCount = count
	}

	var cItems *C.struct_C_StateItem
	if actualCount > 0 {
		cItemsArr := make([]*C.struct_C_StateItem, actualCount)
		for i := uint(0); i < actualCount; i++ {
			if items[i] != nil {
				cItemsArr[i] = C.C_StateItem_New()
				defer C.C_StateItem_Delete(cItemsArr[i])
				cItemsArr[i].key = C.CString(items[i].Key)
				cItemsArr[i].value = C.CString(items[i].Value)
				defer C.free(unsafe.Pointer(cItemsArr[i].key))
				defer C.free(unsafe.Pointer(cItemsArr[i].value))
			}
		}

		for _, cItem := range cItemsArr {
			if cItem != nil {
				cItems = cItem
				break
			}
		}
	}

	ret := int(C.agora_rtm_presence_set_state(this_.rtmPresence,
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		cItems,
		C.size_t(actualCount),
		(*C.uint64_t)(requestId),
	))
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
func (this_ *IRtmPresence) RemoveState(channelName string, channelType RtmChannelType, keys []string, count uint, requestId *uint64) int {
	cChannelName := C.CString(channelName)
	defer C.free(unsafe.Pointer(cChannelName))

	actualCount := uint(len(keys))
	if count < actualCount {
		actualCount = count
	}

	var cKeysArr [](*C.char)
	if actualCount > 0 {
		cKeysArr = make([](*C.char), actualCount)
		for i := uint(0); i < actualCount; i++ {
			cKeysArr[i] = C.CString(keys[i])
			defer C.free(unsafe.Pointer(cKeysArr[i]))
		}
	}

	ret := int(C.agora_rtm_presence_remove_state(this_.rtmPresence,
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		unsafe.SliceData(cKeysArr),
		C.size_t(actualCount),
		(*C.uint64_t)(requestId),
	))
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
func (this_ *IRtmPresence) GetState(channelName string, channelType RtmChannelType, userId string, requestId *uint64) int {
	cChannelName := C.CString(channelName)
	defer C.free(unsafe.Pointer(cChannelName))

	cUserId := C.CString(userId)
	defer C.free(unsafe.Pointer(cUserId))

	ret := int(C.agora_rtm_presence_get_state(this_.rtmPresence,
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		cUserId,
		(*C.uint64_t)(requestId),
	))

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
func (this_ *IRtmPresence) GetOnlineUsers(channelName string, channelType RtmChannelType, options *GetOnlineUsersOptions, requestId *uint64) int {
	cChannelName := C.CString(channelName)
	defer C.free(unsafe.Pointer(cChannelName))

	cOptions := C.C_GetOnlineUsersOptions_New()
	defer C.C_GetOnlineUsersOptions_Delete(cOptions)
	if options != nil {
		cOptions.includeUserId = C.bool(options.IncludeUserId)
		cOptions.includeState = C.bool(options.IncludeState)
		cOptions.page = C.CString(options.Page)
		defer C.free(unsafe.Pointer(cOptions.page))
	} else {
		cOptions.includeUserId = C.bool(false)
		cOptions.includeState = C.bool(false)
		cOptions.page = nil
	}

	ret := int(C.agora_rtm_presence_get_online_users(this_.rtmPresence,
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		cOptions,
		(*C.uint64_t)(requestId),
	))
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
	defer C.free(unsafe.Pointer(cUserId))

	ret := int(C.agora_rtm_presence_get_user_channels(this_.rtmPresence,
		cUserId,
		(*C.uint64_t)(requestId),
	))
	return ret
}

// #endregion IRtmPresence

// #endregion agora::rtm

// #endregion agora
