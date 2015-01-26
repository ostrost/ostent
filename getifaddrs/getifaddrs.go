/*
Package getifaddrs does getifaddrs(3) for Go.
*/
package getifaddrs

// +build unix

/*
#include <sys/socket.h>
#include <net/if.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <ifaddrs.h>

#ifndef AF_LINK
#define AF_LINK AF_PACKET
#endif

#ifndef __linux__ // NOT LINUX
u_int32_t Ibytes(void *data) { return ((struct if_data *)data)->ifi_ibytes; }
u_int32_t Obytes(void *data) { return ((struct if_data *)data)->ifi_obytes; }

u_int32_t Ipackets(void *data) { return ((struct if_data *)data)->ifi_ipackets; }
u_int32_t Opackets(void *data) { return ((struct if_data *)data)->ifi_opackets; }

u_int32_t Ierrors(void *data) { return ((struct if_data *)data)->ifi_ierrors; }
u_int32_t Oerrors(void *data) { return ((struct if_data *)data)->ifi_oerrors; }

#else
#include <linux/if_link.h>
u_int32_t Ibytes(void *data) { return ((struct rtnl_link_stats *)data)->rx_bytes; }
u_int32_t Obytes(void *data) { return ((struct rtnl_link_stats *)data)->tx_bytes; }

u_int32_t Ipackets(void *data) { return ((struct rtnl_link_stats *)data)->rx_packets; }
u_int32_t Opackets(void *data) { return ((struct rtnl_link_stats *)data)->tx_packets; }

u_int32_t Ierrors(void *data) { return ((struct rtnl_link_stats *)data)->rx_errors; }
u_int32_t Oerrors(void *data) { return ((struct rtnl_link_stats *)data)->tx_errors; }
#endif

char ADDR[INET_ADDRSTRLEN];
*/
import "C"
import "unsafe"

// IfData is a struct with interface info.
type IfData struct {
	IP         string
	Name       string
	InBytes    uint
	OutBytes   uint
	InPackets  uint
	OutPackets uint
	InErrors   uint
	OutErrors  uint
}

func ntop(fi *C.struct_ifaddrs) (string, bool) {
	if fi.ifa_addr == nil {
		return "", false
	}
	if fi.ifa_addr.sa_family != C.AF_INET {
		return "", false
	}
	saIn := (*C.struct_sockaddr_in)(unsafe.Pointer(fi.ifa_addr))
	if nil == C.inet_ntop(
		C.int(fi.ifa_addr.sa_family), // C.AF_INET,
		unsafe.Pointer(&saIn.sin_addr),
		&C.ADDR[0],
		C.socklen_t(unsafe.Sizeof(C.ADDR))) {
		return "", false
	}
	return C.GoString((*C.char)(unsafe.Pointer(&C.ADDR))), true
}

// Getifaddrs returns a list of IfData. Unlike with getifaddrs(3) the
// IfData has merged link level and interface address data.
func Getifaddrs() ([]IfData, error) {
	var ifaces *C.struct_ifaddrs
	if rc, err := C.getifaddrs(&ifaces); rc != 0 {
		return []IfData{}, err
	}
	defer C.freeifaddrs(ifaces)

	ips := make(map[string]string)
	ifs := []IfData{}

	for fi := ifaces; fi != nil; fi = fi.ifa_next {
		if fi.ifa_addr == nil {
			continue
		}

		ifaName := C.GoString(fi.ifa_name)
		if ip, ok := ntop(fi); ok {
			ips[ifaName] = ip
			continue // fi.ifa_addr.sa_family == C.AF_INET
		}

		if fi.ifa_addr.sa_family != C.AF_LINK {
			continue
		}

		data := fi.ifa_data
		it := IfData{
			Name:       ifaName,
			InBytes:    uint(C.Ibytes(data)),
			OutBytes:   uint(C.Obytes(data)),
			InPackets:  uint(C.Ipackets(data)),
			OutPackets: uint(C.Opackets(data)),
			InErrors:   uint(C.Ierrors(data)),
			OutErrors:  uint(C.Oerrors(data)),
		}
		if ip, ok := ips[ifaName]; ok {
			it.IP = ip
		}
		ifs = append(ifs, it)
	}
	for i, ifdata := range ifs {
		if ifdata.IP != "" {
			continue
		}
		if ip, ok := ips[ifdata.Name]; ok {
			ifs[i].IP = ip
		}
	}
	return ifs, nil
}
