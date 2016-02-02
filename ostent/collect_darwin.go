// +build darwin

package ostent

import "unsafe"

func IfaDropsOut(unsafe.Pointer) uint { return 0 }
