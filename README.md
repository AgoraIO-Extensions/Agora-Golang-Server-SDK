# Required OS and go version
- supported linux version: 
  - Ubuntu 18.04 LTS and above
- macOS version: only for coding, not for production environment
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

## How to use RTM?

1. download the RTM SDK

    ```bash
    git clone https://github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK.git
    cd Agora-Golang-Server-SDK
    ```

2. install the dependencies

    ```bash
    make deps
    ```

3. build the RTM demo

    ```bash
    ./scripts/rtmbuild.sh
    ```

4. run the RTM demo

    ```bash
    cd bin
    ./rtmdemo <appid> <channelname> <userid>
    ```

    > please replace `<appid>`、`<channelname>` and `<userid>` with your own info.

5. more RTM sample ref to：`cmd/main.go`


# Intergrate into your project
- Clone this repository and checkout to the target branch, and install. recommend to use latest version of go.
The following takes "release/2.1.0" as an example, you can use the latest version to replace release/2.1.0
and for each branch, there is a corresponding tag
for example, 'release/2.1.0' is tagged as 'v2.1.0'
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
  agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/rtc"
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

```markdown
The overall logic of external audio processor is:
APM must be enabled (enableapm=true), otherwise it will not process audio.
In this case, you can use algorithms; or you can only use VAD without processing algorithms.

How to use ExternalAudioProcessor:
  -1. How to enable audio processing algorithms?
    a. Set service configure apmModel=1;
    b. In externalAudioProcessor.Initialize, apmConfig must not be nil, and set the required algorithms to true;
    c. If you also need to use VAD, set vadConfig to non-nil as well!
  
  -2. How to use ExternalAudioProcessor with VAD only?
    a. Set service configure apmModel=1
    b. In externalAudioProcessor.Initialize(nil, inputSampleRate, inputChannels, vadConfig, observer), set apmConfig to nil, and vadConfig must not be nil!

Performance test comparison:
When printing every callback result:
VAD only: Processing 1840ms of input data takes 16ms. 115x speedup
APM+VAD+APM dump disabled: Processing 1840ms of input data takes 46ms. 40x speedup

## 2025.12.17 Release Version 2.4.3

- New: Added support for the `SendIntraRequest` method, which allows you to actively request an encoded key frame (key frame) from a remote user.

### Test Results

#### Case 1: Encoding on Web

If an intra request is sent to the remote user every 1 second, the key frame is encoded and sent approximately every 1 second as well.

```text
send intra request to user 1782469624, time 1010
key frame received, time 1001
send intra request to user 1782469624, time 1008
key frame received, time 998
send intra request to user 1782469624, time 1009
key frame received, time 1033
send intra request to user 1782469624, time 1009
key frame received, time 999
```

#### Case 2: Encoding on Android

If an intra request is sent to the remote user every 1 second, the key frame is encoded and sent approximately every 2 seconds.

```text
send intra request to user 4797, time 1002
send intra request to user 4797, time 1006
key frame received, time 1692
send intra request to user 4797, time 1010
send intra request to user 4797, time 1012
key frame received, time 2023
send intra request to user 4797, time 1010
key frame received, time 1733
send intra request to user 4797, time 1010
send intra request to user 4797, time 1005
key frame received, time 1388
send intra request to user 4797, time 1011
send intra request to user 4797, time 1011
key frame received, time 1971
```

## 2025.12.15 Release 2.4.2

- Default changed to `idleModel=true`
- Added `externalaudioprocessor` in the callback of `externalAudioProcessor`
- Added "Incremental Sending Mode" for supported scenarios

### Features of Incremental Sending Mode

1. Specifically designed for the 'tts' scenario; not recommended for other scenarios. For 'iot' with tts, usage is recommended.
2. Configurable at the connection level, meaning different connections can have different settings.
3. If a connection enables "Incremental Sending Mode", data will be sent at the configured rate during the specified period; after the time elapses, it will revert to normal sending speed mode.

### Usage

1. When initializing the `connection`, set the `scenario` parameter to either `AudioScenarioAIServer` or `AudioScenarioDefault`.
2. Add the parameter `SendExternalAudioParameters` in your `publishConfigure`:
    - `Enabled` (default: false): whether to enable incremental sending mode.
    - `SendSpeed` (range: 1~5, recommended: 2, default: 0): set the transmission speed.
    - `SendMS` (default: 0): duration for fast transmission.
    - `DeliverMuteDataForFakeAdm` (default: false): whether to send silent packets when there is no tts data.

### Product Behavior

- At the beginning of each round, a segment of data is sent quickly
- Afterwards, to send data at the normal speed,i.e 1x speed

> The SDK automatically handles "rounds" and calls a new round internally.

## 2025.11.27 Release 2.4.1
-- Added support for VAD and audio processing algorithms (background voice removal, noise suppression, echo cancellation, automatic gain control, 3A algorithms) on external audio sources. This means you can now use audio processing algorithms and VAD without joining an RTC channel.
-- Modified the return values of LocalUser to distinguish them from the SDK's native error codes.
-- Added example_externalAudioProcessor to demonstrate how to use external audio data.
-- Added Lock-free Ring Buffer implementation in utils.go.
```

