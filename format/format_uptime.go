//+build none

package format

import (
	"fmt"
	"time"
)

func format_uptime(seconds uint64) string {
	d := time.Duration(seconds) * time.Second
	s := ""
	if d > time.Duration(24)*time.Hour {
		days := d / time.Hour / 24
		end := ""
		if days > 1 {
			end = "s"
		}
		s += fmt.Sprintf("%d day%s, ", days, end)
	}
	t := time.Unix(int64(seconds), 0).UTC()
	tf := t.Format("15:04")
	if tf[0] == '0' {
		tf = " " + tf[1:]
	}
	s += tf
	return s
}
