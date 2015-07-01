package templatefunc

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/ostrost/ostent/params"
	"github.com/ostrost/ostent/system/operating"
	"github.com/ostrost/ostent/templateutil/templatepipe"
)

// colspanWord returns "colSpan".
func (f JSXFuncs) colspanWord() string { return "colSpan" } // f is unused

// colspanWord returns "colspan".
func (f HTMLFuncs) colspanWord() string { return "colspan" } // f is unused

// classWord returns "className".
func (f JSXFuncs) classWord() string { return "className" } // f is unused

// classWord returns "class".
func (f HTMLFuncs) classWord() string { return "class" } // f is unused

// forWord returns "htmlFor".
func (f JSXFuncs) forWord() string { return "htmlFor" } // f is unused

// forWord returns "for".
func (f HTMLFuncs) forWord() string { return "for" } // f is unused

// jsxClose returns close tag markup as template.HTML.
func (f JSXFuncs) jsxClose(tag string) template.HTML { return template.HTML("</" + tag + ">") } // f is unused

// jsxClose returns empty template.HTML.
func (f HTMLFuncs) jsxClose(string) (empty template.HTML) { return } // f is unused

func (f JSXFuncs) toggleHrefAttr(value interface{}) (interface{}, error) {
	// f is unused
	return fmt.Sprintf(" href={%s.Href} onClick={this.handleClick}", uncurlv(value)), nil
}

func (f HTMLFuncs) toggleHrefAttr(value interface{}) (interface{}, error) {
	if bp, ok := value.(*params.BoolParam); ok {
		return template.HTMLAttr(fmt.Sprintf(" href=\"%s\"", bp.EncodeToggle())), nil
	}
	return nil, f.CastError("*params.BoolParam")
}

func (f JSXFuncs) periodNameAttr(value interface{}) (interface{}, error) {
	// f is unused
	prefix, _ := DotSplitHash(value)
	_, pname := DotSplit(prefix)
	return fmt.Sprintf(" name=%q", pname), nil
}

func (f HTMLFuncs) periodNameAttr(value interface{}) (interface{}, error) {
	if period, ok := value.(*params.PeriodParam); ok {
		return template.HTMLAttr(fmt.Sprintf(" name=%q",
			period.Pname)), nil
	}
	return nil, f.CastError("*params.PeriodParam")
}

func (f JSXFuncs) periodValueAttr(value interface{}) (interface{}, error) {
	// f is unused
	prefix, _ := DotSplitHash(value)
	return fmt.Sprintf(" onChange={this.handleChange} value={%s.Input}", prefix), nil
}

func (f HTMLFuncs) periodValueAttr(value interface{}) (interface{}, error) {
	if period, ok := value.(*params.PeriodParam); ok {
		if period.Input != "" {
			return template.HTMLAttr(fmt.Sprintf(" value=\"%s\"", period.Input)), nil
		}
		return template.HTMLAttr(""), nil
	}
	return nil, f.CastError("*params.PeriodParam")
}

func (f JSXFuncs) refreshClass(value interface{}, classes string) (interface{}, error) {
	prefix, _ := DotSplitHash(value)
	return fmt.Sprintf(" %s={%q + (%s.InputErrd ? \" has-warning\" : \"\")}",
		f.classWord(), classes, prefix), nil
}

func (f HTMLFuncs) refreshClass(value interface{}, classes string) (interface{}, error) {
	if period, ok := value.(*params.PeriodParam); ok {
		if period.InputErrd {
			classes += " " + "has-warning"
		}
		return template.HTMLAttr(fmt.Sprintf(" %s=%q", f.classWord(), classes)), nil
	}
	return nil, f.CastError("*params.PeriodParam")
}

func (f JSXFuncs) lessHrefAttr(value interface{}) (interface{}, error) {
	// f is unused
	return fmt.Sprintf(" href={%s.LessHref} onClick={this.handleClick}", uncurlv(value)), nil
}
func (f HTMLFuncs) lessHrefAttr(value interface{}) (interface{}, error) {
	if lp, ok := value.(*params.LimitParam); ok {
		return template.HTMLAttr(fmt.Sprintf(" href=\"%s\"", lp.EncodeLess())), nil
	}
	return nil, f.CastError("*params.LimitParam")
}

