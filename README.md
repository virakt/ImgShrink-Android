# ImgShrink Mobile 🖼️📱

Android app for compressing images with a minimal glassy UI inspired by iOS design.

## Features

- **Native File Selection** - Pick images from device storage using Android content URIs
- **Image Preview** - View selected images before compression
- **Quality Control** - Adjust JPEG quality (1-100) with real-time slider
- **Resolution Control** - Set custom width/height or resize by percentage
- **Output Location** - Automatically saves to app's temp directory (Android storage restrictions)
- **Real-time Stats** - See original vs compressed size and reduction percentage
- **Dark Glassy UI** - Pitch black with moonish blue tints and frosted glass effects
- **Batch Processing** - Compress multiple images at once

## Screenshots

[Coming soon]

## Build

### Quick Start

```bash
# Install dependencies
go install fyne.io/fyne/v2/cmd/fyne@latest
go install github.com/fyne-io/fyne-cross@latest

# Build for Android (ARM64)
fyne-cross android -arch=arm64 -app-id=com.imgshrink.mobile

# The APK will be in the project root as ImgShrink.apk
```

### Detailed Build Instructions

See [BUILD.md](BUILD.md) for comprehensive build instructions including:
- Desktop builds (Linux, Windows, macOS)
- Android builds (ARM64, x86_64)
- Build optimizations
- Troubleshooting
- Release signing

### Requirements

- **Go 1.21+**
- **Fyne CLI v2.7.2+**
- **Docker** (for Android builds with fyne-cross)
- **fyne-cross** (recommended for Android builds)

### Build Desktop (for testing)

```bash
go build -o imgshrink-mobile
./imgshrink-mobile
```

## Architecture

- **Supported Architectures**: ARM64-v8a, x86_64
- **APK Size**: ~44MB (ARM64 only)
- **Minimum Android**: API 19 (Android 4.4 KitKat)
- **Target Android**: API 34 (Android 14)

## Technology Stack

- **Language**: Go 1.21+
- **UI Framework**: Fyne v2.7.2
- **Image Processing**: 
  - `github.com/virakt/imgshrink` (core compression API)
  - `github.com/disintegration/imaging` (image manipulation)
  - `github.com/nfnt/resize` (resizing)
- **Build Tool**: fyne-cross (Docker-based cross-compilation)

## Usage

1. **Launch the app** on your Android device
2. **Tap "Select Image"** to choose an image from your device
3. **Adjust settings**:
   - Quality slider (1-100)
   - Width/Height inputs for custom resolution
   - Or use percentage resize
4. **Tap "Compress Image"** to process
5. **View results** showing original size, compressed size, and reduction percentage
6. **Compressed image** is saved to app's temp directory

## Android-Specific Features

### File Handling

The app uses Android's Storage Access Framework (SAF) to:
- Read images from content URIs (not direct file paths)
- Handle scoped storage restrictions (Android 10+)
- Save compressed images to app-specific temp directory

### Permissions

No special permissions required! The app uses:
- `READ_EXTERNAL_STORAGE` (implicit via SAF)
- App-specific storage (no permission needed)

## Development

### Project Structure

```
imgshrink-mobile/
├── main.go              # Main application code
├── Icon.png             # App icon (512x512)
├── FyneApp.toml         # Fyne app metadata
├── go.mod               # Go dependencies
├── README.md            # This file
├── BUILD.md             # Detailed build instructions
└── fyne-cross/          # Build artifacts
    └── dist/
        └── android-arm64/
            └── ImgShrink.apk
```

### Key Components

1. **UI Layout**: Vertical box with cards for each section
2. **Image Selection**: Native file picker with URI handling
3. **Compression Engine**: Reuses imgshrink API from parent project
4. **Settings**: Quality slider, dimension inputs, percentage resize
5. **Results Display**: Shows before/after stats with visual feedback

## Roadmap

- [ ] Add image preview thumbnails
- [ ] Support batch compression
- [ ] Add more output format options (WebP, AVIF)
- [ ] Implement share functionality
- [ ] Add compression presets (Low, Medium, High, Custom)
- [ ] Support for GIF and WebP input formats
- [ ] Dark/Light theme toggle
- [ ] Compression history

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.

## Acknowledgments

- [Fyne](https://fyne.io/) - Cross-platform UI framework
- [ImgShrink](https://github.com/virakt/imgshrink) - Core compression library
- iOS design inspiration for the glassy UI aesthetic

## Support

For issues, questions, or suggestions, please open an issue on GitHub.
