package ostential
import (
	"ostential/types"
	"ostential/view"

	"fmt"
	"sort"
	"sync"
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
	if nowin < previn { // counters got reset
		return ""
	}
	return humanUnitless(uint64(nowin - previn))
}

const TOPROWS = 2

func interfaceMeta(ii InterfaceInfo) types.InterfaceMeta {
	return types.InterfaceMeta{
		NameKey:  ii.Name,
		NameHTML: tooltipable(12, ii.Name),
	}
}

type interfaceFormat interface {
	Current(*types.Interface, InterfaceInfo)
	Delta  (*types.Interface, InterfaceInfo, InterfaceInfo)
}
type interfaceInout interface {
	InOut(InterfaceInfo) (uint, uint)
}

type interfaceBytes struct{}
func (_ interfaceBytes) Current(id *types.Interface, ii InterfaceInfo) {
	id.In  = humanB(uint64(ii. InBytes))
	id.Out = humanB(uint64(ii.OutBytes))
}
func (_ interfaceBytes) Delta(id *types.Interface, ii, pi InterfaceInfo) {
	id.DeltaIn  = bps(8, ii. InBytes, pi. InBytes)
	id.DeltaOut = bps(8, ii.OutBytes, pi.OutBytes)
}

type interfaceInoutErrors struct{}
func (_ interfaceInoutErrors) InOut(ii InterfaceInfo) (uint, uint) {
	return ii.InErrors, ii.OutErrors
}
type interfaceInoutPackets struct{}
func (_ interfaceInoutPackets) InOut(ii InterfaceInfo) (uint, uint) {
	return ii.InPackets, ii.OutPackets
}

type interfaceNumericals struct{interfaceInout}
func (ie interfaceNumericals) Current(id *types.Interface, ii InterfaceInfo) {
	in, out := ie.InOut(ii)
	id.In  = humanUnitless(uint64(in))
	id.Out = humanUnitless(uint64(out))
}
func (ie interfaceNumericals) Delta(id *types.Interface, ii, previousi InterfaceInfo) {
	in, out                   := ie.InOut(ii)
	previous_in, previous_out := ie.InOut(previousi)
	id.DeltaIn  = ps(in,  previous_in)
	id.DeltaOut = ps(out, previous_out)
}

func InterfacesDelta(format interfaceFormat, current, previous []InterfaceInfo, client clientState) *types.Interfaces {
	ifs := make([]types.Interface, len(current))

	for i := range ifs {
		di := types.Interface{
			InterfaceMeta: interfaceMeta(current[i]),
		}
		format.Current(&di, current[i])

		if len(previous) > i {
			format.Delta(&di, current[i], previous[i])
		}

		ifs[i] = di
	}
	if len(ifs) > 1 {
		sort.Sort(interfaceOrder(ifs))
		if !*client.ExpandIF && len(ifs) > TOPROWS {
			ifs = ifs[:TOPROWS]
		}
	}
	ni := new(types.Interfaces)
	ni.List = ifs
	return ni
}

func(li lastinfo) cpuListDelta() sigar.CpuList {
	prev := li.Previous.CPU
	if len(prev.List) == 0 {
		return li.CPU
	}
	coreno := len(li.CPU.List)
	if coreno == 0 { // wait, what?
		return sigar.CpuList{}
	}
	cls := sigar.CpuList{List: make([]sigar.Cpu, coreno) }
	copy(cls.List, li.CPU.List)
	for i := range cls.List {
		cls.List[i].User -= prev.List[i].User
		cls.List[i].Nice -= prev.List[i].Nice
		cls.List[i].Sys  -= prev.List[i].Sys
		cls.List[i].Idle -= prev.List[i].Idle
	}
	return cls
}

