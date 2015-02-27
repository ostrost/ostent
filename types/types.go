//go:generate gen

// Package types has some types of ostent library.
package types

import (
	"container/ring"
	"errors"
	"html/template"
	"math"
	"net/http"
	"net/url"
	"sync"

	metrics "github.com/rcrowley/go-metrics"
	sigar "github.com/rzab/gosigar"
)

// SEQ is a distinct int type for consts and other uses.
type SEQ int

// SeqNReverse holds a SEQ and a Reverse bool.
type SeqNReverse struct {
	SEQ     SEQ
	Reverse bool
}

// AnyOf returns true if the seq is present in the list.
func (seq SEQ) AnyOf(list []SEQ) bool {
	for _, s := range list {
		if s == seq {
			return true
		}
	}
	return false
}

// Sign is a logical operator, useful for sorting.
func (seq SEQ) Sign(t bool) bool { // used in sortable_*.go
	if seq < 0 {
		return t
	}
	return !t
}

// Memory type is a struct of memory metrics.
type Memory struct {
	Kind           string
	Total          string
	Used           string
	Free           string
	UsePercentHTML template.HTML
}

// MEM type has a list of Memory.
type MEM struct {
	List []Memory
}

type RAM struct {
	Memory
	Extra1 uint64 // linux:buffered // darwin:wired
	Extra2 uint64 // linux:cached   // darwin:active
}

// DiskMeta type has common for DiskBytes and DiskInodes fields.
type DiskMeta struct {
	DiskNameHTML template.HTML
	DirNameHTML  template.HTML
	DirNameKey   string
	DevName      string `json:"-"`
}

// DiskBytes type is a struct of disk bytes metrics.
type DiskBytes struct {
	DiskMeta
	Total           string // with units
	Used            string // with units
	Avail           string // with units
	UsePercent      string // as a string, with "%"
	UsePercentClass string
}

// DiskInodes type is a struct of disk inodes metrics.
type DiskInodes struct {
	DiskMeta
	Inodes           string // with units
	Iused            string // with units
	Ifree            string // with units
	IusePercent      string // as a string, with "%"
	IusePercentClass string
}

// DFbytes type has a list of DiskBytes.
type DFbytes struct {
	List []DiskBytes
}

// DFinodes type has a list of DiskInodes.
type DFinodes struct {
	List []DiskInodes
}

// type DiskTable struct {
// 	List  []DiskData
// 	Links *DiskLinkattrs `json:",omitempty"`
// 	HaveCollapsed bool
// }

// Attr type keeps link attributes.
type Attr struct {
	Href, Class, CaretClass string
}

// Attr returns a seq applied Attr taking the la link and updating/setting the parameter.
func (la Linkattrs) Attr(seq SEQ) Attr {
	base := url.Values{}
	for k, v := range la.Base {
		base[k] = v
	}
	attr := Attr{Class: "state"}
	if ascp := la._attr(base, seq); ascp != nil {
		attr.CaretClass = "caret"
		attr.Class += " current"
		if *ascp {
			attr.Class += " dropup"
		}
	}
	attr.Href = "?" + base.Encode() // la._attr modifies base, DO NOT use prior to the call
	return attr
}

// _attr side effect: modifies the base
func (la Linkattrs) _attr(base url.Values, seq SEQ) *bool {
	unlessreverse := func(t bool) *bool {
		if la.Bimap.SEQ2REVERSE[seq] {
			t = !t
		}
		return &t
	}

	if la.Pname == "" {
		if seq == la.Bimap.DefaultSeq {
			return unlessreverse(false)
		}
		return nil
	}

	seqstring := la.Bimap.SEQ2STRING[seq]
	values, haveParam := base[la.Pname]
	base.Set(la.Pname, seqstring)

	if !haveParam { // no parameter in url
		if seq == la.Bimap.DefaultSeq {
			return unlessreverse(false)
		}
		return nil
	}

	pos, neg := values[0], values[0]
	if neg[0] == '-' {
		pos = neg[1:]
		neg = neg[1:]
	} else {
		neg = "-" + neg
	}

	var ascr *bool
	if pos == seqstring {
		t := neg[0] != '-'
		if seq == la.Bimap.DefaultSeq {
			t = true
		}
		ascr = unlessreverse(t)
		base.Set(la.Pname, neg)
	}
	if seq == la.Bimap.DefaultSeq {
		base.Del(la.Pname)
	}
	return ascr
}

