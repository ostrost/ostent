package templatefunc

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/ostrost/ostent/params"
)

func (f JSXFuncs) Class() string    { return "className" }
func (f HTMLFuncs) Class() string   { return "class" }
func (f JSXFuncs) Colspan() string  { return "colSpan" }
func (f HTMLFuncs) Colspan() string { return "colspan" }

func (f JSXFuncs) ClassNonzero(x interface{}, class, sndclass string) template.HTMLAttr {
	return SprintfAttr(` className={%s.Absolute != 0 ? %q : %q}`,
		x.(Uncurler).Uncurl(), class, sndclass)
}

func (f HTMLFuncs) ClassNonzero(x interface{}, class, sndclass string) template.HTMLAttr {
	if x.(params.Num).Absolute == 0 {
		class = sndclass
	}
	return SprintfAttr(" class=%q", class)
}

func (f JSXFuncs) ClassPositive(x interface{}, class, sndclass string) template.HTMLAttr {
	return SprintfAttr(` className={!%s.Negative ? %q : %q}`,
		x.(Uncurler).Uncurl(), class, sndclass)
}

func (f HTMLFuncs) ClassPositive(x interface{}, class, sndclass string) template.HTMLAttr {
	if x.(params.Num).Negative {
		class = sndclass
	}
	return SprintfAttr(" class=%q", class)
}

func (f JSXFuncs) ClassMutext(x interface{}) template.HTMLAttr {
	return SprintfAttr(" className={%s == \"0\" ? %q : \"\"}",
		x.(Uncurler).Uncurl(), "mutext")
}

func (f HTMLFuncs) ClassMutext(x interface{}) template.HTMLAttr {
	if x.(string) == "0" {
		return SprintfAttr(" class=%q", "mutext")
	}
	return SprintfAttr("")
}

func (f JSXFuncs) ClassMutext2(x, y interface{}) template.HTMLAttr {
	return SprintfAttr(" className={%s == \"0\" || %s == \"0\" ? %q : \"\"}",
		x.(Uncurler).Uncurl(), x.(Uncurler).Uncurl(), "mutext")
}

func (f HTMLFuncs) ClassMutext2(x, y interface{}) template.HTMLAttr {
	if x.(string) == "0" || y.(string) == "0" {
		return SprintfAttr(" class=%q", "mutext")
	}
	return SprintfAttr("")
}

// Key returns key attribute: prefix + uncurled x being an Uncurler.
func (f JSXFuncs) Key(prefix string, x interface{}) template.HTMLAttr {
	return SprintfAttr(" key={%q+%s}", prefix+"-", x.(Uncurler).Uncurl())
}

// Key returns empty attribute.
func (f HTMLFuncs) Key(_ string, x interface{}) (empty template.HTMLAttr) { return }

func (f JSXFuncs) FuncHrefT() interface{} {
	return func(_, n Uncurler) (template.HTMLAttr, error) {
		base, last := f.Split(n)
		return SprintfAttr(" href={%s.Tlinks.%s} onClick={this.handleClick}",
			base, last), nil
	}
}

func (f HTMLFuncs) FuncHrefT() interface{} { return f.ParamsFuncs.HrefT }

func (f JSXFuncs) FuncLessD() interface{} {
	return func(_, dur Uncurler, bclass string) (params.ALink, error) {
		return f.Dlink(dur, bclass, "Less", "-")
	}
}

func (f JSXFuncs) FuncMoreD() interface{} {
	return func(_, dur Uncurler, bclass string) (params.ALink, error) {
		return f.Dlink(dur, bclass, "More", "+")
	}
}

func (f HTMLFuncs) FuncLessD() interface{} { return f.ParamsFuncs.LessD }
func (f HTMLFuncs) FuncMoreD() interface{} { return f.ParamsFuncs.MoreD }

func (f JSXFuncs) FuncLessN() interface{} {
	return func(_, num Uncurler, bclass string) (params.ALink, error) {
		return f.Nlink(num, bclass, "Less", "-")
	}
}

func (f JSXFuncs) FuncMoreN() interface{} {
	return func(_, num Uncurler, bclass string) (params.ALink, error) {
		return f.Nlink(num, bclass, "More", "+")
	}
}

func (f HTMLFuncs) FuncLessN() interface{} { return f.ParamsFuncs.LessN }
func (f HTMLFuncs) FuncMoreN() interface{} { return f.ParamsFuncs.MoreN }

