package agorartm

/*


#include "C_IAgoraRtmStorage.h"
#include <stdlib.h>
#include <string.h>

*/
import "C"
import (
	"unsafe"
)

// #region agora
// #region agora::rtm

/**
 * Metadata options.
 */
type MetadataOptions struct {
	RecordTs     bool
	RecordUserId bool
}

// #region MetadataOptions

/**
 * Indicates whether or not to notify server update the modify timestamp of metadata
 */
func (this_ *MetadataOptions) GetRecordTs() bool {
	return this_.RecordTs
}

/**
 * Indicates whether or not to notify server update the modify timestamp of metadata
 */
func (this_ *MetadataOptions) SetRecordTs(recordTs bool) {
	this_.RecordTs = recordTs
}

/**
 * Indicates whether or not to notify server update the modify user id of metadata
 */
func (this_ *MetadataOptions) GetRecordUserId() bool {
	return this_.RecordUserId
}

/**
 * Indicates whether or not to notify server update the modify user id of metadata
 */
func (this_ *MetadataOptions) SetRecordUserId(recordUserId bool) {
	this_.RecordUserId = recordUserId
}

func NewMetadataOptions() *MetadataOptions {
	options := &MetadataOptions{
		RecordTs:     true,
		RecordUserId: true,
	}

	return options
}

// #endregion MetadataOptions

type MetadataItem struct {
	Key          string
	Value        string
	AuthorUserId string
	Revision     uint64
	UpdateTs     uint64
}

// #region MetadataItem

/**
 * The key of the metadata item.
 */
func (this_ *MetadataItem) GetKey() string {
	return this_.Key
}

/**
 * The key of the metadata item.
 */
func (this_ *MetadataItem) SetKey(key string) {
	this_.Key = key
}

/**
 * The value of the metadata item.
 */
func (this_ *MetadataItem) GetValue() string {
	return this_.Value
}

/**
 * The value of the metadata item.
 */
func (this_ *MetadataItem) SetValue(value string) {
	this_.Value = value
}

/**
 * The User ID of the user who makes the latest update to the metadata item.
 */
func (this_ *MetadataItem) GetAquthorUserId() string {
	return this_.AuthorUserId
}

/**
 * The User ID of the user who makes the latest update to the metadata item.
 */
func (this_ *MetadataItem) SetAuthorUserId(authorUserId string) {
	this_.AuthorUserId = authorUserId
}

/**
 * The revision of the metadata item.
 */
func (this_ *MetadataItem) GetRevision() int64 {
	return int64(this_.Revision)
}

/**
 * The revision of the metadata item.
 */
func (this_ *MetadataItem) SetRevision(revision int64) {
	this_.Revision = uint64(revision)
}

/**
 * The Timestamp when the metadata item was last updated.
 */
func (this_ *MetadataItem) GetUpdateTs() int64 {
	return int64(this_.UpdateTs)
}

/**
 * The Timestamp when the metadata item was last updated.
 */
func (this_ *MetadataItem) SetUpdateTs(updateTs int64) {
	this_.UpdateTs = uint64(updateTs)
}

func NewMetadataItem() *MetadataItem {
	item := &MetadataItem{
		Key:          "",
		Value:        "",
		AuthorUserId: "",
		Revision:     0,
		UpdateTs:     0,
	}

	return item
}

// #endregion MetadataItem

// added by wei only for backward compatibility to 0.0.1 version
type IMetadata struct {
	majorRevision int64
	items         []MetadataItem
	itemCount     uint
}

// #region IMetadata

/**
 * Set the major revision of metadata.
 *
 * @param [in] revision The major revision of the metadata.
 */
func (this_ *IMetadata) SetMajorRevision(revision int64) {
	this_.majorRevision = revision
}

/**
 * Get the major revision of metadata.
 *
 * @return the major revision of metadata.
 */
