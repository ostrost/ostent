package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"github.com/ostrost/ostent/flags"
	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/share/assets"
	"github.com/ostrost/ostent/share/templates"
)

var DelayFlags = flags.DelayBounds{
	Max: flags.Delay{Duration: 10 * time.Minute},
	Min: flags.Delay{Duration: time.Second},
	// 10m and 1s are corresponding defaults
}

func init() {
	flag.Var(&DelayFlags.Max, "max-delay", "Maximum for UI `delay`")
	flag.Var(&DelayFlags.Min, "min-delay", "Collect and minimum for UI `delay`")
	flag.Var(&DelayFlags.Min, "d", "Short for min-delay")
	ostent.AddBackground(ostent.ConnectionsLoop)
	ostent.AddBackground(ostent.CollectLoop)
}

func Serve(listener net.Listener, taggedbin bool, extra func(*httprouter.Router, alice.Chain)) error {
	// post-flag.Parse() really
	if DelayFlags.Max.Duration < DelayFlags.Min.Duration {
		DelayFlags.Max.Duration = DelayFlags.Min.Duration
	}

	r, achain, access := ostent.NewServery(taggedbin)

	ostentLog := log.New(os.Stderr, "[ostent] ", 0)
	errlog, errclose := ostent.NewErrorLog()
	defer errclose()

	var (
		serve1 = ostent.NewServeSSE(access, DelayFlags)
		serve2 = ostent.NewServeWS(serve1, errlog)
		serve3 = ostent.NewServeIndex(serve2, taggedbin, templates.IndexTemplate)
		serve4 = ostent.ServeAssets{
			Log:                 ostentLog,
			AssetFunc:           assets.Asset,
			AssetInfoFunc:       assets.AssetInfo,
			AssetAltModTimeFunc: AssetAltModTimeFunc, // from main.*.go
		}
	)

	m, n, chain0 := ostent.NewRoute, ostent.NewHandle, alice.New()

	var (
		panicp = "/panic"
		panics = func(http.ResponseWriter, *http.Request) { panic(panicp) }
		panicr = m(panicp, r.GET, r.HEAD)
	)

	routes := map[*ostent.Route]ostent.Handle{
		m("/index.sse", r.GET): n(chain0, serve1.IndexSSE),
		m("/index.ws", r.GET):  n(chain0, serve2.IndexWS),
		m("/", r.GET, r.HEAD):  n(achain, serve3.Index),
		panicr:                 n(achain, panics),
	}
	for _, path := range assets.AssetNames() {
		p := "/" + path
		if path != "favicon.ico" && path != "robots.txt" {
			p = "/" + ostent.VERSION + "/" + path // the Version prefix
		}
		rr := m(p, r.GET, r.HEAD)
		rr.Asset = true
		routes[rr] = n(achain.Append(context.ClearHandler,
			ostent.AddAssetPathContextFunc(path)), serve4.Serve)
	}
	if taggedbin { // no panicp in bin
		delete(routes, panicr)
	}

	// now bind
	ostent.ApplyRoutes(r, routes, nil)
	if extra != nil {
		extra(r, achain)
	}

	ostent.Banner(listener.Addr().String(), "ostent", ostentLog)
	s := &http.Server{
		ErrorLog: errlog,
		Handler:  r,
	}
	return s.Serve(listener)
}
