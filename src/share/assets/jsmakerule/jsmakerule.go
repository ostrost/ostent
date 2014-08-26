// +build !production

package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ostrost/ostent/src/share/assets"
)

const packageName = "src/share/assets"

func main() {
	flag.Parse()
	target := flag.Arg(0)

	var lines []string
	for _, line := range assets.JsAssetNames() {
		lines = append(lines, filepath.Join(packageName, line))
	}

	if target != "" {
		fmt.Printf("%s: %s\n", target, strings.Join(lines, " "))
		return
	}
	for _, line := range lines {
		fmt.Println(line)
	}
}
