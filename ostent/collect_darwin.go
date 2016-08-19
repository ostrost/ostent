// +build darwin

package ostent

/*
#include <sys/param.h> // for MAXCOMLEN
*/
import "C"

import (
	"path/filepath"
	"unsafe"

	sigar "github.com/ostrost/gosigar"
)

// ProcName returns argv[0].
// pbiComm originating from ProcState may be chopped, in which case
// sigar.ProcExe gives absolute executable path and the basename of it is returned.
func ProcName(pid int, pbiComm string) string {
	if len(pbiComm)+1 < C.MAXCOMLEN {
		return pbiComm
	}
	exe := sigar.ProcExe{}
	if err := exe.Get(pid); err != nil {
		return pbiComm
	}
	return filepath.Base(exe.Name)
}

func IfaDropsOut(unsafe.Pointer) uint { return 0 }
