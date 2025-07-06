package gui

import(
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	// "fyne.io/fyne/v2/data/binding"
)

func (ui *UI) makeTagcontent() fyne.CanvasObject {
	things := []string{"a", "b", "c"}
	list := widget.NewList(
		func() int {
			return len(things)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewCheck("", func(b bool){
				fmt.Println("ich weiß nicht wo ich bin :(")
			}), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			fmt.Println("ka leck mich doch", id)
			// item.(*fyne.Container).Objects[1].(*widget.Label).Text = things[id]
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		fmt.Println("selected something...")
	}
	
	
	return container.NewGridWithColumns(4,
		list, 
		widget.NewLabel("Frame Value"), 
		widget.NewLabel("New Frame Value"),
		widget.NewLabel("Edit Actions"),
	)
}
