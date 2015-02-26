// +build freebsd

package system

import (
	"strings"
	unix "syscall"

	// "golang.org/x/sys/unix"
)

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

// ProcName returns argv[0].
// Actually, unless it's darwin, the procName itself is returned.
func ProcName(_ int, procName string) string {
	return procName
}
