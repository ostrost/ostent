// +build !production

package main

import (
	"flag"
	"log"
	"net"
	"net/http/pprof"

	"github.com/ostrost/ostent"
)

func main() {
	flag.Parse()

	go ostent.Loop()
	// go ostent.CollectdLoop()

	listen, err := net.Listen("tcp", ostent.OstentBindFlag.String())
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(ostent.Serve(listen, false, ostent.Muxmap{
		"/debug/pprof/{name}":  pprof.Index,
		"/debug/pprof/cmdline": pprof.Cmdline,
		"/debug/pprof/profile": pprof.Profile,
		"/debug/pprof/symbol":  pprof.Symbol,
	}))
}
