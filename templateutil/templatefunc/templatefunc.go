package templatefunc

import (
	"bytes"
	"encoding/json"
	"fmt"
	templatehtml "html/template"
	"net/url"
	"strings"
	templatetext "text/template"

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

func CloseTagFunc(noclose []string) func(string) templatehtml.HTML {
	return func(tn string) templatehtml.HTML {
		for _, nc := range noclose {
			if tn == nc {
				return templatehtml.HTML("")
			}
		}
		return templatehtml.HTML("</" + tn + ">")
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

// DotSplitHash returns DotSplit of first (in no particular order) value from h.
func DotSplitHash(i interface{}) (string, string) {
	var curled string
	for _, v := range i.(templatepipe.Hash) {
		curled = v.(string)
		break
		// First (no particular order) value is fine.
	}
	return DotSplit(uncurl(curled))
}

func toggleHrefAttr(value interface{}) interface{} {
	if JSX {
		return fmt.Sprintf(" href={%s.Href} onClick={this.handleClick}", uncurl(value.(string)))
	}
	return templatehtml.HTMLAttr(fmt.Sprintf(" href=\"%s\"",
		value.(*client.BoolParam).EncodeToggle()))
}

func formActionAttr(query interface{}) interface{} {
	if JSX {
		return fmt.Sprintf(" action={\"/form/\"+%s}", uncurl(query.(string)))
	}
	return templatehtml.HTMLAttr(fmt.Sprintf(" action=\"/form/%s\"",
		url.QueryEscape(query.(*client.Query).ValuesEncode(nil))))
}

func periodNameAttr(pparam interface{}) interface{} {
	if JSX {
		prefix, _ := DotSplitHash(pparam)
		_, pname := DotSplit(prefix)
		return fmt.Sprintf(" name=%q", pname)
	}
	period := pparam.(*client.PeriodParam)
	return templatehtml.HTMLAttr(fmt.Sprintf(" name=%q", period.Pname))
}

func periodValueAttr(pparam interface{}) interface{} {
	if JSX {
		prefix, _ := DotSplitHash(pparam)
		return fmt.Sprintf(" onChange={this.handleChange} value={%s.Input}", prefix)
	}
	if p := pparam.(*client.PeriodParam); p.Input != "" {
		return templatehtml.HTMLAttr(fmt.Sprintf(" value=\"%s\"", p.Input))
	}
	return templatehtml.HTMLAttr("")
}

func refreshClass(pparam interface{}, classes string) interface{} {
	if JSX {
		prefix, _ := DotSplitHash(pparam)
		return fmt.Sprintf(" %s={%q + (%s.InputErrd ? \" has-warning\" : \"\")}", classword(), classes, prefix)
	}
	if p := pparam.(*client.PeriodParam); p.InputErrd {
		classes += " " + "has-warning"
	}
	return templatehtml.HTMLAttr(fmt.Sprintf(" %s=%q", classword(), classes))
}

func ifDisabledAttr(value interface{}) templatehtml.HTMLAttr {
	if JSX {
		return templatehtml.HTMLAttr(fmt.Sprintf("disabled={%s.Value ? \"disabled\" : \"\" }", uncurl(value.(string))))
	}
	if value.(*client.BoolParam).Value {
		return templatehtml.HTMLAttr("disabled=\"disabled\"")
	}
	return templatehtml.HTMLAttr("")
}

func ifClassAttr(value interface{}, classes ...string) (templatehtml.HTMLAttr, error) {
	s, err := ifClass(value, classes...)
	if err != nil {
		return templatehtml.HTMLAttr(""), err
	}
	if !JSX {
		s = fmt.Sprintf("%q", s)
	}
	return templatehtml.HTMLAttr(fmt.Sprintf(" %s=%s", classword(), s)), nil
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
	if value.(*client.BoolParam).Value {
		return fstclass, nil
	}
	return sndclass, nil
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
	ep := value.(*client.EnumParam)
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
		ClassAttr templatehtml.HTMLAttr
	}{
		Value:     val,
		ClassAttr: templatehtml.HTMLAttr(ca),
	}
}

func key(prefix, val string) templatehtml.HTMLAttr {
	if !JSX {
		return templatehtml.HTMLAttr("")
	}
	return templatehtml.HTMLAttr(fmt.Sprintf(" key={%q+%s}", prefix+"-", uncurl(val)))
}

type Clipped struct {
	IDAttr      templatehtml.HTMLAttr
	ForAttr     templatehtml.HTMLAttr
	MWStyleAttr templatehtml.HTMLAttr
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
		IDAttr:      templatehtml.HTMLAttr("id=" + key),
		ForAttr:     templatehtml.HTMLAttr(forword() + "=" + key),
		MWStyleAttr: templatehtml.HTMLAttr("style=" + mws),
		Text:        val,
	}, nil
}

type dotValue struct {
	s     string
	hashp *templatepipe.Hash
}

// func (dv dotValue) GoString() string { return dv.GoString() } // WTF?

func (dv dotValue) String() string {
	v := dv.s
	delete(*dv.hashp, "OVERRIDE")
	return v
}

func dot(v interface{}, key string) templatepipe.Hash {
	h := v.(templatepipe.Hash)
	h["OVERRIDE"] = dotValue{s: templatepipe.Curly(key), hashp: &h}
	return h
}

// AceFuncs features functions for templates. In use in acepp and templates.
var AceFuncs = templatehtml.FuncMap{
	"dot":        dot,
	"key":        key,
	"clip":       clip,
	"droplink":   droplink,
	"usepercent": usepercent,

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

// StringExecuteHTML does t.Execute into string returned. Does not clone.
func StringExecuteHTML(t *templatehtml.Template, data interface{}) (string, error) {
	buf := new(bytes.Buffer)
	if err := t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// StringExecute does t.Execute into string returned. Does not clone.
func StringExecute(t *templatetext.Template, data interface{}) (string, error) {
	buf := new(bytes.Buffer)
	if err := t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
