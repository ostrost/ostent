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

	sigar "github.com/ostrost/gosigar"
	metrics "github.com/rcrowley/go-metrics"

	"github.com/ostrost/ostent/flags"
	"github.com/ostrost/ostent/format"
	"github.com/ostrost/ostent/params"
	"github.com/ostrost/ostent/params/enums"
	"github.com/ostrost/ostent/system"
	"github.com/ostrost/ostent/system/operating"
	"github.com/ostrost/ostent/templateutil"
)

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

func diskMeta(disk operating.MetricDF) operating.DiskMeta {
	devname := disk.DevName.Snapshot().Value()
	dirname := disk.DirName.Snapshot().Value()
	return operating.DiskMeta{
		DevName: devname,
		DirName: dirname,
	}
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

func (procs ProcSlice) Ordered(para *params.Params) *PStable {
	para.Psn.Limit = len(procs)
	limitPS := para.Psn.Absolute
	if limitPS > para.Psn.Limit {
		limitPS = para.Psn.Limit
	}

	pst := &PStable{}
	pst.N = new(int)
	*pst.N = limitPS

	if para.Psn.Absolute == 0 {
		return pst
	}

	uids := map[uint]string{}
	sort.Sort(ProcSort{ // not .Stable
		Psk:       &para.Psk,
		ProcSlice: procs,
		UIDs:      uids,
	})

	for _, proc := range procs[:limitPS] {
		pst.List = append(pst.List, operating.ProcData{
			PID:      proc.PID,
			UID:      proc.UID,
			Priority: proc.Priority,
			Nice:     proc.Nice,
			Time:     format.FormatTime(proc.Time),
			Name:     proc.Name,
			User:     username(uids, proc.UID),
			Size:     format.HumanB(proc.Size),
			Resident: format.HumanB(proc.Resident),
		})
	}
	return pst
}

type IndexData struct {
	Generic // inline non-pointer

	CPU     operating.CPUInfo
	MEM     operating.MEM
	Params  *params.Params `json:",omitempty"`
	PStable PStable

	DFbytes  operating.DFbytes  `json:",omitempty"`
	DFinodes operating.DFinodes `json:",omitempty"`

	IFbytes   operating.Interfaces
	IFerrors  operating.Interfaces
	IFpackets operating.Interfaces

	VagrantMachines VagrantMachines
	VagrantError    string
	VagrantErrord   bool

	DISTRIB string
	VERSION string
}

type PStable struct {
	List []operating.ProcData `json:",omitempty"`
	N    *int                 `json:",omitempty"`
}

type IndexUpdate struct {
	Generic // inline non-pointer

	CPU     *operating.CPUInfo `json:",omitempty"`
	MEM     *operating.MEM     `json:",omitempty"`
	Params  *params.Params     `json:",omitempty"`
	PStable *PStable           `json:",omitempty"`

	DFbytes  *operating.DFbytes  `json:",omitempty"`
	DFinodes *operating.DFinodes `json:",omitempty"`

	IFbytes   *operating.Interfaces `json:",omitempty"`
	IFerrors  *operating.Interfaces `json:",omitempty"`
	IFpackets *operating.Interfaces `json:",omitempty"`

	VagrantMachines *VagrantMachines `json:",omitempty"`
	VagrantError    string
	VagrantErrord   bool

	Location *string `json:",omitempty"`
}

type Generic struct {
	Hostname string `json:",omitempty"`
	Uptime   string `json:",omitempty"`
	LA       string `json:",omitempty"`
	IP       string `json:",omitempty"`
}

type last struct {
	MU       sync.Mutex
	ProcList ProcSlice
}

var lastInfo last

func (la *last) collect(c Collector) {
	var wg sync.WaitGroup
	wg.Add(8)                            // EIGHT:
	go c.CPU(&Reg1s, &wg)                // one
	go c.RAM(&Reg1s, &wg)                // two
	go c.Swap(&Reg1s, &wg)               // three
	go c.Disks(&Reg1s, &wg)              // four
	go c.Hostname(RegMSS, &wg)           // five
	go c.Uptime(RegMSS, &wg)             // six
	go c.LA(&Reg1s, &wg)                 // seven
	go c.Interfaces(&Reg1s, RegMSS, &wg) // eight

	pch := make(chan ProcSlice, 1)
	go c.Procs(pch)

	la.MU.Lock()
	defer la.MU.Unlock()
	la.ProcList = <-pch
	wg.Wait()
}

func (la *last) CopyPS() ProcSlice {
	la.MU.Lock()
	defer la.MU.Unlock()
	psCopy := make(ProcSlice, len(la.ProcList))
	copy(psCopy, la.ProcList)
	return psCopy
}

func (mss *MSS) HN(para *params.Params, iu *IndexUpdate) bool {
	// HN has no delay, always updates iu
	iu.Hostname = mss.GetString("hostname")
	return true
}

func (mss *MSS) IP(para *params.Params, iu *IndexUpdate) bool {
	// IP has no delay, always updates iu
	iu.IP = mss.GetString("ip")
	return true
}

func (mss *MSS) UP(para *params.Params, iu *IndexUpdate) bool {
	// UP has no delay, always updates iu
	iu.Uptime = mss.GetString("uptime")
	return true
}

// InterfaceSlice is a list of MetricInterface.
type InterfaceSlice []operating.MetricInterface

// Len, Swap and Less satisfy sorting interface.
func (is InterfaceSlice) Len() int      { return len(is) }
func (is InterfaceSlice) Swap(i, j int) { is[i], is[j] = is[j], is[i] }
func (is InterfaceSlice) Less(i, j int) bool {
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

func FormatInterface(mi operating.MetricInterface, ip InterfaceParts) operating.InterfaceInfo {
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
	return operating.InterfaceInfo{
		Name:     mi.Name,
		In:       form(uint64(in)),            // format.HumanB(uint64(in)),  // with units
		Out:      form(uint64(out)),           // format.HumanB(uint64(out)), // with units
		DeltaIn:  deltaForm(uint64(deltain)),  // format.Bps64(8, in, 0),     // with units
		DeltaOut: deltaForm(uint64(deltaout)), // format.Bps64(8, out, 0),    // with units
	}
}

type InterfaceParts func(operating.MetricInterface) (*operating.GaugeDiff, *operating.GaugeDiff, bool)

func (_ *IndexRegistry) InterfaceBytes(mi operating.MetricInterface) (*operating.GaugeDiff, *operating.GaugeDiff, bool) {
	return mi.BytesIn, mi.BytesOut, true
}
func (_ *IndexRegistry) InterfaceErrors(mi operating.MetricInterface) (*operating.GaugeDiff, *operating.GaugeDiff, bool) {
	return mi.ErrorsIn, mi.ErrorsOut, false
}
func (_ *IndexRegistry) InterfacePackets(mi operating.MetricInterface) (*operating.GaugeDiff, *operating.GaugeDiff, bool) {
	return mi.PacketsIn, mi.PacketsOut, false
}

func (ir *IndexRegistry) Interfaces(para *params.Params, ip InterfaceParts) []operating.InterfaceInfo {
	private := ir.ListPrivateInterface()
	para.Ifn.Limit = private.Len()

	sort.Sort(private) // not .Stable

	var public []operating.InterfaceInfo
	for i, mi := range private {
		if i >= para.Ifn.Absolute {
			break
		}
		public = append(public, FormatInterface(mi, ip))
	}
	return public
}

// ListPrivateInterface returns list of MetricInterface's by traversing the PrivateIFRegistry.
func (ir *IndexRegistry) ListPrivateInterface() (is InterfaceSlice) {
	ir.PrivateIFRegistry.Each(func(name string, i interface{}) {
		is = append(is, i.(operating.MetricInterface))
	})
	return is
}

// GetOrRegisterPrivateInterface produces a registered in PrivateIFRegistry operating.MetricInterface.
func (ir *IndexRegistry) GetOrRegisterPrivateInterface(name string) operating.MetricInterface {
	ir.PrivateMutex.Lock()
	defer ir.PrivateMutex.Unlock()
	if metric := ir.PrivateIFRegistry.Get(name); metric != nil {
		return metric.(operating.MetricInterface)
	}
	i := operating.MetricInterface{
		Interface: &operating.Interface{
			Name:       name,
			BytesIn:    operating.NewGaugeDiff("interface-"+name+".if_octets.rx", ir.Registry),
			BytesOut:   operating.NewGaugeDiff("interface-"+name+".if_octets.tx", ir.Registry),
			ErrorsIn:   operating.NewGaugeDiff("interface-"+name+".if_errors.rx", ir.Registry),
			ErrorsOut:  operating.NewGaugeDiff("interface-"+name+".if_errors.tx", ir.Registry),
			PacketsIn:  operating.NewGaugeDiff("interface-"+name+".if_packets.rx", ir.Registry),
			PacketsOut: operating.NewGaugeDiff("interface-"+name+".if_packets.tx", ir.Registry),
		},
	}
	ir.PrivateIFRegistry.Register(name, i) // error is ignored
	// errs when the type is not derived from (go-)metrics types
	return i
}

func (ir *IndexRegistry) GetOrRegisterPrivateDF(fs sigar.FileSystem) operating.MetricDF {
	ir.PrivateMutex.Lock()
	defer ir.PrivateMutex.Unlock()
	if fs.DirName == "/" {
		fs.DevName = "root"
	} else {
		fs.DevName = strings.Replace(strings.TrimPrefix(fs.DevName, "/dev/"), "/", "-", -1)
	}
	if metric := ir.PrivateDFRegistry.Get(fs.DevName); metric != nil {
		return metric.(operating.MetricDF)
	}
	label := func(tail string) string {
		return fmt.Sprintf("df-%s.df_complex-%s", fs.DevName, tail)
	}
	r, unusedr := ir.Registry, metrics.NewRegistry()
	i := operating.MetricDF{
		DF: &operating.DF{
			DevName:     &operating.StandardMetricString{}, // unregistered
			DirName:     &operating.StandardMetricString{}, // unregistered
			Free:        metrics.NewRegisteredGaugeFloat64(label("free"), r),
			Reserved:    metrics.NewRegisteredGaugeFloat64(label("reserved"), r),
			Total:       metrics.NewRegisteredGauge(label("total"), unusedr),
			Used:        metrics.NewRegisteredGaugeFloat64(label("used"), r),
			Avail:       metrics.NewRegisteredGauge(label("avail"), unusedr),
			UsePercent:  metrics.NewRegisteredGaugeFloat64(label("usepercent"), unusedr),
			Inodes:      metrics.NewRegisteredGauge(label("inodes"), unusedr),
			Iused:       metrics.NewRegisteredGauge(label("iused"), unusedr),
			Ifree:       metrics.NewRegisteredGauge(label("ifree"), unusedr),
			IusePercent: metrics.NewRegisteredGaugeFloat64(label("iusepercent"), unusedr),
		},
	}
	ir.PrivateDFRegistry.Register(fs.DevName, i) // error is ignored
	// errs when the type is not derived from (go-)metrics types
	return i
}

// CPUSlice is a list of MetricCPU.
type CPUSlice []operating.MetricCPU

// Len, Swap and Less satisfy sorting interface.
func (cs CPUSlice) Len() int      { return len(cs) }
func (cs CPUSlice) Swap(i, j int) { cs[i], cs[j] = cs[j], cs[i] }
func (cs CPUSlice) Less(i, j int) bool {
	a, b := cs[i], cs[j]
	var (
		aidle = a.Idle.Percent.Snapshot().Value()
		bidle = b.Idle.Percent.Snapshot().Value()
		/*
			auser = a.User.Percent.Snapshot().Value()
			anice = a.Nice.Percent.Snapshot().Value()
			asys  = a.Sys.Percent.Snapshot().Value()
			buser = b.User.Percent.Snapshot().Value()
			bnice = b.Nice.Percent.Snapshot().Value()
			bsys  = b.Sys.Percent.Snapshot().Value()
		*/
	)
	return aidle < bidle
	// return (auser + anice + asys) > (buser + bnice + bsys)
}

func (ir *IndexRegistry) DF(para *params.Params, iu *IndexUpdate) bool {
	if !para.Dfd.Expired() {
		return false
	}
	switch para.Dft.Absolute {
	case enums.DFBYTES:
		iu.DFbytes = &operating.DFbytes{List: ir.DFbytes(para)}
	case enums.INODES:
		iu.DFinodes = &operating.DFinodes{List: ir.DFinodes(para)}
	default:
		return false
	}
	return true
}

func (ir *IndexRegistry) DFbytes(para *params.Params) []operating.DiskBytes {
	private := ir.ListPrivateDisk()
	para.Dfn.Limit = len(private)

	sort.Stable(DiskSort{
		Dfk:       &para.Dfk,
		DiskSlice: private,
	})

	var public []operating.DiskBytes
	for i, disk := range private {
		if i >= para.Dfn.Absolute {
			break
		}
		public = append(public, FormatDFbytes(disk))
	}
	return public
}

func FormatDFbytes(md operating.MetricDF) operating.DiskBytes {
	var (
		diskTotal = md.Total.Snapshot().Value()
		diskUsed  = md.Used.Snapshot().Value()
		diskAvail = md.Avail.Snapshot().Value()
	)
	total, approxtotal, _ := format.HumanBandback(uint64(diskTotal))
	used, approxused, _ := format.HumanBandback(uint64(diskUsed))
	return operating.DiskBytes{
		DiskMeta:   diskMeta(md),
		Total:      total,
		Used:       used,
		Avail:      format.HumanB(uint64(diskAvail)),
		UsePercent: format.FormatPercent(approxused, approxtotal),
	}
}

func (ir *IndexRegistry) DFinodes(para *params.Params) []operating.DiskInodes {
	private := ir.ListPrivateDisk()
	para.Dfn.Limit = len(private)

	sort.Stable(DiskSort{
		Dfk:       &para.Dfk,
		DiskSlice: private,
	})

	var public []operating.DiskInodes
	for i, disk := range private {
		if i >= para.Dfn.Absolute {
			break
		}
		public = append(public, FormatDFinodes(disk))
	}
	return public
}

func FormatDFinodes(md operating.MetricDF) operating.DiskInodes {
	var (
		diskInodes = md.Inodes.Snapshot().Value()
		diskIused  = md.Iused.Snapshot().Value()
		diskIfree  = md.Ifree.Snapshot().Value()
	)
	itotal, approxitotal, _ := format.HumanBandback(uint64(diskInodes))
	iused, approxiused, _ := format.HumanBandback(uint64(diskIused))
	return operating.DiskInodes{
		DiskMeta:    diskMeta(md),
		Inodes:      itotal,
		Iused:       iused,
		Ifree:       format.HumanB(uint64(diskIfree)),
		IusePercent: format.FormatPercent(approxiused, approxitotal),
	}
}

func (ir *IndexRegistry) VG(para *params.Params, iu *IndexUpdate) bool {
	if !para.Vgd.Expired() {
		return false
	}
	if para.Vgn.Absolute == 0 {
		para.Vgn.Limit = 1
		return false
	}
	machines, err := vagrantmachines(para.Dfn.Absolute)
	if err != nil {
		iu.VagrantErrord, iu.VagrantError = true, err.Error()
		return true
	}
	para.Vgn.Limit = len(machines.List)
	iu.VagrantErrord, iu.VagrantMachines = false, machines
	return true
}

// ProcSlice is a list of ProcInfo.
type ProcSlice []operating.ProcInfo

func (procs ProcSlice) IU(para *params.Params, iu *IndexUpdate) bool {
	if !para.Psd.Expired() {
		return false
	}
	iu.PStable = procs.Ordered(para)
	return true
}

func (ir *IndexRegistry) IF(para *params.Params, iu *IndexUpdate) bool {
	if !para.Ifd.Expired() {
		return false
	}
	switch para.Ift.Absolute {
	case enums.IFBYTES:
		iu.IFbytes = &operating.Interfaces{List: ir.Interfaces(para, ir.InterfaceBytes)}
	case enums.ERRORS:
		iu.IFerrors = &operating.Interfaces{List: Reg1s.Interfaces(para, ir.InterfaceErrors)}
	case enums.PACKETS:
		iu.IFpackets = &operating.Interfaces{List: Reg1s.Interfaces(para, ir.InterfacePackets)}
	default:
		return false
	}
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
	iu.CPU = ir.CPUInternal(para)
	return true
}

func (ir *IndexRegistry) CPUInternal(para *params.Params) *operating.CPUInfo {
	cpu := &operating.CPUInfo{}
	private := ir.ListPrivateCPU()

	if private.Len() == 1 {
		cpu.List = []operating.CoreInfo{FormatCPU("", private[0])}
		para.CPUn.Limit = 1
		return cpu
	}
	para.CPUn.Limit = private.Len() + 1
	sort.Sort(private)

	allabel := fmt.Sprintf("all %d", private.Len())
	public := []operating.CoreInfo{FormatCPU(allabel, ir.PrivateCPUAll)} // first: "all N"

	for i, mc := range private {
		if i >= para.CPUn.Absolute-1 {
			break
		}
		public = append(public, FormatCPU("", mc))
	}
	cpu.List = public
	return cpu
}

func FormatCPU(label string, mc operating.MetricCPU) operating.CoreInfo {
	if label == "" {
		label = "#" + strings.TrimPrefix(mc.N, "cpu-") // A non-"all" mc.
	}
	return operating.CoreInfo{
		N:    label,
		User: mc.User.SnapshotValueUint(),
		Sys:  mc.Sys.SnapshotValueUint(),
		Wait: mc.Wait.SnapshotValueUint(),
		Idle: mc.Idle.SnapshotValueUint(),
	}
}

// ListPrivateCPU returns list of operating.MetricCPU's by traversing the PrivateCPURegistry.
func (ir *IndexRegistry) ListPrivateCPU() (cs CPUSlice) {
	ir.PrivateCPURegistry.Each(func(name string, i interface{}) {
		cs = append(cs, i.(operating.MetricCPU))
	})
	return cs
}

// DiskSlice is a list of MetricDF.
type DiskSlice []operating.MetricDF

// ListPrivateDisk returns list of operating.MetricDF's by traversing the PrivateDFRegistry.
func (ir *IndexRegistry) ListPrivateDisk() (ds DiskSlice) {
	ir.PrivateDFRegistry.Each(func(name string, i interface{}) {
		ds = append(ds, i.(operating.MetricDF))
	})
	return ds
}

// GetOrRegisterPrivateCPU produces a registered in PrivateCPURegistry MetricCPU.
func (ir *IndexRegistry) GetOrRegisterPrivateCPU(coreno int) operating.MetricCPU {
	ir.PrivateMutex.Lock()
	defer ir.PrivateMutex.Unlock()
	name := fmt.Sprintf("cpu-%d", coreno)
	if metric := ir.PrivateCPURegistry.Get(name); metric != nil {
		return metric.(operating.MetricCPU)
	}
	i := *system.NewMetricCPU(ir.Registry, name)
	ir.PrivateCPURegistry.Register(name, i) // error is ignored
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

func (ir *IndexRegistry) UpdateLoadAverage(la sigar.LoadAverage) {
	ir.Mutex.Lock()
	defer ir.Mutex.Unlock()
	ir.Load.Short.Update(la.One)
	ir.Load.Mid.Update(la.Five)
	ir.Load.Long.Update(la.Fifteen)
}

func (ir *IndexRegistry) UpdateCPU(cpus []sigar.Cpu) {
	ir.Mutex.Lock()
	defer ir.Mutex.Unlock()
	all := sigar.Cpu{}
	for coreno, core := range cpus {
		ir.GetOrRegisterPrivateCPU(coreno).Update(core)
		operating.AddSCPU(&all, core)
	}
	ir.PrivateCPUAll.Update(all)
}

func (ir *IndexRegistry) UpdateIFdata(ifdata IfData) {
	ir.Mutex.Lock()
	defer ir.Mutex.Unlock()
	ir.GetOrRegisterPrivateInterface(ifdata.Name).Update(ifdata)
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
	PrivateCPUAll      operating.MetricCPU
	PrivateCPURegistry metrics.Registry // set of MetricCPUs is handled as a metric in this registry
	PrivateIFRegistry  metrics.Registry // set of operating.MetricInterfaces is handled as a metric in this registry
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
		PrivateCPUAll: *system.NewMetricCPU(metrics.NewRegistry(),
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
	psCopy := lastInfo.CopyPS()

	var updated bool
	for _, update := range []func(*params.Params, *IndexUpdate) bool{
		psCopy.IU,
		RegMSS.HN,
		RegMSS.IP,
		RegMSS.UP,
		Reg1s.MEM,
		Reg1s.CPU,
		Reg1s.DF,
		Reg1s.IF,
		Reg1s.LA,
		Reg1s.VG,
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
		DISTRIB: DISTRIB, // value set in init()
		VERSION: VERSION, // value from server.go

		Params:  updates.Params,
		Generic: updates.Generic,
	}

	if updates.CPU != nil {
		data.CPU = *updates.CPU
	}
	if updates.MEM != nil {
		data.MEM = *updates.MEM
	}
	if updates.PStable != nil {
		data.PStable = *updates.PStable
	}
	if updates.DFbytes != nil {
		data.DFbytes = *updates.DFbytes
	}
	if updates.DFinodes != nil {
		data.DFinodes = *updates.DFinodes
	}
	if updates.IFbytes != nil {
		data.IFbytes = *updates.IFbytes
	}
	if updates.IFerrors != nil {
		data.IFerrors = *updates.IFerrors
	}
	if updates.IFpackets != nil {
		data.IFpackets = *updates.IFpackets
	}
	if updates.VagrantMachines != nil {
		data.VagrantMachines = *updates.VagrantMachines
	}
	data.VagrantError = updates.VagrantError
	data.VagrantErrord = updates.VagrantErrord

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
	IndexTemplate *templateutil.LazyTemplate
}

func NewServeSSE(access *Access, mindelay flags.Delay) *ServeSSE {
	return &ServeSSE{Access: access, MinDelay: mindelay}
}

func NewServeWS(ss ServeSSE, errlog *log.Logger, maxdelay flags.Delay) *ServeWS {
	return &ServeWS{ServeSSE: ss, ErrLog: errlog, MaxDelay: maxdelay}
}

func NewServeIndex(sw ServeWS, taggedbin bool, template *templateutil.LazyTemplate) *ServeIndex {
	return &ServeIndex{ServeWS: sw, TaggedBin: taggedbin, IndexTemplate: template}
}

// Index renders index page.
// si is read-only, pointer is for not copying.
func (si *ServeIndex) Index(w http.ResponseWriter, r *http.Request) {
	id, err := indexData(si.MinDelay, si.MaxDelay, r)
	if err != nil {
		if _, ok := err.(params.RenamedConstError); ok {
			http.Redirect(w, r, err.Error(), http.StatusFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := si.IndexTemplate.Response(w, struct {
		TAGGEDbin bool
		Data      IndexData
	}{
		TAGGEDbin: si.TaggedBin,
		Data:      id,
	})
	response.Header().Set("Content-Type", "text/html")
	response.SetContentLength()
	response.Send()
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
// ss is read-only, pointer is for not copying.
func (ss *ServeSSE) IndexSSE(w http.ResponseWriter, r *http.Request) {
	sse := &SSE{Writer: w, MinDelay: ss.MinDelay}
	if ss.Access.Constructor(sse).ServeHTTP(nil, r); sse.Errord { // the request logging
		return
	}
	for { // loop is access-log-free
		SleepTilNextSecond() // TODO is it second?
		if sse.ServeHTTP(nil, r); sse.Errord {
			break
		}
	}
}
