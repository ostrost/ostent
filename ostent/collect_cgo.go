// +build cgo

package ostent

import (
	"sync"

	"github.com/ostrost/ostent/getifaddrs"
)

// IF registers the interfaces with the reg.
func (m Machine) IF(reg Registry, wg *sync.WaitGroup) {
	// m is unused
	if gotifaddrs, err := getifaddrs.Getifaddrs(); err == nil {
		// err is gone
		for _, ifdata := range gotifaddrs {
			if !HardwareIF(ifdata.Name) {
				continue
			}
			reg.UpdateIF(IfData{
				IP:         ifdata.IP,
				Name:       ifdata.Name,
				InBytes:    ifdata.InBytes,
				OutBytes:   ifdata.OutBytes,
				InDrops:    ifdata.InDrops,
				OutDrops:   ifdata.OutDrops,
				InErrors:   ifdata.InErrors,
				OutErrors:  ifdata.OutErrors,
				InPackets:  ifdata.InPackets,
				OutPackets: ifdata.OutPackets,
			})
		}
	}
	wg.Done()
}
