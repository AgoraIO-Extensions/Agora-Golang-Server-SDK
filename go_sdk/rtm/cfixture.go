//go:build test

package agorartm

/*
#include "C_AgoraRtmBase.h"
#include "C_IAgoraRtmClient.h"
#include <stdlib.h>
*/
import "C"

import (
	"testing"
	"unsafe"
)

type userStateSpec struct {
	userID string
	states []StateItem
}

func newCTestCString(t *testing.T, value string) *C.char {
	t.Helper()

	cs := C.CString(value)
	t.Cleanup(func() {
		C.free(unsafe.Pointer(cs))
	})
	return cs
}

func newCTestUserList(t *testing.T, users ...string) C.struct_C_UserList {
	t.Helper()

	var list C.struct_C_UserList
	if len(users) == 0 {
		return list
	}

	n := len(users)
	list.userCount = C.size_t(n)
	raw := C.malloc(C.size_t(n) * C.size_t(unsafe.Sizeof(uintptr(0))))
	if raw == nil {
		t.Fatal("malloc failed for user list")
	}

	charPtrs := (**C.char)(raw)
	cStrings := make([]*C.char, n)
	for i, user := range users {
		cStrings[i] = C.CString(user)
		charPtrSlice := unsafe.Slice(charPtrs, n)
		charPtrSlice[i] = cStrings[i]
	}

	t.Cleanup(func() {
		for _, cs := range cStrings {
			C.free(unsafe.Pointer(cs))
		}
		C.free(raw)
	})

	list.users = charPtrs
	return list
}

func fillCUserState(t *testing.T, dst *C.struct_C_UserState, spec userStateSpec) {
	t.Helper()

	dst.userId = C.CString(spec.userID)
	t.Cleanup(func() {
		C.free(unsafe.Pointer(dst.userId))
	})

	if len(spec.states) == 0 {
		dst.states = nil
		dst.statesCount = 0
		return
	}

	n := len(spec.states)
	dst.statesCount = C.size_t(n)
	ptrRaw := C.malloc(C.size_t(n) * C.size_t(unsafe.Sizeof(uintptr(0))))
	if ptrRaw == nil {
		t.Fatal("malloc failed for state item pointers")
	}

	statePtrs := (**C.struct_C_StateItem)(ptrRaw)
	for i, item := range spec.states {
		si := (*C.struct_C_StateItem)(C.malloc(C.size_t(unsafe.Sizeof(C.struct_C_StateItem{}))))
		if si == nil {
			t.Fatal("malloc failed for state item")
		}
		si.key = C.CString(item.Key)
		si.value = C.CString(item.Value)
		t.Cleanup(func() {
			C.free(unsafe.Pointer(si.key))
			C.free(unsafe.Pointer(si.value))
			C.free(unsafe.Pointer(si))
		})
		statePtrSlice := unsafe.Slice(statePtrs, n)
		statePtrSlice[i] = si
	}

	t.Cleanup(func() {
		C.free(ptrRaw)
	})

	dst.states = (*C.struct_C_StateItem)(unsafe.Pointer(statePtrs))
}

func newCTestUserStateList(t *testing.T, specs ...userStateSpec) (*C.struct_C_UserState, C.size_t) {
	t.Helper()

	if len(specs) == 0 {
		return nil, 0
	}

	n := len(specs)
	raw := C.malloc(C.size_t(n) * C.size_t(unsafe.Sizeof(C.struct_C_UserState{})))
	if raw == nil {
		t.Fatal("malloc failed for user state list")
	}

	t.Cleanup(func() {
		C.free(raw)
	})

	base := (*C.struct_C_UserState)(raw)
	for i, spec := range specs {
		dst := (*C.struct_C_UserState)(unsafe.Pointer(
			uintptr(unsafe.Pointer(base)) + uintptr(i)*unsafe.Sizeof(C.struct_C_UserState{}),
		))
		*dst = C.struct_C_UserState{}
		fillCUserState(t, dst, spec)
	}

	return base, C.size_t(n)
}

