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

// JSX is whether we're doing it for jsx.
var JSX bool

// classword returns either class or className depending on JSX value.
func classword() string {
	return map[bool]string{
		false: "class",     // default
		true:  "className", // jsx case
	}[JSX]
}

// forword returns either for or htmlFor depending on JSX value.
func forword() string {
	return map[bool]string{
		false: "for",     // default
		true:  "htmlFor", // jsx case
	}[JSX]
}

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

func CastError(notype string) error {
	return fmt.Errorf("Cannot convert into %s", notype)
}

func toggleHrefAttr(value interface{}) (interface{}, error) {
	if JSX {
		return fmt.Sprintf(" href={%s.Href} onClick={this.handleClick}",
			uncurl(value.(string))), nil
	}
	if bp, ok := value.(*client.BoolParam); ok {
		return template.HTMLAttr(fmt.Sprintf(" href=\"%s\"", bp.EncodeToggle())), nil
	}
	return nil, CastError("*client.BoolParam")
}

func formActionAttr(value interface{}) (interface{}, error) {
	if JSX {
		return fmt.Sprintf(" action={\"/form/\"+%s}", uncurl(value.(string))), nil
	}
	if query, ok := value.(*client.Query); ok {
		return template.HTMLAttr(fmt.Sprintf(" action=\"/form/%s\"",
			url.QueryEscape(query.ValuesEncode(nil)))), nil
	}
	return nil, CastError("*client.Query")
}

func periodNameAttr(value interface{}) (interface{}, error) {
	if JSX {
		prefix, _ := DotSplitHash(value)
		_, pname := DotSplit(prefix)
		return fmt.Sprintf(" name=%q", pname), nil
	}
	if period, ok := value.(*client.PeriodParam); ok {
		return template.HTMLAttr(fmt.Sprintf(" name=%q",
			period.Pname)), nil
	}
	return nil, CastError("*client.PeriodParam")
}

func periodValueAttr(value interface{}) (interface{}, error) {
	if JSX {
		prefix, _ := DotSplitHash(value)
		return fmt.Sprintf(" onChange={this.handleChange} value={%s.Input}", prefix), nil
	}
	if period, ok := value.(*client.PeriodParam); ok {
		if period.Input != "" {
			return template.HTMLAttr(fmt.Sprintf(" value=\"%s\"", period.Input)), nil
		}
		return template.HTMLAttr(""), nil
	}
	return nil, CastError("*client.PeriodParam")
}

func refreshClass(value interface{}, classes string) (interface{}, error) {
	if JSX {
		prefix, _ := DotSplitHash(value)
		return fmt.Sprintf(" %s={%q + (%s.InputErrd ? \" has-warning\" : \"\")}",
			classword(), classes, prefix), nil
	}
	if period, ok := value.(*client.PeriodParam); ok {
		if period.InputErrd {
			classes += " " + "has-warning"
		}
		return template.HTMLAttr(fmt.Sprintf(" %s=%q", classword(), classes)), nil
	}
	return nil, CastError("*client.PeriodParam")
}

func ifDisabledAttr(value interface{}) (template.HTMLAttr, error) {
	if JSX {
		return template.HTMLAttr(fmt.Sprintf("disabled={%s.Value ? \"disabled\" : \"\" }",
			uncurl(value.(string)))), nil
	}
	if bp, ok := value.(*client.BoolParam); ok {
		if bp.Value {
			return template.HTMLAttr("disabled=\"disabled\""), nil
		}
		return template.HTMLAttr(""), nil
	}
	return template.HTMLAttr(""), CastError("*client.BoolParam")
}

func ifClassAttr(value interface{}, classes ...string) (template.HTMLAttr, error) {
	s, err := ifClass(value, classes...)
	if err != nil {
		return template.HTMLAttr(""), err
	}
	if !JSX {
		s = fmt.Sprintf("%q", s)
	}
	return template.HTMLAttr(fmt.Sprintf(" %s=%s", classword(), s)), nil
}

