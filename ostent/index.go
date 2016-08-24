// Package ostent is the library part of ostent cmd.
package ostent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/user"
	"sort"
	"strings"
	"sync"
	"time"

	metrics "github.com/rcrowley/go-metrics"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"

	"github.com/ostrost/ostent/flags"
	"github.com/ostrost/ostent/format"
	"github.com/ostrost/ostent/internal/plugins/outputs/ostent"
	"github.com/ostrost/ostent/params"
	"github.com/ostrost/ostent/system"
	"github.com/ostrost/ostent/templateutil"
)

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

func (pss PSSlice) Ordered(para *params.Params) *PS {
	para.Psn.Limit = len(pss)
	limitPS := para.Psn.Absolute
	if limitPS > para.Psn.Limit {
		limitPS = para.Psn.Limit
	}

	ps := &PS{}
	ps.N = new(int)
	*ps.N = limitPS

	if para.Psn.Absolute == 0 {
		return ps
	}

	uids := map[uint]string{}
	sort.Sort(PSSort{ // not .Stable
		Psk:     &para.Psk,
		PSSlice: pss,
		UIDs:    uids,
	})

	for _, p := range pss[:limitPS] {
		ps.List = append(ps.List, system.PSData{
			PID:      p.PID,
			UID:      p.UID,
			Priority: p.Priority,
			Nice:     p.Nice,
			Time:     format.Time(p.Time),
			Name:     p.Name,
			User:     username(uids, p.UID),
			Size:     format.HumanB(p.Size),
			Resident: format.HumanB(p.Resident),
		})
	}
	return ps
}

type PS struct {
	List []system.PSData `json:",omitempty"`
	N    *int            `json:",omitempty"`
}

// IndexData is a data map for templates and marshalling.
// Keys (even abbrevs eg CPU) intentionally start with lowercase.
type IndexData map[string]interface{}

type last struct {
	MU            sync.Mutex
	PSSlice       PSSlice
	LastCollected time.Time
	olds          map[string]oldFuncs
}

type oldFuncs struct {
	collectFunc func(*sync.WaitGroup)
	dataFunc    renderFunc
}

type renderFunc func(*params.Params) interface{}

func (la *last) collect(when time.Time, wantprocs bool) {
	la.MU.Lock()
	defer la.MU.Unlock()
	if !la.LastCollected.IsZero() && la.LastCollected.Add(time.Second).After(when) {
		return
	}
	la.LastCollected = when

	var wg sync.WaitGroup
	wg.Add(len(la.olds))
	for _, funcs := range la.olds {
		go funcs.collectFunc(&wg)
	}

	if wantprocs {
		pch := make(chan PSSlice, 1)
		go collectPS(pch)
		la.PSSlice = <-pch
	}
	wg.Wait()
}

func (la *last) CopyPS() PSSlice {
	la.MU.Lock()
	defer la.MU.Unlock()
	dup := make(PSSlice, len(la.PSSlice))
	copy(dup, la.PSSlice)
	return dup
}

// IFSlice is a list of MetricIF.
type IFSlice []*system.MetricIF

// Len, Swap and Less satisfy sorting interface.
func (is IFSlice) Len() int      { return len(is) }
func (is IFSlice) Swap(i, j int) { is[i], is[j] = is[j], is[i] }
func (is IFSlice) Less(i, j int) bool {
	a, b := is[i], is[j]
	amatch := RXlo.Match([]byte(a.Name))
	bmatch := RXlo.Match([]byte(b.Name))
	if !(amatch && bmatch) {
		if amatch {
			return false
		} else if bmatch {
			return true
		}
	}
	return a.Name < b.Name
}

func FormatIF(mi *system.MetricIF) system.IFData {
	ii := system.IFData{
		Name:             mi.Name,
		IP:               mi.IP.Snapshot().Value(),
		DeltaBytesOutNum: mi.BytesOut.DeltaValue(),
	}
	FormatIF1024(mi.BytesIn, &ii.BytesIn, &ii.DeltaBitsIn)
	FormatIF1024(mi.BytesOut, &ii.BytesOut, &ii.DeltaBitsOut)
	FormatIF1000(mi.DropsIn, &ii.DropsIn, &ii.DeltaDropsIn)
	FormatIF1000(mi.DropsOut, &ii.DropsOut, &ii.DeltaDropsOut)
	FormatIF1000(mi.ErrorsIn, &ii.ErrorsIn, &ii.DeltaErrorsIn)
	FormatIF1000(mi.ErrorsOut, &ii.ErrorsOut, &ii.DeltaErrorsOut)
	FormatIF1000(mi.PacketsIn, &ii.PacketsIn, &ii.DeltaPacketsIn)
	FormatIF1000(mi.PacketsOut, &ii.PacketsOut, &ii.DeltaPacketsOut)
	return ii
}

