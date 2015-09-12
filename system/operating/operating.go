// Package operating (as oppose to system) holds platform-independant code.
package operating

import (
	"sync"

	sigar "github.com/ostrost/gosigar"
	metrics "github.com/rcrowley/go-metrics"
)

// Memory type is a struct of memory metrics.
type Memory struct {
	Kind   string
	Total  string
	Used   string
	Free   string
	UsePct string
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

// DFData type is a struct of disk metrics.
type DFData struct {
	DevName string
	DirName string

	// strings with units

	// inodes
	Inodes  string
	Iused   string
	Ifree   string
	IusePct string // without %

	// bytes
	Total  string
	Used   string
	Avail  string
	UsePct string // without %
}

// DF type has a list of DFData.
type DF struct {
	List []DFData
}

// IFData type is a struct of interface metrics.
type IFData struct {
	Name string
	IP   string

	// strings with units

	BytesIn      string
	BytesOut     string
	DeltaBitsIn  string
	DeltaBitsOut string

	DropsIn       string
	DropsOut      *string `json:",omitempty"`
	DeltaDropsIn  string
	DeltaDropsOut *string `json:",omitempty"`

	ErrorsIn       string
	ErrorsOut      string
	DeltaErrorsIn  string
	DeltaErrorsOut string

	PacketsIn       string
	PacketsOut      string
	DeltaPacketsIn  string
	DeltaPacketsOut string
}

// IF type has a list of IFData.
type IF struct {
	List []IFData
}

// PSInfo type is an internal account of a process.
type PSInfo struct {
	PID      uint
	Priority int
	Nice     int
	Time     uint64
	Name     string
	UID      uint
	Size     uint64
	Resident uint64
}

// PSData type is a public (for index context, json marshaling) account of a process.
type PSData struct {
	PID      uint
	UID      uint
	Priority int
	Nice     int
	Time     string
	Name     string
	User     string
	Size     string // with units
	Resident string // with units
}

type RAMUpdater interface {
	UpdateRAM(sigar.Mem, uint64, uint64)
}

type MetricRAM struct {
	Free  metrics.Gauge
	Total metrics.Gauge
	Extra RAMUpdater
}

func ExtraNewMetricRAM(r metrics.Registry, extra RAMUpdater) *MetricRAM {
	return &MetricRAM{
		Free:  metrics.NewRegisteredGauge("memory.memory-free", r),
		Total: metrics.NewRegisteredGauge("memory.memory-total", metrics.NewRegistry()),
		Extra: extra,
	}
}

func (mr *MetricRAM) Update(got sigar.Mem, extra1, extra2 uint64) {
	mr.Free.Update(int64(got.Free))
	mr.Total.Update(int64(got.Total))
	if mr.Extra != nil {
		mr.Extra.UpdateRAM(got, extra1, extra2)
	}
}

func (mr *MetricRAM) UsedValue() uint64 { // Total - Free
	return uint64(mr.Total.Snapshot().Value() - mr.Free.Snapshot().Value())
}

type MetricLoad struct {
	Short metrics.GaugeFloat64
	Mid   metrics.GaugeFloat64
	Long  metrics.GaugeFloat64
}

func NewMetricLoad(r metrics.Registry) *MetricLoad {
	return &MetricLoad{
		Short: metrics.NewRegisteredGaugeFloat64("load.shortterm", r),
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

func (gp *GaugePercent) SnapshotValueUint() uint {
	return uint(gp.Percent.Snapshot().Value())
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

// CPU type has a list of CPUData.
type CPU struct {
	List []CPUData
}

// CPUData type is a struct of cpu metrics.
type CPUData struct {
	N string
	// Pcts without "%"
	UserPct uint
	SysPct  uint
	WaitPct uint
	IdlePct uint
}

type CPUUpdater interface {
	UpdateCPU(sigar.Cpu, int64)
}

type MetricCPU struct {
	// derive from one of (go-)metric types, otherwise it won't be registered
	metrics.Healthcheck
	N       string // The "cpu-N"
	UserPct *GaugePercent
	NicePct *GaugePercent
	SysPct  *GaugePercent
	WaitPct *GaugePercent
	IdlePct *GaugePercent
	Total   *GaugeDiff
	Extra   CPUUpdater
}

func (mc *MetricCPU) Update(scpu sigar.Cpu) {
	totalDelta := mc.Total.UpdateAbsolute(int64(scpu.Total()))
	mc.UserPct.UpdatePercent(totalDelta, scpu.User)
	mc.NicePct.UpdatePercent(totalDelta, scpu.Nice)
	mc.SysPct.UpdatePercent(totalDelta, scpu.Sys)
	mc.WaitPct.UpdatePercent(totalDelta, scpu.Wait)
	mc.IdlePct.UpdatePercent(totalDelta, scpu.Idle)
	if mc.Extra != nil {
		mc.Extra.UpdateCPU(scpu, totalDelta)
	}
}

func ExtraNewMetricCPU(r metrics.Registry, name string, extra CPUUpdater) *MetricCPU {
	return &MetricCPU{
		N:       name,
		UserPct: NewGaugePercent(name+".user", r),
		NicePct: NewGaugePercent(name+".nice", r),
		SysPct:  NewGaugePercent(name+".system", r),
		WaitPct: NewGaugePercent(name+".wait", r),
		IdlePct: NewGaugePercent(name+".idle", r),
		Total:   NewGaugeDiff(name+"-total", metrics.NewRegistry()),
		Extra:   extra,
	}
}

// AddSCPU adds other to dst field by field.
func AddSCPU(dst *sigar.Cpu, other sigar.Cpu) {
	dst.User += other.User
	dst.Nice += other.Nice
	dst.Sys += other.Sys
	dst.Wait += other.Wait
	dst.Idle += other.Idle
	dst.Irq += other.Irq
	dst.SoftIrq += other.SoftIrq
	dst.Stolen += other.Stolen
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

type MetricStringSnapshot string

func (mss MetricStringSnapshot) Snapshot() MetricString { return mss }
func (mss MetricStringSnapshot) Value() string          { return string(mss) }
func (MetricStringSnapshot) Update(string)              { panic("Update called on a MetricStringSnapshot") }

func (sms *StandardMetricString) Snapshot() MetricString {
	sms.Mutex.Lock()
	defer sms.Mutex.Unlock()
	return MetricStringSnapshot(sms.string)
}

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

type MetricDF struct {
	// derive from one of (go-)metric types, otherwise it won't be registered
	metrics.Healthcheck
	DevName  MetricString
	Free     metrics.GaugeFloat64
	Reserved metrics.GaugeFloat64
	Total    metrics.Gauge
	Used     metrics.GaugeFloat64
	Avail    metrics.Gauge
	UsePct   metrics.GaugeFloat64
	Inodes   metrics.Gauge
	Iused    metrics.Gauge
	Ifree    metrics.Gauge
	IusePct  metrics.GaugeFloat64
	DirName  MetricString
}

// Update reads usage and fs and updates the corresponding fields in DF.
func (md *MetricDF) Update(fs sigar.FileSystem, usage sigar.FileSystemUsage) {
	md.DevName.Update(fs.DevName)
	md.DirName.Update(fs.DirName)
	md.Free.Update(float64(usage.Free << 10))
	md.Reserved.Update(float64((usage.Free - usage.Avail) << 10))
	md.Total.Update(int64(usage.Total << 10))
	md.Used.Update(float64(usage.Used << 10))
	md.Avail.Update(int64(usage.Avail << 10))
	md.UsePct.Update(usage.UsePercent())
	md.Inodes.Update(int64(usage.Files))
	md.Iused.Update(int64(usage.Files - usage.FreeFiles))
	md.Ifree.Update(int64(usage.FreeFiles))
	if iusepct := 0.0; usage.Files != 0 {
		iusepct = float64(100) * float64(usage.Files-usage.FreeFiles) / float64(usage.Files)
		md.IusePct.Update(iusepct)
	}
}

// MetricIF set of interface metrics.
type MetricIF struct {
	NewGD func(string) *GaugeDiff // internal
	// derive from one of (go-)metric types, otherwise it won't be registered
	metrics.Healthcheck
	Name       string
	IP         MetricString
	BytesIn    *GaugeDiff
	BytesOut   *GaugeDiff
	DropsIn    *GaugeDiff
	DropsOut   *GaugeDiff
	ErrorsIn   *GaugeDiff
	ErrorsOut  *GaugeDiff
	PacketsIn  *GaugeDiff
	PacketsOut *GaugeDiff
}

func NewMetricIF(reg metrics.Registry, name string) *MetricIF {
	ngd := func(word string) *GaugeDiff {
		return NewGaugeDiff("interface-"+name+"."+word, reg)
	}
	return &MetricIF{
		NewGD:    ngd,
		Name:     name,
		IP:       &StandardMetricString{},
		BytesIn:  ngd("if_octets.rx"),
		BytesOut: ngd("if_octets.tx"),
		DropsIn:  ngd("if_drops.rx"),
		// DropsOut stays nil
		ErrorsIn:   ngd("if_errors.rx"),
		ErrorsOut:  ngd("if_errors.tx"),
		PacketsIn:  ngd("if_packets.rx"),
		PacketsOut: ngd("if_packets.tx"),
	}
}

// Update reads ifaddr and updates the corresponding fields in MetricIF.
func (mi *MetricIF) Update(ifaddr IfAddress) {
	mi.IP.Update(ifaddr.IP())
	mi.BytesIn.UpdateAbsolute(int64(ifaddr.BytesIn()))
	mi.BytesOut.UpdateAbsolute(int64(ifaddr.BytesOut()))
	mi.DropsIn.UpdateAbsolute(int64(ifaddr.DropsIn()))
	// mi.DropsOut dealt below
	mi.ErrorsIn.UpdateAbsolute(int64(ifaddr.ErrorsIn()))
	mi.ErrorsOut.UpdateAbsolute(int64(ifaddr.ErrorsOut()))
	mi.PacketsIn.UpdateAbsolute(int64(ifaddr.PacketsIn()))
	mi.PacketsOut.UpdateAbsolute(int64(ifaddr.PacketsOut()))
	if do := ifaddr.DropsOut(); do != nil {
		if mi.DropsOut == nil {
			mi.DropsOut = mi.NewGD("if_drops.tx")
		}
		mi.DropsOut.UpdateAbsolute(int64(*do))
	}
}

type IfAddress interface {
	Name() string
	IP() string
	BytesIn() uint
	BytesOut() uint
	DropsIn() uint
	DropsOut() *uint
	ErrorsIn() uint
	ErrorsOut() uint
	PacketsIn() uint
	PacketsOut() uint
}
