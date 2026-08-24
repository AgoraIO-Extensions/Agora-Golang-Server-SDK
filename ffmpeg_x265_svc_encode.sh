#!/bin/bash

# FFmpeg libx265 SVC编码命令 - 3层时域分层，无B帧
# 适用于Mac系统
# 
# 参数说明：
# -c:v libx265: 使用x265编码器
# -preset: 编码速度预设（ultrafast, superfast, veryfast, faster, fast, medium, slow, slower, veryslow）
# -crf: 恒定质量因子（18-28，值越小质量越高）
# -x265-params: x265编码器专用参数

# 基本命令（使用CRF质量控制）
ffmpeg -i input.mp4 \
  -c:v libx265 \
  -preset medium \
  -crf 23 \
  -x265-params "keyint=60:min-keyint=60:scenecut=0:bframes=0:b-adapt=0:ref=3:weightp=0:weightb=0:open-gop=0:rc-lookahead=20:no-open-gop=1:no-scenecut=1:no-weightb=1:no-weightp=1:temporal-layers=3" \
  -c:a copy \
  output.mp4

# 或者使用比特率控制（CBR/VBR）
# ffmpeg -i input.mp4 \
#   -c:v libx265 \
#   -preset medium \
#   -b:v 5M \
#   -maxrate 5M \
#   -bufsize 10M \
#   -x265-params "keyint=60:min-keyint=60:scenecut=0:bframes=0:b-adapt=0:ref=3:weightp=0:weightb=0:open-gop=0:rc-lookahead=20:no-open-gop=1:no-scenecut=1:no-weightb=1:no-weightp=1:temporal-layers=3" \
#   -c:a copy \
#   output.mp4

# 关键参数详解：
# keyint=60: 关键帧间隔（GOP大小）
# min-keyint=60: 最小关键帧间隔
# scenecut=0: 禁用场景切换检测（配合no-scenecut=1）
# bframes=0: 禁用B帧
# b-adapt=0: 禁用B帧自适应选择
# ref=3: 参考帧数量
# weightp=0: 禁用P帧加权预测
# weightb=0: 禁用B帧加权预测（虽然已禁用B帧，但保留此参数）
# open-gop=0: 关闭开放GOP
# rc-lookahead=20: 码率控制前瞻帧数
# no-open-gop=1: 强制关闭开放GOP
# no-scenecut=1: 强制禁用场景切换
# no-weightb=1: 强制禁用B帧加权
# no-weightp=1: 强制禁用P帧加权
# temporal-layers=3: 3层时域分层（SVC核心参数）
