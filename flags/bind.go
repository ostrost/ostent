package flags

import (
	"net"
	"strconv"
	"strings"
)

type Bind struct {
	defaultHost string
	defaultPort string
	Host, port  string
}

func NewBind(defhost string, defport int) Bind {
	b := Bind{
		defaultHost: defhost,
		defaultPort: strconv.Itoa(defport),
	}
	_ = b.Set("") // must not err
	return b
}

// String is to conform interfaces (flag.Value, fmt.Stringer).
func (b Bind) String() string { return b.Host + ":" + b.port }

func (b *Bind) Set(input string) error {
	if input == "" {
		input = b.defaultHost + ":" + b.defaultPort
	}
	var err error
	b.Host, b.port, err = net.SplitHostPort(input)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "missing port in address") {
			return err
		}
		b.Host, b.port = input, b.defaultPort
	}
	if _, err = net.LookupPort("tcp", b.port); err != nil {
		return err
	}
	return nil
}

// Type of pflag.Value interface
func (b Bind) Type() string { return "bind" }
