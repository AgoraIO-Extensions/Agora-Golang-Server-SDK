# Audio AED UAP Binding Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add a cross-platform Go wrapper for `aed.h`, with a real Linux `vad_uap` implementation and a safe non-Linux/no-tag stub.

**Architecture:** Keep public AED data types in a common file. Put cgo lifecycle and frame processing in `audio_aed_uap.go` guarded by `linux && vad_uap`; put matching unavailable behavior in `audio_aed_stub.go` guarded by `!linux || !vad_uap`. Inputs are caller-provided float PCM and power spectrum slices, and outputs are copied into Go-owned fields.

**Tech Stack:** Go 1.20, cgo, checked-in Agora UAP AED C ABI, Go tests.

---

### Task 1: Define the cross-platform AED contract

**Files:**
- Create: `go_sdk/rtc/audio_aed_types.go`
- Create: `go_sdk/rtc/audio_aed_stub.go`
- Test: `go_sdk/rtc/audio_aed_test.go`

- [ ] Write tests for stub construction returning `ErrVadUAPNotEnabled`, idempotent release, and invalid input rejection.
- [ ] Run `go test ./go_sdk/rtc` without tags and verify the new tests fail because AED symbols/types are absent.
- [ ] Add `AudioAEDConfig`, `AudioAEDInput`, `AudioAEDOutput`, `AudioAED`, and shared validation/error definitions; implement unavailable stub methods.
- [ ] Run `go test ./go_sdk/rtc` and verify the no-tag/macOS path passes without linking AED.

### Task 2: Implement Linux AED lifecycle and processing

**Files:**
- Create: `go_sdk/rtc/audio_aed_uap.go`
- Test: `go_sdk/rtc/audio_aed_linux_test.go`

- [ ] Add a Linux-tagged smoke test that constructs with defaults, processes a 160-sample frame with 513-bin power input, updates dynamic config, and releases.
- [ ] Run the tagged test and verify it fails before implementation because `NewAudioAED` is unavailable on Linux with `vad_uap`.
- [ ] Bind `aed.h`, map static config, call `create`/`memAllocate`/`init`, expose dynamic-config setter, validate `nBins == len(BinPower)` and `hopSz == len(TimeSignal)`, call `proc`, and copy every `Aed_OutputData` field.
- [ ] Run `go test -tags vad_uap ./go_sdk/rtc` with the checked-in Linux library and verify the smoke test passes.

### Task 3: Verify platform/build behavior

**Files:**
- Modify only if needed: `go_sdk/rtc/audio_aed_types.go`, `go_sdk/rtc/audio_aed_stub.go`, `go_sdk/rtc/audio_aed_uap.go`

- [ ] Run `gofmt` on all new Go files.
- [ ] Run `go test ./go_sdk/rtc` and `go test -tags vad_uap ./go_sdk/rtc` on Linux where the shared library is available.
- [ ] Run a no-tag compile with `GOOS=darwin` if the local Go toolchain supports it; confirm no AED link flags are selected.
- [ ] Inspect `git diff` and report the inability to commit if `.git/index.lock` remains unavailable.
