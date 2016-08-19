package ostent

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

type contextID int

const (
	CAssetPath contextID = iota
	CRouterParams
)

func AddAssetPathContextFunc(path string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			*r = *r.WithContext(context.WithValue(r.Context(),
				CAssetPath, path))
			h.ServeHTTP(w, r)
		})
	}
}

func contextString(r *http.Request, c contextID) (string, bool) {
	v, ok := r.Context().Value(c).(string)
	return v, ok
}

func contextParams(r *http.Request) (httprouter.Params, error) {
	params, ok := r.Context().Value(CRouterParams).(httprouter.Params)
	if !ok {
		return httprouter.Params{}, fmt.Errorf(
			"CRouterParams mistyped/missing in the context")
	}
	return params, nil
}

// ContextParam is to retrieve pname param from context params.
func ContextParam(r *http.Request, pname string) (string, error) {
	params, err := contextParams(r)
	if err != nil {
		return "", err
	}
	pvalue := params.ByName(pname)
	if pvalue == "" {
		return "", fmt.Errorf("%q param is empty", pname)
	}
	return pvalue, nil
}

// HandleFunc wraps hf into handle.
func HandleFunc(hf http.HandlerFunc) httprouter.Handle { return handle(hf) }

func handle(h http.Handler) httprouter.Handle {
	// make a httprouter.Handle from h ignoring httprouter.Params.
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		h.ServeHTTP(w, r)
	}
}
func handleParamSetContext(h http.Handler) httprouter.Handle {
	// make a httprouter.Handle from h addind httprouter.Params into context.
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		*r = *r.WithContext(context.WithValue(r.Context(),
			CRouterParams, p))
		h.ServeHTTP(w, r)
	}
}

// HandleThen wraps then(hf) into handle.
func HandleThen(then func(http.Handler) http.Handler) func(http.HandlerFunc) httprouter.Handle {
	return func(hf http.HandlerFunc) httprouter.Handle { return handle(then(hf)) }
}

// ParamsFunc wraps then(hf) into handle with context setting by handleParamSetContext.
func ParamsFunc(then func(http.Handler) http.Handler) func(http.HandlerFunc) httprouter.Handle {
	return func(hf http.HandlerFunc) httprouter.Handle {
		if then != nil {
			hf = then(hf).ServeHTTP
		}
		return handleParamSetContext(hf)
	}
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
	path, ok := contextString(r, CAssetPath)
	if !ok || path == "" {
		sa.error(w, fmt.Errorf("CAssetPath mistyped/missing in the context"))
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
