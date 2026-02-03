#!/bin/bash

# ImgShrink Mobile - Build Script
# Optimized build for Android (arm64-v8a and x86_64 only)

set -e

echo "🚀 ImgShrink Mobile - Android Build Script"
echo "=========================================="

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if fyne-cross is installed
if ! command -v fyne-cross &> /dev/null; then
    echo -e "${BLUE}Installing fyne-cross...${NC}"
    go install fyne.io/fyne-cross@latest
    echo -e "${GREEN}✓ fyne-cross installed${NC}"
fi

# Install dependencies
echo -e "${BLUE}Installing dependencies...${NC}"
go mod download
go mod tidy
echo -e "${GREEN}✓ Dependencies installed${NC}"

# Build for Android
echo -e "${BLUE}Building for Android (arm64-v8a, x86_64)...${NC}"
echo "This may take a few minutes on first build..."

fyne-cross android \
    -arch=arm64,amd64 \
    -app-id=com.virakt.imgshrink \
    -release \
    -ldflags="-s -w" \
    -tags=release

echo ""
echo -e "${GREEN}✓ Build complete!${NC}"
echo ""
echo "APK files generated:"
echo "  - fyne-cross/dist/android-arm64/ImgShrink.apk"
echo "  - fyne-cross/dist/android-amd64/ImgShrink.apk"
echo ""
echo "Install on your device:"
echo "  adb install fyne-cross/dist/android-arm64/ImgShrink.apk"
