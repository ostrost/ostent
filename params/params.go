package params

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/gorilla/schema"

	"github.com/ostrost/ostent/flags"
)

// Constants for DF sorting criterion.
const (
	_      int = iota
	FS         // 1
	MP         // 2
	AVAIL      // 3
	USEPCT     // 4
	USED       // 5
	TOTAL      // 6
)

// Constants for PS sorting criterion.
const (
	_    int = iota
	PID      // 1
	UID      // 2
	USER     // 3
	PRI      // 4
	NICE     // 5
	VIRT     // 6
	RES      // 7
	TIME     // 8
	NAME     // 9
)

var (
	NumType   = reflect.TypeOf(Num{})
	DelayType = reflect.TypeOf(Delay{})
)

// NewParams constructs new Params.
func NewParams(mindelay, maxdelay flags.Delay) *Params {
	p := &Params{
		Defaults: make(map[interface{}]Num),
		Delays:   make(map[string]*Delay),
		MinDelay: mindelay,
		MaxDelay: maxdelay,
	}

	val := reflect.ValueOf(&p.Schema).Elem()
	for typ, i := val.Type(), 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		tags, ok := TagsOk(sf)
		if !ok || tags[0] == "" {
			continue
		}
		fv := val.Field(i)
		if sft := sf.Type; sft == NumType {
			def, err := NumPrefix(tags[1:], "default")
			if err != nil {
				err = def.UnmarshalText([]byte("1")) // the default "default"
			}
			if err == nil {
				num := fv.Addr().Interface()
				p.Defaults[sf.Name] = def // ?
				p.Defaults[num] = def
			}
		} else if sft == DelayType {
			d := fv.Addr().Interface().(*Delay)
			p.Delays[tags[0]] = d
			// p.Defaults[d] = def
		}
	}
	return p
}

// Expired satisfying receiver interface.
func (p Params) Expired() bool {
	for _, d := range p.Delays {
		if d.Expired() {
			return true
		}
	}
	return false
}

// Expired satisfying receiver interface.
func (d Delay) Expired() bool { return d.Ticks <= 1 }

func (p *Params) Tick() {
	for _, d := range p.Delays {
		d.Tick()
	}
}

func (d *Delay) Tick() {
	d.Ticks++
	if d.Ticks-1 >= int(d.D/time.Second) {
		d.Ticks = 1 // expired
	}
}

type Params struct {
	Schema
	ParamsFuncs
	Defaults map[interface{}]Num `json:"-"`
	Delays   map[string]*Delay   `json:"-"`
	MinDelay flags.Delay         `json:"-"`
	MaxDelay flags.Delay         `json:"-"`
}

type Schema struct {
	// Still is here to be preserved for url encoding.
	// Not in use by Go code, but by js.
	Still Num `url:"still,posonly,default0"`

	// The NewParams must populate .Delays with EACH *Delay

	CPUd Delay `url:"cpud,omitempty"`
	Dfd  Delay `url:"dfd,omitempty"`
	Ifd  Delay `url:"ifd,omitempty"`
	Memd Delay `url:"memd,omitempty"`
	Psd  Delay `url:"psd,omitempty"`

	// Num encodes a number and config toggle.
	// "Negative" value states config displaying and
	// the absolute value still encodes the number.

	CPUn Num `url:"cpun,default2"`
	Dfn  Num `url:"dfn,default2"`
	Ifn  Num `url:"ifn,default2"`
	Memn Num `url:"memn,default2"`
	Psn  Num `url:"psn,default8"`

	Psk Num `url:"psk,default1,enumerate9"` // sort, default PID
	Dfk Num `url:"dfk,default1,enumerate6"` // sort, default FS
}

type Nlinks struct {
	More, Less ALink
}
type Dlinks struct {
	More, Less ALink
}

type ALink struct {
	Href       string
	Text       string
	Badge      string `json:"-"`
	Class      string `json:"-"`
	ExtraClass string `json:",omitempty"`
}

type VLink struct {
	AlignClass string
	CaretClass string
	LinkClass  string
	LinkHref   string
	LinkText   string `json:"-"` // static
}

func (p *Params) MarshalJSON() ([]byte, error) {
	d := struct {
		Schema
		Tlinks map[string]string
		Dlinks map[string]Dlinks
		Nlinks map[string]Nlinks
		Vlinks map[string][]VLink
	}{
		Schema: p.Schema,
		Tlinks: p.Tlinks(),
		Dlinks: p.Dlinks(),
		Nlinks: p.Nlinks(),
		Vlinks: p.Vlinks(),
	}
	return json.Marshal(d)
}

func (p Params) Nlinks() map[string]Nlinks {
	m := make(map[string]Nlinks)
	val := reflect.ValueOf(&p.Schema).Elem()
	for typ, i := val.Type(), 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		if sf.Type != NumType {
			continue
		}
		tags, ok := TagsOk(sf)
		if !ok || tags[0] == "" {
			continue
		}
		num := val.Field(i).Addr().Interface().(*Num)
		nl := Nlinks{}
		// errors are ignored
		nl.More, _ = p.MoreN(&p, num, "")
		nl.Less, _ = p.LessN(&p, num, "")
		m[sf.Name] = nl
	}
	return m
}

