package params

import (
	"fmt"
	"html/template"
	"math"
	"time"

	"github.com/ostrost/ostent/flags"
)

// FuncMapHTML is a FuncMap with param functions.
var FuncMapHTML = template.FuncMap{
	"HrefT": HrefT, // HrefT not used in params package
	"LessD": LessD,
	"MoreD": MoreD,
	"LessN": LessN,
	"MoreN": MoreN,
	"Vlink": Vlink,
}

// HrefT is in the func map.
func HrefT(p *Params, num *Num) (template.HTMLAttr, error) {
	href, err := p.EncodeT(num)
	return template.HTMLAttr(fmt.Sprintf(" href=%q", href)), err
}

// LessD is in the func map.
func LessD(p *Params, d *Delay, bclass string) (ALink, error) {
	return LinkD(p, d, bclass, DelayLess(*d, p.DelayBounds.Min.Duration), "-")
}

// MoreD is in the func map.
func MoreD(p *Params, d *Delay, bclass string) (ALink, error) {
	return LinkD(p, d, bclass, DelayMore(*d, p.DelayBounds.Min.Duration), "+")
}

// LessN is in the func map.
func LessN(p *Params, num *Num, bclass string) (ALink, error) {
	return LinkN(p, num, bclass, Pow2Less(num.Absolute), "-")
}

// MoreN is in the func map.
func MoreN(p *Params, num *Num, bclass string) (ALink, error) {
	return LinkN(p, num, bclass, Pow2More(num.Absolute), "+")
}

// Vlink is in the func map.
func Vlink(p *Params, num *Num, absolute int, text string) (VLink, error) {
	vl := VLink{LinkText: text, LinkClass: "state"}
	negative := new(bool) // EncodeN will use .Negative being false by default
	if num.Absolute == absolute {
		vl.CaretClass = "caret"
		vl.LinkClass += " current"
		if (num.Alpha && !num.Negative) || (!num.Alpha && num.Negative) {
			vl.LinkClass += " dropup"
		}
		*negative = !num.Negative
	}
	qs, err := p.EncodeN(num, absolute, negative)
	if err != nil {
		return VLink{}, err
	}
	vl.LinkHref = qs
	return vl, nil
}

// DelayMore is internal.
func DelayMore(d Delay, step time.Duration) time.Duration {
	const s = time.Second
	const m = time.Second * 60
	var table = map[time.Duration]time.Duration{
		s:      2 * s,
		2 * s:  5 * s,
		5 * s:  10 * s,
		10 * s: 30 * s,
		30 * s: m,
		m:      2 * m,
		2 * m:  5 * m,
		5 * m:  10 * m,
		10 * m: 30 * m,
		30 * m: 60 * m,
	}
	if more, ok := table[d.D]; ok {
		return more
	}
	return d.D + step
}

// DelayLess is internal.
func DelayLess(d Delay, step time.Duration) time.Duration {
	const s = time.Second
	const m = time.Second * 60
	var table = map[time.Duration]time.Duration{
		s:      s,
		2 * s:  s,
		5 * s:  2 * s,
		10 * s: 5 * s,
		30 * s: 10 * s,
		m:      30 * s,
		2 * m:  m,
		5 * m:  2 * m,
		10 * m: 5 * m,
		30 * m: 10 * m,
		60 * m: 30 * m,
	}
	if less, ok := table[d.D]; ok {
		return less
	}
	return d.D - step
}

// LinkD is internal.
func LinkD(p *Params, d *Delay, bclass string, set time.Duration, badge string) (ALink, error) {
	al := ALink{
		Href:       "?",                   // Default
		ExtraClass: " disabled ",          //         Disabled
		Class:      " disabled " + bclass, //                  Values

		Text:  flags.DurationString(set), // Final
		Badge: badge,                     //       Values
	}
	href, err := p.EncodeD(d, set)
	if err != nil {
		return al, err
	}
	switch badge {
	case "-":
		if d.D > p.DelayBounds.Min.Duration {
			al.Href, al.ExtraClass = href, ""
		}
	case "+":
		if d.D < p.DelayBounds.Max.Duration {
			al.Href, al.ExtraClass = href, ""
		}
	}
	al.Class = al.ExtraClass + " " + bclass // Eventually
	return al, nil
}

// Pow2Less is internal.
func Pow2Less(v int) int {
	switch v {
	case 0:
		return 0
	case 1:
		return 0
	case 2:
		return 1
	}
	g := math.Log2(float64(v))
	n := math.Floor(g)
	if n == g {
		n--
	}
	return int(math.Pow(2, n))
}

// Pow2More is internal.
func Pow2More(v int) int {
	switch v {
	case 0:
		return 1
	case 1:
		return 2
	case 2:
		return 4
	}
	if v <= 32768 { // up to 65536
		v = int(math.Pow(2, 1+math.Floor(math.Log2(float64(v)))))
	}
	return v
}

// LinkN is internal.
func LinkN(p *Params, num *Num, bclass string, absolute int, badge string) (ALink, error) {
	al := ALink{ // defaults
		Href:       "?",                   // Default
		ExtraClass: " disabled ",          //         Disabled
		Class:      " disabled " + bclass, //                  Values

		Text:  fmt.Sprintf("%d", absolute), // Final
		Badge: badge,                       //       Values
	}
	href, err := p.EncodeN(num, absolute, nil)
	if err != nil {
		return al, err
	}
	switch badge {
	case "+":
		// when num.Limit is 0, it's unknown, so enable the button
		if num.Limit == 0 || num.Absolute < num.Limit || absolute <= num.Limit {
			al.Href, al.ExtraClass = href, ""
		}
	case "-":
		if absolute > 0 || num.Absolute > 0 {
			al.Href, al.ExtraClass = href, ""
		}
	}
	al.Class = al.ExtraClass + " " + bclass // Eventually
	return al, nil
}
