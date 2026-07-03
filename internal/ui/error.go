package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// RunConfigError shows a small window explaining why the device config
// could not be loaded, then exits when the user closes it. It is used
// instead of log.Fatal so the message is visible in GUI-only builds.
func RunConfigError(configPath string, loadErr error) {
	fyneApp := app.New()
	fyneApp.Settings().SetTheme(newTheme())
	t := translations[LangZH]

	window := fyneApp.NewWindow(t.configErrTitle)

	title := canvas.NewText(t.configErrTitle, color.NRGBA{R: 0xE5, G: 0x39, B: 0x35, A: 0xFF})
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	hint := widget.NewLabel(t.configErrHint + " " + configPath)
	hint.Alignment = fyne.TextAlignCenter
	hint.Wrapping = fyne.TextWrapWord

	detail := widget.NewLabel(loadErr.Error())
	detail.Alignment = fyne.TextAlignCenter
	detail.Wrapping = fyne.TextWrapWord

	quit := widget.NewButton(t.quit, fyneApp.Quit)
	quit.Importance = widget.HighImportance

	window.SetContent(container.NewPadded(container.NewVBox(
		widget.NewLabel(""),
		title,
		hint,
		detail,
		widget.NewLabel(""),
		quit,
	)))
	window.Resize(fyne.NewSize(520, 320))
	window.CenterOnScreen()
	window.ShowAndRun()
}
