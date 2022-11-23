package misc

import (
	"testing"
)

/*

  File:    set_test.go
  Author:  Bob Shofner

*/

func TestSet(t *testing.T) {
	set := NewSet[string]()
	set.Add("item 1")
	set.Add("item 2")
	set.Add("item 3")
	set.Add("item 1")
	c := set.Count()
	if c != 3 {
		t.Errorf("Add Count = %d; want 3", c)
	}
	ix := set.Contains("item 1")
	if ix != 1 {
		t.Errorf("index = %d; want 1", ix)
	}
	ix = set.Contains("item 3")
	if ix != 3 {
		t.Errorf("index = %d; want 3", ix)
	}
	set.Remove(ix)
	c = set.Count()
	if c != 2 {
		t.Errorf("Remove Count = %d; want 2", c)
	}
	ix = set.Contains("item 1")
	if ix != 1 {
		t.Errorf("item 1 index = %d; want 1", ix)
	}
	ix = set.Contains("item 2")
	if ix != 2 {
		t.Errorf("item 2 index = %d; want 2", ix)
	}
	item, b := set.Get(1)
	if !b || item != "item 1" {
		t.Error("item 1 missing")
	}
}
