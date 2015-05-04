package templates

import (
	"github.com/ostrost/ostent/acepp/templatep"
	"github.com/ostrost/ostent/templateutil"
)

var (
	// IndexTemplate is a templateutil.LazyTemplate of "index.html" asset.
	IndexTemplate = templateutil.NewLT(Asset, "index.html", templatep.AceFuncs)
	// DefinesTemplate is a templateutil.LazyTemplate of "defines.html" asset.
	DefinesTemplate = templateutil.NewLT(Asset, "defines.html", templatep.AceFuncs)
)

func InitTemplates(done chan<- struct{}) {
	templateutil.MustInit(IndexTemplate)
	templateutil.MustInit(DefinesTemplate)
	if done != nil {
		done <- struct{}{}
	}
}
