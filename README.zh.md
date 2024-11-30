# Required OS and go version
- 支持的 Linux 版本:
  - Ubuntu 18.04 LTS 及以上版本
- Go 版本:
  - Go 1.21 及以上版本
  - Go 1.20 及以下版本不支持
- **高级示例** 需要:
  - ffmpeg 版本 6.x 及以上，并安装相应的开发库 (libavformat-dev, libavcodec-dev, libswscale-dev...)
  - 安装 pkg-config

# 构建 SDK
```
# 这一步将下载 Agora SDK 到 agora_sdk 和 agora_sdk_mac 目录
make deps
# 如果没有错误发生，现在可以在你的项目中使用 SDK。
make build
```

# 构建并运行基础示例
- 下载测试数据用于测试示例。你也可以使用自己的数据并修改示例中的数据路径。
```
curl -o test_data.zip https://download.agora.io/demo/test/server_sdk_test_data_202410252135.zip
unzip test_data.zip
```
- 构建基础示例
```
make examples
```
- 运行基础示例
```
cd ./bin
# 你需要从 Agora 控制台获取 AGORA_APP_ID 和 AGORA_APP_CERTIFICATE
export AGORA_APP_ID=xxxx
# 如果你在 Agora 控制台为项目开启了认证，则需要 AGORA_APP_CERTIFICATE。
export AGORA_APP_CERTIFICATE=xxx
# 对于 Linux，使用以下命令加载 Agora SDK 库
export LD_LIBRARY_PATH=../agora_sdk
# 对于 Mac，使用以下命令
# export DYLD_LIBRARY_PATH=../agora_sdk_mac
./send_recv_pcm
```

# 构建并运行高级示例
- 按照本 README 第一章的要求安装 ffmpeg 开发库
- 以与基础示例相同的方式下载测试数据。如果你已经下载过，可以跳过这一步。
- 构建高级示例
```
make advanced-examples
```
- 运行高级示例
```
cd ./bin
# 你需要从 Agora 控制台获取 AGORA_APP_ID 和 AGORA_APP_CERTIFICATE
export AGORA_APP_ID=xxxx
# 如果你在 Agora 控制台为项目开启了认证，则需要 AGORA_APP_CERTIFICATE。
export AGORA_APP_CERTIFICATE=xxx
# 对于 Linux，使用以下命令加载 Agora SDK 库
export LD_LIBRARY_PATH=../agora_sdk
# 对于 Mac，使用以下命令
# export DYLD_LIBRARY_PATH=../agora_sdk_mac
./send_h264
```

# 集成到你的项目
- 克隆此仓库并切换到目标分支，然后安装
```
git clone git@github.com:AgoraIO-Extensions/Agora-Golang-Server-SDK.git
cd Agora-Golang-Server-SDK
git checkout release/2.1.0
make install
```
- 在你的 go.mod 文件中添加以下依赖
```
replace github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2 => /path/to/Agora-Golang-Server-SDK

require github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2 v2.1.0
```
- 在你的 Go 文件中添加 import
```
import (
  agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/agoraservice"
)
```
- 在代码中调用 agoraservice 接口
```
  svcCfg := agoraservice.NewAgoraServiceConfig()
  svcCfg.AppId = appid
  agoraservice.Initialize(svcCfg)
```
- 运行项目时，记得将 **agora_sdk**目录 (或 Mac 上的 **agora_sdk_mac** 目录) 路径添加到 LD_LIBRARY_PATH (或 Mac 上的 DYLD_LIBRARY_PATH) 环境变量中。

# 常见问题
## 编译错误
### 未定义符号引用 GLIBC_xxx
- libagora_rtc_sdk 依赖 GLIBC 2.16 及以上版本
- libagora_uap_aed 依赖 GLIBC 2.27 及以上版本
- 解决方案:
  - 如果可能，你可以升级你的 glibc，或者你需要将运行系统升级到 **所需的操作系统版本**
  - 如果你不使用 VAD，并且你的 glibc 版本在 2.16 和 2.27 之间，你可以通过将 go_sdk/agoraserver/ 中的 **audio_vad.go** 文件重命名为 **audio_vad.go.bak** 来禁用 VAD

# 更新日志
## 2024.10.29 发布 2.1.1
- 添加V2版本的音频 VAD 接口及相应的示例。
## 2024.10.24 发布 2.1.0
- 修复了一些 bug
## 2024.10.12 发布 2.0.0
- 简化 SDK 安装。
- 使 Golang SDK 接口与其他编程语言的 Agora SDK 接口一致。
## 2024.08.14 发布 1.1
- 为 RtcConnection 添加 RenewToken 接口
- 更新 VAD 算法
## 2024.06.28
- 支持在 Mac 上构建和测试
## 2024.06.25
- 使 VAD 可用，详情见 agoraservice/rtc_connection_test.go 中的 UT 案例 TestVadCase
- 通过将 AudioScenario 设置为 AudioScenarioChorus 来减少音频延迟，详情见 agoraservice/rtc_connection_test.go 中的 UT 案例 TestBaseCase
- 添加调整发布音频音量的接口，示例见 sample.go
- 通过设置音频发送缓冲区并确保音频数据比实际应发送的时间提前发送，解决接收端接收到的音频噪声问题。

