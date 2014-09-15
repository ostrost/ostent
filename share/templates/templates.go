package templates

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"sync"

	// "github.com/ostrost/ostent/templates"
	"github.com/rzab/amber"
)

var (
	UsePercentTemplate  = BinTemplate{filename: "usepercent.html"}
	TooltipableTemplate = BinTemplate{filename: "tooltipable.html"}
	IndexTemplate       = BinTemplate{filename: "index.html"}
)

func InitTemplates() {
	UsePercentTemplate.Init()
	TooltipableTemplate.Init()
	IndexTemplate.Init()
}

type BinTemplate struct {
	template *template.Template
	filename string
	mutex    sync.Mutex
}

func (bt *BinTemplate) Init() {
	bt.mutex.Lock()
	defer bt.mutex.Unlock()
	bt.initUnlocked()
}

func (bt *BinTemplate) initUnlocked() { // panics (explicit and template.Must) on any error
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

func (bt *BinTemplate) Execute(data interface{}) (*bytes.Buffer, error) {
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

type templateWriter struct {
	writer http.ResponseWriter
	buf    *bytes.Buffer
	err    error
}

func (bt *BinTemplate) Response(w http.ResponseWriter, data interface{}) templateWriter {
	tw := templateWriter{writer: w}
	tw.buf, tw.err = bt.Execute(data)
	return tw
}

func (tw *templateWriter) SetHeader(name, value string) {
	if tw.err != nil {
		return
	}
	tw.writer.Header().Set(name, value)
}

func (tw *templateWriter) SetContentLength() {
	tw.SetHeader("Content-Length", strconv.Itoa(tw.buf.Len()))
}

func (tw *templateWriter) Send() {
	if tw.err != nil {
		http.Error(tw.writer, tw.err.Error(), http.StatusInternalServerError)
	} else {
		io.Copy(tw.writer, tw.buf) // or w.Write(buf.Bytes())
	}
}
