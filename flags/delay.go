package flags

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type DelayBounds struct{ Max, Min Delay }

type Delay struct {
	time.Duration
	Above *time.Duration // optional
	Below *time.Duration // optional
}

// DurationString returns string representation of dur.
func DurationString(dur time.Duration) string {
	s := dur.String()
	if strings.HasSuffix(s, "m0s") {
		s = strings.TrimSuffix(s, "0s")
	}
	if strings.HasSuffix(s, "h0m") {
		s = strings.TrimSuffix(s, "0m")
	}
	return s
}

// String returns Delay string representation
func (d Delay) String() string { return DurationString(d.Duration) }

// MarshalJSON is for encoding/json marshaling into Delay string representation
func (d Delay) MarshalJSON() ([]byte, error) { return json.Marshal(d.String()) }

func (d *Delay) Set(input string) error {
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
	if d.Above != nil && v < *d.Above {
		return fmt.Errorf("Should be above %s: %s", *d.Above, v)
	}
	if d.Below != nil && v > *d.Below {
		return fmt.Errorf("Should be below %s: %s", *d.Below, v)
	}
	d.Duration = v
	return nil
}

// Type of pflag.Value interface
func (d Delay) Type() string { return "delay" }
