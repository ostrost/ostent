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

var (
	// MinDelayFlag is a minimum refresh period for collection.
	MinDelayFlag = flags.Delay{Duration: time.Second} // default
	// MaxDelayFlag is a maximum refresh period for collection.
	MaxDelayFlag = flags.Delay{Duration: 10 * time.Minute} // default
)

func init() {
	flag.Var(&MinDelayFlag, "d", "Short for min-delay")
	flag.Var(&MinDelayFlag, "min-delay", "Collection and minimum for UI `delay`")
	flag.Var(&MaxDelayFlag, "max-delay", "Maximum for UI `delay`")
	ostent.AddBackground(ostent.ConnectionsLoop)
	ostent.AddBackground(ostent.CollectLoop)
}

func Serve(listener net.Listener, taggedbin bool, extra func(*httprouter.Router, alice.Chain)) error {
	mux, chain, access := ostent.NewServery(taggedbin)
	errlog, errclose := ostent.NewErrorLog()
	defer errclose()

	if MaxDelayFlag.Duration < MinDelayFlag.Duration {
		MaxDelayFlag.Duration = MinDelayFlag.Duration
	}

	ss := ostent.NewServeSSE(access, MinDelayFlag)
	mux.Handler("GET", "/index.sse", http.HandlerFunc(ss.IndexSSE))
	sw := ostent.NewServeWS(ss, errlog, MaxDelayFlag)
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
