package ostent

import (
	"regexp"
	"runtime"

	"github.com/ostrost/ostent/getifaddrs"
)

var (
	rx_lo      = regexp.MustCompile("lo\\d*") // "lo" & lo\d+; used in interfaces_unix.go, sortable.go
	RX_fw      = regexp.MustCompile("fw\\d+")
	RX_gif     = regexp.MustCompile("gif\\d+")
	RX_stf     = regexp.MustCompile("stf\\d+")
	RX_wdl     = regexp.MustCompile("awdl\\d+")
	RX_bridge  = regexp.MustCompile("bridge\\d+")
	RX_vboxnet = regexp.MustCompile("vboxnet\\d+")
	RX_airdrop = regexp.MustCompile("p2p\\d+")
)

func realInterface(name string) bool {
	if RX_bridge.MatchString(name) ||
		RX_vboxnet.MatchString(name) {
		return false
	}
	if runtime.GOOS == "darwin" {
		if RX_fw.MatchString(name) ||
			RX_gif.MatchString(name) ||
			RX_stf.MatchString(name) ||
			RX_wdl.MatchString(name) ||
			RX_airdrop.MatchString(name) {
			return false
		}
	}
	return true
}

// getInterfaces registers the interfaces with the reg and send first non-loopback IP to the chan
func getInterfaces(reg Register, CH chan<- string) {
	iflist, _ := getifaddrs.Getifaddrs()
	IP := ""
	for _, ifdata := range iflist {
		if !realInterface(ifdata.Name) {
			continue
		}
		if ifdata.InBytes == 0 &&
			ifdata.OutBytes == 0 &&
			ifdata.InPackets == 0 &&
			ifdata.OutPackets == 0 &&
			ifdata.InErrors == 0 &&
			ifdata.OutErrors == 0 {
			continue
		}
		if IP == "" && !rx_lo.MatchString(ifdata.Name) {
			// first non-loopback IP
			IP = ifdata.IP
		}
		reg.UpdateIFdata(ifdata)
	}
	CH <- IP
}

/* using net.Interfaces
func netinterface_ipaddr() (string, error) {
	// list of the system's network interfaces.
	list_iface, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	var addr []string

	for _, iface := range list_iface {
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if !realInterfaceName(iface.Name) {
			continue
		}
		if aa, err := iface.Addrs(); err == nil {
			if len(aa) == 0 {
				continue
			}
			for _, a := range aa {
				ipnet, ok := a.(*net.IPNet)
				if !ok {
					return "", fmt.Errorf("Not an IP: %v", a)
					continue
				}
				if ipnet.IP.IsLinkLocalUnicast() {
					continue
				}
				addr = append(addr, ipnet.IP.String())
			}
		}
	}
	if len(addr) == 0 {
		return "", nil
	}
	return addr[0], nil
} // */
