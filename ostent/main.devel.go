// +build !production

package main

import (
	"flag"
	"go/build"
	"log"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/ostrost/ostent"
	"github.com/ostrost/ostent/commands"
	"github.com/ostrost/ostent/share/templates"
)

func main() {
	flag.Usage = commands.UsageFunc(flag.CommandLine)
	webserver := commands.NewWebserver().AddCommandLine()
	flag.Parse()

	errd, atexit := commands.ArgCommands()
	defer atexit()

	if errd {
		return
	}

	// Chdir into ostent package directory for templates loading
	if pkg, err := build.Import("github.com/ostrost/ostent", "", build.FindOnly); err != nil {
		log.Fatal(err)
	} else if err := os.Chdir(pkg.Dir); err != nil {
		log.Fatal(err)
	}
	// the background job(s)
	go ostent.Loop()
	// go ostent.CollectdLoop()

	go templates.InitTemplates() // ServeFunc; NB after chdir

	listen := webserver.NetListen()
	errch := make(chan error, 2)
	go func(ch chan<- error) {
		ch <- Serve(listen, false, ostent.Muxmap{
			"/debug/pprof/{name}":  pprof.Index,
			"/debug/pprof/cmdline": pprof.Cmdline,
			"/debug/pprof/profile": pprof.Profile,
			"/debug/pprof/symbol":  pprof.Symbol,
		})
	}(errch)
	sigch := make(chan os.Signal, 2)
	signal.Notify(sigch,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
wait:
	for {
		select {
		case _ = <-sigch:
			break wait
		case err := <-errch:
			log.Fatal(err)
		}
	}
}
