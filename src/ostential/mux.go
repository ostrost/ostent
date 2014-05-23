package ostential
import (
	"log"
	"net/http"

	"github.com/rcrowley/go-tigertonic"
)

type ServeMux interface {
	HandleFunc(string, string, http.HandlerFunc)
	// Handle(string, http.Handler) // intentionally disabled
}

type TrieServeMux struct {
	*tigertonic.TrieServeMux
	production    bool
	logger        *log.Logger
	access_logger *log.Logger
}

func NewMux(production bool, logger, access_logger *log.Logger) *TrieServeMux {
	return &TrieServeMux{
		TrieServeMux:  tigertonic.NewTrieServeMux(),
		production:    production,
		logger:        logger,
		access_logger: access_logger,
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

func (mux *TrieServeMux) wrap(handlerFunc http.HandlerFunc) http.Handler {
	return LoggedFunc(mux.production, mux.access_logger)(
		RecoveringFunc(mux.production, mux.logger)(handlerFunc))
}

func (mux *TrieServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, _ := mux.TrieServeMux.Handler(r)
	if h := mux.handlerFunc(handler); h != nil {
		handler = mux.wrap(h)
	}
	handler.ServeHTTP(w, r)
}

func (mux *TrieServeMux) HandleFunc(method, pattern string, handlerFunc http.HandlerFunc) {
	handler := mux.wrap(handlerFunc)
	if method != "" {
		mux.TrieServeMux.Handle(method, pattern, handler)
		return
	}
	for _, method := range []string{"HEAD", "GET", "POST"} {
		mux.TrieServeMux.Handle(method, pattern, handler)
	}
}

func (mux *TrieServeMux) Handle(string, http.Handler) {
	panic("Unexpected to be used")
}
