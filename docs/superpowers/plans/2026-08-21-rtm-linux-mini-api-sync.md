# RTM Linux MINI API Sync Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Update `go_sdk/rtm` to consume the already-synchronized Linux MINI RTM C API, including token events, new configuration fields, and new error-code constants.

**Architecture:** Keep the existing cgo bridge architecture. Convert borrowed C callback data into owned Go event values before invoking user callbacks, register only callbacks present in the new `C_IRtmEventHandler`, and pass new config strings through temporary C allocations during client creation. The user supplies matching Linux headers and `.so` files; this plan does not modify SDK resources.

**Tech Stack:** Go, cgo, C RTM headers, existing Go `testing` package, Makefile RTM test target.

---

### Task 1: Add failing token-event conversion tests

**Files:**
- Create: `go_sdk/rtm/token_event_test.go`
- Modify: `go_sdk/rtm/cfixture.go`

- [ ] **Step 1: Add a C fixture that allocates the new event shape**

In the existing cgo preamble in `go_sdk/rtm/cfixture.go`, add helpers that allocate a `C_TokenEvent`, its two channel strings, and the channel pointer array, then add a matching free helper. The fixture must populate `eventType=RTM_TOKEN_EVENT_TYPE_READ_PERMISSION_REVOKED`, `reason="permission revoked"`, channels `alpha`/`beta`, and timestamp `123456789`.

- [ ] **Step 2: Write the failing Go conversion test**

Add a test that calls the fixture, converts the C pointer through `NewTokenEvent().fromC`, and asserts event type, reason, timestamp, and both channel names. Add a second test using a zero-valued C event and assert empty reason plus an empty channel slice.

- [ ] **Step 3: Run the focused test and verify the expected failure**

Run:

```bash
go test -C go_sdk/rtm -run 'TestTokenEventFromC' -count=1
```

Expected result before implementation: compile failure because `TokenEvent` and its conversion do not yet exist.

### Task 2: Implement token-event types and conversion

**Files:**
- Modify: `go_sdk/rtm/AgoraRtmBase.go`
- Modify: `go_sdk/rtm/IAgoraRtmClient.go`
- Modify: `go_sdk/rtm/type_convert.go`

- [ ] **Step 1: Add the token-event enum and new error constants**

In `AgoraRtmBase.go`, define `type RtmTokenEventType C.enum_C_RTM_TOKEN_EVENT_TYPE` and constants mapped to `C.RTM_TOKEN_EVENT_TYPE_WILL_EXPIRE` and `C.RTM_TOKEN_EVENT_TYPE_READ_PERMISSION_REVOKED`. Add these untyped Go constants while keeping callback signatures as `int`:

```go
const (
	RtmErrorDuplicateUserID                     = -10027
	RtmErrorChannelSubscribePermissionDenied    = -11038
	RtmErrorChannelPublishPermissionDenied      = -11039
	RtmErrorChannelSubscribeCanceled            = -11040
	RtmErrorStoragePermissionDenied             = -12020
	RtmErrorPresenceInactive                    = -13014
	RtmErrorPresencePendingRequestCanceled      = -13015
	RtmErrorPresencePendingRequestExceedLimit   = -13016
	RtmErrorLockPermissionDenied                = -14010
	RtmErrorHistoryPermissionDenied             = -15006
)
```

- [ ] **Step 2: Add owned Go token-event structures**

In `IAgoraRtmClient.go`, add:

```go
type AffectedResources struct {
	MessageChannels []string
}

type TokenEvent struct {
	EventType         RtmTokenEventType
	Reason            string
	AffectedResources AffectedResources
	Timestamp         uint64
}
```

Provide `NewAffectedResources`, `NewTokenEvent`, and a `fromC(*C.struct_C_TokenEvent)` method. `fromC` must return safely for nil/invalid pointers, copy `reason` with `FastSafeCGoString`, copy each valid C channel string, and leave `MessageChannels` as a non-nil empty slice for nil or zero-count input.

- [ ] **Step 3: Add a reusable channel-list converter**

In `type_convert.go`, add a helper that accepts `*C.struct_C_ChannelList`, validates the pointer and count, reads the `const char** channels` array with `unsafe.Slice`, copies each string, and returns `[]string{}` for nil/zero input. Use this helper from `TokenEvent.fromC`.

- [ ] **Step 4: Run the focused tests and verify they pass**

Run:

```bash
go test -C go_sdk/rtm -run 'TestTokenEventFromC' -count=1
```

Expected result: both conversion tests pass.

### Task 3: Add config-field coverage and mapping

