package templatefunc

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/ostrost/ostent/params"
	"github.com/ostrost/ostent/system/operating"
	"github.com/ostrost/ostent/templateutil/templatepipe"
)

func (f JSXFuncs) classWord() string   { return "className" }
func (f JSXFuncs) colspanWord() string { return "colSpan" }
func (f JSXFuncs) forWord() string     { return "htmlFor" }

func (f HTMLFuncs) classWord() string   { return "class" }
func (f HTMLFuncs) colspanWord() string { return "colspan" }
func (f HTMLFuncs) forWord() string     { return "for" }

// jsxClose returns close tag markup as template.HTML.
func (f JSXFuncs) jsxClose(tag string) template.HTML { return template.HTML("</" + tag + ">") } // f is unused

// jsxClose returns empty template.HTML.
func (f HTMLFuncs) jsxClose(string) (empty template.HTML) { return } // f is unused

func (f JSXFuncs) ifDisabledAttr(value interface{}) (template.HTMLAttr, error) {
	return template.HTMLAttr(fmt.Sprintf("disabled={%s.Value ? \"disabled\" : \"\" }",
		f.uncurlv(value))), nil
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

func (f JSXFuncs) ifNeClassAttr(value interface{}, named string, class string) (template.HTMLAttr, error) {
	vstring := ToString(value)
	_, pname := DotSplit(vstring)
	enums := params.NewParamsENUM(nil)
	ed := enums[pname].EnumDecodec
	_, uptr := ed.Unew()
	if err := uptr.Unmarshal(named, new(bool)); err != nil {
		return template.HTMLAttr(""), err
	}
	return template.HTMLAttr(fmt.Sprintf(" %s={(%s.Uint != %d) ? %q : \"\"} data-tabid=\"%d\" data-title=%q",
		f.classWord(), vstring, uptr.Touint(), class, uptr.Touint(), ed.Text(named))), nil
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
	vstring := ToString(value)
	_, pname := DotSplit(vstring)
	enums := params.NewParamsENUM(nil)
	ed := enums[pname].EnumDecodec
	_, uptr := ed.Unew()
	if err := uptr.Unmarshal(named, new(bool)); err != nil {
		return template.HTMLAttr(""), err
	}
	return template.HTMLAttr(fmt.Sprintf(" %s={(%s.Uint == %d) ? %q : \"\"} data-tabid=\"%d\"",
		f.classWord(), vstring, uptr.Touint(), class, uptr.Touint())), nil
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
	fstclass, sndclass, err := classesChoices("ifBPClass*", classes)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("{%s.Value ? %q : %q }", f.uncurlv(value), fstclass, sndclass), nil
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
	fstclass, sndclass, err := classesChoices("ifBClass*", classes)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("{%s ? %q : %q }", f.uncurlv(value), fstclass, sndclass), nil
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
	vstring := ToString(value)
	_, pname := DotSplit(vstring)
	enums := params.NewParamsENUM(nil)
	ed := enums[pname].EnumDecodec
	return params.DropLink{
		AlignClass: aclass,
		Text:       ed.Text(named), // always static
		Href:       fmt.Sprintf("{%s.%s.%s}", vstring, named, "Href"),
		Class:      fmt.Sprintf("{%s.%s.%s}", vstring, named, "Class"),
		CaretClass: fmt.Sprintf("{%s.%s.%s}", vstring, named, "CaretClass"),
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

func (f JSXFuncs) usepercent(value interface{}) interface{} {
	ca := fmt.Sprintf(" %s={LabelClassColorPercent(%s)}", f.classWord(), f.uncurlv(value))
	return struct {
		Value     string
		ClassAttr template.HTMLAttr
	}{
		Value:     ToString(value),
		ClassAttr: template.HTMLAttr(ca),
	}
}

func (f HTMLFuncs) usepercent(value interface{}) interface{} {
	vstring := ToString(value)
	ca := fmt.Sprintf(" %s=%q", f.classWord(), LabelClassColorPercent(vstring))
	return struct {
		Value     string
		ClassAttr template.HTMLAttr
	}{
		Value:     vstring,
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
		"jsxClose":          f.jsxClose,
		"class":             f.classWord,
		"colspan":           f.colspanWord,
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
	jsxClose(string) template.HTML
	classWord() string
	colspanWord() string
	forWord() string
}

func init() {
	// check for Nota's interfaces compliance
	_ = interface {
		// operating (multiple types):
		Clip(int, string, ...operating.ToStringer) (*operating.Clipped, error)
		KeyAttr(string) template.HTMLAttr

		FormActionAttr() interface{}         // Query
		ToggleHrefAttr() interface{}         // BoolParam
		PeriodNameAttr() interface{}         // PeriodParam
		PeriodValueAttr() interface{}        // PeriodParam
		RefreshClassAttr(string) interface{} // PeriodParam
		LessHrefAttr() interface{}           // LimitParam
		MoreHrefAttr() interface{}           // LimitParam
	}(templatepipe.Nota(nil))
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

// DotSplitV calls DotSplit with value's string.
func DotSplitV(value interface{}) (string, string) {
	return DotSplit(ToString(value))
}

func ToString(value interface{}) string {
	if s, ok := value.(string); ok {
		return s
	}
	return value.(templatepipe.Nota).ToString()
}

func (f HTMLFuncs) CastError(notype string) error {
	// f is unused
	return fmt.Errorf("Cannot convert into %s", notype)
}

func (f JSXFuncs) uncurl(s string) string {
	// f is unused
	return strings.TrimSuffix(strings.TrimPrefix(s, "{"), "}")
}

func (f JSXFuncs) uncurlv(v interface{}) string {
	return f.uncurl(ToString(v))
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
// in passed interface{} (v) being a templatepipe.Nota.
// SetKFunc is used by acepp only.
func SetKFunc(k string) func(interface{}, string) interface{} {
	return func(v interface{}, n string) interface{} {
		if args := strings.Split(n, " "); len(args) > 1 {
			var list []string
			for _, arg := range args {
				list = append(list, templatepipe.Curl(arg))
			}
			v.(templatepipe.Nota)[k] = list
			return v
		}
		v.(templatepipe.Nota)[k] = templatepipe.Curl(n)
		return v
	}
}

// GetKFunc constructs a func which
// gets, deletes and returns k key
// in passed interface{} (v) being a templatepipe.Nota.
// GetKFunc is used by acepp only.
func GetKFunc(k string) func(interface{}) interface{} {
	return func(v interface{}) interface{} {
		h, ok := v.(templatepipe.Nota)
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
