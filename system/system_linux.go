// +build linux

package system

import (
	"fmt"
	"os/exec"
	"strings"
)

// Distrib returns system distribution and version string.
func Distrib() (string, error) {
	// https://unix.stackexchange.com/q/6345
	std, err := exec.Command("lsb_release", "--id", "--release").CombinedOutput()
	if err != nil {
		return "", err
	}
	distrib, release := "", ""
	// strings.TrimRight(string(std), "\n\t ")
	for _, s := range strings.Split(string(std), "\n") {
		b := strings.Split(s, "\t")
		if len(b) == 2 {
			if b[0] == "Distributor ID:" {
				distrib = b[1]
				continue
			} else if b[0] == "Release:" {
				release = b[1]
				continue
			}
		}
	}
	if distrib == "" {
		return "", fmt.Errorf("Could not identify Distributor ID from lsb_release output")
	}
	if release == "" {
		return distrib, fmt.Errorf("Could not identify Release from lsb_release output")
	}
	return distrib + " " + release, nil
}

// ProcName returns argv[0].
// Actually, unless it's darwin, the procName itself is returned.
func ProcName(_ int, procName string) string {
	return procName // from /proc/_/stat, may be shortened
}
