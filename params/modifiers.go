package params

import (
	"fmt"
	"math"
)

// LessN is in the func map.
func LessN(p *Params, num *Num) (ALink, error) {
	return LinkN(p, num, Pow2Less(num.Absolute), false)
}

// MoreN is in the func map.
func MoreN(p *Params, num *Num) (ALink, error) {
	return LinkN(p, num, Pow2More(num.Absolute), true)
}

// Vlink is in the func map.
func Vlink(p *Params, num *Num, absolute int, text string) (VLink, error) {
	vl := VLink{LinkClass: "state"}
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
func LinkN(p *Params, num *Num, absolute int, more bool) (ALink, error) {
	al := ALink{ // defaults
		Href:       "?",          // Default
		ExtraClass: " disabled ", //         Disabled

		Text: fmt.Sprintf("%d", absolute), // Final
	}
	href, err := p.EncodeN(num, absolute, nil)
	if err != nil {
		return al, err
	}
	if more {
		// when num.Limit is 0, it's unknown, so enable the button
		if num.Limit == 0 || num.Absolute < num.Limit || absolute <= num.Limit {
			al.Href, al.ExtraClass = href, ""
		}
	} else {
		if absolute > 0 || num.Absolute > 0 {
			al.Href, al.ExtraClass = href, ""
		}
	}
	return al, nil
}
