package misc

/*

  File:    stack.go
  Author:  Bob Shofner

  Copyright (c) 2022. BSD 3-Clause License
	https://opensource.org/licenses/BSD-3-Clause

  The this permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: generic stack implementation. LIFO
*/

type Stack[T any] struct {
	items []T
}

//goland:noinspection GoUnusedExportedFunction
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{nil}
}

//goland:noinspection GoUnusedExportedFunction
func (s *Stack[T]) Empty() bool {
	return len(s.items) == 0
}

//goland:noinspection GoUnusedExportedFunction
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

//goland:noinspection GoUnusedExportedFunction
func (s *Stack[T]) Peek() (T, bool) {
	var item T
	if len(s.items) > 0 {
		return s.items[len(s.items)-1], true
	}
	return item, false
}

//goland:noinspection GoUnusedExportedFunction
func (s *Stack[T]) Pop() (T, bool) {
	var item T
	if len(s.items) > 0 {
		item = s.items[len(s.items)-1]
		s.items = s.items[:len(s.items)-1]
		return item, true
	}
	return item, false
}
