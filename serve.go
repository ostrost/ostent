package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/context"
	"github.com/justinas/alice"

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

	indexChain := chain.Append(context.ClearHandler,
		ostent.AddContext(ostent.CIndexTemplate, templates.IndexTemplate),
		ostent.AddContext(ostent.CMinDelay, DelayFlag),
		ostent.AddContext(ostent.CTaggedBin, taggedbin))

	indexHandler := indexChain.ThenFunc(ostent.Index)
	mux.Handler("GET", "/", indexHandler)
	mux.Handler("HEAD", "/", indexHandler)

	// chain is not used -- access to log with is passed in context.
	wschain := alice.New(context.ClearHandler,
		ostent.AddContext(ostent.CAccess, access),
		ostent.AddContext(ostent.CErrorLog, errlog),
		ostent.AddContext(ostent.CMinDelay, DelayFlag))
	mux.Handler("GET", "/index.ws", wschain.ThenFunc(ostent.IndexWS))
	mux.Handler("GET", "/index.sse", wschain.ThenFunc(ostent.IndexSSE))

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
