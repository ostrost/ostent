package templatepipe

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/ostrost/ostent/params"
	"github.com/ostrost/ostent/system/operating"
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
	return template.HTMLAttr(fmt.Sprintf(" key={%q+%s}", prefix+"-", n.Uncurl()))
}

func (n Nota) BoolClassAttr(classes ...string) (template.HTMLAttr, error) {
	fstclass, sndclass, err := operating.ClassesChoices("BoolClassAttr", classes)
	if err != nil {
		return template.HTMLAttr(""), err
	}
	return template.HTMLAttr(fmt.Sprintf(" className={%s ? %q : %q}",
		n.Uncurl(), fstclass, sndclass)), nil
}

func (n Nota) BoolParamClassAttr(classes ...string) (template.HTMLAttr, error) {
	fstclass, sndclass, err := operating.ClassesChoices("BoolParamClassAttr", classes)
	if err != nil {
		return template.HTMLAttr(""), err
	}
	return template.HTMLAttr(fmt.Sprintf(" className={%s.Value ? %q : %q}",
		n.Uncurl(), fstclass, sndclass)), nil
}

func (n Nota) Clip(width int, prefix string, id ...fmt.Stringer) (*operating.Clipped, error) {
	k, err := operating.ClipArgs(id, n.Uncurl())
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("{%q+%s}", prefix+"-", Uncurl(k))
	return &operating.Clipped{
		IDAttr:      operating.SprintfAttr(" id=%s", key),
		ForAttr:     operating.SprintfAttr(" htmlFor=%s", key),
		MWStyleAttr: operating.SprintfAttr(" style={{maxWidth: '%dch'}}", width),
		Text:        n.String(),
	}, nil
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
	return operating.SprintfAttr(" className={(%s.Uint == %d) ? %q : %q}",
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
