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



#ver2.2.1
#linux_sdk="https://download.agora.io/sdk/release/agora_rtc_sdk-x86_64-linux-gnu-v4.4.31-20241223_111509-491956-aed.zip"
#mac_sdk="https://download.agora.io/sdk/release/agora_sdk_mac_v4.4.31_23136_FULL_20241223_1245_492039-aed.zip"
# ver2.2.0
#linux_sdk="https://download.agora.io/sdk/release/agora_rtc_sdk-x86_64-linux-gnu-v4.4.31-20241223_111509-491956.zip"
#mac_sdk="https://download.agora.io/sdk/release/agora_sdk_mac_v4.4.31_23136_FULL_20241223_1245_492039.zip"

# old version
#linux_sdk="https://download.agora.io/sdk/release/agora_rtc_sdk-x86_64-linux-gnu-v4.4.30-20241024_101940-398537.zip"
#mac_sdk="https://download.agora.io/sdk/release/agora_rtc_sdk_mac_rel.v4.4.30_22472_FULL_20241024_1224_398653.zip"

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
