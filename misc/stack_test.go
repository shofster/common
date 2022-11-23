package misc

import "testing"

/*

  File:    stack_test.go
  Author:  Bob Shofner

*/

func TestStack(t *testing.T) {
	s := NewStack[string]()
	s.Push("item 3")
	s.Push("item 2")
	s.Push("item 1")
	item, b := s.Peek()
	if !b || item != "item 1" {
		t.Errorf("peeK; item 1 got = %s", item)
	}
	item, b = s.Pop()
	if !b || item != "item 1" {
		t.Errorf("pop; item 1 got = %s", item)
	}
	item, b = s.Pop()
	if !b || item != "item 2" {
		t.Errorf("pop; item 2 got = %s", item)
	}
	item, b = s.Pop()
	if !b || item != "item 3" {
		t.Errorf("pop; item 3 got = %s", item)
	}
	e := s.Empty()
	if !e {
		t.Error("Stack wasn't empty")
	}
	item, b = s.Pop()
	if b || item != "" {
		t.Errorf("empty pop; got = %s", item)
	}
}
