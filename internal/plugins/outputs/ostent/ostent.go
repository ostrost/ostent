package ostent

import (
	"fmt"
	"os/user"
	"sort"
	"strings"
	"sync"

	"github.com/shirou/gopsutil/disk"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/outputs"
	// "github.com/influxdata/telegraf/plugins/serializers"

	"github.com/ostrost/ostent/format"
	"github.com/ostrost/ostent/params"
)

type convert struct {
	k string // k for key; required, others optional
	v *int64 // v for value
	// for humanB conversion:
	s *string // s for string
	r *uint64 // r for round
	f func(uint64) string
}

// decode returns true if everything went ok.
func decode(fields map[string]interface{}, converts []convert) bool {
	for _, c := range converts {
		if f, ok := fields[c.k]; ok {
			if v, ok := f.(int64); ok {
				decodeInt(v, &c)
				continue
			} else if v, ok := f.(float64); ok {
				decodeInt(int64(v), &c)
				continue
			}
		}
		return false // either fields lookup or casting failed
	}
	return true
}

func decodeInt(v int64, c *convert) {
	if c.v != nil {
		*c.v = v
	}
	if c.s != nil {
		if c.f != nil {
			*c.s = c.f(uint64(v))
		} else {
			*c.s = humanB(v, c.r)
		}
	}
}

var humanUnitless = format.HumanUnitless

func humanBits(n uint64) string { return format.HumanBits(n * 8) }

func humanB(value int64, round *uint64) string {
	if round == nil {
		return format.HumanB(uint64(value))
	}
	var s string
	s, *round, _ = format.HumanBandback(uint64(value)) // err is ignored
	return s
}

func newFields(m telegraf.Metric) *fields { return &fields{m.Fields()} }

// methods have pointer receiver not to copy every time.
type fields struct{ mfields map[string]interface{} }

func (fs *fields) decodeFloat64(key string, p *float64) bool {
	v, ok := fs.mfields[key]
	if !ok {
		return false
	}
	n, ok := v.(float64)
	if !ok {
		return false
	}
	*p = n
	return true
}

func (fs *fields) decodeInt64(key string, p *int64) bool {
	v, ok := fs.mfields[key]
	if !ok {
		return false
	}
	n, ok := v.(int64)
	if !ok {
		return false
	}
	*p = n
	return true
}

type cpuData struct {
	N string

	// percents without "%"
	UserPct int64
	SysPct  int64
	WaitPct int64
	IdlePct int64
}

type diskData struct {
	DevName string
	DirName string

	// values for End func:

	rounds, roundInodes struct{ used, free uint64 }

	// values for sorting:

	total int64
	used  int64
	avail int64

	// strings with units

	// bytes
	Total  string
	Used   string
	Avail  string
	UsePct uint // percent without "%"

	// inodes
	Inodes  string
	Iused   string
	Ifree   string
	IusePct uint // percent without "%"
}

type laData struct{ Period, Value string }

type memData struct {
	Kind string

	// strings with units
	Total string
	Used  string
	Free  string

	UsePct uint // percent without "%"
}

type netData struct {
	loopback bool

	Name string
	IP   string `json:",omitempty"` // may be empty

	// strings with units

	BytesIn    string
	BytesOut   string
	DropsIn    string
	DropsOut   string `json:",omitempty"`
	ErrorsIn   string
	ErrorsOut  string
	PacketsIn  string
	PacketsOut string

	DeltaBitsIn      string
	DeltaBitsOut     string
	DeltaBytesOutNum int64
	DeltaDropsIn     string
	DeltaDropsOut    string `json:",omitempty"`
	DeltaErrorsIn    string
	DeltaErrorsOut   string
	DeltaPacketsIn   string
	DeltaPacketsOut  string
}

type procData struct {
	PID      int64 // int32 from fields
	UID      int64 // int32 from fields
	Priority int64 // int32 from fields // missing with gopsutil
	Nice     int64 // int32 from fields // gopsutil miscalcs this

	time_user, time_system float64 // float64 from fields

	time     float64 // float64 from fields
	size     int64   // uint64 from fields
	resident int64   // uint64 from fields

	Time     string // formatted
	Name     string
	User     string // username from .UID
	Size     string // with units
	Resident string // with units
}

type cpuList []*cpuData

// Len, Swap and Less satisfy sorting interface.
func (cl cpuList) Len() int           { return len(cl) }
func (cl cpuList) Swap(i, j int)      { cl[i], cl[j] = cl[j], cl[i] }
func (cl cpuList) Less(i, j int) bool { return cl[i].IdlePct < cl[j].IdlePct }

