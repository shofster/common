package element

import "fyne.io/fyne/v2/widget"

/*

  File:    choice.go
  Author:  Bob Shofner

  Copyright (c) 2022. BSD 3-Clause License
	https://opensource.org/licenses/BSD-3-Clause

  The this permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:	Choice is a widget that represents a Select with a pre-selected entry
		and no PlaceHolder

*/

type Choice struct {
	widget.Select
	def string
}

//goland:noinspection GoUnusedExportedFunction
func NewChoice(options []string, def string) *Choice {
	c := &Choice{}
	c.ExtendBaseWidget(c)
	c.Options = options
	c.def = def
	c.Selected = def
	return c
}

func (c *Choice) Reset() {
	c.Selected = c.def
	c.Refresh()
}

func (c *Choice) Choose(def string) {
	c.def = def
	c.Selected = c.def
	c.Refresh()
}
