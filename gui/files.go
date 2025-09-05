package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	// "fyne.io/fyne/v2/driver/desktop"
	// "fyne.io/fyne/v2/theme"
)

func (ui *UI) initUiFiles() fyne.CanvasObject {
	ui.files = widget.NewList(
		func() int {
			return len(ui.thefiles)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(id widget.ListItemID, o fyne.CanvasObject) {
			usebold := ui.thefiles[id].Selected
			o.(*widget.Label).TextStyle = fyne.TextStyle{Bold: usebold}
			o.(*widget.Label).SetText(ui.thefiles[id].Name)
		})

	ui.files.OnSelected = func(id widget.ListItemID) {
		if ui.thefiles[id].Selected {
			ui.thefiles[id].Selected = false
		} else {
			ui.thefiles[id].Selected = true
		}
		ui.files.UnselectAll()
		ui.files.Refresh()
	}

	lbl := widget.NewLabel("File List")
	lbl.Alignment = fyne.TextAlignCenter
	lbl.TextStyle.Bold = true

	return container.NewBorder(nil, nil, nil, nil,
		ui.files,
	)
}
