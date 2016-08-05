// +build linux

package ostent

import (
	"net"
	"strings"
	"sync"

	pshost "github.com/shirou/gopsutil/host"
	psnet "github.com/shirou/gopsutil/net"
)

// ProcName returns procName back.
func ProcName(_ int, procName string) string { return procName }

func Distrib() (string, error) {
	platform, _, version, err := pshost.PlatformInformation()
	if err != nil {
		return "", err
	}
	if platform == "" {
		return "Docker", nil // Docker is a good guess.
	}
	platform = LSBID(platform)
	if version == "" {
		return platform, nil
	}
	return platform + " " + version, nil
}

// LSBID is to convert gopsutil platform identifier back to LSB ID form.
func LSBID(platform string) string {
	switch platform {
	case "redhat":
		return "RedHat"
	case "linuxmint":
		return "LinuxMint"
	case "scientific":
		return "ScientificSL"
	case "xenserver":
		return "XenServer"
	case "centos":
		return "CentOS"
	case "cloudlinux":
		return "CloudLinux"
	case "opensuse":
		return "OpenSUSE"
	case "suse":
		return "SUSE"
	}
	return strings.Title(platform)
}

// IF registers the interfaces with the reg.
func (ir *IndexRegistry) collectIF(wg *sync.WaitGroup) {
	ifaddrs, err := net.Interfaces()
	if err != nil {
		// err is gone
		wg.Done()
		return
	}
	if list, err := psnet.IOCounters(true); err == nil {
		// err is gone
		for _, iocounter := range list {
			if !HardwareIF(iocounter.Name) {
				continue
			}
			stat := &NetIO{IOCountersStat: iocounter}
			for _, ia := range ifaddrs {
				if ia.Name != iocounter.Name {
					continue
				}
				if addrs, err := ia.Addrs(); err == nil && len(addrs) > 0 {
					// err is gone
					stat.IP = addrs[0].String() // take just the first
				}
				break
			}
			ir.UpdateIF(stat)
		}
	}
	wg.Done()
}

type NetIO struct {
	psnet.IOCountersStat
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
