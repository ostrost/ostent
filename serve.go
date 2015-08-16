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

// DelayFlag is a minimum refresh period for collection.
var DelayFlag = flags.Delay{Duration: time.Second} // default

func init() {
	flag.Var(&DelayFlag, "d", "Short for delay")
	flag.Var(&DelayFlag, "delay", "Collection `delay`")
}

func Serve(listener net.Listener, taggedbin bool, extramap ostent.Muxmap) error {
	mux, chain, access := ostent.NewServery(taggedbin, extramap)
	errlog, errclose := ostent.NewErrorLog()
	defer errclose()

	if index := chain.ThenFunc(ostent.IndexFunc(taggedbin, templates.IndexTemplate,
		DelayFlag)); true {
		mux.Handler("GET", "/", index)
		mux.Handler("HEAD", "/", index)
	}

	// chain is not used -- access is passed to log with.
	mux.HandlerFunc("GET", "/index.ws", ostent.IndexWSFunc(access, errlog, DelayFlag))
	mux.HandlerFunc("GET", "/index.sse", ostent.IndexSSEFunc(access, DelayFlag))

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
