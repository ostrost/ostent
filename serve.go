package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"github.com/ostrost/ostent/cmd"
	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/share/assets"
	"github.com/ostrost/ostent/share/templates"
)

func init() {
	ostent.AddBackground(ostent.ConnectionsLoop)
	ostent.AddBackground(ostent.CollectLoop)
}

// Serve acceps incoming HTTP connections on the listener.
// Routes are set here, extra may be a hook to finalize the router.
// taggedbin is required for some handlers.
func Serve(listener net.Listener, taggedbin bool, extra func(*httprouter.Router, alice.Chain)) error {
	r, achain, access := ostent.NewServery(taggedbin)

	ostentLog := log.New(os.Stderr, "[ostent] ", 0)
	errlog, errclose := ostent.NewErrorLog()
	defer errclose()

	var (
		serve1 = ostent.NewServeSSE(access, cmd.DelayFlags)
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
