package main

import (
	"fmt"
	"image/color"
	"io"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/virakt/imgshrink-mobile/compressor"
)

// Custom dark theme with iOS-inspired glassy design
type glassyTheme struct{}

var (
	// Pitch black with moonish blue tints
	bgColor        = color.NRGBA{R: 10, G: 12, B: 18, A: 255}    // Almost black with blue tint
	surfaceColor   = color.NRGBA{R: 20, G: 25, B: 35, A: 230}    // Dark surface with transparency
	primaryColor   = color.NRGBA{R: 100, G: 120, B: 180, A: 255} // Moonish blue
	accentColor    = color.NRGBA{R: 140, G: 160, B: 220, A: 255} // Lighter blue
	textColor      = color.NRGBA{R: 240, G: 245, B: 255, A: 255} // Off-white
	mutedTextColor = color.NRGBA{R: 160, G: 170, B: 190, A: 255} // Muted blue-gray
	successColor   = color.NRGBA{R: 100, G: 200, B: 150, A: 255} // Soft green
	errorColor     = color.NRGBA{R: 220, G: 100, B: 120, A: 255} // Soft red
)

func (g glassyTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return bgColor
	case theme.ColorNameButton:
		return surfaceColor
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 30, G: 35, B: 45, A: 150}
	case theme.ColorNamePrimary:
		return primaryColor
	case theme.ColorNameHover:
		return accentColor
	case theme.ColorNameFocus:
		return accentColor
	case theme.ColorNameForeground:
		return textColor
	case theme.ColorNameDisabled:
		return mutedTextColor
	case theme.ColorNameError:
		return errorColor
	case theme.ColorNameSuccess:
		return successColor
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (g glassyTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (g glassyTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (g glassyTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 6
	case theme.SizeNameInlineIcon:
		return 20
	case theme.SizeNameScrollBar:
		return 0 // Hide scrollbar
	case theme.SizeNameScrollBarSmall:
		return 0 // Hide scrollbar
	default:
		return theme.DefaultTheme().Size(name)
	}
}

// App state
type ImgShrinkApp struct {
	window              fyne.Window
	selectedFile        string
	outputDirURI        fyne.URI
	outputFilename      string
	outputExt           string
	outputFilenameEntry *widget.Entry
	outputExtGroup      *widget.RadioGroup
	imagePreview        *canvas.Image
	fileLabel           *widget.Label
	sizeLabel           *widget.Label
	qualitySlider       *widget.Slider
	qualityLabel        *widget.Label
	resizeSlider        *widget.Slider
	resizeLabel         *widget.Label
	widthEntry          *widget.Entry
	heightEntry         *widget.Entry
	resizeMode          *widget.RadioGroup
	outputLabel         *widget.Label
	progressBar         *widget.ProgressBar
	statusLabel         *widget.Label
	compressBtn         *widget.Button
	jpegCompressor      *compressor.JPEGCompressor
	pngCompressor       *compressor.PNGCompressor
}

func main() {
	a := app.New()
	a.Settings().SetTheme(&glassyTheme{})

	w := a.NewWindow("ImgShrink")

	// Set initial size, will be adjusted by makeUI based on screen
	w.Resize(fyne.NewSize(400, 800))

	imgApp := &ImgShrinkApp{
		window:         w,
		jpegCompressor: compressor.NewJPEGCompressor(),
		pngCompressor:  compressor.NewPNGCompressor(),
	}

	w.SetContent(imgApp.makeUI())
	w.ShowAndRun()
}

func (app *ImgShrinkApp) makeUI() fyne.CanvasObject {
	// Header - compact
	title := widget.NewLabelWithStyle("ImgShrink", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// Image preview - takes 50% of screen space
	app.imagePreview = canvas.NewImageFromResource(theme.FileImageIcon())
	app.imagePreview.FillMode = canvas.ImageFillContain
	app.imagePreview.SetMinSize(fyne.NewSize(360, 360))

	previewCard := container.NewStack(
		canvas.NewRectangle(surfaceColor),
		container.NewPadded(app.imagePreview),
	)

	// File info - single line compact
	app.fileLabel = widget.NewLabel("No image selected")
	app.fileLabel.Wrapping = fyne.TextTruncate
	app.sizeLabel = widget.NewLabel("")

	// Select button - directly opens file picker
	selectBtn := widget.NewButton("Select Image", func() {
		app.selectImage()
	})
	selectBtn.Importance = widget.HighImportance

	// Quality slider - compact
	app.qualityLabel = widget.NewLabel("Quality: 85%")
	app.qualitySlider = widget.NewSlider(1, 100)
	app.qualitySlider.Value = 85
	app.qualitySlider.Step = 5
	app.qualitySlider.OnChanged = func(value float64) {
		app.qualityLabel.SetText(fmt.Sprintf("Quality: %.0f%%", value))
	}

	// Resize controls - compact
	app.resizeLabel = widget.NewLabel("Resize: 100%")
	app.resizeSlider = widget.NewSlider(10, 100)
	app.resizeSlider.Value = 100
	app.resizeSlider.Step = 10
	app.resizeSlider.OnChanged = func(value float64) {
		app.resizeLabel.SetText(fmt.Sprintf("Resize: %.0f%%", value))
	}
	app.resizeSlider.Hide()

	app.widthEntry = widget.NewEntry()
	app.widthEntry.SetPlaceHolder("W")
	app.widthEntry.Hide()

	app.heightEntry = widget.NewEntry()
	app.heightEntry.SetPlaceHolder("H")
	app.heightEntry.Hide()

	app.resizeMode = widget.NewRadioGroup([]string{"%", "WxH", "None"}, func(value string) {
		switch value {
		case "%":
			app.resizeSlider.Show()
			app.widthEntry.Hide()
			app.heightEntry.Hide()
		case "WxH":
			app.resizeSlider.Hide()
			app.widthEntry.Show()
			app.heightEntry.Show()
		case "None":
			app.resizeSlider.Hide()
			app.widthEntry.Hide()
			app.heightEntry.Hide()
		}
	})
	app.resizeMode.Horizontal = true
	app.resizeMode.SetSelected("None")

	// Output controls - minimal design
	app.outputFilenameEntry = widget.NewEntry()
	app.outputFilenameEntry.SetPlaceHolder("filename")

	app.outputExtGroup = widget.NewRadioGroup([]string{".jpg", ".jpeg", ".png"}, func(value string) {
		app.outputExt = value
	})
	app.outputExtGroup.Horizontal = true
	app.outputExtGroup.SetSelected(".jpg")
	app.outputExt = ".jpg"

	// Compact controls card with proper margins
	controlsCard := container.NewStack(
		canvas.NewRectangle(surfaceColor),
		container.NewPadded(
			container.NewVBox(
				widget.NewLabel("Quality"),
				app.qualitySlider,
				app.qualityLabel,

				widget.NewSeparator(),

				widget.NewLabel("Resize"),
				app.resizeMode,
				app.resizeSlider,
				app.resizeLabel,
				container.NewGridWithColumns(2, app.widthEntry, app.heightEntry),

				widget.NewSeparator(),

				widget.NewLabel("Output Format"),
				app.outputExtGroup,
			),
		),
	)

	// Compress button - full width and much larger
	app.compressBtn = widget.NewButton("Compress Image", func() {
		app.compressImage()
	})
	app.compressBtn.Importance = widget.HighImportance
	app.compressBtn.Disable()

	// Progress
	app.progressBar = widget.NewProgressBar()
	app.progressBar.Hide()

	app.statusLabel = widget.NewLabel("")
	app.statusLabel.Wrapping = fyne.TextWrapWord
	app.statusLabel.Alignment = fyne.TextAlignCenter
	app.statusLabel.Hide() // Hide by default, only show in dialogs

	// Create a large button with spacer to make it bigger
	btnSpacer := canvas.NewRectangle(color.Transparent)
	btnSpacer.SetMinSize(fyne.NewSize(0, 80))
	largeBtn := container.NewStack(
		btnSpacer,
		container.NewPadded(
			container.NewPadded(app.compressBtn),
		),
	)

	// Main layout with proper spacing - no scroll needed
	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		previewCard,
		widget.NewSeparator(),
		container.NewGridWithColumns(2, app.fileLabel, app.sizeLabel),
		selectBtn,
		widget.NewSeparator(),
		controlsCard,
		widget.NewSeparator(),
		largeBtn,
		app.progressBar,
	)

	return content
}

func (app *ImgShrinkApp) selectImage() {
	// Directly open file picker
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, app.window)
			return
		}
		if reader == nil {
			return
		}
		app.handleImageFile(reader)
	}, app.window)
}

