# ImgShrink Mobile 📱

A beautiful Android image compression app built with Fyne, featuring an iOS-inspired glassy dark UI with moonish blue tints.

## Features

- 🎨 **Glassy Dark UI** - Pitch black background with moonish blue accents
- 📱 **Mobile-First Design** - Optimized for Android devices
- 🖼️ **Native File Selection** - Pick images from your device
- 👁️ **Image Preview** - See your image before compression
- ⚙️ **Adjustable Settings** - Control quality and resize percentage
- 📊 **Real-time Stats** - View compression results instantly
- 🚀 **Fast Compression** - Efficient JPEG and PNG compression

## Supported Formats

- **JPEG** (.jpg, .jpeg)
- **PNG** (.png)

## Supported Architectures

- **ARM64-v8a** (64-bit ARM)
- **x86_64** (64-bit Intel/AMD)

## Building

### Prerequisites

1. Install Go 1.21 or later
2. Install fyne-cross for Android builds:
   ```bash
   make install-fyne-cross
   ```

### Build Commands

```bash
# Install dependencies
make deps

# Build optimized Android APK (recommended)
make build-android-fast

# Build standard Android APK
make build-android

# Build debug version
make build-android-debug

# Test on desktop
make run

# Clean build artifacts
make clean

# Full build pipeline
make all
```

### Manual Build

```bash
# Install dependencies
go mod download

# Build for Android
fyne-cross android -arch=arm64,amd64 -app-id=com.virakt.imgshrink -release -ldflags="-s -w"
```

## Build Output

The APK will be generated in:
```
fyne-cross/dist/android-arm64/
fyne-cross/dist/android-amd64/
```

## Installation

1. Transfer the APK to your Android device
2. Enable "Install from Unknown Sources" in Settings
3. Install the APK
4. Grant storage permissions when prompted

## Usage

1. **Select Image** - Tap "Select Image" to choose a photo
2. **Adjust Quality** - Use the quality slider (1-100%)
3. **Resize (Optional)** - Use the resize slider to reduce dimensions
4. **Compress** - Tap "Compress Image" to process
5. **View Results** - See compression stats and saved file location

## UI Design

The app features a minimal, glassy interface inspired by iOS design:

- **Color Palette**:
  - Background: Pitch black (#0A0C12)
  - Surface: Dark blue-tinted (#141923)
  - Primary: Moonish blue (#6478B4)
  - Accent: Light blue (#8CA0DC)
  - Text: Off-white (#F0F5FF)

- **Design Elements**:
  - Rounded corners
  - Subtle transparency
  - Smooth animations
  - Clean typography
  - Generous padding

## Performance Optimizations

- Compiled with `-ldflags="-s -w"` for smaller binary size
- Optimized image processing algorithms
- Efficient memory management
- Fast compression with quality preservation

## Integration with ImgShrink API

This mobile app uses the same compression engine as the desktop ImgShrink tool, ensuring consistent quality and performance across platforms.

## License

MIT License - See LICENSE file for details

## Credits

Built with:
- [Fyne](https://fyne.io/) - Cross-platform GUI toolkit
- [imaging](https://github.com/disintegration/imaging) - Image processing library
- ImgShrink compression engine
