package templatep

import (
	"testing"
	"text/template"

	"github.com/ostrost/ostent/templateutil/templatepipe"
)

func ExecutewDottedTest(t *testing.T, tm *template.Template, expected string) {
	d := templatepipe.Dotted{}
	d.Append([]string{"a"}, nil) // non-string result from templatepipe.Mkmap
	h := templatepipe.Mkmap(d, 0)
	if _, ok := h.(string); ok {
		t.Errorf("Mkmap expected to return non-string on %+v input", d)
	}
	s, err := StringExecute(tm, h)
	if err != nil {
		t.Fatal(err)
	}
	if s != expected {
		t.Errorf("Execute with Mkmap: %q (expected %q)", s, expected)
	}
}

func TestDot(t *testing.T) {
	define := `[[define "define_withdot"]][[with .OVERRIDE]][[.]][[else]][[end]][[end]]`
	text1 := define + `[[template "define_withdot" .]]`
	text2 := define + `[[template "define_withdot" (dot . "rows")]]`
	df := template.FuncMap{"dot": dot}
	tm1, err1 := template.New("withdot1").Funcs(df).Delims("[[", "]]").Parse(text1)
	tm2, err2 := template.New("withdot2").Funcs(df).Delims("[[", "]]").Parse(text2)
	if err1 != nil {
		t.Fatal(err1)
	}
	if err2 != nil {
		t.Fatal(err2)
	}
	ExecutewDottedTest(t, tm1, "")
	ExecutewDottedTest(t, tm2, "{rows}")
}

func DotValueText(t *testing.T, in, expected string) {
	h := dot(templatepipe.Hash{}, in)
	d, ok := h["OVERRIDE"]
	if !ok {
		t.Errorf("Getting \"dot\" from `dot' result is not okd.")
	}
	dv, ok := d.(dotValue)
	if !ok {
		t.Error("Cannot cast to dotValue")
	}
	if s := dv.String(); s != expected {
		t.Errorf("dotValue mismatch: %q (expected %q)", s, expected)
	}
	if v, ok := (*dv.hashp)["OVERRIDE"]; ok || v != nil {
		t.Errorf("Getting \"dot\" from `dot' result is okd: %+v.", v)
	}
}

func TestDotValue(t *testing.T) {
	DotValueText(t, "KEY", "{KEY}")
	DotValueText(t, "aHTML", "<span dangerouslySetInnerHTML={{__html: aHTML}} />")
}
