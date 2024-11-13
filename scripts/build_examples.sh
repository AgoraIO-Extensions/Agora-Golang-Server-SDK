#!/bin/bash

PACKAGE_HOME=$(
    cd $(dirname $0)/..
    pwd
)
examples_path="${PACKAGE_HOME}/go_sdk/examples"

# Check if at least one argument is passed
if [ "$#" -eq 0 ]; then
  echo "No files provided. Usage: $0 file1 file2 ..."
  exit 1
fi

# Iterate over each file passed as an argument
for file in "$@"; 
do
  example_path="$examples_path/$file"
  if [ -d $example_path ]; then
    echo "building example: $file"
    go mod tidy -C $example_path
    go build -C $example_path -o "${PACKAGE_HOME}/bin/${file}"
  else
    echo "File not found: $file"
  fi
done