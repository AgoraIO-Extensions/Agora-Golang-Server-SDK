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

# 常见问题
## 编译错误
### 未定义符号引用 GLIBC_xxx
- libagora_rtc_sdk 依赖 GLIBC 2.16 及以上版本
- libagora_uap_aed 依赖 GLIBC 2.27 及以上版本
- 解决方案:
  - 如果可能，你可以升级你的 glibc，或者你需要将运行系统升级到 **所需的操作系统版本**
  - 如果你不使用 VAD，并且你的 glibc 版本在 2.16 和 2.27 之间，你可以通过将 go_sdk/agoraserver/ 中的 **audio_vad.go** 文件重命名为 **audio_vad.go.bak** 来禁用 VAD

# 更新日志
## 2025.02.23 发布 2.2.1(on development)
-- 增加：
  -- onVolumeIndication的接口？？
  -- 接收编码音频的接口？
  -- 支持rtm？
  -- 增加对push audio结束的api？python已经ok
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

