package cpu

import (
	"fmt"
	"sort"

	"github.com/ostrost/ostent/client"
	"github.com/ostrost/ostent/format"
	"github.com/ostrost/ostent/registry"
	sigar "github.com/rzab/gosigar"
)

type CPUData struct {
	sigarList sigar.CpuList
	deltaList *sigar.CpuList
	cpuInfo   *CPUInfo // cached CPUInfo
	coreno    int      // number of cores, viable when .cpuInfo is non-nil
}

func NewCPUData() CPUData {
	cl := sigar.CpuList{}
	cl.Get()
	return CPUData{sigarList: cl}
}

func (cd CPUData) List() sigar.CpuList {
	if cd.deltaList != nil {
		return *cd.deltaList
	}
	return cd.sigarList
}

func (cd CPUData) SigarList() *sigar.CpuList {
	return &cd.sigarList
}

func (cd *CPUData) CalculateDelta(other []sigar.Cpu) {
	if len(other) == 0 {
		return
	}
	cores := cd.sigarList.List
	coreno := len(cores)
	if len(other) != coreno {
		return
	}
	deltaList := make([]sigar.Cpu, coreno)
	for i, sigarCpu := range cores {
		deltaList[i] = sigarCpu.Delta(other[i])
	}
	cd.deltaList = &sigar.CpuList{List: deltaList}
}

// CPUInfo type has a list of CoreInfo.
type CPUInfo struct {
	List    []CoreInfo     // TODO rename to Cores
	RawInfo []CoreInfo     `json:"-"`
	RawList *sigar.CpuList `json:"-"`
}

// CoreInfo type is a struct of core metrics.
type CoreInfo struct {
	N         string
	User      uint // percent without "%"
	Sys       uint // percent without "%"
	Idle      uint // percent without "%"
	UserClass string
	SysClass  string
	IdleClass string
	// UserSpark string
	// SysSpark  string
	// IdleSpark string
}

func add(sc *sigar.Cpu, other sigar.Cpu) {
	sc.User += other.User
	sc.Sys += other.Sys
	sc.Idle += other.Idle
}

// CPUInfo returns a CPUInfo which is a struct for templates.
func (cd *CPUData) CPUInfo(client client.Client) (*CPUInfo, int) {
	if cd.cpuInfo == nil {
		cd.cpuInfo, cd.coreno = cd.newcpuInfo()
	}
	cl := cd.cpuInfo.List
	cp := &CPUInfo{RawInfo: cl, RawList: cd.deltaList}
	if cd.coreno != 1 { // all but "all"
		cp.RawInfo = cl[1:len(cl)]
		if *client.ExpandCPU {
			cl = cl[1:len(cl)]
		}
	}
	if !*client.ExpandCPU && cd.coreno > client.Toprows-1 {
		cl = cl[:client.Toprows] // "collapsed" view, head of the list
	}
	cp.List = cl
	return cp, cd.coreno
}

// newcpuInfo produces a cpuInfo
func (cd *CPUData) newcpuInfo() (*CPUInfo, int) {
	list, coreno := cd.newCores()
	return &CPUInfo{List: list}, coreno
}

// newCores
func (cd *CPUData) newCores() ([]CoreInfo, int) {
	sum := sigar.Cpu{}

	sigarCores := cd.List().List // .deltaList || .sigarList
	coreno := len(sigarCores)

	cores := make([]CoreInfo, coreno)
	for i, core := range sigarCores {
		cores[i] = coreInfo(core, fmt.Sprintf("#%d", i))
		if coreno != 1 {
			add(&sum, core)
		}
	}
	if coreno == 1 {
		return cores, coreno
	}
	sort.Sort(_cores(cores))
	// TODO finish the implementation
	cores = append([]CoreInfo{
		coreInfo(sum, fmt.Sprintf("all %d", coreno)),
	}, cores...)
	return cores, coreno
}

func coreInfo(core sigar.Cpu, n string) CoreInfo {
	total := core.Total()
	user := format.Percent(core.User, total)
	sys := format.Percent(core.Sys, total)
	idle := format.Percent(core.Idle, total)
	if idle > 100-user-sys { // rounding happened
		idle = 100 - user - sys
	}
	return CoreInfo{
		N:         n,
		User:      user,
		Sys:       sys,
		Idle:      idle,
		UserClass: format.TextClassColorPercent(user),
		SysClass:  format.TextClassColorPercent(sys),
		IdleClass: format.TextClassColorPercent(100 - idle),
	}
}

type _cores []CoreInfo          // just for sorting
func (cs _cores) Len() int      { return len(cs) }
func (cs _cores) Swap(i, j int) { cs[i], cs[j] = cs[j], cs[i] }
func (cs _cores) Less(i, j int) bool {
	return (cs[j].User + cs[j].Sys) < (cs[i].User + cs[i].Sys)
}

func CollectCPU(reg registry.Registry, CH chan<- CPUData, prevcl *sigar.CpuList) {
	cd := NewCPUData()
	if prevcl != nil {
		cd.CalculateDelta(prevcl.List)
	}
	reg.UpdateCPU(cd.sigarList.List)
	CH <- cd
}

type Send struct {
	cpu   *sigar.Cpu
	total *uint64
}

func NewSend(ci CPUInfo, coreno int) Send {
	var cpu *sigar.Cpu
	if ci.RawList != nil {
		cpu = &ci.RawList.List[coreno]
	}
	return Send{cpu: cpu}
}

func (se Send) raw() sigar.Cpu {
	if se.cpu != nil {
		return *se.cpu
	}
	return sigar.Cpu{}
}

func (se *Send) fraction(value uint64) string {
	if se.cpu == nil {
		return "nan"
	}
	if se.total == nil {
		se.total = new(uint64)
		*se.total = se.calcTotal()
	}
	return fmt.Sprintf("%f", float64(value)/float64(*se.total))
}

func (se Send) calcTotal() uint64 {
	return CalcTotal(*se.cpu)
}
