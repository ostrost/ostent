// +build darwin

package internal

import (
	"unsafe"
)

func ifaDropsOut(unsafe.Pointer) uint64 { return 0 }