func FormatIF1024(diff *system.GaugeDiff, info1, info2 *string) {
	delta, abs := diff.Values()
	*info1 = format.HumanB(uint64(abs))
	*info2 = format.HumanBits(uint64(8 * delta))
}

func FormatIF1000(diff *system.GaugeDiff, info1, info2 *string) {
	delta, abs := diff.Values()
	*info1 = format.HumanUnitless(uint64(abs))
	*info2 = format.HumanUnitless(uint64(delta))
}

func (ir *IndexRegistry) GetIF(para *params.Params) []system.IFData {
	private := ir.ListPrivateIF()
	para.Ifn.Limit = private.Len()

	sort.Sort(private) // not .Stable

	var public []system.IFData
	for i, mi := range private {
		if i >= para.Ifn.Absolute {
			break
		}
		public = append(public, FormatIF(mi))
	}
	return public
}

// ListPrivateIF returns list of MetricIF's by traversing the PrivateIFRegistry.
func (ir *IndexRegistry) ListPrivateIF() (is IFSlice) {
	ir.PrivateIFRegistry.Each(func(name string, i interface{}) {
		is = append(is, i.(*system.MetricIF))
	})
	return is
}

// GetOrRegisterPrivateIF produces a registered in PrivateIFRegistry system.MetricIF.
func (ir *IndexRegistry) GetOrRegisterPrivateIF(name string) *system.MetricIF {
	ir.PrivateMutex.Lock()
	defer ir.PrivateMutex.Unlock()
	if metric := ir.PrivateIFRegistry.Get(name); metric != nil {
		return metric.(*system.MetricIF)
	}
	i := system.NewMetricIF(ir.Registry, name)
	ir.PrivateIFRegistry.Register(name, i) // error is ignored
	// errs when the type is not derived from (go-)metrics types
	return i
}

func (ir *IndexRegistry) GetOrRegisterPrivateDF(part disk.PartitionStat) *system.MetricDF {
	ir.PrivateMutex.Lock()
	defer ir.PrivateMutex.Unlock()
	if part.Mountpoint == "/" {
		part.Device = "root"
	} else {
		part.Device = strings.Replace(strings.TrimPrefix(part.Device, "/dev/"), "/", "-", -1)
	}
	if metric := ir.PrivateDFRegistry.Get(part.Device); metric != nil {
		return metric.(*system.MetricDF)
	}
	label := func(tail string) string {
		return fmt.Sprintf("df-%s.df_complex-%s", part.Device, tail)
	}
	r, unusedr := ir.Registry, metrics.NewRegistry()
	i := &system.MetricDF{
		DevName:  &system.StandardMetricString{}, // unregistered
		DirName:  &system.StandardMetricString{}, // unregistered
		Free:     metrics.NewRegisteredGaugeFloat64(label("free"), r),
		Reserved: metrics.NewRegisteredGaugeFloat64(label("reserved"), r),
		Total:    metrics.NewRegisteredGauge(label("total"), unusedr),
		Used:     metrics.NewRegisteredGaugeFloat64(label("used"), r),
		Avail:    metrics.NewRegisteredGauge(label("avail"), unusedr),
		UsePct:   metrics.NewRegisteredGaugeFloat64(label("usepercent"), unusedr),
		Inodes:   metrics.NewRegisteredGauge(label("inodes"), unusedr),
		Iused:    metrics.NewRegisteredGauge(label("iused"), unusedr),
		Ifree:    metrics.NewRegisteredGauge(label("ifree"), unusedr),
		IusePct:  metrics.NewRegisteredGaugeFloat64(label("iusepercent"), unusedr),
	}
	ir.PrivateDFRegistry.Register(part.Device, i) // error is ignored
	// errs when the type is not derived from (go-)metrics types
	return i
}

