package ostential
import (
	"testing"
)

func TestParseArgs(t *testing.T) {
	const defport = "9050"
	for _, v := range []struct{
		a string
		cmp string
	}{
		{   "8001", ":8001"},
		{  ":8001", ":8001"},
		{ "*:8001", ":8001"},
		{ "127.1:8001",     "127.1:8001"},
		{ "127.0.0.1:8001", "127.0.0.1:8001"},
		{ "127.0.0.1",      "127.0.0.1:"+ defport},
		{ "127",            "127.0.0.1:"+ defport},
		{ "127.1",          "127.1:"    + defport},
	} {
		cmp, err := parseaddr(v.a, defport)
		if err != nil {
			t.Error(err)
		}
		if cmp != v.cmp {
			t.Errorf("Mismatch: parseaddr(%v) == %v != %v\n", v.a, v.cmp, cmp)
		}
	}
}
