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

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

type ContextID int

const (
	CPanicError ContextID = iota
	CAssetPath
)

func AddAssetPathContextFunc(path string) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler { // Constructor
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			context.Set(r, CAssetPath, path)
			handler.ServeHTTP(w, r)
		})
	}
}

// Muxmap is a type of a map of pattern to HandlerFunc.
type Muxmap map[string]http.HandlerFunc

func NewServery(taggedbin bool, extramap Muxmap) (*httprouter.Router, alice.Chain, *Access) {
	access := NewAccess(taggedbin)
	chain := alice.New(access.Constructor)
	mux := httprouter.New()
	mux.NotFound = chain.ThenFunc(http.NotFound)
	phandler := chain.Append(context.ClearHandler).
		ThenFunc(NewServePanic(taggedbin).PanicHandler)
	mux.PanicHandler = func(w http.ResponseWriter, r *http.Request, recd interface{}) {
		context.Set(r, CPanicError, recd)
		phandler.ServeHTTP(w, r)
	}
	for path, handler := range extramap {
		h := chain.Then(handler)
		GEAD(mux, path, h)
		mux.Handler("POST", path, h)
	}
	return mux, chain, access
}

func GEAD(mux *httprouter.Router, path string, handler http.Handler) {
	mux.Handler("GET" /* */, path, handler)
	mux.Handler("HEAD" /**/, path, handler)
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

type ServeAssets struct {
	Log                 *log.Logger
	AssetFunc           func(string) ([]byte, error)
	AssetInfoFunc       func(string) (os.FileInfo, error)
	AssetAltModTimeFunc func() time.Time // may be nil
}

// Serve does http.ServeContent with asset content and info.
func (sa ServeAssets) Serve(w http.ResponseWriter, r *http.Request) {
	p := context.Get(r, CAssetPath)
	if p == nil {
		err := fmt.Errorf("ServeAssets.Serve must receive CAssetPath in context")
		sa.Log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	path, ok := p.(string)
	if !ok {
		err := fmt.Errorf("ServeAssets.Serve received non-string CAssetPath in context")
		sa.Log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	text, err := sa.AssetFunc(path)
	if err != nil {
		sa.Log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var mt time.Time
	if sa.AssetAltModTimeFunc != nil {
		mt = sa.AssetAltModTimeFunc()
	} else {
		info, err := sa.AssetInfoFunc(path)
		if err != nil {
			sa.Log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		mt = info.ModTime()
	}
	http.ServeContent(w, r, path, mt, bytes.NewReader(text))
}

type loggerPrint interface {
	Print(v ...interface{})
}

// Banner prints a banner with the logger.
func Banner(listenaddr string, suffix string, logger loggerPrint) {
	hostname, _ := GetHN()
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
const VERSION = "0.4.1"
