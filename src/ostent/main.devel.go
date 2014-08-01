// +build !production

package main
import (
	"ostential"

	"net"
	"log"
	"flag"
	pprof "net/http/pprof"
)

func main() {
	flag.Parse()

	go ostential.Loop()
	// go ostential.CollectdLoop()

	listen, err := net.Listen("tcp", ostential.OstentBindFlag.String())
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(ostential.Serve(listen, false, ostential.Muxmap{
		"/debug/pprof/{name}":  pprof.Index,
		"/debug/pprof/cmdline": pprof.Cmdline,
		"/debug/pprof/profile": pprof.Profile,
		"/debug/pprof/symbol":  pprof.Symbol,
	}))
}




