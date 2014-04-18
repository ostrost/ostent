package ostential
import (
	"ostential/types"
	"ostential/view"

	"fmt"
	"sort"
	"sync"
	"bytes"
	"os/user"
	"net/url"
	"net/http"
	"html/template"

	"github.com/rzab/gosigar"
	"github.com/codegangsta/martini"
)

func bps(factor int, nowin, previn uint) string {
	if nowin < previn { // counters got reset
		return ""
	}
	n := (nowin - previn) * uint(factor) // bits now?
	return humanbits(uint64(n))
}
func unitless(s string) string {
	if s == "" {
		return s
	}
	i := len(s) - 1
	if s[i] == 'b' || s[i] == 'B' {
		s = s[:i]
	} else if s[i] == 'K' {
		s = s[:i] + "k"
	}
	return s
}
func ps(nowin, previn uint) string {
	return unitless(bps(1, nowin, previn))
}

func(s state) InterfacesDelta() types.Interfaces {
	ifs := make([]types.DeltaInterface, len(s.InterfacesTotal))

	for i := range ifs {
		di := types.DeltaInterface{
			NameKey:                    s.InterfacesTotal[i].Name,
			NameHTML:   tooltipable(12, s.InterfacesTotal[i].Name),

			InBytes:    humanB(uint64(s.InterfacesTotal[i]. InBytes)),
			OutBytes:   humanB(uint64(s.InterfacesTotal[i].OutBytes)),
			InPackets:  unitless(humanB(uint64(s.InterfacesTotal[i]. InPackets))),
			OutPackets: unitless(humanB(uint64(s.InterfacesTotal[i].OutPackets))),
			InErrors:   unitless(humanB(uint64(s.InterfacesTotal[i]. InErrors))),
			OutErrors:  unitless(humanB(uint64(s.InterfacesTotal[i].OutErrors))),
		}
		if len(s.PrevInterfacesTotal) > i {
			di.DeltaInBytes    = bps(8, s.InterfacesTotal[i]. InBytes,   s.PrevInterfacesTotal[i]. InBytes)
			di.DeltaOutBytes   = bps(8, s.InterfacesTotal[i].OutBytes,   s.PrevInterfacesTotal[i].OutBytes)
			di.DeltaInPackets  =  ps(   s.InterfacesTotal[i]. InPackets, s.PrevInterfacesTotal[i]. InPackets)
			di.DeltaOutPackets =  ps(   s.InterfacesTotal[i].OutPackets, s.PrevInterfacesTotal[i].OutPackets)
			di.DeltaInErrors   =  ps(   s.InterfacesTotal[i]. InErrors,  s.PrevInterfacesTotal[i]. InErrors)
			di.DeltaOutErrors  =  ps(   s.InterfacesTotal[i].OutErrors,  s.PrevInterfacesTotal[i].OutErrors)
		}
		ifs[i] = di
	}
	var haveCollapsed bool
	if len(ifs) > 1 {
		sort.Sort(interfaceOrder(ifs))
		for i := range ifs { // set collapse after sort
			haveCollapsed = i > 1
			ifs[i].CollapseClass = map[bool]string{true: "collapse"}[haveCollapsed]
		}
	}
	return types.Interfaces{List: ifs, HaveCollapsed: haveCollapsed}
}

func(s state) cpudelta() sigar.CpuList {
	prev := s.PREVCPU
	if len(prev.List) == 0 {
		return s.RAWCPU
	}
	listlen := len(s.RAWCPU.List)
	if listlen == 0 { // wait, what?
		return sigar.CpuList{}
	}
// 	cls := s.RAWCPU
	cls := sigar.CpuList{List: make([]sigar.Cpu, len(s.RAWCPU.List)) }
	copy(cls.List, s.RAWCPU.List)
	for i := range cls.List {
		cls.List[i].User -= prev.List[i].User
		cls.List[i].Nice -= prev.List[i].Nice
		cls.List[i].Sys  -= prev.List[i].Sys
		cls.List[i].Idle -= prev.List[i].Idle
	}
	return cls
}

