package templates

import (
	"github.com/ostrost/ostent/templates"
	"github.com/rzab/amber"
)

var (
	UsePercentTemplate  = templates.BinTemplate{Readfunc: Asset, Filename: "usepercent.html"}
	TooltipableTemplate = templates.BinTemplate{Readfunc: Asset, Filename: "tooltipable.html"}
	IndexTemplate       = templates.BinTemplate{Readfunc: Asset, Filename: "index.html", Cascade: true, Funcmap: amber.FuncMap}
)

func InitTemplates() {
	UsePercentTemplate.MustInit()
	TooltipableTemplate.MustInit()
	IndexTemplate.MustInit()
}
