package ostent

import "strings"

// Distrib is to return distribution identifier string with version.
func Distrib() (string, error) {
	platform, version, err := hostPlatformVersion()
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