func (f JSXFuncs) moreHrefAttr(value interface{}) (interface{}, error) {
	// f is unused
	return fmt.Sprintf(" href={%s.MoreHref} onClick={this.handleClick}", uncurlv(value)), nil
}
func (f HTMLFuncs) moreHrefAttr(value interface{}) (interface{}, error) {
	if lp, ok := value.(*params.LimitParam); ok {
		return template.HTMLAttr(fmt.Sprintf(" href=\"%s\"", lp.EncodeMore())), nil
	}
	return nil, f.CastError("*params.LimitParam")
}

func (f JSXFuncs) ifDisabledAttr(value interface{}) (template.HTMLAttr, error) {
	// f is unused
	return template.HTMLAttr(fmt.Sprintf("disabled={%s.Value ? \"disabled\" : \"\" }",
		uncurlv(value))), nil
}

func (f HTMLFuncs) ifDisabledAttr(value interface{}) (template.HTMLAttr, error) {
	if bp, ok := value.(*params.BoolParam); ok {
		if bp.Value {
			return template.HTMLAttr("disabled=\"disabled\""), nil
		}
		return template.HTMLAttr(""), nil
	}
	return template.HTMLAttr(""), f.CastError("*params.BoolParam")
}

func (f JSXFuncs) ifBPClassAttr(value interface{}, classes ...string) (template.HTMLAttr, error) {
	s, err := f.ifBPClass(value, classes...)
	if err != nil {
		return template.HTMLAttr(""), err
	}
	return template.HTMLAttr(fmt.Sprintf(" %s=%s", f.classWord(), s)), nil
}

func (f HTMLFuncs) ifBPClassAttr(value interface{}, classes ...string) (template.HTMLAttr, error) {
	s, err := f.ifBPClass(value, classes...)
	if err != nil {
		return template.HTMLAttr(""), err
	}
	return template.HTMLAttr(fmt.Sprintf(" %s=%q", f.classWord(), s)), nil
}

/*
func (_ HashValue) ClassAttrUnless(dot interface{}, cmp uint, class string) (template.HTMLAttr, error) {}
func (_ Uint) ClassAttrUnless(dot interface{}, cmp uint, class string) (template.HTMLAttr, error) {
	if z, ok := dot.(HashValue); ok { return z.ClassAttrUnless(dot, cmp, class) }
} // */

func (f JSXFuncs) ifNeClassAttr(value interface{}, named string, class string) (template.HTMLAttr, error) {
	prefix, _ := DotSplitHash(value) // "Data.Links.Params.ENUM.ift" ?
	_, pname := DotSplit(prefix)
	enums := params.NewParamsENUM(nil)
	ed := enums[pname].EnumDecodec
	_, uptr := ed.Unew()
	if err := uptr.Unmarshal(named, new(bool)); err != nil {
		return template.HTMLAttr(""), err
	}
	return template.HTMLAttr(fmt.Sprintf(" %s={(%s.Uint != %d) ? %q : \"\"} data-tabid=\"%d\" data-title=%q",
		f.classWord(), prefix, uptr.Touint(), class, uptr.Touint(), ed.Text(named))), nil
}

func (f HTMLFuncs) ifNeClassAttr(value interface{}, named string, class string) (template.HTMLAttr, error) {
	ep, ok := value.(*params.EnumParam)
	if !ok {
		return template.HTMLAttr(""), f.CastError("*params.EnumParams")
	}
	_, uptr := ep.EnumDecodec.Unew()
	if err := uptr.Unmarshal(named, new(bool)); err != nil {
		return template.HTMLAttr(""), err
	}
	ucmp := uptr.Touint()
	if ep.Number.Uint == ucmp {
		class = ""
	}
	return template.HTMLAttr(fmt.Sprintf(" %s=%q data-tabid=\"%d\" data-title=%q",
		f.classWord(), class, ucmp, ep.EnumDecodec.Text(named))), nil
}

