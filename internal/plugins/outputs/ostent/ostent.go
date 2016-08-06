package ostent

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/shirou/gopsutil/disk"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/outputs"
	// "github.com/influxdata/telegraf/plugins/serializers"

	"github.com/ostrost/ostent/format"
)

type group struct {
	mutex sync.Mutex
	kv    map[string]interface{}
}

type diskData struct {
	DevName string
	DirName string

	// strings with units

	// bytes
	Total  string
	Used   string
	Avail  string
	UsePct uint

	// inodes
	Inodes  string
	Iused   string
	Ifree   string
	IusePct uint
}

type groupDisk struct {
	mutex sync.Mutex
	list  []diskData
}

type cpu struct {
	N       string
	UserPct int
	SysPct  int
	WaitPct int
	IdlePct int
}

type groupCPU struct {
	mutex sync.Mutex
	list  []cpu
}

type cpuList []cpu

// Len, Swap and Less satisfy sorting interface.
func (cl cpuList) Len() int           { return len(cl) }
func (cl cpuList) Swap(i, j int)      { cl[i], cl[j] = cl[j], cl[i] }
func (cl cpuList) Less(i, j int) bool { return cl[i].IdlePct < cl[j].IdlePct }

type ostent struct {
	// serializer serializers.Serializer
	// Metrics map[string]*Metric
	systemDisk   groupDisk
	systemCPU    groupCPU
	systemOstent group

	diskParts dparts
}

type dparts struct {
	mutex sync.Mutex
	parts map[string]string
}

func (o *ostent) SystemOstentCopy() (map[string]string, lalist) {
	o.systemOstent.mutex.Lock()
	defer o.systemOstent.mutex.Unlock()
	dup := make(map[string]string, len(o.systemOstent.kv))
	for k, v := range o.systemOstent.kv {
		if !strings.HasPrefix(k, "load") {
			if s, ok := v.(string); ok {
				dup[k] = s
			}
		}
	}
	periods := [3]string{"1", "5", "15"}
	lal := lalist{make([]la, len(periods))}
	for i, period := range periods {
		if v, ok := o.systemOstent.kv["load"+period]; ok {
			if f, ok := v.(float64); ok {
				lal.List[i] = la{period, fmt.Sprintf("%.2f", f)}
			}
		}
	}
	return dup, lal
}

type la struct {
	Period, Value string
}

type lalist struct{ List []la }
type list struct{ List interface{} }

func (o *ostent) SystemCPUCopyL() interface{}  { return list{o.systemCPUCopy()} }
func (o *ostent) SystemDiskCopyL() interface{} { return list{o.systemDiskCopy()} }

func (o *ostent) systemDiskCopy() []diskData {
	o.systemDisk.mutex.Lock()
	defer o.systemDisk.mutex.Unlock()
	llen := len(o.systemDisk.list)
	if llen == 0 {
		return []diskData{}
	}
	dup := make([]diskData, llen)
	copy(dup, o.systemDisk.list)
	// for i, c := range o.systemDisk.list { copy[i] = c }
	// sort.Sort(diskList(copy))
	return dup
}

