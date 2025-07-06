package gui

import(
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	// "fyne.io/fyne/v2/data/binding"
)

func (ui *UI) makeFiles() fyne.CanvasObject {
	ui.filelist = widget.NewList(
		func() int {
			return len(ui.files)
		},
		func () fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("template"),
				widget.NewEntry(),
			)
		},
		func(id widget.ListItemID, o fyne.CanvasObject) {
			o.(*fyne.Container).Objects[0].(*widget.Label).SetText(ui.files[id].Name)
		})
	ui.filelist.OnSelected = func(id widget.ListItemID){
		val := ui.files[id].Path
		fmt.Println(val)
	}
	return container.NewBorder(nil, nil, nil, nil,
		ui.filelist, 
	)
}
