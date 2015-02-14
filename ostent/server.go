package ostent

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/justinas/alice"
	"github.com/ostrost/ostent/assets"
)

// Muxmap is a type of a map of pattern to HandlerFunc.
type Muxmap map[string]http.HandlerFunc

// Server is a http.Server with auxiliaries
type Server struct {
	http.Server
	Access *logger       // the access_log
	MUX    *TrieServeMux // the mux
	Chain  alice.Chain   // the chain
}

// ServeExtra .Serve's with the extramap in the mux.
func (s *Server) ServeExtra(listener net.Listener, extramap Muxmap) error {
	if extramap != nil {
		for path, handler := range extramap {
			for _, METH := range []string{"HEAD", "GET", "POST"} {
				s.MUX.Handle(METH, path, s.Chain.Then(handler))
			}
		}
	}
	return s.Serve(listener)
}

// NewServer creates a Server.
func NewServer(listener net.Listener, production bool) *Server {
	access := newLogged(production, log.New(os.Stdout, "", 0))
	recovery := recovery(production)
	chain := alice.New(
		access.Constructor,
		recovery.Constructor,
	)
	mux := NewMux(recovery, chain.Then)
	return &Server{
		Server: http.Server{Addr: listener.Addr().String(), Handler: mux},
		Access: access,
		Chain:  chain,
		MUX:    mux,
	}
}

// ServeContentFunc does http.ServeContent the readfunc (Asset or UncompressedAsset) result
func ServeContentFunc(prefix string, readfunc func(string) ([]byte, error), path string, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		text, err := readfunc(path)
		if err != nil {
			panic(err)
		}
		modtime, err := assets.ModTime(prefix, path)
		if err != nil {
			logger.Println(err)
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			// return
		}

		http.ServeContent(w, r, path, modtime, bytes.NewReader(text))
	}
}

type loggerPrint interface {
	Print(v ...interface{})
}

// Banner prints a banner with the logger.
func Banner(listenaddr string, suffix string, logger loggerPrint) {
	hostname, _ := (&Machine{}).Hostname()
	var addrsp *[]net.Addr
	if addrs, err := net.InterfaceAddrs(); err == nil {
		addrsp = &addrs
	} else {
		logger.Print(fmt.Sprintf("%s\n", err))
	}
	bannerText(listenaddr, hostname, suffix, addrsp, logger)
}

func bannerText(listenaddr, hostname, suffix string, addrsp *[]net.Addr, logger loggerPrint) {
	if limit := 32 /* width */ - 6 /* const chars */ - len(suffix); len(hostname) >= limit {
		hostname = hostname[:limit-4] + "..."
	}
	logger.Print(fmt.Sprintf("   %s\n", strings.Repeat("-", len(hostname)+1 /* space */ +len(suffix))))
	logger.Print(fmt.Sprintf(" / %s %s \\\n", hostname, suffix))
	logger.Print("+------------------------------+\n")

	if h, port, err := net.SplitHostPort(listenaddr); err == nil && h == "::" && addrsp != nil {
		// wildcard bind
		fst := true
		for _, a := range *addrsp {
			ip := a.String()
			if ipnet, ok := a.(*net.IPNet); ok {
				ip = ipnet.IP.String()
			}
			if strings.Contains(ip, ":") { // IPv6, skip for now
				continue
			}
			f := fmt.Sprintf("http://%s:%s", ip, port)
			if len(f) < 28 {
				f += strings.Repeat(" ", 28-len(f))
			}
			if !fst {
				logger.Print("|------------------------------|\n")
			}
			fst = false
			logger.Print(fmt.Sprintf("| %s |\n", f))
		}
	} else {
		f := fmt.Sprintf("http://%s", listenaddr)
		if len(f) < 28 {
			f += strings.Repeat(" ", 28-len(f))
		}
		logger.Print(fmt.Sprintf("| %s |\n", f))
	}
	logger.Print("+------------------------------+\n")
}

// VERSION of the latest known release.
// Unused in non-production mode.
// Compared with in github.com/ostrost/ostent/ostent[+build production]
const VERSION = "0.2.0"
