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

var (
	NumType   = reflect.TypeOf(Num{})
	DelayType = reflect.TypeOf(Delay{})
)

// NewParams constructs new Params.
func NewParams(mindelay flags.Delay) *Params {
	p := &Params{
		Defaults: make(map[interface{}]Num),
		Delays:   make(map[string]*Delay),
		MinDelay: mindelay,
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
			if def, err := NumPrefix(tags[1:], "default"); err == nil { // otherwise err is ignored
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
	Defaults map[interface{}]Num `json:"-"`
	Delays   map[string]*Delay   `json:"-"`
	MinDelay flags.Delay         `json:"-"`
}

type Schema struct {
	// Still is here to be preserved for url encoding.
	// Not in use by Go code, but by js.
	Still Num `url:"still,omitempty"`

	// The NewParams must populate .Delays with EACH *Delay

	CPUd Delay `url:"cpud,omitempty"`
	Dfd  Delay `url:"dfd,omitempty"`
	Ifd  Delay `url:"ifd,omitempty"`
	Memd Delay `url:"memd,omitempty"`
	Psd  Delay `url:"psd,omitempty"`
	Vgd  Delay `url:"vgd,omitempty"`

	// Num encodes a number and config toggle.
	// "Negative" value states config displaying and
	// the absolute value still encodes the number.

	CPUn Num `url:"cpun,default2"`
	Dfn  Num `url:"dfn,default2"`
	Ifn  Num `url:"ifn,default2"`
	Memn Num `url:"memn,default2"`
	Psn  Num `url:"psn,default8"`
	Vgn  Num `url:"vgn,default2"`

	Dft Num `url:"dft,default2,enumerate2,posonly"` // tab, default DFBYTES
	Ift Num `url:"ift,default3,enumerate3,posonly"` // tab, default IFBYTES
	Psk Num `url:"psk,default1,enumerate9"`         // sort, default PID
	Dfk Num `url:"dfk,default1,enumerate5"`         // sort, default FS
}

type Nlinks struct {
	Zero, More, Less ALink
}
type Dlinks struct {
	More, Less ALink
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
		nl.Zero, _ = p.ZeroN(num)
		nl.More, _ = p.MoreN(num)
		nl.Less, _ = p.LessN(num)
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
		dl.More, _ = p.MoreD(d)
		dl.Less, _ = p.LessD(d)
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
			if v, err := p.Vlink(v, j, "", ""); err == nil { // err is gone
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
		return DecodeNum(s)
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
			if href, err := p.HrefToggleNegative(num); err == nil {
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

func DecodePositive(value string) (int, error) {
	i, err := strconv.Atoi(value)
	if err == nil && i < 0 {
		return i, fmt.Errorf("Integer decoded may not be negative")
	}
	return i, err
}

func DecodeNum(value string) (num Num, err error) {
	if len(value) > 0 && (value[0] == '-' || value[0] == '!') {
		num.Negative, value = true, value[1:]
	}
	num.Absolute, err = DecodePositive(value)
	return num, err
}

// ConvertNum is a schema decoder's converter into Num.
func ConvertNum(value string) reflect.Value {
	num, err := DecodeNum(value)
	if err != nil { // err is lost
		return reflect.Value{}
	}
	return reflect.ValueOf(num)
}

// Delay has it's own MarshalJSON.
type Delay struct {
	D       time.Duration
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

// ConvertDelayFunc creates a schema decoder's converter into Delay.
func ConvertDelayFunc(mindelay flags.Delay) func(string) reflect.Value {
	return func(value string) reflect.Value {
		d := flags.Delay{Above: &mindelay.Duration}
		if err := d.Set(value); err != nil {
			return reflect.Value{}
		}
		return reflect.ValueOf(Delay{D: d.Duration})
	}
}

func (p *Params) ResetSchema() {
	val := reflect.ValueOf(&p.Schema).Elem()
	for typ, i := val.Type(), 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		fv := val.Field(i)
		switch sf.Type {
		case NumType:
			fv.Set(reflect.ValueOf(Num{}))
		case DelayType:
			fv.Set(reflect.ValueOf(Delay{}))
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

	dec := schema.NewDecoder()
	dec.SetAliasTag("url")
	dec.IgnoreUnknownKeys(true)
	dec.ZeroEmpty(true)
	dec.RegisterConverter(Num{}, ConvertNum)
	dec.RegisterConverter(Delay{}, ConvertDelayFunc(p.MinDelay))

	p.ResetSchema()
	err := dec.Decode(&p.Schema, req.Form)
	if err != nil {
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

// RenamedConstError denotes an error.
type RenamedConstError string

func (rc RenamedConstError) Error() string { return string(rc) }
