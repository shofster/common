package misc

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

/*

  File:    str.go
  Author:  Bob Shofner

  Copyright (c) 2022. BSD 3-Clause License
	https://opensource.org/licenses/BSD-3-Clause

  The this permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: Various string & string list functions
*/

// StringList is the type of array
type StringList []string

// Remove deletes a string from a []string.
//goland:noinspection GoUnusedExportedFunction
func Remove(sl StringList, r string) StringList {
	for i, v := range sl {
		if v == r {
			return append(sl[:i], sl[i+1:]...)
		}
	}
	return sl
}

// Add appends a string to the end of a []string.
//goland:noinspection GoUnusedExportedFunction
func Add(sl StringList, a string, max int) StringList {
	sl = append([]string{a}, sl...)
	if len(sl) > max {
		sl = sl[0:max]
	}
	return sl
}

// Index finds the index of a string's index in an array. -1 == not found.
//goland:noinspection GoUnusedExportedFunction
func Index(sl StringList, s string) int {
	for i, v := range sl {
		if v == s {
			return i
		}
	}
	return -1
}

// Find finds a string using specified "equals" func.
//goland:noinspection GoUnusedExportedFunction
func Find(sl StringList, eq func(int, string) bool) bool {
	for i, v := range sl {
		if eq(i, v) {
			return true
		}
	}
	return false
}

// Prepend removes a string (if present) and prepends that string to the array.
//goland:noinspection GoUnusedExportedFunction
func Prepend(sl StringList, p string) StringList {
	if Index(sl, p) != -1 {
		sl = Remove(sl, p)
	}
	x := append(StringList{p}, sl...)
	return x
}

// Replace removes a string (if present) and prepends that string to the array.
// The maximum size of the list is preserved.
//goland:noinspection GoUnusedExportedFunction
func Replace(sl StringList, p string, max int) StringList {
	slr := Remove(sl, p)
	sla := Prepend(slr, p)
	if len(sla) > max {
		sla = sla[0:max]
	}
	return sla
}

// PrettyMemoryInfo generates a ?meaningful string of memory use.
//goland:noinspection GoUnusedExportedFunction
func PrettyMemoryInfo() string {
	var stats runtime.MemStats
	var mb uint64 = 1024 * 1024
	runtime.ReadMemStats(&stats)
	s := fmt.Sprintf(", Mallocs: %d, Frees: %d, Live: %d", stats.Mallocs, stats.Frees,
		stats.Mallocs-stats.Frees)
	// Sys is the total bytes of memory obtained from the OS.
	// HeapAlloc is bytes of allocated heap objects.
	s += fmt.Sprintf(", Memory- sys: %d, heap: %d (MB)", stats.Sys/mb, stats.HeapAlloc/mb)
	return s
}

// PrettyFileInfo builds a formatted string of a FileInfo type
//goland:noinspection GoUnusedExportedFunction
func PrettyFileInfo(fi os.FileInfo, dtFormat string) string {
	name := fi.Name()
	dt := fi.ModTime().Format(dtFormat)
	m := fmt.Sprintf("%s", fi.Mode())
	f := "%4s %10d %14s %s"
	s := fmt.Sprintf(f, m[0:4], fi.Size(), dt, name)
	return s
}

// PrettyUint builds a string  from a uint4. put a "sep" (probably a comma) between 3 digits.
//goland:noinspection GoUnusedExportedFunction
func PrettyUint(n uint64, sep string) string {
	if n == 0 {
		return "0"
	}
	thousands := make(StringList, 0)
	for n > 0 {
		q := n / 1000
		r := n - q*1000
		d := ""
		if q == 0 {
			d = fmt.Sprintf("%d", r)
		} else {
			d = fmt.Sprintf("%03d", r)
		}
		thousands = Prepend(thousands, d)
		n = q
	}
	s := strings.Join(thousands, sep)
	return s
}

// SizeString - insure text is exactly size characters.
//goland:noinspection GoUnusedExportedFunction
func SizeString(v string, size int) string {
	v = strings.TrimSpace(v)
	if len(v) == size {
		return v
	}
	if len(v) > size {
		return v[:size]
	}
	f := fmt.Sprintf("%%-%ds", size)
	return fmt.Sprintf(f, v)
}
