package element

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

/*

  File:    console.go
  Author:  Bob Shofner

  Copyright (c) 2022. BSD 3-Clause License
	https://opensource.org/licenses/BSD-3-Clause

  The this permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: Widget to add functionality to widget.Accordion.
*/

const windBoxDividerHeight = 1

var _ fyne.Widget = (*WindBox)(nil)

// WindBox replaces Accordion displays a list of WindBoxItems.
// Each item is represented by a button that reveals a detailed view when tapped.
type WindBox struct {
	widget.BaseWidget
	Items     []*WindBoxItem
	MultiOpen bool
}

// NewWindBox creates a new WindBox widget.
//goland:noinspection GoUnusedExportedFunction,GoUnusedExportedFunction
func NewWindBox(items ...*WindBoxItem) *WindBox {
	a := &WindBox{
		Items: items,
	}
	a.ExtendBaseWidget(a)
	return a
}

// Append adds the given item to this Accordion.
func (w *WindBox) Append(item *WindBoxItem) {
	w.Items = append(w.Items, item)
	w.Refresh()
}

// Close collapses the item at the given index.
func (w *WindBox) Close(index int) {
	if index < 0 || index >= len(w.Items) {
		return
	}
	w.Items[index].Open = false
	w.Items[index].Detail.Hide()
	w.Refresh()
}

// CloseAll collapses all items.
func (w *WindBox) CloseAll() {
	for _, i := range w.Items {
		i.Open = false
	}
	w.Refresh()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (w *WindBox) CreateRenderer() fyne.WidgetRenderer {
	w.ExtendBaseWidget(w)
	r := &windBoxRenderer{
		container: w,
	}
	r.updateObjects()
	return r
}

// MinSize returns the size that this widget should not shrink below.
func (w *WindBox) MinSize() fyne.Size {
	w.ExtendBaseWidget(w)
	return w.BaseWidget.MinSize()
}

// Open expands the item at the given index.
// and collapses the others
func (w *WindBox) Open(index int) {
	if index < 0 || index >= len(w.Items) {
		return
	}
	for i, ai := range w.Items {
		if i == index {
			ai.Open = true
			ai.Detail.Show()
			ai.Detail.Refresh()
		} else if !w.MultiOpen {
			ai.Open = false
			ai.Detail.Hide()
		}
	}
	w.Refresh()
}

type windBoxRenderer struct {
	objects   []fyne.CanvasObject
	container *WindBox
	headers   []*fyne.Container
	dividers  []fyne.CanvasObject
}

func (r *windBoxRenderer) Layout(size fyne.Size) {
	x := float32(0)
	y := float32(0)
	for i, ai := range r.container.Items {
		if i != 0 {
			div := r.dividers[i-1]
			div.Move(fyne.NewPos(x, y))
			div.Resize(fyne.NewSize(size.Width, windBoxDividerHeight))
			y += windBoxDividerHeight
		}
		y += theme.Padding()

		h := r.headers[i]
		h.Move(fyne.NewPos(x, y))
		min := h.MinSize().Height
		h.Resize(fyne.NewSize(size.Width, min))
		y += min
		if ai.Open {
			y += theme.Padding()
			d := ai.Detail
			d.Move(fyne.NewPos(x, y))
			min := d.MinSize().Height
			d.Resize(fyne.NewSize(size.Width, min))
			y += min
		}

		y += theme.Padding()
	}
}

func (r windBoxRenderer) MinSize() (size fyne.Size) {
	for i, ai := range r.container.Items {
		size.Height += theme.Padding() * 2
		if i != 0 {
			size.Height += windBoxDividerHeight
		}
		min := r.headers[i].MinSize()
		size.Width = fyne.Max(size.Width, min.Width)
		size.Height += min.Height
		min = ai.Detail.MinSize()
		size.Width = fyne.Max(size.Width, min.Width)
		if ai.Open {
			size.Height += min.Height
			size.Height += theme.Padding()
		}
	}
	return
}

func (r *windBoxRenderer) Refresh() {
	r.updateObjects()
	r.Layout(r.container.Size())
	canvas.Refresh(r.container)
}
func (r *windBoxRenderer) Destroy() {
}
func (r *windBoxRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}
func (r *windBoxRenderer) SetObjects(objects []fyne.CanvasObject) {
	r.objects = objects
}

func (r *windBoxRenderer) updateObjects() {
	is := len(r.container.Items)
	hs := len(r.headers)
	ds := len(r.dividers)
	i := 0
	for ; i < is; i++ {
		ci := r.container.Items[i]
		var h = container.NewHBox()
		if i < hs { // already exists
			h = r.headers[i]
		} else { // a new Item
			b := &widget.Button{}
			b.Alignment = widget.ButtonAlignLeading
			b.IconPlacement = widget.ButtonIconLeadingText
			b.Hidden = false
			b.Importance = widget.LowImportance
			b.Text = ci.Title
			index := i // capture
			b.OnTapped = func() {
				if ci.Open {
					r.container.Close(index)
				} else {
					r.container.Open(index)
				}
			}
			if ci.Open {
				b.Icon = theme.MenuDropUpIcon()
				ci.Detail.Show()
			} else {
				b.Icon = theme.MenuDropDownIcon()
				ci.Detail.Hide()
			}
			h = container.NewHBox(b)
			ac := ci.Action
			if ac != nil {
				h.Objects = append([]fyne.CanvasObject{ac}, h.Objects...)
			}
			b.Show()
			r.headers = append(r.headers, h)
			hs++
		}
	}
	// Hide extras
	for ; i < hs; i++ {
		r.headers[i].Hide()
	}
	// Set objects
	objects := make([]fyne.CanvasObject, hs+is+ds)
	for i, header := range r.headers {
		objects[i] = header
	}
	for i, item := range r.container.Items {
		objects[hs+i] = item.Detail
	}
	// add dividers
	for i = 0; i < ds; i++ {
		if i < len(r.container.Items)-1 {
			r.dividers[i].Show()
		} else {
			r.dividers[i].Hide()
		}
		objects[hs+is+i] = r.dividers[i]
	}
	// make new dividers
	for ; i < is-1; i++ {
		div := widget.NewSeparator()
		r.dividers = append(r.dividers, div)
		objects = append(objects, div)
	}
	r.SetObjects(objects)
}

// WindBoxItem AccordionItem represents a single item in an Accordion.
// it has a Button as a header
type WindBoxItem struct {
	Title  string
	Action *fyne.Container
	Detail fyne.CanvasObject
	Open   bool
}

// NewWindBoxItem NewAccordionItem creates a new item for an Accordion.
//goland:noinspection GoUnusedExportedFunction
func NewWindBoxItem(title string, detail fyne.CanvasObject) *WindBoxItem {
	return &WindBoxItem{
		Title:  title,
		Detail: detail,
	}
}
