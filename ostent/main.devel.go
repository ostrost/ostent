// +build !production

package main

import (
	"flag"
	"go/build"
	"log"
	"net"
	"net/http/pprof"
	"os"

	"github.com/ostrost/ostent"
	"github.com/ostrost/ostent/commands"
	"github.com/ostrost/ostent/share/templates"
)

func main() {
	flag.Parse()

	// MAYBE the only command extract-assets is for production only
	if command := commands.ArgCommand(); command != nil {
		command()
		return
	}

	if pkg, err := build.Import("github.com/ostrost/ostent", "", build.FindOnly); err != nil {
		log.Fatal(err)
	} else if err := os.Chdir(pkg.Dir); err != nil {
		log.Fatal(err)
	}
	go templates.InitTemplates() // after chdir
	go ostent.Loop()
	// go ostent.CollectdLoop()

	listen, err := net.Listen("tcp", ostentBindFlag.String())
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(Serve(listen, false, ostent.Muxmap{
		"/debug/pprof/{name}":  pprof.Index,
		"/debug/pprof/cmdline": pprof.Cmdline,
		"/debug/pprof/profile": pprof.Profile,
		"/debug/pprof/symbol":  pprof.Symbol,
	}))
}
