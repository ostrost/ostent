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

	if total := links.EncodeNU("df", TOTAL); total.Href != "?df=total" || total.Class != "state" || total.CaretClass != "" {
		t.Fatalf("Encode failed: total: %+v", total)
	}
	if mp := links.EncodeNU("df", MP); mp.Href != "?df=-mp" || mp.Class != "state current dropup" || mp.CaretClass != "caret" {
		t.Fatalf("Encode failed: mp: %+v", mp)
	}

	if true {
		req, err := http.NewRequest("GET", "http://localhost/index?df=size", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.ParseForm()
		links := NewLinks()
		err = DF.Decode(req.Form, "df", links, new(Number), new(UintDF))
		if err == nil || err.Error() != "" {
			t.Fatalf("Error expected (%q)", err)
		}
		if s := links.Values.Encode(); s != "df=total" {
			t.Fatalf("Expected Encode: %q", s)

		}
	}
	checklinks := NewLinks()
	CheckRedirect(t, checklinks, new(UintDF), DF, "df", "fs", "df=-fs")
	CheckRedirect(t, checklinks, new(UintPS), PS, "ps", "pid", "df=-fs&ps=-pid")
	CheckRedirect(t, NewLinks(), new(UintPS), PS, "ps", "pid", "ps=-pid")
}

func CheckRedirect(t *testing.T, linker LinkerEncoder, uptr Upointer, decoder Decoder, name, qsend, moved string) {
	req, err := http.NewRequest("GET", "http://localhost/index?"+name+"="+qsend, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.ParseForm()
	err = decoder.Decode(req.Form, name, linker, new(Number), uptr)
	if err == nil {
		t.Fatalf("RenamedConstError expected, got nil")
	}
	if _, ok := err.(RenamedConstError); !ok {
		t.Fatalf("RenamedConstError expected, got: %s", err)
	}
	if s := linker.Encode(); s != moved {
		t.Fatalf("Redirect mismatch (%q): %q", moved, s)
	}
}

type LinkerEncoder interface {
	Linker
	Encode() string // from url.Values
}
