package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

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

var LimitDecodecs = map[string]LimitDecodec{"psn": {Default: 8}}

var BoolDecodecs = map[string]BoolDecodec{
	"still": {Default: false},

	"hidecpu": {Default: false},
	// "hidedf":   {Default: false},
	// "hideif":   {Default: false},
	"hidemem":  {Default: false},
	"hideps":   {Default: false},
	"hideswap": {Default: false},
	"hidevg":   {Default: false},
	// commented-out hide* to be un-commented

	"showconfigcpu": {Default: false},
	"showconfigmem": {Default: false},
	"showconfigps":  {Default: false},
	"showconfigvg":  {Default: false},
	// rest of showconfig* to follow

	"expandcpu": {Default: false},
}

var PeriodParanames = []string{
	"refreshcpu",
	"refreshmem",
	"refreshps",
	"refreshvg",
}

func (bp *BoolParam) Decode(form url.Values) {
	values, ok := form[bp.BoolDecodec.Pname]
	if !ok {
		bp.Value = bp.BoolDecodec.Default
		bp.Query.Del(bp.BoolDecodec.Pname)
		return
	}
	var v string
	if len(values) > 0 {
		v = values[0]
		if v == "1" || v == "true" || v == "True" || v == "TRUE" || v == "on" {
			bp.QuerySet(nil, true)
			bp.Query.Moved = true
			return
		}
		if v != "" { // any value is invalid. Absence of parameter designates "false"
			bp.Query.Del(bp.BoolDecodec.Pname)
			return
		}
		bp.Value = true // v == ""
	}
	bp.QuerySet(nil, bp.Value)
}

// QuerySet modifies the values.
func (bp BoolParam) QuerySet(values url.Values, value bool) {
	if values == nil {
		values = bp.Query.Values
	}
	if value == bp.BoolDecodec.Default {
		values.Del(bp.BoolDecodec.Pname)
	} else {
		values.Set(bp.BoolDecodec.Pname, "")
	}
}

func (ep *EnumParam) Decode(form url.Values) error {
	_, uptr := ep.EnumDecodec.Unew()
	n, spec, err := ep.Find(form[ep.EnumDecodec.Pname], uptr)
	if err != nil {
		return err
	}
	ep.Number = n
	ep.Specified = spec
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
	p := &Params{}
	p.NewQuery()
	p.ENUM = NewParamsENUM(p)
	p.BOOL = make(map[string]*BoolParam)
	for k, v := range BoolDecodecs {
		v.Pname = k
		p.BOOL[k] = &BoolParam{
			Query:       p.Query,
			BoolDecodec: v,
		}
	}
	p.LIMIT = make(map[string]*LimitParam)
	for k, v := range LimitDecodecs {
		v.Pname = k
		p.LIMIT[k] = &LimitParam{
			Query:        p.Query,
			LimitDecodec: v,
		}
	}
	p.PERIOD = make(map[string]*PeriodParam)
	for _, k := range PeriodParanames {
		p.PERIOD[k] = &PeriodParam{
			Query: p.Query,
			PeriodDecodec: PeriodDecodec{
				Pname: k,
			},
			Placeholder: minperiod,
			Period:      flags.Period{Above: &minperiod.Duration},
		}
	}
	return p
}

// NewParamsENUM returns ENUM part of Params.
// Global var EnumDecodecs is ranged.
func NewParamsENUM(p *Params) map[string]*EnumParam {
	if p == nil {
		p = &Params{}
		p.NewQuery()
	}
	enumap := make(map[string]*EnumParam)
	for k, v := range EnumDecodecs {
		v.Pname = k
		enumap[k] = &EnumParam{
			Query:       p.Query,
			EnumDecodec: v,
		}
	}
	return enumap
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
	for _, v := range p.BOOL {
		v.Query = p.Query
	}
	for _, v := range p.ENUM {
		v.Query = p.Query
	}
	for _, v := range p.LIMIT {
		v.Query = p.Query
	}
	for _, v := range p.PERIOD {
		v.Query = p.Query
	}
}

func (p *Params) Decode(form url.Values) {
	for _, v := range p.BOOL {
		v.Decode(form)
	}
	for _, v := range p.ENUM {
		v.Decode(form)
	}
	for _, v := range p.LIMIT {
		v.Decode(form)
	}
	for _, v := range p.PERIOD {
		v.Decode(form)
	}
}

func (p *Params) Refresh(force bool) bool {
	if force {
		return true
	}
	for _, v := range p.PERIOD {
		if v.Tick <= 1 {
			return true
		}
	}
	return false
}

