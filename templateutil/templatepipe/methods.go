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

func (_ Nota) AttrClassP(v Uncurler, fstclass, sndclass string) template.HTMLAttr {
	return SprintfAttr(` className={!%s.Negative ? %q : %q}`,
		v.Uncurl(), fstclass, sndclass)
}

func (_ Nota) AttrClassNonzero(v Uncurler, fstclass, sndclass string) template.HTMLAttr {
	return SprintfAttr(` className={%s.Absolute != 0 ? %q : %q}`,
		v.Uncurl(), fstclass, sndclass)
}

func (_ Nota) AttrClassTab(num, tab Uncurler, cmp int, fstclass, sndclass string) template.HTMLAttr {
	return SprintfAttr(` className={%s.Absolute != 0 && %s.Absolute == %d ? %q : %q}`,
		num.Uncurl(), tab.Uncurl(), cmp, fstclass, sndclass)
}

func (_ Nota) Vlink(this Uncurler, cmp int, text, alignClass string) params.VLink {
	base, last := Split(this)
	return params.VLink{
		AlignClass: alignClass,
		CaretClass: fmt.Sprintf("{%s.Vlinks.%s[%d-1].%s}", base, last, cmp, "CaretClass"),
		LinkClass:  fmt.Sprintf("{%s.Vlinks.%s[%d-1].%s}", base, last, cmp, "LinkClass"),
		LinkHref:   fmt.Sprintf("{%s.Vlinks.%s[%d-1].%s}", base, last, cmp, "LinkHref"),
		LinkText:   text, // always static
	}
}

// ALink is a shadow of params.ALink: has it's own Class method.
// The .ExtraClass must contain uncurled string.
type ALink params.ALink

func (al ALink) Class(base string) string {
	return fmt.Sprintf("{%q + \" \" + (%s != null ? %s : \"\")}",
		base, al.ExtraClass, al.ExtraClass)
}

func (_ Nota) ZeroN(num Uncurler) (ALink, error) { return Nlink(num, "Zero", "") }
func (_ Nota) LessN(num Uncurler) (ALink, error) { return Nlink(num, "Less", "-") }
func (_ Nota) MoreN(num Uncurler) (ALink, error) { return Nlink(num, "More", "+") }

func Nlink(v Uncurler, which, badge string) (ALink, error) {
	base, last := Split(v)
	var (
		href  = fmt.Sprintf( /**/ "{%s.Nlinks.%s.%s.Href}", base, last, which)
		text  = fmt.Sprintf( /**/ "{%s.Nlinks.%s.%s.Text}", base, last, which)
		class = fmt.Sprintf( /* */ "%s.Nlinks.%s.%s.Class", base, last, which) // not curled
	)
	return ALink{APlain: params.APlain{Href: href, Text: text, Badge: badge}, ExtraClass: class}, nil
}

func (_ Nota) LessD(dur Uncurler) (ALink, error) { return Dlink(dur, "Less", "-") }
func (_ Nota) MoreD(dur Uncurler) (ALink, error) { return Dlink(dur, "More", "+") }

func Dlink(v Uncurler, which, badge string) (ALink, error) {
	base, last := Split(v)
	var (
		href  = fmt.Sprintf( /**/ "{%s.Dlinks.%s.%s.Href}", base, last, which)
		text  = fmt.Sprintf( /**/ "{%s.Dlinks.%s.%s.Text}", base, last, which)
		class = fmt.Sprintf( /* */ "%s.Dlinks.%s.%s.Class", base, last, which)
	)
	return ALink{APlain: params.APlain{Href: href, Text: text, Badge: badge}, ExtraClass: class}, nil
}

func (_ Nota) AttrHrefToggleNegative(v Uncurler) interface{} {
	base, last := Split(v)
	return fmt.Sprintf(" href={%s.Tlinks.%s} onClick={this.handleClick}", base, last)
}

// Split splits uncurled v by last ".".
func Split(v Uncurler) (string, string) {
	split := strings.Split(v.Uncurl(), ".")
	return strings.Join(split[:len(split)-1], "."), split[len(split)-1]
	// return split[len(split)-1]
}

func SprintfAttr(format string, args ...interface{}) template.HTMLAttr {
	return template.HTMLAttr(fmt.Sprintf(format, args...))
}
