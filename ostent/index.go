// Package ostent is the library part of ostent cmd.
package ostent

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/user"
	"sort"
	"strings"
	"sync"
	"time"

	sigar "github.com/ostrost/gosigar"
	metrics "github.com/rcrowley/go-metrics"

	"github.com/ostrost/ostent/flags"
	"github.com/ostrost/ostent/format"
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
			Time:     format.FormatTime(p.Time),
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

func (data IndexData) SetString(k, v string) {
	if v != "" {
		data[k] = v
	}
}

type last struct {
	MU            sync.Mutex
	PSSlice       PSSlice
	LastCollected time.Time
}

var lastInfo last

func (la *last) collect(when time.Time, wantprocs bool) {
	la.MU.Lock()
	defer la.MU.Unlock()
	if !la.LastCollected.IsZero() && la.LastCollected.Add(time.Second).After(when) {
		return
	}
	la.LastCollected = when

	c := Machine{}
	var wg sync.WaitGroup
	wg.Add(8)              // EIGHT:
	go c.CPU(&Reg1s, &wg)  // one
	go c.RAM(&Reg1s, &wg)  // two
	go c.Swap(&Reg1s, &wg) // three
	go c.DF(&Reg1s, &wg)   // four
	go c.HN(RegMSS, &wg)   // five
	go c.UP(RegMSS, &wg)   // six
	go c.LA(&Reg1s, &wg)   // seven
	go c.IF(&Reg1s, &wg)   // eight

	if wantprocs {
		pch := make(chan PSSlice, 1)
		go c.PS(pch)
		la.PSSlice = <-pch
	}
	wg.Wait()
}

func (la *last) CopyPS() PSSlice {
	la.MU.Lock()
	defer la.MU.Unlock()
	psCopy := make(PSSlice, len(la.PSSlice))
	copy(psCopy, la.PSSlice)
	return psCopy
}

func (mss *MSS) HN(para *params.Params, data IndexData) bool {
	// HN has no delay, always updates data
	data.SetString("hostname", mss.GetString("hostname"))
	return true
}

