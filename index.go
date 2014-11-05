package ostent

import (
	"container/ring"
	"errors"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"net/url"
	"os/user"
	"sort"
	"strings"
	"sync"

	"github.com/ostrost/ostent/client"
	"github.com/ostrost/ostent/cpu"
	"github.com/ostrost/ostent/format"
	"github.com/ostrost/ostent/getifaddrs"
	"github.com/ostrost/ostent/templates"
	"github.com/ostrost/ostent/types"
	sigar "github.com/rzab/gosigar"
)

func interfaceMeta(ifdata getifaddrs.IfData) types.InterfaceMeta {
	return types.InterfaceMeta{
		NameKey:  ifdata.Name,
		NameHTML: tooltipable(12, ifdata.Name),
	}
}

type interfaceFormat interface {
	Current(*types.Interface, getifaddrs.IfData)
	Delta(*types.Interface, getifaddrs.IfData, getifaddrs.IfData)
}

type interfaceInout interface {
	InOut(getifaddrs.IfData) (uint, uint)
}

type interfaceBytes struct{}

func (_ interfaceBytes) Current(id *types.Interface, ifdata getifaddrs.IfData) {
	id.In = format.HumanB(uint64(ifdata.InBytes))
	id.Out = format.HumanB(uint64(ifdata.OutBytes))
}

func (_ interfaceBytes) Delta(id *types.Interface, ii, ifdata getifaddrs.IfData) {
	id.DeltaIn = format.Bps(8, ii.InBytes, ifdata.InBytes)
	id.DeltaOut = format.Bps(8, ii.OutBytes, ifdata.OutBytes)
}

type interfaceInoutErrors struct{}

func (_ interfaceInoutErrors) InOut(ifdata getifaddrs.IfData) (uint, uint) {
	return ifdata.InErrors, ifdata.OutErrors
}

type interfaceInoutPackets struct{}

func (_ interfaceInoutPackets) InOut(ifdata getifaddrs.IfData) (uint, uint) {
	return ifdata.InPackets, ifdata.OutPackets
}

type interfaceNumericals struct{ interfaceInout }

func (ie interfaceNumericals) Current(id *types.Interface, ifdata getifaddrs.IfData) {
	in, out := ie.InOut(ifdata)
	id.In = format.HumanUnitless(uint64(in))
	id.Out = format.HumanUnitless(uint64(out))
}

func (ie interfaceNumericals) Delta(id *types.Interface, ii, previousIfdata getifaddrs.IfData) {
	in, out := ie.InOut(ii)
	previousIn, previousOut := ie.InOut(previousIfdata)
	id.DeltaIn = format.Ps(in, previousIn)
	id.DeltaOut = format.Ps(out, previousOut)
}

func interfacesDelta(format interfaceFormat, current, previous []getifaddrs.IfData, client client.Client) *types.Interfaces {
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
		if !*client.ExpandIF && len(ifs) > client.Toprows {
			ifs = ifs[:client.Toprows]
		}
	}
	ni := new(types.Interfaces)
	ni.List = ifs
	return ni
}

func (li lastinfo) MEM(client client.Client) *types.MEM {
	mem := new(types.MEM)
	mem.RawRAM = li.RAM
	mem.List = append(mem.List, li.RAM.Memory)
	if !*client.HideSWAP {
		mem.List = append(mem.List, li.Swap)
	}
	return mem
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
	return bimap.DefaultSeq
}

var TooltipableTemplate *templates.BinTemplate

func tooltipable(limit int, full string) template.HTML {
	html := "ERROR"
	if len(full) > limit {
		short := full[:limit]
		if TooltipableTemplate == nil {
			log.Printf("tooltipableTemplate hasn't been set")
		} else if buf, err := TooltipableTemplate.CloneExecute(struct {
			Full, Short string
		}{
			Full:  full,
			Short: short,
		}); err == nil {
			html = buf.String()
		}
	} else {
		html = template.HTMLEscapeString(full)
	}
	return template.HTML(html)
}

