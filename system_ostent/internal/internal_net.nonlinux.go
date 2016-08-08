// +build !linux

package internal

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
u_int32_t bytesRecv(void *data) { return ((struct rtnl_link_stats *)data)->rx_bytes; }
u_int32_t bytesSent(void *data) { return ((struct rtnl_link_stats *)data)->tx_bytes; }

u_int32_t packetsRecv(void *data) { return ((struct rtnl_link_stats *)data)->rx_packets; }
u_int32_t packetsSent(void *data) { return ((struct rtnl_link_stats *)data)->tx_packets; }

u_int32_t errorsIn(void *data) { return ((struct rtnl_link_stats *)data)->rx_errors; }
u_int32_t errorsOut(void *data) { return ((struct rtnl_link_stats *)data)->tx_errors; }

u_int32_t dropsIn(void *data) { return ((struct rtnl_link_stats *)data)->rx_dropped; }

#else // freebsd, darwin
u_int32_t bytesRecv(void *data) { return ((struct if_data *)data)->ifi_ibytes; }
u_int32_t bytesSent(void *data) { return ((struct if_data *)data)->ifi_obytes; }

u_int32_t packetsRecv(void *data) { return ((struct if_data *)data)->ifi_ipackets; }
u_int32_t packetsSent(void *data) { return ((struct if_data *)data)->ifi_opackets; }

u_int32_t errorsIn(void *data) { return ((struct if_data *)data)->ifi_ierrors; }
u_int32_t errorsOut(void *data) { return ((struct if_data *)data)->ifi_oerrors; }

u_int32_t dropsIn(void *data) { return ((struct if_data *)data)->ifi_iqdrops; }
#endif
*/
import "C"
import "github.com/shirou/gopsutil/net"

// IOCounters is to call getifaddrs(3).
func IOCounters(bool) ([]net.IOCountersStat, error) { return getifaddrs() }

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

			BytesRecv:   uint64(C.bytesRecv(data)),
			BytesSent:   uint64(C.bytesSent(data)),
			PacketsRecv: uint64(C.packetsRecv(data)),
			PacketsSent: uint64(C.packetsSent(data)),
			Errin:       uint64(C.errorsIn(data)),
			Errout:      uint64(C.errorsOut(data)),
			Dropin:      uint64(C.dropsIn(data)),
			Dropout:     dropsOut(data),
		})
	}
	return list, nil
}