// CPUSlice is a list of MetricCPU.
type CPUSlice []*system.MetricCPU

// Len, Swap and Less satisfy sorting interface.
func (cs CPUSlice) Len() int      { return len(cs) }
func (cs CPUSlice) Swap(i, j int) { cs[i], cs[j] = cs[j], cs[i] }
func (cs CPUSlice) Less(i, j int) bool {
	aidle := cs[i].IdlePct.Percent.Snapshot().Value()
	bidle := cs[j].IdlePct.Percent.Snapshot().Value()
	return aidle < bidle
}

func (ir *IndexRegistry) dataDF(para *params.Params) interface{} {
	if !para.Dfd.Expired() {
		return nil
	}
	return &system.DF{List: ir.GetDF(para)}
}

func (ir *IndexRegistry) GetDF(para *params.Params) []system.DFData {
	private := ir.ListPrivateDF()
	para.Dfn.Limit = len(private)

	sort.Stable(DFSort{
		Dfk:     &para.Dfk,
		DFSlice: private,
	})

	var public []system.DFData
	for i, df := range private {
		if i >= para.Dfn.Absolute {
			break
		}
		public = append(public, FormatDF(df))
	}
	return public
}

func FormatDF(md *system.MetricDF) system.DFData {
	var (
		vdevname = md.DevName.Snapshot().Value()
		vdirname = md.DirName.Snapshot().Value()
		vinodes  = md.Inodes.Snapshot().Value()
		viused   = md.Iused.Snapshot().Value()
		vifree   = md.Ifree.Snapshot().Value()
		vtotal   = md.Total.Snapshot().Value()
		vused    = md.Used.Snapshot().Value()
		vavail   = md.Avail.Snapshot().Value()
	)
	itotal, approxitotal, _ := format.HumanBandback(uint64(vinodes))
	iused, approxiused, _ := format.HumanBandback(uint64(viused))
	total, approxtotal, _ := format.HumanBandback(uint64(vtotal))
	used, approxused, _ := format.HumanBandback(uint64(vused))
	return system.DFData{
		DevName: vdevname,
		DirName: vdirname,
		Inodes:  itotal,
		Iused:   iused,
		Ifree:   format.HumanB(uint64(vifree)),
		IusePct: format.Percent(approxiused, approxitotal),
		Total:   total,
		Used:    used,
		Avail:   format.HumanB(uint64(vavail)),
		UsePct:  format.Percent(approxused, approxtotal),
	}
}

// PSSlice is a list of PSInfo.
type PSSlice []*system.PSInfo

func (pss PSSlice) data(para *params.Params) interface{} {
	if !para.Psd.Expired() {
		return nil
	}
	return pss.Ordered(para)
}

func (ir *IndexRegistry) dataIF(para *params.Params) interface{} {
	if !para.Ifd.Expired() {
		return nil
	}
	return &system.IF{List: ir.GetIF(para)}
}

func (ir *IndexRegistry) dataCPU(para *params.Params) interface{} {
	if !para.CPUd.Expired() {
		return nil
	}
	if para.CPUn.Absolute == 0 {
		para.CPUn.Limit = 1
		return nil
	}
	return &system.CPU{List: ir.GetCPU(para)}
}

func (ir *IndexRegistry) GetCPU(para *params.Params) []system.CPUData {
	private := ir.ListPrivateCPU()

	if private.Len() == 1 {
		para.CPUn.Limit = 1
		return []system.CPUData{FormatCPU("", private[0])}
	}
	para.CPUn.Limit = private.Len() + 1
	sort.Sort(private)

	allabel := fmt.Sprintf("%d cores", private.Len())
	public := []system.CPUData{FormatCPU(allabel, ir.PrivateCPUAll)} // first: "N cores"

	for i, mc := range private {
		if i >= para.CPUn.Absolute-1 {
			break
		}
		public = append(public, FormatCPU("", mc))
	}
	return public
}

func FormatCPU(label string, mc *system.MetricCPU) system.CPUData {
	if label == "" { // A non-"N cores" mc.
		label = "#" + strings.TrimPrefix(mc.N, "cpu-")
	}
	return system.CPUData{
		N:       label,
		UserPct: mc.UserPct.SnapshotValueUint(),
		SysPct:  mc.SysPct.SnapshotValueUint(),
		WaitPct: mc.WaitPct.SnapshotValueUint(),
		IdlePct: mc.IdlePct.SnapshotValueUint(),
	}
}

