package assets

import (
	"fmt"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
)

func FQscripts(list []string, r *http.Request) (scripts []string) {
	for _, s := range list {
		if !strings.HasPrefix(string(s), "//") {
			s = "//" + r.Host + s
		}
		scripts = append(scripts, s)
	}
	return scripts
}

type sortassets struct {
	names            []string
	substr_indexfrom []string
}

func (sa sortassets) Len() int {
	return len(sa.names)
}

func (sa sortassets) Less(i, j int) bool {
	ii, jj := sa.Len(), sa.Len()
	for w, v := range sa.substr_indexfrom {
		if strings.Contains(sa.names[i], v) {
			ii = w
		}
		if strings.Contains(sa.names[j], v) {
			jj = w
		}
	}
	return ii < jj
}

func (sa sortassets) Swap(i, j int) {
	sa.names[i], sa.names[j] = sa.names[j], sa.names[i]
}

func JsAssetNames(assetnames []string, develreact bool) []string {
	sa := sortassets{
		substr_indexfrom: []string{
			"jquery",
			"bootstrap",
			"react",
			"headroom",

			"gen", "jsript", // either /gen/ or /jscript/
			"milk", // from coffee script
		},
	}

	for _, name := range assetnames {
		const dotjs = ".js"
		if !strings.HasSuffix(name, dotjs) {
			continue
		}
		src := "/" + name
		if develreact && strings.Contains(src, "react") {
			ver := filepath.Base(filepath.Dir(src))
			base := filepath.Base(src)

			cutlen := len(dotjs) // asserted strings.HasSuffix(base, dotjs)
			cutlen += map[bool]int{true: len(".min")}[strings.HasSuffix(base[:len(base)-cutlen], ".min")]
			src = fmt.Sprintf("//fb.me/%s-%s%s", base[:len(base)-cutlen], ver, dotjs)
		}
		sa.names = append(sa.names, src)
	}

	sort.Stable(sa)
	return sa.names
}