func (p Params) Dlinks() map[string]Dlinks {
	m := make(map[string]Dlinks)
	val := reflect.ValueOf(&p.Schema).Elem()
	for typ, i := val.Type(), 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		if sf.Type != DelayType {
			continue
		}
		tags, ok := TagsOk(sf)
		if !ok || tags[0] == "" {
			continue
		}
		d := val.Field(i).Addr().Interface().(*Delay)
		dl := Dlinks{}
		// errors are ignored
		dl.More, _ = p.MoreD(&p, d, "")
		dl.Less, _ = p.LessD(&p, d, "")
		m[sf.Name] = dl
	}
	return m
}

func (p *Params) Vlinks() map[string][]VLink {
	m := make(map[string][]VLink)
	val := reflect.ValueOf(&p.Schema).Elem()
	for typ, i := val.Type(), 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		if sf.Type != NumType {
			continue
		}
		tags, ok := TagsOk(sf)
		if !ok || tags[0] == "" {
			continue
		}
		v := val.Field(i).Addr().Interface().(*Num)
		var vl []VLink
		maxn, err := NumPrefix(tags[1:], "enumerate")
		if err != nil { // err is gone
			continue
		}
		for j := 1; j < maxn.Absolute+1; j++ { // indexed from 1
			if v, err := p.Vlink(p, v, j, "", ""); err == nil { // err is gone
				vl = append(vl, v)
			}
		}
		m[sf.Name] = vl
	}
	return m
}

func ContainsPrefix(words []string, prefix string) (string, bool) {
	for _, w := range words {
		if strings.HasPrefix(w, prefix) {
			return w[len(prefix):], true // string may be ""
		}
	}
	return "", false
}

func NumPrefix(words []string, prefix string) (Num, error) {
	if s, ok := ContainsPrefix(words, prefix); ok && s != "" {
		num := Num{}
		err := num.UnmarshalText([]byte(s))
		return num, err
	}
	return Num{}, fmt.Errorf("%q not prefixing with anything in %+v", prefix, words)
}

func (p *Params) Tlinks() map[string]string {
	m := make(map[string]string)
	val := reflect.ValueOf(&p.Schema).Elem()
	for typ, i := val.Type(), 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		if tags, ok := TagsOk(sf); !ok || tags[0] == "" {
			continue
		}
		if sf.Type == NumType {
			num := val.Field(i).Addr().Interface().(*Num)
			if href, err := p.EncodeT(num); err == nil {
				m[sf.Name] = href
			}
		}
	}
	return m
}

func TagsOk(sf reflect.StructField) ([]string, bool) {
	if sf.PkgPath != "" { // unexported
		return nil, false
	}

	tag := sf.Tag.Get("url")
	if tag == "" || tag == "-" {
		return nil, false
	}
	return strings.Split(tag, ","), true
}

// Num has no MarshalJSON.
type Num struct {
	Negative        bool
	Absolute        int
	DefaultNegative bool `json:"-"`
	DefaultAbsolute int  `json:"-"`
	Limit           int  `json:"-"`
	Alpha           bool `json:"-"`
	PositiveOnly    bool `json:"-"`
}

// EncodeString returns string repr of Num.
// Templates render .Absolute value explicitly.
func (num Num) EncodeString() string {
	var sym string
	if !num.PositiveOnly && num.Negative {
		sym = "-"
		if num.Absolute == 0 {
			sym = "!"
		}
	}
	return fmt.Sprintf("%s%d", sym, num.Absolute)
}

func (num Num) EncodeValues(key string, values *url.Values) error {
	if (!num.PositiveOnly && num.Negative != num.DefaultNegative) || num.Absolute != num.DefaultAbsolute {
		(*values)[key] = []string{num.EncodeString()}
	}
	return nil
}

func (num *Num) UnmarshalText(text []byte) error {
	var negative bool
	if len(text) > 0 && text[0] == '!' {
		negative, text = true, text[1:]
	}
	i, err := strconv.Atoi(string(text))
	if err != nil {
		return err
	}
	num.Negative = i < 0
	if !num.Negative && negative {
		num.Negative = true
	}
	if num.Negative {
		if num.PositiveOnly {
			return fmt.Errorf("Integer decoded may not be negative")
		}
		num.Absolute = -i
	} else {
		num.Absolute = i
	}
	return nil
}

// Delay has it's own MarshalJSON.
type Delay struct {
	D       time.Duration
	Above   *time.Duration
	Below   *time.Duration
	Default time.Duration
	Ticks   int
}

func (d Delay) String() string { return flags.DurationString(d.D) }

func (d Delay) MarshalJSON() ([]byte, error) { return json.Marshal(d.String()) }

