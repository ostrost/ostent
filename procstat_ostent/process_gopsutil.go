package procstat_ostent

import "github.com/shirou/gopsutil/process"

func newProcess(pid int32) (proc, error) { return process.NewProcess(pid) }
