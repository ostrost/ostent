package templatefunc

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/ostrost/ostent/templateutil/templatepipe"
)

func ExecuteWithHashTest(t *testing.T, tm *template.Template, expected string) {
	d := templatepipe.Dotted{}
	d.Append([]string{"a"}, nil) // non-string result from templatepipe.Encurl
	h := templatepipe.Encurl(templatepipe.CurlyX, d, 0)
	if _, ok := h.(string); ok {
		t.Errorf("Encurl expected to return non-string on %+v input", d)
	}
	buf := new(bytes.Buffer)
	if err := tm.Execute(buf, h); err != nil {
		t.Fatal(err)
	}
	if s := buf.String(); s != expected {
		t.Errorf("Execute with Encurl: %q (expected %q)", s, expected)
	}
}

func TestKFuncs(t *testing.T) {
	define := `{{define "defines::define_table"}}{{with columnsset .}}{{.}}{{else}}no columns{{end}}{{end}}`
	text1 := define + `{{template "defines::define_table" .}}`
	text2 := define + `{{template "defines::define_table" setcolumns . "columns"}}`
	fs := template.FuncMap{
		"columnsset": GetKFunc(".OverrideColumns"),
		"setcolumns": SetKFunc(".OverrideColumns"),
	}
	tm1, err1 := template.New("withdot1").Funcs(fs).Parse(text1)
	tm2, err2 := template.New("withdot2").Funcs(fs).Parse(text2)
	if err1 != nil {
		t.Fatal(err1)
	}
	if err2 != nil {
		t.Fatal(err2)
	}
	ExecuteWithHashTest(t, tm1, "no columns")
	ExecuteWithHashTest(t, tm2, "{columns}")
}

func SetKText(t *testing.T, in, expected string) {
	override := ".Override"
	h := SetKFunc(override)(templatepipe.Hash{}, in)
	d, ok := h.(templatepipe.Hash)[override]
	if !ok {
		t.Errorf("SetK[%q] is not okd.", override)
	}
	dv, ok := d.(string)
	if !ok {
		t.Error("Cannot cast to string")
	}
	if dv != expected {
		t.Errorf("SetK[%q] mismatch: %q (expected %q)", override, dv, expected)
	}
}

func TestSetK(t *testing.T) {
	SetKText(t, "KEY", "{KEY}")
	SetKText(t, "aHTML", "<span dangerouslySetInnerHTML={{__html: aHTML}} />")
}