func (app *ImgShrinkApp) handleImageFile(reader fyne.URIReadCloser) {
	uri := reader.URI()

	// Get file extension from URI
	uriStr := uri.String()
	ext := filepath.Ext(uriStr)

	// Create temp file path first
	tmpDir := os.TempDir()
	tmpPath := filepath.Join(tmpDir, "imgshrink_temp"+ext)

	// Write data to temp file directly from reader
	tmpF, err := os.Create(tmpPath)
	if err != nil {
		reader.Close()
		dialog.ShowError(fmt.Errorf("failed to create temp file: %w", err), app.window)
		return
	}

	_, err = io.Copy(tmpF, reader)
	tmpF.Close()
	reader.Close()

	if err != nil {
		os.Remove(tmpPath)
		dialog.ShowError(fmt.Errorf("failed to write temp file: %w", err), app.window)
		return
	}

	// Detect extension from content if not available
	if ext == "" {
		fileData, err := os.ReadFile(tmpPath)
		if err == nil && len(fileData) > 7 {
			// Check for JPEG magic bytes
			if fileData[0] == 0xFF && fileData[1] == 0xD8 {
				ext = ".jpg"
				newPath := filepath.Join(tmpDir, "imgshrink_temp.jpg")
				os.Rename(tmpPath, newPath)
				tmpPath = newPath
			} else if fileData[0] == 0x89 && fileData[1] == 0x50 &&
				fileData[2] == 0x4E && fileData[3] == 0x47 {
				ext = ".png"
				newPath := filepath.Join(tmpDir, "imgshrink_temp.png")
				os.Rename(tmpPath, newPath)
				tmpPath = newPath
			} else if fileData[0] == 0x42 && fileData[1] == 0x4D {
				ext = ".bmp"
				newPath := filepath.Join(tmpDir, "imgshrink_temp.bmp")
				os.Rename(tmpPath, newPath)
				tmpPath = newPath
			} else if len(fileData) > 12 && fileData[0] == 0x00 && fileData[1] == 0x00 {
				// WEBP or other formats
				ext = ".webp"
				newPath := filepath.Join(tmpDir, "imgshrink_temp.webp")
				os.Rename(tmpPath, newPath)
				tmpPath = newPath
			}
		}
	}

	// Get image info first
	info, err := compressor.GetImageInfo(tmpPath)
	if err != nil {
		os.Remove(tmpPath)
		dialog.ShowError(fmt.Errorf("unsupported image format"), app.window)
		return
	}

	// Try to get format from compressor, fallback to extension
	format, err := compressor.GetImageFormat(tmpPath)
	if err != nil {
		// Image is valid but compressor doesn't support it, use generic format
		format = compressor.ImageFormat(strings.TrimPrefix(ext, "."))
	}

	// Update UI
	app.selectedFile = tmpPath
	app.fileLabel.SetText(uri.Name())
	app.sizeLabel.SetText(fmt.Sprintf("%s • %dx%d • %s",
		compressor.FormatBytes(info.Size),
		info.Width,
		info.Height,
		format,
	))

	// Load preview from temp file
	app.imagePreview.File = tmpPath
	app.imagePreview.Resource = nil
	app.imagePreview.Refresh()

	app.compressBtn.Enable()
}

