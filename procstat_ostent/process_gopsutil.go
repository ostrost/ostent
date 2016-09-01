// +build freebsd,darwin

package procstat_ostent

import "github.com/shirou/gopsutil/process"

func newProcess(pid int32) (proc, error) {
	np, err := process.NewProcess(pid)
	return &psproc{np}, err
}

type psproc struct{ *process.Process }

func (pp *psproc) Prio() (int32, error) { return 0, nil }
