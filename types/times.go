package types

import (
	"encoding/json"
	"fmt"
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