type diskList struct {
	k    *params.Num // a pointer to set .Alpha
	list []*diskData
}

func (dl diskList) Len() int      { return len(dl.list) }
func (dl diskList) Swap(i, j int) { dl.list[i], dl.list[j] = dl.list[j], dl.list[i] }
func (dl diskList) Less(i, j int) (r bool) {
	if match, isa, cmpr := ddCmp(dl.k.Absolute, dl.list[i], dl.list[j]); match {
		dl.k.Alpha, r = isa, cmpr
	}
	if dl.k.Negative {
		return !r
	}
	return r
}

func ddCmp(k int, a, b *diskData) (bool, bool, bool) {
	switch k {
	case params.FS:
		return true, true, a.DevName < b.DevName
	case params.MP:
		return true, true, a.DirName < b.DirName

	case params.TOTAL:
		return true, false, a.total > b.total
	case params.USED:
		return true, false, a.used > b.used
	case params.AVAIL:
		return true, false, a.avail > b.avail
	case params.USEPCT:
		cmp := float64(a.used)/float64(a.total) > float64(b.used)/float64(b.total)
		return true, false, cmp
	}
	return false, false, false
}

type netList []*netData

func (nl netList) Len() int      { return len(nl) }
func (nl netList) Swap(i, j int) { nl[i], nl[j] = nl[j], nl[i] }
func (nl netList) Less(i, j int) bool {
	a, b := nl[i], nl[j]
	if !(a.loopback && b.loopback) {
		if a.loopback {
			return false
		} else if b.loopback {
			return true
		}
	}
	return a.Name < b.Name
}

type procList struct {
	k             *params.Num // a pointer to set .Alpha
	list          []*procData
	usernames     map[int64]string
	copyUsernames map[int64]string
}

func (pl procList) Len() int      { return len(pl.list) }
func (pl procList) Swap(i, j int) { pl.list[i], pl.list[j] = pl.list[j], pl.list[i] }
func (pl procList) Less(i, j int) (r bool) {
	k, a, b := pl.k.Absolute, pl.list[i], pl.list[j]
	if match, isa, cmpr := pdCmp(k, a, b); match {
		pl.k.Alpha, r = isa, cmpr
	} else if k == params.USER {
		if pl.usernames == nil {
			pl.usernames = make(map[int64]string)
			for k, v := range pl.copyUsernames {
				pl.usernames[k] = v
			}
		}
		pl.k.Alpha, r = true, username(pl.usernames, a.UID) < username(pl.usernames, b.UID)
	}
	if pl.k.Negative {
		return !r
	}
	return r
}

func pdCmp(k int, a, b *procData) (bool, bool, bool) {
	switch k {
	case params.PID:
		return true, false, a.PID > b.PID
	case params.PRI:
		return true, false, a.Priority > b.Priority
	case params.NICE:
		return true, false, a.Nice > b.Nice
	case params.VIRT:
		return true, false, a.size > b.size
	case params.RES:
		return true, false, a.resident > b.resident
	case params.TIME:
		return true, false, a.time > b.time
	case params.UID:
		return true, false, a.UID > b.UID

	case params.NAME: // alpha
		return true, true, a.Name < b.Name
	}
	return false, false, false
}

type ostent struct {
	// serializer serializers.Serializer
}

type dparts struct {
	mutex sync.Mutex
	parts map[string]string
}

var diskParts = &dparts{parts: make(map[string]string)}

func positiveLimit(n *params.Num) { n.Limit = 1 }
func whenZero(n *params.Num) bool {
	if n.Absolute == 0 {
		positiveLimit(n)
		return true
	}
	return false
}
func limit(n *params.Num, lim int) int {
	n.Limit = lim
	if n.Absolute > n.Limit {
		n.Absolute = n.Limit
	}
	return n.Absolute
}

func (o *ostent) CopyCPU(up *Update, p *params.Params) interface{} { return o.copyCPU(up, p) }
func (o *ostent) copyCPU(up *Update, para *params.Params) []*cpuData {
	n := &para.CPUn
	if whenZero(n) {
		return nil
	}

	llen := len(up.cpuData)
	if llen == 0 {
		positiveLimit(n)
		return nil
	}

	tshift := 0 // "total shift"
	if up.cpuData[llen-1].N == "cpu-total" {
		tshift = 1
	}

	dup := make([]*cpuData, llen)
	copy(dup[tshift:], up.cpuData[:llen-tshift])
	sort.Sort(cpuList(dup[tshift:]))
	if tshift != 0 {
		dup[0] = up.cpuData[llen-tshift] // last, "cpu-total", becomes first
	}

	return dup[:limit(n, llen)]
}

