package templatefunc

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/url"
	"strings"

	"github.com/ostrost/ostent/client"
	"github.com/ostrost/ostent/templateutil/templatepipe"
)

// classWord returns "className".
func (f JSXFuncs) classWord() string { return "className" } // f is unused

// classWord returns "class".
func (f HTMLFuncs) classWord() string { return "class" } // f is unused

// forWord returns "htmlFor".
func (f JSXFuncs) forWord() string { return "htmlFor" } // f is unused

// forWord returns "for".
func (f HTMLFuncs) forWord() string { return "for" } // f is unused

// CloseTagFunc constructs a func returning close tag markup unless the tag is in noclose.
func CloseTagFunc(noclose []string) func(string) template.HTML {
	return func(tag string) template.HTML {
		for _, nc := range noclose {
			if tag == nc {
				return template.HTML("")
			}
		}
		return template.HTML("</" + tag + ">")
	}
}

// closeTagFunc constructs a func returning close tag markup.
func (f JSXFuncs) closeTagFunc() func(string) template.HTML {
	// f is unused
	return CloseTagFunc(nil)
}

// closeTagFunc constructs a func returning close tag markup.
func (f HTMLFuncs) closeTagFunc() func(string) template.HTML {
	// f is unused
	return CloseTagFunc([]string{ // same set as in github.com/yosssi/ace
		"br",
		"hr",
		"img",
		"input",
		"link",
		"meta",
	})
}

func (f JSXFuncs) toggleHrefAttr(value interface{}) (interface{}, error) {
	// f is unused
	return fmt.Sprintf(" href={%s.Href} onClick={this.handleClick}",
		uncurl(value.(string))), nil
}

func (f HTMLFuncs) toggleHrefAttr(value interface{}) (interface{}, error) {
	if bp, ok := value.(*client.BoolParam); ok {
		return template.HTMLAttr(fmt.Sprintf(" href=\"%s\"", bp.EncodeToggle())), nil
	}
	return nil, f.CastError("*client.BoolParam")
}

func (f JSXFuncs) formActionAttr(value interface{}) (interface{}, error) {
	// f is unused
	return fmt.Sprintf(" action={\"/form/\"+%s}", uncurl(value.(string))), nil
}

func (f HTMLFuncs) formActionAttr(value interface{}) (interface{}, error) {
	if query, ok := value.(*client.Query); ok {
		return template.HTMLAttr(fmt.Sprintf(" action=\"/form/%s\"",
			url.QueryEscape(query.ValuesEncode(nil)))), nil
	}
	return nil, f.CastError("*client.Query")
}

func (f JSXFuncs) periodNameAttr(value interface{}) (interface{}, error) {
	// f is unused
	prefix, _ := DotSplitHash(value)
	_, pname := DotSplit(prefix)
	return fmt.Sprintf(" name=%q", pname), nil
}

func (f HTMLFuncs) periodNameAttr(value interface{}) (interface{}, error) {
	if period, ok := value.(*client.PeriodParam); ok {
		return template.HTMLAttr(fmt.Sprintf(" name=%q",
			period.Pname)), nil
	}
	return nil, f.CastError("*client.PeriodParam")
}

func (f JSXFuncs) periodValueAttr(value interface{}) (interface{}, error) {
	// f is unused
	prefix, _ := DotSplitHash(value)
	return fmt.Sprintf(" onChange={this.handleChange} value={%s.Input}", prefix), nil
}

func (f HTMLFuncs) periodValueAttr(value interface{}) (interface{}, error) {
	if period, ok := value.(*client.PeriodParam); ok {
		if period.Input != "" {
			return template.HTMLAttr(fmt.Sprintf(" value=\"%s\"", period.Input)), nil
		}
		return template.HTMLAttr(""), nil
	}
	return nil, f.CastError("*client.PeriodParam")
}

func (f JSXFuncs) refreshClass(value interface{}, classes string) (interface{}, error) {
	prefix, _ := DotSplitHash(value)
	return fmt.Sprintf(" %s={%q + (%s.InputErrd ? \" has-warning\" : \"\")}",
		f.classWord(), classes, prefix), nil
}

