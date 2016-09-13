package params

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/gorilla/schema"
	// "github.com/spf13/pflag"

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

var NumType = reflect.TypeOf(Num{})

// NewParams constructs new Params.
func NewParams() *Params {
	p := &Params{Defaults: make(map[interface{}]Num)}

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
				err = def.UnmarshalText([]byte("-1")) // the default "default"
			}
			if err == nil {
				num := fv.Addr().Interface()
				p.Defaults[sf.Name] = def // ?
				p.Defaults[num] = def
			}
		}
	}
	return p
}

type Params struct {
	Schema
	Defaults map[interface{}]Num `json:"-"`
}

type Schema struct {
	// Still is here to be preserved for url encoding.
	// Not in use by Go code, but by js.
	Still Num `url:"still,posonly,default0"`

	// Num encodes a number and config toggle.
	// "Negative" value states config displaying and
	// the absolute value still encodes the number.

	CPUn Num `url:"cpun,default-2"`
	Dfn  Num `url:"dfn,default-2"`
	Ifn  Num `url:"ifn,default-2"`
	Lan  Num `url:"lan,default-3"`
	Memn Num `url:"memn,default-2"`
	Psn  Num `url:"psn,default-8"`

	Psk Num `url:"psk,default1,enumerate9"` // sort, default PID
	Dfk Num `url:"dfk,default1,enumerate6"` // sort, default FS
}

type Nlinks struct {
	More, Less ALink
}

type ALink struct {
	Href       string
	Text       string
	ExtraClass string `json:",omitempty"`
}

type VLink struct {
	CaretClass string
	LinkClass  string
	LinkHref   string
}

func (p *Params) MarshalJSON() ([]byte, error) {
	d := struct {
		Schema
		Tlinks map[string]string
		Nlinks map[string]Nlinks
		Vlinks map[string][]VLink
	}{
		Schema: p.Schema,
		Tlinks: p.Tlinks(),
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
		nl.More, _ = MoreN(&p, num)
		nl.Less, _ = LessN(&p, num)
		m[sf.Name] = nl
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
			if v, err := Vlink(p, v, j, ""); err == nil { // err is gone
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

func (p *Params) ResetSchema() {
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
		}
	}
}

func (p *Params) SetDefaults(form url.Values) {
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
		}
	}
}

func (p *Params) Decode(req *http.Request) error {
	if err := req.ParseForm(); err != nil { // do ParseForm even if req.Form == nil
		return err
	}
	var moved bool
	for _, name := range []string{
		"df",
		"ps",
		"dft",
		"ift",
		"vgd",
		"vgn",

		"cpud",
		"dfd",
		"ifd",
		"lad",
		"memd",
		"psd",
	} {
		if _, ok := req.Form[name]; ok {
			req.Form.Del(name)
			moved = true
		}
	}

	dec := schema.NewDecoder()
	dec.ZeroEmpty(true)
	dec.SetAliasTag("url") // single tag for decoding AND encoding
	dec.IgnoreUnknownKeys(true)

	p.ResetSchema()
	if err := dec.Decode(&p.Schema, req.Form); err != nil {
		return err
	}
	p.SetDefaults(req.Form)
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

// SetURL sets the .URL.
func (ep *Endpoint) SetURL(u url.URL) { ep.URL = u }

// String return string repr.
func (ep Endpoint) String() string { return strings.TrimPrefix(ep.URL.String(), "http://") }

// Endpoint has an URL and other fields decoded from it.
type Endpoint struct {
	// URL is the base.
	URL url.URL `url:"-"`

	// ServerAddr is server part (host[:port]) of URL.
	ServerAddr flags.Bind `url:"-"`
}

// FetchKey encloses an Endpoint and has extra params.
type FetchKey struct {
	Endpoint
	Schema

	// url tag not used til encoding with query.Values for normalization.
	Times int // `url:"times"`
}

func NewFetchKeys(bind flags.Bind) *FetchKeys {
	def := FetchKey{} // Endpoint: Endpoint{URL: url.URL{Host: bind.String()}}}
	def.URL.Scheme = "http"
	def.URL.Host = bind.String()
	def.URL.Path = "/index.ws"
	return &FetchKeys{Default: def}
}

type FetchKeys struct {
	Values    []FetchKey
	Fragments [][]string
	Default   FetchKey
}

// Set is a flag.Value method.
func (fkeys *FetchKeys) Set(input string) error {
	values := strings.Split(input, ",")
	fkeys.Values = make([]FetchKey, len(values))
	fkeys.Fragments = make([][]string, len(values))
	for i, value := range values {
		newkey := fkeys.Default // copy
		if _, err := Decode(&fkeys.Default.URL, value, false, &newkey, nil); err != nil {
			return err
		}
		if newkey.URL.Path == "" {
			newkey.URL.Path = fkeys.Default.URL.Path
		}
		fkeys.Values[i] = newkey
		fkeys.Fragments[i] = strings.Split(newkey.URL.Fragment, "#")
	}
	return nil
}

// String is a flag.Value method.
func (fkeys FetchKeys) String() string {
	values := fkeys.Values // shortcut
	ss := make([]string, len(values))
	for i, v := range values {
		ss[i] = v.String()
	}
	return strings.Join(ss, ",")
}

// Type is a pflag.Value method.
func (fkeys FetchKeys) Type() string { return "fetchKeys" }

func AddScheme(input string) string { return "http://" + input }

// Decode does url parsing and schema decoding.
func Decode(base *url.URL, input string, ignoreUnknownKeys bool,
	into interface {
		// pflag.Value
		SetURL(url.URL)
	},
	urluser interface {
		UseURL(url.URL) error
	}) (map[string][]string, error) {

	u, err := url.Parse(input)
	if err != nil {
		return nil, err
	}
	if base != nil {
		u = base.ResolveReference(u)
	}
	into.SetURL(*u)
	if urluser != nil {
		if err = urluser.UseURL(*u); err != nil {
			return nil, err
		}
	}
	dec := schema.NewDecoder()
	dec.ZeroEmpty(true)
	dec.SetAliasTag("url") // single tag for decoding AND encoding
	if ignoreUnknownKeys {
		dec.IgnoreUnknownKeys(true)
	}
	values := u.Query()
	err = dec.Decode(into, values)
	return values, err
}
