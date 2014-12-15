package ostent

import (
	"html/template"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/ostrost/ostent/format"
	"github.com/ostrost/ostent/registry"
	"github.com/ostrost/ostent/templates"
	"github.com/ostrost/ostent/types"
	sigar "github.com/rzab/gosigar"
)

type generic struct {
	Hostname string
	Uptime   string
	IP       string // not filled by getGeneric
	LA       string // not filled by getGeneric
}

func getHostname() (string, error) {
	hostname, err := os.Hostname()
	if err == nil {
		hostname = strings.Split(hostname, ".")[0]
	}
	return hostname, err
}

func getGeneric(reg registry.Registry, CH chan<- generic) {
	hostname, _ := getHostname()

	uptime := sigar.Uptime{}
	uptime.Get()

	la := sigar.LoadAverage{}
	la.Get()

	reg.UpdateLoadAverage(la)

	g := generic{
		Hostname: hostname,
		Uptime:   format.FormatUptime(uptime.Length),
	}
	// IP, _ := netinterface_ipaddr(); CH <- g
	CH <- g
}

var UsePercentTemplate *templates.BinTemplate

func _getmem(kind string, in sigar.Swap) types.Memory {
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

	return types.Memory{
		Kind:           kind,
		Total:          total,
		Used:           used,
		Free:           format.HumanB(in.Free),
		UsePercentHTML: template.HTML(html),
	}
}

func getRAM(reg registry.Registry, wg *sync.WaitGroup) {
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

func getSwap(reg registry.Registry, wg *sync.WaitGroup) {
	got := sigar.Swap{}
	got.Get()
	reg.UpdateSwap(got)
	wg.Done()
}

func read_disks(reg registry.Registry, wg *sync.WaitGroup) {
	fls := sigar.FileSystemList{}
	fls.Get()

	// devnames := map[string]bool{}
	dirnames := map[string]bool{}

	for _, fs := range fls.List {

		usage := sigar.FileSystemUsage{}
		usage.Get(fs.DirName)

		if fs.DevName == "shm" ||
			fs.DevName == "none" ||
			fs.DevName == "proc" ||
			fs.DevName == "udev" ||
			fs.DevName == "devfs" ||
			fs.DevName == "sysfs" ||
			fs.DevName == "tmpfs" ||
			fs.DevName == "devpts" ||
			fs.DevName == "cgroup" ||
			fs.DevName == "rootfs" ||
			fs.DevName == "rpc_pipefs" ||

			fs.DirName == "/dev" ||
			strings.HasPrefix(fs.DevName, "map ") {
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

func read_procs(CH chan<- []types.ProcInfo) {
	var procs []types.ProcInfo
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

		procs = append(procs, types.ProcInfo{
			PID:      uint(pid),
			Priority: state.Priority,
			Nice:     state.Nice,
			Time:     time.Total,
			Name:     procname(pid, state.Name), // `procname' defined proc_{darwin,linux}.go
			// Name:     strings.Join(append([]string{procname(pid, state.Name)}, args.List[1:]...), " "),
			UID:      state.Uid,
			Size:     mem.Size,
			Resident: mem.Resident,
		})
	}
	CH <- procs
}

func CollectCPU(reg registry.Registry, wg *sync.WaitGroup) {
	cl := sigar.CpuList{}
	cl.Get()
	reg.UpdateCPU(cl.List)
	wg.Done()
}
