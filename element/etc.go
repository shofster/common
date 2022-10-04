package element

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"net/url"
	"os/exec"
)

/*

  File:    etc.go
  Author:  Bob Shofner

  Copyright (c) 2022. BSD 3-Clause License
	https://opensource.org/licenses/BSD-3-Clause

  The this permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:
*/

// Notify provides system level Notification
//goland:noinspection GoUnusedExportedFunction
func Notify(app fyne.App, title, content string) {
	n := fyne.NewNotification(title, content)
	app.SendNotification(n)
}

// VerifyDelete gets a "delete" confirmation. Modal to the fyne.Window.
//   normally the application's window.
//goland:noinspection GoUnusedExportedFunction
func VerifyDelete(appWindow fyne.Window, content string, res func(bool)) {
	dialog.ShowConfirm("Is it OK to completely remove?", content, res, appWindow)
}
func urlFromString(str string) (*url.URL, error) {
	u, err := url.Parse(str)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Browse sends a file to the App
//goland:noinspection GoUnusedExportedFunction
func Browse(file string) error {
	u, err := urlFromString(file)
	if err != nil {
		return err
	}
	return fyne.CurrentApp().OpenURL(u)
}

// PlayFile sends a file to the player command
//goland:noinspection GoUnusedExportedFunction
func PlayFile(cmd, file string) error {
	command := exec.Command(cmd, file)
	return command.Start()
}
