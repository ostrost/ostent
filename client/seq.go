package client

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/ostrost/ostent/client/enums"
)

// Number is an enums.Uint with sign.
type Number struct {
	enums.Uint
	Negative bool
}

// Upointer defines required (incl. pointer-) methods
// for all enums.Uint-derived types interface.
type Upointer interface {
	enums.Uinter
	Unmarshal(string, *bool) error
	// UnmarshalJSON([]byte) error
}

// DropLink has drop{down,up} link attributes.
type DropLink struct {
	CLASSNAME  string // required for jsx
	AlignClass string
	Text       string
	Href       string
	Class      string
	CaretClass string
}

// EncodeUint returns enums.uinter applied DropLink. .AlignClass is not filled.
func (ep EnumParam) EncodeUint(pname string, uinter enums.Uinter) DropLink {
	base := url.Values{}
	for k, v := range ep.Query.Values {
		base[k] = v
	}
	text, cur := ep.SetBase(base, pname, uinter)
	dl := DropLink{Text: text, Class: "state"}
	if cur != nil {
		dl.CaretClass = "caret"
		dl.Class += " current"
		if *cur {
			dl.Class += " dropup"
		}
	}
	dl.Href = "?" + base.Encode() // sorted by key
	return dl
}

// SetBase modifies the base.
func (ep EnumParam) SetBase(base url.Values, pname string, uinter enums.Uinter) (string, *bool) {
	this := uinter.Touint()
	_, low, err := uinter.Marshal()
	if err != nil { // ignoring the error
		return "", nil
	}

	text := ep.EnumDecodec.Text(strings.ToUpper(low))
	ddef := ep.EnumDecodec.Default.Uint
	dnum := ep.Decoded.Number

	// Default ordering is desc (values are numeric most of the time).
	// Alpha values ordering: asc.
	desc := !ep.IsAlpha(this)
	if dnum.Negative {
		desc = !desc
	}
	var ret *bool
	if this == dnum.Uint {
		ret = new(bool)
		*ret = !desc
	}
	// for default, opposite of having a parameter is it's absence.
	if this == ddef && ep.Decoded.Specified {
		base.Del(pname)
		return text, ret
	}
	if this == dnum.Uint && !dnum.Negative {
		low = "-" + low
	}
	base.Set(pname, low)
	return text, ret
}

// Query type for link making.
type Query struct {
	Values url.Values
	Moved  bool
}

type EnumDecodec struct {
	Default  Number
	Alphas   []enums.Uint
	Unew     func() (string, Upointer) `json:"-"`
	TextFunc func(string) string       `json:"-"`
	Texts    map[string]string         `json:"-"`
	Pname    string
}

func (ec EnumParam) IsAlpha(p enums.Uint) bool {
	for _, u := range ec.EnumDecodec.Alphas {
		if u == p {
			return true
		}
	}
	return false
}

var EnumDecodecs = map[string]EnumDecodec{
	"ps": {
		Default: Number{Uint: enums.Uint(enums.PID)},
		Alphas:  []enums.Uint{enums.Uint(enums.NAME), enums.Uint(enums.USER)},
		// Unew:  func() (string, Upointer) { return "ps", interface{}(new(enums.UintPS)).(Upointer) },
		Unew:     func() (string, Upointer) { return "ps", new(enums.UintPS) },
		TextFunc: strings.ToUpper,
		Texts:    map[string]string{"PRI": "PR", "NICE": "NI", "NAME": "COMMAND"},
		Pname:    "ps",
	},
	"df": {
		Default: Number{Uint: enums.Uint(enums.FS)},
		Alphas:  []enums.Uint{enums.Uint(enums.FS), enums.Uint(enums.MP)},
		// Unew:  func() (string, Upointer) { return "df", interface{}(new(enums.UintDF)).(Upointer) },
		Unew:     func() (string, Upointer) { return "df", new(enums.UintDF) },
		TextFunc: func(s string) string { return strings.Title(strings.ToLower(s)) },
		Texts:    map[string]string{"FS": "Device", "MP": "Mounted"},
		Pname:    "df",
	},
}

