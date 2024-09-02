# 在linux上编译运行
## 运行环境要求
- 支持的Linux版本：
  - Ubuntu 18.04 LTS及以上版本
  - CentOS 7.0及以上版本
- go版本：
  - go 1.20及以上版本
  - 未在go 1.19及以下版本上进行测试

## 准备Agora RTC SDK的C版本
- 下载并解压缩 [agora_sdk.zip](https://download.agora.io/sdk/release/agora_rtc_sdk_linux_20240902_320567.zip)
```
unzip agora_sdk.zip
```
- 在与 **go_wrapper** 目录相同的目录中创建 **agora_sdk** 目录
- **agora_sdk** 目录中应包含 **libagora_rtc_sdk.so** 和 **include_c** 文件夹

## 编译示例代码
```
cd go_wrapper
go mod tidy
go build main
```

## 测试
- 下载并解压缩 [test_data.zip](https://download.agora.io/demo/test/test_data_202409021506.zip)
- 在与 **go_wrapper** 目录相同的目录中创建 **test_data** 目录
- 运行 main
```
export LD_LIBRARY_PATH=/path/to/agora_sdk
cd go_wrapper
./main
```
- 在Linux上运行单元测试
```
export LD_LIBRARY_PATH=/path/to/agora_sdk
cd go_wrapper/agoraservice
go mod tidy
go test -v -count=1 -timeout 20s -run ^TestBaseCase$ agoraservice
```

# 在Mac上编译和运行
- 下载并解压缩 [agora_sdk_mac.zip](https://download.agora.io/sdk/release/agora_rtc_sdk_mac_20240902_320567.zip)
```
unzip agora_sdk_mac.zip
```
- 在与 **go_wrapper** 目录相同的目录中创建 **agora_sdk_mac** 目录
- 下载并解压缩 [test_data.zip](https://download.agora.io/demo/test/test_data_202409021506.zip)
- 在与 **go_wrapper** 目录相同的目录中创建 **test_data** 目录
- 编译并运行 main
```
cd go_wrapper
./build_for_mac.sh
./main
```
- 编译并进行单元测试
```
export CGO_LDFLAGS_ALLOW="-Wl,-rpath,.*"
cd go_wrapper/agoraservice
go mod tidy
go test -v -count=1 -timeout 20s -run ^TestBaseCase$ agoraservice
```

# vad用法
- 调用 NewAudioVad 创建一个新的vad实例
- 调用 ProcessPcmFrame 处理一个音频帧，该帧是一个16位、16kHz和单声道的pcm数据
- ProcessPcmFrame(frame *PcmInAudioFrame) (*PcmOutAudioFrame, int) 是一个具有双返回值的函数，表示语音状态的更新。int返回值表示当前语音活动检测（VAD）状态，其中：

- 0 表示未检测到语音；
- 1 表示语音开始；
- 2 表示正在进行语音；
- 3 表示当前语音段结束；
- -1 表示发生错误。
- 如果函数处于状态1、2或3，则 PcmOutAudioFrame 将包含与VAD状态相对应的PCM数据。

- 当用户想要执行ASR/TTS处理时，应将 PcmOutAudioFrame 的数据发送到ASR系统。
- 当不再需要vad实例时，调用 Release()

-- 注意：
- 在不再需要ASR系统时，应释放VAD实例。
- 一个VAD实例对应一个音频流

# 常见问题解答
## 编译错误
### 未定义符号引用到GLIBC_xxx
- libagora_rtc_sdk 依赖于 GLIBC 2.16及以上版本
- libagora_uap_aed 依赖于 GLIBC 2.27及以上版本
- 解决方案：
  - 如果可能，您可以升级您的glibc，或者您需要将运行系统升级到**所需的操作系统版本**
  - 如果您不使用VAD，并且您的glibc版本在2.16和2.27之间，您可以通过将 go_wrapper/agoraserver/ 目录中的 **audio_vad.go** 文件重命名为 **audio_vad.go.bak** 来禁用VAD

# 更新日志
## 2024.08.14 发布 1.1 版本
- 为 RtcConnection 添加 RenewToken 接口
- 更新 VAD 算法
## 2024.06.28
- 支持在Mac上构建和测试
## 2024.06.25
- 使 VAD 可用，详细信息请参见 agoraservice/rtc_connection_test.go 中的 UT case TestVadCase
- 通过将 AudioScenario 设置为 AUDIO_SCENARIO_CHORUS 来减少音频延迟，详细信息请参见 agoraservice/rtc_connection_test.go 中的 UT case TestBaseCase
- 添加用于调整发布音频音量的接口，示例请参见 sample.go
- 通过设置音频发送缓冲区并确保音频数据在实际发送时间之前发送，解决接收端接收到的音频中的噪音问题。详细信息请参见 sample.go
