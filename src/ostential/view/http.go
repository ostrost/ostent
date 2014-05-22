package view
import (
	"ostential/assets"

	"fmt"
	"bytes"
	"strings"
	"net/http"
)

func StatusLine(status int) string {
	return fmt.Sprintf("%d %s", status, http.StatusText(status))
}

func AssetsHandlerFunc(prefix string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" && req.Method != "HEAD" {
			http.Error(w, StatusLine(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		path := req.URL.Path
		if len(path) >= 3 && path[len(path) - 3:] != ".go" && // cover the bindata.go
			strings.HasPrefix(path, prefix) {
			path = path[len(prefix):]

			if text, err := assets.Asset(path); err == nil {
				reader := bytes.NewReader(text)
				http.ServeContent(w, req, path, assets.ModTime(), reader)
				return
			}
		}
		http.Error(w, StatusLine(http.StatusNotFound), http.StatusNotFound)
	}
}
