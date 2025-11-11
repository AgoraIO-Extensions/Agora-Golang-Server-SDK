package agorartm

/*
#include <sys/mman.h>
#include <unistd.h>
#include <errno.h>
#include <signal.h>
#include <setjmp.h>
#include <string.h>
#include <stdint.h>

int is_valid_memory(const void* ptr) {
    if (ptr == NULL) return 0;

    // get page size
    long page_size = sysconf(_SC_PAGESIZE);
    if (page_size <= 0) return 0;

    // calculate the start address of the page where the pointer is located
    void* page_start = (void*)((uintptr_t)ptr & ~(page_size - 1));

    // check if the page is accessible
    if (msync(page_start, page_size, MS_ASYNC) == -1) {
        if (errno == ENOMEM) {
            return 0; // page does not exist or is not accessible
        }
    }

    return 1;
}

static jmp_buf segv_buf;

void segv_handler(int sig) {
    longjmp(segv_buf, 1);
}

int safe_strlen(const char* str) {
    struct sigaction old_action, new_action;

    // set the signal handler
    new_action.sa_handler = segv_handler;
    sigemptyset(&new_action.sa_mask);
    new_action.sa_flags = 0;

    sigaction(SIGSEGV, &new_action, &old_action);

    int result = 0;
    if (setjmp(segv_buf) == 0) {
        result = strlen(str);
    } else {
        result = -1; // segment fault
    }

    // restore the original signal handler
    sigaction(SIGSEGV, &old_action, NULL);

    return result;
}
*/
import "C"
import "unsafe"

func IsValidMemory(ptr unsafe.Pointer) bool {
	return C.is_valid_memory(ptr) != 0
}

// FastSafeCGoString - fast version, only do basic check
func FastSafeCGoString(cstr *C.char) string {
	if cstr == nil {
		return ""
	}

	// only do simple memory page check, avoid signal processing overhead
	if C.is_valid_memory(unsafe.Pointer(cstr)) == 0 {
		return ""
	}

	return C.GoString(cstr)
}

// SafeCGoString - full security check version
func SafeCGoString(cstr *C.char) string {
	if cstr == nil {
		return ""
	}

	// check if the memory is accessible
	if C.is_valid_memory(unsafe.Pointer(cstr)) == 0 {
		return ""
	}

	length := C.safe_strlen(cstr)
	if length < 0 {
		return ""
	}

	return C.GoString(cstr)
}

// C.struct_C_UserList to UserList
func CUserListToUserList(cUserList *C.struct_C_UserList) *UserList {
	if cUserList == nil || cUserList.users == nil || cUserList.userCount == 0 {
		return nil
	}

	if !IsValidMemory(unsafe.Pointer(cUserList.users)) {
		return nil
	}

	userCount := int(cUserList.userCount)
	users := make([]string, userCount)

	// use unsafe.Slice to create a string pointer slice
	cUsers := unsafe.Slice((**C.char)(unsafe.Pointer(cUserList.users)), userCount)

	for i := 0; i < userCount; i++ {
		if cUsers[i] != nil {
			users[i] = FastSafeCGoString(cUsers[i])
		} else {
			users[i] = ""
		}
	}

	return &UserList{
		Users: users,
	}
}

// C.struct_C_Metadata to IMetadata
func CMetadataToIMetadata(cMetadata *C.struct_C_Metadata) *IMetadata {
	if cMetadata == nil {
		return nil
	}

	if !IsValidMemory(unsafe.Pointer(cMetadata)) {
		return nil
	}

	itemCount := uint32(cMetadata.itemCount)

	var itemsPtr unsafe.Pointer
	if cMetadata.items != nil {
		itemsPtr = unsafe.Pointer(cMetadata.items)
	}

	majorRevision := int64(cMetadata.majorRevision)

	if itemsPtr == nil || itemCount == 0 || !IsValidMemory(itemsPtr) {
		return &IMetadata{
			majorRevision: majorRevision,
			items:         make([]MetadataItem, 0),
			itemCount:     uint(itemCount),
		}
	}

	items := make([]MetadataItem, itemCount)

	// use unsafe.Slice to create MetadataItem pointer slice
	cItems := unsafe.Slice((**C.struct_C_MetadataItem)(itemsPtr), itemCount)

	for i := 0; i < int(itemCount); i++ {
		if cItems[i] != nil && IsValidMemory(unsafe.Pointer(cItems[i])) {
			items[i] = MetadataItem{
				Key:          FastSafeCGoString(cItems[i].key),
				Value:        FastSafeCGoString(cItems[i].value),
				AuthorUserId: FastSafeCGoString(cItems[i].authorUserId),
				Revision:     uint64(cItems[i].revision),
				UpdateTs:     uint64(cItems[i].updateTs),
			}
		}
	}

	return &IMetadata{
		majorRevision: majorRevision,
		items:         items,
		itemCount:     uint(itemCount),
	}
}