func (f JSXFuncs) FuncVlink() interface{} {
	return func(_, this Uncurler, cmp int, text, alignClass string) params.VLink {
		base, last := f.Split(this)
		return params.VLink{
			AlignClass: alignClass,
			CaretClass: fmt.Sprintf("{%s.Vlinks.%s[%d-1].%s}", base, last, cmp, "CaretClass"),
			LinkClass:  fmt.Sprintf("{%s.Vlinks.%s[%d-1].%s}", base, last, cmp, "LinkClass"),
			LinkHref:   fmt.Sprintf("{%s.Vlinks.%s[%d-1].%s}", base, last, cmp, "LinkHref"),
			LinkText:   text, // always static
		}
	}
}

func (f HTMLFuncs) FuncVlink() interface{} { return f.ParamsFuncs.Vlink }

func (f JSXFuncs) Dlink(v Uncurler, bclass, which, badge string) (params.ALink, error) {
	base, last := f.Split(v)
	var (
		href   = fmt.Sprintf( /**/ "{%s.Dlinks.%s.%s.Href}", base, last, which)
		text   = fmt.Sprintf( /**/ "{%s.Dlinks.%s.%s.Text}", base, last, which)
		eclass = fmt.Sprintf( /* */ "%s.Dlinks.%s.%s.ExtraClass", base, last, which) // not curled
	)
	return params.ALink{
		Href:  href,
		Text:  text,
		Badge: badge,
		Class: f.ConcatClass(bclass, eclass),
	}, nil
}

func (f JSXFuncs) ConcatClass(bclass, eclass string) string {
	return fmt.Sprintf("{%q + \" \" + (%s != null ? %s : \"\")}", bclass, eclass, eclass)
}

func (f JSXFuncs) Nlink(v Uncurler, bclass, which, badge string) (params.ALink, error) {
	base, last := f.Split(v)
	var (
		href   = fmt.Sprintf( /**/ "{%s.Nlinks.%s.%s.Href}", base, last, which)
		text   = fmt.Sprintf( /**/ "{%s.Nlinks.%s.%s.Text}", base, last, which)
		eclass = fmt.Sprintf( /* */ "%s.Nlinks.%s.%s.ExtraClass", base, last, which) // not curled
	)
	return params.ALink{
		Href:  href,
		Text:  text,
		Badge: badge,
		Class: f.ConcatClass(bclass, eclass),
	}, nil
}

func (f JSXFuncs) Split(v Uncurler) (string, string) {
	split := strings.Split(v.Uncurl(), ".")
	return strings.Join(split[:len(split)-1], "."), split[len(split)-1]
}

// JSXFuncs has methods implementing Functor.
type JSXFuncs struct{}

// HTMLFuncs has methods implementing Functor.
type HTMLFuncs struct{ params.ParamsFuncs }

// MakeMap is dull but required.
func (f JSXFuncs) MakeMap() template.FuncMap { return MakeMap(f) }

// MakeMap is dull but required.
func (f HTMLFuncs) MakeMap() template.FuncMap { return MakeMap(f) }

// MakeMap constructs template.FuncMap off f implementation.
func MakeMap(f Functor) template.FuncMap {
	return template.FuncMap{
		"HTML":    func(s string) template.HTML { return template.HTML(s) },
		"rowsset": func(interface{}) string { return "" }, // empty pipeline
		// acepp overrides rowsset and adds setrows

		"class":   f.Class,
		"colspan": f.Colspan,

		"AttrKey":       f.Key,
		"ClassNonzero":  f.ClassNonzero,
		"ClassPositive": f.ClassPositive,
		"ClassMutext":   f.ClassMutext,
		"ClassMutext2":  f.ClassMutext2,

		"HrefT": f.FuncHrefT(),
		"LessD": f.FuncLessD(),
		"MoreD": f.FuncMoreD(),
		"LessN": f.FuncLessN(),
		"MoreN": f.FuncMoreN(),
		"Vlink": f.FuncVlink(),
	}
}

// Funcs features functions for templates. In use in acepp and templates.
var Funcs = HTMLFuncs{}.MakeMap()

type Functor interface {
	MakeMap() template.FuncMap
	Class() string
	Colspan() string
	ClassNonzero(interface{}, string, string) template.HTMLAttr
	ClassPositive(interface{}, string, string) template.HTMLAttr
	ClassMutext(interface{}) template.HTMLAttr
	ClassMutext2(interface{}, interface{}) template.HTMLAttr
	Key(string, interface{}) template.HTMLAttr

	FuncHrefT() interface{}
	FuncLessD() interface{}
	FuncMoreD() interface{}
	FuncLessN() interface{}
	FuncMoreN() interface{}
	FuncVlink() interface{}
}

func SprintfAttr(format string, args ...interface{}) template.HTMLAttr {
	return template.HTMLAttr(fmt.Sprintf(format, args...))
}

type Uncurler interface {
	Uncurl() string
}
