//go:generate gen

// Package operating (as oppose to system) holds platform-independant code.
package operating

import (
	"fmt"
	"html/template"
	"sync"

	metrics "github.com/rcrowley/go-metrics"
	sigar "github.com/rzab/gosigar"
)

// Memory type is a struct of memory metrics.
type Memory struct {
	Kind       Field
	Total      string
	Used       string
	Free       string
	UsePercent string
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
	DevName string
	DirName Field
}

// DiskBytes type is a struct of disk bytes metrics.
type DiskBytes struct {
	DiskMeta
	Total      string // with units
	Used       string // with units
	Avail      string // with units
	UsePercent string // without %
}

// DiskInodes type is a struct of disk inodes metrics.
type DiskInodes struct {
	DiskMeta
	Inodes      string // with units
	Iused       string // with units
	Ifree       string // with units
	IusePercent string // without %
}

// DFbytes type has a list of DiskBytes.
type DFbytes struct {
	List []DiskBytes
}

// DFinodes type has a list of DiskInodes.
type DFinodes struct {
	List []DiskInodes
}

// InterfaceInfo type is a struct of interface metrics.
type InterfaceInfo struct {
	Name     Field
	In       string // with units
	Out      string // with units
	DeltaIn  string // with units
	DeltaOut string // with units
}

// Interfaces type has a list of Interface.
type Interfaces struct {
	List []InterfaceInfo
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
	PID      Field
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

// CPUInfo type has a list of CoreInfo.
type CPUInfo struct {
	List          []CoreInfo
	ExpandCPU     *bool   `json:",omitempty"`
	ExpandableCPU *Bool   `json:",omitempty"`
	ExpandtextCPU *string `json:",omitempty"`
}

type Bool bool

func (b Bool) BoolClassAttr(fstclass, sndclass string) template.HTMLAttr {
	s := fstclass
	if !b {
		s = sndclass
	}
	return template.HTMLAttr(fmt.Sprintf(" class=%q", s))
}

// CoreInfo type is a struct of core metrics.
type CoreInfo struct {
	N    Field
	User uint // percent without "%"
	Sys  uint // percent without "%"
	Idle uint // percent without "%"
}

type CPUUpdater interface {
	UpdateCPU(sigar.Cpu, int64)
}

// +gen slice:"PkgSortBy"
type MetricCPU struct {
	*CPU
}

type CPU struct {
	metrics.Healthcheck       // derive from one of (go-)metric types, otherwise it won't be registered
	N                   Field // The "cpu-N"
	User                *GaugePercent
	Nice                *GaugePercent
	Sys                 *GaugePercent
	Idle                *GaugePercent
	Total               *GaugeDiff
	Extra               CPUUpdater
}

func (mc MetricCPU) Update(scpu sigar.Cpu) {
	totalDelta := mc.Total.UpdateAbsolute(int64(scpu.Total()))
	mc.User.UpdatePercent(totalDelta, scpu.User)
	mc.Nice.UpdatePercent(totalDelta, scpu.Nice)
	mc.Sys.UpdatePercent(totalDelta, scpu.Sys)
	mc.Idle.UpdatePercent(totalDelta, scpu.Idle)
	if mc.Extra != nil {
		mc.Extra.UpdateCPU(scpu, totalDelta)
	}
}

func ExtraNewMetricCPU(r metrics.Registry, name string, extra CPUUpdater) *MetricCPU {
	return &MetricCPU{
		CPU: &CPU{
			N:     Field(name),
			User:  NewGaugePercent(name+".user", r),
			Nice:  NewGaugePercent(name+".nice", r),
			Sys:   NewGaugePercent(name+".system", r),
			Idle:  NewGaugePercent(name+".idle", r),
			Total: NewGaugeDiff(name+"-total", metrics.NewRegistry()),
			Extra: extra,
		},
	}
}

// AddSCPU adds other to dst field by field.
func AddSCPU(dst *sigar.Cpu, other sigar.Cpu) {
	dst.User += other.User
	dst.Nice += other.Nice
	dst.Sys += other.Sys
	dst.Idle += other.Idle
	dst.Wait += other.Wait
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

type MetricStringSnapshot StandardMetricString

func (mss *MetricStringSnapshot) Snapshot() MetricString { return mss }
func (mss *MetricStringSnapshot) Value() string          { return mss.string }
func (*MetricStringSnapshot) Update(string)              { panic("Update called on a MetricStringSnapshot") }

func (sms *StandardMetricString) Snapshot() MetricString {
	sms.Mutex.Lock()
	defer sms.Mutex.Unlock()
	return ((*MetricStringSnapshot)(sms))
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

// MetricInterface hold a pointer to Interface.
// +gen slice:"PkgSortBy"
type MetricInterface struct {
	*Interface
}

// Interface is a set of interface metrics.
type Interface struct {
	metrics.Healthcheck // derive from one of (go-)metric types, otherwise it won't be registered
	Name                Field
	BytesIn             *GaugeDiff
	BytesOut            *GaugeDiff
	ErrorsIn            *GaugeDiff
	ErrorsOut           *GaugeDiff
	PacketsIn           *GaugeDiff
	PacketsOut          *GaugeDiff
}

// Update reads ifdata and updates the corresponding fields in Interface.
func (i *Interface) Update(ifdata Getifdata) {
	i.BytesIn.UpdateAbsolute(int64(ifdata.GetInBytes()))
	i.BytesOut.UpdateAbsolute(int64(ifdata.GetOutBytes()))
	i.ErrorsIn.UpdateAbsolute(int64(ifdata.GetInErrors()))
	i.ErrorsOut.UpdateAbsolute(int64(ifdata.GetOutErrors()))
	i.PacketsIn.UpdateAbsolute(int64(ifdata.GetInPackets()))
	i.PacketsOut.UpdateAbsolute(int64(ifdata.GetOutPackets()))
}

type Getifdata interface {
	GetInBytes() uint
	GetOutBytes() uint
	GetInErrors() uint
	GetOutErrors() uint
	GetInPackets() uint
	GetOutPackets() uint
}

// +gen slice:"PkgSortBy"
type Vgmachine struct {
	UUID             Field
	Name             string
	Provider         string
	State            string
	Vagrantfile_path string // Directory
}

// Field is a string type with KeyAttr method.
type Field string

func (f Field) String() string { return string(f) }

// KeyAttr returns empty attribute. In jsx that would be "key=%s".
func (f Field) KeyAttr(string) (empty template.HTMLAttr) { return }
