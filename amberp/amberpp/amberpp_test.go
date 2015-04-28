package main

import (
	"strings"
	"testing"
	"text/template/parse"
)

func TestKeysSorted(t *testing.T) {
	m := map[string]*parse.Tree{"9": nil, "8": nil, "7": nil, "6": nil, "5": nil, "4": nil, "3": nil, "2": nil, "1": nil, "0": nil}
	if s, expected := strings.Join(KeysSorted(m), ""), "0123456789"; s != expected {
		t.Errorf("KeysSorted mismatch: %q (expected %q)", s, expected)
	}
}

func TestCompile(t *testing.T) {
	s, err := compile([]byte(`
html
  body
    div.container#first Hello
`), true, true)
	if err != nil {
		t.Error(err)
	}
	if expected := `<html>
	<body>
		<div className="container" id="first">Hello</div>
	</body>
</html>
`; s != expected {
		t.Errorf("compile mismatch: %q (expected %q)", s, expected)
	}
}
