package ostent

import (
	"strings"

	"github.com/shirou/gopsutil/host"
)

/*
TODO Check for platform (&version?) naming
- gopsutil/host in !linux executes
  - `uname -s` for platform
  - `uname -r` for version
- ostent used to
  - in freebsd: `sysctl kern.version | cut -f1,2`
  - in darwin: "Mac OS X " + `sw_vers -productVersion`
    other darwin commands for version retrieval:
    - sw_vers
    - sw_vers -productVersion
    - system_profiler SPSoftwareDataType
    - defaults read loginwindow SystemVersionStampAsString
    - defaults read /System/Library/CoreServices/SystemVersion ProductVersion
*/

// Distrib is to return distribution identifier string with version.
func Distrib() (string, error) {
	platform, _, version, err := host.PlatformInformation()
	if err != nil {
		return "", err
	}
	if platform == "" {
		return "Docker", nil // Docker is a good guess.
	}
	distid := distribID(platform)
	if version != "" {
		return distid + " " + version, nil
	}
	return distid, nil
}

// distribID is to convert gopsutil platform identifier back to LSB Distributor ID form.
func distribID(platform string) string {
	switch platform {
	case "redhat":
		return "RedHat"
	case "linuxmint":
		return "LinuxMint"
	case "scientific":
		return "ScientificSL"
	case "xenserver":
		return "XenServer"
	case "centos":
		return "CentOS"
	case "cloudlinux":
		return "CloudLinux"
	case "opensuse":
		return "OpenSUSE"
	case "suse":
		return "SUSE"
	}
	return strings.Title(platform)
}
