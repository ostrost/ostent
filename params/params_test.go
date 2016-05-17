package params

import (
	"net/http"
	"testing"
	"time"

	"github.com/ostrost/ostent/flags"
)

var DelayFlags = flags.DelayBounds{
	Max: flags.Delay{Duration: time.Second * 2},
	Min: flags.Delay{Duration: time.Second},
}

func TestBoolLinks(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost/index?memn=0", nil)
	if err != nil {
		t.Fatal(err)
	}
	para := NewParams(DelayFlags)
	if err := para.Decode(req); err != nil {
		t.Fatal(err)
	}
	if para.Memn.Absolute != 0 {
		t.Errorf("Decode failed: %+v, expected %+v", para.Memn.Absolute, 0)
	}
}

func TestLinks(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost/index?dfk=2", nil)
	if err != nil {
		t.Fatal(err)
	}
	para := NewParams(DelayFlags)
	if err = para.Decode(req); err != nil {
		t.Fatal(err)
	}
	if para.Dfk.Negative || para.Dfk.Absolute != 2 {
		t.Errorf("Decode failed: %+v\n", para.Dfk)
	}
	total, err := Vlink(para, &para.Dfk, TOTAL, "")
	if err != nil {
		t.Fatal(err)
	}
	if total.LinkHref != "?dfk=6" || total.LinkClass != "state" || total.CaretClass != "" {
		t.Fatalf("Encode failed: total: %+v", total)
	}
	mp, err := Vlink(para, &para.Dfk, MP, "")
	if err != nil {
		t.Fatal(err)
	}
	if mp.LinkHref != "?dfk=-2" || mp.LinkClass != "state current" || mp.CaretClass != "caret" {
		t.Fatalf("Encode failed: mp: %+v", mp)
	}

	if true {
		req, err := http.NewRequest("GET", "http://localhost/index?dfk=2&df=anypreviousdfv", nil)
		if err != nil {
			t.Fatal(err)
		}
		para := NewParams(DelayFlags)
		if err := para.Decode(req); err == nil || err.Error() != "?dfk=2" {
			t.Fatalf("Error expected (%q)", err)
		}
	}
}
