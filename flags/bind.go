package flags

import (
	"fmt"
	"net"
	"strings"
)

type Bind struct {
	string
	defport string // const
	Host    string // available after flag.Parse()
	Port    string // available after flag.Parse()
}

func NewBind(defport int) Bind {
	b := Bind{defport: fmt.Sprintf("%d", defport)}
	b.Set("")
	return b
}

// satisfying flag.Value interface
func (b Bind) String() string { return string(b.string) }
func (b *Bind) Set(input string) error {
	if input == "" {
		b.Port = b.defport
	} else {
		if !strings.Contains(input, ":") {
			input = ":" + input
		}
		var err error
		b.Host, b.Port, err = net.SplitHostPort(input)
		if err != nil {
			return err
		}
		if b.Host == "*" {
			b.Host = ""
		} else if b.Port == "127" {
			b.Host = "127.0.0.1"
			b.Port = b.defport
		}
		if _, err = net.LookupPort("tcp", b.Port); err != nil {
			if b.Host != "" {
				return err
			}
			b.Host, b.Port = b.Port, b.defport
		}
	}
	b.string = b.Host + ":" + b.Port
	return nil
}
