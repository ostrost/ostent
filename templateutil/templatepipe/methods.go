package templatepipe

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/ostrost/ostent/params"
)

type Uncurler interface {
	Uncurl() string
}

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

func (_ Nota) AttrClassP(v Uncurler, fstclass, sndclass string) template.HTMLAttr {
	s := v.Uncurl()
	return SprintfAttr(" className={(%s != \"!0\" && %s >= 0) ? %q : %q}",
		s, s, fstclass, sndclass)
}

func (n Nota) Body() string {
	s := n.Uncurl()
	return fmt.Sprintf("{%s == \"!0\" ? 0 : (%s < 0 ? -%s : %s)}", s, s, s, s)
}

func (_ Nota) AttrClassZero(v Uncurler, fstclass, sndclass string) template.HTMLAttr {
	s := v.Uncurl()
	return SprintfAttr(" className={(%s == \"!0\" || %s == 0) ? %q : %q}",
		s, s, fstclass, sndclass)
}

func (_ Nota) AttrClassN(v Uncurler, fstclass, sndclass string) template.HTMLAttr {
	return SprintfAttr(" className={%s ? %q : %q}", v.Uncurl(), fstclass, sndclass)
}

func (_ Nota) AttrClassT(v Uncurler, cmp int, fstclass, sndclass string) template.HTMLAttr {
	return SprintfAttr(" className={%s == %d ? %q : %q}",
		v.Uncurl(), cmp, fstclass, sndclass)
}

func (_ Nota) AttrClassParamsError(errs Uncurler, name, fstclass, sndclass string) template.HTMLAttr {
	serrs := errs.Uncurl() // Uncurl(errs.(string)) // errs used to be interface{}
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

func (n Nota) Variate(this Uncurler, cmp int, text, alignClass string) interface{} {
	split := strings.Split(this.Uncurl(), ".")
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

type ALink params.ALink

func (al ALink) Class(base string) string {
	add := Uncurl(al.ExtraClass)
	return fmt.Sprintf("{%q + \" \" + (%s != null ? %s : \"\")}", base, add, add)
}

func (n Nota) ZeroN(num Uncurler) (ALink, error) { return n.Numbered(num, "Zero", "") }
func (n Nota) LessN(num Uncurler) (ALink, error) { return n.Numbered(num, "Less", "-") }
func (n Nota) MoreN(num Uncurler) (ALink, error) { return n.Numbered(num, "More", "+") }

func (n Nota) Numbered(v Uncurler, which, badge string) (ALink, error) {
	split := strings.Split(v.Uncurl(), ".")
	base := strings.Join(split[:len(split)-1], ".")
	last := split[len(split)-1]
	var (
		href  = fmt.Sprintf("{%s.Numbered.%s.%s.Href}", base, last, which)
		text  = fmt.Sprintf("{%s.Numbered.%s.%s.Text}", base, last, which)
		class = fmt.Sprintf("{%s.Numbered.%s.%s.Class}", base, last, which)
	)
	return ALink{APlain: params.APlain{Href: href, Text: text, Badge: badge}, ExtraClass: class}, nil
}

func (n Nota) LessD(dur Uncurler) (ALink, error) { return n.Delayed(dur, "Less", "-") }
func (n Nota) MoreD(dur Uncurler) (ALink, error) { return n.Delayed(dur, "More", "+") }

func (n Nota) Delayed(v Uncurler, which, badge string) (ALink, error) {
	split := strings.Split(v.Uncurl(), ".")
	base := strings.Join(split[:len(split)-1], ".")
	last := split[len(split)-1]
	var (
		href  = fmt.Sprintf("{%s.Delayed.%s.%s.Href}", base, last, which)
		text  = fmt.Sprintf("{%s.Delayed.%s.%s.Text}", base, last, which)
		class = fmt.Sprintf("{%s.Delayed.%s.%s.Class}", base, last, which)
	)
	return ALink{APlain: params.APlain{Href: href, Text: text, Badge: badge}, ExtraClass: class}, nil
}

func (_ Nota) AttrHrefToggle(v Uncurler) interface{} {
	split := strings.Split(v.Uncurl(), ".")
	base := strings.Join(split[:len(split)-1], ".")
	last := split[len(split)-1]
	return fmt.Sprintf(" href={%s.Toggle.%s} onClick={this.handleClick}", base, last)
}

func (n Nota) AttrHrefToggleHead(v Uncurler) interface{} {
	return n.AttrHrefToggle(v)
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
