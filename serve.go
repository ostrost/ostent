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
	mux, access, accesslog := ostent.NewServer(taggedbin, extramap)

	if index := access.ThenFunc(ostent.IndexFunc(taggedbin, templates.IndexTemplate,
		PeriodFlag)); true {
		mux.Handler("GET", "/", index)
		mux.Handler("HEAD", "/", index)
	}

	if formred := ostent.FormRedirectFunc(PeriodFlag, access.ThenFunc); true {
		mux.GET("/form/*Q", formred)
		mux.POST("/form/*Q", formred)
	}

	// access chain is not used.
	// accesslog is passed to log with.
	mux.HandlerFunc("GET", "/index.ws", ostent.IndexWSFunc(accesslog, PeriodFlag))
	mux.HandlerFunc("GET", "/index.sse", ostent.IndexSSEFunc(accesslog, PeriodFlag))

	if !taggedbin { // dev-only
		mux.Handler("GET", "/panic", access.ThenFunc(
			func(w http.ResponseWriter, r *http.Request) {
				panic("/panic")
			}))
	}

	logger := log.New(os.Stderr, "[ostent] ", 0)
	for _, path := range assets.AssetNames() {
		hf := access.Then(ostent.ServeContentFunc(
			AssetReadFunc(assets.Asset),
			AssetInfoFunc(assets.AssetInfo),
			path, logger))
		mux.Handler("GET", "/"+path, hf)
		mux.Handler("HEAD", "/"+path, hf)
	}

	ostent.Banner(listener.Addr().String(), "ostent", logger)
	return http.Serve(listener, mux)
}
