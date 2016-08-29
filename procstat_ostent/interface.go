package procstat_ostent

import (
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/process"
)

type proc interface {
	IOCounters() (*process.IOCountersStat, error)
	MemoryInfo() (*process.MemoryInfoStat, error)
	Name() (string, error)
	Nice() (int32, error)
	NumCtxSwitches() (*process.NumCtxSwitchesStat, error)
	NumFDs() (int32, error)
	NumThreads() (int32, error)
	Percent(time.Duration) (float64, error)
	Times() (*cpu.TimesStat, error)
	Uids() ([]int32, error)
}
