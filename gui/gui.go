package gui

import(
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
)

type UI struct {
	app fyne.App
	win fyne.Window

	filelist *widget.List
	files binding.StringList
}

func InitGui(args []string) {
	ui := new(UI)
	ui.app = app.New()
	ui.app.Settings().SetTheme(newYatTheme())
	ui.win = ui.app.NewWindow("this is yat...")

	ui.files = binding.BindStringList(
		&[]string{},
	)

	maincontent := ui.makeMaincontent()

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.HomeIcon(), func(){fmt.Println("feeling like home...")}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.LogoutIcon(), func(){ui.app.Quit()}),
	)

	mainbox := container.NewBorder(toolbar, nil, nil, nil,
		maincontent,
	)

	ui.win.SetOnDropped(
		func(p fyne.Position, uris []fyne.URI){
			for _, uri := range uris {
				fmt.Println(uri)
			}
		})
	
	ui.win.SetContent(mainbox)
	ui.win.Resize(fyne.NewSize(1024, 768))
	ui.win.CenterOnScreen()
	ui.win.ShowAndRun()
}
