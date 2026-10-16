# Video Frame Observer Implementation Plan

**Goal:** Replace observer2 with all ten synchronous observer callbacks and safe native ownership.

**Architecture:** Keep the connection-owned native handle and existing service handle map. Bridge C function pointers through `video_observer_cgo.[ch]` into Go callbacks. Return independent Go frame copies; forward processing-mode queries without caching or applying them in Go.

**Tech Stack:** Go, cgo, C, bundled Agora native SDK.

## 1. Regression coverage

- [x] Add an explicitly build-tagged native test double and Go helpers in `go_sdk/rtc/video_observer_testbridge*`.
- [x] Add `go_sdk/rtc/video_observer_test.go` covering all ten callback slots, defaults, nil frames, argument/return conversion, isolation of Go frame edits, creation failure, set/unset failure, replacement, and repeated unregister.
- [x] Run `go test -vet=off -tags agora_video_observer_test ./go_sdk/rtc -run TestVideoFrameObserver -count=1`; confirm missing new API fails before implementation.

## 2. Native and Go implementation

- [x] Replace `VideoFrameObserver` fields in `rtc_connection.go`, add process/position types in `video_consts.go`, and keep `unregisterVideoFrameObserver` private for connection-owned cleanup.
- [x] Replace old create/destroy wiring in `type_convert.go`; assign every new C callback slot.
- [x] Declare and define matching trampolines in `video_observer_cgo.h` and `video_observer_cgo.c`, using native numeric `uid_t` for rendering.
- [x] Replace the raw frame callback in `video_observer.go` with ten exported callbacks. Preserve encoded-video code.
- [x] Support copied frame conversion without modifying native storage. Forward processing-mode return values without caching or implementing mode behavior. Correct `VideoFrame.Type` to `VideoPixelFormat`.
- [x] Change register/unregister/release to set/unset, propagate native errors, and clean up failed registrations. Preserve native ownership on failed unset.

## 3. Consumers and verification

- [x] Migrate `send_recv_yuv`, `send_recv_yuv_pcm`, `send_h264`, and `multi_cons_rtx` examples to numeric `OnRenderVideoFrame`.
- [x] Format new Go code with gofmt, preserve unrelated pre-existing formatting, and run `git diff --check`.
- [x] Run tagged regression tests and race checks. Use `go build ./go_sdk/rtc` to verify production compilation independently of unavailable macOS symbols.
- [x] Check for remaining observer2/OnFrame references in changed paths. Review ownership, error paths, ABI signatures, and example UID formatting.
- [x] Document test evidence, baseline vet failures, and the macOS native SDK mismatch in both READMEs and this plan.

## Verification results

- Regression tests failed on the missing callback fields before implementation.
- Scope revision tests reproduced native-frame mutations before removal of processing-mode behavior. They now verify isolated Go copies across all five frame callbacks and both reported native modes, including media-player frames, plus direct forwarding of processing-mode callback results.
- `go test -vet=off -race -tags agora_video_observer_test ./go_sdk/rtc -run TestVideoFrameObserver -count=1` passes, including frame format conversion and Release with failed unset.
- `go build ./go_sdk/rtc` passes. All four migrated examples build and link with the explicit `agora_video_observer_test` tag.
- Production example linking fails on precisely the two missing macOS set/unset symbols. The bundled Linux library exports both, but Linux runtime integration was not exercised on this macOS host.
- Baseline `go test ./go_sdk/rtc` fails vet for a misplaced build comment and existing printf-format issues in `agora_utils.go`, `rtc_callbacks.go`, and `rtc_connection.go`. Regression runs disable vet; these unrelated failures were left unchanged.
- Independent review found no important issues. No commits or external publication were performed.