## 2025.11.14 Release 2.4.0
-- Major update: RTC & RTM merged into a single SDK
-- Major update: RTM updated to version 1.0
-- Major update: Added support for algorithm modules such as AI-NS, AGC, BGHVS
## 2025.08.29 Release 2.3.3
-- update: update rtc sdk to 4.4.32.0829, which fix a bug when set parameter to do auduio dump could case audio-dump thread to leak
-- add: add best practice for audio dump usage: use connection.getparameter to set audio dump parameter
2025.08.25 Release 2.3.2 - 2025.08.25
#Merged Version: first merge version of RTC and RTM sdk
​​  RTC & RTM Integration​​: The RTC and RTM sdk are now merged into a single release. Developers who do not wish to use RTM can skip executing scripts/rtminstall.sh and  scripts/rtmbuild.sh, ensuring no RTM dependencies are included.
#Updates
​-- ​RTC SDK Update​​: Upgraded to version ​​44.32.0820​​.
-- Added ptsparameter to pushAudioFramefor setting the presentation timestamp (PTS) of audio frames.
-- Added ptsparameter to onPlaybackBeforeMixingand other audio frame callbacks.
​​-- Video Frame Enhancements​​, Added ptsparameter to pushVideoFramefor setting the presentation timestamp (PTS) of video frames.
​-- ​Sample Code Enhancements (sample/ai_send_recv_pcm.go)​​
​​  - Protocol Support​​: Added ​​v2/v4 protocols​​, which have been verified and are ready for use.
​​  - PTS Management​​: Introduced PTSAllocatorfor managing PTS allocation.
​​  - Audio Data Management​​: Added PcmRawDataManager to handle audio byte management, allowing developers to pop audio data in bytesinmslength without manual calculations.
​ -​Session Parsing​​: Added SessionParserfor PTS parsing, providing notifications for session events:
​​    - Session ID Change​​: Notifies when a previous session ends and a new session begins.
​​    - Timeout Handling​​: Notifies when a session ends due to timeout.
# 2025-07-21 release 2.3.1
-- ​​Added:​​ OnAudioVolumeIndication, which allows you to get the UID of the user currently speaking.
​​Usage:​​

a.Enable audio_volume_indication and set parameters as follows:
localUser.SetAudioVolumeIndicationParameters(intervalInMs, smooth:default=3, isvad)
b.In the OnAudioVolumeIndication callback, get the UID of the user currently speaking.


# 2025-07-16 release 2.3.0 (IMPORTANT, please upgrade)
-- update rtc sdk,fixed 2 bugs
## New Features & Improvements

1. **Audio Scenario Support**  
   - Added support for `AudioScenarioAiServer` scenario type.  
   - Default `AudioScenario` is now `AudioScenarioAiServer`.  
   - Enabled configuring different `scenario` and `profile` for connections within the same process.  
   - Introduced `publishConfig` in connection creation to set `{scenario, profile, publishAudio, publishVideo, etc.}` [1,2](@ref).

2. **Connection Enhancements**  
   Added to `connection`:  
   - Observer registration:  
     - `RegisterLocalUserObserver`  
     - `RegisterAudioFrameObserver`  
     - `RegisterVideoFrameObserver`  
     - `RegisterVideoEncodedFrameObserver`  
   - Stream control:  
     - `PublishAudio`/`UnpublishAudio`  
     - `PublishVideo`/`UnpublishVideo`  
   - Data pushing:  
     - `PushAudioPcmData`/`PushAudioEncodedData`  
     - `PushVideoFrame`/`PushVideoEncodedData`  
   - Audio interruption: `InterruptAudio`  
   - Status check: `IsPushToRTCCompleted`  
   - Callback implementation: `OnAIQoSCapabilityMissing`  
   - Metadata: `SendAudioMetaData`  
   - **Note**: Automatic `CreateDataStream` (no manual calls needed) [1,3](@ref).

