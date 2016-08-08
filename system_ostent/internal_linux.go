// +build linux

package system_ostent

import (
	"github.com/shirou/gopsutil/net"
)

// iocounters is to call net.IOCounters(pernic).
func iocounters(pernic bool) ([]net.IOCountersStat, error) {
	return net.IOCounters(pernic)
}
