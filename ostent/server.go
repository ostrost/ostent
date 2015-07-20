package ostent

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// Muxmap is a type of a map of pattern to HandlerFunc.
type Muxmap map[string]http.HandlerFunc

func NewServery(taggedbin bool, extramap Muxmap) (*httprouter.Router, alice.Chain, *Access) {
	access := NewAccess(taggedbin)
	chain := alice.New(access.Constructor)
	mux := httprouter.New()
	mux.PanicHandler = access.PanicHandler
	for path, handler := range extramap {
		h := chain.Then(handler)
		mux.Handler("GET", path, h)
		mux.Handler("HEAD", path, h)
		mux.Handler("POST", path, h)
	}
	return mux, chain, access
}

// TimeInfo is for AssetInfoFunc: a reduced os.FileInfo.
type TimeInfo interface {
	ModTime() time.Time
}

// AssetInfoFunc wraps bindata's AssetInfo func. Returns typecasted infofunc.
func AssetInfoFunc(infofunc func(string) (os.FileInfo, error)) func(string) (TimeInfo, error) {
	return func(name string) (TimeInfo, error) {
		return infofunc(name)
	}
}

// AssetReadFunc wraps bindata's Asset func. Returns readfunc itself.
func AssetReadFunc(readfunc func(string) ([]byte, error)) func(string) ([]byte, error) {
	return readfunc
}

// ServeContentFunc does http.ServeContent the readFunc (Asset or UncompressedAsset) result.
// infofunc is typically AssetInfo. modtimefunc may override info.Modtime() result.
func ServeContentFunc(
	readfunc func(string) ([]byte, error),
	infofunc func(string) (TimeInfo, error),
	path string, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		text, err := readfunc(path)
		if err != nil {
			logger.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		info, err := infofunc(path)
		if err != nil {
			logger.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.ServeContent(w, r, path, info.ModTime(), bytes.NewReader(text))
	}
}

type loggerPrint interface {
	Print(v ...interface{})
}

// Banner prints a banner with the logger.
func Banner(listenaddr string, suffix string, logger loggerPrint) {
	hostname, _ := (&Machine{}).GetHostname()
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
// Unused in non-bin mode.
// Compared with in github.com/ostrost/ostent/ostent[+build bin]
const VERSION = "0.2.0"