// ListPrivateCPU returns list of system.MetricCPU's by traversing the PrivateCPURegistry.
func (ir *IndexRegistry) ListPrivateCPU() (cs CPUSlice) {
	ir.PrivateCPURegistry.Each(func(name string, i interface{}) {
		cs = append(cs, i.(*system.MetricCPU))
	})
	return cs
}

// DFSlice is a list of MetricDF.
type DFSlice []*system.MetricDF

// ListPrivateDF returns list of system.MetricDF's by traversing the PrivateDFRegistry.
func (ir *IndexRegistry) ListPrivateDF() (dfs DFSlice) {
	ir.PrivateDFRegistry.Each(func(name string, i interface{}) {
		dfs = append(dfs, i.(*system.MetricDF))
	})
	return dfs
}

// GetOrRegisterPrivateCPU produces a registered in PrivateCPURegistry MetricCPU.
func (ir *IndexRegistry) GetOrRegisterPrivateCPU(coreno int) *system.MetricCPU {
	ir.PrivateMutex.Lock()
	defer ir.PrivateMutex.Unlock()
	name := fmt.Sprintf("cpu-%d", coreno)
	if metric := ir.PrivateCPURegistry.Get(name); metric != nil {
		return metric.(*system.MetricCPU)
	}
	i := system.NewMetricCPU(ir.Registry, name) // of type *system.MetricCPU
	ir.PrivateCPURegistry.Register(name, i)     // error is ignored
	// errs when the type is not derived from (go-)metrics types
	return i
}

type memoryValues struct{ Total, Free, Used uint64 }