// C.struct_C_LockDetail to LockDetail
func CLockDetailToLockDetail(cLockDetail *C.struct_C_LockDetail) *LockDetail {
	if cLockDetail == nil {
		return nil
	}

	if !IsValidMemory(unsafe.Pointer(cLockDetail)) {
		return nil
	}

	return &LockDetail{
		LockName: FastSafeCGoString(cLockDetail.lockName),
		Owner:    FastSafeCGoString(cLockDetail.owner),
		Ttl:      uint32(cLockDetail.ttl),
	}
}

// C.struct_C_UserState to UserState
func CUserStateToUserState(cUserState *C.struct_C_UserState) *UserState {
	if cUserState == nil {
		return nil
	}

	if !IsValidMemory(unsafe.Pointer(cUserState)) {
		return nil
	}

	states := make([]StateItem, cUserState.statesCount)

	if cUserState.states == nil || cUserState.statesCount == 0 || !IsValidMemory(unsafe.Pointer(cUserState.states)) {
		return &UserState{
			UserId: C.GoString(cUserState.userId),
			States: make([]StateItem, 0),
		}
	}

	// use unsafe.Slice to create StateItem pointer slice
	cStates := unsafe.Slice((**C.struct_C_StateItem)(unsafe.Pointer(cUserState.states)), cUserState.statesCount)

	for i := 0; i < int(cUserState.statesCount); i++ {
		if cStates[i] != nil && IsValidMemory(unsafe.Pointer(cStates[i])) {
			states[i] = StateItem{
				Key:   FastSafeCGoString(cStates[i].key),
				Value: FastSafeCGoString(cStates[i].value),
			}
		}
	}

	return &UserState{
		UserId: FastSafeCGoString(cUserState.userId),
		States: states,
	}
}

// C.struct_C_ChannelInfo to ChannelInfo
func CChannelInfoToChannelInfo(cChannelInfo *C.struct_C_ChannelInfo) *ChannelInfo {
	if cChannelInfo == nil {
		return nil
	}

	if !IsValidMemory(unsafe.Pointer(cChannelInfo)) {
		return nil
	}

	return &ChannelInfo{
		ChannelName: FastSafeCGoString(cChannelInfo.channelName),
		ChannelType: RtmChannelType(cChannelInfo.channelType),
	}
}

// C.struct_C_LinkStateEvent to LinkStateEvent
func CLinkStateEventToLinkStateEvent(cLinkStateEvent *C.struct_C_LinkStateEvent) *LinkStateEvent {
	if cLinkStateEvent == nil {
		return nil
	}

	if !IsValidMemory(unsafe.Pointer(cLinkStateEvent)) {
		return nil
	}

	affectedChannels := make([]string, cLinkStateEvent.affectedChannelCount)

	// use unsafe.Slice to create string pointer slice
	cAffectedChannels := unsafe.Slice((**C.char)(unsafe.Pointer(cLinkStateEvent.affectedChannels)), cLinkStateEvent.affectedChannelCount)

	for i := 0; i < int(cLinkStateEvent.affectedChannelCount); i++ {
		if cAffectedChannels[i] != nil {
			affectedChannels[i] = FastSafeCGoString(cAffectedChannels[i])
		} else {
			affectedChannels[i] = ""
		}
	}

	unrestoredChannels := make([]string, cLinkStateEvent.unrestoredChannelCount)

	// use unsafe.Slice to create string pointer slice
	cUnrestoredChannels := unsafe.Slice((**C.char)(unsafe.Pointer(cLinkStateEvent.unrestoredChannels)), cLinkStateEvent.unrestoredChannelCount)

	for i := 0; i < int(cLinkStateEvent.unrestoredChannelCount); i++ {
		if cUnrestoredChannels[i] != nil {
			unrestoredChannels[i] = FastSafeCGoString(cUnrestoredChannels[i])
		} else {
			unrestoredChannels[i] = ""
		}
	}

	return &LinkStateEvent{
		CurrentState:       int(cLinkStateEvent.currentState),
		PreviousState:      int(cLinkStateEvent.previousState),
		ServiceType:        RtmServiceType(cLinkStateEvent.serviceType),
		Operation:          int(cLinkStateEvent.operation),
		ReasonCode:         int(cLinkStateEvent.reasonCode),
		Reason:             FastSafeCGoString(cLinkStateEvent.reason),
		AffectedChannels:   affectedChannels,
		UnrestoredChannels: unrestoredChannels,
		IsResumed:          bool(cLinkStateEvent.isResumed),
		Timestamp:          uint64(cLinkStateEvent.timestamp),
	}
}

// C.struct_C_HistoryMessage to HistoryMessage
func CHistoryMessageToHistoryMessage(cHistoryMessage *C.struct_C_HistoryMessage) *HistoryMessage {
	if cHistoryMessage == nil {
		return nil
	}

	if !IsValidMemory(unsafe.Pointer(cHistoryMessage)) {
		return nil
	}

	return &HistoryMessage{
		MessageType:   RtmMessageType(cHistoryMessage.messageType),
		Message:       FastSafeCGoString(cHistoryMessage.message),
		MessageLength: uint(cHistoryMessage.messageLength),
		Timestamp:     uint64(cHistoryMessage.timestamp),
		Publisher:     FastSafeCGoString(cHistoryMessage.publisher),
		CustomType:    FastSafeCGoString(cHistoryMessage.customType),
	}
}