func(li lastinfo) CPUDelta(client clientState) *types.CPU {
	cls := li.cpuListDelta()
	coreno := len(cls.List)
	if coreno == 0 { // wait, what?
		return &types.CPU{}
	}

	sum := sigar.Cpu{}
	cores := make([]types.Core, coreno)
	for i, each := range cls.List {

		total := each.User + each.Nice + each.Sys + each.Idle

		user := percent(each.User, total)
		sys  := percent(each.Sys,  total)

		idle := uint(0)
		if user + sys < 100 {
			idle = 100 - user - sys
		}

		cores[i] = types.Core{
			N: fmt.Sprintf("#%d", i),
			User: user,
			Sys:  sys,
			Idle: idle,
			UserClass:  textClass_colorPercent(user),
			SysClass:   textClass_colorPercent(sys),
			IdleClass:  textClass_colorPercent(100 - idle),
		}

		sum.User += each.User + each.Nice
		sum.Sys  += each.Sys
		sum.Idle += each.Idle
	}

	cpu := new(types.CPU)
	cpu.DataMeta = types.NewDataMeta()

	if coreno == 1 {
		cores[0].N = "#0"
		*cpu.DataMeta.ExpandText = "1"
		cpu.List = cores
		return cpu
	}
	sort.Sort(cpuOrder(cores))

	*cpu.DataMeta.ExpandText = fmt.Sprintf("%d", coreno)
	*cpu.DataMeta.Expandable = coreno > TOPROWS-1 // one row reserved for "all N"

	if !*client.ExpandCPU {
		if coreno > TOPROWS-1 {
			cores = cores[:TOPROWS-1] // first core(s)
		}

		total := sum.User + sum.Sys + sum.Idle // + sum.Nice

		user := percent(sum.User, total)
		sys  := percent(sum.Sys,  total)
		idle := uint(0)
		if user + sys < 100 {
			idle = 100 - user - sys
		}
		cores = append([]types.Core{{ // "all N"
			N: fmt.Sprintf("all %d", coreno),
			User: user,
			Sys:  sys,
			Idle: idle,
			UserClass: textClass_colorPercent(user),
			SysClass:  textClass_colorPercent(sys),
			IdleClass: textClass_colorPercent(100 - idle),
		}}, cores...)
	}

	cpu.List = cores
	return cpu
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
	Total          string
	Used           string
	Free           string
	UsePercentHTML template.HTML
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

func tooltipable(limit int, full string) template.HTML {
	if len(full) > limit {
		short := full[:limit]
		if html, err := view.TooltipableTemplate.Execute(struct {
			Full, Short string
		}{
			Full: full,
			Short: short,
		}); err == nil {
			return html
		}
	}
	return template.HTML(template.HTMLEscapeString(full))
}

func orderDisks(disks []diskInfo, seq types.SEQ) []diskInfo {
	if len(disks) > 1 {
		sort.Stable(diskOrder{
			disks: disks,
			seq: seq,
			reverse: _DFBIMAP.SEQ2REVERSE[seq],
		})
	}
	return disks
}

func diskMeta(disk diskInfo) types.DiskMeta {
	return types.DiskMeta{
		DiskNameHTML: tooltipable(12, disk.DevName),
		DirNameHTML:  tooltipable(6, disk.DirName),
		DirNameKey:   disk.DirName,
	}
}

func dfbytes(diskinfos []diskInfo, client clientState) *types.DFbytes {
	var disks []types.DiskBytes
	for i, disk := range diskinfos {
		if !*client.ExpandDF && i > 1 {
			break
		}
		total,  approxtotal  := humanBandback(disk.Total)
		used,   approxused   := humanBandback(disk.Used)
		disks = append(disks, types.DiskBytes{
			DiskMeta: diskMeta(disk),
			Total:       total,
			Used:        used,
			Avail:       humanB(disk.Avail),
			UsePercent:  formatPercent(approxused, approxtotal),
			UsePercentClass: labelClass_colorPercent(percent(approxused,  approxtotal)),
		})
	}
	// if !*client.ExpandDF && len(disks) > TOPROWS {
	// 	disks = disks[:TOPROWS]
	// }
	dsb := new(types.DFbytes)
	dsb.List = disks
	return dsb
}

func dfinodes(diskinfos []diskInfo, client clientState) *types.DFinodes {
	var disks []types.DiskInodes
	for i, disk := range diskinfos {
		if !*client.ExpandDF && i > 1 {
			break
		}
		itotal, approxitotal := humanBandback(disk.Inodes)
		iused,  approxiused  := humanBandback(disk.Iused)
		disks = append(disks, types.DiskInodes{
			DiskMeta: diskMeta(disk),
			Inodes:      itotal,
			Iused:       iused,
			Ifree:       humanB(disk.Ifree),
			IusePercent: formatPercent(approxiused, approxitotal),
			IusePercentClass: labelClass_colorPercent(percent(approxiused, approxitotal)),
		})
	}
	dsi := new(types.DFinodes)
	dsi.List = disks
	return dsi
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

func orderProc(procs []types.ProcInfo, seq types.SEQ, clientptr *clientState) ([]types.ProcData, string) {
	client := *clientptr
	sort.Sort(procOrder{ // not sort.Stable
		procs: procs,
		seq: seq,
		reverse: _PSBIMAP.SEQ2REVERSE[seq],
	})

	limitPS := client.psLimit

	if len(procs) <= limitPS {
		limitPS = len(procs)

		if client.psNotexpandable == nil || !*client.psNotexpandable {
			clientptr.psNotexpandable = newtrue()
		}
	} else if clientptr.psNotexpandable != nil {

		if *client.psNotexpandable {
			*clientptr.psNotexpandable = false
		} else {
			clientptr.psNotexpandable = nil
		}
	}

	if len(procs) > limitPS {
		procs = procs[:limitPS]
	}
	plustext := fmt.Sprintf("%d+", limitPS)

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
	return list, plustext
}

type Previous struct {
	CPU        sigar.CpuList
	Interfaces []InterfaceInfo
}

type lastinfo struct {
    Generic generic
	CPU     sigar.CpuList
	RAM     memory
	Swap    memory
	DiskList   []diskInfo
	ProcList   []types.ProcInfo
	Interfaces []InterfaceInfo
	Previous Previous
}

type PageData struct {
    Generic generic
	CPU     types.CPU
	RAM     memory
	Swap    memory

	PStable PStable

	DFlinks *DFlinks        `json:",omitempty"`
	DFbytes  types.DFbytes  `json:",omitempty"`
	DFinodes types.DFinodes `json:",omitempty"`

	DF types.DataMeta
	IF types.DataMeta

	IFbytes   types.Interfaces
	IFerrors  types.Interfaces
	IFpackets types.Interfaces

	DISTRIBHTML template.HTML
	VERSION     string
	HTTP_HOST   string

    ClientState clientState
}

type pageUpdate struct {
    Generic  generic

	RAM      memory // TODO empty on HideMEM
	Swap     memory // TODO empty on HideMEM

	CPU      *types.CPU      `json:",omitempty"`

	DFlinks  *DFlinks        `json:",omitempty"`
	DFbytes  *types.DFbytes  `json:",omitempty"`
	DFinodes *types.DFinodes `json:",omitempty"`

	DF types.DataMeta // a pointer and `json:",omitempty"` ?
	IF types.DataMeta // a pointer and `json:",omitempty"` ?

	// PSlinks *PSlinks `json:",omitempty"`
	PStable *PStable `json:",omitempty"`

	IFbytes   *types.Interfaces `json:",omitempty"`
	IFerrors  *types.Interfaces `json:",omitempty"`
	IFpackets *types.Interfaces `json:",omitempty"`

	ClientState *clientState `json:",omitempty"`
}

var (
	lastLock sync.Mutex
	lastInfo lastinfo
)

func reset_prev() {
	lastLock.Lock()
	defer lastLock.Unlock()

	lastInfo.Previous.CPU        = sigar.CpuList{}
	lastInfo.Previous.Interfaces = []InterfaceInfo{}
}

func collect() {
	lastLock.Lock()
	defer lastLock.Unlock()

	previous := Previous{
		CPU:        lastInfo.CPU,
		Interfaces: lastInfo.Interfaces,
	}

	ifs, ip := NewInterfaces()
	generic := getGeneric()
	generic.IP = ip

	lastInfo = lastinfo{
		Generic:  generic,
		RAM:      getRAM(),
		Swap:     getSwap(),
		DiskList: read_disks(),
		ProcList: read_procs(),
	}
	cl := sigar.CpuList{}; cl.Get()
	lastInfo.CPU  = cl

	lastInfo.Interfaces = filterInterfaces(ifs)
	lastInfo.Previous = previous
}

func linkattrs(req *http.Request, base url.Values, pname string, bimap types.Biseqmap) types.Linkattrs {
	return types.Linkattrs{
		Base:  base,
		Pname: pname,
		Bimap: bimap,
		Seq:   valuesSet(req, base, pname, bimap),
	}
}

func getUpdates(req *http.Request, new_search bool, clientptr *clientState, clientdiff *clientState) (pageUpdate, url.Values, types.SEQ, types.SEQ) {
	client := *clientptr

	req.ParseForm()
	base := url.Values{}

	var (
		df_copy []diskInfo
		ps_copy []types.ProcInfo
		if_copy     []InterfaceInfo
		previf_copy []InterfaceInfo
	)

	var pu pageUpdate
	func() {
		lastLock.Lock()
		defer lastLock.Unlock()

		df_copy = make([]diskInfo,       len(lastInfo.DiskList))
		ps_copy = make([]types.ProcInfo, len(lastInfo.ProcList))
		copy(df_copy, lastInfo.DiskList)
		copy(ps_copy, lastInfo.ProcList)

		if_copy     = make([]InterfaceInfo, len(lastInfo.Interfaces))
		previf_copy = make([]InterfaceInfo, len(lastInfo.Previous.Interfaces))
		copy(if_copy,     lastInfo.Interfaces)
		copy(previf_copy, lastInfo.Previous.Interfaces)

		pu = pageUpdate{
			Generic: lastInfo.Generic,
		}
		if !*client.HideMEM {
			pu.RAM  = lastInfo.RAM
			pu.Swap = lastInfo.Swap
		}
		if !*client.HideCPU {
			pu.CPU = lastInfo.CPUDelta(client)
		}
	}()

	 pu.IF = types.NewDataMeta()
	*pu.IF.Expandable = len(if_copy) > TOPROWS
	*pu.IF.ExpandText = fmt.Sprintf("%d", len(if_copy))

	pslinks := PSlinks(linkattrs(req, base, "ps", _PSBIMAP))
	dflinks := DFlinks(linkattrs(req, base, "df", _DFBIMAP))

	 pu.DF = types.NewDataMeta()
	*pu.DF.Expandable = len(df_copy) > TOPROWS
	*pu.DF.ExpandText = fmt.Sprintf("%d", len(df_copy))

	if !*client.HideDF {
		orderedDisks := orderDisks(df_copy, dflinks.Seq)

		       if *client.TabDF == DFBYTES_TABID  { pu.DFbytes  = dfbytes (orderedDisks, client)
		} else if *client.TabDF == DFINODES_TABID { pu.DFinodes = dfinodes(orderedDisks, client)
		}
	}

	if !*client.HideIF {
		switch *client.TabIF {
		case IFBYTES_TABID:   pu.IFbytes   = InterfacesDelta(interfaceBytes{},                             if_copy, previf_copy, client)
		case IFERRORS_TABID:  pu.IFerrors  = InterfacesDelta(interfaceNumericals{interfaceInoutErrors{}},  if_copy, previf_copy, client)
		case IFPACKETS_TABID: pu.IFpackets = InterfacesDelta(interfaceNumericals{interfaceInoutPackets{}}, if_copy, previf_copy, client)
		}
	}

	if !*client.HidePS {
		pu.PStable = new(PStable)
		pu.PStable.List, pu.PStable.PlusText = orderProc(ps_copy, pslinks.Seq, clientptr)
		pu.PStable.NotExpandable = clientptr.psNotexpandable
	}
	if new_search {
		pu.PStable.Links = &pslinks
		pu.DFlinks       = &dflinks
	}

	if clientdiff != nil {
		 pu.ClientState = new(clientState)
		*pu.ClientState = *clientdiff // client
	}
	return pu, base, dflinks.Seq, pslinks.Seq
}

var DISTRIB string // set with init from init_*.go
func pageData(req *http.Request) PageData {
	client := defaultClientState()
	updates, base, dfseq, psseq := getUpdates(req, false, &client, &client)

	dla := &DFlinks{
		Base: base,
		Pname: "df",
		Bimap: _DFBIMAP,
		Seq: dfseq,
	}
	pla := &PSlinks{
		Base: base,
		Pname: "ps",
		Bimap: _PSBIMAP,
		Seq: psseq,
	}

	data := PageData{
		ClientState: *updates.ClientState,
		Generic:      updates.Generic,
		CPU:         *updates.CPU,
		RAM:          updates.RAM,
		Swap:         updates.Swap,

		DFlinks:  dla,

		PStable: PStable{
			List:          updates.PStable.List,
			PlusText:      updates.PStable.PlusText,
			NotExpandable: updates.PStable.NotExpandable,
			Links: pla,
		},

		DISTRIBHTML: tooltipable(11, DISTRIB), // value from init_*.go
		VERSION:     VERSION,                  // value from server.go
		HTTP_HOST:   req.Host,
	}
	data.DF = updates.DF
	data.IF = updates.IF

	       if updates.DFbytes  != nil { data.DFbytes  = *updates.DFbytes
	} else if updates.DFinodes != nil { data.DFinodes = *updates.DFinodes
	}

	       if updates.IFbytes   != nil { data.IFbytes   = *updates.IFbytes
	} else if updates.IFerrors  != nil { data.IFerrors  = *updates.IFerrors
	} else if updates.IFpackets != nil { data.IFpackets = *updates.IFpackets
	}

	return data
}

func index(req *http.Request, r view.Render) {
	r.HTML(200, "index.html", struct{Data interface{}}{Data: pageData(req),})
}

type Modern struct {
	*martini.Martini
	 martini.Router // the router functions for convenience
}