func (f JSXFuncs) iftEnumAttrs(value interface{}, named string, class string) (template.HTMLAttr, error) {
	prefix, _ := DotSplitHash(value) // "Data.Links.Params.ENUM.ift" ?
	_, pname := DotSplit(prefix)
	enums := params.NewParamsENUM(nil)
	ed := enums[pname].EnumDecodec
	_, uptr := ed.Unew()
	if err := uptr.Unmarshal(named, new(bool)); err != nil {
		return template.HTMLAttr(""), err
	}
	return template.HTMLAttr(fmt.Sprintf(" %s={(%s.Uint == %d) ? %q : \"\"} data-tabid=\"%d\"",
		f.classWord(), prefix, uptr.Touint(), class, uptr.Touint())), nil
}

func (f HTMLFuncs) iftEnumAttrs(value interface{}, named string, class string) (template.HTMLAttr, error) {
	ep, ok := value.(*params.EnumParam)
	if !ok {
		return template.HTMLAttr(""), f.CastError("*params.EnumParams")
	}
	_, uptr := ep.EnumDecodec.Unew()
	if err := uptr.Unmarshal(named, new(bool)); err != nil {
		return template.HTMLAttr(""), err
	}
	ucmp := uptr.Touint()
	if ep.Number.Uint != ucmp {
		class = ""
	}
	return template.HTMLAttr(fmt.Sprintf(" %s=%q data-tabid=\"%d\"", f.classWord(), class, ucmp)), nil
}

func (f JSXFuncs) ifExpandClassAttr(value interface{}, classes ...string) (template.HTMLAttr, error) {
	return "", fmt.Errorf("Not implemented yet")
}
func (f HTMLFuncs) ifExpandClassAttr(value interface{}, classes ...string) (template.HTMLAttr, error) {
	// if ei, ok := value.(*params.ExpandInfo); ok {}
	return "", f.CastError("*ExpandInfo")
}

func (f JSXFuncs) ifBClassAttr(value interface{}, classes ...string) (template.HTMLAttr, error) {
	s, err := f.ifBClass(value, classes...)
	if err != nil {
		return template.HTMLAttr(""), err
	}
	return template.HTMLAttr(fmt.Sprintf(" %s=%s", f.classWord(), s)), nil
}

func (f HTMLFuncs) ifBClassAttr(value interface{}, classes ...string) (template.HTMLAttr, error) {
	s, err := f.ifBClass(value, classes...)
	if err != nil {
		return template.HTMLAttr(""), err
	}
	return template.HTMLAttr(fmt.Sprintf(" %s=%q", f.classWord(), s)), nil
}

func (f JSXFuncs) ifBPClass(value interface{}, classes ...string) (string, error) {
	// f is unused
	fstclass, sndclass, err := classesChoices("ifBPClass*", classes)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("{%s.Value ? %q : %q }", uncurlv(value), fstclass, sndclass), nil
}

func (f HTMLFuncs) ifBPClass(value interface{}, classes ...string) (string, error) {
	fstclass, sndclass, err := classesChoices("ifBPClass*", classes)
	if err != nil {
		return "", err
	}
	if bp, ok := value.(*params.BoolParam); ok {
		if bp.Value {
			return fstclass, nil
		}
		return sndclass, nil
	}
	return "", f.CastError("*bool")
}

func (f JSXFuncs) ifBClass(value interface{}, classes ...string) (string, error) {
	// f is unused
	fstclass, sndclass, err := classesChoices("ifBClass*", classes)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("{%s ? %q : %q }", uncurlv(value), fstclass, sndclass), nil
}

func (f HTMLFuncs) ifBClass(value interface{}, classes ...string) (string, error) {
	fstclass, sndclass, err := classesChoices("ifBClass*", classes)
	if err != nil {
		return "", err
	}
	if bp, ok := value.(*bool); ok {
		if bp != nil && *bp {
			return fstclass, nil
		}
		return sndclass, nil
	}
	return "", f.CastError("*bool")
}

func classesChoices(caller string, classes []string) (string, string, error) {
	if len(classes) == 0 || len(classes) > 3 {
		return "", "", fmt.Errorf("number of args for %s: either 2 or 3 or 4 got %d", caller, 1+len(classes))
	}
	sndclass := ""
	if len(classes) > 1 {
		sndclass = classes[1]
	}
	fstclass := classes[0]
	if len(classes) > 2 {
		fstclass = classes[2] + " " + fstclass
		sndclass = classes[2] + " " + sndclass
	}
	return fstclass, sndclass, nil
}

