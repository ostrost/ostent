package main

import "testing"

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
