package ostential
import (
	"runtime"
	"net/http"
	"html/template"
)

type Recovery bool // true stands for production

func(RC Recovery) Constructor(HANDLER http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(panicstatuscode) // NB

				var description string
				if err, ok := err.(error); ok {
					description = err.Error()
				}
				var stack string
				if !RC { // if !production
					sbuf := make([]byte, 4096 - len(panicstatustext) - len(description))
					size := runtime.Stack(sbuf, false)
					stack = string(sbuf[:size])
				}
				rctemplate.Execute(w, struct {
					Title, Description, Stack string
				}{
					Title:       panicstatustext,
					Description: description,
					Stack:       stack,
				})
			}
		}()
		HANDLER.ServeHTTP(w, r)
	})
}

const panicstatuscode = http.StatusInternalServerError
var   panicstatustext = statusLine(panicstatuscode)

var rctemplate = template.Must(template.New("recovery.html").Parse(`
<html>
<head><title>{{.Title}}</title></head>
<body bgcolor="white">
<center><h1>{{.Description}}</h1></center>
<hr><pre>{{.Stack}}</pre>
</body>
</html>
`))