func (f JSXFuncs) droplink(value interface{}, args ...string) (interface{}, error) {
	// f is unused
	named, aclass := DropLinkArgs(args)
	prefix, _ := DotSplitHash(value)
	_, pname := DotSplit(prefix)
	enums := params.NewParamsENUM(nil)
	ed := enums[pname].EnumDecodec
	return params.DropLink{
		AlignClass: aclass,
		Text:       ed.Text(named), // always static
		Href:       fmt.Sprintf("{%s.%s.%s}", prefix, named, "Href"),
		Class:      fmt.Sprintf("{%s.%s.%s}", prefix, named, "Class"),
		CaretClass: fmt.Sprintf("{%s.%s.%s}", prefix, named, "CaretClass"),
	}, nil
}

func (f HTMLFuncs) droplink(value interface{}, args ...string) (interface{}, error) {
	named, aclass := DropLinkArgs(args)
	ep, ok := value.(*params.EnumParam)
	if !ok {
		return nil, f.CastError("*params.EnumParam")
	}
	pname, uptr := ep.EnumDecodec.Unew()
	if err := uptr.Unmarshal(named, new(bool)); err != nil {
		return nil, err
	}
	l := ep.EncodeUint(pname, uptr)
	l.AlignClass = aclass
	return l, nil
}

