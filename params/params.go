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

var (
	NumType   = reflect.TypeOf(Num{})
	DelayType = reflect.TypeOf(Delay{})
)

// NewParams constructs new Params.
func NewParams(dbounds flags.DelayBounds) *Params {
	p := &Params{
		Defaults:    make(map[interface{}]Num),
		Delays:      make(map[string]*Delay),
		DelayBounds: dbounds,
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
				err = def.UnmarshalText([]byte("-1")) // the default "default"
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
	Defaults    map[interface{}]Num `json:"-"`
	Delays      map[string]*Delay   `json:"-"`
	DelayBounds flags.DelayBounds   `json:"-"`
}

type Schema struct {
	// Still is here to be preserved for url encoding.
	// Not in use by Go code, but by js.
	Still Num `url:"still,posonly,default0"`

	// The NewParams must populate .Delays with EACH *Delay

	CPUd Delay `url:"cpud,omitempty"`
	Dfd  Delay `url:"dfd,omitempty"`
	Ifd  Delay `url:"ifd,omitempty"`
	Lad  Delay `url:"lad,omitempty"`
	Memd Delay `url:"memd,omitempty"`
	Psd  Delay `url:"psd,omitempty"`

	// Num encodes a number and config toggle.
	// "Negative" value states config displaying and
	// the absolute value still encodes the number.

	CPUn Num `url:"cpun,default-2"`
	Dfn  Num `url:"dfn,default-2"`
	Ifn  Num `url:"ifn,default-2"`
	Lan  Num `url:"lan,default-3"`
	Memn Num `url:"memn,default-2"`
	Psn  Num `url:"psn,default0"`

	Psk Num `url:"psk,default1,enumerate9"` // sort, default PID
	Dfk Num `url:"dfk,default1,enumerate6"` // sort, default FS
}

func (p Params) NonZeroPsn() bool { return p.Psn.Absolute != 0 }

type Nlinks struct {
	More, Less ALink
}
type Dlinks struct {
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
		nl.More, _ = MoreN(&p, num)
		nl.Less, _ = LessN(&p, num)
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
		dl.More, _ = MoreD(&p, d)
		dl.Less, _ = LessD(&p, d)
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

func (d Delay) Type() string { return "delay" }

func (d *Delay) Set(input string) error { return d.UnmarshalText([]byte(input)) }

func (d *Delay) UnmarshalText(text []byte) error {
	f := flags.Delay{Above: d.Above, Below: d.Below}
	if err := f.Set(string(text)); err != nil {
		return err
	}
	d.D = f.Duration
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
		case DelayType:
			fv.Set(reflect.ValueOf(Delay{
				Above: &p.DelayBounds.Min.Duration,
				Below: &p.DelayBounds.Max.Duration,
			}))
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
		case DelayType:
			d := fv.Addr().Interface().(*Delay)
			d.Default = p.DelayBounds.Min.Duration
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
			d.D = p.DelayBounds.Min.Duration
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
	} {
		if _, ok := req.Form[name]; ok {
			req.Form.Del(name)
			moved = true
		}
	}

	dec := schema.NewDecoder()
	dec.ZeroEmpty(true)
	dec.IgnoreUnknownKeys(true)
	dec.SetAliasTag("url") // single tag for decoding AND encoding

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

// SetURL sets the .URL.
func (ep *Endpoint) SetURL(u url.URL) { ep.URL = u }

// String return string repr.
func (ep Endpoint) String() string { return strings.TrimPrefix(ep.URL.String(), "http://") }

// Endpoint has an URL and other fields decoded from it.
type Endpoint struct {
	// URL is the base.
	URL url.URL `schema:"-"`

	// ServerAddr is server part (host[:port]) of URL.
	ServerAddr flags.Bind `schema:"-"`

	// The schema fields:
	// Delay is the delay param.
	Delay Delay `schema:"delay,omitempty"`
}

// UseURL sets ep.ServerAddr based on u.
func (ep *Endpoint) UseURL(u url.URL) error {
	u.RawQuery = "" // won't use query string in ServerAddr
	return ep.ServerAddr.Set(u.Host)
}

func NewGraphiteEndpoints(delay time.Duration, bind flags.Bind) GraphiteEndpoints {
	return GraphiteEndpoints{Default: Endpoint{ServerAddr: bind, Delay: Delay{D: delay}}}
}

// GraphiteEndpoints holds graphite endpoints list.
type GraphiteEndpoints struct {
	Values  []Endpoint
	Default Endpoint
}

// Set is a flag.Value method.
func (gends *GraphiteEndpoints) Set(input string) error {
	values := strings.Split(input, ",")
	gends.Values = make([]Endpoint, len(values))
	for i, value := range values {
		gends.Values[i] = gends.Default // copy
		if _, err := Decode(nil, AddScheme(value), false, &gends.Values[i], &gends.Values[i]); err != nil {
			return err
		}
		if gends.Values[i].ServerAddr.Host == "" {
			return fmt.Errorf("server address required for Graphite exporting")
		}
	}
	return nil
}

// String is a flag.Value method.
func (gends GraphiteEndpoints) String() string {
	values := gends.Values // shortcut
	ss := make([]string, len(values))
	for i, v := range values {
		ss[i] = v.String()
	}
	return strings.Join(ss, ",")
}

// Type is a pflag.Value method.
func (gends GraphiteEndpoints) Type() string { return "graphiteEndpoints" }

// InfluxEndpoint holds influxdb params.
type InfluxEndpoint struct {
	Endpoint
	Database string            `schema:"database,omitempty"`
	Username string            `schema:"username,omitempty"`
	Password string            `schema:"password,omitempty"`
	Tags     map[string]string `schema:"-"`
}

// String is a fmt.Stringer method.
func (iend InfluxEndpoint) String() string {
	urlcopy := iend.Endpoint.URL
	qvalues := urlcopy.Query()
	if v := qvalues.Get("password"); v != "" {
		qvalues.Set("password", strings.Repeat(" ", len(v)))
		urlcopy.RawQuery = qvalues.Encode()
	}
	return urlcopy.String()
}

func NewInfluxEndpoints(delay time.Duration, database string) InfluxEndpoints {
	return InfluxEndpoints{Default: InfluxEndpoint{
		Endpoint: Endpoint{Delay: Delay{D: delay}},
		Database: database,
	}}
}

// InfluxEndpoints holds infuxdb endpoints list.
type InfluxEndpoints struct {
	Values  []InfluxEndpoint
	Default InfluxEndpoint
}

// Set is a flag.Value method.
func (iends *InfluxEndpoints) Set(input string) error {
	values := strings.Split(input, ",")
	iends.Values = make([]InfluxEndpoint, len(values))
	for i, value := range values {
		iends.Values[i] = iends.Default // copy
		uvalues, err := Decode(nil, value, true, &iends.Values[i], &iends.Values[i])
		if err != nil {
			return err
		}
		if iends.Values[i].ServerAddr.Host == "" {
			return fmt.Errorf("server address required (prefixed with explicit scheme e.g. http://) for InfluxDB exporting")
		}

		// now .Tags to be set if any
		for _, k := range []string{
			"delay",
			"database",
			"username",
			"password",
		} {
			delete(uvalues, k)
		}
		if klen := len(uvalues); klen != 0 {
			iends.Values[i].Tags = make(map[string]string, klen)
			for k, v := range uvalues {
				iends.Values[i].Tags[k] = v[0]
			}
			// map[string][]string -> map[string]string
		}
	}
	return nil
}

// String is a flag.Value method.
func (iends InfluxEndpoints) String() string {
	values := iends.Values // shortcuts
	ss := make([]string, len(values))
	for i, v := range values {
		ss[i] = v.String()
	}
	return strings.Join(ss, ",")
}

// Type is a pflag.Value method.
func (iends InfluxEndpoints) Type() string { return "infuxEndpoints" }

// LibratoEndpoint holds librato params.
type LibratoEndpoint struct {
	Endpoint
	Email, Token, Source string
}

func NewLibratoEndpoints(delay time.Duration, source string) LibratoEndpoints {
	return LibratoEndpoints{Default: LibratoEndpoint{
		Endpoint: Endpoint{Delay: Delay{D: delay}},
		Source:   source,
	}}
}

// LibratoEndpoints holds librato endpoints list.
type LibratoEndpoints struct {
	Values  []LibratoEndpoint
	Default LibratoEndpoint
}

// Set is a flag.Value method.
func (lends *LibratoEndpoints) Set(input string) error {
	values := strings.Split(input, ",")
	lends.Values = make([]LibratoEndpoint, len(values))
	for i, value := range values {
		lends.Values[i] = lends.Default // copy
		if _, err := Decode(nil, AddScheme(value), false, &lends.Values[i], nil); err != nil {
			return err
		}
		l := &lends.Values[i] // shortcut
		if l.Email == "" {
			return fmt.Errorf("email param required for Librato exporting")
		}
		if l.Token == "" {
			return fmt.Errorf("token param required for Librato exporting")
		}
		if l.Source == "" {
			return fmt.Errorf("source param required for Librato exporting")
		}
	}
	return nil
}

// String is a flag.Value method.
func (lends LibratoEndpoints) String() string {
	values := lends.Values // shortcut
	ss := make([]string, len(values))
	for i, v := range values {
		ss[i] = v.String()
	}
	return strings.Join(ss, ",")
}

// Type is a pflag.Value method.
func (lends LibratoEndpoints) Type() string { return "libratoEndpoints" }

// FetchKey encloses an Endpoint and has extra params.
type FetchKey struct {
	Endpoint
	Schema
	Times int `schema:"times"`
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
		if err := urluser.UseURL(*u); err != nil {
			return nil, err
		}
	}
	dec := schema.NewDecoder()
	dec.ZeroEmpty(true)
	if ignoreUnknownKeys {
		dec.IgnoreUnknownKeys(true)
	}
	values := u.Query()
	err = dec.Decode(into, values)
	return values, err
}
