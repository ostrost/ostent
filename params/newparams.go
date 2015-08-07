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

var DurationKind = reflect.TypeOf(Duration(0)).Kind()

// NewParams constructs new Params.
func NewParams(minperiod flags.Period) *Params {
	p := &Params{
		Defaults:  make(map[interface{}]int),
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
		if k := sf.Type.Kind(); k == reflect.Int {
			if d := TagPrefixedInt(tags[1:], "default"); d != 0 {
				v := fv.Addr().Interface()
				p.Defaults[sf.Name] = d // ?
				p.Defaults[v] = d
			}
		}
		if sf.Type.Kind() == DurationKind {
			fv.Set(reflect.ValueOf(Duration(0)))
			p.Ticks[tags[0]] = NewTicks(fv.Addr().Interface().(*Duration))
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
	if ti.Ticks-1 >= int(time.Duration(*ti.Duration)/time.Second) {
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

type Defaults map[interface{}]int

type Params struct {
	Schema
	Defaults  `json:"-"`       // encoded in MarshalJSON
	Errors    MultiError       `json:",omitempty"`
	Ticks     map[string]Ticks `json:"-"`
	Toprows   int              `json:"-"`
	MinPeriod flags.Period     `json:"-"`
}

type Schema struct {
	Still     bool `url:"still,omitempty"`
	Hidecpu   bool `url:"hidecpu,omitempty"`
	Hidedf    bool `url:"hidedf,omitempty"`
	Hideif    bool `url:"hideif,omitempty"`
	Hidemem   bool `url:"hidemem,omitempty"`
	Hideps    bool `url:"hideps,omitempty"`
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

	// Memn int
	// Cpun int
	// Dfn int
	// Ifn int

	// Psn encodes number of proccesses and ps config toggle.
	// Negative value states config displaying and
	// the absolute value still encodes the ps number.
	Psn int `url:"psn,omitempty,default8"`            // limit
	Psk int `url:"psk,omitempty,default1,enumerate9"` // sort, default PID
	Dfk int `url:"dfk,omitempty,default1,enumerate5"` // sort, default FS
	Dft int `url:"dft,omitempty,default2,enumerate2"` // tab, default DFBYTES
	Ift int `url:"ift,omitempty,default3,enumerate3"` // tab, default IFBYTES

	// The NewParams must populate .Ticks with EACH Refresh*
	Refreshcpu Duration `url:"refreshcpu,omitempty"`
	Refreshdf  Duration `url:"refreshdf,omitempty"`
	Refreshif  Duration `url:"refreshif,omitempty"`
	Refreshmem Duration `url:"refreshmem,omitempty"`
	Refreshps  Duration `url:"refreshps,omitempty"`
	Refreshvg  Duration `url:"refreshvg,omitempty"`
}

func (p *Params) MarshalJSON() ([]byte, error) {
	d := struct {
		Schema
		Toggle     map[string]string
		Variations map[string][]Varlink
		Defaults   map[string]int
	}{
		Schema:     p.Schema,
		Toggle:     p.Toggles(),
		Variations: p.Variations(),
		Defaults:   p.Defaults.StringKeysOnly(),
	}
	return json.Marshal(d)
}

func (def Defaults) StringKeysOnly() map[string]int {
	m := make(map[string]int)
	for k, v := range def {
		if s, ok := k.(string); ok {
			m[s] = v
		}
	}
	return m
}

func (def Defaults) Nonzero(i *int) (int, error) {
	if i != nil && *i != 0 {
		return *i, nil
	}
	if v, ok := def[i]; ok {
		return v, nil
	}
	return 0, fmt.Errorf("Cannot find default for %+v\n", i)
}

func (def Defaults) ZeroForDefault(i *int) int {
	if v, ok := def[i]; ok && *i == v {
		return 0
	}
	return *i
}

func (p *Params) Variations() map[string][]Varlink {
	m := make(map[string][]Varlink)
	val := reflect.ValueOf(&p.Schema).Elem()
	for typ, i := val.Type(), 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		if kind := sf.Type.Kind(); kind != reflect.Int {
			continue
		}
		tags, ok := TagsOk(sf)
		if !ok || tags[0] == "" {
			continue
		}
		v := val.Field(i).Addr().Interface()
		var links []Varlink
		if vv, ok := v.(*int); ok { // better be
			max := TagPrefixedInt(tags[1:], "enumerate")
			for j := 1; j < max+1; j++ { // indexed from 1
				if vl, err := p.Variate(vv, j, "", ""); err == nil {
					links = append(links, vl)
				} // err ignored
			}
		}
		m[sf.Name] = links
	}
	return m
}

func TagPrefixedInt(words []string, prefix string) int {
	for _, w := range words {
		if strings.HasPrefix(w, prefix) {
			if i64, err := strconv.ParseInt(w[len(prefix):], 10, 0); err == nil {
				return int(i64)
			}
		}
	}
	return 0
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
		if sf.Type.Kind() == reflect.Int {
			v := val.Field(i).Addr().Interface().(*int)
			if s, err := p.HrefToggleN(v); err == nil {
				m[sf.Name] = s
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

type Duration time.Duration

func (d Duration) EncodeValues(key string, v *url.Values) error {
	if d != 0 {
		(*v)[key] = []string{time.Duration(d).String()}
	}
	return nil
}

func (d Duration) MarshalJSON() ([]byte, error) {
	if d != 0 {
		return json.Marshal(d)
	}
	return json.Marshal(nil)
}

func (d Duration) String() string {
	if d != 0 {
		return time.Duration(d).String()
	}
	return ""
}

// ConvertDurationFunc creates a schema decoder's converter into Duration.
func ConvertDurationFunc(minperiod flags.Period) func(string) reflect.Value {
	return func(value string) reflect.Value {
		p := flags.Period{Above: &minperiod.Duration}
		if err := p.Set(value); err != nil {
			return reflect.Value{}
		}
		return reflect.ValueOf(Duration(p.Duration))
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
		case DurationKind:
			fv.Set(reflect.ValueOf(Duration(0)))
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
	dec.RegisterConverter(Duration(time.Second), ConvertDurationFunc(p.MinPeriod))

	p.ResetSchema()
	derr := dec.Decode(&p.Schema, req.Form)
	if !moved || derr != nil {
		return derr
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