func DropLinkArgs(args []string) (string, string) {
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

func (f JSXFuncs) usepercent(val interface{}) interface{} {
	ca := fmt.Sprintf(" %s={LabelClassColorPercent(%s)}", f.classWord(), uncurlv(val))
	return struct {
		Value     string
		ClassAttr template.HTMLAttr
	}{
		Value:     stringv(val),
		ClassAttr: template.HTMLAttr(ca),
	}
}

func (f HTMLFuncs) usepercent(val interface{}) interface{} {
	ca := fmt.Sprintf(" %s=%q", f.classWord(), LabelClassColorPercent(stringi(val)))
	return struct {
		Value     string
		ClassAttr template.HTMLAttr
	}{
		Value:     stringi(val),
		ClassAttr: template.HTMLAttr(ca),
	}
}

// JSXFuncs has methods implementing Functor.
type JSXFuncs struct{}

// HTMLFuncs has methods implementing Functor.
type HTMLFuncs struct{}

// MakeMap is dull but required.
func (f JSXFuncs) MakeMap() template.FuncMap { return MakeMap(f) }

// MakeMap is dull but required.
func (f HTMLFuncs) MakeMap() template.FuncMap { return MakeMap(f) }

// MakeMap constructs template.FuncMap off f implementation.
func MakeMap(f Functor) template.FuncMap {
	return template.FuncMap{
		"rowsset": func(interface{}) string { return "" }, // empty pipeline
		// acepp overrides rowsset and adds setrows

		"droplink":          f.droplink,
		"usepercent":        f.usepercent,
		"ifBClass":          f.ifBClass,
		"ifBClassAttr":      f.ifBClassAttr,
		"ifBPClass":         f.ifBPClass,
		"ifBPClassAttr":     f.ifBPClassAttr,
		"ifNeClassAttr":     f.ifNeClassAttr,
		"iftEnumAttrs":      f.iftEnumAttrs,
		"ifExpandClassAttr": f.ifExpandClassAttr,
		"ifDisabledAttr":    f.ifDisabledAttr,
		"toggleHrefAttr":    f.toggleHrefAttr,
		"periodNameAttr":    f.periodNameAttr,
		"periodValueAttr":   f.periodValueAttr,
		"refreshClass":      f.refreshClass,
		"lessHrefAttr":      f.lessHrefAttr,
		"moreHrefAttr":      f.moreHrefAttr,
		"jsxClose":          f.jsxClose,
		"colspan":           f.colspanWord,
		"class":             f.classWord,
		"for":               f.forWord,
	}
}

// Funcs features functions for templates. In use in acepp and templates.
var Funcs = HTMLFuncs{}.MakeMap()

type Functor interface {
	MakeMap() template.FuncMap
	droplink(interface{}, ...string) (interface{}, error)
	usepercent(interface{}) interface{}
	ifBClass(interface{}, ...string) (string, error)
	ifBClassAttr(interface{}, ...string) (template.HTMLAttr, error)
	ifBPClass(interface{}, ...string) (string, error)
	ifBPClassAttr(interface{}, ...string) (template.HTMLAttr, error)
	ifNeClassAttr(interface{}, string, string) (template.HTMLAttr, error)
	iftEnumAttrs(interface{}, string, string) (template.HTMLAttr, error)
	ifExpandClassAttr(interface{}, ...string) (template.HTMLAttr, error)
	ifDisabledAttr(interface{}) (template.HTMLAttr, error)
	toggleHrefAttr(interface{}) (interface{}, error)
	periodNameAttr(interface{}) (interface{}, error)
	periodValueAttr(interface{}) (interface{}, error)
	refreshClass(interface{}, string) (interface{}, error)
	lessHrefAttr(interface{}) (interface{}, error)
	moreHrefAttr(interface{}) (interface{}, error)
	jsxClose(string) template.HTML
	colspanWord() string
	classWord() string
	forWord() string
}

type FormActioner interface {
	FormActionAttr() (interface{}, error)
}
type Keyer interface {
	KeyAttr(string) template.HTMLAttr
}
type Clipper interface {
	Clip(int, string, ...operating.ToStringer) (*operating.Clipped, error)
}

func init() {
	v := templatepipe.Value("")
	// check for Value's interfaces compliance
	_ = FormActioner(v)
	_ = Keyer(v)
	_ = Clipper(v)
}

// DotSplit splits s by last ".".
func DotSplit(s string) (string, string) {
	if s == "" {
		return "", ""
	}
	i := len(s) - 1
	for i > 0 && s[i] != '.' {
		i--
	}
	return s[:i], s[i+1:]
}

// DotSplitHash returns DotSplit of first (in no particular order)
// value from value being a templatepipe.Hash.
func DotSplitHash(value interface{}) (string, string) {
	for _, v := range value.(templatepipe.Hash) {
		return DotSplit(uncurlv(v))
	}
	return "", ""
}

func (f HTMLFuncs) CastError(notype string) error {
	// f is unused
	return fmt.Errorf("Cannot convert into %s", notype)
}

func uncurl(s string) string {
	return strings.TrimSuffix(strings.TrimPrefix(s, "{"), "}")
}

func uncurlv(v interface{}) string {
	return uncurl(stringv(v))
}

func stringv(v interface{}) string {
	return string(v.(templatepipe.Value))
}

func stringi(v interface{}) string {
	return v.(string)
}

func LabelClassColorPercent(p string) string {
	if len(p) > 2 { // 100% and more
		return "label label-danger"
	}
	if len(p) > 1 {
		if p[0] == '9' {
			return "label label-danger"
		}
		if p[0] == '8' {
			return "label label-warning"
		}
		if p[0] == '1' {
			return "label label-success"
		}
		return "label label-info"
	}
	return "label label-success"
}

// SetKFunc constructs a func which
// sets k key to templatepipe.Curly(string (n))
// in passed interface{} (v) being a templatepipe.Hash.
// SetKFunc is used by acepp only.
func SetKFunc(k string) func(interface{}, string) interface{} {
	return func(v interface{}, n string) interface{} {
		if args := strings.Split(n, " "); len(args) > 1 {
			var list []string
			for _, arg := range args {
				list = append(list, templatepipe.Curl(arg))
			}
			v.(templatepipe.Hash)[k] = list
			return v
		}
		v.(templatepipe.Hash)[k] = templatepipe.Curl(n)
		return v
	}
}

// GetKFunc constructs a func which
// gets, deletes and returns k key
// in passed interface{} (v) being a templatepipe.Hash.
// GetKFunc is used by acepp only.
func GetKFunc(k string) func(interface{}) interface{} {
	return func(v interface{}) interface{} {
		h, ok := v.(templatepipe.Hash)
		if !ok {
			return "" // empty pipeline, affects dispatch
		}
		n := h[k]
		if args, ok := n.([]string); ok {
			if len(args) > 1 {
				h[k] = args[1:]
			} else {
				delete(h, k)
			}
			return args[0]
		}
		delete(h, k)
		return n // may also be empty, affects dispatch
	}
}