func (d Delay) EncodeValues(key string, values *url.Values) error {
	if d.D != d.Default {
		(*values)[key] = []string{d.String()}
	}
	return nil
}

func (d *Delay) UnmarshalText(text []byte) error {
	f := flags.Delay{Above: d.Above, Below: d.Below}
	if err := f.Set(string(text)); err != nil {
		return err
	}
	d.D = f.Duration
	return nil
}

func (p *Params) ResetSchema(mindelay, maxdelay flags.Delay) {
	val := reflect.ValueOf(&p.Schema).Elem()
	for typ, i := val.Type(), 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		fv := val.Field(i)
		switch sf.Type {
		case NumType:
			var posonly bool
			if tags, ok := TagsOk(sf); ok {
				if _, ok := ContainsPrefix(tags[1:], "posonly"); ok {
					posonly = true
				}
			}
			fv.Set(reflect.ValueOf(Num{PositiveOnly: posonly}))
		case DelayType:
			fv.Set(reflect.ValueOf(Delay{
				Above: &mindelay.Duration,
				Below: &maxdelay.Duration,
			}))
		}
	}
}

func (p *Params) SetDefaults(form url.Values, mindelay flags.Delay) {
	val := reflect.ValueOf(&p.Schema).Elem()
	for typ, i := val.Type(), 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		fv := val.Field(i)
		switch sf.Type {
		case NumType:
			num := fv.Addr().Interface().(*Num)
			def, ok := p.Defaults[num]
			if !ok {
				continue
			}
			tags, havetags := TagsOk(sf)
			if havetags {
				if _, ok := ContainsPrefix(tags[1:], "posonly"); ok {
					num.PositiveOnly = true
				}
			}
			num.DefaultNegative = def.Negative
			num.DefaultAbsolute = def.Absolute
			if num.Negative && num.Absolute != 0 { // all values specified, no need for defaults
				continue
			}
			if !havetags {
				continue
			}
			if _, ok := form[tags[0]]; ok { // have parameter
				continue
			}
			if !num.Negative { // not allow false init value
				num.Negative = def.Negative
			}
			if num.Absolute == 0 { // not allow 0 init value
				num.Absolute = def.Absolute
			}
		case DelayType:
			d := fv.Addr().Interface().(*Delay)
			d.Default = mindelay.Duration
			if d.D != time.Duration(0) { // value specified, no need for defaults
				continue
			}
			tags, ok := TagsOk(sf)
			if !ok {
				continue
			}
			if _, ok := form[tags[0]]; ok { // have parameter
				continue
			}
			d.D = mindelay.Duration
		}
	}
}

func (p *Params) Decode(req *http.Request) error {
	if err := req.ParseForm(); err != nil { // do ParseForm even if req.Form == nil
		return err
	}
	var moved bool
	if _, moved = req.Form["df"]; moved {
		req.Form.Del("df")
	}
	if _, ok := req.Form["ps"]; ok {
		req.Form.Del("ps")
		moved = true
	}
	if _, ok := req.Form["dft"]; ok {
		req.Form.Del("dft")
		moved = true
	}
	if _, ok := req.Form["ift"]; ok {
		req.Form.Del("ift")
		moved = true
	}
	if _, ok := req.Form["vgd"]; ok {
		req.Form.Del("vgd")
		moved = true
	}
	if _, ok := req.Form["vgn"]; ok {
		req.Form.Del("vgn")
		moved = true
	}

	dec := schema.NewDecoder()
	dec.SetAliasTag("url")
	dec.IgnoreUnknownKeys(true)
	dec.ZeroEmpty(true)

	p.ResetSchema(p.MinDelay, p.MaxDelay)
	if err := dec.Decode(&p.Schema, req.Form); err != nil {
		return err
	}
	p.SetDefaults(req.Form, p.MinDelay)
	if !moved {
		return nil
	}
	s, err := p.Encode()
	if err != nil {
		return err
	}
	return RenamedConstError("?" + s)
}

func (p Params) Encode() (string, error) {
	values, err := query.Values(p.Schema)
	if err != nil {
		return "", err
	}
	return values.Encode(), nil
}

func (p *Params) EncodeD(d *Delay, set time.Duration) (string, error) {
	copy := d.D
	d.D = set
	qs, err := p.Encode()
	d.D = copy
	return "?" + qs, err
}

func (p *Params) EncodeT(num *Num) (string, error) {
	num.Negative = !num.Negative
	qs, err := p.Encode()
	num.Negative = !num.Negative
	return "?" + qs, err
}

func (p *Params) EncodeN(num *Num, absolute int, setNegative *bool) (string, error) {
	copy, ncopy := num.Absolute, num.Negative
	num.Absolute = absolute
	if setNegative != nil {
		num.Negative = *setNegative
	}
	qs, err := p.Encode()
	num.Absolute = copy
	if setNegative != nil {
		num.Negative = ncopy
	}
	return "?" + qs, err
}

// RenamedConstError denotes an error.
type RenamedConstError string

func (rc RenamedConstError) Error() string { return string(rc) }
