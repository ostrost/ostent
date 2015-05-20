package templates

import (
	"github.com/ostrost/ostent/acepp/templatep"
	"github.com/ostrost/ostent/templateutil"
)

var (
	// IndexTemplate is a templateutil.LazyTemplate of "index.html" asset.
	IndexTemplate = templateutil.NewLT(Asset, "index.html", templatep.AceFuncs)
)

// InitTemplates inits must-have templates and signals done when finished.
func InitTemplates(done chan<- struct{}) {
	templateutil.MustInit(IndexTemplate)
	if done != nil {
		done <- struct{}{}
	}
}
