package templatefunc

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/ostrost/ostent/params"
	"github.com/ostrost/ostent/templateutil"
)

func (f JSXLFuncs) ClassNonZero(x interface{}, class, sndclass string) template.HTMLAttr {
	return SprintfAttr(` %s={%s.Absolute != 0 ? %q : %q}`,
		f.Class(),
		x.(Uncurler).Uncurl(), class, sndclass)
}

func (f HTMLFuncs) ClassNonZero(x interface{}, class, sndclass string) template.HTMLAttr {
	if x.(params.Num).Absolute == 0 {
		class = sndclass
	}
	return SprintfAttr(" class=%q", class)
}

func (f JSXLFuncs) ClassPositive(x interface{}, class, sndclass string) template.HTMLAttr {
	return SprintfAttr(` %s={!%s.Negative ? %q : %q}`,
		f.Class(),
		x.(Uncurler).Uncurl(), class, sndclass)
}

func (f HTMLFuncs) ClassPositive(x interface{}, class, sndclass string) template.HTMLAttr {
	if x.(params.Num).Negative {
		class = sndclass
	}
	return SprintfAttr(" class=%q", class)
}

func (f JSXLFuncs) FuncHrefT() interface{} {
	return func(_, n Uncurler) (template.HTMLAttr, error) {
		base, last := f.Split(n)
		return SprintfAttr(" href={%s.Tlinks.%s} onClick={this.handleClick}",
			base, last), nil
	}
}

func (f HTMLFuncs) FuncHrefT() interface{} { return f.ParamsFuncs.HrefT }

func (f JSXLFuncs) FuncLessD() interface{} {
	return func(_, dur Uncurler, bclass string) (params.ALink, error) {
		return f.Dlink(dur, bclass, "Less", "-")
	}
}

func (f JSXLFuncs) FuncMoreD() interface{} {
	return func(_, dur Uncurler, bclass string) (params.ALink, error) {
		return f.Dlink(dur, bclass, "More", "+")
	}
}

func (f HTMLFuncs) FuncLessD() interface{} { return f.ParamsFuncs.LessD }
func (f HTMLFuncs) FuncMoreD() interface{} { return f.ParamsFuncs.MoreD }

func (f JSXLFuncs) FuncLessN() interface{} {
	return func(_, num Uncurler, bclass string) (params.ALink, error) {
		return f.Nlink(num, bclass, "Less", "-")
	}
}

func (f JSXLFuncs) FuncMoreN() interface{} {
	return func(_, num Uncurler, bclass string) (params.ALink, error) {
		return f.Nlink(num, bclass, "More", "+")
	}
}

func (f HTMLFuncs) FuncLessN() interface{} { return f.ParamsFuncs.LessN }
func (f HTMLFuncs) FuncMoreN() interface{} { return f.ParamsFuncs.MoreN }

func (f JSXLFuncs) FuncVlink() interface{} {
	return func(_, this Uncurler, cmp int, text string) params.VLink {
		base, last := f.Split(this)
		return params.VLink{
			CaretClass: fmt.Sprintf("{%s.Vlinks.%s[%d-1].%s}", base, last, cmp, "CaretClass"),
			LinkClass:  fmt.Sprintf("{%s.Vlinks.%s[%d-1].%s}", base, last, cmp, "LinkClass"),
			LinkHref:   fmt.Sprintf("{%s.Vlinks.%s[%d-1].%s}", base, last, cmp, "LinkHref"),
			LinkText:   text, // always static
		}
	}
}

func (f HTMLFuncs) FuncVlink() interface{} { return f.ParamsFuncs.Vlink }

func (f JSXLFuncs) Dlink(v Uncurler, bclass, which, badge string) (params.ALink, error) {
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

// ConcatClass is internal (not required for interface)
func (f JSXLFuncs) ConcatClass(bclass, eclass string) string {
	return fmt.Sprintf("{%q + \" \" + (%s != null ? %s : \"\")}", bclass, eclass, eclass)
}

// Nlink is internal (not required for interface)
func (f JSXLFuncs) Nlink(v Uncurler, bclass, which, badge string) (params.ALink, error) {
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

// Split is internal (not required for interface)
func (f JSXLFuncs) Split(v Uncurler) (string, string) {
	split := strings.Split(v.Uncurl(), ".")
	return strings.Join(split[:len(split)-1], "."), split[len(split)-1]
}

// JSXLFuncs has methods implementing Functor.
type JSXLFuncs struct{ templateutil.JSXLFuncs }

// HTMLFuncs has methods implementing Functor.
type HTMLFuncs struct {
	// templateutil.Functor
	templateutil.HTMLFuncs
	params.ParamsFuncs
}

// ConstructMaps constructs template.FuncMap off f implementation.
func ConstructMap(f Functor) template.FuncMap {
	return templateutil.CombineMaps(f, template.FuncMap{
		"ClassNonZero":  f.ClassNonZero,
		"ClassPositive": f.ClassPositive,

		"HrefT": f.FuncHrefT(),
		"LessD": f.FuncLessD(),
		"MoreD": f.FuncMoreD(),
		"LessN": f.FuncLessN(),
		"MoreN": f.FuncMoreN(),
		"Vlink": f.FuncVlink(),
	})
}

func FuncMapJSXL() template.FuncMap {
	return ConstructMap(JSXLFuncs{templateutil.NewJSXLFuncs()})
}

func FuncMapHTML() template.FuncMap {
	return ConstructMap(HTMLFuncs{
		HTMLFuncs:   templateutil.NewHTMLFuncs(),
		ParamsFuncs: params.ParamsFuncs{},
	})
}

type Functor interface {
	templateutil.Functor

	ClassNonZero(interface{}, string, string) template.HTMLAttr
	ClassPositive(interface{}, string, string) template.HTMLAttr

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
