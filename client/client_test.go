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

	if size := links.EncodeNU("df", DFSIZE); size.Href != "?df=dfsize" || size.Class != "state" || size.CaretClass != "" {
		t.Fatalf("Encode failed: size: %+v", size)
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
		unused := Number{}
		err = DF.Decode(req.Form, "df", links, &unused, new(UintDF))
		if err == nil || err.Error() != "" {
			t.Fatalf("Error expected (%q)", err)
		}
		if s := links.Values.Encode(); s != "df=dfsize" {
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
	unused := Number{}
	err = decoder.Decode(req.Form, name, linker, &unused, uptr)
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