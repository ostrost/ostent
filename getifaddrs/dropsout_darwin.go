// +build darwin

package getifaddrs

import "unsafe"

func IfaDropsOut(unsafe.Pointer) *uint { return nil }
