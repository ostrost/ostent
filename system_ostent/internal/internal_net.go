package internal

import (
	"net"
	"sync"

	psnet "github.com/shirou/gopsutil/net"
)

func interfaceIPByName(interfaces []psnet.InterfaceStat, name string) (string, bool) {
	for _, fi := range interfaces {
		if addrs := fi.Addrs; fi.Name == name && len(addrs) > 0 {
			if ip, _, err := net.ParseCIDR(addrs[0].Addr); err == nil { // err is ignored
				return ip.String(), true
			}
		}
	}
	return "", false
}

func AddTags(interfaces []psnet.InterfaceStat, ioname string, isLoopback bool, tags map[string]string) {
	if ip, ok := interfaceIPByName(interfaces, ioname); ok {
		tags["ip"] = ip
	}
	if isLoopback {
		tags["nonemptyifLoopback"] = "nonempty"
	}
}

type LastNetIOStats struct {
	mutexLast sync.Mutex
	last      map[string]psnet.IOCountersStat
}

func (s *LastNetIOStats) AddDeltaFields(io psnet.IOCountersStat, fields map[string]interface{}) {
	s.mutexLast.Lock()
	defer s.mutexLast.Unlock()
	if last, ok := s.last[io.Name]; ok {
		for k, v := range map[string]interface{}{
			"delta_bytes_sent":   io.BytesSent - last.BytesSent,
			"delta_bytes_recv":   io.BytesRecv - last.BytesRecv,
			"delta_packets_sent": io.PacketsSent - last.PacketsSent,
			"delta_packets_recv": io.PacketsRecv - last.PacketsRecv,
			"delta_err_in":       io.Errin - last.Errin,
			"delta_err_out":      io.Errout - last.Errout,
			"delta_drop_in":      io.Dropin - last.Dropin,
			"delta_drop_out":     io.Dropout - last.Dropout,
		} {
			fields[k] = v
		}
	}
	if s.last == nil {
		s.last = make(map[string]psnet.IOCountersStat)
	}
	s.last[io.Name] = io
}
