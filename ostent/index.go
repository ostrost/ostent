// Package ostent is the library part of ostent cmd.
package ostent

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/user"
	"strings"
	"sync"

	"github.com/ostrost/ostent/client"
	"github.com/ostrost/ostent/flags"
	"github.com/ostrost/ostent/format"
	"github.com/ostrost/ostent/getifaddrs"
	"github.com/ostrost/ostent/system"
	"github.com/ostrost/ostent/system/operating"
	"github.com/ostrost/ostent/templateutil"
	metrics "github.com/rcrowley/go-metrics"
	sigar "github.com/rzab/gosigar"
)

func interfaceMeta(ifdata getifaddrs.IfData) operating.InterfaceMeta {
	return interfaceMetaFromString(ifdata.Name)
}

func interfaceMetaFromString(name string) operating.InterfaceMeta {
	return operating.InterfaceMeta{
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

var DefinesTemplate *templateutil.LazyTemplate

func tooltipable(limit int, full string) template.HTML {
	var html string
	if len(full) > limit {
		var err error
		html, err = DefinesTemplate.LookupApply("define_tooltipable", struct {
			Full, Short string
		}{
			Full:  full,
			Short: full[:limit],
		})
		if err != nil {
			html = template.HTMLEscapeString(err.Error())
		}
	} else {
		html = template.HTMLEscapeString(full)
	}
	return template.HTML(html)
}

func diskMeta(disk operating.MetricDF) operating.DiskMeta {
	devname := disk.DevName.Snapshot().Value()
	dirname := disk.DirName.Snapshot().Value()
	return operating.DiskMeta{
		DiskNameHTML: tooltipable(12, devname),
		DirNameHTML:  tooltipable(6, dirname),
		DirNameKey:   dirname,
		DevName:      devname,
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

func (procs MPSlice) Ordered(cl *client.Client, send *client.SendClient) []operating.ProcData {
	uids := map[uint]string{}
	operating.MetricProcSlice(procs).SortSortBy(LessProcFunc(uids, cl.PSSEQ)) // not .StableSortBy

	pslen := len(procs)
	limitPS := cl.PSlimit
	notdec := limitPS <= 1
	notexp := limitPS >= pslen

	if limitPS >= pslen { // notexp
		limitPS = pslen // NB modified limitPS
	} else {
		procs = procs[:limitPS]
	}

	send.SetBool(&send.PSnotDecreasable, &cl.PSnotDecreasable, notdec)
	send.SetBool(&send.PSnotExpandable, &cl.PSnotExpandable, notexp)
	send.SetString(&send.PSplusText, &cl.PSplusText, fmt.Sprintf("%d+", limitPS))

	var list []operating.ProcData
	for _, proc := range procs {
		list = append(list, operating.ProcData{
			PID:      proc.PID,
			UID:      proc.UID,
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

type clientData struct {
	client.Client
	HideMEM    *bool           `json:",omitempty"` // for template only
	RefreshMEM *client.Refresh `json:",omitempty"` // for template only
}

type IndexData struct {
	Generic // inline non-pointer

	CPU     operating.CPUInfo
	MEM     operating.MEM
	Links   *Links `json:",omitempty"`
	PStable PStable

	DFbytes  operating.DFbytes  `json:",omitempty"`
	DFinodes operating.DFinodes `json:",omitempty"`

	IFbytes   operating.Interfaces
	IFerrors  operating.Interfaces
	IFpackets operating.Interfaces

	VagrantMachines *VagrantMachines
	VagrantError    string
	VagrantErrord   bool

	DISTRIB string
	VERSION string
	IFTABS  client.Tabs
	DFTABS  client.Tabs

	Client clientData
}

type IndexUpdate struct {
	Generic // inline non-pointer

	CPU     *operating.CPUInfo `json:",omitempty"`
	MEM     *operating.MEM     `json:",omitempty"`
	Links   *Links             `json:",omitempty"`
	PStable *PStable           `json:",omitempty"`

	DFbytes  *operating.DFbytes  `json:",omitempty"`
	DFinodes *operating.DFinodes `json:",omitempty"`

	IFbytes   *operating.Interfaces `json:",omitempty"`
	IFerrors  *operating.Interfaces `json:",omitempty"`
	IFpackets *operating.Interfaces `json:",omitempty"`

	VagrantMachines *VagrantMachines `json:",omitempty"`
	VagrantError    string
	VagrantErrord   bool

	Client *client.SendClient `json:",omitempty"`
}

type Generic struct {
	Hostname string `json:",omitempty"`
	Uptime   string `json:",omitempty"`
	LA       string `json:",omitempty"`
	IP       string `json:",omitempty"`
}

type last struct {
	MU       sync.Mutex
	ProcList operating.MetricProcSlice
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

	pch := make(chan operating.MetricProcSlice, 1)
	go c.Procs(pch)

	la.MU.Lock()
	defer la.MU.Unlock()
	la.ProcList = <-pch
	wg.Wait()
}

func (la *last) CopyPS() MPSlice {
	la.MU.Lock()
	defer la.MU.Unlock()
	psCopy := make(MPSlice, len(la.ProcList))
	copy(psCopy, la.ProcList)
	return psCopy
}

func (mss *MSS) HN(cli *client.Client, send *client.SendClient, iu *IndexUpdate) interface{} {
	iu.Hostname = mss.GetString("hostname")
	generic := iu.Generic
	generic.Hostname = iu.Hostname
	return IndexUpdate{Generic: generic}
}

func (mss *MSS) UP(cli *client.Client, send *client.SendClient, iu *IndexUpdate) interface{} {
	iu.Uptime = mss.GetString("uptime")
	generic := iu.Generic
	generic.Uptime = iu.Uptime
	return IndexUpdate{Generic: generic}
}

func (mss *MSS) IP(cli *client.Client, send *client.SendClient, iu *IndexUpdate) interface{} {
	iu.IP = mss.GetString("ip")
	generic := iu.Generic
	generic.IP = iu.IP
	return IndexUpdate{Generic: generic}
}

func LessInterface(a, b operating.MetricInterface) bool {
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
		InterfaceMeta: interfaceMetaFromString(mi.Name),
		In:            form(uint64(in)),            // format.HumanB(uint64(in)),  // with units
		Out:           form(uint64(out)),           // format.HumanB(uint64(out)), // with units
		DeltaIn:       deltaForm(uint64(deltain)),  // format.Bps64(8, in, 0),     // with units
		DeltaOut:      deltaForm(uint64(deltaout)), // format.Bps64(8, out, 0),    // with units
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

func (ir *IndexRegistry) Interfaces(cli *client.Client, send *client.SendClient, ip InterfaceParts) []operating.InterfaceInfo {
	private := ir.ListPrivateInterface()

	send.SetBool(&send.ExpandableIF, &cli.ExpandableIF, len(private) > cli.Toprows)
	send.SetString(&send.ExpandtextIF, &cli.ExpandtextIF, fmt.Sprintf("Expanded (%d)", len(private)))

	private.SortSortBy(LessInterface)
	var public []operating.InterfaceInfo
	for i, mi := range private {
		if !*cli.ExpandIF && i >= cli.Toprows {
			break
		}
		public = append(public, FormatInterface(mi, ip))
	}
	return public
}

// ListPrivateInterface returns list of MetricInterface's by traversing the PrivateInterfaceRegistry.
func (ir *IndexRegistry) ListPrivateInterface() (lmi operating.MetricInterfaceSlice) {
	ir.PrivateInterfaceRegistry.Each(func(name string, i interface{}) {
		lmi = append(lmi, i.(operating.MetricInterface))
	})
	return lmi
}

// GetOrRegisterPrivateInterface produces a registered in PrivateInterfaceRegistry operating.MetricInterface.
func (ir *IndexRegistry) GetOrRegisterPrivateInterface(name string) operating.MetricInterface {
	ir.PrivateMutex.Lock()
	defer ir.PrivateMutex.Unlock()
	if metric := ir.PrivateInterfaceRegistry.Get(name); metric != nil {
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
	ir.PrivateInterfaceRegistry.Register(name, i) // error is ignored
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

func LessCPU(a, b operating.MetricCPU) bool {
	var (
		auser = a.User.Percent.Snapshot().Value()
		anice = a.Nice.Percent.Snapshot().Value()
		asys  = a.Sys.Percent.Snapshot().Value()
		buser = b.User.Percent.Snapshot().Value()
		bnice = b.Nice.Percent.Snapshot().Value()
		bsys  = b.Sys.Percent.Snapshot().Value()
	)
	return (auser + anice + asys) > (buser + bnice + bsys)
}

func (ir *IndexRegistry) DF(cli *client.Client, send *client.SendClient, iu *IndexUpdate) interface{} {
	switch cli.TabDF.Uint {
	case client.DFBYTES:
		iu.DFbytes = &operating.DFbytes{List: ir.DFbytes(cli, send)}
		return IndexUpdate{DFbytes: iu.DFbytes}
	case client.DFINODES:
		iu.DFinodes = &operating.DFinodes{List: ir.DFinodes(cli, send)}
		return IndexUpdate{DFinodes: iu.DFinodes}
	}
	return nil
}

func (ir *IndexRegistry) DFbytes(cli *client.Client, send *client.SendClient) []operating.DiskBytes {
	private := ir.ListPrivateDisk()

	send.SetBool(&send.ExpandableDF, &cli.ExpandableDF, len(private) > cli.Toprows)
	send.SetString(&send.ExpandtextDF, &cli.ExpandtextDF, fmt.Sprintf("Expanded (%d)", len(private)))

	private.StableSortBy(LessDiskFunc(cli.DFSEQ))

	var public []operating.DiskBytes
	for i, disk := range private {
		if !*cli.ExpandDF && i > cli.Toprows-1 {
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
		DiskMeta:        diskMeta(md),
		Total:           total,
		Used:            used,
		Avail:           format.HumanB(uint64(diskAvail)),
		UsePercent:      format.FormatPercent(approxused, approxtotal),
		UsePercentClass: format.LabelClassColorPercent(format.Percent(approxused, approxtotal)),
	}
}

func (ir *IndexRegistry) DFinodes(cli *client.Client, send *client.SendClient) []operating.DiskInodes {
	private := ir.ListPrivateDisk()

	send.SetBool(&send.ExpandableDF, &cli.ExpandableDF, len(private) > cli.Toprows)
	send.SetString(&send.ExpandtextDF, &cli.ExpandtextDF, fmt.Sprintf("Expanded (%d)", len(private)))

	private.StableSortBy(LessDiskFunc(cli.DFSEQ))

	var public []operating.DiskInodes
	for i, disk := range private {
		if !*cli.ExpandDF && i > cli.Toprows-1 {
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
		DiskMeta:         diskMeta(md),
		Inodes:           itotal,
		Iused:            iused,
		Ifree:            format.HumanB(uint64(diskIfree)),
		IusePercent:      format.FormatPercent(approxiused, approxitotal),
		IusePercentClass: format.LabelClassColorPercent(format.Percent(approxiused, approxitotal)),
	}
}

func (ir *IndexRegistry) VG(cli *client.Client, send *client.SendClient, iu *IndexUpdate) interface{} {
	machines, err := vagrantmachines()
	if err != nil {
		iu.VagrantErrord = true
		iu.VagrantError = err.Error()
		return IndexUpdate{
			VagrantErrord: iu.VagrantErrord,
			VagrantError:  iu.VagrantError,
		}
	}
	iu.VagrantErrord = false
	iu.VagrantMachines = machines
	return IndexUpdate{
		VagrantErrord:   iu.VagrantErrord,
		VagrantMachines: iu.VagrantMachines,
	}
}

// MPSlice is a operating.MetricProcSlice with some methods.
type MPSlice operating.MetricProcSlice

func (procs MPSlice) IU(cli *client.Client, send *client.SendClient, iu *IndexUpdate) interface{} {
	iu.PStable = &PStable{List: procs.Ordered(cli, send)}
	return IndexUpdate{PStable: iu.PStable}
}

func (ir *IndexRegistry) IF(cli *client.Client, send *client.SendClient, iu *IndexUpdate) interface{} {
	switch cli.TabIF.Uint {
	case client.IFBYTES:
		iu.IFbytes = &operating.Interfaces{List: ir.Interfaces(cli, send, ir.InterfaceBytes)}
		return IndexUpdate{IFbytes: iu.IFbytes}
	case client.IFERRORS:
		iu.IFerrors = &operating.Interfaces{List: Reg1s.Interfaces(cli, send, ir.InterfaceErrors)}
		return IndexUpdate{IFerrors: iu.IFerrors}
	case client.IFPACKETS:
		iu.IFpackets = &operating.Interfaces{List: Reg1s.Interfaces(cli, send, ir.InterfacePackets)}
		return IndexUpdate{IFpackets: iu.IFpackets}
	}
	return nil
}

func (ir *IndexRegistry) CPU(cli *client.Client, send *client.SendClient, iu *IndexUpdate) interface{} {
	list := ir.CPUInternal(cli, send)
	iu.CPU = &operating.CPUInfo{List: list}
	return IndexUpdate{CPU: iu.CPU}
}

func (ir *IndexRegistry) CPUInternal(cli *client.Client, send *client.SendClient) []operating.CoreInfo {
	private := ir.ListPrivateCPU()

	send.SetBool(&send.ExpandableCPU, &cli.ExpandableCPU, len(private) > cli.Toprows) // one row reserved for "all N"
	send.SetString(&send.ExpandtextCPU, &cli.ExpandtextCPU, fmt.Sprintf("Expanded (%d)", len(private)))

	if len(private) == 1 {
		return []operating.CoreInfo{FormatCPU(private[0])}
	}
	private.SortSortBy(LessCPU)
	var public []operating.CoreInfo
	if !*cli.ExpandCPU {
		public = []operating.CoreInfo{FormatCPU(ir.PrivateCPUAll)}
	}
	for i, mc := range private {
		if !*cli.ExpandCPU && i > cli.Toprows-2 {
			// "collapsed" view, head of the list
			break
		}
		public = append(public, FormatCPU(mc))
	}
	return public
}

func FormatCPU(mc operating.MetricCPU) operating.CoreInfo {
	user := uint(mc.User.Percent.Snapshot().Value()) // rounding
	// .Nice is unused
	sys := uint(mc.Sys.Percent.Snapshot().Value())   // rounding
	idle := uint(mc.Idle.Percent.Snapshot().Value()) // rounding
	N := mc.N
	if prefix := "cpu-"; strings.HasPrefix(N, prefix) { // true for all but "all"
		N = "#" + N[len(prefix):] // fmt.Sprintf("#%d", n)
	}
	return operating.CoreInfo{
		N:         N,
		User:      user,
		Sys:       sys,
		Idle:      idle,
		UserClass: format.TextClassColorPercent(user),
		SysClass:  format.TextClassColorPercent(sys),
		IdleClass: format.TextClassColorPercent(100 - idle),
	}
}

// ListPrivateCPU returns list of operating.MetricCPU's by traversing the PrivateCPURegistry.
func (ir *IndexRegistry) ListPrivateCPU() (lmc operating.MetricCPUSlice) {
	ir.PrivateCPURegistry.Each(func(name string, i interface{}) {
		lmc = append(lmc, i.(operating.MetricCPU))
	})
	return lmc
}

// ListPrivateDisk returns list of operating.MetricDF's by traversing the PrivateDFRegistry.
func (ir *IndexRegistry) ListPrivateDisk() (lmd operating.MetricDFSlice) {
	ir.PrivateDFRegistry.Each(func(name string, i interface{}) {
		lmd = append(lmd, i.(operating.MetricDF))
	})
	return lmd
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

func (ir *IndexRegistry) SWAP(client *client.Client, send *client.SendClient, iu *IndexUpdate) interface{} {
	// client is unused
	// send is unused
	if iu.MEM == nil {
		iu.MEM = new(operating.MEM)
	}
	if iu.MEM.List == nil {
		iu.MEM.List = []operating.Memory{}
	}
	gs := ir.Swap
	iu.MEM.List = append(iu.MEM.List,
		_getmem("swap", sigar.Swap{
			Total: gs.TotalValue(),
			Free:  uint64(gs.Free.Snapshot().Value()),
			Used:  uint64(gs.Used.Snapshot().Value()),
		}))
	// did modify iu.MEM
	return IndexUpdate{MEM: iu.MEM}
}

func (ir *IndexRegistry) MEM(client *client.Client, send *client.SendClient, iu *IndexUpdate) interface{} {
	// client is unused
	// send is unused
	gr := ir.RAM
	mem := new(operating.MEM)
	mem.List = []operating.Memory{
		_getmem("RAM", sigar.Swap{
			Total: uint64(gr.Total.Snapshot().Value()),
			Free:  uint64(gr.Free.Snapshot().Value()),
			Used:  gr.UsedValue(), // == .Total - .Free
		}),
	}
	iu.MEM = mem
	return IndexUpdate{MEM: iu.MEM}
}

func (ir *IndexRegistry) LA(cli *client.Client, send *client.SendClient, iu *IndexUpdate) interface{} {
	gl := ir.Load
	iu.LA = gl.Short.Sparkline() + " " + fmt.Sprintf("%.2f %.2f %.2f",
		gl.Short.Snapshot().Value(),
		gl.Mid.Snapshot().Value(),
		gl.Long.Snapshot().Value())
	generic := iu.Generic
	generic.LA = iu.LA
	return IndexUpdate{Generic: generic}
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
	if ir.PrivateCPUAll.N == "all" {
		ir.PrivateCPUAll.N = fmt.Sprintf("all %d", len(cpus))
	}
	ir.PrivateCPUAll.Update(all)
}

func (ir *IndexRegistry) UpdateIFdata(ifdata getifaddrs.IfData) {
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
	Registry                 metrics.Registry
	PrivateCPUAll            operating.MetricCPU
	PrivateCPURegistry       metrics.Registry // set of MetricCPUs is handled as a metric in this registry
	PrivateInterfaceRegistry metrics.Registry // set of operating.MetricInterfaces is handled as a metric in this registry
	PrivateDFRegistry        metrics.Registry // set of operating.MetricDFs is handled as a metric in this registry
	PrivateMutex             sync.Mutex

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
	Reg1s = IndexRegistry{
		Registry:                 metrics.NewRegistry(),
		PrivateCPURegistry:       metrics.NewRegistry(),
		PrivateInterfaceRegistry: metrics.NewRegistry(),
		PrivateDFRegistry:        metrics.NewRegistry(),
	}
	Reg1s.PrivateCPUAll = /* *Reg1s.RegisterCPU */ *system.NewMetricCPU(
		/* pcreg := */ metrics.NewRegistry(), "all")
	// pcreg.Register("all", Reg1s.PrivateCPUAll)

	Reg1s.RAM = system.NewMetricRAM(Reg1s.Registry)
	Reg1s.Swap = operating.NewMetricSwap(Reg1s.Registry)
	Reg1s.Load = operating.NewMetricLoad(Reg1s.Registry)
}

type Set struct {
	Hide    bool
	Refresh *client.Refresh `json:",omitempty"`
	Update  func(*client.Client, *client.SendClient, *IndexUpdate) interface{}
}

func (s Set) Hidden() bool { return s.Hide }
func (s *Set) Expired(forcerefresh bool) bool {
	return s.Refresh.Refresh(forcerefresh)
}

/*
type SetInterface interface {
	Hidden() bool
	Expired(bool) bool
	Update(*client.Client, *client.SendClient, *IndexUpdate) interface{}
}
// */

func getUpdates(req *http.Request, cl *client.Client, send client.SendClient, forcerefresh bool) (iu IndexUpdate, err error) {
	cl.RecalcRows() // before anything

	psCopy := lastInfo.CopyPS()

	if req != nil {
		req.ParseForm() // do ParseForm even if req.Form == nil
		iu.Links = &Links{*client.NewLinks()}
		client.DF.Decode(req.Form, "df", iu.Links, &cl.DFSEQ, new(client.UintDF))
		client.DF.Decode(req.Form, "ps", iu.Links, &cl.PSSEQ, new(client.UintPS))

		if iu.Links.Decodes.RCError != nil {
			return iu, client.RenamedConstError("?" + iu.Links.Values.Encode())
		}
	}

	set := []Set{
		{*cl.HideRAM, cl.RefreshRAM, Reg1s.MEM},
		// if RAM is hidden, so is SWAP:
		{*cl.HideRAM || *cl.HideSWAP, cl.RefreshSWAP, Reg1s.SWAP},
		{*cl.HideCPU, cl.RefreshCPU, Reg1s.CPU},
		{*cl.HideDF, cl.RefreshDF, Reg1s.DF},
		{*cl.HideIF, cl.RefreshIF, Reg1s.IF},
		{*cl.HidePS, cl.RefreshPS, psCopy.IU},
		{*cl.HideVG, cl.RefreshVG, Reg1s.VG},

		// always-shown bits:
		{false, cl.RefreshHN, RegMSS.HN},
		{false, cl.RefreshHN, RegMSS.UP},
		{false, cl.RefreshHN, RegMSS.IP},
		{false, cl.RefreshHN, Reg1s.LA},
	}

	// var additions []interface{}
	for _, x := range set {
		if !x.Expired(forcerefresh) { // this has side effect
			continue
		}
		if x.Hidden() {
			continue
		}
		if add := x.Update(cl, &send, &iu); add != nil {
			// additions = append(additions, add)
		}
	}

	if send.Modified {
		iu.Client = &send
	}
	return iu, nil
}

func indexData(minperiod flags.Period, req *http.Request) (IndexData, error) {
	if Connections.Len() == 0 {
		// collect when there're no active connections, so Loop does not collect
		lastInfo.collect(&Machine{})
	}

	cl := client.DefaultClient(minperiod)
	updates, err := getUpdates(req, &cl, client.SendClient{}, true)
	if err != nil {
		return IndexData{}, err
	}

	data := IndexData{
		Generic: updates.Generic,

		CPU:     *updates.CPU,
		MEM:     *updates.MEM,
		Links:   updates.Links,
		PStable: *updates.PStable,

		DISTRIB: DISTRIB,       // value set in init()
		VERSION: VERSION,       // value from server.go
		DFTABS:  client.DFTABS, // "const"
		IFTABS:  client.IFTABS, // "const"

		Client: clientData{Client: cl, HideMEM: cl.HideRAM, RefreshMEM: cl.RefreshRAM},
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

func IndexFunc(production bool, template *templateutil.LazyTemplate, minperiod flags.Period) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		index(production, template, minperiod, w, r)
	}
}

func index(production bool, template *templateutil.LazyTemplate, minperiod flags.Period, w http.ResponseWriter, r *http.Request) {
	id, err := indexData(minperiod, r)
	if err != nil {
		if _, ok := err.(client.RenamedConstError); ok {
			http.Redirect(w, r, err.Error(), http.StatusMovedPermanently)
			return
		}
		http.Error(w, err.Error(), panicstatuscode)
		return
	}

	response := template.Response(w, struct {
		CLASSNAME  string // MUST HAVE
		PRODUCTION bool
		Data       IndexData
	}{
		PRODUCTION: production,
		Data:       id,
	})
	response.Header().Set("Content-Type", "text/html")
	response.SetContentLength()
	response.Send()
}

type SSE struct {
	Writer      http.ResponseWriter // points to the writer
	MinPeriod   flags.Period
	SentHeaders bool
	Errord      bool
}

// ServeHTTP is a regular serve func except the first argument,
// passed as a copy, is unused. sse.Writer is there for writes.
func (sse *SSE) ServeHTTP(_ http.ResponseWriter, r *http.Request) {
	w := sse.Writer
	id, err := indexData(sse.MinPeriod, r)
	if err != nil {
		http.Error(w, err.Error(), panicstatuscode)
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

func IndexSSEFunc(access *logger, minperiod flags.Period) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		IndexSSE(access, minperiod, w, r)
	}
}

func IndexSSE(access *logger, minperiod flags.Period, w http.ResponseWriter, r *http.Request) {
	sse := &SSE{Writer: w, MinPeriod: minperiod}
	// The request is logged just once.
	if access.Constructor(sse).ServeHTTP(nil, r); sse.Errord {
		return
	}
	// Loop is access-log-free.
	for {
		SleepTilNextSecond() // TODO is it second?
		if sse.ServeHTTP(nil, r); sse.Errord {
			break
		}
	}
}
