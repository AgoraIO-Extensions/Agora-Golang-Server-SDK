#!/bin/bash

# set error exit
set -e

# Create output directory
if [ -d "bin" ]; then
    echo "bin dir already exists, checking for existing rtmdemo file..."
    if [ -f "./bin/rtmdemo" ]; then
        echo "Removing existing rtmdemo file..."
        rm ./bin/rtmdemo
    else
        echo "rtmdemo file does not exist, skipping deletion"
    fi
else
    mkdir -p bin
fi

# build command
echo "Building example application..."
go build -o ./bin/rtmdemo ./cmd/example/



# check if build is successful
if [ $? -eq 0 ]; then
    echo "Build successful! Binary is located at ./bin/example"
else
    echo "Build failed!"
    exit 1
fi 