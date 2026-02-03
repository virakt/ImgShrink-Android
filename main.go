package main

import (
	"fmt"
	"image/color"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
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
		return 4
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
	window         fyne.Window
	selectedFile   string
	outputURI      fyne.URI
	imagePreview   *canvas.Image
	fileLabel      *widget.Label
	sizeLabel      *widget.Label
	qualitySlider  *widget.Slider
	qualityLabel   *widget.Label
	resizeSlider   *widget.Slider
	resizeLabel    *widget.Label
	widthEntry     *widget.Entry
	heightEntry    *widget.Entry
	resizeMode     *widget.RadioGroup
	outputLabel    *widget.Label
	progressBar    *widget.ProgressBar
	statusLabel    *widget.Label
	compressBtn    *widget.Button
	jpegCompressor *compressor.JPEGCompressor
	pngCompressor  *compressor.PNGCompressor
}

func main() {
	a := app.New()
	a.Settings().SetTheme(&glassyTheme{})

	w := a.NewWindow("ImgShrink")
	w.Resize(fyne.NewSize(400, 650))

	imgApp := &ImgShrinkApp{
		window:         w,
		jpegCompressor: compressor.NewJPEGCompressor(),
		pngCompressor:  compressor.NewPNGCompressor(),
	}

	w.SetContent(imgApp.makeUI())
	w.ShowAndRun()
}

func (app *ImgShrinkApp) makeUI() fyne.CanvasObject {
	// Compact header
	title := widget.NewLabelWithStyle("ImgShrink", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// Compact image preview
	app.imagePreview = canvas.NewImageFromResource(theme.FileImageIcon())
	app.imagePreview.FillMode = canvas.ImageFillContain
	app.imagePreview.SetMinSize(fyne.NewSize(200, 150))

	previewCard := container.NewStack(
		canvas.NewRectangle(surfaceColor),
		container.NewPadded(app.imagePreview),
	)

	// Compact file info
	app.fileLabel = widget.NewLabel("No image selected")
	app.fileLabel.Wrapping = fyne.TextTruncate
	app.sizeLabel = widget.NewLabel("")

	// Select button
	selectBtn := widget.NewButton("Select Image", func() {
		app.selectImage()
	})
	selectBtn.Importance = widget.HighImportance

	// Compact quality slider
	app.qualityLabel = widget.NewLabel("Quality: 85%")
	app.qualitySlider = widget.NewSlider(1, 100)
	app.qualitySlider.Value = 85
	app.qualitySlider.Step = 5
	app.qualitySlider.OnChanged = func(value float64) {
		app.qualityLabel.SetText(fmt.Sprintf("Quality: %.0f%%", value))
	}

	// Compact resize controls
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

	// Output location
	app.outputLabel = widget.NewLabel("Output: Choose location")
	app.outputLabel.Wrapping = fyne.TextTruncate

	selectOutputBtn := widget.NewButton("Output Location", func() {
		app.selectOutputLocation()
	})

	// Compact controls card
	controlsCard := container.NewStack(
		canvas.NewRectangle(surfaceColor),
		container.NewPadded(
			container.NewVBox(
				app.qualityLabel,
				app.qualitySlider,
				container.NewGridWithColumns(3,
					widget.NewLabel("Resize:"),
					app.resizeMode,
				),
				app.resizeLabel,
				app.resizeSlider,
				container.NewGridWithColumns(2, app.widthEntry, app.heightEntry),
				selectOutputBtn,
				app.outputLabel,
			),
		),
	)

	// Compress button
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

	// Compact main layout - no scrolling needed
	content := container.NewVBox(
		title,
		previewCard,
		app.fileLabel,
		app.sizeLabel,
		selectBtn,
		controlsCard,
		app.compressBtn,
		app.progressBar,
		app.statusLabel,
	)

	return container.NewPadded(content)
}

func (app *ImgShrinkApp) selectImage() {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, app.window)
			return
		}
		if reader == nil {
			return
		}

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
				}
			}
		}

		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			os.Remove(tmpPath)
			dialog.ShowError(fmt.Errorf("unsupported file type. Please select a JPEG or PNG image"), app.window)
			return
		}

		// Validate image format
		format, err := compressor.GetImageFormat(tmpPath)
		if err != nil {
			os.Remove(tmpPath)
			dialog.ShowError(fmt.Errorf("unsupported image format"), app.window)
			return
		}

		// Get image info
		info, err := compressor.GetImageInfo(tmpPath)
		if err != nil {
			os.Remove(tmpPath)
			dialog.ShowError(err, app.window)
			return
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
		app.statusLabel.SetText("")

	}, app.window)
}

func (app *ImgShrinkApp) selectOutputLocation() {
	// Use file save dialog to let user choose where to save
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, app.window)
			return
		}
		if writer == nil {
			return
		}
		writer.Close()

		app.outputURI = writer.URI()
		app.outputLabel.SetText(fmt.Sprintf("Output: %s", writer.URI().Name()))
	}, app.window)
}

func (app *ImgShrinkApp) compressImage() {
	if app.selectedFile == "" {
		return
	}

	app.compressBtn.Disable()
	app.progressBar.Show()
	app.statusLabel.SetText("Compressing...")

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
			app.statusLabel.SetText(fmt.Sprintf("❌ Error: %v", err))
			app.progressBar.Hide()
			app.compressBtn.Enable()
			return
		}

		// If user selected output location, copy file there
		var finalPath string
		if app.outputURI != nil {
			// On Android, write to user-selected URI
			if runtime.GOOS == "android" {
				writer, err := storage.Writer(app.outputURI)
				if err == nil {
					compressedData, err := os.ReadFile(result.OutputPath)
					if err == nil {
						_, err = writer.Write(compressedData)
						writer.Close()
						if err == nil {
							finalPath = app.outputURI.String()
						}
					}
				}
			} else {
				// On desktop, use the URI path directly
				finalPath = app.outputURI.Path()
				os.Rename(result.OutputPath, finalPath)
			}
		} else {
			finalPath = result.OutputPath
		}

		if finalPath == "" {
			finalPath = result.OutputPath
		}

		successMsg := fmt.Sprintf("✓ Success!\n"+
			"%s → %s\n"+
			"Saved: %s (%.1f%%)\n"+
			"Location: %s",
			compressor.FormatBytes(result.InputSize),
			compressor.FormatBytes(result.OutputSize),
			compressor.FormatBytes(result.InputSize-result.OutputSize),
			result.Reduction,
			filepath.Base(finalPath),
		)

		app.statusLabel.SetText(successMsg)
		app.progressBar.Hide()
		app.compressBtn.Enable()

		// Show success dialog
		dialog.ShowInformation("Success", successMsg, app.window)
	}()
}
