package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"time"

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
	sw := ostent.NewServeWS(*ss, errlog, MaxDelayFlag)
	mux.Handler("GET", "/index.ws", http.HandlerFunc(sw.IndexWS))

	si := ostent.NewServeIndex(*sw, taggedbin, templates.IndexTemplate)
	indexHandler := chain.ThenFunc(si.Index)
	mux.Handler("GET", "/", indexHandler)
	mux.Handler("HEAD", "/", indexHandler)

	if !taggedbin { // dev-only
		mux.Handler("GET", "/panic", chain.ThenFunc(
			func(w http.ResponseWriter, r *http.Request) {
				panic("/panic")
			}))
	}

	logger := log.New(os.Stderr, "[ostent] ", 0)
	for _, path := range assets.AssetNames() {
		hf := chain.Then(ostent.ServeContentFunc(
			AssetReadFunc(assets.Asset),
			AssetInfoFunc(assets.AssetInfo),
			path, logger))
		mux.Handler("GET", "/"+path, hf)
		mux.Handler("HEAD", "/"+path, hf)
	}

	ostent.Banner(listener.Addr().String(), "ostent", logger)
	s := &http.Server{
		ErrorLog: errlog,
		Handler:  mux,
	}
	return s.Serve(listener)
}
