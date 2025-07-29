#!/bin/bash

# set error exit
set -e

# define variables
RTM_URL="https://download.agora.io/sdk/release/rtm_agora_sdk.zip"
TEMP_DIR="/tmp/rtm_install_$$"
AGORA_SDK_DIR="./agora_sdk"

echo "start download RTM SDK..."

# 创建临时目录
mkdir -p "$TEMP_DIR"

# download RTM SDK
echo "downloading from $RTM_URL..."
curl -L -o "$TEMP_DIR/rtm_agora_sdk.zip" "$RTM_URL"

# check download result
if [ ! -f "$TEMP_DIR/rtm_agora_sdk.zip" ]; then
    echo "download failed"
    rm -rf "$TEMP_DIR"
    exit 1
fi

echo "download completed, start unzip..."

# unzip to temporary directory
cd "$TEMP_DIR"
unzip -q rtm_agora_sdk.zip


# check unzip result
if [ ! -d "agora_sdk" ]; then
    echo "unzip failed, agora_sdk directory not found"
    rm -rf "$TEMP_DIR"
    exit 1
fi

echo "unzip completed, start copy files..."

# back to project root directory
cd - > /dev/null

# create target directory: if exists, do not modify, otherwise create
mkdir -p "$AGORA_SDK_DIR"

# check write permission for target directory
if [ ! -w "$AGORA_SDK_DIR" ]; then
    echo "error: no write permission for $AGORA_SDK_DIR"
    echo "please check the directory permission or use sudo to run the script"
    rm -rf "$TEMP_DIR"
    exit 1
fi

# copy rtm_include directory and all .h files, note: overwrite if exists
if [ -d "$TEMP_DIR/agora_sdk/agora_rtm_sdk_c" ]; then
    echo "copy agora_rtm_sdk_c directory..."
    cp -r "$TEMP_DIR/agora_sdk/agora_rtm_sdk_c" "$AGORA_SDK_DIR/"
else
    echo "warning: agora_rtm_sdk_c directory not found"
fi

# copy .so files, not overwrite if exists
echo "copy .so files..."
find "$TEMP_DIR/agora_sdk" -name "*.so" -exec cp -n {} "$AGORA_SDK_DIR/" \;

# copy .dylib files, not overwrite if exists
echo "copy .dylib files..."
find "$TEMP_DIR/agora_sdk" -name "*.dylib" -exec cp -n {} "$AGORA_SDK_DIR/" \;

# 清理临时目录
echo "clean temporary files..."
rm -rf "$TEMP_DIR"

echo "RTM SDK installation completed!"
echo "files copied to $AGORA_SDK_DIR"

