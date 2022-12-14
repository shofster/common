package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/shofster/common/fileutil"
	"log"
)

func main() {
	_ = path.Set("c:")
	a := app.New()
	appWindow = a.NewWindow("MENUS")
	item1 := fyne.NewMenuItem("item1", func() {
		log.Println("item1")
	})
	item2 := fyne.NewMenuItem("item2", func() {
		log.Println("item2")
	})
	file := fyne.NewMenu("File", item1, item2)
	mainMenu := fyne.NewMainMenu(file)
	appWindow.SetMainMenu(mainMenu)
	buttDir := widget.NewButton("select dir", func() {
		testDirOpen()
	})
	fileDir := widget.NewButton("select file", func() {
		testFileOpen()
	})
	fileutil.DefaultIconType = fileutil.FileIconType
	c := container.NewVBox(buttDir, fileDir)
	appWindow.SetContent(c)
	appWindow.Resize(fyne.NewSize(800, 600))
	appWindow.ShowAndRun()
}

var path = binding.NewString()
var appWindow fyne.Window

func testDirOpen() {
	placeSel := fileutil.FileSelectFilter{
		Title:      "Select a Place",
		FileType:   fileutil.Dir,
		FileSelect: fileutil.Open,
		Multiple:   false,
		Hidden:     ""}
	last := "E:\\charon\\z"
	_ = path.Set(last)
	fileutil.FileSelect(placeSel, path, appWindow, func(dirs []string) {
		log.Printf("%d selected\n", len(dirs))
		if len(dirs) > 0 {
			fmt.Println("dir", dirs[0])
		}
	})
}
func testFileOpen() {
	placeSel := fileutil.FileSelectFilter{
		Title:      "Select a File",
		FileType:   fileutil.File,
		FileSelect: fileutil.Open,
		Multiple:   false,
		Hidden:     ""}
	last := "E:\\charon\\z"
	_ = path.Set(last)
	fileutil.FileSelect(placeSel, path, appWindow, func(files []string) {
		log.Printf("%d selected\n", len(files))
		if len(files) > 0 {
			fmt.Println("dir", files[0])
		}
	})
}
