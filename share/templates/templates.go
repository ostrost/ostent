package templates

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"sync"

	"github.com/rzab/amber"
)

var (
	UsePercentTemplate  = binTemplate{filename: "usepercent.html"}
	TooltipableTemplate = binTemplate{filename: "tooltipable.html"}
	IndexTemplate       = binTemplate{filename: "index.html"}
)

func InitTemplates() {
	UsePercentTemplate.init()
	TooltipableTemplate.init()
	IndexTemplate.init()
}

type binTemplate struct {
	template *template.Template
	filename string
	mutex    sync.Mutex
}

func (bt *binTemplate) init() {
	bt.mutex.Lock()
	defer bt.mutex.Unlock()
	bt.initUnlocked()
}

func (bt *binTemplate) initUnlocked() { // panics (explicit and template.Must) on any error
	text, err := Asset(bt.filename)
	if err != nil {
		panic(err)
	}
	if bt.filename != "index.html" { // the simple case
		bt.template = template.Must(template.New(bt.filename).Parse(string(text)))
		return
	}
	// index.html specifics:
	// 1. `t' may be .New'd multiple times for cascaded templates
	// 2. custom .Funcs

	t := template.New("templates.html") // root template, MUST NOT t.New("templates.html") later, causes redefinition of the template
	template.Must(t.Parse("Empty"))     // initial template in sudden case we won't have any

	// repeat if necessary, `name' for .New must be new
	subt := t.New(bt.filename)
	subt.Funcs(amber.FuncMap)
	template.Must(subt.Parse(string(text)))

	bt.template = t
}

func (bt *binTemplate) Execute(data interface{}) (*bytes.Buffer, error) {
	var (
		filename string
		clone    *template.Template
		err      error
	)
	func() {
		bt.mutex.Lock()
		defer bt.mutex.Unlock()
		if bt.template == nil {
			bt.initUnlocked()
		}
		clone, err = bt.template.Clone()
		filename = bt.filename
	}()
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	if err := clone.ExecuteTemplate(buf, bt.filename, data); err != nil {
		return nil, err
	}
	return buf, nil
}

type httpResponse struct {
	writer http.ResponseWriter
	buf    *bytes.Buffer
	err    error
}

func (bt *binTemplate) Response(w http.ResponseWriter, data interface{}) httpResponse {
	r := httpResponse{writer: w}
	r.buf, r.err = bt.Execute(data)
	return r
}

func (r *httpResponse) SetHeader(name, value string) {
	if r.err != nil {
		return
	}
	r.writer.Header().Set(name, value)
}

func (r *httpResponse) SetContentLength() {
	r.SetHeader("Content-Length", strconv.Itoa(r.buf.Len()))
}

func (r *httpResponse) Send() {
	if r.err != nil {
		http.Error(r.writer, r.err.Error(), http.StatusInternalServerError)
	} else {
		io.Copy(r.writer, r.buf) // or w.Write(buf.Bytes())
	}
}
