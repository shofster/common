package element

import "fyne.io/fyne/v2/data/binding"

/*

  File:    bind.go
  Author:  Bob Shofner

  Copyright (c) 2022. BSD 3-Clause License
	https://opensource.org/licenses/BSD-3-Clause

  The this permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: Binding getter and setter that drop error.
*/

// StringGetter gets bound string. swallows error.
//goland:noinspection GoUnusedExportedFunction
func StringGetter(bs binding.String) string {
	s, _ := bs.Get()
	return s
}

// StringSetter sets bound string. swallows error.
//goland:noinspection GoUnusedExportedFunction
func StringSetter(s string, bs binding.String) {
	_ = bs.Set(s)
}

// BoolGetter gets bound bool. swallows error.
//goland:noinspection GoUnusedExportedFunction
func BoolGetter(bb binding.Bool) bool {
	s, _ := bb.Get()
	return s
}

// BoolSetter gets bound bool. swallows error.
//goland:noinspection GoUnusedExportedFunction
func BoolSetter(s bool, bb binding.Bool) {
	_ = bb.Set(s)
}

// IntGetter gets bound int. swallows error.
//goland:noinspection GoUnusedExportedFunction
func IntGetter(bi binding.Int) int {
	s, _ := bi.Get()
	return s
}

// IntSetter gets bound int. swallows error.
//goland:noinspection GoUnusedExportedFunction
func IntSetter(s int, bi binding.Int) {
	_ = bi.Set(s)
}
