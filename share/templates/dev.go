// +build !bin

package templates

import (
	"go/build"
	"log"
)

// ThisPkgPath defined for looking up the package directory.
const ThisPkgPath = "github.com/ostrost/ostent/share/templates"

var rootDir string

func init() {
	pkg, err := build.Import(ThisPkgPath, "", build.FindOnly)
	if err != nil {
		log.Fatal(err)
	}
	rootDir = pkg.Dir
}
