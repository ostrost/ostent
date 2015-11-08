// +build !bin

package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"github.com/ostrost/ostent/commands"
	_ "github.com/ostrost/ostent/commands/ostent"
	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/share/templates"
)

// AssetAltModTimeFunc is nil.
var AssetAltModTimeFunc func() time.Time

func main() {
	webserver := commands.NewWebserver(8050).AddCommandLine()
	commands.Parse(flag.CommandLine, os.Args[1:])

	errd, atexit := commands.ArgCommands()
	defer atexit()

	if errd {
		return
	}

	ostent.RunBackground()

	templatesLoaded := make(chan struct{}, 1)
	go templates.InitTemplates(templatesLoaded)

	listen := webserver.NetListen()
	errch := make(chan error, 2)
	go func(ch chan<- error) {
		<-templatesLoaded
		ch <- Serve(listen, false, PprofExtra)
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

func (ps PprofServe) Serve(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	name := params.ByName("name")
	var handler func(http.ResponseWriter, *http.Request)
	switch name {
	case "cmdline":
		handler = pprof.Cmdline
	case "profile":
		handler = pprof.Profile
	case "symbol":
		handler = pprof.Symbol
	// TODO case "trace": handler = pprof.Trace // in go1.5
	default:
		handler = pprof.Index
	}
	if handler == nil {
		http.NotFound(w, r)
		return
	}
	ps.Chain.ThenFunc(handler).ServeHTTP(w, r)
}

type PprofServe struct{ Chain alice.Chain }

func PprofExtra(r *httprouter.Router, chain alice.Chain) {
	handle := PprofServe{chain}.Serve
	r.GET("/debug/pprof/:name", handle)
	r.HEAD("/debug/pprof/:name", handle)
	r.POST("/debug/pprof/:name", handle)
}
