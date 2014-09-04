package ostent

import (
	"fmt"
	"testing"
	// "strconv"
	// "github.com/dustin/go-humanize"
	sigar "github.com/rzab/gosigar"
)

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

func Test_humanB(t *testing.T) {
	for _, v := range []struct {
		a    uint64
		cmp  string
		back uint64
	}{
		{1023, "1023B", 1023},
		{1024, "1.0K", 1024},
		{117649480 * 1024, "112G", 120259084288},
	} {
		cmp := humanB(v.a)
		if cmp[0] == ' ' {
			t.Errorf("Unexpected: starts with a space: %q", cmp)
		}
		if cmp != v.cmp {
			t.Errorf("Mismatch: humanB(%v) == %v != %v\n", v.a, v.cmp, cmp)
		}
		bcmp, back, err := humanBandback(v.a)
		if err != nil {
			t.Error(err)
		}
		if bcmp[0] == ' ' {
			t.Errorf("Unexpected: starts with a space: %q", bcmp)
		}
		if bcmp != v.cmp {
			t.Errorf("Mismatch: humanBandback(%v) == %v != %v\n", v.a, v.cmp, bcmp)
		}
		if back != v.back {
			t.Errorf("Mismatch: humanBandback(%v) == %v != %v\n", v.a, v.back, back)
		}
	}
}

func Test_humanbits(t *testing.T) {
	for _, v := range []struct {
		a   uint64
		cmp string
	}{
		{1023, "1023b"},
		{1024, "1.0k"},
	} {
		cmp := humanbits(v.a)
		if cmp[0] == ' ' {
			t.Errorf("Unexpected: starts with a space: %q", cmp)
		}
		if cmp != v.cmp {
			t.Errorf("Mismatch: humanB(%v) == %v != %v\n", v.a, v.cmp, cmp)
		}
	}
}

func Test_humanUnitless(t *testing.T) {
	for _, v := range []struct {
		a   uint64
		cmp string
	}{
		{999, "999"},
		{1000, "1.0k"},
		{1001, "1.0k"},
		{1050, "1.1k"},
	} {
		cmp := humanUnitless(v.a)
		if cmp[0] == ' ' {
			t.Errorf("Unexpected: starts with a space: %q", cmp)
		}
		if cmp != v.cmp {
			t.Errorf("Mismatch: humanUnitless(%v) == %v != %v\n", v.a, v.cmp, cmp)
		}
	}
}

func Test_percent(t *testing.T) {
	for _, v := range []struct {
		a, b uint64
		cmp  uint
		fcmp string
	}{
		{1, 0, 0, "0"},
		{201, 1000, 21, "21"},
		{800, 1000, 80, "80"},
		{890, 1000, 89, "89"},
		{891, 1000, 90, "90"},
		{899, 1000, 90, "90"},
		{900, 1000, 90, "90"},
		{901, 1000, 91, "91"},
		{990, 1000, 99, "99"},
		{991, 1000, 99, "99"},
		{995, 1000, 99, "99"},
		{996, 1000, 99, "99"},
		{999, 1000, 99, "99"},
		{1000, 1000, 100, "100"},
	} {
		cmp := percent(v.a, v.b)
		if cmp != v.cmp {
			t.Errorf("Mismatch: percent(%v, %v) == %v != %v\n", v.a, v.b, v.cmp, cmp)
		}
		fcmp := formatPercent(v.a, v.b)
		if fcmp != v.fcmp {
			t.Errorf("Mismatch: formatPercent(%v, %v) == %v != %v\n", v.a, v.b, v.fcmp, fcmp)
		}
	}
}

func Test_formatTime(t *testing.T) {
	for _, v := range []struct {
		a   int
		cmp string
	}{
		{1000 * 62, "   01:02"},
		{1000 * 60 * 60, "01:00:00"},
	} {
		cmp := formatTime(uint64(v.a))
		if cmp != v.cmp {
			t.Errorf("Mismatch: formatTime(%v) == %v != %v\n", v.a, v.cmp, cmp)
		}
	}
}

func Test_formatUptime(t *testing.T) {
	for _, v := range []struct {
		a   float64
		cmp string
	}{
		{1080720, "12 days, 12:12"},
		{1069920, "12 days,  9:12"},
		{43920, "12:12"},
		{33120, " 9:12"},
	} {
		cmp := formatUptime(v.a)
		if cmp != v.cmp {
			t.Errorf("Mismatch: formatUptime(%v) == %v != %v\n", v.a, v.cmp, cmp)
		}
	}
}

func sigarUptime(t *testing.B) *sigar.Uptime {
	return &sigar.Uptime{Length: 1080720 + float64(t.N)}
}

func BenchmarkUptimeFormat(t *testing.B)      { sigarUptime(t).Format() }
func BenchmarkFormatUptime(t *testing.B)      { formatUptime((*sigarUptime(t)).Length) }
func BenchmarkSigarUptimeFormat(t *testing.B) { sigarUptimeFormatString(*sigarUptime(t)) }

// the way sigar.Uptime.Format implemented
// sans bytes.Buffer, bufio.NewWriter stuff
func sigarUptimeFormatString(u sigar.Uptime) string {
	uptime := uint64(u.Length)
	days := uptime / (60 * 60 * 24)

	s := ""
	if days != 0 {
		end := ""
		if days > 1 {
			end = "s"
		}
		s = fmt.Sprintf("%d day%s, ", days, end)
	}

	minutes := uptime / 60
	hours := minutes / 60
	hours %= 24
	minutes %= 60

	s += fmt.Sprintf("%2d:%02d", hours, minutes)
	return s
}
