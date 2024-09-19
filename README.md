# Build and run on linux
## Required OS and go version
- supported linux version: 
  - Ubuntu 18.04 LTS and above
  - CentOS 7.0 and above
- go version:
  - go 1.20 and above
  - not tested on go 1.19 and below

## Prepare C version of agora rtc sdk
- download and unzip [agora_sdk.zip](https://download.agora.io/sdk/release/agora_rtc_sdk_linux_20240902_320567.zip)
```
unzip agora_sdk.zip
```
- make **agora_sdk** directory in the same directory with **go_wrapper**
- there should be **libagora_rtc_sdk.so** and **include_c** in **agora_sdk** directory

## Build sample
```
cd go_wrapper/examples/send_recv_pcm
go mod tidy
go build main
```

## Test
- download and unzip [test_data.zip](https://download.agora.io/demo/test/test_data_202409021506.zip)
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
- download and unzip [agora_sdk_mac.zip](https://download.agora.io/sdk/release/agora_rtc_sdk_mac_20240902_320567.zip)
```
unzip agora_sdk_mac.zip
```
- make **agora_sdk_mac** directory in the same directory with **go_wrapper**
- download and unzip [test_data.zip](https://download.agora.io/demo/test/test_data_202409021506.zip)
- make **test_data** directory in the same directory with **go_wrapper**
- build and run main
```
cd go_wrapper/examples/send_recv_pcm
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

# FAQ
## compile error
### undefined symbol reference to GLIBC_xxx
- libagora_rtc_sdk depends on GLIBC 2.16 and above
- libagora_uap_aed depends on GLIBC 2.27 and above
- solutions are:
  - you can upgrade your glibc if possible, or you will need to upgrade your running system to **required os version**
  - if you don't use VAD, and your glibc version is between 2.16 and 2.27, you can disable VAD by rename **audio_vad.go** file in go_wrapper/agoraserver/ to **audio_vad.go.bak**

# Change log
## 2024.09.19 release 1.3
- fix possible crash
## 2024.09.13 release 1.2.2
- there is no need to config ffmpeg path for **send_mp4** example any more
- make video pixel format support for other than yuv420p
## 2024.09.04 release 1.2.1
- move examples to **go_wrapper/examples** directory
- add **send_mp4** example combined with ffmpeg
  - make sure ffmpeg installed if you build **send_mp4** example
  - and replace **cgo flags** of ffmpeg path in send_mp4.go
## 2024.09.02 release 1.2
- libuap_aed.so library update
- Add parameters to AudioVadConfig
- Add AreaCode to AgoraServiceConfig
## 2024.08.14 release 1.1
- Add RenewToken interface for RtcConnection
- Update VAD algorithm
## 2024.06.28
- Support build and test on mac
## 2024.06.25
- Make VAD available, for details see UT case TestVadCase in agoraservice/rtc_connection_test.go
- Reduce audio delay, by setting AudioScenario to AUDIO_SCENARIO_CHORUS, for details see UT case TestBaseCase in agoraservice/rtc_connection_test.go
- Add interface for adjusting publish audio volume with examples in sample.go
- Solve the problem of noise in the received audio on the receiver side, by setting the audio send buffer and ensuring that the audio data is sent earlier than the actual time it should be sent. For details, see sample.go
