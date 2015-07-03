package format

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

func FormatUptime(seconds float64) string { // "seconds" is expected to be sigar.Uptime.Length
	d := time.Duration(seconds) * time.Second
	s := ""
	if d > time.Duration(24)*time.Hour {
		days := d / time.Hour / 24
		end := ""
		if days > 1 {
			end = "s"
		}
		s += fmt.Sprintf("%d day%s, ", days, end)
	}
	t := time.Unix(int64(seconds), 0).UTC()
	tf := t.Format("15:04")
	if tf[0] == '0' {
		tf = " " + tf[1:]
	}
	s += tf
	return s
}

func HumanUnitless(n uint64) string {
	sizes := []string{"", "k", "M", "G", "T", "P", "E"}
	base := float64(1000)
	if float64(n) < base { // small number
		return fmt.Sprintf("%d%s", n, sizes[0])
	}
	e := math.Floor(math.Log(float64(n)) / math.Log(base))
	pow := math.Pow(base, math.Floor(e))
	val := float64(n) / pow
	f := "%.0f"
	if val < 10 {
		f = "%.1f"
	}
	return fmt.Sprintf(f+"%s", val, sizes[int(e)])
}

func _formatOctet(n uint64, bits bool) (string, string, float64, float64) { // almost humanize.IBytes
	sizes := []string{"B", "K", "M", "G", "T", "P", "E"}
	if bits { // bits instead of bytes
		sizes = []string{"b", "k", "m", "g", "t", "p", "e"}
	}
	base := float64(1024)
	if float64(n) < base { // small number
		return fmt.Sprintf("%d%s", n, sizes[0]), "%.0f", float64(n), float64(1)
	}
	e := math.Floor(math.Log(float64(n)) / math.Log(base))
	pow := math.Pow(base, math.Floor(e))
	val := float64(n) / pow
	f := "%.0f"
	if val < 10 {
		f = "%.1f"
	}
	s := fmt.Sprintf(f+"%s", val, sizes[int(e)])
	return s, f, val, pow
}

func HumanBits(n uint64) string {
	s, _, _, _ := _formatOctet(n, true)
	return s
}

func HumanB(n uint64) string {
	s, _, _, _ := _formatOctet(n, false)
	return s
}

func HumanBandback(n uint64) (string, uint64, error) {
	s, f, val, pow := _formatOctet(n, false)
	d, err := strconv.ParseFloat(fmt.Sprintf(f, val), 64)
	return s, uint64(d * pow), err
}

func Percent(used, total uint64) uint {
	if total == 0 {
		return 0
	}
	used *= 100
	pct := uint64(used / total)
	if pct != 99 && used%total != 0 {
		pct++
	}
	return uint(pct)
}

func FormatPercent(used, total uint64) string {
	return strconv.Itoa(int(Percent(used, total))) // without "%"
}

func FormatTime(T uint64) string {
	// 	ms := T % 60
	t := T / 1000
	ss := t % 60
	t /= 60 // fst t shift
	mm := t % 60
	t /= 60 // snd t shift
	hh := t % 24
	if hh > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hh, mm, ss)
	}
	return fmt.Sprintf("   %02d:%02d", mm, ss)
}

/* unused
func Bps(factor int, current, previous uint64) string {
	if current < previous { // counters got reset
		return ""
	}
	diff := (current - previous) * uint64(factor) // bits now if the factor is 8
	return HumanBits(diff)
}

func Ps(current, previous uint64) string {
	if current < previous { // counters got reset
		return ""
	}
	return HumanUnitless(current - previous)
}
// */