func (app *ImgShrinkApp) compressImage() {
	if app.selectedFile == "" {
		return
	}

	app.compressBtn.Disable()
	app.progressBar.Show()

	go func() {
		// Get compression options
		options := compressor.DefaultOptions()
		options.Quality = int(app.qualitySlider.Value)

		// Use temp directory for compression
		options.OutputDir = os.TempDir()

		// Handle resize options based on mode
		switch app.resizeMode.Selected {
		case "%":
			if app.resizeSlider.Value < 100 {
				options.ResizePercent = app.resizeSlider.Value
			}
		case "WxH":
			// Parse width and height
			if app.widthEntry.Text != "" {
				fmt.Sscanf(app.widthEntry.Text, "%d", &options.ResizeWidth)
			}
			if app.heightEntry.Text != "" {
				fmt.Sscanf(app.heightEntry.Text, "%d", &options.ResizeHeight)
			}
		case "None":
			// No resizing
			options.ResizePercent = 0
			options.ResizeWidth = 0
			options.ResizeHeight = 0
		}

		// Detect format and compress
		format, _ := compressor.GetImageFormat(app.selectedFile)

		var result *compressor.CompressionResult
		var err error

		if format == compressor.FormatJPEG {
			result, err = app.jpegCompressor.Compress(app.selectedFile, options)
		} else {
			// For PNG, use compression level instead of quality
			options.CompressionLevel = int(app.qualitySlider.Value / 11) // Map 0-100 to 0-9
			result, err = app.pngCompressor.Compress(app.selectedFile, options)
		}

		// Update UI on main thread
		if err != nil || !result.Success {
			app.progressBar.Hide()
			app.compressBtn.Enable()
			dialog.ShowError(fmt.Errorf("compression failed: %v", err), app.window)
			return
		}

		// Get output extension
		ext := app.outputExt
		if ext == "" {
			ext = ".jpg"
		}

		// Determine output filename - use original name if not specified
		outputFilename := "compressed" + ext
		if app.selectedFile != "" {
			baseName := filepath.Base(app.selectedFile)
			baseName = strings.TrimSuffix(baseName, filepath.Ext(baseName))
			outputFilename = baseName + "_compressed" + ext
		}

		// Read compressed data
		compressedData, err := os.ReadFile(result.OutputPath)
		if err != nil {
			app.progressBar.Hide()
			app.compressBtn.Enable()
			dialog.ShowError(fmt.Errorf("failed to read compressed file: %v", err), app.window)
			return
		}

		// Always use file save dialog to let user choose location
		saveDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				app.progressBar.Hide()
				app.compressBtn.Enable()
				dialog.ShowError(err, app.window)
				return
			}
			if writer == nil {
				app.progressBar.Hide()
				app.compressBtn.Enable()
				return
			}

			_, err = writer.Write(compressedData)
			writer.Close()

			if err != nil {
				app.progressBar.Hide()
				app.compressBtn.Enable()
				dialog.ShowError(fmt.Errorf("failed to write file: %v", err), app.window)
				return
			}

			// Show success
			successMsg := fmt.Sprintf("✓ Compression Successful!\n\n"+
				"Original: %s\n"+
				"Compressed: %s\n"+
				"Saved: %s (%.1f%% reduction)",
				compressor.FormatBytes(result.InputSize),
				compressor.FormatBytes(result.OutputSize),
				compressor.FormatBytes(result.InputSize-result.OutputSize),
				result.Reduction,
			)

			app.progressBar.Hide()
			app.compressBtn.Enable()
			dialog.ShowInformation("Success", successMsg, app.window)
		}, app.window)

		saveDialog.SetFileName(outputFilename)
		saveDialog.Show()
	}()
}
