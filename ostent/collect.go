package ostent

import (
	"regexp"
	"runtime"
	"sync"

	sigar "github.com/ostrost/gosigar"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"

	"github.com/ostrost/ostent/system"
)

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

func (ir *IndexRegistry) collectLA(wg *sync.WaitGroup) {
	if stat, err := load.Avg(); err == nil {
		ir.UpdateLA(*stat)
	}
	wg.Done()
}

func (ir *IndexRegistry) collectMEM(wg *sync.WaitGroup) {
	var (
		ram, err1  = mem.VirtualMemory()
		swap, err2 = mem.SwapMemory()
	)
	if err1 == nil && err2 == nil {
		ir.UpdateMEM(ram, swap)
	}
	wg.Done()
}

func (ir *IndexRegistry) collectDF(wg *sync.WaitGroup) {
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
		ir.UpdateDF(part, usage)
	}
	wg.Done()
}

func collectPS(CH chan<- PSSlice) {
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

func (ir *IndexRegistry) collectCPU(wg *sync.WaitGroup) {
	var (
		aggs, err1 = cpu.Times(false)
		list, err2 = cpu.Times(true)
	)
	if err1 == nil && err2 == nil {
		ir.UpdateCPU(aggs[0], list)
	}
	wg.Done()
}