**Files:**
- Create: `go_sdk/rtm/rtm_config_test.go`
- Modify: `go_sdk/rtm/IAgoraRtmClient.go`

- [ ] **Step 1: Write the failing config mapping test**

Add a cgo test helper that allocates a `C_RtmConfig`, and a Go helper used only by the test to populate the C struct from a `RtmConfig`. Assert that `ReconnectTimeout` reaches `cConfig.reconnectTimeout` and that `Parameters` reaches the expected C string. Free all C allocations in the test.

- [ ] **Step 2: Add Go config fields and defaults**

Extend `RtmConfig` with `ReconnectTimeout uint32` and `Parameters string`. Initialize them to `0` and `""` in `NewRtmConfig`.

- [ ] **Step 3: Map and free the fields in `NewRtmClient`**

Set `cConfig.reconnectTimeout = C.uint32_t(config.ReconnectTimeout)`. Allocate `cConfig.parameters = C.CString(config.Parameters)` and defer `C.free(unsafe.Pointer(cConfig.parameters))` alongside `appId` and `userId`. Do not pass Go string memory to C.

- [ ] **Step 4: Run the config test and verify it passes**

Run:

```bash
go test -C go_sdk/rtm -run 'TestRtmConfig' -count=1
```

Expected result: the new mapping test passes without leaking the temporary C string.

### Task 4: Update the event-handler bridge for the new callback ABI

**Files:**
- Modify: `go_sdk/rtm/RtmEventHandlerBridge.go`
- Modify: `go_sdk/rtm/rtm_event_handler_bridge_test.go`

- [ ] **Step 1: Write the failing callback test**

Extend the bridge test handler with `OnTokenEvent`, invoke `cgo_RtmEventHandlerBridge_onTokenEvent` using a fixture event, and assert that the callback receives the copied values. Add a structural assertion that `CRtmEventHandler()` assigns the new callback pointer when the matching Linux `.so` is installed.

- [ ] **Step 2: Update the cgo declarations and public handler**

Declare `cgo_RtmEventHandlerBridge_onTokenEvent` with `struct C_TokenEvent *event`. Add `OnTokenEvent func(event *TokenEvent)` to `RtmEventHandler`. In `CRtmEventHandler`, assign `ret.onTokenEvent` to the new bridge function.

- [ ] **Step 3: Implement the exported token callback**

Add `cgo_RtmEventHandlerBridge_onTokenEvent`, guard nil handler/client/event-handler/callback values, construct `NewTokenEvent`, call `fromC`, then invoke `OnTokenEvent`.

- [ ] **Step 4: Remove the obsolete callback completely**

Delete `OnConnectionStateChanged` from `RtmEventHandler`, remove its cgo declaration, remove its registration assignment, and delete the exported `cgo_RtmEventHandlerBridge_onConnectionStateChanged` function. Do not leave references to `C_RTM_CONNECTION_STATE` or `C_RTM_CONNECTION_CHANGE_REASON` in `go_sdk/rtm`.

- [ ] **Step 5: Run bridge tests and verify they pass**

Run:

```bash
go test -C go_sdk/rtm -run 'TestRtmEventHandlerBridge' -count=1
```

Expected result: token callback and existing callback tests pass with the synchronized Linux headers and libraries.

### Task 5: Format, compile, and run the complete RTM suite

**Files:**
- Modify: `go_sdk/rtm/*.go` only as required by formatting.

- [ ] **Step 1: Format the changed Go files**

Run:

```bash
gofmt -w go_sdk/rtm/AgoraRtmBase.go go_sdk/rtm/IAgoraRtmClient.go go_sdk/rtm/RtmEventHandlerBridge.go go_sdk/rtm/type_convert.go go_sdk/rtm/cfixture.go go_sdk/rtm/token_event_test.go go_sdk/rtm/rtm_config_test.go go_sdk/rtm/rtm_event_handler_bridge_test.go
```

- [ ] **Step 2: Compile the RTM package for Linux**

Run:

```bash
go test -C go_sdk/rtm -run '^$'
```

Expected result: package compiles and links against the user-installed matching Linux RTM `.so` files.

- [ ] **Step 3: Run the complete RTM test target**

Run:

```bash
make test-rtm
```

Expected result: all RTM tests pass with no references to the removed callback or old C enum names.

- [ ] **Step 4: Verify the final diff**

Run:

```bash
git diff --check
git status --short
git diff --stat
```

Confirm only intended Go RTM files and the implementation plan/spec changes are present; do not stage unrelated user files.
