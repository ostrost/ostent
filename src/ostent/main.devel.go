// +build !production

package main
import (
	"ostential"

	"os"
	"net"
	"log"
	"flag"
	pprof "net/http/pprof"

	"github.com/codegangsta/martini"
)

func main() {
	flag.Parse()

	os.Setenv("HOST", ostential.BindFlag.Host) // for martini
	os.Setenv("PORT", ostential.BindFlag.Port) // for martini

	martini.Env = martini.Dev
	go ostential.Loop()

	listen, err := net.Listen("tcp", ostential.BindFlag.String())
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(ostential.Serve(listen, ostential.LogAll, func(m *ostential.Modern) {
		m.Any("/debug/pprof/cmdline", pprof.Cmdline)
		m.Any("/debug/pprof/profile", pprof.Profile)
		m.Any("/debug/pprof/symbol",  pprof.Symbol)
		m.Any("/debug/pprof/.*",      pprof.Index)
	}))
}