func ifClass(value interface{}, classes ...string) (string, error) {
	if len(classes) == 0 || len(classes) > 3 {
		return "", fmt.Errorf("number of args for ifClass*: either 2 or 3 or 4 got %d", 1+len(classes))
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
	if JSX {
		return fmt.Sprintf("{%s.Value ? %q : %q }", uncurl(value.(string)), fstclass, sndclass), nil
	}
	if bp, ok := value.(*client.BoolParam); ok {
		if bp.Value {
			return fstclass, nil
		}
		return sndclass, nil
	}
	return "", CastError("*client.BoolParam")
}

func uncurl(s string) string {
	return strings.TrimSuffix(strings.TrimPrefix(s, "{"), "}")
}

func droplink(value interface{}, ss ...string) (interface{}, error) {
	if value == nil {
		return nil, fmt.Errorf("value supplied for droplink is nil")
	}
	var named string
	if len(ss) > 0 {
		named = ss[0]
	}
	AC := "text-right" // default
	if len(ss) > 1 {
		AC = ""
		if ss[1] != "" {
			AC = "text-" + ss[1]
		}
	}
	if JSX {
		prefix, _ := DotSplitHash(value)
		_, pname := DotSplit(prefix)
		enums := client.NewParamsENUM(nil)
		ed := enums[pname].EnumDecodec
		return client.DropLink{
			AlignClass: AC,
			Text:       ed.Text(named), // always static
			Href:       fmt.Sprintf("{%s.%s.%s}", prefix, named, "Href"),
			Class:      fmt.Sprintf("{%s.%s.%s}", prefix, named, "Class"),
			CaretClass: fmt.Sprintf("{%s.%s.%s}", prefix, named, "CaretClass"),
		}, nil
	}
	ep, ok := value.(*client.EnumParam)
	if !ok {
		return nil, CastError("*client.EnumParam")
	}
	pname, uptr := ep.EnumDecodec.Unew()
	if err := uptr.Unmarshal(named, new(bool)); err != nil {
		return nil, err
	}
	l := ep.EncodeUint(pname, uptr)
	l.AlignClass = AC
	return l, nil
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

func usepercent(val string) interface{} {
	var ca string
	if JSX {
		ca = " className={LabelClassColorPercent(" + uncurl(val) + ")}"
	} else {
		ca = fmt.Sprintf(" class=%q", LabelClassColorPercent(val))
	}
	return struct {
		Value     string
		ClassAttr template.HTMLAttr
	}{
		Value:     val,
		ClassAttr: template.HTMLAttr(ca),
	}
}

func key(prefix, val string) template.HTMLAttr {
	if !JSX {
		return template.HTMLAttr("")
	}
	return template.HTMLAttr(fmt.Sprintf(" key={%q+%s}", prefix+"-", uncurl(val)))
}

type Clipped struct {
	IDAttr      template.HTMLAttr
	ForAttr     template.HTMLAttr
	MWStyleAttr template.HTMLAttr
	Text        string
}

func clip(width int, prefix, val string, rest ...string) (*Clipped, error) {
	var key, mws string
	if JSX {
		key = fmt.Sprintf("{%q+%s}", prefix+"-", uncurl(val))
		mws = fmt.Sprintf("{{maxWidth: '%dch'}}", width)
	} else { // quote everything
		key = fmt.Sprintf("%q", url.QueryEscape(prefix+"-"+val))
		mws = fmt.Sprintf("\"max-width: %dch \"", width)
	}
	if len(rest) == 1 {
		val = rest[0]
	} else if len(rest) > 0 {
		return nil, fmt.Errorf("clip expects either 5 or 6 arguments")
	}
	return &Clipped{
		IDAttr:      template.HTMLAttr("id=" + key),
		ForAttr:     template.HTMLAttr(forword() + "=" + key),
		MWStyleAttr: template.HTMLAttr("style=" + mws),
		Text:        val,
	}, nil
}

// Funcs features functions for templates. In use in acepp and templates.
var Funcs = template.FuncMap{
	"rowsset": func(interface{}) string { return "" }, // empty pipeline
	// acepp overrides rowsset and adds setrows

	"key":             key,
	"clip":            clip,
	"droplink":        droplink,
	"usepercent":      usepercent,
	"ifClass":         ifClass,
	"ifClassAttr":     ifClassAttr,
	"ifDisabledAttr":  ifDisabledAttr,
	"toggleHrefAttr":  toggleHrefAttr,
	"formActionAttr":  formActionAttr,
	"periodNameAttr":  periodNameAttr,
	"periodValueAttr": periodValueAttr,
	"refreshClass":    refreshClass,
	"closeTag":        CloseTagFunc(nil),
	"class":           classword,
	"for":             forword,

	"json": func(v interface{}) (string, error) {
		j, err := json.Marshal(v)
		return string(j), err
	},
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
