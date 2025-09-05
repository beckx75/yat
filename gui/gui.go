package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	// "fyne.io/fyne/v2/data/binding"
	"beckx.online/butils/fileutils"
)

type UI struct {
	app fyne.App
	win fyne.Window

	thefiles []*fileutils.TheFile
	files    *widget.List
}

func InitGui(args []string) {
	ui := new(UI)

	ui.app = app.New()
	ui.app.Settings().SetTheme(newYatTheme())
	ui.win = ui.app.NewWindow("this is yat...")

	var err error
	if len(args) > 0 {
		ui.thefiles, _, err = fileutils.GetFiles(args, []string{".mp3", ".flac"})
		if err != nil {
			dialog.ShowError(err, ui.win)
		}
	} else {
		ui.thefiles = []*fileutils.TheFile{}
	}

	cntFiles := ui.initUiFiles()
	cntFrames := ui.makeTagcontent()

	maincontent := container.NewGridWithRows(2,
		cntFiles, cntFrames,
	)

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.HomeIcon(), func() { fmt.Println("feeling like home...") }),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.LogoutIcon(), func() { ui.app.Quit() }),
	)

	mainbox := container.NewBorder(toolbar, nil, nil, nil,
		maincontent,
	)

	ui.win.SetOnDropped(
		func(p fyne.Position, uris []fyne.URI) {
			for _, uri := range uris {
				fmt.Println(uri)
				fv := fileutils.NewTheFile(uri.Path())
				ui.thefiles = append(ui.thefiles, fv)
			}
		})

	ui.win.SetContent(mainbox)
	ui.win.Resize(fyne.NewSize(1024, 768))
	ui.win.CenterOnScreen()
	ui.win.ShowAndRun()
}
