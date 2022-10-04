package element

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

/*

  File:    theme.go
  Author:  Bob Shofner

  Copyright (c) 2022. BSD 3-Clause License
	https://opensource.org/licenses/BSD-3-Clause

  The this permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:
*/

var _ fyne.Theme = (*ScsiTheme)(nil)
var ScsiSizeNameText float32 = 12
var ScsiSizeNamePadding float32 = 2

type ScsiTheme struct {
}

//goland:noinspection GoUnusedExportedFunction
func NewTheme(kc KeyChain) *ScsiTheme {
	t := &ScsiTheme{}
	kp := make(chan KeyPressed)
	kc.Register(kp)
	go func() { // process key event
		for p := range kp {
			t.keyHandler(p)
		}

	}()
	return t
}
func (t *ScsiTheme) keyHandler(key KeyPressed) {
	switch key.Name {
	case fyne.KeyPlus:
		ScsiSizeNameText++
	case fyne.KeyMinus:
		if ScsiSizeNameText > 6 {
			ScsiSizeNameText--
		}
	}
	return
}
func (t ScsiTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameScrollBar:
		return theme.PrimaryColor()
	}
	return theme.DefaultTheme().Color(name, variant)
}
func (t ScsiTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t ScsiTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}
func (t ScsiTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return ScsiSizeNamePadding
	case theme.SizeNameText:
		return ScsiSizeNameText
	}
	return theme.DefaultTheme().Size(name)
}
