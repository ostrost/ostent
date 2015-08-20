package flags

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestDelayStr(t *testing.T) {
	for _, v := range []struct {
		dstr string
		cmp  string
		json []byte
	}{
		{"1h", "1h", []byte("\"1h\"")},
	} {
		td, err := time.ParseDuration(v.dstr)
		if err != nil {
			t.Error(err)
		}
		d := Delay{Duration: td}
		cmp := d.String()
		if cmp != v.cmp {
			t.Errorf("Mismatch: %q vs %q", cmp, v.cmp)
		}
		j, err := d.MarshalJSON()
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(j, v.json) {
			t.Errorf("Mismatch: %q vs %q", j, v.json)
		}
	}
}

func TestDelaySet(t *testing.T) {
	for _, v := range []struct {
		dstr  string
		above string
		str   string
		err   error
	}{
		{"abc", "", "", fmt.Errorf("time: invalid duration abc")},
		{"0.5s", "", "", fmt.Errorf("Less than a second: 500ms")},
		{"1.5s", "", "", fmt.Errorf("Not a multiple of a second: 1.5s")},
		{"1s1s", "", "2s", nil},
		{"3s", "4s", "", fmt.Errorf("Should be above 4s: 3s")},
	} {
		td, err := time.ParseDuration(v.dstr)
		if err != nil && err.Error() != v.err.Error() {
			t.Error(err)
		}
		d := Delay{Duration: td}
		if v.above != "" {
			if ad, err := time.ParseDuration(v.above); err != nil {
				t.Error(err)
			} else {
				d.Above = &ad
			}
		}
		err = d.Set(v.dstr)
		if err != nil {
			if err.Error() == v.err.Error() {
				continue
			}
			t.Error(err)
		}
		if d.String() != v.str {
			t.Errorf("Mismatch: %q vs %q", d.String(), v.str)
		}
	}
}

func TestBindSet(t *testing.T) {
	for _, v := range []struct {
		a    string
		cmp  string
		errs []error
	}{
		{"a:b:", "", []error{fmt.Errorf("too many colons in address a:b:")}},
		{"localhost:nonexistent", "", []error{
			fmt.Errorf("lookup tcp/nonexistent: Servname not supported for ai_socktype"),
			fmt.Errorf("lookup tcp/nonexistent: nodename nor servname provided, or not known"),
			fmt.Errorf("unknown port tcp/nonexistent"),
		}},
		{"localhost", "localhost:9050", nil},
		{"", ":9050", nil},
		{":8001", ":8001", nil},
		{"*:8001", "*:8001", nil},
		{"127", "127:9050", nil},
		{"127.1", "127.1:9050", nil},
		{"127.0.0.1", "127.0.0.1:9050", nil},
		{"127.1:8001", "127.1:8001", nil},
		{"127.0.0.1:8001", "127.0.0.1:8001", nil},
	} {
		b := NewBind(9050)
		if err := b.Set(v.a); err != nil {
			unknownerr := true
			for _, x := range v.errs {
				if err.Error() == x.Error() {
					unknownerr = false
					break
				}
			}
			if unknownerr {
				t.Errorf("Error: %q\nExpected errors: %+v\n", err, v.errs)
			}
			continue
		}
		if b.string != v.cmp {
			t.Errorf("Mismatch: Bind %v == %v != %v\n", v.a, v.cmp, b.string)
		}
	}
}
