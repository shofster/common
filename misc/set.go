package misc

/*

  File:    set.go
  Author:  Bob Shofner

  Copyright (c) 2022. BSD 3-Clause License
	https://opensource.org/licenses/BSD-3-Clause

  The this permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: generic set implementation. Index is insertion order (1 relqtive).
*/

type Set[T comparable] struct {
	items []T
}

//goland:noinspection GoUnusedExportedFunction
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{nil}
}

//goland:noinspection GoUnusedExportedFunction
func (s *Set[T]) Count() int {
	return len(s.items)
}

//goland:noinspection GoUnusedExportedFunction
func (s *Set[T]) Contains(t T) int {
	for ix := 0; ix < len(s.items); ix++ {
		if s.items[ix] == t {
			return ix + 1
		}
	}
	return 0
}

//goland:noinspection GoUnusedExportedFunction
func (s *Set[T]) Add(t T) {
	if s.Contains(t) == 0 {
		s.items = append(s.items, t)
	}
}

//goland:noinspection GoUnusedExportedFunction
func (s *Set[T]) Remove(ix int) {
	ix--
	if ix > -1 && ix < len(s.items) {
		s.items = append(s.items[:ix], s.items[ix+1:]...)
	}
	return
}

//goland:noinspection GoUnusedExportedFunction
func (s *Set[T]) Get(ix int) (t T, b bool) {
	ix--
	if ix > -1 && ix < len(s.items) {
		t = s.items[ix]
		b = true
	}
	return
}