func(s state) CPU() types.CPU {
	sum := sigar.Cpu{}
	cls := s.cpudelta()
	coreno := len(cls.List)
	if coreno == 0 { // wait, what?
		return types.CPU{}
	}

	c := make([]types.Core, coreno + 1) // + total
	for i, cp := range cls.List {

		total := cp.User + cp.Nice + cp.Sys + cp.Idle

		user := percent(cp.User, total)
		sys  := percent(cp.Sys,  total)

		idle := uint(0)
		if user + sys < 100 {
			idle = 100 - user - sys
		}

		i++ // won't use c[0], it's for totals
		c[i].N    = fmt.Sprintf("#%d", i-1)
 		c[i].User, c[i].UserClass = user, textClass_colorPercent(user)
 		c[i].Sys,  c[i]. SysClass = sys,  textClass_colorPercent(sys)
		c[i].Idle, c[i].IdleClass = idle, textClass_colorPercent(100 - idle)

		sum.User += cp.User + cp.Nice
		sum.Sys  += cp.Sys
		sum.Idle += cp.Idle
	}
	if coreno == 1 {
		c[1].N = "#0"
		return types.CPU{List: c[1:]}
	}
	sort.Sort(cpuOrder(c[1:]))
	// collapse after sorting

	var haveCollapsed bool
	for i := range c[1:] { // c[0] is for totals
		haveCollapsed = i > 0 // collapse all but one
		c[i + 1].CollapseClass = map[bool]string{true: "collapse"}[haveCollapsed]
	}

	total := sum.User + sum.Sys + sum.Idle // + sum.Nice

	user := percent(sum.User, total)
	sys  := percent(sum.Sys,  total)
	idle := uint(0)
	if user + sys < 100 {
		idle = 100 - user - sys
	}

	c[0].N                    = fmt.Sprintf("all %d", coreno)
 	c[0].User, c[0].UserClass = user, textClass_colorPercent(user)
 	c[0].Sys,  c[0]. SysClass = sys,  textClass_colorPercent(sys)
	c[0].Idle, c[0].IdleClass = idle, textClass_colorPercent(100 - idle)

	return types.CPU{List: c, HaveCollapsed: haveCollapsed}
}

func textClass_colorPercent(p uint) string {
	return "text-" + colorPercent(p)
}

func labelClass_colorPercent(p uint) string {
	return "label label-" + colorPercent(p)
}

func colorPercent(p uint) string {
	if p > 90 { return "danger"  }
	if p > 80 { return "warning" }
	if p > 20 { return "info"    }
	return "success"
}

type memory struct {
	Total       string
	Used        string
	Free        string
	UsePercent  string

	UsePercentClass string
}

type diskInfo struct {
	DevName     string
	Total       uint64
	Used        uint64
	Avail       uint64
	UsePercent  float64
	Inodes      uint64
	Iused       uint64
	Ifree       uint64
	IusePercent float64
	DirName     string
}

func valuesSet(req *http.Request, base url.Values, pname string, bimap types.Biseqmap) types.SEQ {
	if params, ok := req.Form[pname]; ok && len(params) > 0 {
		if seq, ok := bimap.STRING2SEQ[params[0]]; ok {
			base.Set(pname, params[0])
			return seq
		}
	}
	return bimap.Default_seq
}

var (
	_attr_start = "<span title=\""
	_attr_end   = "\" />"
	_attr_template = template.Must(template.New("attr").Parse(_attr_start +"{{.}}"+ _attr_end))
)

func attribute_escape(data string) string {
	if _template, err := _attr_template.Clone(); err == nil {
		buf := new(bytes.Buffer)
		if err := _template.Execute(buf, data); err == nil {
			s := buf.String()
			return s[len(_attr_start):len(s) - len(_attr_end)]
		}
	}
	return ""
}

func tooltipable(limit int, devname string) template.HTML {
	if len(devname) <= limit {
		return template.HTML(devname)
	}
	title_attr := attribute_escape(devname)
	short := template.HTMLEscapeString(devname[:limit])
	s := template.HTML(fmt.Sprintf(`
<span title="%s" class="tooltipable" data-toggle="tooltip" data-placement="left">%s<span class="inlinecode">...</span></span>`,
		title_attr, short))
	return s
}

