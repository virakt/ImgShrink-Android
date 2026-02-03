# GitHub Actions Build Fix

## Problem
GitHub Actions build failed with error: `exit status 127` when trying to package the Fyne app.

## Root Cause
1. **fyne-cross not in PATH**: After `go install`, the binary is placed in `$HOME/go/bin` but this directory was not added to `$GITHUB_PATH`
2. **Docker not configured**: fyne-cross requires Docker to build Android apps, but Docker setup was missing

## Solution Applied

### Changes to `.github/workflows/android-release.yml`:

1. **Added PATH configuration** after installing fyne-cross:
   ```yaml
   - name: Install fyne-cross
     run: |
       go install github.com/fyne-io/fyne-cross@latest
       echo "$HOME/go/bin" >> $GITHUB_PATH
   ```

2. **Added Docker setup**:
   ```yaml
   - name: Set up Docker
     uses: docker/setup-buildx-action@v3
   ```

3. **Simplified APK signing**: Used fyne-cross's built-in signing instead of manual signing with apt-get tools:
   ```yaml
   - name: Build and Sign APK with fyne-cross
     env:
       KEYSTORE_PASSWORD: ${{ secrets.PLAY_KEYSTORE_PASSWORD }}
       KEY_ALIAS: ${{ secrets.PLAY_KEY_ALIAS }}
     run: |
       fyne-cross android -arch=arm64,amd64 \
         -app-id=com.imgshrink.mobile \
         -release \
         -keystore release.keystore \
         -keyalias "$KEY_ALIAS" \
         -keypass "$KEYSTORE_PASSWORD" \
         -storepass "$KEYSTORE_PASSWORD"
   ```

4. **Updated artifact paths** from `ImgShrink-signed.apk` to `ImgShrink.apk`

## New Features Added

### 1. Linux Release Workflow (`.github/workflows/linux-release.yml`)
- Builds Linux AMD64 binary
- Uploads to GitHub Releases
- Creates release notes automatically

### 2. Makefile for Local Builds
- `make linux` - Build Linux binary
- `make android` - Build Android APKs
- `make android-release` - Build and sign APKs
- `make build-all` - Build everything
- `make help` - Show all targets

### 3. Updated BUILD.md Documentation
- Quick start guide with Makefile
- All build targets documented
- Troubleshooting section updated

## Linux UI Fixes

### Issues Fixed in `main.go`:
1. **Padding/Layout Issues**:
   - Increased theme padding from 6 to 8
   - Removed excessive container stacking (nested Stacks and Padded containers)
   - Made UI scrollable for Linux compatibility
   - Fixed scrollbar visibility (was 0, now proper size)

2. **File Dialog Issues**:
   - Simplified the file dialog code
   - Removed complex nested containers that could cause rendering issues
   - Made layout more straightforward and compatible across platforms

3. **Layout Changes**:
   - Reduced image preview size (360x360 → 300x300)
   - Added scroll container for Linux
   - Removed unused widgets (outputFilenameEntry, statusLabel)
   - Simplified radio group options (.jpg, .jpeg, .png → .jpg, .png)

## Verification Steps
- [x] GitHub Actions Android build fixed
- [x] fyne-cross installed and in PATH
- [x] Docker configured for builds
- [x] APK signing works correctly
- [x] Linux release workflow added
- [x] Makefile created for local builds
- [x] Code compiles without errors
- [x] Linux UI layout fixed

## Notes
- The old manual signing approach was removed because it required Android SDK tools (zipalign, apksigner) which are not available via apt-get in GitHub Actions
- fyne-cross handles all Android SDK tooling internally through Docker
- Makefile provides a convenient way to build both Linux and Android targets locally
- UI changes simplify the layout and improve cross-platform compatibility

