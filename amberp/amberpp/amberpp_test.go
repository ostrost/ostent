package main

import (
	"strings"
	"testing"
)

func TestDotted(t *testing.T) {
	abc := "a.b.c"
	words := strings.Split(abc, ".")
	d := Dotted{}
	d.Append(words)
	d.Append(strings.Split(abc+".1", "."))
	d.Append(strings.Split(abc+".2", "."))
	if x := d.Find(append(words, "z")); x != nil {
		t.Errorf("Find returned non-nil")
	}
	if x := d.Find(words); x == nil {
		t.Errorf("Find returned nil")
	} else if s := x.Notation(); s != abc {
		t.Errorf("Notation: %q (expected %q)", s, abc)
	}
	if expected, s := `[]
  [a]
    [b]
      [c]
        [1]
        [2]
`, d.GoString(); s != expected {
		t.Errorf("GoString: %s (expected %s)", s, expected)
	}
}
