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
##2024.12.11 Release 2.1.3
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
### How to release resources?