func (o *ostent) systemCPUCopy() []cpu {
	o.systemCPU.mutex.Lock()
	defer o.systemCPU.mutex.Unlock()
	llen := len(o.systemCPU.list)
	if llen == 0 {
		return []cpu{}
	}
	tshift := 0 // "total shift"
	if o.systemCPU.list[llen-1].N == "cpu-total" {
		tshift = 1
	}

	dup := make([]cpu, llen)
	copy(dup[tshift:], o.systemCPU.list[:llen-tshift])
	// for i, c := range o.systemCPU.list[:llen-tshift] { dup[tshift+i] = c }
	sort.Sort(cpuList(dup[tshift:]))
	if tshift != 0 {
		dup[0] = o.systemCPU.list[llen-tshift] // last, "cpu-total", becomes first
	}
	return dup
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

func (o *ostent) writeSystemDisk(diskno int, m telegraf.Metric) {
	fields := m.Fields()
	dd := diskData{DirName: m.Tags()["path"]}
	dd.DevName, _ = o.diskParts.mpDevice(dd.DirName) // err is ignored

	var aused, afree, iaused, iafree uint64
	for _, pair := range []struct {
		name  string
		value *string
		back  *uint64
	}{
		{"total", &dd.Total, nil},
		{"used", &dd.Used, &aused},
		{"free", &dd.Avail, &afree},
		{"inodes_total", &dd.Inodes, nil},
		{"inodes_used", &dd.Iused, &iaused},
		{"inodes_free", &dd.Ifree, &iafree},
	} {
		if field, ok := fields[pair.name]; ok {
			if v, ok := field.(int64); ok {
				if pair.back != nil {
					*pair.value, *pair.back, _ = format.HumanBandback(uint64(v))
				} else {
					*pair.value = format.HumanB(uint64(v))
				}
			}
		}
	}
	dd.UsePct = format.Percent(aused, aused+afree)
	dd.IusePct = format.Percent(iaused, iaused+iafree)

	o.systemDisk.mutex.Lock()
	defer o.systemDisk.mutex.Unlock()
	if len(o.systemDisk.list) < diskno {
		list := make([]diskData, diskno)
		copy(list, o.systemDisk.list)
		o.systemDisk.list = list
	}
	o.systemDisk.list[diskno-1] = dd
}

func (o *ostent) writeSystemCPU(cpuno int, m telegraf.Metric) {
	fields := m.Fields()
	id := m.Tags()["cpu"]
	c := cpu{N: id}

	for _, pair := range []struct {
		name  string
		value *int
	}{
		{"usage_user", &c.UserPct},
		{"usage_system", &c.SysPct},
		{"usage_iowait", &c.WaitPct},
		{"usage_idle", &c.IdlePct},
	} {
		if field, ok := fields[pair.name]; ok {
			if v, ok := field.(float64); ok {
				*pair.value = int(v)
			}
		}
	}

	o.systemCPU.mutex.Lock()
	defer o.systemCPU.mutex.Unlock()
	if len(o.systemCPU.list) < cpuno {
		list := make([]cpu, cpuno)
		copy(list, o.systemCPU.list)
		o.systemCPU.list = list
	}
	o.systemCPU.list[cpuno-1] = c
}

func (o *ostent) setDiskno(diskno int) {
	o.systemDisk.mutex.Lock()
	defer o.systemDisk.mutex.Unlock()
	if len(o.systemDisk.list) > diskno {
		o.systemDisk.list = o.systemDisk.list[:diskno]
	}
}

func (o *ostent) setCPUno(cpuno int) {
	o.systemCPU.mutex.Lock()
	defer o.systemCPU.mutex.Unlock()
	if len(o.systemCPU.list) > cpuno {
		o.systemCPU.list = o.systemCPU.list[:cpuno]
	}
}

func (o *ostent) writeSystemOstent(m telegraf.Metric) {
	o.systemOstent.mutex.Lock()
	defer o.systemOstent.mutex.Unlock()
	for k, v := range m.Fields() {
		o.systemOstent.kv[k] = v
	}
}

/*
type Metric struct{ value string }
func (m Metric) String() string { return m.value }

func (o *ostent) writeMetric(m telegraf.Metric) error {
	name := m.Name()
	for k, field := range m.Fields() {
		o.Metrics[name+"."+k] = &Metric{
			value: fmt.Sprintf("%#v", field),
		}
	}
	return nil
} // */

func (o *ostent) Connect() error       { return nil }
func (o *ostent) Close() error         { return nil }
func (o *ostent) SampleConfig() string { return `` }
func (o *ostent) Description() string  { return "in-memory output" }

// func (o *ostent) SetSerializer(s serializers.Serializer) { o.serializer = s }

func (o *ostent) Write(ms []telegraf.Metric) error {
	if len(ms) == 0 {
		return nil
	}

	cpus, disks := 0, 0
	for _, m := range ms {
		switch m.Name() {
		case "system_ostent":
			o.writeSystemOstent(m)
		case "cpu":
			cpus++
			o.writeSystemCPU(cpus, m)
		case "disk":
			disks++
			o.writeSystemDisk(disks, m)
		}
		// default: if err := o.writeMetric(m); err != nil { return err }
	}
	o.setDiskno(disks)
	o.setCPUno(cpus)
	return nil
}

var Output = &ostent{
	// Metrics: make(map[string]*Metric),
	systemOstent: group{kv: make(map[string]interface{})},
	systemCPU:    groupCPU{},
	systemDisk:   groupDisk{},
	diskParts:    dparts{parts: make(map[string]string)},
}

func init() { outputs.Add("ostent", func() telegraf.Output { return Output }) }
