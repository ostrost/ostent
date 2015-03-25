package templates

import (
	"github.com/ostrost/ostent/templateutil"
	"github.com/rzab/amber"
)

var (
	// IndexTemplate is a templateutil.LazyTemplate of "index.html" asset.
	IndexTemplate = templateutil.NewLT(Asset, "index.html", amber.FuncMap)
	// DefinesTemplate is a templateutil.LazyTemplate of "defines.html" asset.
	DefinesTemplate = templateutil.NewLT(Asset, "defines.html", amber.FuncMap)
)

func InitTemplates(done chan<- struct{}) {
	templateutil.MustInit(IndexTemplate)
	templateutil.MustInit(DefinesTemplate)
	if done != nil {
		done <- struct{}{}
	}
}
