package templatepipe

import (
	"strings"
	"testing"
)

func TestDotted(t *testing.T) {
	abc := "a.b.c"
	words := strings.Split(abc, ".")
	d := dotted{}
	d.append(words, nil)
	d.append(strings.Split(abc+".1", "."), nil)
	d.append(strings.Split(abc+".2", "."), nil)
	if x := d.find(append(words, "z")); x != nil {
		t.Errorf("find returned non-nil")
	}
	if x := d.find(words); x == nil {
		t.Errorf("find returned nil")
	} else if _, _, s := x.notation(); s != abc {
		t.Errorf("notation: %q (expected %q)", s, abc)
	}
}

func TestEncurl(t *testing.T) {
	words := strings.Split("a.b.c", ".")
	d := dotted{}
	d.append(words, nil)
	l1 := d.find(words)
	if l1 == nil {
		t.Errorf("dotted.find returned nil")
	}
	l1.ranged = true
	l1.keys = []string{"z"}
	l1.decl = "DECL"

	w2 := strings.Split("a.b.z", ".")
	d.append(w2, nil)
	if l2 := d.find(w2); l2 != nil {
		l2.ranged = true
	}

	v := encurl(curly, d, -1)
	c := v.(Nota)["a"].(Nota)["b"].(Nota)["c"].([]map[string]Nota)[0]
	if expected := "{DECL.z}"; c["z"].String() != expected {
		t.Errorf("encurl result mismatch: %q (expected %q)", c["z"], expected)
	}
	if z := v.(Nota)["a"].(Nota)["b"].(Nota)["z"].([]string); len(z) != 0 {
		t.Errorf("z is expected to be empty: %+v\n", z)
	}
}
