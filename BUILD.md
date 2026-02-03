 # Build Instructions

## Quick Start with Makefile (Recommended)

Use the Makefile for easy builds on all platforms:

```bash
# Show available targets
make help

# Build Linux binary
make linux

# Build Android APKs
make android

# Build everything
make build-all
```

## Prerequisites

1. **Go 1.21+**
2. **Docker** (for Android builds)
3. **fyne-cross**: `go install github.com/fyne-io/fyne-cross@latest`

## Desktop Build (Linux)

### Using Makefile
```bash
make linux          # Build with version
make linux-simple  # Build simple binary
make run           # Build and run
```

### Manual Build
```bash
go build -ldflags="-s -w" -o imgshrink-mobile
```

### Using Fyne CLI (Cross-platform)
```bash
# Windows
fyne package -os windows -icon Icon.png

# macOS
fyne package -os darwin -icon Icon.png
```

## Android Build

### Using Makefile (Recommended)
```bash
# Build ARM64 APK
make android-arm64

# Build x86_64 APK
make android-amd64

# Build both architectures
make android

# Build and sign APKs (requires keystore)
make android-release
```

### Manual Build with fyne-cross
```bash
# Build for ARM64
fyne-cross android -arch=arm64 -app-id=com.imgshrink.mobile -icon Icon.png

# Build for x86_64
fyne-cross android -arch=amd64 -app-id=com.imgshrink.mobile -icon Icon.png

# Build for both
fyne-cross android -arch=arm64,amd64 -app-id=com.imgshrink.mobile -icon Icon.png
```

### Using Fyne CLI (Alternative)
```bash
fyne package -os android/arm64 -appID com.imgshrink.mobile -icon Icon.png
```

## Build Optimizations

### Reduce Binary Size
```bash
go build -ldflags="-s -w" -o imgshrink-mobile
```

### Use UPX Compression (Optional)
```bash
upx --best --lzma imgshrink-mobile
```

### Faster Builds
```bash
# Enable Go build cache
export GOCACHE=$HOME/.cache/go-build

# Use parallel builds
go build -p 4
```

## Build Artifacts

### Desktop
- **Linux**: `imgshrink-mobile-VERSION-linux-amd64` (executable)
- **Windows**: `ImgShrink.exe`
- **macOS**: `ImgShrink.app`

### Android
- **ARM64 APK**: `fyne-cross/dist/android-arm64/ImgShrink.apk`
- **x86_64 APK**: `fyne-cross/dist/android-amd64/ImgShrink.apk`
- **Size**: ~44MB for ARM64

## Troubleshooting

### Android Build Issues

1. **fyne-cross not found**:
   ```bash
   make check-fyne-cross
   ```

2. **Docker permission errors**:
   ```bash
   sudo usermod -aG docker $USER
   newgrp docker
   ```

3. **Out of memory during build**:
   - Increase Docker memory limit
   - Build single architecture: `make android-arm64`

### Desktop Build Issues

1. **Missing dependencies on Linux**:
   ```bash
   # Ubuntu/Debian
   sudo apt-get install gcc libgl1-mesa-dev xorg-dev
   
   # Fedora
   sudo dnf install gcc mesa-libGL-devel libXcursor-devel libXrandr-devel libXinerama-devel libXi-devel
   ```

2. **CGO errors**:
   - Ensure GCC/Clang is installed
   - On Windows, install TDM-GCC or MinGW-w64

## Testing the APK

### Install on Device
```bash
adb install fyne-cross/dist/android-arm64/ImgShrink.apk
```

### Install on Emulator
```bash
adb -e install fyne-cross/dist/android-arm64/ImgShrink.apk
```

### Check APK Info
```bash
aapt dump badging fyne-cross/dist/android-arm64/ImgShrink.apk
```

## Release Build

### Generate a Keystore (First Time)
```bash
keytool -genkey -v -keystore imgshrink.keystore -alias imgshrink -keyalg RSA -keysize 2048 -validity 10000
```

### Build and Sign with Makefile
```bash
make android-release
```

### Manual Signing with fyne-cross
```bash
fyne-cross android -arch=arm64,amd64 \
  -app-id=com.imgshrink.mobile \
  -release \
  -keystore imgshrink.keystore \
  -key-name imgshrink \
  -keystore-pass YOUR_PASSWORD \
  -key-pass YOUR_PASSWORD \
  -icon Icon.png
```

## Build Times

Approximate build times on a modern system:

- **Desktop (Linux)**: 10-15 seconds
- **Android (ARM64)**: 2-3 minutes (first build), 30-60 seconds (cached)
- **Android (ARM64 + x86_64)**: 4-5 minutes (first build), 1-2 minutes (cached)
