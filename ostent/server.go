package ostent

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
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

type Handle struct {
	Chain alice.Chain
	HFunc http.HandlerFunc
}

func NewHandle(chain alice.Chain, hfunc http.HandlerFunc) Handle {
	return Handle{Chain: chain, HFunc: hfunc}
}

type RouteInfo struct {
	Path   string
	Asset  bool
	Post   bool // filled after ApplyRoutes
	Params bool // filled after ApplyRoutes
}

type Route struct {
	*RouteInfo
	RouterFuncs *[]func(string, httprouter.Handle)
}

func NewRoute(path string, rfs ...func(string, httprouter.Handle)) *Route {
	return &Route{&RouteInfo{Path: path}, &rfs}
}

func ApplyRoutes(r *httprouter.Router, routes map[*Route]Handle, conv func(http.Handler) httprouter.Handle) {
	if conv == nil {
		// blank converter: does nothing with httprouter.Params
		conv = func(handler http.Handler) httprouter.Handle {
			return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
				handler.ServeHTTP(w, r)
			}
		}
	}
	// .RouterFuncs bind the path with the router
	for route, handle := range routes {
		for _, rfunc := range *route.RouterFuncs {
			rfunc(route.Path, conv(handle.Chain.ThenFunc(handle.HFunc)))
		}
	}
	// fill .RouteInfo.{Post,Params}
	for route := range routes {
		// tsr stands for trailing slash redirect
		if handle, params, tsr := r.Lookup("POST", route.Path); handle != nil && !tsr {
			route.Post = true
			route.Params = params != nil
		}
		// there should not be a params-less POST and params-featured GET
		if !route.Params {
			if handle, params, tsr := r.Lookup("GET", route.Path); handle != nil && !tsr {
				route.Params = params != nil
			}
		}
	}
}

func NewServery(taggedbin bool) (*httprouter.Router, alice.Chain, *Access) {
	access := NewAccess(taggedbin)
	achain := alice.New(access.Constructor)
	r := httprouter.New()
	r.NotFound = achain.ThenFunc(http.NotFound)
	phandler := achain.Append(context.ClearHandler).
		ThenFunc(NewServePanic(taggedbin).PanicHandler)
	r.PanicHandler = func(w http.ResponseWriter, r *http.Request, recd interface{}) {
		context.Set(r, CPanicError, recd)
		phandler.ServeHTTP(w, r)
	}
	return r, achain, access
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

// VERSION of the latest known release.
// Unused in non-bin mode.
// Compared with in github.com/ostrost/ostent/main.go
// MUST BE semver compatible: no two digits ("X.Y") allowed.
const VERSION = "0.6.1"