func (p *Params) Tick() {
	for _, v := range p.PERIOD {
		v.Tick++
		if v.Tick-1 >= int(time.Duration(v.Period.Duration)/time.Second) {
			v.Tick = 1 // expired
		}
	}
}

// TODO .Refresh method is used in ostent.Set/Refresher only. To be removed.
func (p PeriodParam) Refresh(forcerefresh bool) bool {
	if forcerefresh {
		return true
	}
	return p.Tick <= 1 // p.expired()
}

type Params struct {
	BOOL   map[string]*BoolParam
	ENUM   map[string]*EnumParam
	LIMIT  map[string]*LimitParam
	PERIOD map[string]*PeriodParam
	Query  *Query
}

type Query struct {
	url.Values
	Moved          bool
	UpdateLocation bool
}

func (q Query) MarshalJSON() ([]byte, error) {
	return json.Marshal(url.QueryEscape(q.ValuesEncode(nil)))
}

func (q Query) ValuesCopy() url.Values {
	copy := url.Values{}
	for k, v := range q.Values {
		copy[k] = v
	}
	return copy
}

type PeriodDecodec struct {
	Pname string `json:"-"`
}

// PeriodParam represents period parameter.
type PeriodParam struct {
	Query         *Query `json:"-"` // Explicitly non-marshallable url.Values.
	PeriodDecodec        // Read-only an entry from global var BoolDecoders.
	// PeriodDecodec is inlined for access to .Pname for templates
	Placeholder flags.Period
	Period      flags.Period
	Input       string
	InputErrd   bool

	// Params.Tick() must be called once per second; .tick is 1 when the period expired.
	Tick int `json:"-"` // Not to be marshaled.
}

func (pp *PeriodParam) Decode(form url.Values) {
	values, ok := form[pp.PeriodDecodec.Pname]
	if !ok || len(values) == 0 {
		return
	}
	// empty values[0] is ok
	if err := pp.Period.Set(values[0]); err != nil {
		pp.Input, pp.InputErrd = values[0], true
	} else {
		pp.Input, pp.InputErrd = pp.Period.String(), false
	}
	if pp.Period.Duration != pp.Placeholder.Duration {
		pp.Query.Set(pp.PeriodDecodec.Pname, pp.Input)
	} else {
		pp.Query.Del(pp.PeriodDecodec.Pname)
	}
	pp.Query.UpdateLocation = true // New location.
}

type LimitDecodec struct {
	Pname   string
	Default uint
}

type LimitParam struct {
	Query        *Query       `json:"-"` // url.Values here.
	LimitDecodec LimitDecodec `json:"-"` // Read-only, an entry from global var LimitDecoders.
	Value        uint         `json:"-"` // Decoded value.
}

func (lp *LimitParam) Decode(form url.Values) {
	values, ok := form[lp.LimitDecodec.Pname]
	if !ok {
		lp.Value = lp.LimitDecodec.Default
		return
	}
	if len(values) == 0 {
		lp.Query.Del(lp.LimitDecodec.Pname)
		return
	}
	if i, err := strconv.Atoi(values[0]); err == nil && i > 0 && i <= 65536 {
		lp.Value = uint(i)
		lp.Query.Set(lp.LimitDecodec.Pname, fmt.Sprintf("%d", lp.Value))
		return
	}
	lp.Query.Del(lp.LimitDecodec.Pname)
}

func (lp LimitParam) EncodeLess() template.HTMLAttr {
	value := lp.Value
	if value >= 2 {
		g := math.Log2(float64(value))
		n := math.Floor(g)
		if n == g {
			n--
		}
		value = uint(math.Pow(2, n))
	}
	return lp.Encode(&value)
}

func (lp LimitParam) EncodeMore() template.HTMLAttr {
	value := lp.Value
	if value <= 32768 { // up to 65536
		value = uint(math.Pow(2, 1+math.Floor(math.Log2(float64(value)))))
	}
	return lp.Encode(&value)
}

func (lp LimitParam) Encode(this *uint) template.HTMLAttr {
	if this == nil {
		this = &lp.Value
	}
	values := lp.Query.ValuesCopy()
	if *this != lp.LimitDecodec.Default {
		values.Set(lp.LimitDecodec.Pname, fmt.Sprintf("%d", *this))
	} else {
		values.Del(lp.LimitDecodec.Pname)
	}
	return template.HTMLAttr("?" + lp.Query.ValuesEncode(values))
}

func (lp LimitParam) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LessHref, MoreHref template.HTMLAttr
	}{
		LessHref: lp.EncodeLess(),
		MoreHref: lp.EncodeMore(),
	})
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
	bp.QuerySet(values, !bp.Value)
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
