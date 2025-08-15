package agorartm

/*


#include "C_IAgoraRtmStorage.h"
#include <stdlib.h>
#include <string.h>

*/
import "C"
import (
	"fmt"
	"unsafe"
)

// #region agora
// #region agora::rtm

/**
 * Metadata options.
 */
type MetadataOptions C.struct_C_MetadataOptions

// #region MetadataOptions

/**
 * Indicates whether or not to notify server update the modify timestamp of metadata
 */
func (this_ *MetadataOptions) GetRecordTs() bool {
	return bool(this_.recordTs)
}

/**
 * Indicates whether or not to notify server update the modify timestamp of metadata
 */
func (this_ *MetadataOptions) SetRecordTs(recordTs bool) {
	this_.recordTs = C.bool(recordTs)
}

/**
 * Indicates whether or not to notify server update the modify user id of metadata
 */
func (this_ *MetadataOptions) GetRecordUserId() bool {
	return bool(this_.recordUserId)
}

/**
 * Indicates whether or not to notify server update the modify user id of metadata
 */
func (this_ *MetadataOptions) SetRecordUserId(recordUserId bool) {
	this_.recordUserId = C.bool(recordUserId)
}

func NewMetadataOptions() *MetadataOptions {
	return (*MetadataOptions)(C.C_MetadataOptions_New())
}
func (this_ *MetadataOptions) Delete() {
	C.C_MetadataOptions_Delete((*C.struct_C_MetadataOptions)(this_))
}

// #endregion MetadataOptions

type MetadataItem C.struct_C_MetadataItem

// #region MetadataItem

/**
 * The key of the metadata item.
 */
func (this_ *MetadataItem) GetKey() string {
	return C.GoString(this_.key)
}

/**
 * The key of the metadata item.
 */
func (this_ *MetadataItem) SetKey(key string) {
	this_.key = C.CString(key)
}

/**
 * The value of the metadata item.
 */
func (this_ *MetadataItem) GetValue() string {
	return C.GoString(this_.value)
}

/**
 * The value of the metadata item.
 */
func (this_ *MetadataItem) SetValue(value string) {
	this_.value = C.CString(value)
}

/**
 * The User ID of the user who makes the latest update to the metadata item.
 */
func (this_ *MetadataItem) GetAquthorUserId() string {
	return C.GoString(this_.authorUserId)
}

/**
 * The User ID of the user who makes the latest update to the metadata item.
 */
func (this_ *MetadataItem) SetAuthorUserId(authorUserId string) {
	this_.authorUserId = C.CString(authorUserId)
}

/**
 * The revision of the metadata item.
 */
func (this_ *MetadataItem) GetRevision() int64 {
	return int64(this_.revision)
}

/**
 * The revision of the metadata item.
 */
func (this_ *MetadataItem) SetRevision(revision int64) {
	this_.revision = C.int64_t(revision)
}

/**
 * The Timestamp when the metadata item was last updated.
 */
func (this_ *MetadataItem) GetUpdateTs() int64 {
	return int64(this_.updateTs)
}

/**
 * The Timestamp when the metadata item was last updated.
 */
func (this_ *MetadataItem) SetUpdateTs(updateTs int64) {
	this_.updateTs = C.int64_t(updateTs)
}

func NewMetadataItem() *MetadataItem {
	return (*MetadataItem)(C.C_MetadataItem_New())
}
func (this_ *MetadataItem) Delete() {
	C.C_MetadataItem_Delete((*C.struct_C_MetadataItem)(this_))
}

// #endregion MetadataItem

// added by wei only 向后兼容到0.0.1
type IMetadata C.struct_C_Metadata

// #region IMetadata

/**
 * Set the major revision of metadata.
 *
 * @param [in] revision The major revision of the metadata.
 */
func (this_ *IMetadata) SetMajorRevision(revision int64) {
	this_.majorRevision = C.int64_t(revision)
}

/**
 * Get the major revision of metadata.
 *
 * @return the major revision of metadata.
 */
func (this_ *IMetadata) GetMajorRevision() int64 {
	return int64(this_.majorRevision)
}

/**
 * Add or revise a metadataItem to current metadata.
 */
func (this_ *IMetadata) SetMetadataItem(item *MetadataItem) {
	this_.items = (*C.struct_C_MetadataItem)(item)
}

/**
 * Get the metadataItem array of current metadata.
 *
 * @param [out] items The address of the metadataItem array.
 * @param [out] size The size the metadataItem array.
 */
func (this_ *IMetadata) GetMetadataItems(item **MetadataItem, size *uint) {
	//do nothing, can access directly
	*size = uint(this_.itemCount)
}

