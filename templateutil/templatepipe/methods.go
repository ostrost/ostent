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

func (n Nota) AttrKey(prefix string) template.HTMLAttr {
	return SprintfAttr(" key={%q+%s}", prefix+"-", n.Uncurl())
}

func (_ Nota) AttrClassP(v Uncurler, fstclass, sndclass string) template.HTMLAttr {
	s := v.Uncurl()
	return SprintfAttr(` className={(%s != "!0" && %s.substr(0, 1) != "-") ? %q : %q}`,
		s, s, fstclass, sndclass)
}

func (n Nota) Body() string {
	s := n.Uncurl()
	return fmt.Sprintf(`{%s == "!0" ? "0" : %s.replace(/^-/, "")}`, s, s)
}

func (_ Nota) AttrClassNonzero(v Uncurler, fstclass, sndclass string) template.HTMLAttr {
	s := v.Uncurl()
	return SprintfAttr(` className={(%s != "!0" && %s != "0") ? %q : %q}`,
		s, s, fstclass, sndclass)
}

func (_ Nota) AttrClassTab(num, tab Uncurler, cmp int, fstclass, sndclass string) template.HTMLAttr {
	n := num.Uncurl()
	return SprintfAttr(` className={(%s != "!0" && %s != "0" && %s == "%d") ? %q : %q}`,
		n, n, tab.Uncurl(), cmp, fstclass, sndclass)
}

func (n Nota) Vlink(this Uncurler, cmp int, text, alignClass string) params.VLink {
	base, last := Split(this)
	return params.VLink{
		AlignClass: alignClass,
		CaretClass: fmt.Sprintf("{%s.Vlinks.%s[%d-1].%s}", base, last, cmp, "CaretClass"),
		LinkClass:  fmt.Sprintf("{%s.Vlinks.%s[%d-1].%s}", base, last, cmp, "LinkClass"),
		LinkHref:   fmt.Sprintf("{%s.Vlinks.%s[%d-1].%s}", base, last, cmp, "LinkHref"),
		LinkText:   text, // always static
	}
}

type ALink params.ALink

func (al ALink) Class(base string) string {
	add := Uncurl(al.ExtraClass)
	return fmt.Sprintf("{%q + \" \" + (%s != null ? %s : \"\")}", base, add, add)
}

func (n Nota) ZeroN(num Uncurler) (ALink, error) { return n.Nlink(num, "Zero", "") }
func (n Nota) LessN(num Uncurler) (ALink, error) { return n.Nlink(num, "Less", "-") }
func (n Nota) MoreN(num Uncurler) (ALink, error) { return n.Nlink(num, "More", "+") }

func (n Nota) Nlink(v Uncurler, which, badge string) (ALink, error) {
	base, last := Split(v)
	var (
		href  = fmt.Sprintf("{%s.Nlinks.%s.%s.Href}", base, last, which)
		text  = fmt.Sprintf("{%s.Nlinks.%s.%s.Text}", base, last, which)
		class = fmt.Sprintf("{%s.Nlinks.%s.%s.Class}", base, last, which)
	)
	return ALink{APlain: params.APlain{Href: href, Text: text, Badge: badge}, ExtraClass: class}, nil
}

func (n Nota) LessD(dur Uncurler) (ALink, error) { return n.Dlink(dur, "Less", "-") }
func (n Nota) MoreD(dur Uncurler) (ALink, error) { return n.Dlink(dur, "More", "+") }

func (n Nota) Dlink(v Uncurler, which, badge string) (ALink, error) {
	base, last := Split(v)
	var (
		href  = fmt.Sprintf("{%s.Dlinks.%s.%s.Href}", base, last, which)
		text  = fmt.Sprintf("{%s.Dlinks.%s.%s.Text}", base, last, which)
		class = fmt.Sprintf("{%s.Dlinks.%s.%s.Class}", base, last, which)
	)
	return ALink{APlain: params.APlain{Href: href, Text: text, Badge: badge}, ExtraClass: class}, nil
}

func (_ Nota) AttrHrefToggle(v Uncurler) interface{} {
	base, last := Split(v)
	return fmt.Sprintf(" href={%s.Tlinks.%s} onClick={this.handleClick}", base, last)
}

func (n Nota) AttrHrefToggleHead(v Uncurler) interface{} {
	return n.AttrHrefToggle(v)
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
