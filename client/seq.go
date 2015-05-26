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
func (p Param) EncodeUint(pname string, uinter enums.Uinter) DropLink {
	base := url.Values{}
	for k, v := range p.Query.Values {
		base[k] = v
	}
	text, cur := p.SetBase(base, pname, uinter)
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
func (p Param) SetBase(base url.Values, pname string, uinter enums.Uinter) (string, *bool) {
	this := uinter.Touint()
	_, low, err := uinter.Marshal()
	if err != nil { // ignoring the error
		return "", nil
	}

	text := p.Decodec.Text(strings.ToUpper(low))
	ddef := p.Decodec.Default.Uint
	dnum := p.Decoded.Number

	// Default ordering is desc (values are numeric most of the time).
	// Alpha values ordering: asc.
	desc := !p.Decodec.IsAlpha(this)
	if dnum.Negative {
		desc = !desc
	}
	var ret *bool
	if this == dnum.Uint {
		ret = new(bool)
		*ret = !desc
	}
	// for default, opposite of having a parameter is it's absence.
	if this == ddef && p.Decoded.Specified {
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

type Decodec struct {
	Default  Number
	Alphas   []enums.Uint
	Unew     func() (string, Upointer) `json:"-"`
	TextFunc func(string) string       `json:"-"`
	Texts    map[string]string         `json:"-"`
	Pname    string
}

func (d Decodec) IsAlpha(p enums.Uint) bool {
	for _, u := range d.Alphas {
		if u == p {
			return true
		}
	}
	return false
}

var Decodecs = map[string]Decodec{
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

func (p *Param) Decode(form url.Values, setn *Number) error {
	d := p.Decodec
	_, uptr := d.Unew()
	n, spec, err := p.Find(form[d.Pname], uptr)
	if err != nil {
		return err
	}
	*setn = n
	p.Decoded.Number = n
	p.Decoded.Specified = spec
	return nil
}

// Find side effects: uptr.Unmarshal and p.Query.Set (eg url.Values{}.Set())
func (p *Param) Find(values []string, uptr Upointer) (Number, bool, error) {
	d := p.Decodec
	if len(values) == 0 || values[0] == "" {
		return d.Default, false, nil
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
				p.Query.Values.Set(d.Pname, l)
			}
			p.Moved = true
		}
		return Number{}, true, err
	}
	n := Number{
		Uint:     uptr.Touint(),
		Negative: negate,
	}
	p.Query.Values.Set(d.Pname, values[0])
	return n, true, nil
}

// NewParams constructs new Params.
// Global var Decodecs is ranged.
func NewParams(req *http.Request) Params {
	if req != nil {
		req.ParseForm() // do ParseForm even if req.Form == nil
		_ = req.Form    // TODO use this
	}
	query := &Query{Values: make(url.Values)}
	p := make(Params)
	for k, v := range Decodecs {
		p[k] = &Param{
			Decodec: v,
			Query:   query,
		}
	}
	return p
}

type Param struct {
	Decodec // Read-only, an entry from global var Decodecs.
	Decoded struct {
		Number
		Specified bool
	}
	Query *Query // contains current url.Values
	Moved bool
}

// MarshalJSON goes over all defined constants
// (by the means of p.Decodec.Unew() & .Marshal method of Uinter)
// to returns a map of constants to DropLink.
func (p Param) MarshalJSON() ([]byte, error) {
	m := map[string]DropLink{}
	name, uptr := p.Decodec.Unew()
	uter := uptr.(enums.Uinter)
	marshal := uptr.Marshal
	for i := 0; i < 100; i++ {
		nextuter, s, err := marshal()
		if err != nil {
			break
		}
		m[strings.ToUpper(s)] = p.EncodeUint(name, uter)
		marshal = nextuter.Marshal
		uter = nextuter
	}
	return json.Marshal(m)
}

type Params map[string]*Param

func (p Params) Moved() bool {
	for _, v := range p {
		if v.Moved {
			return true
		}
	}
	return false
}

// Encode picks first Param and uses it's .Query.Values.
func (p Params) Encode() string {
	for _, v := range p {
		return v.Query.Values.Encode()
	}
	return ""
}

func (d Decodec) Text(in string) string {
	if s, ok := d.Texts[in]; ok {
		return s
	}
	return d.TextFunc(in)
}
