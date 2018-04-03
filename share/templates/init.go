package templates

import "github.com/ostrost/ostent/templateutil"

var (
	// IndexTemplate is a templateutil.LazyTemplate of "index.html" asset.
	IndexTemplate = templateutil.NewLT(Asset, AssetInfo,
		[]string{"define_page.html", "index.html"})
)

// InitTemplates inits must-have templates.
func InitTemplates() {
	templateutil.MustInit(IndexTemplate)
}
