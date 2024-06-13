#!/bin/bash

PACKAGE_HOME=$(
    cd $(dirname $0)/..
    pwd
)

UNAME_S=`uname -s`
OS=unknown

linux_sdk="https://download.agora.io/sdk/release/agora_rtc_sdk-linux-gnu-v4.4.30-20241012_114642-379681.zip"
mac_sdk="https://download.agora.io/sdk/release/agora_rtc_sdk_mac_v4.4.30_22304_FULL_20241015_1616_384496.zip"

if [[ $UNAME_S == Linux ]]; then
    OS=linux
elif [[ $UNAME_S == Darwin ]]; then
    OS=mac
# else ifeq ($(UNAME_S),CYGWIN) # Cygwin is Unix-like environment on Windows
#     OS := windows
# else ifeq ($(UNAME_S),MINGW32) # MinGW is a native Windows port of GNU tools
#     OS := windows
# else ifeq ($(UNAME_S),MSYS) # MSYS is a collection of GNU utilities for Windows
#     OS := windows
else
    echo "Unsupported OS: ${UNAME_S}"
    exit 1
fi

echo "OS: ${OS}"

# Download the Agora RTC SDK
if [[ -d agora_sdk ]]; then
    rm -rf agora_sdk
fi
if [[ -d agora_sdk_mac ]]; then
    rm -rf agora_sdk_mac
fi
if [[ -f agora_sdk.zip ]]; then
    rm -f agora_sdk.zip
fi
if [[ -f agora_sdk_mac.zip ]]; then
    rm -f agora_sdk_mac.zip
fi

if [[ $OS == mac ]]; then
    curl -o agora_sdk_mac.zip $mac_sdk
    unzip agora_sdk_mac.zip
    mv agora_sdk agora_sdk_mac
    rm -f agora_sdk_mac.zip
fi
curl -o agora_sdk.zip $linux_sdk
unzip agora_sdk.zip
rm -f agora_sdk.zip