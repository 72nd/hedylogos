package gui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func Run() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Toolbar Widget")
	myWindow.SetMainMenu(fyne.NewMainMenu(fyne.NewMenu("File", fyne.NewMenuItem("Open...", nil))))

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.FolderOpenIcon(), func() {
			chooser := dialog.NewFileOpen(callback, myWindow)
			chooser.SetFilter(storage.NewExtensionFileFilter([]string{".graphml"}))
			chooser.Show()
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.MediaPlayIcon(), func() {}),
		widget.NewToolbarAction(theme.MediaStopIcon(), func() {}),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			log.Println("Display help")
		}),
	)

	input := widget.NewEntry()
	bottom := container.NewBorder(nil, nil, nil, widget.NewButton("Go", nil), input)
	prompt := widget.NewEntry()
	prompt.SetText("[Open a graphml file to start]")
	prompt.Disable()
	content := container.NewBorder(toolbar, bottom, nil, nil, prompt)
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

func callback(foo fyne.URIReadCloser, err error) {
	fmt.Println(foo.URI())
}
