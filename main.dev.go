// +build !bin

package main

import (
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/spf13/cobra"

	"github.com/ostrost/ostent/cmd"
	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/share/templates"
)

// AssetAltModTimeFunc is nil.
var AssetAltModTimeFunc func() time.Time

func OstentRunE(*cobra.Command, []string) error {
	listen, err := net.Listen("tcp", cmd.OstentBind.String())
	if err != nil {
		return err
	}

	ostent.RunBackground()

	templatesLoaded := make(chan struct{}, 1)
	go templates.InitTemplates(templatesLoaded)

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
			return err
		}
	}
	return nil
}

func main() {
	cmd.OstentCmd.RunE = OstentRunE
	cmd.Execute()
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