3. **Internal Simplifications**  
   The following are now handled automatically:  
   - `newMediaNodeFactory`  
   - Track creation (e.g., `NewCustomAudioTrackPcm`, `NewCustomVideoTrack`)  
   - Sender initialization (e.g., `NewAudioPcmDataSender`)  
   - Observer unregistration (e.g., `unregisterAudioFrameObserver`) [1,4](@ref).

---

## Integration Workflow
```go
1. svcCfg := agoraservice.NewAgoraServiceConfig()
2. agoraservice.Initialize(svcCfg)
3. con := agoraservice.NewRtcConnection(conCfg, publishConfig)
4. // Register observers:
   con.RegisterObserver(conHandler)
   con.RegisterAudioFrameObserver(...)
   con.RegisterLocalUserObserver(...)
5. con.Connect(token, channelName, userId)
6. con.PublishAudio() // or con.PublishVideo()
7. con.PushAudioPcmData(...) // or con.PushVideoFrame(...)
8. con.Disconnect()
9. con.Release()
10. agoraservice.Release() // Call on process exit

Note 1​​:

Steps 1–2 are called ​​once​​ during process startup.
Steps 3–9 can loop for multiple connections.
Step 10 is called ​​once​​ on process exit .
​​Note 2​​:

​​For AI scenarios​​:
Use AudioScenarioAiServer (optimized for lower latency + better weak-network performance).
​​Requirement​​: Client-side SDK must use AIClient scenario (consult Agora SA for compatible SDK versions).
​​For non-AI scenarios​​:
Contact Agora Support to match scenarios to your use case.
​​Connection cleanup​​:
Observers auto-unregister on Release() (no manual unregister calls) 
### Key Changes Summary
| **Before**                     | **After 2.3.0**                     |
|-------------------------------|-------------------------------------|
| Manual `CreateDataStream`     | ✅ Automatic                        |
| Manual observer unregistration | ✅ Automatic on `Release()`         |
| Fixed per-process scenario    | ✅ Multi-scenario per process      |
| Client/Server scenario mismatch| ❗ `AIClient` mandatory for AI use |  

> ⚠️ **Critical**: Client-side SDK compatibility is required for `AudioScenarioAiServer` – verify with Agora support 

## 2025.07.02 release 2.2.10
-- remove： Removed the setting of scenario in the enableSteroEncodeMode function, as this scenario can be passed as a parameter in NewRtcConnection, making it no longer necessary to set within enableSteroEncodeMode.
## 2025.07.01 release 2.2.9
-- Added log/data/config directory for log/data and config files
-- Added: ai_server senario, and direct custom audio track
-- Supports callback when aiqos missing mode: aiqosmissing refers to when the client-side version does not support aiqos capabilities (it's the version, not the specific scenario chosen by the client-side), which needs to notify the server through aiqosmissing. The server will have OnAIQoSCapabilityMissing in this callback, where developers can return the expected scenario setting, and the SDK will automatically switch to this scenario. If the developer does not want automatic switching, they can return -1. Then the developer can call localuse.UpdateAudioTrack(senario) to switch!

-- Added connection-level scenario configuration:  
  Different connections in the same process can use distinct scenarios (e.g., one `ai_server`, another `chorus`).  
-- Add senario parameter for NewCustomAudioTrack 
-- add state in publish and unplbish
-- fix a bug in send_mp4 sample

⚠️ **Critical Note**  
Parameters for `NewcustomAudioTrack`  **must match** the `scenario` in the connection's creation api: NewRTCConnection. i.e, must set same senario paramter for both NewCustomAudioTrack and NewRTCConnection. Otherwise, the audio track's behavior will be unpredictable.

## 2025.11.11 release
-- update rtm version to 1.0
-- update rtm sample
## 2025.11.04 release 2.3.4
-- update: udpate sdk
-- add: add amp to include ns algorithm
-- update: update vad algorithm
-- add: added idlemode, which can release connection's c handle when timedout
## 2025.05.26 release 2.2.8
--  fix: fix a bug in sterom mode, the custome bitrate is not work
--  update: update mac sdk version to 4.4.32
## 2025.05.21 发布 2.2.7
-- 更新：vad v1 lib removed libc++ and libstdc++ dependence，fix a bug that the vad lib can not be loaded in some os which has not libc++ and libstdc++.
## 2025.05.19 发布 2.2.6
-- update vadv1 lib，fix a bug in parameter passing
-- remove debug information in sample_vad.go
## 2025.04.28 Release 2.2.5
-- Update rtc sdk to 4.4.32
## 2025.04.28 Release 2.2.4
-- Added EnableEncryption api, to enable encryption
-- Added onEncryptionError calback, to notify encryption error;
## 2025.04.15 发布 2.2.3
-- Added: add a Dequeue class for thread-safe and go chan，referto: send_recv_yuv_pcm.go
-- Added: add onlocalauiodstats, onremoveaudiostats, onlocalvideostats, ref to: send_recv_yuv_pcm.go
## 2025.03.26 Release 2.2.2
-- FIX：fix a bug in vad Release, if call vad.Release() for more than once, it will cause a crash
-- ADD：add a func in sample_vad.go, to test stero vad for stero pcm audio file
##  2025.03.05 Release 2.2.1
-- In the configuration, added support for dual-channel encoding for water devices. By default, this is not supported. If enabled, internal private parameters will be automatically modified to support dual-channel encoding. However, at the app layer, when setting playback parameters, channel=2 must be specified. ​OK
-- Provided the Vad v1 algorithm. ​OK
-- Modified the Vad v1 algorithm interface: Exposed control parameters such as RMS, etc. ​OK
-- Included the Vad v1 algorithm library files in the SDK package. ​OK
-- Provided the SteroVad class. ​OK
-- AudioConsumer:

Set the default buffer size to 100ms. ​OK
Added the isPushToRtcCompleted interface. ​OK
Fixed the incorrect calculation of AudioConsumer::_samples_per_channel in stereo mode. ​OK
-- VideoObserver:
Intercepted frames with userid=0 in OnFrame. This was a bug where, if the local device also sent video, OnFrame would callback with uid=0 frames, which are invalid and needed to be filtered out. ​OK
-- Documentation:
Added instructions for private deployment and proxy deployment. ​OK
## 2024.12.23 Release 2.2.0
-- Update:
  - Update the sdk version to 4.4.31.
-- Add:
  - Add the SendAudioMetaData interface in LocalUser.
  - Usage: Call it directly, with a frequency limit within 100; the length of each data is within 64 bytes.

  - Add the onAudioMetaDataResult interface in LocauUserObserver.
  - Add the domainLimit member in serviceconfigure, which is used to determine whether to limit to the case where the Domain is a url. The default is 0, indicating no limit.
-- Modify:
  - ExternalVideoFrame adds support for additional yuv colorspace, enabling the encoding of solid - color background images. It is usually used in digital human scenarios.
## 2024.12.18 Release 2.1.4
-- Change: 
  - default to support vad v2 mode
  - change default stopRecognizeCount from 20 to 50 for better experience
  - add vadConfigure parameter description
## 2024.12.11 Release 2.1.3
-- New Features:
  -AudioVadManager: Introduced to manage VAD (Voice Activity Detection) for multiple audio sources. In practice, each audio source requires its own VAD instance. To simplify the developer experience, we provide a unified interface to manage these VAD instances.
-- Changes:
  - In AudioFrameObserver::OnPlaybackAudioFrameBeforeMixing, two new parameters have been added: vadResultState and vadResultFrame, which return the VAD detection results.
  - In LocalUser::RegisterAudioFrameObserver, two new parameters have been added to control whether VAD recognition should be enabled. If VAD recognition is enabled, the recognition results will be returned in OnPlaybackAudioFrameBeforeMixing.
-- VAD Usage from Version 2.1.3 onwards:
To enable VAD recognition, pass the VAD configuration when calling LocalUser::RegisterAudioFrameObserver. Example:
vadConfigure := &agoraservice.AudioVadConfigV2{
  PreStartRecognizeCount: 16,
  StartRecognizeCount:    30,
  StopRecognizeCount:     20,
  ActivePercent:          0.7,
  InactivePercent:        0.5,
  StartVoiceProb:         70,
  StartRms:               -50.0,
  StopVoiceProb:          70,
  StopRms:                -50.0,
}
localUser.RegisterAudioFrameObserver(audioObserver, 1, vadConfigure)
In the AudioFrameObserver::OnPlaybackAudioFrameBeforeMixing callback, you can retrieve the VAD recognition results and process them as follows:
you_audio_frame_observer::OnPlaybackAudioFrameBeforeMixing(channelId string, uid uint64, frame *agoraservice.AudioFrame, vadResultState agoraservice.VadState, vadResultFrame *agoraservice.AudioFrame) {
  if vadResultState == agoraservice.VadStateStartSpeaking {
    // Start speaking: process vadResultFrame for ASR/STT
  }
  if vadResultState == agoraservice.VadStateSpeaking {
    // Speaking: process vadResultFrame for ASR/STT
  }
  if vadResultState == agoraservice.VadStateStopSpeaking {
    // Stop speaking: process vadResultFrame for ASR/STT and handle end of speech business logic
  }
}
Notes:
  If VAD is enabled, always use the vadResultFrame data instead of the frame data. Using the frame data will result in lost recognition results for ASR/STT.
## Release 2.1.2 on December 2, 2024
- The AudioConsumer has been added, which can be used to push audio data.
- The VadDump method has been added, which can be used to debug VAD.
- When the profile of the service is chorus:
  - In localuser, setSenddelayinms will be automatically added.
  - In audiotrack, the track will be automatically set to chorus mode.
  - Developers no longer need to manually set setSenddelayinms and set the track to chorus mode in AI scenarios.
- The global mutex in the service has been modified to the sync.map mode to split the granularity of the mutex.
- The audioloopback sample has been added for audio loopback testing.
- The sample has been modified to provide a mode of entering the appid and channelname via the command line.
- Support for the AudioLable plugin has been added, and developers no longer need to call EableExtension at the app layer.
- The onpublishstatechanged interface has been added.
- The return status of VAD has been modified to be unified into three states: NoSpeakong, Speaking, StopSpeaking; moreover, when it is in the StopSpeaking state, the current frame data will also be returned.

## 2024.10.29 release 2.1.1
- Add audio VAD interface of version 2 and corresponding example.
## 2024.10.24 release 2.1.0
- Fixed some bug.
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


### Common Usage Q&A
## The relationship between service and process?
- A process can only have one service, and the service can only be initialized once.
- A service can only have one media_node_factory.
- A service can have multiple connections.
- Release media_node_factory.release() and service.release() when the process exits.
## If using Docker with one user per Docker, when the user starts Docker and logs out, how should Docker be released?
- In this case, create service/media_node_factory and connection when the process starts.
- Release service/media_node_factory and connection when the process exits, ensuring that...
## If Docker is used to support multiple users and Docker runs for a long time, what should be done?
- In this case, we recommend using the concept of a connection pool.
- Create service/media_node_factory and a connection pool (only new connections, without initialization) when the process starts.
- When a user logs in, get a connection from the connection pool, initialize it, execute con.connect() and set up callbacks, and then join the channel.
- Handle business operations.
- When a user logs out, execute con.disconnect() and release the audio/video tracks and observers associated with the connection, but do not call con.release(); then put the connection back into the connection pool.
- When the process exits, release the connection pool (release each con.release()), service/media_node_factory, and the connection pool (release each con.release()) to ensure resource release and optimal performance.
## Use of VAD
# Source code: voice_detection.py
# Sample code: example_audio_vad.py
# It is recommended to use VAD V2 version, and the class is: AudioVadV2; Reference: voice_detection.py.
# Use of VAD:
  1. Call _vad_instance.init(AudioVadConfigV2) to initialize the vad instance. Reference: voice_detection.py. Assume the instance is: _vad_instance
  2. In audio_frame_observer::on_playback_audio_frame_before_mixing(audio_frame):

  3. Call the process of the vad module: state, bytes = _vad_instance.process(audio_frame)
Judge the value of state according to the returned state, and do corresponding processing.

    A. If state is _vad_instance._vad_state_startspeaking, it indicates that the user is "starting to speak", and speech recognition (STT/ASR) operations can be started. Remember: be sure to pass the returned bytes to the recognition module instead of the original audio_frame, otherwise the recognition result will be incorrect.
    B. If state is _vad_instance._vad_state_stopspeaking, it indicates that the user is "stopping speaking", and speech recognition (STT/ASR) operations can be stopped. Remember: be sure to pass the returned bytes to the recognition module instead of the original audio_frame, otherwise the recognition result will be incorrect.
    C. If state is _vad_instance._vad_state_speaking, it indicates that the user is "speaking", and speech recognition (STT/ASR) operations can be continued. Remember: be sure to pass the returned bytes to the recognition module instead of the original audio_frame, otherwise the recognition result will be incorrect.
# Note: 
  If the vad module is used and it is expected to use the vad module for speech recognition (STT/ASR) and other operations, then be sure to pass the returned bytes to the recognition module instead of the original audio_frame, otherwise the recognition result will be incorrect.
# How to better troubleshoot VAD issues: It includes two aspects, configuration and debugging.
  1. Ensure that the initialization parameters of the vad module are correct. Reference: voice_detection.py.
  2. In state, bytes = on_playback_audio_frame_before_mixing(audio_frame):

    - A . Save the data of audio_frame to a local file, reference: example_audio_pcm_send.py. This is to record the original audio data. For example, it can be named: source_{time.time()*1000}.pcm
    - B.Save the result of each vad processing:

      - a When state == start_speaking: create a new binary file, for example, named: vad_{time.time()*1000}.pcm, and write bytes to the file.
      - b When state == speaking: write bytes to the file.
      - c When state == stop_speaking: write bytes to the file and close the file.
    Note: In this way, problems can be troubleshot based on the original audio file and the audio file processed by vad. This function can be disabled in the production environment.

## Debugging of VAD
- Use the utility class VadDump. Refer to examples/sample_vad, which can help troubleshoot problems.
- This method will generate three types of files:
- sourec.pcm: The original information of the audio.
- vad_{index}.pcm: The audio information after VAD processing, where index is an integer that increments from 0.
- label.txt: The label information of the audio.
Note: If there are any issues with VAD, please send these files to the Agora team for troubleshooting.
## How to Push the Audio Generated by TTS into the Channel?
    localuser.unpublish_audio(audio_track)
    localuser.unpublish_video(video_track)
    audio_track.set_enabled(0)
    video_track.set_enabled(0)

    localuser.unregister_audio_frame_observer()
    localuser.unregister_video_frame_observer()
    localuser.unregister_local_user_observer()

    connection.disconnect()
    connection.unregister_observer()

    localuser.release()
    connection.release()

    
    audio_track.release()
    video_track.release()
    pcm_data_sender.release()
    video_data_sender.release()
    audio_consumer.release()

    media_node_factory.release()
    agora_service.release()
    
    #set to None
    audio_track = None
    video_track = None
    audio_observer = None
    video_observer = None
    local_observer = None
    localuser = None
    connection = None
    agora_service = None

### How to push the audio generated by TTS into the channel?
  # Source code: audio_consumer.py
  # Sample code: example_audio_consumer.py
### ## Interrupt Handling in AI Scenarios
# Definition of Interrupt
In human-machine dialogue, an interrupt refers to the situation where a user suddenly interrupts the robot's response, requesting the robot to stop its current response immediately and shift to answer the user's new question. This behavior is called an interrupt.

# Trigger Conditions for Interrupts
Interrupts can be defined in different ways depending on the product. There are generally two modes:

- Mode 1: Voice Activation Mode
When it detects that the user is speaking, the interrupt strategy is triggered. For example, when the system recognizes speech, it triggers the interrupt strategy to stop the robot's response.

- Mode 2: ASR Activation Mode
When the system detects that the user is speaking and receives a result from ASR (Automatic Speech Recognition) or STT (Speech-to-Text), the interrupt strategy is triggered.

# Advantages of Different Interrupt Strategies
Voice Activation Interrupt

Advantages:
Reduces the user's wait time and the likelihood of interrupts, as the robot will stop its response immediately when the user starts speaking, eliminating the need for the user to wait for the robot to finish speaking.
Disadvantages:
Since this is voice-activated, it may be triggered by meaningless audio signals, depending on the accuracy of the VAD (Voice Activity Detection). For example, if someone is typing on the keyboard while the AI is speaking, it might trigger the interrupt incorrectly.
ASR Activation Interrupt

Advantages:
Reduces the probability of unnecessary interrupts because the interrupt strategy is triggered only after ASR or STT has recognized the user’s speech.
Disadvantages:
Since this is ASR/STT-triggered, it requires converting the audio signal into text, which introduces a delay before the interrupt can be processed.
- Recommended Mode
If the VAD can filter out non-speech signals and only triggers when human speech is detected, the Voice Activation Mode is recommended. This mode is also suitable when the delay in processing the interrupt is not a major concern.

If the interrupt delay is not sensitive, the ASR Activation Mode is recommended. This mode can filter out non-speech signals more effectively and reduce the probability of an unintended interrupt.

How to Implement Interrupts? What Actions Are Required?
In a human-machine dialogue system, conversations are typically structured in "rounds," where each round consists of a question from the user, followed by a response from the robot, and so on. For each round, we can assign a roundId, incrementing it with each new round. A round consists of the following stages:

VAD (Voice Activity Detection):
This marks the start of the dialogue, where the system detects the beginning and end of the user's speech. It then passes this information to the ASR for further processing.

ASR (Automatic Speech Recognition):
This phase involves recognizing the user's speech and converting it into text, which is then passed to the LLM (Large Language Model).

LLM (Large Language Model):
This is the generation phase, where the LLM processes the recognized user input and generates a response.

TTS (Text-to-Speech):
In this phase, the LLM’s response is converted into an audio format.

RTC Streaming:
The generated audio is streamed via RTC (Real-Time Communication) to be played back to the user.

Therefore, an interrupt happens when, in the next round (roundId+1), either through Voice Activation (triggered by the VAD phase) or ASR Activation (triggered when ASR recognizes the user’s speech), the following actions must be performed:

Stop the LLM Generation in the current round (roundId).
Stop the TTS Synthesis in the current round (roundId).
Stop the RTC Streaming in the current round (roundId).
API Call References:
Call: AudioConsumer.clear()
Call: LocalAudioTrack.clear_sender_buffer()
Business Layer: Clear any remaining TTS-related data (if applicable)


## When to Pass LLM Results to TTS for Synthesis?
LLM (Large Language Model) results are returned asynchronously and in a streaming manner. When should the results from the LLM be passed to TTS (Text-to-Speech) for synthesis?

Two main factors need to be considered:

Ensure that the TTS synthesized speech is unambiguous:
The speech synthesized by TTS must be clear, complete, and continuous. For example, if the LLM returns the text: "中间的首都是北京吗？", and we pass it to TTS as:

"中",
"国首",
"是北",
"京吗？",
This would result in ambiguous synthesis because there are no spaces between certain words (e.g., between "中" and "国", "首" and "是", and "京" and "吗"). Proper segmentation must be ensured to avoid such ambiguities.
Minimize overall processing delay:
If the LLM results are passed to TTS only after the entire response is generated, the speech synthesis will be unambiguous and continuous. However, this approach introduces significant delay, which negatively affects the user experience.

Recommended Approach
To achieve a balance between clarity and minimal delay, the following steps should be followed:

Store the LLM results in a cache as they are received.
Perform a reverse scan of the cached data to find the most recent punctuation mark.
Truncate the data from the start to the most recent punctuation mark and pass it to TTS for synthesis.
Remove the truncated data from the cache. The remaining data should be moved to the beginning of the cache and continue waiting for additional data from the LLM.

##VAD Configuration Parameters
AgoraAudioVadConfigV2 Properties

Property Name	Type	Description	Default Value	Value Range
preStartRecognizeCount	int	Number of audio frames saved before detecting speech	16	[0, ]
startRecognizeCount	int	Total number of audio frames to detect speech start	30	[1, max]
stopRecognizeCount	int	Number of audio frames to detect speech stop	50	[1, max]
activePercent	float	Percentage of active frames in startRecognizeCount frames	0.7	[0.0, 1.0]
inactivePercent	float	Percentage of inactive frames in stopRecognizeCount frames	0.5	[0.0, 1.0]
startVoiceProb	int	Probability that an audio frame contains human voice	70	[0, 100]
stopVoiceProb	int	Probability that an audio frame contains human voice	70	[0, 100]
startRmsThreshold	int	Energy dB threshold for detecting speech start	-50	[-100, 0]
stopRmsThreshold	int	Energy dB threshold for detecting speech stop	-50	[-100, 0]
Notes:
startRmsThreshold and stopRmsThreshold:

The higher the value, the louder the speaker's voice needs to be compared to the surrounding background noise.
In quiet environments, it is recommended to use the default value of -50.
In noisy environments, you can increase the threshold to between -40 and -30 to reduce false positives.
Adjusting these thresholds based on the actual use case and audio characteristics can achieve optimal performance.
stopRecognizeCount:

This value reflects how long to wait after detecting non-human voice before concluding that the user has stopped speaking. It controls the gap between consecutive speech utterances. Within this gap, VAD will treat adjacent sentences as part of the same speech.
A shorter gap will increase the likelihood of adjacent sentences being recognized as separate speech segments. Typically, it is recommended to set this value between 50 and 80.
For example: "Good afternoon, [interval_between_sentences] what are some fun places to visit in Beijing?"

If the interval_between_sentences between the speaker's phrases is greater than the stopRecognizeCount, the VAD will recognize the above as two separate VADs:

VAD1: Good afternoon
VAD2: What are some fun places to visit in Beijing?
If the interval_between_sentences is less than stopRecognizeCount, the VAD will recognize the above as a single VAD:

VAD: Good afternoon, what are some fun places to visit in Beijing?



If latency is a concern, you can lower this value, or consult with the development team to determine how to manage latency while ensuring semantic continuity in speech recognition. This will help avoid the AI being interrupted too sensitively.

Usage of Dual-Channel Encoding: Not Recommended Unless Necessary! If Needed, Please Consult with R&D for Confirmation!
Applicable Scenarios
Scenarios where the client must have dual-channel audio, meaning the data in the left and right channels must be different. Remember: Not recommended unless the left and right channel data are definitely different!

Client Usage
Refer to the documentation.

Server Usage
​Configure in serviceConfigure:
Set audioScenario to gameStreaming and enable dual-channel encoding.

go
svcCfg := agoraservice.NewAgoraServiceConfig()
svcCfg.AppId = appid
// Change audio scenario
svcCfg.AudioScenario = agoraservice.AudioScenarioGameStreaming
svcCfg.EnableSteroEncodeMode = 1

agoraservice.Initialize(svcCfg)
​Set Callback Parameters to Dual-Channel:
Since dual-channel VAD only supports a 16k sample rate, set the callback parameters to dual-channel with a 16k sample rate.

go
localUser.SetPlaybackAudioFrameBeforeMixingParameters(2, 16000)
​Disable VAD in audioFrameObserver:

go
localUser.RegisterAudioFrameObserver(audioObserver, 0, nil)
​Use SteroVad for Dual-Channel VAD Check in audioFrameObserver Callback:
​4.1 Initialize StereoVad (Before Callback):

go
// VAD v1 for stereo
vadConfigV1 := &agoraservice.AudioVadConfig{
    StartRecognizeCount:    10,
    StopRecognizeCount:     6,
    PreStartRecognizeCount: 10,
    ActivePercent:          0.6,
    InactivePercent:        0.2,
    // Other parameters can be adjusted as needed or left as default
}
// Generate stereo VAD
steroVadInst := agoraservice.NewSteroVad(vadConfigV1, vadConfigV1)
​4.2 Dual-Channel VAD Check:
Perform dual-channel VAD check in the audioFrameObserver callback.

go
audioObserver := &agoraservice.AudioFrameObserver{
    OnPlaybackAudioFrameBeforeMixing: func(localUser *agoraservice.LocalUser, channelId string, userId string, frame *agoraservice.AudioFrame, vadResultState agoraservice.VadState, vadResultFrame *agoraservice.AudioFrame) bool {
        // Do something...

        // Perform stereo VAD processing
        start := time.Now().Local().UnixMilli()
        leftFrame, leftState, rightFrame, rightState := steroVadInst.ProcessAudioFrame(frame)
        end := time.Now().UnixMilli()
        leftLen := 0
        if leftFrame != nil {
            leftLen = len(leftFrame.Buffer)
        }
        rightLen := 0
        if rightFrame != nil {
            rightLen = len(rightFrame.Buffer)
        }
        fmt.Printf("Left VAD state %d, left len %d, right VAD state %d, right len: %d, diff = %d\n", leftState, leftLen, rightState, rightLen, end-start)

        // Dump VAD frame: Only for debug, remove this when releasing
        dumpSteroVadResult(1, leftFrame, leftState)
        dumpSteroVadResult(0, rightFrame, rightState)

        // Do something...
        // Return
        return true
    }
}
​Release:
Must be called after agoraservice::Release() to take effect. Example code:

go
agoraservice.Release()
if steroVadInst != nil {
    steroVadInst.Release()
}
steroVadInst = nil
Difference Between Proxy and Private Deployment
​Proxy: Addresses network restrictions in LAN environments. Through the proxy, devices in the LAN can access Agora's RTC services via the public network.
​Private Deployment: Deploys Agora's RTC services on your own servers. Private deployment allows better control over Agora's RTC services and enhances user privacy protection.
How to Support Private Deployment
​Call Sequence:
​***Must be called after creating the connection but before con.Connect to take effect. If called after agoraservice::Init(), it will not take effect.

​API Call:

go
connectionCon := agoraservice.NewRtcConnection(&conCfg)
localUser := con.GetLocalUser()

// Test LAN
params1 := `{"rtc.enable_nasa2": false}`
params2 := `{"rtc.force_local": true}`
params3 := `{"rtc.local_domain": "ap.1452738.agora.local"}`
params4 := `{"rtc.local_ap_list": ["20.1.125.55"]}`

// Set parameters
agoraservice.GetAgoraParameter().SetParameters(params1)
agoraservice.GetAgoraParameter().SetParameters(params2)
agoraservice.GetAgoraParameter().SetParameters(params3)
agoraservice.GetAgoraParameter().SetParameters(params4)
How to Support Proxy
​Call Sequence:
Call after agoraService::Init() (but it is not verified whether it actually takes effect).

​API Call:

go
parameter.setBool("rtc.enable_proxy", true);
s->setBool("rtc.force_local", true);
s->setBool("rtc.local_ap_low_level", true);

s->setBool("rtc.enable_nasa2", true);
s->setParameters("{\"rtc.vos_list\":[\"10.62.0.95:4701\"]}");
s->setParameters("{\"rtc.local_ap_list\":[\"10.62.0.95\"]}");
