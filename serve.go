package main

import (
	"flag"
	"log"
	"net"
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
	server := ostent.NewServer(listener, taggedbin)
	access := server.Access
	chain := server.Chain
	mux := server.MUX
	recovery := mux.Recovery

	logger := log.New(os.Stderr, "[ostent] ", 0)
	assetnames := assets.AssetNames()
	for _, path := range assetnames {
		hf := chain.Then(ostent.ServeContentFunc(
			AssetReadFunc(assets.Asset),
			AssetInfoFunc(assets.AssetInfo),
			path, logger))
		mux.Handle("GET", "/"+path, hf)
		mux.Handle("HEAD", "/"+path, hf)
	}

	// access is passed to Index*Func for them to log with.
	// recovery.ConstructorFunc used to bypass the chain so no double log.
	mux.Handle("GET", "/index.ws", recovery.
		ConstructorFunc(ostent.IndexWSFunc(access, PeriodFlag)))
	mux.Handle("GET", "/index.sse", recovery.
		ConstructorFunc(ostent.IndexSSEFunc(access, PeriodFlag)))

	index := chain.ThenFunc(ostent.IndexFunc(taggedbin,
		templates.IndexTemplate, PeriodFlag))
	mux.Handle("GET", "/", index)
	mux.Handle("HEAD", "/", index)

	/* panics := func(http.ResponseWriter, *http.Request) {
		panic(fmt.Errorf("I'm panicing"))
	}
	mux.Handle("GET", "/panic", chain.ThenFunc(panics)) // */

	ostent.Banner(listener.Addr().String(), "ostent", logger)
	return server.ServeExtra(listener, extramap)
}
