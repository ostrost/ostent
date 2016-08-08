// +build darwin

package internal

import (
	"unsafe"
)

func dropsOut(unsafe.Pointer) uint64 { return 0 }
