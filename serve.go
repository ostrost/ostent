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

// PeriodFlag is a minimum refresh period for collection.
var PeriodFlag = flags.Period{Duration: time.Second} // default

func init() {
	flag.Var(&PeriodFlag, "u", "Collection (update) interval")
	flag.Var(&PeriodFlag, "update", "Collection (update) interval")
}

func Serve(listener net.Listener, taggedbin bool, extramap ostent.Muxmap) error {
	mux, chain, access := ostent.NewServery(taggedbin, extramap)
	errlog, errclose := ostent.NewErrorLog()
	defer errclose()

	if index := chain.ThenFunc(ostent.IndexFunc(taggedbin, templates.IndexTemplate,
		PeriodFlag)); true {
		mux.Handler("GET", "/", index)
		mux.Handler("HEAD", "/", index)
	}

	if formred := ostent.FormRedirectFunc(PeriodFlag, chain.ThenFunc); true {
		mux.GET("/form/*Q", formred)
		mux.POST("/form/*Q", formred)
	}

	// chain is not used -- access is passed to log with.
	mux.HandlerFunc("GET", "/index.ws", ostent.IndexWSFunc(access, errlog, PeriodFlag))
	mux.HandlerFunc("GET", "/index.sse", ostent.IndexSSEFunc(access, PeriodFlag))

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