// Linkattrs type for link making.
type Linkattrs struct {
	Base  url.Values
	Pname string
	Bimap Biseqmap
}

func valuesSet(req *http.Request, base url.Values, pname string, bimap Biseqmap) SEQ {
	if params, ok := req.Form[pname]; ok && len(params) > 0 {
		if seq, ok := bimap.STRING2SEQ[params[0]]; ok {
			base.Set(pname, params[0])
			return seq
		}
	}
	return bimap.DefaultSeq
}

func NewLinkAttrs(req *http.Request, base url.Values, pname string, bimap Biseqmap, seq *SEQ) *Linkattrs {
	*seq = valuesSet(req, base, pname, bimap)
	return &Linkattrs{
		Base:  base,
		Pname: pname,
		Bimap: bimap,
	}
}

// InterfaceMeta type has common Interface fields.
type InterfaceMeta struct {
	NameKey  string
	NameHTML template.HTML
}

// Interface type is a struct of interface metrics.
type Interface struct {
	InterfaceMeta
	In       string // with units
	Out      string // with units
	DeltaIn  string // with units
	DeltaOut string // with units
}

// Interfaces type has a list of Interface.
type Interfaces struct {
	List []Interface
}

// MetricProc hold a pointer to ProcInfo.
// +gen slice:"PkgSortBy"
type MetricProc struct {
	*ProcInfo
}

// ProcInfo type is an internal account of a process.
type ProcInfo struct {
	PID      uint
	Priority int
	Nice     int
	Time     uint64
	Name     string
	UID      uint
	Size     uint64
	Resident uint64
}

// ProcData type is a public (for index context, json marshaling) account of a process.
type ProcData struct {
	PID      uint
	Priority int
	Nice     int
	Time     string
	NameRaw  string
	NameHTML template.HTML
	UserHTML template.HTML
	Size     string // with units
	Resident string // with units
}

type MetricRAMCommon struct {
	Total metrics.Gauge
}

func NewMetricRAMCommon() MetricRAMCommon {
	return MetricRAMCommon{
		Total: metrics.NewRegisteredGauge("memory.memory-total", metrics.NewRegistry()),
	}
}

func (mrc *MetricRAMCommon) UpdateCommon(got sigar.Mem) {
	mrc.Total.Update(int64(got.Total))
}

type GaugeShortLoad struct {
	metrics.GaugeFloat64
	Ring  *ring.Ring
	Min   int
	Max   int
	Mutex sync.Mutex
}

