package ostential
import (
	"ostential/assets"

	"os"
	"fmt"
	"log"
	"net"
	"flag"
	"bytes"
	"strings"
	"net/http"
)

type bindValue struct {
	string
	defport string // const
	Host, Port string // available after flag.Parse()
}

func newBind(defstring, defport string) bindValue {
	bv := bindValue{defport: defport}
	bv.Set(defstring)
	return bv
}

// satisfying flag.Value interface
func(bv bindValue) String() string { return string(bv.string); }
func(bv *bindValue) Set(input string) error {
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

var BindFlag = newBind(":8050", "8050")
func init() {
	flag.Var(&BindFlag, "b",    "Bind address")
	flag.Var(&BindFlag, "bind", "Bind address")
}

func Serve(listen net.Listener, production bool, cbservemux func(ServeMux)) error {
	logger := log.New(os.Stderr, "[ostent] ", 0)
	mux := NewMux(production, logger, log.New(os.Stdout, "", 0))

	for _, filename := range assets.BindataKeys() {
		hf := serveContentFunc(filename)
		mux.HandleFunc("GET",  "/"+ filename, hf)
		mux.HandleFunc("HEAD", "/"+ filename, hf)
	}

	mux.HandleFunc("GET",  "/ws", slashws)
	mux.HandleFunc("GET",  "/",   index)
	mux.HandleFunc("HEAD", "/",   index)

	// mux.HandleFunc("/panic", ... func(http.ResponseWriter, *http.Request) { panic(fmt.Errorf("I'm panicing")) })

	if cbservemux != nil {
		cbservemux(mux)
	}

	banner(listen, logger)

	server := &http.Server{Addr: listen.Addr().String(), Handler: mux}
	return server.Serve(listen)
}

func serveContentFunc(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		text, err := assets.Asset(path)
		if err != nil {
			panic(err)
		}
		reader := bytes.NewReader(text)
		modtime, err := assets.ModTime("assets", path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.ServeContent(w, r, path, modtime, reader)
	}
}

func banner(listen net.Listener, logger *log.Logger) {
	hostname := getGeneric().Hostname
	logger.Printf("   %s\n", strings.Repeat("-", len(hostname) + 7))
	if len(hostname) > 19 {
		hostname = hostname[:16] +"..."
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
					f += strings.Repeat(" ", 28 - len(f))
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
			f += strings.Repeat(" ", 28 - len(f))
		}
		logger.Printf("| %s |", f)
	}
	logger.Printf("+------------------------------+")
}

const VERSION = "0.1.6"