func (o *ostent) CopyDisk(up *Update, p *params.Params) interface{} { return o.copyDisk(up, p) }
func (o *ostent) copyDisk(up *Update, para *params.Params) []*diskData {
	n := &para.Dfn
	if whenZero(n) {
		return nil
	}

	llen := len(up.diskData)
	if llen == 0 {
		positiveLimit(n)
		return nil
	}

	dup := make([]*diskData, llen)
	copy(dup, up.diskData)
	sort.Stable(diskList{k: &para.Dfk, list: dup})

	return dup[:limit(n, llen)]
}

func (o *ostent) CopyLA(up *Update, p *params.Params) interface{} { return o.copyLA(up, p) }
func (o *ostent) copyLA(up *Update, para *params.Params) []laData {
	n := &para.Lan
	if whenZero(n) {
		return nil
	}

	periods := []string{"1", "5", "15"}[:limit(n, 3)]
	dup := make([]laData, len(periods))
	for i, period := range periods {
		if v, ok := up.kv["load"+period]; ok {
			if f, ok := v.(float64); ok {
				dup[i] = laData{
					Period: period,
					Value:  fmt.Sprintf("%.2f", f),
				}
			}
		}
	}
	return dup
}

func (o *ostent) CopyMem(up *Update, p *params.Params) interface{} { return o.copyMem(up, p) }
func (o *ostent) copyMem(up *Update, para *params.Params) []*memData {
	n := &para.Memn
	if whenZero(n) {
		return nil
	}

	dup := make([]*memData, len(up.memData))
	copy(dup, up.memData[:])

	return dup[:limit(n, 2)]
}

func (o *ostent) CopyNet(up *Update, p *params.Params) interface{} { return o.copyNet(up, p) }
func (o *ostent) copyNet(up *Update, para *params.Params) []*netData {
	n := &para.Ifn
	if whenZero(n) {
		return nil
	}

	llen := len(up.netData)
	if llen == 0 {
		positiveLimit(n)
		return nil
	}

	dup := make([]*netData, llen)
	copy(dup, up.netData)
	sort.Sort(netList(dup)) // not .Stable

	return dup[:limit(n, llen)]
}

func username(usernames map[int64]string, uid int64) string {
	if s, ok := usernames[uid]; ok {
		return s
	}
	s := fmt.Sprintf("%d", uid)
	if usr, err := user.LookupId(s); err == nil {
		s = usr.Username
	}
	usernames[uid] = s
	return s
}

func (o *ostent) CopyProc(up *Update, p *params.Params) interface{} { return o.copyProc(up, p) }
func (o *ostent) copyProc(up *Update, para *params.Params) []*procData {
	n := &para.Psn
	if whenZero(n) {
		return nil
	}

	llen := len(up.procData)
	if llen == 0 {
		positiveLimit(n)
		return nil
	}

	dup := make([]*procData, llen)
	copy(dup, up.procData)

	pl := procList{k: &para.Psk, list: dup, copyUsernames: up.usernames}
	sort.Sort(pl) // not .Stable

	return dup[:limit(n, llen)]
}

func (o *ostent) CopySO(up *Update, _ *params.Params) map[string]string {
	const skipprefix = "load"

	mlen := len(up.kv)
	for k := range up.kv {
		if strings.HasPrefix(k, skipprefix) {
			mlen--
		}
	}

	dup := make(map[string]string, mlen)
	for k, v := range up.kv {
		if strings.HasPrefix(k, skipprefix) {
			continue
		}
		if s, ok := v.(string); ok {
			dup[k] = s
		}
	}
	return dup
}

func writeProcstat(m telegraf.Metric, up *Update) bool {
	fs := newFields(m)

	var pid int64
	if !fs.decodeInt64("pid", &pid) || pid == 0 {
		return false
	}

	push := true
	var pd *procData
	/*
		for _, p := range up.procData {
			if p.PID == pid {
				pd, push = p, false
				break
			}
		}
	*/
	tags := m.Tags()
	if pd == nil {
		pd = &procData{PID: pid, Name: tags["process_name"]}
	}

	if fs.decodeInt64("uid", &pd.UID) {
		pd.User = username(up.usernames, pd.UID)
	}

	fs.decodeInt64("prio", &pd.Priority)
	fs.decodeInt64("nice", &pd.Nice)

	fs.decodeFloat64("cpu_time_user", &pd.time_user)
	fs.decodeFloat64("cpu_time_system", &pd.time_system)

	if fs.decodeInt64("memory_vms", &pd.size) {
		pd.Size = format.HumanB(uint64(pd.size))
	}
	if fs.decodeInt64("memory_rss", &pd.resident) {
		pd.Resident = format.HumanB(uint64(pd.resident))
	}

	if push {
		up.procData = append(up.procData, pd)
	}
	_, ok := tags["elapsed"]
	return ok
}

