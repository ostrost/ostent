package ostent

import (
	"fmt"
	"net"
	"testing"
)

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
		cmp: `   --------------------------
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
		bannerText(v.listenip, v.hostname, "ostent", v.addrsp, logger)
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