#### 常见用法Q&A
## serveice和进程的关系？
- 一个进程只能有一个service，只能对service做一次初始化；
- 一个service，只能有一个media_node_factory；
- 一个service，可以有多个connection；
- 在进程退出去的时候，再释放：media_node_factory.release() 和 service.release()

## 如果对docker的使用是一个docker一个用户，用户的时候，启动docker，用户退出去的时候，就释放docker，那么应该这么来做？
- 这个时候就在进程启动的时候，创建service/media_node_factory 和 connection；
- 在进程退出的时候，释放service/media_node_factory 和 connedtion，这样就可以保证

## 如果docker的使用是一个docker支持多个用户的时候，docker会长时间运行，应该怎么做？
- 这个情况下，我们推荐用connection pool的概念
- 在进程启动的时候，创建service/media_node_factory 和 connection pool（只是new connection，并不初始化）；
- 当有用户进来的时候，就从connection pool中获取一个connection，然后初始化，执行 con.connect()并且设置好回调，然后加入频道；
- 处理业务
- 当用户退出的时候，con.disconnect()并释放跟随该conn 的audio/video track，observer等，但不调用con.release();然后将该con 放回connection pool中；
- 在进程退出的时候，释放和 connedtion pool（对每一个con.release()
释放 service/media_node_factory 和 connedtion pool（对每一个con.release()），这样就可以保证资源的释放和性能最优




## VAD的使用
# source code: voice_detection.py
# sample code: example_audio_vad.py
- 推荐用VAD V2版本，类为： AudioVadV2； 参考：voice_detection.py；

- VAD 的使用： 
  - 1. 调用 _vad_instance.init(AudioVadConfigV2) 初始化vad实例.参考：voice_detection.py。 实例假如为： _vad_instance
  - 2. 在audio_frame_observer::on_playback_audio_frame_before_mixing(audio_frame) 中:
    - 1. 调用 vad模块的process:  state, bytes = _vad_instance.process(audio_frame)
    - 2. 根据返回的state，判断state的值，并做相应的处理
       - A. 如果state为 _vad_instance._vad_state_startspeaking，则表明当前“开始说话”，可以开始进行语音识别（STT/ASR）等操作。记住：一定要将返回的bytes 交给识别模块，而不是原始的audio_frame，否则会导致识别结果不正确。
       - B. 如果state为 _vad_instance._vad_state_stopspeaking，则表明当前“停止说话”，可以停止语音识别（STT/ASR）等操作。记住：一定要将返回的bytes 交给识别模块，而不是原始的audio_frame，否则会导致识别结果不正确。
       - C. 如果state为 _vad_instance._vad_state_speaking，则表明当前“说话中”，可以继续进行语音识别（STT/ASR）等操作。记住：一定要将返回的bytes 交给识别模块，而不是原始的audio_frame，否则会导致识别结果不正确。
  备注：如果使用了vad模块，并且希望用vad模块进行语音识别（STT/ASR）等操作，那么一定要将返回的bytes 交给识别模块，而不是原始的audio_frame，否则会导致识别结果不正确。
- 如何更好的排查VAD的问题：包含2个方面，配置和调试。
  - 1. 确保vad模块的初始化参数正确，参考：voice_detection.py。
  - 2. 在state,bytes = on_playback_audio_frame_before_mixing(audio_frame) 中，
    - 1. 将audio_frame的data 的data 保存到本地文件，参考：example_audio_pcm_send.py。这个就是录制原始的音频数据。比如可以命名为：source_{time.time()*1000}.pcm
    - 2. 保存每一次vad 处理的结果：
      - A state==start_speaking的时候：新建一个二进制文件，比如命名为：vad_{time.time()*1000}.pcm，并将bytes 写入到文件中。
      - B state==speaking的时候：将bytes 写入到文件中。
      - C state==stop_speaking的时候：将bytes 写入到文件中。并关闭文件。
  备注：这样就可以根据原始音频文件和vad处理后的音频文件，进行排查问题。生产环境的时候，可以关闭这个功能
# VAD 的调试
- 使用工具类 VadDump，参考：examples/sample_vad，可以帮助排查问题。
- 该方法会生成3类文件：
- sourec.pcm : 音频的原始信息
- vad_{index}.pcm : vad处理后的音频信息,index 是从0开始递增的整数
- label.txt : 音频的标签信息
备注：如果VAD有问题，请将这些文件发送给Agora团队，以便排查问题。


## 如何将TTS生成的音频推入到频道中？
# source code: audio_consumer.py
# sample code: example_audio_consumer.py

## 如何释放资源？
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

## 如何执行打断？