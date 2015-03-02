//go:generate gen

// Package assetutil provides asset utilities.
package assetutil

import (
	"net/http"
	"strings"
	"time"
)

// TimeInfo is for *Asset{Info,Read}Func: a reduced os.FileInfo.
type TimeInfo interface {
	ModTime() time.Time
}

// FQJSANSlice forms absolute urls with r.Host.
func FQJSANSlice(js JSANSlice, r *http.Request) JSANSlice {
	return js.SelectJSAN(func(n JSAN) JSAN {
		if !strings.HasPrefix(string(n), "//") {
			return JSAN("//"+r.Host) + n
		}
		return n
	})
}

// LessJSANFunc makes a Less func for JSAN.
func LessJSANFunc(max int) func(JSAN, JSAN) bool {
	return func(a, b JSAN) bool {
		ii, jj := max, max
		for w, v := range JSORDER {
			if strings.Contains(string(a), v) {
				ii = w
			}
			if strings.Contains(string(b), v) {
				jj = w
			}
		}
		return ii < jj
	}
}

// JSORDER defines weight in ordering js filenames/assetnames.
var JSORDER = []string{
	"jquery",
	"bootstrap",       // depends on jquery
	"d3",              // no dependencies?
	"metricsgraphics", // depends on d3, bootstrap

	"react",    // no dependencies
	"headroom", // no dependencies? maybe jquery

	"gen", "jsript", // either /gen/ or /jscript/
	"milk", // from coffee script
}

// JSassetNames filters and sorts JSANs.
func JSassetNames(assetnames []string) (js JSANSlice) {
	for _, name := range assetnames { // convert from []string
		js = append(js, JSAN(name))
	}
	js = js.Where(func(n JSAN) bool { return strings.HasSuffix(string(n), ".js") })
	js.SortSortBy(LessJSANFunc(len(js))) // sort after filter and before the map
	return js.SelectJSAN(func(n JSAN) JSAN { return JSAN("/") + n })
}

// JSAN derives string for filtering and sorting JS asset names.
// +gen slice:"Select[JSAN],Where,PkgSortBy"
type JSAN string
