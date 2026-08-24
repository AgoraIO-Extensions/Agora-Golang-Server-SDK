# RTM Linux MINI API Sync Design

## Goal

Adapt the repository's Go RTM binding to the Linux MINI RTM C API supplied in:

- `/Users/weihognqin/Documents/work/agorartmsdkforc/include`
- `/Users/weihognqin/Documents/work/agorartmsdkforc/Shengwang_Native_SDK_for_Linux_MINI_RTM/rtm/sdk/high_level_api/include`

The Go package must target the latest Linux C API. Existing Mac RTM support remains separate. The
caller is responsible for synchronizing the matching Linux headers and shared libraries into the
repository before building or running the Linux RTM package.

## Scope

### SDK resource prerequisite

The caller updates the repository's Linux RTM C headers and matching `.so` files before using this
binding. The existing cgo include and library paths remain the source of truth. This task does not
copy, replace, or otherwise edit SDK headers, bridge headers, or binaries.

### Go binding

Update `go_sdk/rtm` to match the Linux C ABI:

- expose the new token event type and error codes;
- add Go representations and conversion for `C_TokenEvent`, affected message channels, reason, and timestamp;
- add `RtmEventHandler.OnTokenEvent` and wire it through the C callback bridge;
- map `RtmConfig.ReconnectTimeout` and `RtmConfig.Parameters` into `C_RtmConfig` with correct ownership/lifetime handling;
- remove `OnConnectionStateChanged`, its C declarations, bridge implementation, registration, and related old connection-state enum usage;
- continue exposing `OnLinkStateEvent` as the connection-state callback for the current API.

No compatibility alias or deprecated replacement for `OnConnectionStateChanged` is required.

## Data flow

The C callback bridge receives a borrowed `C_TokenEvent` valid for the callback duration. The Go bridge converts all strings and channel names into owned Go values before invoking `OnTokenEvent`, so callbacks do not retain C memory. Nil pointers and zero channel counts produce empty Go strings/slices rather than panics.

Configuration strings are allocated with `C.CString` only for the duration needed to populate the C config and are released after client creation/config teardown according to the existing package pattern. The Go API owns its input strings; the C SDK must not receive pointers to Go memory.

## Testing and verification

Add focused tests before implementation for:

1. token event conversion, including nil/empty affected-channel lists;
2. new config fields being represented and passed through the C config conversion path;
3. callback registration and dispatch using `OnTokenEvent`.

After the caller has installed matching Linux MINI headers and `.so` files, run `gofmt`,
`make test-rtm`, and a Linux cgo package build/compile check. Verify header symbol references
against the installed C headers. Runtime behavior is in scope only with those matching binaries.

## Non-goals

- Copying or replacing Linux `.so` files or C headers in this repository.
- Changing Mac RTM headers, libraries, or cgo behavior.
- Retaining the removed `OnConnectionStateChanged` Go API.
- Unrelated RTM or RTC refactors.
