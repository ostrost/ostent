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
	NumType      = reflect.TypeOf(Num{})
	DurationType = reflect.TypeOf(Duration{})
)

// NewParams constructs new Params.
func NewParams(minperiod flags.Period) *Params {
	p := &Params{
		Defaults:  make(map[interface{}]Num),
		Ticks:     make(map[string]Ticks),
		MinPeriod: minperiod,
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
		} else if sft == DurationType {
			dur := fv.Addr().Interface().(*Duration)
			p.Ticks[tags[0]] = NewTicks(dur)
			// p.Defaults[dur] = def
		}
	}
	return p
}

func (p *Params) Decode(req *http.Request) error {
	err := p.DecodeDecode(req)
	if merr, ok := err.(MultiError); ok {
		p.Errors = merr
	}
	p.Toprows = map[bool]int{true: 1, false: 2}[p.Hideswap]
	return err
}

func (p Params) RefreshFunc(dp *Duration) func(bool) bool {
	return func(force bool) bool {
		if force {
			return true
		}
		for _, ti := range p.Ticks {
			if dp == ti.Duration && ti.Expired() {
				return true
			}
		}
		return false
	}
}

func (p Params) Refresh(force bool) bool {
	if force {
		return true
	}
	for _, ti := range p.Ticks {
		if ti.Expired() {
			return true
		}
	}
	return false
}

func (p Params) Expired() bool {
	return p.Refresh(false)
}

func (p *Params) Tick() {
	for _, v := range p.Ticks {
		v.Tick()
	}
}

func (ti Ticks) Expired() bool { return ti.Ticks <= 1 }

func (ti *Ticks) Tick() {
	ti.Ticks++
	if ti.Ticks-1 >= int(ti.Duration.D/time.Second) {
		ti.Ticks = 1 // expired
	}
}

type Ticks struct {
	Ticks    int
	Duration *Duration
}

func NewTicks(dp *Duration) Ticks {
	return Ticks{Duration: dp}
}

type Params struct {
	Schema
	Defaults  map[interface{}]Num `json:"-"`
	Errors    MultiError          `json:",omitempty"`
	Ticks     map[string]Ticks    `json:"-"`
	Toprows   int                 `json:"-"`
	MinPeriod flags.Period        `json:"-"`
}

type Schema struct {
	Still     bool `url:"still,omitempty"`
	Hidecpu   bool `url:"hidecpu,omitempty"`
	Hidedf    bool `url:"hidedf,omitempty"`
	Hideif    bool `url:"hideif,omitempty"`
	Hidemem   bool `url:"hidemem,omitempty"`
	Hideswap  bool `url:"hideswap,omitempty"`
	Hidevg    bool `url:"hidevg,omitempty"`
	Configcpu bool `url:"configcpu,omitempty"`
	Configdf  bool `url:"configdf,omitempty"`
	Configif  bool `url:"configif,omitempty"`
	Configmem bool `url:"configmem,omitempty"`
	Configvg  bool `url:"configvg,omitempty"`
	Expanddf  bool `url:"expanddf,omitempty"`
	Expandif  bool `url:"expandif,omitempty"`
	Expandcpu bool `url:"expandcpu,omitempty"`

	// Memn Num
	// Cpun Num
	// Dfn Num
	// Ifn Num

	// Psn encodes number of proccesses and ps config toggle.
	// Negative value states config displaying and
	// the absolute value still encodes the ps number.
	Psn Num `url:"psn,default8"`            // limit
	Psk Num `url:"psk,default1,enumerate9"` // sort, default PID
	Dfk Num `url:"dfk,default1,enumerate5"` // sort, default FS
	Dft Num `url:"dft,default2,enumerate2"` // tab, default DFBYTES
	Ift Num `url:"ift,default3,enumerate3"` // tab, default IFBYTES

	// The NewParams must populate .Ticks with EACH *d
	Psd Duration `url:"psd,omitempty"`

	Refreshcpu Duration `url:"refreshcpu,omitempty"`
	Refreshdf  Duration `url:"refreshdf,omitempty"`
	Refreshif  Duration `url:"refreshif,omitempty"`
	Refreshmem Duration `url:"refreshmem,omitempty"`
	Refreshvg  Duration `url:"refreshvg,omitempty"`
}

