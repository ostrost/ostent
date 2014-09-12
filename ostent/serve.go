package main

import (
	"log"
	"net"
	"os"

	"github.com/ostrost/ostent"
	"github.com/ostrost/ostent/share/assets"
)

func Serve(listener net.Listener, production bool, extramap ostent.Muxmap) error {
	server := ostent.NewServer(listener, production)
	access := server.Access
	chain := server.Chain
	mux := server.MUX
	recovery := mux.Recovery

	logger := log.New(os.Stderr, "[ostent] ", 0)
	for _, path := range assets.AssetNames() {
		hf := chain.Then(ostent.ServeContentFunc("share/assets", path, logger))
		mux.Handle("GET", "/"+path, hf)
		mux.Handle("HEAD", "/"+path, hf)
	}

	// no logger-wrapping for slashws, because it logs by itself once a query received via websocket
	mux.Handle("GET", "/ws", recovery.ConstructorFunc(ostent.SlashwsFunc(access, periodFlag.Duration)))

	index := chain.ThenFunc(ostent.IndexFunc(periodFlag.Duration))
	mux.Handle("GET", "/", index)
	mux.Handle("HEAD", "/", index)

	/* panics := func(http.ResponseWriter, *http.Request) {
		panic(fmt.Errorf("I'm panicing"))
	}
	mux.Handle("GET", "/panic", chain.ThenFunc(panics)) // */

	ostent.Banner(listener.Addr().String(), "ostent", logger)
	return server.ServeExtra(listener, extramap)
}
