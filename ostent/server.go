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
	CRouterParams
)

func AddAssetPathContextFunc(path string) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler { // Constructor
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			context.Set(r, CAssetPath, path)
			handler.ServeHTTP(w, r)
		})
	}
}

func ContextParams(r *http.Request) (httprouter.Params, error) {
	pinterface := context.Get(r, CRouterParams)
	if pinterface == nil {
		return httprouter.Params{}, fmt.Errorf("CRouterParams is missing in meta request")
	}
	params, ok := pinterface.(httprouter.Params)
	if !ok {
		return httprouter.Params{}, fmt.Errorf("CRouterParams not of type httprouter.Params")
	}
	return params, nil
}

// HandleFunc wraps http.HandlerFunc(h) into handle.
func HandleFunc(hf http.HandlerFunc) httprouter.Handle {
	return handle(http.HandlerFunc(hf))
}

func handle(h http.Handler) httprouter.Handle {
	// make a httprouter.Handle from h ignoring httprouter.Params.
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		h.ServeHTTP(w, r)
	}
}
func handleParamSetContext(h http.Handler) httprouter.Handle {
	// make a httprouter.Handle from h addind httprouter.Params into context.
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		context.Set(r, CRouterParams, p)
		h.ServeHTTP(w, r)
	}
}

// HandleThen wraps then(hf) into handle.
func HandleThen(then func(http.Handler) http.Handler) func(http.HandlerFunc) httprouter.Handle {
	return func(hf http.HandlerFunc) httprouter.Handle {
		return handle(then(hf))
	}
}

// ParamsFunc wraps then(hf) into handle with context setting by handleParamSetContext.
func ParamsFunc(then func(http.Handler) http.Handler) func(http.HandlerFunc) httprouter.Handle {
	return func(hf http.HandlerFunc) httprouter.Handle {
		return handleParamSetContext(then(hf))
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
