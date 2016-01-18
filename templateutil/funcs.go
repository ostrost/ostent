package templateutil

import "html/template"

func (f HTMLFuncs) Class() string   { return "class" }
func (f JSXLFuncs) Class() string   { return "lcassName" }
func (f HTMLFuncs) Colspan() string { return "colspan" }
func (f JSXLFuncs) Colspan() string { return "colSpan" }

// HTMLFuncs has methods implementing Functor.
type HTMLFuncs struct{}

// JSXLFuncs has methods implementing Functor.
type JSXLFuncs struct{}

type Functor interface {
	Class() string
	Colspan() string
}

// CombineMaps makes new template.FuncMap off f implementation and extra.
func CombineMaps(f Functor, extra template.FuncMap) template.FuncMap {
	combined := template.FuncMap{
		"HTML": func(s string) template.HTML { return template.HTML(s) },

		"class":   f.Class,
		"colspan": f.Colspan,
	}
	for k, f := range extra {
		combined[k] = f
	}
	return combined
}

func FuncMapHTML() template.FuncMap { return CombineMaps(NewHTMLFuncs(), nil) }
func FuncMapJSXL() template.FuncMap { return CombineMaps(NewJSXLFuncs(), nil) }

func NewHTMLFuncs() HTMLFuncs { return HTMLFuncs{} }
func NewJSXLFuncs() JSXLFuncs { return JSXLFuncs{} }
