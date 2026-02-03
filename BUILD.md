# Build Instructions

## Prerequisites

1. **Go 1.21+**
2. **Fyne CLI**: `go install fyne.io/fyne/v2/cmd/fyne@latest`
3. **Docker** (for Android builds)
4. **fyne-cross**: `go install github.com/fyne-io/fyne-cross@latest`

## Desktop Build

### Linux
```bash
go build -o imgshrink-mobile
```

### Windows
```bash
fyne package -os windows -icon Icon.png
```

### macOS
```bash
fyne package -os darwin -icon Icon.png
```

## Android Build

### Using fyne-cross (Recommended)

Build for ARM64 (most modern Android devices):
```bash
fyne-cross android -arch=arm64 -app-id=com.imgshrink.mobile
```

Build for both ARM64 and x86_64:
```bash
fyne-cross android -arch=arm64,amd64 -app-id=com.imgshrink.mobile
```

The APK will be generated in the project root directory as `ImgShrink.apk`.

### Using fyne package (Alternative)

```bash
fyne package -os android/arm64 -appID com.imgshrink.mobile -icon Icon.png
```

**Note**: This requires Android SDK and NDK to be installed locally.

## Build Optimizations

### Reduce Build Size

1. **Strip debug symbols**:
```bash
go build -ldflags="-s -w" -o imgshrink-mobile
```

2. **Use UPX compression** (optional):
```bash
upx --best --lzma imgshrink-mobile
```

### Faster Builds

1. **Enable Go build cache**:
```bash
export GOCACHE=$HOME/.cache/go-build
```

2. **Use parallel builds**:
```bash
go build -p 4
```

3. **For fyne-cross, reuse Docker cache**:
The cache is automatically stored in `~/.cache/fyne-cross`

## Build Artifacts

### Desktop
- **Linux**: `imgshrink-mobile` (executable)
- **Windows**: `ImgShrink.exe`
- **macOS**: `ImgShrink.app`

### Android
- **APK**: `ImgShrink.apk` (44MB for ARM64)
- **Supported architectures**: ARM64-v8a, x86_64

## Troubleshooting

### Android Build Issues

1. **"No such file or directory" for APK**:
   - The APK is built in the project root, not in fyne-cross/tmp
   - Check for `ImgShrink.apk` in the current directory

2. **Docker permission errors**:
   ```bash
   sudo usermod -aG docker $USER
   newgrp docker
   ```

3. **Out of memory during build**:
   - Increase Docker memory limit in Docker Desktop settings
   - Or build for single architecture at a time

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
adb install ImgShrink.apk
```

### Install on Emulator
```bash
adb -e install ImgShrink.apk
```

### Check APK Info
```bash
aapt dump badging ImgShrink.apk
```

## Release Build

For production release:

1. **Generate a keystore** (first time only):
```bash
keytool -genkey -v -keystore imgshrink.keystore -alias imgshrink -keyalg RSA -keysize 2048 -validity 10000
```

2. **Sign the APK**:
```bash
jarsigner -verbose -sigalg SHA256withRSA -digestalg SHA-256 -keystore imgshrink.keystore ImgShrink.apk imgshrink
```

3. **Align the APK**:
```bash
zipalign -v 4 ImgShrink.apk ImgShrink-aligned.apk
```

## Build Times

Approximate build times on a modern system:

- **Desktop (Linux)**: 10-15 seconds
- **Android (ARM64)**: 2-3 minutes (first build), 30-60 seconds (cached)
- **Android (ARM64 + x86_64)**: 4-5 minutes (first build), 1-2 minutes (cached)
