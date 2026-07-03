// Package ui builds the Fyne window and wires it to the controller.
package ui

import (
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"obs-group-control/internal/obs"
	"obs-group-control/internal/state"
)

// App owns the Fyne application and widgets.
type App struct {
	fyneApp    fyne.App
	window     fyne.Window
	controller *obs.Controller
	store      *state.Store
	lang       Lang

	titleText    *canvas.Text
	subtitleText *canvas.Text
	statusText   *canvas.Text
	countLabel   *widget.Label
	failedText   *canvas.Text
	msgLabel     *widget.Label
	startBtn     *widget.Button
	stopBtn      *widget.Button
	langSelect   *widget.Select
}

// New creates the application window for the given controller.
func New(controller *obs.Controller) *App {
	a := &App{
		fyneApp:    app.New(),
		controller: controller,
		lang:       LangZH,
	}
	a.window = a.fyneApp.NewWindow(translations[a.lang].windowTitle)
	a.store = state.NewStore(controller.DeviceCount(), a.render)
	a.buildUI()
	a.render(a.store.Get())
	a.controller.StartMonitor(a.store.SetFailedCount)
	return a
}

// Run shows the window and blocks until it is closed.
func (a *App) Run() {
	a.window.Resize(fyne.NewSize(520, 400))
	a.window.ShowAndRun()
}

func (a *App) buildUI() {
	a.titleText = canvas.NewText("", color.White)
	a.titleText.TextSize = 30
	a.titleText.TextStyle = fyne.TextStyle{Bold: true}
	a.titleText.Alignment = fyne.TextAlignCenter

	a.subtitleText = canvas.NewText("", color.NRGBA{R: 0xB0, G: 0xB0, B: 0xB0, A: 0xFF})
	a.subtitleText.TextSize = 15
	a.subtitleText.Alignment = fyne.TextAlignCenter

	a.statusText = canvas.NewText("", color.White)
	a.statusText.TextSize = 26
	a.statusText.TextStyle = fyne.TextStyle{Bold: true}
	a.statusText.Alignment = fyne.TextAlignCenter

	a.countLabel = widget.NewLabel("")
	a.countLabel.Alignment = fyne.TextAlignCenter

	a.failedText = canvas.NewText("", color.NRGBA{R: 0xE5, G: 0x39, B: 0x35, A: 0xFF})
	a.failedText.TextSize = 14
	a.failedText.Alignment = fyne.TextAlignCenter

	a.msgLabel = widget.NewLabel("")
	a.msgLabel.Alignment = fyne.TextAlignCenter
	a.msgLabel.Wrapping = fyne.TextWrapWord

	a.startBtn = widget.NewButton("", a.onStart)
	a.startBtn.Importance = widget.HighImportance
	a.stopBtn = widget.NewButton("", a.onStop)

	a.langSelect = widget.NewSelect([]string{"中文", "English"}, func(sel string) {
		if sel == "English" {
			a.lang = LangEN
		} else {
			a.lang = LangZH
		}
		a.applyLang()
	})
	a.langSelect.SetSelectedIndex(0)

	counts := container.NewHBox(
		layout.NewSpacer(),
		a.countLabel,
		container.NewCenter(a.failedText),
		layout.NewSpacer(),
	)

	buttons := container.NewGridWithColumns(2, a.startBtn, a.stopBtn)
	content := container.NewVBox(
		container.NewHBox(layout.NewSpacer(), a.langSelect),
		widget.NewLabel(""), // breathing room below the menu/language row
		a.titleText,
		a.subtitleText,
		widget.NewLabel(""),
		a.statusText,
		counts,
		a.msgLabel,
		buttons,
	)
	pad := layout.NewCustomPaddedLayout(16, 16, 24, 24)
	a.window.SetContent(container.New(pad, content))
	a.applyLang()
}

// applyLang re-renders every static text in the current language.
func (a *App) applyLang() {
	t := translations[a.lang]
	a.window.SetTitle(t.windowTitle)
	a.titleText.Text = t.windowTitle
	a.titleText.Refresh()
	a.subtitleText.Text = t.subtitle
	a.subtitleText.Refresh()
	a.startBtn.SetText(t.start)
	a.stopBtn.SetText(t.stop)
	a.render(a.store.Get())
}

// onStart starts recording on all devices in the background.
func (a *App) onStart() {
	a.store.Set(state.Starting, "")
	go func() {
		if err := a.controller.StartAll(); err != nil {
			a.store.Set(state.Failed, err.Error())
			return
		}
		// Only reached after every device confirmed it is recording.
		a.store.Set(state.Recording, "")
	}()
}

// onStop stops recording on all devices in the background.
func (a *App) onStop() {
	a.store.Set(state.Stopping, "")
	go func() {
		if err := a.controller.StopAll(); err != nil {
			a.store.Set(state.Failed, err.Error())
			return
		}
		a.store.Set(state.Idle, "")
	}()
}

// render updates widgets from a snapshot; safe to call from any goroutine.
func (a *App) render(snap state.Snapshot) {
	fyne.Do(func() {
		t := translations[a.lang]

		a.statusText.Text = t.phases[snap.Phase]
		a.statusText.Color = phaseColor(snap.Phase)
		a.statusText.Refresh()

		a.countLabel.SetText(t.devices + ": " + strconv.Itoa(snap.DeviceCount))
		if snap.FailedCount > 0 {
			a.failedText.Text = t.failed + ": " + strconv.Itoa(snap.FailedCount)
		} else {
			a.failedText.Text = ""
		}
		a.failedText.Refresh()

		a.msgLabel.SetText(snap.Message)

		busy := snap.Phase == state.Starting || snap.Phase == state.Stopping
		if busy {
			a.startBtn.Disable()
			a.stopBtn.Disable()
		} else if snap.Phase == state.Recording {
			a.startBtn.Disable()
			a.stopBtn.Enable()
		} else {
			a.startBtn.Enable()
			a.stopBtn.Enable()
		}
	})
}

func phaseColor(p state.Phase) color.Color {
	switch p {
	case state.Recording:
		return color.NRGBA{R: 0xE5, G: 0x39, B: 0x35, A: 0xFF}
	case state.Failed:
		return color.NRGBA{R: 0xFF, G: 0xA0, B: 0x00, A: 0xFF}
	default:
		return color.White
	}
}