type Numbered struct {
	Zero, More, Less ALink
}
type Delayed struct {
	More, Less ALink
}

func (p *Params) MarshalJSON() ([]byte, error) {
	d := struct {
		Schema
		Toggle     map[string]string
		Delayed    map[string]Delayed
		Numbered   map[string]Numbered
		Variations map[string][]Varlink
		// Defaults   map[string]Num
	}{
		Schema:     p.Schema,
		Toggle:     p.Toggles(),
		Delayed:    p.Delayed(),
		Numbered:   p.Numbered(),
		Variations: p.Variations(),
		// Defaults:   p.Defaults.StringKeysOnly(),
	}
	return json.Marshal(d)
}

/*
func (p Params) StringKeysOnly() map[string]Num {
	m := make(map[string]Num)
	for k, v := range p.Defaults {
		if s, ok := k.(string); ok {
			m[s] = v
		}
	}
	return m
}
*/

func (p Params) Numbered() map[string]Numbered {
	m := make(map[string]Numbered)
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
		red := Numbered{}
		// errors are ignored
		red.Zero, _ = p.ZeroN(num)
		red.More, _ = p.MoreN(num)
		red.Less, _ = p.LessN(num)
		m[sf.Name] = red
	}
	return m
}

func (p Params) Delayed() map[string]Delayed {
	m := make(map[string]Delayed)
	val := reflect.ValueOf(&p.Schema).Elem()
	for typ, i := val.Type(), 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		if sf.Type != DurationType {
			continue
		}
		tags, ok := TagsOk(sf)
		if !ok || tags[0] == "" {
			continue
		}
		dur := val.Field(i).Addr().Interface().(*Duration)
		dly := Delayed{}
		// errors are ignored
		dly.More, _ = p.MoreD(dur)
		dly.Less, _ = p.LessD(dur)
		m[sf.Name] = dly
	}
	return m
}

func (p *Params) Variations() map[string][]Varlink {
	m := make(map[string][]Varlink)
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
		var links []Varlink
		maxn, err := NumPrefix(tags[1:], "enumerate")
		if err != nil { // err is gone
			continue
		}
		for j := 1; j < maxn.Body+1; j++ { // indexed from 1
			if vl, err := p.Variate(v, j, "", ""); err == nil { // err is gone
				links = append(links, vl)
			}
		}
		m[sf.Name] = links
	}
	return m
}

func NumPrefix(words []string, prefix string) (Num, error) {
	for _, w := range words {
		if strings.HasPrefix(w, prefix) {
			return DecodeNum(w[len(prefix):])
		}
	}
	return Num{}, fmt.Errorf("%q not prefixing in %+v", prefix, words)
}

