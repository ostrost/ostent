// +build freebsd

package ostent

import (
	"strings"
	unix "syscall"

	// "golang.org/x/sys/unix"
)

func getDistrib() string {
	if uname, err := unix.Sysctl("kern.version"); err == nil {
		split := strings.Split(uname, " ")
		if len(split) > 1 {
			return strings.Join(split[:2], " ")
		}
	}
	return "FreeBSD"
}

func procname(_ int, procName string) string {
	return procName
}