func orderDisk(disks []diskInfo, seq types.SEQ) DiskTable {
	if len(disks) > 1 {
		sort.Stable(diskOrder{
			disks: disks,
			seq: seq,
			reverse: _DFBIMAP.SEQ2REVERSE[seq],
		})
	}

	var dd []types.DiskData
	var haveCollapsed bool
	for i, disk := range disks {
		total,  approxtotal  := humanBandback(disk.Total)
		used,   approxused   := humanBandback(disk.Used)
		itotal, approxitotal := humanBandback(disk.Inodes)
		iused,  approxiused  := humanBandback(disk.Iused)

		haveCollapsed = i > 1
		dd = append(dd, types.DiskData{
			DiskNameKey:  disk.DevName,
			DiskNameHTML: tooltipable(12, disk.DevName),
			DirNameHTML:  tooltipable(8, disk.DirName),

			Total:       total,
			Used:        used,
			Avail:       humanB(disk.Avail),
			UsePercent:  formatPercent(approxused, approxtotal),

			Inodes:      itotal,
			Iused:       iused,
			Ifree:       humanB(disk.Ifree),
			IusePercent: formatPercent(approxiused, approxitotal),

			 UsePercentClass: labelClass_colorPercent(percent(approxused,  approxtotal)),
			IusePercentClass: labelClass_colorPercent(percent(approxiused, approxitotal)),

			CollapseClass: map[bool]string{true: "collapse"}[haveCollapsed],
		})
	}
	return DiskTable{List: dd, HaveCollapsed: haveCollapsed}
}

var _DFBIMAP = types.Seq2bimap(DFFS, // the default seq for ordering
	types.Seq2string{
		DFFS:      "fs",
		DFSIZE:    "size",
		DFUSED:    "used",
		DFAVAIL:   "avail",
		DFMP:      "mp",
	}, []types.SEQ{
		DFFS, DFMP,
	})

var _PSBIMAP = types.Seq2bimap(PSPID, // the default seq for ordering
	types.Seq2string{
		PSPID:   "pid",
		PSPRI:   "pri",
		PSNICE:  "nice",
		PSSIZE:  "size",
		PSRES:   "res",
		PSTIME:  "time",
		PSNAME:  "name",
		PSUID:   "user",
	}, []types.SEQ{
		PSNAME, PSUID,
	})

func username(uids map[uint]string, uid uint) string {
	if s, ok := uids[uid]; ok {
		return s
	}
	s := fmt.Sprintf("%d", uid)
	if usr, err := user.LookupId(s); err == nil {
		s = usr.Username
	}
	uids[uid] = s
	return s
}

func orderProc(procs []types.ProcInfo, seq types.SEQ) []types.ProcData {
	sort.Sort(procOrder{ // not sort.Stable
		procs: procs,
		seq: seq,
		reverse: _PSBIMAP.SEQ2REVERSE[seq],
	})

	if len(procs) > 20 {
		procs = procs[:20]
	}

	uids := map[uint]string{}
	var list []types.ProcData
	for _, proc := range procs {
		list = append(list, types.ProcData{
			PID:        proc.PID,
			Priority:   proc.Priority,
			Nice:       proc.Nice,
			Time:       formatTime(proc.Time),
			NameHTML:   tooltipable(42, proc.Name),
			UserHTML:   tooltipable(12, username(uids, proc.Uid)),
			Size:       humanB(proc.Size),
			Resident:   humanB(proc.Resident),
		})
	}
	return list
}

type state struct {
    About    about
    System   system
	RAWCPU   sigar.CpuList
	PREVCPU  sigar.CpuList
	RAM      memory
	Swap     memory
	DiskList []diskInfo
	ProcList []types.ProcInfo

	InterfacesTotal     []InterfaceTotal
	PrevInterfacesTotal []InterfaceTotal
}

type Page struct {
    About     about
    System    system
	CPU       types.CPU
	RAM       memory
	Swap      memory
	DiskTable DiskTable
	ProcTable ProcTable

	Interfaces types.Interfaces

	DISTRIBHTML template.HTML
	VERSION     string
	HTTP_HOST   string
}
type pageUpdate struct {
    About    about
    System   system
	CPU      types.CPU
	RAM      memory
	Swap     memory

	DiskTable DiskTable
	ProcTable ProcTable

	Interfaces types.Interfaces
}

