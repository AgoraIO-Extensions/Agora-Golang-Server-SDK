#!/bin/sh

export CGO_LDFLAGS_ALLOW="-Wl,-rpath,.*"
export CGO_LDFLAGS="-Wl,-rpath,../../../agora_sdk_mac"
go mod tidy
go build main