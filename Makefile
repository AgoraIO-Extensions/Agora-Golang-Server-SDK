# Makefile for Go project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
CURRENT_PATH=$(shell pwd)
UNAME_S := $(shell uname -s)
OS := unknown
BASIC_EXAMPLES := send_recv_pcm send_recv_yuv recv_h264 recv_pcm_loopback sample_vad send_recv_yuv_pcm ai_send_recv_pcm
BASIC_EXAMPLES =  send_recv_pcm send_recv_yuv recv_h264 recv_pcm_loopback sample_vad send_recv_yuv_pcm ai_send_recv_pcm example_externalaudio_processor
ADVANCED_EXAMPLES := multi_cons_rtx send_encoded_audio send_h264 send_mp4
ADVANCED_EXAMPLES =  send_encoded_audio send_h264 send_mp4 multi_cons_rtx

ifeq ($(UNAME_S),Linux)
    OS := linux
else ifeq ($(UNAME_S),Darwin)
    OS := mac
# else ifeq ($(UNAME_S),CYGWIN) # Cygwin is Unix-like environment on Windows
#     OS := windows
# else ifeq ($(UNAME_S),MINGW32) # MinGW is a native Windows port of GNU tools
#     OS := windows
# else ifeq ($(UNAME_S),MSYS) # MSYS is a collection of GNU utilities for Windows
#     OS := windows
else
    $(error Unsupported OS: $(UNAME_S))
endif

# Install dependencies
.PHONY: deps
deps:
	./scripts/install_agora_sdk.sh
	$(GOGET) -v ./...

# Build the project
.PHONY: build
build:
	mkdir -p $(CURRENT_PATH)/bin
	$(GOBUILD) -C $(CURRENT_PATH)/go_sdk/rtc -v -o $(CURRENT_PATH)/bin/rtc.a

# Clean the project
.PHONY: clean
clean:
	rm -rf $(CURRENT_PATH)/bin
	$(GOCLEAN) -v

# Test the project
.PHONY: test
test:
	$(GOTEST) -v ./...

# Install the project
.PHONY: install
install: deps
	$(GOCMD) install -C $(CURRENT_PATH)/go_sdk/rtc -v

# Build examples
.PHONY: examples
examples:
	./scripts/build_examples.sh $(BASIC_EXAMPLES)

.PHONY: advanced-examples
advanced-examples:
	./scripts/build_examples.sh $(ADVANCED_EXAMPLES)