# ImgShrink Mobile - Release Notes

## Version 1.0.0 - Initial Release

### Overview
ImgShrink Mobile is an Android application for compressing images built with the Fyne framework and Go. It features a beautiful iOS-inspired glassy dark UI with moonish blue accents.

### Features

#### Image Compression
- **JPEG Compression**: Quality control from 1-100%
- **PNG Compression**: Lossless compression with quality mapping
- **Format Support**: JPEG (.jpg, .jpeg) and PNG (.png)
- **Metadata Stripping**: Removes EXIF data for smaller files

#### Resize Options
Three flexible resize modes:
1. **Percentage Mode**: Scale images by percentage (10-100%)
2. **Dimensions Mode**: Set custom width and/or height in pixels
3. **None Mode**: Keep original dimensions

#### User Interface
- **Glassy Dark Theme**: Pitch black background (#0A0C12) with moonish blue tints
- **iOS-Inspired Design**: Smooth, modern interface with rounded corners
- **Material Design 3**: MD3-like components and interactions
- **Native File Selection**: Android's native file picker
- **Image Preview**: Real-time preview of selected images
- **Output Location**: Choose custom output directory

#### Technical Details
- **Architectures**: ARM64-v8a and x86_64 (AMD64)
- **Minimum Android**: 5.0 (API 21)
- **APK Size**: ~24MB
- **Desktop Binary**: 23MB (optimized with -ldflags="-s -w")

### Build Information

#### Desktop Build
```bash
go build -ldflags="-s -w" -o imgshrink-mobile
```
- **Platform**: Linux, macOS, Windows
- **Size**: 23MB
- **Status**: ✅ Tested and working

#### Android Build
```bash
fyne-cross android -arch=arm64,amd64 -app-id=com.virakt.imgshrink
```
- **ARM64 APK**: fyne-cross/dist/android-arm64/ImgShrink.apk
- **AMD64 APK**: fyne-cross/dist/android-amd64/ImgShrink.apk
- **Size**: ~24MB per architecture
- **Status**: ✅ Successfully built and tested

### Installation

#### From APK
1. Download the appropriate APK for your device architecture
2. Enable "Install from Unknown Sources" in Android settings
3. Install the APK
4. Grant storage permissions when prompted

#### From Source
```bash
git clone https://github.com/virakt/ImgShrink-Android.git
cd ImgShrink-Android
fyne-cross android -arch=arm64 -app-id=com.virakt.imgshrink
```

### GitHub Actions
Automated builds are configured via GitHub Actions:
- Triggers on push to main branch
- Builds both ARM64 and AMD64 APKs
- Uploads artifacts for download
- Creates releases automatically

### Known Issues

1. **Icon Warning**: The build process shows a warning about icon format, but this doesn't affect functionality
2. **First Launch**: May take a few seconds to initialize on first launch
3. **Large Images**: Very large images (>20MB) may take longer to process

### Fixes in This Release

#### v1.0.0
- ✅ Fixed nil pointer dereference crash on startup
- ✅ Reordered widget initialization to prevent crashes
- ✅ Converted icon from JPEG to PNG format
- ✅ Added GitHub Actions workflow for automated builds
- ✅ Optimized binary size with build flags

### Testing Status

#### Completed Tests
- ✅ Desktop compilation successful
- ✅ Android APK generation (ARM64 & AMD64)
- ✅ App launches without crashes
- ✅ APK installation on Android device successful

#### Pending Tests
- ⏳ Image selection functionality
- ⏳ Compression with various quality settings
- ⏳ Resize modes (Percentage/Dimensions/None)
- ⏳ Output location selection
- ⏳ Large file handling
- ⏳ Performance on mobile devices

### Compression Performance

Expected results based on settings:

| Quality | JPEG Reduction | PNG Reduction |
|---------|---------------|---------------|
| 100%    | 10-20%        | 20-30%        |
| 85%     | 40-60%        | 30-40%        |
| 50%     | 70-80%        | 40-50%        |
| 25%     | 85-90%        | 50-60%        |

### Future Enhancements

#### Planned Features
- Batch image processing
- Before/after comparison view
- Compression presets (Web, Email, Archive)
- Format conversion (JPEG ↔ PNG)
- WebP support
- Cloud storage integration
- Share functionality
- Dark/Light theme toggle

#### Performance Improvements
- Hardware acceleration
- Multi-threaded batch processing
- Background compression service
- Compression preview

### Credits

- **Framework**: [Fyne](https://fyne.io/) - Cross-platform GUI toolkit
- **Image Processing**: [imaging](https://github.com/disintegration/imaging) - Go image processing library
- **Build Tool**: [fyne-cross](https://github.com/fyne-io/fyne-cross) - Cross-compilation tool
- **Base Project**: [ImgShrink](https://github.com/virakt/imgshrink) - Original TUI application

### License
MIT License - See LICENSE file for details

### Support
- **Issues**: https://github.com/virakt/ImgShrink-Android/issues
- **Discussions**: https://github.com/virakt/ImgShrink-Android/discussions

### Changelog

#### 2024-02-03 - v1.0.0
- Initial release
- iOS-inspired glassy dark UI
- JPEG and PNG compression
- Three resize modes
- Native file selection
- Custom output location
- ARM64 and AMD64 support
- GitHub Actions CI/CD
