// +build darwin

package system_ostent

import (
	"unsafe"
)

func ifaDropsOut(unsafe.Pointer) uint64 { return 0 }
