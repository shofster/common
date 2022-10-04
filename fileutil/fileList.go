package fileutil

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"io/fs"
	"strings"
	"time"
)

/*

  File:    fileList.go
  Author:  Bob Shofner

  Copyright (c) 2022. BSD 3-Clause License
	https://opensource.org/licenses/BSD-3-Clause

  The this permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*

 */

var DefaultDoubleClickTime = time.Millisecond * 500
var DefaultMonospace = false
var DefaultIconType = CheckBoxType

//goland:noinspection GoUnusedGlobalVariable,SpellCheckingInspection
var RFCtypes = []string{"Default", "UNIX", "RFC822", "RFC822Z", "RFC1123", "RFC1123Z", "RFC3339"}
var defaultDateTimeFormat = "01/02/06 15:04"
var dtFormat = ""

func NewFileList(path string, sel FileSelectFilter, action FileSelectAction,
	pretty func(name string, info fs.FileInfo, err error) string) (*widget.List, *DirectoryEntry, error) {
	// get specific DirectoryEntry
	dir, err := NewDirectoryView(path, sel)
	if err != nil {
		return nil, nil, err
	}
	dtFormat = sel.DtFormat
	if pretty == nil { // use standard display
		pretty = ls_al
	}
	lastSelected := -1
	lastTime := time.Now()
	fileList := widget.NewList(
		// length
		func() int {
			return dir.Count()
		},
		// create
		func() fyne.CanvasObject {
			switch DefaultIconType {
			case FileIconType:
				return container.NewHBox(
					widget.NewIcon(theme.FileIcon()),
					widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: DefaultMonospace}))
			}
			return container.NewHBox(
				widget.NewCheck("", nil),
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: DefaultMonospace}))
		},
		// update
		func(id widget.ListItemID, item fyne.CanvasObject) {
			file := dir.File(id)
			//log.Println("UPDATE ON ", file.String())
			info, err := file.Info()
			item.(*fyne.Container).Objects[1].(*widget.Label).Text = pretty(file.DisplayName(), info, err)
			switch DefaultIconType {
			case FileIconType:
				if file.IsSelected() {
					item.(*fyne.Container).Objects[0].(*widget.Icon).SetResource(theme.ConfirmIcon())
				} else {
					if file.IsDir() {
						item.(*fyne.Container).Objects[0].(*widget.Icon).SetResource(theme.FolderIcon())
					} else {
						item.(*fyne.Container).Objects[0].(*widget.Icon).SetResource(theme.FileIcon())
					}
				}
			default:
				item.(*fyne.Container).Objects[0].(*widget.Check).SetChecked(file.IsSelected())
			}
			item.(*fyne.Container).Objects[0].Refresh()
			item.(*fyne.Container).Objects[1].Refresh()
		})
	fileList.OnSelected = func(id widget.ListItemID) {
		file := dir.File(id) // current file POINTER
		//log.Printf("Selected %d %s, IsDir %t, IsSelected %t, FileType %d, Multiple %t, lastId %d",
		//	id, file.DisplayName(), file.IsDir(), file.IsSelected(), sel.FileType, sel.Multiple,
		//	lastSelected)
		duration := time.Now().Sub(lastTime)
		if file.IsDir() && lastSelected == id && duration <= DefaultDoubleClickTime {
			//		log.Printf(" double click ON %d", id)
			if action.OnDoubleClick != nil {
				action.OnDoubleClick(*file)
			}
			// double click - clear all selected
			files := dir.files
			for ix := range files {
				fp := dir.File(ix)
				if fp.selected {
					fp.SetSelected(false)
					fileList.Unselect(ix)
				}
			}
			lastSelected = -1
		} else if !file.IsDir() && lastSelected == id && duration <= DefaultDoubleClickTime {
			if action.OnDoubleClick != nil {
				action.OnDoubleClick(*file)
			}
		} else {
			if sel.FileType == File && file.IsDir() { // must have File
				fileList.Unselect(id)
			} else {
				if sel.Multiple {
					file.SetSelected(!file.selected)
				} else { // only a single File or Dir, toggle select
					if !file.selected {
						// clear all others
						files := dir.files
						for ix := range files {
							fp := dir.File(ix)
							if fp.selected {
								fp.SetSelected(false)
								fileList.Unselect(ix)
							}
						}
						file.SetSelected(true)
					} else {
						file.SetSelected(false)
					}
				}
			}
			// log.Printf("%s -> set to %t", file.DisplayName(), file.IsSelected())
			lastSelected = id
			lastTime = time.Now()
			if action.OnClick != nil {
				file.index = id
				action.OnClick(*file)
			}
		}
		fileList.Unselect(id) // allow multiple AND double click
		fileList.Refresh()
	}
	return fileList, dir, nil
}

// ls_al. LINUX ls -Al
//goland:noinspection GoSnakeCaseUsage,SpellCheckingInspection
func ls_al(name string, info fs.FileInfo, err error) string {
	if err != nil {
		return fmt.Sprintf("Unable to get FileInfo for %s", name)
	}
	dfmt := defaultDateTimeFormat
	switch strings.ToUpper(dtFormat) {
	case "UNIX":
		dfmt = time.UnixDate
	case "RFC822":
		dfmt = time.RFC822
	case "RFC822Z":
		dfmt = time.RFC822Z
	case "RFC1123":
		dfmt = time.RFC1123
	case "RFC1123Z":
		dfmt = time.RFC1123Z
	case "RFC3339":
		dfmt = time.RFC3339
	}
	dt := info.ModTime().Format(dfmt)
	m := fmt.Sprintf("%s", info.Mode())
	f := "%4s %10d %s %s"
	return fmt.Sprintf(f, m[0:4], info.Size(), dt, name)
}
