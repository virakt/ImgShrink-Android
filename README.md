# ImgShrink Mobile

Android app for compressing images with a minimal glassy UI inspired by iOS design.

## Features

- **Native File Selection** - Pick images from device storage
- **Image Preview** - View selected images before compression
- **Quality Control** - Adjust JPEG quality (1-100)
- **Resolution Control** - Set custom width/height or resize by percentage
- **Output Location** - Choose where to save compressed images
- **Real-time Stats** - See original vs compressed size and reduction percentage
- **Dark Glassy UI** - Pitch black with moonish blue tints

## Build

### Requirements
- Go 1.21+
- Fyne CLI v2.7.2+

### Install Fyne CLI
```bash
go install fyne.io/fyne/v2/cmd/fyne@latest
```

### Build APK (ARM64 only)
```bash
fyne package --target android/arm64 --appID com.virakt.imgshrink --icon Icon.png --release
```

### Build Desktop (for testing)
```bash
go build -o imgshrink-mobile
```

## Architecture

- **ARM64-v8a only** (24MB APK)
- Optimized for modern Android devices
