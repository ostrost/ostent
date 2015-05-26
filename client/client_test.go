package client

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/ostrost/ostent/client/enums"
)

func TestLinks(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost/index?df=mp", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := NewParams(req)
	err = params["df"].Decode(req.Form, new(Number))
	if err != nil {
		t.Fatal(err)
	}
	num := params["df"].Decoded.Number
	if num.Negative || num.Uint != enums.Uint(enums.MP) {
		t.Errorf("Decode failed: %+v\n", num)
	}

	if total := params["df"].EncodeUint("df", enums.TOTAL); total.Href != "?df=total" || total.Class != "state" || total.CaretClass != "" {
		t.Fatalf("Encode failed: total: %+v", total)
	}
	if mp := params["df"].EncodeUint("df", enums.MP); mp.Href != "?df=-mp" || mp.Class != "state current dropup" || mp.CaretClass != "caret" {
		t.Fatalf("Encode failed: mp: %+v", mp)
	}

	if true {
		req, err := http.NewRequest("GET", "http://localhost/index?df=size", nil)
		if err != nil {
			t.Fatal(err)
		}
		params := NewParams(req)
		err = params["df"].Decode(req.Form, new(Number))
		if err == nil || err.Error() != "" {
			t.Fatalf("Error expected (%q)", err)
		}
		if s := params.Encode(); s != "df=total" {
			t.Fatalf("Expected Encode: %q", s)

		}
	}
	CheckRedirect(t, new(Number), NewForm(t, "df=fs&ps=pid"), []string{"df"}, "df=-fs")
	CheckRedirect(t, new(Number), NewForm(t, "df=fs&ps=pid"), []string{"df", "ps"}, "df=-fs&ps=-pid")
	CheckRedirect(t, new(Number), NewForm(t, "df=fs&ps=pid"), []string{"ps"}, "ps=-pid")
}

func CheckRedirect(t *testing.T, num *Number, form Form, names []string, moved string) {
	for _, name := range names {
		err := form.Params[name].Decode(form.Values, num)
		if err == nil {
			t.Fatalf("RenamedConstError expected, got nil")
		}
		if _, ok := err.(enums.RenamedConstError); !ok {
			t.Fatalf("RenamedConstError expected, got: %s", err)
		}
	}
	if s := form.Params.Encode(); s != moved {
		t.Fatalf("Redirect mismatch (%q): %q", moved, s)
	}
}

type Form struct {
	url.Values
	Params
}

func NewForm(t *testing.T, qs string) Form {
	req, err := http.NewRequest("GET", "http://localhost/index?"+qs, nil)
	if err != nil {
		t.Fatal(err)
	}
	return Form{req.Form, NewParams(req)}
}
