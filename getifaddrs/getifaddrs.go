// +build linux freebsd darwin

// Package getifaddrs does getifaddrs(3) for Go.
package getifaddrs

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

char ADDR[INET_ADDRSTRLEN];
*/
import "C"
import "unsafe"

// IfAddr is a struct with interface info.
type IfAddr struct {
	IfaIP         string
	IfaName       string
	IfaBytesIn    uint
	IfaBytesOut   uint
	IfaPacketsIn  uint
	IfaPacketsOut uint
	IfaDropsIn    uint
	IfaDropsOut   *uint // nil in darwin
	IfaErrorsIn   uint
	IfaErrorsOut  uint
}

// GetName and other methods may be combined into an interface.
func (ia IfAddr) IP() string       { return ia.IfaIP }
func (ia IfAddr) Name() string     { return ia.IfaName }
func (ia IfAddr) BytesIn() uint    { return ia.IfaBytesIn }
func (ia IfAddr) BytesOut() uint   { return ia.IfaBytesOut }
func (ia IfAddr) DropsIn() uint    { return ia.IfaDropsIn }
func (ia IfAddr) DropsOut() *uint  { return ia.IfaDropsOut }
func (ia IfAddr) ErrorsIn() uint   { return ia.IfaErrorsIn }
func (ia IfAddr) ErrorsOut() uint  { return ia.IfaErrorsOut }
func (ia IfAddr) PacketsIn() uint  { return ia.IfaPacketsIn }
func (ia IfAddr) PacketsOut() uint { return ia.IfaPacketsOut }

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

// Getifaddrs returns a list of IfAddr. Unlike with getifaddrs(3) the
// IfAddr has merged link level and interface address data.
func Getifaddrs() ([]IfAddr, error) {
	ret := []IfAddr{}

	var ifaces *C.struct_ifaddrs
	if rc, err := C.getifaddrs(&ifaces); rc != 0 {
		return ret, err
	}
	defer C.freeifaddrs(ifaces)

	ips := make(map[string]string)
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
		ia := IfAddr{
			IfaName:       ifaName,
			IfaBytesIn:    uint(C.Ibytes(data)),
			IfaBytesOut:   uint(C.Obytes(data)),
			IfaPacketsIn:  uint(C.Ipackets(data)),
			IfaPacketsOut: uint(C.Opackets(data)),
			IfaDropsIn:    uint(C.Idrops(data)),
			IfaDropsOut:   IfaDropsOut(data), // may return nil
			IfaErrorsIn:   uint(C.Ierrors(data)),
			IfaErrorsOut:  uint(C.Oerrors(data)),
		}
		if ip, ok := ips[ifaName]; ok {
			ia.IfaIP = ip
		}
		ret = append(ret, ia)
	}
	for i, ifaddr := range ret {
		if ifaddr.IfaIP == "" {
			if ip, ok := ips[ifaddr.IfaName]; ok {
				ret[i].IfaIP = ip
			}
		}
	}
	return ret, nil
}