func (f HTMLFuncs) refreshClass(value interface{}, classes string) (interface{}, error) {
	if period, ok := value.(*client.PeriodParam); ok {
		if period.InputErrd {
			classes += " " + "has-warning"
		}
		return template.HTMLAttr(fmt.Sprintf(" %s=%q", f.classWord(), classes)), nil
	}
	return nil, f.CastError("*client.PeriodParam")
}

func (f JSXFuncs) ifDisabledAttr(value interface{}) (template.HTMLAttr, error) {
	// f is unused
	return template.HTMLAttr(fmt.Sprintf("disabled={%s.Value ? \"disabled\" : \"\" }",
		uncurl(value.(string)))), nil
}

func (f HTMLFuncs) ifDisabledAttr(value interface{}) (template.HTMLAttr, error) {
	if bp, ok := value.(*client.BoolParam); ok {
		if bp.Value {
			return template.HTMLAttr("disabled=\"disabled\""), nil
		}
		return template.HTMLAttr(""), nil
	}
	return template.HTMLAttr(""), f.CastError("*client.BoolParam")
}

func (f JSXFuncs) ifClassAttr(value interface{}, classes ...string) (template.HTMLAttr, error) {
	s, err := f.ifClass(value, classes...)
	if err != nil {
		return template.HTMLAttr(""), err
	}
	return template.HTMLAttr(fmt.Sprintf(" %s=%s", f.classWord(), s)), nil
}

func (f HTMLFuncs) ifClassAttr(value interface{}, classes ...string) (template.HTMLAttr, error) {
	s, err := f.ifClass(value, classes...)
	if err != nil {
		return template.HTMLAttr(""), err
	}
	return template.HTMLAttr(fmt.Sprintf(" %s=%q", f.classWord(), s)), nil
}

func (f JSXFuncs) ifClass(value interface{}, classes ...string) (string, error) {
	// f is unused
	fstclass, sndclass, err := classesChoices(classes)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("{%s.Value ? %q : %q }", uncurl(value.(string)), fstclass, sndclass), nil
}

func (f HTMLFuncs) ifClass(value interface{}, classes ...string) (string, error) {
	fstclass, sndclass, err := classesChoices(classes)
	if err != nil {
		return "", err
	}
	if bp, ok := value.(*client.BoolParam); ok {
		if bp.Value {
			return fstclass, nil
		}
		return sndclass, nil
	}
	return "", f.CastError("*client.BoolParam")
}

