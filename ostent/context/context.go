package context

import (
	"net/http"

	gorillaContext "github.com/gorilla/context"
	"golang.org/x/net/context"
)

type contextContext int

const contextKey contextContext = iota

// WithContext is to be used as r.WithContext in go1.7.
func WithContext(r *http.Request, ctx context.Context) *http.Request {
	gorillaContext.Set(r, contextKey, ctx)
	return r
}

// Context is to be used as r.Context in go1.7.
func Context(r *http.Request) context.Context {
	if p := gorillaContext.Get(r, contextKey); p != nil {
		if ctx, ok := p.(context.Context); ok {
			return ctx
		}
	}
	return context.TODO()
}

// Write is a shortcut and NOT to be used in go1.7.
func Write(r *http.Request, key, value interface{}) *http.Request {
	ctx := context.WithValue(Context(r), key, value)
	// := context.WithValue(r.Context(), ...) // go1.7

	*r = *WithContext(r, ctx) // = *r.WithContext(ctx) // go1.7
	return r
}

// Read is a shortcut and NOT to be used in go1.7.
func Read(r *http.Request, key interface{}) interface{} {
	return Context(r).Value(key) // := r.Context().Value(...) // go1.7
}
