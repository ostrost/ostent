// +build !darwin

package ostent

import "github.com/shirou/gopsutil/host"

func hostPlatformVersion() (string, string, error) {
	platform, _, version, err := host.PlatformInformation()
	return platform, version, err
}
