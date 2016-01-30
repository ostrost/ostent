package flags

import (
	"fmt"
	"net"
	"strings"
)

type Bind struct {
	defport string // const
	Host    string // available after flag.Parse()
	Port    string // available after flag.Parse()
}

func NewBind(defport int) Bind {
	b := Bind{defport: fmt.Sprintf("%d", defport)}
	b.Set("")
	return b
}

// String is to conform interfaces (flag.Value, fmt.Stringer).
func (b Bind) String() string { return string(b.Host + ":" + b.Port) }

// ClientString returns b string suitable for client.
// If the host part of b is empty, 127.0.0.1 is assumed.
func (b Bind) ClientString() string {
	if b.Host == "" {
		b.Host = "127.0.0.1" // b is a copy (not a pointer)
	}
	return b.String()
}

func (b *Bind) Set(input string) error {
	if input == "" {
		input = ":" + b.defport
	}
	var err error
	b.Host, b.Port, err = net.SplitHostPort(input)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "missing port in address") {
			return err
		}
		b.Host, b.Port = input, b.defport
	}
	if _, err = net.LookupPort("tcp", b.Port); err != nil {
		return err
	}
	return nil
}

// Type of pflag.Value interface
func (b Bind) Type() string { return "bind" }
