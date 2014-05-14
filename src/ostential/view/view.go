package view

import (
	"io"
	"fmt"
	"bytes"
	"strconv"
	"net/http"
	"html/template"
	"github.com/rzab/amber"
	"github.com/codegangsta/martini"
)

type stringTemplate struct {
	*template.Template
}

var UsePercentTemplate  = mustTemplate("usepercent.html")
var TooltipableTemplate = mustTemplate("tooltipable.html")

func mustTemplate(filename string) stringTemplate {
	reader, ok := _bindata[filename]
	if !ok {
		panic(fmt.Errorf("No %q in bindata\n", filename))
	}
	text, err := reader()
	if err != nil {
		panic(err)
	}
	return stringTemplate{template.Must(template.New(filename).Parse(string(text)))}
}

func(st stringTemplate) Execute(data interface{}) (template.HTML, error) {
	clone, err := st.Template.Clone()
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err := clone.Execute(buf, data); err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

func bincompile() *template.Template {
	t := template.New("templates.html")
	template.Must(t.Parse("Empty")) // initial template in case we won't have any

	for filename, reader := range _bindata { // from bindata.go
		text, err := reader()
		if err != nil {
			panic(err)
		}
		subt := t.New(filename)
		subt.Funcs(amber.FuncMap)
		template.Must(subt.Parse(string(text)))
	}
	return t
}

type Render interface {
	HTML(int, string, interface{})
}

func (render *render) HTML(status int, name string, data interface{}) {
	buf := new(bytes.Buffer)
	err := render.Template.ExecuteTemplate(buf, name, data)
	if err != nil {
		http.Error(render, err.Error(), http.StatusInternalServerError)
	}
	render.WriteHeader(status)
	render.Header().Set("Content-Type", "text/html")
 	render.Header().Set("Content-Length", strconv.Itoa(len(buf.String())))
// 	io.WriteString(render.ResponseWriter, buf.String())
	io.Copy(render.ResponseWriter, buf)
}

func BinTemplates_MartiniHandler() martini.Handler {
	empl := bincompile()
	return func(res http.ResponseWriter, c martini.Context) {
		emplone, _ := empl.Clone()
		c.MapTo(&render{
			ResponseWriter: res,
			Template: emplone,
		}, (*Render)(nil))
	}
}

type render struct {
	http.ResponseWriter
	*template.Template
}
