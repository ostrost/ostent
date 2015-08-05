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

func (n Nota) AttrActionForm() template.HTMLAttr {
	return SprintfAttr(" action={\"/form/\"+%s}", n.Uncurl())
}

func (n Nota) AttrKey(prefix string) template.HTMLAttr {
	return SprintfAttr(" key={%q+%s}", prefix+"-", n.Uncurl())
}

func (_ Nota) AttrClassN(v, fstclass, sndclass string) template.HTMLAttr {
	return SprintfAttr(" className={%s ? %q : %q}", Uncurl(v), fstclass, sndclass)
}

func (_ Nota) AttrClassT(defaults, v string, cmp int, fstclass, sndclass string) template.HTMLAttr {
	defaults, v = Uncurl(defaults), Uncurl(v)
	split := strings.Split(v, ".")
	last := split[len(split)-1]
	return SprintfAttr(" className={%d == %s || (%s == 0 && %d == %s.%s) ? %q : %q}",
		cmp, v, v, cmp, defaults, last, fstclass, sndclass)
}

func (_ Nota) AttrClassParamsError(errs interface{}, name, fstclass, sndclass string) template.HTMLAttr {
	serrs := Uncurl(errs.(string))
	return SprintfAttr(" className={%s && %s.%s ? %q : %q}",
		serrs, serrs, name, fstclass, sndclass)
}

/*
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
} // */

func (n Nota) Variate(this string, cmp int, text, alignClass string) interface{} {
	split := strings.Split(Uncurl(this), ".")
	base := strings.Join(split[:len(split)-1], ".")
	last := split[len(split)-1]
	// param = Uncurl(param)
	return params.Varlink{
		AlignClass: alignClass,
		CaretClass: fmt.Sprintf("{%s.Variations.%s[%d-1].%s}", base, last, cmp, "CaretClass"),
		LinkClass:  fmt.Sprintf("{%s.Variations.%s[%d-1].%s}", base, last, cmp, "LinkClass"),
		LinkHref:   fmt.Sprintf("{%s.Variations.%s[%d-1].%s}", base, last, cmp, "LinkHref"),
		LinkText:   text, // always static
	}
}

/*
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
// */

func (_ Nota) AttrHrefToggle(s string) interface{} {
	split := strings.Split(Uncurl(s), ".")
	base := strings.Join(split[:len(split)-1], ".")
	last := split[len(split)-1]
	return fmt.Sprintf(" href={%s.Toggle.%s} onClick={this.handleClick}", base, last)
}

// Data.Params.RefreshXX RefreshXX
func (n Nota) AttrNameRefresh(fieldName string) interface{} {
	// TODO fieldName is "Refreshsmth", ought to have a map to actual parameter
	// lowercase fieldName suffices for now
	return fmt.Sprintf(" name=%q", strings.ToLower(fieldName))
}

func (n Nota) AttrValueRefresh(fieldName string) interface{} {
	return fmt.Sprintf(" value={%s.%s} onChange={this.handleChange}",
		n, fieldName)
}

/*
func (n Nota) RefreshClassAttr(classes string) interface{} {
 	return fmt.Sprintf(" className={%q + (%s.InputErrd ? %q : \"\")}",
 		classes, n, " has-warning")
}
*/

func (n Nota) AttrHrefLess(s string) interface{} {
	return fmt.Sprintf(" href={%s.LessHref} onClick={this.handleClick}", n.Uncurl())
}

func (n Nota) AttrHrefMore(s string) interface{} {
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
