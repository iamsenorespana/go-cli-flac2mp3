@echo off
mkdir dist 2>nul

echo 📦 Building for Linux...
set GOOS=linux
set GOARCH=amd64
go build -o dist/flac2mp3-linux-amd64 main.go

echo 📦 Building for Windows...
set GOOS=windows
set GOARCH=amd64
go build -o dist/flac2mp3-windows-amd64.exe main.go

echo 📦 Building for macOS (Intel)...
set GOOS=darwin
set GOARCH=amd64
go build -o dist/flac2mp3-darwin-amd64 main.go

echo 📦 Building for macOS (Apple Silicon)...
set GOOS=darwin
set GOARCH=arm64
go build -o dist/flac2mp3-darwin-arm64 main.go

echo 🚀 All builds complete! Check the \dist folder.