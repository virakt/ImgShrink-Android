# Installation Guide

## For Developers

### Prerequisites

1. **Go 1.21 or later**
   ```bash
   go version
   ```

2. **Docker** (required for fyne-cross)
   ```bash
   docker --version
   ```

3. **Git** (optional, for cloning)
   ```bash
   git version
   ```

### Step 1: Get the Code

If you have the source code already, skip to Step 2.

```bash
cd /path/to/imgshrink-mobile
```

### Step 2: Install Dependencies

```bash
go mod download
go mod tidy
```

### Step 3: Test on Desktop (Optional)

Before building for Android, test the app on your desktop:

```bash
go run .
```

This will launch the app in a desktop window. You can test all features except native Android file picking.

### Step 4: Install fyne-cross

```bash
go install fyne.io/fyne-cross@latest

# Add to PATH if needed
export PATH=$PATH:$(go env GOPATH)/bin
```

### Step 5: Build for Android

Choose one of these methods:

#### Method A: Using build script (Easiest)
```bash
./build.sh
```

#### Method B: Using Make
```bash
make build-android-fast
```

#### Method C: Manual build
```bash
fyne-cross android \
    -arch=arm64,amd64 \
    -app-id=com.virakt.imgshrink \
    -release \
    -ldflags="-s -w" \
    -tags=release
```

### Step 6: Locate the APK

After successful build, find your APK files:

```bash
# For ARM64 devices (most phones)
ls -lh fyne-cross/dist/android-arm64/ImgShrink.apk

# For x86_64 devices (emulators, some tablets)
ls -lh fyne-cross/dist/android-amd64/ImgShrink.apk
```

## For End Users

### Prerequisites

- Android device running Android 5.0 or later
- ARM64-v8a or x86_64 processor
- ~50 MB free storage

### Installation Steps

#### Option 1: Via ADB (Developer Method)

1. **Enable USB Debugging** on your Android device:
   - Go to Settings → About Phone
   - Tap "Build Number" 7 times
   - Go back to Settings → Developer Options
   - Enable "USB Debugging"

2. **Connect your device** to computer via USB

3. **Install the APK**:
   ```bash
   # For most devices (ARM64)
   adb install fyne-cross/dist/android-arm64/ImgShrink.apk
   
   # For emulators/x86 devices
   adb install fyne-cross/dist/android-amd64/ImgShrink.apk
   ```

#### Option 2: Direct Installation (User Method)

1. **Transfer APK** to your device:
   - Via USB cable
   - Via email
   - Via cloud storage
   - Via messaging app

2. **Enable Unknown Sources**:
   - Go to Settings → Security
   - Enable "Unknown Sources" or "Install Unknown Apps"
   - Grant permission to your file manager

3. **Install**:
   - Open the APK file on your device
   - Tap "Install"
   - Wait for installation to complete
   - Tap "Open" or find "ImgShrink" in your app drawer

4. **Grant Permissions**:
   - When first opening the app, grant storage permissions
   - This allows the app to read and save images

### First Run

1. Open ImgShrink
2. Tap "Select Image"
3. Choose an image from your gallery
4. Adjust quality and resize settings
5. Tap "Compress Image"
6. View results and find compressed image in the same folder

## Troubleshooting

### Build Issues

**"fyne-cross: command not found"**
```bash
go install fyne.io/fyne-cross@latest
export PATH=$PATH:$(go env GOPATH)/bin
```

**"Docker daemon not running"**
```bash
# Linux
sudo systemctl start docker

# macOS
open -a Docker

# Windows
Start Docker Desktop
```

**"Build takes forever"**
- First build downloads Docker images (~2GB)
- Subsequent builds are much faster
- Use `build-android-fast` for optimized builds

### Installation Issues

**"App not installed"**
- Check if you have enough storage
- Ensure you're using the correct architecture APK
- Try uninstalling any previous version

**"Parse error"**
- APK file may be corrupted
- Re-download or rebuild the APK
- Ensure your Android version is 5.0+

**"Unknown sources blocked"**
- Enable "Install Unknown Apps" for your file manager
- Settings → Apps → Special Access → Install Unknown Apps

### Runtime Issues

**"App crashes on startup"**
- Grant storage permissions
- Check Android version (minimum 5.0)
- Clear app data and restart

**"Can't select images"**
- Grant storage permissions in Settings → Apps → ImgShrink
- Check if you have images in your gallery

**"Compression fails"**
- Ensure image format is JPEG or PNG
- Check available storage space
- Try with a smaller image first

## Uninstallation

### Via Settings
1. Go to Settings → Apps
2. Find "ImgShrink"
3. Tap "Uninstall"

### Via ADB
```bash
adb uninstall com.virakt.imgshrink
```

## Updates

To update the app:
1. Build new version
2. Install new APK (will replace old version)
3. Your settings and data are preserved

## Support

For issues or questions:
- Check [README.md](README.md) for documentation
- Review [QUICKSTART.md](QUICKSTART.md) for build help
- See [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md) for technical details

## Security

- App only requests storage permissions
- No internet connection required
- No data collection or tracking
- All processing done locally on device
- Open source - review the code yourself

---

**Happy Compressing! 📱✨**
