@echo off
setlocal enabledelayedexpansion

echo Building CodeStatistics...

REM Check if Go is installed
go version >nul 2>&1
if errorlevel 1 (
    echo Error: Go is not installed or not in PATH
    exit /b 1
)

REM Clean previous builds
if exist "bin" rmdir /s /q bin
mkdir bin

REM Get dependencies
echo Getting dependencies...
go mod tidy
if errorlevel 1 (
    echo Error: Failed to get dependencies
    exit /b 1
)

REM Set common build flags
set CGO_ENABLED=0
set BUILD_FLAGS=-ldflags "-s -w -X 'main.version=1.0.0' -X 'main.buildTime=' -X 'main.gitCommit='"

REM Build for Windows x64
echo Building for Windows x64...
set GOOS=windows
set GOARCH=amd64
go build %BUILD_FLAGS% -o bin\CodeStatistics.exe main.go
if errorlevel 1 (
    echo Error: Failed to build for Windows x64
    exit /b 1
)

REM Build for Linux x64
echo Building for Linux x64...
set GOOS=linux
set GOARCH=amd64
go build %BUILD_FLAGS% -o bin\CodeStatistics_linux main.go
if errorlevel 1 (
    echo Error: Failed to build for Linux x64
    exit /b 1
)

echo.
echo Build completed successfully!
echo Output files:
dir /b bin\*
echo.
echo Usage: bin\CodeStatistics.exe -h
