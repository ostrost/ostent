package ostential
import (
	"ostential/view"

	"os"
	"fmt"
	"log"
	"net"
	"flag"
	"time"
	"sync"
	"strings"
	"net/http"

	"github.com/codegangsta/martini"
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

func newModern() *Modern { // customized martini.Classic
	r := martini.NewRouter()
	m := martini.New()
	m.Use(martini.Recovery())
	m.Action(r.Handle)

	return &Modern{
		Martini: m,
		Router: r,
	}
}

func Serve(listen net.Listener, logfunc Logfunc, cb func(*Modern)) error {
	m := newModern() // as oppose to classic
	if cb != nil {
		cb(m)
	}

	logger := log.New(os.Stderr, "[ostent] ", 0)
	m.Map(logger) // log.Logger object

	m.Use(logfunc) // log middleware

	// m.Use(assets_bindata())
	m.Any("/robots.txt",        view.AssetsHandlerFunc("/"))        // http.HandleFunc("/robots.txt", ...)
	m.Any("/assets/robots.txt", http.NotFound)                      // http.HandleFunc("/assets/robots.txt", http.NotFound)
	m.Any("/assets/.*",         view.AssetsHandlerFunc("/assets/")) // http.HandleFunc("/assets/",    ...)

	// a martini.Handler, handles all non-asset requests e.g. "/" and "/ws"
	m.Use(view.BinTemplates_MartiniHandler())

	m.Get("/",   index)
	m.Get("/ws", slashws)

	hostname := getGeneric().HostnameString
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

	server := &http.Server{Addr: listen.Addr().String(), Handler: m}
	return server.Serve(listen)
}

/* func assets_bindata() martini.Handler {
	return func(res http.ResponseWriter, req *http.Request, log *log.Logger) {
		if req.Method != "GET" && req.Method != "HEAD" {
			return
		}
		path := req.URL.Path
		if path == "/" || path == "" || filepath.Ext(path) == ".go" { // cover the bindata.go
			return
		}
		if path[0] == '/' {
			path = path[1:]
		}
		text, err := assets.Asset(path)
		if err != nil {
			return
		}
		reader := bytes.NewReader(text)
		http.ServeContent(res, req, path, assets.ModTime(), reader)
	}
} // */

type Logfunc martini.Handler

var logOneLock sync.Mutex
var logged = map[string]bool{}
func LogOne(res http.ResponseWriter, req *http.Request, c martini.Context, logger *log.Logger) {
	start := time.Now()
	c.Next()

	rw := res.(martini.ResponseWriter)
	status := rw.Status()
	if status != 200 && status != 304 && req.URL.Path != "/ws" {
		logThis(start, res, req, logger)
		return
	}

	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		host = req.RemoteAddr
	}
	logOneLock.Lock()
	if _, ok := logged[host]; ok {
		logOneLock.Unlock()
		return
	}
	logged[host] = true
	logOneLock.Unlock()

	logger.Printf("%s\tRequested from %s; subsequent successful requests will not be logged\n", time.Now().Format("15:04:05"), host)
}

func LogAll(res http.ResponseWriter, req *http.Request, c martini.Context, logger *log.Logger) {
	start := time.Now()
	c.Next()
	logThis(start, res, req, logger)
}

var ZEROTIME, _ = time.Parse("15:04:05", "00:00:00")

func logThis(start time.Time, res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	diff := time.Since(start)
	since := ZEROTIME.Add(diff).Format("5.0000s")

	rw := res.(martini.ResponseWriter)
	status := rw.Status()
	code := fmt.Sprintf("%d", status)
	if status != 200 {
		text := http.StatusText(status)
		if text != "" {
			code += fmt.Sprintf(" %s", text)
		}
	}
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		host = req.RemoteAddr
	}

	logger.Printf("%s\t%s\t%s\t%v\t%s\t%s\n", start.Format("15:04:05"), host, since, code, req.Method, req.URL.Path)
}

const VERSION = "0.1.6"
