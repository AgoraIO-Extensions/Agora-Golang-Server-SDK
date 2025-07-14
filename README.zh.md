# Required OS and go version
- 支持的 Linux 版本:
  - Ubuntu 18.04 LTS 及以上版本
- macOS 支持(仅仅用于开发测试，不能用于线上环境)
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
- 克隆此仓库并切换到目标分支，然后安装.推荐用最新的分支
```下面用release/2.1.0为例，你可以根据需要选择其他分支
对应每一分支，有相关的tag存在，比如‘release/2.1.0’ 对应的tag为：V2.1.0
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

## 逻辑关系
在一个进程中，只能有一个agoraservice实例.
该实例可以创建多个connection实例.
因此：只能创建/初始化一次agoraservice实例，但可以创建多个connection实例.

# 常见问题
## 编译错误
### 未定义符号引用 GLIBC_xxx
- libagora_rtc_sdk 依赖 GLIBC 2.16 及以上版本
- libagora_uap_aed 依赖 GLIBC 2.27 及以上版本
- 解决方案:
  - 如果可能，你可以升级你的 glibc，或者你需要将运行系统升级到 **所需的操作系统版本**
  - 如果你不使用 VAD，并且你的 glibc 版本在 2.16 和 2.27 之间，你可以通过将 go_sdk/agoraserver/ 中的 **audio_vad.go** 文件重命名为 **audio_vad.go.bak** 来禁用 VAD

# 更新日志
## todo
-- factory 内置在agoraservice中，不再需要单独的factory
--  connection pool实现：在service中实现一个connection pool，用于管理connection的生命周期
  - 从servcie中创建connection
  - connection不在具备release 功能，在service中统一管理
-- datastream：
  - 不需要在conneciton中创建，直接提供sendstreammessge接口
  -- 是否不需要onplayback之类的？但因为我们的客户有pcm，encoded，yuv等，如果单一的用livekit::AudioStream可能不是特别灵活，所以需要解决的是如何防止阻塞操作？而不是设计
  一个新的？？？--建议不需要
  -- 增加一个getstatics的接口
  -- 增加了一个Dequeue，用于线程安全的chan机制。参考： sendrcvpcmyuv done

##Todo：2025.07.03
**SDK**
-- mediafactory 放入到service中，不再需要单独的factory
-- connection:
  - NewRtcConnection(.., publishOption):增加一个publishOption，用于设置publish的参数，包括profile/senario/publish audio和publish video 的类型：pcm/encoded/yuv/encodedimage;内部隐藏对track/sender的实现
  - con.PublishAudio
  - con.PublishVideo
  - con.UnpublishAudio
  - con.UnpublishVideo
  - con.InteruptAudio
  - con.UpdateScenario ???
 

  - con.PushAudioFrame(..):增加一个PushAudioFrame接口，用于推送audio frame
  - con.PushVideoFrame(..):增加一个PushVideoFrame接口，用于推送video frame
  - con.PushEncodedImageFrame(..):增加一个PushEncodedImageFrame接口，用于推送encoded image frame
  - con.PushEncodedVideoFrame(..):增加一个PushEncodedVideoFrame接口，用于推送encoded video frame
**Usage**
process lifecycle: begin, only once
  1. 创建一个service
  2. Initialize service
busisness lifecycle: multiple times
  1. 创建一个connection
  2. register observer: all or only connedction observer??
  2. connect
  3. Publish audio/video
  4. Push audio/video frame
  5. Unpublish audio/video
  6. InterruptAudio
  7. Disconnect: inner auto to call unregister all related observers, no need to call unregister observer manually
  8. Release connection
Process lifecycle: end, only once
  1. Release service

NOTE:
 1. 不在需要Audioconsumer，内部默认设置到16min的内存。
 2. 在audioConsuemer中，用conneciton来代替
 3. 还没有考虑好怎么来处理AudioConsumer？？？其实也就是有可能处理暂停的情况 todo？？从而需要决定是否从pcmsender中去除scenario的字段
 con增加publishconfigure字段
 con：对外隐藏track，sender
 con：提供pushaudio/video接口
 con：不在需要audioconsumer，提供ispushcomplted接口
 4.con.Unreg 全部变成私有的，对开发者来说不需要调用，也不能调用@
 5.con.Reg替代locause.reg，也就是说对开发者loclauser.regxx不再需要
 OnAiqosCapabilityMissing接口
 conn.audioTrack.SetEnabled(true):只需要在conn.create的时候，设置一次就可以；后续不在需要设置
 conn.videoTrack.SetEnabled(true):只需要在conn.create的时候，设置一次就可以；后续不在需要设置??todo: 需要验证下

 localuser：去掉publish，去掉regiser**

 ** 对AudioConsumer的修改**
 1.开发者不需要调用AudioSonsumer，直接内部设置在connetion中
 2.AudioConsumer的设计意义：或者解决的问题
 2.1 避免直接调用pcmsender出现的静音，根因是生产和消费不对齐的问题 ==》 ok
 2.2 提供业务逻辑的暂停和恢复  ==》从端侧native来做，就是提供一个pause的功能，从server端好像不能
 2.3 提供业务逻辑上是否推送完成的判断 ==》将这个逻辑，集成放在connetion中，从而可以达到目标。
  

## todo：2025.05.31 ai_server senario 版本
-- 增加了log/data/config 目录，用于存放日志，数据，配置文件
-- 支持ai_server senario + 支持direct custom audio track
-- 支持对aiqosmissing模式下的通知：所谓aiqosmissing，就是当端侧的版本不支持aiqos能力的时候（是版本，不是端侧所选择的具体senario），需要通过aiqosmissing来通知服务端，服务端在这个会有OnAIQoSCapabilityMissing，在这个回调中，开发者可以返回期望设置的senario，sdk会自动切换到这个senario。如果开发者不希望自动切换，可以返回-1.然后开发者可以自己调用localuse.UpdateAudioTrack(senario)来做切换！

-- server端：需要用directAudioTrack来做推送，不再用customAudioTrack；也不需要按照10ms的速率来发送，就是有多少就发送多少
-- server端：不在需要audioConsumer，或者在audioconsumer内部做对senario的适配，这样能做兼容
-- server端：vad打断的时候，需要调用publish/unpublish方法，替代localaudiotrack的clearsendbuffer方法
是否可以unpub后，多久可以在做pub？
DirectCustomAudioTrackPcm：
- 内部是否还需要设置senddelayinms？内部是否有缓存？--内部没有缓存。
- 内部原理大概是啥？内部会拆分成10个ms一个包，不停的发送。总体的发送速度取决于编码/上行带宽等综合因素。
- 如果app层也是10ms的频率push数据，和cuscomaudiotrack是否行为一致？--基本一致
- 是否可以用在别的senario？==如果用在别的senario，需要按照10ms来推。否则不太推荐。
- ai client目前不限制
- 打断需要用unpub 机制
-- 增加对日志的管理
-- 增加能力协商
-- steromode 给去掉，全部用设置profile，audioscenario + 私有参数的形式来解决

--设计思想
-- 1. 对外来说，并没有customAudioTrack和DirectAudioTrack的区别，都是AudioTrack，只是内部根据senario的不同而做不同的实现；
-- 2. 在service中，根据senario的不同，创建不同的AudioTrack，并做适配；
-- 3. 在service中，根据senario的不同，创建不同的AudioConsumer，并做适配；
-- 4. 在service中，根据senario的不同，对pcmdatasender做不同的适配
-- 5. 在service中，无论是那种senario，对外都是一致的，都是AudioTrack和AudioConsumer、PcmDataSender，但内部根据senario的不同，做不同的实现
--AI server的用法：
-- 1. 在创建service的时候，指定senario为ai_server
-- 2. 别的用法都不用修改
-- 3. 当有ai对话，需要做打断的时候，调用unpub方法，打断后，调用pub方法，重新开始？？？
pub/unpub: 只是api调用，大概是0～1ms；


--测试结果： 0603
  -1. 不能在onplaybackbeformixing中，直接调用sendpcmdata，否则会block！原来的版本是可以的。参考/recv_pcm_loopback
  -2. 能力协商的case：
    case1: server 先加入，别的另外加入？
    case2: server后加入，别的先加入？

done:
--1. publish和unplbish中增加了一个bool判断当前的状态，允许多次pub/unpub。

todo：用法
1、发送编码音频/视频的频率如何来控制？？？好像没有啥控制的，通过pts？？
如果要发送yuv：需要调用9哦那个
	con.SetVideoEncoderConfiguration 来这只编码属性，这个可以在con的周期内，随时调用，来改变编码属性
如果要发送编码视频，需要在publishConfigure中，设置publishVideo为true，同时设置encoderConfiguration，来设置编码属性

## 2025.07.10 发布 2.3.0
-- 增加对AudioScenarioAiServer 类型scenario的支持
-- 默认的AudioScenario是AudioScenarioAiServer
-- 增加对同一个进程中的conneciton，允许配置为不同的scenario，profile
-- 在创建connection的时候，增加publishconfigure，通过该configure来设置{scenario, profile,publishAudio, publishVideo,ect}
-- conneciton中，增加了：
  - RegisterLocalUserObserver, RegisterAudioFrameObserver, RegisterVideoFrameObserver,RegisterVideoEncodedFrameObserver
  - PublishAudio/UnpublsihAudio, PublishVideo/UnpublsihVideo
  - PushAudioPcmData/PushAudioEncodedData, PushVideoFrame/PushVideoEncodedData
  - InterruptAudio
  - IsPushToRTCCompleted方法
  - OnAIQoSCapabilityMissing回调接口的实现
  - SendAudioMetaData方法
  - 不在需要人工调用CreateDataStream，内部自动默认
-- 集成方式：
1. svcCfg := agoraservice.NewAgoraServiceConfig()
2.agoraservice.Initialize(svcCfg)
3.con := agoraservice.NewRtcConnection(conCfg, publishConfig)
4. register observer;
  -4.1 con.RegisterObserver(conHandler)
  -4.2 con.RegisterAudioFrameObserver
  -4.3 con.RegisterLocalUserObserver(localUserObserver)
5.con.Connect(token, channelName, userId)
6. con.PublisAudio/Video
7. con.PushAudioPcmData/EncodedData
8. con.Disconnect
9. con.Release
10.agoraservice.Release()

Note：
1、对单纯的音频场景，api 调用数目从上一个版本的21个，降低到了当前版本的12个api调用，减少了9个。关键的在做释放操作的时候，不在依赖开发者调用api的时序，内部自动处理，降低了出错概率
medianode/track/sender，unreg（3个），release 3个，总共减少了9个调用
2、对AI场景，推荐用AudioScenarioAIServer，这个场景，针对ai的模式，内部做了优化，在降低延迟的同时，能提高弱网体验。(相比chorus，在iphone下，回环延迟低20～30ms，同时弱网下，体验更好)
对非AI场景，可以配备别的scenario，推荐咨询agora技术支持，以确保在该设置下，能和客户的业务场景匹配

## 2025.07.08 发布 2.2.11
-- 增加对AudioScenarioAiServer 类型scenario的支持
-- 默认的AudioScenario是AudioScenarioAiServer
-- 增加对同一个进程中的conneciton，允许配置为不同的scenario，profile

## 2025.07.02 发布 2.2.10
-- remove： 去除了在enableSteroEncodeMode函数中设置senraio的代码，因为这个senario的设置，可以在NewRtcConnection中作为一个参数传入，不在需要在这个enableSteroEncodeMode中设置。

## 2025.07.01 发布 2.2.9（re-update)
-- 增加了log/data/config 目录，用于存放日志，数据，配置文件
-- 支持ai_server senario + 支持direct custom audio track
-- 支持对aiqosmissing模式下的通知：所谓aiqosmissing，就是当端侧的版本不支持aiqos能力的时候（是版本，不是端侧所选择的具体senario），需要通过aiqosmissing来通知服务端，服务端在这个会有OnAIQoSCapabilityMissing，在这个回调中，开发者可以返回期望设置的senario，sdk会自动切换到这个senario。如果开发者不希望自动切换，可以返回-1.然后开发者可以自己调用localuse.UpdateAudioTrack(senario)来做切换！
-- 支持connection级别的设置senario，也就是说在一个进程内，connection可以设置为不同的senario，比如一个connection设置为ai_server，另一个connection设置为chorus
-- 支持customerAudioTrack根据senario来做创建
-- 打断：localUser.InterruptAudio(audioconsumer*)，这个api是用于打断ai对话的。建议在做ai打断的时候，用这个新的方式来实现！
-- UpdateAudioSenario：LocalUser.UpdateAudioSenario,用来更新connection的senario以及和conncection关联的audioTrack的senario。通常不需要主动调用！
NOTE：customAudioTrack创建参数里面的senario需要和connection创建参数的senario保持一致！！否则audioTrack的行为可能和预期不符合！
-- publish和unplbish中增加了一个bool判断当前的状态，允许多次pub/unpub
-- fix a bug in send_mp4 sample



## 2025.05.26 发布 2.2.8
-- 修改：修改一个bug，在立体声模式下，编码码率不生效的bug
-- 更新：将mac sdk 版本更新到4.4.32
## 2025.05.21 发布 2.2.7
-- 更新：vad v1 lib removed libc++ and libstdc++ dependence，fix a bug that the vad lib can not be loaded in some os which has not libc++ and libstdc++.
## 2025.05.19 发布 2.2.6
-- 更新vadv1 库，fix 一个参数传递在某些case下会传递不成功的bug
-- 去除sample_vad.go中的打印信息
## 2025.04.28 发布 2.2.5
-- 替换：更新rtc sdk 到4.4.32
## 2025.04.28 发布 2.2.4
-- 增加了EnableEncryption接口，用来设置是否启用加密;
-- 增加了onEncryptionError接口，用于处理加密错误;
## 2025.04.15 发布 2.2.3
-- 增加了一个Dequeue，用于线程安全的chan机制。参考： send_recv_yuv_pcm.go
-- 增加localauiodstats, removeaudiostats, localvideostats, remotevideostats。参考：send_recv_yuv_pcm.go
## 2025.03.26 发布 2.2.2
-- 修复：vad Release中的一个不严谨的地方，多次调用Release会导致crash
-- 增加：sample_vad.go中增加一个函数，用于测试stero 音频的vad
## 2025.03.05 发布 2.2.1
-- configure中水装置是否支持双声道编码，默认是不支持。如果要支持，内部会默认修改私有参数，从而可以支持双声道编码。但app层，在设置playback参数的时候，需要设置channel=2. ok
-- 提供Vad v1 算法 ok
-- vad v1算法接口改变：对外提供rms等控制参数？？ok
-- sdk包中需要带上vad v1 算法的库文件。ok
-- 提供SteroVad 类  ok
-- audioconsumer：
  - 默认设置为100ms的缓冲区大小 ok
  - 增加：isPushToRtcCompleted接口 ok
  - AudioConsumer::_samples_per_channel 在立体声下计算错误的，fixed ok
-- videoObserver
  - OnFrame中拦截userid为0的frame。目前是一个bug。就是如果本地也发送视频的时候，onFrame会回调uid=0的frame，这个frame是无效的，需要拦截掉。ok
-- 文档：
  - 增加对私有部署、proxy部署的调用说明 ok
-- SDK 更新：sdk 中增加vad v1库，并重新命名，但保持原始版本号，也兼容原始版本。

-- 增加：
  -- onVolumeIndication的接口？？
  -- 接收编码音频的接口？
  -- 支持rtm？
  -- 增加一个对kit的封装，方便使用。管理
    -- {service/factory,connection的创建等}等， 
    -- 然后client，管理{connection, localuser, track,pcm_sender,audio_consumer, vad_manager}
## 2024.12.23 发布 2.2.0
-- 更新：
  -- 更新sdk 版本到4.4.31
-- 增加：
  -- 在LocalUser中增加SendAudioMetaData接口
    -- 用法是：直接调用，频率限制在100内；每一次的数据长度在64 bytes以内
    -- 测试数据？？？
  -- 在LocauUserObserver中增加onAudioMetaDataResult接口
  -- serviceconfigure中增加domainLimit成员，用在是否限制在Domain的为url的情况。默认是0，表示不限制
-- 修改：
  -- ExternalVideoFrame增加对额外yuv colorspace的支持，从而可以编码纯色的背景图。通常用在数字人场景
  考虑如果方便的做audio meta的输入？？释放在pcm数据里面加入？？？但如果这样加入，又如何做audioconsumer？？除非内部设置audioconsumer是10ms的timer？？
## 2024.12.18 发布 2.1.4
-- 修改：
  -- 默认支持VAD v2模块
  -- VAD的参数默认从20，修改到50
  -- 添加了VAD识别的参数配置说明。参考下面。
## 2024.12.11 发布 2.1.3
-- 增加：
  - 增加了AuduioVadManager，用来管理多个音频源的VAD。在实际用的时候，一个音频源就需要一个Vad实例子，为了简化开发者的开发难度，提供一个统一的接口来管理这些Vad实例。
-- 修改：
  -- 在AudioFrameObserver::OnPlaybackAudioFrameBeforeMixing中增加了2个参数：vadResultState和vadResultFrame，用来返回VAD识别结果。
  -- 在LocalUser::RegisterAudioFrameObserver中，增加2个参数，用来设置是否需要启动VAD识别。如果启动VAD识别，VAD识别后的结果会在OnPlaybackAudioFrameBeforeMixing中返回。
2.1.3版本后VAD的用法：
  1、在LocalUser::RegisterAudioFrameObserver中,启动vad。如：
      vadConfigure := &agoraservice.AudioVadConfigV2{
        PreStartRecognizeCount: 16,
        StartRecognizeCount:    30,
        StopRecognizeCount:     50,
        ActivePercent:          0.7,
        InactivePercent:        0.5,
        StartVoiceProb:         70,
        StartRms:               -50.0,
        StopVoiceProb:          70,
        StopRms:                -50.0,
      }
      localUser.RegisterAudioFrameObserver(audioObserver, 1, vadConfigure)
  2、在AudioFrameObserver::OnPlaybackAudioFrameBeforeMixing中，可以获取到VAD识别结果，并对结果做处理
  you_audio_frame_observer::OnPlaybackAudioFrameBeforeMixing(channelId string, uid uint64, frame *agoraservice.AudioFrame, vadResultState agoraservice.VadState, vadResultFraem *agoraservice.AudioFrame){
    if vadResultState == agoraservice.VadStateStartSpeaking { // 开始说话，可以将vadResultFrame 的数据交给ASR/STT}
    if vadResultState == agoraservice.VadStateSpeaking{ // 正在说话，可以将vadResultFrame 的数据交给ASR/STT}
    if vadResultState == agoraservice.VadStateStopSpeaking{ // 停止说话，可以将vadResultFrame 的数据交给ASR/STT；//并且做结束说话的业务处理}
  }
  备注：如果启用了VAD，就一定要用到vadResultFrame的数据，而不用用frame的数据，否则stt/ars 做识别的时候，结果会丢失。
## 2024.12.02 发布 2.1.2
- 增加了AudioConsumer,可以用来对音频数据做push
- 增加了VadDump方法，可以用来对vad做debug
- 当service的profile是chorus的时候，
  - 在localuser中，自动增加setSenddelayinms
  - 在audiotrack中，自动设置track为chorus模式
  - 开发者在Ai场景中，不在需要人工设置setsenddelayinms 和设置track为chorus模式
- 修改service中的全局mutex 为sync.map模式，对mutex的粒度做拆分
- 增加了audioloopback sample，用来做音频回环测试
- 修改sample，提供命令行输入appid，channelname的模式
- 增加对AudioLable 插件的支持，开发者不在需要在app层调用EableExtension
- 增加了onpublishstatechanged 接口
- 修改了VAD的返回状态，统一为：NoSpeakong, Speaking, StopSpeaking 3种状态；而且在StopSpeaking的时候，也会返回当前的frame 数据；
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

## 在AI场景下，如何做打断？
- 打断的定义
  在人机对话中，打断是指用户在对话过程中突然打断机器人的回答，要求机器人立即停止当前回答并转而回答用户的新问题。这个行为就叫打断
- 打断触发的条件
  打断根据不同的产品定义，一般有两种方式：
  - 模式1：语音激励模式. 当检测到有用户说话，就执行打断策略，比如用户说话时，识别到用户说话，就执行打断策略。
  - 模式2：ASR激励模式. 当检测到有用户说话、并且asr/STT的识别返回有结果的时候，就执行打断策略。
- 不同打断策略的优点
  - 1. 语音激励打断：
    - 优点：
    - 1. 减少用户等待时间，减少用户打断的概率。因为用户说话时，机器人会立即停止回答，用户不需要等待机器人回答完成。
    - 缺点：
    - 1. 因为是语音激励模式，有可能会被无意义的语音信号给打断，依赖于VAD判断的准确性。比如AI在回答的时候，如果有人敲击键盘，就可能触发语音激励，将AI打断。
  -2 . ASR激励打断：
    - 优点：
    - 1. 降低用户打断的概率。因为用户说话时，asr/STT识别到用户说话，才会触发打断策略。
    - 缺点：
    - 1. 因为是asr/STT激励模式，需要将语音信号转换成文本，会增加打断的延迟。

- 推荐模式
  如果VAD能过滤掉非人声，只是在有人声的时候，才触发VAD判断，建议用语音激励模式；或者是对打断要求延迟敏感的时候，用改模式
  如果对打断延迟不敏感，建议用ASR激励模式，因为ASR激励模式，可以过滤掉非人声，降低用户打断的概率。
- 如何实现打断？打断需要做哪些操作？
  定义：人机对话，通常可以理解为对话轮的方式来进行。比如用户问一个问题，机器人回答一个问题；然后用户再问一个问题，机器人再回答一个问题。这样的模式就是对话轮。我们假设给对话轮一个roundId,每轮对话，roundId+1。 一个对话轮包含了这样的3个阶段/组成部分：vad、asr、LLM、TTS、rtc推流。
  1. vad： 是指人机对话的开始，通过vad识别出用户说话的开始和结束，然后根据用户说话的开始和结束，交给后续的ASR。
  2. asr： 是指人机对话的识别阶段，通过asr识别出用户说的话，然后交给LLM。
  3. LLM： 是指人机对话的生成阶段，LLM根据用户说的话，生成一个回答。
  4. TTS： 是指人机对话的合成阶段，LLM根据生成的回答，合成一个音频。
  5. rtc推流： 是指人机对话的推流阶段，将合成后的音频推流到rtc，然后机器人播放音频。

  因此，所谓的打断，就是在（roundid+1）轮的时候，无论是用语音激励（VAD阶段触发）还是用ASR激励（就是在ASR识别出用户说的话）打断，都需要做如下的操作：
  1. 停止当前轮roundID轮的LLM生成。
  2. 停止当前轮roundID轮的TTS合成。
  3. 停止当前轮roundID轮的RTC推流。
   API调用参考：
    a 调用:AudioConsumer.clear()；
    b 调用:LocalAudioTrack.ClearSenderBuffer()；
    c 业务层：清除TTS返回来保留的数据（如果有）
## LLM的结果什么时候交给TTS做合成？
  LLM的结果是异步返回的，而且都是流式返回的。应该按照什么时机将LLM的结果交给TTS做合成呢？
  需要考虑2个因素：
  1. 无歧义、连续、流畅：确保TTS合成的语音是没有歧义、而且是完整、连续的。比如LLM返回的文本是："中间的首都是北京吗？"如果我们给TTS的是：中  然后是：国首  然后是：是北   然后是：京吗？  这样合成会有歧义，因为"中"和"国"之间没有空格，"首"和"是"之间没有空格，"京"和"吗"之间没有空格。
  2. 确保整个流程延迟最低。LLM 生成完成后，在交给TTS，这样的处理方式，合成的语音一定是没有歧义，而且是连续的。但延迟会很大，对用户体验不友好。
  推荐的方案：
    将有LLM返回数据的时候：
    a LLM返回的结果存放在缓存中
    b 对缓存中的数据做逆序扫描，找到最近的一个标点符号
    c 将缓存中的数据，从头开始到最尾的一个标点符号截断，然后交给TTS做合成。
    d 将截断后的数据，从缓存中删除。剩余的数据，移动到缓存头位置，继续等待LLM返回数据。

## VAD配置参数的含义
AgoraAudioVadConfigV2 属性
属性名	                 类型	    描述	                         默认值	     取值范围
preStartRecognizeCount	int	    开始说话状态前保存的音频帧数	      16	       [0, ]  
startRecognizeCount	    int   	判断是否开始说话状态的音频帧总数     30	    [1, max]
stopRecognizeCount	    int	    判断停止说话状态的音频帧数	        50	    [1, max]
activePercent	          float	  在 startRecognizeCount 
                                  帧中活跃帧的百分比	            0.7	    0.0, 1.0]
inactivePercent       	float	  在 stopRecognizeCount
                                 帧中非活跃帧的百分比	             0.5     [0.0, 1.0]
startVoiceProb	        int	    音频帧是人声的概率	              70	    [0, 100]
stopVoiceProb	          int	     音频帧是人声的概率	               70	    [0, 100]
startRmsThreshold	      int	     音频帧的能量分贝阈值	            -50	     [-100, 0]
stopRmsThreshold	      int	    音频帧的能量分贝阈值            	-50	    [-100, 0]
注意事项

startRmsThreshold 和 stopRmsThreshold:
值越高，就需要说话人的声音相比周围环境的环境音的音量越大
在安静环境中推荐使用默认值 -50。
在嘈杂环境中可以调高到 -40 到 -30 之间，以减少误检。
根据实际使用场景和音频特征进行微调可获得最佳效果。

stopReecognizeCount: 反映在识别到非人声的情况下，需要等待多长时间才认为用户已经停止说话。可以用来控制说话人相邻语句的间隔，在该间隔内，VAD会将相邻的语句当作一段话。如果该值短，相邻语句就越容易被识别为2段话。通常推荐50～80。
比如：下午好，[interval_between_sentences]北京有哪些好玩的地方？
如果说话人语气之间的间隔interval_between_sentences 大于stopReecognizeCount，那么VAD就会将上述识别为2个vad：
vad1: 下午好
vad2: 北京有哪些好玩的地方？
如果interval_between_sentences 小于 stopReecognizeCount，那么VAD就会将上述识别为1个vad：
vad： 下午好，北京有哪些好玩的地方？

如果对延迟敏感，可以调低该值，或者咨询研发，在降低该值的情况下，应该如何在应用层做处理，在保障延迟的情况下，还能确保语意的连续性，不会产生AI被敏感的打断的感觉。

## 双声道编码的用法：除非必要，不推荐!如果需要使用，请联系研发确认！
 # 适用场景
 需要客户端一定是双声道的场景，就是左/右声道的数据一定不一样。记住：不推荐，一定是左右声道数据不一样的情下！
 # 客户端用法
 参考文档
 # 服务端用法
1、在serviceConfigure设置：设置audoSenariao为：gameStreaming；并且enable双声道编码
svcCfg := agoraservice.NewAgoraServiceConfig()
	svcCfg.AppId = appid
	// change audio senario
	svcCfg.AudioScenario = agoraservice.AudioScenarioGameStreaming
	svcCfg.EnableSteroEncodeMode =  1

	agoraservice.Initialize(svcCfg)

2、回调参数设置为双声道：！！因为双声道的vad只支持16k的采样率，所以需要设置回调参数为双声道，采样率为16k
localUser.SetPlaybackAudioFrameBeforeMixingParameters(2, 16000)

3、audioframeobserver中，取消vad
localUser.RegisterAudioFrameObserver(audioObserver, 0, nil)

4、在audioframeobserver的回调中，用SteroVad来做双声道的vad检查
4.1 StereoVad初始化（回调前做）
//vad v1 for stero 
	vadConfigV1 := &agoraservice.AudioVadConfig{
		StartRecognizeCount:    10,
		StopRecognizeCount:     6,
		PreStartRecognizeCount: 10,
		ActivePercent:          0.6,
		InactivePercent:        0.2,
    // 别的参数可以根据需要调节，也可以采用默认的
	}
	// generate stero vad
	steroVadInst := agoraservice.NewSteroVad(vadConfigV1, vadConfigV1)

4.2 双声道vad的检查：在audioframeobserver的回调中，用SteroVad来做双声道的vad检查
audioObserver := &agoraservice.AudioFrameObserver{
		OnPlaybackAudioFrameBeforeMixing: func(localUser *agoraservice.LocalUser, channelId string, userId string, frame *agoraservice.AudioFrame, vadResultState agoraservice.VadState, vadResultFraem *agoraservice.AudioFrame) bool {
			// do something...
			
			// do stero vad process
			
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
			fmt.Printf("left vad state %d, left len %d, right vad state %d, right len: %d,diff = %d\n", leftState, leftLen, rightState, rightLen, end-start)
			
			// dump vad frame: only for debug, when release, please remove this
			dumpSteroVadResult(1, leftFrame, leftState)
			dumpSteroVadResult(0, rightFrame, rightState)

      // do something...
      // return 
      return true
    }
  5、释放
  在agoraservice::release()之后调用，才能生效。如下代码：
  	agoraservice.Release()
	if steroVadInst != nil {
		steroVadInst.Release()
	}
	steroVadInst = nil

## 代理和私有化的区别
  -代理应对的是LAN环境的网络限制，通过proxy，能够让局域网内的设备通过公网访问agora的rtc服务。
  -私有化部署，则是将agora的rtc服务部署到自己的服务器上，通过私有化部署，可以更好的控制agora的rtc服务，同时也可以更好的保护用户的隐私。

## 如何支持私有部署
1、调用时序：
***就是需要在创建conneton之后、在con.Connect 之前就调用才能生效。
如果在agoraservice::init()之后调用，则不会生效。

2、Api调用：如下
connetioncon := agoraservice.NewRtcConnection(&conCfg)
localUser := con.GetLocalUser()

	
		// test lan 
	params1 := `{"rtc.enable_nasa2": false}`
	params2 := `{"rtc.force_local": true}`
	params3 := `{"rtc.local_domain": "ap.1452738.agora.local"}`
	params4 := `{"rtc.local_ap_list": ["20.1.125.55"]}`

	// 设置参数
	agoraservice.GetAgoraParameter().SetParameters(params1)
	agoraservice.GetAgoraParameter().SetParameters(params2)
	agoraservice.GetAgoraParameter().SetParameters(params3)
	agoraservice.GetAgoraParameter().SetParameters(params4)
	//
## 如何支持代理
1、调用时序：
在agoraService::init()之后调用(但 没有验证是否真的生效？)

2、Api调用：如下
parameter.setBool("rtc.enable_proxy", true);
s->setBool("rtc.force_local",true);
s->setBool("rtc.local_ap_low_level",true);

s->setBool("rtc.enable_nasa2", true);
s->setParameters("{\"rtc.vos_list\":[\"10.62.0.95:4701\"]}");
s->setParameters("{\"rtc.local_ap_list\":[\"10.62.0.95\"]}");

## 如何使用ai server/ai client？
-- 如何做打断：调用localUser.UnpublishAudio(track)
-- 在推送的时候，调用localUser.PublishAudio(track)，可以多次调用也可以的，没有啥影响。
-- pub/unpub：api调用是在0～1ms，callback在：1~2ms内！
-- upub：api是在0～1ms，callback在：1~2ms内！

## Ai server/ai client如何使用
-- 1、ai server：在init之后，调用aiService.InitAiServer()，然后调用aiService.StartAiServer()，即可启动ai server。

## 有关AudioFrameObserver
-- playbackbeforemixing 和playbackAudioFrame是可以同时触发的，但不建议这样使用，因为这样会导致回调的频率很高，从而影响性能。
-- onMixed 不会触发：因为server sdk并没有采集，所以不会触发。在sever sdk可以用onPlaybackAudioFrame来获取。
			

