package main

import (
	"fmt"
	"image/color"
	"path/filepath"

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
		return 8
	case theme.SizeNameInlineIcon:
		return 24
	case theme.SizeNameScrollBar:
		return 12
	case theme.SizeNameScrollBarSmall:
		return 6
	default:
		return theme.DefaultTheme().Size(name)
	}
}

// App state
type ImgShrinkApp struct {
	window         fyne.Window
	selectedFile   string
	outputPath     string
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
	w.Resize(fyne.NewSize(400, 700))

	imgApp := &ImgShrinkApp{
		window:         w,
		jpegCompressor: compressor.NewJPEGCompressor(),
		pngCompressor:  compressor.NewPNGCompressor(),
	}

	w.SetContent(imgApp.makeUI())
	w.ShowAndRun()
}

func (app *ImgShrinkApp) makeUI() fyne.CanvasObject {
	// Header
	title := widget.NewLabelWithStyle("ImgShrink", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	subtitle := widget.NewLabelWithStyle("Compress your images", fyne.TextAlignCenter, fyne.TextStyle{})
	subtitle.TextStyle.Italic = true

	header := container.NewVBox(
		container.NewPadded(title),
		subtitle,
	)

	// Image preview
	app.imagePreview = canvas.NewImageFromResource(theme.FileImageIcon())
	app.imagePreview.FillMode = canvas.ImageFillContain
	app.imagePreview.SetMinSize(fyne.NewSize(300, 300))

	previewCard := container.NewStack(
		canvas.NewRectangle(surfaceColor),
		container.NewPadded(app.imagePreview),
	)

	// File info
	app.fileLabel = widget.NewLabel("No image selected")
	app.fileLabel.Wrapping = fyne.TextWrapWord
	app.sizeLabel = widget.NewLabel("")

	fileInfo := container.NewVBox(
		app.fileLabel,
		app.sizeLabel,
	)

	// Select button
	selectBtn := widget.NewButton("Select Image", func() {
		app.selectImage()
	})
	selectBtn.Importance = widget.HighImportance

	// Quality slider
	app.qualityLabel = widget.NewLabel("Quality: 85%")
	app.qualitySlider = widget.NewSlider(1, 100)
	app.qualitySlider.Value = 85
	app.qualitySlider.Step = 5
	app.qualitySlider.OnChanged = func(value float64) {
		app.qualityLabel.SetText(fmt.Sprintf("Quality: %.0f%%", value))
	}

	qualityControl := container.NewVBox(
		app.qualityLabel,
		app.qualitySlider,
	)

	// Resize slider (for percentage mode) - CREATE FIRST
	app.resizeLabel = widget.NewLabel("Resize: 100%")
	app.resizeSlider = widget.NewSlider(10, 100)
	app.resizeSlider.Value = 100
	app.resizeSlider.Step = 5
	app.resizeSlider.OnChanged = func(value float64) {
		app.resizeLabel.SetText(fmt.Sprintf("Resize: %.0f%%", value))
	}
	app.resizeSlider.Hide()

	// Dimension entries (for custom dimensions mode) - CREATE SECOND
	app.widthEntry = widget.NewEntry()
	app.widthEntry.SetPlaceHolder("Width (px)")
	app.widthEntry.Hide()

	app.heightEntry = widget.NewEntry()
	app.heightEntry.SetPlaceHolder("Height (px)")
	app.heightEntry.Hide()

	// Resize mode selection - CREATE LAST (after widgets it references)
	app.resizeMode = widget.NewRadioGroup([]string{"Percentage", "Dimensions", "None"}, func(value string) {
		switch value {
		case "Percentage":
			app.resizeSlider.Show()
			app.widthEntry.Hide()
			app.heightEntry.Hide()
		case "Dimensions":
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

	resizeControl := container.NewVBox(
		widget.NewLabel("Resize Options:"),
		app.resizeMode,
		app.resizeLabel,
		app.resizeSlider,
		container.NewGridWithColumns(2,
			app.widthEntry,
			app.heightEntry,
		),
	)

	// Output location
	app.outputLabel = widget.NewLabel("Output: Same as input")
	app.outputLabel.Wrapping = fyne.TextWrapWord

	selectOutputBtn := widget.NewButton("Choose Output Location", func() {
		app.selectOutputLocation()
	})

	outputControl := container.NewVBox(
		widget.NewLabel("Output Location:"),
		app.outputLabel,
		selectOutputBtn,
	)

	// Controls card
	controlsCard := container.NewStack(
		canvas.NewRectangle(surfaceColor),
		container.NewPadded(
			container.NewVBox(
				qualityControl,
				widget.NewSeparator(),
				resizeControl,
				widget.NewSeparator(),
				outputControl,
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

	// Main layout
	content := container.NewVBox(
		header,
		widget.NewSeparator(),
		previewCard,
		fileInfo,
		selectBtn,
		widget.NewSeparator(),
		controlsCard,
		app.compressBtn,
		app.progressBar,
		app.statusLabel,
	)

	return container.NewPadded(
		container.NewVScroll(content),
	)
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
		defer reader.Close()

		uri := reader.URI()
		path := uri.Path()

		// Validate image format
		format, err := compressor.GetImageFormat(path)
		if err != nil {
			dialog.ShowError(fmt.Errorf("unsupported image format"), app.window)
			return
		}

		// Get image info
		info, err := compressor.GetImageInfo(path)
		if err != nil {
			dialog.ShowError(err, app.window)
			return
		}

		// Update UI
		app.selectedFile = path
		app.fileLabel.SetText(filepath.Base(path))
		app.sizeLabel.SetText(fmt.Sprintf("%s • %dx%d • %s",
			compressor.FormatBytes(info.Size),
			info.Width,
			info.Height,
			format,
		))

		// Load preview
		img := canvas.NewImageFromURI(uri)
		img.FillMode = canvas.ImageFillContain
		app.imagePreview.Resource = nil
		app.imagePreview.File = path
		app.imagePreview.Refresh()

		app.compressBtn.Enable()
		app.statusLabel.SetText("")

	}, app.window)
}

func (app *ImgShrinkApp) selectOutputLocation() {
	dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
		if err != nil {
			dialog.ShowError(err, app.window)
			return
		}
		if uri == nil {
			return
		}

		app.outputPath = uri.Path()
		app.outputLabel.SetText(fmt.Sprintf("Output: %s", filepath.Base(app.outputPath)))
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

		// Set output directory if selected
		if app.outputPath != "" {
			options.OutputDir = app.outputPath
		}

		// Handle resize options based on mode
		switch app.resizeMode.Selected {
		case "Percentage":
			if app.resizeSlider.Value < 100 {
				options.ResizePercent = app.resizeSlider.Value
			}
		case "Dimensions":
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

		successMsg := fmt.Sprintf("✓ Compressed successfully!\n"+
			"%s → %s\n"+
			"Saved: %s (%.1f%% reduction)\n"+
			"Output: %s",
			compressor.FormatBytes(result.InputSize),
			compressor.FormatBytes(result.OutputSize),
			compressor.FormatBytes(result.InputSize-result.OutputSize),
			result.Reduction,
			result.OutputPath,
		)

		app.statusLabel.SetText(successMsg)
		app.progressBar.Hide()
		app.compressBtn.Enable()

		// Show success dialog
		dialog.ShowInformation("Success", successMsg, app.window)
	}()
}
