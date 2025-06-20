package gui

import(
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
)

func (ui *UI)makeMaincontent() fyne.CanvasObject {
	return container.NewGridWithRows(2,
		ui.makeFilescontent(),
		ui.makeTagcontent(),
	)
}

func (ui *UI) makeFilescontent() fyne.CanvasObject {
	ui.files.Append("eins")
	ui.files.Append("eins")
	ui.files.Append("eins")
	ui.filelist = widget.NewListWithData(ui.files,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})
	ui.filelist.OnSelected = func(id widget.ListItemID){
		val, _ := ui.files.GetValue(id)
		fmt.Println(val)
	}
	return container.NewGridWithColumns(3,
		ui.filelist, 
		widget.NewLabel("Taginfo"), 
		widget.NewLabel("Cover"), 
	)
}

func (ui *UI) makeTagcontent() fyne.CanvasObject {
	things := make([]string, 1000)
	for i := range things{
		things[i] = fmt.Sprintf("%d",i)
	}
	list := widget.NewList(
		func() int {
			return len(things)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewCheck("template", func(b bool){}), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			fmt.Println(item.(*fyne.Container).Objects[0].(*widget.Check).Checked)
		},
	)
	
	
	return container.NewGridWithColumns(4,
		list, 
		widget.NewLabel("Frame Value"), 
		widget.NewLabel("New Frame Value"),
		widget.NewLabel("Edit Actions"),
	)
}
