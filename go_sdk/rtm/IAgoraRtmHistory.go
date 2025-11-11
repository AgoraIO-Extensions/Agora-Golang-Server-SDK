package agorartm

/*

#include "C_IAgoraRtmHistory.h"
#include <stdlib.h>

*/
import "C"
import "unsafe"

// #region agora

// #region agora::rtm

/**
 * The IRtmHistory class.
 *
 * This class provides the rtm history methods that can be invoked by your app.
 */
type IRtmHistory struct {
	rtmHistory unsafe.Pointer
}

// #region IRtmHistory

/**
 * Get History Messages Options
 */
type GetHistoryMessagesOptions struct {
	Start int64
	End   int64
	Count int
}

// #region GetHistoryMessagesOptions

/**
 * Start timestamp
 */
func (this_ *GetHistoryMessagesOptions) GetStart() int64 {
	return this_.Start
}

/**
 * Start timestamp
 */
func (this_ *GetHistoryMessagesOptions) SetStart(start int64) {
	this_.Start = start
}

/**
 * End timestamp
 */
func (this_ *GetHistoryMessagesOptions) GetEnd() int64 {
	return this_.End
}

/**
 * End timestamp
 */
func (this_ *GetHistoryMessagesOptions) SetEnd(end int64) {
	this_.End = end
}

/**
 * Message count
 */
func (this_ *GetHistoryMessagesOptions) GetCount() int {
	return this_.Count
}

/**
 * Message count
 */
func (this_ *GetHistoryMessagesOptions) SetCount(count int) {
	this_.Count = count
}

func NewGetHistoryMessagesOptions() *GetHistoryMessagesOptions {
	return &GetHistoryMessagesOptions{
		Start: 0,
		End:   0,
		Count: 100,
	}
}

// #endregion GetHistoryMessagesOptions

/**
 * Gets history messages in the channel.
 *
 * @param [in] channelName The name of the channel.
 * @param [in] channelType The type of the channel.
 * @param [in] options The query options.
 * @param [out] requestId The related request id of this operation.
 * @return
 * - 0: Success.
 * - < 0: Failure.
 */
func (this_ *IRtmHistory) GetMessages(channelName string, channelType RtmChannelType, options *GetHistoryMessagesOptions, requestId *uint64) int {
	cChannelName := C.CString(channelName)
	defer C.free(unsafe.Pointer(cChannelName))

	cOptions := &C.C_GetHistoryMessagesOptions{}
	if options != nil {
		cOptions.start = C.int64_t(options.Start)
		cOptions.end = C.int64_t(options.End)
		cOptions.count = C.int(options.Count)
	} else {
		cOptions.start = 0
		cOptions.end = 0
		cOptions.count = 100
	}

	ret := int(C.agora_rtm_history_get_messages(
		this_.rtmHistory,
		cChannelName,
		C.enum_C_RTM_CHANNEL_TYPE(channelType),
		cOptions,
		(*C.uint64_t)(requestId),
	))
	return ret
}

// #endregion IRtmHistory

// #endregion agora::rtm

// #endregion agora
