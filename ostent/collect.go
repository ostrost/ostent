package ostent

import (
	"html/template"
	"log"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/ostrost/ostent/format"
	"github.com/ostrost/ostent/getifaddrs"
	"github.com/ostrost/ostent/registry"
	"github.com/ostrost/ostent/system"
	"github.com/ostrost/ostent/system/operating"
	"github.com/ostrost/ostent/templateutil"
	sigar "github.com/rzab/gosigar"
)

// Collector is collection interface.
type Collector interface {
	GetHostname() (string, error)
	Hostname(S2SRegistry, *sync.WaitGroup)
	Uptime(S2SRegistry, *sync.WaitGroup)
	LA(registry.Registry, *sync.WaitGroup)
	RAM(registry.Registry, *sync.WaitGroup)
	Swap(registry.Registry, *sync.WaitGroup)
	Interfaces(registry.Registry, S2SRegistry, *sync.WaitGroup)
	Procs(chan<- operating.MetricProcSlice)
	Disks(registry.Registry, *sync.WaitGroup)
	CPU(registry.Registry, *sync.WaitGroup)
}

var (
	// RXlo is a regexp to match loopback network interface
	RXlo = regexp.MustCompile("^lo\\d*$")

	// RXfw is a regexp to match non-hardware network interface
	RXfw = regexp.MustCompile("^fw\\d+$")
	// RXgif is a regexp to match non-hardware network interface
	RXgif = regexp.MustCompile("^gif\\d+$")
	// RXstf is a regexp to match non-hardware network interface
	RXstf = regexp.MustCompile("^stf\\d+$")
	// RXwdl is a regexp to match non-hardware network interface
	RXwdl = regexp.MustCompile("^awdl\\d+$")
	// RXbridge is a regexp to match non-hardware network interface
	RXbridge = regexp.MustCompile("^bridge\\d+$")
	// RXvboxnet is a regexp to match non-hardware network interface
	RXvboxnet = regexp.MustCompile("^vboxnet\\d+$")
	// RXairdrop is a regexp to match non-hardware network interface
	RXairdrop = regexp.MustCompile("^p2p\\d+$")
)

// HardwareInterface returns false for known virtual/software network interface name.
func HardwareInterface(name string) bool {
	if RXbridge.MatchString(name) ||
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

// ApplyperInterface calls apply for each found hardware interface.
func (m *Machine) ApplyperInterface(apply func(getifaddrs.IfData) bool) error {
	// m is unused
	gotifaddrs, err := getifaddrs.Getifaddrs()
	if err != nil {
		return err
	}
	for _, ifdata := range gotifaddrs {
		if !HardwareInterface(ifdata.Name) {
			continue
		}
		if !apply(ifdata) {
			break
		}
	}
	return nil
}

type FoundIP struct {
	string
}

func (fip *FoundIP) Next(ifdata getifaddrs.IfData) bool {
	if fip.string != "" {
		return false
	}
	if !RXlo.MatchString(ifdata.Name) { // non-loopback
		fip.string = ifdata.IP
		return false
	}
	return true
}

// Interfaces registers the interfaces with the reg and send first non-loopback IP to the chan
func (m *Machine) Interfaces(reg registry.Registry, sreg S2SRegistry, wg *sync.WaitGroup) {
	fip := FoundIP{}
	m.ApplyperInterface(func(ifdata getifaddrs.IfData) bool {
		fip.Next(ifdata)
		if ifdata.InBytes == 0 &&
			ifdata.OutBytes == 0 &&
			ifdata.InPackets == 0 &&
			ifdata.OutPackets == 0 &&
			ifdata.InErrors == 0 &&
			ifdata.OutErrors == 0 {
			// nothing
		} else {
			reg.UpdateIFdata(ifdata)
		}
		return true
	})
	sreg.SetString("ip", fip.string)
	wg.Done()
}

func (m *Machine) GetHostname() (string, error) {
	// m is unused
	hostname, err := os.Hostname()
	if err == nil {
		hostname = strings.Split(hostname, ".")[0]
	}
	return hostname, err
}

func (m *Machine) Hostname(sreg S2SRegistry, wg *sync.WaitGroup) {
	if hostname, err := m.GetHostname(); err == nil {
		sreg.SetString("hostname", hostname)
	}
	wg.Done()
}

func (m *Machine) Uptime(sreg S2SRegistry, wg *sync.WaitGroup) {
	uptime := sigar.Uptime{}
	uptime.Get()
	sreg.SetString("uptime", format.FormatUptime(uptime.Length))
	wg.Done()
}

func (m *Machine) LA(reg registry.Registry, wg *sync.WaitGroup) {
	la := sigar.LoadAverage{}
	la.Get()
	reg.UpdateLoadAverage(la)
	wg.Done()
}

var UsePercentTemplate *templateutil.BinTemplate

func _getmem(kind string, in sigar.Swap) operating.Memory {
	total, approxtotal, _ := format.HumanBandback(in.Total)
	used, approxused, _ := format.HumanBandback(in.Used)
	usepercent := format.Percent(approxused, approxtotal)

	html := "ERROR"
	if TooltipableTemplate == nil {
		log.Printf("TooltipableTemplate hasn't been set")
	} else if buf, err := UsePercentTemplate.CloneExecute(struct {
		Class, Value, CLASSNAME string
	}{
		Value: strconv.Itoa(int(usepercent)), // without "%"
		Class: format.LabelClassColorPercent(usepercent),
	}); err == nil {
		html = buf.String()
	}

	return operating.Memory{
		Kind:           kind,
		Total:          total,
		Used:           used,
		Free:           format.HumanB(in.Free),
		UsePercentHTML: template.HTML(html),
	}
}

func (m *Machine) RAM(reg registry.Registry, wg *sync.WaitGroup) {
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

func (m *Machine) Swap(reg registry.Registry, wg *sync.WaitGroup) {
	// m is unused
	got := sigar.Swap{}
	got.Get()
	reg.UpdateSwap(got)
	wg.Done()
}

func (m *Machine) Disks(reg registry.Registry, wg *sync.WaitGroup) {
	// m is unused
	fls := sigar.FileSystemList{}
	fls.Get()

	// devnames := map[string]bool{}
	dirnames := map[string]bool{}

	for _, fs := range fls.List {

		usage := sigar.FileSystemUsage{}
		usage.Get(fs.DirName)

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

func (m *Machine) Procs(CH chan<- operating.MetricProcSlice) {
	// m is unused
	var procs operating.MetricProcSlice
	pls := sigar.ProcList{}
	pls.Get()

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

		procs = append(procs, operating.MetricProc{
			ProcInfo: &operating.ProcInfo{
				PID:      uint(pid),
				Priority: state.Priority,
				Nice:     state.Nice,
				Time:     time.Total,
				Name:     system.ProcName(pid, state.Name),
				// Name:  strings.Join(append([]string{system.ProcName(pid, state.Name)}, args.List[1:]...), " "),
				UID:      state.Uid,
				Size:     mem.Size,
				Resident: mem.Resident,
			},
		})
	}
	CH <- procs
}

func (m *Machine) CPU(reg registry.Registry, wg *sync.WaitGroup) {
	// m is unused
	cl := sigar.CpuList{}
	cl.Get()
	reg.UpdateCPU(cl.List)
	wg.Done()
}
