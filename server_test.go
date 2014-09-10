package ostent

import (
	"errors"
	"fmt"
	"net"
	"testing"
)

func TestParseArgs(t *testing.T) {
	const defport = "9050"
	for _, v := range []struct {
		a   string
		cmp string
		err error
	}{
		{"a:b:", "", errors.New("too many colons in address a:b:")},
		{"localhost:nonexistent", "", errors.New("unknown port tcp/nonexistent")},
		{"localhost", "localhost:9050", nil},
		{"", ":9050", nil},
		{"8001", ":8001", nil},
		{"8001", ":8001", nil},
		{":8001", ":8001", nil},
		{"*:8001", ":8001", nil},
		{"127.1:8001", "127.1:8001", nil},
		{"127.0.0.1:8001", "127.0.0.1:8001", nil},
		{"127.0.0.1", "127.0.0.1:" + defport, nil},
		{"127", "127.0.0.1:" + defport, nil},
		{"127.1", "127.1:" + defport, nil},
	} {
		bv := newBind(v.a, defport) // double Set, should be ok
		if err := bv.Set(v.a); err != nil {
			if err.Error() != v.err.Error() {
				t.Errorf("Error: %s\nMismatch: %s\n", err, v.err)
			}
			continue
		}
		if bv.string != v.cmp {
			t.Errorf("Mismatch: bindFlag %v == %v != %v\n", v.a, v.cmp, bv.string)
		}
	}
}

type testLogger struct {
	tester *testing.T
	buffer string
}

func (l *testLogger) Print(v ...interface{}) {
	s := fmt.Sprint(v...)
	if len(s) > 0 && s[len(s)-1] != '\n' {
		l.tester.Error("Banner line should end with \\n")
	}
	l.buffer += s
}

func TestBanner(t *testing.T) {
	_, net1, _ := net.ParseCIDR("127.0.0.1/32")
	_, net2, _ := net.ParseCIDR("127.0.0.2/32")
	_, net3, _ := net.ParseCIDR("fe80::/10")
	localaddrs := []net.Addr{net1, net2, net3}

	for _, v := range []struct {
		hostname string
		listenip string
		cmp      string
		addrsp   *[]net.Addr
	}{{
		hostname: "testhostname24charswidth",
		listenip: "127.0.0.1",
		cmp: `   -------------------------------
 / testhostname24ch... ostent \
+------------------------------+
| http://127.0.0.1             |
+------------------------------+
`,
	}, {
		hostname: "abc",
		listenip: "[::]:7050",
		addrsp:   &localaddrs,
		cmp: `   ----------
 / abc ostent \
+------------------------------+
| http://127.0.0.1:7050        |
|------------------------------|
| http://127.0.0.2:7050        |
+------------------------------+
`,
	}} {
		logger := &testLogger{tester: t}
		bannerText(v.listenip, v.hostname, v.addrsp, logger)
		if logger.buffer != v.cmp {
			t.Error(Mismatch{logger.buffer, v.cmp})
		}
	}
}

type Mismatch struct {
	output, expect interface{}
}

func (m Mismatch) Error() string {
	return fmt.Sprintf("Failed.\nOutput: %#v\nExpect: %#v", m.output, m.expect)
}
