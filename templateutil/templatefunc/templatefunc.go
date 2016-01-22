package templatefunc

import (
	"html/template"

	"github.com/ostrost/ostent/params"
	"github.com/ostrost/ostent/templateutil"
)

func FuncMapHTML() template.FuncMap {
	pfuncs := params.ParamsFuncs{}
	return templateutil.CombineMaps(templateutil.NewHTMLFuncs(), template.FuncMap{
		"HrefT": pfuncs.HrefT,
		"LessD": pfuncs.LessD,
		"MoreD": pfuncs.MoreD,
		"LessN": pfuncs.LessN,
		"MoreN": pfuncs.MoreN,
		"Vlink": pfuncs.Vlink,
	})
}
