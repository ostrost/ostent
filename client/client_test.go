package client

import (
	"net/http"
	"net/url"
	"testing"
)

func TestLinks(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost/index?df=mp", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.ParseForm()

	la := &Linkattrs{
		Bimaps: map[string]Biseqmap{
			"df": DFBIMAP,
		},
	}
	DFSEQ := la.Param(req, url.Values{}, "df")
	if DFSEQ != DFMP {
		t.Fatalf("Linkattrs.Param failed")
	}
	if size := la.Attr("df", DFSIZE); size.Href != "?df=size" || size.Class != "state" || size.CaretClass != "" {
		t.Fatalf("Attr failed: size: %+v", size)
	}
	if mp := la.Attr("df", DFMP); mp.Href != "?df=-mp" || mp.Class != "state current dropup" || mp.CaretClass != "caret" {
		t.Fatalf("Attr failed: mp: %+v", mp)
	}
}