func newCTestStateItems(t *testing.T, items ...StateItem) (*C.struct_C_StateItem, C.size_t) {
	t.Helper()

	if len(items) == 0 {
		return nil, 0
	}

	n := len(items)
	raw := C.malloc(C.size_t(n) * C.size_t(unsafe.Sizeof(C.struct_C_StateItem{})))
	if raw == nil {
		t.Fatal("malloc failed for state items")
	}

	t.Cleanup(func() {
		C.free(raw)
	})

	base := (*C.struct_C_StateItem)(raw)
	for i, item := range items {
		dst := (*C.struct_C_StateItem)(unsafe.Pointer(
			uintptr(unsafe.Pointer(base)) + uintptr(i)*unsafe.Sizeof(C.struct_C_StateItem{}),
		))
		dst.key = C.CString(item.Key)
		dst.value = C.CString(item.Value)
		t.Cleanup(func() {
			C.free(unsafe.Pointer(dst.key))
			C.free(unsafe.Pointer(dst.value))
		})
	}

	return base, C.size_t(n)
}

func newCTestPresenceEvent(t *testing.T) *C.struct_C_PresenceEvent {
	t.Helper()

	ev := (*C.struct_C_PresenceEvent)(C.malloc(C.size_t(unsafe.Sizeof(C.struct_C_PresenceEvent{}))))
	if ev == nil {
		t.Fatal("malloc failed for presence event")
	}

	t.Cleanup(func() {
		C.free(unsafe.Pointer(ev))
	})

	*ev = C.struct_C_PresenceEvent{}
	return ev
}

func newCTestTokenEvent(t *testing.T) *C.struct_C_TokenEvent {
	t.Helper()

	ev := (*C.struct_C_TokenEvent)(C.malloc(C.size_t(unsafe.Sizeof(C.struct_C_TokenEvent{}))))
	if ev == nil {
		t.Fatal("malloc failed for token event")
	}
	*ev = C.struct_C_TokenEvent{}
	ev.eventType = C.RTM_TOKEN_EVENT_TYPE_READ_PERMISSION_REVOKED
	ev.reason = C.CString("permission revoked")
	ev.timestamp = 123456789

	channelNames := []string{"alpha", "beta"}
	channelRaw := C.malloc(C.size_t(len(channelNames)) * C.size_t(unsafe.Sizeof(uintptr(0))))
	if channelRaw == nil {
		t.Fatal("malloc failed for token event channels")
	}
	channelPtrs := (**C.char)(channelRaw)
	channelSlice := unsafe.Slice(channelPtrs, len(channelNames))
	for i, channelName := range channelNames {
		channelSlice[i] = C.CString(channelName)
	}
	ev.affectedResources.messageChannels.channels = channelPtrs
	ev.affectedResources.messageChannels.channelCount = C.size_t(len(channelNames))

	t.Cleanup(func() {
		for _, channel := range channelSlice {
			C.free(unsafe.Pointer(channel))
		}
		C.free(channelRaw)
		C.free(unsafe.Pointer(ev.reason))
		C.free(unsafe.Pointer(ev))
	})

	return ev
}

func newCTestEmptyTokenEvent(t *testing.T) *C.struct_C_TokenEvent {
	t.Helper()

	ev := (*C.struct_C_TokenEvent)(C.malloc(C.size_t(unsafe.Sizeof(C.struct_C_TokenEvent{}))))
	if ev == nil {
		t.Fatal("malloc failed for empty token event")
	}
	*ev = C.struct_C_TokenEvent{}
	t.Cleanup(func() { C.free(unsafe.Pointer(ev)) })
	return ev
}

func readCTestRtmConfigOptions(t *testing.T, config *RtmConfig) (uint32, string) {
	t.Helper()

	cConfig := C.C_RtmConfig_New()
	if cConfig == nil {
		t.Fatal("C_RtmConfig_New returned nil")
	}
	cleanup := applyRtmConfigOptions(cConfig, config)
	t.Cleanup(func() {
		cleanup()
		C.C_RtmConfig_Delete(cConfig)
	})

	return uint32(cConfig.reconnectTimeout), C.GoString(cConfig.parameters)
}

func cTestEventHandlerHasTokenCallback(t *testing.T) bool {
	t.Helper()

	handler := CRtmEventHandler()
	if handler == nil {
		t.Fatal("CRtmEventHandler returned nil")
	}
	t.Cleanup(func() { C.C_IRtmEventHandler_Delete(handler) })
	return handler.onTokenEvent != nil
}

func newCTestEventHandler(t *testing.T, client *IRtmClient) *C.struct_C_IRtmEventHandler {
	t.Helper()

	handler := C.C_IRtmEventHandler_New(nil)
	if handler == nil {
		t.Fatal("C_IRtmEventHandler_New returned nil")
	}

	t.Cleanup(func() {
		C.C_IRtmEventHandler_Delete(handler)
	})

	handler.userData = unsafe.Pointer(client)
	return handler
}
