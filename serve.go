package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/context"

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

func Serve(listener net.Listener, taggedbin bool, extramap ostent.Muxmap) error {
	mux, chain, access := ostent.NewServery(taggedbin, extramap)
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
	ostent.GEAD(mux, "/", chain.ThenFunc(si.Index))

	if !taggedbin { // dev-only
		mux.Handler("GET", "/panic", chain.ThenFunc(
			func(w http.ResponseWriter, r *http.Request) {
				panic("/panic")
			}))
	}

	sa := ostent.ServeAssets{
		Log:                 log.New(os.Stderr, "[ostent] ", 0),
		AssetFunc:           assets.Asset,
		AssetInfoFunc:       assets.AssetInfo,
		AssetAltModTimeFunc: AssetAltModTimeFunc, // from main.*.go
	}
	for _, path := range assets.AssetNames() {
		pattern := path
		if path != "favicon.ico" && path != "robots.txt" {
			pattern = ostent.VERSION + "/" + path // the Version prefix
		}
		cchain := chain.Append(context.ClearHandler, ostent.AddAssetPathContextFunc(path))
		ostent.GEAD(mux, "/"+pattern, cchain.ThenFunc(sa.Serve))
	}

	ostent.Banner(listener.Addr().String(), "ostent", sa.Log)
	s := &http.Server{
		ErrorLog: errlog,
		Handler:  mux,
	}
	return s.Serve(listener)
}
