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

如何使用rtm？
1 下载仓库包
2 进入目录，执行make deps
3 执行./scripts/rtmbuild.sh
4 进入目录bin，执行./rtmdemo appid channelname userid
5 rtm sample 参考：cmd/main.go

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
  agoraservice "github.com/AgoraIO-Extensions/Agora-Golang-Server-SDK/v2/go_sdk/rtc"
)
```
- 在代码中调用 agoraservice 接口
```
  svcCfg := agoraservice.NewAgoraServiceConfig()
  svcCfg.AppId = appid
  agoraservice.Initialize(svcCfg)
```
- 运行项目时，记得将 **agora_sdk**目录 (或 Mac 上的 **agora_sdk_mac** 目录) 路径添加到 LD_LIBRARY_PATH (或 Mac 上的 DYLD_LIBRARY_PATH) 环境变量中。

## 如何使用 RTM

1. **构建 RTC**  
   按照上文步骤（`make deps` / `make build`）完成 RTC 的依赖安装和构建。

2. **安装 RTM 相关 SDK**  
   ```bash
   ./script/rtminstall.sh
   ```

3. **构建 RTM Demo**  
   ```bash
   ./script/rtmbuild.sh
   ```

4. **运行 RTM Demo**  
   参考 `cmd/example` 下的 demo 实现，或在 `bin` 目录下运行：
   ```bash
   cd bin
   ./rtmdemo <appid> <channelname> <userid>
   ```
   > 请将 `<appid>`、`<channelname>` 和 `<userid>` 替换为你自己的信息。

5. **更多用法**  
   可参考源码 `cmd/example` 进行集成和开发。

##  ❗ ❗逻辑关系，非常重要 ❗ ❗
- 一个进程只能有一个service instance；在进程开始的时候，创建service；在进程结束的时候，销毁service。
- 一个实例，可以有多个connection，connection可以根据业务需要，随时建立和销毁
- 所有的observer或者是回调中，都不能在调用sdk自身的api，也不能在回调中做cpu耗时的工作，数据拷贝是可以的。
- video codec支持情况：
  - H264： 编码/解码都支持
  - H265： 解码支持，编码不支持
  - AV1： 编码/解码都支持，但编码分辨率必须大于360p，否则会回退到H264
  - VP8： 编码/解码都支持
  - VP9： 编码/解码都支持

# 常见问题
## 编译错误
### 未定义符号引用 GLIBC_xxx
- libagora_rtc_sdk 依赖 GLIBC 2.16 及以上版本
- libagora_uap_aed 依赖 GLIBC 2.27 及以上版本
- 解决方案:
  - 如果可能，你可以升级你的 glibc，或者你需要将运行系统升级到 **所需的操作系统版本**
  - 如果你不使用 VAD，并且你的 glibc 版本在 2.16 和 2.27 之间，你可以通过将 go_sdk/agoraserver/ 中的 **audio_vad.go** 文件重命名为 **audio_vad.go.bak** 来禁用 VAD

#todo
- [ ] 增加对agora_local_user_send_intra_request的支持
- [ ] 增加对on_intra_request_received 的支持
- [ ] 增加onaudio_volume_indication的支持

# 更新日志
ains 成功的标记是：
get ai-ns control extension success
[10/28/25 21:08:54:974][5635][W]:load ai-ns weight resource success


external auido processor 的总体逻辑是：
必须enableamp 是true，否则不处理
 在这样的情况下，可以做算法；也可以只是做vad，不处理算法

externalaudioprocessor 使用方法：
  -1、如何启动算法？  
    a、servcie configure apmModel=1；
    b、externalAudioProcessor.Initialize中，apmconfigure不能为空，并且将所需要的算法设置为true；
    c、如果还需要使用vad，同事设置vadconfigure不能为空！
  
  -2、extern audio processor 怎么只使用vad？
    a、servcie configure apmModel=1
    b、在externalAudioProcessor.Initialize(nil, inputsamplerate, inputchannels, vadConfig,observer)中，设置apmconfigure为空，vadconfigure不能为空！

性能测试对比：
每一个回调结果都打印出来的情况下：
只启动vad：输入1840ms的数据，处理耗时16ms。x 115 倍数
启动apm+vad+关闭apmdump： 输入1840ms的数据，处理耗时 46ms。x 40 倍数

## 2025.12.17 发布 2.4.3 版本

- 新增：支持 `SendIntraRequest` 方法，可主动向远端用户请求编码关键帧（key frame）。

### 测试结果

#### case1: 编码为 Web

如果按照 1s 发送一次请求，编码 Key frame 基本也是 1s。

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

#### case2: 编码为 Android

如果按照 1s 发送一次请求，编码 Key frame 基本为 2s。

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



## 2025.12.15 发布 2.4.2 版本

- 默认变更为 `idleModel=true`
- 在 `externalAudioProcessor` 的回调中，增加 `externalaudioprocessor`
- 增加 `增量发送模式` ，用于增量发送模式。


### 增量发送模式方案特点

1. 仅仅针对'tts'场景，其他场景不建议使用；如果是'iot'下的tts场景，则建议使用
2. 针对connection 级别来做设置，就是说不同的 connection 可以有不同的设置
3. 如果connecton设置为‘增量发送模式',则在设置的时间内，用设置的速率来发送数据，如果超过设置的时间，则恢复为正常发送模式

### 使用方法

1. 在 `connection` 初始化时，设置 `scenario` 参数为 `AudioScenarioAIServer` 或者 `AudioScenarioDefault`
2. 在 `publishConfigure` 内增加参数'SendExternalAudioParameters'：
    - `Enabled` (默认值 false),是否启用增量发送模式
    - `SendSpeed`（1~5，推荐 2，默认值为0）
    - `SendMS`（默认值为0, 快速发送的时间段）
    - `DeliverMuteDataForFakeAdm`（默认值为false，在没有tts 数据的时候，是否发送静音包）

### 产品行为

- 每一轮开始时，先发送一段快发数据
- 之后按照正常的速度（也就是1倍速）发送

> 内部在每一新轮自动调用：`轮次`


## 2025.11.27 发布 2.4.1 版本
-- 增加了对外部来源音频数据的vad、{背景人声消除, 噪声抑制, 回声消除, 自动增益控制, 3a算法}的支持。也就是说不需要加入rtc频道，也可以使用音频处理算法和vad。
-- 修改了LocalUser的返回值，从而和sdk本身的错误码做区分。
-- 增加了example_externalAudioProcessor，用于演示如何使用外部音频数据。
-- utils.go中增加了Lock-freeRingBuffer的实现。
## 2025.11.14 发布 2.4.0 版本
-- 主要更新：rtc& rtm 融合为一个sdk
-- 主要更新：rtm 更新到1.0版本
-- 主要更新：增加ai-ns，agc,bghvs等算法模块支持，可以在内部实现vad、{背景人声消除, 噪声抑制, 回声消除, 自动增益控制, 3a算法}等算法。
## 2025.11.11 发布 2.3.5 版本
-- 更新rtm到1.0 版本
-- 更新rtm sample
## 2025.11.04 发布 2.3.4 版本
-- 更新：update rtc sdk 版本
-- 增加：增加apm 模块，支持下行链路的ns等处理
-- 更新：更新vad 算法模块
-- 增加：增加idleMode 可以做到delay 释放conneciton的c handle
## 2025.08.29 发布 2.3.3 版本
-- 更新：update rtc sdk 到44.32.0829版本,fix 一个启动audiodump 会导致audio-dump 线程泄漏的问题
-- 最佳实现：用connection级别的方式来设置aduiodump！！
## 2025.08.25 发布 2.3.2 版本
-- 融合版本：将rtc 和rtm融合在一起，后续都这样来发布。开发者使用的时候，如果不想使用rtm，就不要去执行scripts/rtminstall.sh和rtmbuild.sh，这样就不会有rtm的依赖。
-- 更新：update rtc sdk 到44.32.0820版本
-- 增加：在pushAudioFrame中，增加pts参数，用于设置帧的pts
-- 增加：在pushVideoFrame中，增加pts参数，用于设置帧的pts
-- 增加：在onplaybackBeforMixing等audioFrame的回调中，增加pts参数
-- 增加：在sample,ai_send_recv_pcm.go中，
  -- 增加v2/v4协议，这个是用验证过可以使用的协议；
  -- 增加PTSAllocator，用来管理PTS的分配
  -- 增加PcmRawDataManager，用来对音频字节做管理，可以pop出来bytesinms长度的音频数据，避免开发者自己去计算
  -- 增加SessionParser，用来对pts做解析，根据具体的协议，可以在session开始和结束的时候提供通知。通知的机制为：
    - sessionid改变的时候：通知上一次session结束；下一个session开始
    - 在超时的时候，通知sesion结束
## 2025.07.21 发布 2.3.1
-- 增加：OnAudioVolumeIndication，可以通过该接口获取当前正在说话的用户的uid
用法：
1、启用audio_volume_indication,并设置参数,如下：localUser.SetAudioVolumeIndicationParameters(intervalInMs, smooth:default=3, isvad)
2、在onAudioVolumeIndication中获取当前正在说话的用户的uid

## 2025.07.16 发布 2.3.0
-- update rtc sdk, fixed 2 bugs
-- 增加对AudioScenarioAiServer 类型scenario的支持
-- 默认的AudioScenario是AudioScenarioAiServer
-- 增加对同一个进程中的conneciton，允许配置为不同的scenario，profile
-- 在创建connection的时候，增加publishconfigure，通过该configure来设置{scenario, profile,publishAudio, publishVideo,ect}
-- conneciton中，增加了：
  - RegisterLocalUserObserver, RegisterAudioFrameObserver, RegisterVideoFrameObserver,RegisterVideoEncodedFrameObserver
  - PublishAudio/UnpublsihAudio, PublishVideo/UnpublsihVideo
  - PushAudioPcmData/PushAudioEncodedData, PushVideoFrame/PushVideoEncodedData
  - InterruptAudio：支持打断功能
  - IsPushToRTCCompleted方法
  - OnAIQoSCapabilityMissing回调接口的实现
  - SendAudioMetaData方法
  - 不在需要人工调用CreateDataStream，内部自动默认
-- 不对外公开的方法：
  - 不再需要开发者人工带用 newMediaNodeFactory
  - 不在需要开发者人工调用：NewCustomAudioTrackPcm, NewCustomAudioTrackEncoded, NewCustomVideoTrack等
  - 不在需要开发者人工调用：NewAudioPcmDataSender等sender
  - 不在需要开发者人工调用 unregisterAudioFrameObserver, unregisterVideoFrameObserver等observer
-- 集成方式： ❗具体集成方式，如何做升级，请咨询SA
1. svcCfg := agoraservice.NewAgoraServiceConfig()
2.agoraservice.Initialize(svcCfg)
3.con := agoraservice.NewRtcConnection(conCfg, publishConfig)
4. register observer;
  -4.1 con.RegisterObserver(conHandler)
  -4.2 con.RegisterAudioFrameObserver
  -4.3 con.RegisterLocalUserObserver(localUserObserver)
5. con.Connect(token, channelName, userId)
6. con.PublisAudio or con.PublishVideo
7. con.PushAudioPcmData/EncodedData or con.PushVideoFrame/EncodedData
8. con.Disconnect
9. con.Release
10.agoraservice.Release()
Note1: 步骤1、2在进程启动的时候调用一次，后面不在需要调用；步骤10在进程退出的时候调用一次，后面不在需要调用。
循环步骤3-9，可以支持多connection。

Note2：
  1、对AI场景，推荐用AudioScenarioAIServer，该模式针对ai的模式，内部做了优化，在降低延迟的同时，能提高弱网体验。(相比chorus，在iphone下，回环延迟低20～30ms，同时弱网下，体验更好)。服务端用AIServer scenario，客户端一定要用AIClient Scenario，否则，语音会有异常。（支持AIClient scenario的sdk版本请咨询SA）
2、对非AI场景，可以配备别的scenario，推荐咨询agora技术支持，以确保在该设置下，能和客户的业务场景匹配
3、在释放conneciton的时候，不再需要人工调用unregister observer，内部会自动unregister
核心变更汇总：
| **Before**                     | **After 2.3.0**                     |
|-------------------------------|-------------------------------------|
| Manual `CreateDataStream`     | ✅ Automatic                        |
| Manual observer unregistration | ✅ Automatic on `Release()`         |
| Fixed per-process scenario    | ✅ Multi-scenario per process      |
| Client/Server scenario mismatch| ❗ `AIClient` mandatory for AI use |  
 ❗ ❗关键提示：支持AICLient scenario的sdk版本请咨询SA

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
			

