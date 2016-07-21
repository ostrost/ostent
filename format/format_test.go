package format

import "testing"

/* func Test_humanizeParseBytes(t *testing.T) {
	_, err := humanize.ParseBytes("70GiB")
	if err != nil {
		t.Error(err)
	}
	cmp := strconv.FormatFloat(float64(10.6), 'f', 0, 64)
	if cmp != "11" {
		t.Errorf("Mismatch, cmp: \"%v\"\n", cmp)
	}
} // */

func Test_HumanB(t *testing.T) {
	for _, v := range []struct {
		a    uint64
		cmp  string
		back uint64
	}{
		{1023, "1023B", 1023},
		{1024, "1.0K", 1024},
		{117649480 * 1024, "112G", 120259084288},
	} {
		cmp := HumanB(v.a)
		if cmp[0] == ' ' {
			t.Errorf("Unexpected: starts with a space: %q", cmp)
		}
		if cmp != v.cmp {
			t.Errorf("Mismatch: HumanB(%v) == %v != %v\n", v.a, v.cmp, cmp)
		}
		bcmp, back, err := HumanBandback(v.a)
		if err != nil {
			t.Error(err)
		}
		if bcmp[0] == ' ' {
			t.Errorf("Unexpected: starts with a space: %q", bcmp)
		}
		if bcmp != v.cmp {
			t.Errorf("Mismatch: HumanBandback(%v) == %v != %v\n", v.a, v.cmp, bcmp)
		}
		if back != v.back {
			t.Errorf("Mismatch: HumanBandback(%v) == %v != %v\n", v.a, v.back, back)
		}
	}
}

func Test_HumanBits(t *testing.T) {
	for _, v := range []struct {
		a   uint64
		cmp string
	}{
		{1023, "1023b"},
		{1024, "1.0k"},
	} {
		cmp := HumanBits(v.a)
		if cmp[0] == ' ' {
			t.Errorf("Unexpected: starts with a space: %q", cmp)
		}
		if cmp != v.cmp {
			t.Errorf("Mismatch: HumanBits(%v) == %v != %v\n", v.a, v.cmp, cmp)
		}
	}
}

func Test_HumanUnitless(t *testing.T) {
	for _, v := range []struct {
		a   uint64
		cmp string
	}{
		{999, "999"},
		{1000, "1.0k"},
		{1001, "1.0k"},
		{1050, "1.1k"},
	} {
		cmp := HumanUnitless(v.a)
		if cmp[0] == ' ' {
			t.Errorf("Unexpected: starts with a space: %q", cmp)
		}
		if cmp != v.cmp {
			t.Errorf("Mismatch: HumanUnitless(%v) == %v != %v\n", v.a, v.cmp, cmp)
		}
	}
}

func Test_Percent(t *testing.T) {
	for _, v := range []struct {
		a, b uint64
		cmp  uint
	}{
		{1, 0, 0},
		{201, 1000, 21},
		{800, 1000, 80},
		{890, 1000, 89},
		{891, 1000, 90},
		{899, 1000, 90},
		{900, 1000, 90},
		{901, 1000, 91},
		{990, 1000, 99},
		{991, 1000, 99},
		{995, 1000, 99},
		{996, 1000, 99},
		{999, 1000, 99},
		{1000, 1000, 100},
	} {
		if cmp := Percent(v.a, v.b); cmp != v.cmp {
			t.Errorf("Mismatch: Percent(%v, %v) == %v != %v\n", v.a, v.b, v.cmp, cmp)
		}
	}
}

func TestTime(t *testing.T) {
	for _, v := range []struct {
		a   int
		cmp string
	}{
		{1000 * 62, "   01:02"},
		{1000 * 60 * 60, "01:00:00"},
	} {
		cmp := Time(uint64(v.a))
		if cmp != v.cmp {
			t.Errorf("Mismatch: Time(%v) == %v != %v\n", v.a, v.cmp, cmp)
		}
	}
}

func TestUptime(t *testing.T) {
	for _, v := range []struct {
		a   uint64
		cmp string
	}{
		{1080720, "12 days, 12:12"},
		{1069920, "12 days,  9:12"},
		{43920, "12:12"},
		{33120, " 9:12"},
	} {
		if cmp := Uptime(v.a); cmp != v.cmp {
			t.Errorf("Mismatch: Uptime(%v) == %v != %v\n", v.a, v.cmp, cmp)
		}
	}
}