func classesChoices(classes []string) (string, string, error) {
	if len(classes) == 0 || len(classes) > 3 {
		return "", "", fmt.Errorf("number of args for ifClass*: either 2 or 3 or 4 got %d", 1+len(classes))
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
	enums := client.NewParamsENUM(nil)
	ed := enums[pname].EnumDecodec
	return client.DropLink{
		AlignClass: aclass,
		Text:       ed.Text(named), // always static
		Href:       fmt.Sprintf("{%s.%s.%s}", prefix, named, "Href"),
		Class:      fmt.Sprintf("{%s.%s.%s}", prefix, named, "Class"),
		CaretClass: fmt.Sprintf("{%s.%s.%s}", prefix, named, "CaretClass"),
	}, nil
}

func (f HTMLFuncs) droplink(value interface{}, args ...string) (interface{}, error) {
	named, aclass := DropLinkArgs(args)
	ep, ok := value.(*client.EnumParam)
	if !ok {
		return nil, f.CastError("*client.EnumParam")
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

func (f JSXFuncs) usepercent(val string) interface{} {
	ca := fmt.Sprintf(" %s={LabelClassColorPercent(%s)}", f.classWord(), uncurl(val))
	return struct {
		Value     string
		ClassAttr template.HTMLAttr
	}{
		Value:     val,
		ClassAttr: template.HTMLAttr(ca),
	}
}

func (f HTMLFuncs) usepercent(val string) interface{} {
	ca := fmt.Sprintf(" %s=%q", f.classWord(), LabelClassColorPercent(val))
	return struct {
		Value     string
		ClassAttr template.HTMLAttr
	}{
		Value:     val,
		ClassAttr: template.HTMLAttr(ca),
	}
}

func (f JSXFuncs) key(prefix, val string) template.HTMLAttr {
	// f is unused
	return template.HTMLAttr(fmt.Sprintf(" key={%q+%s}", prefix+"-", uncurl(val)))
}

func (f HTMLFuncs) key(string, string) template.HTMLAttr { return template.HTMLAttr("") } // f is unused

func SprintfAttr(format string, args ...interface{}) template.HTMLAttr {
	return template.HTMLAttr(fmt.Sprintf(format, args...))
}

type Clipped struct {
	IDAttr      template.HTMLAttr
	ForAttr     template.HTMLAttr
	MWStyleAttr template.HTMLAttr
	Text        string
}

func (f JSXFuncs) clip(width int, prefix, val string, rest ...string) (*Clipped, error) {
	text, err := ClipArgs(val, rest)
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("{%q+%s}", prefix+"-", uncurl(val))
	return &Clipped{
		IDAttr:      SprintfAttr("id=%s", key),
		ForAttr:     SprintfAttr("%s=%s", f.forWord(), key),
		MWStyleAttr: SprintfAttr("style={{maxWidth: '%dch'}}", width),
		Text:        text,
	}, nil
}
func (f HTMLFuncs) clip(width int, prefix, val string, rest ...string) (*Clipped, error) {
	text, err := ClipArgs(val, rest)
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("%q", url.QueryEscape(prefix+"-"+val))
	return &Clipped{
		IDAttr:      SprintfAttr("id=%s", key),
		ForAttr:     SprintfAttr("%s=%s", f.forWord(), key),
		MWStyleAttr: SprintfAttr("style=\"max-width: %dch \"", width),
		Text:        text,
	}, nil
}

func ClipArgs(fst string, rest []string) (string, error) {
	if len(rest) == 1 {
		return rest[0], nil
	} else if len(rest) > 0 {
		return "", fmt.Errorf("clip expects either 5 or 6 arguments")
	}
	return fst, nil
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

		"key":             f.key,
		"clip":            f.clip,
		"droplink":        f.droplink,
		"usepercent":      f.usepercent,
		"ifClass":         f.ifClass,
		"ifClassAttr":     f.ifClassAttr,
		"ifDisabledAttr":  f.ifDisabledAttr,
		"toggleHrefAttr":  f.toggleHrefAttr,
		"formActionAttr":  f.formActionAttr,
		"periodNameAttr":  f.periodNameAttr,
		"periodValueAttr": f.periodValueAttr,
		"refreshClass":    f.refreshClass,
		"closeTag":        f.closeTagFunc(),
		"class":           f.classWord,
		"for":             f.forWord,

		"json": func(v interface{}) (string, error) {
			j, err := json.Marshal(v)
			return string(j), err
		},
	}
}

// Funcs features functions for templates. In use in acepp and templates.
var Funcs = HTMLFuncs{}.MakeMap()

type Functor interface {
	MakeMap() template.FuncMap
	key(string, string) template.HTMLAttr
	clip(int, string, string, ...string) (*Clipped, error)
	droplink(interface{}, ...string) (interface{}, error)
	usepercent(string) interface{}
	ifClass(interface{}, ...string) (string, error)
	ifClassAttr(interface{}, ...string) (template.HTMLAttr, error)
	ifDisabledAttr(interface{}) (template.HTMLAttr, error)
	toggleHrefAttr(interface{}) (interface{}, error)
	formActionAttr(interface{}) (interface{}, error)
	periodNameAttr(interface{}) (interface{}, error)
	periodValueAttr(interface{}) (interface{}, error)
	refreshClass(interface{}, string) (interface{}, error)
	closeTagFunc() func(string) template.HTML
	classWord() string
	forWord() string
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
		return DotSplit(uncurl(v.(string)))
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
		v.(templatepipe.Hash)[k] = templatepipe.Curly(n)
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
		delete(h, k)
		return n // may also be empty, affects dispatch
	}
}
