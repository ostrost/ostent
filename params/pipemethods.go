package params

import (
	"fmt"
	"html/template"
	"math"
	"net/url"
	"reflect"
	"strings"

	"github.com/google/go-querystring/query"
)

func SprintfAttr(format string, args ...interface{}) template.HTMLAttr {
	return template.HTMLAttr(fmt.Sprintf(format, args...))
}

// AttrActionForm is for template.
func (p Params) AttrActionForm() (interface{}, error) {
	s, err := p.Encode()
	if err != nil {
		return nil, err
	}
	return SprintfAttr(" action=\"/form/%s\"", url.QueryEscape(s)), nil
}

func (p Params) AttrClassN(b bool, fstclass, sndclass string) template.HTMLAttr {
	// p is unused
	s := fstclass
	if !b {
		s = sndclass
	}
	return SprintfAttr(" class=%q", s)
}

/*
func (bp BoolParam) DisabledAttr() interface{} {
	if !bp.Value {
		return template.HTMLAttr("")
	}
	return SprintfAttr(" disabled=%q", "disabled")
}

func (bp BoolParam) ToggleHrefAttr() interface{} {
	return SprintfAttr(" href=\"%s\"", bp.EncodeToggle())
}
*/

func (p *Params) AttrClassParamsError(m MultiError, name, fstclass, sndclass string) template.HTMLAttr {
	_, ok := m[name]
	return p.AttrClassN(ok, fstclass, sndclass)
}

func (p *Params) AttrClassT(defaults map[interface{}]int, vp *int, cmp int, fstclass, sndclass string) template.HTMLAttr {
	v := *vp
	if v == 0 {
		if d, ok := defaults[vp]; ok {
			v = d
		}
	}
	return p.AttrClassN(v == cmp, fstclass, sndclass)
}

// p is a pointer to flip (twice) the b.
func (p *Params) HrefToggle(b *bool) (string, error) {
	*b = !*b
	s, err := p.Encode()
	*b = !*b
	return "?" + s, err
}

// p is a pointer to flip (twice) the b.
func (p *Params) AttrHrefToggle(b *bool) (interface{}, error) {
	s, err := p.HrefToggle(b)
	return SprintfAttr(" href=%q", s), err
}

// TODO In Decoder, have a cache of type=>fieldName=>tag(got,splitted)
func (p Params) AttrNameRefresh(fieldName string) (interface{}, error) {
	field, ok := reflect.TypeOf(p).FieldByName(fieldName)
	if !ok {
		return nil, fmt.Errorf("Params has no field %q", fieldName)
	}
	tag := strings.Split(field.Tag.Get("schema"), ",")[0]
	return SprintfAttr(" name=%q", tag), nil
}

func (p Params) AttrValueRefresh(fieldName string) (interface{}, error) {
	field, ok := reflect.TypeOf(p).FieldByName(fieldName)
	if !ok {
		return nil, fmt.Errorf("Params has no field %q", fieldName)
	}
	tag := strings.Split(field.Tag.Get("schema"), ",")[0]
	values, err := query.Values(p)
	if err != nil {
		return nil, err
	}
	v, ok := values[tag]
	if !ok || len(v) == 0 || v[0] == "" {
		return template.HTMLAttr(""), nil
	}
	return SprintfAttr(" value=%q", v[0]), nil
}

/*
func (ep EnumParam) EnumClassAttr(named, classif string, optelse ...string) (template.HTMLAttr, error) {
	classelse, err := EnumClassAttrArgs(optelse)
	if err != nil {
		return template.HTMLAttr(""), err
	}
	_, uptr := ep.EnumDecodec.Unew()
	if err := uptr.Unmarshal(named, new(bool)); err != nil {
		return template.HTMLAttr(""), err
	}
	if ep.Number.Uint != uptr.Touint() {
		classif = classelse
	}
	return SprintfAttr(" class=%q", classif), nil
}

func EnumClassAttrArgs(opt []string) (string, error) {
	if len(opt) == 1 {
		return opt[0], nil
	} else if len(opt) > 1 {
		return "", fmt.Errorf("number of args for EnumClassAttr: either 2 or 3 got %d",
			2+len(opt))
	}
	return "", nil
}
// */

type Varlink struct {
	AlignClass string
	CaretClass string
	LinkClass  string
	LinkHref   string
	LinkText   string `json:"-"` // static
}

func (p *Params) Variate(this *int, cmp int, text, alignClass string) (Varlink, error) {
	i := p.Nonzero(this)
	vl := Varlink{LinkText: text, LinkClass: "state"}
	if i == cmp || i == -cmp {
		vl.CaretClass = "caret"
		vl.LinkClass += " current"
		if i == cmp {
			cmp = -cmp
			vl.LinkClass += " dropup"
		}
	}
	copy := *this
	*this = cmp // set
	s, err := p.Encode()
	*this = copy // revert
	if err != nil {
		return Varlink{}, err
	}
	vl.LinkHref = "?" + s
	return vl, nil
}

/*
func (ep EnumParam) EnumLink(args ...string) (EnumLink, error) {
	named, aclass := EnumLinkArgs(args)
	pname, uptr := ep.EnumDecodec.Unew()
	if err := uptr.Unmarshal(named, new(bool)); err != nil {
		return EnumLink{}, err
	}
	l := ep.EncodeUint(pname, uptr)
	l.AlignClass = aclass
	return l, nil
}

func EnumLinkArgs(args []string) (string, string) {
	var named string
	if len(args) > 0 {
		named = args[0]
	}
	aclass := "text-right" // default
	if len(args) > 1 {
		aclass = ""
		if args[1] != "" {
			aclass = "text-" + args[1]
		}
	}
	return named, aclass
}

func (pp PeriodParam) PeriodNameAttr() interface{} {
	return SprintfAttr(" name=%q", pp.Pname)
}

func (pp PeriodParam) PeriodValueAttr() interface{} {
	if pp.Input == "" {
		return template.HTMLAttr("")
	}
	return SprintfAttr(" value=\"%s\"", pp.Input)
}

func (pp PeriodParam) RefreshClassAttr(classes string) interface{} {
	if pp.InputErrd {
		classes += " has-warning"
	}
	return SprintfAttr(" class=%q", classes)
}

func (lp LimitParam) LessHrefAttr() interface{} {
	return SprintfAttr(" href=%q", lp.EncodeLess())
}

func (lp LimitParam) MoreHrefAttr() interface{} {
	return SprintfAttr(" href=%q", lp.EncodeMore())
}
*/

func (p *Params) AttrHrefLess(value *int) (template.HTMLAttr, error) {
	old := *value
	if *value < 0 {
		*value = -*value
	}
	if *value >= 2 {
		g := math.Log2(float64(*value))
		n := math.Floor(g)
		if n == g {
			n--
		}
		*value = int(math.Pow(2, n))
	}
	s, err := p.Encode()
	*value = old
	return SprintfAttr(" href=%q", s), err
}

func (p *Params) AttrHrefMore(value *int) (template.HTMLAttr, error) {
	old := *value
	if *value < 0 {
		*value = -*value
	}
	if *value <= 32768 { // up to 65536
		*value = int(math.Pow(2, 1+math.Floor(math.Log2(float64(*value)))))
	}
	s, err := p.Encode()
	*value = old
	return SprintfAttr(" href=%q", s), err
}
