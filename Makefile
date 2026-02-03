# ImgShrink Mobile - Makefile for Android builds
# Optimized for armv8/x64 architectures only

.PHONY: all deps build-android clean install-fyne-cross

# Install fyne-cross for cross-compilation
install-fyne-cross:
	@echo "Installing fyne-cross..."
	go install fyne.io/fyne-cross@latest

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Build for Android (arm64-v8a and x86_64 only)
build-android: deps
	@echo "Building for Android (arm64-v8a, x86_64)..."
	fyne-cross android -arch=arm64,amd64 -app-id=com.virakt.imgshrink -release

# Build for Android with optimizations (faster builds)
build-android-fast: deps
	@echo "Building for Android with optimizations..."
	fyne-cross android -arch=arm64,amd64 -app-id=com.virakt.imgshrink -release \
		-ldflags="-s -w" \
		-tags=release

# Build debug version
build-android-debug: deps
	@echo "Building debug version for Android..."
	fyne-cross android -arch=arm64,amd64 -app-id=com.virakt.imgshrink

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf fyne-cross
	rm -f *.apk
	rm -f *.aab

# Run on desktop for testing
run:
	@echo "Running on desktop..."
	go run .

# Full build pipeline
all: install-fyne-cross build-android-fast

# Help
help:
	@echo "ImgShrink Mobile - Build Commands"
	@echo ""
	@echo "make install-fyne-cross  - Install fyne-cross tool"
	@echo "make deps                - Install Go dependencies"
	@echo "make build-android       - Build Android APK (arm64, x64)"
	@echo "make build-android-fast  - Build with optimizations (faster)"
	@echo "make build-android-debug - Build debug version"
	@echo "make run                 - Run on desktop for testing"
	@echo "make clean               - Clean build artifacts"
	@echo "make all                 - Full build pipeline"
