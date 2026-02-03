
# ImgShrink Mobile - Makefile
# For building Linux binary and Android APKs

# Variables
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -s -w
BINARY_NAME := imgshrink-mobile
APK_NAME := ImgShrink

# Default target
.PHONY: all
all: linux

# Build Linux binary
.PHONY: linux
linux:
	@echo "Building Linux binary..."
	go build -ldflags="$(LDFLAGS)" -o $(BINARY_NAME)-$(VERSION)-linux-amd64
	@echo "✓ Built: $(BINARY_NAME)-$(VERSION)-linux-amd64"
	@chmod +x $(BINARY_NAME)-$(VERSION)-linux-amd64

# Build Linux binary (simple)
.PHONY: linux-simple
linux-simple:
	go build -ldflags="$(LDFLAGS)" -o $(BINARY_NAME)

# Build Android ARM64 APK
.PHONY: android-arm64
android-arm64: check-fyne-cross
	@echo "Building Android ARM64 APK..."
	fyne-cross android -arch=arm64 -app-id=com.imgshrink.mobile -icon Icon.png
	@echo "✓ Built: fyne-cross/dist/android-arm64/$(APK_NAME).apk"

# Build Android x86_64 APK
.PHONY: android-amd64
android-amd64: check-fyne-cross
	@echo "Building Android x86_64 APK..."
	fyne-cross android -arch=amd64 -app-id=com.imgshrink.mobile -icon Icon.png
	@echo "✓ Built: fyne-cross/dist/android-amd64/$(APK_NAME).apk"

# Build Android APKs for both architectures
.PHONY: android
android: android-arm64 android-amd64
	@echo ""
	@echo "✓ All Android APKs built successfully!"

# Build Android APKs and sign them
.PHONY: android-release
android-release: check-fyne-cross
	@echo "Building and signing Android APKs..."
	@read -p "Keystore path: " KEYSTORE && \
	read -p "Key alias: " KEYALIAS && \
	read -s -p "Keystore password: " KEYSTORE_PASS && echo && \
	read -s -p "Key password: " KEY_PASS && echo && \
	fyne-cross android -arch=arm64,amd64 \
		-app-id=com.imgshrink.mobile \
		-release \
		-keystore $$KEYSTORE \
		-keyalias $$KEYALIAS \
		-storepass $$KEYSTORE_PASS \
		-keypass $$KEY_PASS \
		-icon Icon.png
	@echo ""
	@echo "✓ Signed APKs built:"
	@echo "  - fyne-cross/dist/android-arm64/$(APK_NAME).apk"
	@echo "  - fyne-cross/dist/android-amd64/$(APK_NAME).apk"

# Check if fyne-cross is installed
.PHONY: check-fyne-cross
check-fyne-cross:
	@which fyne-cross > /dev/null 2>&1 || (echo "fyne-cross not found. Installing..." && go install github.com/fyne-io/fyne-cross@latest)

# Install dependencies
.PHONY: deps
deps:
	go mod download
	@echo "Dependencies installed."

# Clean build artifacts
.PHONY: clean
clean:
	rm -f $(BINARY_NAME)*
	rm -rf fyne-cross/dist/*
	rm -rf fyne-cross/tmp/*
	@echo "✓ Cleaned build artifacts"

# Full clean including cache
.PHONY: distclean
distclean: clean
	go clean -cache
	@echo "✓ Full clean completed"

# Run the application (Linux only)
.PHONY: run
run: linux-simple
	./$(BINARY_NAME)

# Build everything
.PHONY: build-all
build-all: linux android

# Show help
.PHONY: help
help:
	@echo "ImgShrink Mobile - Build Targets"
	@echo ""
	@echo "Linux targets:"
	@echo "  linux          - Build Linux binary with version"
	@echo "  linux-simple   - Build Linux binary (simple)"
	@echo "  run            - Build and run Linux binary"
	@echo ""
	@echo "Android targets:"
	@echo "  android        - Build APKs for all architectures"
	@echo "  android-arm64  - Build ARM64 APK"
	@echo "  android-amd64  - Build x86_64 APK"
	@echo "  android-release - Build and sign APKs (requires keystore)"
	@echo ""
	@echo "Utility targets:"
	@echo "  deps           - Install dependencies"
	@echo "  clean          - Remove build artifacts"
	@echo "  distclean      - Full clean including cache"
	@echo "  help           - Show this help"
	@echo ""
	@echo "Version: $(VERSION)"

