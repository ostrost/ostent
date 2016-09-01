package procstat_ostent

import (
	"os"
	"strings"
	"syscall"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/process"
)

type proc interface {
	MemoryInfo() (*process.MemoryInfoStat, error)
	Name() (string, error)
	Nice() (int32, error)
	Prio() (int32, error)
	Times() (*cpu.TimesStat, error)
	Uids() ([]int32, error)
}

func isNotExist(err error) bool {
	if err == nil {
		return false
	}
	if os.IsNotExist(err) { // gopsutil error
		return true
	}
	switch pe := err.(type) {
	case syscall.Errno:
		return pe == syscall.ESRCH // gosigar linux error
	}
	if strings.Contains(err.Error(), "no such") { // gosigar generic? error
		return true
	}
	return false
}
