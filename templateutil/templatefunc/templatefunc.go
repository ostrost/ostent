package templatefunc

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/ostrost/ostent/params"
	"github.com/ostrost/ostent/templateutil"
)

func (f JSXLFuncs) FuncHrefT() interface{} { return func() int { panic(fmt.Errorf("Do not use 1")) } }
func (f JSXLFuncs) FuncLessD() interface{} { return func() int { panic(fmt.Errorf("Do not use 2")) } }
func (f JSXLFuncs) FuncMoreD() interface{} { return func() int { panic(fmt.Errorf("Do not use 3")) } }
func (f JSXLFuncs) FuncLessN() interface{} { return func() int { panic(fmt.Errorf("Do not use 4")) } }
func (f JSXLFuncs) FuncMoreN() interface{} { return func() int { panic(fmt.Errorf("Do not use 5")) } }

func (f HTMLFuncs) FuncHrefT() interface{} { return f.ParamsFuncs.HrefT }
func (f HTMLFuncs) FuncLessD() interface{} { return f.ParamsFuncs.LessD }
func (f HTMLFuncs) FuncMoreD() interface{} { return f.ParamsFuncs.MoreD }
func (f HTMLFuncs) FuncLessN() interface{} { return f.ParamsFuncs.LessN }
func (f HTMLFuncs) FuncMoreN() interface{} { return f.ParamsFuncs.MoreN }
func (f HTMLFuncs) FuncVlink() interface{} { return f.ParamsFuncs.Vlink }

func (f JSXLFuncs) FuncVlink() interface{} {
	type Uncurler interface {
		Uncurl() string
	}
	return func(_, this Uncurler, cmp int, text string) params.VLink {
		split := strings.Split(this.Uncurl(), ".")
		base, last := strings.Join(split[:len(split)-1], "."), split[len(split)-1]
		return params.VLink{
			CaretClass: fmt.Sprintf("{%s.Vlinks.%s[%d-1].%s}", base, last, cmp, "CaretClass"),
			LinkClass:  fmt.Sprintf("{%s.Vlinks.%s[%d-1].%s}", base, last, cmp, "LinkClass"),
			LinkHref:   fmt.Sprintf("{%s.Vlinks.%s[%d-1].%s}", base, last, cmp, "LinkHref"),
			LinkText:   text, // always static
		}
	}
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

	FuncHrefT() interface{}
	FuncLessD() interface{}
	FuncMoreD() interface{}
	FuncLessN() interface{}
	FuncMoreN() interface{}
	FuncVlink() interface{}
}
