// +build darwin

package ostent

// #include <sys/param.h>
import "C"
import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	sigar "github.com/rzab/gosigar"
)

type getsUptime interface {
	Get() error
	Get32() error
}

var getUptime = func(gu getsUptime) error {
	return gu.Get()
}

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
	ver := strings.TrimRight(string(std), "\n\t ")
	if ver == "10.10" {
		getUptime = func(gu getsUptime) error {
			return gu.Get32()
		}
	}
	return "Mac OS X " + ver
}

// ProcState returns chopped proc name, in which case
// get the ProcExe and return basename of executable path
func procname(pid int, pbi_comm string) string {
	if len(pbi_comm)+1 < C.MAXCOMLEN {
		return pbi_comm
	}
	exe := sigar.ProcExe{}
	if err := exe.Get(pid); err != nil {
		return pbi_comm
	}
	return filepath.Base(exe.Name)
}
