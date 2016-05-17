package flags

import (
	"net"
	"strconv"
	"strings"
)

type Bind struct {
	DefaultHost string
	DefaultPort string
	Host, Port  string
}

func NewBind(defhost string, defport int) Bind {
	b := Bind{
		DefaultHost: defhost,
		DefaultPort: strconv.Itoa(defport),
	}
	_ = b.Set("") // must not err
	return b
}

// String is to conform interfaces (flag.Value, fmt.Stringer).
func (b Bind) String() string { return b.Host + ":" + b.Port }

func (b *Bind) Set(input string) error {
	if input == "" {
		input = b.DefaultHost + ":" + b.DefaultPort
	}
	var err error
	b.Host, b.Port, err = net.SplitHostPort(input)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "missing port in address") {
			return err
		}
		b.Host, b.Port = input, b.DefaultPort
	}
	if _, err = net.LookupPort("tcp", b.Port); err != nil {
		return err
	}
	return nil
}

// Type of pflag.Value interface
func (b Bind) Type() string { return "bind" }
