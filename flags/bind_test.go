package flags

import (
	"fmt"
	"testing"
)

func TestBindSet(t *testing.T) {
	for _, v := range []struct {
		a    string
		cmp  string
		errs []error
	}{
		{"a:b:", "", []error{fmt.Errorf("too many colons in address " + "a:b:")}},
		{"localhost:nonexistent", "localhost:nonexistent", []error{
			fmt.Errorf("unknown port tcp/nonexistent"),
			fmt.Errorf("lookup tcp/nonexistent: Servname not supported for ai_socktype"),
			fmt.Errorf("lookup tcp/nonexistent: nodename nor servname provided, or not known"),
		}},
		{"localhost", "localhost:9050", nil},
		{"", ":9050", nil},
		{":8001", ":8001", nil},
		{"*:8001", "*:8001", nil},
		{"127", "127:9050", nil},
		{"127.1", "127.1:9050", nil},
		{"127.0.0.1", "127.0.0.1:9050", nil},
		{"127.1:8001", "127.1:8001", nil},
		{"127.0.0.1:8001", "127.0.0.1:8001", nil},
	} {
		b := NewBind("", 9050)
		if err := b.Set(v.a); err != nil {
			unknownerr := true
			for _, x := range v.errs {
				if err.Error() == x.Error() {
					unknownerr = false
					break
				}
			}
			if !unknownerr {
				continue
			}
			t.Errorf("Error: %q\nExpected errors: %+v\n", err, v.errs)
		}
		if s := b.String(); s != v.cmp {
			t.Errorf("Mismatch: Bind %v == %v != %v\n", v.a, v.cmp, s)
		}
	}
}