func (ep *EnumParam) Decode(form url.Values) error {
	_, uptr := ep.EnumDecodec.Unew()
	n, spec, err := ep.Find(form[ep.EnumDecodec.Pname], uptr)
	if err != nil {
		return err
	}
	ep.Decoded.Number = n
	ep.Decoded.Specified = spec
	return nil
}

// Find side effects: uptr.Unmarshal and p.Query.Set (eg url.Values{}.Set())
func (ep *EnumParam) Find(values []string, uptr Upointer) (Number, bool, error) {
	if len(values) == 0 || values[0] == "" {
		return ep.EnumDecodec.Default, false, nil
	}
	var negate bool
	in := values[0]
	if in[0] == '-' {
		in = in[1:]
		negate = true
	}
	err := uptr.Unmarshal(in, &negate) // .UnmarshalJSON([]byte(fmt.Sprintf("%q", strings.ToUpper(in))))
	if err != nil {
		if _, ok := err.(enums.RenamedConstError); ok {
			// The case when err (of type RenamedConstError) is set
			// AND uptr actually holds corresponding ("renamed") value.
			if _, l, err := uptr.Marshal(); err == nil {
				if negate {
					l = "-" + l
				}
				ep.Query.Values.Set(ep.EnumDecodec.Pname, l)
			}
			ep.Moved = true
		}
		return Number{}, true, err
	}
	n := Number{
		Uint:     uptr.Touint(),
		Negative: negate,
	}
	ep.Query.Values.Set(ep.EnumDecodec.Pname, values[0])
	return n, true, nil
}

// NewParams constructs new Params.
// Global var Decodecs is ranged.
func NewParams(req *http.Request) Params {
	if req != nil {
		req.ParseForm() // do ParseForm even if req.Form == nil
		_ = req.Form    // TODO use this
	}
	q := &Query{Values: make(url.Values)}
	enum := make(Enums)
	for k, v := range EnumDecodecs {
		enum[k] = &EnumParam{
			EnumDecodec: v,
			Query:       q,
		}
	}
	return Params{
		ENUM: enum,
	}
}

type EnumParam struct {
	// EnumDecodec is read-only, an entry from global var Decodecs.
	EnumDecodec
	Decoded struct {
		Number
		Specified bool
	}
	Query *Query // contains current url.Values
	Moved bool
}

func (ec EnumParam) LessorMore(r bool) bool {
	// numeric values: flip r
	if !ec.IsAlpha(ec.Decoded.Number.Uint) {
		r = !r
	}
	if ec.Decoded.Number.Negative {
		r = !r
	}
	return r
}

// MarshalJSON goes over all defined constants
// (by the means of p.EnumDecodec.Unew() & .Marshal method of Uinter)
// to returns a map of constants to DropLink.
func (ep EnumParam) MarshalJSON() ([]byte, error) {
	m := map[string]DropLink{}
	name, uptr := ep.EnumDecodec.Unew()
	uter := uptr.(enums.Uinter)
	marshal := uptr.Marshal
	for i := 0; i < 100; i++ {
		nextuter, s, err := marshal()
		if err != nil {
			break
		}
		m[strings.ToUpper(s)] = ep.EncodeUint(name, uter)
		marshal = nextuter.Marshal
		uter = nextuter
	}
	return json.Marshal(m)
}

type Enums map[string]*EnumParam

type Params struct {
	ENUM Enums
}

func (ps Params) Moved() bool {
	for _, v := range ps.ENUM {
		if v.Moved {
			return true
		}
	}
	return false
}

// Encode picks first EnumParam and uses it's .Query.Values.
func (ps Params) Encode() string {
	for _, v := range ps.ENUM {
		return v.Query.Values.Encode()
	}
	return ""
}

func (ed EnumDecodec) Text(in string) string {
	if s, ok := ed.Texts[in]; ok {
		return s
	}
	return ed.TextFunc(in)
}
