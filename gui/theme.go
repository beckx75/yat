package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

type yatTheme struct {
	fyne.Theme
}

func newYatTheme() fyne.Theme {
	return &yatTheme{Theme: theme.DefaultTheme()}
}

func (yt *yatTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return yt.Theme.Color(name, theme.VariantLight)
}

func (yt *yatTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameText {
		return 20
	}
	return yt.Theme.Size(name)
}
