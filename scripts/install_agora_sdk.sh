#!/bin/bash

PACKAGE_HOME=$(
    cd $(dirname $0)/..
    pwd
)

UNAME_S=`uname -s`
OS=unknown


#ver2.2.4
#linux_sdk="https://download.agora.io/sdk/release/agora_rtc_sdk-x86_64-linux-gnu-v4.4.32-20250425_144419-675648-aed-0521.zip"
#mac_sdk="https://download.agora.io/sdk/release/agora_sdk_mac_v4.4.32_24464_FULL_20250522_1619_711331-aed.zip"

#version 2.2.9 for aiqos 
#linux_sdk="https://download.agora.io/sdk/release/agora_rtc_sdk-x86_64-linux-gnu-v4.4.32-20250528_175520-720646-aed.zip"
#mac_sdk="https://download.agora.io/sdk/release/agora_sdk_mac_v4.4.32_24517_FULL_20250528_1905_720662-aed.zip"

linux_sdk="https://download.agora.io/sdk/release/agora_rtc_sdk-x86_64-linux-gnu-v4.4.32-20250612_142144-741749-aed.zip"
mac_sdk="https://download.agora.io/sdk/release/agora_sdk_mac_v4.4.32_24654_FULL_20250612_1503_741767-aed.zip"

#date: 20250617 for fix oncapability issue

linux_sdk="https://download.agora.io/sdk/release/agora_rtc_sdk-x86_64-linux-gnu-v4.4.32-20250617_170645-749554-aed.zip"
mac_sdk="https://download.agora.io/sdk/release/agora_sdk_mac_v4.4.32_24690_FULL_20250617_1817_749572-aed.zip"

linux_sdk="https://download.agora.io/sdk/release/agora_rtc_sdk-x86_64-linux-gnu-v4.4.32-20250715_161625-791246-aed.zip"
mac_sdk="https://download.agora.io/sdk/release/agora_sdk_mac_v4.4.32_24915_FULL_20250715_1710_791284-aed.zip"

#date:20250814 for present timestamp issue
linux_sdk="https://download.agora.io/sdk/release/agora_rtc_sdk-x86_64-linux-gnu-v4.4.32-20250814_105012-837202-aed.zip"
mac_sdk="https://download.agora.io/sdk/release/agora_sdk_mac_v4.4.32_25242_FULL_20250814_1133_837220-aed.zip"

#date:20250818 for oncallback with pts
linux_sdk="https://download.agora.io/sdk/release/agora_rtc_sdk-x86_64-linux-gnu-v4.4.32-20250818_170541-843836-aed.zip"
mac_sdk="https://download.agora.io/sdk/release/agora_sdk_mac_v4.4.32_25295_FULL_20250818_1850_843923-aed.zip"
#mac_sdk="https://download.agora.io/sdk/release/agora_sdk_mac_v4.4.32_25295_FULL_20250818_1850_843923-test.zip"
#date: 20250820 update sdk with pts callback, and add sessionParser helper for parse pts
linux_sdk="https://download.agora.io/sdk/release/agora_rtc_sdk-x86_64-linux-gnu-v4.4.32-20250820_100926-846510-aed.zip"
mac_sdk="https://download.agora.io/sdk/release/agora_sdk_mac_v4.4.32_25321_FULL_20250820_1052_846534-aed.zip"

#date: 20250829 for update sdk with framedump
linux_sdk="https://download.agora.io/sdk/release/agora_rtc_sdk-x86_64-linux-gnu-v4.4.32-20250829_160340-860733-aed.zip"
mac_sdk="https://download.agora.io/sdk/release/agora_sdk_mac_v4.4.32_25418_FULL_20250829_1647_860754-aed.zip"

#date: 20251028 update sdk for apm filter only test purpose for mac
mac_sdk="https://download.agora.io/sdk/release/agora_sdk_mac_v4.4.32_25776_FULL_20251021_1824_939915-aed.zip"
mac_sdk="https://download.agora.io/sdk/release/agora_sdk_mac_v4.4.30_25869_FULL_20251030_1836_953684-aed.zip"
linux_sdk="https://download.agora.io/sdk/release/agora_rtc_sdk_x86_64-linux-gnu-v4.4.32.150_26715_SERVER_20251030_1807-aed.zip"

#date:20251121 update sdk for auto load so and auto load ains resource, change audio-sink with userdata
linux_sdk="https://download.agora.io/sdk/release/agora_rtc_sdk_x86_64-linux-gnu-Agora_Native_SDK_for_Linux_x64_zhourui_26895_SERVER_20251121_1628_987405_20251021_1427-3a.zip"
mac_sdk="https://download.agora.io/sdk/release/agora_sdk_mac_Agora_Native_SDK_for_Mac_zhourui_26100_FULL_20251121_1643_987406_20251021_1427-3a.zip"

#date:20251124 update sdk for fix the issue of auto load so and auto load ains resource
linux_sdk="https://download.agora.io/sdk/release/agora_rtc_sdk_x86_64-linux-gnu-Agora_Native_SDK_for_Linux_x64_zhourui_26895_SERVER_20251121_1628_987405_20251021_1427-3a.zip"
#mac_sdk="https://download.agora.io/sdk/release/agora_sdk_mac_Agora_Native_SDK_for_Mac_zhourui_26100_FULL_20251121_1643_987406_20251021_1427-3a.zip"
mac_sdk="https://download.agora.io/sdk/release/agora_sdk_mac_Agora_Native_SDK_for_Mac_zhourui_26101_FULL_20251121_2135_987705_20251021_1427-3a.zip"

#date: 20251211 update sdk for 2.4.2
linux_sdk="https://download.agora.io/sdk/release/agora_rtc_sdk_x86_64-linux-gnu-v4.4.32.154_26982_SERVER_20251210_1745_994155_20251021_1427-3a.zip"
mac_sdk="https://download.agora.io/sdk/release/agora_sdk_mac_v4.4.32.154_26308_FULL_20251210_1756_994156_20251021_1427-3a.zip"
#end usrl

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

check_and_download() {
    download_url=$1
    dst_dir=$2
    version_file="${dst_dir}/sdk_version"
    if [[ -f "${version_file}" ]]; then
        sdk_version=$(cat "${version_file}")
        if [[ "${sdk_version}" == $download_url ]]; then
            echo "${dst_dir} downloaded"
            return 0
        fi
    fi
    zip_file="${dst_dir}.zip"
    if [[ -f $zip_file ]]; then
        rm -f $zip_file
    fi
    if [[ -d $dst_dir ]]; then
        rm -rf $dst_dir
    fi
    echo "Downloading ${download_url} to ${zip_file}"
    curl -o $zip_file $download_url
    if [[ $? -ne 0 ]]; then
        echo "Failed to download ${download_url}"
        return 1
    fi
    unzip $zip_file
    rm -f $zip_file
    if [[ "${dst_dir}" != "agora_sdk" ]]; then
        mv agora_sdk $dst_dir
    fi
    echo $download_url > "${dst_dir}/sdk_version"
    echo "$dst_dir downloaded"
}

# Download the Agora RTC SDK
if [[ $OS == mac ]]; then
    check_and_download $mac_sdk agora_sdk_mac
fi
check_and_download $linux_sdk agora_sdk
