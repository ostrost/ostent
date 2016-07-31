package format

import (
	"fmt"
	"math"
	"strconv"
)

func _formatOctet(sizes []string, baseInt int, n uint64) (string, string, float64, float64) { // almost humanize.IBytes
	base := float64(baseInt)
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
	return fmt.Sprintf(f+"%s", val, sizes[int(e)]), f, val, pow
}

var (
	unitlessSizes = []string{"", "k", "M", "G", "T", "P", "E"}
	bytesSizes    = []string{"B", "K", "M", "G", "T", "P", "E"}
	bitsSizes     = []string{"b", "k", "m", "g", "t", "p", "e"}
)

func HumanUnitless(n uint64) string {
	s, _, _, _ := _formatOctet(unitlessSizes, 1000, n)
	return s
}

func HumanBits(n uint64) string {
	s, _, _, _ := _formatOctet(bitsSizes, 1024, n)
	return s
}

func HumanB(n uint64) string {
	s, _, _, _ := _formatOctet(bytesSizes, 1024, n)
	return s
}

func HumanBandback(n uint64) (string, uint64, error) {
	s, f, val, pow := _formatOctet(bytesSizes, 1024, n)
	d, err := strconv.ParseFloat(fmt.Sprintf(f, val), 64)
	return s, uint64(d * pow), err
}

func Percent(used, total uint64) uint {
	if total == 0 {
		return 0
	}
	used *= 100
	pct := used / total
	if pct != 99 && used%total != 0 {
		pct++
	}
	return uint(pct)
}

func Time(T uint64) string {
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
