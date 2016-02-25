// +build !bin

package main

import (
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/spf13/cobra"

	"github.com/ostrost/ostent/cmd"
)

// AssetAltModTimeFunc is nil.
var AssetAltModTimeFunc func() time.Time

func main() {
	cmd.OstentCmd.RunE = func(*cobra.Command, []string) error {
		return Serve(cmd.OstentBind.String(), false, PprofExtra)
	}
	cmd.Execute()
}

// Serve is PprofServe handler.
func (ps PprofServe) Serve(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	name := params.ByName("name")
	var handler func(http.ResponseWriter, *http.Request)
	switch name {
	case "cmdline":
		handler = pprof.Cmdline
	case "profile":
		handler = pprof.Profile
	case "symbol":
		handler = pprof.Symbol
	case "trace":
		handler = pprof.Trace
	default:
		handler = pprof.Index
	}
	if handler == nil {
		http.NotFound(w, r)
		return
	}
	ps.Chain.ThenFunc(handler).ServeHTTP(w, r)
}

// PprofServe is pprof data serving handler.
type PprofServe struct{ Chain alice.Chain }

// PprofExtra is a hook to add PprofServe-handled routes to r.
func PprofExtra(r *httprouter.Router, chain alice.Chain) {
	handle := PprofServe{chain}.Serve
	r.GET("/debug/pprof/:name", handle)
	r.HEAD("/debug/pprof/:name", handle)
	r.POST("/debug/pprof/:name", handle)
}
