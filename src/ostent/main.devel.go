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
	log.Fatal(ostential.Serve(listen, false, func(mux *http.ServeMux) {
		mux.HandleFunc("/debug/pprof/",        pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol",  pprof.Symbol)
	}))
}




