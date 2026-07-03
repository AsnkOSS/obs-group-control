// Package ui builds the Fyne window and wires it to the controller.
package ui

import (
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
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

	statusText *canvas.Text
	countLabel *widget.Label
	msgLabel   *widget.Label
	startBtn   *widget.Button
	stopBtn    *widget.Button
}

// New creates the application window for the given controller.
func New(controller *obs.Controller) *App {
	a := &App{
		fyneApp:    app.New(),
		controller: controller,
	}
	a.window = a.fyneApp.NewWindow("OBS Group Control")
	a.store = state.NewStore(controller.DeviceCount(), a.render)
	a.buildUI()
	a.render(a.store.Get())
	return a
}

// Run shows the window and blocks until it is closed.
func (a *App) Run() {
	a.window.Resize(fyne.NewSize(360, 220))
	a.window.ShowAndRun()
}

func (a *App) buildUI() {
	a.statusText = canvas.NewText("", color.White)
	a.statusText.TextSize = 22
	a.statusText.TextStyle = fyne.TextStyle{Bold: true}
	a.statusText.Alignment = fyne.TextAlignCenter

	a.countLabel = widget.NewLabel("")
	a.countLabel.Alignment = fyne.TextAlignCenter

	a.msgLabel = widget.NewLabel("")
	a.msgLabel.Alignment = fyne.TextAlignCenter
	a.msgLabel.Wrapping = fyne.TextWrapWord

	a.startBtn = widget.NewButton("Start", a.onStart)
	a.stopBtn = widget.NewButton("Stop", a.onStop)

	buttons := container.NewGridWithColumns(2, a.startBtn, a.stopBtn)
	content := container.NewVBox(
		a.statusText,
		a.countLabel,
		a.msgLabel,
		buttons,
	)
	a.window.SetContent(container.NewPadded(content))
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
		a.statusText.Text = snap.Phase.String()
		a.statusText.Color = phaseColor(snap.Phase)
		a.statusText.Refresh()

		a.countLabel.SetText(deviceCountText(snap.DeviceCount))
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

func deviceCountText(n int) string {
	return "Devices: " + strconv.Itoa(n)
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
