// +build !production

package main
import (
	"ostential"

	"net"
	"log"
	"flag"
	"net/http"
	pprof "net/http/pprof"
)

func main() {
	flag.Parse()

	go ostential.Loop()

	listen, err := net.Listen("tcp", ostential.BindFlag.String())
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(ostential.Serve(listen, false, map[string]http.HandlerFunc{
		"/debug/pprof/{name}":  pprof.Index,
		"/debug/pprof/cmdline": pprof.Cmdline,
		"/debug/pprof/profile": pprof.Profile,
		"/debug/pprof/symbol":  pprof.Symbol,
	}))
}




