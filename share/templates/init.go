package templates

import (
	"github.com/ostrost/ostent/templateutil"
	"github.com/rzab/amber"
)

var (
	UsePercentTemplate  = &templateutil.BinTemplate{Readfunc: Asset, Filename: "usepercent.html"}
	TooltipableTemplate = &templateutil.BinTemplate{Readfunc: Asset, Filename: "tooltipable.html"}
	IndexTemplate       = &templateutil.BinTemplate{Readfunc: Asset, Filename: "index.html", Cascade: true, Funcmap: amber.FuncMap}
)

func InitTemplates(done chan<- struct{}) {
	templateutil.MustInit(UsePercentTemplate)
	templateutil.MustInit(TooltipableTemplate)
	templateutil.MustInit(IndexTemplate)
	if done != nil {
		done <- struct{}{}
	}
}
