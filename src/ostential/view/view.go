package view
import (
	"fmt"
	"bytes"
	"html/template"
	"github.com/rzab/amber"
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

func Bincompile() *template.Template {
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
