package ostential
import (
	"ostential/types"
	"ostential/view"

	"os"
	"fmt"
	"strconv"
	"strings"
	"html/template"
	"github.com/rzab/gosigar"
)

type generic struct { // ex about
	HostnameHTML   template.HTML
	HostnameString string // for page title
	IP             string

	// ex system
	Uptime string
//	La1    string
//	La5    string
//	La15   string
	LA     string
}

func getGeneric() generic { // ex getAbout
	hostname, err := os.Hostname()
	if err == nil {
		hostname = strings.Split(hostname, ".")[0]
	}
	// IP, _ := netinterface_ipaddr()
	g := generic{
		HostnameString: hostname,
		HostnameHTML: tooltipable(11, hostname),
		// IP: IP,
	}
	uptime := sigar.Uptime{}; uptime.Get()
	g.Uptime = formatUptime(uptime.Length)

	la := sigar.LoadAverage{}
	la.Get()

//	g.La1  = fmt.Sprintf("%.2f", la.One)
//	g.La5  = fmt.Sprintf("%.2f", la.Five)
//	g.La15 = fmt.Sprintf("%.2f", la.Fifteen)
	g.LA   = fmt.Sprintf("%.2f %.2f %.2f", la.One, la.Five, la.Fifteen)

	return g
}

func _getmem(in sigar.Swap) memory {
	total, approxtotal := humanBandback(in.Total)
	used,  approxused  := humanBandback(in.Used)
	usepercent := percent(approxused, approxtotal)

	UPhtml, _ := view.UsePercentTemplate.Execute(struct{
		Class, Value string
	}{
		Value: strconv.Itoa(int(usepercent)), // without "%"
		Class: labelClass_colorPercent(usepercent),
	})

	return memory{
		Total: total,
		Used:  used,
		Free:  humanB(in.Free),
		UsePercentHTML: UPhtml,
	}
}
func getRAM() memory {
	got := sigar.Mem{}; got.Get()
	return _getmem(sigar.Swap{
		Total: got.Total,
		Used:  got.Used,
		Free:  got.Free,
	})
}
func getSwap() memory {
	got := sigar.Swap{}; got.Get()
	return _getmem(got)
}

func read_disks() (disks []diskInfo) {
	fls := sigar.FileSystemList{}
	fls.Get()

// 	devnames := map[string]bool{}
	dirnames := map[string]bool{}

	for _, fs := range fls.List {

		usage := sigar.FileSystemUsage{}
		usage.Get(fs.DirName)

		if  fs.DevName == "shm"    ||
			fs.DevName == "none"   ||
			fs.DevName == "proc"   ||
			fs.DevName == "udev"   ||
			fs.DevName == "devfs"  ||
			fs.DevName == "sysfs"  ||
			fs.DevName == "tmpfs"  ||
			fs.DevName == "devpts" ||
			fs.DevName == "cgroup" ||
			fs.DevName == "rootfs" ||
			fs.DevName == "rpc_pipefs" ||

			fs.DirName == "/dev" ||
			strings.HasPrefix(fs.DevName, "map ") {
			continue
		}
	// 	if _, ok := devnames[fs.DevName]; ok { continue }
		if _, ok := dirnames[fs.DirName]; ok { continue }
	// 	devnames[fs.DevName] = true
		dirnames[fs.DirName] = true

		iusePercent := 0.0
		if usage.Files != 0 {
			iusePercent = float64(100) * float64(usage.Files - usage.FreeFiles) / float64(usage.Files)
		}
		disks = append(disks, diskInfo{
			DevName:     fs.DevName,

			Total:       usage.Total << 10, // * 1024
			Used:        usage.Used  << 10, // == Total - Free
			Avail:       usage.Avail << 10,
			UsePercent:  usage.UsePercent(),

			Inodes:      usage.Files,
			Iused:       usage.Files - usage.FreeFiles,
			Ifree:       usage.FreeFiles,
			IusePercent: iusePercent,

			DirName:     fs.DirName,
		})
	}
	return disks
}

func read_procs() (procs []types.ProcInfo) {
	pls := sigar.ProcList{}
	pls.Get()

	for _, pid := range pls.List {

		state := sigar.ProcState{}
		// args  := sigar.ProcArgs{}
		time  := sigar.ProcTime{}
		mem   := sigar.ProcMem{}

		if err := state.Get(pid); err != nil { continue }
		// if err :=  args.Get(pid); err != nil { continue }
		if err :=  time.Get(pid); err != nil { continue }
		if err :=   mem.Get(pid); err != nil { continue }

		procs = append(procs, types.ProcInfo{
			PID:      uint(pid),
			Priority: state.Priority,
			Nice:     state.Nice,
			Time:     time.Total,
			// `procname' defined proc_{darwin,linux}.go
			Name:     procname(pid, state.Name),
			// Name:     strings.Join(append([]string{procname(pid, state.Name)}, args.List[1:]...), " "),
			Uid:      state.Uid,
			Size:     mem.Size,
			Resident: mem.Resident,
		})
	}
	return procs
}
