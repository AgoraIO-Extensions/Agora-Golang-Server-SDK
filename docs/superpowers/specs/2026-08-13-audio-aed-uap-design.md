# Audio AED UAP Binding Design

## Goal

Add a Go binding for the newer `agora_sdk/include/c/api2/aed.h` interface without changing the existing `AudioVad`/`audio_vad_uap.go` behavior.

## Platform and Build Tags

- `go_sdk/rtc/audio_aed_uap.go` is built only for `linux && vad_uap` and links `libagora_uap_aed.so`.
- `go_sdk/rtc/audio_aed_stub.go` is built for every other platform/tag combination. It exposes the same public types and methods but reports `ErrVadUAPNotEnabled` and never references the AED header or library.
- macOS therefore remains buildable and runnable without an AED `.dylib`.

## Public API

The new binding uses distinct names so it cannot be confused with the legacy VAD wrapper:

```go
type AudioAED struct { ... }

func NewAudioAED(cfg *AudioAEDConfig) (*AudioAED, error)
func (aed *AudioAED) Release()
func (aed *AudioAED) Process(input *AudioAEDInput) (*AudioAEDOutput, error)
```

`AudioAEDConfig` maps the static AED configuration (`Enable`, `FFTSz`, `HopSz`, `AnaWindowSz`, `FreqInputAvailable`, `UseCVersionAI`). Dynamic configuration is set through an optional `SetDynamicConfig` method and defaults come from the C library.

`AudioAEDInput` accepts caller-owned, precomputed `BinPower []float32` and the current PCM `TimeSignal []float32`. The wrapper passes these buffers to C only for the duration of `Process`; callers may reuse or mutate them after the method returns. The input also includes `NBins`/`HopSz` validation through slice lengths and optional explicit values are avoided to prevent contradictory metadata.

`AudioAEDOutput` copies all fields from `Aed_OutputData` into Go-owned memory, including energy, RMS, voice/music probabilities, binary decisions, and pitch frequency. No processed PCM frame is synthesized because AED does not return PCM.

## Lifecycle and Errors

The Linux implementation follows the required C lifecycle:

1. `create`
2. `memAllocate` with static config
3. `init`
4. optional dynamic-config update
5. `proc` for each frame
6. `destroy` during `Release`

Every non-zero C return is surfaced as a Go error. `Process` rejects a nil/released instance, missing or mismatched input buffers, and unsupported FFT/hop dimensions before entering C. `Release` is idempotent.

The stub returns `ErrVadUAPNotEnabled` from construction and processing, while `Release` remains safe. The existing `ErrVadUAPNotEnabled` variable is reused.

## Testing

- Add platform-independent tests for input validation, output field copying helpers, and stub behavior when built without Linux AED support.
- Compile the package with no tags on macOS-compatible settings to prove no AED symbols are required.
- On Linux with `-tags vad_uap`, build against the checked-in shared library and run a smoke test covering construction, one valid frame, dynamic config, and release.
