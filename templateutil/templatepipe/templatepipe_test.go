package templatepipe

import (
	"strings"
	"testing"
)

func TestDotted(t *testing.T) {
	abc := "a.b.c"
	words := strings.Split(abc, ".")
	d := Dotted{}
	d.Append(words, nil)
	d.Append(strings.Split(abc+".1", "."), nil)
	d.Append(strings.Split(abc+".2", "."), nil)
	if x := d.Find(append(words, "z")); x != nil {
		t.Errorf("Find returned non-nil")
	}
	if x := d.Find(words); x == nil {
		t.Errorf("Find returned nil")
	} else if _, _, s := x.Notation(); s != abc {
		t.Errorf("Notation: %q (expected %q)", s, abc)
	}
}

func TestEncurl(t *testing.T) {
	words := strings.Split("a.b.c", ".")
	d := Dotted{}
	d.Append(words, nil)
	l1 := d.Find(words)
	if l1 == nil {
		t.Errorf("Dotted.Find returned nil")
	}
	l1.Ranged = true
	l1.Keys = []string{"z"}
	l1.Decl = "DECL"

	w2 := strings.Split("a.b.z", ".")
	d.Append(w2, nil)
	if l2 := d.Find(w2); l2 != nil {
		l2.Ranged = true
	}

	v := Encurl(CurlyX, d, -1)
	c := v.(Nota)["a"].(Nota)["b"].(Nota)["c"].([]map[string]Nota)[0]
	if expected := "{DECL.z}"; c["z"].String() != expected {
		t.Errorf("Encurl result mismatch: %q (expected %q)", c["z"], expected)
	}
	if z := v.(Nota)["a"].(Nota)["b"].(Nota)["z"].([]string); len(z) != 0 {
		t.Errorf("z is expected to be empty: %+v\n", z)
	}
}
