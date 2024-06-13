# build sdk
## required OS and go version
- supported linux version: 
  - Ubuntu 18.04 LTS and above
  - CentOS 7.0 and above
- go version:
  - go 1.20 and above
  - not tested on go 1.19 and below

## prepare C version of agora rtc sdk
- download and unzip [agora_sdk.zip](https://share.weiyun.com/YfMpLtfR)
```
unzip agora_sdk.zip
```
- make **agora_sdk** directory in the same directory with **go_wrapper**
- there should be **libagora_rtc_sdk.so** and **include_c** in **agora_sdk** directory

## build sample
```
cd go_wrapper
go build main
```

# test
- open premium, input channel name with **lhztest**, and join channel
- run main
```
export LD_LIBRARY_PATH=/path/to/agora_sdk
cd go_wrapper
./main
```
