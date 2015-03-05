package flags

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Period struct {
	time.Duration
	Above *time.Duration // optional
}

// String returns Period string representation
func (p Period) String() string {
	s := p.Duration.String()
	if strings.HasSuffix(s, "m0s") {
		s = strings.TrimSuffix(s, "0s")
	}
	if strings.HasSuffix(s, "h0m") {
		s = strings.TrimSuffix(s, "0m")
	}
	return s
}

// MarshalJSON is for encoding/json marshaling into Period string representation
func (p Period) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func (p *Period) Set(input string) error {
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
	if p.Above != nil && v < *p.Above {
		return fmt.Errorf("Should be above %s: %s", *p.Above, v)
	}
	p.Duration = v
	return nil
}
