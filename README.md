# Build and run on linux
## Required OS and go version
- supported linux version: 
  - Ubuntu 18.04 LTS and above
  - CentOS 7.0 and above
- go version:
  - go 1.20 and above
  - not tested on go 1.19 and below

## Prepare C version of agora rtc sdk
- download and unzip [agora_sdk.zip](https://share.weiyun.com/1tuBWw6O)
```
unzip agora_sdk.zip
```
- make **agora_sdk** directory in the same directory with **go_wrapper**
- there should be **libagora_rtc_sdk.so** and **include_c** in **agora_sdk** directory

## Build sample
```
cd go_wrapper
go build main
```

## Test
- download and unzip [test_data.zip](https://share.weiyun.com/d3BuJNkZ)
- make **test_data** directory in the same directory with **go_wrapper**
- run main
```
export LD_LIBRARY_PATH=/path/to/agora_sdk
cd go_wrapper
./main
```
- run ut on linux
```
export LD_LIBRARY_PATH=/path/to/agora_sdk
cd go_wrapper/agoraservice
go test -v -count=1 -timeout 20s -run ^TestBaseCase$ agoraservice
```

# Build and run on mac
- download and unzip [agora_sdk_mac.zip](https://share.weiyun.com/jgvFzRI0)
```
unzip agora_sdk_mac.zip
```
- make **agora_sdk_mac** directory in the same directory with **go_wrapper**
- download and unzip [test_data.zip](https://share.weiyun.com/d3BuJNkZ)
- make **test_data** directory in the same directory with **go_wrapper**
- build and run main
```
cd go_wrapper
./build_for_mac.sh
./main
```
- build and test ut
```
export CGO_LDFLAGS_ALLOW="-Wl,-rpath,.*"
cd go_wrapper/agoraservice
go test -v -count=1 -timeout 20s -run ^TestBaseCase$ agoraservice
```

# vad usage
- call NewAudioVad to create a new vad instance
- call ProcessPcmFrame to process one audio frame, the frame is a slice of 16 bits,16khze and mono pcm
data,
- ProcessPcmFrame(frame *PcmInAudioFrame) (*PcmOutAudioFrame, int) is a function with a dual return - value, indicating an update of the speech state. The int parameter represents the current Voice Activity Detection (VAD) state, where:

- 0 represents no speech detected;
- 1 represents the start of speech;
- 2 represents speech in progress;
- 3 represents the end of the current speech segment;
- -1 represents an error occurred.
- If the function is in states 1, 2, or 3, the PcmOutAudioFrame will contain the PCM data corresponding to the VAD state.

- When a user wants to perform ASR/TTS processing, they should send the data from the PcmOutAudioFrame to the ASR system.
- call Release() when the vad instance is no longer needed
-- NOTE： 
- The VAD instance should be released after the ASR system is no longer needed.
- One VAD instance corresponds to one audio stream

# Change log
## 2024.06.28
- support build and test on mac
## 2024.06.25
- Make VAD available, for details see UT case TestVadCase in agoraservice/rtc_connection_test.go
- Reduce audio delay, by setting AudioScenario to AUDIO_SCENARIO_CHORUS, for details see UT case TestBaseCase in agoraservice/rtc_connection_test.go
- Add interface for adjusting publish audio volume with examples in sample.go
- Solve the problem of noise in the received audio on the receiver side, by setting the audio send buffer and ensuring that the audio data is sent earlier than the actual time it should be sent. For details, see sample.go