func (this_ *IMetadata) GetMajorRevision() int64 {
	return this_.majorRevision
}

/**
 * Add or revise a metadataItem to current metadata.
 */
func (this_ *IMetadata) SetMetadataItem(item *MetadataItem) {
	if item != nil {
		this_.items = append(this_.items, *item)
		this_.itemCount = uint(len(this_.items))
	}
}

/**
 * Get the metadataItem array of current metadata.
 *
 * @param [out] items The address of the metadataItem array.
 * @param [out] size The size the metadataItem array.
 */
func (this_ *IMetadata) GetMetadataItems(item **MetadataItem, size *uint) {
	if len(this_.items) > 0 {
		*item = &this_.items[0]
	}
	*size = this_.itemCount
}

// #endregion IMetadata

// VoidPtr represents a C void* pointer
type VoidPtr unsafe.Pointer

// IRtmStorage wraps a C void* pointer for the RTM storage interface
type IRtmStorage struct {
	rtmStorage unsafe.Pointer
}

// #region IRtmStorage

/** Creates the metadata object and returns the pointer.
 * @return Pointer of the metadata object.
 */
func (this_ *IRtmStorage) CreateMetadata() *IMetadata {
	metadata := &IMetadata{
		majorRevision: 0,
		items:         make([]MetadataItem, 0),
		itemCount:     0,
	}

	return metadata
}

