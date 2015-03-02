// +build !production

package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ostrost/ostent/assets"
	shareassets "github.com/ostrost/ostent/share/assets"
)

const packageName = "share/assets"

func main() {
	flag.Parse()
	target := flag.Arg(0)

	var lines []string
	for _, line := range assets.JSassetNames(shareassets.AssetNames()) {
		lines = append(lines, filepath.Join(packageName, string(line)))
	}

	if target != "" {
		fmt.Printf("%s: %s\n", target, strings.Join(lines, " "))
		return
	}
	for _, line := range lines {
		fmt.Println(line)
	}
}
