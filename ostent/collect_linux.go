// +build linux

package ostent

import (
	"net"
	"sync"

	psnet "github.com/shirou/gopsutil/net"
)

// IF registers the interfaces with the reg.
func (m Machine) IF(reg Registry, wg *sync.WaitGroup) {
	// m is unused
	ifaddrs, err := net.Interfaces()
	if err != nil {
		// err is gone
		wg.Done()
		return
	}
	if list, err := psnet.NetIOCounters(true); err == nil {
		// err is gone
		for _, netio := range list {
			if !HardwareIF(netio.Name) {
				continue
			}
			stat := &NetIO{NetIOCountersStat: netio}
			for _, ia := range ifaddrs {
				if ia.Name != netio.Name {
					continue
				}
				if addrs, err := ia.Addrs(); err == nil && len(addrs) > 0 {
					// err is gone
					stat.IP = addrs[0].String() // take just the first
				}
				break
			}
			reg.UpdateIF(stat)
		}
	}
	wg.Done()
}

type NetIO struct {
	psnet.NetIOCountersStat
	IP string
}

func (io NetIO) GetName() string  { return io.Name }
func (io NetIO) GetIP() string    { return io.IP }
func (io NetIO) BytesIn() uint    { return uint(io.BytesRecv) }
func (io NetIO) BytesOut() uint   { return uint(io.BytesSent) }
func (io NetIO) DropsIn() uint    { return uint(io.Dropin) }
func (io NetIO) DropsOut() uint   { return uint(io.Dropout) }
func (io NetIO) ErrorsIn() uint   { return uint(io.Errin) }
func (io NetIO) ErrorsOut() uint  { return uint(io.Errout) }
func (io NetIO) PacketsIn() uint  { return uint(io.PacketsRecv) }
func (io NetIO) PacketsOut() uint { return uint(io.PacketsSent) }
