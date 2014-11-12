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
	metrics "github.com/rcrowley/go-metrics"
	sigar "github.com/rzab/gosigar"
)

func interfaceMeta(ifdata getifaddrs.IfData) types.InterfaceMeta {
	return interfaceMetaFromString(ifdata.Name)
}

func interfaceMetaFromString(name string) types.InterfaceMeta {
	return types.InterfaceMeta{
		NameKey:  name,
		NameHTML: tooltipable(12, name),
	}
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
		DevName:      disk.DevName,
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
			RawTotal:        disk.Total,
			RawUsed:         disk.Used,
			RawAvail:        disk.Avail,
		})
	}
	return &types.DFbytes{List: disks}
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
	return &types.DFinodes{List: disks}
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
	CPU *sigar.CpuList
}

type last struct {
	lastinfo
	mutex sync.Mutex
}

type lastinfo struct {
	Generic  generic
	CPU      cpu.CPUData
	DiskList []diskInfo
	ProcList []types.ProcInfo
	Previous *Previous
	lastfive lastfive
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
}

func (la *last) collect() {
	gch := make(chan generic, 1)
	cch := make(chan cpu.CPUData, 1)
	dch := make(chan []diskInfo, 1)
	pch := make(chan []types.ProcInfo, 1)
	ifch := make(chan string, 1)

	var wg sync.WaitGroup
	wg.Add(2) // two so far
	go getRAM(&wg)
	go getSwap(&wg)

	go getGeneric(gch)
	go read_disks(dch)
	go read_procs(pch)
	go getInterfaces(Reg1s, ifch)

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
			CPU: la.CPU.SigarList(),
		},
		Generic:  <-gch,
		CPU:      <-cch,
		DiskList: <-dch,
		ProcList: <-pch,
	}

	la.Generic.IP = <-ifch
	wg.Wait()

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

// GaugeDiff holds two Gauge metrics: the first is the exported one.
// Caveat: The exported metric value is 0 initially, not "nan", until updated.
type GaugeDiff struct {
	Delta    metrics.Gauge // Delta as the primary metric.
	Absolute metrics.Gauge // Absolute keeps the absolute value, not exported as it's registered in private registry.
	Previous metrics.Gauge // Previous keeps the previous absolute value, not exported as it's registered in private registry.
	Mutex    sync.Mutex
}

func NewGaugeDiff(name string, r metrics.Registry) GaugeDiff {
	return GaugeDiff{
		Delta:    metrics.NewRegisteredGauge(name, r),
		Absolute: metrics.NewRegisteredGauge(name+"-absolute", metrics.NewRegistry()),
		Previous: metrics.NewRegisteredGauge(name+"-previous", metrics.NewRegistry()),
	}
}

func (gd *GaugeDiff) Values() (int64, int64) {
	gd.Mutex.Lock()
	defer gd.Mutex.Unlock()
	return gd.Delta.Snapshot().Value(), gd.Absolute.Snapshot().Value()
}

func (gd *GaugeDiff) UpdateAbsolute(absolute int64) {
	gd.Mutex.Lock()
	defer gd.Mutex.Unlock()
	previous := gd.Previous.Snapshot().Value()
	gd.Absolute.Update(absolute)
	gd.Previous.Update(absolute)
	if previous != 0 { // otherwise do not update
		if absolute < previous { // counters got reset
			previous = 0
		}
		gd.Delta.Update(absolute - previous)
	}
}

type ListMetricInterface []MetricInterface  // satisfying sort.Interface
func (x ListMetricInterface) Len() int      { return len(x) }
func (x ListMetricInterface) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x ListMetricInterface) Less(i, j int) bool {
	if rx_lo.Match([]byte(x[i].Name)) {
		return false
	}
	return x[i].Name < x[j].Name
}

type MetricInterface struct {
	metrics.Healthcheck // derive from one of (go-)metric types, otherwise it won't be registered
	Name                string
	BytesIn             GaugeDiff
	BytesOut            GaugeDiff
	ErrorsIn            GaugeDiff
	ErrorsOut           GaugeDiff
	PacketsIn           GaugeDiff
	PacketsOut          GaugeDiff
}

func (mi *MetricInterface) Update(ifdata getifaddrs.IfData) {
	mi.BytesIn.UpdateAbsolute(int64(ifdata.InBytes))
	mi.BytesOut.UpdateAbsolute(int64(ifdata.OutBytes))
	mi.ErrorsIn.UpdateAbsolute(int64(ifdata.InErrors))
	mi.ErrorsOut.UpdateAbsolute(int64(ifdata.OutErrors))
	mi.PacketsIn.UpdateAbsolute(int64(ifdata.InPackets))
	mi.PacketsOut.UpdateAbsolute(int64(ifdata.OutPackets))
}

