package ostent

import (
	"html/template"
	"net/http"
	"runtime"
)

func PanicHandlerFunc(taggedbin bool, recd interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(PanicStatusCode)

		var description string
		if err, ok := recd.(error); ok {
			description = err.Error()
		} else if s, ok := recd.(string); ok {
			description = s
		}
		var stack string
		if !taggedbin {
			sbuf := make([]byte, 4096-len(PanicStatusText)-len(description))
			size := runtime.Stack(sbuf, false)
			stack = string(sbuf[:size])
		}
		if tpl, err := PanicTemplate.Clone(); err == nil { // otherwise bail out
			tpl.Execute(w, struct {
				Title, Description, Stack string
			}{
				Title:       PanicStatusText,
				Description: description,
				Stack:       stack,
			})
		}
	}
}

const PanicStatusCode = http.StatusInternalServerError

var (
	PanicStatusText = statusLine(PanicStatusCode)
	PanicTemplate   = template.Must(template.New("recovery.html").Parse(`
<html>
<head><title>{{.Title}}</title></head>
<body bgcolor="white">
<center><h1>{{.Description}}</h1></center>
<hr><pre>{{.Stack}}</pre>
</body>
</html>
`))
)
