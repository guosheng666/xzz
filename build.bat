@echo off
setlocal

set APP_NAME=mysql-sync-tool
set VERSION=1.0.0

echo Building %APP_NAME% v%VERSION% for all platforms...

mkdir bin 2>nul

echo Building for Windows x86...
set GOOS=windows
set GOARCH=386
go build -ldflags="-s -w" -o bin/%APP_NAME%-windows-x86.exe .
if %ERRORLEVEL% EQU 0 echo   Windows x86: OK

echo Building for Windows x64...
set GOARCH=amd64
go build -ldflags="-s -w" -o bin/%APP_NAME%-windows-x64.exe .
if %ERRORLEVEL% EQU 0 echo   Windows x64: OK

echo Building for Linux x86...
set GOOS=linux
set GOARCH=386
go build -ldflags="-s -w" -o bin/%APP_NAME%-linux-x86 .
if %ERRORLEVEL% EQU 0 echo   Linux x86: OK

echo Building for Linux x64...
set GOARCH=amd64
go build -ldflags="-s -w" -o bin/%APP_NAME%-linux-x64 .
if %ERRORLEVEL% EQU 0 echo   Linux x64: OK

echo Building for Linux arm64...
set GOARCH=arm64
go build -ldflags="-s -w" -o bin/%APP_NAME%-linux-arm64 .
if %ERRORLEVEL% EQU 0 echo   Linux arm64: OK

echo Building for Linux loongarch64...
set GOARCH=loongarch64
go build -ldflags="-s -w" -o bin/%APP_NAME%-linux-loongarch64 .
if %ERRORLEVEL% EQU 0 echo   Linux loongarch64: OK

echo.
echo Build complete! Binaries are in the bin directory.
dir bin

endlocal
