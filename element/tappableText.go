package element

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

/*

  File:    tappableText.go
  Author:  Bob Shofner

  Copyright (c) 2022. BSD 3-Clause License
	https://opensource.org/licenses/BSD-3-Clause

  The this permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
	An EmptyWidget is a widget that wraps CanvasObjects to make them "Tappable" and "Mouseable"
	TappableText is a canvas.Text wrapped in an EmptyWidget
	MouseableHover is a canvas.CanvasObject wrapped in an EmptyWidget
*/

type EmptyWidget struct {
	widget.BaseWidget
	Objects           []fyne.CanvasObject
	OnTapped          func(*fyne.PointEvent)
	OnTappedSecondary func(*fyne.PointEvent)
	OnDoubleTapped    func(*fyne.PointEvent)
	OnMouseIn         func(*desktop.MouseEvent)
	OnMouseOut        func()
}

//goland:noinspection GoUnusedExportedFunction
func NewEmptyWidget(onTapped func(event *fyne.PointEvent), onTappedSecondary func(event *fyne.PointEvent)) *EmptyWidget {
	em := &EmptyWidget{OnTapped: onTapped, OnTappedSecondary: onTappedSecondary}
	em.Objects = make([]fyne.CanvasObject, 0)
	return em
}
func (e *EmptyWidget) AddObject(o fyne.CanvasObject) {
	e.Objects = append(e.Objects, o)
}

func (e *EmptyWidget) Tapped(pe *fyne.PointEvent) {
	if e.OnTapped != nil {
		e.OnTapped(pe)
	}
}
func (e *EmptyWidget) TappedSecondary(pe *fyne.PointEvent) {
	if e.OnTappedSecondary != nil {
		e.OnTappedSecondary(pe)
	}
}
func (e *EmptyWidget) DoubleTapped(pe *fyne.PointEvent) {
	if e.OnDoubleTapped != nil {
		e.OnDoubleTapped(pe)
	}
}

// Implements: desktop.Hoverable
var _ desktop.Hoverable = (*EmptyWidget)(nil)

// Implements: fyne.Tappable
var _ fyne.Tappable = (*EmptyWidget)(nil)

// Implements: fyne.Focusable
var _ fyne.Focusable = (*EmptyWidget)(nil)

func (e *EmptyWidget) MouseIn(ev *desktop.MouseEvent) {
	if e.OnMouseIn != nil {
		e.OnMouseIn(ev)
	}
}
func (e *EmptyWidget) MouseOut() {
	if e.OnMouseOut != nil {
		e.OnMouseOut()
	}
}
func (e *EmptyWidget) MouseMoved(_ *desktop.MouseEvent) {
	//	log.Println("implement MouseMoved")
}
func (e *EmptyWidget) FocusGained() {
	//	log.Println("focus gained")
}
func (e *EmptyWidget) FocusLost() {
	//	log.Println("focus lost")
}
func (e *EmptyWidget) TypedRune(rune) {
	//	log.Println("typed rune")
}
func (e *EmptyWidget) TypedKey(_ *fyne.KeyEvent) {
	//	log.Println("typed key")
}

// //////////////////////////////////////////////////

type emptyWidgetRenderer struct {
	objects []fyne.CanvasObject
	layout  fyne.Layout
}

func (e *EmptyWidget) CreateRenderer() fyne.WidgetRenderer {
	e.ExtendBaseWidget(e)
	var objects []fyne.CanvasObject
	objects = append(objects, e.Objects[:1]...)
	r := &emptyWidgetRenderer{objects, layout.NewMaxLayout()}
	r.applyTheme()
	return r
}
func (r *emptyWidgetRenderer) Destroy() {
}
func (r *emptyWidgetRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}
func (r *emptyWidgetRenderer) SetObjects(objects []fyne.CanvasObject) {
	r.objects = objects
}
func (r *emptyWidgetRenderer) MinSize() (size fyne.Size) {
	var w float32
	var h float32
	for _, o := range r.objects {
		s := o.MinSize()
		w += s.Width
		h += s.Height
	}
	return fyne.NewSize(w, h)
}
func (r *emptyWidgetRenderer) Layout(_ fyne.Size) {
	var inset fyne.Position
	for _, o := range r.objects {
		o.Move(inset)
		inset.Y += o.MinSize().Height
	}
}
func (r *emptyWidgetRenderer) Refresh() {
	for _, o := range r.objects {
		o.Refresh()
	}
}
func (r *emptyWidgetRenderer) applyTheme() {
}

//goland:noinspection GoUnusedExportedFunction
func NewTappableText(text *canvas.Text, onTapped func(event *fyne.PointEvent), onTappedSecondary func(*fyne.PointEvent)) *EmptyWidget {
	em := &EmptyWidget{OnTapped: onTapped, OnTappedSecondary: onTappedSecondary}
	em.Objects = make([]fyne.CanvasObject, 0)
	em.AddObject(text)
	return em
}
