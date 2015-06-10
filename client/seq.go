package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/url"
	"sort"
	"strings"

	"github.com/ostrost/ostent/client/enums"
	"github.com/ostrost/ostent/flags"
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
	AlignClass string
	Text       string `json:"-"` // static
	Href       string
	Class      string
	CaretClass string
}

// EncodeUint returns enums.uinter applied DropLink. .AlignClass is not filled.
func (ep EnumParam) EncodeUint(pname string, uinter enums.Uinter) DropLink {
	values := ep.Query.ValuesCopy()
	text, cur := ep.SetValue(values, pname, uinter)
	dl := DropLink{Text: text, Class: "state"}
	if cur != nil {
		dl.CaretClass = "caret"
		dl.Class += " current"
		if *cur {
			dl.Class += " dropup"
		}
	}
	dl.Href = "?" + ep.Query.ValuesEncode(values)
	return dl
}

// SetValue modifies the values.
func (ep EnumParam) SetValue(values url.Values, pname string, uinter enums.Uinter) (string, *bool) {
	this := uinter.Touint()
	_, low, err := uinter.Marshal()
	if err != nil { // ignoring the error
		return "", nil
	}

	text := ep.EnumDecodec.Text(strings.ToUpper(low))
	ddef := ep.EnumDecodec.Default.Uint
	dnum := ep.Number

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
	if this == ddef && ep.Specified {
		values.Del(pname)
		return text, ret
	}
	if this == dnum.Uint && !dnum.Negative {
		low = "-" + low
	}
	values.Set(pname, low)
	return text, ret
}

type EnumDecodec struct {
	Pname   string
	Default Number
	Alphas  []enums.Uint
	Unew    func() (string, Upointer) `json:"-"` // func cannot be marshaled
	Text    func(string) string       `json:"-"` // func cannot be marshaled
}

func (ep EnumParam) IsAlpha(p enums.Uint) bool {
	for _, u := range ep.EnumDecodec.Alphas {
		if u == p {
			return true
		}
	}
	return false
}

func (ed EnumDecodec) DefaultParam(params *Params) EnumParam {
	return EnumParam{
		Query:       params.Query,
		EnumDecodec: ed,
	}
}

// TextFunc constructs string replacement func.
// ab defines exact mapping, miss-case: fs funcs chain-apply.
func TextFunc(ab map[string]string, fs ...func(string) string) func(string) string {
	return func(s string) string {
		if n, ok := ab[s]; ok {
			return n
		}
		for _, f := range fs {
			s = f(s)
		}
		return s
	}
}

var EnumDecodecs = map[string]EnumDecodec{
	"ps": {
		Default: Number{Uint: enums.Uint(enums.PID)},
		Alphas:  []enums.Uint{enums.Uint(enums.NAME), enums.Uint(enums.USER)},
		Unew:    func() (string, Upointer) { return "ps", new(enums.UintPS) },
		Text:    TextFunc(map[string]string{"PRI": "PR", "NICE": "NI", "NAME": "COMMAND"}, strings.ToUpper),
	},
	"df": {
		Default: Number{Uint: enums.Uint(enums.FS)},
		Alphas:  []enums.Uint{enums.Uint(enums.FS), enums.Uint(enums.MP)},
		Unew:    func() (string, Upointer) { return "df", new(enums.UintDF) },
		Text:    TextFunc(map[string]string{"FS": "Device", "MP": "Mounted"}, strings.ToLower, strings.Title),
	},
}

var BoolDecodecs = map[string]BoolDecodec{
	"still": {Default: false},

	// "hidecpu":  {Default: false},
	// "hidedf":   {Default: false},
	// "hideif":   {Default: false},
	"hidemem": {Default: false},
	// "hideps":   {Default: false},
	"hideswap": {Default: false},
	// "hidevg":   {Default: false},
	// commented-out hide* to be un-commented

	"showconfigmem": {Default: false},
	// rest of showconfig* to follow
}

var PeriodParanames = []string{
	"refreshmem",
}

func (bp *BoolParam) Decode(form url.Values) {
	values, ok := form[bp.BoolDecodec.Pname]
	if !ok {
		bp.Value = bp.BoolDecodec.Default
		bp.Query.Del(bp.BoolDecodec.Pname)
		return
	}
	if len(values) != 0 || values[0] == "" || values[0] == "1" || values[0] == "true" || values[0] == "TRUE" {
		bp.Value = true
	} // else .Value stays false
	bp.Query.Set(bp.BoolDecodec.Pname, bp.StringValue(bp.Value))
}

func (bp BoolParam) StringValue(value bool) string {
	if value != bp.BoolDecodec.Default {
		return ""
	}
	return fmt.Sprintf("%t", value)
}

func (ep *EnumParam) Decode(form url.Values, setep *EnumParam) error {
	_, uptr := ep.EnumDecodec.Unew()
	n, spec, err := ep.Find(form[ep.EnumDecodec.Pname], uptr)
	if err != nil {
		return err
	}
	ep.Number = n
	ep.Specified = spec
	if setep != nil {
		*setep = *ep
	}
	return nil
}

// Find side effects: uptr.Unmarshal and p.Values.Set
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
				ep.Query.Set(ep.EnumDecodec.Pname, l)
			}
			ep.Query.Moved = true
		}
		return Number{}, true, err
	}
	n := Number{
		Uint:     uptr.Touint(),
		Negative: negate,
	}
	ep.Query.Set(ep.EnumDecodec.Pname, values[0])
	return n, true, nil
}

