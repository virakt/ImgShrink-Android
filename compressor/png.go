package compressor

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
)

// PNGCompressor handles PNG compression
type PNGCompressor struct{}

// NewPNGCompressor creates a new PNG compressor
func NewPNGCompressor() *PNGCompressor {
	return &PNGCompressor{}
}

// Compress compresses a PNG image
func (c *PNGCompressor) Compress(inputPath string, options CompressionOptions) (*CompressionResult, error) {
	result := &CompressionResult{
		InputPath: inputPath,
		Success:   false,
	}

	// Get input file info
	inputInfo, err := os.Stat(inputPath)
	if err != nil {
		result.Error = fmt.Errorf("failed to stat input file: %w", err)
		return result, result.Error
	}
	result.InputSize = inputInfo.Size()

	// Open and decode image
	file, err := os.Open(inputPath)
	if err != nil {
		result.Error = fmt.Errorf("failed to open input file: %w", err)
		return result, result.Error
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		result.Error = fmt.Errorf("failed to decode PNG: %w", err)
		return result, result.Error
	}

	// Apply resizing if needed
	if options.ResizePercent > 0 && options.ResizePercent < 100 {
		bounds := img.Bounds()
		newWidth := int(float64(bounds.Dx()) * options.ResizePercent / 100)
		newHeight := int(float64(bounds.Dy()) * options.ResizePercent / 100)
		img = imaging.Resize(img, newWidth, newHeight, imaging.Lanczos)
	} else if options.ResizeWidth > 0 || options.ResizeHeight > 0 {
		if options.ResizeWidth > 0 && options.ResizeHeight > 0 {
			img = imaging.Resize(img, options.ResizeWidth, options.ResizeHeight, imaging.Lanczos)
		} else if options.ResizeWidth > 0 {
			img = imaging.Resize(img, options.ResizeWidth, 0, imaging.Lanczos)
		} else {
			img = imaging.Resize(img, 0, options.ResizeHeight, imaging.Lanczos)
		}
	}

	// Store dimensions
	bounds := img.Bounds()
	result.Width = bounds.Dx()
	result.Height = bounds.Dy()

	// Generate output path
	result.OutputPath = GenerateOutputPath(inputPath, options)

	// Create output directory if needed
	outputDir := filepath.Dir(result.OutputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		result.Error = fmt.Errorf("failed to create output directory: %w", err)
		return result, result.Error
	}

	// Create output file
	outFile, err := os.Create(result.OutputPath)
	if err != nil {
		result.Error = fmt.Errorf("failed to create output file: %w", err)
		return result, result.Error
	}
	defer outFile.Close()

	// Encode with compression settings
	encoder := &png.Encoder{
		CompressionLevel: png.CompressionLevel(options.CompressionLevel),
	}

	if err := encoder.Encode(outFile, img); err != nil {
		result.Error = fmt.Errorf("failed to encode PNG: %w", err)
		return result, result.Error
	}

	// Get output file size
	outputInfo, err := os.Stat(result.OutputPath)
	if err != nil {
		result.Error = fmt.Errorf("failed to stat output file: %w", err)
		return result, result.Error
	}
	result.OutputSize = outputInfo.Size()
	result.Reduction = CalculateReduction(result.InputSize, result.OutputSize)
	result.Success = true

	return result, nil
}

// EstimateSize estimates the compressed size
func (c *PNGCompressor) EstimateSize(inputPath string, options CompressionOptions) (int64, error) {
	info, err := GetImageInfo(inputPath)
	if err != nil {
		return 0, err
	}

	// Rough estimation
	estimatedSize := int64(float64(info.Size) * 0.8)

	// Adjust for resizing
	if options.ResizePercent > 0 && options.ResizePercent < 100 {
		factor := options.ResizePercent / 100
		estimatedSize = int64(float64(estimatedSize) * factor * factor)
	}

	return estimatedSize, nil
}
