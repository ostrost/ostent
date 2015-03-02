package assets

import (
	"net/http"
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

func JsAssetNames(assetnames []string) []string {
	sa := sortassets{
		substr_indexfrom: []string{
			"jquery",
			"bootstrap",       // depends on jquery
			"d3",              // no dependencies?
			"metricsgraphics", // depends on d3, bootstrap

			"react",    // no dependencies
			"headroom", // no dependencies? maybe jquery

			"gen", "jsript", // either /gen/ or /jscript/
			"milk", // from coffee script
		},
	}
	for _, name := range assetnames {
		if strings.HasSuffix(name, ".js") {
			sa.names = append(sa.names, "/"+name)
		}
	}
	sort.Stable(sa)
	return sa.names
}
