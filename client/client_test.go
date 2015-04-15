package client

import (
	"net/http"
	"testing"
)

func TestLinks(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost/index?df=mp", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.ParseForm()
	links := NewLinks()
	num := Number{}
	err = DF.Decode(req.Form, "df", links, &num, new(UintDF))
	if err != nil {
		t.Fatal(err)
	}
	if num.Negative || num.Uint != Uint(MP) {
		t.Errorf("Decode failed: %+v\n", num)
	}

	if size := links.Encode("df", DFSIZE); size.Href != "?df=dfsize" || size.Class != "state" || size.CaretClass != "" {
		t.Fatalf("Encode failed: size: %+v", size)
	}
	if mp := links.Encode("df", MP); mp.Href != "?df=-mp" || mp.Class != "state current dropup" || mp.CaretClass != "caret" {
		t.Fatalf("Encode failed: mp: %+v", mp)
	}

	if true {
		req, err := http.NewRequest("GET", "http://localhost/index?df=size", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.ParseForm()
		links := NewLinks()
		num := Number{}
		err = DF.Decode(req.Form, "df", links, &num, new(UintDF))
		if err == nil || err.Error() != "" {
			t.Fatalf("Error expected (%q)", err)
		}
		if s := links.Values.Encode(); s != "df=dfsize" {
			t.Fatalf("Expected Encode: %q", s)

		}
	}
}
