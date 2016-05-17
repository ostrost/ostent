package ostent

import (
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"

	sigar "github.com/ostrost/gosigar"

	"github.com/ostrost/ostent/format"
	"github.com/ostrost/ostent/system"
)

// Registry has updates with sigar values.
type Registry interface {
	UpdateIF(system.IfAddress)
	UpdateCPU(sigar.Cpu, []sigar.Cpu)
	UpdateLA(sigar.LoadAverage)
	UpdateSwap(sigar.Swap)
	UpdateRAM(sigar.Mem, uint64, uint64)
	UpdateDF(sigar.FileSystem, sigar.FileSystemUsage)
}

// Collector is collection interface.
type Collector interface {
	GetHN() (string, error)
	HN(S2SRegistry, *sync.WaitGroup)
	UP(S2SRegistry, *sync.WaitGroup)
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

func (m Machine) UP(sreg S2SRegistry, wg *sync.WaitGroup) {
	// m is unused
	uptime := sigar.Uptime{}
	if err := uptime.Get(); err == nil {
		sreg.SetString("uptime", format.Uptime(uptime.Length))
	}
	wg.Done()
}

func (m Machine) LA(reg Registry, wg *sync.WaitGroup) {
	// m is unused
	la := sigar.LoadAverage{}
	if err := la.Get(); err == nil {
		reg.UpdateLA(la)
	}
	wg.Done()
}

func _getmem(kind string, in sigar.Swap) system.Memory {
	total, approxtotal, _ := format.HumanBandback(in.Total)
	used, approxused, _ := format.HumanBandback(in.Used)

	return system.Memory{
		Kind:   kind,
		Total:  total,
		Used:   used,
		Free:   format.HumanB(in.Free),
		UsePct: format.Percent(approxused, approxtotal),
	}
}

func (m Machine) RAM(reg Registry, wg *sync.WaitGroup) {
	// m is unused
	got := sigar.Mem{}
	extra1, extra2, _ := sigar.GetExtra(&got)
	reg.UpdateRAM(got, extra1, extra2)
	wg.Done()

	// inactive := got.ActualFree - got.Free // == got.Used - got.ActualUsed // "kern"
	// _ = inactive

	// Used = .Total - .Free
	// | Free |           Used +%         | Total
	// | Free | Inactive | Active | Wired | Total

	// TODO active := vm_data.active_count << 12 (pagesize)
	// TODO wired  := vm_data.wire_count   << 12 (pagesoze)
}

func (m Machine) Swap(reg Registry, wg *sync.WaitGroup) {
	// m is unused
	got := sigar.Swap{}
	if err := got.Get(); err == nil {
		reg.UpdateSwap(got)
	}
	wg.Done()
}

func (m Machine) DF(reg Registry, wg *sync.WaitGroup) {
	// m is unused
	fls := sigar.FileSystemList{}
	if err := fls.Get(); err != nil {
		wg.Done()
		return
	}

	// devnames := map[string]bool{}
	dirnames := map[string]bool{}

	for _, fs := range fls.List {

		usage := sigar.FileSystemUsage{}
		if err := usage.Get(fs.DirName); err != nil {
			continue
		}

		if !strings.HasPrefix(fs.DevName, "/") {
			continue
		}
		// if _, ok := devnames[fs.DevName]; ok
		if _, ok := dirnames[fs.DirName]; ok {
			continue
		}
		// devnames[fs.DevName] = true
		dirnames[fs.DirName] = true

		reg.UpdateDF(fs, usage)
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
	all, list := sigar.Cpu{}, sigar.CpuList{}
	err1 := all.Get()
	err2 := list.Get()
	if err1 == nil && err2 == nil {
		reg.UpdateCPU(all, list.List)
	}
	wg.Done()
}
