// +build !linux

package system_ostent

/*
#include <sys/socket.h>
#include <net/if.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <ifaddrs.h>

#ifndef AF_LINK
#define AF_LINK AF_PACKET
#endif

#ifdef __linux__
#include <linux/if_link.h>
u_int32_t Ibytes(void *data) { return ((struct rtnl_link_stats *)data)->rx_bytes; }
u_int32_t Obytes(void *data) { return ((struct rtnl_link_stats *)data)->tx_bytes; }

u_int32_t Ipackets(void *data) { return ((struct rtnl_link_stats *)data)->rx_packets; }
u_int32_t Opackets(void *data) { return ((struct rtnl_link_stats *)data)->tx_packets; }

u_int32_t Ierrors(void *data) { return ((struct rtnl_link_stats *)data)->rx_errors; }
u_int32_t Oerrors(void *data) { return ((struct rtnl_link_stats *)data)->tx_errors; }

u_int32_t Idrops(void *data) { return ((struct rtnl_link_stats *)data)->rx_dropped; }

#else // freebsd, darwin
u_int32_t Ibytes(void *data) { return ((struct if_data *)data)->ifi_ibytes; }
u_int32_t Obytes(void *data) { return ((struct if_data *)data)->ifi_obytes; }

u_int32_t Ipackets(void *data) { return ((struct if_data *)data)->ifi_ipackets; }
u_int32_t Opackets(void *data) { return ((struct if_data *)data)->ifi_opackets; }

u_int32_t Ierrors(void *data) { return ((struct if_data *)data)->ifi_ierrors; }
u_int32_t Oerrors(void *data) { return ((struct if_data *)data)->ifi_oerrors; }

u_int32_t Idrops(void *data) { return ((struct if_data *)data)->ifi_iqdrops; }
#endif
*/
import "C"
import "github.com/shirou/gopsutil/net"

// iocounters is to call getifaddrs(3).
func iocounters(bool) ([]net.IOCountersStat, error) { return getifaddrs() }

func getifaddrs() ([]net.IOCountersStat, error) {
	list := []net.IOCountersStat{}

	var ifaddrs *C.struct_ifaddrs
	if rc, err := C.getifaddrs(&ifaddrs); rc != 0 {
		return list, err
	}
	defer C.freeifaddrs(ifaddrs)

	for fi := ifaddrs; fi != nil; fi = fi.ifa_next {
		if fi.ifa_addr == nil || fi.ifa_addr.sa_family != C.AF_LINK {
			continue
		}
		data := fi.ifa_data
		list = append(list, net.IOCountersStat{
			Name: C.GoString(fi.ifa_name),

			BytesRecv:   uint64(C.Ibytes(data)),
			BytesSent:   uint64(C.Obytes(data)),
			PacketsRecv: uint64(C.Ipackets(data)),
			PacketsSent: uint64(C.Opackets(data)),
			Errin:       uint64(C.Ierrors(data)),
			Errout:      uint64(C.Oerrors(data)),
			Dropin:      uint64(C.Idrops(data)),
			Dropout:     ifaDropsOut(data),
		})
	}
	return list, nil
}
