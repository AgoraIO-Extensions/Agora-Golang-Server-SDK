# Makefile for Go project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
EXAMPLES_PATH=./go_sdk/examples
CURRENT_PATH=$(shell pwd)
UNAME_S := $(shell uname -s)
OS := unknown

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
	$(GOBUILD) -C $(CURRENT_PATH)/go_sdk/agoraservice -v -o $(CURRENT_PATH)/bin/agoraservice.a

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
	$(GOCMD) install -C $(CURRENT_PATH)/go_sdk/agoraservice -v

# Build examples
.PHONY: examples
examples:
	for dir in $(EXAMPLES_PATH)/*; do \
		if [ -d "$$dir" ]; then \
			go mod tidy -C $$dir; \
			$(GOBUILD) -C $$dir -o $(CURRENT_PATH)/bin/$$(basename $$dir); \
		fi \
	done