# Synchronous video frame observer

The approved change replaces `VideoFrameObserver.OnFrame` and the native
`video_frame_observer2` implementation with all ten callbacks from
`agora_sdk/include/c/api2/agora_local_user.h:video_frame_observer`.

## API and bridge

Expose `OnCaptureVideoFrame`, `OnPreEncodeVideoFrame`,
`OnMediaPlayerVideoFrame`, `OnRenderVideoFrame`, `OnTranscodedVideoFrame`,
`GetVideoFrameProcessMode`, `GetVideoFormatPreference`, `GetRotationApplied`,
`GetMirrorApplied`, and `GetObservedFramePosition`. Render UIDs are `uint32`;
source types and media player IDs are `int`. Frame callbacks return `bool`.
Use typed process-mode and observation-position constants and the existing
`VideoPixelFormat` type. Correct `VideoFrame.Type` to `VideoPixelFormat`, matching
the native field. Remove the old UID-zero workaround.

Missing frame callbacks accept frames. Missing configuration callbacks default
to read-only processing, I420 (matching existing Go YUV consumers), no rotation
or mirroring, and post-capture plus pre-render observation. Nil native frames
are rejected. Unknown handles never invoke application callbacks.

Keep frame buffers copied into Go memory so applications can retain frames.
All frame callbacks, including media-player callbacks, operate on Go copies.
Frame edits remain local to Go and are not written back to native storage.
`GetVideoFrameProcessMode` only forwards the application callback's return value
to the SDK, defaulting to read-only. Do not cache that value on the connection
or implement processing-mode behavior in this change.

## Lifecycle

Create/destroy with `agora_video_frame_observer_create/destroy`, register with
`agora_local_user_set_video_frame_observer`, and unregister with
`agora_local_user_unset_video_frame_observer`. Publish the handle mapping before
registration because the SDK may synchronously query configuration. Return
native error codes; creation or registration failure leaves no new observer
or handle mapping. Do not destroy a registered observer when unset fails.
Keep `unregisterVideoFrameObserver` private. The connection owns cleanup during
replacement, disconnect, and release; developers only register observers and
cannot explicitly unregister them. Re-registering the same object is a no-op.
Release unsets the new observer and destroys it after connection teardown.

## Migration and verification

Migrate the four in-repository observer consumers to `OnRenderVideoFrame` and
numeric UIDs. Tests exercise the actual C function-pointer table, argument and
return conversion, safe defaults, retained frame memory, isolation of Go edits,
processing-mode return-value forwarding, and native
lifecycle failures using an explicitly tagged native test double.

The checked-in Linux SDK exports all required symbols. The checked-in macOS
SDK exports create/destroy but lacks set/unset; it cannot run this new API
without a matching SDK update. Do not silently fall back to observer2 or add
production registration stubs. Verify compilation separately from tagged
bridge/lifecycle tests and report this runtime limitation. Existing unrelated
vet errors are baseline failures and should be reported separately.
