package client

import (
	"html/template"
	"net/http"
	"net/url"
	"testing"

	"github.com/ostrost/ostent/client/enums"
)

func TestBoolLinks(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost/index?showconfigmem", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.ParseForm()
	params := NewParams()
	scm := params.BOOL["showconfigmem"]
	scm.Decode(req.Form)
	if scm.Value != true {
		t.Errorf("Decode failed: %t, expected %t", scm.Value, true)
	}
	if s := ValuesEncode(params.Query.Values); s != "showconfigmem" {
		t.Fatalf("Unexpected Values.Encode: %q", s)
	}
	if h := scm.EncodeToggle(); h != template.HTMLAttr("?") {
		t.Fatalf("Unexpected EncodeToggle: %q", h)
	}
	if s := ValuesEncode(params.Query.Values); s != "showconfigmem" {
		t.Fatalf("Unexpected Values.Encode (changed after EncodeToggle): %q", s)
	}
}

func TestLinks(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost/index?df=mp", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.ParseForm()
	df := NewParams().ENUM["df"]
	err = df.Decode(req.Form, nil)
	if err != nil {
		t.Fatal(err)
	}
	if num := df.Number; num.Negative || num.Uint != enums.Uint(enums.MP) {
		t.Errorf("Decode failed: %+v\n", num)
	}

	if total := df.EncodeUint("df", enums.TOTAL); total.Href != "?df=total" || total.Class != "state" || total.CaretClass != "" {
		t.Fatalf("Encode failed: total: %+v", total)
	}
	if mp := df.EncodeUint("df", enums.MP); mp.Href != "?df=-mp" || mp.Class != "state current dropup" || mp.CaretClass != "caret" {
		t.Fatalf("Encode failed: mp: %+v", mp)
	}

	if true {
		req, err := http.NewRequest("GET", "http://localhost/index?df=size", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.ParseForm()
		params := NewParams()
		err = params.ENUM["df"].Decode(req.Form, nil)
		if err == nil || err.Error() != "" {
			t.Fatalf("Error expected (%q)", err)
		}
		if s := params.Query.Values.Encode(); s != "df=total" {
			t.Fatalf("Expected Encode: %q", s)
		}
	}
	CheckRedirect(t, NewForm(t, "df=fs&ps=pid"), []string{"df"}, "df=-fs")
	CheckRedirect(t, NewForm(t, "df=fs&ps=pid"), []string{"df", "ps"}, "df=-fs&ps=-pid")
	CheckRedirect(t, NewForm(t, "df=fs&ps=pid"), []string{"ps"}, "ps=-pid")

	form := NewForm(t, "df=fs&ps=pid")
	CheckRedirect(t, form, []string{"ps"}, "ps=-pid")
	if err := form.Params.ENUM["df"].Decode(url.Values{"df": []string{"mp"}}, nil); err != nil {
		t.Fatalf("Decoding errd unexpectedly: %s", err)
	}
	if s, moved := form.Params.Query.Values.Encode(), "df=mp&ps=-pid"; s != moved {
		t.Fatalf("Redirect mismatch (%q): %q", moved, s)
	}
}

func CheckRedirect(t *testing.T, form Form, names []string, moved string) {
	for _, name := range names {
		err := form.Params.ENUM[name].Decode(form.Values, nil)
		if err == nil {
			t.Fatalf("RenamedConstError expected, got nil")
		}
		if _, ok := err.(enums.RenamedConstError); !ok {
			t.Fatalf("RenamedConstError expected, got: %s", err)
		}
	}
	if s := form.Params.Query.Values.Encode(); s != moved {
		t.Fatalf("Redirect mismatch (%q): %q", moved, s)
	}
}

type Form struct {
	url.Values
	*Params
}

func NewForm(t *testing.T, qs string) Form {
	req, err := http.NewRequest("GET", "http://localhost/index?"+qs, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.ParseForm()
	return Form{req.Form, NewParams()}
}
