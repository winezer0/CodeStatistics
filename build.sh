#!/bin/bash

set -e

echo "Building CodeStatistics..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed or not in PATH"
    exit 1
fi

# Clean previous builds
if [ -d "bin" ]; then
    rm -rf bin
fi
mkdir -p bin

# Get dependencies
echo "Getting dependencies..."
go mod tidy
if [ $? -ne 0 ]; then
    echo "Error: Failed to get dependencies"
    exit 1
fi

# Set common build flags
export CGO_ENABLED=0
BUILD_FLAGS='-ldflags "-s -w -X main.version=1.0.0 -X main.buildTime= -X main.gitCommit="'

# Build for Linux x64
echo "Building for Linux x64..."
GOOS=linux GOARCH=amd64 go build $BUILD_FLAGS -o bin/CodeStatistics_linux main.go
if [ $? -ne 0 ]; then
    echo "Error: Failed to build for Linux x64"
    exit 1
fi

# Build for Windows x64
echo "Building for Windows x64..."
GOOS=windows GOARCH=amd64 go build $BUILD_FLAGS -o bin/CodeStatistics.exe main.go
if [ $? -ne 0 ]; then
    echo "Error: Failed to build for Windows x64"
    exit 1
fi

echo ""
echo "Build completed successfully!"
echo "Output files:"
ls -la bin/
echo ""
echo "Usage: ./bin/CodeStatistics_linux -h"
