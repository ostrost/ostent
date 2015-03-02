package main

import (
	"log"
	"net"
	"os"

	"github.com/ostrost/ostent/assetutil"
	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/share/assets"
	sharetemplates "github.com/ostrost/ostent/share/templates"
)

func init() {
	ostent.UsePercentTemplate = sharetemplates.UsePercentTemplate
	ostent.TooltipableTemplate = sharetemplates.TooltipableTemplate
}

func Serve(listener net.Listener, production bool, extramap ostent.Muxmap) error {
	server := ostent.NewServer(listener, production)
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

	// access is passed to log every query received via websocket
	// recovery.ConstructorFunc used to bypass chain so no double log
	mux.Handle("GET", "/ws", recovery.
		ConstructorFunc(ostent.SlashwsFunc(access, periodFlag.Duration)))

	index := chain.ThenFunc(ostent.IndexFunc(sharetemplates.IndexTemplate,
		assetutil.JSassetNames(assetnames), periodFlag.Duration))
	mux.Handle("GET", "/", index)
	mux.Handle("HEAD", "/", index)

	/* panics := func(http.ResponseWriter, *http.Request) {
		panic(fmt.Errorf("I'm panicing"))
	}
	mux.Handle("GET", "/panic", chain.ThenFunc(panics)) // */

	ostent.Banner(listener.Addr().String(), "ostent", logger)
	return server.ServeExtra(listener, extramap)
}
