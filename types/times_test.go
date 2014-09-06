package types

import (
	"bytes"
	"errors"
	"testing"
	"time"
)

func TestDuration(t *testing.T) {
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
		d := Duration(td)
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

func TestPeriodValue(t *testing.T) {
	for _, v := range []struct {
		dstr  string
		above string
		str   string
		err   error
	}{
		{"abc", "", "", errors.New("time: invalid duration abc")},
		{"0.5s", "", "", errors.New("Less than a second: 500ms")},
		{"1.5s", "", "", errors.New("Not a multiple of a second: 1.5s")},
		{"1s1s", "", "2s", nil},
		{"3s", "4s", "", errors.New("Should be above 4s: 3s")},
	} {
		td, err := time.ParseDuration(v.dstr)
		if err != nil && err.Error() != v.err.Error() {
			t.Error(err)
		}
		p := PeriodValue{Duration: Duration(td)}
		if v.above != "" {
			if ad, err := time.ParseDuration(v.above); err != nil {
				t.Error(err)
			} else {
				ap := Duration(ad)
				p.Above = &ap
			}
		}
		err = p.Set(v.dstr)
		if err != nil {
			if err.Error() == v.err.Error() {
				continue
			}
			t.Error(err)
		}
		if p.String() != v.str {
			t.Errorf("Mismatch: %q vs %q", p.String(), v.str)
		}
	}
}
