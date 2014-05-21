package view
import (
	"ostential/assets"

	"bytes"
	"strings"
	"net/http"
)

func AssetsHandlerFunc(prefix string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" && req.Method != "HEAD" {
			http.Error(w, "405 Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		path := req.URL.Path
		if len(path) >= 3 && path[len(path) - 3:] != ".go" && // cover the bindata.go
			strings.HasPrefix(path, prefix) {

			path = path[len(prefix):]
			text, err := assets.Asset(path)
			if err == nil {
				reader := bytes.NewReader(text)
				http.ServeContent(w, req, path, assets.ModTime(), reader)
				return
			}
		}
		http.Error(w, "404 Not Found", http.StatusNotFound)
	}
}
