# Required OS and go version
- supported linux version: 
  - Ubuntu 18.04 LTS and above
- go version:
  - go 1.21 and above
  - go 1.20 and below are not supported
- **advanced examples** require:
  - ffmpeg version 6.x and above, and corresponding development library installed (libavformat-dev, libavcodec-dev, libswscale-dev...)
  - pkg-config installed

# Build SDK
```
# This step download agora sdk to agora_sdk and agora_sdk_mac directory
make deps
# If no error occured, you can use sdk in your project now.
make build
```

# Build and run basic examples
- Download testing data for testing examples. You can also use your own data and modify the data path in examples.
```
curl -o test_data.zip https://download.agora.io/demo/test/server_sdk_test_data_202410252135.zip
unzip test_data.zip
```
- Build basic examples
```
make examples
```
- Run basic example
```
cd ./bin
# you should get AGORA_APP_ID and AGORA_APP_CERTIFICATE from agora console
export AGORA_APP_ID=xxxx
# if you turn on authentication for your project on agora console, then AGORA_APP_CERTIFICATE is needed.
export AGORA_APP_CERTIFICATE=xxx
# for linux, use the command bellow to load agora sdk library
export LD_LIBRARY_PATH=../agora_sdk
# for mac, use command bellow
# export DYLD_LIBRARY_PATH=../agora_sdk_mac
./send_recv_pcm
```

# Build and run advanced examples
- Install ffmpeg development library as required in the first chapter in this README
- Download test data in the same way with basic examples. If you already downloaded it, this step can be skipped.
- Build advanced examples
```
make advanced-examples
```
- Run advanced examples
```
cd ./bin
# you should get AGORA_APP_ID and AGORA_APP_CERTIFICATE from agora console
export AGORA_APP_ID=xxxx
# if you turn on authentication for your project on agora console, then AGORA_APP_CERTIFICATE is needed.
export AGORA_APP_CERTIFICATE=xxx
# for linux, use the command bellow to load agora sdk library
export LD_LIBRARY_PATH=../agora_sdk
# for mac, use command bellow
# export DYLD_LIBRARY_PATH=../agora_sdk_mac
./send_h264
```

# Intergrate into your project
- Clone this repository and checkout to the target branch, and install
```
git clone git@github.com:AgoraIO-Extensions/Agora-Golang-Server-SDK.git
cd Agora-Golang-Server-SDK
git checkout release/2.1.0
make install
```
- Add the following dependency to your go.mod
```
replace github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2 => /path/to/Agora-Golang-Server-SDK

require github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2 v2.1.0
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
- When you run your project, remember to add **agora_sdk directory** (or **agora_sdk_mac directory** for mac) path to your **LD_LIBRARY_PATH** (or **DYLD_LIBRARY_PATH** for mac) environment variable.


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