func writeProcstatEnd(up *Update) {
	for _, pd := range up.procData {
		if pd == nil {
			continue
		}
		pd.time = 1000.0 * (pd.time_user + pd.time_system)
		pd.Time = format.Time(uint64(pd.time))
	}
}

func writeSystemCPU(m telegraf.Metric, up *Update, cpui int) bool {
	n := m.Tags()["cpu"]

	push := true
	var cd *cpuData
	/*
		for _, c := range up.cpuData {
			if c != nil && c.N == n {
				cd, push = c, false
				break
			}
		}
	*/
	if cd == nil {
		cd = &cpuData{N: n}
	}

	if !decode(m.Fields(), []convert{
		{k: "usage_user", v: &cd.UserPct},
		{k: "usage_system", v: &cd.SysPct},
		{k: "usage_iowait", v: &cd.WaitPct},
		{k: "usage_idle", v: &cd.IdlePct},
	}) {
		return false // either fields lookup or casting failed
	}

	if push {
		up.cpuData[cpui] = cd
		return true
	}
	return false
}

func (dp *dparts) mpDevice(mountpoint string) (string, error) {
	dp.mutex.Lock()
	defer dp.mutex.Unlock()

	if device, ok := dp.parts[mountpoint]; ok {
		return device, nil
	}
	parts, err := disk.Partitions(true)
	if err != nil {
		return "", err
	}
	for _, p := range parts {
		dp.parts[p.Mountpoint] = p.Device
	}
	return dp.parts[mountpoint], nil
}

func writeSystemDisk(m telegraf.Metric, up *Update, diski int) bool {
	path := m.Tags()["path"]

	push := true
	var dd *diskData
	/*
		for _, d := range up.diskData {
			if d != nil && d.DirName == path {
				dd, push = d, false
				break
			}
		}
	*/
	if dd == nil {
		device, _ := diskParts.mpDevice(path) // err is ignored
		dd = &diskData{DirName: path, DevName: device}
	}

	if !decode(m.Fields(), []convert{
		{k: "total", s: &dd.Total, v: &dd.total},
		{k: "used", s: &dd.Used, v: &dd.used, r: &dd.rounds.used},
		{k: "free", s: &dd.Avail, v: &dd.avail, r: &dd.rounds.free},
		{k: "inodes_total", s: &dd.Inodes},
		{k: "inodes_used", s: &dd.Iused, r: &dd.roundInodes.used},
		{k: "inodes_free", s: &dd.Ifree, r: &dd.roundInodes.free},
	}) {
		return false // either fields lookup or casting failed
	}

	if push {
		up.diskData[diski] = dd
		return true
	}
	return false
}

func writeSystemDiskEnd(up *Update) {
	for _, dd := range up.diskData {
		if dd == nil {
			continue
		}
		dd.UsePct = format.Percent(dd.rounds.used, dd.rounds.used+dd.rounds.free)
		dd.IusePct = format.Percent(dd.roundInodes.used, dd.roundInodes.used+dd.roundInodes.free)
	}
}

func writeSystemMem(m telegraf.Metric, up *Update, memi int) {
	md := up.memData[memi]

	md.Kind = m.Name()
	isRAM := md.Kind == "mem"
	if isRAM {
		md.Kind = "RAM"
	}

	var values struct{ total, free int64 }
	var rounds struct{ total, used uint64 }

	fields := m.Fields()
	if !decode(fields, []convert{
		{k: "total", v: &values.total, s: &md.Total, r: &rounds.total},
		{k: "free", v: &values.free, s: &md.Free},
	}) {
		return // either fields lookup or casting failed
	}
	if isRAM {
		md.Used = humanB(values.total-values.free, &rounds.used)
	} else if !decode(fields, []convert{
		{k: "used", s: &md.Used, r: &rounds.used},
	}) {
		return // either fields lookup or casting failed
	}
	md.UsePct = format.Percent(rounds.used, rounds.total)
}