func orderDisks(disks []diskInfo, seq types.SEQ) []diskInfo {
	if len(disks) > 1 {
		sort.Stable(diskOrder{
			disks:   disks,
			seq:     seq,
			reverse: client.DFBIMAP.SEQ2REVERSE[seq],
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

func dfbytes(diskinfos []diskInfo, client client.Client) *types.DFbytes {
	var disks []types.DiskBytes
	for i, disk := range diskinfos {
		if !*client.ExpandDF && i > client.Toprows-1 {
			break
		}
		total, approxtotal, _ := format.HumanBandback(disk.Total)
		used, approxused, _ := format.HumanBandback(disk.Used)
		disks = append(disks, types.DiskBytes{
			DiskMeta:        diskMeta(disk),
			Total:           total,
			Used:            used,
			Avail:           format.HumanB(disk.Avail),
			UsePercent:      format.FormatPercent(approxused, approxtotal),
			UsePercentClass: format.LabelClassColorPercent(format.Percent(approxused, approxtotal)),
		})
	}
	dsb := new(types.DFbytes)
	dsb.List = disks
	return dsb
}

func dfinodes(diskinfos []diskInfo, client client.Client) *types.DFinodes {
	var disks []types.DiskInodes
	for i, disk := range diskinfos {
		if !*client.ExpandDF && i > client.Toprows-1 {
			break
		}
		itotal, approxitotal, _ := format.HumanBandback(disk.Inodes)
		iused, approxiused, _ := format.HumanBandback(disk.Iused)
		disks = append(disks, types.DiskInodes{
			DiskMeta:         diskMeta(disk),
			Inodes:           itotal,
			Iused:            iused,
			Ifree:            format.HumanB(disk.Ifree),
			IusePercent:      format.FormatPercent(approxiused, approxitotal),
			IusePercentClass: format.LabelClassColorPercent(format.Percent(approxiused, approxitotal)),
		})
	}
	dsi := new(types.DFinodes)
	dsi.List = disks
	return dsi
}

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

func orderProc(procs []types.ProcInfo, cl *client.Client, send *client.SendClient) []types.ProcData {
	if len(procs) > 1 {
		sort.Sort(procOrder{ // not sort.Stable
			procs:   procs,
			seq:     cl.PSSEQ,
			reverse: client.PSBIMAP.SEQ2REVERSE[cl.PSSEQ],
		})
	}

	limitPS := cl.PSlimit
	notdec := limitPS <= 1
	notexp := limitPS >= len(procs)

	if limitPS >= len(procs) { // notexp
		limitPS = len(procs) // NB modified limitPS
	} else {
		procs = procs[:limitPS]
	}

	client.SetBool(&cl.PSnotDecreasable, &send.PSnotDecreasable, notdec)
	client.SetBool(&cl.PSnotExpandable, &send.PSnotExpandable, notexp)
	client.SetString(&cl.PSplusText, &send.PSplusText, fmt.Sprintf("%d+", limitPS))

	uids := map[uint]string{}
	var list []types.ProcData
	for _, proc := range procs {
		list = append(list, types.ProcData{
			PID:      proc.PID,
			Priority: proc.Priority,
			Nice:     proc.Nice,
			Time:     format.FormatTime(proc.Time),
			NameHTML: tooltipable(42, proc.Name),
			UserHTML: tooltipable(12, username(uids, proc.UID)),
			Size:     format.HumanB(proc.Size),
			Resident: format.HumanB(proc.Resident),
		})
	}
	return list
}

type Previous struct {
	CPU        *sigar.CpuList
	Interfaces []getifaddrs.IfData
}

type last struct {
	lastinfo
	mutex sync.Mutex
}

type lastinfo struct {
	Generic    generic
	CPU        cpu.CPUData
	RAM        types.RAM
	Swap       types.Memory
	DiskList   []diskInfo
	ProcList   []types.ProcInfo
	Interfaces []getifaddrs.IfData
	Previous   *Previous
	lastfive   lastfive
}

type lastfive struct {
	// CPU []*fiveCPU
	milliLA1 *five
}

type fiveCPU struct {
	user, sys, idle *five
}

type five struct {
	*ring.Ring
	min, max int
}

func newFive() *five {
	return &five{Ring: ring.New(5), min: -1, max: -1}
}

func (f *five) push(v int) {
	push(&f, v)
}

func push(ff **five, v int) {
	if *ff == nil {
		*ff = newFive()
	}
	f := *ff
	setmin := f.min == -1 || v < f.min
	setmax := f.max == -1 || v > f.max
	if setmin {
		f.min = v
	}
	if setmax {
		f.max = v
	}

	if f.Len() != 0 {
		prev := f.Prev().Value
		if prev != nil {
			// Don't push if the bars for the current and previous are equal

			i, _, e1 := f.bar(prev.(int))
			j, _, e2 := f.bar(v)
			if e1 == nil && e2 == nil && i == j {
				return
			}
		}
	}

	r := f.Move(1)
	r.Move(4).Value = v
	f.Ring = r // gc please

	// recalc min, max of the remained values

	if !setmin {
		if f.Ring != nil && f.Ring.Value != nil {
			f.min = f.Ring.Value.(int)
		}
		f.Do(func(o interface{}) {
			if o == nil {
				return
			}
			v := o.(int)
			if f.min > v {
				f.min = v
			}
		})
	}
	if !setmax {
		if f.Ring != nil && f.Ring.Value != nil {
			f.max = f.Ring.Value.(int)
		}
		f.Do(func(o interface{}) {
			if o == nil {
				return
			}
			v := o.(int)
			if f.max < v {
				f.max = v
			}
		})
	}
}

var bARS = []string{
	"▁",
	"▂",
	"▃",
	// "▄", // looks bad in browsers
	"▅",
	"▆",
	"▇",
	// "█", // looks bad in browsers
}

func (f five) bar(v int) (int, string, error) {
	if f.max == -1 || f.min == -1 { // || f.max == f.min {
		return -1, "", errors.New("Unknown min or max")
	}
	spread := f.max - f.min

	fi := 0.0
	if spread != 0 {
		// fi = float64(v-f.min) / float64(spread)
		fi = float64(f.round(v)-float64(f.min)) / float64(spread)
		if fi > 1.0 {
			// panic("impossible") // ??
			fi = 1.0
		}
	}
	i := int(round(fi * float64(len(bARS)-1)))
	return i, bARS[i], nil
}

func (f five) round(v int) float64 {
	unit := float64(f.max-f.min) /* spread */ / float64(len(bARS)-1)
	times := round((float64(v) - float64(f.min)) / unit)
	return float64(f.min) + unit*times
}

func round(val float64) float64 {
	_, d := math.Modf(val)
	return map[bool]func(float64) float64{true: math.Ceil, false: math.Floor}[d >= 0.5](val)
}

func (f five) spark() string {
	if f.max == -1 || f.min == -1 { // || f.max == f.min {
		return ""
	}

	s := ""
	f.Do(func(o interface{}) {
		if o == nil {
			return
		}
		if _, c, err := f.bar(o.(int)); err == nil {
			s += c
		}
	})
	return s
}

type IndexData struct {
	Generic generic
	CPU     cpu.CPUInfo
	MEM     types.MEM

	PStable PStable
	PSlinks *PSlinks `json:",omitempty"`

	DFlinks  *DFlinks       `json:",omitempty"`
	DFbytes  types.DFbytes  `json:",omitempty"`
	DFinodes types.DFinodes `json:",omitempty"`

	IFbytes   types.Interfaces
	IFerrors  types.Interfaces
	IFpackets types.Interfaces

	VagrantMachines *vagrantMachines
	VagrantError    string
	VagrantErrord   bool

	DISTRIB        string
	VERSION        string
	PeriodDuration types.Duration // default refresh value for placeholder

	Client client.Client

	IFTABS client.IFtabs
	DFTABS client.DFtabs
}

type indexUpdate struct {
	Generic  *generic        `json:",omitempty"`
	CPU      *cpu.CPUInfo    `json:",omitempty"`
	MEM      *types.MEM      `json:",omitempty"`
	DFlinks  *DFlinks        `json:",omitempty"`
	DFbytes  *types.DFbytes  `json:",omitempty"`
	DFinodes *types.DFinodes `json:",omitempty"`
	PSlinks  *PSlinks        `json:",omitempty"`
	PStable  *PStable        `json:",omitempty"`

	IFbytes   *types.Interfaces `json:",omitempty"`
	IFerrors  *types.Interfaces `json:",omitempty"`
	IFpackets *types.Interfaces `json:",omitempty"`

	VagrantMachines *vagrantMachines `json:",omitempty"`
	VagrantError    string
	VagrantErrord   bool

	Client *client.SendClient `json:",omitempty"`
}

var lastInfo last

func (la *last) reset_prev() {
	la.mutex.Lock()
	defer la.mutex.Unlock()

	if la.Previous == nil {
		return
	}
	la.Previous.CPU = nil
	la.Previous.Interfaces = []getifaddrs.IfData{}
}

func (la *last) collect() {
	gch := make(chan generic, 1)
	rch := make(chan types.RAM, 1)
	sch := make(chan types.Memory, 1)
	cch := make(chan cpu.CPUData, 1)
	dch := make(chan []diskInfo, 1)
	pch := make(chan []types.ProcInfo, 1)
	ifch := make(chan IfInfo, 1)

	go getRAM(rch)
	go getSwap(sch)
	go getGeneric(gch)
	go read_disks(dch)
	go read_procs(pch)
	go newInterfaces(ifch)

	func() {
		la.mutex.Lock()
		defer la.mutex.Unlock()

		var prevcl *sigar.CpuList
		if la.Previous != nil {
			prevcl = la.Previous.CPU
		}
		go cpu.CollectCPU(cch, prevcl)
	}()

	la.mutex.Lock()
	defer la.mutex.Unlock()

	// NB .mutex unchanged
	la.lastinfo = lastinfo{
		lastfive: la.lastfive,
		Previous: &Previous{
			CPU:        la.CPU.SigarList(),
			Interfaces: la.Interfaces,
		},
		Generic:  <-gch,
		RAM:      <-rch,
		Swap:     <-sch,
		CPU:      <-cch,
		DiskList: <-dch,
		ProcList: <-pch,
	}

	ii := <-ifch
	la.Generic.IP = ii.IP
	la.Interfaces = ii.List

	push(&la.lastfive.milliLA1, int(float64(100)*la.Generic.LoadAverage.One))
	la.Generic.LA1spark = la.lastfive.milliLA1.spark()

	/* delta, isdelta := la.cpuListDelta()
	for i, core := range delta.List {
		var fcpu *fiveCPU
		if i >= len(la.lastfive.CPU) {
			fcpu = &fiveCPU{
				user: newFive(),
				sys:  newFive(),
				idle: newFive(),
			}
			la.lastfive.CPU = append(la.lastfive.CPU, fcpu)
		} else {
			fcpu = la.lastfive.CPU[i]
		}
		if isdelta {
			_ = core
			fcpu.user.push(int(core.User))
			fcpu.sys .push(int(core.Sys))
			fcpu.idle.push(int(core.Idle))
		}
	} // */
}

func linkattrs(req *http.Request, base url.Values, pname string, bimap types.Biseqmap, seq *types.SEQ) *types.Linkattrs {
	*seq = valuesSet(req, base, pname, bimap)
	return &types.Linkattrs{
		Base:  base,
		Pname: pname,
		Bimap: bimap,
	}
}

func getUpdates(req *http.Request, cl *client.Client, send client.SendClient, forcerefresh bool) indexUpdate {

	cl.RecalcRows() // before anything

	var (
		coreno      int
		df_copy     []diskInfo
		ps_copy     []types.ProcInfo
		if_copy     []getifaddrs.IfData
		previf_copy []getifaddrs.IfData
	)
	iu := indexUpdate{}
	func() {
		lastInfo.mutex.Lock()
		defer lastInfo.mutex.Unlock()

		df_copy = make([]diskInfo, len(lastInfo.DiskList))
		ps_copy = make([]types.ProcInfo, len(lastInfo.ProcList))
		if_copy = make([]getifaddrs.IfData, len(lastInfo.Interfaces))

		copy(df_copy, lastInfo.DiskList)
		copy(ps_copy, lastInfo.ProcList)
		copy(if_copy, lastInfo.Interfaces)

		if lastInfo.lastinfo.Previous != nil {
			previf_copy = make([]getifaddrs.IfData, len(lastInfo.Previous.Interfaces))
			copy(previf_copy, lastInfo.Previous.Interfaces)
		}

		if true { // cl.RefreshGeneric.Refresh(forcerefresh)
			g := lastInfo.Generic
			g.LA = g.LA1spark + " " + g.LA
			iu.Generic = &g // &lastInfo.Generic
		}
		if !*cl.HideMEM && cl.RefreshMEM.Refresh(forcerefresh) {
			iu.MEM = lastInfo.MEM(*cl)
		}
		if !*cl.HideCPU && cl.RefreshCPU.Refresh(forcerefresh) {
			iu.CPU, coreno = lastInfo.CPU.CPUInfo(*cl)
		}
	}()

	if req != nil {
		req.ParseForm() // do ParseForm even if req.Form == nil, otherwise *links won't be set for index requests without parameters
		base := url.Values{}
		iu.PSlinks = (*PSlinks)(linkattrs(req, base, "ps", client.PSBIMAP, &cl.PSSEQ))
		iu.DFlinks = (*DFlinks)(linkattrs(req, base, "df", client.DFBIMAP, &cl.DFSEQ))
	}

	if iu.CPU != nil { // TODO Is it ok to update the *cl.Expand*CPU when the CPU is shown only?
		client.SetBool(&cl.ExpandableCPU, &send.ExpandableCPU, coreno > cl.Toprows-1) // one row reserved for "all N"
		client.SetString(&cl.ExpandtextCPU, &send.ExpandtextCPU, fmt.Sprintf("Expanded (%d)", coreno))
	}

	if true {
		client.SetBool(&cl.ExpandableIF, &send.ExpandableIF, len(if_copy) > cl.Toprows)
		client.SetString(&cl.ExpandtextIF, &send.ExpandtextIF, fmt.Sprintf("Expanded (%d)", len(if_copy)))

		client.SetBool(&cl.ExpandableDF, &send.ExpandableDF, len(df_copy) > cl.Toprows)
		client.SetString(&cl.ExpandtextDF, &send.ExpandtextDF, fmt.Sprintf("Expanded (%d)", len(df_copy)))
	}

	if !*cl.HideDF && cl.RefreshDF.Refresh(forcerefresh) {
		orderedDisks := orderDisks(df_copy, cl.DFSEQ)

		if *cl.TabDF == client.DFBYTES_TABID {
			iu.DFbytes = dfbytes(orderedDisks, *cl)
		} else if *cl.TabDF == client.DFINODES_TABID {
			iu.DFinodes = dfinodes(orderedDisks, *cl)
		}
	}

	if !*cl.HideIF && cl.RefreshIF.Refresh(forcerefresh) {
		switch *cl.TabIF {
		case client.IFBYTES_TABID:
			iu.IFbytes = interfacesDelta(interfaceBytes{}, if_copy, previf_copy, *cl)
		case client.IFERRORS_TABID:
			iu.IFerrors = interfacesDelta(interfaceNumericals{interfaceInoutErrors{}}, if_copy, previf_copy, *cl)
		case client.IFPACKETS_TABID:
			iu.IFpackets = interfacesDelta(interfaceNumericals{interfaceInoutPackets{}}, if_copy, previf_copy, *cl)
		}
	}

	if !*cl.HidePS && cl.RefreshPS.Refresh(forcerefresh) {
		iu.PStable = new(PStable)
		iu.PStable.List = orderProc(ps_copy, cl, &send)
	}

	if !*cl.HideVG && cl.RefreshVG.Refresh(forcerefresh) {
		machines, err := vagrantmachines()
		if err != nil {
			iu.VagrantError = err.Error()
			iu.VagrantErrord = true
		} else {
			iu.VagrantMachines = machines
			iu.VagrantErrord = false
		}
	}

	if send != (client.SendClient{}) {
		iu.Client = &send
	}
	return iu
}

func indexData(minrefresh types.Duration, req *http.Request) IndexData {
	if Connections.Len() == 0 {
		// collect when there're no active connections, so Loop does not collect
		lastInfo.collect()
	}

	cl := client.DefaultClient(minrefresh)
	updates := getUpdates(req, &cl, client.SendClient{}, true)

	data := IndexData{
		Client:  cl,
		Generic: *updates.Generic,
		CPU:     *updates.CPU,
		MEM:     *updates.MEM,

		DFlinks: updates.DFlinks,
		PSlinks: updates.PSlinks,

		PStable: *updates.PStable,

		DISTRIB: DISTRIB, // value set in init()
		VERSION: VERSION, // value from server.go

		PeriodDuration: minrefresh, // default refresh value for placeholder
	}

	if updates.DFbytes != nil {
		data.DFbytes = *updates.DFbytes
	} else if updates.DFinodes != nil {
		data.DFinodes = *updates.DFinodes
	}

	if updates.IFbytes != nil {
		data.IFbytes = *updates.IFbytes
	} else if updates.IFerrors != nil {
		data.IFerrors = *updates.IFerrors
	} else if updates.IFpackets != nil {
		data.IFpackets = *updates.IFpackets
	}
	data.VagrantMachines = updates.VagrantMachines
	data.VagrantError = updates.VagrantError
	data.VagrantErrord = updates.VagrantErrord

	data.DFTABS = client.DFTABS // const
	data.IFTABS = client.IFTABS // const

	return data
}

func statusLine(status int) string {
	return fmt.Sprintf("%d %s", status, http.StatusText(status))
}

func init() {
	DISTRIB = getDistrib()
}

var DISTRIB string

func fqscripts(list []string, r *http.Request) (scripts []string) {
	for _, s := range list {
		if !strings.HasPrefix(string(s), "//") {
			s = "//" + r.Host + s
		}
		scripts = append(scripts, s)
	}
	return scripts
}

func IndexFunc(template templates.BinTemplate, scripts []string, minrefresh types.Duration) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		index(template, scripts, minrefresh, w, r)
	}
}

func index(template templates.BinTemplate, scripts []string, minrefresh types.Duration, w http.ResponseWriter, r *http.Request) {
	response := template.Response(w, struct {
		Data      IndexData
		SCRIPTS   []string
		CLASSNAME string
	}{
		Data:    indexData(minrefresh, r),
		SCRIPTS: fqscripts(scripts, r),
	})
	response.SetHeader("Content-Type", "text/html")
	response.SetContentLength()
	response.Send()
}
