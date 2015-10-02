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
	"github.com/ostrost/ostent/system/operating"
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
		ps.List = append(ps.List, operating.PSData{
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

type IndexData struct {
	Params *params.Params `json:",omitempty"`

	Generic // inline non-pointer

	MEM operating.MEM
	DF  operating.DF
	CPU operating.CPU
	IF  operating.IF
	PS  PS

	DISTRIB string
}

type PS struct {
	List []operating.PSData `json:",omitempty"`
	N    *int               `json:",omitempty"`
}

type IndexUpdate struct {
	Params *params.Params `json:",omitempty"`

	Generic // inline non-pointer

	MEM *operating.MEM `json:",omitempty"`
	DF  *operating.DF  `json:",omitempty"`
	CPU *operating.CPU `json:",omitempty"`
	IF  *operating.IF  `json:",omitempty"`
	PS  *PS            `json:",omitempty"`

	Location *string `json:",omitempty"`
}

type Generic struct {
	HN string `json:",omitempty"`
	UP string `json:",omitempty"`
	LA string `json:",omitempty"`
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

func (mss *MSS) HN(para *params.Params, iu *IndexUpdate) bool {
	// HN has no delay, always updates iu
	iu.HN = mss.GetString("hostname")
	return true
}

func (mss *MSS) UP(para *params.Params, iu *IndexUpdate) bool {
	// UP has no delay, always updates iu
	iu.UP = mss.GetString("uptime")
	return true
}

// IFSlice is a list of MetricIF.
type IFSlice []*operating.MetricIF

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

func FormatIF(mi *operating.MetricIF) operating.IFData {
	ii := operating.IFData{
		Name: mi.Name,
		IP:   mi.IP.Snapshot().Value(),
	}
	FormatIF1024(mi.BytesIn, &ii.BytesIn, &ii.DeltaBitsIn)
	FormatIF1024(mi.BytesOut, &ii.BytesOut, &ii.DeltaBitsOut)
	FormatIF1000(mi.DropsIn, &ii.DropsIn, &ii.DeltaDropsIn)
	FormatIF1000(mi.ErrorsIn, &ii.ErrorsIn, &ii.DeltaErrorsIn)
	FormatIF1000(mi.ErrorsOut, &ii.ErrorsOut, &ii.DeltaErrorsOut)
	FormatIF1000(mi.PacketsIn, &ii.PacketsIn, &ii.DeltaPacketsIn)
	FormatIF1000(mi.PacketsOut, &ii.PacketsOut, &ii.DeltaPacketsOut)
	if mi.DropsOut != nil {
		ii.DropsOut, ii.DeltaDropsOut = new(string), new(string)
		FormatIF1000(mi.DropsOut, ii.DropsOut, ii.DeltaDropsOut)
	}
	return ii
}

func FormatIF1024(diff *operating.GaugeDiff, info1, info2 *string) {
	delta, abs := diff.Values()
	*info1 = format.HumanB(uint64(abs))
	*info2 = format.HumanBits(uint64(8 * delta))
}

func FormatIF1000(diff *operating.GaugeDiff, info1, info2 *string) {
	delta, abs := diff.Values()
	*info1 = format.HumanUnitless(uint64(abs))
	*info2 = format.HumanUnitless(uint64(delta))
}

func (ir *IndexRegistry) GetIF(para *params.Params) []operating.IFData {
	private := ir.ListPrivateIF()
	para.Ifn.Limit = private.Len()

	sort.Sort(private) // not .Stable

	var public []operating.IFData
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
		is = append(is, i.(*operating.MetricIF))
	})
	return is
}

// GetOrRegisterPrivateIF produces a registered in PrivateIFRegistry operating.MetricIF.
func (ir *IndexRegistry) GetOrRegisterPrivateIF(name string) *operating.MetricIF {
	ir.PrivateMutex.Lock()
	defer ir.PrivateMutex.Unlock()
	if metric := ir.PrivateIFRegistry.Get(name); metric != nil {
		return metric.(*operating.MetricIF)
	}
	i := operating.NewMetricIF(ir.Registry, name)
	ir.PrivateIFRegistry.Register(name, i) // error is ignored
	// errs when the type is not derived from (go-)metrics types
	return i
}

func (ir *IndexRegistry) GetOrRegisterPrivateDF(fs sigar.FileSystem) *operating.MetricDF {
	ir.PrivateMutex.Lock()
	defer ir.PrivateMutex.Unlock()
	if fs.DirName == "/" {
		fs.DevName = "root"
	} else {
		fs.DevName = strings.Replace(strings.TrimPrefix(fs.DevName, "/dev/"), "/", "-", -1)
	}
	if metric := ir.PrivateDFRegistry.Get(fs.DevName); metric != nil {
		return metric.(*operating.MetricDF)
	}
	label := func(tail string) string {
		return fmt.Sprintf("df-%s.df_complex-%s", fs.DevName, tail)
	}
	r, unusedr := ir.Registry, metrics.NewRegistry()
	i := &operating.MetricDF{
		DevName:  &operating.StandardMetricString{}, // unregistered
		DirName:  &operating.StandardMetricString{}, // unregistered
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
type CPUSlice []*operating.MetricCPU

// Len, Swap and Less satisfy sorting interface.
func (cs CPUSlice) Len() int      { return len(cs) }
func (cs CPUSlice) Swap(i, j int) { cs[i], cs[j] = cs[j], cs[i] }
func (cs CPUSlice) Less(i, j int) bool {
	aidle := cs[i].IdlePct.Percent.Snapshot().Value()
	bidle := cs[j].IdlePct.Percent.Snapshot().Value()
	return aidle < bidle
}

func (ir *IndexRegistry) DF(para *params.Params, iu *IndexUpdate) bool {
	if !para.Dfd.Expired() {
		return false
	}
	iu.DF = &operating.DF{List: ir.GetDF(para)}
	return true
}

func (ir *IndexRegistry) GetDF(para *params.Params) []operating.DFData {
	private := ir.ListPrivateDF()
	para.Dfn.Limit = len(private)

	sort.Stable(DFSort{
		Dfk:     &para.Dfk,
		DFSlice: private,
	})

	var public []operating.DFData
	for i, df := range private {
		if i >= para.Dfn.Absolute {
			break
		}
		public = append(public, FormatDF(df))
	}
	return public
}

func FormatDF(md *operating.MetricDF) operating.DFData {
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
	return operating.DFData{
		DevName: vdevname,
		DirName: vdirname,
		Inodes:  itotal,
		Iused:   iused,
		Ifree:   format.HumanB(uint64(vifree)),
		IusePct: format.FormatPercent(approxiused, approxitotal),
		Total:   total,
		Used:    used,
		Avail:   format.HumanB(uint64(vavail)),
		UsePct:  format.FormatPercent(approxused, approxtotal),
	}
}

// PSSlice is a list of PSInfo.
type PSSlice []*operating.PSInfo

func (pss PSSlice) IU(para *params.Params, iu *IndexUpdate) bool {
	if !para.Psd.Expired() {
		return false
	}
	iu.PS = pss.Ordered(para)
	return true
}

func (ir *IndexRegistry) IF(para *params.Params, iu *IndexUpdate) bool {
	if !para.Ifd.Expired() {
		return false
	}
	iu.IF = &operating.IF{List: ir.GetIF(para)}
	return true
}

func (ir *IndexRegistry) CPU(para *params.Params, iu *IndexUpdate) bool {
	if !para.CPUd.Expired() {
		return false
	}
	if para.CPUn.Absolute == 0 {
		para.CPUn.Limit = 1
		return false
	}
	iu.CPU = &operating.CPU{List: ir.GetCPU(para)}
	return true
}

func (ir *IndexRegistry) GetCPU(para *params.Params) []operating.CPUData {
	private := ir.ListPrivateCPU()

	if private.Len() == 1 {
		para.CPUn.Limit = 1
		return []operating.CPUData{FormatCPU("", private[0])}
	}
	para.CPUn.Limit = private.Len() + 1
	sort.Sort(private)

	allabel := fmt.Sprintf("all %d", private.Len())
	public := []operating.CPUData{FormatCPU(allabel, ir.PrivateCPUAll)} // first: "all N"

	for i, mc := range private {
		if i >= para.CPUn.Absolute-1 {
			break
		}
		public = append(public, FormatCPU("", mc))
	}
	return public
}

func FormatCPU(label string, mc *operating.MetricCPU) operating.CPUData {
	if label == "" {
		label = "#" + strings.TrimPrefix(mc.N, "cpu-") // A non-"all" mc.
	}
	return operating.CPUData{
		N:       label,
		UserPct: mc.UserPct.SnapshotValueUint(),
		SysPct:  mc.SysPct.SnapshotValueUint(),
		WaitPct: mc.WaitPct.SnapshotValueUint(),
		IdlePct: mc.IdlePct.SnapshotValueUint(),
	}
}

// ListPrivateCPU returns list of operating.MetricCPU's by traversing the PrivateCPURegistry.
func (ir *IndexRegistry) ListPrivateCPU() (cs CPUSlice) {
	ir.PrivateCPURegistry.Each(func(name string, i interface{}) {
		cs = append(cs, i.(*operating.MetricCPU))
	})
	return cs
}

// DFSlice is a list of MetricDF.
type DFSlice []*operating.MetricDF

// ListPrivateDF returns list of operating.MetricDF's by traversing the PrivateDFRegistry.
func (ir *IndexRegistry) ListPrivateDF() (dfs DFSlice) {
	ir.PrivateDFRegistry.Each(func(name string, i interface{}) {
		dfs = append(dfs, i.(*operating.MetricDF))
	})
	return dfs
}

// GetOrRegisterPrivateCPU produces a registered in PrivateCPURegistry MetricCPU.
func (ir *IndexRegistry) GetOrRegisterPrivateCPU(coreno int) *operating.MetricCPU {
	ir.PrivateMutex.Lock()
	defer ir.PrivateMutex.Unlock()
	name := fmt.Sprintf("cpu-%d", coreno)
	if metric := ir.PrivateCPURegistry.Get(name); metric != nil {
		return metric.(*operating.MetricCPU)
	}
	i := system.NewMetricCPU(ir.Registry, name) // of type *operating.MetricCPU
	ir.PrivateCPURegistry.Register(name, i)     // error is ignored
	// errs when the type is not derived from (go-)metrics types
	return i
}

func (ir *IndexRegistry) MEM(para *params.Params, iu *IndexUpdate) bool {
	if !para.Memd.Expired() {
		return false
	}
	para.Memn.Limit = 2
	if para.Memn.Absolute < 1 {
		return false
	}
	iu.MEM = new(operating.MEM)
	iu.MEM.List = []operating.Memory{}
	iu.MEM.List = append(iu.MEM.List,
		_getmem("RAM", sigar.Swap{
			Total: uint64(ir.RAM.Total.Snapshot().Value()),
			Free:  uint64(ir.RAM.Free.Snapshot().Value()),
			Used:  ir.RAM.UsedValue(), // == .Total - .Free
		}))

	if para.Memn.Absolute < 2 {
		return true
	}
	iu.MEM.List = append(iu.MEM.List,
		_getmem("swap", sigar.Swap{
			Total: ir.Swap.TotalValue(),
			Free:  uint64(ir.Swap.Free.Snapshot().Value()),
			Used:  uint64(ir.Swap.Used.Snapshot().Value()),
		}))
	return true
}

func (ir *IndexRegistry) LA(para *params.Params, iu *IndexUpdate) bool {
	// LA has no delay, always updates iu
	iu.LA = fmt.Sprintf("%.2f %.2f %.2f",
		ir.Load.Short.Snapshot().Value(),
		ir.Load.Mid.Snapshot().Value(),
		ir.Load.Long.Snapshot().Value())
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

func (ir *IndexRegistry) UpdateCPU(cpuslice []sigar.Cpu) {
	ir.Mutex.Lock()
	defer ir.Mutex.Unlock()
	all := sigar.Cpu{}
	for coreno, cpu := range cpuslice {
		ir.GetOrRegisterPrivateCPU(coreno).Update(cpu)
		operating.AddSCPU(&all, cpu)
	}
	ir.PrivateCPUAll.Update(all)
}

func (ir *IndexRegistry) UpdateIF(ifaddr operating.IfAddress) {
	ir.Mutex.Lock()
	defer ir.Mutex.Unlock()
	ir.GetOrRegisterPrivateIF(ifaddr.Name()).Update(ifaddr)
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
	PrivateCPUAll      *operating.MetricCPU
	PrivateCPURegistry metrics.Registry // set of MetricCPUs is handled as a metric in this registry
	PrivateIFRegistry  metrics.Registry // set of operating.MetricIFs is handled as a metric in this registry
	PrivateDFRegistry  metrics.Registry // set of operating.MetricDFs is handled as a metric in this registry
	PrivateMutex       sync.Mutex

	RAM  *operating.MetricRAM
	Swap operating.MetricSwap
	Load *operating.MetricLoad

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
		PrivateCPUAll: system.NewMetricCPU(metrics.NewRegistry(),
			"all" /* This "all" never used or referenced by */),
		PrivateCPURegistry: metrics.NewRegistry(),
		PrivateDFRegistry:  metrics.NewRegistry(),
		PrivateIFRegistry:  metrics.NewRegistry(),
		Load:               operating.NewMetricLoad(reg),
		Swap:               operating.NewMetricSwap(reg),
		RAM:                system.NewMetricRAM(reg),
	}
}

func getUpdates(req *http.Request, para *params.Params) (IndexUpdate, bool, error) {
	iu := IndexUpdate{}
	if req != nil {
		err := para.Decode(req)
		if err != nil {
			return iu, false, err
		}
		// iu.Location = newloc // may be nil
		iu.Params = para
	}
	lastInfo.collect(NextSecond(), para.NonZeroPsn())
	psCopy := lastInfo.CopyPS()

	var updated bool
	for _, update := range []func(*params.Params, *IndexUpdate) bool{
		psCopy.IU,
		RegMSS.HN,
		RegMSS.UP,
		Reg1s.MEM,
		Reg1s.CPU,
		Reg1s.DF,
		Reg1s.IF,
		Reg1s.LA,
	} {
		if update(para, &iu) {
			updated = true
		}
	}
	return iu, updated, nil
}

func indexData(mindelay, maxdelay flags.Delay, req *http.Request) (IndexData, error) {
	para := params.NewParams(mindelay, maxdelay)
	updates, _, err := getUpdates(req, para)
	if err != nil {
		return IndexData{}, err
	}

	data := IndexData{
		Params:  updates.Params,
		Generic: updates.Generic,

		DISTRIB: DISTRIB, // value set in init()
	}

	if updates.CPU != nil {
		data.CPU = *updates.CPU
	}
	if updates.MEM != nil {
		data.MEM = *updates.MEM
	}
	if updates.PS != nil {
		data.PS = *updates.PS
	}
	if updates.DF != nil {
		data.DF = *updates.DF
	}
	if updates.IF != nil {
		data.IF = *updates.IF
	}
	return data, nil
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

/*
func FormRedirectFunc(mindelay, maxdelay flags.Delay, wrap func(http.HandlerFunc) http.Handler) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, muxpara httprouter.Params) {
		wrap(func(w http.ResponseWriter, req *http.Request) {
			where := "/"
			if q := muxpara.ByName("Q"); q != "" {
				req.URL.RawQuery = req.Form.Encode() + "&" + strings.TrimPrefix(q, "?")
				req.Form = nil // reset the .Form for .ParseForm() to parse new r.URL.RawQuery.
				para := params.NewParams(mindelay, maxdelay)
				para.Decode(req) // OR err.Error()
				if s, err := para.Encode(); err == nil {
					where = "/?" + s
				}
			}
			http.Redirect(w, req, where, http.StatusFound)
		}).ServeHTTP(w, req)
	}
}
*/

type ServeSSE struct {
	Access   *Access
	MinDelay flags.Delay
}

type ServeWS struct {
	ServeSSE
	ErrLog   *log.Logger
	MaxDelay flags.Delay
}

type ServeIndex struct {
	ServeWS
	TaggedBin     bool
	OstentVersion string
	IndexTemplate *templateutil.LazyTemplate
}

func NewServeSSE(access *Access, mindelay flags.Delay) ServeSSE {
	return ServeSSE{Access: access, MinDelay: mindelay}
}

func NewServeWS(ss ServeSSE, errlog *log.Logger, maxdelay flags.Delay) ServeWS {
	return ServeWS{ServeSSE: ss, ErrLog: errlog, MaxDelay: maxdelay}
}

func NewServeIndex(sw ServeWS, taggedbin bool, template *templateutil.LazyTemplate) ServeIndex {
	return ServeIndex{ServeWS: sw, TaggedBin: taggedbin, IndexTemplate: template, OstentVersion: VERSION /* value from server.go */}
}

// Index renders index page.
func (si ServeIndex) Index(w http.ResponseWriter, r *http.Request) {
	id, err := indexData(si.MinDelay, si.MaxDelay, r)
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
		OstentVersion string
		Data          IndexData
	}{
		TAGGEDbin:     si.TaggedBin,
		OstentVersion: si.OstentVersion,
		Data:          id,
	})
}

type SSE struct {
	Writer      http.ResponseWriter // points to the writer
	MinDelay    flags.Delay
	MaxDelay    flags.Delay
	SentHeaders bool
	Errord      bool
}

// ServeHTTP is a regular serve func except the first argument,
// passed as a copy, is unused. sse.Writer is there for writes.
func (sse *SSE) ServeHTTP(_ http.ResponseWriter, r *http.Request) {
	w := sse.Writer
	id, err := indexData(sse.MinDelay, sse.MaxDelay, r)
	if err != nil {
		if _, ok := err.(params.RenamedConstError); ok {
			http.Redirect(w, r, err.Error(), http.StatusFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	text, err := json.Marshal(id)
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
	sse := &SSE{Writer: w, MinDelay: ss.MinDelay}
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
