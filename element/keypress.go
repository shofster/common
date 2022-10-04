package element

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"sync"
	"time"
)

/*

  File:    keypress.go
  Author:  Bob Shofner

  Copyright (c) 2022. BSD 3-Clause License
	https://opensource.org/licenses/BSD-3-Clause

  The this permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:
*/

const shift = 1
const control = 2
const alt = 4
const cmd = 8

type Keypress struct {
	OnKeyUp  func(k KeyPressed)
	modifier uint8
	pressed  int64
	long     int64
	mu       sync.Mutex
}
type KeyPressed struct {
	Name     fyne.KeyName
	modifier uint8
}

//goland:noinspection GoUnusedExportedFunction
func NewKeyPress(long int64, onKeyUp func(k KeyPressed)) *Keypress {
	k := &Keypress{long: long, OnKeyUp: onKeyUp}
	go func() {
		for {
			time.Sleep(time.Duration(long) * time.Millisecond)
			k.mu.Lock()
			// clear modifier if no regular key pressed in time
			if k.modifier != 0 {
				now := time.Now().UnixMilli()
				d := now - k.pressed
				if d > k.long {
					// fmt.Println("cleared then", k.pressed, "now", now, " diff", d)
					k.modifier = 0
					k.pressed = 0
				}
			}
			k.mu.Unlock()
		}
	}()
	return k
}
func (k *Keypress) PressedKey(name fyne.KeyName) {
	var start = func(b uint8) {
		k.mu.Lock()
		defer k.mu.Unlock()
		k.modifier |= b
		k.pressed = time.Now().UnixMilli()
		// fmt.Println("set bit", b, "on", name, " now", k.pressed)
	}
	switch name {
	case desktop.KeyShiftLeft, desktop.KeyShiftRight:
		start(shift)
	case desktop.KeyControlLeft, desktop.KeyControlRight:
		start(control)
	case desktop.KeyAltLeft, desktop.KeyAltRight:
		start(alt)
	case desktop.KeySuperLeft, desktop.KeySuperRight:
		start(cmd)
	default:
		k.OnKeyUp(KeyPressed{Name: name, modifier: k.modifier})
		k.modifier = 0
	}
}
func (k KeyPressed) IsModified() bool {
	return k.modifier != 0
}
func (k KeyPressed) IsShift() bool {
	return (k.modifier & shift) != 0
}
func (k KeyPressed) IsControl() bool {
	return (k.modifier & control) != 0
}
func (k KeyPressed) IsAlt() bool {
	return (k.modifier & alt) != 0
}
func (k KeyPressed) IsCmd() bool {
	return (k.modifier & cmd) != 0
}
