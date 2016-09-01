// +build linux

package procstat_ostent

import (
	"bytes"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/process"

	sigar "github.com/cloudfoundry/gosigar"
)

func newProcess(pid int32) (proc, error) { return &sigarProcess{pid: int(pid)}, nil }

type sigarProcess struct {
	pid int

	state    sigar.ProcState // plain struct
	stateErr error
	stateGot bool

	ptime    sigar.ProcTime // plain struct
	ptimeErr error
	ptimeGot bool

	uid    uid // plain value
	uidErr error
	uidGot bool
}

func (sp *sigarProcess) getState() error {
	if !sp.stateGot && sp.stateErr == nil {
		sp.state = sigar.ProcState{}
		sp.stateErr = sp.state.Get(sp.pid)
		sp.stateGot = true
	}
	return sp.stateErr
}

func (sp *sigarProcess) getPtime() error {
	if !sp.ptimeGot && sp.ptimeErr == nil {
		sp.ptime = sigar.ProcTime{}
		sp.ptimeErr = sp.ptime.Get(sp.pid)
		sp.ptimeGot = true
	}
	return sp.ptimeErr
}

func (sp *sigarProcess) getUid() error {
	if !sp.uidGot && sp.uidErr == nil {
		sp.uidErr = sp.uid.Get(sp.pid)
		sp.uidGot = true
	}
	return sp.uidErr
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

func (sp *sigarProcess) Prio() (int32, error) {
	err := sp.getState()
	return int32(sp.state.Priority), err
}

func (sp *sigarProcess) Times() (*cpu.TimesStat, error) {
	err := sp.getPtime()
	return &cpu.TimesStat{
		CPU:    "cpu",
		User:   float64(sp.ptime.User / 1000),
		System: float64(sp.ptime.Sys / 1000),
	}, err
}

func (sp *sigarProcess) Uids() ([]int32, error) {
	err := sp.getUid()
	return []int32{int32(sp.uid)}, err
}

type uid uint32

func (self *uid) Get(pid int) error {
	procdir := "/proc" // a la gopsutil
	if v := os.Getenv("HOST_PROC"); v != "" {
		procdir = v
	}

	status, err := ioutil.ReadFile(procdir + "/" + strconv.Itoa(pid) + "/status")
	if err != nil {
		if perr, ok := err.(*os.PathError); ok && perr.Err == syscall.ENOENT {
			return syscall.ESRCH // sigar type of error
		}
		return err
	}

	for _, line := range bytes.Split(status, []byte("\n")) {
		fields := bytes.Split(line, []byte(":"))
		if !bytes.Equal(fields[0], []byte("Uid")) || len(fields) < 2 {
			continue
		}

		if v, err := strconv.ParseUint(string(bytes.Fields(bytes.TrimLeft(
			fields[1], " "))[0]), 10, 32); err == nil { // err is gone
			*self = uid(v) // cast uint64 to uint32 to uid
			break
		}
	}
	return nil
}
