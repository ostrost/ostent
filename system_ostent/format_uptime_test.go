package system_ostent

import "testing"

func TestFormatUptime(t *testing.T) {
	for _, v := range []struct {
		a   uint64
		cmp string
	}{
		{1080720, "12 days, 12:12"},
		{1069920, "12 days,  9:12"},
		{43920, "12:12"},
		{33120, " 9:12"},
	} {
		if cmp := format_uptime(v.a); cmp != v.cmp {
			t.Errorf("Mismatch: format_uptime(%v) == %v != %v\n", v.a, v.cmp, cmp)
		}
	}
}
