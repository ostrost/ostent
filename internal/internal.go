package internal

import (
	"strconv"
	"time"
)

type Duration struct{ Duration time.Duration }

func (d *Duration) UnmarshalTOML(b []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(b[1 : len(b)-1]))
	if err == nil {
		return nil
	}
	if n, err := strconv.ParseInt(string(b), 10, 64); err == nil {
		d.Duration = time.Second * time.Duration(n)
	} else if f, err := strconv.ParseFloat(string(b), 64); err == nil {
		d.Duration = time.Second * time.Duration(f)
	}
	return nil
}
