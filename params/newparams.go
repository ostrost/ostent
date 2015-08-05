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

type Params struct {
	Schema
	Errors    MultiError          `schema:"-" url:"-" json:",omitempty"`
	Ticks     map[string]Ticks    `schema:"-" url:"-" json:"-"`
	Toprows   int                 `schema:"-" url:"-" json:"-"`
	MinPeriod flags.Period        `schema:"-" url:"-" json:"-"`
	Defaults  map[interface{}]int `schema:"-" url:"-" json:"-"` // encoded in MarshalJSON
}

type Schema struct {
	Still     bool `schema:"still"     url:"still,omitempty"`
	Hidecpu   bool `schema:"hidecpu"   url:"hidecpu,omitempty"`
	Hidedf    bool `schema:"hidedf"    url:"hidedf,omitempty"`
	Hideif    bool `schema:"hideif"    url:"hideif,omitempty"`
	Hidemem   bool `schema:"hidemem"   url:"hidemem,omitempty"`
	Hideps    bool `schema:"hideps"    url:"hideps,omitempty"`
	Hideswap  bool `schema:"hideswap"  url:"hideswap,omitempty"`
	Hidevg    bool `schema:"hidevg"    url:"hidevg,omitempty"`
	Configcpu bool `schema:"configcpu" url:"configcpu,omitempty"`
	Configdf  bool `schema:"configdf"  url:"configdf,omitempty"`
	Configif  bool `schema:"configif"  url:"configif,omitempty"`
	Configmem bool `schema:"configmem" url:"configmem,omitempty"`
	Configps  bool `schema:"configps"  url:"configps,omitempty"`
	Configvg  bool `schema:"configvg"  url:"configvg,omitempty"`
	Expanddf  bool `schema:"expanddf"  url:"expanddf,omitempty"`
	Expandif  bool `schema:"expandif"  url:"expandif,omitempty"`
	Expandcpu bool `schema:"expandcpu" url:"expandcpu,omitempty"`

	// Memn int
	// Cpun int
	// Dfn int
	// Ifn int
	Psn int `schema:"psn" url:"psn,omitempty,default8"`            // limit
	Psk int `schema:"psk" url:"psk,omitempty,default1,enumerate9"` // sort, default PID
	Dfk int `schema:"dfk" url:"dfk,omitempty,default1,enumerate5"` // sort, default FS
	Dft int `schema:"dft" url:"dft,omitempty,default2,enumerate2"` // tab, default DFBYTES
	Ift int `schema:"ift" url:"ift,omitempty,default3,enumerate3"` // tab, default IFBYTES

	// The NewParams must populate .Ticks with EACH Refresh*
	Refreshcpu Duration `schema:"refreshcpu" url:"refreshcpu,omitempty"`
	Refreshdf  Duration `schema:"refreshdf"  url:"refreshdf,omitempty"`
	Refreshif  Duration `schema:"refreshif"  url:"refreshif,omitempty"`
	Refreshmem Duration `schema:"refreshmem" url:"refreshmem,omitempty"`
	Refreshps  Duration `schema:"refreshps"  url:"refreshps,omitempty"`
	Refreshvg  Duration `schema:"refreshvg"  url:"refreshvg,omitempty"`
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
		Defaults:   p.StringDefaults(),
	}
	return json.Marshal(d)
}

func (p Params) StringDefaults() map[string]int {
	m := make(map[string]int)
	for k, v := range p.Defaults {
		if s, ok := k.(string); ok {
			m[s] = v
		}
	}
	return m
}

func (p Params) Nonzero(v *int) int {
	if v != nil && *v != 0 {
		return *v
	}
	if d, ok := p.Defaults[v]; ok {
		return d
	}
	fmt.Printf("Cannot find default for %+v\n", v)
	return 0
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
		if sf.Type.Kind() != reflect.Bool {
			continue
		}
		if tags, ok := TagsOk(sf); !ok || tags[0] == "" {
			continue
		}
		b := val.Field(i).Addr().Interface().(*bool)
		if s, err := p.HrefToggle(b); err == nil {
			m[sf.Name] = s
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
	values, err := query.Values(p)
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
