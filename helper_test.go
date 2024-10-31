package sqlhelper

import "testing"

func TestInNumbers(t *testing.T) {
	ls := []int{1, 2, 3}
	s := InString(ls)
	if s != "1, 2, 3" {
		t.Fatalf("in ints: %v", s)
	}
}

func TestInStrings(t *testing.T) {
	ls := []string{"a", "b", "c"}
	s := InString(ls)
	if s != `'a', 'b', 'c'` {
		t.Fatalf("in strings: %v", s)
	}
}
