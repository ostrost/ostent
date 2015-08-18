package templatefunc

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/ostrost/ostent/templateutil/templatepipe"
)

func (f JSXFuncs) Class() string    { return "className" }
func (f HTMLFuncs) Class() string   { return "class" }
func (f JSXFuncs) Colspan() string  { return "colSpan" }
func (f HTMLFuncs) Colspan() string { return "colspan" }

// Key returns empty attribute.
func (f HTMLFuncs) Key(_ string, x interface{}) (empty template.HTMLAttr) { return }

// Key returns key attribute: prefix + uncurled x being an Uncurler.
func (f JSXFuncs) Key(prefix string, x interface{}) template.HTMLAttr {
	return SprintfAttr(" key={%q+%s}", prefix+"-", x.(Uncurler).Uncurl())
}

func SprintfAttr(format string, args ...interface{}) template.HTMLAttr {
	return template.HTMLAttr(fmt.Sprintf(format, args...))
}

type Uncurler interface {
	Uncurl() string
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

		"AttrKey": f.Key,
		"class":   f.Class,
		"colspan": f.Colspan,
	}
}

// Funcs features functions for templates. In use in acepp and templates.
var Funcs = HTMLFuncs{}.MakeMap()

type Functor interface {
	MakeMap() template.FuncMap
	Class() string
	Colspan() string
	Key(string, interface{}) template.HTMLAttr
}

/*
func init() {
	// check for Nota's interfaces compliance
	_ = interface { ... }(templatepipe.Nota(nil))
} */

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
