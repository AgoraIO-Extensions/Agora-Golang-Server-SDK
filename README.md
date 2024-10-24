# Required OS and go version
- supported linux version: 
  - Ubuntu 18.04 LTS and above
  - CentOS 7.0 and above
- go version:
  - go 1.21 and above
  - not tested on go 1.19 and below
- examples require:
  - ffmpeg version 6.x and above, and corresponding development library installed (libavformat-dev, libavcodec-dev, libswscale-dev...)
  - pkg-config installed

# Build and run examples
```
make deps
make build
make examples
cd ./bin
export LD_LIBRARY_PATH=../agora_sdk
# for mac, use command bellow
# export DYLD_LIBRARY_PATH=../agora_sdk_mac
./multi_cons_rtx -h
```

# Intergrate into your project
- Clone this repository and checkout to the target branch, and install
```
git clone git@github.com:AgoraIO-Extensions/Agora-Golang-Server-SDK.git
cd Agora-Golang-Server-SDK
git checkout dev/2.0
make install
```
- Add the following dependency to your go.mod
```
replace github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2 => /path/to/Agora-Golang-Server-SDK

require github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2 v2.0.4
```
- Add import in your go file
```
import (
  agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/agoraservice"
)
```
- Call agoraservice interface in your code
```
	svcCfg := agoraservice.NewAgoraServiceConfig()
	svcCfg.AppId = appid
	agoraservice.Initialize(svcCfg)
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
  - if you don't use VAD, and your glibc version is between 2.16 and 2.27, you can disable VAD by rename **audio_vad.go** file in go_sdk/agoraserver/ to **audio_vad.go.bak**

# Change log
## 2024.10.24 release 2.1.0
- Fixed some bug
## 2024.10.12 release 2.0.0
- Simplify SDK install.
- Make golang SDK interfaces consistent with agora SDK interfaces of other program languages.
## 2024.08.14 release 1.1
- Add RenewToken interface for RtcConnection
- Update VAD algorithm
## 2024.06.28
- Support build and test on mac
## 2024.06.25
- Make VAD available, for details see UT case TestVadCase in agoraservice/rtc_connection_test.go
- Reduce audio delay, by setting AudioScenario to AudioScenarioChorus, for details see UT case TestBaseCase in agoraservice/rtc_connection_test.go
- Add interface for adjusting publish audio volume with examples in sample.go
- Solve the problem of noise in the received audio on the receiver side, by setting the audio send buffer and ensuring that the audio data is sent earlier than the actual time it should be sent. For details, see sample.go
