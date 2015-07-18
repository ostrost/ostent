package ostent

import (
	"html/template"
	"net/http"
	"runtime"
)

type Recovery bool // true means stack inclusion in error output

func (rc Recovery) ConstructorFunc(hf http.HandlerFunc) http.Handler {
	return rc.Constructor(http.HandlerFunc(hf))
}

func (rc Recovery) PanicHandle(recd interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(panicstatuscode) // NB

		var description string
		if err, ok := recd.(error); ok {
			description = err.Error()
		} else if s, ok := recd.(string); ok {
			description = s
		}
		var stack string
		if !rc {
			sbuf := make([]byte, 4096-len(panicstatustext)-len(description))
			size := runtime.Stack(sbuf, false)
			stack = string(sbuf[:size])
		}
		if tpl, err := rctemplate.Clone(); err == nil { // otherwise bail out
			tpl.Execute(w, struct {
				Title, Description, Stack string
			}{
				Title:       panicstatustext,
				Description: description,
				Stack:       stack,
			})
		}
	}
}

func (rc Recovery) Constructor(HANDLER http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recd := recover(); recd != nil {
				rc.PanicHandle(recd)(w, r)
			}
		}()
		HANDLER.ServeHTTP(w, r)
	})
}

const panicstatuscode = http.StatusInternalServerError

var panicstatustext = statusLine(panicstatuscode)

var rctemplate = template.Must(template.New("recovery.html").Parse(`
<html>
<head><title>{{.Title}}</title></head>
<body bgcolor="white">
<center><h1>{{.Description}}</h1></center>
<hr><pre>{{.Stack}}</pre>
</body>
</html>
`))
