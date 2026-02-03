# Quick Start Guide 🚀

## Build the Android App

### Option 1: Using the build script (Recommended)
```bash
./build.sh
```

### Option 2: Using Make
```bash
make all
```

### Option 3: Manual build
```bash
# Install fyne-cross
go install fyne.io/fyne-cross@latest

# Install dependencies
go mod download

# Build for Android
fyne-cross android -arch=arm64,amd64 -app-id=com.virakt.imgshrink -release -ldflags="-s -w"
```

## Install on Android Device

### Via ADB
```bash
# For ARM64 devices (most modern phones)
adb install fyne-cross/dist/android-arm64/ImgShrink.apk

# For x86_64 devices (emulators, some tablets)
adb install fyne-cross/dist/android-amd64/ImgShrink.apk
```

### Via File Transfer
1. Copy the APK to your device
2. Open the APK file on your device
3. Allow installation from unknown sources if prompted
4. Install the app

## Test on Desktop (Before Building for Android)

```bash
# Run the app on your desktop to test UI
go run .
```

## Build Optimizations

The build is optimized for:
- **Smaller binary size**: `-ldflags="-s -w"` strips debug info
- **Faster builds**: Only builds for arm64 and amd64 (no arm32)
- **Release mode**: Optimized compilation

## Troubleshooting

### "fyne-cross not found"
```bash
go install fyne.io/fyne-cross@latest
export PATH=$PATH:$(go env GOPATH)/bin
```

### "Docker not running"
fyne-cross requires Docker. Install and start Docker:
```bash
# On Linux
sudo systemctl start docker

# On macOS
open -a Docker
```

### Build takes too long
First build downloads Docker images and dependencies. Subsequent builds are much faster.

## App Features

- Select images from device storage
- Adjust compression quality (1-100%)
- Resize images (10-100%)
- Real-time preview
- View compression statistics
- Save compressed images

## Supported Devices

- **Minimum Android**: 5.0 (API 21)
- **Architectures**: ARM64-v8a, x86_64
- **Permissions**: Storage access for reading/writing images

## Next Steps

1. Build the app using one of the methods above
2. Install on your Android device
3. Grant storage permissions
4. Start compressing images!

For more details, see [README.md](README.md)
