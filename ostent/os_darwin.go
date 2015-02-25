// +build darwin

package ostent

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	sigar "github.com/rzab/gosigar"
)

func getDistrib() string {
	/* various cli to show the mac version
	sw_vers
	sw_vers -productVersion
	system_profiler SPSoftwareDataType
	defaults read loginwindow SystemVersionStampAsString
	defaults read /System/Library/CoreServices/SystemVersion ProductVersion
	*/
	std, err := exec.Command("sw_vers", "-productVersion").CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "sw_vers: %s\n", err)
		return ""
	}
	return "Mac OS X " + strings.TrimRight(string(std), "\n\t ")
}

// ProcState returns chopped proc name, in which case
// get the ProcExe and return basename of executable path
func procname(pid int, pbi_comm string) string {
	if len(pbi_comm)+1 < sigar.CMAXCOMLEN {
		return pbi_comm
	}
	exe := sigar.ProcExe{}
	if err := exe.Get(pid); err != nil {
		return pbi_comm
	}
	return filepath.Base(exe.Name)
}
