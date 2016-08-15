package system_ostent

import (
	"fmt"
	"net"

	psnet "github.com/shirou/gopsutil/net"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"

	"github.com/ostrost/ostent/system_ostent/internal"
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

func (s *NetIOStats) addDeltaFields(io psnet.IOCountersStat, fields map[string]interface{}) {
	if last, ok := s.last[io.Name]; ok {
		deltaFields := map[string]interface{}{
			"delta_bytes_sent":   io.BytesSent - last.BytesSent,
			"delta_bytes_recv":   io.BytesRecv - last.BytesRecv,
			"delta_packets_sent": io.PacketsSent - last.PacketsSent,
			"delta_packets_recv": io.PacketsRecv - last.PacketsRecv,
			"delta_err_in":       io.Errin - last.Errin,
			"delta_err_out":      io.Errout - last.Errout,
			"delta_drop_in":      io.Dropin - last.Dropin,
			"delta_drop_out":     io.Dropout - last.Dropout,
		}
		for k, v := range deltaFields {
			fields[k] = v
		}
	}
	if s.last == nil {
		s.last = make(map[string]psnet.IOCountersStat)
	}
	s.last[io.Name] = io
}

type PS interface{}
type systemPS struct{}

type NetIOStats struct {
	last map[string]psnet.IOCountersStat

	ps PS

	skipChecks bool
	Interfaces []string
}

func (_ *NetIOStats) Description() string {
	return "Read metrics about network interface usage"
}

var netSampleConfig = `
  ## By default, telegraf gathers stats from any up interface (excluding loopback)
  ## Setting interfaces will tell it to gather these explicit interfaces,
  ## regardless of status.
  ##
  # interfaces = ["eth0"]
`

func (_ *NetIOStats) SampleConfig() string {
	return netSampleConfig
}

func (s *NetIOStats) Gather(acc telegraf.Accumulator) error {
	netio, err := internal.IOCounters(true)
	if err != nil {
		return fmt.Errorf("error getting net io info: %s", err)
	}

	interfaces, err := psnet.Interfaces()
	if err != nil {
		return err
	}

	for _, io := range netio {
		var isLoopback bool

		if len(s.Interfaces) != 0 {
			var found bool

			for _, name := range s.Interfaces {
				if name == io.Name {
					found = true
					break
				}
			}

			if !found {
				continue
			}
		} else if !s.skipChecks {
			iface, err := net.InterfaceByName(io.Name)
			if err != nil {
				continue
			}

			if iface.Flags&net.FlagLoopback == net.FlagLoopback {
				// continue // DO NOT skip loopback interface
				isLoopback = true
			}

			if iface.Flags&net.FlagUp == 0 {
				continue
			}
		}

		tags := map[string]string{
			"interface": io.Name,
		}
		if isLoopback {
			tags["nonemptyifLoopback"] = "nonempty"
		}
		if ip, ok := interfaceIPByName(interfaces, io.Name); ok {
			tags["ip"] = ip
		}

		fields := map[string]interface{}{
			"bytes_sent":   io.BytesSent,
			"bytes_recv":   io.BytesRecv,
			"packets_sent": io.PacketsSent,
			"packets_recv": io.PacketsRecv,
			"err_in":       io.Errin,
			"err_out":      io.Errout,
			"drop_in":      io.Dropin,
			"drop_out":     io.Dropout,
		}
		s.addDeltaFields(io, fields)
		acc.AddFields("net", fields, tags)
	}

	/*
		// Get system wide stats for different network protocols
		// (ignore these stats if the call fails)
		netprotos, _ := s.ps.NetProto()
		fields := make(map[string]interface{})
		for _, proto := range netprotos {
			for stat, value := range proto.Stats {
				name := fmt.Sprintf("%s_%s", strings.ToLower(proto.Protocol),
					strings.ToLower(stat))
				fields[name] = value
			}
		}
		tags := map[string]string{
			"interface": "all",
		}
		acc.AddFields("net", fields, tags)
	*/

	return nil
}

func init() {
	inputs.Add("net_ostent", func() telegraf.Input {
		return &NetIOStats{ps: &systemPS{}}
	})
}