func (mss *MSS) UP(para *params.Params, data IndexData) bool {
	// UP has no delay, always updates data
	data.SetString("uptime", mss.GetString("uptime"))
	return true
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
	FormatIF1000(mi.ErrorsIn, &ii.ErrorsIn, &ii.DeltaErrorsIn)
	FormatIF1000(mi.ErrorsOut, &ii.ErrorsOut, &ii.DeltaErrorsOut)
	FormatIF1000(mi.PacketsIn, &ii.PacketsIn, &ii.DeltaPacketsIn)
	FormatIF1000(mi.PacketsOut, &ii.PacketsOut, &ii.DeltaPacketsOut)
	if mi.DropsOut != nil {
		FormatIF1000(mi.DropsOut, &ii.DropsOut, &ii.DeltaDropsOut)
	} else {
		ii.DropsOut, ii.DeltaDropsOut = "-1", "-1"
	}
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

func (ir *IndexRegistry) GetOrRegisterPrivateDF(fs sigar.FileSystem) *system.MetricDF {
	ir.PrivateMutex.Lock()
	defer ir.PrivateMutex.Unlock()
	if fs.DirName == "/" {
		fs.DevName = "root"
	} else {
		fs.DevName = strings.Replace(strings.TrimPrefix(fs.DevName, "/dev/"), "/", "-", -1)
	}
	if metric := ir.PrivateDFRegistry.Get(fs.DevName); metric != nil {
		return metric.(*system.MetricDF)
	}
	label := func(tail string) string {
		return fmt.Sprintf("df-%s.df_complex-%s", fs.DevName, tail)
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
	ir.PrivateDFRegistry.Register(fs.DevName, i) // error is ignored
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

func (ir *IndexRegistry) DF(para *params.Params, data IndexData) bool {
	if !para.Dfd.Expired() {
		return false
	}
	data["df"] = &system.DF{List: ir.GetDF(para)}
	return true
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

func (pss PSSlice) IU(para *params.Params, data IndexData) bool {
	if !para.Psd.Expired() {
		return false
	}
	data["procs"] = pss.Ordered(para)
	return true
}

func (ir *IndexRegistry) IF(para *params.Params, data IndexData) bool {
	if !para.Ifd.Expired() {
		return false
	}
	data["netio"] = &system.IF{List: ir.GetIF(para)}
	return true
}

func (ir *IndexRegistry) CPU(para *params.Params, data IndexData) bool {
	if !para.CPUd.Expired() {
		return false
	}
	if para.CPUn.Absolute == 0 {
		para.CPUn.Limit = 1
		return false
	}
	data["cpu"] = &system.CPU{List: ir.GetCPU(para)}
	return true
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

func (ir *IndexRegistry) MEM(para *params.Params, data IndexData) bool {
	if !para.Memd.Expired() {
		return false
	}
	para.Memn.Limit = 2
	if para.Memn.Absolute < 1 {
		return false
	}
	mem := new(system.MEM)
	mem.List = []system.Memory{}
	mem.List = append(mem.List,
		_getmem("RAM", sigar.Swap{
			Total: uint64(ir.RAM.Total.Snapshot().Value()),
			Free:  uint64(ir.RAM.Free.Snapshot().Value()),
			Used:  ir.RAM.UsedValue(), // == .Total - .Free
		}))
	data["mem"] = mem

	if para.Memn.Absolute < 2 {
		return true
	}
	mem.List = append(mem.List,
		_getmem("swap", sigar.Swap{
			Total: ir.Swap.TotalValue(),
			Free:  uint64(ir.Swap.Free.Snapshot().Value()),
			Used:  uint64(ir.Swap.Used.Snapshot().Value()),
		}))
	return true
}

func (ir *IndexRegistry) LA(para *params.Params, data IndexData) bool {
	if !para.Lad.Expired() {
		return false
	}
	if para.Lan.Absolute < 1 {
		return false
	}
	para.Lan.Limit = 3
	if para.Lan.Absolute > para.Lan.Limit {
		para.Lan.Absolute = para.Lan.Limit
	}
	type LA struct {
		Period, Value string
	}
	data["la"] = &struct{ List []LA }{[]LA{
		{"1", fmt.Sprintf("%.2f", ir.Load.Short.Snapshot().Value())},
		{"5", fmt.Sprintf("%.2f", ir.Load.Mid.Snapshot().Value())},
		{"15", fmt.Sprintf("%.2f", ir.Load.Long.Snapshot().Value())},
	}[:para.Lan.Absolute]}
	return true
}

func (ir *IndexRegistry) UpdateDF(fs sigar.FileSystem, usage sigar.FileSystemUsage) {
	ir.Mutex.Lock()
	defer ir.Mutex.Unlock()
	ir.GetOrRegisterPrivateDF(fs).Update(fs, usage)
}

func (ir *IndexRegistry) UpdateRAM(got sigar.Mem, extra1, extra2 uint64) {
	ir.Mutex.Lock()
	defer ir.Mutex.Unlock()
	ir.RAM.Update(got, extra1, extra2)
}

// UpdateSwap reads got and updates the ir.Swap. TODO Bad description.
func (ir *IndexRegistry) UpdateSwap(got sigar.Swap) {
	ir.Mutex.Lock()
	defer ir.Mutex.Unlock()
	ir.Swap.Update(got)
}

func (ir *IndexRegistry) UpdateLA(la sigar.LoadAverage) {
	ir.Mutex.Lock()
	defer ir.Mutex.Unlock()
	ir.Load.Short.Update(la.One)
	ir.Load.Mid.Update(la.Five)
	ir.Load.Long.Update(la.Fifteen)
}

func (ir *IndexRegistry) UpdateCPU(all sigar.Cpu, list []sigar.Cpu) {
	ir.Mutex.Lock()
	defer ir.Mutex.Unlock()
	ir.PrivateCPUAll.Update(all)
	for coreno, cpu := range list {
		ir.GetOrRegisterPrivateCPU(coreno).Update(cpu)
	}
}

func (ir *IndexRegistry) UpdateIF(ifaddr system.IfAddress) {
	ir.Mutex.Lock()
	defer ir.Mutex.Unlock()
	ir.GetOrRegisterPrivateIF(ifaddr.GetName()).Update(ifaddr)
}

// S2SRegistry is for string kv storage.
type S2SRegistry interface {
	SetString(string, string)
	GetString(string) string
}

// MSS implements S2SRegistry in a map[string]string.
type MSS struct {
	MU sync.Mutex
	KV map[string]string
}

func (mss *MSS) SetString(k, v string) {
	mss.MU.Lock()
	defer mss.MU.Unlock()
	mss.KV[k] = v
}

func (mss *MSS) GetString(k string) string {
	mss.MU.Lock()
	defer mss.MU.Unlock()
	return mss.KV[k]
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

var (
	Reg1s  IndexRegistry
	RegMSS = &MSS{KV: map[string]string{}}
)

func init() {
	reg := metrics.NewRegistry()
	Reg1s = IndexRegistry{
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
	for _, update := range []func(*params.Params, IndexData) bool{
		// These are updaters:
		psCopy.IU,
		RegMSS.HN,
		RegMSS.UP,
		Reg1s.MEM,
		Reg1s.CPU,
		Reg1s.DF,
		Reg1s.IF,
		Reg1s.LA,
	} {
		if update(para, data) {
			updated = true
		}
	}
	return data, updated, nil
}

func statusLine(status int) string {
	return fmt.Sprintf("%d %s", status, http.StatusText(status))
}

func init() {
	var err error
	DISTRIB, err = system.Distrib()
	if err != nil {
		log.Printf("WARN %s\n", err)
	}
}

// DISTRIB is distribution string and it's version.
// Set at init, result of system.Distrib.
var DISTRIB string

type ServeSSE struct {
	Access *Access
	flags.DelayBounds
}

type ServeWS struct {
	ServeSSE
	ErrLog *log.Logger
}

type ServeIndex struct {
	ServeWS
	TaggedBin     bool
	IndexTemplate *templateutil.LazyTemplate
}

func NewServeSSE(access *Access, dbounds flags.DelayBounds) ServeSSE {
	return ServeSSE{Access: access, DelayBounds: dbounds}
}

func NewServeWS(ss ServeSSE, errlog *log.Logger) ServeWS {
	return ServeWS{ServeSSE: ss, ErrLog: errlog}
}

func NewServeIndex(sw ServeWS, taggedbin bool, template *templateutil.LazyTemplate) ServeIndex {
	return ServeIndex{ServeWS: sw, TaggedBin: taggedbin, IndexTemplate: template}
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
		TAGGEDbin     bool
		Distrib       string
		OstentVersion string
		Exporting     ExportingList
		Data          IndexData
	}{
		TAGGEDbin:     si.TaggedBin,
		Distrib:       DISTRIB,   // value set in init()
		OstentVersion: VERSION,   // from ./server.go
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
	if ss.Access.Constructor(sse).ServeHTTP(nil, r); sse.Errord { // the request logging
		return
	}
	for { // loop is access-log-free
		_, sleep := NextSecondDelta()
		time.Sleep(sleep)
		if sse.ServeHTTP(nil, r); sse.Errord {
			break
		}
	}
}
