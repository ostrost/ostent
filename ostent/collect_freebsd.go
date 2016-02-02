// +build freebsd

package ostent

import (
	"strings"
	unix "syscall"

	// "golang.org/x/sys/unix"
)

// ProcName returns procName back.
func ProcName(_ int, procName string) string { return procName }

// Distrib returns system distribution and version string.
func Distrib() (string, error) {
	uname, err := unix.Sysctl("kern.version")
	if err != nil {
		return "", err
	}
	split := strings.Split(uname, " ")
	if len(split) > 1 {
		return strings.Join(split[:2], " "), nil
	}
	return "FreeBSD", nil
}