func (gsl *GaugeShortLoad) Update(floatValue float64) {
	gsl.Mutex.Lock()
	defer gsl.Mutex.Unlock()
	gsl.GaugeFloat64.Update(floatValue)
	value := int(float64(100) * floatValue)
	// func push(ff **five, v int)
	setmin := gsl.Min == -1.0 || value < gsl.Min
	setmax := gsl.Max == -1.0 || value > gsl.Max
	if setmin {
		gsl.Min = value
	}
	if setmax {
		gsl.Max = value
	}

	if gsl.Ring.Len() != 0 {
		if prev := gsl.Ring.Prev().Value; prev != nil {
			// Don't push if the bars for the current and previous are equal
			i, _, e1 := gsl.Bar(prev.(int))
			j, _, e2 := gsl.Bar(value)
			if e1 == nil && e2 == nil && i == j {
				return
			}
		}
	}

	ring := gsl.Ring.Move(1)
	ring.Move(4).Value = value
	gsl.Ring = ring // gc please

	// recalc min, max of the remained values

	if !setmin {
		if gsl.Ring != nil && gsl.Ring.Value != nil {
			gsl.Min = gsl.Ring.Value.(int)
		}
		gsl.Ring.Do(func(o interface{}) {
			if o == nil {
				return
			}
			if v := o.(int); gsl.Min > v {
				gsl.Min = v
			}
		})
	}
	if !setmax {
		if gsl.Ring != nil && gsl.Ring.Value != nil {
			gsl.Max = gsl.Ring.Value.(int)
		}
		gsl.Ring.Do(func(o interface{}) {
			if o == nil {
				return
			}
			if v := o.(int); gsl.Max < v {
				gsl.Max = v
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

func (gsl *GaugeShortLoad) Bar(v int) (int, string, error) {
	if gsl.Max == -1 || gsl.Min == -1 { // || f.max == f.min {
		return -1, "", errors.New("Unknown min or max")
	}
	spread := gsl.Max - gsl.Min

	fi := 0.0
	if spread != 0 {
		// fi = float64(v-f.min) / float64(spread)
		fi = float64(gsl.round(v)-float64(gsl.Min)) / float64(spread)
		if fi > 1.0 {
			// panic("impossible") // ??
			fi = 1.0
		}
	}
	i := int(round(fi * float64(len(bARS)-1)))
	return i, bARS[i], nil
}

func (gsl *GaugeShortLoad) round(v int) float64 {
	unit := float64(gsl.Max-gsl.Min) /* spread */ / float64(len(bARS)-1)
	times := round((float64(v) - float64(gsl.Min)) / unit)
	return float64(gsl.Min) + unit*times
}

func round(val float64) float64 {
	_, d := math.Modf(val)
	return map[bool]func(float64) float64{true: math.Ceil, false: math.Floor}[d >= 0.5](val)
}

func (gsl *GaugeShortLoad) Sparkline() string {
	if gsl.Max == -1 || gsl.Min == -1 { // || gsl.Max == gsl.Min {
		return ""
	}
	s := ""
	gsl.Ring.Do(func(o interface{}) {
		if o == nil {
			return
		}
		if _, c, err := gsl.Bar(o.(int)); err == nil {
			s += c
		}
	})
	return s
}

type MetricLoad struct {
	Short GaugeShortLoad
	Mid   metrics.GaugeFloat64
	Long  metrics.GaugeFloat64
}

func NewMetricLoad(r metrics.Registry) *MetricLoad {
	short := GaugeShortLoad{
		GaugeFloat64: metrics.NewGaugeFloat64(),
		Ring:         ring.New(5), // 5 values
		Min:          -1.0,
		Max:          -1.0,
	}
	// short := metrics.NewRegisteredGaugeFloat64("load.shortterm", r)
	r.Register("load.shortterm", short.GaugeFloat64)
	return &MetricLoad{
		Short: short,
		Mid:   metrics.NewRegisteredGaugeFloat64("load.midterm", r),
		Long:  metrics.NewRegisteredGaugeFloat64("load.longterm", r),
	}
}

type MetricSwap struct {
	Free metrics.Gauge
	Used metrics.Gauge
}

func NewMetricSwap(r metrics.Registry) MetricSwap {
	return MetricSwap{
		Free: metrics.NewRegisteredGauge("swap.swap-free", r),
		Used: metrics.NewRegisteredGauge("swap.swap-used", r),
	}
}

func (ms *MetricSwap) TotalValue() uint64 { // Free + Used
	return uint64(ms.Free.Snapshot().Value() + ms.Used.Snapshot().Value())
}

func (ms *MetricSwap) Update(got sigar.Swap) {
	ms.Free.Update(int64(got.Free))
	ms.Used.Update(int64(got.Used))
}

// GaugeDiff holds two Gauge metrics: the first is the exported one.
// Caveat: The exported metric value is 0 initially, not "nan", until updated.
type GaugeDiff struct {
	Delta    metrics.Gauge // Delta as the primary metric.
	Absolute metrics.Gauge // Absolute keeps the absolute value, not exported as it's registered in private registry.
	Previous metrics.Gauge // Previous keeps the previous absolute value, not exported as it's registered in private registry.
	Mutex    sync.Mutex
}

func NewGaugeDiff(name string, r metrics.Registry) *GaugeDiff {
	return &GaugeDiff{
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

func (gd *GaugeDiff) UpdateAbsolute(absolute int64) int64 {
	gd.Mutex.Lock()
	defer gd.Mutex.Unlock()
	previous := gd.Previous.Snapshot().Value()
	gd.Absolute.Update(absolute)
	gd.Previous.Update(absolute)
	if previous == 0 { // do not .Update
		return 0
	}
	if absolute < previous { // counters got reset
		previous = 0
	}
	delta := absolute - previous
	gd.Delta.Update(delta)
	return delta
}

type GaugePercent struct {
	Percent  metrics.GaugeFloat64 // Percent as the primary metric.
	Previous metrics.Gauge
	Mutex    sync.Mutex
}

func NewGaugePercent(name string, r metrics.Registry) *GaugePercent {
	return &GaugePercent{
		Percent:  metrics.NewRegisteredGaugeFloat64(name, r),
		Previous: metrics.NewRegisteredGauge(name+"-previous", metrics.NewRegistry()),
	}
}

func (gp *GaugePercent) UpdatePercent(totalDelta int64, uabsolute uint64) {
	gp.Mutex.Lock()
	defer gp.Mutex.Unlock()
	previous := gp.Previous.Snapshot().Value()
	absolute := int64(uabsolute)
	gp.Previous.Update(absolute)
	if previous != 0 /* otherwise do not update */ &&
		absolute >= previous /* otherwise counters got reset */ &&
		totalDelta != 0 /* otherwise there were no previous value for Total */ {
		percent := float64(100) * float64(absolute-previous) / float64(totalDelta) // TODO rounding good?
		if percent > 100.0 {
			percent = 100.0
		}
		gp.Percent.Update(percent)
	}
}

type MetricCPUCommon struct {
	metrics.Healthcheck        // derive from one of (go-)metric types, otherwise it won't be registered
	N                   string // The "cpu-N"
	User                *GaugePercent
	Nice                *GaugePercent
	Sys                 *GaugePercent
	Idle                *GaugePercent
	Total               *GaugeDiff
}

func (mcc *MetricCPUCommon) UpdateCommon(sigarCpu sigar.Cpu, total uint64) int64 {
	totalDelta := mcc.Total.UpdateAbsolute(int64(total))
	mcc.User.UpdatePercent(totalDelta, sigarCpu.User)
	mcc.Nice.UpdatePercent(totalDelta, sigarCpu.Nice)
	mcc.Sys.UpdatePercent(totalDelta, sigarCpu.Sys)
	mcc.Idle.UpdatePercent(totalDelta, sigarCpu.Idle)
	return totalDelta
}

func NewMetricCPUCommon(r metrics.Registry, name string) *MetricCPUCommon {
	return &MetricCPUCommon{
		N:     name,
		User:  NewGaugePercent(name+".user", r),
		Nice:  NewGaugePercent(name+".nice", r),
		Sys:   NewGaugePercent(name+".system", r),
		Idle:  NewGaugePercent(name+".idle", r),
		Total: NewGaugeDiff(name+"-total", metrics.NewRegistry()),
	}
}

type MetricString interface {
	Snapshot() MetricString
	Value() string
	Update(string)
}

type StandardMetricString struct {
	string
	Mutex sync.Mutex
}

type MetricStringSnapshot StandardMetricString

func (mss *MetricStringSnapshot) Snapshot() MetricString { return mss }
func (mss *MetricStringSnapshot) Value() string          { return mss.string }
func (*MetricStringSnapshot) Update(string)              { panic("Update called on a MetricStringSnapshot") }

func (sms *StandardMetricString) Snapshot() MetricString { return ((*MetricStringSnapshot)(sms)) }
func (sms *StandardMetricString) Value() string {
	sms.Mutex.Lock()
	defer sms.Mutex.Unlock()
	return sms.string
}

func (sms *StandardMetricString) Update(new string) {
	sms.Mutex.Lock()
	defer sms.Mutex.Unlock()
	sms.string = new
}

// +gen slice:"PkgSortBy"
type MetricDF struct {
	*DF
}

type DF struct {
	metrics.Healthcheck // derive from one of (go-)metric types, otherwise it won't be registered
	DevName             MetricString
	Free                metrics.GaugeFloat64
	Reserved            metrics.GaugeFloat64
	Total               metrics.Gauge
	Used                metrics.GaugeFloat64
	Avail               metrics.Gauge
	UsePercent          metrics.GaugeFloat64
	Inodes              metrics.Gauge
	Iused               metrics.Gauge
	Ifree               metrics.Gauge
	IusePercent         metrics.GaugeFloat64
	DirName             MetricString
}

// Update reads usage and fs and updates the corresponding fields in DF.
func (md *DF) Update(fs sigar.FileSystem, usage sigar.FileSystemUsage) {
	md.DevName.Update(fs.DevName)
	md.DirName.Update(fs.DirName)
	md.Free.Update(float64(usage.Free << 10))
	md.Reserved.Update(float64((usage.Free - usage.Avail) << 10))
	md.Total.Update(int64(usage.Total << 10))
	md.Used.Update(float64(usage.Used << 10))
	md.Avail.Update(int64(usage.Avail << 10))
	md.UsePercent.Update(usage.UsePercent())
	md.Inodes.Update(int64(usage.Files))
	md.Iused.Update(int64(usage.Files - usage.FreeFiles))
	md.Ifree.Update(int64(usage.FreeFiles))
	if iusePercent := 0.0; usage.Files != 0 {
		iusePercent = float64(100) * float64(usage.Files-usage.FreeFiles) / float64(usage.Files)
		md.IusePercent.Update(iusePercent)
	}
}
