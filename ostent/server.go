package ostent

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/julienschmidt/httprouter"
)

func LogHandler(logRequests bool, h http.Handler) http.Handler {
	if !logRequests {
		return h
	}
	return handlers.CombinedLoggingHandler(os.Stderr, h)
}

func ServerHandler(logRequests bool, h http.Handler) http.Handler {
	h = handlers.RecoveryHandler(
		handlers.RecoveryLogger(log.New(os.Stderr, "[panic recovery] ", log.LstdFlags)),
		handlers.PrintRecoveryStack(true),
	)(h)
	return LogHandler(logRequests, h)
}

type ContextID int

const (
	CAssetPath ContextID = iota
	CRouterParams
)

func AddAssetPathContextFunc(path string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			context.Set(r, CAssetPath, path)
			h.ServeHTTP(w, r)
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
	ReadFunc       func(string) ([]byte, error)
	InfoFunc       func(string) (os.FileInfo, error)
	AltModTimeFunc func() time.Time // may be nil
}

func (sa ServeAssets) error(w http.ResponseWriter, err error) {
	logru.Println(err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// Serve does http.ServeContent with asset content and info.
func (sa ServeAssets) Serve(w http.ResponseWriter, r *http.Request) {
	p := context.Get(r, CAssetPath)
	if p == nil {
		sa.error(w, fmt.Errorf("ServeAssets.Serve must receive CAssetPath in context"))
		return
	}
	path, ok := p.(string)
	if !ok {
		sa.error(w, fmt.Errorf("ServeAssets.Serve received non-string CAssetPath in context"))
		return
	}
	text, err := sa.ReadFunc(path)
	if err != nil {
		sa.error(w, err)
		return
	}
	var mt time.Time
	if sa.AltModTimeFunc != nil {
		mt = sa.AltModTimeFunc()
	} else {
		info, err := sa.InfoFunc(path)
		if err != nil {
			sa.error(w, err)
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
