package ostent

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/justinas/alice"
	"github.com/ostrost/ostent/src/share/assets"
)

type bindValue struct {
	string
	defport string // const
	Host    string // available after flag.Parse()
	Port    string // available after flag.Parse()
}

func newBind(defstring, defport string) bindValue {
	bv := bindValue{defport: defport}
	bv.Set(defstring)
	return bv
}

// satisfying flag.Value interface
func (bv bindValue) String() string { return string(bv.string) }
func (bv *bindValue) Set(input string) error {
	if input == "" {
		return nil
	}
	if !strings.Contains(input, ":") {
		input = ":" + input
	}
	var err error
	bv.Host, bv.Port, err = net.SplitHostPort(input)
	if err != nil {
		return err
	}
	if bv.Host == "*" {
		bv.Host = ""
	} else if bv.Port == "127" {
		bv.Host = "127.0.0.1"
		bv.Port = bv.defport
	}
	if _, err = net.LookupPort("tcp", bv.Port); err != nil {
		if bv.Host != "" {
			return err
		}
		bv.Host, bv.Port = bv.Port, bv.defport
	}

	bv.string = bv.Host + ":" + bv.Port
	return nil
}

// OstentBindFlag is a bindValue hoding the ostent bind address.
var OstentBindFlag = newBind(":8050", "8050")

// CollectdBindFlag is a bindValue hoding the ostent collectd bind address.
// var CollectdBindFlag = newBind("", "8051") // "" by default meaning DO NOT BIND
func init() {
	flag.Var(&OstentBindFlag, "b", "short for bind")
	flag.Var(&OstentBindFlag, "bind", "Bind address")
	// flag.Var(&CollectdBindFlag, "collectdb",    "short for collectdbind")
	// flag.Var(&CollectdBindFlag, "collectdbind", "Bind address for collectd receiving")
}

// Muxmap is a type of a map of pattern to HandlerFunc.
type Muxmap map[string]http.HandlerFunc

var stdaccess *logger // a global, available after Serve call

// Serve does http.Serve with the listener l and constructed *TrieServeMux.
// production is passed to logging and recoverying middleware.
// Non-nil extramap is passed to the mux.
// Returns http.Serve result.
func Serve(l net.Listener, production bool, extramap Muxmap) error {
	logger := log.New(os.Stderr, "[ostent] ", 0)
	access := log.New(os.Stdout, "", 0)

	stdaccess = newLogged(production, access)
	recovery := recovery(production)

	chain := alice.New(
		stdaccess.Constructor,
		recovery.Constructor,
	)
	mux := NewMux(chain.Then)

	for _, path := range assets.AssetNames() {
		hf := chain.Then(serveContentFunc(path, logger))
		mux.Handle("GET", "/"+path, hf)
		mux.Handle("HEAD", "/"+path, hf)
	}

	//	chain.ThenFunc(slashws)) handler would include stdlogger
	//  slashws uses stdlogger itself
	mux.Handle("GET", "/ws", recovery.ConstructorFunc(slashws))

	mux.Handle("GET", "/", chain.ThenFunc(index))
	mux.Handle("HEAD", "/", chain.ThenFunc(index))

	/* panics := func(http.ResponseWriter, *http.Request) {
		panic(fmt.Errorf("I'm panicing"))
	}
	mux.Handle("GET", "/panic", chain.ThenFunc(panics)) // */

	if extramap != nil {
		for path, handler := range extramap {
			for _, METH := range []string{"HEAD", "GET", "POST"} {
				mux.Handle(METH, path, chain.Then(handler))
			}
		}
	}

	banner(l, logger)

	server := &http.Server{Addr: l.Addr().String(), Handler: mux}
	return server.Serve(l)
}

func serveContentFunc(path string, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		text, err := assets.Uncompressedasset(path)
		if err != nil {
			panic(err)
		}
		modtime, err := assets.ModTime("src/share/assets", path)
		if err != nil {
			logger.Println(err)
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			// return
		}

		http.ServeContent(w, r, path, modtime, bytes.NewReader(text))
	}
}

func banner(listen net.Listener, logger *log.Logger) {
	hostname, _ := getHostname()
	logger.Printf("   %s\n", strings.Repeat("-", len(hostname)+7))
	if len(hostname) > 19 {
		hostname = hostname[:16] + "..."
	}
	logger.Printf(" / %s ostent \\ \n", hostname)
	logger.Printf("+------------------------------+")

	addr := listen.Addr()
	if h, port, err := net.SplitHostPort(addr.String()); err == nil && h == "::" {
		// wildcard bind

		/* _, IP := NewInterfaces()
		logger.Printf("        http://%s", IP) // */
		addrs, err := net.InterfaceAddrs()
		if err == nil {
			fst := true
			for _, a := range addrs {
				ipnet, ok := a.(*net.IPNet)
				if !ok || strings.Contains(ipnet.IP.String(), ":") {
					continue // no IPv6 for now
				}
				f := fmt.Sprintf("http://%s:%s", ipnet.IP.String(), port)
				if len(f) < 28 {
					f += strings.Repeat(" ", 28-len(f))
				}
				if !fst {
					logger.Printf("|------------------------------|")
				}
				fst = false
				logger.Printf("| %s |", f)
			}
		}
	} else {
		f := fmt.Sprintf("http://%s", addr.String())
		if len(f) < 28 {
			f += strings.Repeat(" ", 28-len(f))
		}
		logger.Printf("| %s |", f)
	}
	logger.Printf("+------------------------------+")
}

// VERSION of the latest known release.
// Unused in non-production mode.
// Compared with in main.production.go.
const VERSION = "0.1.9"
