// +build !linux,!darwin

package internal

/*
#include <stdio.h>
#include <sys/socket.h>
#include <net/if.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <ifaddrs.h>

#ifdef __linux__
#include <linux/if_link.h>
u_int32_t cdropsOut(void *data) { return ((struct rtnl_link_stats *)data)->tx_dropped; }

#else // freebsd
u_int32_t cdropsOut(void *data) { return ((struct if_data *)data)->ifi_oqdrops; }
#endif
*/
// #cgo CFLAGS: -D_IFI_OQDROPS
import "C"
import "unsafe"

func dropsOut(data unsafe.Pointer) uint64 { return uint64(C.cdropsOut(data)) }
