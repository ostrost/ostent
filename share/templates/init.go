package templates

import "github.com/ostrost/ostent/templateutil"

var (
	// IndexTemplate is a templateutil.LazyTemplate of "index.html" asset.
	IndexTemplate = templateutil.NewLT(Asset, AssetInfo, "index.html", nil)
)

// InitTemplates inits must-have templates.
func InitTemplates() {
	templateutil.MustInit(IndexTemplate)
}
