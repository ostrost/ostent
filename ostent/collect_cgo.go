// +build sigar,cgo

package ostent

import (
	"sync"

	"github.com/ostrost/ostent/getifaddrs"
)

// IF registers the interfaces with the reg.
func (m Machine) IF(reg Registry, wg *sync.WaitGroup) {
	// m is unused
	if ifaddrs, err := getifaddrs.Getifaddrs(); err == nil {
		// err is gone
		for _, ifaddr := range ifaddrs {
			if !HardwareIF(ifaddr.GetName()) {
				continue
			}
			reg.UpdateIF(&ifaddr) // pointer not to copy everywhere
		}
	}
	wg.Done()
}
