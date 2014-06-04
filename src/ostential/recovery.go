package ostential
import (
	"runtime"
	"net/http"
	"html/template"
)

type Recovery bool // true stands for production

func(rc Recovery) Constructor(HANDLER http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if E := recover(); E != nil {
				rc.recov(w, E)
			}
		}()
		HANDLER.ServeHTTP(w, r)
	})
}

const panicstatuscode = http.StatusInternalServerError
var   panicstatustext = statusLine(panicstatuscode)

func (RC Recovery) recov(w http.ResponseWriter, E interface{}) {
	w.WriteHeader(panicstatuscode) // NB

	var description string
	if err, ok := E.(error); ok {
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

var rctemplate = template.Must(template.New("recovery.html").Parse(`
<html>
<head><title>{{.Title}}</title></head>
<body bgcolor="white">
<center><h1>{{.Description}}</h1></center>
<hr><pre>{{.Stack}}</pre>
</body>
</html>
`))
