// +build linux

package internal

import (
	"github.com/shirou/gopsutil/net"
)

// IOCounters is to call net.IOCounters(pernic).
func IOCounters(pernic bool) ([]net.IOCountersStat, error) {
	return net.IOCounters(pernic)
}
