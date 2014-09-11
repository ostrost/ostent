package types

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"
)

// Duration type derives time.Duration
type Duration time.Duration

// String returns Duration string representation
func (d Duration) String() string {
	s := time.Duration(d).String()
	if strings.HasSuffix(s, "m0s") {
		s = strings.TrimSuffix(s, "0s")
	}
	if strings.HasSuffix(s, "h0m") {
		s = strings.TrimSuffix(s, "0m")
	}
	return s
}

// MarshalJSON is for encoding/json marshaling into Duration string representation
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

type PeriodValue struct {
	Duration
	Above *Duration // optional
}

func (pv *PeriodValue) Set(input string) error {
	v, err := time.ParseDuration(input)
	if err != nil {
		return err
	}
	if v < time.Second { // hard coded
		return fmt.Errorf("Less than a second: %s", v)
	}
	if v%time.Second != 0 {
		return fmt.Errorf("Not a multiple of a second: %s", v)
	}
	if pv.Above != nil && v < time.Duration(*pv.Above) {
		return fmt.Errorf("Should be above %s: %s", *pv.Above, v)
	}
	pv.Duration = Duration(v)
	return nil
}

type BindValue struct {
	string
	defport string // const
	Host    string // available after flag.Parse()
	Port    string // available after flag.Parse()
}

func NewBindValue(defstring, defport string) BindValue {
	bv := BindValue{defport: defport}
	bv.Set(defstring)
	return bv
}

// satisfying flag.Value interface
func (bv BindValue) String() string { return string(bv.string) }
func (bv *BindValue) Set(input string) error {
	if input == "" {
		bv.Port = bv.defport
	} else {
		if !strings.Contains(input, ":") {
			input = ":" + input
		}
		var err error
		bv.Host, bv.Port, err = net.SplitHostPort(input)
		if err != nil {
			return err
		}
		if bv.Host == "*" {
			bv.Host = ""
		} else if bv.Port == "127" {
			bv.Host = "127.0.0.1"
			bv.Port = bv.defport
		}
		if _, err = net.LookupPort("tcp", bv.Port); err != nil {
			if bv.Host != "" {
				return err
			}
			bv.Host, bv.Port = bv.Port, bv.defport
		}
	}

	bv.string = bv.Host + ":" + bv.Port
	return nil
}