var stateLock sync.Mutex
var lastState state
func reset_prev() {
	stateLock.Lock()
	defer stateLock.Unlock()
	lastState.PrevInterfacesTotal = []InterfaceTotal{}
	lastState.PREVCPU.List        = []sigar.Cpu{}
}
func collect() { // state
	stateLock.Lock()
	defer stateLock.Unlock()

	prev_ifstotal := lastState.InterfacesTotal
	prev_cpu      := lastState.RAWCPU

	ifs, ip := NewInterfaces()
	about := getAbout()
	about.IP = ip

	lastState = state{
		About:    about,
		System:   getSystem(),
		RAM:      getRAM(),
		Swap:     getSwap(),
		DiskList: read_disks(),
		ProcList: read_procs(),
	}
	cl := sigar.CpuList{}; cl.Get()
	lastState.PREVCPU = prev_cpu
	lastState.RAWCPU  = cl

	ifstotal := filterInterfaces(ifs)
	lastState.PrevInterfacesTotal = prev_ifstotal
	lastState.InterfacesTotal     = ifstotal

//	return lastState
}

func linkattrs(req *http.Request, base url.Values, pname string, bimap types.Biseqmap) types.Linkattrs {
	return types.Linkattrs{
		Base:  base,
		Pname: pname,
		Bimap: bimap,
		Seq:   valuesSet(req, base, pname, bimap),
	}
}

func updates(req *http.Request, new_search bool) (pageUpdate, url.Values, types.SEQ, types.SEQ) {
	req.ParseForm()
	base := url.Values{}

	dflinks := DiskLinkattrs(linkattrs(req, base, "df", _DFBIMAP))
	pslinks := ProcLinkattrs(linkattrs(req, base, "ps", _PSBIMAP))

	var pu pageUpdate
	var disks_copy []diskInfo
	var procs_copy []types.ProcInfo
	func() {
		stateLock.Lock()
		defer stateLock.Unlock()

		disks_copy = make([]diskInfo, len(lastState.DiskList))
		procs_copy = make([]types.ProcInfo, len(lastState.ProcList))
		copy(disks_copy, lastState.DiskList)
		copy(procs_copy, lastState.ProcList)

		pu = pageUpdate{
			About:    lastState.About,
			System:   lastState.System,
			CPU:      lastState.CPU(),
			RAM:      lastState.RAM,
			Swap:     lastState.Swap,
			Interfaces: lastState.InterfacesDelta(),
		}
	}()
	pu.DiskTable      = orderDisk(disks_copy, dflinks.Seq)
	pu.ProcTable.List = orderProc(procs_copy, pslinks.Seq)
	if new_search {
		pu.ProcTable.Links = &pslinks
		pu.DiskTable.Links = &dflinks
	}
	return pu, base, dflinks.Seq, pslinks.Seq
}

var DISTRIB string // set with init from init_*.go
func collected(req *http.Request) Page {
	latest, base, dfseq, psseq := updates(req, false)
	return Page{
		About:   latest.About,
		System:  latest.System,
		CPU:     latest.CPU,
		RAM:     latest.RAM,
		Swap:    latest.Swap,
		DiskTable: DiskTable{
			List:          latest.DiskTable.List,
			HaveCollapsed: latest.DiskTable.HaveCollapsed,
			Links: &DiskLinkattrs{
				Base: base,
				Pname: "df",
				Bimap: _DFBIMAP,
				Seq: dfseq,
			},
		},
		ProcTable: ProcTable{
			List: latest.ProcTable.List,
			Links: &ProcLinkattrs{
				Base: base,
				Pname: "ps",
				Bimap: _PSBIMAP,
				Seq: psseq,
			},
		},
		Interfaces: latest.Interfaces, // types.Interfaces{List: latest.Interfaces },
		DISTRIBHTML: tooltipable(11, DISTRIB), // value from init_*.go
		VERSION: VERSION,                      // value from server.go
		HTTP_HOST: req.Host,
	}
}

func index(req *http.Request, r view.Render) {
	r.HTML(200, "index.html", struct{Data interface{}}{collected(req)})
}

type Modern struct {
	*martini.Martini
	 martini.Router // the router functions for convenience
}
