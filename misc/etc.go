package misc

import (
	"os/exec"
	"strings"
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
  Description: Various random functions.
*/

// LittleEndianToInt - build Int from slice bytes in Intel order.
//goland:noinspection GoUnusedExportedFunction
func LittleEndianToInt(bytes []uint8) int {
	s := 0
	var v int
	for i := 0; i < len(bytes); i++ {
		v |= int(bytes[i]) << s
		s += 8
	}
	return v
}

// BigEndianToInt - build Int from slice bytes in Network (Motorola) order.
//goland:noinspection GoUnusedExportedFunction
func BigEndianToInt(bytes []uint8) int {
	var v int
	for i := 0; i < len(bytes); i++ {
		v <<= 8
		v |= int(bytes[i])
	}
	return v
}

// IfElse - call yes or no function based on cond.
// warning - closure potentially can cause side effects within its scope.
//goland:noinspection GoUnusedExportedFunction
func IfElse(cond bool, yes, no func()) {
	if cond {
		yes()
	} else {
		no()
	}
}

// Camel convert to Camelcase
//goland:noinspection GoUnusedExportedFunction
func Camel(text string) string {
	text = strings.TrimSpace(text)
	t := ""
	if len(text) > 0 {
		t = strings.ToUpper(text[:1])
		if len(text) > 1 {
			t += strings.ToLower(text[1:])
		}
	}
	return t
}

// GetSystemCommand creates an executable "Cmd"
//goland:noinspection GoUnusedExportedFunction
func GetSystemCommand(cmd, path string) *exec.Cmd {
	batch := make([]string, 0)
	batch = append(batch, path)
	command := exec.Command(cmd, batch[0:]...)
	return command
}
