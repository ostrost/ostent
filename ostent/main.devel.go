// +build !production

package main

import (
	"flag"
	"go/build"
	"log"
	"net/http/pprof"
	"os"

	"github.com/ostrost/ostent"
	"github.com/ostrost/ostent/commands"
	"github.com/ostrost/ostent/share/templates"
)

func main() {
	webserver := commands.FlagSetNewWebserver(flag.CommandLine)
	flag.Parse()

	if errd := commands.ArgCommands(); errd { // explicit commands
		return
	}

	if pkg, err := build.Import("github.com/ostrost/ostent", "", build.FindOnly); err != nil {
		log.Fatal(err)
		// chdir for templates loading
	} else if err := os.Chdir(pkg.Dir); err != nil {
		log.Fatal(err)
	}
	// the background job(s)
	go ostent.Loop()
	// go ostent.CollectdLoop()

	go templates.InitTemplates() // ServeFunc; NB after chdir

	listen := webserver.NetListen()
	log.Fatal(Serve(listen, false, ostent.Muxmap{
		"/debug/pprof/{name}":  pprof.Index,
		"/debug/pprof/cmdline": pprof.Cmdline,
		"/debug/pprof/profile": pprof.Profile,
		"/debug/pprof/symbol":  pprof.Symbol,
	}))
}