/**
 * Set the metadata of a specified channel.
 *
 * @param [in] channelName The name of the channel.
 * @param [in] channelType Which channel type, RTM_CHANNEL_TYPE_STREAM or RtmChannelTypeMESSAGE.
 * @param [in] data Metadata data.
 * @param [in] options The options of operate metadata.
 * @param [in] lock lock for operate channel metadata.
 * @param [out] requestId The unique ID of this request.
 *
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmStorage) SetChannelMetadata(channelName string, channelType RtmChannelType, data *IMetadata, options *MetadataOptions, lockName string, requestId *uint64) int {
	cChannelName := C.CString(channelName)
	defer C.free(unsafe.Pointer(cChannelName))
	cLockName := C.CString(lockName)
	defer C.free(unsafe.Pointer(cLockName))

	var cOptions *C.struct_C_MetadataOptions
	if options != nil {
		cOptions = C.C_MetadataOptions_New()
		defer C.C_MetadataOptions_Delete(cOptions)
		cOptions.recordTs = C.bool(options.RecordTs)
		cOptions.recordUserId = C.bool(options.RecordUserId)
	}

	var cData *C.struct_C_Metadata
	if data != nil {
		cData = C.C_Metadata_New()
		defer C.C_Metadata_Delete(cData)
		cData.majorRevision = C.int64_t(data.majorRevision)
		cData.itemCount = C.size_t(data.itemCount)
		cItems := make([]*C.struct_C_MetadataItem, data.itemCount)
		for i := 0; i < int(data.itemCount); i++ {
			cItems[i] = C.C_MetadataItem_New()
			defer C.C_MetadataItem_Delete(cItems[i])
			cItems[i].key = C.CString(data.items[i].Key)
			defer C.free(unsafe.Pointer(cItems[i].key))
			cItems[i].value = C.CString(data.items[i].Value)
			defer C.free(unsafe.Pointer(cItems[i].value))
			cItems[i].authorUserId = C.CString(data.items[i].AuthorUserId)
			defer C.free(unsafe.Pointer(cItems[i].authorUserId))
			cItems[i].revision = C.int64_t(data.items[i].Revision)
			cItems[i].updateTs = C.int64_t(data.items[i].UpdateTs)
		}
		if len(cItems) > 0 {
			cData.items = cItems[0]
		}
	}

	C.agora_rtm_storage_set_channel_metadata(this_.rtmStorage,
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		cData,
		cOptions,
		cLockName,
		(*C.uint64_t)(requestId),
	)
	return 0
}

/**
 * Update the metadata of a specified channel.
 *
 * @param [in] channelName The channel Name of the specified channel.
 * @param [in] channelType Which channel type, RTM_CHANNEL_TYPE_STREAM or RtmChannelTypeMESSAGE.
 * @param [in] data Metadata data.
 * @param [in] options The options of operate metadata.
 * @param [in] lock lock for operate channel metadata.
 * @param [out] requestId The unique ID of this request.
 *
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmStorage) UpdateChannelMetadata(channelName string, channelType RtmChannelType, data *IMetadata, options *MetadataOptions, lockName string, requestId *uint64) {
	cChannelName := C.CString(channelName)
	defer C.free(unsafe.Pointer(cChannelName))
	cLockName := C.CString(lockName)
	defer C.free(unsafe.Pointer(cLockName))

	var cOptions *C.struct_C_MetadataOptions
	if options != nil {
		cOptions = C.C_MetadataOptions_New()
		defer C.C_MetadataOptions_Delete(cOptions)
		cOptions.recordTs = C.bool(options.RecordTs)
		cOptions.recordUserId = C.bool(options.RecordUserId)
	}

	var cData *C.struct_C_Metadata
	if data != nil {
		cData = C.C_Metadata_New()
		defer C.C_Metadata_Delete(cData)
		cData.majorRevision = C.int64_t(data.majorRevision)
		cData.itemCount = C.size_t(data.itemCount)
		cItems := make([]*C.struct_C_MetadataItem, data.itemCount)
		for i := 0; i < int(data.itemCount); i++ {
			cItems[i] = C.C_MetadataItem_New()
			defer C.C_MetadataItem_Delete(cItems[i])
			cItems[i].key = C.CString(data.items[i].Key)
			defer C.free(unsafe.Pointer(cItems[i].key))
			cItems[i].value = C.CString(data.items[i].Value)
			defer C.free(unsafe.Pointer(cItems[i].value))
			cItems[i].authorUserId = C.CString(data.items[i].AuthorUserId)
			defer C.free(unsafe.Pointer(cItems[i].authorUserId))
			cItems[i].revision = C.int64_t(data.items[i].Revision)
			cItems[i].updateTs = C.int64_t(data.items[i].UpdateTs)
		}
		if len(cItems) > 0 {
			cData.items = cItems[0]
		}
	}

	C.agora_rtm_storage_update_channel_metadata(this_.rtmStorage,
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		cData,
		cOptions,
		cLockName,
		(*C.uint64_t)(requestId),
	)
}

/**
 * Remove the metadata of a specified channel.
 *
 * @param [in] channelName The channel Name of the specified channel.
 * @param [in] channelType Which channel type, RTM_CHANNEL_TYPE_STREAM or RtmChannelTypeMESSAGE.
 * @param [in] data Metadata data.
 * @param [in] options The options of operate metadata.
 * @param [in] lock lock for operate channel metadata.
 * @param [out] requestId The unique ID of this request.
 *
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmStorage) RemoveChannelMetadata(channelName string, channelType RtmChannelType, data *IMetadata, options *MetadataOptions, lockName string, requestId *uint64) {
	cChannelName := C.CString(channelName)
	defer C.free(unsafe.Pointer(cChannelName))
	cLockName := C.CString(lockName)
	defer C.free(unsafe.Pointer(cLockName))

	var cOptions *C.struct_C_MetadataOptions
	if options != nil {
		cOptions = C.C_MetadataOptions_New()
		defer C.C_MetadataOptions_Delete(cOptions)
		cOptions.recordTs = C.bool(options.RecordTs)
		cOptions.recordUserId = C.bool(options.RecordUserId)
	}

	var cData *C.struct_C_Metadata
	if data != nil {
		cData = C.C_Metadata_New()
		defer C.C_Metadata_Delete(cData)
		cData.majorRevision = C.int64_t(data.majorRevision)
		cData.itemCount = C.size_t(data.itemCount)
		cItems := make([]*C.struct_C_MetadataItem, data.itemCount)
		for i := 0; i < int(data.itemCount); i++ {
			cItems[i] = C.C_MetadataItem_New()
			defer C.C_MetadataItem_Delete(cItems[i])
			cItems[i].key = C.CString(data.items[i].Key)
			defer C.free(unsafe.Pointer(cItems[i].key))
			cItems[i].value = C.CString(data.items[i].Value)
			defer C.free(unsafe.Pointer(cItems[i].value))
			cItems[i].authorUserId = C.CString(data.items[i].AuthorUserId)
			defer C.free(unsafe.Pointer(cItems[i].authorUserId))
			cItems[i].revision = C.int64_t(data.items[i].Revision)
			cItems[i].updateTs = C.int64_t(data.items[i].UpdateTs)
		}
		if len(cItems) > 0 {
			cData.items = cItems[0]
		}
	}

	C.agora_rtm_storage_remove_channel_metadata(this_.rtmStorage,
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		cData,
		cOptions,
		cLockName,
		(*C.uint64_t)(requestId),
	)
}

/**
 * Get the metadata of a specified channel.
 *
 * @param [in] channelName The channel Name of the specified channel.
 * @param [in] channelType Which channel type, RTM_CHANNEL_TYPE_STREAM or RtmChannelTypeMESSAGE.
 * @param requestId The unique ID of this request.
 *
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmStorage) GetChannelMetadata(channelName string, channelType RtmChannelType, requestId *uint64) {
	cChannelName := C.CString(channelName)
	defer C.free(unsafe.Pointer(cChannelName))

	C.agora_rtm_storage_get_channel_metadata(this_.rtmStorage,
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		(*C.uint64_t)(requestId),
	)
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
	defer C.free(unsafe.Pointer(cUserId))

	var cOptions *C.struct_C_MetadataOptions
	if options != nil {
		cOptions = C.C_MetadataOptions_New()
		defer C.C_MetadataOptions_Delete(cOptions)
		cOptions.recordTs = C.bool(options.RecordTs)
		cOptions.recordUserId = C.bool(options.RecordUserId)
	}

	var cData *C.struct_C_Metadata
	if data != nil {
		cData = C.C_Metadata_New()
		defer C.C_Metadata_Delete(cData)
		cData.majorRevision = C.int64_t(data.majorRevision)
		cData.itemCount = C.size_t(data.itemCount)
		cItems := make([]*C.struct_C_MetadataItem, data.itemCount)
		for i := 0; i < int(data.itemCount); i++ {
			cItems[i] = C.C_MetadataItem_New()
			defer C.C_MetadataItem_Delete(cItems[i])
			cItems[i].key = C.CString(data.items[i].Key)
			defer C.free(unsafe.Pointer(cItems[i].key))
			cItems[i].value = C.CString(data.items[i].Value)
			defer C.free(unsafe.Pointer(cItems[i].value))
			cItems[i].authorUserId = C.CString(data.items[i].AuthorUserId)
			defer C.free(unsafe.Pointer(cItems[i].authorUserId))
			cItems[i].revision = C.int64_t(data.items[i].Revision)
			cItems[i].updateTs = C.int64_t(data.items[i].UpdateTs)
		}
		if len(cItems) > 0 {
			cData.items = cItems[0]
		}
	}

	C.agora_rtm_storage_set_user_metadata(this_.rtmStorage,
		cUserId,
		cData,
		cOptions,
		(*C.uint64_t)(requestId),
	)
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
	defer C.free(unsafe.Pointer(cUserId))

	var cOptions *C.struct_C_MetadataOptions
	if options != nil {
		cOptions = C.C_MetadataOptions_New()
		defer C.C_MetadataOptions_Delete(cOptions)
		cOptions.recordTs = C.bool(options.RecordTs)
		cOptions.recordUserId = C.bool(options.RecordUserId)
	}

	var cData *C.struct_C_Metadata
	if data != nil {
		cData = C.C_Metadata_New()
		defer C.C_Metadata_Delete(cData)
		cData.majorRevision = C.int64_t(data.majorRevision)
		cData.itemCount = C.size_t(data.itemCount)
		cItems := make([]*C.struct_C_MetadataItem, data.itemCount)
		for i := 0; i < int(data.itemCount); i++ {
			cItems[i] = C.C_MetadataItem_New()
			defer C.C_MetadataItem_Delete(cItems[i])
			cItems[i].key = C.CString(data.items[i].Key)
			defer C.free(unsafe.Pointer(cItems[i].key))
			cItems[i].value = C.CString(data.items[i].Value)
			defer C.free(unsafe.Pointer(cItems[i].value))
			cItems[i].authorUserId = C.CString(data.items[i].AuthorUserId)
			defer C.free(unsafe.Pointer(cItems[i].authorUserId))
			cItems[i].revision = C.int64_t(data.items[i].Revision)
			cItems[i].updateTs = C.int64_t(data.items[i].UpdateTs)
		}
		if len(cItems) > 0 {
			cData.items = cItems[0]
		}
	}

	C.agora_rtm_storage_update_user_metadata(this_.rtmStorage,
		cUserId,
		cData,
		cOptions,
		(*C.uint64_t)(requestId),
	)
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
	defer C.free(unsafe.Pointer(cUserId))

	var cOptions *C.struct_C_MetadataOptions
	if options != nil {
		cOptions = C.C_MetadataOptions_New()
		defer C.C_MetadataOptions_Delete(cOptions)
		cOptions.recordTs = C.bool(options.RecordTs)
		cOptions.recordUserId = C.bool(options.RecordUserId)
	}

	var cData *C.struct_C_Metadata
	if data != nil {
		cData = C.C_Metadata_New()
		defer C.C_Metadata_Delete(cData)
		cData.majorRevision = C.int64_t(data.majorRevision)
		cData.itemCount = C.size_t(data.itemCount)
		cItems := make([]*C.struct_C_MetadataItem, data.itemCount)
		for i := 0; i < int(data.itemCount); i++ {
			cItems[i] = C.C_MetadataItem_New()
			defer C.C_MetadataItem_Delete(cItems[i])
			cItems[i].key = C.CString(data.items[i].Key)
			defer C.free(unsafe.Pointer(cItems[i].key))
			cItems[i].value = C.CString(data.items[i].Value)
			defer C.free(unsafe.Pointer(cItems[i].value))
			cItems[i].authorUserId = C.CString(data.items[i].AuthorUserId)
			defer C.free(unsafe.Pointer(cItems[i].authorUserId))
			cItems[i].revision = C.int64_t(data.items[i].Revision)
			cItems[i].updateTs = C.int64_t(data.items[i].UpdateTs)
		}
		if len(cItems) > 0 {
			cData.items = cItems[0]
		}
	}

	C.agora_rtm_storage_remove_user_metadata(this_.rtmStorage,
		cUserId,
		cData,
		cOptions,
		(*C.uint64_t)(requestId),
	)
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
	defer C.free(unsafe.Pointer(cUserId))

	C.agora_rtm_storage_get_user_metadata(this_.rtmStorage,
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
	defer C.free(unsafe.Pointer(cUserId))

	C.agora_rtm_storage_subscribe_user_metadata(this_.rtmStorage,
		cUserId,
		(*C.uint64_t)(requestId),
	)
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
	defer C.free(unsafe.Pointer(cUserId))

	var requestId uint64
	C.agora_rtm_storage_unsubscribe_user_metadata(this_.rtmStorage,
		cUserId,
		(*C.uint64_t)(&requestId),
	)
}

// #endregion IRtmStorage

// #endregion agora::rtm

// #endregion agora