func newMemory(kind string, in memoryValues) system.Memory {
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

func (ir *IndexRegistry) dataMEM(para *params.Params) interface{} {
	if !para.Memd.Expired() {
		return nil
	}
	para.Memn.Limit = 2
	if para.Memn.Absolute < 1 {
		return nil
	}
	var (
		mtotal = ir.RAM.Total.Snapshot().Value()
		mfree  = ir.RAM.Free.Snapshot().Value()
	)
	mem := new(system.MEM)
	mem.List = []system.Memory{}
	mem.List = append(mem.List,
		newMemory("RAM", memoryValues{
			Total: uint64(mtotal),
			Free:  uint64(mfree),
			Used:  uint64(mtotal - mfree),
		}))

	if para.Memn.Absolute < 2 {
		return mem
	}
	var (
		sfree = ir.Swap.Free.Snapshot().Value()
		sused = ir.Swap.Used.Snapshot().Value()
	)
	mem.List = append(mem.List,
		newMemory("swap", memoryValues{
			Total: uint64(sfree + sused),
			Free:  uint64(sfree),
			Used:  uint64(sused),
		}))
	return mem
}

func (ir *IndexRegistry) dataLA(para *params.Params) interface{} {
	if !para.Lad.Expired() {
		return nil
	}
	if para.Lan.Absolute < 1 {
		return nil
	}
	para.Lan.Limit = 3
	if para.Lan.Absolute > para.Lan.Limit {
		para.Lan.Absolute = para.Lan.Limit
	}
	type LA struct {
		Period, Value string
	}
	return &struct{ List []LA }{[]LA{
		{"1", fmt.Sprintf("%.2f", ir.Load.Short.Snapshot().Value())},
		{"5", fmt.Sprintf("%.2f", ir.Load.Mid.Snapshot().Value())},
		{"15", fmt.Sprintf("%.2f", ir.Load.Long.Snapshot().Value())},
	}[:para.Lan.Absolute]}
}

func (ir *IndexRegistry) UpdateDF(part disk.PartitionStat, usage *disk.UsageStat) {
	ir.Mutex.Lock()
	defer ir.Mutex.Unlock()
	ir.GetOrRegisterPrivateDF(part).Update(part, usage)
}

// UpdateMEM reads stat and updates ir.RAM and ir.Swap
func (ir *IndexRegistry) UpdateMEM(ram *mem.VirtualMemoryStat, swap *mem.SwapMemoryStat) {
	ir.Mutex.Lock()
	defer ir.Mutex.Unlock()
	ir.RAM.Update(ram)
	ir.Swap.Update(swap)
}

func (ir *IndexRegistry) UpdateLA(stat load.AvgStat) {
	ir.Mutex.Lock()
	defer ir.Mutex.Unlock()
	ir.Load.Short.Update(stat.Load1)
	ir.Load.Mid.Update(stat.Load5)
	ir.Load.Long.Update(stat.Load15)
}

func (ir *IndexRegistry) UpdateCPU(agg cpu.TimesStat, list []cpu.TimesStat) {
	ir.Mutex.Lock()
	defer ir.Mutex.Unlock()
	ir.PrivateCPUAll.Update(agg)
	for coreno, cpu := range list {
		ir.GetOrRegisterPrivateCPU(coreno).Update(cpu)
	}
}

func (ir *IndexRegistry) UpdateIF(ifaddr system.IfAddress) {
	ir.Mutex.Lock()
	defer ir.Mutex.Unlock()
	ir.GetOrRegisterPrivateIF(ifaddr.GetName()).Update(ifaddr)
}

type IndexRegistry struct {
	Registry           metrics.Registry
	PrivateCPUAll      *system.MetricCPU /// metrics.Registry
	PrivateCPURegistry metrics.Registry  // set of MetricCPUs is handled as a metric in this registry
	PrivateIFRegistry  metrics.Registry  // set of system.MetricIFs is handled as a metric in this registry
	PrivateDFRegistry  metrics.Registry  // set of system.MetricDFs is handled as a metric in this registry
	PrivateMutex       sync.Mutex

	RAM  *system.MetricRAM
	Swap system.MetricSwap
	Load *system.MetricLoad

	Mutex sync.Mutex
}

type UpgradeInfo struct {
	RWMutex       sync.RWMutex
	LatestVersion string
}

func (ui *UpgradeInfo) Set(lv string) {
	ui.RWMutex.Lock()
	defer ui.RWMutex.Unlock()
	ui.LatestVersion = lv
}

func (ui *UpgradeInfo) Get() string {
	ui.RWMutex.RLock()
	s := ui.LatestVersion
	ui.RWMutex.RUnlock()
	if s == "" {
		return ""
	}
	return s + " release available"
}

var (
	lastInfo      last
	news          map[string]renderFunc
	OstentUpgrade = new(UpgradeInfo)
	Reg1s         *IndexRegistry
)

func init() {
	reg := metrics.NewRegistry()
	Reg1s = &IndexRegistry{
		Registry: reg,
		PrivateCPUAll: system.NewMetricCPU(reg, // metrics.NewRegistry(),
			"cpu"),
		PrivateCPURegistry: metrics.NewRegistry(),
		PrivateDFRegistry:  metrics.NewRegistry(),
		PrivateIFRegistry:  metrics.NewRegistry(),
		Load:               system.NewMetricLoad(reg),
		Swap:               system.NewMetricSwap(reg),
		RAM:                system.NewMetricRAM(reg),
	}

	lastInfo.olds = map[string]oldFuncs{
	/* if a key is commented out (or missing from predefined set),
	   func Updates may fill data[key] with a ostent.Output.Copy* */

	// "cpu":   {Reg1s.collectCPU, Reg1s.dataCPU},
	// "df":    {Reg1s.collectDF, Reg1s.dataDF},
	// "la":    {Reg1s.collectLA, Reg1s.dataLA},
	// "mem":   {Reg1s.collectMEM, Reg1s.dataMEM},
	// "netio": {Reg1s.collectIF, Reg1s.dataIF},
	}
	news = map[string]renderFunc{
		"cpu": ostent.Output.CopyCPU,
		"df":  ostent.Output.CopyDisk,
		// "la" is copied with ostent.Output.CopySO
		"mem":   ostent.Output.CopyMem,
		"netio": ostent.Output.CopyNet,
	}

	news["procs"] = ostent.Output.CopyProc
}

func Updates(req *http.Request, para *params.Params) (IndexData, bool, error) {
	data := IndexData{}
	if req != nil {
		if err := para.Decode(req); err != nil {
			return data, false, err
		}
		// data features "params" only when req is not nil (new request).
		// So updaters do not read data for it, but expect non-nil para as an argument.
		data["params"] = para
	}

	lastInfo.collect(NextSecond(), para.NonZeroPsn())
	psCopy := lastInfo.CopyPS()

	var updated bool
	if value := psCopy.data(para); value != nil {
		data["procs"] = value
		updated = true
	}

	for key, funcs := range lastInfo.olds {
		if value := funcs.dataFunc(para); value != nil {
			data[key] = value
			updated = true
		}
	}

	for key, dataFunc := range news {
		if _, ok := lastInfo.olds[key]; ok {
			continue
		}
		data[key] = dataFunc(para)
		updated = true
	}
	sodup, sola := ostent.Output.CopySO(para)
	data["system_ostent"] = sodup
	if _, ok := lastInfo.olds["la"]; !ok {
		data["la"] = sola
	}
	return data, updated, nil
}

type ServeSSE struct {
	logRequests bool
	flags.DelayBounds
}

type ServeWS struct {
	ServeSSE
	logger logger
}

type ServeIndex struct {
	ServeWS
	StaticData
	IndexTemplate *templateutil.LazyTemplate
}

type StaticData struct {
	TAGGEDbin     bool
	Distrib       string
	OstentVersion string
}

func NewServeSSE(logRequests bool, dbounds flags.DelayBounds) ServeSSE {
	return ServeSSE{logRequests: logRequests, DelayBounds: dbounds}
}

func NewServeWS(se ServeSSE, lg logger) ServeWS { return ServeWS{ServeSSE: se, logger: lg} }

func NewServeIndex(sw ServeWS, template *templateutil.LazyTemplate, sd StaticData) ServeIndex {
	return ServeIndex{ServeWS: sw, StaticData: sd, IndexTemplate: template}
}

// Index renders index page.
func (si ServeIndex) Index(w http.ResponseWriter, r *http.Request) {
	para := params.NewParams(si.DelayBounds)
	data, _, err := Updates(r, para)
	if err != nil {
		if _, ok := err.(params.RenamedConstError); ok {
			http.Redirect(w, r, err.Error(), http.StatusFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	si.IndexTemplate.Apply(w, struct {
		StaticData
		OstentUpgrade string
		Exporting     ExportingList
		Data          IndexData
	}{
		StaticData:    si.StaticData,
		OstentUpgrade: OstentUpgrade.Get(),
		Exporting:     Exporting, // from ./ws.go
		Data:          data,
	})
}

type SSE struct {
	Writer      http.ResponseWriter // points to the writer
	Params      *params.Params
	SentHeaders bool
	Errord      bool
}

// ServeHTTP is a regular serve func except the first argument,
// passed as a copy, is unused. sse.Writer is there for writes.
func (sse *SSE) ServeHTTP(_ http.ResponseWriter, r *http.Request) {
	w := sse.Writer
	data, _, err := Updates(r, sse.Params)
	if err != nil {
		if _, ok := err.(params.RenamedConstError); ok {
			http.Redirect(w, r, err.Error(), http.StatusFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	text, err := json.Marshal(data)
	if err != nil {
		sse.Errord = true
		// what would http.Error do
		if sse.SetHeader("Content-Type", "text/plain; charset=utf-8") {
			w.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Fprintln(w, err.Error())
		return
	}
	sse.SetHeader("Content-Type", "text/event-stream")
	if _, err := w.Write(append(append([]byte("data: "), text...), []byte("\n\n")...)); err != nil {
		sse.Errord = true
	}
}

func (sse *SSE) SetHeader(name, value string) bool {
	if sse.SentHeaders {
		return false
	}
	sse.SentHeaders = true
	sse.Writer.Header().Set(name, value)
	return true
}

// IndexSSE serves SSE updates.
func (ss ServeSSE) IndexSSE(w http.ResponseWriter, r *http.Request) {
	sse := &SSE{Writer: w, Params: params.NewParams(ss.DelayBounds)}
	if LogHandler(ss.logRequests, sse).ServeHTTP(nil, r); sse.Errord {
		return
	}
	for { // loop is log-requests-free
		_, sleep := NextSecondDelta()
		time.Sleep(sleep)
		if sse.ServeHTTP(nil, r); sse.Errord {
			break
		}
	}
}
