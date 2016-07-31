package ostent

import (
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"

	sigar "github.com/ostrost/gosigar"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"

	"github.com/ostrost/ostent/system"
)

// Registry has updates with gopsutil stats.
type Registry interface {
	UpdateIF(system.IfAddress)
	UpdateCPU(cpu.TimesStat, []cpu.TimesStat)
	UpdateLA(load.AvgStat)
	UpdateSwap(*mem.SwapMemoryStat)
	UpdateRAM(*mem.VirtualMemoryStat)
	UpdateDF(disk.PartitionStat, *disk.UsageStat)
}

// Collector is collection interface.
type Collector interface {
	GetHN() (string, error)
	HN(S2SRegistry, *sync.WaitGroup)
	LA(Registry, *sync.WaitGroup)
	RAM(Registry, *sync.WaitGroup)
	Swap(Registry, *sync.WaitGroup)
	IF(Registry, *sync.WaitGroup)
	PS(chan<- PSSlice)
	DF(Registry, *sync.WaitGroup)
	CPU(Registry, *sync.WaitGroup)
}

// These are regexps to match network interfaces.
var (
	RXlo      = regexp.MustCompile(`^lo\d*$`)
	RXvbr     = regexp.MustCompile(`^virbr\d+$`)
	RXvbrnic  = regexp.MustCompile(`^virbr\d+-nic$`)
	RXbridge  = regexp.MustCompile(`^bridge\d+$`)
	RXvboxnet = regexp.MustCompile(`^vboxnet\d+$`)
	RXfw      = regexp.MustCompile(`^fw\d+$`)
	RXgif     = regexp.MustCompile(`^gif\d+$`)
	RXstf     = regexp.MustCompile(`^stf\d+$`)
	RXwdl     = regexp.MustCompile(`^awdl\d+$`)
	RXairdrop = regexp.MustCompile(`^p2p\d+$`)
)

// HardwareIF returns false for known virtual/software network interface name.
func HardwareIF(name string) bool {
	if RXvbr.MatchString(name) ||
		RXvbrnic.MatchString(name) ||
		RXbridge.MatchString(name) ||
		RXvboxnet.MatchString(name) {
		return false
	}
	if runtime.GOOS == "darwin" {
		if RXfw.MatchString(name) ||
			RXgif.MatchString(name) ||
			RXstf.MatchString(name) ||
			RXwdl.MatchString(name) ||
			RXairdrop.MatchString(name) {
			return false
		}
	}
	return true
}

// Machine implements Collector by collecting the maching metrics.
type Machine struct{}

func (m Machine) GetHN() (string, error) {
	// m is unused
	return GetHN()
}

func GetHN() (string, error) {
	hostname, err := os.Hostname()
	if err == nil {
		hostname = strings.Split(hostname, ".")[0]
	}
	return hostname, err
}

func (m Machine) HN(sreg S2SRegistry, wg *sync.WaitGroup) {
	if hostname, err := m.GetHN(); err == nil {
		sreg.SetString("hostname", hostname)
	}
	wg.Done()
}

func (m Machine) LA(reg Registry, wg *sync.WaitGroup) {
	// m is unused
	if stat, err := load.Avg(); err == nil {
		reg.UpdateLA(*stat)
	}
	wg.Done()
}

func (m Machine) RAM(reg Registry, wg *sync.WaitGroup) {
	// m is unused
	if stat, err := mem.VirtualMemory(); err == nil {
		reg.UpdateRAM(stat)
	}
	wg.Done()
}

func (m Machine) Swap(reg Registry, wg *sync.WaitGroup) {
	// m is unused
	if stat, err := mem.SwapMemory(); err == nil {
		reg.UpdateSwap(stat)
	}
	wg.Done()
}

func (m Machine) DF(reg Registry, wg *sync.WaitGroup) {
	// m is unused
	parts, err := disk.Partitions(false)
	if err != nil {
		wg.Done()
		return
	}

	devices := map[string]struct{}{}
	for _, part := range parts {
		usage, err := disk.Usage(part.Mountpoint)
		if err != nil {
			continue
		}
		if _, ok := devices[part.Device]; ok {
			continue
		}
		devices[part.Device] = struct{}{}
		reg.UpdateDF(part, usage)
	}
	wg.Done()
}

func (m Machine) PS(CH chan<- PSSlice) {
	// m is unused
	var pss PSSlice
	pls := sigar.ProcList{}
	if err := pls.Get(); err != nil {
		CH <- nil
		return
	}

	for _, pid := range pls.List {

		state := sigar.ProcState{}
		// args := sigar.ProcArgs{}
		time := sigar.ProcTime{}
		mem := sigar.ProcMem{}

		if err := state.Get(pid); err != nil {
			continue
		}
		// if err :=  args.Get(pid); err != nil { continue }
		if err := time.Get(pid); err != nil {
			continue
		}
		if err := mem.Get(pid); err != nil {
			continue
		}

		pss = append(pss, &system.PSInfo{
			PID:      uint(pid),
			Priority: state.Priority,
			Nice:     state.Nice,
			Time:     time.Total,
			Name:     ProcName(pid, state.Name),
			// Name:  strings.Join(append([]string{ProcName(pid, state.Name)}, args.List[1:]...), " "),
			UID:      state.Uid,
			Size:     mem.Size,
			Resident: mem.Resident,
		})
	}
	CH <- pss
}

func (m Machine) CPU(reg Registry, wg *sync.WaitGroup) {
	// m is unused
	var (
		aggs, err1 = cpu.Times(false)
		list, err2 = cpu.Times(true)
	)
	if err1 == nil && err2 == nil {
		reg.UpdateCPU(aggs[0], list)
	}
	wg.Done()
}
