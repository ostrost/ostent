// +build cgo

package ostent

import (
	"sync"

	"github.com/ostrost/ostent/getifaddrs"
)

// ApplyperInterface calls apply for each found hardware interface.
func (m Machine) ApplyperInterface(apply func(getifaddrs.IfData) bool) error {
	// m is unused
	gotifaddrs, err := getifaddrs.Getifaddrs()
	if err != nil {
		return err
	}
	for _, ifdata := range gotifaddrs {
		if !HardwareInterface(ifdata.Name) {
			continue
		}
		if !apply(ifdata) {
			break
		}
	}
	return nil
}

// Interfaces registers the interfaces with the reg and first non-loopback IP with the sreg.
func (m Machine) Interfaces(reg Registry, sreg S2SRegistry, wg *sync.WaitGroup) {
	fip := FoundIP{}
	m.ApplyperInterface(func(ifdata getifaddrs.IfData) bool {
		fip.Next(ifdata)
		if ifdata.InBytes == 0 &&
			ifdata.OutBytes == 0 &&
			ifdata.InPackets == 0 &&
			ifdata.OutPackets == 0 &&
			ifdata.InErrors == 0 &&
			ifdata.OutErrors == 0 {
			// nothing
		} else {
			reg.UpdateIFdata(IfData{
				Name:       ifdata.Name,
				InBytes:    ifdata.InBytes,
				OutBytes:   ifdata.OutBytes,
				InErrors:   ifdata.InErrors,
				OutErrors:  ifdata.OutErrors,
				InPackets:  ifdata.InPackets,
				OutPackets: ifdata.OutPackets,
			})
		}
		return true
	})
	sreg.SetString("ip", fip.string)
	wg.Done()
}

type FoundIP struct {
	string
}

func (fip *FoundIP) Next(ifdata getifaddrs.IfData) bool {
	if fip.string != "" {
		return false
	}
	if !RXlo.MatchString(ifdata.Name) { // non-loopback
		fip.string = ifdata.IP
		return false
	}
	return true
}