func (mi MetricInterface) FormatInterface(ip InterfaceParts) types.Interface {
	ing, outg, isbytes := ip(mi)
	deltain, in := ing.Values()
	deltaout, out := outg.Values()
	form := format.HumanUnitless
	deltaForm := format.HumanUnitless // format.Ps
	if isbytes {
		form = format.HumanB
		deltaForm = func(c uint64) string { // , p uint64
			// return format.Bps(8, c, p) // format.Bps64(8, {in,out}, 0)
			return format.HumanBits(c * 8) // passing the bits
		}
	}
	return types.Interface{
		InterfaceMeta: interfaceMetaFromString(mi.Name),
		In:            form(uint64(in)),            // format.HumanB(uint64(in)),  // with units
		Out:           form(uint64(out)),           // format.HumanB(uint64(out)), // with units
		DeltaIn:       deltaForm(uint64(deltain)),  // format.Bps64(8, in, 0),     // with units
		DeltaOut:      deltaForm(uint64(deltaout)), // format.Bps64(8, out, 0),    // with units
	}
}

type InterfaceParts func(MetricInterface) (GaugeDiff, GaugeDiff, bool)

func (_ Registry) InterfaceBytes(mi MetricInterface) (GaugeDiff, GaugeDiff, bool) {
	return mi.BytesIn, mi.BytesOut, true
}
func (_ Registry) InterfaceErrors(mi MetricInterface) (GaugeDiff, GaugeDiff, bool) {
	return mi.ErrorsIn, mi.ErrorsOut, false
}
func (_ Registry) InterfacePackets(mi MetricInterface) (GaugeDiff, GaugeDiff, bool) {
	return mi.PacketsIn, mi.PacketsOut, false
}

func (r Registry) Interfaces(cli *client.Client, send *client.SendClient, ip InterfaceParts) []types.Interface {
	private := r.ListPrivateInterface()
	var public []types.Interface

	client.SetBool(&cli.ExpandableIF, &send.ExpandableIF, len(private) > cli.Toprows)
	client.SetString(&cli.ExpandtextIF, &send.ExpandtextIF, fmt.Sprintf("Expanded (%d)", len(private)))

	if len(private) == 0 || len(private) == 1 {
		return public
	}
	sort.Sort(ListMetricInterface(private))
	// MAYBE have a measure not to make full list in r.ListPrivateInterface
	for i, p := range private {
		if !*cli.ExpandIF && i >= cli.Toprows {
			break
		}
		public = append(public, p.FormatInterface(ip))
	}
	return public
}

func (r *Registry) ListPrivateInterface() []MetricInterface {
	var mi []MetricInterface
	r.PrivateRegistry.Each(func(name string, i interface{}) {
		mi = append(mi, i.(MetricInterface))
	})
	return mi
}

func (r *Registry) GetOrRegisterPrivateInterface(name string) *MetricInterface {
	r.PrivateMutex.Lock()
	defer r.PrivateMutex.Unlock()
	if metric := r.PrivateRegistry.Get(name); metric != nil {
		i := metric.(MetricInterface)
		return &i
	}
	i := MetricInterface{
		Name:       name,
		BytesIn:    NewGaugeDiff("interface-"+name+".if_octets.rx", r.Registry),
		BytesOut:   NewGaugeDiff("interface-"+name+".if_octets.tx", r.Registry),
		ErrorsIn:   NewGaugeDiff("interface-"+name+".if_errors.rx", r.Registry),
		ErrorsOut:  NewGaugeDiff("interface-"+name+".if_errors.tx", r.Registry),
		PacketsIn:  NewGaugeDiff("interface-"+name+".if_packets.rx", r.Registry),
		PacketsOut: NewGaugeDiff("interface-"+name+".if_packets.tx", r.Registry),
	}
	r.PrivateRegistry.Register(name, i) // error is ignored
	// errs when the type is not derived from (go-)metrics types
	return &i
}