func writeSystemNet(m telegraf.Metric, up *Update, neti int) bool {
	tags := m.Tags()
	name := tags["interface"]
	if name == "all" { // uninterested NetProto stats
		return false
	}

	push := true
	var nd *netData
	/*
		for _, n := range up.netData {
			if n != nil && n.Name == name {
				nd, push = n, false
				break
			}
		}
	*/
	if nd == nil {
		nd = &netData{Name: name, IP: tags["ip"]}
		if _, ok := tags["nonemptyifLoopback"]; ok {
			nd.loopback = true
		}
	}

	if !decode(m.Fields(), []convert{
		{k: "bytes_sent", s: &nd.BytesOut},
		{k: "bytes_recv", s: &nd.BytesIn},
		{f: humanUnitless, k: "packets_sent", s: &nd.PacketsOut},
		{f: humanUnitless, k: "packets_recv", s: &nd.PacketsIn},
		{f: humanUnitless, k: "err_in", s: &nd.ErrorsIn},
		{f: humanUnitless, k: "err_out", s: &nd.ErrorsOut},
		{f: humanUnitless, k: "drop_in", s: &nd.DropsIn},
		{f: humanUnitless, k: "drop_out", s: &nd.DropsOut},

		{f: humanBits, k: "delta_bytes_sent", s: &nd.DeltaBitsOut, v: &nd.DeltaBytesOutNum},
		{f: humanBits, k: "delta_bytes_recv", s: &nd.DeltaBitsIn},

		{f: humanUnitless, k: "delta_packets_sent", s: &nd.DeltaPacketsOut},
		{f: humanUnitless, k: "delta_packets_recv", s: &nd.DeltaPacketsIn},
		{f: humanUnitless, k: "delta_err_in", s: &nd.DeltaErrorsIn},
		{f: humanUnitless, k: "delta_err_out", s: &nd.DeltaErrorsOut},
		{f: humanUnitless, k: "delta_drop_in", s: &nd.DeltaDropsIn},
		{f: humanUnitless, k: "delta_drop_out", s: &nd.DeltaDropsOut},
	}) {
		return false // either fields lookup or casting failed
	}

	if push {
		up.netData[neti] = nd
		return true
	}
	return false
}

func writeSystemOstent(m telegraf.Metric, up *Update) {
	for k, v := range m.Fields() {
		up.kv[k] = v
	}
}

func (o *ostent) Close() error         { return nil }
func (o *ostent) Connect() error       { return nil }
func (o *ostent) Description() string  { return "in-memory output" }
func (o *ostent) SampleConfig() string { return `` }

// func (o *ostent) SetSerializer(s serializers.Serializer) { o.serializer = s }

func (o *ostent) Write(ms []telegraf.Metric) error {
	if len(ms) == 0 {
		return nil
	}

	cpuno, diskno, netno := 0, 0, 0
	for _, m := range ms {
		switch m.Name() {
		case "cpu":
			cpuno++
		case "disk":
			diskno++
		case "net":
			netno++
		}
	}
	up := &Update{
		usernames: make(map[int64]string),

		kv:       make(map[string]interface{}),
		cpuData:  make([]*cpuData, cpuno),
		diskData: make([]*diskData, diskno),
		netData:  make([]*netData, netno),
	}
	for i := range up.memData {
		up.memData[i] = new(memData)
	}

	var procElapsed bool
	cpui, diski, neti := 0, 0, 0
	for _, m := range ms {
		switch m.Name() {
		case "system_ostent":
			writeSystemOstent(m, up)

		case "cpu":
			if writeSystemCPU(m, up, cpui) {
				cpui++
			}
		case "disk":
			if writeSystemDisk(m, up, diski) {
				diski++
			}
		case "mem":
			writeSystemMem(m, up, 0)
		case "swap":
			writeSystemMem(m, up, 1)
		case "net":
			if writeSystemNet(m, up, neti) {
				neti++
			}
		case "procstat_ostent":
			if writeProcstat(m, up) {
				procElapsed = true
			}
		}
	}
	up.cpuData = up.cpuData[:cpui]
	up.diskData = up.diskData[:diski]
	up.netData = up.netData[:neti]

	writeProcstatEnd(up)
	writeSystemDiskEnd(up)

	if false {
		fmt.Printf("Written  %#v pids; elapsed %t\n", len(up.procData), procElapsed)
	}

	Updates.set(up)
	return nil
}

type Update struct {
	usernames map[int64]string

	kv       map[string]interface{}
	cpuData  []*cpuData
	diskData []*diskData
	laData   []laData
	memData  [2]*memData
	netData  []*netData
	procData []*procData
}

func (u *updates) set(up *Update) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	close(u.ch)
	u.ch = make(chan *Update, 1)
	u.ch <- up
}

func (u *updates) Get() chan *Update {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	return u.ch
}

type updates struct {
	// mutex to protect ch (it's being closed and recreated in set, read in Get)
	mutex sync.Mutex
	ch    chan *Update
}

var Updates = &updates{ch: make(chan *Update, 1)}

var Output = &ostent{}

func init() { outputs.Add("ostent", func() telegraf.Output { return Output }) }
