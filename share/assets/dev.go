// +build !bin

package assets

import (
	"go/build"
	"log"
	"time"
)

// ThisPkgPath defined for looking up the package directory.
const ThisPkgPath = "github.com/ostrost/ostent/share/assets"

var (
	// rootDir is this package directory.
	rootDir string
	// AssetAltModTimeFunc is nil.
	AssetAltModTimeFunc func() time.Time
)

func init() {
	pkg, err := build.Import(ThisPkgPath, "", build.FindOnly)
	if err != nil {
		log.Fatal(err)
	}
	rootDir = pkg.Dir
}
