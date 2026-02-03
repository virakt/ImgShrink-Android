# Testing Checklist

## Desktop Testing (Completed ✓)

- [x] App compiles successfully
- [x] Binary size reasonable (~30MB)
- [x] Dependencies resolved correctly
- [ ] UI renders correctly
- [ ] Image selection dialog works
- [ ] Image preview displays
- [ ] Quality slider functions
- [ ] Resize mode selection works
- [ ] Width/Height inputs work
- [ ] Output location selection works
- [ ] Compression executes successfully
- [ ] Results display correctly

## Android Build Testing (In Progress)

- [ ] ARM64 APK builds successfully
- [ ] x86_64 APK builds successfully
- [ ] APK sizes are reasonable (~15-20MB)
- [ ] Build optimizations applied (-s -w flags)

## Android Runtime Testing (Requires Device)

### Installation
- [ ] APK installs on ARM64 device
- [ ] APK installs on x86_64 emulator
- [ ] Storage permissions requested
- [ ] App icon displays correctly
- [ ] App appears in launcher

### UI/UX
- [ ] Glassy dark theme renders correctly
- [ ] Moonish blue colors display properly
- [ ] Text is readable
- [ ] Touch targets are appropriate size
- [ ] Scrolling works smoothly
- [ ] Buttons respond to touch
- [ ] Sliders work with touch input

### Functionality
- [ ] Native file picker opens
- [ ] Can select JPEG images
- [ ] Can select PNG images
- [ ] Image preview displays correctly
- [ ] Quality slider adjusts (1-100%)
- [ ] Resize mode switches work
  - [ ] Percentage mode
  - [ ] Dimensions mode (width/height)
  - [ ] None mode
- [ ] Width/height inputs accept numbers
- [ ] Output location picker works
- [ ] Compress button triggers compression
- [ ] Progress indicator shows during compression
- [ ] Success message displays
- [ ] Compressed file is saved correctly

### Image Compression Tests

#### JPEG Tests
- [ ] Small JPEG (<1MB)
- [ ] Medium JPEG (1-5MB)
- [ ] Large JPEG (>5MB)
- [ ] Quality 100% (minimal compression)
- [ ] Quality 50% (medium compression)
- [ ] Quality 10% (high compression)
- [ ] Resize by percentage (50%)
- [ ] Resize by dimensions (800x600)
- [ ] Custom output location

#### PNG Tests
- [ ] Small PNG (<1MB)
- [ ] Medium PNG (1-5MB)
- [ ] Large PNG (>5MB)
- [ ] PNG with transparency
- [ ] Quality levels (mapped to compression 0-9)
- [ ] Resize by percentage
- [ ] Resize by dimensions
- [ ] Custom output location

### Edge Cases
- [ ] Very large image (>20MB)
- [ ] Very small image (<100KB)
- [ ] Square image (1:1 aspect ratio)
- [ ] Wide image (16:9 aspect ratio)
- [ ] Tall image (9:16 aspect ratio)
- [ ] Invalid file selection
- [ ] Insufficient storage space
- [ ] App backgrounding during compression
- [ ] Multiple compressions in sequence

### Performance
- [ ] Compression completes in reasonable time
- [ ] App remains responsive during compression
- [ ] Memory usage is acceptable
- [ ] No crashes or freezes
- [ ] Battery usage is reasonable

### Results Verification
- [ ] Output file size is smaller than input
- [ ] Reduction percentage is accurate
- [ ] Output file path is correct
- [ ] Output file is valid image
- [ ] Image quality is acceptable
- [ ] Dimensions match expected values

## Regression Testing

After any code changes:
- [ ] Desktop build still works
- [ ] Android build still works
- [ ] All core features still function
- [ ] No new crashes introduced
- [ ] Performance hasn't degraded

## Test Results Summary

### Desktop
- Build: ✓ Success
- Binary Size: 30MB
- Status: Compiled successfully, UI testing pending

### Android
- Build: In Progress
- Target Architectures: ARM64, x86_64
- Optimizations: -ldflags="-s -w"
- Status: Building...

### Notes
- First build takes longer due to Docker image downloads
- Subsequent builds will be faster
- Icon placeholder created (custom icon recommended)
- Go version adjusted to 1.23 for compatibility
