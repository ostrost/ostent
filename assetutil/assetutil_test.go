package assetutil

import (
	"net/http"
	"testing"
)

// TestFQJSANSlice tests for FQJSANSlice transformation.
func TestFQJSANSlice(t *testing.T) {
	for _, test := range []struct {
		Fst, JS string
		Req     *http.Request
	}{
		{"//localhost/a", "/a", &http.Request{Host: "localhost"}},
		{"//localhost/b", "//localhost/b", nil},
	} {
		fq := FQJSANSlice(JSANSlice{JSAN(test.JS)}, test.Req)
		if len(fq) != 1 {
			t.Errorf("FQJSANSlice: length %d instead of 1\n", len(fq))
		}
		if string(fq[0]) != test.Fst {
			t.Errorf("FQJSANSlice: [0] %q instead of %q\n", fq[0], test.Fst)
		}
	}
}

// TestLessJSANFunc tests for JSAN comparison.
func TestLessJSANFunc(t *testing.T) {
	for n, test := range []struct {
		Asnd, Bfst string
	}{
		{"a", "b"},
		{"a", "jquery"},
		{"bootstrap", "jquery"},
	} {
		if LessJSANFunc(100)(JSAN(test.Asnd), JSAN(test.Bfst)) {
			t.Errorf("#%d LessJSANFunc(%q, %q) failed. Expected true\n", n, test.Asnd, test.Bfst)
		}
	}
}

// TestJSassetNames tests JSassetNames results.
func TestJSassetNames(t *testing.T) {
	sorted := JSassetNames([]string{"a.js", "bootstrap.js", "jquery.js", "b.js"})
	sorted.StableSortBy(LessJSANFunc(len(sorted))) // 2nd stable sort
	if string(sorted[0]) != "/jquery.js" ||
		string(sorted[1]) != "/bootstrap.js" ||
		string(sorted[2]) != "/a.js" ||
		string(sorted[3]) != "/b.js" {
		t.Errorf("JSassetNames sorted: %+v\n", sorted)
	}
}
