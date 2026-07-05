# flac2mp3

A fast, cross-platform CLI tool written in Go that converts FLAC audio files to MP3 (320 kbps) using [FFmpeg](https://ffmpeg.org/). Supports single-file and batch directory conversion, with concurrent processing for speed.

## Features

- Convert a single `.flac` file to `.mp3`
- Batch convert an entire folder of FLAC files
- Up to 4 concurrent conversions for faster batch processing
- Outputs to an `mp3/` sub-folder alongside the source files
- Built-in `doctor` command to verify your FFmpeg installation
- Cross-platform: Linux, Windows, macOS (Intel & Apple Silicon)

## Requirements

- [Go](https://go.dev/) 1.18+ (to build from source)
- [FFmpeg](https://ffmpeg.org/download.html) installed and available in your `PATH`

## Installation

### Download a pre-built binary

Pre-built binaries for all platforms are produced by the build scripts and placed in the `dist/` folder.

| Platform         | Binary                          |
|------------------|---------------------------------|
| Linux (amd64)    | `dist/flac2mp3-linux-amd64`     |
| Windows (amd64)  | `dist/flac2mp3-windows-amd64.exe` |
| macOS (Intel)    | `dist/flac2mp3-darwin-amd64`    |
| macOS (Apple Silicon) | `dist/flac2mp3-darwin-arm64` |

### Build from source

**Linux / macOS**
```bash
chmod +x build.sh
./build.sh
```

**Windows**
```bat
build.bat
```

Binaries are written to the `dist/` directory.

### Run directly with Go

```bash
go run main.go <args>
```

## Usage

```
flac2mp3 doctor                  # Check if ffmpeg is installed and working
flac2mp3 <path/to/file.flac>     # Convert a single FLAC file
flac2mp3 <path/to/folder>        # Convert all FLAC files in a folder
```

### Examples

Check that FFmpeg is available:
```bash
flac2mp3 doctor
```

Convert a single file:
```bash
flac2mp3 ~/Music/track.flac
```

Batch convert a folder:
```bash
flac2mp3 ~/Music/AlbumFolder/
```

Output is always saved to an `mp3/` sub-folder inside the source directory. For example:

```
~/Music/AlbumFolder/
├── track01.flac
├── track02.flac
└── mp3/
    ├── track01.mp3
    └── track02.mp3
```

## How it works

- Single file: creates an `mp3/` directory next to the source file and runs `ffmpeg -ab 320k`.
- Directory: scans for all `.flac` files, creates an `mp3/` sub-folder, then converts up to 4 files concurrently using goroutines and a buffered channel as a semaphore.
- All conversions use 320 kbps bitrate and will overwrite existing output files (`-y`).

## License

MIT — see [LICENSE](LICENSE).
