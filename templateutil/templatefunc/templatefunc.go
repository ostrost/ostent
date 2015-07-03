package templatefunc

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/ostrost/ostent/system/operating"
	"github.com/ostrost/ostent/templateutil/templatepipe"
)

func (f JSXFuncs) Class() string    { return "className" }
func (f HTMLFuncs) Class() string   { return "class" }
func (f JSXFuncs) Colspan() string  { return "colSpan" }
func (f HTMLFuncs) Colspan() string { return "colspan" }

// JSXClose returns empty template.HTML.
func (f HTMLFuncs) JSXClose(string) (empty template.HTML) { return }

// JSXClose returns close tag markup as template.HTML.
func (f JSXFuncs) JSXClose(tag string) template.HTML {
	return template.HTML("</" + tag + ">")
}

// JSXFuncs has methods implementing Functor.
type JSXFuncs struct{}

// HTMLFuncs has methods implementing Functor.
type HTMLFuncs struct{}

// MakeMap is dull but required.
func (f JSXFuncs) MakeMap() template.FuncMap { return MakeMap(f) }

// MakeMap is dull but required.
func (f HTMLFuncs) MakeMap() template.FuncMap { return MakeMap(f) }

// MakeMap constructs template.FuncMap off f implementation.
func MakeMap(f Functor) template.FuncMap {
	return template.FuncMap{
		"rowsset": func(interface{}) string { return "" }, // empty pipeline
		// acepp overrides rowsset and adds setrows

		"class":    f.Class,
		"colspan":  f.Colspan,
		"jsxClose": f.JSXClose,
	}
}

// Funcs features functions for templates. In use in acepp and templates.
var Funcs = HTMLFuncs{}.MakeMap()

type Functor interface {
	MakeMap() template.FuncMap
	Class() string
	Colspan() string
	JSXClose(string) template.HTML
}

func init() {
	// check for Nota's interfaces compliance
	_ = interface {
		// operating (multiple types):
		BoolClassAttr(...string) (template.HTMLAttr, error)
		Clip(int, string, ...fmt.Stringer) (*operating.Clipped, error)
		KeyAttr(string) template.HTMLAttr

		FormActionAttr() interface{}                                        // Query
		BoolParamClassAttr(...string) (template.HTMLAttr, error)            // BoolParam
		DisabledAttr() interface{}                                          // BoolParam
		ToggleHrefAttr() interface{}                                        // BoolParam
		EnumClassAttr(string, string, ...string) (template.HTMLAttr, error) // EnumParam
		EnumLink(...string) (interface{}, error)                            // EnumParam
		PeriodNameAttr() interface{}                                        // PeriodParam
		PeriodValueAttr() interface{}                                       // PeriodParam
		RefreshClassAttr(string) interface{}                                // PeriodParam
		LessHrefAttr() interface{}                                          // LimitParam
		MoreHrefAttr() interface{}                                          // LimitParam
	}(templatepipe.Nota(nil))
}

// SetKFunc constructs a func which
// sets k key to templatepipe.Curly(string (n))
// in passed interface{} (v) being a templatepipe.Nota.
// SetKFunc is used by acepp only.
func SetKFunc(k string) func(interface{}, string) interface{} {
	return func(v interface{}, n string) interface{} {
		if args := strings.Split(n, " "); len(args) > 1 {
			var list []string
			for _, arg := range args {
				list = append(list, templatepipe.Curl(arg))
			}
			v.(templatepipe.Nota)[k] = list
			return v
		}
		v.(templatepipe.Nota)[k] = templatepipe.Curl(n)
		return v
	}
}

// GetKFunc constructs a func which
// gets, deletes and returns k key
// in passed interface{} (v) being a templatepipe.Nota.
// GetKFunc is used by acepp only.
func GetKFunc(k string) func(interface{}) interface{} {
	return func(v interface{}) interface{} {
		h, ok := v.(templatepipe.Nota)
		if !ok {
			return "" // empty pipeline, affects dispatch
		}
		n := h[k]
		if args, ok := n.([]string); ok {
			if len(args) > 1 {
				h[k] = args[1:]
			} else {
				delete(h, k)
			}
			return args[0]
		}
		delete(h, k)
		return n // may also be empty, affects dispatch
	}
}
