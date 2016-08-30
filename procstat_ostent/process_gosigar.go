//+build none

package procstat_ostent

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/process"

	sigar "github.com/.../gosigar"
)

func newProcess(pid int32) (proc, error) { return &sigarProcess{pid: int(pid)}, nil }

type sigarProcess struct {
	pid int

	state    sigar.ProcState
	stateErr error
	stateGot bool

	ptime    sigar.ProcTime
	ptimeErr error
	ptimeGot bool
}

func (sp *sigarProcess) getState() error {
	if !sp.stateGot && sp.stateErr == nil {
		sp.state = sigar.ProcState{}
		sp.stateErr = sp.state.Get(sp.pid)
		sp.stateGot = true
	}
	return sp.stateErr
}

func (sp *sigarProcess) getTime() error {
	if !sp.ptimeGot && sp.ptimeErr == nil {
		sp.ptime = sigar.ProcTime{}
		sp.ptimeErr = sp.ptime.Get(sp.pid)
		sp.ptimeGot = true
	}
	return sp.ptimeErr
}

func (sp *sigarProcess) MemoryInfo() (*process.MemoryInfoStat, error) {
	mem := sigar.ProcMem{}
	err := mem.Get(sp.pid)
	return &process.MemoryInfoStat{
		RSS: mem.Resident,
		VMS: mem.Size,
		// .Swap is omitted even by gopsutil
	}, err
}

func (sp *sigarProcess) Name() (string, error) {
	// .Name acquiring resets sp.{state,ptime}*
	sp.stateGot, sp.stateErr = false, nil
	sp.ptimeGot, sp.ptimeErr = false, nil
	err := sp.getState()
	return sp.state.Name, err
}

func (sp *sigarProcess) Nice() (int32, error) {
	err := sp.getState()
	return int32(sp.state.Nice), err
}

func (sp *sigarProcess) Times() (*cpu.TimesStat, error) {
	err := sp.getTime()
	return &cpu.TimesStat{
		CPU:    "cpu",
		User:   float64(sp.ptime.User / 1000),
		System: float64(sp.ptime.Sys / 1000),
	}, err
}

func (sp *sigarProcess) Uids() ([]int32, error) {
	err := sp.getState()
	return []int32{int32(sp.state.Uid)}, err
}