/**
 * Clear the metadataItem array & reset major revision
 */
func (this_ *IMetadata) ClearMetadata() {
	//do nothing, can access directly
}

/**
 * Release the metadata instance.
 */
func (this_ *IMetadata) Release() {
	//C.C_IMetadata_release(unsafe.Pointer(this_))
	fmt.Println("Release： NOT IMPLEMENTED")
}

// #endregion IMetadata

// VoidPtr represents a C void* pointer
type VoidPtr unsafe.Pointer

// IRtmStorage wraps a C void* pointer for the RTM storage interface
type IRtmStorage struct {
	ptr VoidPtr
}

// #region IRtmStorage

/** Creates the metadata object and returns the pointer.
 * @return Pointer of the metadata object.
 */
func (this_ *IRtmStorage) CreateMetadata() *IMetadata {
	// 通过c的方式来申请一个内存快
	size := C.size_t(unsafe.Sizeof(C.struct_C_Metadata{}))
	metadata := C.malloc(size)
	//metadata = C.malloc(C.sizeof__C.struct_C_Metadata)
	// memset to 0
	C.memset(unsafe.Pointer(metadata), 0, size)
	return (*IMetadata)(metadata)
}

/**
 * Set the metadata of a specified channel.
 *
 * @param [in] channelName The name of the channel.
 * @param [in] channelType Which channel type, RTM_CHANNEL_TYPE_STREAM or RTM_CHANNEL_TYPE_MESSAGE.
 * @param [in] data Metadata data.
 * @param [in] options The options of operate metadata.
 * @param [in] lock lock for operate channel metadata.
 * @param [out] requestId The unique ID of this request.
 *
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmStorage) SetChannelMetadata(channelName string, channelType RTM_CHANNEL_TYPE, data *IMetadata, options *MetadataOptions, lockName string, requestId *uint64) int {
	cChannelName := C.CString(channelName)
	cLockName := C.CString(lockName)
	C.agora_rtm_storage_set_channel_metadata(unsafe.Pointer(this_.ptr),
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		(*C.struct_C_Metadata)(unsafe.Pointer(data)),
		(*C.struct_C_MetadataOptions)(options),
		cLockName,
		(*C.uint64_t)(requestId),
	)
	C.free(unsafe.Pointer(cChannelName))
	C.free(unsafe.Pointer(cLockName))
	return 0
}

/**
 * Update the metadata of a specified channel.
 *
 * @param [in] channelName The channel Name of the specified channel.
 * @param [in] channelType Which channel type, RTM_CHANNEL_TYPE_STREAM or RTM_CHANNEL_TYPE_MESSAGE.
 * @param [in] data Metadata data.
 * @param [in] options The options of operate metadata.
 * @param [in] lock lock for operate channel metadata.
 * @param [out] requestId The unique ID of this request.
 *
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmStorage) UpdateChannelMetadata(channelName string, channelType RTM_CHANNEL_TYPE, data *IMetadata, options *MetadataOptions, lockName string, requestId *uint64) {
	cChannelName := C.CString(channelName)
	cLockName := C.CString(lockName)
	C.agora_rtm_storage_update_channel_metadata(unsafe.Pointer(this_.ptr),
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		(*C.struct_C_Metadata)(unsafe.Pointer(data)),
		(*C.struct_C_MetadataOptions)(options),
		cLockName,
		(*C.uint64_t)(requestId),
	)
	C.free(unsafe.Pointer(cChannelName))
	C.free(unsafe.Pointer(cLockName))
}

/**
 * Remove the metadata of a specified channel.
 *
 * @param [in] channelName The channel Name of the specified channel.
 * @param [in] channelType Which channel type, RTM_CHANNEL_TYPE_STREAM or RTM_CHANNEL_TYPE_MESSAGE.
 * @param [in] data Metadata data.
 * @param [in] options The options of operate metadata.
 * @param [in] lock lock for operate channel metadata.
 * @param [out] requestId The unique ID of this request.
 *
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmStorage) RemoveChannelMetadata(channelName string, channelType RTM_CHANNEL_TYPE, data *IMetadata, options *MetadataOptions, lockName string, requestId *uint64) {
	cChannelName := C.CString(channelName)
	cLockName := C.CString(lockName)
	C.agora_rtm_storage_remove_channel_metadata(unsafe.Pointer(this_.ptr),
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		(*C.struct_C_Metadata)(unsafe.Pointer(data)),
		(*C.struct_C_MetadataOptions)(options),
		cLockName,
		(*C.uint64_t)(requestId),
	)
	C.free(unsafe.Pointer(cChannelName))
	C.free(unsafe.Pointer(cLockName))
}

/**
 * Get the metadata of a specified channel.
 *
 * @param [in] channelName The channel Name of the specified channel.
 * @param [in] channelType Which channel type, RTM_CHANNEL_TYPE_STREAM or RTM_CHANNEL_TYPE_MESSAGE.
 * @param requestId The unique ID of this request.
 *
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmStorage) GetChannelMetadata(channelName string, channelType RTM_CHANNEL_TYPE, requestId *uint64) {
	cChannelName := C.CString(channelName)
	C.agora_rtm_storage_get_channel_metadata(unsafe.Pointer(this_.ptr),
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		(*C.uint64_t)(requestId),
	)
	C.free(unsafe.Pointer(cChannelName))
}

/**
 * Set the metadata of a specified user.
 *
 * @param [in] userId The user ID of the specified user.
 * @param [in] data Metadata data.
 * @param [in] options The options of operate metadata.
 * @param [out] requestId The unique ID of this request.
 *
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmStorage) SetUserMetadata(userId string, data *IMetadata, options *MetadataOptions, requestId *uint64) {
	cUserId := C.CString(userId)
	C.agora_rtm_storage_set_user_metadata(unsafe.Pointer(this_.ptr),
		cUserId,
		(*C.struct_C_Metadata)(unsafe.Pointer(data)),
		(*C.struct_C_MetadataOptions)(options),
		(*C.uint64_t)(requestId),
	)
	C.free(unsafe.Pointer(cUserId))
}

/**
 * Update the metadata of a specified user.
 *
 * @param [in] userId The user ID of the specified user.
 * @param [in] data Metadata data.
 * @param [in] options The options of operate metadata.
 * @param [out] requestId The unique ID of this request.
 *
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmStorage) UpdateUserMetadata(userId string, data *IMetadata, options *MetadataOptions, requestId *uint64) {
	cUserId := C.CString(userId)
	C.agora_rtm_storage_update_user_metadata(unsafe.Pointer(this_.ptr),
		cUserId,
		(*C.struct_C_Metadata)(unsafe.Pointer(data)),
		(*C.struct_C_MetadataOptions)(options),
		(*C.uint64_t)(requestId),
	)
	C.free(unsafe.Pointer(cUserId))
}

/**
 * Remove the metadata of a specified user.
 *
 * @param [in] userId The user ID of the specified user.
 * @param [in] data Metadata data.
 * @param [in] options The options of operate metadata.
 * @param [out] requestId The unique ID of this request.
 *
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmStorage) RemoveUserMetadata(userId string, data *IMetadata, options *MetadataOptions, requestId *uint64) {
	cUserId := C.CString(userId)
	C.agora_rtm_storage_remove_user_metadata(unsafe.Pointer(this_.ptr),
		cUserId,
		(*C.struct_C_Metadata)(unsafe.Pointer(data)),
		(*C.struct_C_MetadataOptions)(options),
		(*C.uint64_t)(requestId),
	)
	C.free(unsafe.Pointer(cUserId))
}

/**
 * Get the metadata of a specified user.
 *
 * @param [in] userId The user ID of the specified user.
 * @param [out] requestId The unique ID of this request.
 *
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmStorage) GetUserMetadata(userId string, requestId *uint64) {
	cUserId := C.CString(userId)
	C.agora_rtm_storage_get_user_metadata(unsafe.Pointer(this_.ptr),
		cUserId,
		(*C.uint64_t)(requestId),
	)

}

/**
 * Subscribe the metadata update event of a specified user.
 *
 * @param [in] userId The user ID of the specified user.
 *
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmStorage) SubscribeUserMetadata(userId string, requestId *uint64) {
	cUserId := C.CString(userId)
	C.agora_rtm_storage_subscribe_user_metadata(unsafe.Pointer(this_.ptr),
		cUserId,
		(*C.uint64_t)(requestId),
	)
	C.free(unsafe.Pointer(cUserId))
}

/**
 * unsubscribe the metadata update event of a specified user.
 *
 * @param [in] userId The user ID of the specified user.
 *
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmStorage) UnsubscribeUserMetadata(userId string) {
	cUserId := C.CString(userId)
	var requestId uint64
	C.agora_rtm_storage_unsubscribe_user_metadata(unsafe.Pointer(this_.ptr),
		cUserId,
		(*C.uint64_t)(&requestId),
	)
	C.free(unsafe.Pointer(cUserId))
}

// #endregion IRtmStorage

// #endregion agora::rtm

// #endregion agora
