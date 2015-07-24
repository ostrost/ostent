package templatepipe

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/ostrost/ostent/params"
)

func Uncurl(s string) string {
	return strings.TrimSuffix(strings.TrimPrefix(s, "{"), "}")
}

func (n Nota) Uncurl() string {
	return Uncurl(n.String())
}

func (n Nota) FormActionAttr() interface{} {
	return fmt.Sprintf(" action={\"/form/\"+%s}", n.Uncurl())
}

func (n Nota) KeyAttr(prefix string) template.HTMLAttr {
	return SprintfAttr(" key={%q+%s}", prefix+"-", n.Uncurl())
}

func (n Nota) BoolClassAttr(fstclass, sndclass string) template.HTMLAttr {
	return SprintfAttr(" className={%s ? %q : %q}", n.Uncurl(), fstclass, sndclass)
}

func (n Nota) BoolParamClassAttr(fstclass, sndclass string) template.HTMLAttr {
	return SprintfAttr(" className={%s.Value ? %q : %q}", n.Uncurl(), fstclass, sndclass)
}

func (n Nota) DisabledAttr() interface{} {
	return fmt.Sprintf(" disabled={%s.Value ? %q : \"\" }", n.Uncurl(), "disabled")
}

func (n Nota) EnumClassAttr(named, classif string, optelse ...string) (template.HTMLAttr, error) {
	classelse, err := params.EnumClassAttrArgs(optelse)
	if err != nil {
		return template.HTMLAttr(""), err
	}
	eparams := params.NewParamsENUM(nil)
	ed := eparams[n.Base()].EnumDecodec
	_, uptr := ed.Unew()
	if err := uptr.Unmarshal(named, new(bool)); err != nil {
		return template.HTMLAttr(""), err
	}
	return SprintfAttr(" className={(%s.Uint == %d) ? %q : %q}",
		n, uptr.Touint(), classif, classelse), nil
}

func (n Nota) EnumLink(args ...string) (interface{}, error) {
	named, aclass := params.EnumLinkArgs(args)
	eparams := params.NewParamsENUM(nil)
	ed := eparams[n.Base()].EnumDecodec
	return params.EnumLink{
		AlignClass: aclass,
		Text:       ed.Text(named), // always static
		Href:       fmt.Sprintf("{%s.%s.%s}", n, named, "Href"),
		Class:      fmt.Sprintf("{%s.%s.%s}", n, named, "Class"),
		CaretClass: fmt.Sprintf("{%s.%s.%s}", n, named, "CaretClass"),
	}, nil
}

func (n Nota) ToggleHrefAttr() interface{} {
	return fmt.Sprintf(" href={%s.Href} onClick={this.handleClick}", n.Uncurl())
}

func (n Nota) PeriodNameAttr() interface{} {
	return fmt.Sprintf(" name=%q", n.Base())
}

func (n Nota) PeriodValueAttr() interface{} {
	return fmt.Sprintf(" value={%s.Input} onChange={this.handleChange}", n)
}

func (n Nota) RefreshClassAttr(classes string) interface{} {
	return fmt.Sprintf(" className={%q + (%s.InputErrd ? %q : \"\")}",
		classes, n, " has-warning")
}

func (n Nota) LessHrefAttr() interface{} {
	return fmt.Sprintf(" href={%s.LessHref} onClick={this.handleClick}", n.Uncurl())
}

func (n Nota) MoreHrefAttr() interface{} {
	return fmt.Sprintf(" href={%s.MoreHref} onClick={this.handleClick}", n.Uncurl())
}

// Base is like filepath.Base on n with "." separator.
func (n Nota) Base() string {
	split := strings.Split(n.String(), ".")
	return split[len(split)-1]
}

func SprintfAttr(format string, args ...interface{}) template.HTMLAttr {
	return template.HTMLAttr(fmt.Sprintf(format, args...))
}
