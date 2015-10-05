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
	mux, chain, access := ostent.NewServery(taggedbin)
	errlog, errclose := ostent.NewErrorLog()
	defer errclose()

	if DelayFlags.Max.Duration < DelayFlags.Min.Duration {
		DelayFlags.Max.Duration = DelayFlags.Min.Duration
	}

	ss := ostent.NewServeSSE(access, DelayFlags)
	mux.Handler("GET", "/index.sse", http.HandlerFunc(ss.IndexSSE))
	sw := ostent.NewServeWS(ss, errlog)
	mux.Handler("GET", "/index.ws", http.HandlerFunc(sw.IndexWS))

	si := ostent.NewServeIndex(sw, taggedbin, templates.IndexTemplate)
	if p, h := "/", chain.ThenFunc(si.Index); true {
		mux.Handler("GET", p, h)
		mux.Handler("HEAD", p, h)
	}

	if !taggedbin { // dev-only
		if p, h := "/panic", chain.ThenFunc(
			func(w http.ResponseWriter, r *http.Request) {
				panic("/panic")
			}); true {
			mux.Handler("GET", p, h)
			mux.Handler("HEAD", p, h)
		}
	}

	sa := ostent.ServeAssets{
		Log:                 log.New(os.Stderr, "[ostent] ", 0),
		AssetFunc:           assets.Asset,
		AssetInfoFunc:       assets.AssetInfo,
		AssetAltModTimeFunc: AssetAltModTimeFunc, // from main.*.go
	}
	for _, path := range assets.AssetNames() {
		p := "/" + path
		if path != "favicon.ico" && path != "robots.txt" {
			p = "/" + ostent.VERSION + "/" + path // the Version prefix
		}
		cchain := chain.Append(context.ClearHandler, ostent.AddAssetPathContextFunc(path))
		h := cchain.ThenFunc(sa.Serve)
		mux.Handler("GET", p, h)
		mux.Handler("HEAD", p, h)
	}
	if extra != nil {
		extra(mux, chain)
	}

	ostent.Banner(listener.Addr().String(), "ostent", sa.Log)
	s := &http.Server{
		ErrorLog: errlog,
		Handler:  mux,
	}
	return s.Serve(listener)
}
