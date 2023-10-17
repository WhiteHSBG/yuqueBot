#!/bin/bash

# 检查参数个数是否正确
if [ $# -ne 2 ]; then
    echo "Usage: $0 <input_file> <output_file>"
    exit 1
fi

input_file=$1
output_file=$2

# Windows x64
GOOS=windows GOARCH=amd64 go build -trimpath -o ${output_file}-windows-amd64.exe ${input_file}

# Linux x64
GOOS=linux GOARCH=amd64 go build -trimpath -o ${output_file}-linux-amd64 ${input_file}

# Linux arm64
GOOS=linux GOARCH=arm64 go build -trimpath -o ${output_file}-linux-arm64 ${input_file}

# macOS arm64
GOOS=darwin GOARCH=arm64 go build -trimpath -o ${output_file}-macos-arm64 ${input_file}

# macOS amd64
GOOS=darwin GOARCH=amd64 go build -trimpath -o ${output_file}-macos-amd64 ${input_file}