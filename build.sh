#!/bin/bash

# Create output directory
mkdir -p dist

echo "📦 Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o dist/flac2mp3-linux-amd64 main.go

echo "📦 Building for Windows..."
GOOS=windows GOARCH=amd64 go build -o dist/flac2mp3-windows-amd64.exe main.go

echo "📦 Building for macOS (Intel)..."
GOOS=darwin GOARCH=amd64 go build -o dist/flac2mp3-darwin-amd64 main.go

echo "📦 Building for macOS (Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build -o dist/flac2mp3-darwin-arm64 main.go

echo "🚀 All builds complete! Check the /dist folder."