package ostent

import (
	"net"
	"testing"
)

// NetAddressIP returns IP address of the first found hardware network interface
func NetAddressIP() (string, error) {
	// list of the system's network interfaces.
	listIface, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range listIface {
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if !HardwareInterface(iface.Name) {
			continue
		}
		if aa, err := iface.Addrs(); err == nil {
			if len(aa) == 0 {
				continue
			}
			for _, a := range aa {
				ipnet, ok := a.(*net.IPNet)
				if !ok {
					// log.Fatalf("Unable to cast: %v", a)
					continue
				}
				if !ipnet.IP.IsLinkLocalUnicast() {
					return ipnet.IP.String(), nil
				}
			}
		}
	}
	return "", nil
}

func TestInterfaceIP(t *testing.T) {
	fip := &FoundIP{}
	if err := (&Machine{}).ApplyperInterface(fip.Next); err != nil {
		t.Error(err)
		return
	}
	nip, err := NetAddressIP()
	if err != nil {
		t.Error(err)
		return
	}
	if fip.string != nip && nip != "127.0.0.2" { // travis(linux) has just lo?
		t.Errorf("Mismatch:\nExpected: %+v\nGot     : %+v\n", nip, fip.string)
	}
}