func (r Registry) MEM(client client.Client) *types.MEM {
	gr := r.RAM
	mem := new(types.MEM)
	mem.List = []types.Memory{
		_getmem("RAM", sigar.Swap{
			Total: uint64(gr.Total.Snapshot().Value()),
			Free:  uint64(gr.Free.Snapshot().Value()),
			Used:  gr.UsedValue(), // == .Total - .Free
		}),
	}
	if !*client.HideSWAP {
		gs := r.Swap
		mem.List = append(mem.List,
			_getmem("swap", sigar.Swap{
				Total: gs.TotalValue(),
				Free:  uint64(gs.Free.Snapshot().Value()),
				Used:  uint64(gs.Used.Snapshot().Value()),
			}))
	}
	return mem
}

func (r Registry) LA() string {
	gl := r.Load
	return fmt.Sprintf("%.2f %.2f %.2f",
		gl.Short.Snapshot().Value(),
		gl.Mid.Snapshot().Value(),
		gl.Long.Snapshot().Value())
}

type Registry struct {
	Registry        metrics.Registry
	PrivateRegistry metrics.Registry
	PrivateMutex    sync.Mutex

	// set of MetricInterfaces is handled as a metric in PrivateRegistry

	RAM  types.GaugeRAM
	Swap types.GaugeSwap
	Load types.GaugeLoad
}

var Reg1s Registry

func init() {
	reg1s := metrics.NewRegistry()
	Reg1s = Registry{
		Registry:        reg1s,
		PrivateRegistry: metrics.NewRegistry(),
	}
	Reg1s.RAM = types.NewGaugeRAM(Reg1s.Registry)
	Reg1s.Swap = types.NewGaugeSwap(Reg1s.Registry)
	Reg1s.Load = types.NewGaugeLoad(Reg1s.Registry)

	// addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:2003")
	// go metrics.Graphite(reg, 1*time.Second, "ostent", addr)
}

func getUpdates(req *http.Request, cl *client.Client, send client.SendClient, forcerefresh bool) indexUpdate {

	cl.RecalcRows() // before anything

	var (
		coreno  int
		df_copy []diskInfo
		ps_copy []types.ProcInfo
	)
	iu := indexUpdate{}
	func() {
		lastInfo.mutex.Lock()
		defer lastInfo.mutex.Unlock()

		df_copy = make([]diskInfo, len(lastInfo.DiskList))
		ps_copy = make([]types.ProcInfo, len(lastInfo.ProcList))

		copy(df_copy, lastInfo.DiskList)
		copy(ps_copy, lastInfo.ProcList)

		if true { // cl.RefreshGeneric.Refresh(forcerefresh)
			g := lastInfo.Generic
			g.LA = g.LA1spark + " " + Reg1s.LA()
			iu.Generic = &g // &lastInfo.Generic
		}
		if !*cl.HideMEM && cl.RefreshMEM.Refresh(forcerefresh) {
			iu.MEM = Reg1s.MEM(*cl)
		}
		if !*cl.HideCPU && cl.RefreshCPU.Refresh(forcerefresh) {
			iu.CPU, coreno = lastInfo.CPU.CPUInfo(*cl)
		}
	}()

	if req != nil {
		req.ParseForm() // do ParseForm even if req.Form == nil, otherwise *links won't be set for index requests without parameters
		base := url.Values{}
		iu.PSlinks = (*PSlinks)(types.NewLinkAttrs(req, base, "ps", client.PSBIMAP, &cl.PSSEQ))
		iu.DFlinks = (*DFlinks)(types.NewLinkAttrs(req, base, "df", client.DFBIMAP, &cl.DFSEQ))
	}

	if iu.CPU != nil { // TODO Is it ok to update the *cl.Expand*CPU when the CPU is shown only?
		client.SetBool(&cl.ExpandableCPU, &send.ExpandableCPU, coreno > cl.Toprows-1) // one row reserved for "all N"
		client.SetString(&cl.ExpandtextCPU, &send.ExpandtextCPU, fmt.Sprintf("Expanded (%d)", coreno))
	}

	if true {
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
			iu.IFbytes = &types.Interfaces{List: Reg1s.Interfaces(cl, &send, Reg1s.InterfaceBytes)}
		case client.IFERRORS_TABID:
			iu.IFerrors = &types.Interfaces{List: Reg1s.Interfaces(cl, &send, Reg1s.InterfaceErrors)}
		case client.IFPACKETS_TABID:
			iu.IFpackets = &types.Interfaces{List: Reg1s.Interfaces(cl, &send, Reg1s.InterfacePackets)}
		}
	}

	if !*cl.HidePS && cl.RefreshPS.Refresh(forcerefresh) {
		iu.PStable = &PStable{List: orderProc(ps_copy, cl, &send)}
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