// NewParams constructs new Params.
// Global var BoolDecodecs, PeriodParanames are ranged.
func NewParams(minperiod flags.Period) *Params {
	p := &Params{Query: &Query{Values: make(url.Values)}}
	bools := make(map[string]*BoolParam)
	for k, v := range BoolDecodecs {
		v.Pname = k
		bools[k] = &BoolParam{
			Query:       p.Query,
			BoolDecodec: v,
		}
	}
	periods := make(map[string]*PeriodParam)
	for _, k := range PeriodParanames {
		periods[k] = &PeriodParam{
			Query: p.Query,
			PeriodDecodec: PeriodDecodec{
				Pname:       k,
				Placeholder: minperiod,
			},
			Value: flags.Period{Above: &minperiod.Duration},
		}
	}
	p.ENUM = NewParamsENUM(p)
	p.BOOL = bools
	p.PERIOD = periods
	return p
}

// NewParamsENUM returns ENUM part of Params.
// Global var EnumDecodecs is ranged.
func NewParamsENUM(p *Params) map[string]*EnumParam {
	if p == nil {
		p = &Params{}
		p.NewQuery()
	}
	enums := make(map[string]*EnumParam)
	for k, v := range EnumDecodecs {
		v.Pname = k
		enums[k] = &EnumParam{
			Query:       p.Query,
			EnumDecodec: v,
		}
	}
	return enums
}

// EnumParam represents enum parameter. Features MarshalJSON method
// thus all fields are explicitly marked as non-marshallable.
type EnumParam struct {
	Query       *Query      `json:"-"` // url.Values here.
	EnumDecodec EnumDecodec `json:"-"` // Read-only, an entry from global var EnumDecodecs.
	Number      Number      `json:"-"` // Decoded Number.
	Specified   bool        `json:"-"` // True if a valid value was specified for decoding.
}

func (ep EnumParam) LessorMore(r bool) bool {
	// numeric values: flip r
	if !ep.IsAlpha(ep.Number.Uint) {
		r = !r
	}
	if ep.Number.Negative {
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
	uter := uptr.(enums.Uinter) // downcast. Upointer inlines Uinter.
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

func (bp BoolParam) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Href  template.HTMLAttr
		Value bool
	}{
		Href:  bp.EncodeToggle(),
		Value: bp.Value,
	})
}

func (p *Params) NewQuery() {
	p.Query = &Query{ // new Query
		Values: make(url.Values),
	}
}

func (p *Params) Decode(form url.Values) {
	// for _, v := range p.ENUM { v.Decode(form) }
	for _, v := range p.BOOL {
		v.Decode(form)
	}
	for _, v := range p.PERIOD {
		v.Decode(form)
	}
	/*
		ruri := "?" + p.Query.ValuesEncode(nil)
		for _, v := range p.PERIOD {
			v.FormAction = ruri
		} // */
}

type Params struct {
	ENUM   map[string]*EnumParam
	BOOL   map[string]*BoolParam
	PERIOD map[string]*PeriodParam
	Query  *Query `json:"-"`
}

type Query struct {
	url.Values
	Moved bool
}

// func (q Query) RequestURI() string { return "?" + q.ValuesEncode(nil) }
// func (q Query) MarshalJSON() ([]byte, error) { return json.Marshal(struct{ RequestURI string }{q.RequestURI()}) }

func (q Query) ValuesCopy() url.Values {
	copy := url.Values{}
	for k, v := range q.Values {
		copy[k] = v
	}
	return copy
}

type PeriodDecodec struct {
	Pname       string `json:"-"`
	Placeholder flags.Period
}

// PeriodParam represents period parameter.
type PeriodParam struct {
	Query         *Query `json:"-"` // Explicitly non-marshallable url.Values.
	PeriodDecodec        // Read-only an entry from global var BoolDecoders.
	Value         flags.Period
	// FormAction string
}

func (pp *PeriodParam) Decode(form url.Values) {
	values, ok := form[pp.PeriodDecodec.Pname]
	if ok && len(values) > 0 {
		if err := pp.Value.Set(values[0]); err == nil {
			pp.Query.Set(pp.PeriodDecodec.Pname, pp.Value.String())
			return
		}
	}
	pp.Query.Del(pp.PeriodDecodec.Pname)
}

type BoolDecodec struct {
	Pname   string
	Default bool
}

// BoolParam represents bool parameter. Features MarshalJSON method
// thus all fields are explicitly marked as non-marshallable.
type BoolParam struct {
	Query       *Query      `json:"-"` // url.Values here.
	BoolDecodec BoolDecodec `json:"-"` // Read-only, an entry from global var BoolDecoders.
	Value       bool        `json:"-"` // Decoded value.
}

// EncodeToggle returns template.HTMLAttr having the bp value inverted and encoded.
// The other values are copied from bp.Query.Values.
func (bp BoolParam) EncodeToggle() template.HTMLAttr {
	values := bp.Query.ValuesCopy()
	value := !bp.Value // here's the toggle
	if value == bp.BoolDecodec.Default {
		values.Del(bp.BoolDecodec.Pname)
	} else {
		values.Set(bp.BoolDecodec.Pname, bp.StringValue(value))
	}
	return template.HTMLAttr("?" + bp.Query.ValuesEncode(values))
}

func (q Query) ValuesEncode(v url.Values) string {
	if v == nil {
		v = q.Values
	}
	if v == nil {
		return ""
	}
	var buf bytes.Buffer
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		prefix := url.QueryEscape(k)
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(prefix)
			if v == "" {
				continue
			}
			buf.WriteString("=" + url.QueryEscape(v))
		}
	}
	return buf.String()
}
