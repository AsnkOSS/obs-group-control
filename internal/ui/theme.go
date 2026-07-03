package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// largeTheme wraps the default theme with bigger text and spacing so the
// app stays readable on high-DPI (4K) displays.
type largeTheme struct {
	fyne.Theme
}

func newTheme() fyne.Theme {
	return &largeTheme{Theme: theme.DefaultTheme()}
}

func (t *largeTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 16
	case theme.SizeNameHeadingText:
		return 30
	case theme.SizeNameSubHeadingText:
		return 22
	case theme.SizeNamePadding:
		return 6
	case theme.SizeNameInnerPadding:
		return 10
	default:
		return t.Theme.Size(name)
	}
}

func (t *largeTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return t.Theme.Color(name, variant)
}
