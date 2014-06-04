package ostential
import (
	"net/http"

	"github.com/rcrowley/go-tigertonic"
)

type ServeMux interface {
	HandleFunc(string, string, http.HandlerFunc)
	// Handle(string, http.Handler) // intentionally disabled
}

type TrieServeMux struct {
	*tigertonic.TrieServeMux
	newhandler func(func(http.ResponseWriter, *http.Request)) http.HandlerFunc
}

func NewMux(newhandler func(func(http.ResponseWriter, *http.Request)) http.HandlerFunc) *TrieServeMux {
	return &TrieServeMux{
		TrieServeMux: tigertonic.NewTrieServeMux(),
		newhandler:   newhandler,
	}
}

// catch tigertonic error handlers, override
func (mux *TrieServeMux) handlerFunc(handler http.Handler) http.HandlerFunc {
	NA := tigertonic.MethodNotAllowedHandler{}
	NF := tigertonic.NotFoundHandler{}
	if handler == NF {
		return http.NotFound
	}
	if handler == NA {
		return func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, statusLine(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
	return nil
}

func (mux *TrieServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, _ := mux.TrieServeMux.Handler(r)
	if h := mux.handlerFunc(handler); h != nil {
		handler = mux.newhandler(h)
	}
	handler.ServeHTTP(w, r)
}

func (mux *TrieServeMux) HandleFunc(method, pattern string, handlerFunc http.HandlerFunc) {
	mux.TrieServeMux.Handle(method, pattern, http.HandlerFunc(handlerFunc))
}

func (mux *TrieServeMux) Handle(string, http.Handler) {
	panic("Unexpected to be used")
}
