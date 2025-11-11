package agorartm

/*


#include "C_IAgoraRtmLock.h"
#include <stdlib.h>

*/
import "C"
import "unsafe"

// #region agora

// #region agora::rtm

/**
 * The IRtmLock class.
 *
 * This class provides the rtm lock methods that can be invoked by your app.
 */
type IRtmLock struct {
	rtmLock unsafe.Pointer
}

// #region IRtmLock

/**
 * sets a lock
 *
 * @param [in] channelName The name of the channel.
 * @param [in] channelType The type of the channel.
 * @param [in] lockName The name of the lock.
 * @param [in] ttl The lock ttl.
 * @param [out] requestId The related request id of this operation.
 */
func (this_ *IRtmLock) SetLock(channelName string, channelType RtmChannelType, lockName string, ttl uint32, requestId *uint64) {
	cChannelName := C.CString(channelName)
	defer C.free(unsafe.Pointer(cChannelName))
	cLockName := C.CString(lockName)
	defer C.free(unsafe.Pointer(cLockName))

	C.agora_rtm_lock_set_lock(this_.rtmLock,
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		cLockName,
		C.uint32_t(ttl),
		(*C.uint64_t)(requestId),
	)
}

/**
 * gets locks in the channel
 *
 * @param [in] channelName The name of the channel.
 * @param [in] channelType The type of the channel.
 * @param [out] requestId The related request id of this operation.
 */
func (this_ *IRtmLock) GetLocks(channelName string, channelType RtmChannelType, requestId *uint64) {
	cChannelName := C.CString(channelName)
	C.agora_rtm_lock_get_locks(this_.rtmLock,
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		(*C.uint64_t)(requestId),
	)
	C.free(unsafe.Pointer(cChannelName))
}

/**
 * removes a lock
 *
 * @param [in] channelName The name of the channel.
 * @param [in] channelType The type of the channel.
 * @param [in] lockName The name of the lock.
 * @param [out] requestId The related request id of this operation.
 */
func (this_ *IRtmLock) RemoveLock(channelName string, channelType RtmChannelType, lockName string, requestId *uint64) {
	cChannelName := C.CString(channelName)
	cLockName := C.CString(lockName)
	C.agora_rtm_lock_remove_lock(this_.rtmLock,
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		cLockName,
		(*C.uint64_t)(requestId),
	)
	C.free(unsafe.Pointer(cChannelName))
	C.free(unsafe.Pointer(cLockName))
}

/**
 * acquires a lock
 *
 * @param [in] channelName The name of the channel.
 * @param [in] channelType The type of the channel.
 * @param [in] lockName The name of the lock.
 * @param [in] retry Whether to automatically retry when acquires lock failed
 * @param [out] requestId The related request id of this operation.
 */
func (this_ *IRtmLock) AcquireLock(channelName string, channelType RtmChannelType, lockName string, retry bool, requestId *uint64) {
	cChannelName := C.CString(channelName)
	cLockName := C.CString(lockName)
	C.agora_rtm_lock_acquire_lock(this_.rtmLock,
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		cLockName,
		C.bool(retry),
		(*C.uint64_t)(requestId),
	)
	C.free(unsafe.Pointer(cChannelName))
	C.free(unsafe.Pointer(cLockName))
}

/**
 * releases a lock
 *
 * @param [in] channelName The name of the channel.
 * @param [in] channelType The type of the channel.
 * @param [in] lockName The name of the lock.
 * @param [out] requestId The related request id of this operation.
 */
func (this_ *IRtmLock) ReleaseLock(channelName string, channelType RtmChannelType, lockName string, requestId *uint64) {
	cChannelName := C.CString(channelName)
	cLockName := C.CString(lockName)
	C.agora_rtm_lock_release_lock(this_.rtmLock,
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		cLockName,
		(*C.uint64_t)(requestId),
	)
	C.free(unsafe.Pointer(cChannelName))
	C.free(unsafe.Pointer(cLockName))
}

/**
 * disables a lock
 *
 * @param [in] channelName The name of the channel.
 * @param [in] channelType The type of the channel.
 * @param [in] lockName The name of the lock.
 * @param [in] owner The lock owner.
 * @param [out] requestId The related request id of this operation.
 */
func (this_ *IRtmLock) RevokeLock(channelName string, channelType RtmChannelType, lockName string, owner string, requestId *uint64) {
	cChannelName := C.CString(channelName)
	cLockName := C.CString(lockName)
	cOwner := C.CString(owner)
	C.agora_rtm_lock_revoke_lock(this_.rtmLock,
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		cLockName,
		cOwner,
		(*C.uint64_t)(requestId),
	)
	C.free(unsafe.Pointer(cChannelName))
	C.free(unsafe.Pointer(cLockName))
	C.free(unsafe.Pointer(cOwner))
}

// #endregion IRtmLock

// #endregion agora::rtm

// #endregion agora