func (p *Params) Toggles() map[string]string {
	m := make(map[string]string)
	val := reflect.ValueOf(&p.Schema).Elem()
	for typ, i := val.Type(), 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		if tags, ok := TagsOk(sf); !ok || tags[0] == "" {
			continue
		}
		if sf.Type.Kind() == reflect.Bool {
			v := val.Field(i).Addr().Interface().(*bool)
			if s, err := p.HrefToggle(v); err == nil {
				m[sf.Name] = s
			}
		}
		if sf.Type == NumType {
			num := val.Field(i).Addr().Interface().(*Num)
			if href, err := p.HrefToggleHead(num); err == nil {
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

type Num struct {
	Head        bool
	Body        int
	DefaultHead bool
	DefaultBody int
	Limit       int
}

func (num Num) EncodeValues(key string, values *url.Values) error {
	if num.Head != num.DefaultHead || num.Body != num.DefaultBody {
		(*values)[key] = []string{num.String()}
	}
	return nil
}

func (num Num) MarshalJSON() ([]byte, error) { return json.Marshal(num.String()) }

func (num Num) String() string {
	var sym string
	if num.Head {
		sym = "-"
		if num.Body == 0 {
			sym = "!"
		}
	}
	return fmt.Sprintf("%s%d", sym, num.Body)
}

func DecodePositive(value string) (int, error) {
	i, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	if i < 0 {
		return 0, fmt.Errorf("Integer decoded may not be negative")
	}
	return i, nil
}

func DecodeNum(value string) (Num, error) {
	var head bool
	if len(value) > 0 && (value[0] == '-' || value[0] == '!') {
		head, value = true, value[1:]
	}
	body, err := DecodePositive(value)
	if err != nil {
		return Num{}, err
	}
	return Num{Head: head, Body: body}, nil
}

// ConvertNum is a schema decoder's converter into Num.
func ConvertNum(value string) reflect.Value {
	num, err := DecodeNum(value)
	if err != nil { // err is lost
		return reflect.Value{}
	}
	return reflect.ValueOf(num)
}

type Duration struct {
	D       time.Duration
	Default time.Duration
}

func (dur Duration) EncodeValues(key string, values *url.Values) error {
	if dur.D != dur.Default {
		(*values)[key] = []string{dur.String()}
	}
	return nil
}

func (dur Duration) MarshalJSON() ([]byte, error) { return json.Marshal(dur.String()) }

func (dur Duration) String() string { return flags.DurationString(dur.D) }

// ConvertDurationFunc creates a schema decoder's converter into Duration.
func ConvertDurationFunc(minperiod flags.Period) func(string) reflect.Value {
	return func(value string) reflect.Value {
		p := flags.Period{Above: &minperiod.Duration}
		if err := p.Set(value); err != nil {
			return reflect.Value{}
		}
		return reflect.ValueOf(Duration{D: p.Duration})
	}
}

func (p *Params) ResetSchema() {
	val := reflect.ValueOf(&p.Schema).Elem()
	for typ, i := val.Type(), 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		fv := val.Field(i)
		switch sf.Type.Kind() {
		case reflect.Bool:
			fv.SetBool(false)
		case reflect.Int:
			fv.SetInt(0)
		}
		switch sf.Type {
		case NumType:
			fv.Set(reflect.ValueOf(Num{}))
		case DurationType:
			fv.Set(reflect.ValueOf(Duration{}))
		}
	}
}

func (p *Params) SetDefaults(form url.Values, minperiod flags.Period) {
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
			num.DefaultHead = def.Head
			num.DefaultBody = def.Body
			if num.Head && num.Body != 0 { // all values specified, no need for defaults
				continue
			}
			tags, ok := TagsOk(sf)
			if !ok {
				continue
			}
			if _, ok := form[tags[0]]; ok { // have parameter
				continue
			}
			if !num.Head { // not allow false init value
				num.Head = def.Head
			}
			if num.Body == 0 { // not allow 0 init value
				num.Body = def.Body
			}
		case DurationType:
			dur := fv.Addr().Interface().(*Duration)
			dur.Default = minperiod.Duration
			if dur.D != time.Duration(0) { // value specified, no need for defaults
				continue
			}
			tags, ok := TagsOk(sf)
			if !ok {
				continue
			}
			if _, ok := form[tags[0]]; ok { // have parameter
				continue
			}
			dur.D = minperiod.Duration
		}
	}
}

func (p *Params) DecodeDecode(req *http.Request) error {
	if err := req.ParseForm(); err != nil {
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
	dec.RegisterConverter(Duration{}, ConvertDurationFunc(p.MinPeriod))

	p.ResetSchema()
	err := dec.Decode(&p.Schema, req.Form)
	if err != nil {
		return err
	}
	p.SetDefaults(req.Form, p.MinPeriod)
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

type MultiError map[string]error

func (e MultiError) Error() string {
	return fmt.Sprintf("%d error(s)", len(e))
}

// RenamedConstError denotes an error.
type RenamedConstError string

func (rc RenamedConstError) Error() string { return string(rc) }
