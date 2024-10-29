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