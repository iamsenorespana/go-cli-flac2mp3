package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	// 1. Run the "Doctor" check first
	if len(os.Args) > 1 && os.Args[1] == "doctor" {
		runDoctorCheck()
		return
	}

	// Validate basic arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  go run main.go doctor                 # Check if ffmpeg is installed")
		fmt.Println("  go run main.go <path/to/file.flac>    # Convert a single file")
		fmt.Println("  go run main.go <path/to/folder>       # Convert all FLACs in a folder")
		os.Exit(1)
	}

	targetPath := os.Args[1]
	
	// Ensure ffmpeg is available before processing any files
	if !isFFmpegInstalled() {
		fmt.Println("❌ Error: ffmpeg is not installed or not in your PATH.")
		fmt.Println("Run 'go run main.go doctor' for more details.")
		os.Exit(1)
	}

	// 2. Check if the input is a file or a directory
	fileInfo, err := os.Stat(targetPath)
	if err != nil {
		fmt.Printf("❌ Error accessing path: %v\n", err)
		os.Exit(1)
	}

	if fileInfo.IsDir() {
		processDirectory(targetPath)
	} else {
		processSingleFile(targetPath)
	}
}

// runDoctorCheck verifies if ffmpeg is accessible and prints its version
func runDoctorCheck() {
	fmt.Println("🔍 Running system doctor check...")
	path, err := exec.LookPath("ffmpeg")
	if err != nil {
		fmt.Println("❌ ffmpeg: NOT FOUND")
		fmt.Println("   Please install FFmpeg and ensure it is added to your system's PATH environmental variable.")
		return
	}
	fmt.Printf("✅ ffmpeg: FOUND at %s\n", path)

	// Get version info just to be sure it responds correctly
	cmd := exec.Command("ffmpeg", "-version")
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("⚠️  ffmpeg found, but failed to execute: %v\n", err)
		return
	}
	// Print just the first line of the version output
	firstLine := strings.Split(string(out), "\n")[0]
	fmt.Printf("   Details: %s\n", firstLine)
}

// isFFmpegInstalled is a quick boolean check used before running conversions
func isFFmpegInstalled() bool {
	_, err := exec.LookPath("ffmpeg")
	return err == nil
}

// processSingleFile converts a single file and outputs it to an 'mp3' folder in the same directory
func processSingleFile(filePath string) {
	if !strings.HasSuffix(strings.ToLower(filePath), ".flac") {
		fmt.Println("❌ Error: The file must have a .flac extension.")
		return
	}

	dir := filepath.Dir(filePath)
	mp3Dir := filepath.Join(dir, "mp3")

	// Create the mp3 directory if it doesn't exist
	if err := os.MkdirAll(mp3Dir, os.ModePerm); err != nil {
		fmt.Printf("❌ Error creating mp3 directory: %v\n", err)
		return
	}

	fileName := filepath.Base(filePath)
	convertFile(dir, mp3Dir, fileName)
}

// processDirectory scans a folder for FLACs and batch converts them concurrently
func processDirectory(dirPath string) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("❌ Error reading directory: %v\n", err)
		return
	}

	// Filter for FLAC files
	var flacFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".flac") {
			flacFiles = append(flacFiles, file.Name())
		}
	}

	if len(flacFiles) == 0 {
		fmt.Println("📂 No FLAC files found in this directory.")
		return
	}

	mp3Dir := filepath.Join(dirPath, "mp3")
	if err := os.MkdirAll(mp3Dir, os.ModePerm); err != nil {
		fmt.Printf("❌ Error creating mp3 directory: %v\n", err)
		return
	}

	fmt.Printf("🚀 Found %d FLAC files. Starting conversion...\n", len(flacFiles))

	// Concurrency control
	var wg sync.WaitGroup
	maxGoroutines := 4 // Adjust based on your CPU cores
	guard := make(chan struct{}, maxGoroutines)

	for _, fileName := range flacFiles {
		wg.Add(1)
		guard <- struct{}{} // Block if channel is full

		go func(name string) {
			defer wg.Done()
			defer func() { <-guard }() // Free up space in channel
			convertFile(dirPath, mp3Dir, name)
		}(fileName)
	}

	wg.Wait()
	fmt.Printf("🎉 All files converted successfully! Saved to: %s\n", mp3Dir)
}

// convertFile handles the actual execution of the ffmpeg command
func convertFile(sourceDir, outputDir, fileName string) {
	inputPath := filepath.Join(sourceDir, fileName)
	
	// Swap .flac extension for .mp3
	mp3Name := strings.TrimSuffix(fileName, filepath.Ext(fileName)) + ".mp3"
	outputPath := filepath.Join(outputDir, mp3Name)

	fmt.Printf("⏳ Converting: %s -> mp3/%s\n", fileName, mp3Name)

	// -i input -ab bitrate -y (overwrite existing)
	cmd := exec.Command("ffmpeg", "-i", inputPath, "-ab", "320k", outputPath, "-y")
	
	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Failed to convert %s: %v\n", fileName, err)
	} else {
		fmt.Printf("✅ Finished: %s\n", mp3Name)
	}
}