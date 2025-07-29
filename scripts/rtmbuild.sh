#!/bin/bash

# 设置错误时退出
set -e

# 创建输出目录
if [ -d "bin" ]; then
    echo "bin dir already exists, deleting all files inside..."
    rm -rf bin/*
else
    mkdir -p bin
fi

# 构建命令
echo "Building example application..."
go build -o ./bin/ ./cmd/example/



# 检查构建是否成功
if [ $? -eq 0 ]; then
    echo "Build successful! Binary is located at ./bin/example"
else
    echo "Build failed!"
    exit 1
fi 