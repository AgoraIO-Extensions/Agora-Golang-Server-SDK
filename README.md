# build sdk
## prepare C version of agora rtc sdk
- make agora_sdk directory in the same directory with go_wrapper
- there should be libagora_rtc_sdk.so and include_c in agora_sdk directory
```
unzip agora_sdk.zip
```

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
